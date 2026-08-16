package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/yukihito-jokyu/knowledge/internal/domain"
	"github.com/yukihito-jokyu/knowledge/internal/persistence/sqlite"
)

func TestExecuteRetrievalWithStore(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	tests := []struct {
		name        string
		arguments   []string
		store       retrievalStoreStub
		wantHandled bool
		wantCode    errorCode
		wantData    bool
		wantContext bool
	}{
		{
			name: "search-text",
			arguments: []string{
				"search-text",
				"--query",
				"channel",
			},
			wantHandled: true,
			wantData:    true,
			wantContext: true,
		},
		{
			name: "search-concept",
			arguments: []string{
				"search-concept",
				"--concept",
				"channel",
			},
			wantHandled: true,
			wantData:    true,
			wantContext: true,
		},
		{
			name: "get",
			arguments: []string{
				"get",
				"--assertion-id",
				"asrt_01",
			},
			wantHandled: true,
			wantData:    true,
			wantContext: true,
		},
		{
			name: "get-evidence",
			arguments: []string{
				"get-evidence",
				"--assertion-id",
				"asrt_01",
			},
			wantHandled: true,
			wantData:    true,
			wantContext: true,
		},
		{
			name: "search-related",
			arguments: []string{
				"search-related",
				"--seed-kind",
				"assertion",
				"--seed-id",
				"asrt_01",
			},
			wantHandled: true,
			wantData:    true,
			wantContext: true,
		},
		{
			name: "search-contradictions",
			arguments: []string{
				"search-contradictions",
				"--concept",
				"channel",
			},
			wantHandled: true,
			wantData:    true,
			wantContext: true,
		},
		{
			name: "search-temporal",
			arguments: []string{
				"search-temporal",
				"--concept",
				"channel",
				"--scope-key",
				"language",
				"--scope-value",
				"Go",
			},
			wantHandled: true,
			wantData:    true,
			wantContext: true,
		},
		{
			name: "search-relatedのnot_found",
			arguments: []string{
				"search-related",
				"--seed-kind",
				"assertion",
				"--seed-id",
				"asrt_01",
			},
			store:       retrievalStoreStub{relatedError: domain.ErrRelationSeedNotFound},
			wantHandled: true,
			wantCode:    notFoundError,
		},
		{
			name: "search-contradictionsのstorage_error",
			arguments: []string{
				"search-contradictions",
				"--assertion-id",
				"asrt_01",
			},
			store:       retrievalStoreStub{contradictionError: errors.New("read failure")},
			wantHandled: true,
			wantCode:    storageError,
		},
		{
			name: "search-temporalのstorage_error",
			arguments: []string{
				"search-temporal",
				"--concept",
				"channel",
			},
			store:       retrievalStoreStub{temporalError: errors.New("read failure")},
			wantHandled: true,
			wantCode:    storageError,
		},
		{
			name: "getのnot_found",
			arguments: []string{
				"get",
				"--assertion-id",
				"asrt_01",
			},
			store:       retrievalStoreStub{getError: domain.ErrAssertionNotFound},
			wantHandled: true,
			wantCode:    notFoundError,
		},
		{
			name: "getのstorage_error",
			arguments: []string{
				"get",
				"--assertion-id",
				"asrt_01",
			},
			store:       retrievalStoreStub{getError: errors.New("read failure")},
			wantHandled: true,
			wantCode:    storageError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, parseError := parseCommand(tt.arguments)
			if parseError.code != "" {
				t.Fatalf("%v をparseする: %v", tt.arguments, parseError)
			}
			var received context.Context
			store := tt.store
			store.receivedContext = &received
			data, executionError, handled := executeRetrievalWithStore(ctx, parsed, store)
			if handled != tt.wantHandled || executionError.code != tt.wantCode || (data != nil) != tt.wantData {
				t.Fatalf("executeRetrievalWithStore(%v) = %#v, %#v, %t", tt.arguments, data, executionError, handled)
			}
			if tt.wantContext && received != ctx {
				t.Fatal("ContextがStoreへ伝播しません")
			}
		})
	}
}

func TestExecuteRetrievalWithStoreDoesNotHandleUnknownOperation(t *testing.T) {
	_, _, handled := executeRetrievalWithStore(context.Background(), command{operation: "unknown"}, retrievalStoreStub{})
	if handled {
		t.Fatal("未知の操作を処理済みにしました")
	}
}

