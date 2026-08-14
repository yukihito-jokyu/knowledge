// SQLiteパッケージは知識ストアの接続とスキーマ移行を管理する。
package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/yukihito-jokyu/knowledge/internal/domain"
	_ "modernc.org/sqlite"
)

const driverName = "sqlite"

var (
	// ErrInconsistentSchemaは、空でも完全でもないスキーマを示す。
	ErrInconsistentSchema = errors.New("inconsistent knowledge store schema")

	// 埋込みスキーマ移行SQL。
	//go:embed migrations/*.sql
	migrationFiles         embed.FS
	migrationFileSystem    fs.FS = migrationFiles
	openDatabase                 = sql.Open
	loadEmbeddedMigrations       = embeddedMigrations
	inspectSchemaState           = func(store *Store, ctx context.Context) (schemaStatus, []int, error) { return store.schemaState(ctx) }
	commitMigration              = func(tx *sql.Tx) error { return tx.Commit() }
	waitForIntegrationGate       = waitIntegrationGate
)

var schemaObjects = []string{
	"schema_migrations",
	"assertions",
	"assertion_revisions",
	"revision_scopes",
	"evidence",
	"concepts",
	"concept_terms",
	"concept_aliases",
	"assertion_concepts",
	"assertion_aliases",
	"temporal_metadata",
	"relations",
	"assertion_lexical_index",
}

var schemaIndexes = []string{
	"idx_scopes_key_value",
	"idx_evidence_assertion",
	"idx_aliases_concept",
	"idx_concept_terms_concept",
	"idx_assertion_concepts_concept",
	"idx_assertion_aliases_value",
	"idx_relations_source",
	"idx_relations_target",
	"idx_temporal_window",
}

var migrationName = regexp.MustCompile(`^([0-9]+)_.+\.sql$`)

// Storeは永続化アダプター用の内部SQLite接続を保持する。
type Store struct {
	db *sql.DB
}

// OpenはSQLiteデータベースを開き、外部キーを有効化して埋込み移行を適用する。
func Open(ctx context.Context, dsn string) (*Store, error) {
	db, err := openDatabase(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite store: %w", err)
	}
	// SQLiteの外部キー有効化は接続単位のため、初期版では接続を1本に制限する。
	db.SetMaxOpenConns(1)
	store := &Store{db: db}
	if err := store.initialize(ctx); err != nil {
		_ = db.Close()

		return nil, err
	}

	return store, nil
}

// Closeは内部DB接続を閉じる。
func (s *Store) Close() error {
	return s.db.Close()
}

const enableForeignKeysSQL = `
PRAGMA foreign_keys = ON
`

func (s *Store) initialize(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, enableForeignKeysSQL); err != nil {
		return fmt.Errorf("enable sqlite foreign keys: %w", err)
	}

	migrations, err := loadEmbeddedMigrations()
	if err != nil {
		return err
	}
	state, applied, err := inspectSchemaState(s, ctx)
	if err != nil {
		return err
	}
	pending, err := pendingMigrations(migrations, applied)
	if err != nil {
		return ErrInconsistentSchema
	}
	switch state {
	case schemaCurrent:
		return s.applyMigrationScripts(ctx, pending)
	case schemaInconsistent:
		return ErrInconsistentSchema
	case schemaEmpty:
		return s.applyMigrationScripts(ctx, pending)
	default:
		return fmt.Errorf("unknown schema state")
	}
}

type schemaStatus int

const (
	schemaEmpty schemaStatus = iota
	schemaCurrent
	schemaInconsistent
)

const schemaObjectsSQL = `
SELECT name
FROM sqlite_schema
WHERE type IN ('table', 'virtual table')
  AND name IN (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
  )
ORDER BY name
`

const schemaIndexSQL = `
SELECT name
FROM sqlite_schema
WHERE type = 'index'
  AND name = ?
`

const schemaMigrationsSQL = `
SELECT version
FROM schema_migrations
ORDER BY version
`

