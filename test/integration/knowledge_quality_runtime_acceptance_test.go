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
	Candidates         json.RawMessage      `json:"candidates"`
	UpdateStatus       string               `json:"update_status"`
	Decisions          json.RawMessage      `json:"decisions"`
	Operations         json.RawMessage      `json:"cli_operations"`
	StoreDiff          json.RawMessage      `json:"store_diff"`
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

// runtimeCandidateとruntimeDecisionは、Runtimeが通常Workflowから取得した
// 一時成果物をテスト側で独立に検証するための最小表現である。
type runtimeCandidate struct {
	ID                string            `json:"candidate_id"`
	EvidenceKind      string            `json:"evidence_kind"`
	Strength          string            `json:"strength"`
	ProposedAssertion string            `json:"proposed_assertion"`
	EvidenceRawText   string            `json:"evidence_raw_text"`
	Scope             map[string]string `json:"scope"`
	Temporal          json.RawMessage   `json:"temporal"`
	SearchQueries     []string          `json:"search_queries"`
}

type runtimeDecision struct {
	CandidateID       string             `json:"candidate_id"`
	Action            string             `json:"action"`
	TargetAssertionID string             `json:"target_assertion_id"`
	ExecutionStatus   string             `json:"execution_status"`
	Operations        []runtimeOperation `json:"cli_operations"`
}

type runtimeStoreDiff struct {
	None   bool     `json:"none"`
	Retain []string `json:"retain"`
	Add    []string `json:"add"`
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
			before := qualityStateOf(t, store.Path)
			result, invocations := runRuntimeAcceptance(t, binary, store, testCase)
			assertRuntimeCaseResult(t, testCase, result)
			assertRuntimeCLIInvocations(t, invocations, testCase, result)
			assertRuntimeStoreDiff(t, store.Path, before, qualityStateOf(t, store.Path), testCase)
		})
	}
	if !found {
		t.Fatalf("選択caseがfixtureにありません: %s", selected)
	}
}

func runRuntimeAcceptance(t *testing.T, binary string, store defaultStore, testCase qualityCase) (runtimeCaseResult, string) {
	t.Helper()
	output := filepath.Join(t.TempDir(), "runtime-case-result.json")
	runtimeBinary := runtimeStoreBinary(t, binary, store, testCase)
	prompt := runtimeAcceptancePrompt(t, runtimeBinary.Path, testCase)
	// #nosec G204 -- 全引数はfixture、テストが生成した隔離パス、または固定値である。
	cmd := exec.CommandContext(context.Background(), "codex", "exec", "--ephemeral", "-s", "workspace-write", "--add-dir", filepath.Dir(store.Path), "-C", filepath.Clean("../.."), "-o", output, prompt)
	t.Logf("Codex Runtime起動: case=%s skills=%s binary=%s", testCase.CaseID, strings.Join(runtimeSkillPaths(testCase), ","), runtimeBinary.Path)
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

	return result, runtimeBinary.InvocationLog
}