func TestExecuteRetrievalWithStoreReturnsTemporalValidationError(t *testing.T) {
	_, executionError, handled := executeRetrievalWithStore(
		context.Background(),
		command{
			operation: "search-temporal",
			options: []option{
				{
					name:  "at",
					value: "invalid",
				},
			},
		},
		retrievalStoreStub{},
	)
	if !handled || executionError.code != validationError || executionError.field != "at" {
		t.Fatalf("無効な時刻の結果 = %#v, handled=%t", executionError, handled)
	}
}

func TestExecuteRetrievalUsesDefaultStore(t *testing.T) {
	originalUserConfigDir := userConfigDir
	originalOpenSQLiteStore := openSQLiteStore
	t.Cleanup(func() {
		userConfigDir = originalUserConfigDir
		openSQLiteStore = originalOpenSQLiteStore
	})
	tests := []struct {
		name        string
		open        func(*context.Context) func(context.Context, string) (*sqlite.Store, error)
		wantCode    errorCode
		wantData    bool
		wantContext bool
	}{
		{
			name: "ContextをSQLiteへ伝播する",
			open: func(received *context.Context) func(context.Context, string) (*sqlite.Store, error) {
				return func(got context.Context, path string) (*sqlite.Store, error) {
					*received = got

					return sqlite.Open(got, path)
				}
			},
			wantData:    true,
			wantContext: true,
		},
		{
			name: "Store open失敗",
			open: func(*context.Context) func(context.Context, string) (*sqlite.Store, error) {
				return func(context.Context, string) (*sqlite.Store, error) {
					return nil, errors.New("open failure")
				}
			},
			wantCode: storageError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userConfigDir = func() (string, error) { return t.TempDir(), nil }
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			var received context.Context
			openSQLiteStore = tt.open(&received)
			parsed, parseError := parseCommand([]string{"search-text", "--query", "channel"})
			if parseError.code != "" {
				t.Fatalf("search-textをparseする: %v", parseError)
			}
			data, executionError, handled := executeRetrieval(ctx, parsed)
			if !handled || executionError.code != tt.wantCode || (data != nil) != tt.wantData {
				t.Fatalf("executeRetrieval() = %#v, %#v, %t", data, executionError, handled)
			}
			if tt.wantContext && received != ctx {
				t.Fatal("composition rootがContextをSQLiteへ伝播しません")
			}
		})
	}
}

func TestDefaultStorePathAndOpenFailures(t *testing.T) {
	originalUserConfigDir := userConfigDir
	originalMakeStoreDirectory := makeStoreDirectory
	originalOpenSQLiteStore := openSQLiteStore
	t.Cleanup(func() {
		userConfigDir = originalUserConfigDir
		makeStoreDirectory = originalMakeStoreDirectory
		openSQLiteStore = originalOpenSQLiteStore
	})

	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "既定Store path",
			test: func(t *testing.T) {
				userConfigDir = func() (string, error) { return "/tmp/config", nil }
				path, err := defaultStorePath()
				if err != nil || path != filepath.Join("/tmp/config", "knowledge-cli", "knowledge.db") {
					t.Fatalf("defaultStorePath() = %q, %v", path, err)
				}
			},
		},
		{
			name: "設定ディレクトリ解決失敗",
			test: func(t *testing.T) {
				userConfigDir = func() (string, error) { return "", errors.New("no config directory") }
				if _, err := defaultStorePath(); err == nil {
					t.Fatal("設定ディレクトリ解決失敗を返しません")
				}
				if _, err := openDefaultRetrievalStore(context.Background()); err == nil {
					t.Fatal("Store初期化時に設定ディレクトリ解決失敗を返しません")
				}
			},
		},
		{
			name: "Storeディレクトリ作成失敗",
			test: func(t *testing.T) {
				userConfigDir = func() (string, error) { return t.TempDir(), nil }
				makeStoreDirectory = func(string, os.FileMode) error { return errors.New("mkdir failure") }
				if _, err := openDefaultRetrievalStore(context.Background()); err == nil {
					t.Fatal("Storeディレクトリ作成失敗を返しません")
				}
			},
		},
		{
			name: "開始前の中断",
			test: func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				if _, err := openDefaultRetrievalStore(ctx); !errors.Is(err, context.Canceled) {
					t.Fatalf("開始前の中断 = %v", err)
				}
			},
		},
		{
			name: "設定ディレクトリ解決後の中断",
			test: func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				userConfigDir = func() (string, error) {
					cancel()

					return t.TempDir(), nil
				}
				makeStoreDirectory = func(string, os.FileMode) error {
					t.Fatal("設定ディレクトリ解決後に中断したのにディレクトリを作成しました")

					return nil
				}
				if _, err := openDefaultRetrievalStore(ctx); !errors.Is(err, context.Canceled) {
					t.Fatalf("設定ディレクトリ解決後の中断 = %v", err)
				}
			},
		},
		{
			name: "ディレクトリ作成後の中断",
			test: func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				userConfigDir = func() (string, error) { return t.TempDir(), nil }
				makeStoreDirectory = func(string, os.FileMode) error {
					cancel()

					return nil
				}
				openSQLiteStore = func(context.Context, string) (*sqlite.Store, error) {
					t.Fatal("ディレクトリ作成後に中断したのにStoreを開きました")

					return nil, nil
				}
				if _, err := openDefaultRetrievalStore(ctx); !errors.Is(err, context.Canceled) {
					t.Fatalf("ディレクトリ作成後の中断 = %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userConfigDir = originalUserConfigDir
			makeStoreDirectory = originalMakeStoreDirectory
			openSQLiteStore = originalOpenSQLiteStore
			tt.test(t)
		})
	}
}

