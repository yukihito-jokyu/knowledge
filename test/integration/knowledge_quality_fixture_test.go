package integration_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/yukihito-jokyu/knowledge/internal/persistence/sqlite"
	_ "modernc.org/sqlite"
)

type qualityFixture struct {
	Feature string                 `json:"feature"`
	Seeds   map[string]qualitySeed `json:"seeds"`
	Cases   []qualityCase          `json:"cases"`
}
type qualitySeed struct {
	Assertions []qualityAssertion `json:"assertions"`
	Relations  []qualityRelation  `json:"relations"`
}
type qualityAssertion struct {
	ID       string            `json:"id"`
	Revision int               `json:"revision"`
	Text     string            `json:"text"`
	Prior    []qualityRevision `json:"prior_revisions"`
	Evidence []qualityEvidence `json:"evidence"`
	Scope    map[string]string `json:"scope"`
	Temporal struct {
		VersionScope string `json:"version_scope"`
		ValidFrom    string `json:"valid_from"`
	} `json:"temporal"`
}
type qualityRevision struct {
	Revision int    `json:"revision"`
	Text     string `json:"text"`
}
type qualityEvidence struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Text       string `json:"text"`
	ObservedAt string `json:"observed_at"`
	Temporal   struct {
		VersionScope string `json:"version_scope"`
	} `json:"temporal"`
}
type qualityRelation struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Type   string `json:"type"`
	Target string `json:"target"`
}
type qualityCase struct {
	CaseID   string          `json:"case_id"`
	Layers   []string        `json:"layers"`
	SeedID   string          `json:"seed_id"`
	Input    json.RawMessage `json:"input"`
	Expected struct {
		Assessment   string            `json:"assessment"`
		Confidence   string            `json:"confidence"`
		TraceStop    string            `json:"trace_stop"`
		Assessments  map[string]string `json:"assessments"`
		CandidateIDs []string          `json:"candidate_ids"`
		UpdateStatus string            `json:"update_status"`
		Operations   []string          `json:"cli_operations"`
		PartialTrace struct {
			Operation string `json:"operation"`
			Stop      string `json:"stop"`
			ExitCode  int    `json:"exit_code"`
		} `json:"partial_trace"`
		Trace struct {
			Operations  []string `json:"operations"`
			Queries     []string `json:"queries"`
			ResultIDs   []string `json:"result_ids"`
			EvidenceIDs []string `json:"evidence_ids"`
			BudgetUsed  int      `json:"budget_used"`
			Stop        string   `json:"stop"`
		} `json:"search_trace"`
		CLI struct {
			Arguments []string `json:"arguments"`
			ResultIDs []string `json:"result_ids"`
			ErrorCode string   `json:"error_code"`
			ExitCode  int      `json:"exit_code"`
		} `json:"cli"`
		StoreDiff struct {
			None   bool     `json:"none"`
			Retain []string `json:"retain"`
			Add    []string `json:"add"`
		} `json:"store_diff"`
		NotExecuted []string `json:"not_executed_layers"`
	} `json:"expected"`
	Contracts []string `json:"contract_references"`
}

func TestKnowledgeQualityFixtureCasesAtProcessBoundary(t *testing.T) {
	fixture := readQualityFixture(t)
	validateQualityFixture(t, fixture)
	binary, selected, found := buildCLI(t, true), os.Getenv("KNOWLEDGE_ACCEPTANCE_CASE_ID"), false
	results := map[qualityCaseResultKey]qualityCaseResult{}
	for _, testCase := range fixture.Cases {
		t.Run(testCase.CaseID, func(t *testing.T) {
			if selected != "" && selected != testCase.CaseID {
				t.Skip("選択caseではありません")
			}
			found = true
			root := t.TempDir()
			store := defaultStoreConfiguration(t, root)
			seedQualityStore(t, store.Path, fixture.Seeds[testCase.SeedID])
			before := qualityStateOf(t, store.Path)
			result, caseErr := executeQualityCaseLayers(testCase, qualityCaseLayerChecks(t, binary, store, testCase, before))
			if err := qualityStoreDiffMatches(t, store, before, testCase); err != nil {
				caseErr = errors.Join(caseErr, err)
			}
			if err := os.RemoveAll(root); err != nil {
				caseErr = errors.Join(caseErr, err)
			}
			if _, err := os.Stat(root); !os.IsNotExist(err) {
				caseErr = errors.Join(caseErr, fmt.Errorf("case Storeを破棄できません: %w", err))
			}
			if err := recordQualityCaseResult(results, result); err != nil {
				caseErr = errors.Join(caseErr, err)
			}
			t.Logf("CLI Case Result: case_id=%s status=%s first_mismatch_layer=%s not_executed_layers=%v", result.CaseID, result.Status, result.FirstMismatchLayer, result.NotExecutedLayers)
			if caseErr != nil {
				t.Error(caseErr)
			}
		})
	}
	if selected != "" && !found {
		t.Fatalf("選択caseがfixtureにありません: %s", selected)
	}
}

