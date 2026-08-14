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
