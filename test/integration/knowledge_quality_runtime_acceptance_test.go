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
	CaseID             string               `json:"case_id"`
	ExecutionMode      string               `json:"execution_mode"`
	Status             string               `json:"status"`
	ExecutedLayers     []string             `json:"executed_layers"`
	FirstMismatchLayer *string              `json:"first_mismatch_layer"`
	NotExecutedLayers  []string             `json:"not_executed_layers"`
	Assessment         json.RawMessage      `json:"assessment"`
	Confidence         string               `json:"confidence"`
	TraceStop          string               `json:"trace_stop"`
	CandidateIDs       []string             `json:"candidate_ids"`
	UpdateStatus       string               `json:"update_status"`
	Operations         json.RawMessage      `json:"cli_operations"`
	PartialTrace       *runtimePartialTrace `json:"partial_trace"`
	SearchTrace        *runtimeSearchTrace  `json:"search_trace"`
	Assessments        json.RawMessage      `json:"assessments"`
	Markdown           string               `json:"markdown"`
}

type runtimeAssessment struct {
	Status     string `json:"status"`
	Confidence string `json:"confidence"`
}

type runtimeOperation struct {
	Operation string `json:"operation"`
}

type runtimePartialTrace struct {
	Operation string `json:"operation"`
	Stop      string `json:"stop"`
	ExitCode  int    `json:"exit_code"`
	ErrorCode string `json:"error_code"`
}

type runtimeSearchTrace struct {
	Operations  []string `json:"operations"`
	Queries     []string `json:"queries"`
	ResultIDs   []string `json:"result_ids"`
	EvidenceIDs []string `json:"evidence_ids"`
	BudgetUsed  int      `json:"budget_used"`
	Stop        string   `json:"stop"`
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
	prompt := runtimeAcceptancePrompt(t, runtimeBinary, testCase)
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
	result := decodeRuntimeCaseResult(t, content)
	t.Logf("Runtime Case Result: %s", content)

	return result
}

func decodeRuntimeCaseResult(t *testing.T, content []byte) runtimeCaseResult {
	t.Helper()
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(content, &fields); err != nil {
		t.Fatalf("Runtime Case ResultがJSONではありません: %v: %s", err, content)
	}
	for _, key := range []string{"case_id", "execution_mode", "status", "executed_layers", "first_mismatch_layer", "not_executed_layers", "assessment", "confidence", "trace_stop", "search_trace", "assessments", "candidate_ids", "update_status", "cli_operations", "partial_trace", "markdown"} {
		if _, ok := fields[key]; !ok {
			t.Fatalf("Runtime Case Resultに必須fieldがありません: %s: %s", key, content)
		}
	}
	var result runtimeCaseResult
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("Runtime Case Resultがschemaに一致しません: %v: %s", err, content)
	}

	return result
}

