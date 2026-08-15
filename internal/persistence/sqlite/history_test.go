package sqlite

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yukihito-jokyu/knowledge/internal/domain"
)

func TestHistoryMutationsPreserveHistoryAndCurrentIndex(t *testing.T) {
	store, err := Open(context.Background(), filepath.Join(t.TempDir(), "knowledge.db"))
	if err != nil {
		t.Fatalf("Storeを開く: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	first := createHistoryAssertion(t, store, "first")
	second := createHistoryAssertion(t, store, "second")
	attach := domain.AttachEvidenceRequest{
		AssertionID: first,
		Evidence: domain.CreateEvidence{
			Kind:       "correction",
			RawText:    "evidence",
			ObservedAt: "2026-08-14T00:00:00.000000000Z",
		},
	}
	attached, err := store.AttachEvidence(context.Background(), attach)
	if err != nil || attached.EvidenceID == "" {
		t.Fatalf("AttachEvidence() = %#v, %v", attached, err)
	}
	if _, err := store.AttachEvidence(context.Background(), domain.AttachEvidenceRequest{AssertionID: "missing"}); !errors.Is(err, domain.ErrAssertionNotFound) {
		t.Fatalf("not found = %v", err)
	}
	if _, err := store.AttachEvidence(context.Background(), attach); !errors.Is(err, domain.ErrMutationConflict) {
		t.Fatalf("duplicate evidence = %v", err)
	}
	revised, err := store.ReviseAssertion(context.Background(), domain.ReviseRequest{
		AssertionID:    first,
		NormalizedText: "revised",
		Scope: []domain.Scope{{
			Key:   "language",
			Value: "Go",
		}},
	})
	if err != nil || revised.PreviousRevision != 1 || revised.Revision != 2 {
		t.Fatalf("ReviseAssertion() = %#v, %v", revised, err)
	}
	detail, err := store.GetAssertion(context.Background(), first)
	if err != nil || detail.CurrentRevision != 2 || len(detail.Revisions) != 2 {
		t.Fatalf("revision history = %#v, %v", detail, err)
	}
	results, err := store.SearchText(context.Background(), "revised")
	if err != nil || len(results) != 1 || results[0].ID != first {
		t.Fatalf("current index = %#v, %v", results, err)
	}
	replacement := domain.SupersedeRequest{
		SupersededAssertionID:  first,
		ReplacementAssertionID: second,
	}
	relation, err := store.Supersede(context.Background(), replacement)
	if err != nil || relation.RelationID == "" {
		t.Fatalf("Supersede() = %#v, %v", relation, err)
	}
	if _, err := store.Supersede(context.Background(), replacement); !errors.Is(err, domain.ErrMutationConflict) {
		t.Fatalf("duplicate relation = %v", err)
	}
	if _, err := store.Supersede(context.Background(), domain.SupersedeRequest{
		SupersededAssertionID:  second,
		ReplacementAssertionID: first,
	}); !errors.Is(err, domain.ErrMutationConflict) {
		t.Fatalf("cycle = %v", err)
	}
	if _, err := store.Supersede(context.Background(), domain.SupersedeRequest{
		SupersededAssertionID:  first,
		ReplacementAssertionID: first,
	}); !errors.Is(err, domain.ErrMutationConflict) {
		t.Fatalf("self = %v", err)
	}
}

func createHistoryAssertion(t *testing.T, store *Store, text string) string {
	t.Helper()
	result, err := store.CreateAssertion(context.Background(), domain.CreateRequest{
		NormalizedText: text,
		Evidence: []domain.CreateEvidence{{
			Kind:       "user_code",
			RawText:    text,
			ObservedAt: "2026-08-13T00:00:00.000000000Z",
		}},
	})
	if err != nil {
		t.Fatalf("Assertionを作成する: %v", err)
	}

	return result.AssertionID
}

func TestHistoryComparisonHelpers(t *testing.T) {
	value := "value"
	if nullString(sql.NullString{}) != nil || nullString(sql.NullString{
		String: value,
		Valid:  true,
	}) == nil {
		t.Fatal("nullString")
	}
	if !sameScopes(
		[]domain.Scope{
			{
				Key:   "b",
				Value: "2",
			},
			{
				Key:   "a",
				Value: "1",
			},
		},
		[]domain.Scope{
			{
				Key:   "a",
				Value: "1",
			},
			{
				Key:   "b",
				Value: "2",
			},
		},
	) {
		t.Fatal("Scope集合比較")
	}
	if sameScopes([]domain.Scope{{
		Key:   "a",
		Value: "1",
	}}, nil) {
		t.Fatal("異なるScopeを同一視しました")
	}
	if !sameTemporal(nil, nil) || sameTemporal(nil, &domain.Temporal{}) {
		t.Fatal("nil temporal")
	}
	if !sameTemporal(&domain.Temporal{ValidFrom: &value}, &domain.Temporal{ValidFrom: &value}) {
		t.Fatal("同じtemporal")
	}
	if equalString(&value, nil) || !equalString(&value, &value) {
		t.Fatal("string比較")
	}
}

func TestHistoryMutationWriteFailures(t *testing.T) {
	originalQuery, originalExec, originalCommit := coverageQuery, coverageExec, coverageCommit
	t.Cleanup(func() { coverageQuery, coverageExec, coverageCommit = originalQuery, originalExec, originalCommit })
	for _, testCase := range []struct {
		name string
		call func(*Store) error
	}{
		{"attach", func(store *Store) error {
			_, err := store.AttachEvidence(context.Background(), domain.AttachEvidenceRequest{
				AssertionID: "asrt",
				Evidence:    domain.CreateEvidence{},
			})

			return err
		}},
		{"revise", func(store *Store) error {
			_, err := store.ReviseAssertion(context.Background(), domain.ReviseRequest{
				AssertionID:    "asrt",
				NormalizedText: "new",
			})

			return err
		}},
		{"supersede", func(store *Store) error {
			_, err := store.Supersede(context.Background(), domain.SupersedeRequest{
				SupersededAssertionID:  "old",
				ReplacementAssertionID: "new",
			})

			return err
		}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			coverageQuery = historyCoverageQuery
			coverageExec = func(string) error { return errors.New("write") }
			database, err := sql.Open(coverageDriverName, "")
			if err != nil {
				t.Fatal(err)
			}
			defer database.Close()
			if err := testCase.call(&Store{db: database}); err == nil {
				t.Fatal("write failureを返しません")
			}
		})
	}
	coverageQuery = historyCoverageQuery
	coverageExec = nil
	coverageCommit = func() error { return errors.New("commit") }
	database, err := sql.Open(coverageDriverName, "")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if _, err := (&Store{db: database}).AttachEvidence(context.Background(), domain.AttachEvidenceRequest{AssertionID: "asrt"}); err == nil {
		t.Fatal("commit failureを返しません")
	}
}

func historyCoverageQuery(query string) (driver.Rows, error) {
	switch {
	case strings.Contains(query, "current_revision"):
		return &coverageRows{
			columns: []string{"current_revision"},
			values:  [][]driver.Value{{int64(1)}},
		}, nil
	case strings.Contains(query, "assertion_revisions AS r"):
		return &coverageRows{
			columns: []string{"normalized_text", "scope_key", "scope_value", "valid_from", "valid_until", "version_scope", "observed_at", "last_verified"},
			values:  [][]driver.Value{{"old", nil, nil, nil, nil, nil, nil, nil}},
		}, nil
	case strings.Contains(query, "assertions WHERE"):
		return &coverageRows{
			columns: []string{"found"},
			values:  [][]driver.Value{{int64(1)}},
		}, nil
	default:
		return nil, sql.ErrNoRows
	}
}

// TestHistoryMutationFailureBoundaries keeps each mutation's transaction error
// boundary explicit.  The coverage driver has no persistent state, so a returned
// error also proves the following write/commit boundary cannot be reached.
func TestHistoryMutationFailureBoundaries(t *testing.T) {
	originalQuery, originalExec, originalCommit, originalBegin, originalGate, originalRandom := coverageQuery, coverageExec, coverageCommit, coverageBegin, waitForIntegrationGate, readRandom
	t.Cleanup(func() {
		coverageQuery, coverageExec, coverageCommit, coverageBegin = originalQuery, originalExec, originalCommit, originalBegin
		waitForIntegrationGate, readRandom = originalGate, originalRandom
	})
	attach := func(ctx context.Context, store *Store) error {
		_, err := store.AttachEvidence(ctx, domain.AttachEvidenceRequest{
			AssertionID: "asrt",
			Evidence:    domain.CreateEvidence{},
		})

		return err
	}
	revise := func(ctx context.Context, store *Store) error {
		_, err := store.ReviseAssertion(ctx, historyRevisionRequest())

		return err
	}
	supersede := func(ctx context.Context, store *Store) error {
		_, err := store.Supersede(ctx, domain.SupersedeRequest{
			SupersededAssertionID:  "old",
			ReplacementAssertionID: "new",
		})

		return err
	}
	tests := []struct {
		name      string
		call      func(context.Context, *Store) error
		configure func(context.CancelFunc)
	}{
		{"attach begin", attach, func(context.CancelFunc) { coverageBegin = historyError }},
		{"revise begin", revise, func(context.CancelFunc) { coverageBegin = historyError }},
		{"supersede begin", supersede, func(context.CancelFunc) { coverageBegin = historyError }},
		{"attach gate", attach, func(context.CancelFunc) { waitForIntegrationGate = historyGateError }},
		{"revise gate", revise, func(context.CancelFunc) { waitForIntegrationGate = historyGateError }},
		{"supersede gate", supersede, func(context.CancelFunc) { waitForIntegrationGate = historyGateError }},
		{"attach assertion query", attach, func(context.CancelFunc) { coverageQuery = historyQueryError("assertions WHERE") }},
		{"attach duplicate query", attach, func(context.CancelFunc) { coverageQuery = historyQueryError("FROM evidence") }},
		{"revise current query", revise, func(context.CancelFunc) { coverageQuery = historyQueryError("current_revision") }},
		{"revise content query", revise, func(context.CancelFunc) { coverageQuery = historyQueryError("assertion_revisions AS r") }},
		{"supersede first assertion", supersede, func(context.CancelFunc) { coverageQuery = historyQueryError("assertions WHERE") }},
		{"supersede cycle query", supersede, func(context.CancelFunc) { coverageQuery = historyQueryError("WITH RECURSIVE") }},
		{"supersede duplicate query", supersede, func(context.CancelFunc) { coverageQuery = historyQueryError("FROM relations") }},
		{"attach random", attach, func(context.CancelFunc) { readRandom = func([]byte) (int, error) { return 0, errors.New("random") } }},
		{"supersede random", supersede, func(context.CancelFunc) { readRandom = func([]byte) (int, error) { return 0, errors.New("random") } }},
		{"attach commit", attach, func(context.CancelFunc) { coverageCommit = historyError }},
		{"revise commit", revise, func(context.CancelFunc) { coverageCommit = historyError }},
		{"supersede commit", supersede, func(context.CancelFunc) { coverageCommit = historyError }},
		{"attach cancellation", attach, func(cancel context.CancelFunc) {
			coverageExec = func(string) error {
				cancel()

				return nil
			}
		}},
		{"revise cancellation", revise, func(cancel context.CancelFunc) {
			coverageExec = func(string) error {
				cancel()

				return nil
			}
		}},
		{"supersede cancellation", supersede, func(cancel context.CancelFunc) {
			coverageExec = func(string) error {
				cancel()

				return nil
			}
		}},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			coverageQuery, coverageExec, coverageCommit, coverageBegin = historyCoverageQuery, nil, nil, nil
			waitForIntegrationGate, readRandom = originalGate, originalRandom
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			testCase.configure(cancel)
			store := historyCoverageStore(t)
			if err := testCase.call(ctx, store); err == nil {
				t.Fatal("storage errorを返しません")
			}
		})
	}
}

