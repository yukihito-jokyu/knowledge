# FEAT-002 実装タスク

> **前提:** 利用者の「はい」（2026-08-15）を、更新済み詳細設計の明示承認およびTask分解開始指示として記録した。独立設計レビューは `pass` である。

## Design Readiness Audit

| 監査項目 | 結果 | 根拠 |
| --- | --- | --- |
| 利用者による詳細設計承認 | pass | `.ai/workflow/state.yaml` の `features.FEAT-002.human_design_review: approved` |
| 独立設計レビュー | pass | [design-review.md](design-review.md) の総合判定 `pass` |
| Claim入力、状態分類、停止規則 | pass | [design.md](design.md)「Behavioral Scenarios」「State Classification」「Budget / Stop Conditions」 |
| 実行形態、AssessmentとTraceの成果物契約 | pass | [knowledge-search-workflow.md](design/knowledge-search-workflow.md)、[assessment-template.md](design/assessment-template.md)、[search-trace.md](design/search-trace.md)、[DEC-FEAT-009](decisions/DEC-FEAT-009.md) |
| CLI operation・JSON・error・exit契約 | pass | [design.md](design.md)「既存CLI利用契約」、FEAT-001の[Command Catalog](../FEAT-001/design/command-catalog.md) |
| SQLite・migration・公開設定の変更 | not_applicable | 本Featureは既存CLIを消費するだけで、DB書込み、migration、公開設定を追加しない。 |
| flowchart・sequence diagram | pass | [design.md](design.md)「Interaction」 |
| operation別資料とDB access map | pass | 既存operationを変更せず利用するため、FEAT-001の[operation資料](../FEAT-001/design/commands/)と[DB access map](../FEAT-001/design/database-access-map.md)を正規根拠とする。 |
| Fixture・受入観点 | pass | [design.md](design.md)「Acceptance / Test Design」およびFEAT-005が所有するScenario A〜J横断評価 |

**結論:** pass。実装者が新しい業務規則、公開CLI/JSON、永続化形式、または探索上限を再決定する未確定事項はない。

## 固定済み制約

- Knowledge CLIの既存11操作だけを利用し、CLIへ意味判断、検索戦略、Evidence強度、Knowledge Assessmentを移さない。
- `search-semantic`、Embedding、Vector Indexは利用しない。
- Claimごとの上限はDEC-FEAT-008どおり、CLI呼出し12回、探索4回、補助探索4回、Evidence取得4回、Relation深さ1、矛盾・時点探索の合計2回とする。
- AssessmentとSearch TraceはMarkdown成果物として分離し、初期提供では永続化、新しいJSON API、公開Budget設定を追加しない。
- CLIの技術失敗は知識状態へ変換しない。response開始前のexit 130は`canceled`として、error JSONを捏造せず中断を伝播する。

## タスク一覧

| ID | タイトル | 論理領域 | 依存 |
| --- | --- | --- | --- |
| TASK-002-01 | Codex Knowledge Searchの探索・評価を実現する | Codex workflow: exploration / evaluation | FEAT-001 implementation handoff |
| TASK-002-02 | 一時成果物と失敗・中断境界を実現する | Codex workflow: artifact / failure boundary | TASK-002-01 |
| TASK-002-03 | 固定Budget・状態分類・中断を横断検証する | verification | TASK-002-01, TASK-002-02 |

## TASK-002-01: Codex Knowledge Searchの探索・評価を実現する

- **目的:** 一つの独立評価可能なClaimに対して、既存Knowledge CLIを有限に反復利用し、Evidenceに基づく一つのKnowledge Assessment状態へ到達できるようにする。
- **関連要件:** REQ-004、REQ-005、REQ-019〜021、BR-002、BR-003、BR-006、BR-007、BR-010、NFR-001〜003、NFR-005〜006、DEC-REQ-001、DEC-FEAT-008、DEC-FEAT-009。
- **論理領域:** Codex workflow: exploration / evaluation。
- **作業内容:** 呼出側から一つのTarget Claimと明示されたScope・version・時点文脈を受けて、Codex AI Runtime内でClaimごとに同期実行する。既存Knowledge CLIを決定論的な検索・取得capabilityとして利用し、直接字句探索から構成要素またはConcept、必要なRelation、Evidence、矛盾候補、時点差分へ限定拡張する。Evidence原文、Scope、時点、Claimとの意味対応をCodexが評価し、優先規則に従って7状態、Confidence、Known、Knowledge Gap、矛盾・時点差分を決定する。探索を停止する条件と固定Budgetを一実行内で一貫して制御する。
- **受入条件:**
  - 最初に`search-text`を利用し、未解決なときだけ構成要素または明示Concept、Relation、Evidence、矛盾、時点差分を必要最小限に探索する。
  - 既存CLIを同期的に利用し、各呼出し前に既存operation契約のoption、enum、selector、値の組を確認する。ClaimにないScope・version・時点filterは推測して追加しない。
  - `known`、`partially_known`、`inferable`、`contradicted`、`outdated`、`no_evidence`、`uncertain`のいずれか一つを、設計の優先規則とEvidenceに基づいて決定する。
  - Evidence未観測をユーザーの未知へ変換せず、Concept一致・Relationの存在・検索候補の類似だけを`known`の根拠にしない。
  - `inferable`は必要な構成知識と明示Relationによる導出をEvidenceで支持できる場合にだけ選び、直接Evidenceのある`known`と区別する。
  - `search-temporal`はConceptまたはScope selectorがある場合だけ利用し、時点条件だけで実行しない。
  - 同一論理queryを再実行せず、強い結論、飽和、増分なし、探索経路なし、またはBudget到達で停止する。
