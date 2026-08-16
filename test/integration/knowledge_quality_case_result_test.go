package integration_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

const (
	qualityCLIBoundaryExecutionMode = "cli_boundary"
	qualityRuntimeExecutionMode     = "runtime_acceptance"
)

var qualityLayerOrder = []string{
	"cli_store",
	"knowledge_search",
	"knowledge_acquisition",
	"knowledge_update",
	"end_to_end",
}

// qualityCaseResultは、同じ固定caseを異なる受入境界で照合した結果を表す。
// 結果はテスト出力だけに残し、FixtureやStoreへ保存しない。
type qualityCaseResult struct {
	CaseID             string
	ExecutionMode      string
	Status             string
	ExecutedLayers     []string
	FirstMismatchLayer string
	NotExecutedLayers  []string
	NotRunReason       string
}

type qualityCaseResultKey struct {
	CaseID        string
	ExecutionMode string
}

type qualityLayerObservation struct {
	Layer   string
	Matched bool
}

func diagnoseQualityCaseResult(testCase qualityCase, executionMode string, observations []qualityLayerObservation, notRunReason string) qualityCaseResult {
	result := qualityCaseResult{
		CaseID:             testCase.CaseID,
		ExecutionMode:      executionMode,
		Status:             "pass",
		ExecutedLayers:     []string{},
		FirstMismatchLayer: "none",
		NotExecutedLayers:  append([]string{}, testCase.Expected.NotExecuted...),
	}
	if notRunReason != "" {
		result.Status = "not_run"
		result.NotRunReason = notRunReason
		result.NotExecutedLayers = qualityOrderedNotExecutedLayers(testCase, "")

		return result
	}
	for _, observation := range observations {
		if !containsQualityLayer(testCase.Layers, observation.Layer) || containsQualityLayer(result.ExecutedLayers, observation.Layer) {
			continue
		}
		result.ExecutedLayers = append(result.ExecutedLayers, observation.Layer)
		if !observation.Matched {
			result.Status = "failed"
			result.FirstMismatchLayer = observation.Layer
			result.NotExecutedLayers = qualityOrderedNotExecutedLayers(testCase, observation.Layer)

			return result
		}
	}

	return result
}

func qualityOrderedNotExecutedLayers(testCase qualityCase, mismatch string) []string {
	result := []string{}
	include := mismatch == ""
	for _, layer := range qualityLayerOrder {
		if layer == mismatch {
			include = true

			continue
		}
		if containsQualityLayer(testCase.Expected.NotExecuted, layer) || (include && containsQualityLayer(testCase.Layers, layer)) {
			result = append(result, layer)
		}
	}

	return result
}

type qualityLayerCheck struct {
	Layer string
	Check func() error
}

func executeQualityCaseLayers(testCase qualityCase, checks []qualityLayerCheck) (qualityCaseResult, error) {
	observations := make([]qualityLayerObservation, 0, len(checks))
	for _, check := range checks {
		err := check.Check()
		observations = append(observations, qualityLayerObservation{
			Layer:   check.Layer,
			Matched: err == nil,
		})
		if err != nil {
			return diagnoseQualityCaseResult(testCase, qualityCLIBoundaryExecutionMode, observations, ""), err
		}
	}

	return diagnoseQualityCaseResult(testCase, qualityCLIBoundaryExecutionMode, observations, ""), nil
}

func containsQualityLayer(layers []string, want string) bool {
	for _, layer := range layers {
		if layer == want {
			return true
		}
	}

	return false
}

func qualityCaseResultKeyOf(result qualityCaseResult) qualityCaseResultKey {
	return qualityCaseResultKey{
		CaseID:        result.CaseID,
		ExecutionMode: result.ExecutionMode,
	}
}

func recordQualityCaseResult(results map[qualityCaseResultKey]qualityCaseResult, result qualityCaseResult) error {
	key := qualityCaseResultKeyOf(result)
	if _, exists := results[key]; exists {
		return fmt.Errorf("Case Resultを上書きできません: case_id=%s execution_mode=%s", key.CaseID, key.ExecutionMode)
	}
	results[key] = result

	return nil
}

