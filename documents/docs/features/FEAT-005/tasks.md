# FEAT-005 実装タスク

> **前提:** 利用者の「taskごとにissueを登録して下さい」（2026-08-15）を、独立レビュー`pass`後の修正版詳細設計の明示承認およびTask分解開始指示として記録した。

## Design Readiness Audit

| 監査項目 | 結果 | 根拠 |
| --- | --- | --- |
| 利用者による詳細設計承認 | pass | `.ai/workflow/state.yaml` の `features.FEAT-005.human_design_review: approved` |
| 独立設計レビュー | pass | [design-review.md](design-review.md) の総合判定 `pass` |
| Scenario A〜Jと派生Xの固定入力・seed・oracle | pass | [scenario-catalog.md](design/scenario-catalog.md) |
| 層別oracleと最初の不一致層 | pass | [design.md](design.md#acceptance--test-design) |
| 隔離Store・実CLIプロセス境界 | pass | [design.md](design.md#implementation-deliverables--placement)、[design.md](design.md#security--nfr-considerations) |
| 新規DB・migration・公開CLI/API | not_applicable | 新設・変更せず、既存Knowledge CLIとWorkflowを消費・観測する。 |
| 既存CLI operation資料 | complete | [FEAT-001 Command Catalog](../FEAT-001/design/command-catalog.md)とoperation資料を正規根拠として参照する。 |

**結論:** pass。ケースの意味、seed、期待結果、停止時の扱い、Store差分、既存公開境界を実装者が再決定する未確定事項はない。

## 固定済みの実装対象と配置

実装本体は親リポジトリのテスト専用領域に置く。既存Workflow Skillは被観測対象であり、通常利用時のコンテキストを増やさないため変更しない。`documents/docs/features/FEAT-005/`は正規の設計・handoff資料であり、実行する受入suite本体ではない。

| 配置先（親リポジトリ基準） | 作成・変更 | 責務 |
| --- | --- | --- |
| `testdata/fixtures/` 配下のFEAT-005 Fixture領域 | 作成 | A〜J／Xの固定入力、隔離Store seed、期待oracle、期待Store差分を保持する。 |
| `test/integration/` | 作成・変更 | 実`knowledge` binary、隔離SQLite Store、stdout/stderr、exit code、履歴・差分をプロセス境界で観測する。 |

`cmd/knowledge/`、`internal/`、SQLite migration、Knowledge CLIの公開operation・JSON、公開設定は変更しない。

## タスク一覧

| ID | タイトル | 論理領域 | 依存 |
| --- | --- | --- | --- |
| TASK-005-01 | 固定受入Fixtureと隔離Store契約を提供する | fixture / test data | FEAT-001〜004 implementation handoff |
| TASK-005-02 | Knowledge Searchの層別oracleとReading Value参照を検証する | workflow verification: search / contract reference | TASK-005-01 |
| TASK-005-03 | Knowledge AcquisitionとKnowledge Updateの層別oracleを検証する | workflow verification: acquisition / update | TASK-005-01 |
| TASK-005-04 | プロセス境界の横断受入suiteと診断を完成する | integration verification | TASK-005-02、TASK-005-03 |

## TASK-005-01: 固定受入Fixtureと隔離Store契約を提供する

- **目的:** Scenario A〜Jおよび派生Xを、再現可能かつ互いに独立した入力・seed・期待oracleとして提供する。
- **関連要件:** REQ-021、BR-002、BR-003、BR-005、BR-008、BR-009、NFR-001〜006、DEC-FEAT-016。
- **論理領域:** fixture / test data。
- **作業内容:** [固定受入Caseカタログ](design/scenario-catalog.md)を正として、各`case_id`のEpisode、Claim、隔離Store seed、期待成果物、期待Store差分、契約参照をテスト専用Fixtureへ実体化する。caseごとの一時Store初期化と終了後の破棄を可能にする。
- **受入条件:**
  - A〜Jおよび`FEAT005-X-SEARCH-TECHNICAL-FAILURE`／`FEAT005-X-SEARCH-CANCELED`が一意な`case_id`で選択・全件実行でき、catalogと別の入力・seed・合否規則を割り当てない。
  - 各caseは入力、seed、実行対象層、期待oracle、期待Store差分、参照契約を持つ。`store_diff:none`、`retain`、`add`はcatalogの意味で判定できる。
  - Fixtureは合成データだけを用い、外部記事、実ユーザーStore、認証情報、既定Store、共有Storeを利用しない。
  - 各実行は独立した一時Storeから開始し、終了後に破棄する。あるcaseのmutation、失敗、中断が別caseへ影響しない。
  - Hの訂正前seed、訂正後seed、旧Evidence／revision／Relation保持、新revision／訂正Evidence追加を別々に再現できる。
- **依存:** FEAT-001〜004 implementation handoff。
- **対象外／注記:** FixtureのJSON入れ子、test helper、実行コードの構造はImplementation領域が決める。製品公開入力・出力・保存形式は追加しない。

## TASK-005-02: Knowledge Searchの層別oracleとReading Value参照を検証する

- **目的:** 固定Claimに対するSearch Traceの結論・根拠・停止と、A／G／I／Jの既存Reading Value検証契約への参照完全性を検証する。
- **関連要件:** REQ-021、BR-002、BR-003、BR-009、NFR-002〜006、DEC-FEAT-016〜018。
- **論理領域:** workflow verification: search / contract reference。
- **作業内容:** テスト専用Fixtureと統合suiteから既存`knowledge-search`の成果物を観測し、A〜E、H再評価、I、J、派生XのSearch／Trace oracleを照合する。A、G、I、JはReading Valueを起動せず、`reading_value_reference`がFEAT-003の指定節へ一意に対応することだけを確認する。
- **受入条件:**
  - A〜EおよびH再評価で、catalogに固定された`no_evidence`、`known`、`partially_known`、`contradicted`、`outdated`、Confidence、根拠ID、Search Traceのoperation・query・結果ID・Evidence ID・Budget・停止理由を照合できる。
  - A、G、I、Jの`reading_value_reference`が、catalogで指定されたFEAT-003 V-002またはV-004へ一意に遡れる。Reading Valueを起動せず、その成果物・推奨・Assessment MapをFixture、oracle、Case Resultへ含めない。
  - 技術失敗Xでは`storage_error`、中断Xではresponse開始前・無出力のexit 130、Trace停止位置、Assessment不在、後続層未実行を観測できる。成功結論・推奨へ変換しない。
  - 検索不要のF／GにSearch Traceを要求せず、対象外層の未実行を合格へ読み替えない。
- **依存:** TASK-005-01。
- **対象外／注記:** Knowledge SearchおよびReading Valueの業務規則、Reading Value Skill、公開CLI／JSON、DB schemaは変更しない。

## TASK-005-03: Knowledge AcquisitionとKnowledge Updateの層別oracleを検証する

- **目的:** 取得禁止入力で更新しないこと、および訂正時に既存履歴を保持したまま既存CLI操作列で更新することを検証する。
- **関連要件:** REQ-021、BR-005、BR-008、NFR-001、NFR-003〜006、DEC-FEAT-016。
- **論理領域:** workflow verification: acquisition / update。
- **作業内容:** テスト専用Fixtureと統合suiteから既存`knowledge-acquisition`および`knowledge-update`の成果物を観測し、F／Gの空Candidate・空Decision・mutationなし、Hの訂正Candidate、検索・取得・revision・Evidence追加・履歴保全を既存Workflow／CLI契約へ対応付ける。
- **受入条件:**
  - Fの質問だけ、およびGのAI説明だけではCandidateが空であり、Update Decision、CLI operation、mutation、Store差分がすべてないことを照合できる。
  - Hでは`cand-h-1`のkind、strength、Scope、Temporal、検索入力を固定どおり観測し、`search-text → get → get-evidence → revise → attach-evidence`の操作列と各既存CLI結果を照合できる。
  - Hは旧Evidence、旧revision、Relationを削除せず、訂正Evidenceと新revisionを追加するStore差分を照合できる。更新後seedのSearchが訂正Evidenceを根拠に`known`となることを追跡できる。
  - Updateのread/searchとmutationは同じUpdate層のprocess記録・Store差分で診断し、AcquisitionとUpdateを単一の層へ混在させない。
  - F／G／Hを含む対象caseで、失敗または中断後の後続層・後続operationを実行せず、case resultに未実行を記録できる。
- **依存:** TASK-005-01。
- **対象外／注記:** 新CLI operation、意味検索、補償削除、再実行ledger、既存Workflowの判断規則変更は行わない。

## TASK-005-04: プロセス境界の横断受入suiteと診断を完成する

- **目的:** 同じcase IDでFixture入力から各層の観測までを連結し、実CLIプロセス境界、Store差分、最初の不一致層、未実行層を再現可能に報告する。
- **関連要件:** REQ-021、BR-002、BR-003、BR-005、BR-008、BR-009、NFR-001〜006、DEC-FEAT-016。
- **論理領域:** integration verification。
- **作業内容:** 実`knowledge` binaryとcaseごとの隔離SQLite Storeを使う統合suiteを完成し、CLIのargv、stdout、stderr、exit code、実行順、mutation数、履歴・差分を既存公開契約に照らして観測する。Case Resultへ層別照合結果、最初の不一致層、根拠参照、未実行層または`not_run`理由を記録する。
- **受入条件:**
  - A、F、G、H、I、JのEnd-to-End必須caseを、同じ`case_id`で入力から最終観測まで追跡できる。B〜EおよびH再評価は実CLIと隔離Storeで既定oracleを検証する。
  - CLI／Store、Knowledge Search、Knowledge Acquisition、Knowledge Update、End-to-Endを設計上の照合順で評価し、最初の不一致だけを原因層として記録する。Reading Value参照不備は別途記録し、対象外層はpassへ変換しない。
  - 派生Xの`storage_error`とexit 130を実プロセス境界で確認し、部分Trace、停止、Assessment不在、後続未実行、mutation 0、`store_diff:none`を報告する。親Orchestration中断時だけ部分Traceを出力・保存しない既存契約も保持する。
  - `failed`または`not_run`はScenario合格に数えず、最初の不一致層または実行不能理由を残す。
  - 公開CLIのargv、stdout/stderr、JSON、exit codeを観測対象として保持しつつ、`cmd/knowledge/`、`internal/`、SQLite migration、公開CLI operation・JSON、公開設定を変更しない。
- **依存:** TASK-005-02、TASK-005-03。
- **対象外／注記:** 性能ベンチマーク、外部記事へのライブアクセス、プロダクト機能追加、Fixture外の自動回復・再実行は扱わない。

## Dependency Notes

- TASK-005-01が全caseで共有する固定入力・seed・oracleと隔離を提供する。
- TASK-005-02とTASK-005-03はFixture契約をそれぞれSearch／参照完全性、Acquisition／Updateの責務境界で消費するため、相互に独立して進められる。
- TASK-005-04は両方の層別oracleを同じcase IDで統合し、最初の不一致層とプロセス境界を確認するため、TASK-005-02とTASK-005-03に依存する。

## Implementation Readiness Review

| 項目 | 記録 |
| --- | --- |
| 判定 | pending |
| reviewer | 再実施待ち |
| 親リポジトリ俯瞰の範囲 | 実行用Skillを変更しない配置修正後に、テスト領域、被観測Workflow、Knowledge CLI、既存Fixture、検証コマンド、依存Feature handoffを独立・読み取り専用で再確認する。 |
| 根拠 | 先行レビューは実行用Skillの変更を含む草案を対象にしていたため、この配置修正に対しては再利用しない。 |
