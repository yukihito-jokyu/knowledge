package integration_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
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
			runQualityCase(t, binary, store, testCase)
			if testCase.Expected.StoreDiff.None && !reflect.DeepEqual(before, qualityStateOf(t, store.Path)) {
				t.Fatal("store_diff:none のcaseがStoreを変更しました")
			}
			if err := os.RemoveAll(root); err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(root); !os.IsNotExist(err) {
				t.Fatalf("case Storeを破棄できません: %v", err)
			}
		})
	}
	if selected != "" && !found {
		t.Fatalf("選択caseがfixtureにありません: %s", selected)
	}
}

func runQualityCase(t *testing.T, binary string, store defaultStore, c qualityCase) {
	t.Helper()
	assertFixtureOracle(t, c)
	if c.CaseID == "FEAT005-X-SEARCH-TECHNICAL-FAILURE" {
		root := filepath.Join(t.TempDir(), "not-a-directory")
		if err := os.WriteFile(root, []byte("blocked"), 0o600); err != nil {
			t.Fatal(err)
		}
		bad := defaultStoreConfiguration(t, root)
		stdout, stderr, err := runCommand(binary, bad.Environment, c.Expected.CLI.Arguments)
		if !isExitCode(err, c.Expected.CLI.ExitCode) {
			t.Fatalf("exit = %v", err)
		}
		assertStdout(t, stdout, "")
		assertQualityErrorCode(t, stderr, c.Expected.CLI.ErrorCode)

		return
	}
	if c.CaseID == "FEAT005-H-UPDATE-CORRECTION" {
		assertQualityCorrection(t, binary, store, c)

		return
	}
	var stdout, stderr string
	var err error
	if len(c.Expected.CLI.Arguments) == 0 {
		assertNoCandidateEpisode(t, c)

		return
	}
	if c.CaseID == "FEAT005-X-SEARCH-CANCELED" {
		stdout, stderr, err = runInterruptedCommand(t, binary, store.Environment, c.Expected.CLI.Arguments, "search-text")
	} else {
		stdout, stderr, err = runCommand(binary, store.Environment, c.Expected.CLI.Arguments)
	}
	if !isExitCode(err, c.Expected.CLI.ExitCode) {
		t.Fatalf("exit = %v, want %d", err, c.Expected.CLI.ExitCode)
	}
	if c.Expected.CLI.ExitCode == 130 {
		assertStdout(t, stdout, "")
		assertStderr(t, stderr, nil)

		return
	}
	assertStderr(t, stderr, nil)
	assertQualityResultIDs(t, stdout, c.Expected.CLI.ResultIDs)
	if c.CaseID == "FEAT005-E-SEARCH-OUTDATED" {
		assertQualityEvidenceTemporal(t, store.Path)
		assertQualityEvidenceTemporalResponse(t, binary, store)
	}
}

func assertQualityEvidenceTemporal(t *testing.T, path string) {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	rows, err := db.QueryContext(context.Background(), "SELECT evidence_id, version_scope FROM evidence_temporal_metadata ORDER BY evidence_id")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			t.Error(err)
		}
	}()
	got := map[string]string{}
	for rows.Next() {
		var id, version string
		if err := rows.Scan(&id, &version); err != nil {
			t.Fatal(err)
		}
		got[id] = version
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, map[string]string{
		"ev-e-old":        "go1.21",
		"ev-e-correction": "go1.22+",
	}) {
		t.Fatalf("Evidence Temporal = %#v", got)
	}
}

func assertQualityEvidenceTemporalResponse(t *testing.T, binary string, store defaultStore) {
	t.Helper()
	stdout, stderr, err := runCommand(binary, store.Environment, []string{"get-evidence", "--assertion-id", "as-e"})
	if err != nil {
		t.Fatal(err)
	}
	assertStderr(t, stderr, nil)
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
		t.Fatal(err)
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
		t.Fatalf("get-evidence temporal = %#v, want %#v", got, want)
	}
}

func assertFixtureOracle(t *testing.T, c qualityCase) {
	t.Helper()
	if len(c.Expected.NotExecuted) == 0 {
		t.Fatal("not_executed_layersがありません")
	}
	if c.CaseID == "FEAT005-X-SEARCH-TECHNICAL-FAILURE" || c.CaseID == "FEAT005-X-SEARCH-CANCELED" {
		if c.Expected.PartialTrace.Operation != "search-text" || c.Expected.PartialTrace.Stop == "" || !reflect.DeepEqual(c.Expected.NotExecuted, []string{"knowledge_acquisition", "knowledge_update", "end_to_end"}) {
			t.Fatalf("X partial traceが不正: %#v", c.Expected.PartialTrace)
		}

		return
	}
	if c.Expected.Assessment != "" && (c.Expected.Confidence == "" || c.Expected.TraceStop == "") {
		t.Fatalf("search oracleが不完全: %s", c.CaseID)
	}
}