func TestQualityCaseResultDiagnosis(t *testing.T) {
	testCase := qualityCase{
		CaseID: "FEAT005-H-UPDATE-CORRECTION",
		Layers: []string{
			"cli_store",
			"knowledge_acquisition",
			"knowledge_update",
			"end_to_end",
		},
	}
	testCase.Expected.NotExecuted = []string{"knowledge_search"}

	failed := diagnoseQualityCaseResult(testCase, qualityCLIBoundaryExecutionMode, []qualityLayerObservation{
		{
			Layer:   "cli_store",
			Matched: true,
		},
		{
			Layer:   "knowledge_acquisition",
			Matched: true,
		},
		{
			Layer:   "knowledge_update",
			Matched: false,
		},
		{
			Layer:   "end_to_end",
			Matched: true,
		},
	}, "")
	if failed.Status != "failed" || !reflect.DeepEqual(failed.ExecutedLayers, []string{"cli_store", "knowledge_acquisition", "knowledge_update"}) || failed.FirstMismatchLayer != "knowledge_update" || !reflect.DeepEqual(failed.NotExecutedLayers, []string{"knowledge_search", "end_to_end"}) {
		t.Fatalf("最初の不一致と後続未実行を記録できません: %#v", failed)
	}

	notRun := diagnoseQualityCaseResult(testCase, qualityRuntimeExecutionMode, nil, "Codex Runtimeが利用できません")
	if notRun.Status != "not_run" || len(notRun.ExecutedLayers) != 0 || notRun.NotRunReason == "" || notRun.FirstMismatchLayer != "none" || !reflect.DeepEqual(notRun.NotExecutedLayers, []string{"cli_store", "knowledge_search", "knowledge_acquisition", "knowledge_update", "end_to_end"}) {
		t.Fatalf("not_run理由を記録できません: %#v", notRun)
	}

	results := map[qualityCaseResultKey]qualityCaseResult{}
	if err := recordQualityCaseResult(results, failed); err != nil {
		t.Fatal(err)
	}
	if err := recordQualityCaseResult(results, notRun); err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("同じcase IDの異なる実行境界を分離できません: %#v", results)
	}
	if err := recordQualityCaseResult(results, failed); err == nil {
		t.Fatal("同一論理キーのCase Result上書きを拒否できません")
	}
}

func TestQualityCaseResultRecordsProcessBoundaryFailure(t *testing.T) {
	fixture := readQualityFixture(t)
	var testCase qualityCase
	for _, candidate := range fixture.Cases {
		if candidate.CaseID == "FEAT005-B-SEARCH-EXACT" {
			testCase = candidate

			break
		}
	}
	if testCase.CaseID == "" {
		t.Fatal("B caseがfixtureにありません")
	}
	root := t.TempDir()
	store := defaultStoreConfiguration(t, root)
	seedQualityStore(t, store.Path, fixture.Seeds[testCase.SeedID])
	stdout, stderr, err := runCommand(buildCLI(t, true), store.Environment, []string{"unknown-operation"})
	if err == nil || stdout != "" || stderr == "" {
		t.Fatalf("失敗を注入した実CLIプロセス境界を観測できません: stdout=%q stderr=%q err=%v", stdout, stderr, err)
	}
	result := diagnoseQualityCaseResult(testCase, qualityCLIBoundaryExecutionMode, []qualityLayerObservation{
		{
			Layer:   "cli_store",
			Matched: false,
		},
	}, "")
	if result.Status != "failed" || !reflect.DeepEqual(result.ExecutedLayers, []string{"cli_store"}) || result.FirstMismatchLayer != "cli_store" || !reflect.DeepEqual(result.NotExecutedLayers, []string{"knowledge_search", "knowledge_acquisition", "knowledge_update", "end_to_end"}) {
		t.Fatalf("実CLI失敗のCase Resultが不正です: %#v", result)
	}
}

func TestQualityCaseResultRecordsActualLayerOracleMismatch(t *testing.T) {
	fixture := readQualityFixture(t)
	var testCase qualityCase
	for _, candidate := range fixture.Cases {
		if candidate.CaseID == "FEAT005-B-SEARCH-EXACT" {
			testCase = candidate

			break
		}
	}
	if testCase.CaseID == "" {
		t.Fatal("B caseがfixtureにありません")
	}
	root := t.TempDir()
	store := defaultStoreConfiguration(t, root)
	seedQualityStore(t, store.Path, fixture.Seeds[testCase.SeedID])
	binary := buildCLI(t, true)
	result, err := executeQualityCaseLayers(testCase, []qualityLayerCheck{
		{
			Layer: "cli_store",
			Check: func() error {
				_, stderr, commandErr := runCommand(binary, store.Environment, testCase.Expected.CLI.Arguments)
				if commandErr != nil || stderr != "" {
					return fmt.Errorf("CLI/Store oracle: stderr=%q err=%w", stderr, commandErr)
				}

				return nil
			},
		},
		{
			Layer: "knowledge_search",
			Check: func() error {
				stdout, _, commandErr := runCommand(binary, store.Environment, testCase.Expected.CLI.Arguments)
				if commandErr != nil {
					return commandErr
				}
				var response struct {
					Data struct {
						Results []struct {
							ID string `json:"assertion_id"`
						} `json:"results"`
					} `json:"data"`
				}
				if err := json.Unmarshal([]byte(stdout), &response); err != nil {
					return err
				}
				if len(response.Data.Results) != 1 || response.Data.Results[0].ID != "intentionally-wrong-id" {
					return fmt.Errorf("Search oracle mismatch: results=%#v", response.Data.Results)
				}

				return nil
			},
		},
	})
	if err == nil {
		t.Fatal("期待値不一致を注入できません")
	}
	if result.Status != "failed" || !reflect.DeepEqual(result.ExecutedLayers, []string{"cli_store", "knowledge_search"}) || result.FirstMismatchLayer != "knowledge_search" || !reflect.DeepEqual(result.NotExecutedLayers, []string{"knowledge_acquisition", "knowledge_update", "end_to_end"}) {
		t.Fatalf("実層oracle不一致のCase Resultが不正です: %#v", result)
	}
}

