package integration_test

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// runtimeCaseResultはRuntimeが呼出しsessionへ返す一時的なCase Resultである。
type runtimeCaseResult struct {
	CaseID             string          `json:"case_id"`
	ExecutionMode      string          `json:"execution_mode"`
	Status             string          `json:"status"`
	ExecutedLayers     []string        `json:"executed_layers"`
	FirstMismatchLayer *string         `json:"first_mismatch_layer"`
	Assessment         json.RawMessage `json:"assessment"`
	Confidence         string          `json:"confidence"`
	TraceStop          string          `json:"trace_stop"`
	CandidateIDs       []string        `json:"candidate_ids"`
	UpdateStatus       string          `json:"update_status"`
	Operations         json.RawMessage `json:"cli_operations"`
	PartialTraceStop   *string         `json:"partial_trace_stop"`
	Assessments        json.RawMessage `json:"assessments"`
	Markdown           string          `json:"markdown"`
}

type runtimeAssessment struct {
	Status     string `json:"status"`
	Confidence string `json:"confidence"`
}

type runtimeOperation struct {
	Operation string `json:"operation"`
}

type runtimeClaimAssessment struct {
	AssertionID string `json:"assertion_id"`
	Status      string `json:"status"`
}

func TestKnowledgeQualityRuntimeAcceptance(t *testing.T) {
	if os.Getenv("KNOWLEDGE_RUNTIME_ACCEPTANCE") != "1" {
		t.Skip("Codex Runtimeを使う明示受入実行ではありません")
	}
	fixture := readQualityFixture(t)
	validateQualityFixture(t, fixture)
	binary := buildCLI(t, true)
	selected := os.Getenv("KNOWLEDGE_ACCEPTANCE_CASE_ID")
	found := selected == ""
	for _, testCase := range fixture.Cases {
		t.Run(testCase.CaseID, func(t *testing.T) {
			if selected != "" && selected != testCase.CaseID {
				t.Skip("選択caseではありません")
			}
			found = true
			root := t.TempDir()
			store := defaultStoreConfiguration(t, root)
			seedQualityStore(t, store.Path, fixture.Seeds[testCase.SeedID])
			result := runRuntimeAcceptance(t, binary, store, testCase)
			assertRuntimeCaseResult(t, testCase, result)
		})
	}
	if !found {
		t.Fatalf("選択caseがfixtureにありません: %s", selected)
	}
}

func runRuntimeAcceptance(t *testing.T, binary string, store defaultStore, testCase qualityCase) runtimeCaseResult {
	t.Helper()
	output := filepath.Join(t.TempDir(), "runtime-case-result.json")
	runtimeBinary := runtimeStoreBinary(t, binary, store, testCase)
	prompt := "これはFEAT-005のテスト専用Runtime受入評価です。repositoryを変更せず、与えた隔離Storeだけを使ってください。最初に次の既存Skillを読み、その通常契約に従って評価してください: " + strings.Join(runtimeSkillPaths(testCase), ", ") + "。Reading Value、外部URL、共有Storeは使わないでください。指定したknowledge binaryを実際に呼び、必要なCLI操作とWorkflow判断を観測してください。\n" +
		"最終回答はJSONだけにし、case_id、execution_mode=runtime_acceptance、status=pass/failed/not_run、executed_layers、first_mismatch_layer、assessment、confidence、trace_stop、assessments、candidate_ids、update_status、cli_operations、partial_trace_stop、markdownを必ず含めてください。markdownには今回の一時評価Markdownを文字列で入れてください。\n" +
		"case_id: " + testCase.CaseID + "\ninput: " + string(testCase.Input) + "\nexpected CLI arguments: " + strings.Join(testCase.Expected.CLI.Arguments, " ") + "\nexpected operations: " + strings.Join(testCase.Expected.Operations, ",") + "\nknowledge binary: " + runtimeBinary + "\n"
	// #nosec G204 -- 全引数はfixture、テストが生成した隔離パス、または固定値である。
	cmd := exec.CommandContext(context.Background(), "codex", "exec", "--ephemeral", "-s", "workspace-write", "--add-dir", filepath.Dir(store.Path), "-C", filepath.Clean("../.."), "-o", output, prompt)
	t.Logf("Codex Runtime起動: case=%s skills=%s binary=%s", testCase.CaseID, strings.Join(runtimeSkillPaths(testCase), ","), runtimeBinary)
	combined, err := cmd.CombinedOutput()
	if len(combined) > 0 {
		t.Logf("Codex Runtime実行ログ:\n%s", combined)
	}
	if err != nil {
		t.Fatalf("Runtime launcher: %v: %s", err, combined)
	}
	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	var result runtimeCaseResult
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("Runtime Case ResultがJSONではありません: %v: %s", err, content)
	}
	t.Logf("Runtime Case Result: %s", content)

	return result
}

func runtimeSkillPaths(testCase qualityCase) []string {
	paths := []string{}
	for _, layer := range testCase.Layers {
		switch layer {
		case "knowledge_search":
			paths = append(paths, "skills/knowledge-search/SKILL.md")
		case "knowledge_acquisition":
			paths = append(paths, "skills/knowledge-acquisition/SKILL.md")
		case "knowledge_update":
			paths = append(paths, "skills/knowledge-update/SKILL.md")
		}
	}

	return paths
}