- **依存:** FEAT-001のimplementation handoff（既存CLI operation・JSON/error/exit契約が利用可能であること）。
- **対象外／注記:** Articleの取得・Claim分解・Reading Valueの決定、Semantic Search、CLI/DB公開契約の変更は対象外。実装先はCodex側Knowledge Search workflowであり、Knowledge CLIのGo実装ではない。正規根拠は[design.md](design.md)、[knowledge-search-workflow.md](design/knowledge-search-workflow.md)、[DEC-FEAT-008](decisions/DEC-FEAT-008.md)、[DEC-FEAT-009](decisions/DEC-FEAT-009.md)。

## TASK-002-02: 一時成果物と失敗・中断境界を実現する

- **目的:** 評価結果と診断過程を別々のMarkdown成果物として出力し、CLI技術失敗と利用者中断を評価結果から明確に分離する。
- **関連要件:** REQ-005、REQ-019〜021、BR-002、BR-010、NFR-001〜003、NFR-005、CON-002、CON-003、DEC-FEAT-008、DEC-FEAT-009。
- **論理領域:** Codex workflow: artifact / failure boundary。
- **作業内容:** 正常終了時に状態、Confidence、Known、Knowledge Gap、Supporting Evidence、Relation、矛盾、時点差分、結論、Trace参照をAssessmentへ構成し、別のTraceには操作、最小限の入力要約、結果/Evidence ID、増分、Budget、停止理由を残す。両者は実行内で一時的に受け渡し、Codex側の返却としてAssessmentと対応Trace参照を呼出側へ渡す。初期提供では保存、公開I/O、JSON APIを追加しない。既存CLIの成功JSON、既知error、exit 130・無出力を区別して呼出側へ伝播する。
- **受入条件:**
  - 正常終了時だけAssessmentとTrace参照を返し、Assessmentは[assessment-template.md](design/assessment-template.md)の必須節を持ち、Traceへ全操作ログを重複させない。
  - Traceは[search-trace.md](design/search-trace.md)のBudget、Steps、Stop Reason、必要時のTechnical FailureまたはCancellationを記録し、Evidence原文を不要に複製しない。
  - `validation_error`、`not_found`、`storage_error`、`internal_error`、JSONプロトコル不整合ではAssessmentを出力せず、`technical_failure`のTraceと失敗を呼出側へ返す。`validation_error`は再試行しない。
  - response開始前のexit 130・無出力ではerror codeを作らず、Assessmentなしで`canceled`、operation、exit code、実行済みstepを部分Traceへ記録して中断を伝播する。Parent Orchestration自身も中断済みならTraceを出力・保存しない。
  - 空検索結果、空Store、未登録Conceptは成功結果として扱い、技術失敗と混同しない。
- **依存:** TASK-002-01。
- **対象外／注記:** Trace/Assessmentの永続保存、新しいCLI operation・JSON API・公開設定、Reading Valueによる再評価は対象外。

## TASK-002-03: 固定Budget・状態分類・中断を横断検証する

- **目的:** FEAT-002の評価規則、有限探索、成果物境界、失敗伝播が、既存CLIとの実行境界で再現可能に確認できるようにする。
- **関連要件:** REQ-004、REQ-005、REQ-019〜021、BR-002、BR-003、BR-006、BR-007、BR-010、NFR-001〜003、NFR-005〜006、DEC-FEAT-008、DEC-FEAT-009。
- **論理領域:** verification。
- **作業内容:** 既存CLIのプロセス境界をstubまたはfixtureで観測し、空Store、完全一致、構成知識、導出可能、矛盾、旧知識、競合Evidence、EvidenceなしRelation、技術失敗、キャンセルを含むClaim別fixtureに対して、呼出し順序・入力・上限・既存error JSON・exit code、Assessment、Trace、Parentへの結果が設計どおりであることを検証する。Scenario A〜Jをまたぐ横断Fixture評価はFEAT-005の所有範囲とし、本Taskは個別Claim評価のoracleを提供する。
- **受入条件:**
  - 空Storeは`no_evidence`となり、未知断定を含まず、直接探索と停止理由をTraceから追跡できる。
  - strong Evidence、構成知識、明示Relationによる導出、優位な反対Evidence、新しい訂正、解消不能な競合Evidenceが、それぞれ設計どおりの状態とConfidenceになる。
  - Relationのみ、質問のみ、AI説明のみ、EvidenceなしAssertionを`known`の根拠にできない。
  - CLI呼出し13回目、Evidence取得5件目、Relation深さ2、同一論理queryの再実行、矛盾・時点探索の合計3回目を実行しない。Traceから各Budgetと停止理由を確認できる。
  - CLIの各技術失敗はAssessmentなしで`technical_failure`として、exit 130・無出力はerror codeなしの`canceled`として、互いに区別してParentへ伝播する。
  - AssessmentとTraceが分離され、Traceが不要なEvidence原文を複製しないことを確認する。
- **依存:** TASK-002-01、TASK-002-02。
- **対象外／注記:** FEAT-005所有のScenario A〜Jをまたぐ評価基盤そのもの、Semantic Searchの評価、Article入力UIは対象外。

## Dependency Notes

- TASK-002-01は、FEAT-001が提供する既存CLIの実装可能状態を前提とする。
- TASK-002-02は、TASK-002-01が決定した探索結果・停止理由を成果物化するため後続とする。
- TASK-002-03は、探索制御と成果物/失敗境界をともに観測するため、両Taskに依存する。