func qualityCaseLayerChecks(t *testing.T, binary string, store defaultStore, c qualityCase, before qualityState) []qualityLayerCheck {
	t.Helper()
	checks := []qualityLayerCheck{}
	if c.CaseID == "FEAT005-X-SEARCH-TECHNICAL-FAILURE" {
		return qualityTechnicalFailureChecks(t, binary, c)
	}
	if c.CaseID == "FEAT005-H-UPDATE-CORRECTION" {
		return qualityCorrectionChecks(t, binary, store, c, before)
	}
	if len(c.Expected.CLI.Arguments) == 0 {
		return qualityNoCandidateChecks(t, store, c, before)
	}
	checks = append(checks, qualityCLIAndSearchChecks(t, binary, store, c)...)
	if containsQualityLayer(c.Layers, "end_to_end") {
		checks = append(checks, qualityLayerCheck{
			Layer: "end_to_end",
			Check: qualityEndToEndCheck(t, store, before, c),
		})
	}

	return checks
}

func qualityEndToEndCheck(t *testing.T, store defaultStore, before qualityState, c qualityCase) func() error {
	t.Helper()

	return func() error {
		return qualityStoreDiffMatches(t, store, before, c)
	}
}

// qualityStoreDiffMatchesは、層の有無にかかわらずfixtureのStore不変条件を照合する。
func qualityStoreDiffMatches(t *testing.T, store defaultStore, before qualityState, c qualityCase) error {
	t.Helper()
	if c.CaseID == "FEAT005-H-UPDATE-CORRECTION" {
		return qualityCorrectionStoreMatches(t, store, before, c)
	}
	if c.Expected.StoreDiff.None && !reflect.DeepEqual(before, qualityStateOf(t, store.Path)) {
		return errors.New("store_diff:none のcaseがStoreを変更しました")
	}

	return nil
}

// qualityCorrectionStoreMatchesは、訂正のrevisionとEvidenceを実Storeから照合する。
func qualityCorrectionStoreMatches(t *testing.T, store defaultStore, before qualityState, c qualityCase) error {
	t.Helper()
	var input struct {
		Episode struct {
			UserContributions []struct {
				SourceText string `json:"source_text"`
			} `json:"user_contributions"`
		} `json:"episode"`
	}
	if err := json.Unmarshal(c.Input, &input); err != nil {
		return err
	}
	if c.Expected.StoreDiff.None || len(input.Episode.UserContributions) != 1 {
		return fmt.Errorf("H Store diff oracleが不正: %#v", c.Expected.StoreDiff)
	}
	after := qualityStateOf(t, store.Path)
	for _, id := range c.Expected.StoreDiff.Retain {
		if id == "rev-h-1" {
			if !qualityRevisionExists(t, store.Path, "as-h", 1) {
				return errors.New("保持revisionが消失しました")
			}

			continue
		}
		if !after.IDs[id] {
			return fmt.Errorf("保持対象が消失しました: %s", id)
		}
	}
	if !reflect.DeepEqual(c.Expected.StoreDiff.Add, []string{
		"rev-h-2",
		"ev-h-correction",
	}) {
		return fmt.Errorf("H add oracle = %v", c.Expected.StoreDiff.Add)
	}
	for _, id := range c.Expected.StoreDiff.Add {
		switch id {
		case "rev-h-2":
			if !qualityRevisionExists(t, store.Path, "as-h", 2) {
				return errors.New("追加revisionがありません")
			}
		case "ev-h-correction":
			// Evidence IDは公開契約上不透明なので、fixtureの論理IDを入力由来の実IDへ対応付ける。
			evidenceID, err := qualityCorrectionEvidenceID(store.Path, input.Episode.UserContributions[0].SourceText)
			if err != nil {
				return err
			}
			if !after.IDs[evidenceID] {
				return fmt.Errorf("追加Evidenceがありません: %s", evidenceID)
			}
		default:
			return fmt.Errorf("未知のH add対象: %s", id)
		}
	}
	if after.Revisions["as-h"] != before.Revisions["as-h"]+1 {
		return fmt.Errorf("revision = %d, want %d", after.Revisions["as-h"], before.Revisions["as-h"]+1)
	}
	if after.Evidence["as-h\x00correction\x00"+input.Episode.UserContributions[0].SourceText] != 1 {
		return errors.New("訂正Evidenceが追加されません")
	}

	return nil
}

