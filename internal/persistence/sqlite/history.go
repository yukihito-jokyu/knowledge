package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/yukihito-jokyu/knowledge/internal/domain"
)

const historyAssertionExistsSQL = `SELECT 1 FROM assertions WHERE assertion_id = ?`

const evidenceExistsSQL = `
SELECT 1 FROM evidence
WHERE assertion_id = ? AND kind = ? AND raw_text = ? AND observed_at = ?
`

const insertHistoryEvidenceSQL = `
INSERT INTO evidence (evidence_id, assertion_id, kind, raw_text, observed_at, created_at)
VALUES (?, ?, ?, ?, ?, ?)
`

const currentRevisionSQL = `SELECT current_revision FROM assertions WHERE assertion_id = ?`

const revisionContentSQL = `
SELECT r.normalized_text, s.scope_key, s.scope_value,
  t.valid_from, t.valid_until, t.version_scope, t.observed_at, t.last_verified
FROM assertion_revisions AS r
LEFT JOIN revision_scopes AS s ON s.assertion_id = r.assertion_id AND s.revision = r.revision
LEFT JOIN temporal_metadata AS t ON t.assertion_id = r.assertion_id AND t.revision = r.revision
WHERE r.assertion_id = ? AND r.revision = ?
ORDER BY s.scope_key
`

const insertRevisionSQL = `
INSERT INTO assertion_revisions (assertion_id, revision, normalized_text, created_at)
VALUES (?, ?, ?, ?)
`

const insertRevisionScopeSQL = `
INSERT INTO revision_scopes (assertion_id, revision, scope_key, scope_value)
VALUES (?, ?, ?, ?)
`

const updateCurrentRevisionSQL = `UPDATE assertions SET current_revision = ? WHERE assertion_id = ?`

const deleteLexicalIndexSQL = `DELETE FROM assertion_lexical_index WHERE assertion_id = ?`

const supersedeExistsSQL = `
SELECT 1 FROM relations
WHERE source_kind = 'assertion' AND source_id = ? AND relation_type = 'supersedes'
  AND target_kind = 'assertion' AND target_id = ?
`

const supersedeCycleSQL = `
WITH RECURSIVE lineage(assertion_id) AS (
  SELECT target_id FROM relations
  WHERE source_kind = 'assertion' AND source_id = ? AND relation_type = 'supersedes'
  UNION
  SELECT r.target_id FROM relations AS r JOIN lineage AS l ON r.source_id = l.assertion_id
  WHERE r.source_kind = 'assertion' AND r.relation_type = 'supersedes'
)
SELECT 1 FROM lineage WHERE assertion_id = ?
`

const insertSupersedeSQL = `
INSERT INTO relations (relation_id, source_kind, source_id, relation_type, target_kind, target_id, created_at)
VALUES (?, 'assertion', ?, 'supersedes', 'assertion', ?, ?)
`