func (s *Store) schemaState(ctx context.Context) (schemaStatus, []int, error) {
	args := make([]any, len(schemaObjects))
	for i, name := range schemaObjects {
		args[i] = name
	}
	rows, err := s.db.QueryContext(ctx, schemaObjectsSQL, args...)
	if err != nil {
		return schemaInconsistent, nil, fmt.Errorf("inspect sqlite schema: %w", err)
	}
	defer rows.Close()

	var found []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return schemaInconsistent, nil, fmt.Errorf("read sqlite schema: %w", err)
		}
		found = append(found, name)
	}
	if err := rows.Err(); err != nil {
		return schemaInconsistent, nil, fmt.Errorf("iterate sqlite schema: %w", err)
	}
	if len(found) == 0 {
		return schemaEmpty, nil, nil
	}
	if len(found) != len(schemaObjects) {
		return schemaInconsistent, nil, nil
	}
	for _, index := range schemaIndexes {
		var name string
		err := s.db.QueryRowContext(ctx, schemaIndexSQL, index).Scan(&name)
		if errors.Is(err, sql.ErrNoRows) {
			return schemaInconsistent, nil, nil
		}
		if err != nil {
			return schemaInconsistent, nil, fmt.Errorf("inspect sqlite index: %w", err)
		}
	}

	rows, err = s.db.QueryContext(ctx, schemaMigrationsSQL)
	if err != nil {
		return schemaInconsistent, nil, fmt.Errorf("read schema migrations: %w", err)
	}
	defer rows.Close()
	var versions []int
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return schemaInconsistent, nil, fmt.Errorf("read schema migration version: %w", err)
		}
		versions = append(versions, version)
	}
	if err := rows.Err(); err != nil {
		return schemaInconsistent, nil, fmt.Errorf("iterate schema migrations: %w", err)
	}
	if len(versions) == 0 {
		return schemaInconsistent, nil, nil
	}

	return schemaCurrent, versions, nil
}

const insertMigrationRecordSQL = `
INSERT INTO schema_migrations (
  version,
  applied_at
) VALUES (?, ?)
`