func qualityCorrectionEvidenceID(path string, text string) (string, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return "", err
	}
	defer db.Close()
	var evidenceID string
	if err := db.QueryRowContext(context.Background(), "SELECT evidence_id FROM evidence WHERE assertion_id = ? AND kind = ? AND raw_text = ?", "as-h", "correction", text).Scan(&evidenceID); err != nil {
		return "", err
	}

	return evidenceID, nil
}

func qualityCLIAndSearchChecks(t *testing.T, binary string, store defaultStore, c qualityCase) []qualityLayerCheck {
	t.Helper()
	var stdout string

	return []qualityLayerCheck{
		{
			Layer: "cli_store",
			Check: func() error {
				var stderr string
				var err error
				if c.CaseID == "FEAT005-X-SEARCH-CANCELED" {
					stdout, stderr, err = runInterruptedCommand(t, binary, store.Environment, c.Expected.CLI.Arguments, "search-text")
				} else {
					stdout, stderr, err = runCommand(binary, store.Environment, c.Expected.CLI.Arguments)
				}
				if !isExitCode(err, c.Expected.CLI.ExitCode) {
					if err != nil {
						return fmt.Errorf("exit = %w, want %d", err, c.Expected.CLI.ExitCode)
					}

					return fmt.Errorf("exit = 0, want %d", c.Expected.CLI.ExitCode)
				}
				if c.Expected.CLI.ExitCode == 130 {
					if stdout != "" || stderr != "" {
						return fmt.Errorf("exit 130 output = stdout=%q stderr=%q", stdout, stderr)
					}

					return nil
				}
				if stderr != "" {
					return fmt.Errorf("stderr = %q", stderr)
				}

				return nil
			},
		},
		{
			Layer: "knowledge_search",
			Check: func() error {
				if c.Expected.CLI.ExitCode == 130 {
					return nil
				}
				if err := qualityResultIDsMatch(stdout, c.Expected.CLI.ResultIDs); err != nil {
					return err
				}
				if c.CaseID == "FEAT005-E-SEARCH-OUTDATED" {
					if err := qualityEvidenceTemporalMatches(store.Path); err != nil {
						return err
					}
					if err := qualityEvidenceTemporalResponseMatches(binary, store); err != nil {
						return err
					}
				}

				return nil
			},
		},
	}
}