// AttachEvidence は既存Assertionへ不変のEvidenceを追加する。
func (s *Store) AttachEvidence(ctx context.Context, request domain.AttachEvidenceRequest) (domain.AttachEvidenceResult, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.AttachEvidenceResult{}, fmt.Errorf("begin attach evidence: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if err := waitForIntegrationGate(ctx, "mutation"); err != nil {
		return domain.AttachEvidenceResult{}, err
	}
	if err := requireAssertion(ctx, tx, request.AssertionID); err != nil {
		return domain.AttachEvidenceResult{}, err
	}
	var duplicate int
	err = tx.QueryRowContext(ctx, evidenceExistsSQL, request.AssertionID, request.Evidence.Kind, request.Evidence.RawText, request.Evidence.ObservedAt).Scan(&duplicate)
	if err == nil {
		return domain.AttachEvidenceResult{}, domain.ErrMutationConflict
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return domain.AttachEvidenceResult{}, fmt.Errorf("check evidence conflict: %w", err)
	}
	id, err := createID("evd_")
	if err != nil {
		return domain.AttachEvidenceResult{}, err
	}
	if _, err := tx.ExecContext(ctx, insertHistoryEvidenceSQL, id, request.AssertionID, request.Evidence.Kind, request.Evidence.RawText, request.Evidence.ObservedAt, time.Now().UTC().Format(time.RFC3339Nano)); err != nil {
		return domain.AttachEvidenceResult{}, fmt.Errorf("insert evidence: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return domain.AttachEvidenceResult{}, err
	}
	if err := tx.Commit(); err != nil {
		return domain.AttachEvidenceResult{}, fmt.Errorf("commit attach evidence: %w", err)
	}

	return domain.AttachEvidenceResult{
		AssertionID: request.AssertionID,
		EvidenceID:  id,
	}, nil
}

// ReviseAssertion は新revisionと現行参照、字句Indexを同時に更新する。
func (s *Store) ReviseAssertion(ctx context.Context, request domain.ReviseRequest) (domain.ReviseResult, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.ReviseResult{}, fmt.Errorf("begin revise assertion: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if err := waitForIntegrationGate(ctx, "mutation"); err != nil {
		return domain.ReviseResult{}, err
	}
	var current int
	if err := tx.QueryRowContext(ctx, currentRevisionSQL, request.AssertionID).Scan(&current); errors.Is(err, sql.ErrNoRows) {
		return domain.ReviseResult{}, domain.ErrAssertionNotFound
	} else if err != nil {
		return domain.ReviseResult{}, fmt.Errorf("read current revision: %w", err)
	}
	if same, err := sameRevision(ctx, tx, request, current); err != nil {
		return domain.ReviseResult{}, err
	} else if same {
		return domain.ReviseResult{}, domain.ErrMutationConflict
	}
	next := current + 1
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := tx.ExecContext(ctx, insertRevisionSQL, request.AssertionID, next, request.NormalizedText, now); err != nil {
		return domain.ReviseResult{}, fmt.Errorf("insert revision: %w", err)
	}
	for _, scope := range request.Scope {
		if _, err := tx.ExecContext(ctx, insertRevisionScopeSQL, request.AssertionID, next, scope.Key, scope.Value); err != nil {
			return domain.ReviseResult{}, fmt.Errorf("insert revision scope: %w", err)
		}
	}
	if request.Temporal != nil {
		if _, err := tx.ExecContext(ctx, insertCreateTemporalSQL, request.AssertionID, next, request.Temporal.ValidFrom, request.Temporal.ValidUntil, request.Temporal.VersionScope, request.Temporal.ObservedAt, request.Temporal.LastVerified); err != nil {
			return domain.ReviseResult{}, fmt.Errorf("insert revision temporal: %w", err)
		}
	}
	if _, err := tx.ExecContext(ctx, updateCurrentRevisionSQL, next, request.AssertionID); err != nil {
		return domain.ReviseResult{}, fmt.Errorf("update current revision: %w", err)
	}
	if _, err := tx.ExecContext(ctx, deleteLexicalIndexSQL, request.AssertionID); err != nil {
		return domain.ReviseResult{}, fmt.Errorf("delete lexical index: %w", err)
	}
	if _, err := tx.ExecContext(ctx, insertCreateLexicalIndexSQL, request.AssertionID); err != nil {
		return domain.ReviseResult{}, fmt.Errorf("insert lexical index: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return domain.ReviseResult{}, err
	}
	if err := tx.Commit(); err != nil {
		return domain.ReviseResult{}, fmt.Errorf("commit revise assertion: %w", err)
	}

	return domain.ReviseResult{
		AssertionID:      request.AssertionID,
		PreviousRevision: current,
		Revision:         next,
	}, nil
}

func requireAssertion(ctx context.Context, tx *sql.Tx, id string) error {
	var found int
	if err := tx.QueryRowContext(ctx, historyAssertionExistsSQL, id).Scan(&found); errors.Is(err, sql.ErrNoRows) {
		return domain.ErrAssertionNotFound
	} else if err != nil {
		return fmt.Errorf("check assertion: %w", err)
	}

	return nil
}

func sameRevision(ctx context.Context, tx *sql.Tx, request domain.ReviseRequest, current int) (bool, error) {
	rows, err := tx.QueryContext(ctx, revisionContentSQL, request.AssertionID, current)
	if err != nil {
		return false, fmt.Errorf("read current revision content: %w", err)
	}
	defer rows.Close()
	var text string
	var scopes []domain.Scope
	var temporal *domain.Temporal
	for rows.Next() {
		var key, value sql.NullString
		var validFrom, validUntil, versionScope, observedAt, lastVerified sql.NullString
		if err := rows.Scan(&text, &key, &value, &validFrom, &validUntil, &versionScope, &observedAt, &lastVerified); err != nil {
			return false, fmt.Errorf("read current revision row: %w", err)
		}
		if key.Valid {
			scopes = append(scopes, domain.Scope{
				Key:   key.String,
				Value: value.String,
			})
		}
		if validFrom.Valid || validUntil.Valid || versionScope.Valid || observedAt.Valid || lastVerified.Valid {
			temporal = &domain.Temporal{
				ValidFrom:    nullString(validFrom),
				ValidUntil:   nullString(validUntil),
				VersionScope: nullString(versionScope),
				ObservedAt:   nullString(observedAt),
				LastVerified: nullString(lastVerified),
			}
		}
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("iterate current revision: %w", err)
	}

	return text == request.NormalizedText && sameScopes(scopes, request.Scope) && sameTemporal(temporal, request.Temporal), nil
}

func nullString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	result := value.String

	return &result
}
func sameScopes(left, right []domain.Scope) bool {
	left = append([]domain.Scope(nil), left...)
	right = append([]domain.Scope(nil), right...)
	sort.Slice(left, func(i, j int) bool { return left[i].Key+"\x00"+left[i].Value < left[j].Key+"\x00"+left[j].Value })
	sort.Slice(right, func(i, j int) bool { return right[i].Key+"\x00"+right[i].Value < right[j].Key+"\x00"+right[j].Value })
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}

	return true
}
func sameTemporal(left, right *domain.Temporal) bool {
	if left == nil || right == nil {
		return left == right
	}

	return equalString(left.ValidFrom, right.ValidFrom) && equalString(left.ValidUntil, right.ValidUntil) && equalString(left.VersionScope, right.VersionScope) && equalString(left.ObservedAt, right.ObservedAt) && equalString(left.LastVerified, right.LastVerified)
}
func equalString(left, right *string) bool {
	if left == nil || right == nil {
		return left == right
	}

	return *left == *right
}

