# Skill Tests Changelog

## 2026-08-14

- `review-knowledge-cli`に、candidate差分から変更production packageを独立抽出し、temporary profileでcoverageを再測定するreview専用Python scriptと回帰scenarioを追加。
- `impl-knowledge-cli`群のSkill補助scriptをPythonへ統一し、工程同一性scriptはオーケストレーター、実装品質gateは実装Skillへ配置。CI callerも新配置の独立したliteral・coverage gateへ更新し、validatorと実装最終チェック、独立レビューに旧参照・file配置matrix・孤立・重複・owner確認を追加。
- Knowledge CLI実装workflowに、正規sourceのauthority、baseline、candidate ID、構造化packet、finding ID、verdict、同一candidate gateを定める共通契約を追加。
- Issue更新を実装途中からreview・最終整合性監査の両方がPASSした後へ移し、更新後のread-backを必須化。
- Feature Decisionは実装Skillが直接作成せず、`planning-orchestrator`からartifact ownerの`feature-design`へ戻すownershipを明確化。
- 実装固有のlint・coverage・comment・process fixture scenarioをimplementerへ移し、各roleのscenario ownershipを修正。
- skillset validatorへscenario schema・target一致・role coverage・共通契約・candidate fingerprint・metadataの構造検査を追加。
- candidate IDの変更集合を自動列挙し、source ID、Implementation Report verdict、共通finding lifecycleを追加。coverageとliteral layout検査を独立gateへ分離。
- source IDへinventoryの種別とroot pathを含め、Feature directory全体から現存fileだけへの監査範囲縮小を検知する回帰scenarioを追加。
- `impl-knowledge-cli`: 完了報告前に対象Issueの受入条件を実施・検証状況へ更新する手順と回帰scenarioを追加。未検証項目を完了扱いしない規約を明確化。
- `impl-knowledge-cli-implementation`: 利用者確認コマンドの棚卸しを、`testdata/fixtures`の期待I/Oと`test/integration`の実バイナリ検証へ対応付ける手順を追加。Taskfileを期待結果の正本にしない規約と回帰scenarioを追加。
- `impl-knowledge-cli-implementation`: Goおよびmigration SQLの説明コメント、literal整形、table-driven test、package別statement coverage 100%、SQLite SQL配置とplaceholderの規約・回帰scenarioを追加。
- `impl-knowledge-cli-implementation`: 最終candidateで`task lint`を含むfull gateを実行し、Implementation Reportへ記録する規約と回帰scenarioを追加。

## 2026-08-15

- `audit-knowledge-cli-conformance`: 単一command fixtureを、設計済みuse caseの検証証拠として誤認しないよう、response由来のIDを同一隔離Store上の次commandへ渡す実binary操作列を確認する規約と回帰scenarioを追加。#185
