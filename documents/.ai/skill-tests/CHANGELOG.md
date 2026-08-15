# Skill Maintenance Changelog

Skill修正時に追記する。

## 2026-08-15

- Target Skill: feature-design
- Symptom: 人間判断が必要な設計上の未確定事項について、Decisionを作成・提示する前に会話だけで選択を求めた。
- Root Cause: Human Decision PointsがDecisionの作成・説明・会話での参照順序を必須化していなかった。
- Change: 人間への質問前に、平易な未決Decisionへ機能・事実・未確定理由・選択肢の影響・推奨・反映先を記録し、会話ではそのDecisionを参照する手順を追加した。
- Regression Scenario: SKILLTEST-025
- Notes: 利用者のProduct DecisionをAIが決めず、回答後は同一Decisionを更新する。

- Target Skill: feature-design
- Symptom: 詳細設計が論理責務だけで完了し、実装成果物の具体的な配置先を利用者が後から質問しなければ確認できなかった。
- Root Cause: Implementation Deliverables / Placementを、利用者が明示的に求めた場合だけ作る任意規則にしていた。
- Change: 同節を詳細設計の必須節とし、既存規約で裏付けられるroot-relative配置・責務・依存・変更しない領域を常に記録する。配置を導けない場合はdesign_readyにしない。
- Regression Scenario: SKILLTEST-024
- Notes: function、symbol、framework API、根拠のない配置の指定は引き続き禁止する。

- Target Skill: task-breakdown, planning-orchestrator, implementation-handoff schema
- Symptom: Task分割とhandoffの後に、実装担当者が親リポジトリの実装先・既存依存・適用規約を俯瞰して実装可能性を確認する工程がなかった。
- Root Cause: Task Breakdownの完了条件が設計の完全性までで、実際の実装基盤・配置・依存契約との接続を独立に反証するpost-task reviewを定義していなかった。
- Change: Task分割後、親リポジトリrootから全体を俯瞰するよう明示したfreshなsubagentによる実装者視点レビューを必須化し、結果をtasks/handoffへ記録してからreadyにする規則とhandoff項目を追加した。
- Regression Scenario: SKILLTEST-023
- Notes: reviewerは読み取り専用であり、Taskやhandoffの編集、未承認の実装設計の追加を行わない。

- Target Skill: feature-design
- Symptom: Codex workflow Featureで、論理責務だけを記録し、実装担当者が作成するSkill成果物と固定済みの配置先を判断できなかった。
- Root Cause: implementation-neutralの規則が、Initial DesignまたはAGENTS.mdで固定済みの配置規約と、利用者が明示的に求めた実装配置の記録まで抑止していた。
- Change: 固定済み配置規約または利用者の明示要求がある場合、root-relative配置・成果物・責務・依存を`design.md`へ記録する規則を追加した。コード構造の指定は引き続き禁止する。
- Regression Scenario: SKILLTEST-022
- Notes: Product固有の配置名をSkillへ固定せず、既存規約または利用者要求を根拠にする。

- Target Skill: knowledge-search
- Symptom: ユースケース検証のfixture seedとcancel stubを、将来のプロダクト実行にも使うSkill配下へ置いた。
- Root Cause: 実行Skillの再利用資源と、Task-002-03の検証資源のArtifact Ownershipを分離せず、検証のための便利さを優先した。
- Change: fixture生成・中断stubを`documents/.ai/skill-tests/knowledge-search/scripts/`へ移し、実行Skillには読取りworkflowと同梱契約だけを残した。
- Regression Scenario: SKILLTEST-021
- Notes: 実行SkillはCLI write、fixture生成、検証script参照を行わない。

- Target Skill: knowledge-search
- Symptom: 実行に必要な評価規則とCLI操作契約を、Skill外の設計資料を読む`Required Reading`へ委ねていた。
- Root Cause: Generated Skillの自己完結性（Input Contract / Procedure）を検査せず、同梱`references/`を操作契約の配布先にしなかった。
- Change: 外部資料の必須参照を禁止し、評価手順は`SKILL.md`、CLI operation・option・JSON・exit codeは`references/cli-operations.md`へ固定した。
- Regression Scenario: SKILLTEST-019
- Notes: 実行時に読む正規契約はSkill本文と同梱リファレンスだけであり、CLIの実行結果は探索対象データとして扱う。

## 2026-08-13

