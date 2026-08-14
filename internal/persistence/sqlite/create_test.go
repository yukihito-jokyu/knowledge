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

func TestCreateAssertion(t *testing.T) {
	store, err := Open(context.Background(), filepath.Join(t.TempDir(), "knowledge.db"))
	if err != nil {
		t.Fatalf("Storeを開く: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	request := domain.CreateRequest{
		NormalizedText: "channel send",
		Scope: []domain.Scope{
			{
				Key:   "language",
				Value: "Go",
			},
		},
		Concepts: []domain.CreateConcept{
			{
				Name: "channel",
				Aliases: []string{
					"chan",
				},
			},
		},
		Aliases: []domain.AssertionAlias{
			{
				Kind:  "identifier",
				Value: "ch",
			},
		},
		Evidence: []domain.CreateEvidence{
			{
				Kind:       "user_code",
				RawText:    "send(ch)",
				ObservedAt: "2026-08-14T00:00:00.000000000Z",
			},
		},
	}
	result, err := store.CreateAssertion(context.Background(), request)
	if err != nil {
		t.Fatalf("CreateAssertion() error = %v", err)
	}
	if result.AssertionID == "" || result.Revision != 1 || len(result.EvidenceIDs) != 1 || len(result.Concepts) != 1 {
		t.Fatalf("CreateAssertion() = %#v", result)
	}
	results, err := store.SearchText(context.Background(), "channel")
	if err != nil || len(results) != 1 || results[0].ID != result.AssertionID {
		t.Fatalf("SearchText() = %#v, %v", results, err)
	}
	if _, err := store.CreateAssertion(context.Background(), request); !errors.Is(err, domain.ErrCreateConflict) {
		t.Fatalf("重複CreateAssertion() error = %v", err)
	}
	if _, err := store.CreateAssertion(context.Background(), domain.CreateRequest{
		NormalizedText: "missing target",
		Relations: []domain.CreateRelation{
			{
				Type:       "causes",
				TargetKind: "assertion",
				TargetID:   "missing",
			},
		},
		Evidence: request.Evidence,
	}); !errors.Is(err, domain.ErrCreateRelationTargetNotFound) {
		t.Fatalf("不存在Relation targetのerror = %v", err)
	}
}

func TestCreateAssertionFailureBoundaries(t *testing.T) {
	originalGate := waitForIntegrationGate
	originalMarshalScope := marshalScope
	t.Cleanup(func() {
		waitForIntegrationGate = originalGate
		marshalScope = originalMarshalScope
	})
	request := domain.CreateRequest{NormalizedText: "text"}
	coverageQuery = func(string) (driver.Rows, error) { return nil, sql.ErrNoRows }
	database, err := sql.Open(coverageDriverName, "")
	if err != nil {
		t.Fatalf("coverage DB: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	for _, testCase := range []struct {
		name  string
		setup func()
	}{
		{
			name: "gate",
			setup: func() {
				waitForIntegrationGate = func(context.Context, string) error { return errors.New("gate") }
			},
		},
		{
			name: "marshal",
			setup: func() {
				marshalScope = func(any) ([]byte, error) { return nil, errors.New("marshal") }
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			waitForIntegrationGate = originalGate
			marshalScope = originalMarshalScope
			testCase.setup()
			if _, err := (&Store{db: database}).CreateAssertion(context.Background(), request); err == nil {
				t.Fatal("failureを返しません")
			}
		})
	}
	if _, err := (&Store{db: database}).CreateAssertion(context.Background(), request); err == nil {
		t.Fatal("Exec failureを返しません")
	}
}

func TestCreateAssertionWriteFailures(t *testing.T) {
	request := domain.CreateRequest{
		NormalizedText: "text",
		Scope: []domain.Scope{{
			Key:   "k",
			Value: "v",
		}},
		Concepts: []domain.CreateConcept{{
			Name:    "concept",
			Aliases: []string{"alias"},
		}},
		Aliases: []domain.AssertionAlias{{
			Kind:  "identifier",
			Value: "id",
		}},
		Temporal: &domain.Temporal{},
		Evidence: []domain.CreateEvidence{{
			Kind:       "user_code",
			RawText:    "e",
			ObservedAt: "2026-08-14T00:00:00Z",
		}},
		Relations: []domain.CreateRelation{{
			Type:       "causes",
			TargetKind: "assertion",
			TargetID:   "asrt_target",
		}},
	}
	for _, marker := range []string{"INSERT INTO assertions", "INSERT INTO assertion_revisions", "INSERT INTO revision_scopes", "INSERT INTO concepts", "INSERT INTO concept_terms", "INSERT INTO concept_aliases", "INSERT INTO assertion_concepts", "INSERT INTO assertion_aliases", "INSERT INTO temporal_metadata", "INSERT INTO evidence", "INSERT INTO relations", "INSERT INTO assertion_lexical_index"} {
		t.Run(marker, func(t *testing.T) {
			coverageQuery = func(query string) (driver.Rows, error) {
				if strings.Contains(query, "SELECT 1 FROM assertions") {
					return &coverageRows{
						columns: []string{"found"},
						values:  [][]driver.Value{{int64(1)}},
					}, nil
				}

				return nil, sql.ErrNoRows
			}
			coverageExec = func(query string) error {
				if strings.Contains(query, marker) {
					return errors.New("write failure")
				}

				return nil
			}
			database, err := sql.Open(coverageDriverName, "")
			if err != nil {
				t.Fatalf("coverage DB: %v", err)
			}
			t.Cleanup(func() { _ = database.Close() })
			if _, err := (&Store{db: database}).CreateAssertion(context.Background(), request); err == nil {
				t.Fatal("write failureを返しません")
			}
		})
	}
}

func TestCreateAssertionIdentifierAndCommitFailures(t *testing.T) {
	originalReadRandom := readRandom
	originalCommit := coverageCommit
	originalExec := coverageExec
	t.Cleanup(func() {
		readRandom = originalReadRandom
		coverageCommit = originalCommit
		coverageExec = originalExec
	})
	coverageQuery = func(string) (driver.Rows, error) { return nil, sql.ErrNoRows }
	database, err := sql.Open(coverageDriverName, "")
	if err != nil {
		t.Fatalf("coverage DB: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	readRandom = func([]byte) (int, error) { return 0, errors.New("random") }
	if _, err := (&Store{db: database}).CreateAssertion(context.Background(), domain.CreateRequest{NormalizedText: "text"}); err == nil {
		t.Fatal("ID failureを返しません")
	}
	readRandom = originalReadRandom
	coverageExec = nil
	coverageCommit = func() error { return errors.New("commit") }
	if _, err := (&Store{db: database}).CreateAssertion(context.Background(), domain.CreateRequest{NormalizedText: "text"}); err == nil {
		t.Fatal("commit failureを返しません")
	}
}

func TestCreateHelperIdentifierFailures(t *testing.T) {
	originalReadRandom := readRandom
	t.Cleanup(func() { readRandom = originalReadRandom })
	readRandom = func([]byte) (int, error) { return 0, errors.New("random") }
	store, err := Open(context.Background(), filepath.Join(t.TempDir(), "knowledge.db"))
	if err != nil {
		t.Fatalf("Storeを開く: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	tx, err := store.db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback() })
	if _, err := createEvidence(context.Background(), tx, "asrt", []domain.CreateEvidence{{}}, ""); err == nil {
		t.Fatal("evidence ID failureを返しません")
	}
	if _, err := createRelations(context.Background(), tx, "asrt", []domain.CreateRelation{{}}, ""); err == nil {
		t.Fatal("relation ID failureを返しません")
	}
}

func TestCreateRelationAndConceptLookupFailures(t *testing.T) {
	coverageQuery = func(string) (driver.Rows, error) { return nil, errors.New("query") }
	database, err := sql.Open(coverageDriverName, "")
	if err != nil {
		t.Fatalf("coverage DB: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	tx, err := database.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback() })
	if err := verifyCreateRelationTargets(context.Background(), tx, []domain.CreateRelation{{
		TargetKind: "assertion",
		TargetID:   "id",
	}}); err == nil {
		t.Fatal("relation lookup failureを返しません")
	}
	if _, err := resolveCreateConcept(context.Background(), tx, domain.CreateConcept{
		Name: "name",
	}, ""); err == nil {
		t.Fatal("concept lookup failureを返しません")
	}
}

func TestCreateRemainingFailureBranches(t *testing.T) {
	originalReadRandom := readRandom
	t.Cleanup(func() { readRandom = originalReadRandom })
	coverageQuery = func(string) (driver.Rows, error) { return nil, errors.New("query") }
	database, err := sql.Open(coverageDriverName, "")
	if err != nil {
		t.Fatalf("coverage DB: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	if _, err := (&Store{db: database}).CreateAssertion(context.Background(), domain.CreateRequest{
		NormalizedText: "text",
	}); err == nil {
		t.Fatal("conflict query failureを返しません")
	}
	lookupIDs := []string{"one", "two"}
	lookupCount := 0
	coverageQuery = func(string) (driver.Rows, error) {
		value := lookupIDs[lookupCount]
		lookupCount++

		return &coverageRows{
			columns: []string{"concept_id"},
			values:  [][]driver.Value{{value}},
		}, nil
	}
	tx, err := database.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback() })
	if _, err := resolveCreateConcept(context.Background(), tx, domain.CreateConcept{
		Name:    "name",
		Aliases: []string{"alias"},
	}, ""); !errors.Is(err, domain.ErrCreateConflict) {
		t.Fatalf("concept conflict = %v", err)
	}
	coverageQuery = func(string) (driver.Rows, error) { return nil, sql.ErrNoRows }
	readRandom = func([]byte) (int, error) { return 0, errors.New("random") }
	if _, err := resolveCreateConcept(context.Background(), tx, domain.CreateConcept{
		Name: "new",
	}, ""); err == nil {
		t.Fatal("concept ID failureを返しません")
	}
}

func TestCreateAssertionRollsBackWhenCanceledBeforeCommit(t *testing.T) {
	originalExec := coverageExec
	t.Cleanup(func() { coverageExec = originalExec })
	coverageQuery = func(string) (driver.Rows, error) { return nil, sql.ErrNoRows }
	database, err := sql.Open(coverageDriverName, "")
	if err != nil {
		t.Fatalf("coverage DB: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	coverageExec = func(query string) error {
		if strings.Contains(query, "assertion_lexical_index") {
			cancel()
		}

		return nil
	}
	if _, err := (&Store{db: database}).CreateAssertion(ctx, domain.CreateRequest{NormalizedText: "text"}); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancel error = %v", err)
	}
}

func TestCreateAssertionConcurrentConflict(t *testing.T) {
	path := filepath.Join(t.TempDir(), "knowledge.db")
	first, err := Open(context.Background(), path)
	if err != nil {
		t.Fatalf("first store: %v", err)
	}
	t.Cleanup(func() { _ = first.Close() })
	second, err := Open(context.Background(), path)
	if err != nil {
		t.Fatalf("second store: %v", err)
	}
	t.Cleanup(func() { _ = second.Close() })
	request := domain.CreateRequest{NormalizedText: "same"}
	results := make(chan error, 2)
	go func() { _, err := first.CreateAssertion(context.Background(), request); results <- err }()
	go func() { _, err := second.CreateAssertion(context.Background(), request); results <- err }()
	createErrors := []error{<-results, <-results}
	if (createErrors[0] == nil) == (createErrors[1] == nil) {
		t.Fatalf("concurrent errors = %#v", createErrors)
	}
	if errors.Is(createErrors[0], domain.ErrCreateConflict) || errors.Is(createErrors[1], domain.ErrCreateConflict) {
		t.Fatalf("external lockをconflictに分類しました: %#v", createErrors)
	}
}

func TestCreateAssertionBeginFailures(t *testing.T) {
	originalBegin := coverageBegin
	t.Cleanup(func() { coverageBegin = originalBegin })
	database, err := sql.Open(coverageDriverName, "")
	if err != nil {
		t.Fatalf("coverage DB: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	for _, beginErr := range []error{errors.New("database is locked"), errors.New("begin failure")} {
		coverageBegin = func() error { return beginErr }
		_, err := (&Store{db: database}).CreateAssertion(context.Background(), domain.CreateRequest{NormalizedText: "text"})
		if beginErr.Error() == "database is locked" && (err == nil || errors.Is(err, domain.ErrCreateConflict)) {
			t.Fatalf("locked = %v", err)
		}
		if beginErr.Error() == "begin failure" && err == nil {
			t.Fatal("begin failureを返しません")
		}
	}
}

func TestCreateConceptTermConstraintAndQueryDSN(t *testing.T) {
	originalExec := coverageExec
	t.Cleanup(func() { coverageExec = originalExec })
	coverageQuery = func(string) (driver.Rows, error) { return nil, sql.ErrNoRows }
	coverageExec = func(query string) error {
		if strings.Contains(query, "INSERT INTO concept_terms") {
			return errors.New("constraint failed")
		}

		return nil
	}
	database, err := sql.Open(coverageDriverName, "")
	if err != nil {
		t.Fatalf("coverage DB: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	tx, err := database.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback() })
	if _, err := resolveCreateConcept(context.Background(), tx, domain.CreateConcept{Name: "name"}, ""); !errors.Is(err, domain.ErrCreateConflict) {
		t.Fatalf("constraint = %v", err)
	}
	store, err := Open(context.Background(), filepath.Join(t.TempDir(), "query.db")+"?cache=shared")
	if err != nil {
		t.Fatalf("query dsn store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
}

func TestCreateAssertionRollsBackCanceledContext(t *testing.T) {
	store, err := Open(context.Background(), filepath.Join(t.TempDir(), "knowledge.db"))
	if err != nil {
		t.Fatalf("Storeを開く: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = store.CreateAssertion(ctx, domain.CreateRequest{NormalizedText: "canceled"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("CreateAssertion() error = %v", err)
	}
	var count int
	if err := store.db.QueryRowContext(context.Background(), "SELECT count(*) FROM assertions").Scan(&count); err != nil || count != 0 {
		t.Fatalf("rollback count/error = %d, %v", count, err)
	}
}

func TestCreateHelpers(t *testing.T) {
	store, err := Open(context.Background(), filepath.Join(t.TempDir(), "knowledge.db"))
	if err != nil {
		t.Fatalf("Storeを開く: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	seed, err := store.CreateAssertion(context.Background(), domain.CreateRequest{
		NormalizedText: "target",
		Concepts: []domain.CreateConcept{
			{
				Name: "existing",
			},
		},
		Evidence: []domain.CreateEvidence{
			{
				Kind:       "user_code",
				RawText:    "seed",
				ObservedAt: "2026-08-14T00:00:00.000000000Z",
			},
		},
	})
	if err != nil {
		t.Fatalf("seed作成: %v", err)
	}
	tx, err := store.db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("transaction開始: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback() })
	if err := verifyCreateRelationTargets(context.Background(), tx, []domain.CreateRelation{
		{
			TargetKind: "concept",
			TargetID:   seed.Concepts[0].ID,
		},
	}); err != nil {
		t.Fatalf("concept relation target: %v", err)
	}
	if _, err := tx.ExecContext(context.Background(), insertCreateAssertionSQL, "asrt_helper", "2026-08-14T00:00:00Z"); err != nil {
		t.Fatalf("helper assertion挿入: %v", err)
	}
	concepts, err := createConcepts(context.Background(), tx, "asrt_helper", []domain.CreateConcept{
		{
			Name: "existing",
			Aliases: []string{
				"existing-alias",
			},
		},
	}, "2026-08-14T00:00:00Z")
	if err != nil || len(concepts) != 1 || concepts[0].ID != seed.Concepts[0].ID {
		t.Fatalf("createConcepts() = %#v, %v", concepts, err)
	}
	if _, err := createEvidence(context.Background(), tx, "asrt_helper", []domain.CreateEvidence{
		{
			Kind:       "user_code",
			RawText:    "helper",
			ObservedAt: "2026-08-14T00:00:00Z",
		},
	}, "2026-08-14T00:00:00Z"); err != nil {
		t.Fatalf("createEvidence(): %v", err)
	}
	if _, err := createRelations(context.Background(), tx, "asrt_helper", []domain.CreateRelation{
		{
			Type:       "causes",
			TargetKind: "assertion",
			TargetID:   seed.AssertionID,
		},
	}, "2026-08-14T00:00:00Z"); err != nil {
		t.Fatalf("createRelations(): %v", err)
	}
}