// runtimeStoreBinaryはCodex本体の認証用HOMEを変えず、子CLIだけを隔離Storeへ向ける。
func runtimeStoreBinary(t *testing.T, binary string, store defaultStore, testCase qualityCase) string {
	t.Helper()
	var assignment string
	for _, entry := range store.Environment {
		if key, value, found := strings.Cut(entry, "="); found && (key == "HOME" || key == "APPDATA" || key == "XDG_CONFIG_HOME") {
			assignment = key + "='" + strings.ReplaceAll(value, "'", "'\\''") + "'\n"

			break
		}
	}
	if assignment == "" {
		t.Fatal("隔離Store設定がありません")
	}
	// CodexにはStoreの親ディレクトリだけを追加許可しているため、wrapperもそこへ置く。
	parent := filepath.Dir(store.Path)
	path := filepath.Join(parent, "knowledge-runtime-store")
	script := "#!/bin/sh\n" + assignment + "export " + strings.Split(assignment, "=")[0] + "\n"
	if testCase.CaseID == "FEAT005-X-SEARCH-TECHNICAL-FAILURE" {
		blockedRoot := filepath.Join(parent, "blocked-store-root")
		if err := os.WriteFile(blockedRoot, []byte("blocked"), 0o600); err != nil {
			t.Fatal(err)
		}
		key := strings.Split(assignment, "=")[0]
		script = "#!/bin/sh\n" + key + "='" + strings.ReplaceAll(blockedRoot, "'", "'\\''") + "'\nexport " + key + "\nexec '" + strings.ReplaceAll(binary, "'", "'\\''") + "' \"$@\"\n"
	}
	if testCase.CaseID == "FEAT005-X-SEARCH-CANCELED" {
		ready := filepath.Join(parent, "search-text-ready")
		script += "export KNOWLEDGE_TEST_INTEGRATION_GATE_STAGE='search-text'\n" +
			"export KNOWLEDGE_TEST_INTEGRATION_GATE_READY='" + strings.ReplaceAll(ready, "'", "'\\''") + "'\n" +
			"'" + strings.ReplaceAll(binary, "'", "'\\''") + "' \"$@\" &\npid=$!\n" +
			"for i in $(seq 1 500); do\n  if [ -s \"$KNOWLEDGE_TEST_INTEGRATION_GATE_READY\" ]; then kill -INT \"$pid\"; break; fi\n  sleep 0.01\ndone\nwait \"$pid\"\nexit $?\n"
	} else if testCase.CaseID != "FEAT005-X-SEARCH-TECHNICAL-FAILURE" {
		script += "exec '" + strings.ReplaceAll(binary, "'", "'\\''") + "' \"$@\"\n"
	}
	// #nosec G306 -- 実行可能なテスト用wrapperであり、親TempDirは0700で隔離される。
	if err := os.WriteFile(path, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}

	return path
}

func assertRuntimeCaseResult(t *testing.T, testCase qualityCase, got runtimeCaseResult) {
	t.Helper()
	if got.CaseID != testCase.CaseID || got.ExecutionMode != "runtime_acceptance" || got.Status != "pass" || (got.FirstMismatchLayer != nil && *got.FirstMismatchLayer != "none") || len(got.ExecutedLayers) == 0 || got.Markdown == "" {
		t.Fatalf("Runtime Case Resultが不正: %#v", got)
	}
	assessment := runtimeAssessmentOf(t, got.Assessment)
	if assessment.Confidence == "" {
		assessment.Confidence = got.Confidence
	}
	if testCase.Expected.Assessment != "" && (assessment.Status != testCase.Expected.Assessment || assessment.Confidence != testCase.Expected.Confidence || normalizeRuntimeValue(got.TraceStop) != testCase.Expected.TraceStop) {
		t.Fatalf("Search Runtime oracle = %#v", got)
	}
	if len(testCase.Expected.Assessments) > 0 && !reflect.DeepEqual(runtimeAssessments(t, got.Assessments), testCase.Expected.Assessments) {
		t.Fatalf("複数Claim Runtime oracle = %#v", got)
	}
	if (len(testCase.Expected.CandidateIDs) > 0 && !reflect.DeepEqual(got.CandidateIDs, testCase.Expected.CandidateIDs)) || (testCase.Expected.UpdateStatus != "" && got.UpdateStatus != testCase.Expected.UpdateStatus) || (len(testCase.Expected.Operations) > 0 && !reflect.DeepEqual(runtimeOperationNames(t, got.Operations), testCase.Expected.Operations)) {
		t.Fatalf("Acquisition/Update Runtime oracle = %#v", got)
	}
	if testCase.Expected.PartialTrace.Stop != "" && (got.PartialTraceStop == nil || normalizeRuntimeValue(*got.PartialTraceStop) != testCase.Expected.PartialTrace.Stop) {
		t.Fatalf("partial Trace Runtime oracle = %#v", got)
	}
}

func runtimeAssessmentOf(t *testing.T, raw json.RawMessage) runtimeAssessment {
	t.Helper()
	var value runtimeAssessment
	if err := json.Unmarshal(raw, &value); err == nil && value.Status != "" {
		return value
	}
	if err := json.Unmarshal(raw, &value.Status); err != nil {
		t.Fatalf("Runtime assessmentを読む: %v", err)
	}

	return value
}

func runtimeAssessments(t *testing.T, raw json.RawMessage) map[string]string {
	t.Helper()
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	var values []runtimeClaimAssessment
	if err := json.Unmarshal(raw, &values); err != nil {
		t.Fatalf("Runtime assessmentsを読む: %v", err)
	}
	result := make(map[string]string, len(values))
	for _, value := range values {
		result[value.AssertionID] = value.Status
	}

	return result
}

func runtimeOperationNames(t *testing.T, raw json.RawMessage) []string {
	t.Helper()
	var values []runtimeOperation
	if err := json.Unmarshal(raw, &values); err != nil {
		t.Fatalf("Runtime operationsを読む: %v", err)
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, value.Operation)
	}

	return result
}

func normalizeRuntimeValue(value string) string {
	normalized := strings.ReplaceAll(value, " ", "_")
	if normalized == "no_viable_path" {
		return "no_expandable_path"
	}

	return normalized
}