// Supersede は置換Relationを一度だけ追加する。
func (s *Store) Supersede(ctx context.Context, request domain.SupersedeRequest) (domain.SupersedeResult, error) {
	if request.SupersededAssertionID == request.ReplacementAssertionID {
		return domain.SupersedeResult{}, domain.ErrMutationConflict
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.SupersedeResult{}, fmt.Errorf("begin supersede: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if err := waitForIntegrationGate(ctx, "mutation"); err != nil {
		return domain.SupersedeResult{}, err
	}
	if err := requireAssertion(ctx, tx, request.SupersededAssertionID); err != nil {
		return domain.SupersedeResult{}, err
	}
	if err := requireAssertion(ctx, tx, request.ReplacementAssertionID); err != nil {
		return domain.SupersedeResult{}, err
	}
	cycle, err := hasConflict(ctx, tx, supersedeCycleSQL, request.SupersededAssertionID, request.ReplacementAssertionID)
	if err != nil {
		return domain.SupersedeResult{}, err
	}
	duplicate, err := hasConflict(ctx, tx, supersedeExistsSQL, request.ReplacementAssertionID, request.SupersededAssertionID)
	if err != nil {
		return domain.SupersedeResult{}, err
	}
	if cycle || duplicate {
		return domain.SupersedeResult{}, domain.ErrMutationConflict
	}
	id, err := createID("rel_")
	if err != nil {
		return domain.SupersedeResult{}, err
	}
	if _, err := tx.ExecContext(ctx, insertSupersedeSQL, id, request.ReplacementAssertionID, request.SupersededAssertionID, time.Now().UTC().Format(time.RFC3339Nano)); err != nil {
		return domain.SupersedeResult{}, fmt.Errorf("insert supersede: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return domain.SupersedeResult{}, err
	}
	if err := tx.Commit(); err != nil {
		return domain.SupersedeResult{}, fmt.Errorf("commit supersede: %w", err)
	}

	return domain.SupersedeResult{
		RelationID:             id,
		SupersededAssertionID:  request.SupersededAssertionID,
		ReplacementAssertionID: request.ReplacementAssertionID,
	}, nil
}

func hasConflict(ctx context.Context, tx *sql.Tx, query string, args ...any) (bool, error) {
	var found int
	err := tx.QueryRowContext(ctx, query, args...).Scan(&found)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check supersede conflict: %w", err)
	}

	return true, nil
}