- Target Skill: design-review, planning-orchestrator, task-breakdown, artifact-map, implementation-handoff schema
- Symptom: 設計者自身の監査だけで、要件・承認Decisionに根拠のない公開契約を含む詳細設計がhandoffへ進んだ。
- Root Cause: 独立した原典照合、AIの導出と人間Decisionの分類、資料間整合を必須にするreview gateがなかった。
- Change: freshなsubagentで実行するdesign-review Skillとレビュー成果物を追加し、pass済みの独立レビューをTask Breakdownとhandoffの必須条件にした。
- Regression Scenario: SKILLTEST-001, SKILLTEST-009, SKILLTEST-010
- Notes: reviewは設計を修正せず、根拠不明な契約をFeature Designまたは人間Decisionへ差し戻す。

- Target Skill: feature-design, task-breakdown, planning-orchestrator, implementation-handoff schema
- Symptom: 永続化・CLI Featureで論理契約が未完成のままTask Breakdownとimplementation_readyへ遷移した。
- Root Cause: Feature Designに契約完成度ゲートと図の選択基準がなく、Task BreakdownとOrchestratorに未完成設計を差し戻す監査がなかった。handoff schemaも契約の完成状態を表せなかった。
- Change: Feature性質に応じた契約完成度、フローチャート・シーケンス図の作成条件、Design Readiness Audit、handoffの契約参照を追加した。
- Regression Scenario: SKILLTEST-001, SKILLTEST-002, SKILLTEST-003
- Notes: Feature固有の要件・設計・Task成果物は変更していない。

## 2026-08-13

- Target Skill: feature-design, task-breakdown, planning-orchestrator, implementation-handoff schema
- Symptom: SQLiteとJSON CLIが採用済みでも、概念契約だけで設計完了となり、DDL相当とoperation別JSON Schemaが出力されなかった。
- Root Cause: Contract Completeness Gateが論理概念・操作の粒度で止まり、採用済みtechnology/transportに必要な形式契約を必須化していなかった。
- Change: SQLite／関係DBのDDL相当、JSON CLI／APIのoperation別wire schema・CLI exit codeをFeature Designの必須契約とし、Task BreakdownとOrchestratorの監査・handoff schemaへ追加した。
- Regression Scenario: SKILLTEST-004, SKILLTEST-005
- Notes: file path、symbol、framework API、library-specific implementationは引き続きImplementation領域に残す。

## 2026-08-13

- Target Skill: feature-design, task-breakdown, planning-orchestrator, artifact-map, implementation-handoff schema
- Symptom: DDLとoperation別wire schemaが存在しても、コマンドごとのI/O、DB read/write、transaction境界、図を実装者が逆引きできず、設計資料の可読性と実装可能性が不足した。
- Root Cause: 形式契約の完成条件が単一の設計節に留まり、operation単位の資料構造、DBリファレンス、access map、要求されたoperation別図の網羅性を規定していなかった。
- Change: Operation Documentation Coverage Gateと補助設計資料のOwnershipを追加し、Task Breakdown／Orchestrator／handoffでoperation資料とDB参照の監査・引継ぎを必須化した。
- Regression Scenario: SKILLTEST-006, SKILLTEST-007, SKILLTEST-008
- Notes: 図は全Featureや全operationに一律強制せず、Featureの性質または利用者の明示要求に応じて作成する。

## Entry Template

- Date:
- Target Skill:
- Symptom:
- Root Cause:
- Change:
- Regression Scenario:
- Notes:

## 2026-08-14

- Target Skill: task-breakdown, planning-orchestrator
- Symptom: 新規SQLite JSON CLIの公開契約と論理タスクが完成していても、実装言語、driver、配置責務、共通検証手順が未決定のままhandoffがreadyになった。
- Root Cause: Design Readiness Auditとimplementation_ready監査が論理契約・wire schemaだけを確認し、実装に必要な横断的技術基盤を確認していなかった。
- Change: 新規実行可能成果物、driver、migration、統合testを含むFeatureでは、既存コードまたはInitial Designに採用技術、依存方針、配置責務、共通検証手順があることをTask BreakdownとOrchestratorの必須監査に追加した。未決定ならInitial Designへ差し戻す。
- Regression Scenario: SKILLTEST-016, SKILLTEST-017
- Notes: 既存コードの確立済み規約は根拠として認め、技術スタックをFeature Taskへ推測で埋めることは引き続き禁止する。

## 2026-08-14

