# Codex Knowledge Search Workflow 契約

## 実行形態

このワークフローはCodex AI Runtime内で、Claimごとに一回実行する。Knowledge CLIを子ワークフローや評価器として置き換えず、既存JSON CLIを決定論的な検索・取得capabilityとして同期的に呼び出す。CLIのsuccess/error JSON、stdout/stderr、exit codeはFEAT-001の既存operation契約をそのまま消費する。

実行単位は一つのTarget Claimであり、複数Claimの並列化、再試行の共通キュー、実行履歴の永続化は本Featureの範囲外とする。

## 入力と開始条件

- 呼出側は、独立評価可能なTarget Claim本文と、Claimに明示されたScope、対象version、時点文脈だけを渡す。
- Scope、version、時点が未指定の場合、ワークフローはfilterやselectorを推測して追加しない。
- 実行前に、利用予定のCLI operationの既存option、enum、selector、値の組を確認する。確認できない入力はCLI呼出しを開始せず、`technical_failure`として扱う。

## 実行とCLI利用

1. `search-text`で直接候補を一度検索する。
2. 取得結果とClaimの意味をCodexが比較し、未解決部分に限り構成要素または明示Conceptを選ぶ。
3. `search-text`または`search-concept`、必要時だけ`get`、`get-evidence`、`search-related`、`search-contradictions`、`search-temporal`を既存契約どおりに呼ぶ。
4. 候補、Evidence、Scope、時点、明示RelationをCodexが解釈し、[design.md](../design.md#state-classification)の優先規則で一つの状態を決める。

CLIは候補・Assertion・Evidenceを返すだけであり、同等性、Evidence強度、次query、停止、状態を判断しない。各呼出し直後に、結果ID、Evidence ID、増分、Budget消費、次の判断をSearch Traceの実行中状態へ追加する。同じ論理queryを再実行しない。

## Budgetと停止

Budget、探索順序、停止理由は[DEC-FEAT-008](../decisions/DEC-FEAT-008.md)と[design.md](../design.md#budget--stop-conditions)に従う。ワークフローは各CLI呼出しの前に上限を確認し、上限到達後に追加呼出しを行わない。停止時は、未解決部分をKnowledge Gapとして残す。

## 成果物と呼出側への返却

| 終了 | Codex側の返却 | 成果物 |
| --- | --- | --- |
| 正常終了 | 一つのKnowledge Assessmentと対応するSearch Trace参照 | [assessment-template.md](assessment-template.md)および[search-trace.md](search-trace.md)に従うMarkdown成果物。初期提供では永続化しない。 |
| `technical_failure` | Assessmentなしの評価失敗 | Search Traceに実行済みstep、operation、既存error code、再試行しない理由を残し、呼出側へ失敗を伝播する。 |
| `canceled` | Assessmentなしの中断 | response開始前のexit 130・無出力をoperationとexit codeのみの部分Traceとして扱い、中断を伝播する。Parent Orchestration自身も中断済みならTraceを出力・保存しない。 |

Reading Valueは正常終了時のAssessmentだけを入力として利用する。`technical_failure`または`canceled`を知識状態や`no_evidence`へ変換しない。TraceはAssessmentの根拠追跡と診断用であり、Evidence原文を不要に複製しない。

## 検証oracle

- 既存CLIのプロセス境界をstubまたはfixtureで観測し、呼出し順序、入力、上限、exit code、既存error JSONの伝播を確認する。
- Claim別fixtureで、7状態、Confidence、Known、Knowledge Gap、Evidence ID、停止理由が[design.md](../design.md#acceptance--test-design)に一致することを確認する。
- 13回目のCLI呼出し、5件目のEvidence取得、Relation深さ2、矛盾・時点探索合計3回目、同一論理queryの再実行がないことをTraceから確認する。
- `technical_failure`と`canceled`のどちらもAssessmentを返さず、空結果は成功結果として区別されることを確認する。

この契約は新しい公開I/Oを定義しない。Codex内の論理的なワークフロー境界と、その観測oracleを定める。