func decodeRuntimeCaseResult(t *testing.T, content []byte) runtimeCaseResult {
	t.Helper()
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(content, &fields); err != nil {
		t.Fatalf("Runtime Case ResultがJSONではありません: %v: %s", err, content)
	}
	for _, key := range []string{"case_id", "execution_mode", "status", "executed_layers", "first_mismatch_layer", "not_executed_layers", "assessment", "confidence", "trace_stop", "search_trace", "assessments", "candidate_ids", "candidates", "update_status", "decisions", "cli_operations", "store_diff", "partial_trace", "markdown"} {
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
	layers, err := json.Marshal(testCase.Layers)
	if err != nil {
		t.Fatal(err)
	}

	return "これはFEAT-005のテスト専用Runtime受入評価です。repositoryを変更せず、与えた隔離Storeだけを使ってください。最初に次の既存Skillを読み、その通常契約に従って評価してください: " + strings.Join(runtimeSkillPaths(testCase), ", ") + "。Reading Value、外部URL、共有Storeは使わないでください。指定したknowledge binaryを実際に呼び、必要なCLI操作とWorkflow判断を観測してください。\n" +
		"最終回答は次の固定schemaに一致するJSONオブジェクトだけにしてください。Markdownの前後に説明やcode fenceを置いてはいけません。期待oracleは渡しません。statusは通常Skillを実行した観測結果に基づき、正常完了ならpass、実行不能ならnot_run、想定外の不一致ならfailedにしてください。期待されたXのCLI失敗そのものをstatus=failedにしてはいけません。\n" +
		"{\"case_id\":string,\"execution_mode\":\"runtime_acceptance\",\"status\":\"pass\"|\"failed\"|\"not_run\",\"executed_layers\":[string],\"first_mismatch_layer\":\"none\"|string,\"not_executed_layers\":[string],\"assessment\":null|{\"status\":string,\"confidence\":string},\"confidence\":string,\"trace_stop\":string,\"search_trace\":null|{\"operations\":[string],\"queries\":[string],\"result_ids\":[string],\"evidence_ids\":[string],\"budget_used\":number,\"stop\":string},\"assessments\":null|[{\"assertion_id\":string,\"status\":string}],\"candidate_ids\":[string],\"candidates\":[{\"candidate_id\":string,\"evidence_kind\":string,\"strength\":string,\"proposed_assertion\":string,\"evidence_raw_text\":string,\"scope\":{string:string},\"temporal\":null|object,\"search_queries\":[string]}],\"update_status\":string,\"decisions\":[{\"candidate_id\":string,\"action\":string,\"target_assertion_id\":string,\"execution_status\":string,\"cli_operations\":[{\"operation\":string}]}],\"cli_operations\":[{\"operation\":string}],\"store_diff\":{\"none\":bool,\"retain\":[string],\"add\":[string]},\"partial_trace\":null|{\"operation\":string,\"stop\":string,\"exit_code\":number,\"error_code\":string},\"markdown\":string}\n" +
		"executed_layersは次の固定配列を同じ順序で報告してください: " + string(layers) + "。not_executed_layersは実行しなかった層を順序どおり報告し、両配列を排他的にしてください。検索結果が空ならそのIDを推測してget/get-evidenceや更新を実行せず、最初の不一致層を報告してください。assessmentがnullならconfidenceとtrace_stopは空文字列にしてください。assessmentがある場合はassessment.confidenceとconfidenceを同じ値にしてください。assessmentsはclaim_idやselection_statusを使わず、必ずassertion_idとstatusだけを持つ配列にしてください。\n" +
		"AcquisitionとUpdateは別層として診断してください。AcquisitionがCandidateを出さない場合はcandidates、candidate_ids、decisions、cli_operationsをすべて空にし、Updateはcompletedを報告してください。Candidateを出す場合は、Acquisition Markdownから候補fieldsをそのままcandidatesへ、Update Markdownからdecisionとそのcli_operationsをそのままdecisionsへ転記し、そのdecision操作列をcli_operationsにも同順で報告してください。手動でCandidateやDecisionを作らず、UpdateへはAcquisitionのCandidateだけを渡してください。store_diffは隔離Storeの実前後差分を報告してください。Searchだけのcaseは実際にsearch-textを使ってもcli_operations、candidates、decisionsを[]にしてください。partial_traceはX以外ではnullです。Xでは最初のsearch-textの失敗／中断後に後続operationを実行せず、cli_operations、candidates、decisionsを[]にし、観測したexit_codeとoperation、stop、error_codeをpartial_traceへ入れてください。markdownは空文字列にせず、観測結果を要約するMarkdown見出しと本文を必ず入れてください。\n" +
		"case_id: " + testCase.CaseID + "\ninput: " + string(input) + "\nknowledge binary: " + runtimeBinary + "\n"
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

type runtimeBinary struct {
	Path          string
	InvocationLog string
}

// runtimeStoreBinaryはCodex本体の認証用HOMEを変えず、子CLIだけを隔離Storeへ向ける。
// 全起動argvを記録し、Runtimeの自己申告でなくプロセス境界で操作列を検証する。
func runtimeStoreBinary(t *testing.T, binary string, store defaultStore, testCase qualityCase) runtimeBinary {
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
	invocationLog := filepath.Join(parent, "knowledge-runtime-invocations")
	script := "#!/bin/sh\n" + assignment + "export " + strings.Split(assignment, "=")[0] + "\n"
	if testCase.CaseID == "FEAT005-X-SEARCH-TECHNICAL-FAILURE" {
		blockedRoot := filepath.Join(parent, "blocked-store-root")
		if err := os.WriteFile(blockedRoot, []byte("blocked"), 0o600); err != nil {
			t.Fatal(err)
		}
		key := strings.Split(assignment, "=")[0]
		script = "#!/bin/sh\n" + key + "='" + strings.ReplaceAll(blockedRoot, "'", "'\\''") + "'\nexport " + key + "\n"
	}
	// Record separator (0x1e) makes every invocation distinct; unit separator
	// (0x1f) preserves argv boundaries without depending on shell quoting.
	script += "printf '\\036' >> '" + strings.ReplaceAll(invocationLog, "'", "'\\''") + "'\nprintf '%s\\037' \"$@\" >> '" + strings.ReplaceAll(invocationLog, "'", "'\\''") + "'\n"
	if testCase.CaseID == "FEAT005-X-SEARCH-CANCELED" {
		ready := filepath.Join(t.TempDir(), "search-text-ready")
		script += "export KNOWLEDGE_TEST_INTEGRATION_GATE_STAGE='search-text'\n" +
			"export KNOWLEDGE_TEST_INTEGRATION_GATE_READY='" + strings.ReplaceAll(ready, "'", "'\\''") + "'\n" +
			"'" + strings.ReplaceAll(binary, "'", "'\\''") + "' \"$@\" &\npid=$!\n" +
			"for i in $(seq 1 500); do\n  if [ -s \"$KNOWLEDGE_TEST_INTEGRATION_GATE_READY\" ]; then kill -INT \"$pid\"; break; fi\n  sleep 0.01\ndone\nwait \"$pid\"\nexit $?\n"
	} else {
		script += "exec '" + strings.ReplaceAll(binary, "'", "'\\''") + "' \"$@\"\n"
	}
	// #nosec G306 -- 実行可能なテスト用wrapperであり、親TempDirは0700で隔離される。
	if err := os.WriteFile(path, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}

	return runtimeBinary{
		Path:          path,
		InvocationLog: invocationLog,
	}
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
	if !reflect.DeepEqual(got.CandidateIDs, runtimeCandidateIDs(testCase)) || got.UpdateStatus != testCase.Expected.UpdateStatus || !runtimeOperationsMatch(t, got.Operations, testCase.Expected.Operations) || !runtimeStoreDiffMatches(t, got.StoreDiff, runtimeExpectedStoreDiff(testCase)) {
		t.Fatalf("Acquisition/Update Runtime oracle = %#v", got)
	}
	assertRuntimeWorkflowArtifacts(t, testCase, got)
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

func assertRuntimeWorkflowArtifacts(t *testing.T, testCase qualityCase, got runtimeCaseResult) {
	t.Helper()
	candidates := runtimeCandidates(t, got.Candidates)
	decisions := runtimeDecisions(t, got.Decisions)
	if testCase.CaseID != "FEAT005-H-UPDATE-CORRECTION" {
		if len(candidates) != 0 || len(decisions) != 0 {
			t.Fatalf("Candidate対象外caseのWorkflow成果物 = candidates=%#v decisions=%#v", candidates, decisions)
		}

		return
	}
	if len(candidates) != 1 || len(decisions) != 1 {
		t.Fatalf("HのWorkflow成果物件数 = candidates=%#v decisions=%#v", candidates, decisions)
	}
	candidate := candidates[0]
	if candidate.ID != "cand-h-1" || candidate.EvidenceKind != "correction" || candidate.Strength != "strong" || candidate.ProposedAssertion != "unbuffered channelへのsendはreceiverが受信可能になったときに完了する" || candidate.EvidenceRawText != qualityCorrectionEvidenceText(t, testCase) || !reflect.DeepEqual(candidate.Scope, map[string]string{"language": "Go"}) || !runtimeValueIsNull(candidate.Temporal) || !reflect.DeepEqual(candidate.SearchQueries, []string{"unbuffered channelへのsendはreceiverが受信可能になったときに完了する", "unbuffered channelのsendはreceiverが受信可能になる前にも完了する"}) {
		t.Fatalf("H CandidateがWorkflow契約に一致しません: %#v", candidate)
	}
	decision := decisions[0]
	if decision.CandidateID != candidate.ID || decision.Action != "revise" || decision.TargetAssertionID != "as-h" || decision.ExecutionStatus != "applied" || !reflect.DeepEqual(runtimeOperationNames(t, got.Operations), runtimeOperationNamesOf(decision.Operations)) {
		t.Fatalf("H DecisionがWorkflow契約に一致しません: %#v", decision)
	}
}

func runtimeCandidates(t *testing.T, raw json.RawMessage) []runtimeCandidate {
	t.Helper()
	var values []runtimeCandidate
	if err := json.Unmarshal(raw, &values); err != nil {
		t.Fatalf("Runtime candidatesを読む: %v", err)
	}

	return values
}

func runtimeDecisions(t *testing.T, raw json.RawMessage) []runtimeDecision {
	t.Helper()
	var values []runtimeDecision
	if err := json.Unmarshal(raw, &values); err != nil {
		t.Fatalf("Runtime decisionsを読む: %v", err)
	}

	return values
}

func runtimeOperationNamesOf(values []runtimeOperation) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, value.Operation)
	}

	return result
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

func runtimeStoreDiffMatches(t *testing.T, raw json.RawMessage, expected runtimeStoreDiff) bool {
	t.Helper()
	var actual runtimeStoreDiff
	if err := json.Unmarshal(raw, &actual); err != nil {
		t.Fatalf("store_diffがschemaに一致しません: %v", err)
	}

	return reflect.DeepEqual(actual, expected)
}

func runtimeCandidateIDs(testCase qualityCase) []string {
	if testCase.Expected.CandidateIDs == nil {
		return []string{}
	}

	return testCase.Expected.CandidateIDs
}

func runtimeExpectedStoreDiff(testCase qualityCase) runtimeStoreDiff {
	retain := testCase.Expected.StoreDiff.Retain
	if retain == nil {
		retain = []string{}
	}
	add := testCase.Expected.StoreDiff.Add
	if add == nil {
		add = []string{}
	}

	return runtimeStoreDiff{
		None:   testCase.Expected.StoreDiff.None,
		Retain: retain,
		Add:    add,
	}
}

func assertRuntimeStoreDiff(t *testing.T, path string, before, after qualityState, testCase qualityCase) {
	t.Helper()
	if testCase.Expected.StoreDiff.None {
		if !reflect.DeepEqual(before, after) {
			t.Fatalf("Runtime caseがstore_diff:noneに反してStoreを変更しました: %s", testCase.CaseID)
		}

		return
	}
	if testCase.CaseID != "FEAT005-H-UPDATE-CORRECTION" {
		t.Fatalf("Runtime Store差分のoracleがありません: %s", testCase.CaseID)
	}
	for _, id := range testCase.Expected.StoreDiff.Retain {
		if id == "rev-h-1" {
			if !qualityRevisionExists(t, path, "as-h", 1) {
				t.Fatal("Runtime Hで保持revisionが消失しました")
			}

			continue
		}
		if !after.IDs[id] {
			t.Fatalf("Runtime Hで保持対象が消失しました: %s", id)
		}
	}
	if after.Revisions["as-h"] != before.Revisions["as-h"]+1 {
		t.Fatal("Runtime Hでrevisionが追加されません")
	}
	if after.Evidence["as-h\x00correction\x00"+qualityCorrectionEvidenceText(t, testCase)] != 1 {
		t.Fatal("Runtime Hで訂正Evidenceが追加されません")
	}
}

func qualityCorrectionEvidenceText(t *testing.T, testCase qualityCase) string {
	t.Helper()
	var input struct {
		Episode struct {
			UserContributions []struct {
				SourceText string `json:"source_text"`
			} `json:"user_contributions"`
		} `json:"episode"`
	}
	if err := json.Unmarshal(testCase.Input, &input); err != nil {
		t.Fatal(err)
	}
	if len(input.Episode.UserContributions) != 1 || input.Episode.UserContributions[0].SourceText == "" {
		t.Fatal("Hの訂正Evidence入力が不正です")
	}

	return input.Episode.UserContributions[0].SourceText
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

func assertRuntimeCLIInvocations(t *testing.T, path string, testCase qualityCase, got runtimeCaseResult) {
	t.Helper()
	invocations := runtimeCLIInvocations(t, path)
	if testCase.CaseID == "FEAT005-X-SEARCH-TECHNICAL-FAILURE" || testCase.CaseID == "FEAT005-X-SEARCH-CANCELED" {
		if !reflect.DeepEqual(invocations, [][]string{testCase.Expected.CLI.Arguments}) {
			t.Fatalf("Xが最初のsearch-text以外を起動しました: %#v", invocations)
		}

		return
	}
	if testCase.CaseID == "FEAT005-F-ACQUISITION-QUESTION" || testCase.CaseID == "FEAT005-G-ACQUISITION-AI-ONLY" {
		if len(invocations) != 0 {
			t.Fatalf("F/GはCLI operationを起動してはいけません: %#v", invocations)
		}

		return
	}
	if testCase.CaseID != "FEAT005-H-UPDATE-CORRECTION" {
		return
	}
	decisions := runtimeDecisions(t, got.Decisions)
	if len(decisions) != 1 || !reflect.DeepEqual(runtimeOperationNamesOf(decisions[0].Operations), runtimeInvocationOperationNames(invocations)) {
		t.Fatalf("HのWorkflow Decision操作と実CLI起動が一致しません: decision=%#v invocations=%#v", decisions, invocations)
	}
}

func runtimeCLIInvocations(t *testing.T, path string) [][]string {
	t.Helper()
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	var invocations [][]string
	for _, record := range strings.Split(string(content), "\x1e") {
		if record == "" {
			continue
		}
		args := strings.Split(record, "\x1f")
		if len(args) > 0 && args[len(args)-1] == "" {
			args = args[:len(args)-1]
		}
		invocations = append(invocations, args)
	}

	return invocations
}

func runtimeInvocationOperationNames(invocations [][]string) []string {
	operations := make([]string, 0, len(invocations))
	for _, args := range invocations {
		if len(args) == 0 {
			operations = append(operations, "")

			continue
		}
		operations = append(operations, args[0])
	}

	return operations
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
		if !strings.Contains(prompt, "期待oracleは渡しません") || !strings.Contains(prompt, "candidates") || !strings.Contains(prompt, "decisions") || !strings.Contains(prompt, "Xでは最初のsearch-text") {
			t.Fatal("Runtime出力の独立Workflow/X契約がpromptにありません")
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
	candidates, decisions := passingRuntimeWorkflowArtifacts(t, testCase, operations)
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
		CandidateIDs:       runtimeCandidateIDs(testCase),
		Candidates:         candidates,
		UpdateStatus:       testCase.Expected.UpdateStatus,
		Decisions:          decisions,
		Operations:         operationsJSON,
		StoreDiff:          runtimeStoreDiffJSON(t, testCase),
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

func passingRuntimeWorkflowArtifacts(t *testing.T, testCase qualityCase, operations []runtimeOperation) (json.RawMessage, json.RawMessage) {
	t.Helper()
	candidates := []runtimeCandidate{}
	decisions := []runtimeDecision{}
	if testCase.CaseID == "FEAT005-H-UPDATE-CORRECTION" {
		candidates = append(candidates, runtimeCandidate{
			ID:                "cand-h-1",
			EvidenceKind:      "correction",
			Strength:          "strong",
			ProposedAssertion: "unbuffered channelへのsendはreceiverが受信可能になったときに完了する",
			EvidenceRawText:   qualityCorrectionEvidenceText(t, testCase),
			Scope:             map[string]string{"language": "Go"},
			Temporal:          json.RawMessage("null"),
			SearchQueries: []string{
				"unbuffered channelへのsendはreceiverが受信可能になったときに完了する",
				"unbuffered channelのsendはreceiverが受信可能になる前にも完了する",
			},
		})
		decisions = append(decisions, runtimeDecision{
			CandidateID:       "cand-h-1",
			Action:            "revise",
			TargetAssertionID: "as-h",
			ExecutionStatus:   "applied",
			Operations:        operations,
		})
	}
	candidatesJSON, err := json.Marshal(candidates)
	if err != nil {
		t.Fatal(err)
	}
	decisionsJSON, err := json.Marshal(decisions)
	if err != nil {
		t.Fatal(err)
	}

	return candidatesJSON, decisionsJSON
}

func runtimeStoreDiffJSON(t *testing.T, testCase qualityCase) json.RawMessage {
	t.Helper()
	value, err := json.Marshal(runtimeExpectedStoreDiff(testCase))
	if err != nil {
		t.Fatal(err)
	}

	return value
}