func TestReviseAssertionWriteFailureBoundaries(t *testing.T) {
	originalQuery, originalExec := coverageQuery, coverageExec
	t.Cleanup(func() { coverageQuery, coverageExec = originalQuery, originalExec })
	for _, marker := range []string{
		"INSERT INTO assertion_revisions",
		"INSERT INTO revision_scopes",
		"INSERT INTO temporal_metadata",
		"UPDATE assertions SET current_revision",
		"DELETE FROM assertion_lexical_index",
		"INSERT INTO assertion_lexical_index",
	} {
		t.Run(marker, func(t *testing.T) {
			coverageQuery = historyCoverageQuery
			coverageExec = func(query string) error {
				if strings.Contains(query, marker) {
					return errors.New("write")
				}

				return nil
			}
			_, err := historyCoverageStore(t).ReviseAssertion(context.Background(), historyRevisionRequest())
			if err == nil {
				t.Fatal("storage errorを返しません")
			}
		})
	}
}

func TestHistoryComparisonAndConflictErrors(t *testing.T) {
	originalQuery := coverageQuery
	t.Cleanup(func() { coverageQuery = originalQuery })
	store := historyCoverageStore(t)
	tests := []struct {
		name  string
		query func(string) (driver.Rows, error)
		call  func() error
	}{
		{"revision scan", historyMalformedRevisionQuery, func() error {
			_, err := store.ReviseAssertion(context.Background(), historyRevisionRequest())

			return err
		}},
		{"revision rows", historyRevisionRowsErrorQuery, func() error {
			_, err := store.ReviseAssertion(context.Background(), historyRevisionRequest())

			return err
		}},
		{"has conflict", historyQueryError("WITH RECURSIVE"), func() error {
			_, err := store.Supersede(context.Background(), domain.SupersedeRequest{
				SupersededAssertionID:  "old",
				ReplacementAssertionID: "new",
			})

			return err
		}},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			coverageQuery = testCase.query
			if err := testCase.call(); err == nil {
				t.Fatal("storage errorを返しません")
			}
		})
	}
	if sameScopes([]domain.Scope{{
		Key:   "a",
		Value: "1",
	}}, []domain.Scope{{
		Key:   "a",
		Value: "2",
	}}) {
		t.Fatal("異なる同数Scopeを同一視しました")
	}
}

