# FEAT-005 要件

## 目的

Knowledge Search、Knowledge Acquisition／Knowledge Updateを、同じ再現可能な受入Fixtureと実CLIプロセス境界から層別に検証し、期待結果と実際の差異を責務境界まで追跡できるようにする。Reading ValueのI/Jは既存FEAT-003検証契約の必須観測節への対応付けとして保持し、その通常URL入力Workflowを本Featureから起動・変更しない。品質保証は既存の製品挙動を変更せず、その正しさを観測する。

## 含む要件

- REQ-021
- BR-002、BR-003、BR-005、BR-008、BR-009
- NFR-006
- NFR-001、NFR-002、NFR-003、NFR-004、NFR-005（Fixtureの隔離、根拠・履歴・境界の検証に必要な範囲）
- CON-002、CON-003、CON-004、CON-006
- Issue #175 §43〜§45のScenario A〜J

## Feature Acceptance

1. Scenario A〜Jのそれぞれについて、入力、初期Knowledge Store、必要な記事またはEpisode、期待する観測結果、検証層を一意に定義できる。
2. 同じFixtureを再実行しても、隔離Storeと固定入力により同じ判定・CLI境界・更新後状態を観測できる。
3. Knowledge Searchを実行するケースでは、Search Traceから検索operation、query、結果ID、Evidence ID、Budget、停止理由、技術失敗または中断の位置を追跡できる。
4. 失敗時は、CLI／Store、Knowledge Search、Knowledge Acquisition、Knowledge Updateのどの層で期待との差異が起きたかを、当該層の成果物・プロセス記録・Store差分から切り分けられる。Reading Valueは既存FEAT-003検証契約の必須観測節を参照し、本Featureで新たな実行結果を作らない。
5. Scenario FとGではKnowledge Storeを更新せず、Scenario Hでは過去Evidenceを削除せずに訂正後の評価可能な履歴を残す。
6. Scenario Iは`read_selected`、Scenario Jは未知Claim数だけを根拠に`read_full`にしないことを、既存FEAT-003検証契約のV-002に一意に対応付けられる。
7. このFeatureの実装は、Knowledge CLIの公開operation、JSON、stdout/stderr、exit code、SQLite schema・migration、公開設定、および各Workflowの業務判断を追加・変更しない。

## 範囲外

- Knowledge Searchのquery生成、検索順序、状態分類、Budget、停止規則の変更（FEAT-002）。
- 記事取得、Article Claim分解、推奨規則、公開ネットワーク境界の変更（FEAT-003）。
- Candidate抽出、更新操作選択、同期更新の振る舞いの変更（FEAT-004）。
- Knowledge CLI、SQLite schema・migration、公開CLI JSON／option、保存先・公開設定の変更（FEAT-001）。
- Semantic Search、Embedding、Vector Indexの評価（FEAT-006）。
- 任意の外部記事サイトを使う非決定論的な受入試験、性能ベンチマーク、一般的な記事品質評価。
- Reading Valueの通常URL入力Workflowを固定Article Analysis／Assessment Mapで起動するテスト用入口の追加、またはそのSkill本文・参照資料へのFixture／oracle追記。

## 受入Scenarioの正規対応

| Scenario | 主題 | 主な検証層 | 必須の観測結果 |
| --- | --- | --- | --- |
| A | 空Store | Knowledge Search、Reading Value参照 | `no_evidence`を未知・初心者へ変換せず、Traceに停止理由を残す。Reading Valueの結論はFEAT-003検証契約を参照する。 |
| B | Target Claim完全一致 | CLI／Knowledge Search | 強いEvidenceを根拠に`known`となる。 |
| C | 構成知識 | Knowledge Search | 直接一致なしでも`partially_known`または`inferable`となり、KnownとGapを分離する。 |
| D | 誤った理解 | CLI／Knowledge Search | Evidenceを根拠に`contradicted`となる。 |
| E | 古い知識 | CLI／Knowledge Search | 時点・訂正根拠から`outdated`となる。 |
| F | 質問だけ | Knowledge Acquisition／Knowledge Update | Candidate、Decision、CLI更新がない。 |
| G | AIによる説明だけ | Knowledge Acquisition／Knowledge Update／Reading Value参照 | 記事・AI説明をEvidenceへ保存せず、更新がない。Reading Value本文は本Featureで再実行しない。 |
| H | 後日の訂正 | Update／CLI／Store | 過去Evidenceを保持し、訂正Evidenceと履歴から再評価可能である。 |
| I | 大半が既知・一部重要な未知 | Reading Value検証契約参照 | `read_selected`、読む／飛ばす位置、Gain、根拠参照を求めるFEAT-003 V-002へ対応付ける。 |
| J | 些末な未知 | Reading Value検証契約参照 | 未知候補だけで`read_full`にせず、Attention Cost・重要度・根拠を求めるFEAT-003 V-002へ対応付ける。 |

## 依存

- FEAT-001が提供するCLI、SQLite、process境界Fixture。
- FEAT-002が提供するKnowledge Assessment、Search Trace、状態分類とBudgetの契約。
- FEAT-003が提供するArticle Analysis、Assessment Map、Reading Value Assessmentの契約。
- FEAT-004が提供するEpisode、Candidate、Update Decision／Resultおよび同期更新の契約。