func (s *Store) applyMigrationScripts(ctx context.Context, migrations []migration) error {
	if len(migrations) == 0 {
		return nil
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin schema migration: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if err := waitForIntegrationGate(ctx, "migration"); err != nil {
		return err
	}
	for _, migration := range migrations {
		if _, err := tx.ExecContext(ctx, migration.sql); err != nil {
			return fmt.Errorf("apply schema migration: %w", err)
		}
		if _, err := tx.ExecContext(ctx, insertMigrationRecordSQL, migration.version, time.Now().UTC().Format(time.RFC3339Nano)); err != nil {
			return fmt.Errorf("record schema migration: %w", err)
		}
	}
	if err := commitMigration(tx); err != nil {
		return fmt.Errorf("commit schema migration: %w", err)
	}

	return nil
}

type migration struct {
	version int
	sql     string
}

func embeddedMigrations() ([]migration, error) {
	entries, err := fs.ReadDir(migrationFileSystem, "migrations")
	if err != nil {
		return nil, fmt.Errorf("read embedded migrations: %w", err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	migrations := make([]migration, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		matches := migrationName.FindStringSubmatch(entry.Name())
		if matches == nil {
			return nil, fmt.Errorf("invalid migration name %q", entry.Name())
		}
		version, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("parse migration version %q: %w", entry.Name(), err)
		}
		contents, err := fs.ReadFile(migrationFileSystem, "migrations/"+entry.Name())
		if err != nil {
			return nil, fmt.Errorf("read embedded migration %q: %w", entry.Name(), err)
		}
		migrations = append(migrations, migration{
			version: version,
			sql:     string(contents),
		})
	}
	for i, migration := range migrations {
		if migration.version != i+1 {
			return nil, fmt.Errorf("migrations must be consecutive from version 1")
		}
	}
	if len(migrations) == 0 {
		return nil, fmt.Errorf("no embedded migrations")
	}

	return migrations, nil
}

func pendingMigrations(migrations []migration, applied []int) ([]migration, error) {
	if len(applied) > len(migrations) {
		return nil, fmt.Errorf("migration history exceeds embedded migrations")
	}
	for i, version := range applied {
		if version != i+1 || migrations[i].version != version {
			return nil, fmt.Errorf("inconsistent migration history")
		}
	}

	return migrations[len(applied):], nil
}

const searchTextSQL = `
WITH candidate AS (
  SELECT assertion_id FROM assertion_lexical_index
  WHERE assertion_lexical_index MATCH ?
), ranked AS (
  SELECT candidate.assertion_id,
         EXISTS (SELECT 1 FROM assertion_lexical_index WHERE assertion_id = candidate.assertion_id AND assertion_lexical_index MATCH ('normalized_text : ' || ?)) AS assertion_text_hit,
         (EXISTS (SELECT 1 FROM assertion_lexical_index WHERE assertion_id = candidate.assertion_id AND assertion_lexical_index MATCH ('concept_name : ' || ?)) OR EXISTS (SELECT 1 FROM assertion_lexical_index WHERE assertion_id = candidate.assertion_id AND assertion_lexical_index MATCH ('concept_alias : ' || ?))) AS concept_hit,
         EXISTS (SELECT 1 FROM assertion_lexical_index WHERE assertion_id = candidate.assertion_id AND assertion_lexical_index MATCH ('assertion_alias : ' || ?)) AS alias_hit,
         EXISTS (SELECT 1 FROM assertion_lexical_index WHERE assertion_id = candidate.assertion_id AND assertion_lexical_index MATCH ('scope_key : ' || ?)) AS scope_key_hit,
         EXISTS (SELECT 1 FROM assertion_lexical_index WHERE assertion_id = candidate.assertion_id AND assertion_lexical_index MATCH ('scope_value : ' || ?)) AS scope_value_hit
  FROM candidate
)
SELECT a.assertion_id, r.normalized_text, a.current_revision,
       assertion_text_hit, concept_hit, alias_hit, scope_key_hit, scope_value_hit
FROM ranked
JOIN assertions AS a ON a.assertion_id = ranked.assertion_id
JOIN assertion_revisions AS r ON r.assertion_id = a.assertion_id AND r.revision = a.current_revision
ORDER BY (assertion_text_hit + concept_hit + alias_hit + scope_key_hit + scope_value_hit) DESC, a.assertion_id ASC
`

const assertionConceptsSQL = `
SELECT c.concept_id, c.name
FROM assertion_concepts AS ac
JOIN concepts AS c ON c.concept_id = ac.concept_id
WHERE ac.assertion_id = ?
ORDER BY c.concept_id ASC
`

const currentScopeSQL = `
SELECT s.scope_key, s.scope_value
FROM revision_scopes AS s
JOIN assertions AS a ON a.assertion_id = s.assertion_id AND a.current_revision = s.revision
WHERE s.assertion_id = ?
ORDER BY s.scope_key ASC
`

// SearchText は現行Assertionを字句検索する。
func (s *Store) SearchText(ctx context.Context, query string) ([]domain.AssertionSummary, error) {
	phrase := ftsPhrase(query)
	rows, err := s.db.QueryContext(ctx, searchTextSQL, phrase, phrase, phrase, phrase, phrase, phrase, phrase)
	if err != nil {
		return nil, fmt.Errorf("search text: %w", err)
	}
	defer rows.Close()

	results := make([]domain.AssertionSummary, 0)
	for rows.Next() {
		var result domain.AssertionSummary
		var textHit, conceptHit, aliasHit, scopeKeyHit, scopeValueHit bool
		if err := rows.Scan(&result.ID, &result.NormalizedText, &result.Revision, &textHit, &conceptHit, &aliasHit, &scopeKeyHit, &scopeValueHit); err != nil {
			return nil, fmt.Errorf("read text search result: %w", err)
		}
		result.MatchedFields = matchedFields(textHit, conceptHit, aliasHit, scopeKeyHit, scopeValueHit)
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate text search result: %w", err)
	}
	for index := range results {
		if err := s.hydrateSummary(ctx, &results[index]); err != nil {
			return nil, err
		}
	}

	return results, nil
}

func ftsPhrase(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
}

func matchedFields(textHit, conceptHit, aliasHit, scopeKeyHit, scopeValueHit bool) []string {
	fields := make([]string, 0, 5)
	if textHit {
		fields = append(fields, "assertion_text")
	}
	if conceptHit {
		fields = append(fields, "concept")
	}
	if aliasHit {
		fields = append(fields, "alias")
	}
	if scopeKeyHit {
		fields = append(fields, "scope_key")
	}
	if scopeValueHit {
		fields = append(fields, "scope_value")
	}

	return fields
}

func (s *Store) hydrateSummary(ctx context.Context, result *domain.AssertionSummary) error {
	concepts, err := s.conceptsForAssertion(ctx, result.ID)
	if err != nil {
		return err
	}
	scope, err := s.currentScopeForAssertion(ctx, result.ID)
	if err != nil {
		return err
	}
	result.Concepts = concepts
	result.Scope = scope

	return nil
}

func (s *Store) conceptsForAssertion(ctx context.Context, assertionID string) ([]domain.Concept, error) {
	rows, err := s.db.QueryContext(ctx, assertionConceptsSQL, assertionID)
	if err != nil {
		return nil, fmt.Errorf("get assertion concepts: %w", err)
	}
	defer rows.Close()
	concepts := make([]domain.Concept, 0)
	for rows.Next() {
		var concept domain.Concept
		if err := rows.Scan(&concept.ID, &concept.Name); err != nil {
			return nil, fmt.Errorf("read assertion concept: %w", err)
		}
		concepts = append(concepts, concept)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate assertion concepts: %w", err)
	}

	return concepts, nil
}

func (s *Store) currentScopeForAssertion(ctx context.Context, assertionID string) ([]domain.Scope, error) {
	return s.scopeForAssertion(ctx, assertionID)
}

func (s *Store) scopeForAssertion(ctx context.Context, assertionID string) ([]domain.Scope, error) {
	rows, err := s.db.QueryContext(ctx, currentScopeSQL, assertionID)
	if err != nil {
		return nil, fmt.Errorf("get assertion scope: %w", err)
	}
	defer rows.Close()
	scope := make([]domain.Scope, 0)
	for rows.Next() {
		var entry domain.Scope
		if err := rows.Scan(&entry.Key, &entry.Value); err != nil {
			return nil, fmt.Errorf("read assertion scope: %w", err)
		}
		scope = append(scope, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate assertion scope: %w", err)
	}

	return scope, nil
}

const searchConceptSQL = `
SELECT c.concept_id, c.name, a.assertion_id, r.normalized_text, a.current_revision
FROM concept_terms AS ct
JOIN concepts AS c ON c.concept_id = ct.concept_id
LEFT JOIN assertion_concepts AS ac ON ac.concept_id = c.concept_id
LEFT JOIN assertions AS a ON a.assertion_id = ac.assertion_id
LEFT JOIN assertion_revisions AS r ON r.assertion_id = a.assertion_id AND r.revision = a.current_revision
WHERE ct.term = ?
ORDER BY a.assertion_id ASC
`

// SearchConcept はConcept名またはAliasから現行Assertionを取得する。
func (s *Store) SearchConcept(ctx context.Context, term string) (domain.ConceptSearchResult, error) {
	rows, err := s.db.QueryContext(ctx, searchConceptSQL, term)
	if err != nil {
		return domain.ConceptSearchResult{}, fmt.Errorf("search concept: %w", err)
	}
	defer rows.Close()
	result := domain.ConceptSearchResult{Results: make([]domain.AssertionSummary, 0)}
	for rows.Next() {
		var concept domain.Concept
		var summary domain.AssertionSummary
		var assertionID sql.NullString
		var normalizedText sql.NullString
		var revision sql.NullInt64
		if err := rows.Scan(&concept.ID, &concept.Name, &assertionID, &normalizedText, &revision); err != nil {
			return domain.ConceptSearchResult{}, fmt.Errorf("read concept search result: %w", err)
		}
		result.Concept = &concept
		if !assertionID.Valid {
			continue
		}
		summary.ID = assertionID.String
		summary.NormalizedText = normalizedText.String
		summary.Revision = int(revision.Int64)
		result.Results = append(result.Results, summary)
	}
	if err := rows.Err(); err != nil {
		return domain.ConceptSearchResult{}, fmt.Errorf("iterate concept search result: %w", err)
	}
	for index := range result.Results {
		scope, err := s.currentScopeForAssertion(ctx, result.Results[index].ID)
		if err != nil {
			return domain.ConceptSearchResult{}, err
		}
		result.Results[index].Scope = scope
	}

	return result, nil
}

const assertionRevisionsSQL = `
SELECT a.current_revision, r.revision, r.normalized_text, t.revision,
       t.valid_from, t.valid_until, t.version_scope, t.observed_at, t.last_verified,
       s.scope_key, s.scope_value
FROM assertions AS a
JOIN assertion_revisions AS r ON r.assertion_id = a.assertion_id
LEFT JOIN temporal_metadata AS t ON t.assertion_id = r.assertion_id AND t.revision = r.revision
LEFT JOIN revision_scopes AS s ON s.assertion_id = r.assertion_id AND s.revision = r.revision
WHERE a.assertion_id = ?
ORDER BY r.revision ASC, s.scope_key ASC
`

const assertionConceptDetailsSQL = `
SELECT c.concept_id, c.name, ca.alias
FROM assertion_concepts AS ac
JOIN concepts AS c ON c.concept_id = ac.concept_id
LEFT JOIN concept_aliases AS ca ON ca.concept_id = c.concept_id
WHERE ac.assertion_id = ?
ORDER BY c.concept_id ASC, ca.alias ASC
`

const assertionAliasesSQL = `
SELECT alias_kind, alias_value
FROM assertion_aliases
WHERE assertion_id = ?
ORDER BY alias_kind ASC, alias_value ASC
`

// GetAssertion はAssertionの全revisionと関連データを取得する。
func (s *Store) GetAssertion(ctx context.Context, assertionID string) (domain.AssertionDetail, error) {
	if err := waitForIntegrationGate(ctx, "read"); err != nil {
		return domain.AssertionDetail{}, err
	}
	rows, err := s.db.QueryContext(ctx, assertionRevisionsSQL, assertionID)
	if err != nil {
		return domain.AssertionDetail{}, fmt.Errorf("get assertion: %w", err)
	}
	defer rows.Close()
	detail := domain.AssertionDetail{
		ID:        assertionID,
		Revisions: make([]domain.Revision, 0),
	}
	revisionIndexes := make(map[int]int)
	for rows.Next() {
		var currentRevision int
		var revisionNumber int
		var normalizedText string
		var temporalRevision sql.NullInt64
		var validFrom, validUntil, versionScope, observedAt, lastVerified sql.NullString
		var scopeKey, scopeValue sql.NullString
		if err := rows.Scan(&currentRevision, &revisionNumber, &normalizedText, &temporalRevision, &validFrom, &validUntil, &versionScope, &observedAt, &lastVerified, &scopeKey, &scopeValue); err != nil {
			return domain.AssertionDetail{}, fmt.Errorf("read assertion revision: %w", err)
		}
		detail.CurrentRevision = currentRevision
		index, exists := revisionIndexes[revisionNumber]
		if !exists {
			revision := domain.Revision{
				Number:         revisionNumber,
				NormalizedText: normalizedText,
				Scope:          make([]domain.Scope, 0),
			}
			if temporalRevision.Valid {
				revision.Temporal = &domain.Temporal{
					ValidFrom:    nullableString(validFrom),
					ValidUntil:   nullableString(validUntil),
					VersionScope: nullableString(versionScope),
					ObservedAt:   nullableString(observedAt),
					LastVerified: nullableString(lastVerified),
				}
			}
			detail.Revisions = append(detail.Revisions, revision)
			index = len(detail.Revisions) - 1
			revisionIndexes[revisionNumber] = index
		}
		if scopeKey.Valid {
			detail.Revisions[index].Scope = append(detail.Revisions[index].Scope, domain.Scope{
				Key:   scopeKey.String,
				Value: scopeValue.String,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return domain.AssertionDetail{}, fmt.Errorf("iterate assertion revisions: %w", err)
	}
	if len(detail.Revisions) == 0 {
		return domain.AssertionDetail{}, domain.ErrAssertionNotFound
	}
	concepts, err := s.conceptDetailsForAssertion(ctx, assertionID)
	if err != nil {
		return domain.AssertionDetail{}, err
	}
	aliases, err := s.aliasesForAssertion(ctx, assertionID)
	if err != nil {
		return domain.AssertionDetail{}, err
	}
	detail.Concepts = concepts
	detail.Aliases = aliases

	return detail, nil
}

func nullableString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}

	return &value.String
}

func (s *Store) conceptDetailsForAssertion(ctx context.Context, assertionID string) ([]domain.ConceptDetail, error) {
	rows, err := s.db.QueryContext(ctx, assertionConceptDetailsSQL, assertionID)
	if err != nil {
		return nil, fmt.Errorf("get assertion concept details: %w", err)
	}
	defer rows.Close()
	concepts := make([]domain.ConceptDetail, 0)
	for rows.Next() {
		var id, name string
		var alias sql.NullString
		if err := rows.Scan(&id, &name, &alias); err != nil {
			return nil, fmt.Errorf("read assertion concept detail: %w", err)
		}
		if len(concepts) == 0 || concepts[len(concepts)-1].ID != id {
			concepts = append(concepts, domain.ConceptDetail{
				ID:      id,
				Name:    name,
				Aliases: make([]string, 0),
			})
		}
		if alias.Valid {
			index := len(concepts) - 1
			concepts[index].Aliases = append(concepts[index].Aliases, alias.String)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate assertion concept details: %w", err)
	}

	return concepts, nil
}

func (s *Store) aliasesForAssertion(ctx context.Context, assertionID string) ([]domain.AssertionAlias, error) {
	rows, err := s.db.QueryContext(ctx, assertionAliasesSQL, assertionID)
	if err != nil {
		return nil, fmt.Errorf("get assertion aliases: %w", err)
	}
	defer rows.Close()
	aliases := make([]domain.AssertionAlias, 0)
	for rows.Next() {
		var alias domain.AssertionAlias
		if err := rows.Scan(&alias.Kind, &alias.Value); err != nil {
			return nil, fmt.Errorf("read assertion alias: %w", err)
		}
		aliases = append(aliases, alias)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate assertion aliases: %w", err)
	}

	return aliases, nil
}

const assertionExistsSQL = `
SELECT 1
FROM assertions
WHERE assertion_id = ?
`

const evidenceSQL = `
SELECT evidence_id, kind, raw_text, observed_at
FROM evidence
WHERE assertion_id = ?
ORDER BY observed_at ASC, evidence_id ASC
`

// GetEvidence はAssertionに紐付くEvidence履歴を取得する。
func (s *Store) GetEvidence(ctx context.Context, assertionID string) (domain.EvidenceResult, error) {
	var exists int
	err := s.db.QueryRowContext(ctx, assertionExistsSQL, assertionID).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.EvidenceResult{}, domain.ErrAssertionNotFound
	}
	if err != nil {
		return domain.EvidenceResult{}, fmt.Errorf("check assertion existence: %w", err)
	}
	rows, err := s.db.QueryContext(ctx, evidenceSQL, assertionID)
	if err != nil {
		return domain.EvidenceResult{}, fmt.Errorf("get evidence: %w", err)
	}
	defer rows.Close()
	result := domain.EvidenceResult{
		AssertionID: assertionID,
		Evidence:    make([]domain.Evidence, 0),
	}
	for rows.Next() {
		var evidence domain.Evidence
		if err := rows.Scan(&evidence.ID, &evidence.Kind, &evidence.RawText, &evidence.ObservedAt); err != nil {
			return domain.EvidenceResult{}, fmt.Errorf("read evidence: %w", err)
		}
		result.Evidence = append(result.Evidence, evidence)
	}
	if err := rows.Err(); err != nil {
		return domain.EvidenceResult{}, fmt.Errorf("iterate evidence: %w", err)
	}

	return result, nil
}

const relationSeedExistsSQL = `
SELECT 1
FROM assertions
WHERE assertion_id = ?
UNION ALL
SELECT 1
FROM concepts
WHERE concept_id = ?
`

const searchRelatedSQL = `
SELECT r.relation_id,
       r.relation_type,
       CASE WHEN r.source_kind = ? AND r.source_id = ? THEN 'outgoing' ELSE 'incoming' END,
       CASE WHEN r.source_kind = ? AND r.source_id = ? THEN r.target_kind ELSE r.source_kind END,
       CASE WHEN r.source_kind = ? AND r.source_id = ? THEN r.target_id ELSE r.source_id END,
       target_revision.normalized_text
FROM relations AS r
LEFT JOIN assertions AS target_assertion
  ON (CASE WHEN r.source_kind = ? AND r.source_id = ? THEN r.target_kind ELSE r.source_kind END) = 'assertion'
 AND target_assertion.assertion_id = (CASE WHEN r.source_kind = ? AND r.source_id = ? THEN r.target_id ELSE r.source_id END)
LEFT JOIN assertion_revisions AS target_revision
  ON target_revision.assertion_id = target_assertion.assertion_id
 AND target_revision.revision = target_assertion.current_revision
WHERE ((r.source_kind = ? AND r.source_id = ?)
    OR (r.target_kind = ? AND r.target_id = ?))
`

// SearchRelated は検索起点から見たRelationの相手を取得する。
func (s *Store) SearchRelated(ctx context.Context, kind string, id string, relationTypes []string) ([]domain.RelatedResult, error) {
	var exists int
	assertionID, conceptID := "", ""
	if kind == "assertion" {
		assertionID = id
	} else {
		conceptID = id
	}
	err := s.db.QueryRowContext(ctx, relationSeedExistsSQL, assertionID, conceptID).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrRelationSeedNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("check relation seed: %w", err)
	}
	query, arguments := relatedQuery(kind, id, relationTypes)
	rows, err := s.db.QueryContext(ctx, query, arguments...)
	if err != nil {
		return nil, fmt.Errorf("search related: %w", err)
	}
	defer rows.Close()
	results := make([]domain.RelatedResult, 0)
	for rows.Next() {
		var result domain.RelatedResult
		var normalizedText sql.NullString
		if err := rows.Scan(&result.RelationID, &result.RelationType, &result.Direction, &result.Target.Kind, &result.Target.ID, &normalizedText); err != nil {
			return nil, fmt.Errorf("read related result: %w", err)
		}
		if normalizedText.Valid {
			result.Target.NormalizedText = &normalizedText.String
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate related result: %w", err)
	}

	return results, nil
}

func relatedQuery(kind string, id string, relationTypes []string) (string, []any) {
	arguments := []any{
		kind,
		id,
		kind,
		id,
		kind,
		id,
		kind,
		id,
		kind,
		id,
		kind,
		id,
		kind,
		id,
	}
	query := searchRelatedSQL
	if len(relationTypes) > 0 {
		placeholders := strings.TrimSuffix(strings.Repeat("?,", len(relationTypes)), ",")
		query += " AND r.relation_type IN (" + placeholders + ")"
		for _, relationType := range relationTypes {
			arguments = append(arguments, relationType)
		}
	}

	return query + " ORDER BY r.relation_id ASC", arguments
}

const searchContradictionsSQL = `
WITH selected_concept(concept_id) AS (
  SELECT c.concept_id
  FROM concepts AS c
  LEFT JOIN concept_aliases AS ca ON ca.concept_id = c.concept_id
  WHERE ? IS NOT NULL AND (c.name = ? OR ca.alias = ?)
), selected_seed(assertion_id) AS (
  SELECT ? WHERE ? IS NOT NULL
  UNION
  SELECT ac.assertion_id
  FROM assertion_concepts AS ac
  JOIN selected_concept AS sc ON sc.concept_id = ac.concept_id
), candidates AS (
  SELECT r.relation_id, 'outgoing' AS direction, ss.assertion_id AS seed_id, r.target_id
  FROM relations AS r
  JOIN selected_seed AS ss ON r.source_kind = 'assertion' AND r.source_id = ss.assertion_id
  WHERE r.relation_type = 'contradicts' AND r.target_kind = 'assertion'
  UNION ALL
  SELECT r.relation_id, 'incoming' AS direction, ss.assertion_id AS seed_id, r.source_id
  FROM relations AS r
  JOIN selected_seed AS ss ON r.target_kind = 'assertion' AND r.target_id = ss.assertion_id
  WHERE r.relation_type = 'contradicts' AND r.source_kind = 'assertion'
)
SELECT c.relation_id, c.direction, c.seed_id, c.target_id, r.normalized_text
FROM candidates AS c
JOIN assertions AS a ON a.assertion_id = c.target_id
JOIN assertion_revisions AS r ON r.assertion_id = a.assertion_id AND r.revision = a.current_revision
ORDER BY c.relation_id ASC, c.seed_id ASC, c.direction ASC
`

// SearchContradictions は指定selectorに一致するAssertionの矛盾候補を取得する。
func (s *Store) SearchContradictions(ctx context.Context, assertionID *string, concept *string) ([]domain.ContradictionResult, error) {
	var assertionValue, conceptValue any
	if assertionID != nil {
		assertionValue = *assertionID
	}
	if concept != nil {
		conceptValue = *concept
	}
	rows, err := s.db.QueryContext(ctx, searchContradictionsSQL, conceptValue, conceptValue, conceptValue, assertionValue, assertionValue)
	if err != nil {
		return nil, fmt.Errorf("search contradictions: %w", err)
	}
	defer rows.Close()
	results := make([]domain.ContradictionResult, 0)
	for rows.Next() {
		var result domain.ContradictionResult
		var normalizedText string
		if err := rows.Scan(&result.RelationID, &result.Direction, &result.SeedID, &result.Target.ID, &normalizedText); err != nil {
			return nil, fmt.Errorf("read contradiction result: %w", err)
		}
		result.Target.Kind = "assertion"
		result.Target.NormalizedText = &normalizedText
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate contradiction result: %w", err)
	}

	return results, nil
}