func assertNoCandidateEpisode(t *testing.T, c qualityCase) {
	t.Helper()
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
		t.Fatal(err)
	}
	if len(c.Expected.CandidateIDs) != 0 || c.Expected.UpdateStatus != "completed" || len(c.Expected.Operations) != 0 {
		t.Fatalf("空Candidate oracleが不正: %#v", c.Expected)
	}
	if c.CaseID == "FEAT005-F-ACQUISITION-QUESTION" && len(input.Episode.UserContributions) != 1 {
		t.Fatal("Fの固定質問入力を消費できません")
	}
	if c.CaseID == "FEAT005-G-ACQUISITION-AI-ONLY" && len(input.Episode.UserContributions) != 0 {
		t.Fatal("GのAI-only入力を消費できません")
	}
	if c.CaseID == "FEAT005-G-ACQUISITION-AI-ONLY" && (input.AIResponse == nil || input.AIResponse.ID != "ai-g" || input.AIResponse.Text == "") {
		t.Fatal("Gの固定AI入力を消費できません")
	}
}

func assertQualityCorrection(t *testing.T, binary string, store defaultStore, c qualityCase) {
	t.Helper()
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
		t.Fatal(err)
	}
	wantOperations := []string{"search-text", "search-text", "get", "get-evidence", "revise", "attach-evidence"}
	if !reflect.DeepEqual(c.Expected.Operations, wantOperations) {
		t.Fatalf("H operations = %v, want %v", c.Expected.Operations, wantOperations)
	}
	if !reflect.DeepEqual(c.Expected.CandidateIDs, []string{"cand-h-1"}) || !reflect.DeepEqual(c.Expected.CLI.Arguments, []string{"get", "--assertion-id", "as-h"}) {
		t.Fatal("H fixture oracleが固定契約と一致しません")
	}
	before := qualityStateOf(t, store.Path)
	for _, operation := range []struct {
		arguments []string
		resultIDs []string
	}{
		{arguments: []string{"search-text", "--query", input.Claim.Text}, resultIDs: []string{}},
		{arguments: []string{"search-text", "--query", "unbuffered channelのsendはreceiverが受信可能になる前にも完了する"}, resultIDs: []string{"as-h"}},
		{arguments: []string{"get", "--assertion-id", "as-h"}},
		{arguments: []string{"get-evidence", "--assertion-id", "as-h"}},
	} {
		stdout, stderr, err := runCommand(binary, store.Environment, operation.arguments)
		if err != nil {
			t.Fatalf("更新前read %v: %v", operation.arguments[0], err)
		}
		assertStderr(t, stderr, nil)
		if operation.resultIDs != nil {
			assertQualityResultIDs(t, stdout, operation.resultIDs)
		}
	}
	if _, _, err := runCommand(binary, store.Environment, []string{"revise", "--assertion-id", "as-h", "--normalized-text", input.Claim.Text}); err != nil {
		t.Fatal(err)
	}
	if _, _, err := runCommand(binary, store.Environment, []string{"attach-evidence", "--assertion-id", "as-h", "--evidence-kind", "correction", "--evidence-text", input.Episode.UserContributions[0].SourceText, "--evidence-observed-at", input.Episode.UserContributions[0].ObservedAt}); err != nil {
		t.Fatal(err)
	}
	after := qualityStateOf(t, store.Path)
	for _, id := range c.Expected.StoreDiff.Retain {
		if id == "rev-h-1" {
			if !qualityRevisionExists(t, store.Path, "as-h", 1) {
				t.Fatal("保持revisionが消失")
			}

			continue
		}
		if !after.IDs[id] {
			t.Fatalf("保持対象が消失: %s", id)
		}
	}
	if !reflect.DeepEqual(c.Expected.StoreDiff.Add, []string{"rev-h-2", "ev-h-correction"}) {
		t.Fatalf("H add oracle = %v", c.Expected.StoreDiff.Add)
	}
	if after.Revisions["as-h"] != before.Revisions["as-h"]+1 {
		t.Fatal("revisionが追加されません")
	}
	if after.Evidence["as-h\x00correction\x00"+input.Episode.UserContributions[0].SourceText] != 1 {
		t.Fatal("訂正Evidenceが追加されません")
	}
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

func assertQualityResultIDs(t *testing.T, stdout string, want []string) {
	t.Helper()
	var response struct {
		Data struct {
			Results []struct {
				ID string `json:"assertion_id"`
			} `json:"results"`
			AssertionID string `json:"assertion_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(stdout), &response); err != nil {
		t.Fatal(err)
	}
	got := make([]string, 0, len(response.Data.Results))
	for _, r := range response.Data.Results {
		got = append(got, r.ID)
	}
	if response.Data.AssertionID != "" {
		got = []string{response.Data.AssertionID}
	}
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("result IDs = %v, want %v, stdout=%s", got, want, stdout)
	}
}
func assertQualityErrorCode(t *testing.T, stderr, want string) {
	t.Helper()
	var response struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(stderr), &response); err != nil {
		t.Fatal(err)
	}
	if response.Error.Code != want {
		t.Fatalf("error code = %q, want %q", response.Error.Code, want)
	}
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