func TestHistoryRemainingBranches(t *testing.T) {
	originalQuery, originalExec := coverageQuery, coverageExec
	t.Cleanup(func() { coverageQuery, coverageExec = originalQuery, originalExec })
	tests := []struct {
		name  string
		query func(string) (driver.Rows, error)
		call  func(*Store) error
	}{
		{"revise not found", historyCurrentNotFoundQuery, func(store *Store) error {
			_, err := store.ReviseAssertion(context.Background(), historyRevisionRequest())

			return err
		}},
		{"same revision conflict with scope temporal", historySameRevisionQuery, func(store *Store) error {
			_, err := store.ReviseAssertion(context.Background(), historyRevisionRequest())

			return err
		}},
		{"supersede replacement assertion", historySecondAssertionErrorQuery(), func(store *Store) error {
			_, err := store.Supersede(context.Background(), domain.SupersedeRequest{
				SupersededAssertionID:  "old",
				ReplacementAssertionID: "new",
			})

			return err
		}},
		{"supersede duplicate storage", historyQueryError("target_kind = 'assertion'"), func(store *Store) error {
			_, err := store.Supersede(context.Background(), domain.SupersedeRequest{
				SupersededAssertionID:  "old",
				ReplacementAssertionID: "new",
			})

			return err
		}},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			coverageQuery, coverageExec = testCase.query, nil
			err := testCase.call(historyCoverageStore(t))
			if err == nil {
				t.Fatal("期待したerrorを返しません")
			}
		})
	}
	t.Run("revise cancellation after final write", func(t *testing.T) {
		coverageQuery = historyCoverageQuery
		ctx := &historyCancellableContext{Context: context.Background()}
		coverageExec = func(query string) error {
			if strings.Contains(query, "INSERT INTO assertion_lexical_index") {
				ctx.cancelled = true
			}

			return nil
		}
		request := domain.ReviseRequest{
			AssertionID:    "asrt",
			NormalizedText: "new",
		}
		if _, err := historyCoverageStore(t).ReviseAssertion(ctx, request); !errors.Is(err, context.Canceled) {
			t.Fatalf("cancellation = %v", err)
		}
	})
}