func runtimeAcceptancePrompt(t *testing.T, runtimeBinary string, testCase qualityCase) string {
	t.Helper()
	input := runtimeInput(t, testCase.Input)
	expected, err := json.Marshal(testCase.Expected)
	if err != nil {
		t.Fatal(err)
	}
	layers, err := json.Marshal(testCase.Layers)
	if err != nil {
		t.Fatal(err)
	}

	return "これはFEAT-005のテスト専用Runtime受入評価です。repositoryを変更せず、与えた隔離Storeだけを使ってください。最初に次の既存Skillを読み、その通常契約に従って評価してください: " + strings.Join(runtimeSkillPaths(testCase), ", ") + "。Reading Value、外部URL、共有Storeは使わないでください。指定したknowledge binaryを実際に呼び、必要なCLI操作とWorkflow判断を観測してください。\n" +
		"最終回答は次の固定schemaに一致するJSONオブジェクトだけにしてください。Markdownの前後に説明やcode fenceを置いてはいけません。statusは、期待された成功又は期待された技術失敗／中断を実観測して全oracleが一致した場合だけpassです。期待と異なる場合はfailed、実行不能ならnot_runです。期待されたXのCLI失敗そのものをstatus=failedにしてはいけません。\n" +
		"{\"case_id\":string,\"execution_mode\":\"runtime_acceptance\",\"status\":\"pass\"|\"failed\"|\"not_run\",\"executed_layers\":[string],\"first_mismatch_layer\":\"none\"|string,\"not_executed_layers\":[string],\"assessment\":null|{\"status\":string,\"confidence\":string},\"confidence\":string,\"trace_stop\":string,\"search_trace\":null|{\"operations\":[string],\"queries\":[string],\"result_ids\":[string],\"evidence_ids\":[string],\"budget_used\":number,\"stop\":string},\"assessments\":null|[{\"assertion_id\":string,\"status\":string}],\"candidate_ids\":[string],\"update_status\":string,\"cli_operations\":[{\"operation\":string}],\"partial_trace\":null|{\"operation\":string,\"stop\":string,\"exit_code\":number,\"error_code\":string},\"markdown\":string}\n" +
		"executed_layersは次の固定配列を同じ順序で報告してください: " + string(layers) + "。not_executed_layersはexpected oracleの配列を同じ順序で報告してください。両配列は排他的です。検索結果が空ならそのIDを推測してget/get-evidenceや更新を実行せず、最初の不一致層をfailedとして報告してください。assessmentがnullならconfidenceとtrace_stopは空文字列にしてください。assessmentがある場合はassessment.confidenceとconfidenceを同じ値にしてください。assessmentsはclaim_idやselection_statusを使わず、必ずassertion_idとstatusだけを持つ配列にしてください。\n" +
		"cli_operationsはKnowledge Update対象のcaseだけに記録します。Searchだけのcaseは実際にsearch-textを使っても[]にしてください。Candidate対象外のcaseのcandidate_idsは[]にしてください。partial_traceはX以外ではnullです。Xでは期待されたCLI exit_code（expected oracleのcli.exit_code）を使い、operation、stop、error_codeを含むobjectにしてください。markdownは空文字列にせず、観測結果を要約するMarkdown見出しと本文を必ず入れてください。\n" +
		"Hではattach-evidenceの公開結果が自動採番のevidence_idでも、対象Assertion・kind・text・observed_atがexpected oracleと一致すれば、Fixtureの論理ID ev-h-correctionに対応した成功として報告してください。訂正Candidateの引用された旧命題を二つ目のsearch-text queryとして使えるのは、通常Skillのsearch_queries契約に従って検索結果からIDを得た場合だけです。Hのcandidate_idsとcli_operationsは期待oracleの固定値・固定順序をそのまま報告し、余分なCLI操作を含めてはいけません。\n" +
		"case_id: " + testCase.CaseID + "\ninput: " + string(input) + "\nexpected oracle: " + string(expected) + "\nknowledge binary: " + runtimeBinary + "\n"
}

// runtimeInputはDEC-FEAT-018どおりReading Value参照をRuntime入力から除外する。
func runtimeInput(t *testing.T, raw json.RawMessage) json.RawMessage {
	t.Helper()
	var input map[string]json.RawMessage
	if err := json.Unmarshal(raw, &input); err != nil {
		t.Fatalf("Runtime inputを読む: %v", err)
	}
	delete(input, "reading_value_reference")
	result, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

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
		ready := filepath.Join(t.TempDir(), "search-text-ready")
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
	if got.CaseID != testCase.CaseID || got.ExecutionMode != "runtime_acceptance" || got.Status != "pass" || got.FirstMismatchLayer == nil || *got.FirstMismatchLayer != "none" || !reflect.DeepEqual(got.ExecutedLayers, testCase.Layers) || !reflect.DeepEqual(got.NotExecutedLayers, testCase.Expected.NotExecuted) || got.Markdown == "" {
		t.Fatalf("Runtime Case Resultが不正: %#v", got)
	}
	assessment := runtimeAssessmentOf(t, got.Assessment)
	if testCase.Expected.Assessment == "" {
		if assessment != (runtimeAssessment{}) || got.Confidence != "" || got.TraceStop != "" {
			t.Fatalf("Search対象外caseのRuntime oracle = %#v", got)
		}
	} else if assessment.Status != testCase.Expected.Assessment || assessment.Confidence != testCase.Expected.Confidence || got.Confidence != assessment.Confidence || got.TraceStop != testCase.Expected.TraceStop {
		t.Fatalf("Search Runtime oracle = %#v", got)
	}
	if len(testCase.Expected.Assessments) > 0 && !reflect.DeepEqual(runtimeAssessments(t, got.Assessments), testCase.Expected.Assessments) {
		t.Fatalf("複数Claim Runtime oracle = %#v", got)
	}
	if len(testCase.Expected.Assessments) == 0 && !runtimeAssessmentsAreEmpty(t, got.Assessments) {
		t.Fatalf("複数Claim対象外caseのRuntime oracle = %#v", got)
	}
	if len(got.CandidateIDs) != len(testCase.Expected.CandidateIDs) || (len(testCase.Expected.CandidateIDs) > 0 && !reflect.DeepEqual(got.CandidateIDs, testCase.Expected.CandidateIDs)) || (testCase.Expected.UpdateStatus != "" && got.UpdateStatus != testCase.Expected.UpdateStatus) || (runtimeCaseHasUpdate(testCase) && !runtimeOperationsMatch(t, got.Operations, testCase.Expected.Operations)) {
		t.Fatalf("Acquisition/Update Runtime oracle = %#v", got)
	}
	if testCase.Expected.PartialTrace.Stop != "" && (got.PartialTrace == nil || got.PartialTrace.Operation != testCase.Expected.PartialTrace.Operation || got.PartialTrace.Stop != testCase.Expected.PartialTrace.Stop || got.PartialTrace.ExitCode != runtimePartialTraceExitCode(testCase) || got.PartialTrace.ErrorCode != testCase.Expected.CLI.ErrorCode) {
		t.Fatalf("partial Trace Runtime oracle = %#v", got)
	}
	if testCase.Expected.PartialTrace.Stop == "" && got.PartialTrace != nil {
		t.Fatalf("partial Trace対象外caseのRuntime oracle = %#v", got)
	}
	if !runtimeSearchTraceMatches(testCase, got.SearchTrace) {
		t.Fatalf("Search Trace対象外caseのRuntime oracle = %#v", got.SearchTrace)
	}
}