func qualityResultIDsMatch(stdout string, want []string) error {
	var response struct {
		Data struct {
			Results []struct {
				ID string `json:"assertion_id"`
			} `json:"results"`
			AssertionID string `json:"assertion_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(stdout), &response); err != nil {
		return err
	}
	got := make([]string, 0, len(response.Data.Results))
	for _, result := range response.Data.Results {
		got = append(got, result.ID)
	}
	if response.Data.AssertionID != "" {
		got = []string{response.Data.AssertionID}
	}
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		return fmt.Errorf("result IDs = %v, want %v", got, want)
	}

	return nil
}

func qualityNoCandidateChecks(t *testing.T, store defaultStore, c qualityCase, before qualityState) []qualityLayerCheck {
	t.Helper()

	return []qualityLayerCheck{
		{
			Layer: "knowledge_acquisition",
			Check: func() error {
				if len(c.Expected.CandidateIDs) != 0 {
					return errors.New("Candidateが空ではありません")
				}
				var input struct {
					Episode struct {
						UserContributions []struct {
							SourceText string `json:"source_text"`
						} `json:"user_contributions"`
					} `json:"episode"`
					AIResponse *struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"ai_response"`
				}
				if err := json.Unmarshal(c.Input, &input); err != nil {
					return err
				}
				if c.CaseID == "FEAT005-F-ACQUISITION-QUESTION" && len(input.Episode.UserContributions) != 1 {
					return errors.New("Fの固定質問入力を消費できません")
				}
				if c.CaseID == "FEAT005-G-ACQUISITION-AI-ONLY" && (len(input.Episode.UserContributions) != 0 || input.AIResponse == nil || input.AIResponse.ID != "ai-g" || input.AIResponse.Text == "") {
					return errors.New("Gの固定AI入力を消費できません")
				}

				return nil
			},
		},
		{
			Layer: "knowledge_update",
			Check: func() error {
				if c.Expected.UpdateStatus != "completed" || len(c.Expected.Operations) != 0 {
					return fmt.Errorf("Update oracle = %#v", c.Expected)
				}

				return nil
			},
		},
		{
			Layer: "end_to_end",
			Check: qualityEndToEndCheck(t, store, before, c),
		},
	}
}

func qualityTechnicalFailureChecks(t *testing.T, binary string, c qualityCase) []qualityLayerCheck {
	t.Helper()
	var stderr string

	return []qualityLayerCheck{
		{
			Layer: "cli_store",
			Check: func() error {
				root := filepath.Join(t.TempDir(), "not-a-directory")
				if err := os.WriteFile(root, []byte("blocked"), 0o600); err != nil {
					return err
				}
				bad := defaultStoreConfiguration(t, root)
				stdout, actualStderr, err := runCommand(binary, bad.Environment, c.Expected.CLI.Arguments)
				stderr = actualStderr
				if !isExitCode(err, c.Expected.CLI.ExitCode) || stdout != "" {
					if err != nil {
						return fmt.Errorf("technical failure = stdout=%q stderr=%q err=%w", stdout, stderr, err)
					}

					return fmt.Errorf("technical failure = stdout=%q stderr=%q exit=0", stdout, stderr)
				}

				return nil
			},
		},
		{
			Layer: "knowledge_search",
			Check: func() error {
				var response struct {
					Error struct {
						Code string `json:"code"`
					} `json:"error"`
				}
				if err := json.Unmarshal([]byte(stderr), &response); err != nil || response.Error.Code != c.Expected.CLI.ErrorCode {
					return fmt.Errorf("error code = %q, want %q", response.Error.Code, c.Expected.CLI.ErrorCode)
				}

				return nil
			},
		},
	}
}

func qualityCorrectionChecks(t *testing.T, binary string, store defaultStore, c qualityCase, before qualityState) []qualityLayerCheck {
	t.Helper()

	return []qualityLayerCheck{
		{
			Layer: "cli_store",
			Check: func() error {
				stdout, stderr, err := runCommand(binary, store.Environment, c.Expected.CLI.Arguments)
				if err != nil || stderr != "" {
					if err != nil {
						return fmt.Errorf("H get: %w stderr=%q", err, stderr)
					}

					return fmt.Errorf("H get stderr=%q", stderr)
				}

				return qualityResultIDsMatch(stdout, c.Expected.CLI.ResultIDs)
			},
		},
		{
			Layer: "knowledge_acquisition",
			Check: func() error {
				if !reflect.DeepEqual(c.Expected.CandidateIDs, []string{"cand-h-1"}) {
					return fmt.Errorf("H candidate = %v", c.Expected.CandidateIDs)
				}

				return nil
			},
		},
		{
			Layer: "knowledge_update",
			Check: func() error {
				var input struct {
					Episode struct {
						UserContributions []struct {
							SourceText string `json:"source_text"`
							ObservedAt string `json:"observed_at"`
						} `json:"user_contributions"`
					} `json:"episode"`
					Claim struct {
						Text string `json:"text"`
					} `json:"claim"`
				}
				if err := json.Unmarshal(c.Input, &input); err != nil {
					return err
				}
				if len(input.Episode.UserContributions) != 1 || input.Claim.Text == "" {
					return fmt.Errorf("H update input = %#v", input)
				}
				if _, stderr, err := runCommand(binary, store.Environment, []string{"revise", "--assertion-id", "as-h", "--normalized-text", input.Claim.Text}); err != nil || stderr != "" {
					if err != nil {
						return fmt.Errorf("H revise: %w stderr=%q", err, stderr)
					}

					return fmt.Errorf("H revise stderr=%q", stderr)
				}
				if _, stderr, err := runCommand(binary, store.Environment, []string{"attach-evidence", "--assertion-id", "as-h", "--evidence-kind", "correction", "--evidence-text", input.Episode.UserContributions[0].SourceText, "--evidence-observed-at", input.Episode.UserContributions[0].ObservedAt}); err != nil || stderr != "" {
					if err != nil {
						return fmt.Errorf("H attach-evidence: %w stderr=%q", err, stderr)
					}

					return fmt.Errorf("H attach-evidence stderr=%q", stderr)
				}

				return nil
			},
		},
		{
			Layer: "end_to_end",
			Check: qualityEndToEndCheck(t, store, before, c),
		},
	}
}