func historyCoverageStore(t *testing.T) *Store {
	t.Helper()
	database, err := sql.Open(coverageDriverName, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = database.Close() })

	return &Store{db: database}
}

func historyRevisionRequest() domain.ReviseRequest {
	value := "2026-08-14T00:00:00.000000000Z"

	return domain.ReviseRequest{
		AssertionID:    "asrt",
		NormalizedText: "new",
		Scope: []domain.Scope{{
			Key:   "key",
			Value: "value",
		}},
		Temporal: &domain.Temporal{ValidFrom: &value},
	}
}

func historyError() error { return errors.New("storage") }

func historyGateError(context.Context, string) error { return errors.New("gate") }

func historyQueryError(marker string) func(string) (driver.Rows, error) {
	return func(query string) (driver.Rows, error) {
		if strings.Contains(query, marker) {
			return nil, errors.New("query")
		}

		return historyCoverageQuery(query)
	}
}

func historyMalformedRevisionQuery(query string) (driver.Rows, error) {
	if strings.Contains(query, "assertion_revisions AS r") {
		return &coverageRows{
			columns: []string{"normalized_text"},
			values:  [][]driver.Value{{"old"}},
		}, nil
	}

	return historyCoverageQuery(query)
}

func historyRevisionRowsErrorQuery(query string) (driver.Rows, error) {
	if strings.Contains(query, "assertion_revisions AS r") {
		return &coverageRows{
			columns: []string{"normalized_text", "scope_key", "scope_value", "valid_from", "valid_until", "version_scope", "observed_at", "last_verified"},
			err:     errors.New("rows"),
		}, nil
	}

	return historyCoverageQuery(query)
}

func historyCurrentNotFoundQuery(query string) (driver.Rows, error) {
	if strings.Contains(query, "current_revision") {
		return nil, sql.ErrNoRows
	}

	return historyCoverageQuery(query)
}

func historySameRevisionQuery(query string) (driver.Rows, error) {
	if strings.Contains(query, "assertion_revisions AS r") {
		return &coverageRows{
			columns: []string{"normalized_text", "scope_key", "scope_value", "valid_from", "valid_until", "version_scope", "observed_at", "last_verified"},
			values:  [][]driver.Value{{"new", "key", "value", "2026-08-14T00:00:00.000000000Z", nil, nil, nil, nil}},
		}, nil
	}

	return historyCoverageQuery(query)
}

func historySecondAssertionErrorQuery() func(string) (driver.Rows, error) {
	count := 0

	return func(query string) (driver.Rows, error) {
		if strings.Contains(query, "assertions WHERE") {
			count++
			if count == 2 {
				return nil, errors.New("second assertion")
			}
		}

		return historyCoverageQuery(query)
	}
}

type historyCancellableContext struct {
	context.Context
	cancelled bool
}

func (c *historyCancellableContext) Err() error {
	if c.cancelled {
		return context.Canceled
	}

	return nil
}