func TestRetrievalResponse(t *testing.T) {
	value := "2026-08-14T00:00:00Z"
	tests := []struct {
		name    string
		value   any
		wantNil bool
	}{
		{
			name: "字句検索結果",
			value: []domain.AssertionSummary{{
				ID:             "asrt_01",
				NormalizedText: "channel send",
				Revision:       1,
				Concepts: []domain.Concept{{
					ID:   "cpt_01",
					Name: "channel",
				}},
				Scope: []domain.Scope{{
					Key:   "language",
					Value: "Go",
				}},
				MatchedFields: []string{"assertion_text"},
			}},
		},
		{
			name: "Concept検索結果",
			value: domain.ConceptSearchResult{
				Concept: &domain.Concept{
					ID:   "cpt_01",
					Name: "channel",
				},
				Results: []domain.AssertionSummary{{
					ID:             "asrt_01",
					NormalizedText: "channel send",
					Revision:       1,
					Scope: []domain.Scope{{
						Key:   "language",
						Value: "Go",
					}},
				}},
			},
		},
		{
			name: "Relation検索結果",
			value: []domain.RelatedResult{{
				RelationID:   "rel_01",
				RelationType: "causes",
				Direction:    "outgoing",
				Target: domain.RelationTarget{
					Kind:           "assertion",
					ID:             "asrt_02",
					NormalizedText: &value,
				},
			}},
		},
		{
			name: "矛盾候補検索結果",
			value: []domain.ContradictionResult{{
				RelationID: "rel_02",
				Direction:  "incoming",
				SeedID:     "asrt_02",
				Target: domain.RelationTarget{
					Kind: "assertion",
					ID:   "asrt_01",
				},
			}},
		},
		{
			name: "Assertion取得結果",
			value: domain.AssertionDetail{
				ID:              "asrt_01",
				CurrentRevision: 1,
				Revisions: []domain.Revision{{
					Number:         1,
					NormalizedText: "channel send",
					Scope: []domain.Scope{{
						Key:   "language",
						Value: "Go",
					}},
					Temporal: &domain.Temporal{
						ValidFrom:    &value,
						ValidUntil:   &value,
						VersionScope: &value,
						ObservedAt:   &value,
						LastVerified: &value,
					},
				}},
				Concepts: []domain.ConceptDetail{{
					ID:      "cpt_01",
					Name:    "channel",
					Aliases: []string{"chan"},
				}},
				Aliases: []domain.AssertionAlias{{
					Kind:  "identifier",
					Value: "ch",
				}},
			},
		},
		{
			name: "Evidence取得結果",
			value: domain.EvidenceResult{
				AssertionID: "asrt_01",
				Evidence: []domain.Evidence{{
					ID:         "evd_01",
					Kind:       "user_code",
					RawText:    "first",
					ObservedAt: value,
				}},
			},
		},
		{
			name: "時点検索結果",
			value: []domain.TemporalSearchResult{{
				AssertionID:    "asrt_01",
				NormalizedText: "channel send",
				Temporal: domain.Temporal{
					ValidFrom:    &value,
					ValidUntil:   &value,
					VersionScope: &value,
					ObservedAt:   &value,
					LastVerified: &value,
				},
			}},
		},
		{
			name:    "未知のresponse",
			value:   "unknown",
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := retrievalResponse(tt.value); (got == nil) != tt.wantNil {
				t.Fatalf("retrievalResponse() = %#v, wantNil %t", got, tt.wantNil)
			}
		})
	}
	response := retrievalResponse(domain.EvidenceResult{
		AssertionID: "asrt_01",
		Evidence: []domain.Evidence{{
			ID: "evd_temporal",
			Temporal: &domain.Temporal{
				VersionScope: &value,
			},
		}},
	}).(map[string]any)
	evidence := response["evidence"].([]map[string]any)
	temporal, ok := evidence[0]["temporal"].(map[string]any)
	versionScope, hasVersionScope := temporal["version_scope"].(*string)
	if !ok || !hasVersionScope || versionScope == nil || *versionScope != value {
		t.Fatalf("Evidence Temporal response = %#v", response)
	}
	for _, tt := range []struct {
		name string
		call func() any
	}{
		{
			name: "nil Concept",
			call: func() any { return conceptResponse(nil) },
		},
		{
			name: "nil Temporal",
			call: func() any { return temporalResponse(nil) },
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.call(); got != nil {
				t.Fatalf("nil値のresponse = %#v, want nil", got)
			}
		})
	}
}