func runtimeCaseHasUpdate(testCase qualityCase) bool {
	for _, layer := range testCase.Layers {
		if layer == "knowledge_update" {
			return true
		}
	}

	return false
}

func runtimePartialTraceExitCode(testCase qualityCase) int {
	if testCase.Expected.CLI.ExitCode != 0 {
		return testCase.Expected.CLI.ExitCode
	}

	return testCase.Expected.PartialTrace.ExitCode
}

func runtimeAssessmentOf(t *testing.T, raw json.RawMessage) runtimeAssessment {
	t.Helper()
	if runtimeValueIsNull(raw) {
		return runtimeAssessment{}
	}
	var value runtimeAssessment
	if err := json.Unmarshal(raw, &value); err != nil || value.Status == "" || value.Confidence == "" {
		t.Fatalf("Runtime assessmentを読む: %v", err)
	}

	return value
}

func runtimeAssessments(t *testing.T, raw json.RawMessage) map[string]string {
	t.Helper()
	if runtimeValueIsNull(raw) {
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

func runtimeValueIsNull(raw json.RawMessage) bool {
	return len(raw) == 0 || string(raw) == "null"
}

func runtimeAssessmentsAreEmpty(t *testing.T, raw json.RawMessage) bool {
	t.Helper()

	return runtimeValueIsNull(raw) || len(runtimeAssessments(t, raw)) == 0
}

func runtimeOperationsMatch(t *testing.T, raw json.RawMessage, expected []string) bool {
	t.Helper()
	actual := runtimeOperationNames(t, raw)

	return len(actual) == len(expected) && (len(expected) == 0 || reflect.DeepEqual(actual, expected))
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

func TestRuntimeCaseResultContract(t *testing.T) {
	fixture := readQualityFixture(t)
	validateQualityFixture(t, fixture)
	for _, testCase := range fixture.Cases {
		t.Run(testCase.CaseID, func(t *testing.T) {
			result := passingRuntimeCaseResult(t, testCase)
			content, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}
			assertRuntimeCaseResult(t, testCase, decodeRuntimeCaseResult(t, content))
		})
	}
	for _, testCase := range fixture.Cases {
		if !runtimeSearchTraceForbidden(testCase) {
			continue
		}
		t.Run(testCase.CaseID+"/search-trace-forbidden", func(t *testing.T) {
			result := passingRuntimeCaseResult(t, testCase)
			result.SearchTrace = &runtimeSearchTrace{Stop: "unexpected"}
			content, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}
			got := decodeRuntimeCaseResult(t, content)
			if got.SearchTrace == nil {
				t.Fatal("否定例のSearch Traceを作成できません")
			}
			if runtimeSearchTraceMatches(testCase, got.SearchTrace) {
				t.Fatalf("Search Trace対象外caseを拒否できません: %s", testCase.CaseID)
			}
		})
	}
}

func runtimeSearchTraceForbidden(testCase qualityCase) bool {
	return !hasQualitySearchLayer(testCase) || testCase.Expected.PartialTrace.Stop != ""
}

func runtimeSearchTraceMatches(testCase qualityCase, trace *runtimeSearchTrace) bool {
	if runtimeSearchTraceForbidden(testCase) {
		return trace == nil
	}

	return trace != nil && reflect.DeepEqual(*trace, runtimeSearchTrace{
		Operations:  testCase.Expected.Trace.Operations,
		Queries:     testCase.Expected.Trace.Queries,
		ResultIDs:   testCase.Expected.Trace.ResultIDs,
		EvidenceIDs: testCase.Expected.Trace.EvidenceIDs,
		BudgetUsed:  testCase.Expected.Trace.BudgetUsed,
		Stop:        testCase.Expected.Trace.Stop,
	})
}

func TestRuntimeAcceptancePromptExcludesReadingValue(t *testing.T) {
	fixture := readQualityFixture(t)
	for _, testCase := range fixture.Cases {
		if testCase.CaseID != "FEAT005-I-READ-SELECTED" {
			continue
		}
		prompt := runtimeAcceptancePrompt(t, "/tmp/knowledge", testCase)
		if strings.Contains(prompt, "FEAT-003 V-002") {
			t.Fatal("Reading Value参照をRuntime入力へ渡しています")
		}
		if !strings.Contains(prompt, "assertion_id") || !strings.Contains(prompt, "partial_trace") {
			t.Fatal("Runtime Case Result schemaがpromptにありません")
		}
		if !strings.Contains(prompt, "[\"cli_store\",\"knowledge_search\",\"end_to_end\"]") {
			t.Fatal("Runtimeへ期待executed_layersを渡していません")
		}
		if !strings.Contains(prompt, "cli.exit_code") || !strings.Contains(prompt, "ev-h-correction") {
			t.Fatal("Runtime出力のX/H契約がpromptにありません")
		}

		return
	}
	t.Fatal("I caseがfixtureにありません")
}

func passingRuntimeCaseResult(t *testing.T, testCase qualityCase) runtimeCaseResult {
	t.Helper()
	operations := make([]runtimeOperation, 0, len(testCase.Expected.Operations))
	for _, operation := range testCase.Expected.Operations {
		operations = append(operations, runtimeOperation{
			Operation: operation,
		})
	}
	operationsJSON, err := json.Marshal(operations)
	if err != nil {
		t.Fatal(err)
	}
	assessments := make([]runtimeClaimAssessment, 0, len(testCase.Expected.Assessments))
	for assertionID, status := range testCase.Expected.Assessments {
		assessments = append(assessments, runtimeClaimAssessment{
			AssertionID: assertionID,
			Status:      status,
		})
	}
	assessmentsJSON, err := json.Marshal(assessments)
	if err != nil {
		t.Fatal(err)
	}
	assessment := json.RawMessage("null")
	if testCase.Expected.Assessment != "" {
		assessment, err = json.Marshal(runtimeAssessment{
			Status:     testCase.Expected.Assessment,
			Confidence: testCase.Expected.Confidence,
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	firstMismatch := "none"
	result := runtimeCaseResult{
		CaseID:             testCase.CaseID,
		ExecutionMode:      "runtime_acceptance",
		Status:             "pass",
		ExecutedLayers:     testCase.Layers,
		FirstMismatchLayer: &firstMismatch,
		NotExecutedLayers:  testCase.Expected.NotExecuted,
		Assessment:         assessment,
		Confidence:         testCase.Expected.Confidence,
		TraceStop:          testCase.Expected.TraceStop,
		CandidateIDs:       testCase.Expected.CandidateIDs,
		UpdateStatus:       testCase.Expected.UpdateStatus,
		Operations:         operationsJSON,
		Assessments:        assessmentsJSON,
		Markdown:           "# Runtime acceptance\n",
	}
	if testCase.Expected.PartialTrace.Stop != "" {
		result.PartialTrace = &runtimePartialTrace{
			Operation: testCase.Expected.PartialTrace.Operation,
			Stop:      testCase.Expected.PartialTrace.Stop,
			ExitCode:  runtimePartialTraceExitCode(testCase),
			ErrorCode: testCase.Expected.CLI.ErrorCode,
		}
	}
	if hasQualitySearchLayer(testCase) && testCase.Expected.PartialTrace.Stop == "" {
		result.SearchTrace = &runtimeSearchTrace{
			Operations:  testCase.Expected.Trace.Operations,
			Queries:     testCase.Expected.Trace.Queries,
			ResultIDs:   testCase.Expected.Trace.ResultIDs,
			EvidenceIDs: testCase.Expected.Trace.EvidenceIDs,
			BudgetUsed:  testCase.Expected.Trace.BudgetUsed,
			Stop:        testCase.Expected.Trace.Stop,
		}
	}

	return result
}