- Target Skill: planning-orchestrator, task-breakdown, feature-design, design-review, implementation-handoff schema
- Symptom: 独立レビューのpassだけでTask Breakdownへ進める余地があり、operation資料では利用者向け説明・バリデーション・結果条件・SQLコメント・日本語図の可読性基準が不足していた。
- Root Cause: 独立レビューと利用者承認を別ゲートとして表現しておらず、Operation Documentation Coverage Gateが契約網羅性だけを規定していた。
- Change: 独立レビューpass後の利用者による明示的な詳細設計承認を必須化し、handoff schemaへ承認記録を追加した。operation資料の可読性項目と、design-reviewでの確認項目を追加した。
- Regression Scenario: SKILLTEST-011, SKILLTEST-012, SKILLTEST-013, SKILLTEST-014
- Notes: L3/L4 Decisionの承認規則は変更せず、承認済み設計の確認機会だけを追加した。

- Target Skill: feature-design
- Symptom: operation資料のI/O概要に入力検証・not_found・conflictを再掲し、バリデーション設計と同じ条件を二重に説明していた。
- Root Cause: Operation Documentation Coverage GateがI/Oとバリデーション設計の責務境界を明示していなかった。
- Change: I/Oを正常な入出力と出力意味に限定し、入力不適合・不存在・競合と対応するエラー結果をバリデーション設計へ集約する規則を追加した。
- Regression Scenario: SKILLTEST-013
- Notes: storage failureなど該当する失敗もバリデーション設計に記載し、同じ条件を二重に記載しない。

- Target Skill: feature-design, design-review
- Symptom: 複数のSQL文を実行するoperationでも、シーケンス図がSQLite Storeへの一つの矢印だけで、SQLとの対応を追跡できなかった。
- Root Cause: SQL文単位の連番と、シーケンス図の矢印との対応をOperation Documentation Coverage Gateに定めていなかった。
- Change: 複数SQL文を連番化し、同番号のSQLite Store矢印を文ごとに記載する規則と、design-reviewの照合項目を追加した。
- Regression Scenario: SKILLTEST-013, SKILLTEST-014
- Notes: 一つのSQL文を複数行で記述する場合は同じ番号を継続して使う。

- Target Skill: feature-design, design-review
- Symptom: シーケンス図でnot_found、conflict、storage_errorなど異なる失敗結果を一つの「失敗」分岐にまとめ、戻り値とrollbackの差異を追跡できなかった。
- Root Cause: Operation Documentation Coverage Gateが失敗分岐を結果ごとに分ける規則を定めていなかった。
- Change: 異なる失敗結果を別分岐にし、対応するSQL番号とrollback有無を示す規則とレビュー観点を追加した。
- Regression Scenario: SKILLTEST-013, SKILLTEST-014
- Notes: 同じerror codeでも原因・rollback有無が異なる場合は別分岐とする。

- Target Skill: feature-design, design-review
- Symptom: 複数の書込みSQLが失敗した場合を一つの `storage_error` 分岐にまとめたため、どのSQLで失敗してrollbackしたか追跡できなかった。
- Root Cause: 失敗分岐を異なるerror codeごとに分ける規則だけで、同一error codeに至る複数のSQL失敗を区別する規則が不足していた。
- Change: 同じerror codeでも失敗条件またはSQL番号が異なる場合は別分岐とし、SQL番号とrollback有無を各分岐に記載・監査する規則へ強化した。
- Regression Scenario: SKILLTEST-013, SKILLTEST-014
- Notes: SQL文が一つの場合は、同じ失敗を重複して図示しない。

- Target Skill: feature-design, design-review
- Symptom: シーケンス図のエラー分岐を成功経路の `alt` / `else` に入れ子にしたため、正常処理の流れを追いにくかった。
- Root Cause: 失敗原因を分ける規則はあったが、成功経路から分離する図の構造を規定していなかった。
- Change: エラー系を `break` による早期終了として成功経路の外に置く規則と、レビュー観点を追加した。
- Regression Scenario: SKILLTEST-013
- Notes: 成功パスの選択が必要な場合だけ `alt` を使い、エラーには使わない。

- Target Skill: feature-design, design-review
- Symptom: フローチャートで終了する分岐が左、継続する分岐が右となり、主経路を追いにくかった。また更新系の図でtransaction・rollback・書込み失敗がシーケンス図と一致していなかった。
- Root Cause: フローチャートの分岐配置と、図間の処理遷移照合をOperation Documentation Coverage Gateに定めていなかった。
- Change: 継続分岐を先に左、終了分岐を後に右へ置く規則と、transaction・rollback・書込み失敗の図間整合の確認規則を追加した。
- Regression Scenario: SKILLTEST-013
- Notes: Mermaidの配置は宣言順に依存するため、同一判断ノードのリンク順も監査対象とする。