func TestRetrievalCLIError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want errorCode
	}{
		{
			name: "Relation検索起点なし",
			err:  domain.ErrRelationSeedNotFound,
			want: notFoundError,
		},
		{
			name: "Assertionなし",
			err:  domain.ErrAssertionNotFound,
			want: notFoundError,
		},
		{
			name: "保存失敗",
			err:  errors.New("failure"),
			want: storageError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := retrievalCLIError(tt.err); got.code != tt.want {
				t.Fatalf("retrievalCLIError() = %q, want %q", got.code, tt.want)
			}
		})
	}
}

func TestRunWithExecutor(t *testing.T) {
	originalStdout := processStdout
	t.Cleanup(func() { processStdout = originalStdout })
	execute := func(_ context.Context, parsed command) (any, cliError, bool) {
		if parsed.options[0].value == "missing" {
			return nil, cliError{
				code:    notFoundError,
				message: "Assertionが見つかりません",
			}, true
		}
		if parsed.options[0].value == "broken" {
			return nil, cliError{
				code:    storageError,
				message: "Knowledge Storeの読取に失敗しました",
			}, true
		}

		return map[string]any{"assertion_id": "asrt_01"}, cliError{}, true
	}

	tests := []struct {
		name        string
		assertionID string
		wantCode    int
		wantStdout  bool
		wantStderr  bool
	}{
		{
			name:        "成功",
			assertionID: "asrt_01",
			wantStdout:  true,
		},
		{
			name:        "未検出",
			assertionID: "missing",
			wantCode:    3,
			wantStderr:  true,
		},
		{
			name:        "Storage error",
			assertionID: "broken",
			wantCode:    1,
			wantStderr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			processStdout = &stdout
			if got := runWithExecutor(context.Background(), []string{"get", "--assertion-id", tt.assertionID}, &stderr, execute); got != tt.wantCode {
				t.Fatalf("runWithExecutor() = %d, want %d", got, tt.wantCode)
			}
			if (stdout.Len() > 0) != tt.wantStdout || (stderr.Len() > 0) != tt.wantStderr {
				t.Fatalf("stdout/stderr = %q/%q", stdout.String(), stderr.String())
			}
		})
	}
}

type retrievalStoreStub struct {
	getError           error
	relatedError       error
	contradictionError error
	temporalError      error
	receivedContext    *context.Context
}