func qualityEvidenceTemporalMatches(path string) error {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	defer db.Close()
	rows, err := db.QueryContext(context.Background(), "SELECT evidence_id, version_scope FROM evidence_temporal_metadata ORDER BY evidence_id")
	if err != nil {
		return err
	}
	defer rows.Close()
	got := map[string]string{}
	for rows.Next() {
		var id, version string
		if err := rows.Scan(&id, &version); err != nil {
			return err
		}
		got[id] = version
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if !reflect.DeepEqual(got, map[string]string{
		"ev-e-old":        "go1.21",
		"ev-e-correction": "go1.22+",
	}) {
		return fmt.Errorf("Evidence Temporal = %#v", got)
	}

	return nil
}

func qualityEvidenceTemporalResponseMatches(binary string, store defaultStore) error {
	stdout, stderr, err := runCommand(binary, store.Environment, []string{"get-evidence", "--assertion-id", "as-e"})
	if err != nil {
		return err
	}
	if stderr != "" {
		return fmt.Errorf("get-evidence stderr=%q", stderr)
	}
	var response struct {
		Data struct {
			Evidence []struct {
				ID       string `json:"evidence_id"`
				Temporal struct {
					VersionScope string `json:"version_scope"`
				} `json:"temporal"`
			} `json:"evidence"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(stdout), &response); err != nil {
		return err
	}
	got := map[string]string{}
	for _, evidence := range response.Data.Evidence {
		got[evidence.ID] = evidence.Temporal.VersionScope
	}
	want := map[string]string{
		"ev-e-old":        "go1.21",
		"ev-e-correction": "go1.22+",
	}
	if !reflect.DeepEqual(got, want) {
		return fmt.Errorf("get-evidence temporal = %#v, want %#v", got, want)
	}

	return nil
}

func hasQualitySearchLayer(c qualityCase) bool {
	for _, layer := range c.Layers {
		if layer == "knowledge_search" {
			return true
		}
	}

	return false
}

func qualityRevisionExists(t *testing.T, path, assertionID string, revision int) bool {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var count int
	if err := db.QueryRowContext(context.Background(), "SELECT count(*) FROM assertion_revisions WHERE assertion_id = ? AND revision = ?", assertionID, revision).Scan(&count); err != nil {
		t.Fatal(err)
	}

	return count == 1
}

type qualityState struct {
	IDs       map[string]bool
	Revisions map[string]int
	Evidence  map[string]int
}

func qualityStateOf(t *testing.T, path string) qualityState {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	state := qualityState{
		IDs:       map[string]bool{},
		Revisions: map[string]int{},
		Evidence:  map[string]int{},
	}
	assertionRows, err := db.QueryContext(context.Background(), "SELECT assertion_id, current_revision FROM assertions")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := assertionRows.Close(); err != nil {
			t.Error(err)
		}
	}()
	for assertionRows.Next() {
		var id string
		var revision int
		if err := assertionRows.Scan(&id, &revision); err != nil {
			t.Fatal(err)
		}
		state.IDs[id] = true
		state.Revisions[id] = revision
	}
	if err := assertionRows.Err(); err != nil {
		t.Fatal(err)
	}
	evidenceRows, err := db.QueryContext(context.Background(), "SELECT evidence_id, assertion_id, kind, raw_text FROM evidence")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := evidenceRows.Close(); err != nil {
			t.Error(err)
		}
	}()
	for evidenceRows.Next() {
		var id, assertionID, kind, text string
		if err := evidenceRows.Scan(&id, &assertionID, &kind, &text); err != nil {
			t.Fatal(err)
		}
		state.IDs[id] = true
		state.Evidence[assertionID+"\x00"+kind+"\x00"+text]++
	}
	if err := evidenceRows.Err(); err != nil {
		t.Fatal(err)
	}
	relationRows, err := db.QueryContext(context.Background(), "SELECT relation_id FROM relations")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := relationRows.Close(); err != nil {
			t.Error(err)
		}
	}()
	for relationRows.Next() {
		var id string
		if err := relationRows.Scan(&id); err != nil {
			t.Fatal(err)
		}
		state.IDs[id] = true
	}
	if err := relationRows.Err(); err != nil {
		t.Fatal(err)
	}

	return state
}

func seedQualityStore(t *testing.T, path string, seed qualitySeed) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	store, err := sqlite.Open(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			t.Error(err)
		}
	}()
	for _, a := range seed.Assertions {
		revision := a.Revision
		if revision == 0 {
			revision = 1
		}
		if _, err := tx.ExecContext(context.Background(), "INSERT INTO assertions(assertion_id,current_revision,created_at) VALUES(?,?,?)", a.ID, revision, "2026-01-01T00:00:00Z"); err != nil {
			t.Fatal(err)
		}
		revisions := append(append([]qualityRevision{}, a.Prior...), qualityRevision{revision, a.Text})
		for _, r := range revisions {
			if _, err := tx.ExecContext(context.Background(), "INSERT INTO assertion_revisions(assertion_id,revision,normalized_text,created_at) VALUES(?,?,?,?)", a.ID, r.Revision, r.Text, "2026-01-01T00:00:00Z"); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := tx.ExecContext(context.Background(), "INSERT INTO assertion_lexical_index(assertion_id,normalized_text,concept_name,concept_alias,scope_key,scope_value,assertion_alias) VALUES(?,?,?,?,?,?,?)", a.ID, a.Text, "channel", "", "language", "Go", ""); err != nil {
			t.Fatal(err)
		}
		for key, value := range a.Scope {
			if _, err := tx.ExecContext(context.Background(), "INSERT INTO revision_scopes(assertion_id,revision,scope_key,scope_value) VALUES(?,?,?,?)", a.ID, revision, key, value); err != nil {
				t.Fatal(err)
			}
		}
		if a.Temporal.VersionScope != "" || a.Temporal.ValidFrom != "" {
			if _, err := tx.ExecContext(context.Background(), "INSERT INTO temporal_metadata(assertion_id,revision,valid_from,version_scope) VALUES(?,?,?,?)", a.ID, revision, nullableQualityString(a.Temporal.ValidFrom), nullableQualityString(a.Temporal.VersionScope)); err != nil {
				t.Fatal(err)
			}
		}
		for _, e := range a.Evidence {
			if _, err := tx.ExecContext(context.Background(), "INSERT INTO evidence(evidence_id,assertion_id,kind,raw_text,observed_at,created_at) VALUES(?,?,?,?,?,?)", e.ID, a.ID, e.Kind, e.Text, e.ObservedAt, "2026-01-01T00:00:00Z"); err != nil {
				t.Fatal(err)
			}
			if e.Temporal.VersionScope != "" {
				if _, err := tx.ExecContext(context.Background(), "INSERT INTO evidence_temporal_metadata(evidence_id,version_scope) VALUES(?,?)", e.ID, e.Temporal.VersionScope); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
	for _, r := range seed.Relations {
		if _, err := tx.ExecContext(context.Background(), "INSERT INTO relations(relation_id,source_kind,source_id,relation_type,target_kind,target_id,created_at) VALUES(?,?,?,?,?,?,?)", r.ID, "assertion", r.Source, r.Type, "assertion", r.Target, "2026-01-01T00:00:00Z"); err != nil {
			t.Fatal(err)
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func nullableQualityString(value string) any {
	if value == "" {
		return nil
	}

	return value
}

func validateQualityFixture(t *testing.T, f qualityFixture) {
	t.Helper()
	if f.Feature != "FEAT-005" || len(f.Cases) != 13 {
		t.Fatal("fixture headerまたはcase数が不正")
	}
	seen := map[string]bool{}
	for _, c := range f.Cases {
		if c.CaseID == "" || seen[c.CaseID] || len(c.Layers) == 0 || c.Input == nil || len(c.Contracts) == 0 || len(c.Expected.NotExecuted) == 0 {
			t.Fatalf("fixture契約が欠落: %s", c.CaseID)
		}
		if _, ok := f.Seeds[c.SeedID]; !ok {
			t.Fatalf("seedがない: %s", c.SeedID)
		}
		executed := make(map[string]bool, len(c.Layers))
		for _, layer := range c.Layers {
			executed[layer] = true
		}
		for _, layer := range c.Expected.NotExecuted {
			if executed[layer] {
				t.Fatalf("executed_layersとnot_executed_layersが重複: %s: %s", c.CaseID, layer)
			}
		}
		seen[c.CaseID] = true
	}
	for _, id := range []string{"FEAT005-A-SEARCH-EMPTY", "FEAT005-B-SEARCH-EXACT", "FEAT005-C-SEARCH-PARTIAL", "FEAT005-D-SEARCH-CONTRADICTED", "FEAT005-E-SEARCH-OUTDATED", "FEAT005-F-ACQUISITION-QUESTION", "FEAT005-G-ACQUISITION-AI-ONLY", "FEAT005-H-UPDATE-CORRECTION", "FEAT005-H-SEARCH-CORRECTED", "FEAT005-I-READ-SELECTED", "FEAT005-J-READ-TRIVIAL", "FEAT005-X-SEARCH-TECHNICAL-FAILURE", "FEAT005-X-SEARCH-CANCELED"} {
		if !seen[id] {
			t.Fatalf("catalog固定caseがありません: %s", id)
		}
	}
}

func TestKnowledgeQualityReadingValueReferences(t *testing.T) {
	fixture := readQualityFixture(t)
	verification, err := os.ReadFile(filepath.Join("..", "..", "skills", "reading-value", "references", "verification.md"))
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{
		"FEAT005-A-SEARCH-EMPTY":        "FEAT-003 V-002",
		"FEAT005-G-ACQUISITION-AI-ONLY": "FEAT-003 V-004",
		"FEAT005-I-READ-SELECTED":       "FEAT-003 V-002",
		"FEAT005-J-READ-TRIVIAL":        "FEAT-003 V-002",
	}
	for _, testCase := range fixture.Cases {
		want, required := expected[testCase.CaseID]
		if !required {
			continue
		}
		var input struct {
			ReadingValueReference string `json:"reading_value_reference"`
		}
		if err := json.Unmarshal(testCase.Input, &input); err != nil {
			t.Fatal(err)
		}
		if input.ReadingValueReference != want || !containsQualityContract(testCase.Contracts, strings.ReplaceAll(want, " ", "/")) {
			t.Fatalf("Reading Value参照が一意に対応しません: %s", testCase.CaseID)
		}
		section := strings.Replace(want, "FEAT-003 ", "### ", 1)
		if !strings.Contains(string(verification), section) {
			t.Fatalf("Reading Value検証契約に参照先がありません: %s", want)
		}
	}
}

func containsQualityContract(contracts []string, want string) bool {
	for _, contract := range contracts {
		if contract == want {
			return true
		}
	}

	return false
}
func readQualityFixture(t *testing.T) qualityFixture {
	t.Helper()
	source, err := os.ReadFile(filepath.Join("..", "..", "testdata", "fixtures", "acceptance", "knowledge-quality", "cases.json"))
	if err != nil {
		t.Fatal(err)
	}
	var f qualityFixture
	if err := json.Unmarshal(source, &f); err != nil {
		t.Fatal(err)
	}

	return f
}