- Target Skill: feature-design, design-review
- Symptom: operation別の詳細設計だけがあり、利用目的ごとの操作順、各操作を使う理由、組合せの矛盾・不要操作を確認できなかった。
- Root Cause: Operation Documentation Coverage Gateがoperation単体の完全性だけを求め、ユースケースによる横断的な操作列と監査を要求していなかった。
- Change: 複数operationを組み合わせるFeatureでは、目的・利用者・前提・操作列・操作理由・終了条件を持つユースケース資料、全operationへの対応付け、漏れ・重複・責務混同・範囲外混入・不要操作・資料間矛盾の監査を必須化した。design-reviewにも同じ照合観点を追加した。
- Regression Scenario: SKILLTEST-006
- Notes: ユースケース資料は既存operationの使い方を説明するものであり、根拠のない公開CLI契約、設定、保存先、運用仕様を追加してはならない。

- Target Skill: feature-design
- Symptom: ユースケースが単一資料の操作表に留まり、利用者が実際の CLI 入力、応答、応答値から次操作への受渡しを追えなかった。
- Root Cause: Operation Documentation Coverage Gateがユースケースの分割と、operation間のJSONデータフロー例を要求していなかった。
- Change: ユースケースごとの資料分割を原則とし、JSON CLI／APIでは各順序にoperation、入力JSON、成功応答JSON、次操作へ渡す応答値を記載する規則を追加した。
- Regression Scenario: SKILLTEST-006
- Notes: 例は既存の承認済み公開契約を説明するためだけに使い、新しいCLI引数、field、設定、保存先を定めない。

- Target Skill: feature-design
- Symptom: ユースケースを分割する要求を、既存の要約・整合監査資料を置換する指示として解釈し、正本を削除した。
- Root Cause: Operation Documentation Coverage Gateが、既存ユースケース資料を正本として保持し、分割資料を補足として追加する規則を定めていなかった。
- Change: 既存 `design/use-cases.md` の構造を正本として保持し、操作列・データ受渡し例はユースケース別の分割資料へ配置する規則へ訂正した。
- Regression Scenario: SKILLTEST-006
- Notes: 利用者が明示的に正本の置換または削除を依頼した場合だけ、その指示を優先する。

- Target Skill: feature-design
- Symptom: 利用者がユースケース別資料への完全移行と索引資料の削除を指示したにもかかわらず、既存の索引を残す規則を優先した。
- Root Cause: ユースケース資料の配置方針が、利用者による索引の維持・削除の選択を表現できなかった。
- Change: ユースケース別資料には目的・根拠・操作表・終了条件とJSONデータ受渡し例を完結して置き、索引資料の存否は利用者の指示に従う規則へ変更した。
- Regression Scenario: SKILLTEST-006
- Notes: 索引を削除する場合も、Feature設計から各ユースケース資料への入口を残す。

- Date: 2026-08-14
- Target Skill: feature-design, design-review, task-breakdown
- Symptom: JSON出力CLIという表現を、request JSON標準入力が必須であると解釈し、承認済みの名前付きoption入力契約を表現できなかった。
- Root Cause: 形式契約、operation資料、ユースケース、reviewの規則が入力transportをJSONに固定していた。
- Change: JSON出力と入力transportを分離し、承認済みの入力形式（名前付きoptionまたはrequest JSON）を操作別に定義・監査する規則へ変更した。
- Regression Scenario: SKILLTEST-015
- Notes: JSON出力を維持しても入力JSONを推測しない。複数値の入力は承認済みtransportを使い、JSON文字列の埋込みを既定にしない。

## 2026-08-14

- Target Skill: skill-maintainer
- Symptom: 実装リポジトリ直下で`validate_skillset.py`を実行すると、workflowが`documents/.ai/workflow`にあるため失敗した。
- Root Cause: validatorが単一のカレントディレクトリ配下だけをskill/workflow rootとしていた。
- Change: 親リポジトリと`documents/`のどちらを起点にしても、両方のskill rootと正しいworkflow rootを解決するようにした。
- Regression Scenario: SKILLTEST-018
- Notes: documents側のsymlinkと親階層のskill本体は同一skillとして重複排除する。