func (store retrievalStoreStub) SearchText(ctx context.Context, _ string) ([]domain.AssertionSummary, error) {
	store.receive(ctx)

	return []domain.AssertionSummary{{ID: "asrt_01"}}, nil
}

func (store retrievalStoreStub) SearchConcept(ctx context.Context, _ string) (domain.ConceptSearchResult, error) {
	store.receive(ctx)

	return domain.ConceptSearchResult{}, nil
}

func (store retrievalStoreStub) GetAssertion(ctx context.Context, _ string) (domain.AssertionDetail, error) {
	store.receive(ctx)
	if store.getError != nil {
		return domain.AssertionDetail{}, store.getError
	}

	return domain.AssertionDetail{ID: "asrt_01"}, nil
}

func (store retrievalStoreStub) GetEvidence(ctx context.Context, _ string) (domain.EvidenceResult, error) {
	store.receive(ctx)

	return domain.EvidenceResult{AssertionID: "asrt_01"}, nil
}

func (store retrievalStoreStub) SearchRelated(ctx context.Context, _ string, _ string, _ []string) ([]domain.RelatedResult, error) {
	store.receive(ctx)
	if store.relatedError != nil {
		return nil, store.relatedError
	}

	text := "related"

	return []domain.RelatedResult{{
		RelationID:   "rel_01",
		RelationType: "causes",
		Direction:    "outgoing",
		Target: domain.RelationTarget{
			Kind:           "assertion",
			ID:             "asrt_02",
			NormalizedText: &text,
		},
	}}, nil
}

func (store retrievalStoreStub) SearchContradictions(ctx context.Context, _ *string, _ *string) ([]domain.ContradictionResult, error) {
	store.receive(ctx)
	if store.contradictionError != nil {
		return nil, store.contradictionError
	}

	text := "contradiction"

	return []domain.ContradictionResult{{
		RelationID: "rel_02",
		Direction:  "incoming",
		SeedID:     "asrt_02",
		Target: domain.RelationTarget{
			Kind:           "assertion",
			ID:             "asrt_01",
			NormalizedText: &text,
		},
	}}, nil
}

func (store retrievalStoreStub) SearchTemporal(ctx context.Context, _ *string, _ []domain.Scope, _ domain.TemporalSearchFilter) ([]domain.TemporalSearchResult, error) {
	store.receive(ctx)
	if store.temporalError != nil {
		return nil, store.temporalError
	}

	return []domain.TemporalSearchResult{}, nil
}

func (store retrievalStoreStub) receive(ctx context.Context) {
	if store.receivedContext != nil {
		*store.receivedContext = ctx
	}
}

func TestTemporalSearchFilterCanonicalizesUTC(t *testing.T) {
	tests := []struct {
		name      string
		values    map[string][]string
		wantAt    string
		wantFrom  string
		wantUntil string
		wantCode  errorCode
	}{
		{
			name:   "時点を正規化する",
			values: map[string][]string{"at": {"2026-08-14T00:00:00Z"}},
			wantAt: "2026-08-14T00:00:00.000000000Z",
		},
		{
			name: "期間を正規化する",
			values: map[string][]string{
				"valid-from":  {"2026-08-14T00:00:00Z"},
				"valid-until": {"2026-08-14T00:00:00.1Z"},
			},
			wantFrom:  "2026-08-14T00:00:00.000000000Z",
			wantUntil: "2026-08-14T00:00:00.100000000Z",
		},
		{
			name:     "時刻形式が不正",
			values:   map[string][]string{"at": {"today"}},
			wantCode: validationError,
		},
		{
			name:     "UTCではない",
			values:   map[string][]string{"at": {"2026-08-14T00:00:00+09:00"}},
			wantCode: validationError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := temporalSearchFilter(tt.values)
			if err.code != tt.wantCode {
				t.Fatalf("temporalSearchFilter() error = %q, want %q", err.code, tt.wantCode)
			}
			if tt.wantCode != "" {
				return
			}
			if temporalFilterValue(filter.At) != tt.wantAt || temporalFilterValue(filter.ValidFrom) != tt.wantFrom || temporalFilterValue(filter.ValidUntil) != tt.wantUntil {
				t.Fatalf("filter = %#v", filter)
			}
		})
	}
}

func temporalFilterValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