func TestQualityCaseResultRecordsHUpdateOracleMismatch(t *testing.T) {
	fixture := readQualityFixture(t)
	var testCase qualityCase
	for _, candidate := range fixture.Cases {
		if candidate.CaseID == "FEAT005-H-UPDATE-CORRECTION" {
			testCase = candidate

			break
		}
	}
	if testCase.CaseID == "" {
		t.Fatal("H caseがfixtureにありません")
	}
	root := t.TempDir()
	store := defaultStoreConfiguration(t, root)
	seedQualityStore(t, store.Path, fixture.Seeds[testCase.SeedID])
	binary := buildCLI(t, true)
	result, err := executeQualityCaseLayers(testCase, []qualityLayerCheck{
		{
			Layer: "cli_store",
			Check: func() error {
				_, stderr, commandErr := runCommand(binary, store.Environment, testCase.Expected.CLI.Arguments)
				if commandErr != nil || stderr != "" {
					return fmt.Errorf("CLI/Store oracle: stderr=%q err=%w", stderr, commandErr)
				}

				return nil
			},
		},
		{
			Layer: "knowledge_acquisition",
			Check: func() error {
				var input struct {
					Episode struct {
						UserContributions []struct {
							SourceText string `json:"source_text"`
						} `json:"user_contributions"`
					} `json:"episode"`
				}
				if err := json.Unmarshal(testCase.Input, &input); err != nil {
					return err
				}
				if len(input.Episode.UserContributions) != 1 || input.Episode.UserContributions[0].SourceText == "" {
					return fmt.Errorf("Acquisition oracle: %#v", input.Episode)
				}

				return nil
			},
		},
		{
			Layer: "knowledge_update",
			Check: func() error {
				_, stderr, commandErr := runCommand(binary, store.Environment, []string{"revise", "--assertion-id", "as-h", "--normalized-text", "unbuffered channelのsendはreceiverが受信可能になるまでblockする"})
				if commandErr != nil || stderr != "" {
					return fmt.Errorf("Update CLI oracle: stderr=%q err=%w", stderr, commandErr)
				}

				return fmt.Errorf("Update oracle mismatch: injected expected revision=999")
			},
		},
	})
	if err == nil {
		t.Fatal("H Update期待値不一致を注入できません")
	}
	if result.Status != "failed" || !reflect.DeepEqual(result.ExecutedLayers, []string{"cli_store", "knowledge_acquisition", "knowledge_update"}) || result.FirstMismatchLayer != "knowledge_update" || !reflect.DeepEqual(result.NotExecutedLayers, []string{"knowledge_search", "end_to_end"}) {
		t.Fatalf("H Update不一致のCase Resultが不正です: %#v", result)
	}
}

func TestQualityFailureAndCancellationKeepNormalStoreUnchanged(t *testing.T) {
	fixture := readQualityFixture(t)
	binary := buildCLI(t, true)
	for _, caseID := range []string{
		"FEAT005-X-SEARCH-TECHNICAL-FAILURE",
		"FEAT005-X-SEARCH-CANCELED",
	} {
		t.Run(caseID, func(t *testing.T) {
			var testCase qualityCase
			for _, candidate := range fixture.Cases {
				if candidate.CaseID == caseID {
					testCase = candidate

					break
				}
			}
			if testCase.CaseID == "" {
				t.Fatalf("X caseがfixtureにありません: %s", caseID)
			}
			root := t.TempDir()
			store := defaultStoreConfiguration(t, root)
			seedQualityStore(t, store.Path, fixture.Seeds[testCase.SeedID])
			before := qualityStateOf(t, store.Path)
			result, err := executeQualityCaseLayers(testCase, qualityCaseLayerChecks(t, binary, store, testCase, before))
			if err != nil || result.Status != "pass" {
				t.Fatalf("Xの実プロセスoracleが失敗: result=%#v err=%v", result, err)
			}
			if err := qualityStoreDiffMatches(t, store, before, testCase); err != nil {
				t.Fatalf("Xの失敗または中断後に通常Storeが変化: %v", err)
			}
		})
	}
}
