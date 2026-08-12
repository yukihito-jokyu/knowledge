# Task Map

## 目的

Issue #1から分割したタスクの現在構造、状態、成果物、依存関係、次の議論位置を一覧できるようにする。

この文書は現在のスナップショットだけを示す。提案の変遷、分割理由、承認理由は[タスク分割の議論・承認記録](task-decomposition-decisions.md)に残す。

## Source of Truth

| 情報 | 参照先 |
| --- | --- |
| システム要件・設計原則 | [Issue #1](https://github.com/yukihito-jokyu/knowledge/issues/1) |
| 分割案・変更経緯・承認記録 | [task-decomposition-decisions.md](task-decomposition-decisions.md) |
| 現在のタスク構造・進行位置 | 本文書 |
| leafの後続接続逆引き | [task-connections.md](task-connections.md)（本文書から生成する派生表示） |

## ID規則

- `L<n>`: 大分類
- `L<n>-M<n>`: 中分類
- `L<n>-M<n>-S<n>`: 小分類
- 承認済みIDは原則として維持する。
- 議論中のタスクは承認時に名称・IDが変わる可能性がある。

## 状態

タスク構造の確定度と、実作業の進捗を分離して管理する。

### 計画状態

| 状態 | 意味 |
| --- | --- |
| 承認済み | タスクの存在、名称、親子関係が承認されている |
| 議論中 | 提案・監査済みだが、まだ承認されていない |
| 再検討中 | 承認済みだが、後続分割との矛盾に対する修正案があり再承認を待っている |
| 未分割 | 親タスクのみ承認済みで、子タスクをまだ議論していない |

### 実行状態

| 状態 | 意味 |
| --- | --- |
| 未着手 | 成果物の作成を開始していない |
| 進行中 | 成果物を作成中 |
| 完了 | 完了条件を満たした |
| ブロック | 依存関係や外部判断により進行できない |

## 現在位置

- 現在の議論: `L1-M1-S1`の設計成果物を監査し、後続タスクへ引き継げる状態にする
- 現在の状態: `L1-M1-S1`の論理スキーマを作成・再監査済みで、利用者によるCommit承認を待っている
- 直前の完了: 8状態化、Scopeの結合規則、Evidenceの命題評価と利用者自己申告の分離を設計へ反映
- 次の議論: 利用者の承認後に`L1-M1-S1`をCommitし、`L1-M1-S2`と`L1-M1-S3`へ引き継ぐ
- GitHub Issue対応表: [`docs/github-issue-map.md`](github-issue-map.md)

## 階層マップ

```text
L1 Knowledge論理モデル・CLI公開契約を設計する [承認済み]
├─ L1-M1 Knowledge論理データモデル・整合性制約を設計する [承認済み]
│  ├─ L1-M1-S1 Knowledge論理スキーマを定義する [承認済み]
│  ├─ L1-M1-S2 Knowledge論理関連と参照整合性制約を定義する [承認済み]
│  ├─ L1-M1-S3 Evidence由来Stateの整合性・導出契約を定義する [承認済み]
│  ├─ L1-M1-S4 Knowledge更新操作の論理事前・事後条件と履歴系譜を定義する [承認済み]
│  └─ L1-M1-S5 論理モデルの要件トレーサビリティを確認する [承認済み]
└─ L1-M2 CLI公開契約・JSON Schemaを設計する [承認済み]
   ├─ L1-M2-S1 CLI公開コマンド一覧と動作・呼び出し・引数契約を定義する [承認済み]
   ├─ L1-M2-S2 検索・複数取得結果のコレクション契約を定義する [承認済み]
   ├─ L1-M2-S3 CLIエラー分類と終了コード契約を定義する [承認済み]
   ├─ L1-M2-S4 CLIのJSON入出力Schemaを定義する [承認済み]
   ├─ L1-M2-S5 CLI契約のVersioning・後方互換性規則を定義する [承認済み]
   └─ L1-M2-S6 CLI公開契約の要件トレーサビリティと責務境界遵守を確認する [承認済み]

L2 Knowledge Store／CLI基盤 [承認済み]
├─ L2-M1 Knowledge Storeの物理永続化・更新基盤を設計・実装する [承認済み]
│  ├─ L2-M1-S1 正本永続化の技術スタックを選定しADRを確定する [承認済み]
│  ├─ L2-M1-S2 正本データの物理モデル・履歴・Schema Evolutionを詳細設計する [承認済み]
│  ├─ L2-M1-S3 永続化操作・Transaction・Repository境界を詳細設計する [承認済み]
│  ├─ L2-M1-S4 Baseline物理SchemaとDB制約をMigrationとして実装する [承認済み]
│  ├─ L2-M1-S5 Schema Version管理とMigration実行基盤を実装する [承認済み]
│  ├─ L2-M1-S6 DB接続・Transaction・ID・Record変換Coreを実装する [承認済み]
│  ├─ L2-M1-S7 Evidence由来Stateの導出・再評価経路を実装する [承認済み]
│  ├─ L2-M1-S8 create・attach-evidenceの追加型更新を実装する [承認済み]
│  ├─ L2-M1-S9 revise・supersedeと非破壊履歴系譜を実装する [承認済み]
│  └─ L2-M1-S10 Knowledge詳細・Evidence・更新履歴の取得Repositoryを実装する [承認済み]
├─ L2-M2 複合検索・Index管理基盤を設計・実装する [承認済み]
│  ├─ L2-M2-S1 候補検索・Index Lifecycle要件と正本設計Review Gate基準を定義する [承認済み]
│  ├─ L2-M2-S2 検索・Index技術スタックを選定しADRを確定する [承認済み]
│  ├─ L2-M2-S3 候補検索アーキテクチャと方式別Indexを詳細設計する [承認済み]
│  ├─ L2-M2-S4 正本同期・Index Version・再構築／障害回復を詳細設計する [承認済み]
│  ├─ L2-M2-S5 共通候補検索Coreと正本照合を実装する [承認済み]
│  ├─ L2-M2-S6 Lexical Indexと文字列候補検索を実装する [承認済み]
│  ├─ L2-M2-S7 Embedding生成・Semantic Indexと意味候補検索を実装する [承認済み]
│  ├─ L2-M2-S8 Concept／Relation Indexと構造・Contradiction候補検索を実装する [承認済み]
│  ├─ L2-M2-S9 Temporal候補検索を実装する [承認済み]
│  └─ L2-M2-S10 正本同期・Index Version管理・再構築／障害回復を実装する [承認済み]
└─ L2-M3 Knowledge CLI実行基盤とJSON境界を詳細設計・実装する [承認済み]
   ├─ L2-M3-S1 CLI Framework・JSON検証・Build／配布技術を選定しADRを確定する [承認済み]
   ├─ L2-M3-S2 CLI内部ArchitectureとCommand-to-Port Mappingを詳細設計する [承認済み]
   ├─ L2-M3-S3 CLI Process CoreとCommand Dispatchを実装する [承認済み]
   ├─ L2-M3-S4 JSON Validation・Serialization・Error／終了コード境界を実装する [承認済み]
   ├─ L2-M3-S5 取得・候補検索Command Handler／Adapterを実装する [承認済み]
   ├─ L2-M3-S6 Knowledge更新Command Handler／Adapterを実装する [承認済み]
   ├─ L2-M3-S7 Index管理・再構築／障害回復CommandをCLIへ接続する [承認済み]
   ├─ L2-M3-S8 設定・初期化・Composition Rootと配布Artifactを実装する [承認済み]
   └─ L2-M3-S9 実行ArtifactのL1公開契約適合をBlack-box検証する [承認済み]
L3 Knowledge蓄積Skillsを設計・実装する [承認済み]
├─ L3-M1 Knowledge Acquisition Skillを詳細設計・実装する [承認済み]
│  ├─ L3-M1-S1 Episode入力・Source Reference・Evidence候補採否を詳細設計する [承認済み]
│  ├─ L3-M1-S2 Knowledge Candidate正規化・L1 Field写像を詳細設計する [承認済み]
│  ├─ L3-M1-S3 Candidate Markdown／no-candidate・入力不足引渡し契約を定義する [承認済み]
│  └─ L3-M1-S4 Knowledge Acquisition Skill instructions／referencesを実装し、固有構造契約を確認する [承認済み]
└─ L3-M2 Knowledge Update Skillを詳細設計・実装する [承認済み]
   ├─ L3-M2-S1 Candidate入力受入・評価可否境界を詳細設計する [承認済み]
   ├─ L3-M2-S2 Existing Knowledge探索・意味比較手順を詳細設計する [承認済み]
   ├─ L3-M2-S3 4更新操作選択・非永続化判断とDecision Markdown契約を定義する [承認済み]
   ├─ L3-M2-S4 CLI Command Mapping・実行結果／失敗引渡しを詳細設計する [承認済み]
   ├─ L3-M2-S5 Knowledge Update Skill instructions／referencesをMock CLI契約に基づき実装し、固有構造契約を確認する [承認済み]
   └─ L3-M2-S6 Knowledge Update Skillを実CLI Artifactへ接続し、固有Command連携を確認する [承認済み]
L4 記事価値判定Skillsを設計・実装する [承認済み]
├─ L4-M1 Article Analysis Skillを詳細設計・実装する [承認済み]
│  ├─ L4-M1-S1 記事入力・取得方式・解析可否境界を詳細設計する [承認済み]
│  ├─ L4-M1-S2 Article overview・Claim分解／正規化手順を詳細設計する [承認済み]
│  ├─ L4-M1-S3 Claim Location・記事内Support根拠の追跡方法を詳細設計する [承認済み]
│  ├─ L4-M1-S4 Article Analysis Markdown・局所再分析結果契約を定義する [承認済み]
│  └─ L4-M1-S5 Article Analysis Skill instructions／referencesを実装し、固有構造契約を確認する [承認済み]
├─ L4-M2 Knowledge Search Skillを詳細設計・実装する [承認済み]
│  ├─ L4-M2-S1 技術非依存な探索要求・Query Journey・局所停止理由を詳細設計する [承認済み]
│  ├─ L4-M2-S2 Article Claim受入・Target Claim分解／検索variant生成手順を詳細設計する [承認済み]
│  ├─ L4-M2-S3 CLI検索・取得Command Mappingと部分失敗引渡しを詳細設計する [承認済み]
│  ├─ L4-M2-S4 Evidence意味比較・探索十分性／Knowledge Assessment判定手順を詳細設計する [承認済み]
│  ├─ L4-M2-S5 Knowledge Assessment・raw Search Trace・局所再検索結果契約を定義する [承認済み]
│  ├─ L4-M2-S6 Knowledge Search SkillをMock CLI契約に基づき実装し、固有構造契約を確認する [承認済み]
│  └─ L4-M2-S7 Knowledge Search Skillを実CLI Artifactへ接続し、固有Command連携を確認する [承認済み]
└─ L4-M3 Reading Value Skillを詳細設計・実装する [承認済み]
   ├─ L4-M3-S1 Article Analysis／Knowledge Search入力整合・評価可否境界を詳細設計する [承認済み]
   ├─ L4-M3-S2 Knowledge AssessmentからRecognition Gain Candidateへの適用手順を詳細設計する [承認済み]
   ├─ L4-M3-S3 Claim・記事内SupportからReliability／Applicabilityを判断する手順を詳細設計する [承認済み]
   ├─ L4-M3-S4 Attention Cost・Claim横断統合によるRecommendation／読書範囲決定手順を詳細設計する [承認済み]
   ├─ L4-M3-S5 追加調査要求・Reading Value Assessment Markdown契約を定義する [承認済み]
   └─ L4-M3-S6 Reading Value Skill instructions／referencesを実装し、固有構造契約を確認する [承認済み]
L5 Parent Orchestration Skillと全体Workflowを設計・実装する [承認済み]
├─ L5-M1 Parent entrypoint・共通Orchestration状態／実行制御契約を詳細設計する [承認済み]
│  ├─ L5-M1-S1 Parent entrypoint・Workflow識別／起動受付契約を定義する [承認済み]
│  ├─ L5-M1-S2 子Skill Registry・起動／成果物引渡し契約を定義する [承認済み]
│  ├─ L5-M1-S3 共通実行状態・診断Envelope契約を定義する [承認済み]
│  ├─ L5-M1-S4 Run／Claim／Candidate／Cycle・Search Trace相関契約を定義する [承認済み]
│  └─ L5-M1-S5 Budget・retry／restart／reroute／stop・副作用安全性契約を定義する [承認済み]
├─ L5-M2 記事価値判定Workflowの起動・Claim fan-out／結果統合・再調査・終了を詳細設計する [承認済み]
│  ├─ L5-M2-S1 記事価値判定Workflowを開始し、Article Analysisを起動・初期成果物を受け入れる [承認済み]
│  ├─ L5-M2-S2 canonical Claim単位のKnowledge Search fan-out契約を定義する [承認済み]
│  ├─ L5-M2-S3 Claim別Knowledge Search結果を相関・集約し、Reading Valueへ引き渡す [承認済み]
│  ├─ L5-M2-S4 追加再分析／再検索要求のRouting・再調査Cycle反映を定義する [承認済み]
│  └─ L5-M2-S5 最終Assessment受入・異常／Budget到達時のWorkflow終端契約を定義する [承認済み]
├─ L5-M3 Knowledge蓄積WorkflowのEpisode起動・Candidate分岐・更新終了を詳細設計する [承認済み]
│  ├─ L5-M3-S1 Knowledge蓄積Workflowを開始し、Knowledge Acquisitionを起動・成果物を受け入れる [承認済み]
│  ├─ L5-M3-S2 candidate／no-candidate／input_insufficient分岐を定義する [承認済み]
│  ├─ L5-M3-S3 Candidate更新Work Itemを編成し、Knowledge Updateを起動する [承認済み]
│  ├─ L5-M3-S4 Candidate更新結果を集約し、部分成功後の安全な再開順序を定義する [承認済み]
│  └─ L5-M3-S5 正常／異常／Budget到達時のKnowledge蓄積Workflow終端契約を定義する [承認済み]
└─ L5-M4 Parent Orchestration Skillを実装し、Mock契約確認後にL3／L4実Componentを統合する [承認済み]
   ├─ L5-M4-S1 Parent entrypoint・Registry・共通実行制御を実装する [承認済み]
   ├─ L5-M4-S2 記事価値判定WorkflowのParent instructionsを実装する [承認済み]
   ├─ L5-M4-S3 Knowledge蓄積WorkflowのParent instructionsを実装する [承認済み]
   ├─ L5-M4-S4 Mock ComponentsでParent固有Orchestration契約を確認する [承認済み]
   ├─ L5-M4-S5 Parent Orchestration SkillをL3実Componentsへ接続し、Knowledge蓄積連携を確認する [承認済み]
   └─ L5-M4-S6 Parent Orchestration SkillをL4実Componentsへ接続し、記事価値判定連携を確認する [承認済み]
L6 評価設計・横断品質保証を実装する [承認済み]
├─ L6-M1 横断評価方針・Coverage／判定基準・Search Trace診断／Report契約を詳細設計する [承認済み]
│  ├─ L6-M1-S1 評価レイヤ・対象範囲・L1〜L6テスト所有境界を定義する [承認済み]
│  ├─ L6-M1-S2 Scenario A〜J・Invariant・Done DefinitionのCoverage／検証責務Matrixを定義する [承認済み]
│  ├─ L6-M1-S3 評価指標・合否／評価不能・反復判定基準を定義する [承認済み]
│  ├─ L6-M1-S4 既存Search Trace・相関情報を用いた失敗原因診断手順を定義する [承認済み]
│  └─ L6-M1-S5 評価結果Envelope・診断／集約Report契約を定義する [承認済み]
├─ L6-M2 Scenario A〜J・Invariantを検証する共有Dataset／Fixtureを設計・実装する [承認済み]
│  ├─ L6-M2-S1 共有Dataset／Fixture Schema・Version・Provenance契約を設計する [承認済み]
│  ├─ L6-M2-S2 Scenario A〜J・Invariantの具体評価Case Catalog・期待Oracleを設計する [承認済み]
│  ├─ L6-M2-S3 共有Dataset Manifest・Fixture静的検証基盤を実装する [承認済み]
│  ├─ L6-M2-S4 Scenario A〜E・検索／Knowledge State Invariantの共有Fixtureを実装する [承認済み]
│  ├─ L6-M2-S5 Scenario F〜H・Knowledge Acquisition／Update Invariantの共有Fixtureを実装する [承認済み]
│  └─ L6-M2-S6 Scenario I〜J・Reading Value Invariantの共有Fixtureを実装する [承認済み]
├─ L6-M3 CLI・Skill・Workflow共通評価Harnessを設計・実装する [承認済み]
│  ├─ L6-M3-S1 共通評価Harnessの実行・再現性・隔離要件を定義する [承認済み]
│  ├─ L6-M3-S2 評価Harnessの技術スタック・Runner・依存管理方式を選定しADRを確定する [承認済み]
│  ├─ L6-M3-S3 Dataset Loader・Target Adapter・Oracle・Trace Collector・Reporter Architectureを詳細設計する [承認済み]
│  ├─ L6-M3-S4 評価Harness Core・Dataset Loader・Oracle／Report Pipelineを実装する [承認済み]
│  ├─ L6-M3-S5 検索・取得・更新CLI Target Adapterを実装する [承認済み]
│  ├─ L6-M3-S6 L3／L4 Component Skill Target Adapterを実装する [承認済み]
│  └─ L6-M3-S7 Parent Workflow Target Adapterを実装する [承認済み]
├─ L6-M4 CLI横断品質・L3／L4各Component Skillのレイヤ別／Agent評価を実行し、診断Reportを生成する [承認済み]
│  ├─ L6-M4-S1 共有FixtureでCLI検索・取得・更新の横断品質を評価する [承認済み]
│  ├─ L6-M4-S2 Knowledge Acquisition SkillのAgent評価を実行する [承認済み]
│  ├─ L6-M4-S3 Knowledge Update SkillのAgent評価を実行する [承認済み]
│  ├─ L6-M4-S4 Article Analysis SkillのAgent評価を実行する [承認済み]
│  ├─ L6-M4-S5 Knowledge Search SkillのAgentic Search評価を実行する [承認済み]
│  └─ L6-M4-S6 Reading Value SkillのAgent評価を実行する [承認済み]
└─ L6-M5 記事価値判定・Knowledge蓄積WorkflowのAcceptance／E2E・回帰評価を実行する [承認済み]
   ├─ L6-M5-S1 Knowledge蓄積WorkflowのAcceptance／E2E評価を実行する [承認済み]
   ├─ L6-M5-S2 記事価値判定WorkflowのAcceptance／E2E評価を実行する [承認済み]
   └─ L6-M5-S3 全Scenario回帰Suiteを実行し、最終Acceptance／診断Reportを確定する [承認済み]
```

## 確定した横断境界

詳細な変更理由と旧構造は、議論記録の「決定009」を参照する。

| 対象 | 確定した境界 | 主な影響 |
| --- | --- | --- |
| L1 | Knowledge論理モデルとCLI公開契約に限定 | L2〜L5の詳細設計との重複を解消 |
| L3〜L5 | L3・L4は個別Skills、L5はParent Orchestrationと全体Workflowを所有 | Workflow制御の二重所有を解消 |
| L2-M1〜M3 | 正本、派生Index、CLI Process境界を分離し、検索要求Review Gateを設ける | 設計入力、同期・再構築、CLI内部設計を明確化 |
| L6・各実装 | 各実装は固有テスト、L6は共有評価資産と横断評価を所有 | テストの二重所有を解消 |

## タスク台帳

| ID | タスク名 | 計画状態 | 実行状態 | 主成果物・到達状態 | 直接依存 |
| --- | --- | --- | --- | --- | --- |
| L1 | Knowledge論理モデル・CLI公開契約を設計する | 承認済み | 進行中 | Knowledge論理モデルとCLI公開契約・JSON Schema | なし |
| L1-M1 | Knowledge論理データモデル・整合性制約を設計する | 承認済み | 進行中 | 実装可能で検索要件を満たす論理モデル | なし |
| L1-M1-S1 | Knowledge論理スキーマを定義する | 承認済み | 進行中 | 全論理レコードのデータ辞書 | なし |
| L1-M1-S2 | Knowledge論理関連と参照整合性制約を定義する | 承認済み | 未着手 | 関連種別・方向・多重度・参照整合性の制約表 | L1-M1-S1 |
| L1-M1-S3 | Evidence由来Stateの整合性・導出契約を定義する | 承認済み | 未着手 | Evidence、導出結果、根拠追跡、再計算・無効化の契約 | L1-M1-S1 |
| L1-M1-S4 | Knowledge更新操作の論理事前・事後条件と履歴系譜を定義する | 承認済み | 未着手 | 更新操作の遷移表と履歴系譜規則 | L1-M1-S2、L1-M1-S3 |
| L1-M1-S5 | 論理モデルの要件トレーサビリティを確認する | 承認済み | 未着手 | Issue要件と設計要素の双方向トレーサビリティ表 | L1-M1-S1〜S4 |
| L1-M2 | CLI公開契約・JSON Schemaを設計する | 承認済み | 未着手 | コマンド、入出力Schema、エラー、互換性契約 | L1-M1 |
| L1-M2-S1 | CLI公開コマンド一覧と動作・呼び出し・引数契約を定義する | 承認済み | 未着手 | 公開コマンド台帳 | L1-M1-S5 |
| L1-M2-S2 | 検索・複数取得結果のコレクション契約を定義する | 承認済み | 未着手 | 集合取得規則とコマンド別適用表 | L1-M2-S1 |
| L1-M2-S3 | CLIエラー分類と終了コード契約を定義する | 承認済み | 未着手 | Error台帳と終了コード対応表 | L1-M2-S1 |
| L1-M2-S4 | CLIのJSON入出力Schemaを定義する | 承認済み | 未着手 | 公開操作のJSON Schema一式 | L1-M2-S1〜S3 |
| L1-M2-S5 | CLI契約のVersioning・後方互換性規則を定義する | 承認済み | 未着手 | 契約Versioning・互換性方針 | L1-M2-S4 |
| L1-M2-S6 | CLI公開契約の要件トレーサビリティと責務境界遵守を確認する | 承認済み | 未着手 | Issue要件とCLI契約の双方向対応表 | L1-M2-S1〜S5 |
| L2 | Knowledge Store／CLI基盤 | 承認済み | 未着手 | 決定論的な保存・検索・更新基盤 | L1 |
| L2-M1 | Knowledge Storeの物理永続化・更新基盤を設計・実装する | 承認済み | 未着手 | Migration可能で履歴・整合性を保持するKnowledge Store | 開始入力: L1、L2-M2-S1 |
| L2-M1-S1 | 正本永続化の技術スタックを選定しADRを確定する | 承認済み | 未着手 | L2-M2-S1のReview Gate基準を満たす永続化技術ADR | L1-M1、L2-M2-S1 |
| L2-M1-S2 | 正本データの物理モデル・履歴・Schema Evolutionを詳細設計する | 承認済み | 未着手 | L2-M2-S1の検索・再構築要件を満たす物理ストレージ詳細設計書 | L2-M1-S1 |
| L2-M1-S3 | 永続化操作・Transaction・Repository境界を詳細設計する | 承認済み | 未着手 | Commit済み正本変更境界を含む実行時永続化設計書 | L2-M1-S2、L1-M2-S1、L1-M2-S3 |
| L2-M1-S4 | Baseline物理SchemaとDB制約をMigrationとして実装する | 承認済み | 未着手 | 初回DDL／Migration | L2-M1-S2、L2 Design Freeze Gate |
| L2-M1-S5 | Schema Version管理とMigration実行基盤を実装する | 承認済み | 未着手 | Version管理されたMigration Runner | L2-M1-S4 |
| L2-M1-S6 | DB接続・Transaction・ID・Record変換Coreを実装する | 承認済み | 未着手 | 取得・更新が共有する低水準永続化層 | L2-M1-S3、L2-M1-S5 |
| L2-M1-S7 | Evidence由来Stateの導出・再評価経路を実装する | 承認済み | 未着手 | EvidenceとStateの実行時整合経路 | L2-M1-S6 |
| L2-M1-S8 | create・attach-evidenceの追加型更新を実装する | 承認済み | 未着手 | 原子的な追加型更新Repository | L2-M1-S6、L2-M1-S7 |
| L2-M1-S9 | revise・supersedeと非破壊履歴系譜を実装する | 承認済み | 未着手 | 非破壊更新Repositoryと履歴系譜 | L2-M1-S6、L2-M1-S7 |
| L2-M1-S10 | Knowledge詳細・Evidence・更新履歴の取得Repositoryを実装する | 承認済み | 未着手 | 論理Aggregateの復元、Evidence・履歴追跡、正本Snapshot列挙・Commit済み差分取得Port | L2-M1-S6、L2-M1-S7 |
| L2-M2 | 複合検索・Index管理基盤を設計・実装する | 承認済み | 未着手 | Issue 12章の全検索プリミティブを支える候補検索と、正本から再構築可能な派生Index | 開始入力: L1、L4-M2-S1。完了・Merge: L2-M1 |
| L2-M2-S1 | 候補検索・Index Lifecycle要件と正本設計Review Gate基準を定義する | 承認済み | 未着手 | 技術非依存の検索・鮮度・整合・再構築要件とReview Gate基準 | L1-M1、L1-M2、L3-M2-S2、L4-M2-S1 |
| L2-M2-S2 | 検索・Index技術スタックを選定しADRを確定する | 承認済み | 未着手 | 検索・Embedding・同期方式の技術選定ADR | L2-M2-S1、L2-M1-S1〜S3 |
| L2-M2-S3 | 候補検索アーキテクチャと方式別Indexを詳細設計する | 承認済み | 未着手 | 内部Component、Index物理表現、共通Candidateの詳細設計 | L2-M2-S2 |
| L2-M2-S4 | 正本同期・Index Version・再構築／障害回復を詳細設計する | 承認済み | 未着手 | Commit境界から世代切替・復旧までのIndex Lifecycle設計 | L2-M2-S3、L2-M1-S3 |
| L2-M2-S5 | 共通候補検索Coreと正本照合を実装する | 承認済み | 未着手 | Provider共通契約、Filter・順序・Pagination、正本照合Core | 着手: L2-M2-S3、L2 Design Freeze Gate。完了・Merge: L2-M1-S10 |
| L2-M2-S6 | Lexical Indexと文字列候補検索を実装する | 承認済み | 未着手 | Lexical Indexと文字列候補検索Provider | L2-M2-S4、L2-M2-S5 |
| L2-M2-S7 | Embedding生成・Semantic Indexと意味候補検索を実装する | 承認済み | 未着手 | Version追跡可能なEmbedding・Semantic検索Provider | L2-M2-S4、L2-M2-S5 |
| L2-M2-S8 | Concept／Relation Indexと構造・Contradiction候補検索を実装する | 承認済み | 未着手 | Relation系構造Indexと候補検索Provider | L2-M2-S4、L2-M2-S5 |
| L2-M2-S9 | Temporal候補検索を実装する | 承認済み | 未着手 | 独立物理Indexを必須としないTemporal検索Provider | L2-M2-S4、L2-M2-S5 |
| L2-M2-S10 | 正本同期・Index Version管理・再構築／障害回復を実装する | 承認済み | 未着手 | 全検索Providerを統合するIndex Lifecycle実装 | 着手: L2-M2-S4、L2 Design Freeze Gate。完了・Merge: L2-M2-S6〜S9、L2-M1-S8〜S10 |
| L2-M3 | Knowledge CLI実行基盤とJSON境界を詳細設計・実装する | 承認済み | 未着手 | 内部CLI設計とL1公開契約に適合する実行可能CLI | 着手: L1-M2。正本・検索Adapter本番実装: L2 Design Freeze Gate。完了・Merge: L2-M1、L2-M2 |
| L2-M3-S1 | CLI Framework・JSON検証・Build／配布技術を選定しADRを確定する | 承認済み | 未着手 | CLI Runtime、Schema Validator、Build・配布方式の技術選定ADR | 着手: L1-M2。ADR確定: L2-M1-S1、L2-M2-S2 |
| L2-M3-S2 | CLI内部ArchitectureとCommand-to-Port Mappingを詳細設計する | 承認済み | 未着手 | 全公開CommandのHandler・Port所有表と内部CLI詳細設計 | L2-M3-S1、L2-M1-S3、L2-M2-S3・S4 |
| L2-M3-S3 | CLI Process CoreとCommand Dispatchを実装する | 承認済み | 未着手 | Process Lifecycle、Command Interface、Dispatch、help・version入口 | L2-M3-S2、L2 Design Freeze Gate |
| L2-M3-S4 | JSON Validation・Serialization・Error／終了コード境界を実装する | 承認済み | 未着手 | JSON入出力、stdout・stderr、内部Error変換、終了コード境界 | L2-M3-S3 |
| L2-M3-S5 | 取得・候補検索Command Handler／Adapterを実装する | 承認済み | 未着手 | 全取得・検索Commandの読取りAdapter | 着手: L2-M3-S4。完了・Merge: L2-M1-S10、L2-M2-S5〜S9 |
| L2-M3-S6 | Knowledge更新Command Handler／Adapterを実装する | 承認済み | 未着手 | create・attach-evidence・revise・supersedeの更新Adapter | 着手: L2-M3-S4。完了・Merge: L2-M1-S8・S9 |
| L2-M3-S7 | Index管理・再構築／障害回復CommandをCLIへ接続する | 承認済み | 未着手 | Index Lifecycle操作のCLI Adapter | 着手: L2-M3-S4。完了・Merge: L2-M2-S10 |
| L2-M3-S8 | 設定・初期化・Composition Rootと配布Artifactを実装する | 承認済み | 未着手 | 全Adapterを配線したVersion付き実行Artifact | L2-M3-S3〜S7。完了・Merge: L2-M1、L2-M2 |
| L2-M3-S9 | 実行ArtifactのL1公開契約適合をBlack-box検証する | 承認済み | 未着手 | 配布Artifactに対するCLI単体Contract Suiteと適合結果 | L2-M3-S8、L1-M2-S6 |
| L3 | Knowledge蓄積Skillsを設計・実装する | 承認済み | 未着手 | Knowledge Acquisition・Knowledge Update Skills | 着手: L1。完了・Release: L2-M3-S9 |
| L3-M1 | Knowledge Acquisition Skillを詳細設計・実装する | 承認済み | 未着手 | Acquisition Skill、Candidate Markdown契約、no-candidate・入力不足表現 | 着手: L1-M1。本番実装: L3 Skill Design Freeze Gate。L2依存なし |
| L3-M1-S1 | Episode入力・Source Reference・Evidence候補採否を詳細設計する | 承認済み | 未着手 | Episode入力・Evidence採否境界設計書 | L1-M1 |
| L3-M1-S2 | Knowledge Candidate正規化・L1 Field写像を詳細設計する | 承認済み | 未着手 | Candidate正規化・L1 Field写像仕様 | 着手: L1-M1。確定: L3-M1-S1 |
| L3-M1-S3 | Candidate Markdown／no-candidate・入力不足引渡し契約を定義する | 承認済み | 未着手 | Candidate Markdown契約とcanonical例 | L3-M1-S1・S2 |
| L3-M1-S4 | Knowledge Acquisition Skill instructions／referencesを実装し、固有構造契約を確認する | 承認済み | 未着手 | 固有構造契約に適合するKnowledge Acquisition Skill package | L3 Skill Design Freeze Gate |
| L3-M2 | Knowledge Update Skillを詳細設計・実装する | 承認済み | 未着手 | Update Skill、Candidate受入規則、操作決定Markdown、CLI Command Mapping | 設計着手: L1-M1。契約確定: L3-M1 Candidate契約、L1-M2。本番実装: L3 Gate。実CLI統合: L2-M3-S5・S6・S8 |
| L3-M2-S1 | Candidate入力受入・評価可否境界を詳細設計する | 承認済み | 未着手 | Candidate受入・評価可否境界設計書 | L3-M1-S3、L1 |
| L3-M2-S2 | Existing Knowledge探索・意味比較手順を詳細設計する | 承認済み | 未着手 | 更新対象探索・意味比較設計書 | L3-M1-S3、L1。出力をL2-M2-S1へ入力する |
| L3-M2-S3 | 4更新操作選択・非永続化判断とDecision Markdown契約を定義する | 承認済み | 未着手 | Update Decision Markdown契約 | L3-M2-S1・S2、L1-M1-S4 |
| L3-M2-S4 | CLI Command Mapping・実行結果／失敗引渡しを詳細設計する | 承認済み | 未着手 | CLI利用・結果引渡し契約 | 着手: L1-M2。確定: L3-M2-S2・S3 |
| L3-M2-S5 | Knowledge Update Skill instructions／referencesをMock CLI契約に基づき実装し、固有構造契約を確認する | 承認済み | 未着手 | Mockで固有構造契約を確認済みのKnowledge Update Skill package | L3 Skill Design Freeze Gate |
| L3-M2-S6 | Knowledge Update Skillを実CLI Artifactへ接続し、固有Command連携を確認する | 承認済み | 未着手 | M2固有の実CLI Component連携 | L3-M2-S5、L2-M3-S5・S6・S8。L3 ReleaseはL2-M3-S9後 |
| L4 | 記事価値判定Skillsを設計・実装する | 承認済み | 未着手 | Article Analysis・Knowledge Search・Reading Value Skills | 設計着手: L1。本番Skill実装: L4 Gate。L4-M2実CLI統合: L2-M3-S5・S8。完了・Release: L4-M1〜M3、L2-M3-S9 |
| L4-M1 | Article Analysis Skillを詳細設計・実装する | 承認済み | 未着手 | Article Analysis Skill、Claims／Support Markdown契約 | 設計着手: L1。本番実装: L4 Skill Design Freeze Gate |
| L4-M1-S1 | 記事入力・取得方式・解析可否境界を詳細設計する | 承認済み | 未着手 | 記事入力・取得方式・Source同一性・解析可否境界設計 | Issue、L1 |
| L4-M1-S2 | Article overview・Claim分解／正規化手順を詳細設計する | 承認済み | 未着手 | Overview・Claim分解／正規化設計 | Issue、L1 |
| L4-M1-S3 | Claim Location・記事内Support根拠の追跡方法を詳細設計する | 承認済み | 未着手 | Location・Support根拠追跡設計 | L4-M1-S1・S2 |
| L4-M1-S4 | Article Analysis Markdown・局所再分析結果契約を定義する | 承認済み | 未着手 | Article Analysis Markdown・局所再分析結果契約 | L4-M1-S1〜S3 |
| L4-M1-S5 | Article Analysis Skill instructions／referencesを実装し、固有構造契約を確認する | 承認済み | 未着手 | 実行可能なArticle Analysis Skill package | L4 Skill Design Freeze Gate |
| L4-M2 | Knowledge Search Skillを詳細設計・実装する | 承認済み | 未着手 | Knowledge Search Skill、Assessment／raw Search Trace契約、CLI連携 | 設計着手: L1。実CLI統合: L2-M3-S5・S8。L4 Release: L2-M3-S9 |
| L4-M2-S1 | 技術非依存な探索要求・Query Journey・局所停止理由を詳細設計する | 承認済み | 未着手 | 技術非依存な探索能力・遷移・必要取得情報・局所停止理由 | Issue、L1-M1、L1-M2 |
| L4-M2-S2 | Article Claim受入・Target Claim分解／検索variant生成手順を詳細設計する | 承認済み | 未着手 | Claim受入・検索用部分Claim／variant生成手順 | L4-M2-S1、L4-M1-S4 |
| L4-M2-S3 | CLI検索・取得Command Mappingと部分失敗引渡しを詳細設計する | 承認済み | 未着手 | Knowledge Search固有CLI Mapping・失敗引渡し契約 | L4-M2-S1・S2、L1-M2 |
| L4-M2-S4 | Evidence意味比較・探索十分性／Knowledge Assessment判定手順を詳細設計する | 承認済み | 未着手 | Evidenceから8状態・Known・Gap・Confidenceを導く手順 | L4-M2-S1・S2、L1-M1 |
| L4-M2-S5 | Knowledge Assessment・raw Search Trace・局所再検索結果契約を定義する | 承認済み | 未着手 | Assessment・Trace・実行状態を分離した出力契約 | L4-M2-S1〜S4 |
| L4-M2-S6 | Knowledge Search SkillをMock CLI契約に基づき実装し、固有構造契約を確認する | 承認済み | 未着手 | Mockで固有構造契約を確認済みのKnowledge Search Skill package | L4 Skill Design Freeze Gate |
| L4-M2-S7 | Knowledge Search Skillを実CLI Artifactへ接続し、固有Command連携を確認する | 承認済み | 未着手 | Knowledge Search固有の実CLI Component連携 | L4-M2-S6、L2-M3-S5・S8。L4 ReleaseはL2-M3-S9後 |
| L4-M3 | Reading Value Skillを詳細設計・実装する | 承認済み | 未着手 | Reading Value Skill、最終Assessment Markdown契約 | 設計着手: L1。契約確定: L4-M1・M2。本番実装: L4 Gate |
| L4-M3-S1 | Article Analysis／Knowledge Search入力整合・評価可否境界を詳細設計する | 承認済み | 未着手 | M1・M2入力のClaim整合・評価可否境界 | L4-M1-S4、L4-M2-S5 |
| L4-M3-S2 | Knowledge AssessmentからRecognition Gain Candidateへの適用手順を詳細設計する | 承認済み | 未着手 | 8状態を再判定しないRecognition Gain適用手順 | L4-M3-S1 |
| L4-M3-S3 | Claim・記事内SupportからReliability／Applicabilityを判断する手順を詳細設計する | 承認済み | 未着手 | Supportに追跡可能なReliability・Applicability判断手順 | L4-M3-S1 |
| L4-M3-S4 | Attention Cost・Claim横断統合によるRecommendation／読書範囲決定手順を詳細設計する | 承認済み | 未着手 | 3推奨と読書範囲のClaim横断統合手順 | L4-M3-S2・S3 |
| L4-M3-S5 | 追加調査要求・Reading Value Assessment Markdown契約を定義する | 承認済み | 未着手 | 最終Assessment／追加調査要求の排他的出力契約 | L4-M3-S1・S4 |
| L4-M3-S6 | Reading Value Skill instructions／referencesを実装し、固有構造契約を確認する | 承認済み | 未着手 | 実行可能なReading Value Skill package | L4 Skill Design Freeze Gate |
| L5 | Parent Orchestration Skillと全体Workflowを設計・実装する | 承認済み | 未着手 | Skill起動、再調査、Budget、終了判定、全体Workflow | 設計着手: Issue・L1。設計確定: L3・L4 producer契約。Mock実装: L5 Gate。実Component統合・完了: L3・L4 Release。後続: L6 |
| L5-M1 | Parent entrypoint・共通Orchestration状態／実行制御契約を詳細設計する | 承認済み | 未着手 | 共通entrypoint、状態遷移、Budget、相関、再試行・終了の制御契約 | 着手: Issue・L1。確定: L3-M1-S3、L3-M2-S3・S4、L4-M1-S4、L4-M2-S5、L4-M3-S5 |
| L5-M1-S1 | Parent entrypoint・Workflow識別／起動受付契約を定義する | 承認済み | 未着手 | 外部受付、Workflow識別、Run初期化要求のumbrella契約 | Issue、L1 |
| L5-M1-S2 | 子Skill Registry・起動／成果物引渡し契約を定義する | 承認済み | 未着手 | 子Skill識別・Version・Invocation・成果物引渡し契約 | L3-M1-S3、L3-M2-S3・S4、L4-M1-S4、L4-M2-S5、L4-M3-S5 |
| L5-M1-S3 | 共通実行状態・診断Envelope契約を定義する | 承認済み | 未着手 | producer意味状態を保持する共通実行状態・診断Envelope | L5-M1-S1・S2 |
| L5-M1-S4 | Run／Claim／Candidate／Cycle・Search Trace相関契約を定義する | 承認済み | 未着手 | ID Owner、成果物系譜、Search Trace相関契約 | L5-M1-S1・S2、L4-M2-S5 |
| L5-M1-S5 | Budget・retry／restart／reroute／stop・副作用安全性契約を定義する | 承認済み | 未着手 | Budget、再実行判定、停止理由、副作用安全性契約 | L5-M1-S3・S4、L1-M1-S4 |
| L5-M2 | 記事価値判定Workflowの起動・Claim fan-out／結果統合・再調査・終了を詳細設計する | 承認済み | 未着手 | 記事価値判定Workflowの状態遷移・再調査・終端設計 | L5-M1、L4-M1-S1・S4、L4-M2-S2・S5、L4-M3-S1・S5 |
| L5-M2-S1 | 記事価値判定Workflowを開始し、Article Analysisを起動・初期成果物を受け入れる | 承認済み | 未着手 | 記事Workflow開始・Article Analysis起動／成果物受入 | L5-M1、L4-M1-S1・S4 |
| L5-M2-S2 | canonical Claim単位のKnowledge Search fan-out契約を定義する | 承認済み | 未着手 | canonical Claimから検索Work Itemへの一対多展開契約 | L5-M2-S1、L4-M2-S2 |
| L5-M2-S3 | Claim別Knowledge Search結果を相関・集約し、Reading Valueへ引き渡す | 承認済み | 未着手 | Claim別結果の機械的fan-in・Reading Value入力Bundle | L5-M2-S2、L4-M2-S5、L4-M3-S1 |
| L5-M2-S4 | 追加再分析／再検索要求のRouting・再調査Cycle反映を定義する | 承認済み | 未着手 | 許可済み追加調査要求の固有Routing・Cycle反映契約 | L5-M2-S3、L5-M1-S4・S5、L4-M1-S4、L4-M2-S5、L4-M3-S5 |
| L5-M2-S5 | 最終Assessment受入・異常／Budget到達時のWorkflow終端契約を定義する | 承認済み | 未着手 | 最終Assessment受入、異常・部分・Budget終端契約 | L5-M2-S1〜S4、L5-M1-S3・S5、L4-M3-S5 |
| L5-M3 | Knowledge蓄積WorkflowのEpisode起動・Candidate分岐・更新終了を詳細設計する | 承認済み | 未着手 | Knowledge蓄積Workflowの分岐・再試行・終端設計 | L5-M1、L3-M1-S1・S3、L3-M2-S1・S3・S4 |
| L5-M3-S1 | Knowledge蓄積Workflowを開始し、Knowledge Acquisitionを起動・成果物を受け入れる | 承認済み | 未着手 | Knowledge蓄積Workflow開始・Acquisition起動／成果物受入 | L5-M1、L3-M1-S1・S3 |
| L5-M3-S2 | candidate／no-candidate／input_insufficient分岐を定義する | 承認済み | 未着手 | Acquisition結果を意味変更しない3分岐契約 | L5-M3-S1、L3-M1-S3、L5-M1-S3・S5 |
| L5-M3-S3 | Candidate更新Work Itemを編成し、Knowledge Updateを起動する | 承認済み | 未着手 | Candidate追跡単位・batch／直列既定のUpdate起動契約 | L5-M3-S2、L3-M2-S1・S3、L5-M1-S4 |
| L5-M3-S4 | Candidate更新結果を集約し、部分成功後の安全な再開順序を定義する | 承認済み | 未着手 | Candidate別結果集約・pending集合・安全な再開順序 | L5-M3-S3、L3-M2-S3・S4、L5-M1-S4・S5 |
| L5-M3-S5 | 正常／異常／Budget到達時のKnowledge蓄積Workflow終端契約を定義する | 承認済み | 未着手 | 正常・部分・異常・Budget終端と未処理範囲 | L5-M3-S1〜S4、L5-M1-S3・S5 |
| L5-M4 | Parent Orchestration Skillを実装し、Mock契約確認後にL3／L4実Componentを統合する | 承認済み | 未着手 | 実行可能なParent SkillとParent固有Component連携確認 | Mock実装: L5 Design Freeze Gate。L3統合: L5-M4-S4＋L3 Release。L4統合: L5-M4-S4＋L4 Release。完了: S5・S6 |
| L5-M4-S1 | Parent entrypoint・Registry・共通実行制御を実装する | 承認済み | 未着手 | root Parent Skill、共通制御、Registry実体 | L5 Design Freeze Gate |
| L5-M4-S2 | 記事価値判定WorkflowのParent instructionsを実装する | 承認済み | 未着手 | 記事Workflow実行reference | L5 Design Freeze Gate。MergeはL5-M4-S1後 |
| L5-M4-S3 | Knowledge蓄積WorkflowのParent instructionsを実装する | 承認済み | 未着手 | Knowledge蓄積Workflow実行reference | L5 Design Freeze Gate。MergeはL5-M4-S1後 |
| L5-M4-S4 | Mock ComponentsでParent固有Orchestration契約を確認する | 承認済み | 未着手 | Parent固有Mock・決定論的契約確認 | L5-M4-S1〜S3 |
| L5-M4-S5 | Parent Orchestration SkillをL3実Componentsへ接続し、Knowledge蓄積連携を確認する | 承認済み | 未着手 | L3実SkillsとのKnowledge蓄積Component連携 | L5-M4-S4、L3 Release |
| L5-M4-S6 | Parent Orchestration SkillをL4実Componentsへ接続し、記事価値判定連携を確認する | 承認済み | 未着手 | L4実Skillsとの記事価値判定Component連携 | L5-M4-S4、L4 Release |
| L6 | 評価設計・横断品質保証を実装する | 承認済み | 未着手 | 共有Fixture・Dataset・評価Harness・Acceptance・E2E評価 | 各Workflow Suite着手: 対応するL5-M4-S5またはS6＋Dataset／Harness readiness。Report確定: 対応M4評価。L6完了: L6-M5-S3＋Final Quality Gate |
| L6-M1 | 横断評価方針・Coverage／判定基準・Search Trace診断／Report契約を詳細設計する | 承認済み | 未着手 | Coverage matrix、評価基準、Trace診断写像、Report schema | Issue、L1〜L5の確定契約 |
| L6-M1-S1 | 評価レイヤ・対象範囲・L1〜L6テスト所有境界を定義する | 承認済み | 未着手 | 評価レイヤ・対象・非対象・Owner方針 | Issue、Task Map |
| L6-M1-S2 | Scenario A〜J・Invariant・Done DefinitionのCoverage／検証責務Matrixを定義する | 承認済み | 未着手 | Requirement―評価レイヤ―観測―Owner―Oracle種別の追跡表 | L6-M1-S1、L1〜L5出力契約 |
| L6-M1-S3 | 評価指標・合否／評価不能・反復判定基準を定義する | 承認済み | 未着手 | 決定論的／Agent評価の判定・反復基準 | L6-M1-S2 |
| L6-M1-S4 | 既存Search Trace・相関情報を用いた失敗原因診断手順を定義する | 承認済み | 未着手 | 既存Trace・相関情報からの失敗原因診断手順 | L6-M1-S1、L4-M2-S5、L5-M1-S4 |
| L6-M1-S5 | 評価結果Envelope・診断／集約Report契約を定義する | 承認済み | 未着手 | Caseから最終Reportまでの共通Envelope・集約契約 | L6-M1-S3・S4 |
| L6-M2 | Scenario A〜J・Invariantを検証する共有Dataset／Fixtureを設計・実装する | 承認済み | 未着手 | version・provenance・oracleを持つ共有Dataset／Fixture | 設計: L6-M1。実装: L6 Evaluation Design Freeze Gate |
| L6-M2-S1 | 共有Dataset／Fixture Schema・Version・Provenance契約を設計する | 承認済み | 未着手 | Dataset／Fixture共通Schema・互換性契約 | L6-M1-S1・S2 |
| L6-M2-S2 | Scenario A〜J・Invariantの具体評価Case Catalog・期待Oracleを設計する | 承認済み | 未着手 | 全Scenario・Invariantの具体Case Catalog・期待Oracle | L6-M2-S1、L6-M1-S3 |
| L6-M2-S3 | 共有Dataset Manifest・Fixture静的検証基盤を実装する | 承認済み | 未着手 | Manifest・Schema・参照・Coverage Validator | L6 Evaluation Design Freeze Gate、L6-M2-S1・S2 |
| L6-M2-S4 | Scenario A〜E・検索／Knowledge State Invariantの共有Fixtureを実装する | 承認済み | 未着手 | A〜E・検索／State領域の初期State・Evidence・期待Assessment Fixture | L6-M2-S3 |
| L6-M2-S5 | Scenario F〜H・Knowledge Acquisition／Update Invariantの共有Fixtureを実装する | 承認済み | 未着手 | F〜H・蓄積領域のEpisode・Evidence・更新前後State Fixture | L6-M2-S3 |
| L6-M2-S6 | Scenario I〜J・Reading Value Invariantの共有Fixtureを実装する | 承認済み | 未着手 | I〜J・Reading Value領域の記事・Knowledge・期待推奨Fixture | L6-M2-S3 |
| L6-M3 | CLI・Skill・Workflow共通評価Harnessを設計・実装する | 承認済み | 未着手 | 共通runner、対象別Adapter、隔離／reset、正規化・Report出力 | 設計: L6-M1。実装: L6 Evaluation Design Freeze Gate |
| L6-M3-S1 | 共通評価Harnessの実行・再現性・隔離要件を定義する | 承認済み | 未着手 | 技術非依存なHarness能力・再現性・隔離要件 | L6-M1-S1・S3 |
| L6-M3-S2 | 評価Harnessの技術スタック・Runner・依存管理方式を選定しADRを確定する | 承認済み | 未着手 | Harness技術選定ADR | L6-M3-S1、Repository技術情報 |
| L6-M3-S3 | Dataset Loader・Target Adapter・Oracle・Trace Collector・Reporter Architectureを詳細設計する | 承認済み | 未着手 | Harness Interface・データフロー・失敗境界設計 | L6-M3-S2、L6-M2-S1、L6-M1-S4・S5 |
| L6-M3-S4 | 評価Harness Core・Dataset Loader・Oracle／Report Pipelineを実装する | 承認済み | 未着手 | 共通Harness Core・Loader・判定・Report Pipeline | L6 Evaluation Design Freeze Gate、L6-M3-S3 |
| L6-M3-S5 | 検索・取得・更新CLI Target Adapterを実装する | 承認済み | 未着手 | CLI Process・JSON・Store隔離Adapter | 着手: L6-M3-S4、L1-M2-S6。完了・Merge: L2-M3-S8 |
| L6-M3-S6 | L3／L4 Component Skill Target Adapterを実装する | 承認済み | 未着手 | Component Skill起動・Markdown回収Adapter | 着手: L6-M3-S4、L3-M1-S1・S3、L3-M2-S1・S3、L4-M1-S1・S4、L4-M2-S2・S5、L4-M3-S1・S5 |
| L6-M3-S7 | Parent Workflow Target Adapterを実装する | 承認済み | 未着手 | Parent entrypoint・Run／Cycle・終端回収Adapter | 着手: L6-M3-S4、L5 Design Freeze Gate。完了・Merge: L5-M4-S4 |
| L6-M4 | CLI横断品質・L3／L4各Component Skillのレイヤ別／Agent評価を実行し、診断Reportを生成する | 承認済み | 未着手 | CLI・L3・L4対象別評価Suiteと診断Report | L6-M1〜M3、各Target Readiness Gate。親単位の一括Release待ちは設けない |
| L6-M4-S1 | 共有FixtureでCLI検索・取得・更新の横断品質を評価する | 承認済み | 未着手 | CLI横断評価Suite・診断Report | 共通: L6-M1-S2・S3・S5、L6-M2-S3、L6-M3-S4。対象: L6-M2-S4・S5、L6-M3-S5、L2-M3-S9 |
| L6-M4-S2 | Knowledge Acquisition SkillのAgent評価を実行する | 承認済み | 未着手 | Acquisition Agent評価Suite・診断Report | 共通: L6-M1-S2・S3・S5、L6-M2-S3、L6-M3-S4。対象: L6-M2-S5、L6-M3-S6、L3-M1-S4 |
| L6-M4-S3 | Knowledge Update SkillのAgent評価を実行する | 承認済み | 未着手 | Update Agent評価Suite・診断Report | 共通: L6-M1-S2・S3・S5、L6-M2-S3、L6-M3-S4。対象: L6-M2-S5、L6-M3-S5・S6、L3-M2-S6。Report確定: L6-M4-S1合格 |
| L6-M4-S4 | Article Analysis SkillのAgent評価を実行する | 承認済み | 未着手 | Article Analysis Agent評価Suite・診断Report | 共通: L6-M1-S2・S3・S5、L6-M2-S3、L6-M3-S4。対象: L6-M2-S6、L6-M3-S6、L4-M1-S5 |
| L6-M4-S5 | Knowledge Search SkillのAgentic Search評価を実行する | 承認済み | 未着手 | Knowledge Search Agentic評価Suite・診断Report | 共通: L6-M1-S2・S3・S5、L6-M2-S3、L6-M3-S4。対象: L6-M2-S4、L6-M3-S5・S6、L4-M2-S7。Report確定: L6-M4-S1合格 |
| L6-M4-S6 | Reading Value SkillのAgent評価を実行する | 承認済み | 未着手 | Reading Value Agent評価Suite・診断Report | 共通: L6-M1-S2・S3・S5、L6-M2-S3、L6-M3-S4。対象: L6-M2-S4・S6、L6-M3-S6、L4-M3-S6 |
| L6-M5 | 記事価値判定・Knowledge蓄積WorkflowのAcceptance／E2E・回帰評価を実行する | 承認済み | 未着手 | 2 WorkflowのAcceptance／E2E・回帰SuiteとReport | Suite着手: 対応Dataset／Harness＋L5-M4-S5またはS6。Report確定: 対応M4評価。完了: S3＋Final Quality Gate |
| L6-M5-S1 | Knowledge蓄積WorkflowのAcceptance／E2E評価を実行する | 承認済み | 未着手 | Knowledge蓄積Workflow E2E Suite・Report | 着手: L6-M2-S5、L6-M3-S5・S7、L5-M4-S5。Report確定: L6-M4-S1〜S3 |
| L6-M5-S2 | 記事価値判定WorkflowのAcceptance／E2E評価を実行する | 承認済み | 未着手 | 記事価値判定Workflow E2E Suite・Report | 着手: L6-M2-S4・S6、L6-M3-S5・S7、L5-M4-S6。Report確定: L6-M4-S1・S4〜S6 |
| L6-M5-S3 | 全Scenario回帰Suiteを実行し、最終Acceptance／診断Reportを確定する | 承認済み | 未着手 | 全Scenario回帰Suite・最終Acceptance／診断Report | L6-M5-S1・S2 |

## Gate前設計成果物のPath・単一Owner（2026-08-11承認済み）

本表はL1の全leafとL2 Design Freeze Gate前の設計leafを対象とする。本番Schema、Migration、Repository、Index、CLI Runtime等の実装PathはL2 Design Freeze Gate通過記録で確定する。

| Task ID | 設計成果物Path／Glob | 単一Owner | Merge前提 |
| --- | --- | --- | --- |
| L1-M1-S1 | `docs/design/knowledge-model/logical-schema.md` | L1-M1-S1 | planning baseline |
| L1-M1-S2 | `docs/design/knowledge-model/relations-and-integrity.md` | L1-M1-S2 | L1-M1-S1 |
| L1-M1-S3 | `docs/design/knowledge-model/evidence-derived-state.md` | L1-M1-S3 | L1-M1-S1 |
| L1-M1-S4 | `docs/design/knowledge-model/update-operations-and-lineage.md` | L1-M1-S4 | L1-M1-S2・S3 |
| L1-M1-S5 | `docs/design/knowledge-model/requirements-traceability.md`、`docs/design/knowledge-model/README.md` | L1-M1-S5 | L1-M1-S1〜S4 |
| L1-M2-S1 | `docs/design/cli-contract/command-catalog.md` | L1-M2-S1 | L1-M1-S5 |
| L1-M2-S2 | `docs/design/cli-contract/collection-contract.md` | L1-M2-S2 | L1-M2-S1 |
| L1-M2-S3 | `docs/design/cli-contract/errors-and-exit-codes.md` | L1-M2-S3 | L1-M2-S1 |
| L1-M2-S4 | `docs/design/cli-contract/schemas/**` | L1-M2-S4 | L1-M2-S1〜S3 |
| L1-M2-S5 | `docs/design/cli-contract/versioning-and-compatibility.md` | L1-M2-S5 | L1-M2-S4 |
| L1-M2-S6 | `docs/design/cli-contract/requirements-traceability.md`、`docs/design/cli-contract/README.md` | L1-M2-S6 | L1-M2-S1〜S5 |
| L2-M2-S1 | `docs/design/search-infrastructure/search-and-index-requirements.md` | L2-M2-S1 | L1、L3-M2-S2、L4-M2-S1 |
| L2-M1-S1 | `docs/design/knowledge-store/adr/0001-persistence-stack.md` | L2-M1-S1 | L2-M2-S1 |
| L2-M1-S2 | `docs/design/knowledge-store/physical-model-history-schema-evolution.md` | L2-M1-S2 | L2-M1-S1 |
| L2-M1-S3 | `docs/design/knowledge-store/repository-transaction-boundary.md`、`docs/design/knowledge-store/README.md` | L2-M1-S3 | L2-M1-S2 |
| L2-M2-S2 | `docs/design/search-infrastructure/adr/0001-search-index-stack.md` | L2-M2-S2 | L2-M2-S1、L2-M1-S1〜S3 |
| L2-M2-S3 | `docs/design/search-infrastructure/candidate-search-architecture.md` | L2-M2-S3 | L2-M2-S2 |
| L2-M2-S4 | `docs/design/search-infrastructure/index-lifecycle-recovery.md`、`docs/design/search-infrastructure/README.md` | L2-M2-S4 | L2-M2-S3、L2-M1-S3 |
| L2-M3-S1 | `docs/design/cli-runtime/adr/0001-runtime-validation-distribution.md` | L2-M3-S1 | 調査着手=L1-M2-S6。ADR確定・Merge=L2-M1-S1＋L2-M2-S2 |
| L2-M3-S2 | `docs/design/cli-runtime/command-to-port-architecture.md`、`docs/design/cli-runtime/README.md` | L2-M3-S2 | L2-M3-S1、L2-M1-S3、L2-M2-S3・S4 |

所有規則:

- 各Directoryの`README.md`は表で指定した終端leafだけが更新し、成果物一覧・依存・参照リンクだけを持つ。
- JSON Schemaと`docs/design/cli-contract/schemas/schema-catalog.md`はL1-M2-S4が`docs/design/cli-contract/schemas/**`で単一所有する。
- ADR番号はDirectory-localとし、global ADR indexは作成しない。
- 最初のTaskだけplanning baselineから分岐し、依存Taskは必要な前提がMergeされた統合Commitから開始する。並行兄弟だけが同じ前提Commitを共有する。
- 設計worktreeはTask Map、決定記録、接続逆引き台帳を変更しない。これらはplanning integration ownerだけが更新する。

## 依存関係

### 全体

```text
L1 論理モデル・CLI公開契約
  ├─→ L2-M1の設計
  ├─→ L3・L4のSkill詳細設計
  └─→ L6の評価仕様・Dataset設計

L4-M2-S1 技術非依存探索要求設計
  + L1-M2検索契約
  + L2-M1物理モデル
        ↓
L2-M2 技術選定・詳細設計
        ↓
L2-M2実装 → L2-M3実装
        ↓
L3・L4のCLI統合 → L5統合
        ↓
L6 Acceptance・E2E・回帰評価
```

L6の評価設計・Fixture・Dataset準備は各対象設計と並行し、最終E2E評価はL5統合後に実施する。

### L1-M1内

```text
L1-M1-S1 論理スキーマ
   ├─→ L1-M1-S2 関連・参照整合性 ─┐
   └─→ L1-M1-S3 Evidence・State契約 ├─→ L1-M1-S4 更新・履歴系譜
                                    ↓
                         L1-M1-S5 要件トレーサビリティ
```

### L1-M2内

```text
L1-M1-S5 論理モデル要件トレーサビリティ
   ↓
L1-M2-S1 公開コマンド・呼び出し契約
   ├─→ L1-M2-S2 結果コレクション契約 ─┐
   └─→ L1-M2-S3 Error・終了コード契約 ─┼─→ L1-M2-S4 JSON Schema
                                         ↓
                              L1-M2-S5 Versioning・互換性
                                         ↓
                              L1-M2-S6 要件トレーサビリティ
```

### L2内

```text
L1 論理モデル・CLI公開契約
              ↓
L2-M1 物理永続化・更新基盤
              ↓
L2-M2 複合検索・Index管理基盤
              ↓
L2-M3 CLI実行基盤・JSON境界
```

### L2-M1内

```text
L2-M1-S1 技術選定ADR
        ↓
L2-M1-S2 物理モデル・履歴・Schema Evolution設計
        ├─→ L2-M1-S4 Baseline Schema ─→ L2-M1-S5 Migration実行基盤 ─┐
        └─→ L2-M1-S3 永続化操作・Transaction・Repository設計 ───────┼─→ L2-M1-S6 永続化Core
                                                                          ↓
                                                           L2-M1-S7 Evidence由来State
                                                            ├─→ L2-M1-S8 追加型更新
                                                            ├─→ L2-M1-S9 非破壊更新
                                                            └─→ L2-M1-S10 詳細・Evidence・履歴取得
```

- L2-M1-S1・S2は、L2-M2-S1が定義する技術非依存の検索・Index Lifecycle要件とReview Gate基準を入力にする。
- L2-M1-S3はL1-M2-S1・S3も入力とし、未Commit変更を公開しないCommit済み正本変更境界を定義する。
- L2-M1-S10はL2-M1-S6・S7を直接入力とし、更新実装S8・S9と並行してFixtureにより実装できる。

### L2横断Design Freeze Gate・worktree計画（2026-08-11承認済み）

責務分類は変更せず、設計成果物間の整合性を確認するMilestoneとして確定した。独立成果物を作る小分類ではないため、新しい小分類IDは付与しない。

```text
L3-M2-S2 更新用探索要求 ─┐
L4-M2-S1 技術非依存探索要求 ─┴→ L2-M2-S1 要件・Review Gate基準
        ↓
L2-M1-S1〜S3 正本の技術選定・詳細設計
        ↓
L2-M2-S2〜S4 検索の技術選定・Query Plane・Lifecycle詳細設計
        ↓
L2-M3-S1・S2 CLI技術選定・Command-to-Port詳細設計
        ↓
L2 Design Freeze Gate
        ↓
L2-M1-S4〜S7 基盤・State実装
        ├─→ L2-M1-S8 追加型更新
        ├─→ L2-M1-S9 非破壊更新
        └─→ L2-M1-S10 取得Repository
                    ↓
          L2-M2-S5 共通検索Core
                    ↓
          L2-M2-S6〜S9 方式別Provider
                    ↓
          L2-M2-S10 Lifecycle統合
```

Gate条件は、L3-M2-S2の技術非依存な更新用探索要求、L4-M2-S1、L2-M2-S1、L2-M1-S1〜S3、L2-M2-S2〜S4、L2-M3-S1・S2の完了に加え、正本Field・Relation・Temporal Metadata・Revision、Commit済み変更境界、正本Snapshot列挙／差分取得、再構築、Schema・Index・Embedding Model・CLI Runtime・JSON Validator・Build方式のVersion互換性と技術共存、全公開CommandからRepository・検索・Lifecycle Portへの対応、正本／派生Index／CLI共有資産の所有者が横断的に整合していることである。

Gate通過記録には、実装構成確定後の`Path／Glob`、資産種別、Owner Task、変更禁止Task、起点Commit、Merge前提・順序、生成物・Lockfile更新Ownerを記載する。

- `L2-M1-S4`以降の本番実装はGate通過後に開始する。Gate前は技術選定用の使い捨てPoCだけを許容する。
- `L2-M2-S5`はGate通過後にPortとMockで着手できるが、正本照合Adapterを含む完了・Mergeには`L2-M1-S10`を必要とする。
- `L2-M1-S4〜S7`はSchema、Migration、接続・変換Coreを共有するため直列にMergeする。
- `L2-M1-S8〜S10`はS6・S7のInterfaceを固定し、共有Error、Mapper、Fixture、DI、Schemaを変更しない場合だけ並行化する。M2を早く解放するMerge順はS10、S8、S9とする。
- `L2-M2-S6〜S9`はS5のProvider契約を共通起点とし、方式別Directoryと固有テストだけを所有して並行化する。
- Migration、共通Interface、Provider Registry、DI、依存定義・Lockfile、生成物、共有Fixtureは単一Ownerだけが変更する。
- Gate通過を記録したCommitを各worktreeの共通起点とし、依存TaskのMerge後に追従してから統合する。既存の未Commit変更は起点Commitへ混入させない。
- `L2-M2`親の`L2-M1`依存は完了依存であり、M2-S1〜S4の設計開始を妨げない。

### L2-M2内（2026-08-11承認済み）

```text
L1-M1・L1-M2・L3-M2-S2更新用探索要求・L4-M2-S1技術非依存探索要求
                    ↓
L2-M2-S1 要件・Review Gate基準
                    ↓
           L2-M1-S1〜S3
                    ↓
L2-M2-S2 技術選定ADR
        ↓
L2-M2-S3 検索・Index詳細設計
        ↓
L2-M2-S4 Index Lifecycle詳細設計
        ↓
L2 Design Freeze Gate
        ↓
L2-M2-S5 共通候補検索Core・正本照合
        ↓
                           S6 Lexical・S7 Semantic・
                           S8 Concept/Relation/Contradiction・S9 Temporal
        ↓
                           L2-M2-S10 同期・Version・再構築・障害回復
```

- S6〜S9は並行実装可能であり、S9 TemporalはS8 Relationへ依存しない。
- S5はS4を直接の設計入力にはしないが、Design Freeze GateがS3・S4双方の完了を確認するため、着手順ではGateを待つ。
- L4-M2-S1をL2-M2-S1の先行入力とし、L4-M2全体の完了は待たない。
- L2-M2-S10はL2-M1-S8・S9のCommit済み変更境界と各検索Providerを統合する。

### L2-M3内（2026-08-11承認済み）

```text
L1-M2 ─────────────→ S1 技術調査開始
L2-M1-S1 + L2-M2-S2 → S1 ADR確定
S1 + L2-M1-S3 + L2-M2-S3・S4
                    ↓
             S2 CLI内部設計
                    ↓
        L2 Design Freeze Gate
                    ↓
             S3 Process Core
                    ↓
             S4 JSON・Error境界
            ┌───────┼───────┐
       S5 取得・検索 S6 更新 S7 Index管理
            └───────┼───────┘
                    S8 配線・配布Artifact
                    ↓
                    S9 Black-box契約適合
```

- S1はL1-M2後に調査を開始し、L2-M1-S1・L2-M2-S2との技術共存を確認してADRを確定する。
- S2はM1・M2のPort設計を入力にし、全公開CommandのOwnerと配線境界を実装前に確定する。
- GateはS1・S2を前提に含み、通過Commitを本番実装の共通起点としてS3、S4を直列Mergeする。
- S5〜S7はS4 Merge後、機能別Directoryと固有テストだけを所有して並行化する。
- Command Registry、DI、設定、Build・配布定義、依存定義・Lockfile、生成物はS8だけが変更する。
- S9はBlack-box契約テストDirectoryだけを所有し、本番修正は該当Owner Taskへ戻す。
- Merge順は`S1 → S2 → Gate → S3 → S4 → S5〜S7 → S8 → S9`とする。
- S9の契約適合はCLI固有の実装検証に限定し、共有Dataset、Agent評価、Acceptance・E2EはL6が所有する。

### L3内（2026-08-11承認済み）

```text
L1-M1
  ├─→ L3-M1 Acquisition設計
  └─→ L3-M2 Update判断設計開始

L3-M1 Candidate Markdown契約 + L1-M2
                    ↓
          L3-M2受入・CLI Mapping設計
                    ↓
       L3 Skill Design Freeze Gate
              ┌─────┴─────┐
              ↓           ↓
          M1実装       M2 Mock実装
                          + L2-M3-S5・S6・S8
                              ↓
                         M2実CLI統合
                              ↓
                    L2-M3-S9後にL3 Release
```

- M1はCandidate Markdown契約の単一Owner、M2はread-only consumerとする。引き渡し契約は独立中分類にしない。
- no-candidateを受けてM2を起動せず終了する分岐はL5が所有する。
- Gateは両Skillの詳細設計、Markdown契約、no-candidate、L1 Field・CLI Mapping、責務境界、Path Ownerの整合を確認するMilestoneであり、独立Task IDを付与しない。
- Gate通過Commitを実装worktreeの共通起点とし、M1・M2は別Skill Directoryだけを所有して並行化する。
- 共有Registry・OrchestrationはL5、共有Fixture・Dataset・評価HarnessはL6が単独所有する。
- 実Path／GlobとGateを構成する小分類Task IDは、各中分類の小分類時に確定する。

#### L3 Skill Design Freeze Gateの追加条件（2026-08-11承認済み）

`no-candidate`を、十分なEpisode評価後に保存候補がない正常結果として扱い、入力欠落・切断・Source Reference解決不能による評価不能と区別できる入力・出力契約をGate条件へ追加する。M1は状態・診断を出力し、再試行・終了・M2非起動の制御はL5が所有する。

### L3-M1内（2026-08-11承認済み）

```text
L1-M1 ─→ S1開始
L1-M1 ─→ S2先行設計
S1完了 ─→ S2確定
S1 + S2 ─→ S3
S3 ─→ L3-M2受入設計確定
S3 + L3-M2設計完了
       ↓
L3 Skill Design Freeze Gate
       ↓
      S4

S3 ─→ L6評価設計
S4 ─→ L6 Agent評価実行
```

- S1・S2は別の設計Fileを所有し、S2は先行作業できるが、確定・MergeはS1後とする。
- S3はCandidate契約と規範例の単一Ownerであり、S4・L3-M2はread-onlyで利用する。
- S4はGate通過後にSkill本体と補助Referenceを実装し、静的構造・必須Section・status契約を確認する。
- 共有Registry・OrchestrationはL5、共有Fixture・Dataset・Harness・意味的Agent評価はL6が所有する。
- Merge順は`S1 → S2 → S3 → Gate → S4`とし、Gate通過CommitをM1・M2本番実装worktreeの共通起点にする。

### L3-M2内（2026-08-11承認済み）

```text
L3-M1-S3 + L1 ─→ S1・S2
S1 + S2 ───────→ S3
L1-M2 ─────────→ S4先行設計
S2 + S3 ───────→ S4確定
S1〜S4 ────────→ L3 Skill Design Freeze Gate
Gate ──────────→ S5 Mock前提実装
S5 + L2-M3-S5・S6・S8
              ─→ S6 実CLI統合
S6 + L3-M1-S4 + L2-M3-S9
              ─→ L3 Release
```

- S1・S2は別設計Fileで並行可能とし、S3は両方を入力に確定する。
- S4はL1-M2から先行設計できるが、検索MappingはS2、更新Mappingと結果引渡しはS3を待って確定する。
- S5はGate後にMockで先行し、S6のみL2-M3-S5・S6・S8を待つ。L2-M3-S9はL3 Release条件とする。
- S1〜S4の成果物、S5のSkill package、S6のComponent連携資産はそれぞれ単独Ownerとし、L3-M1契約とL2 CLI実装はread-onlyで利用する。
- Merge順は`L3-M1-S3 → S1・S2 → S3 → S4 → Gate → S5 → S6`とする。

#### 承認済み構造への追加条件（2026-08-11承認済み）

1. L3 Skill Design Freeze Gateへ、「正常な非永続化判断」「Candidate入力不備・評価不能」「CLI実行失敗・部分成功」を区別するM2出力・診断契約を追加する。M2は報告、L5は再試行・停止・Workflow継続を所有する。
2. L3-M2-S2の技術非依存な更新用検索要求を、L2-M2-S1とL2 Design Freeze Gateの入力へ追加する。L3-M2-S2はL1とL3-M1-S3から先行設計し、L2検索基盤を更新利用にも適合させる。

### L4内（2026-08-11承認済み）

```text
Issue + L1 ─→ M1 Claim／Support契約設計
Issue + L1 ─→ M2 技術非依存探索要求の先行設計 ─→ L2-M2-S1
M1 Claims契約 + L1-M2 ─→ M2 Assessment／Trace／CLI Mapping確定
M1 Claims契約 + M2 Assessment契約 ─→ M3 Reading Value契約確定
M1〜M3詳細設計 ─→ L4 Skill Design Freeze Gate
Gate ─→ M1実装・M2 Mock実装・M3実装
M2 Mock実装 + L2-M3-S5・S8 ─→ M2実CLI統合
M1〜M3完了 + L2-M3-S9 ─→ L4 Release
```

- Gateは3 Skillの契約、状態、責務境界、CLI Mapping、Path Ownerを確定するMilestoneであり、独立Task IDを付けない。
- GateではKnowledge Assessmentの8状態と、`complete / input_insufficient / incomplete / partial_failure`の探索実行状態を別軸として確定する。
- L2実行ArtifactをGate条件に含めず、探索要求からL2を設計する依存との循環を防ぐ。
- L4-M1〜M3は別Skill Directoryを所有し、Gate通過Commitから並行実装する。共有Registry・Orchestration・Trace相関はL5、共有評価資産はL6が所有する。
- Merge順は`M1契約 → M2契約 → M3契約 → Gate → M1・M2 Mock・M3実装 → M2実CLI統合 → L4 Release`とする。

#### 承認済み構造への修正（2026-08-11承認済み）

1. L4親依存を`設計着手: L1。本番Skill実装: L4 Skill Design Freeze Gate。L4-M2実CLI統合: L2-M3-S5・S8。完了・Release: L4-M1〜M3、L2-M3-S9`へ細分化した。
2. L2-M2-S1、L2-M2親、L2 Design Freeze Gateの暫定前方参照を、`L4-M2-S1`へ確定した。

### L4-M1内（2026-08-11承認済み）

```text
Issue + L1 ─┬─→ S1 取得方式・入力境界 ─┐
            └─→ S2 Overview・Claim正規化 ┴→ S3 Location・Support追跡
                                                   ↓
                                      S4 Markdown・局所再分析結果契約
                                        ├─→ L4-M2契約確定
                                        ├─→ L4-M3契約確定
                                        ├─→ L5 Routing設計
                                        └─→ L6評価設計

S1〜S4 + L4-M2・M3詳細設計
              ↓
L4 Skill Design Freeze Gate
              ↓
S5 Skill実装・固有構造契約確認
              ↓
L6 Agent評価
```

- S1とS2はL1と承認済みPath Ownerを含むCommitを共通起点として別worktreeで並行し、`S1・S2 → S3 → S4 → L4 Gate → S5`の順にMergeする。
- S1〜S4は各設計Fileの単一Owner、S5はS4契約File以外のArticle Analysis Skill packageの単一Ownerとする。
- S4はM2・M3・L5・L6へ引き渡すproducer契約であり、consumerはread-onlyで利用する。
- S5はL2や実CLIを待たず、L4 Gate通過Commitから本番実装する。

### L4-M2内（2026-08-11承認済み）

```text
Issue + L1-M1・L1-M2
          ↓
         S1 技術非依存探索要求 ─→ L2-M2-S1・L2 Design Freeze Gate

S1 + L4-M1-S4
          ↓
         S2 Claim受入・検索variant生成
        ┌─┴──────────┐
        ↓            ↓
S3 CLI Mapping   S4 Assessment判断
        └─────┬──────┘
              ↓
S5 Assessment・Trace・局所再検索結果契約
              ↓
L4-M1-S4 + S5 → L4-M3詳細設計
              ↓
M1〜M3詳細設計 → L4 Skill Design Freeze Gate
              ↓
S6 Mock CLI Skill実装
              ↓
S6 + L2-M3-S5・S8 → S7 実CLI統合
              ↓
M1〜M3実装 + L2-M3-S9 → L4 Release
```

- S1とL4-M1-S4は別Pathで並行し、S2は両方の統合後に確定する。
- S3とS4はS2統合Commitから別worktreeへ分岐して並行できる。
- S4はKnowledge Assessmentの意味判定規則、S5はAssessment・raw Trace・実行状態の外部表現を所有する。
- S6はL4 Gate後にMockで実装し、S7だけがL2-M3-S5・S8の実Artifactを待つ。L2-M3-S9はL4 Release条件とする。
- Merge順は`S1 → S2 → S3・S4 → S5 → L4 Gate → S6 → S7`とする。
- 各Taskの単一Owner Pathは、S1〜S5が各設計・契約File、S6がS3・S5以外のKnowledge Search Skill package、S7が`tests/component/knowledge-search/**`である。
- 共有Registry、Routing、Budget、再試行、相関はL5、共有Fixture、Dataset、Harness、Agent評価、Acceptance、E2EはL6が所有する。

#### 承認済み構造への修正（2026-08-11承認済み）

1. L2-M2親の開始入力、L2-M2-S1の直接依存、L2 Design Freeze Gateの確認対象にある暫定前方参照を、`L4-M2-S1 / agentic-search-requirements.md`へ置換する。
2. L4 Skill Design Freeze Gateへ、Knowledge Assessmentの8状態と、`complete / input_insufficient / incomplete / partial_failure`の探索実行状態を別軸で確定する条件を追加する。

### L4-M3内（2026-08-11承認済み）

```text
L4-M1-S4 + L4-M2-S5
          ↓
S1 入力整合・評価可否境界
     ┌────┴────┐
     ↓         ↓
S2 Recognition Gain   S3 Reliability・Applicability
     └────┬────┘
          ↓
S4 Attention Cost・Recommendation・読書範囲
          ↓
S1 + S4 → S5 最終Assessment／追加調査要求契約
          ↓
M1〜M3詳細設計 → L4 Skill Design Freeze Gate
          ↓
S6 Reading Value Skill実装
          ↓
L5統合・L6評価
```

- S2とS3はS1統合Commitから別worktreeへ分岐し、並行できる。
- S2はKnowledge Assessmentを再判定せずRecognition Gainへ適用し、S3は記事内SupportをReliability／Applicability判断へ適用する。
- S4のAttention Costは既存入力から得られる定性的proxyを用い、新規の上流Field、精密な読了時間、単一スコアを前提にしない。
- S5は`final_assessment`と`additional_research_request`を排他的に表現する。再調査のRouting、回数・Budget、再試行、相関、終了判定はL5が所有する。
- L4 Skill Design Freeze Gateでは、Knowledge Assessmentの8状態と実行状態の直交、SupportとReliabilityの境界、3 Recommendationの固定、最終Assessmentと追加調査要求の排他性、Path Ownerを確定する。
- Merge順は`L4-M1-S4・L4-M2-S5 → S1 → S2・S3 → S4 → S5 → L4 Gate → S6`とする。
- S1〜S4は各設計File、S5はReading Value出力契約File、S6はS5契約File以外のReading Value Skill packageを単一Ownerとする。共有Registry・OrchestrationはL5、共有評価資産はL6が所有する。
- Reading Value SkillはCLIを直接利用しないため、独立した技術選定、Mock CLI、実CLI統合タスクは追加しない。

### L5内（2026-08-11承認済み）

```text
L3-M1-S3 + L3-M2-S3・S4
+ L4-M1-S4 + L4-M2-S5 + L4-M3-S5
                    ↓
        L5-M1 共通Orchestration設計
              ┌─────┴─────┐
              ↓           ↓
      L5-M2 記事設計   L5-M3 蓄積設計
              └─────┬─────┘
                    ↓
          L5 Design Freeze Gate
                    ↓
            L5-M4 Mock実装
              ┌─────┴─────┐
              ↓           ↓
        L3 Release     L4 Release
              ↓           ↓
         L5-M4-S5      L5-M4-S6
              └─────┬─────┘
                    ↓
             L5-M4 完了
                    ↓
        L6 最終Acceptance・E2E
```

- 中分類は、共通制御設計、記事Workflow設計、Knowledge蓄積Workflow設計、Parent Skill実装・統合の4件を提案する。
- M2とM3はM1統合Commitから別worktreeへ分岐し、別設計Fileを所有して並行できる。
- M4は中分類としてParent Skillの完成責任を一つに保つが、小分類ではMock実装と実Component統合を依存の異なるleafへ必ず分ける。
- L5 Design Freeze Gateは、entrypoint・Registry、状態遷移、ID相関、Budget、再試行・再調査・終了、副作用重複防止、Claim fan-out／fan-in、結果の排他性、Path Owner、Mock差替え境界を確定するMilestoneであり、独立Task IDを付けない。
- 記事価値判定とKnowledge蓄積は別Workflowとし、記事評価結果からKnowledge更新を自動起動しない。
- L5-M1〜M3は`docs/design/orchestration/**`、L5-M4 Mock実装は`.agents/skills/parent-orchestration/**`、実統合は`tests/component/parent-orchestration/**`を単一Ownerとする。L3・L4 Skill packageとL2 CLIはread-onlyで利用し、共有評価資産はL6が所有する。
- Merge順は`producer契約 → M1 → M2・M3 → L5 Gate → M4 Mock実装 → {L3 Release → S5、L4 Release → S6} → M4完了 → L6`とする。S5とS6は独立開始・並行Mergeできる。

#### 承認済み構造への修正（2026-08-11承認済み）

L5親の旧依存`設計着手: L1。完了・Merge: L2-M3、L3、L4`を、`設計着手: Issue・L1。設計確定: L3・L4 producer契約。Mock実装: L5 Design Freeze Gate。実Component統合・完了: L3 Release・L4 Release。後続: L6最終Acceptance・E2E`へ細分化した。

L2-M3-S9はL3・L4 Releaseに含まれるため、L5からの直接依存を削除した。

### L5-M1内（2026-08-11承認済み）

```text
Issue・L1 ───────────────→ S1 外部受付

L3-M1-S3 + L3-M2-S3・S4
+ L4-M1-S4 + L4-M2-S5 + L4-M3-S5
                         ─→ S2 Registry・引渡し

S1 + S2 ─┬─→ S3 共通状態・診断Envelope
         └─→ S4 相関契約
                  ↓
S3 + S4 + L1-M1-S4
                  ↓
S5 Budget・実行制御・副作用安全性
                  ↓
           L5-M2・L5-M3
                  ↓
       L5 Design Freeze Gate
```

- S1はParentの外部受付、S2はParentから子Skillへの内部Invocation境界を所有する。
- S3は実行状態と診断Envelope、S4はRun／Claim／Candidate／CycleとSearch Traceの識別・相関を所有する。
- S5はS3・S4とL1-M1-S4を入力に、Budget、retry／restart／reroute／stop、副作用不明時の非盲目的再実行を定義する。CLIのexactly-onceや冪等性は再設計しない。
- Workflow固有のClaim fan-out／fan-inと再調査配送はL5-M2、Candidate分岐と再開順序はL5-M3、Skill packageとRegistry実体はL5-M4、共有評価資産はL6へ残す。
- 新しいGateは追加せず、S1〜S5を既存L5 Design Freeze Gateの入力にする。
- S1は`docs/design/orchestration/common-control-contract.md`、S2は`child-skill-registry-contract.md`、S3は`execution-state-envelope.md`、S4は`correlation-trace-contract.md`、S5は`execution-control-policy.md`の単一Ownerとする。
- S1はIssue・L1から先行Mergeできる。S2のproducer契約依存が揃った後、`S1・S2 → S3・S4 → S5 → L5-M2・M3 → L5 Gate`の順でMergeする。S3・S4はS1・S2統合Commitから別worktreeへ分岐して並行できる。
- 5件はいずれも小分類leafであり、これ以上子タスクへ分割しない。

### L5-M2内（2026-08-11承認済み）

```text
L5-M1 + L4-M1-S1・S4
              ↓
S1 記事Workflow開始・Article Analysis受入
              ↓
L4-M2-S2 + S2 Claim fan-out
              ↓
L4-M2-S5 + L4-M3-S1 + S3 結果fan-in・引渡し
              ↓
L4-M1-S4 + L4-M2-S5 + L4-M3-S5 + S4 再調査Routing
              ↓
S1〜S4 + L5-M1-S3・S5 + L4-M3-S5
              ↓
S5 最終Assessment受入・Workflow終端
              ↓
L5 Design Freeze Gate
```

- S1はArticle Analysisの一対一起動、S2はcanonical Claimから検索Work Itemへのfan-out、S3はClaim別結果のfan-inを所有する。
- S4はL5-M1-S5で許可された再調査のWorkflow固有Routing、S5は正常・部分・異常・Budget到達を有限に閉じる終端を所有する。
- 実行時にS4からS2・S3相当へ戻る遷移はWorkflow内のCycleであり、Task DAGの循環ではない。
- S1〜S5は別Fileを単一所有する。別worktreeで先行ドラフトできるが、確定・Mergeは`S1 → S2 → S3 → S4 → S5 → L5 Gate`とする。
- 直接依存へ`L4-M1-S1、L4-M2-S2、L4-M3-S1`を追加し、起動入力・Claim受入・Reading Value入力整合の隠れ依存を解消した。
- 5件はいずれも小分類leafであり、これ以上子タスクへ分割しない。

### L5-M3内（2026-08-11承認済み）

```text
L5-M1 + L3-M1-S1・S3
              ↓
S1 蓄積Workflow開始・Acquisition受入
              ↓
L5-M1-S3・S5 + S2 Acquisition結果3分岐
              ↓
L3-M2-S1・S3 + L5-M1-S4 + S3 Update Work Item編成・起動
              ↓
L3-M2-S3・S4 + L5-M1-S4・S5 + S4 結果集約・安全な再開
              ↓
S1〜S4 + L5-M1-S3・S5
              ↓
S5 Knowledge蓄積Workflow終端
              ↓
L5 Design Freeze Gate
```

- S1はKnowledge Acquisitionの一対一起動、S2は`candidate／no-candidate／input_insufficient`の意味を変えない分岐を所有する。
- S3はCandidateを追跡可能なWork Itemへ編成する。副作用競合を避けるためbatchまたは直列を既定とし、L3が非競合を明示した場合だけ並行化する。
- S4は成功済みを除外したpending集合と安全な再開順序、S5は正常・部分・異常・Budget到達時の終端を所有する。
- S1〜S5は別Fileを単一所有し、`S1 → S2 → S3 → S4 → S5 → L5 Gate`でMergeする。L5-M2とは別Directoryで並行可能である。
- 直接依存へ`L3-M1-S1、L3-M2-S1`を追加し、Episode／Source ReferenceとCandidate受入の隠れ依存を解消した。
- 5件はいずれも小分類leafであり、これ以上子タスクへ分割しない。

### L5-M4内（2026-08-11承認済み）

```text
L5-M1〜M3
    ↓
L5 Design Freeze Gate
    ├─→ S1 共通Parent実装 ───┐
    ├─→ S2 記事Workflow実装 ─┼─→ S1〜S3統合 → S4 Mock契約確認
    └─→ S3 蓄積Workflow実装 ─┘
                          ↓
                ┌─────────┴─────────┐
                ↓                   ↓
         L3 Release + S5     L4 Release + S6
                └─────────┬─────────┘
                          ↓
                      L5-M4完了
                          ↓
                L6 Acceptance・E2E
```

- S1はroot `SKILL.md`・共通references・Registryの単一Owner、S2・S3はWorkflow別referenceの単一Ownerとする。
- S1〜S3はGate通過Commitから別worktreeで並行実装できる。Mergeはroot OwnerのS1を先にし、S2・S3は順不同、その後S4へ進む。
- S4は実Component Release前にParent固有の状態遷移・相関・Budget・再調査・終端をMockで確認する。共有Dataset・Agent評価・Acceptance・E2EはL6へ残す。
- S5はS4とL3 Releaseだけ、S6はS4とL4 Releaseだけを待ち、別Directoryで並行できる。L5-M4完了は双方の合格を必要とする。
- 旧構造の「L3・L4 Release後に実Component統合を一括開始」を上記2分岐へ修正した。L5親の最終完了条件は変えない。
- 6件はいずれも小分類leafであり、これ以上子タスクへ分割しない。

### L6内（承認済み）

```text
M1-S1 → M1-S2 → M1-S3 ─┐
   └────→ M1-S4 ─────────┼→ M1-S5
                         ↓
       M2-S1 → M2-S2     M3-S1 → M3-S2 → M3-S3
               └───────────┬──────────────┘
                           ↓
           L6 Evaluation Design Freeze Gate
             ┌─────────────┴─────────────┐
             ↓                           ↓
       M2-S3 Dataset検証           M3-S4 Harness Core
        ┌────┼────┐                ┌────┼────┐
        ↓    ↓    ↓                ↓    ↓    ↓
      M2-S4 S5   S6              M3-S5 S6   S7
        └────┴────┴──────────┬─────┴────┴────┘
                             ↓
                    M4 対象別品質評価
                    ┌────────┴────────┐
                    ↓                 ↓
        M5-S1 Knowledge蓄積     M5-S2 記事価値判定
                    └────────┬────────┘
                             ↓
                 M5-S3 全回帰・最終Report
                             ↓
                    Final Quality Gate
```

- 中分類は、評価規範、共有Dataset／Fixture、共通Harness、対象別Component品質、Workflow Acceptance／E2Eの5件とする。
- 小分類はM1=5件、M2=6件、M3=7件、M4=6件、M5=3件の計27件とし、すべてこれ以上子を持たないleafとする。
- M1-S2は抽象的な検証責務、M2-S2は具体Caseと期待Oracleを所有する。M1-S4はL4・L5のTrace／相関をread-onlyで診断に利用する。
- M2-S4〜S6へInvariantを検索／State、Acquisition／Update、Reading Valueの領域別に配分し、`残余Invariant`を置かない。
- M3は技術非依存要件、技術ADR、ArchitectureをGate前に終え、Gate後にCoreと3 Adapterを実装する。技術選定を実装へ混在させない。
- Gate後はM2-S3とM3-S4を別worktreeで並行し、その統合後にScenario別FixtureとTarget別Adapterを並行できる。
- M4は各TargetのArtifactとHarness Adapterが揃った時点で開始する。M4-S3・S5のSuite準備はCLI評価と並行できるが、Report確定はM4-S1合格後とする。
- M5-S1はKnowledge蓄積、S2は記事価値判定の必要資産だけで着手し、対応M4 sliceの合格後にReportを確定する。S1・S2は別worktreeで並行し、S3が全回帰と最終Reportを単独所有する。

#### L6 Evaluation Design Freeze Gate

`M1-S1〜S5 + M2-S1・S2 + M3-S1〜S3`を入力に、Issue 44・45・49・52章のRequirement Coverage、具体Case、Dataset schema／version／provenance、Harness要件／ADR／Architecture、Agent反復／評価不能、既存Trace診断、Report schema、Path Owner、起点Commit、Merge順を確定する。実装・評価合格は含めず、Task IDを付けない。

#### Final Quality Gate

M5-S3の全回帰・最終Reportが閾値を満たし、未解決blockerがないことだけを確認するMilestoneとし、Task IDを付けない。

#### worktree所有境界

- M1: `docs/design/evaluation/spec/**`
- M2設計: `docs/design/evaluation/datasets/**`
- M2実装: `tests/evaluation/datasets/**`
- M3設計・ADR: `docs/design/evaluation/harness/**`、`docs/design/evaluation/adr/**`
- M3実装: `tests/evaluation/harness/**`
- M4: `tests/evaluation/suites/{cli,l3,l4}/**`、`tests/evaluation/reports/{cli,l3,l4}/**`
- M5: `tests/evaluation/workflows/**`、`tests/evaluation/reports/workflows/**`、`tests/evaluation/reports/final/**`

共有Fixture、Harness、Report schemaはM4・M5から変更せず、対応Ownerへ戻す。CIの一時Reportはartifact出力とし、固定baselineだけをReport Directoryで所有する。

#### 承認済み構造への修正候補

L6親の旧依存`最終Acceptance・E2E: L5`を、`各Workflow Suite着手: 対応するL5-M4-S5またはS6＋Dataset／Harness readiness。Report確定: 対応M4評価。L6完了: M5-S3＋Final Quality Gate`へ具体化する。L1〜L5の階層・責務は変更しない。L6中分類・小分類と一括して再承認する。

### L5-M2〜M4で確定した修正（2026-08-11承認済み）

| 対象 | 現状 | 修正案 | 理由 |
| --- | --- | --- | --- |
| L5-M2直接依存 | L5-M1、L4-M1-S4、L4-M2-S5、L4-M3-S5 | L4-M1-S1、L4-M2-S2、L4-M3-S1を追加 | 起動入力、Claim受入、Reading Value入力整合の隠れ依存を明示する |
| L5-M3直接依存 | L5-M1、L3-M1-S3、L3-M2-S3・S4 | L3-M1-S1、L3-M2-S1を追加 | Episode／Source ReferenceとCandidate受入の隠れ依存を明示する |
| L5-M4実統合開始 | L3・L4 Releaseの双方を待って一括開始 | `S4 + L3 Release → S5`と`S4 + L4 Release → S6`へ分岐 | 別Workflow・別Directoryであり、相手側Releaseを待つ必要がない |

いずれも階層や責務を変更せず、依存を正確化して不要な直列待ちを除く修正として、L5-M2〜M4の小分類と同時に承認された。

## 承認済み構造への修正（2026-08-11承認済み）

L2-M2の小分類監査で、決定009に含まれる依存関係を具体化すると循環することが判明した。ユーザー承認により次の変更を確定した。

| 対象 | 旧状態 | 確定した変更 | 理由 |
| --- | --- | --- | --- |
| L2-M1-S1・S2 | L2-M2の同期・再構築要件に対するReview Gateを待つ | L2-M2-S1の技術非依存要件・Gate基準へ明示的に依存する | L2-M2がL2-M1へ依存する循環を、要件→正本設計→検索技術選定の順序に直す |
| L2-M1-S10 | L2-M1-S7〜S9へ依存する | L2-M1-S6・S7へ依存する | 取得Repositoryは更新処理の完成を待たず、物理CoreとEvidence由来Stateを入力にFixtureで独立実装できる |
| L2-M2の親成果物 | 4種の候補検索 | Issue 12章の全検索プリミティブを支える候補検索 | 6公開操作と4実装系統を混同しない |

## 責務の配置

| 内容 | 配置先 |
| --- | --- |
| 論理レコード、関連、Evidence由来State、論理履歴 | L1-M1 |
| CLI引数、JSON Schema、エラー、終了コード、契約Version | L1-M2 |
| 物理DB、Index、Embedding、Migration実装 | L2 |
| 会話・作業からのEvidence抽出とKnowledge更新判断 | L3 |
| Article Claim分析、Knowledge Search判断、Reading Value判断、Search Trace生成 | L4 |
| Skill起動、Search Trace相関、再調査、Budget、終了判定、全体制御 | L5 |
| 各タスク固有のUnit・Contract・Component Integrationテスト | 各実装タスク |
| 共有Fixture・Dataset、評価Harness、Agent評価、Acceptance・E2E・回帰評価 | L6 |

## 更新規則

1. 分割案を提示した時点で、子タスクを「議論中」として追加する。
2. 承認後に「承認済み」へ変更し、議論経緯は`task-decomposition-decisions.md`へ追記する。
3. 実装開始・完了は「実行状態」だけを更新し、過去の承認記録を書き換えない。
4. タスクを統合・移動・削除する場合は、先に議論記録へ理由を残してから本マップを更新する。
5. 小分類は最終実行単位とし、子タスクを作らない。作業手順やチェック項目は小分類の完了条件として管理する。
6. 承認はその時点の暫定的な構造確定とし、後続分割で矛盾、欠落、重複、粒度不一致、階層誤りが判明した場合は承認済みタスクも修正候補に戻す。
7. 承認済みタスクを変更する場合は、旧構造を議論記録へ残し、影響範囲と移行対応を示した修正案を再承認してからTask Mapへ反映する。
