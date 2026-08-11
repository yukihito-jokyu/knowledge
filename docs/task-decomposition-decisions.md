# タスク分割の議論・承認記録

## この文書の目的

Issue #1「統合設計仕様書」を実装可能なタスクへ段階的に分割する際の、提案、分割理由、議論、承認結果を記録する。

- 対象Issue: https://github.com/yukihito-jokyu/knowledge/issues/1
- 記録開始日: 2026-08-10
- 現在の段階: `L5-M1`〜`L5-M4`の小分類と依存修正は承認済み。`L6`を中分類から小分類まで連続して分割中

## 決定 001: タスクの大分類

### 状態

承認済み（2026-08-10）

### 承認された分類

1. 共通契約・詳細設計
2. Knowledge Store／CLI基盤
3. Knowledge蓄積ワークフロー
4. 記事価値判定ワークフロー
5. オーケストレーション／統合
6. 評価・品質保証

### 分割理由

#### 1. 共通契約・詳細設計

Knowledge Model、JSON Schema、CLIコマンド、判定ステータスなど、CLIと各Skillが共有する契約を先に確定する必要があるため。

#### 2. Knowledge Store／CLI基盤

保存・検索・更新を行う決定論的な基盤を、意味判断を担当するCodexから分離し、「Codexが判断し、CLIが保存・取得する」という責務境界を維持するため。

#### 3. Knowledge蓄積ワークフロー

会話や作業からEvidenceを抽出し、既存Knowledgeとの照合を経て知識状態を更新する処理が、独立した一連のワークフローを形成するため。

#### 4. 記事価値判定ワークフロー

URLからの記事分析、Claim単位の知識検索、読む価値の判定が、利用者へ価値を届ける一つの連続したワークフローを形成するため。

#### 5. オーケストレーション／統合

個別Skill間の成果物受け渡し、再調査、探索回数制御、終了判定は、各Skill固有の責務ではなく全体制御の責務であるため。

#### 6. 評価・品質保証

「未観測は未知ではない」などの不変条件とAcceptance Scenarioを、CLI、Skill、End-to-Endの各レイヤーで横断的に検証する必要があるため。

### 依存関係

基本的な実施順序は次のとおりとする。

```text
共通契約・詳細設計
        ↓
Knowledge Store／CLI基盤
        ↓
Knowledge蓄積ワークフロー・記事価値判定ワークフロー
        ↓
オーケストレーション／統合
        ↓
総合的な評価・品質保証
```

評価・品質保証の準備と各レイヤーのテストは、対応する設計・実装と並行して進める。

## 今後の進め方

大分類ごとに、以下を一つずつ議論して承認記録を追加する。

1. 中分類タスク案
2. 各中分類へ分ける理由
3. 中分類間の依存関係
4. 議論による変更点
5. 承認日と承認結果

### 議論順

1. 共通契約・詳細設計
2. Knowledge Store／CLI基盤
3. Knowledge蓄積ワークフロー
4. 記事価値判定ワークフロー
5. オーケストレーション／統合
6. 評価・品質保証

## 決定 002: 共通契約・詳細設計の中分類

### 状態

承認済み（2026-08-10）

### 初回提案（2026-08-10）

1. ドメインモデルと不変条件の詳細化
2. 検索要件と検索契約の設計
3. 更新・履歴契約の設計
4. CLIコマンド・JSON Schemaの設計
5. Skill間成果物契約の設計
6. 物理構成・互換性方針の決定

### 分割理由

#### 1. ドメインモデルと不変条件の詳細化

Knowledge Assertion、Concept、Evidence、Derived Knowledge State、Relation、Scope、Temporal Metadataの意味と制約は、すべての設計・実装の前提になるため。

#### 2. 検索要件と検索契約の設計

本システムはArticle Claimから関連Knowledgeへ到達する検索要件を起点に保存形式を決める方針であり、検索プリミティブ、結果の意味、Search Traceを独立して具体化する必要があるため。

#### 3. 更新・履歴契約の設計

`create`、`attach-evidence`、`revise`、`supersede`、訂正、Derived State再評価は検索とは異なる整合性・履歴保持の規則を持つため。

#### 4. CLIコマンド・JSON Schemaの設計

CLIとCodexの境界は機械的な契約であり、コマンド引数、入出力、エラー、終了コードを固定することで、CLI実装とSkill実装を独立して進められるため。

#### 5. Skill間成果物契約の設計

Codex Subagent間ではMarkdownを使う一方、必須セクションや判定語彙は共有する必要があるため。CLIのJSON契約とは変更理由と利用者が異なるので分離する。

#### 6. 物理構成・互換性方針の決定

物理データモデル、Index実装、Migration、Schema Versioningは論理契約を実現する技術判断であり、論理要件の確定後にまとめて判断する必要があるため。

### 想定する依存関係

```text
ドメインモデルと不変条件
        ├─→ 検索要件と検索契約 ─┐
        └─→ 更新・履歴契約 ─────┼─→ CLIコマンド・JSON Schema
                                 ├─→ Skill間成果物契約
                                 └─→ 物理構成・互換性方針
```

### 未決事項

- 上記6分類を中分類として採用するか
- 「物理構成・互換性方針」を本大分類に含めるか、Knowledge Store／CLI基盤へ移すか

### Issue記載内容との再照合（2026-08-10）

初回提案は、すべてが「Issue内の情報不足」を理由に独立タスクになるわけではない。Issue #1の各節と照合した結果は次のとおり。

#### ドメインモデルと不変条件の詳細化

Issueの10章と49章で概念、責務、不変条件は十分に定義されている。一方、必須・任意属性、識別子、Cardinality、整合性制約、EvidenceからDerived Stateを導出する規則は未定義である。

結論: 概念定義を繰り返すタスクは不要。実装に必要な論理データモデルと整合性制約の詳細設計へ限定して残す。

#### 検索要件と検索契約の設計

Issueの11、12、14〜16、43章で検索方式、検索プリミティブ、Agentic Search、停止原則、Search Traceの目的まで定義されている。一方、CLIの引数、Filter、Pagination、Top-K、結果順序、Scoreの意味などは未定義である。

結論: 検索要件を再定義する独立タスクは不要。未定義部分は「CLI公開契約・JSON Schemaの詳細設計」へ統合する。

#### 更新・履歴契約の設計

Issueの13、25、39章で更新操作と履歴保持の原則は定義されている。一方、各操作の事前・事後条件、冪等性、RevisionとSupersedeの厳密な違い、履歴取得方法は未定義である。

結論: 独立した概念設計ではなく、データ整合性は論理データモデルへ、コマンド挙動はCLI公開契約へ統合する。

#### CLIコマンド・JSON Schemaの設計

Issueの9、12、13章ではコマンド名、責務、JSON利用方針と出力例が定義されているが、完全な入出力Schema、エラー、終了コード、互換性規則は未定義である。51章でも次フェーズの成果物として明示されている。

結論: 情報不足が明確なため、独立タスクとして残す。

#### Skill間成果物契約の設計

Issueの17〜21、27〜36、47章で成果物の主要セクションと判定語彙が定義されている。詳細は不足しているが、Article AnalysisからReading Valueまでの契約は「記事価値判定ワークフロー」、AcquisitionからUpdateまでの契約は「Knowledge蓄積ワークフロー」の内部に閉じる。

結論: 共通契約として独立させず、各ワークフローの中分類へ移す。

#### 物理構成・互換性方針の決定

Issueの48章で、物理DB、Embedding Engine、Index Library、Migration方式は未決定と明記されている。

結論: 必要なタスクだが、共通契約ではなく「Knowledge Store／CLI基盤」の中分類へ移す。

### 再編案

「共通契約・詳細設計」の中分類は、以下の2タスクに絞る。

1. Knowledge論理データモデル・整合性制約の詳細設計
2. CLI公開契約・JSON Schemaの詳細設計

初回提案の残りは、重複を避けるため上記2タスクへ統合するか、責務を持つ別の大分類へ移す。

### 承認結果

以下の2タスクを中分類として承認した。

1. Knowledge論理データモデル・整合性制約を設計する
2. CLI公開契約・JSON Schemaを設計する

検索・更新の詳細は上記2タスクへ統合する。Skill成果物契約と物理構成は、それぞれ責務を持つ別の大分類で扱う。

## 決定 003: 今後のタスク分割手順

### 状態

承認済み（2026-08-10）

### 進行順

大分類内のすべての中分類を先に決めるのではなく、深さ優先で次の順に進める。

```text
中分類を承認
  ↓
その中分類の小分類を提案・検証・承認
  ↓
次の中分類へ移動
```

現在の次の議論対象は「Knowledge論理データモデル・整合性制約を設計する」の小分類とする。

### サブエージェントの利用

今後のタスク分割では、最低限次の2役を利用する。

1. タスク分割担当: 親タスクを一段下のタスクへ分割する。
2. タスク妥当性確認担当: 各候補がIssue内で既に十分定義済みでないか、重複、粒度、配置を独立して確認する。

Issueですでに十分定義された設計判断を、再検討だけのタスクとして作らない。ただし、仕様が十分でも実装・検証という未実施の成果物が必要ならタスクとして残す。

### Skill化

上記手順を再利用可能なSkillとして作成し、以後の分割議論で使用する。

- Skill: `.agents/skills/decompose-tasks/SKILL.md`
- 分割担当プロンプト: `.agents/skills/decompose-tasks/references/decomposition-agent-prompt.md`
- 妥当性確認担当プロンプト: `.agents/skills/decompose-tasks/references/validity-agent-prompt.md`
- 構造検証: 成功

## 決定 004: Knowledge論理データモデル・整合性制約の小分類

### 状態

承認済み（2026-08-10）

### 提案（2026-08-10）

#### Knowledge論理エンティティと属性を定義する

分割した理由: Issueでは各概念の意味は定義済みだが、識別子、属性、型、必須性、列挙値、一意性、構造的な検証規則が未定義であるため。

#### Knowledge論理エンティティ間の関連と参照整合性を定義する

分割した理由: IssueではRelation種別と探索可能な接続は示されているが、方向、Cardinality、重複可否、参照整合性が未定義であり、単一エンティティの属性とは別の制約として扱うため。

#### Evidence整合性とDerived User Knowledge State導出規則を定義する

分割した理由: EvidenceをSource of Truthとする原則は定義済みだが、複数EvidenceからDerived Stateを導出、再評価、無効化する具体的な規則が未定義であるため。

#### Knowledgeの論理ライフサイクルと履歴保持規則を定義する

分割した理由: `create`、`attach-evidence`、`revise`、`supersede`の目的は定義済みだが、状態遷移、同一性、系譜、Temporal Metadataとの整合性が未定義であるため。

#### 検索要件に対する論理モデルの充足性を検証する

分割した理由: 検索要件から保存形式を逆算する原則と検索プリミティブは定義済みだが、今回設計する属性・関連ですべての検索経路が成立することを確認する成果物が未作成であるため。

### サブエージェント検証結果

- 分割担当: 上記5件を候補として提案した。
- 妥当性確認担当: 1〜4を`keep_missing_detail`、5を`keep_unimplemented`と判定した。
- Issue内で定義済み事項の再検討だけを目的とする候補は含まれていない。
- 物理DB・Index・Migration、CLI引数・JSON Schema、CodexによるEvidence抽出、自動評価実装は対象外とする。

### 依存関係

```text
論理エンティティと属性
  ├─→ 関連と参照整合性 ─────────┐
  └─→ Evidence整合性と導出規則 ─┼─→ 論理ライフサイクルと履歴
                                  └─→ 検索要件の充足性検証
```

検索要件の充足性検証は、他の4タスクの完了後に実施する。

### 最終小分類としての再監査（2026-08-10）

小分類より下へ分割しない前提で、分割担当と妥当性確認担当がそれぞれIssue #1へ戻り、原子性、網羅性、重複、配置、定義済み事項の再議論有無を再監査した。

#### 再監査後の提案

1. Knowledge論理スキーマを定義する
2. Knowledge論理関連と参照整合性制約を定義する
3. Evidence由来Stateの整合性・導出契約を定義する
4. Knowledge更新操作の論理事前・事後条件と履歴系譜を定義する
5. 論理モデルの要件トレーサビリティを確認する

#### 主成果物による境界

| 小分類 | 主成果物 | 境界 |
| --- | --- | --- |
| Knowledge論理スキーマ | 論理データ辞書 | レコード内の属性と構造制約 |
| Knowledge論理関連と参照整合性 | 関連制約表 | レコード間の関連と参照制約 |
| Evidence由来Stateの整合性・導出契約 | EvidenceからStateへの導出契約 | 導出根拠、再計算、無効化 |
| Knowledge更新操作と履歴系譜 | 操作別の事前・事後条件表と系譜規則 | 時間を伴う更新の論理効果 |
| 要件トレーサビリティ | Issue要件と設計要素の双方向対応表 | 1〜4の横断的な設計検証 |

#### 追加で既存タスクへ包含する項目

- 論理スキーマ: 知識主体／User境界、Evidence出典・Episode参照、Alias、Scope、Temporal Metadata、監査メタデータ
- 更新操作: 再試行・重複時の論理効果
- トレーサビリティ: 検索、更新、Evidence追跡、モデル関連Invariant、Acceptance Scenario A〜Hのモデル関連条件

独立した6件目は追加しない。

#### 再議論しない固定入力

- Knowledge Assertion中心とAssertionの粒度
- Conceptは検索アンカーであること
- Raw Evidenceと正規化Assertionの分離
- Evidence Strengthの3分類
- Knowledge Assessmentの7状態
- EvidenceがSource of Truthであること
- `create`、`attach-evidence`、`revise`、`supersede`の4操作
- 過去Evidenceを物理削除しないこと
- Issue 49章のInvariant
- Codexが判断し、CLIが保存・取得する責務境界

上記は各小分類の設計入力として使用し、選び直さない。

#### 再監査結果

- 5件すべて必要であり、追加・統合・削除は不要。
- 1〜4はIssue内で不足する詳細設計、5は未作成の設計検証成果物である。
- 各タスクは1つの主成果物で完了判定できるため、これ以上子タスクへ分割しない。

### 承認結果

再監査後の5件を最終小分類として承認した。これらの下に子タスクは作成しない。

次の議論対象を「CLI公開契約・JSON Schemaを設計する」の小分類へ移す。

## 決定 005: CLI公開契約・JSON Schemaの小分類

### 状態

承認済み（2026-08-10）

### 提案（2026-08-10）

#### CLI公開コマンド一覧と動作・呼び出し・引数契約を定義する

分割した理由: Issueでは検索・取得8操作と更新4操作の名前と目的は定義済みだが、公開コマンドの完全な一覧、構文、入力経路、引数、型、必須性、既定値、排他条件が未定義であるため。Index管理、ID管理、履歴保持を公開操作にするか内部動作にするかも明示する。

#### 検索・複数取得結果のコレクション契約を定義する

分割した理由: IssueではTop-Kを設定可能としているが、Filter、Sort、Scoreの意味と比較範囲、同順位処理、件数上限、Pagination、継続情報が未定義であるため。

#### CLIエラー分類と終了コード契約を定義する

分割した理由: 入力不正、対象なし、競合、参照整合性違反、利用不能などを識別するError code、終了コード、stderr、再試行可否がIssueで未定義であるため。

#### CLIのJSON入出力Schemaを定義する

分割した理由: IssueではCLIからCodexへJSONを返す方針と一例だけが示され、各公開コマンドのRequest、成功Response、Error Response、必須・任意項目、ID・日時・列挙値の完全なSchemaが未定義であるため。

#### CLI契約のVersioning・後方互換性規則を定義する

分割した理由: 物理実装変更で論理契約を壊さない原則は定義済みだが、契約Versionの表現、互換・破壊的変更、廃止、未知Field、Version不一致時の挙動が未定義であるため。

#### CLI公開契約の要件トレーサビリティと責務境界遵守を確認する

分割した理由: IssueのCLI責務とCodexとの境界は定義済みだが、全公開操作について引数、出力、Error、集合規則、Versionが揃い、意味判断がCLIへ混入していないことを確認する成果物が未作成であるため。

### 主成果物による境界

| 小分類 | 主成果物 | 境界 |
| --- | --- | --- |
| 公開コマンド・呼び出し契約 | コマンド台帳 | 操作の公開範囲とCLI上の呼び出し表現 |
| 結果コレクション契約 | 集合取得規則と適用表 | 複数結果の意味と継続方法 |
| Error・終了コード契約 | Error台帳と終了コード対応表 | 失敗の意味とProcess境界 |
| JSON入出力Schema | 機械検証可能なSchema一式 | 現行契約のJSON表現 |
| Versioning・互換性 | 互換性方針 | 将来の契約変更規則 |
| 要件トレーサビリティ | Issue要件とCLI契約の双方向対応表 | 上記成果物の横断検証 |

### 再議論しない固定入力

- CLIとCodexの機械的境界ではJSONを使用すること
- Issue 12章の検索・取得プリミティブの意味
- Issue 13章の4更新操作の意味
- Codexが意味判断し、CLIが決定論的に保存・取得する責務境界
- `no_evidence`等のKnowledge AssessmentをCLIが判断しないこと
- 物理DB、Index、Embedding、Migration実装はL2で扱うこと

### サブエージェント検証結果

- 分割担当と妥当性確認担当は、6件すべてを維持すべきと判断した。
- 1〜5はIssue内で不足する契約詳細、6は未作成の設計検証成果物である。
- 独立した7件目は不要。
- 各タスクは1つの主成果物を持つため、これ以上子タスクへ分割しない。

### 依存関係

```text
公開コマンド一覧・動作・呼び出し契約
  ├─→ 結果コレクション契約 ─┐
  └─→ Error・終了コード契約 ─┼─→ JSON入出力Schema
                               ↓
                       Versioning・互換性
                               ↓
                  要件トレーサビリティ確認
```

### 承認結果

提案した6件を最終小分類として承認した。各タスクは1つの主成果物で完了判定し、これらの下に子タスクは作成しない。

次の議論対象を`L2`「Knowledge Store／CLI基盤」の中分類へ移す。

## 決定 006: Knowledge Store／CLI基盤の中分類

### 状態

承認済み（2026-08-10）

### 初回分割案と妥当性監査（2026-08-10）

分割担当は、物理データモデル、永続化、4種類の検索基盤、Index整合性管理、CLI実行境界の8件を提案した。

妥当性確認担当は、8件の内容自体は必要だが、中分類としては細かすぎると判定した。物理データモデルと永続化は同じ正本データ基盤、個別検索方式とIndex整合性管理は同じ検索サブシステムを構成するため、中分類では3件へ統合する。統合前の項目は、各中分類を承認した後の小分類候補として保持する。

### 監査後の提案（2026-08-10）

#### Knowledge Storeの物理永続化・更新基盤を設計・実装する

分割した理由: IssueではKnowledge論理モデル、4更新操作、Evidence起点、履歴保持が定義済みだが、物理DB、Schema、Transaction、Migration、履歴の物理表現、Derived Stateの再計算・無効化と実装が未決定・未実施であるため。正本データの整合性を所有する単位としてまとめる。

#### 複合検索・Index管理基盤を設計・実装する

分割した理由: IssueではLexical、Semantic、Relationの3系統とTemporal検索が必要と定義済みだが、各Engine・Library、物理Index、検索実装、同期、再構築、障害復旧が未決定・未実施であるため。個別検索方式と派生Indexのライフサイクルを一体で管理する単位としてまとめる。

#### Knowledge CLI実行基盤とJSON境界を実装する

分割した理由: L1でCLI公開コマンド、引数、JSON Schema、Error、互換性を設計するが、それらを永続化・検索基盤へ接続し、Codexから利用可能な実行形式として提供する処理は未実装であるため。入力検証、Dispatch、JSON出力、終了コード、設定・初期化、実行Artifactを担う境界として分離する。

### 主成果物による境界

| 中分類 | 主成果物・到達状態 | 境界 |
| --- | --- | --- |
| 物理永続化・更新基盤 | Migration可能で履歴・整合性を保持するKnowledge Store | 正本データの保存、取得、更新 |
| 複合検索・Index管理基盤 | 各検索方式と再構築可能なIndex群 | 派生検索構造と候補取得 |
| CLI実行基盤・JSON境界 | L1契約に適合する実行可能CLI | Codexと内部基盤のProcess境界 |

### 再議論しない固定入力

- Knowledge Assertion中心、EvidenceがSource of Truthであること
- L1-M1で定義する論理モデル、関連、導出、更新・履歴契約
- L1-M2で定義する公開コマンド、JSON Schema、Error、互換性契約
- Lexical、Semantic、Relationの3系統を最低限持つこと
- Temporal候補を検索可能にすること
- Codexが意味判断し、CLIは決定論的に保存・取得すること
- 類似度、矛盾候補、新旧候補からCLIがKnowledge Assessmentを決定しないこと

### サブエージェント検証結果

- 分割担当の8件案を、妥当性確認担当の監査により中分類3件へ統合した。
- 3件はいずれもIssueで要求されているが、詳細が未確定または実装が未実施である。
- 元の8件に含まれるMigration、個別検索方式、Index整合性管理は削除せず、小分類の議論対象へ移す。
- L2内のUnit・Integration・L1契約適合確認は各実装の完了条件に含め、横断的な評価仕様、Fixture、Acceptance Scenario検証はL6へ配置する。

### 依存関係

```text
L1 論理モデル・CLI公開契約
              ↓
物理永続化・更新基盤
              ↓
複合検索・Index管理基盤
              ↓
CLI実行基盤・JSON境界
```

### 承認結果

監査後の3件を中分類として承認した。初回8件案の各要素は、削除せず対応する中分類の小分類候補として扱う。

次の議論対象を`L2-M1`「Knowledge Storeの物理永続化・更新基盤を設計・実装する」の小分類へ移す。

## 決定 007: Knowledge Store物理永続化・更新基盤の小分類

### 状態

承認済み（2026-08-11）

### 提案（2026-08-10）

#### 正本DBを選定しBaseline物理Schemaを実装する

分割した理由: Issueでは物理DBと物理Schemaが未決定であり、L1の論理レコード、User境界、Evidence、Relation、Scope、Temporal Metadata、監査情報、履歴系譜を正本データ構造へ写像する成果物が未作成であるため。DB機能とSchema制約は不可分なので1件にまとめる。

#### Schema Version管理とMigration実行基盤を実装する

分割した理由: IssueではMigration方式が未決定であり、既存の正本データと履歴を失わずSchemaを更新する仕組みが未実装であるため。Baseline Schemaの内容とは分け、Version台帳と適用Lifecycleを担当する。

#### DB接続・Transaction・ID・Record変換を担う永続化Coreを実装する

分割した理由: Issueでは決定論的保存とID管理がCLI責務として定義済みだが、接続Lifecycle、Transaction、競合制御、ID生成、論理Recordと物理Rowの変換が未決定・未実装であるため。すべての取得・更新が共有する低水準基盤として分離する。

#### Evidence由来Stateの導出・再評価整合性を実装する

分割した理由: EvidenceをSource of Truthとする原則は決定済みだが、Derived Stateを永続化・Cacheするか都度導出するか、および選択した方式で不整合を防ぐ実装は未決定・未実施であるため。更新と取得を横断するInvariantとして独立させる。

#### `create`・`attach-evidence`の追加型更新処理を実装する

分割した理由: 両操作の意味と論理効果はIssueとL1で定義済みだが、Assertion一式の新規保存と既存AssertionへのEvidence追加を原子的に行う処理が未実装であるため。既存履歴を置換しない追加型Transactionとしてまとめる。

#### `revise`・`supersede`の非破壊更新と履歴系譜を実装する

分割した理由: 両操作と過去Evidenceを削除しない原則は定義済みだが、旧Recordを保持した新Version・系譜の保存、競合・再試行時の原子的処理が未実装であるため。履歴を伴う非破壊更新として追加型操作から分離する。

#### Knowledge詳細・Evidence・履歴系譜の取得Repositoryを実装する

分割した理由: `get`と`get-evidence`の公開契約はL1で扱うが、正本DBから論理Aggregateを復元し、Evidence、Temporal Metadata、Revision・Supersede系譜まで追跡する取得処理は未実装であるため。探索検索ではない主キー・参照起点の取得として分離する。

### 主成果物による境界

| 小分類 | 主成果物・到達状態 | 境界 |
| --- | --- | --- |
| 正本DB・Baseline Schema | 利用可能な初期正本DB | 物理配置とDB制約 |
| Schema Version・Migration | 再実行可能で原子的なMigration基盤 | 物理Schemaの変更Lifecycle |
| 永続化Core | 接続・Transaction・ID・Record変換層 | 低水準DB共通処理 |
| Evidence由来State整合性 | 選択方式に応じた導出・無効化・再評価処理 | EvidenceとDerived Stateの整合 |
| create・attach-evidence | 追加型更新Repository | 新規作成とEvidence追加 |
| revise・supersede | 非破壊更新Repositoryと系譜 | 修正・置換と履歴保持 |
| 詳細・Evidence・履歴取得 | 論理Aggregate取得Repository | 主キー・参照起点の復元と追跡 |

### 再議論しない固定入力

- L1-M1で定義する論理Schema、関連、参照整合性、State導出、4更新操作、履歴系譜
- EvidenceがSource of Truthであり、過去Evidenceを物理削除しないこと
- Correctionは独立した第5更新操作ではなく、既存操作とState再評価で扱うこと
- L1-M2で定義するCLI公開契約と契約Version
- CLIが意味的同一性、Knowledge Assessment、矛盾、新旧の正しさを判断しないこと

### サブエージェント検証結果

- 分割担当と妥当性確認担当は、7件すべてを最終小分類として維持すべきと判断した。
- 物理DB、Migration方式、Transaction、ID、Derived Stateの物理方式はIssue内で不足する詳細である。
- 4更新操作と取得は意味が定義済みだが、実装が未実施であるためタスクとして必要である。
- User境界、Relation正本、Temporal Metadata、監査情報、Correction、再試行・重複、競合・Crash時の原子性は7件の完了条件へ包含し、独立した8件目は追加しない。
- Backup、暗号化、削除、一般的な性能最適化はIssueに根拠がないため追加しない。
- 各小分類固有のUnit・Integration確認は完了条件に含め、横断Fixture、CLIテスト仕様、Acceptance Scenario評価はL6へ配置する。

### 依存関係

```text
正本DB・Baseline物理Schema
            ↓
Schema Version・Migration
            ↓
永続化Core
            ↓
Evidence由来State導出・再評価
       ├─→ create・attach-evidence ─┐
       └─→ revise・supersede ──────┼─→ 詳細・Evidence・履歴取得
```

更新2件は並行実装可能だが、最終的な取得確認では両方の保存結果とState整合性を縦断して確認する。

### 技術選定・詳細設計の欠落指摘と再監査（2026-08-11）

ユーザーから、親タスクが「設計・実装」であるにもかかわらず実装中心の小分類であり、技術選定も独立成果物になっていないとの指摘を受けた。

旧案では、DB選定をBaseline Schema実装へ混在させ、Driver、Data Access技術、Migration Tool等の共通技術選定を明示していなかった。また、物理モデル、Transaction、競合制御、Repository境界等の詳細設計を実装タスクへ暗黙に含めていた。このため、設計判断を実装前に承認・追跡できない分割になっていた。

分割担当と妥当性確認担当による再監査の結果、技術選定1件、詳細設計2件、実装7件へ修正する。

### 修正版提案（2026-08-11）

#### 正本永続化の技術スタックを選定しADRを確定する

分割した理由: Issue 48章で物理DBとMigration方式が未決定であり、DB、Driver／Data Access技術、Migration Tool／方式、必要な外部ID生成技術は全実装へ影響する共通判断だからである。比較基準、採否理由、制約を実装前に独立して承認可能にする。

#### 正本データの物理モデル・履歴・Schema Evolutionを詳細設計する

分割した理由: L1の論理契約を選定技術へ写像し、表、型、Key、制約、User境界、Evidence・Relation正本、Temporal・監査情報、非破壊履歴、Derived State保存方針、Schema Versionを確定する成果物が未作成だからである。静的な保存表現を1つの詳細設計としてまとめる。

#### 永続化操作・Transaction・Repository境界を詳細設計する

分割した理由: 4更新操作の意味は定義済みだが、Transaction、競合・再試行・冪等性、Crash時原子性、ID採番Lifecycle、State再評価Hook、Repository・Record変換・取得Aggregateの実行時設計が未作成だからである。物理データ構造とは異なる動的な操作設計として分離する。

#### Baseline物理SchemaとDB制約をMigrationとして実装する

分割した理由: 詳細設計で確定した正本構造を、初回適用可能なDDL／Migrationとして実DBへ実装する成果物が未作成だからである。

#### Schema Version管理とMigration実行基盤を実装する

分割した理由: Version台帳、初期化、Upgrade、二重適用防止、失敗・Crash時の原子性を持つSchema変更Lifecycleが未実装だからである。

#### DB接続・Transaction・ID・Record変換Coreを実装する

分割した理由: 詳細設計に従い、全Repositoryが共有する接続Lifecycle、Transaction primitive、ID生成、Record Mapper、内部Error処理を構築する必要があるため。

#### Evidence由来Stateの導出・再評価経路を実装する

分割した理由: 詳細設計で選択した方式に従い、EvidenceをSource of TruthとしてStateを導出・再評価し、保存・Cache方式では必要な無効化を行う共通経路が未実装だからである。

#### `create`・`attach-evidence`の追加型更新を実装する

分割した理由: L1と詳細設計で確定した契約に従い、新規Knowledge作成とEvidence追加を原子的に行う追加型更新が未実装だからである。

#### `revise`・`supersede`と非破壊履歴系譜を実装する

分割した理由: L1と詳細設計で確定した契約に従い、旧Record・Evidenceを保持して新Versionと系譜を作る非破壊更新が未実装だからである。

#### Knowledge詳細・Evidence・更新履歴の取得Repositoryを実装する

分割した理由: 正本DBからKnowledge Aggregate、Evidence、Relation正本、Temporal・監査情報、Revision・Supersede履歴を復元する取得処理が未実装だからである。

### 修正版の成果物境界

| 小分類 | 主成果物・到達状態 | 境界 |
| --- | --- | --- |
| 技術スタック選定 | 永続化技術ADR | 採用する技術要素と選定根拠 |
| 物理モデル詳細設計 | 物理ストレージ詳細設計書 | 静的な保存表現とSchema Evolution |
| 永続化操作詳細設計 | 実行時永続化設計書 | 動的な操作、Transaction、Repository境界 |
| Baseline Schema | 初回DDL／Migration | 初期物理構造の実装 |
| Migration実行基盤 | Version管理されたMigration Runner | Schema変更の実行Lifecycle |
| 永続化Core | DB共通実装 | 接続・Transaction・ID・変換 |
| Evidence由来State | State導出・再評価経路 | EvidenceとStateの実行時整合 |
| 追加型更新 | create・attach-evidence Repository | 追加型Transaction |
| 非破壊更新 | revise・supersede Repository | Version生成と履歴系譜 |
| 取得Repository | 論理Aggregate取得処理 | 詳細・Evidence・履歴の復元 |

### 再監査結果

- 修正版10件はすべて必要であり、統合・削除・追加は不要である。
- 技術選定はDB、Driver／Data Access、Migration、必要な外部ID生成技術に限定し、Repository／Transaction境界やID採番Lifecycleは詳細設計で扱う。
- 物理モデル設計は静的構造、永続化操作設計は動的振る舞いを扱うため重複しない。
- Embedding Engine・Index LibraryはL2-M2、CLI Framework・配布・JSON境界はL2-M3へ配置する。
- 旧7件案は承認前の案として履歴に残し、Task Mapの現行案は修正版10件へ置き換える。

### 修正版の依存関係

```text
技術選定ADR
    ↓
物理モデル・履歴・Schema Evolution設計
    ├─→ Baseline Schema ─→ Migration実行基盤 ─┐
    └─→ 永続化操作・Transaction・Repository設計 ├─→ 永続化Core
                                                   ↓
                                       Evidence由来State経路
                                        ├─→ 追加型更新 ─┐
                                        └─→ 非破壊更新 ┼─→ 詳細・Evidence・履歴取得
```

### 承認結果（2026-08-11）

ユーザー承認により、修正版10件を`L2-M1`の小分類として確定する。

この承認は現時点の原典・階層・既知の依存関係に基づく。後続タスクの分割によって、責務重複、欠落、粒度不一致、階層配置の不整合が判明した場合は、承認済みであっても履歴を保持したまま修正案を提示し、再承認を受ける。

次の議論対象を`L2-M2`「複合検索・Index管理基盤を設計・実装する」の小分類へ移す。


## 決定 008: 承認済みタスクの継続的な整合性再評価

### 状態

承認済み（2026-08-11）

### 決定内容

タスクの承認は、その時点で得られている原典、親子関係、依存関係に基づく構造確定として扱い、将来にわたる不変条件とはしない。

後続タスクの分割時は、新規候補だけでなく既存の大分類・中分類・小分類も見回し、次を横断確認する。

- 親の名称・主成果物と子タスクの総和が一致するか
- 責務の欠落または複数タスクによる二重所有がないか
- 同一階層の粒度が揃っているか
- 大分類、中分類、小分類の配置が適切か
- 必要な技術選定、詳細設計、実装、検証が欠落していないか
- 原典の設計原則と依存順序が逆転していないか
- 隣接タスク間の受け渡し境界と横断要件の所有者が明確か

矛盾を発見した場合は、承認済みであることを維持理由にせず、旧決定を履歴に残したまま、対象、問題、修正案、理由、影響範囲、再承認単位を提示する。確定Task Mapの名称・親子関係・責務は、再承認後に変更する。

通常の新規分割は一度に一階層だけ行う。一方、整合性監査は全階層を対象にできるが、監査だけで複数階層の新構造を同時確定しない。

### Skillへの反映

`decompose-tasks` Skillと、分割担当・妥当性確認担当の両プロンプトへ次を追加した。

- 承認済みタスクを含む全階層の横断監査
- 承認状態を妥当性の根拠にしない規則
- 親子成果物、粒度、工程、依存、横断責務の確認
- 既存構造の修正候補に必要な記録項目
- 修正案の再承認前は確定構造を書き換えない規則

構造検証は成功した。独立サブエージェントによる利用テストでも、承認済みL1の親子不一致を検出し、旧承認を保持した再承認案を提示できることを確認した。

## 決定 009: 全タスク階層の横断整合性監査

### 状態

承認済み（2026-08-11）

### 監査方法

分割担当がIssue #1、全承認記録、Task Mapを横断監査し、別の妥当性確認担当がIssue原文へ戻って独立検証した。過去の承認は妥当性の根拠に使用していない。

### 修正案1: L1の名称・主成果物を子タスクへ一致させる

#### タスク名

`L1 Knowledge論理モデル・CLI公開契約を設計する`

#### 分割した理由

現在の子タスクはKnowledge論理モデルとCLI公開契約だけであり、Skill間成果物契約と物理詳細設計は別の大分類へ移動済みである。現名称「共通契約・詳細設計」は範囲を広く見せ、L2〜L5の設計責務と重複するため、親の名称と主成果物だけを実際の子の総和へ限定する。

影響範囲はL1の名称、主成果物、全体依存図の表記のみとし、L1-M1・L1-M2の内容は変更しない。

### 修正案2: L3〜L5を個別SkillとParent Orchestrationの責務へ再編する

#### タスク名

- `L3 Knowledge蓄積Skillsを設計・実装する`
- `L4 記事価値判定Skillsを設計・実装する`
- `L5 Parent Orchestration Skillと全体Workflowを設計・実装する`

#### 分割した理由

現構造ではL3・L4がWorkflowを所有し、L5もSkill間制御、再調査、終了判定、E2E統合を所有するため、Workflow制御が二重化する。Issue 8、24、37〜39、46〜47章に合わせ、L3・L4は個別Skill固有の判断・成果物契約、L5はSkill起動、成果物受け渡し、再調査、Budget、終了判定、両Workflowの接続を所有する。

Search Traceは独立大分類を追加せず、L4がQuery・取得結果・継続／停止理由を生成し、L5がRun・Claim・再調査Cycle・Budgetと相関し、L6が評価時の失敗原因診断に使用する。

### 修正案3: L2の設計工程と正本・Index・CLI境界を補正する

#### タスク名

- `L2-M1-S1 正本永続化の技術スタックを選定しADRを確定する`
- `L2-M1-S2 正本データの物理モデル・履歴・Schema Evolutionを詳細設計する`
- `L2-M1-S3 永続化操作・Transaction・Repository境界を詳細設計する`
- `L2-M2 複合検索・Index管理基盤を設計・実装する`
- `L2-M3 Knowledge CLI実行基盤とJSON境界を詳細設計・実装する`

#### 分割した理由

L1-M2はCLI公開契約を扱うだけで、CLI Framework、Handler、Validation、Dispatch、設定、内部Error変換等の内部詳細設計を所有しないため、L2-M3へ詳細設計を明示する必要がある。

また、正本更新から派生Indexへ渡す境界が未定義であるため、L2-M1-S3に「未Commit変更を公開せず、Commit済みの対象ID・Revision等をM2が認識または再照合でき、正本から再構築できる」という技術非依存Invariantを追加する。通知、Polling、Outbox、同一Transaction等の方式選定はL2-M2に残す。

L2-M1-S3の直接依存には、内部Transaction・Errorが公開操作を実現できるよう、L1-M2-S1とL1-M2-S3を追加する。L2-M1-S1・S2は、L4の検索探索要求とL2-M2の同期・再構築要件に対するReview Gateを通してから確定する。

Temporal検索は必須だが、独立した物理Temporal Indexを必須とはしない。`search-contradictions`についてもAI的矛盾判定Engineを追加せず、既存Relation・Concept等から候補を返す責務に限定する。

### 修正案4: 検索起点の工程依存とテスト所有境界を明確にする

#### タスク名

`L6 評価設計・横断品質保証を実装する`

#### 分割した理由

現在の`L1 → L2 → L4`という単純な順序では、Knowledge Search Skillの具体的な探索要求より先に物理Indexを確定し、「Knowledge Storage is designed backward from Retrieval」というIssue 6.7・49章のInvariantと逆転するためである。

L4全体をL2-M2の前提にすると循環するため、L4内のKnowledge Search探索要求・Query Journey・停止条件・必要取得情報の詳細設計だけをL2-M2設計へ先行させる。L2-M2実装後にL4のKnowledge Search Skill実装を接続する。

各実装タスクは固有のUnit、Contract、Component Integrationテストを完了条件として所有する。L6は横断テスト方針、共有Fixture・Dataset、評価Harness、Acceptance Scenario、E2E・回帰実行、診断Reportを所有する。L6の評価仕様・Dataset設計は各設計と並行し、最終E2E実行だけをL5統合後に行う。

### 提案する工程依存

```text
L1 論理モデル・CLI公開契約
  ├─→ L2-M1の設計
  ├─→ L3・L4のSkill詳細設計
  └─→ L6の評価仕様・Dataset設計

L4 Knowledge Searchの探索要求設計
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

### 再承認単位

1. L1の名称・成果物補正
2. L3〜L5の責務再編とSearch Trace配置
3. L2-M1-S1〜S3、L2-M2、L2-M3の境界・工程補正
4. 全体依存とL6の評価工程・テスト所有補正

上記4組は相互に関係するが、変更範囲を追跡できるよう個別に承認する。承認前は既存の名称・親子関係を確定構造として維持し、対象タスクをTask Map上で「再検討中」と表示する。

### 承認結果（2026-08-11）

ユーザー承認により、4組の修正案をすべて確定する。Task Mapへ新しい名称、責務、設計入力、工程依存、テスト所有境界を反映し、対象タスクを「承認済み」へ戻す。

次の議論対象を`L2-M2`「複合検索・Index管理基盤を設計・実装する」の小分類へ移す。

## 決定 010: 複合検索・Index管理基盤の小分類

### 状態

承認済み（2026-08-11）

### 分割・監査方法

分割担当がIssue #1、承認記録、Task Mapを全文確認して小分類案を作成し、別の妥当性確認担当がIssue原文へ戻って独立監査した。Issueで決定済みの検索能力は固定入力とし、未決定の技術・物理設計と未実装成果物だけを対象にした。

妥当性確認により、当初9件へ「候補検索・Index Lifecycle要件と正本設計Review Gate基準」を追加した10件を最終提案とする。この追加は既存のL2-M1とL2-M2の循環依存を解消するために必要であり、単なる資料作成タスクではない。

### L2-M2-S1

#### タスク名

`候補検索・Index Lifecycle要件と正本設計Review Gate基準を定義する`

#### 分割した理由

Issueは検索能力と正本原則を定めているが、検索対象、鮮度、Commit済み正本との整合、再構築可能性を正本設計へ渡す技術非依存のGate基準は未定義である。これを技術選定より前に独立させ、L2-M1がL2-M2の完成を待つ循環を解消する。

### L2-M2-S2

#### タスク名

`検索・Index技術スタックを選定しADRを確定する`

#### 分割した理由

Embedding Engine・Model、Index Library、Relation系構造、正本変更との同期方式はIssueで未決定である。比較軸、採否理由、高リスク項目の検証を一つのADRへ確定し、後続設計で選定を繰り返さないため独立させる。正本DBとMigration ToolはL2-M1の決定を再選定しない。

### L2-M2-S3

#### タスク名

`候補検索アーキテクチャと方式別Indexを詳細設計する`

#### 分割した理由

公開検索契約とKnowledge Search探索要求を、内部Component、方式別Index、共通Candidate、Filter、Score由来、安定順序、Pagination、正本照合境界へ写像する設計がIssueにはない。方式別実装へ設計判断を残さないため独立させる。

### L2-M2-S4

#### タスク名

`正本同期・Index Version・再構築／障害回復を詳細設計する`

#### 分割した理由

候補を読むQuery Planeと、Commit済み変更の反映、Version判定、Checkpoint、部分失敗、世代切替、Rollbackを扱うIndex Lifecycleは異なる責務である。全Indexへ一貫する復旧設計を一つの成果物として確定するため分離する。

### L2-M2-S5

#### タスク名

`共通候補検索Coreと正本照合を実装する`

#### 分割した理由

内部Request・Candidate・Provider、共通Filter、安定順序、Pagination、件数制限、候補と正本の照合は全検索方式で共有する。方式別実装とCLI配線の重複を防ぐため共通Coreにまとめる。正本Aggregateの復元はL2-M1-S10を利用し、再実装しない。

### L2-M2-S6

#### タスク名

`Lexical Indexと文字列候補検索を実装する`

#### 分割した理由

完全一致、識別子、Alias、用語検索はEmbeddingに依存しない固有のIndex・正規化・一致規則を持ち、独立して実装・テスト・完了判定できるため分離する。

### L2-M2-S7

#### タスク名

`Embedding生成・Semantic Indexと意味候補検索を実装する`

#### 分割した理由

Embedding生成、Model Version、Vector Index、再Embeddingは固有の実行経路と障害条件を持つ。Lexical検索から独立して完成させ、Knowledgeの既知・正誤という意味判断は含めない。

### L2-M2-S8

#### タスク名

`Concept／Relation Indexと構造・Contradiction候補検索を実装する`

#### 分割した理由

Concept検索、Relation探索、Contradiction候補は同じ構造Indexを利用できる。Contradictionを別タスクにするとAI的矛盾推論を実装すると誤解されるため、明示済みRelation・Conceptから候補を返す一つの実装へ統合する。

### L2-M2-S9

#### タスク名

`Temporal候補検索を実装する`

#### 分割した理由

Temporal検索はIssue上の独立した公開検索能力であり、時間・Version境界に固有のテストと完了条件を持つ。独立物理Temporal Indexは必須とせず、正本Fieldまたは選定済みIndexで候補を取得し、最新・正誤・陳腐化の意味判断は含めない。

### L2-M2-S10

#### タスク名

`正本同期・Index Version管理・再構築／障害回復を実装する`

#### 分割した理由

各Providerへの冪等反映、Version・Checkpoint、完全再構築、検証後の世代切替、部分障害からの再開は同じLifecycle状態を扱う。個別Indexへ分散させず、一つの統合実装として整合性と復旧を完了判定する。

### 依存関係案

```text
L1-M1・L1-M2・L4 Knowledge Search探索要求
                    ↓
L2-M2-S1 要件・Review Gate基準
                    ↓
           L2-M1-S1〜S3
                    ↓
L2-M2-S2 技術選定ADR → S3 検索・Index設計 → S4 Lifecycle設計
                                                ↓
                                    S5 共通検索Core・正本照合
                                  ┌─────┬─────┬─────┐
                                  S6    S7    S8    S9
                                  └─────┴─────┴─────┘
                                                ↓
                                      S10 Lifecycle実装
```

S6〜S9は並行実装でき、TemporalはConcept／Relation実装へ固定依存しない。L4のKnowledge Search探索要求設計は暫定前方参照とし、L4分割後に具体的なTask IDへ置換して再監査する。

### Issueで固定する入力

- Lexical、Semantic、Concept、Relation、Contradiction候補、Temporalの各検索プリミティブ
- CLIは候補を決定論的に返し、意味判断はCodexが行う責務境界
- Contradictionは既存Relation・Concept等から候補を返し、新しい矛盾を推論しない
- Temporal検索は必要だが、独立した物理Temporal Indexは必須ではない
- Evidenceを含むKnowledge Storeが正本であり、派生Indexは正本から再構築可能である

親成果物の「4種」は公開操作数と誤解されるため、「Issue 12章の全検索プリミティブを支える候補検索と、正本から再構築可能な派生Index」へ明確化する。

### 対象外

- 正本DB、Schema、Migration、Transaction、変更境界の実装: L2-M1
- 公開CLI引数、JSON Schema、Error、終了Code: L1-M2
- CLI Handler、設定、配布: L2-M3
- Query生成、探索継続・停止、意味評価、Search Trace生成: L4
- Trace相関と再調査制御: L5
- 共有Fixture、評価Harness、Acceptance、E2E: L6
- Hybrid融合、Reranking、AI的矛盾判定: IssueまたはL1で要求されるまで追加しない

### 全階層監査で検出した修正案

#### 修正案1

##### タスク名

`L2-M1-S1・S2のReview Gate依存をL2-M2-S1へ変更する`

##### 分割した理由

決定009の「L2-M1-S1・S2がL2-M2の同期・再構築要件を待つ」と「L2-M2がL2-M1へ依存する」を同時に適用すると循環する。新S1が技術非依存の要件とGate基準を先に定義し、その後に正本技術・物理設計、検索技術選定へ進む順序へ改める。

#### 修正案2

##### タスク名

`L2-M1-S10の直接依存をL2-M1-S6・S7へ縮小する`

##### 分割した理由

Knowledge取得Repositoryは物理CoreとEvidence由来Stateを入力にFixtureで独立実装でき、create・revise・supersede実装の完成を待つ必然性がない。L2-M2-S5の正本照合を不要に遅延させないため直接依存を縮小する。

### 妥当性監査結果

- 10件はいずれもIssueの未決定詳細または未実装成果物であり、Issue決定事項の再議論ではない。
- 各タスクはADR、要件・Gate基準、詳細設計、共通Core、方式別Provider、Lifecycle実装のいずれか一つを主成果物とし、これ以上の子分割を必要としない。
- Concept・Relation・Contradictionの統合と、同期・Version・再構築・障害回復の統合により過剰分割を避けた。
- 上記2件以外に、承認済みL1、L2-M3、L3〜L6との新たな責務矛盾は確認されなかった。

### 承認対象

1. L2-M2小分類10件
2. L2-M2親成果物の表現修正
3. L2-M1-S1・S2のReview Gate依存修正
4. L2-M1-S10の直接依存修正

### 一部承認結果（2026-08-11）

ユーザーは小分類10件をまだ承認せず、横断差異として提示した次の3点だけを承認した。

1. L2-M2親成果物を「Issue 12章の全検索プリミティブを支える候補検索と、正本から再構築可能な派生Index」へ明確化する。
2. L2-M1-S1・S2のReview Gate依存をL2-M2-S1へ変更する。
3. L2-M1-S10の直接依存をL2-M1-S6・S7へ縮小する。

上記をTask Mapへ確定反映し、L2-M1-S1・S2・S10とL2-M2親を「承認済み」へ戻した。L2-M2-S1〜S10は、各項目の必要性だけでなく、隣接候補へ統合せずL2-M2直下の独立タスクにする理由を再提示するまで「議論中」を維持する。

### 分割理由の不足指摘と再提示（2026-08-11）

ユーザーから、初回説明は各項目の必要性に留まり、「なぜL2-M2直下の別タスクへ分けるのか」の説明が不足しているとの指摘を受けた。分割担当と妥当性確認担当が、親責務、隣接成果物との差、統合時の問題を再監査し、次の理由へ具体化した。

#### 候補検索・Index Lifecycle要件と正本設計Review Gate基準を定義する

分割した理由：  
検索・鮮度・正本整合・再構築の技術非依存条件を先に確定し、L2-M1の正本設計と後続の技術選定が同じ基準を入力にできるよう独立させる。技術選定へ統合すると、選択した技術に要件を合わせる逆順となり、L2-M1との循環依存が再発する。

#### 検索・Index技術スタックを選定しADRを確定する

分割した理由：  
全検索方式とLifecycle設計を拘束する採用技術と採否理由だけを、一つの意思決定成果物として確定する。詳細設計へ統合すると、技術比較と設計変更を分離できず、採否理由を単独で承認・変更できない。

#### 候補検索アーキテクチャと方式別Indexを詳細設計する

分割した理由：  
検索要求を、候補を読み出すQuery PlaneのComponent、Index表現、Candidate契約へ写像する責務を持つ。正本同期・復旧を扱うLifecycle設計へ統合すると、検索契約変更と障害・Version方針変更という異なる変更理由が混在する。

#### 正本同期・Index Version・再構築／障害回復を詳細設計する

分割した理由：  
全Provider共通の同期、Version、Checkpoint、再構築、障害回復規則を、Query Planeと実装から独立して事前レビュー可能にする。Query設計へ統合すると責務が過大になり、Lifecycle実装へ統合すると復旧規則を実装前に確定できない。

#### 共通候補検索Coreと正本照合を実装する

分割した理由：  
Request、Candidate、Filter、順序、Pagination、正本照合という方式非依存処理を一箇所へ集約する。方式別Providerへ分散すると同じ契約が重複実装され、共通不具合と方式固有不具合を切り分けられない。

#### Lexical Indexと文字列候補検索を実装する

分割した理由：  
文字列正規化、完全一致、Alias、Identifier検索という固有規則を一つの完了単位にする。Semantic・構造・Temporal検索と統合すると、異なる技術の完成を同時に要求し、独立テストと並行実装ができない。

#### Embedding生成・Semantic Indexと意味候補検索を実装する

分割した理由：  
Embedding生成、Model Version、Vector Index、再Embeddingという固有の実行経路と障害条件を所有する。決定論的な文字列・構造検索へ統合すると、Model変更が無関係な検索実装まで巻き込む。

#### Concept／Relation Indexと構造・Contradiction候補検索を実装する

分割した理由：  
Concept、Relation、Contradiction候補を同じ構造IndexとTraversalの完了単位としてまとめる一方、文字列・Vector検索から分離する。Semantic検索へ統合すると明示的Relationと意味類似が混同され、Contradictionを独立させるとAI推論責務と誤解される。

#### Temporal候補検索を実装する

分割した理由：  
時間、Version、有効区間境界という固有規則を、物理Temporal IndexやRelation実装を必須にせず独立検証する。Relation検索へ統合するとTemporal検索がRelation Indexへ依存するように見え、時間境界固有の完了判定が埋没する。

#### 正本同期・Index Version管理・再構築／障害回復を実装する

分割した理由：  
全Providerにまたがる一つのLifecycle状態とCoordinatorを所有し、同期、Checkpoint、世代切替、復旧を方式別実装へ分散させない。各Providerへ統合すると状態所有が重複し、Lifecycle設計へ統合すると設計と横断実装を一つの終端タスクへ詰め込み過ぎる。

再監査の結果、10件はいずれも成果物、変更理由、完了判定が分離できるため維持する。S6〜S9の粒度は公開コマンド数ではなく、共有する検索機構の系統で揃える。

依存関係は、S3からS4とS5へ分岐し、S6〜S9がS4・S5の両方を入力として並行実装され、S10が全Providerを統合する形へ精緻化する。

### 最終承認結果（2026-08-11）

ユーザー承認により、L2-M2-S1〜S10の小分類10件を確定した。先に承認済みの親成果物表現、L2-M1-S1・S2のReview Gate依存、L2-M1-S10の直接依存と合わせ、L2-M2の小分類議論を完了する。

## 決定 011: L2設計先行Gateとgit worktree実行計画

### 状態

承認済み（2026-08-11）

### 指摘

L2-M1とL2-M2はそれぞれ設計と実装を含むが、現行DAGではL2-M1-S2完了直後にL2-M1-S4のBaseline Schema実装を開始できる。Transaction・Commit済み変更境界を定めるL2-M1-S3と、Index技術・物理表現・Lifecycleを定めるL2-M2-S2〜S4の確定前に本番Schemaを作るため、後続設計による手戻り経路が残っていた。

分割担当が依存DAGとworktree競合域を監査し、別の妥当性確認担当が過剰待機、循環、並行実装条件を独立監査した。両者とも責務分類と小分類数は維持し、実装開始前の横断Milestoneを追加する結論となった。

### 修正案1

#### タスク名

`L2 Design Freeze Gateを実装開始条件として追加する`

#### 分割した理由

正本と派生Indexは別の中分類だが、正本Schema、変更境界、Index Schema、Version、再構築は相互に拘束する。M1またはM2単独の設計完了だけで実装すると、後から確定した他方の設計でSchemaやRepositoryを変更する手戻りが起きるため、両方の設計成果物を横断確認してから本番実装へ進むGateを置く。

Gateは独立成果物を作る小分類ではなく、既存設計成果物の整合を判定するMilestoneとする。L6の品質評価とは異なり、実装前の設計整合性だけを対象にする。

Gateの前提はL2-M2-S1、L2-M1-S1〜S3、L2-M2-S2〜S4の完了とし、次を確認する。

- 全検索方式に必要なField、Relation、Temporal Metadata、Revisionを正本から取得できる
- Commit済み変更の対象ID、Revision、変更種別を通知または再照合できる
- 正本Snapshotの列挙または差分取得により全派生Indexを再構築できる
- 正本Schema、Index Schema、Embedding ModelのVersion互換性を判定できる
- Driver、Index Library、Migration、Transaction方式が共存できる
- 正本と派生IndexのSchema・Migration所有者が一意である

Gate後に設計変更が必要な場合は変更を禁止せず、影響する設計成果物を更新してGateを再通過する。Gate前の技術選定用PoCは許容するが、本番実装へ流用しない。

### 修正案2

#### タスク名

`着手依存と完了・Merge依存を分ける`

#### 分割した理由

Design Freeze後の実装をすべて直列化すると安全だが、独立PortとMockで進められる作業まで不要に待たせる。L2-M2-S5はGate通過後にM2-S3のPortとM1-S3のRepository境界を使って着手できる一方、正本照合Adapterを含めて完了・MergeするにはL2-M1-S10が必要である。開始可能性と統合可能性を一つの「依存」で表すと、手戻り防止か並行性のどちらかを失うため分ける。

L2-M1-S4をGateの直接依存先とし、後続M1実装は推移的にGate後とする。L2-M2-S5もGateを着手条件、L2-M1-S10を完了・Merge条件とする。

### 修正案3

#### タスク名

`git worktreeの所有境界と統合順を定義する`

#### 分割した理由

タスクの責務が独立していても、同じSchema、Migration、Interface、Registry、DI、Lockfileを複数worktreeが変更すると、Merge時に設計差異が再混入する。並行化単位を「タスク名」だけで決めず、変更するFile・Directoryの所有とMerge順まで固定して初めて、独立した小分類として安全に実行できる。

推奨実行規則は次のとおりとする。

1. Gate通過内容だけを記録したCommitを共通起点にする。既存の未Commit変更を混入させない。
2. L2-M1-S4〜S7はSchema、Migration、接続・変換Coreを共有するため直列にMergeする。
3. S6・S7のRepository Interfaceを固定後、L2-M1-S8〜S10を分岐する。共有Error、Mapper、Fixture、DI、Schemaを変更する場合は並行化せず、M2を早く解放する順にS10、S8、S9をMergeする。
4. L2-M2-S5のProvider契約をMerge後、L2-M2-S6〜S9を同じCommitから分岐し、各方式のDirectoryと固有テストだけを変更する。
5. Provider Registry、DI、共通Fixture、依存定義・Lockfile、生成物は単一Integratorだけが更新する。
6. L2-M1-S8・S9とL2-M2-S6〜S9のMerge後、L2-M2-S10でLifecycleを統合する。
7. 公開CLIのDI・Command配線はL2-M3だけが所有する。

### 修正後の依存関係案

```text
L2-M2-S1
    ↓
L2-M1-S1〜S3
    ↓
L2-M2-S2〜S4
    ↓
L2 Design Freeze Gate
    ↓
L2-M1-S4 → S5 → S6 → S7
                         ├─→ S8 ──────────────┐
                         ├─→ S9 ──────────────┤
                         └─→ S10 → L2-M2-S5 ─┤
                                                ├─→ L2-M2-S6〜S9
                                                ↓
                                           L2-M2-S10
```

L2-M2-S5はGate後にMockで着手可能とし、図のS10依存は完了・Merge条件を表す。L2-M2-S10はL2-M1-S8・S9、L2-M2-S4、L2-M2-S6〜S9を必須条件とする。正本Snapshot列挙／差分取得の実装所有はGateでL2-M1-S10またはL2-M2-S10のどちらか一方に確定する。

### 監査結果

- 新しい小分類は不要であり、分類の変更ではなくMilestone、依存辺、実行規則の修正で足りる。
- L2-M1-S7をS6後、S8〜S10前に置く順序は妥当である。
- L2-M1-S8〜S10の並行化は共有Fileを変更しない場合だけ妥当である。
- L2-M2-S6〜S9はProvider契約とLifecycle契約の確定後、方式別Directory所有により並行化できる。
- L2-M2親のL2-M1依存は「M2全体の完了依存」であり、M2-S1〜S4の設計開始条件ではない。これを開始依存と解釈すると循環が再発する。

### 承認対象

1. 新規小分類を追加せず、L2 Design Freeze GateをMilestoneとして追加する
2. L2-M1-S4の着手依存へGateを追加する
3. L2-M2-S5はGate後にMockで着手し、L2-M1-S10後に完了・Mergeする
4. 上記のworktree所有境界とMerge順を実行規則にする
5. L2-M2親のL2-M1依存を完了依存と明記する

### 承認結果（2026-08-11）

ユーザー承認により、上記5件を確定した。新しい小分類は追加せず、L2 Design Freeze Gateを横断MilestoneとしてTask Mapへ反映した。L2-M1-S4をGate通過後の本番実装開始点とし、L2-M2-S5はGate後にMockで着手可能、L2-M1-S10後に完了・Merge可能とした。

worktreeはGate通過Commitを共通起点とし、L2-M1-S4〜S7を直列、L2-M1-S8〜S10とL2-M2-S6〜S9を所有境界が分離できる場合だけ並行化する。共有Schema、Migration、Interface、Registry、DI、Lockfile、生成物、共有Fixtureは単一Ownerまたは直列Mergeとする。

## 決定 012: タスク分割スキルへ依存関係監査を追加する

### 状態

実施済み（2026-08-11）

### タスク名

`decompose-tasksスキルへ依存DAGとworktree監査を追加する`

### 分割した理由

責務と粒度だけが妥当でも、設計確定前の実装、親依存の開始・完了条件の混同、Mockで先行可能なTaskの過剰待機、共有Fileを同時変更するworktreeが残ると、実行段階で手戻りやMerge矛盾が発生する。今後の分割ごとに同じ観点を再現可能にするため、SKILL.mdと分割担当・妥当性確認担当の両プロンプトへ必須監査項目として追加した。

追加した観点は次のとおりである。

- 着手依存、完了・Merge依存、後続利用者の区別
- 循環、隠れ依存、不要な直列化、設計確定前の本番実装経路
- 親タスク間の開始依存と完了依存の区別
- 独立成果物を持たないDesign Gate／Milestoneの扱い
- worktreeの共通起点、File・Directory所有、共有資産の単一Owner、Merge順
- 承認済みタスクを含む全階層DAGの継続的な再監査

Skill Creatorの検証Scriptを実行し、Skill構造とFrontmatterが有効であることを確認した。

## 決定 013: 既存タスクの依存関係横断監査

### 状態

承認済み（2026-08-11）

### 監査方法

依存DAG、設計先行Gate、worktree所有を追加した`decompose-tasks`スキルを使い、分割担当がIssue #1とL1〜L6の全台帳を監査した。別の妥当性確認担当が、着手依存と完了・Merge依存、循環、過剰待機、推移的依存の重複を独立検証した。

L1内、L2-M1のS1〜S9、L2-M2のS6〜S9には、新しい循環や設計前実装経路を確認しなかった。L5の具体的なBudget・停止理由・Search Trace契約は、中分類分割時に再監査する。

### 既承認内容の文書整合補正

- L2-M2-S5がS3直後に着手できるように見えた依存図を、S3・S4完了後のDesign Freeze Gateを待つ図へ修正した。
- L2-M1親の開始入力へ、子S1が必要とするL2-M2-S1を明記した。
- L6台帳を、対象設計ごとの評価準備、対象実装後の評価、L5後の最終Acceptance・E2Eへ分けた。
- Gate通過記録に、Path／Glob、Owner Task、変更禁止Task、起点Commit、Merge順、生成物・Lockfile Ownerの実値表を必須化した。

いずれも既承認の工程・worktree規則をTask Mapへ正しく反映する補正であり、タスク責務の移動は行わない。

### 修正候補1

#### タスク名

`L4の着手依存と完了・Merge依存を分離する`

#### 分割した理由

現行台帳ではL4がL2へ依存する一方、L2-M2-S1はL4のKnowledge Search探索要求設計を入力とするため、L2依存を着手条件と読むと循環する。L4のArticle Analysis、Knowledge Search、Reading Value設計はL1から開始し、実CLIとの統合完了だけをL2-M3後にすることで、「検索から保存を逆算する」原則と実装順を両立する。

修正案は`着手: L1`、`完了・Merge: L2-M3`とする。L4分割後、L2-M2-S1の暫定前方参照をKnowledge Search探索要求設計の具体Task IDへ置換する。

### 修正候補2

#### タスク名

`正本Snapshot・差分取得PortをL2-M1-S10へ配置し、L2-M2-S10の完了依存を補う`

#### 分割した理由

派生Indexの完全再構築には、正本を列挙しCommit済み変更差分を取得するPortが必要だが、現在はM1-S10またはM2-S10のどちらが所有するか未確定である。M2が物理DBを直接読む経路を残さないため、正本読取りPortは取得RepositoryであるM1-S10、変更記録の原子的書込みはM1-S8・S9、消費・Checkpoint・再構築はM2-S10へ配置する。

M2-S10は`着手: M2-S4 + Design Freeze Gate（Mock可）`、`完了・Merge: M2-S6〜S9 + M1-S8・S9・S10`とする。M1-S7はM1-S8〜S10から推移的に満たされ、M2をState導出実装へ直接結合するため直接依存には追加しない。

### 修正候補3

#### タスク名

`L2-M3の内部設計・Adapter実装・完了依存を分離する`

#### 分割した理由

現行の`L1-M2、L2-M1、L2-M2`依存は、全依存を着手条件にするとCLI内部設計まで過剰に遅らせ、完了条件と解釈すると未確定Portへ本番Adapterを実装できてしまう。Framework、Validation、JSON境界等の内部設計はL1-M2後に開始し、正本・検索Adapterの本番実装だけをDesign Freeze Gate後、親全体の完了・MergeをL2-M1・L2-M2後にする。

具体的な子Taskへの依存配置はL2-M3小分類時に再監査する。

### 修正候補4

#### タスク名

`L3のSkill設計着手とCLI統合完了の依存を分離する`

#### 分割した理由

Knowledge Acquisition・Updateの判断と成果物契約はL1およびIssue 22〜26、39章から設計でき、L2全体の完成を待つ必要がない。一方、実行可能なKnowledge CLIとの統合完了にはL2-M3が必要である。よって`着手: L1`、`完了・Merge: L2-M3`へ分ける。

### 修正候補5

#### タスク名

`L5のWorkflow設計着手と全Component統合完了の依存を分離する`

#### 分割した理由

Parent OrchestrationのWorkflow、再調査、Budget、終了判定はIssue 8.1、16、37〜39章とL1契約を入力に先行設計できる。L2-M3・L3・L4の全実装を着手条件にすると、各Skillへ渡すBudget・停止理由・Search Trace契約の設計が後回しになり、後続Skill設計の手戻り原因になる。設計着手はL1、完了・MergeはL2-M3・L3・L4とする。

L3・L4・L5分割時に、Budget設定、停止理由返却、Search Trace生成・相関の具体Task間DAGを再監査する。

### 後続監査Trigger

- L4分割時にL2-M2-S1の暫定参照先を具体Task IDへ置換する。
- L2-M3分割時に先行可能なFramework設計とGate後のAdapter実装を別の依存へ落とす。
- 同期技術選定後にSnapshot、delta、cursor、outboxの所有を再確認する。
- L3・L4・L5分割時にBudget、停止理由、Search Trace契約の循環がないか確認する。

### 再承認単位

1. L4親依存とL2-M2-S1の前方参照方針
2. L2-M1-S10とL2-M2-S10のSnapshot・差分取得境界
3. L2-M3の着手・Adapter実装・完了依存
4. L3の着手・完了依存
5. L5の着手・完了依存

### 承認結果（2026-08-11）

ユーザー承認により、修正候補5件をすべて確定した。

- L4はL1後に設計着手し、L2-M3後にCLI統合を完了・Mergeする。
- 正本Snapshot列挙・Commit済み差分取得PortはL2-M1-S10が所有し、L2-M2-S10はL2-M1-S8〜S10とL2-M2-S6〜S9の完了後に統合を完了・Mergeする。L2-M1-S7への直接依存は追加しない。
- L2-M3はL1-M2後に内部設計へ着手し、正本・検索Adapterの本番実装をL2 Design Freeze Gate後、全体の完了・MergeをL2-M1・L2-M2後とする。
- L3はL1後にSkill設計へ着手し、L2-M3後にCLI統合を完了・Mergeする。
- L5はL1後にWorkflow設計へ着手し、L2-M3・L3・L4後に全Component統合を完了・Mergeする。

Task Mapの対象を「再検討中」から承認済みへ戻し、次の議論をL2-M3の小分類へ移す。

## 決定 014: Knowledge CLI実行基盤とJSON境界の小分類

### 状態

承認済み（2026-08-11）

### 分割・監査方法

分割担当がIssue #1、承認記録、Task Mapを確認し、L2-M3を公開契約の再設計ではなく、未確定のCLI技術、内部設計、共通Runtime、機能別Adapter、配線・配布へ分割した。別の妥当性確認担当がIssue原文へ戻り、原典で決定済みの内容の再議論、L1-M2・L2-M1・L2-M2との重複、L6とのテスト所有、依存DAG、Design Freeze Gate、worktree所有を独立監査した。

当初の9件案は、Process CoreとJSON／Error境界、配線・配布と契約適合Suiteをそれぞれ分離していた。独立監査では、いずれも同じ共有境界とFileを変更し、分けても並行化できないため過剰分割と判定した。これらを統合した8件を最終提案とする。

### L2-M3-S1

#### タスク名

`CLI Framework・設定／配布方式を選定しADRを確定する`

#### 分割した理由

CLI Framework、設定読込み、Build・配布方式はIssue 48章で未決定の物理詳細であり、全CLI実装を拘束する共通判断である。内部設計や実装へ埋め込まず、採否理由を独立して変更・承認可能にする。

### L2-M3-S2

#### タスク名

`CLI内部実行境界とCommand-to-Port配線を詳細設計する`

#### 分割した理由

L1-M2は公開契約を所有するが、Handler、Validation、Dispatch、Adapter Port、Error変換、初期化、DIの内部構造は定義しない。全公開CommandをM1・M2のPortへ割り当て、実装前にCLI内部の責務境界を確定するため独立させる。

### L2-M3-S3

#### タスク名

`共通CLI Runtime・JSON検証・Error境界を実装する`

#### 分割した理由

Command受付、Dispatch、JSON Schema検証、Serialization、Error・終了コード変換は全操作が共有する一つのProcess境界である。個別Adapterへ分散するとL1公開契約が重複実装されるため、共通Runtimeとして一つの完了単位にする。

### L2-M3-S4

#### タスク名

`Knowledge更新Command Adapterを実装する`

#### 分割した理由

create、attach-evidence、revise、supersedeはTransaction、競合、非破壊履歴を伴う同じ更新境界を持つ。取得・検索と失敗条件が異なるため、CLIが意味判断せずM1の更新Portへ接続する単位として分離する。

### L2-M3-S5

#### タスク名

`Knowledge詳細・Evidence取得Command Adapterを実装する`

#### 分割した理由

get、get-evidenceは指定IDから正本Aggregate、Evidence、履歴を取得する。候補検索とはBackendと結果契約が異なるため、L2-M1-S10への読取り専用Adapterとして分離する。

### L2-M3-S6

#### タスク名

`候補検索Command Adapterを実装する`

#### 分割した理由

Issue 12章の6検索プリミティブは、集合規則、Pagination、安定順序を共有し、M2の候補検索Portへ接続する同じ変更理由を持つ。方式別Providerを再実装せず、意味判断を加えないCLI境界として一つにまとめる。

### L2-M3-S7

#### タスク名

`Index Lifecycle管理Command Adapterを実装する`

#### 分割した理由

Index Version、再構築、障害回復は、Knowledgeの読取り・更新とは異なる運用境界である。Lifecycle状態をCLIで二重管理せず、L2-M2-S10へ接続する入口だけを所有するため分離する。

### L2-M3-S8

#### タスク名

`公開CLIを配線・配布しL1契約適合を確認する`

#### 分割した理由

Entry Point、Command Registry、DI、設定、Build、配布Artifactと全公開Commandの契約適合は、同じ共有Fileと実行Artifactを完成させる統合責務である。別タスクにするとworktree競合と未検証の中間状態を増やすため一つにまとめる。

### 依存関係

```text
L1-M2
  ↓
S1 技術選定 → S2 内部設計 → S3 共通Runtime
                              + L2 Design Freeze Gate
                              ├─→ S4 更新Adapter
                              ├─→ S5 取得Adapter
                              ├─→ S6 検索Adapter
                              └─→ S7 Lifecycle Adapter
                                           ↓
                                  S8 配線・配布・契約適合
```

- S4の完了・MergeにはL2-M1-S8・S9、S5にはL2-M1-S10、S6にはL2-M2-S5〜S9、S7にはL2-M2-S10を必要とする。
- S8はS3〜S7とL1-M2-S6を入力とし、L2-M1・L2-M2完了後に完了・Mergeする。
- L2 Design Freeze GateはS4〜S7の本番Adapter開始条件として維持し、S1〜S3はL1-M2後に先行可能とする。GateへS1・S2は追加しない。
- S1〜S3は共通Runtime、Error、設定境界を共有するため直列にMergeする。
- S4〜S7はAdapter別Directoryと固有テストだけを所有する場合に並行化する。
- Entry Point、Command Registry、DI、設定既定値、Build・配布定義、依存定義・Lockfile、生成物、共有FixtureはS8の単一Ownerとする。
- S4〜S7のworktreeはS3をMerge済みかつGate通過済みのCommitを共通起点にし、Merge順は`S1 → S2 → S3 → S4〜S7 → S8`とする。

### 妥当性監査結果

- 8件はいずれもIssueの未決定詳細または未実装成果物であり、公開コマンドやJSON Schemaを再決定するタスクではない。
- S3とS8の統合により、共有Fileを触るだけの直列タスクを増やす過剰分割を避けた。
- S8の契約適合はCLI固有テストであり、L6が所有する共有Dataset、Agent評価、Acceptance・E2Eとは重複しない。
- 初回8件監査時点では既存構造の修正候補なしとしたが、後続の追加監査でDesign Freeze Gateの不足を検出した。以下の最終修正を優先する。

### 追加監査による最終修正

8件案に対して別の妥当性確認担当が、成果物の独立性と「L2全体の設計を先に完了させる」という承認済み方針を再監査した。その結果、共有Fileを扱うことは直列Mergeの理由にはなるが、成果物と完了判定が異なるTaskを統合する理由にはならないと判断した。

次の2組を再分離し、最終提案を9件へ修正する。

1. CLI Process Lifecycle・Dispatch実装と、JSON Validation・Serialization・Error／終了コード境界実装
2. 設定・初期化・Composition Root・配布Artifact実装と、実行Artifact外部からのL1公開契約Black-box検証

また、L2-M3-S1・S2でCLI技術とM1・M2のPort配線を確定した後に、L2-M1・M2・M3の本番実装へ進むよう、既存L2 Design Freeze Gateの前提追加を修正候補とする。これは後続分割で判明した承認済み構造の不足であり、再承認前は確定Gateへ反映しない。

### 修正候補

#### タスク名

`L2 Design Freeze GateへL2-M3の技術選定・内部設計を追加する`

#### 分割した理由

現在のGateは正本と派生Indexの設計だけを確認するため、M1・M2実装後にCLI Runtimeとの技術非互換や、公開Commandを接続できないPort不足が判明する手戻り経路が残る。L2-M3-S1の技術ADRとS2のCommand-to-Port MappingをGate前に確定し、L2全体の設計を整合させてから本番実装へ進む。

S1はL1-M2後に調査を開始し、L2-M1-S1・L2-M2-S2後に技術共存を確認してADRを確定する。S2はS1、L2-M1-S3、L2-M2-S3・S4を入力にし、既存Gate前提とS1・S2が揃った時点でGateを判定する。この依存にはM1・M2からS1・S2へ戻る辺がないため循環しない。

### 最終小分類案

#### L2-M3-S1

タスク名：`CLI Framework・JSON検証・Build／配布技術を選定しADRを確定する`

分割した理由：CLI共通技術の採否を一つの意思決定成果物として確定し、M1・M2の採用技術と共存できることを実装前に確認するため。

#### L2-M3-S2

タスク名：`CLI内部ArchitectureとCommand-to-Port Mappingを詳細設計する`

分割した理由：L1の全公開CommandをM1・M2のPortへ漏れなく割り当て、Handler、Validation、初期化、DIの内部責務を本番実装前に確定するため。

#### L2-M3-S3

タスク名：`CLI Process CoreとCommand Dispatchを実装する`

分割した理由：Process Lifecycle、引数解析、共通Command Interface、Dispatch、help・version入口は、JSON変換やBackend処理から独立してFake Handlerで完了判定できるため。

#### L2-M3-S4

タスク名：`JSON Validation・Serialization・Error／終了コード境界を実装する`

分割した理由：JSON Schema検証、stdout・stderr、内部Error変換、終了コードはL1契約へ適合する一つの機械境界であり、Process障害と別に完了判定できるため。

#### L2-M3-S5

タスク名：`取得・候補検索Command Handler／Adapterを実装する`

分割した理由：全読取り操作は副作用を持たず、集合規則、Pagination、安定順序を保ってM1・M2の取得Portへ接続する共通の変更理由を持つため。

#### L2-M3-S6

タスク名：`Knowledge更新Command Handler／Adapterを実装する`

分割した理由：4更新操作はTransaction、競合、非破壊履歴を伴い、読取りとは異なる失敗条件を持つため。

#### L2-M3-S7

タスク名：`Index管理・再構築／障害回復CommandをCLIへ接続する`

分割した理由：Index Lifecycle操作はKnowledge読取り・更新と異なる運用境界であり、M2-S10の状態をCLIで二重管理せず接続するため。

#### L2-M3-S8

タスク名：`設定・初期化・Composition Rootと配布Artifactを実装する`

分割した理由：Registry、DI、設定、Build、Version埋込みを単一Ownerへ集約し、全AdapterをCodexから利用可能な実行Artifactへ統合するため。

#### L2-M3-S9

タスク名：`実行ArtifactのL1公開契約適合をBlack-box検証する`

分割した理由：Artifact生成責任と外部契約検証責任を分け、全公開CommandのJSON、Error、終了コード、Versionを実行物の外側から独立して判定するため。共有Dataset・Skill連携・Acceptance・E2EはL6に残す。

### 最終依存関係

```text
L1-M2 ─────────────→ S1調査開始
L2-M1-S1 + L2-M2-S2 → S1 ADR確定
S1 + L2-M1-S3 + L2-M2-S3・S4
                    ↓
                 S2設計
                    ↓
既存Gate前提 + S1 + S2
                    ↓
        L2 Design Freeze Gate（修正候補）
                    ↓
            S3 → S4 → S5・S6・S7
                              ↓
                             S8 → S9
```

- S5の完了・MergeはL2-M1-S10・L2-M2-S5〜S9、S6はL2-M1-S8・S9、S7はL2-M2-S10を必要とする。
- S8はS3〜S7を統合し、L2-M1・L2-M2後に完了・Mergeする。
- S3、S4は共有Coreを扱うため直列Mergeする。
- S5〜S7はS4 Merge後、方式別Directoryだけを所有して並行化する。
- Registry、DI、設定、Build定義、Lockfile、生成物はS8、Black-box契約テストDirectoryはS9が単独所有する。
- Gate通過Commitを本番実装worktreeの共通起点とし、Merge順は`S3 → S4 → S5〜S7 → S8 → S9`とする。

### 再承認単位

1. L2 Design Freeze GateへL2-M3-S1・S2を追加する修正
2. L2-M3小分類9件と依存DAG・worktree所有境界

### 承認結果（2026-08-11）

ユーザーは、L2 Design Freeze GateへL2-M3-S1・S2を追加する修正と、L2-M3小分類9件、依存DAG、worktree所有境界を承認した。

これにより、L2 Design Freeze Gateは正本・検索・Index設計に加え、CLI技術の共存性と全公開CommandのPort配線まで確認してから本番実装へ進む確定Milestoneとなる。Task MapのL2-M3-S1〜S9を承認済みに更新し、次の議論をL3「Knowledge蓄積Skills」の中分類へ移す。

## 決定 015: Knowledge蓄積Skillsの中分類

### 状態

承認済み（2026-08-11）

### 分割・監査方法

分割担当がIssue #1、承認記録、Task Mapを確認し、L3をKnowledge候補を生成するSkillと、既存Knowledgeとの比較・更新を行うSkillへ分割した。別の妥当性確認担当が、Issue既定事項の再議論、L5のOrchestration責務、L6の横断評価責務、L1・L2との依存、Design Freeze Gate、worktree所有を独立監査した。

既存承認済み構造に先行修正が必要な矛盾は確認されなかった。中分類は次の2件とし、引き渡し契約は独立実行物ではないため第三の中分類にはしない。

### L3-M1

#### タスク名

`Knowledge Acquisition Skillを詳細設計・実装する`

#### 分割した理由

Conversation／Task EpisodeからCandidate Assertion、Concept、Evidence、Source Reference、Scope、Evidence Strengthを抽出し、保存候補がないことも表現するproducer責務は、既存Knowledgeとの照合・更新判断を行うconsumer責務と入力、変更理由、完了判定が異なるため分離する。Issueで既定の抽出原則を選び直さず、未定義のEpisode境界、候補採否・正規化手順、L1 Fieldへの写像、Markdown契約とSkill実装を成果物にする。

主成果物はKnowledge Acquisition Skill一式、L3-M2へ渡すCandidate Markdown契約、no-candidate表現である。

### L3-M2

#### タスク名

`Knowledge Update Skillを詳細設計・実装する`

#### 分割した理由

Candidateを既存Knowledgeと意味的に比較し、create、attach-evidence、revise、supersede、永続化しない判断を選び、Knowledge CLIへ接続するconsumer責務は、Acquisitionと利用契約、障害条件、CLI依存が異なるため分離する。Issueで既定の4更新操作と訂正原則を選び直さず、未定義の検索・同一性・操作選択手順、skip理由、CLI Command Mapping、再試行・Error処理、Markdown契約とSkill実装を成果物にする。skipは第5のCLI更新操作ではなく、永続化しない判断として扱う。

主成果物はKnowledge Update Skill一式、Candidate受入規則、既存Knowledge照合手順、操作決定・理由Markdown、CLI Command Mappingである。

### 妥当性監査結果

- 2件でL3親成果物を充足し、L5のSkill起動・no-candidate分岐・全体制御、L6の共有Fixture・Dataset・評価Harnessとは重複しない。
- M1→M2のMarkdown契約はM1がproducerとして単独所有し、M2がread-only consumerになるため、独立した第三中分類は不要である。
- Knowledge-worthy Evidence、Evidence Strength、AI説明単独・質問単独をEvidence化しない原則、4更新操作、過去Evidence非削除はIssueの固定入力とし、再議論しない。

### 依存関係

```text
L1-M1
  ├─→ M1 Acquisition設計
  └─→ M2 Update判断設計開始

M1 Candidate Markdown契約 + L1-M2
                    ↓
             M2受入・CLI Mapping設計
                    ↓
       L3 Skill Design Freeze Gate
              ┌─────┴─────┐
              ↓           ↓
          M1実装       M2 Mock実装
                          + L2-M3-S5・S6
                          + L2-M3-S8
                              ↓
                         M2実CLI統合
                              ↓
                    L2-M3-S9後にL3 Release
```

- M1はL1-M1を入力に設計・実装でき、L2を完了条件にしない。
- M2の判断設計はL1-M1後に開始できる。M1のCandidate Markdown契約確定とCLI MappingにはM1設計成果物とL1-M2を必要とする。
- M2の実CLI統合はL2-M3-S5・S6を直接入力とし、L2-M3-S8の実行Artifactへ接続する。L3全体のReleaseはL2-M3-S9後とする。

### L3 Skill Design Freeze Gate・worktree計画

Gateは独立成果物を持つ中分類ではなく、両Skillの詳細設計、M1→M2 Markdown契約、no-candidate表現、L1 Field・CLI Mapping、Codex／CLI／L5責務境界、Path Ownerを実装前に確定するMilestoneとする。

- Gate通過Commitを実装worktreeの共通起点にする。
- M1とM2は別Skill Directoryを単独所有し、Gate後に並行実装できる。
- Candidate契約Artifactと契約例はM1だけが変更し、M2はread-onlyで利用する。
- 共有Registry・OrchestrationはL5、共有Fixture・Dataset・評価HarnessはL6が単独所有する。
- 実Path／Globと小分類Task IDは、L3-M1・M2の小分類時に確定する。

### 承認対象

1. L3-M1 `Knowledge Acquisition Skillを詳細設計・実装する`
2. L3-M2 `Knowledge Update Skillを詳細設計・実装する`
3. L3 Skill Design Freeze Gate、依存DAG、worktree所有境界

### 承認結果（2026-08-11）

ユーザーはL3中分類2件、L3 Skill Design Freeze Gate、依存DAG、worktree所有境界を承認した。

Task MapのL3-M1・M2を承認済みに更新し、承認済みの進行規則に従って、次の議論を最初の中分類L3-M1「Knowledge Acquisition Skill」の小分類へ移す。

## 決定 016: Knowledge Acquisition Skillの小分類

### 状態

承認済み（2026-08-11）

### 分割・監査方法

分割担当がIssue #1、承認記録、Task Mapを確認し、L3-M1を入力・Evidence採否設計、Candidate正規化設計、M1→M2引き渡し契約、Skill実装、Skill固有検証の5件へ分割した。別の妥当性確認担当が、各候補をこれ以上子へ分割しない最終実行単位として、Issue既定事項の再議論、L3-M2・L5・L6との重複、依存DAG、Design Freeze Gate、worktree所有を独立監査した。

監査の結果、独立したSkill固有検証Taskは、構造・必須Sectionの確認がSkill実装の完了条件と重複し、意味的Scenario評価がL6の責務と重複するため削除した。最終案は4件とする。

### 既存構造の修正候補

#### 対象

`L3 Skill Design Freeze Gate`

#### 問題

承認済みGateはno-candidate表現を確認するが、「十分なEpisodeを評価した結果、保存候補がない状態」と、「入力欠落・切断・Source Reference解決不能により評価できない状態」の区別を明示していない。両者を同一視すると、L5が入力不備を正常終了として扱う可能性がある。

#### 修正案

Gate条件へ`no-candidateと入力不足を区別できる入力・出力契約`を追加する。M1は状態と診断情報を出力し、再試行・終了・M2非起動の制御はL5が所有する。

#### 影響範囲・再承認単位

L3-M1-S1・S3、L3-M2の受入設計、L5の分岐設計へ影響する。L3の親子分類は変更せず、今回のL3-M1小分類4件と同時に再承認する。

### L3-M1-S1

#### タスク名

`Episode入力・Source Reference・Evidence候補採否を詳細設計する`

#### 分割した理由

L5が起動時点を判断する一方、M1は渡されたEpisodeの必須内容、出典を追跡する位置情報、入力不足の判定、Knowledge-worthyなEvidence候補の採否手順を定義する必要があるため分離する。Evidence Strengthの3分類、質問、AI説明、訂正の原則はIssueの固定入力とし、選び直さない。

主成果物はEpisode入力・Evidence採否境界設計書とする。

### L3-M1-S2

#### タスク名

`Knowledge Candidate正規化・L1 Field写像を詳細設計する`

#### 分割した理由

採用したEvidenceから、独立評価可能なAssertion、検索AnchorであるConcept、Scope、Source Reference、Evidence Strength、Temporal情報を正規化し、L1論理Fieldへ写像する責務は、観測箇所とEvidence採否を決めるS1とは異なるため分離する。既存Knowledgeとの同一性・更新操作はL3-M2へ残す。

主成果物はCandidate正規化・L1 Field写像仕様とする。

### L3-M1-S3

#### タスク名

`Candidate Markdown／no-candidate・入力不足引渡し契約を定義する`

#### 分割した理由

M1が所有しM2とL5が読み取り利用する成果物について、必須Section、複数Candidate境界、Source Reference、順序・省略規則と規範例を一つの契約として確定する必要があるため分離する。候補あり、十分な評価後の候補なし、入力不足による評価不能を区別し、no-candidate後の終了や再試行の制御はL5へ残す。

主成果物はCandidate Markdown契約とcanonical例とする。canonical例は規範文書であり、L6の実行Fixture・Datasetにはしない。

### L3-M1-S4

#### タスク名

`Knowledge Acquisition Skill instructions／referencesを実装し、固有構造契約を確認する`

#### 分割した理由

確定設計を、入力、抽出手順、禁止事項、出力、失敗時挙動を持つ実行可能なSkill packageへ変換する未実装成果物であるため分離する。SKILL.mdと補助Referenceは相互依存が強いため一つにまとめ、静的構造・必須Section・status契約の確認を完了条件へ含める。Knowledge-worthy抽出、AI説明、質問、訂正などの意味的Scenario評価はL6へ残す。

主成果物は実行可能で固有構造契約に適合するKnowledge Acquisition Skill packageとする。

### 妥当性監査結果

- S1は観測箇所とEvidence採否、S2は採用Evidenceから正規化Candidateを作る責務であり、分離しても各成果物を単独判定できる。
- S3のcanonical例は引き渡し契約の規範例としてGate前に確定できる。実行Fixture・Datasetへの転用はしない。
- 独立した検証Taskは置かない。Skill固有の構造契約確認はS4、共有Fixture・Dataset・Harness・意味的Agent評価・AcceptanceはL6が所有する。
- 技術選定Taskは追加しない。Codex SkillとMarkdown成果物はIssueで固定済みであり、新しい技術判断がないためである。

### 依存関係

```text
L1-M1 ─→ S1開始
L1-M1 ─→ S2先行設計
S1完了 ─→ S2確定
S1 + S2 ─→ S3
S3 ─→ L3-M2の受入設計確定
S3 + L3-M2設計完了
       ↓
L3 Skill Design Freeze Gate
       ↓
      S4

S3 ─→ L6のM1評価設計
S4 ─→ L6のM1 Agent評価実行
```

- S1とS2を全面直列にはせず、S2はL1-M1から先行設計し、確定だけS1を待つ。
- S4の本番実装はL3 Skill Design Freeze Gate後に開始する。
- L3-M2も判断設計は先行でき、Candidate受入契約の確定だけS3を待つ。

### worktree所有境界

- S1: `docs/design/knowledge-acquisition/episode-evidence-boundary.md`
- S2: `docs/design/knowledge-acquisition/candidate-normalization-mapping.md`
- S3: `.agents/skills/knowledge-acquisition/references/candidate-markdown-contract.md`
- S4: `.agents/skills/knowledge-acquisition/SKILL.md`と、S3契約を除く同Skillの`references/**`
- L5: 共有Registry、起動・no-candidate・再試行・終了のOrchestration
- L6: 共有Fixture、Dataset、Harness、Agent評価、Acceptance・E2E

S1・S2は別Fileで先行可能だが、S2確定はS1完了後とする。Merge順は`S1 → S2 → S3 → Gate → S4`とし、Gate通過CommitをM1・M2本番実装worktreeの共通起点にする。M2の実Path／GlobはL3-M2小分類時に確定する。

### 承認対象

1. L3 Skill Design Freeze Gateへno-candidateと入力不足の区別を追加する修正
2. L3-M1小分類4件
3. 依存DAGとworktree所有境界

### 承認結果（2026-08-11）

ユーザーはL3 Skill Design Freeze Gateの修正、L3-M1小分類4件、依存DAG、worktree所有境界を承認した。

Task Mapへ確定構造を反映し、次の議論をL3-M2「Knowledge Update Skill」の小分類へ移す。

## 決定 017: Knowledge Update Skillの小分類

### 状態

承認済み（2026-08-11）

### 分割・監査方法

分割担当がIssue #1、承認記録、Task Mapを確認し、L3-M2を入力受入、既存Knowledge探索・比較、更新判断契約、CLI実行契約、Mock前提Skill実装、実CLI統合の6件へ分割した。別の妥当性確認担当が、Issue既定事項の再議論、L3-M1・L4・L5・L6との重複、L1／L2との依存DAG、Design Freeze Gate、worktree所有を独立監査した。

監査の結果、6件を維持しつつ、S1の責務を入力の評価可否へ限定し、S3・S4・S5・S6の名称を成果物に合わせて明確化する。加えて、承認済み構造へ2件の修正候補を提示する。

### 既存構造の修正候補1

#### 対象

`L3 Skill Design Freeze Gate`

#### 問題

承認済みGateは正常なno-candidateと入力不足・評価不能を区別するが、M2で発生する「正常な非永続化判断」「Candidate入力不備・評価不能」「CLI実行失敗・部分成功」の区別までは明示していない。これらを同一状態としてL5へ渡すと、正常終了、再試行、停止を誤る可能性がある。

#### 修正案

Gate条件へ、上記3状態を区別するM2出力・診断契約を追加する。M2は状態と診断を報告し、再試行回数、再起動、停止、Workflow継続はL5が所有する。

#### 影響範囲・再承認単位

L3-M2-S1・S3・S4、L5の分岐設計へ影響する。L3の親子分類は変更せず、今回のL3-M2小分類6件と同時に再承認する。

### 既存構造の修正候補2

#### 対象

`L2-M2-S1 候補検索・Index Lifecycle要件と正本設計Review Gate基準を定義する`

#### 問題

L2-M2の検索基盤は検索要件から逆算して設計するが、現在の入力にはL3-M2が重複Assertion防止・Evidence追加・訂正対象特定に必要とする更新用検索要求が明示されていない。L4の読書価値向け検索だけで設計を確定すると、更新Skillが必要とする検索経路が後から判明し手戻りになる可能性がある。

#### 修正案

L3-M2-S2の技術非依存な既存Knowledge探索要求を、L2-M2-S1の入力とL2 Design Freeze Gateの確認対象へ追加する。L3-M2-S2はL1とL3-M1-S3を入力にL2実装前から設計できるため、循環は生じない。

#### 影響範囲・再承認単位

L2-M2-S1、L2 Design Freeze Gate、L3-M2-S2へ影響する。L2-M2の小分類数・成果物Ownerは変更せず、依存入力だけを今回のL3-M2小分類6件と同時に再承認する。

### L3-M2-S1

#### タスク名

`Candidate入力受入・評価可否境界を詳細設計する`

#### 分割した理由

L3-M1-S3のCandidate Markdown契約を変更せず、正常Candidateだけを更新評価へ受け入れ、入力不備・評価不能を保存処理や`Skipped candidates`へ混入させない境界がIssueでは未定義なため分離する。M2は状態と診断を返すが、起動・再試行・終了制御はL5へ残す。

主成果物は`docs/design/knowledge-update/candidate-acceptance.md`とする。

### L3-M2-S2

#### タスク名

`Existing Knowledge探索・意味比較手順を詳細設計する`

#### 分割した理由

Issueは新規作成前の既存Knowledge検索とCodexによる同一性判断を定めているが、Candidateからの検索順、同一・関連・Evidence重複・矛盾・時点差の比較観点、比較の十分性は未定義なため分離する。L4のユーザー知識状態評価ではなく、更新対象を特定する比較に限定する。

主成果物は`docs/design/knowledge-update/existing-knowledge-comparison.md`とする。

### L3-M2-S3

#### タスク名

`4更新操作選択・非永続化判断とDecision Markdown契約を定義する`

#### 分割した理由

4更新操作、訂正、履歴保持はIssueとL1で固定済みだが、Candidateと比較結果からの操作選択条件、更新しない判断、候補間競合、理由・結果の記録形式は未定義なため分離する。非永続化は第5のCLI操作ではなく、判断と理由として記録する。

主成果物は`.agents/skills/knowledge-update/references/update-decision-contract.md`とする。

### L3-M2-S4

#### タスク名

`CLI Command Mapping・実行結果／失敗引渡しを詳細設計する`

#### 分割した理由

意味的な判断をL1の検索・取得・4更新Commandへ変換する順序、Request・Response対応、成功・部分失敗・再試行可能性の報告方法はIssueでは未定義なため分離する。CLI契約・冪等性を再設計せず入力として利用し、再試行回数・停止制御はL5へ残す。

主成果物は`.agents/skills/knowledge-update/references/cli-command-mapping.md`とする。

### L3-M2-S5

#### タスク名

`Knowledge Update Skill instructions／referencesをMock CLI契約に基づき実装し、固有構造契約を確認する`

#### 分割した理由

確定設計を実行可能なSkill packageへ変換する成果物が未実装であり、実CLI完成前でもL1 Schema準拠Mockにより必須Section、参照、Command名、出力Templateを確認できるため、実CLI統合から分離する。意味的Agent評価と共有評価資産はL6へ残す。

主成果物は`.agents/skills/knowledge-update/SKILL.md`、S3・S4以外の同Skill `references/**`、同Skillの固有Contract確認資産とする。

### L3-M2-S6

#### タスク名

`Knowledge Update Skillを実CLI Artifactへ接続し、固有Command連携を確認する`

#### 分割した理由

Mock前提のSkill完成と、L2の実行Artifactで検索・取得・更新Commandを呼び、判断結果と永続化結果を対応付けられる状態は依存時期と完了条件が異なるため分離する。M2固有のComponent連携だけを確認し、L2-M3-S9の全CLI契約検証やL6の共有Agent評価・E2Eは繰り返さない。

主成果物は`tests/component/knowledge-update/**`とする。

### 妥当性監査結果

- S1は入力の評価可否、S3は正常Candidateに対する永続化・非永続化判断であり、責務を分離できる。
- S2は更新対象特定に限定し、L4のTarget Claimに対するknowledge assessmentや汎用Search Traceを含めない。
- S4はL1／L2のCLI契約を入力としてM2側の利用順・結果引渡しだけを設計する。
- S5のMock確認はSkill固有の構造契約に限定し、意味的Scenario評価はL6が所有する。
- S6はM2固有のComponent連携であり、L2-M3-S9のCLI Black-box契約検証、L6のAcceptance・E2Eと重複させない。
- 独立した技術選定Taskは追加しない。Codex Skill、Markdown、Knowledge CLIはIssueとL1／L2で固定済みである。

### 依存関係

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

S2の検索要求 ─→ L2-M2-S1 ─→ L2 Design Freeze Gate
```

- S1とS2は別設計Fileで並行可能とし、両方の完了後にS3を確定する。
- S4はL1-M2から先行設計できるが、検索MappingはS2、更新Mappingと結果引渡しはS3の完了後に確定する。
- S5はL3 Gate後にMockで先行でき、L2完成を待たない。S6だけがL2-M3-S5・S6・S8を待つ。
- L2-M3-S9はS6の着手条件ではなくL3 Release条件とする。

### worktree所有境界

- S1: `docs/design/knowledge-update/candidate-acceptance.md`
- S2: `docs/design/knowledge-update/existing-knowledge-comparison.md`
- S3: `.agents/skills/knowledge-update/references/update-decision-contract.md`
- S4: `.agents/skills/knowledge-update/references/cli-command-mapping.md`
- S5: `.agents/skills/knowledge-update/SKILL.md`、S3・S4以外の`references/**`、同Skillの固有Contract確認資産
- S6: `tests/component/knowledge-update/**`
- L3-M1-S3のCandidate契約、L2のCLI実装・Schema・DI・Registryはread-onlyとする。
- 共有Registry・起動・再試行・終了制御はL5、共有Fixture・Dataset・Harness・Agent評価・Acceptance・E2EはL6が単独所有する。

Merge順は`L3-M1-S3 → S1・S2 → S3 → S4 → Gate → S5 → S6`とし、Gate通過CommitをL3-M1-S4とL3-M2-S5の共通起点にする。

### 承認対象

1. L3 Skill Design Freeze Gateの状態区別追加
2. L2-M2-S1とL2 Design Freeze Gateへの更新用検索要求入力追加
3. L3-M2小分類6件
4. 依存DAGとworktree所有境界

### 承認結果（2026-08-11）

ユーザーはL3 Skill Design Freeze Gateの状態区別追加、L2-M2-S1とL2 Design Freeze Gateへの更新用検索要求入力追加、L3-M2小分類6件、依存DAG、worktree所有境界を承認した。

Task Mapへ確定構造を反映し、次の議論をL4「記事価値判定Skills」の中分類へ移す。

## 決定 018: 記事価値判定Skillsの中分類

### 状態

承認済み（2026-08-11）

### 分割・監査方法

分割担当がIssue #1、承認記録、Task Mapを確認し、L4をArticle Analysis、Knowledge Search、Reading Valueの3件へ分割した。別の妥当性確認担当が、Issue既定事項の再議論、L3-M2・L5・L6との重複、SupportとReliabilityの責務、L1／L2との依存DAG、Design Freeze Gate、worktree所有を独立監査した。

監査の結果、3件を維持する。Article Analysisは記事内のSupport観測、Knowledge SearchはユーザーKnowledge Assessmentとraw Search Trace、Reading ValueはSupportを含む入力からReliabilityと読む価値を判断する責務へ限定する。加えて、承認済み構造へ2件の修正候補を提示する。

### 既存構造の修正候補1

#### 対象

`L4 記事価値判定Skillsを設計・実装する`の依存関係

#### 問題

現行の`着手: L1。完了・Merge: L2-M3`では、CLIを直接利用しないArticle AnalysisとReading Value、およびMock CLIで実装可能なKnowledge SearchまでL2-M3完了待ちに見え、過剰に直列化される。

#### 修正案

`設計着手: L1。本番Skill実装: L4 Skill Design Freeze Gate。L4-M2実CLI統合: L2-M3-S5・S8。完了・Release: L4-M1〜M3、L2-M3-S9`へ変更する。

#### 影響範囲・再承認単位

L4、L4-M1〜M3、L5、L6へ影響する。中分類3件、L4 Gate、依存DAG、worktree所有と同時に再承認する。

### 既存構造の修正候補2

#### 対象

`L2-M2-S1`、L2-M2親の開始入力、`L2 Design Freeze Gate`

#### 問題

現行の`L4のKnowledge Search探索要求設計（暫定前方参照）`では、成果物Ownerと完了条件が未確定である。L4-M2全体完了を待つと、検索基盤完成後にKnowledge Searchを統合する循環も生じる。

#### 修正案

依存先を`L4-M2配下の技術非依存な探索要求・Query Journey・停止理由・必要取得情報の設計成果物（小分類ID確定まで前方参照）`へ具体化する。L4-M2の先行設計成果物をL2-M2-S1へ入力し、L2 GateではS1への反映を確認する。小分類承認時に具体S-IDへ最終置換する。

#### 影響範囲・再承認単位

L4-M2、L2-M2、L2-M2-S1、L2 Design Freeze Gateへ影響する。L4中分類と同時に再承認する。

### L4-M1

#### タスク名

`Article Analysis Skillを詳細設計・実装する`

#### 分割した理由

URLや記事を、ユーザーKnowledgeと比較可能なArticle overviewとClaim群へ変換するproducer責務は、既存Knowledgeの探索や読書価値判断と入力・変更理由・完了判定が異なるため分離する。Claim Role、Importance、Location、SupportはIssueの固定入力とし、未定義の入力境界、Claim識別、Support根拠、後続へのMarkdown契約、再分析応答、Skill実装を成果物にする。

### L4-M2

#### タスク名

`Knowledge Search Skillを詳細設計・実装する`

#### 分割した理由

Target ClaimからQueryを生成・変更し、CLIを反復利用してEvidence付きKnowledge Assessmentとraw Search Traceを生成する責務は、記事のClaim化や読書価値判断と異なるため分離する。7状態、停止原則、`no_evidence is not unknown`、CLIは候補取得のみというIssueの原則は選び直さず、技術非依存な探索要求、Query Journey、十分性・停止理由、CLI Command Mapping、Markdown契約、Skill実装と実CLI統合を成果物にする。

### L4-M3

#### タスク名

`Reading Value Skillを詳細設計・実装する`

#### 分割した理由

Article AnalysisとKnowledge Assessmentを統合し、Recognition Gain、Reliability、Attention Costから読む範囲を判断する責務は、Claimを抽出する責務やユーザーがClaimを知っているかを評価する責務と異なるため分離する。3推奨と認識利得分類はIssueの固定入力とし、入力対応、Claim・節単位の統合、Reliability判断、追加調査要求、最終Markdown契約、Skill実装を成果物にする。

### 妥当性監査結果

- 3件はIssueの独立Skill成果物と一対一で、L4親成果物を過不足なく満たす。
- L4-M1は記事中のSupport・一次情報・具体例の観測まで、L4-M3はそのSupportを認識更新の根拠として採用可能かというReliability判断を所有する。
- L3-M2-S2はCandidateから更新対象を特定する探索、L4-M2はArticle Target ClaimからユーザーKnowledge状態とGapを評価する探索であり、入力・判断・出力が異なる。
- L4-M2は局所的な探索継続・停止理由とraw Search Traceを所有する。Budget上限、再起動・再試行、Skill間の再調査、Run／Claim／Cycle相関はL5が所有する。
- 共有成果物契約、再調査、Search Trace、評価を独立した第4中分類にしない。各producer、L5、L6の既存責務と重複するためである。
- 独立した技術選定中分類は追加しない。Codex Skill、Markdown、Knowledge CLIという実行境界はIssueで確定済みであり、記事取得方法やCLI利用詳細は各中分類の小分類で実装前に設計する。

### 依存関係

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

- L2-M3-S5とS8はL4-M2の実CLI統合に対する直接依存、L2-M3-S9はL4 Release条件とする。
- L4-M3はL4-M2のMarkdown契約を利用するため、L4-M2の実CLI統合完了は待たない。

### L4 Skill Design Freeze Gate

3 SkillのMarkdown契約、正常・入力不足・部分失敗の状態、追加調査要求と報告signal、Support／Reliability境界、7状態を再判定しない境界、CLI Mapping、raw Search TraceとL5の相関責務境界、Path／Glob Ownerを確定する。

Gateは独立成果物を持たないMilestoneでありTask IDを付けない。L2実行ArtifactをGate条件に含めず、L4の探索要求からL2を設計する依存との循環を防ぐ。

### worktree所有境界

- L4-M1はArticle Analysis Skill DirectoryとClaims／Support契約・canonical例を単独所有する。
- L4-M2はKnowledge Search Skill DirectoryとAssessment契約、raw Trace event、CLI invocation referenceを単独所有する。
- L4-M3はReading Value Skill Directoryと最終Assessment契約を単独所有する。
- consumerは上流契約をread-onlyで利用する。共有Registry、Orchestration、Trace相関はL5、共有Fixture、Dataset、Harness、Agent評価、Acceptance、E2EはL6が所有する。
- 設計Mergeは`M1 → M2 → M3 → Gate`とし、独立部分は並行作業できる。Gate通過Commitを3実装worktreeの共通起点とし、Gate後は別Directoryで並行実装する。
- 実Path／GlobはGate記録で未確定項目なく確定する。

### 承認対象

1. L4親依存の細分化
2. L2-M2-S1、L2-M2親、L2 Design Freeze Gateの前方参照具体化
3. L4中分類3件
4. L4 Skill Design Freeze Gate
5. 依存DAGとworktree所有境界

### 承認結果（2026-08-11）

ユーザーはL4親依存の細分化、L2-M2-S1・L2-M2親・L2 Design Freeze Gateの前方参照具体化、L4中分類3件、L4 Skill Design Freeze Gate、依存DAG、worktree所有境界を承認した。

Task Mapへ確定構造を反映し、次の議論をL4-M1「Article Analysis Skill」の小分類へ移す。

## 決定 019: Article Analysis Skillの小分類

### 状態

承認済み（2026-08-11）

### 分割・監査方法

分割担当がIssue #1、承認記録、Task Mapを確認し、L4-M1を記事入力・取得方式、Claim分解、Location／Support追跡、Markdown／局所再分析結果契約、Skill実装の5件へ分割した。別の妥当性確認担当が、Issue既定事項の再議論、leafとしての完結性、L4-M2・M3、L5、L6との責務重複、技術選定の要否、依存DAG、L4 Gate、worktree所有を独立監査した。

監査の結果、5件を維持する。S1には利用するCodex取得手段の選定を含め、S4は局所再分析の結果契約へ限定する。既存承認済み構造の修正候補はない。

### L4-M1-S1

#### タスク名

`記事入力・取得方式・解析可否境界を詳細設計する`

#### 分割した理由

IssueはURLをArticle Analysisへ渡すことを定めているが、利用する取得手段、受理する入力、記事同一性、本文抽出範囲、取得不能・入力不足・部分取得をどう区別するかは未定義である。不完全な記事を完全に分析したように扱わないため、Claim抽出から独立した入力境界として分離する。

主成果物は`docs/design/article-analysis/article-input-boundary.md`とする。Codexの取得手段、Fallback、リダイレクト後URL・取得時点等のSource識別、動的ページ・認証／課金壁・本文欠落、正常・入力不足・取得不能・部分取得／部分解析の状態を定義する。

### L4-M1-S2

#### タスク名

`Article overview・Claim分解／正規化手順を詳細設計する`

#### 分割した理由

IssueはArticle overview、Claim、Role、ImportanceとClaim単位の検索を定めているが、overviewの範囲、複合Claimの分割、独立評価可能な粒度、重複・言い換えの統合、安定ID、記事内順序、Role・Importanceの適用手順は未定義である。記事を比較可能なClaim群へ変換する意味処理として分離する。

主成果物は`docs/design/article-analysis/overview-claim-decomposition.md`とする。Issue既定の項目とRole語彙は選び直さず、Claimの原子性、正規化、重複・包含、安定ID、順序を定義する。

### L4-M1-S3

#### タスク名

`Claim Location・記事内Support根拠の追跡方法を詳細設計する`

#### 分割した理由

Issueは各ClaimへLocationとSupportを付けることを定めているが、Claimと節・段落・引用・実装例・失敗例・計測・一次情報参照を再追跡可能に結び付ける方法は未定義である。「何をClaimとするか」を扱うS2とは変更理由と完了条件が異なり、M3のReliability最終判断との境界を明確にするため分離する。

主成果物は`docs/design/article-analysis/location-support-traceability.md`とする。Supportの有無、取得範囲外、確認不能、位置が不安定な場合も明示する。M1は記事内の観測事実だけを返し、根拠の十分性・信頼性・採用可能性は判断しない。

### L4-M1-S4

#### タスク名

`Article Analysis Markdown・局所再分析結果契約を定義する`

#### 分割した理由

IssueはMarkdownの主要項目を例示しているが、必須性、複数Claim境界、正常・入力不足・部分失敗、対象を限定した再分析結果の表現は未定義である。S1〜S3をM2・M3がread-onlyで利用できるproducer-owned契約へ固定するため分離する。

主成果物は`.agents/skills/article-analysis/references/article-analysis-markdown-contract.md`とする。対象Claim、確認範囲、観測追加、変更・不変、置換・追加関係を表現し、Claim IDを維持する。Routing、相関管理、再試行回数、Budget、終了判定はL5へ残す。

### L4-M1-S5

#### タスク名

`Article Analysis Skill instructions／referencesを実装し、固有構造契約を確認する`

#### 分割した理由

S1〜S4の確定設計を、記事取得、Overview作成、Claim分解、Location／Support観測、状態報告、局所再分析結果まで実行できるSkill packageへ変換する成果物が未実装である。意味的なAgent評価とは分け、Skill本体と密接な静的・構造Contract確認を同じ完了単位にする。

主成果物は`.agents/skills/article-analysis/SKILL.md`、S4以外の同Skill `references/**`、固有Contract確認資産とする。S1で選定した取得手段とFallbackを実装し、共有Dataset、実URL横断Scenario、Agent評価、Acceptance、E2EはL6へ残す。

### 妥当性監査結果

- S1は入力・取得状態、S2はClaim意味単位、S3は記事箇所とSupport根拠の追跡、S4は外部Markdown表現、S5は実行可能Skill packageを所有し、主成果物と完了条件が重複しない。
- S3をS2へ統合しない。Claimの意味分解と根拠追跡は変更理由が異なり、Support観測とReliability判断の境界が不明瞭になるためである。
- 独立した技術選定Taskは追加しない。Codex Skill内の取得手段はS1の狭い設計判断として完結し、S5で実装できる。利用可能な取得手段で成立しない場合だけ、L4 Gateで別基盤Taskの追加を再検討する。
- S4は局所再分析の結果だけを所有する。M3は追加調査要求を生成し、L5はRouting、再試行、Budget、終了を制御する。
- S5の確認はSkill固有のinstructions／references整合とMarkdown構造契約に限定し、共有評価はL6が所有する。
- 5件はいずれも単一主成果物で完了判定でき、これ以上子へ分割しないleafとして妥当である。

### 依存関係

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

- S1とS2は別成果物で並行可能とし、両方の統合後にS3、S4を直列Mergeする。
- L4-M2の技術非依存探索要求はS4前から先行できるが、Claim受入契約の確定はS4を待つ。
- L4-M3はS4をread-onlyで利用し、M1へRoutingする責務を持たない。
- S5はL2や実CLIを待たず、L4 Gate後に本番実装できる。

### worktree所有境界

- S1: `docs/design/article-analysis/article-input-boundary.md`
- S2: `docs/design/article-analysis/overview-claim-decomposition.md`
- S3: `docs/design/article-analysis/location-support-traceability.md`
- S4: `.agents/skills/article-analysis/references/article-analysis-markdown-contract.md`
- S5: `.agents/skills/article-analysis/SKILL.md`、S4を除く同Skillの`references/**`、固有Contract確認資産
- L5: 共有Registry、起動、Routing、再試行、Budget、終了、Run／Cycle相関
- L6: 共有Fixture、Dataset、Harness、Agent評価、Acceptance、E2E

S1とS2は、L1と承認済みPath Ownerを含むCommitを共通起点として別worktreeで並行する。Merge順は`S1・S2 → S3 → S4 → L4 Gate → S5`とする。S5はGate通過Commitから分岐し、S4の契約Fileをread-onlyで利用する。

### 承認対象

1. L4-M1小分類5件
2. L4-M1の依存DAGとL4 Gateへの接続
3. worktree所有境界とMerge順

### 承認結果（2026-08-11）

ユーザーはL4-M1小分類5件、L4-M1の依存DAGとL4 Gateへの接続、worktree所有境界とMerge順を承認した。

Task Mapへ確定構造を反映し、次の議論をL4-M2「Knowledge Search Skill」の小分類へ移す。

## 決定 020: Knowledge Search Skillの小分類

### 状態

承認済み（2026-08-11）

### 分割・監査方法

分割担当がIssue #1、承認記録、Task Mapを全文確認し、L4-M2を技術非依存探索要求、Article Claim受入、CLI Mapping、Knowledge Assessment判断、出力契約、Mock実装、実CLI統合の7件へ分割した。別の妥当性確認担当が、Issue既定事項の再議論、leafとしての完結性、L3-M2・L4-M3・L5・L6との重複、依存循環、Design Freeze Gate、worktree所有を独立監査した。

監査の結果、7件を維持する。S1とS2は「探索能力・遷移要求」と「特定Article Claimからの部分Claim・検索variant生成」、S4とS5は「意味判定規則」と「外部Envelope・直列表現」で変更理由と主成果物が異なるため統合しない。

### 先に提示する既存構造の修正候補

#### L2の暫定前方参照をL4-M2-S1へ確定する

`L2-M2`親の開始入力、`L2-M2-S1`の直接依存、`L2 Design Freeze Gate`の確認対象にある「L4-M2配下の技術非依存探索要求」を、`L4-M2-S1`の`agentic-search-requirements.md`へ置換する。L4-M2全体へ依存させると循環するため、L2へ先行投入できるS1だけを依存先にする。L2のタスク数と責務は変更しない。

#### L4 Skill Design Freeze Gateへ状態の直交性を追加する

Knowledge Assessmentの7状態と探索実行状態を別軸で確定する条件を追加する。

```text
Assessment status:
known / partially_known / inferable / contradicted /
outdated / no_evidence / uncertain

Execution status:
complete / input_insufficient / incomplete / partial_failure
```

不完全な探索や入力不足を`no_evidence`または`uncertain`へ混入すると、「未観測と未知を分離する」Invariantを破るためである。影響先はL4-M2-S4・S5、L4-M3、L5、L6であり、今回の小分類と同時に再承認する。

### L4-M2-S1

#### タスク名

`技術非依存な探索要求・Query Journey・局所停止理由を詳細設計する`

#### 分割した理由

Knowledge Searchの利用要求から検索基盤を逆算するには、Target Claimを調べるために必要な探索能力、探索段階の遷移、各段階で必要な取得情報、局所的な継続・停止理由をL2の技術選定より先に確定する必要がある。Issue既定の探索方式や停止原則を選び直さず、その適用箇所と記録表現を技術非依存の成果物へ落とすため分離する。

主成果物は`docs/design/knowledge-search/agentic-search-requirements.md`とする。Index・Embedding・Top-K等の物理方式はL2、Budget値・再試行・再起動・Workflow終了はL5へ残す。

### L4-M2-S2

#### タスク名

`Article Claim受入・Target Claim分解／検索variant生成手順を詳細設計する`

#### 分割した理由

L4-M1が生成したClaimを変更せず受け入れる境界と、特定Claimから検索用の部分Claim・Concept・表現variantを生成する処理は、検索基盤へ渡す一般的な探索能力・遷移要求とは入力と完了条件が異なる。元ClaimのIDと同一性を維持し、検索用表現を正式Article Claimへ昇格させない責務として分離する。

主成果物は`docs/design/knowledge-search/target-claim-query-reconstruction.md`とする。L4-M1-S4の正常・入力不足・部分解析状態を受け入れ、再分析の必要性はL5へ報告するだけに限定する。

### L4-M2-S3

#### タスク名

`CLI検索・取得Command Mappingと部分失敗引渡しを詳細設計する`

#### 分割した理由

Codexが選んだ探索操作をL1の検索・取得Commandへ変換する機械境界と、CLIの対象なし・部分結果・Errorを意味判定へ混入させず引き渡す規則は、Query生成やKnowledge Assessmentとは異なる。L1のSchemaを複製せず、各探索段階でRequest・Response・Pagination・Errorをどう利用するかを固定するため分離する。

主成果物は`.agents/skills/knowledge-search/references/cli-command-mapping.md`とする。`retryable`情報はL1-M2から受け取ってL5へ引き渡し、再試行判断は行わない。

### L4-M2-S4

#### タスク名

`Evidence意味比較・探索十分性／Knowledge Assessment判定手順を詳細設計する`

#### 分割した理由

Issueは7状態、Evidence Strength、Confidence等の意味を定めているが、取得済みEvidenceから意味的同等性、部分知識、推論可能性、矛盾、時点差、Known・Knowledge Gapを導く適用手順は未定義である。IssueとL1の語彙を再定義せず、Knowledge Search固有の意味判定規則として外部表現から分離する。

主成果物は`docs/design/knowledge-search/assessment-decision.md`とする。L3-M2-S2の更新対象探索、L4-M3のReading Value・Reliability判断は含めない。

### L4-M2-S5

#### タスク名

`Knowledge Assessment・raw Search Trace・局所再検索結果契約を定義する`

#### 分割した理由

正式Assessment、診断用raw Search Trace、局所再検索結果は同じClaimと探索結果を外部へ渡すproducer契約である一方、raw Traceを正式Assessmentへ混入させてはならない。S4の意味判定規則を変更せず、正常・入力不足・探索未完了・部分失敗を含む別Envelopeとして直列化する責務を分離する。

主成果物は`.agents/skills/knowledge-search/references/knowledge-search-output-contract.md`とする。局所再検索は初回検索と同じEnvelopeへ対象範囲・差分を付け、8番目のKnowledge状態を追加しない。Query／Step IDはM2が生成できるが、Run／Cycle IDはL5から受け取って返すだけとする。

### L4-M2-S6

#### タスク名

`Knowledge Search SkillをMock CLI契約に基づき実装し、固有構造契約を確認する`

#### 分割した理由

S1〜S5の確定設計を反復探索できる実行可能なSkill packageへ変換する成果物が未実装である。実CLI完成前でも同じ公開Command契約を持つMockでSkill固有のinstructions・出力構造・Command呼出しを確認できるため、実CLI Component統合から分離する。

主成果物は`.agents/skills/knowledge-search/SKILL.md`、S3・S5以外のSkill固有`references/**`、Mock構造Contract確認資産とする。共有Dataset、意味品質、Agent評価、Acceptance、E2EはL6へ残す。

### L4-M2-S7

#### タスク名

`Knowledge Search Skillを実CLI Artifactへ接続し、固有Command連携を確認する`

#### 分割した理由

Mockへ適合したSkillが実際の配布Artifactを通じて検索・取得Command、終了コード、JSON、Pagination、部分失敗を扱える状態は、S6とは依存時期と完了条件が異なる。CLI自体のBlack-box契約検証や横断Agent評価を重複せず、L4-M2固有のComponent連携として分離する。

主成果物は`tests/component/knowledge-search/**`とする。S7はS6で確定したCommand記述とL2 CLIをread-onlyで利用し、契約やSkill本体の変更が必要な場合は該当Ownerへ戻す。

### 妥当性監査結果

- 7件はいずれも単一の主成果物と完了判定を持ち、これ以上子へ分割しないleafとして妥当である。
- S1は探索能力・遷移・必要取得情報・停止理由の表現、S2は特定Claimからの検索variant生成を所有し、重複しない。
- S4は取得済みEvidenceの意味判定規則、S5はAssessment・Trace・実行状態の外部表現を所有し、重複しない。
- L3-M2-S2は更新Candidateから既存Knowledgeを特定し、L4-M2はArticle ClaimからユーザーのKnowledge状態を評価するため、入力・判断・出力が異なる。
- L4-M3は記事SupportのReliabilityとReading Value、L5はRouting・Budget・再試行・再起動・相関・Workflow終了、L6は共有Fixture・Dataset・Harness・Agent評価・Acceptance・E2Eを所有する。
- 独立した技術選定Taskは追加しない。Codex Skill、Markdown、Knowledge CLIは原典とL1で固定済みであり、検索・Index技術はL2-M2-S2が所有する。

### 依存関係

```text
Issue + L1-M1・L1-M2
          ↓
         S1 ─────────────→ L2-M2-S1 → L2 Design Freeze Gate

S1 + L4-M1-S4
          ↓
         S2
        ┌─┴──────────┐
        ↓            ↓
S3 CLI Mapping   S4 Assessment判断
  + L1-M2          + L1-M1
        └─────┬──────┘
              ↓
             S5
              ↓
L4-M1-S4 + S5 → L4-M3詳細設計
              ↓
M1〜M3詳細設計 → L4 Skill Design Freeze Gate
              ↓
             S6
              ↓
S6 + L2-M3-S5・S8 → S7
              ↓
M1〜M3実装 + L2-M3-S9 → L4 Release
```

S1とL4-M1-S4は別Pathで並行可能である。S3とS4はS2統合Commitから別worktreeへ分岐し、並行できる。L4 GateはM3詳細設計を待ち、L2の実Artifactは待たないため循環しない。

### worktree所有境界

- S1: `docs/design/knowledge-search/agentic-search-requirements.md`
- S2: `docs/design/knowledge-search/target-claim-query-reconstruction.md`
- S3: `.agents/skills/knowledge-search/references/cli-command-mapping.md`
- S4: `docs/design/knowledge-search/assessment-decision.md`
- S5: `.agents/skills/knowledge-search/references/knowledge-search-output-contract.md`
- S6: `.agents/skills/knowledge-search/SKILL.md`、S3・S5以外のSkill固有reference、Mock構造Contract確認資産
- S7: `tests/component/knowledge-search/**`
- L5: 共有Registry、Routing、Budget、再起動・再試行、Run／Claim／Cycle相関
- L6: 共有Fixture、Dataset、Harness、Agent評価、Acceptance、E2E

Merge順は`S1 → S2 → S3・S4 → S5 → L4 Gate → S6 → S7`とする。ただしS1とL4-M1-S4は並行し、S2は両方の統合後に確定する。S3・S4はS2統合Commit、S6はGate通過Commit、S7はS6とL2-M3-S8を含むCommitを共通起点とする。

### 承認対象

1. L2の前方参照をL4-M2-S1へ確定する修正
2. L4 GateへAssessment状態と実行状態の直交性を追加する修正
3. L4-M2小分類7件
4. 依存DAGとL4 Gateへの接続
5. worktree所有境界とMerge順

### 承認結果（2026-08-11）

ユーザーは、L2の前方参照を`L4-M2-S1`へ確定する修正、L4 GateへAssessment状態と探索実行状態の直交性を追加する修正、L4-M2小分類7件、依存DAG、worktree所有境界とMerge順を承認した。

Task Mapへ確定構造を反映し、次の議論をL4-M3「Reading Value Skill」の役割・必要性の確認へ移す。

## 決定 021: Reading Value Skillの小分類

### 状態

承認済み（2026-08-11）

### 分割・監査方法

分割担当がIssue #1、承認記録、Task Mapを全文確認し、L4-M3を入力整合、Recognition Gain、Reliability／Applicability、Attention Cost・推奨範囲、最終出力・追加調査要求契約、Skill実装の6件へ分割した。別の妥当性確認担当が、Issue既定事項の再議論、leafとしての完結性、L4-M1・M2、L5、L6との重複、依存DAG、L4 Gate、worktree所有を独立監査した。

監査の結果、6件を維持する。S2とS3はユーザーKnowledgeとの差分と記事側根拠・利用可能性、S4とS5は統合判断手順と外部結果契約で変更理由が異なるため統合しない。既存承認済み構造の修正候補はない。

### L4-M3-S1

#### タスク名

`Article Analysis／Knowledge Search入力整合・評価可否境界を詳細設計する`

#### 分割した理由

M1とM2は別producerが契約を所有するため、同一Claimの対応付け、欠落・重複・部分結果、評価可能範囲の扱いが未定義である。Knowledge Assessmentの7状態と探索実行状態を混同せず、評価ロジックへ渡す前に入力の整合性と評価可否を確定する責務として分離する。

主成果物は`docs/design/reading-value/input-alignment-and-evaluability.md`とする。M1のcanonical Claim IDで結合し、M2の部分Claim・検索variantをArticle Claimへ昇格させない。Claim単位の評価可能、限定評価可能、評価不能を区別し、Routingは行わない。

### L4-M3-S2

#### タスク名

`Knowledge AssessmentからRecognition Gain Candidateへの適用手順を詳細設計する`

#### 分割した理由

Issueは7状態と7種のRecognition Gainを定めているが、Known、Knowledge Gap、Role、ImportanceからClaim単位のGain Candidateへ変換する具体的な適用手順は未定義である。M2のKnowledge状態を再判定せず、ユーザー固有の認識差分へ変換する責務として分離する。

主成果物は`docs/design/reading-value/recognition-gain-application.md`とする。NoveltyとImportanceを根拠付きで適用し、`known`でもStructural Gainがあり得ること、`no_evidence`を未知と断定しないこと、単一数値Scoreへ変換しないことを維持する。

### L4-M3-S3

#### タスク名

`Claim・記事内SupportからReliability／Applicabilityを判断する手順を詳細設計する`

#### 分割した理由

M1はSupportを観測するだけであり、その根拠を認識更新へ採用可能か判断する責務はM3にある。ユーザーKnowledgeとの差分を扱うS2とは入力と完了条件が異なるため、記事側のEvidence Quality、Reliability、Applicabilityを判断する手順として分離する。

主成果物は`docs/design/reading-value/reliability-applicability.md`とする。M1のSupport観測を変更せず、一次情報、実装例、失敗例、計測等を区別する。一般的な記事品質や単一Reliability Scoreを主出力にせず、観測できない個人的関心を推測しない。

### L4-M3-S4

#### タスク名

`Attention Cost・Claim横断統合によるRecommendation／読書範囲決定手順を詳細設計する`

#### 分割した理由

Claim単位のGain CandidateとReliabilityだけでは、記事全体を読むか、特定箇所だけ読むか、飛ばすかを決められない。Claimの分布、Role、Location、前提・文脈依存とAttention Costを記事全体で統合する最終判断責務として分離する。

主成果物は`docs/design/reading-value/recommendation-range-aggregation.md`とする。Issue既定の`read_full / read_selected / skip`を選び直さず、`read_selected`の飛び地、前提節、重複・隣接範囲を扱う。Attention Costは利用可能な入力による定性的proxyとし、精密な読了時間や数値式を必須にしない。

### L4-M3-S5

#### タスク名

`追加調査要求・Reading Value Assessment Markdown契約を定義する`

#### 分割した理由

最終Assessmentと追加調査要求は同じReading Value producerからL5へ渡す排他的な結果形態であり、別タスクにするとClaim ID・理由・状態の表現が二重所有になる。S4の判断手順を変更せず、終端結果と非終端要求を外部へ渡す契約として分離する。

主成果物は`.agents/skills/reading-value/references/reading-value-output-contract.md`とする。`final_assessment`と`additional_research_request`を排他的にし、追加要求を第4のRecommendationにしない。Routing、Budget、再試行、Run／Cycle相関、終了判断はL5へ残す。

### L4-M3-S6

#### タスク名

`Reading Value Skill instructions／referencesを実装し、固有構造契約を確認する`

#### 分割した理由

S1〜S5の確定設計を、M1・M2成果物の受入、Gain・Reliability・Attention Cost評価、追加調査要求、最終Assessment生成まで実行できるSkill packageへ変換する成果物が未実装である。意味品質の横断評価とは分け、Skill本体と密接な構造Contract確認を同じ完了単位にする。

主成果物は`.agents/skills/reading-value/SKILL.md`、S5以外の同Skill `references/**`、固有構造Contract確認資産とする。Knowledge CLIを直接呼ばず、共有Registry、Routing、Budget、相関、共有評価資産を所有しない。

### 妥当性監査結果

- 6件はいずれも単一主成果物と独立した完了判定を持ち、これ以上子へ分割しないleafとして妥当である。
- S1はproducer間の入力整合、S2はユーザーKnowledge差分、S3は記事側根拠と利用可能性、S4は記事横断の最終判断、S5は外部結果契約、S6は実行可能Skillを所有する。
- Knowledge Assessmentの7状態、Recognition Gain 7種、3 Recommendation、Markdown主要SectionはIssueの固定入力とし、再選定しない。
- Reading ValueはKnowledge CLIを直接利用しないため、Mock CLI実装や実CLI統合タスクを追加しない。
- L5は起動、Routing、再調査回数、Budget、Run／Cycle相関、終了を所有する。L6は共有Fixture、Dataset、Harness、Agent評価、Acceptance、E2Eを所有する。
- 独立した技術選定タスクは追加しない。Codex SkillとMarkdown境界はIssueで確定済みである。
- 既存承認済み構造に追加修正はない。L4 Gateの既承認条件を今回の成果物で具体化できる。

### 依存関係

```text
L4-M1-S4 + L4-M2-S5
          ↓
S1 入力整合・評価可否
     ┌────┴────┐
     ↓         ↓
S2 Gain     S3 Reliability・Applicability
     └────┬────┘
          ↓
S4 Attention Cost・推奨・範囲統合
          ↓
S1 + S4 → S5 最終結果／追加調査要求契約
          ↓
M1〜M3詳細設計
          ↓
L4 Skill Design Freeze Gate
          ↓
S6 Reading Value Skill実装
          ↓
L5統合・L6 Agent評価
```

- S2とS3はS1統合Commitから別worktreeへ分岐して並行できる。
- S5は評価不能時の結果も直列化するため、S4だけでなくS1を直接入力にする。
- S6はL2の実Artifactを待たず、L4 Gate通過Commitから実装できる。
- L4 Releaseは従来どおりM1・M2・M3実装とL2-M3-S9を待つ。

### L4 Gate確認事項

- Assessment 7状態と探索実行状態が直交している。
- M3がM2の7状態を再判定しない。
- M1のSupport観測とM3のReliability判断が分離されている。
- Attention Costが利用可能な入力だけに基づく。
- 最終Assessmentと追加調査要求が排他的である。
- 追加調査要求をRecommendationへ混入させない。
- Routing、Budget、相関、終了がL5に残っている。

これらは既存Gateの具体化であり、独立タスクや承認済み構造の修正ではない。

### worktree所有境界

- S1: `docs/design/reading-value/input-alignment-and-evaluability.md`
- S2: `docs/design/reading-value/recognition-gain-application.md`
- S3: `docs/design/reading-value/reliability-applicability.md`
- S4: `docs/design/reading-value/recommendation-range-aggregation.md`
- S5: `.agents/skills/reading-value/references/reading-value-output-contract.md`
- S6: `.agents/skills/reading-value/SKILL.md`、S5以外のSkill固有reference、固有構造Contract資産

Merge順は`L4-M1-S4・L4-M2-S5 → S1 → S2・S3 → S4 → S5 → L4 Gate → S6`とする。S6はGate通過Commitを起点とし、M1・M2契約とS5をread-onlyで使用する。共有Registry、Routing、Budget、相関はL5、共有評価資産はL6の単独Ownerとする。

### 承認対象

1. L4-M3小分類6件
2. 依存DAGとL4 Gateへの接続
3. worktree所有境界とMerge順

### 承認結果（2026-08-11）

ユーザーは、L4-M3小分類6件、依存DAG、L4 Gateへの接続、worktree所有境界とMerge順を承認した。

既存承認済み構造への追加修正はなく、Task Mapへ確定構造を反映した。次の議論を`L5`「Parent Orchestration Skillと全体Workflowを設計・実装する」の中分類へ移す。

## 決定 022: Parent Orchestration Skillと全体Workflowの中分類

### 状態

承認済み（2026-08-11）

### 分割・監査方法

分割担当がIssue #1、承認記録、Task Mapを全文確認し、L5を共通Orchestration設計、記事価値判定Workflow設計、Knowledge蓄積Workflow設計、Parent Skill実装・Component統合の4件へ分割した。別の妥当性確認担当が、Issue既定事項の再議論、M1〜M3とM4の工程粒度、Mockと実Component統合の分離位置、両Workflowを接続する第5中分類の要否、L3・L4・L6との責務境界、依存循環、Gate、worktree所有を独立監査した。

監査の結果、中分類4件を維持する。M4のMock実装と実Component統合は依存時期と所有資産が異なるが、同じParent Skillの完成責任に属するため、M4の小分類で必ず分離する。記事価値判定とKnowledge蓄積は共通entrypointから選択される別Workflowであり、直結するPipelineではないため、接続専用の第5中分類は追加しない。

### L5-M1

#### タスク名

`Parent entrypoint・共通Orchestration状態／実行制御契約を詳細設計する`

#### 分割した理由

記事価値判定とKnowledge蓄積では処理順が異なる一方、entrypoint、Skill Registry、成果物引渡し、状態分岐、Budget、retry／restart／reroute／stop、Run／Claim／Cycle相関、終了理由は共通である。各Workflowへ重複させず、各Skillの意味判断を変更しない共通Control Planeとして分離する。

主成果物は`docs/design/orchestration/common-control-contract.md`とする。URL／Episodeのentrypoint、Workflow識別、契約Version不一致、全実行状態の遷移、Budget範囲・消費、ID所有者・伝播、Search Trace相関、Knowledge Update再試行時の副作用重複防止を定義する。Assessment 7状態、Recommendation、更新判断はproducerの固定入力とし、L5で再判定しない。

着手はIssue #1とL1から可能だが、確定・MergeにはL3-M1-S3、L3-M2-S3・S4、L4-M1-S4、L4-M2-S5、L4-M3-S5のproducer契約を必要とする。後続はL5-M2〜M4とL6である。

### L5-M2

#### タスク名

`記事価値判定Workflowの起動・Claim fan-out／結果統合・再調査・終了を詳細設計する`

#### 分割した理由

URLからArticle Analysis、Claim単位のKnowledge Search、Reading Valueへ進む処理には、canonical Claim IDを維持したfan-out／fan-in、部分結果の集約、追加再分析・再検索のRouting、最終Assessmentという固有制御がある。共通実行制御や各Skillの意味判断とは変更理由と完了条件が異なるため分離する。

主成果物は`docs/design/orchestration/article-reading-workflow.md`とする。`URL → Article Analysis → Claim別Knowledge Search → Reading Value`の状態遷移、`final_assessment`と`additional_research_request`の排他制御、Budget到達・入力不足・部分失敗・評価不能時の終端を定義する。Knowledge Assessment、Reliability、Recommendationを再判定せず、記事評価結果だけではKnowledge蓄積を起動しない。

L5-M1の共通契約案から着手し、確定・MergeはL5-M1、L4-M1-S4、L4-M2-S5、L4-M3-S5へ依存する。L5-M3とは別設計Fileで並行可能であり、後続はL5-M4とL6である。

### L5-M3

#### タスク名

`Knowledge蓄積WorkflowのEpisode起動・Candidate分岐・更新終了を詳細設計する`

#### 分割した理由

EpisodeからKnowledge Acquisitionを起動し、候補がある場合だけKnowledge Updateへ渡すWorkflowは、記事評価とは入口、正常終了条件、部分成功後の再試行、副作用が異なる。`no-candidate`、入力不足、正常な非永続化、部分成功、CLI失敗を混同せず、ExposureをKnowledge更新へ自動接続しない制御として分離する。

主成果物は`docs/design/orchestration/knowledge-accumulation-workflow.md`とする。意味のあるEpisode終了を起点とする条件、`candidate / no-candidate / input_insufficient`、Candidate別のUpdate起動、成功済み更新を重複させない再試行、終了状態を定義する。Acquisition／Updateの意味判断と4更新操作をL5で再判定しない。

L5-M1の共通契約案から着手し、確定・MergeはL5-M1、L3-M1-S3、L3-M2-S3・S4へ依存する。L5-M2とは別設計Fileで並行可能であり、後続はL5-M4とL6である。

### L5-M4

#### タスク名

`Parent Orchestration Skillを実装し、Mock契約確認後にL3／L4実Componentを統合する`

#### 分割した理由

M1〜M3の設計を実行可能な単一Parent Skillへ変換し、実際のL3・L4 Skillsと接続する成果物が未実装である。共有entrypoint、Registry、Parent packageを複数中分類が編集すると制御の二重所有とworktree競合が生じるため、Parent Skillの完成責任を一つへ集約する。

主成果物は`.agents/skills/parent-orchestration/**`とParent固有Component連携確認資産とする。小分類では、`Mock Component契約によるSkill実装・固有構造契約確認`と、`L3／L4実Component接続・連携確認`を依存の異なる2工程へ必ず分離する。前者はL5 Design Freeze Gate後に開始し、後者はL3・L4 Release後に行う。共有Fixture、Dataset、Harness、Agent評価、Acceptance、E2EはL6へ残す。

### 妥当性監査結果

- 4件は共通制御設計、記事Workflow設計、Knowledge蓄積Workflow設計、単一Parent Skillの実装・統合という異なる主成果物を持ち、親の範囲を満たす。
- M1〜M3を実装前に確定し、L5 Design Freeze Gate後にM4へ進むため、設計変更による実装手戻りを抑止できる。
- M4のMock実装と実Component統合は中分類を増やさず、小分類で別leafにする。これによりMock実装をL3・L4 Release前に先行できる。
- 2つのWorkflowを直結する第5中分類は追加しない。Issue 6.5の`Exposure is not Knowledge Acquisition`により、記事評価やAI説明だけでKnowledge更新を起動してはならない。
- 独立した技術選定タスクは追加しない。Codex Skill、Markdown受渡し、Knowledge CLI境界はIssueで確定済みである。
- L5はKnowledge CLIを直接呼ばず、producerのAssessment、Recommendation、更新判断を再判定しない。
- L6は共有Fixture、Dataset、Harness、Agent評価、Acceptance、E2Eを所有し、L5-M4はParent固有の構造・状態遷移・Component連携確認に限定する。

### 依存関係

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
                    ↓
          L3 Release + L4 Release
                    ↓
        L5-M4 実Component統合
                    ↓
        L6 最終Acceptance・E2E
```

L3・L4のproducer契約を入力にL5を設計し、L5をL3・L4 Release条件へ戻さないため循環しない。L3・L4 ReleaseはL2-M3-S9を含むため、L5実統合からL2-M3-S9への直接依存は追加しない。

### L5 Design Freeze Gate

独立成果物のないMilestoneとして、entrypoint・Workflow識別・Registry、全状態遷移、ID所有・相関、Budget、retry／restart／reroute／stopと副作用重複防止、Claim fan-out／fan-in、最終結果と追加調査要求の排他性、no-candidateと正常な非永続化、記事評価からKnowledge更新を自動起動しないInvariant、Path Owner、Mock差替え境界を実装前に確定する。独立Task IDは付与しない。

### worktree所有境界

- L5-M1: `docs/design/orchestration/common-control-contract.md`
- L5-M2: `docs/design/orchestration/article-reading-workflow.md`
- L5-M3: `docs/design/orchestration/knowledge-accumulation-workflow.md`
- L5-M4 Mock実装: `.agents/skills/parent-orchestration/**`とParent固有Contract確認資産
- L5-M4実統合: `tests/component/parent-orchestration/**`
- L3・L4 Skill packageとL2 CLIはread-only入力とする。
- L6だけが共有Fixture、Dataset、Harness、Acceptance、E2Eを変更する。

M2・M3はM1統合Commitから別worktreeへ分岐して並行可能とする。Merge順は`producer契約 → M1 → M2・M3 → L5 Gate → M4 Mock実装 → L3・L4 Release → M4実統合 → L6`とする。

### 既存構造の修正候補

#### 対象

`L5`の依存表現

#### 現状

`設計着手: L1。完了・Merge: L2-M3、L3、L4`

#### 問題

共通設計の着手、producer契約による確定、Gate後のMock実装、実Component統合、親完了の依存が一つの表現へ混在している。また、L2-M3はL3・L4 Releaseに既に含まれるため、実統合からの直接依存は重複する。

#### 修正案

`設計着手: Issue・L1。設計確定: L3・L4 producer契約。Mock実装: L5 Design Freeze Gate。実Component統合・完了: L3 Release・L4 Release。後続: L6最終Acceptance・E2E`へ細分化する。

#### 理由

開始依存と完了依存を分離し、設計とMock実装をL3・L4の全実装完了前に先行させつつ、実Component未完成のままL5を完了させないためである。

#### 影響範囲・再承認単位

L5の依存欄と全体DAGだけを変更する。L3・L4・L6の責務・親子構造は変更しない。L5中分類4件、L5 Gate、依存DAG、worktree境界と一括で再承認する。

### 承認対象

1. L5中分類4件
2. L5 Design Freeze Gate
3. 依存DAGとworktree所有・Merge順
4. L5親依存表現の修正案

### 承認結果（2026-08-11）

ユーザーは、L5中分類4件、L5 Design Freeze Gate、依存DAG、worktree所有・Merge順、L5親依存表現の修正案を承認した。

Task Mapへ確定構造と工程別依存を反映し、次の議論を`L5-M1`「Parent entrypoint・共通Orchestration状態／実行制御契約を詳細設計する」の小分類へ移す。

## 決定 023: Parent entrypoint・共通Orchestration状態／実行制御契約の小分類

### 状態

承認済み（2026-08-11）

### 分割・監査方法

分割担当がIssue #1、承認記録、Task Mapを全文確認し、L5-M1を外部受付、子Skill境界、共通実行状態、相関、副作用安全性を含む実行制御の5件へ分割した。別の妥当性確認担当が、Issue既定事項の再議論、leafとしての独立完了性、S1とS2、S3とS4、S4とS5の境界、L5-M2・M3・M4およびL6との責務重複、依存DAG、Gate、worktree所有を独立監査した。

監査の結果、5件を維持する。S4は記事側だけでなくKnowledge蓄積側の部分成功も追跡できるよう`Candidate`を対象へ追加する。S5はL5単独でexactly-onceを保証する表現を避け、L1の冪等性・重複時の論理効果とproducer結果契約を入力にした`副作用安全性`として定義する。

### L5-M1-S1

#### タスク名

`Parent entrypoint・Workflow識別／起動受付契約を定義する`

#### 分割した理由

IssueはURL評価とEpisode処理の2 Workflowを定めているが、Parentが入力をどう識別し、欠落・競合・未対応入力をどう扱い、どのWorkflowのRun開始要求へ変換するかは未定義である。子Skillの選択・起動やWorkflow内の処理順とは異なる外部受付境界として分離する。

主成果物は`docs/design/orchestration/common-control-contract.md`とする。URL／Episodeの排他的識別、必須入力、受付・拒否、未知Workflow、契約Version不一致、Run初期化要求を規範例付きで定義し、S2〜S5の契約文書を参照するumbrella契約としてS1だけが編集する。記事評価やAI説明へのExposureだけでKnowledge蓄積を起動しない。

### L5-M1-S2

#### タスク名

`子Skill Registry・起動／成果物引渡し契約を定義する`

#### 分割した理由

Issueは5つの子SkillとMarkdown受渡しを定めているが、Skill識別、能力、契約Version、起動要求、成果物・診断・部分結果の引渡し、未登録・非互換時の扱い、Mock差替え境界は未定義である。S1の外部受付と分け、producer-owned成果物を変更・再判定しないParentから子Skillへの内部Invocation境界として分離する。

主成果物は`docs/design/orchestration/child-skill-registry-contract.md`とする。L3-M1-S3、L3-M2-S3・S4、L4-M1-S4、L4-M2-S5、L4-M3-S5をread-only入力とし、Workflow固有の起動順とRoutingはL5-M2・M3、実RegistryはL5-M4へ残す。

### L5-M1-S3

#### タスク名

`共通実行状態・診断Envelope契約を定義する`

#### 分割した理由

L3・L4は異なる正常・非実行・失敗状態を返すため、単純な成功／失敗へ潰すと`no-candidate`、正常な非永続化、入力不足、探索未完了、部分成功を誤って扱う。Budgetや再実行の判断より先に、producerの意味状態を保ったままParent Control Planeへ載せる状態・診断契約として分離する。

主成果物は`docs/design/orchestration/execution-state-envelope.md`とする。Parent実行状態とproducer Domain結果を直交させ、正常、入力不足、評価不能、中断、取消、Budget停止、部分失敗、失敗、成功済み・未処理範囲、原因、retryable情報、原結果参照を定義する。状態に応じていつ再実行・終了するかはS5へ残す。

### L5-M1-S4

#### タスク名

`Run／Claim／Candidate／Cycle・Search Trace相関契約を定義する`

#### 分割した理由

Issueは再調査とSearch Traceを要求するが、Run、再調査Cycle、Claim別探索、Candidate別更新、Attemptの対応関係とID所有者は未定義である。状態遷移や再実行判断とは分け、複数Workflowを横断する成果物系譜と診断可能性を保証する相関契約として分離する。

主成果物は`docs/design/orchestration/correlation-trace-contract.md`とする。L4-M1のcanonical Claim ID、L4-M2のQuery／Step ID、L3のCandidate参照を変更せず、L5が所有するRun／Cycle／Attempt／Work Itemとの生成・伝播・再利用・新規発行条件、fan-out／fan-in、再調査結果、Candidate更新、診断Envelopeとの結合を定義する。raw Search Traceの内容はL4-M2、評価利用はL6へ残す。

### L5-M1-S5

#### タスク名

`Budget・retry／restart／reroute／stop・副作用安全性契約を定義する`

#### 分割した理由

無限探索防止には、単なる最大回数ではなく、Budgetの範囲・消費単位・有限な上限と、状態・診断・相関情報からretry、restart、reroute、stopを選ぶ規則が必要である。副作用を伴う更新の再実行条件も同じAttemptとBudgetを利用するため、独立6件目にせず共通実行制御へ含める。

主成果物は`docs/design/orchestration/execution-control-policy.md`とする。Run／Workflow／Claim／Candidate／Cycle等のBudget、再実行判定表、取消・中断、成功済みWork Itemの除外、結果不明時の照合・停止、終了理由を定義する。L1-M1-S4の冪等性・重複時の論理効果を入力とし、CLIのexactly-onceや冪等性を再設計しない。記事の再調査配送先はL5-M2、Candidate別の再開順序はL5-M3へ残す。

### 妥当性監査結果

- 5件は、Parent外部受付、子Skill内部境界、共通状態、相関、実行制御という異なる単一成果物と完了判定を持ち、これ以上子へ分割しないleafとして妥当である。
- S1とS2は外部entrypointと子Skill Invocation、S3とS4は状態と識別・追跡情報、S4とS5は相関の事実と相関を用いる制御判断として分離できる。
- Workflow固有のClaim fan-out／fan-in、追加調査配送、Candidate分岐、終了集約はL5-M2・M3へ残す。
- Skill package、Registry実体、Mock・実Component接続はL5-M4へ残す。
- 規範例は各設計契約へ含められるが、共有Fixture、Dataset、Harness、Agent評価、Acceptance、E2EはL6へ残す。
- 独立した技術選定、実装、検証タスクは追加しない。Codex Skill、Markdown受渡し、Knowledge CLI境界はIssueで確定済みであり、L5-M1は詳細設計だけを所有する。
- 既存承認済み階層、L5 Gate、L3・L4・L6の責務に追加修正はない。

### 依存関係

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

S1はIssue・L1から先行でき、S2はL3・L4 producer契約を待って確定する。S3・S4はS1・S2の統合Commitから別worktreeへ分岐して並行できる。S5は状態、相関、L1の冪等性契約を同時に利用するため直列化する。M2・M3固有設計をS5の入力へ戻さないため循環しない。

### L5 Design Freeze Gateへの接続

新しいGateは追加しない。既存L5 Gateで、外部entrypointと子Skill Registryの境界、producer意味状態とParent実行状態の直交、全ID Owner、Search Traceの非改変、Budgetと再実行判断の共通範囲、副作用不明時の非盲目的再実行、記事評価からKnowledge蓄積を自動起動しないInvariant、Mock差替え境界、S1〜S5のPath Ownerを確認する。

### worktree所有境界

- S1: `docs/design/orchestration/common-control-contract.md`
- S2: `docs/design/orchestration/child-skill-registry-contract.md`
- S3: `docs/design/orchestration/execution-state-envelope.md`
- S4: `docs/design/orchestration/correlation-trace-contract.md`
- S5: `docs/design/orchestration/execution-control-policy.md`
- L5-M2・M3: 各Workflow設計Fileのみ
- L5-M4: `.agents/skills/parent-orchestration/**`とParent固有Component連携資産
- L6: 共有Fixture、Dataset、Harness、Agent評価、Acceptance、E2E

S1は先行Mergeできる。S2の依存が揃った後、`S1・S2 → S3・S4 → S5 → L5-M2・M3 → L5 Gate`の順でMergeする。S3・S4はS1・S2統合Commitを共通起点とする。各Taskは他Taskの契約Fileをread-onlyで参照する。

### 承認対象

1. L5-M1小分類5件
2. 依存DAGとL5 Design Freeze Gateへの接続
3. worktree所有境界とMerge順

### 承認結果（2026-08-11）

ユーザーは、L5-M1小分類5件、依存DAG、L5 Design Freeze Gateへの接続、worktree所有境界とMerge順を承認した。

Task Mapへ確定構造を反映し、次の議論を`L5-M2`「記事価値判定Workflowの起動・Claim fan-out／結果統合・再調査・終了を詳細設計する」の小分類へ移す。ユーザーの指示により、L5-M2からL5-M4までは途中承認を待たず連続して分割・独立監査し、まとめて提示する。

## 決定 024: 記事価値判定Workflowの小分類

### 状態

承認済み（2026-08-11）

### 分割・監査方法

分割担当がIssue #1、承認記録、Task Mapを照合し、L5-M2を記事Workflow開始、Claim fan-out、結果fan-in、追加調査Routing、終端の5件へ分割した。別の妥当性確認担当が、Issue内での定義不足、leaf原子性、L5-M1共通制御とL4各Skillの意味判断との境界、依存DAG、L5 Gate、worktree所有を独立監査した。

監査の結果、5件を維持する。S1は共通Run初期化を再定義せず記事Workflow固有開始へ限定し、S3はKnowledge Assessmentを再判定しない機械的な相関・集約とする。S4はL5-M1-S5が許可した後の固有Routing、S5はL4-M3が決めたAssessmentの受入とWorkflow終端だけを所有する。

### L5-M2-S1

#### タスク名

`記事価値判定Workflowを開始し、Article Analysisを起動・初期成果物を受け入れる`

#### 分割した理由

L5-M1は外部受付と共通Run初期化を定義するが、受付済み記事RunからArticle Analysisを起動し、正常、入力不足、取得不能、部分解析を受け入れてClaim処理へ移る記事固有遷移は未定義である。Claim別並行処理が始まる前の一対一境界として分離する。

主成果物は`docs/design/orchestration/article-reading-workflow.md`とする。記事Workflow全体を参照するumbrella文書としてS1が単独所有し、canonical Claim IDを維持し、記事評価からKnowledge蓄積を起動しないことを完了条件に含める。

### L5-M2-S2

#### タスク名

`canonical Claim単位のKnowledge Search fan-out契約を定義する`

#### 分割した理由

IssueはClaim単位の探索を定めるが、Article Analysis成果物を一意な検索Work Itemへ展開する規則、0件・無効Claim、重複防止、順序・並行性、相関、部分解析時の未処理範囲は未定義である。一対多の制御として起動・結果集約から分離する。

主成果物は`docs/design/orchestration/article-reading/article-search-fanout.md`とする。検索variantや部分Claimをcanonical Claimへ昇格させず、Knowledge Assessmentや検索戦略をL5で判断しない。

### L5-M2-S3

#### タスク名

`Claim別Knowledge Search結果を相関・集約し、Reading Valueへ引き渡す`

#### 分割した理由

複数検索Work Itemの成功、欠落、重複、遅延、部分結果をcanonical Claimへ対応付け、評価可能範囲を保持した入力BundleとしてReading Valueへ渡す多対一制御は未定義である。Knowledge Assessmentの7状態と実行状態を混同せず、終端判断を行わないfan-inとして分離する。

主成果物は`docs/design/orchestration/article-reading/article-search-fanin-handoff.md`とする。Article Analysis結果とKnowledge AssessmentをClaimごとに対応付け、L4-M3の受入契約へ引き渡す。

### L5-M2-S4

#### タスク名

`追加再分析／再検索要求のRouting・再調査Cycle反映を定義する`

#### 分割した理由

IssueはReading ValueからArticle AnalysisまたはKnowledge Searchへ再調査できることを定めるが、要求対象の配送、局所結果の受入、次Cycleへの反映と再集約は未定義である。Budget・reroute可否はL5-M1-S5に残し、許可後の記事Workflow固有循環制御として分離する。

主成果物は`docs/design/orchestration/article-reading/article-research-cycle-routing.md`とする。Cycle／Attemptと元成果物の系譜を維持し、各Skillの局所再分析・再検索結果を意味変更せずS2・S3の経路へ戻す。

### L5-M2-S5

#### タスク名

`最終Assessment受入・異常／Budget到達時のWorkflow終端契約を定義する`

#### 分割した理由

正常な最終Assessmentと、入力不足、評価不能、部分失敗、取消、Budget到達では、返却可能な成果物・診断・未評価範囲が異なる。Reading Valueの`skip`を異常停止と混同せず、すべての実行経路を有限に閉じる責務として分離する。

主成果物は`docs/design/orchestration/article-reading/article-workflow-termination.md`とする。L4-M3の`final_assessment`を再判定せず受け入れ、`additional_research_request`との排他性、残存Work Itemとの整合、部分成果物と終了理由を定義する。

### 妥当性監査結果

- 5件はいずれもIssue 8.1、16、29、37〜38、43、49、52章が要求する未定義のWorkflow詳細であり、`keep_missing_detail`と判定した。
- S1とS2は一対一の初期起動と一対多のWork Item展開、S2とS3はfan-outとfan-in、S3とS5は非終端集約とWorkflow終端、S4とL5-M1-S5は固有配送と共通実行方針として分離できる。
- Claim分解、Knowledge Assessment、Reliability、Recommendation、追加調査内容の意味判断をL5へ移さない。
- 実装はL5-M4、共有Fixture・Dataset・Harness・Acceptance・E2EはL6へ残す。独立技術選定は不要である。

### 依存関係

```text
L5-M1 + L4-M1-S1・S4
              ↓
S1 Workflow開始・Article Analysis受入
              ↓
L4-M2-S2 + S2 Claim fan-out
              ↓
L4-M2-S5 + L4-M3-S1 + S3 結果fan-in・引渡し
              ↓
L4-M1-S4 + L4-M2-S5 + L4-M3-S5 + S4 再調査Routing
              ↓
S1〜S4 + L5-M1-S3・S5 + L4-M3-S5
              ↓
S5 最終Assessment受入・終端
              ↓
L5 Design Freeze Gate
```

Taskの実施DAGは直列とし、実行時にS4からS2・S3相当の状態へ戻ることはWorkflow状態遷移であってTask依存の循環ではない。

### L5 Design Freeze Gateへの接続

新しいGateは追加しない。既存Gateで、canonical Claim ID不変、fan-out／fan-inの欠落・重複防止、producer意味状態と実行状態の直交、`final_assessment`と`additional_research_request`の排他性、有限な再調査Cycle、Budget拒否時の終端、`skip`と異常停止の区別、記事評価からKnowledge蓄積を起動しないInvariantを確認する。

### worktree所有境界

- S1: `docs/design/orchestration/article-reading-workflow.md`
- S2: `docs/design/orchestration/article-reading/article-search-fanout.md`
- S3: `docs/design/orchestration/article-reading/article-search-fanin-handoff.md`
- S4: `docs/design/orchestration/article-reading/article-research-cycle-routing.md`
- S5: `docs/design/orchestration/article-reading/article-workflow-termination.md`

共通起点はL5-M1と必要なL4契約を含むCommitとし、Merge順は`S1 → S2 → S3 → S4 → S5 → L5 Gate`とする。各Taskは単独Fileを所有するため別worktreeで先行ドラフト可能だが、確定とMergeは依存元の統合後に行う。親成果物はumbrella文書と4参照文書からなるdocument setとして扱う。

### 既存構造の修正候補

#### 対象

`L5-M2`の直接依存

#### 現状

`L5-M1、L4-M1-S4、L4-M2-S5、L4-M3-S5`

#### 問題

最終出力契約だけでは、Article Analysisの起動入力、Knowledge SearchのClaim受入、Reading Valueの入力整合に隠れ依存が残る。

#### 修正案

現行依存へ`L4-M1-S1、L4-M2-S2、L4-M3-S1`を追加する。

#### 理由・影響範囲・再承認単位

いずれもread-onlyの入力契約で循環は生じない。L5-M2台帳・DAG・Gate入力だけを変更し、階層数、L4責務、L5-M1は変更しない。L5-M2小分類5件と一括して再承認する。

### 承認対象

1. L5-M2小分類5件
2. 依存DAGとL5 Design Freeze Gateへの接続
3. worktree所有境界とMerge順
4. L5-M2直接依存の修正案

### 承認結果（2026-08-11）

ユーザーは、L5-M2小分類5件、依存DAG、L5 Design Freeze Gateへの接続、worktree所有境界・Merge順、直接依存へ`L4-M1-S1、L4-M2-S2、L4-M3-S1`を追加する修正案を承認した。

## 決定 025: Knowledge蓄積Workflowの小分類

### 状態

承認済み（2026-08-11）

### 分割・監査方法

分割担当がIssue #1、承認記録、Task Mapを照合し、L5-M3をKnowledge蓄積Workflow開始、Acquisition結果分岐、Update Work Item編成、部分成功後の安全な再開、終端の5件へ分割した。別の妥当性確認担当が、`no-candidate`と入力不足、正常な非永続化と失敗、L5-M1の共通副作用安全性との境界、Candidate並行実行条件、依存DAG、Gate、worktree所有を独立監査した。

監査の結果、5件を維持する。S1は共通Episode受付ではなく受付済みEpisodeからのWorkflow固有開始に限定する。S3は1 Candidateにつき無条件に1 Skillを並行起動する契約にせず、Candidate単位で追跡可能なWork Itemを編成し、Episode内batchまたは直列を既定とする。L3が非競合を明示できる場合だけ並行化を許可する。

### L5-M3-S1

#### タスク名

`Knowledge蓄積Workflowを開始し、Knowledge Acquisitionを起動・成果物を受け入れる`

#### 分割した理由

L5-M1はEpisode入力受付と共通Run初期化を定義するが、受付済みの完了EpisodeをKnowledge Acquisitionへ渡し、成果物・診断を受け取るKnowledge蓄積固有遷移は未定義である。Evidence価値をL5で再判断しない一対一の開始境界として分離する。

主成果物は`docs/design/orchestration/knowledge-accumulation-workflow.md`とする。umbrella文書としてS1が単独所有し、Episode ID、Source Reference、Run／Cycleを伝播する。未完了・欠落・重複Episodeを`no-candidate`へ変換せず、記事評価やAI説明へのExposureだけでは起動しない。

### L5-M3-S2

#### タスク名

`candidate／no-candidate／input_insufficient分岐を定義する`

#### 分割した理由

Knowledge Acquisitionの3結果は、Update経路、正常な候補なし、評価不能という異なる制御へ分かれる。producer成果物を再判定せず、`no-candidate`と入力不足を混同しない機械的分岐として分離する。

主成果物は`docs/design/orchestration/knowledge-accumulation/acquisition-result-branching.md`とする。`candidate`だけを更新経路へ、`no-candidate`をUpdate非起動の正常終端候補へ、`input_insufficient`を診断付き制御判断へ渡す。

### L5-M3-S3

#### タスク名

`Candidate更新Work Itemを編成し、Knowledge Updateを起動する`

#### 分割した理由

Candidateを追跡可能な更新Work Itemへ編成してKnowledge Updateへ渡す処理は、Acquisition結果の分岐や更新結果集約とは異なる一対多の論理制御である。副作用と候補間競合があり得るため、論理fan-outを無条件な物理並行実行と同一視せず独立させる。

主成果物は`docs/design/orchestration/knowledge-accumulation/update-work-items.md`とする。Candidate ID、Source Reference、Run／Cycle／Attemptを維持し、L3-M2の入力不備・4操作・非永続化判断をL5で選び直さない。既定はEpisode内batchまたは直列とし、L3が非競合を明示する場合だけ並行可能とする。

### L5-M3-S4

#### タスク名

`Candidate更新結果を集約し、部分成功後の安全な再開順序を定義する`

#### 分割した理由

成功、正常な非永続化、入力不備、部分成功、失敗、結果不明をCandidate単位で集約し、成功済みを除外したpending集合と再開順序へ写像する処理は、起動時のWork Item編成とは異なる。共通retry／restart／stop可否をKnowledge蓄積Workflowへ適用する単位として分離する。

主成果物は`docs/design/orchestration/knowledge-accumulation/update-result-resumption.md`とする。L5-M1-S5の副作用安全性方針を利用し、結果不明時は照合または停止を選び、盲目的に再実行しない。CLIのexactly-once、冪等性、4更新操作は再設計しない。

### L5-M3-S5

#### タスク名

`正常／異常／Budget到達時のKnowledge蓄積Workflow終端契約を定義する`

#### 分割した理由

`no-candidate`、全Candidate成功、全件正常非永続化、部分完了、入力不足、更新失敗、取消、Budget停止では終了結果が異なる。全経路を有限に閉じ、成功済み成果物、未処理範囲、結果不明、再開可否を失わない終端責務として分離する。

主成果物は`docs/design/orchestration/knowledge-accumulation/workflow-termination.md`とする。全Work Itemを終端または未処理として説明し、未確定の副作用を成功・未実行へ誤分類しない。

### 妥当性監査結果

- 5件はいずれもIssue 8.1、23〜26、39〜40、49、52章が要求する未定義のWorkflow詳細であり、`keep_missing_detail`と判定した。
- S1とS2はWorkflow開始とproducer結果分岐、S2とS5は途中分岐と最終Envelope、S3とS4はWork Item編成と結果集約・再開として分離できる。
- S4はL5-M1-S5の共通方針をKnowledge蓄積へ適用するだけで、冪等性や更新操作を再設計しない。
- 実装はL5-M4、共有評価はL6へ残す。統合・移動・削除、独立技術選定の追加は不要である。

### 依存関係

```text
L5-M1 + L3-M1-S1・S3
              ↓
S1 Workflow開始・Acquisition受入
              ↓
L5-M1-S3・S5 + S2 Acquisition結果分岐
              ↓
L3-M2-S1・S3 + L5-M1-S4 + S3 Update Work Item編成・起動
              ↓
L3-M2-S3・S4 + L5-M1-S4・S5 + S4 結果集約・安全な再開
              ↓
S1〜S4 + L5-M1-S3・S5
              ↓
S5 Workflow終端
              ↓
L5 Design Freeze Gate
```

### L5 Design Freeze Gateへの接続

新しいGateは追加しない。既存Gateで、Exposureから自動起動しないこと、3分岐の非混同、Candidate ID・順序・系譜、正常な非永続化と失敗の分離、成功済み除外、結果不明時の非盲目的再実行、無根拠なCandidate並行実行禁止、有限Budget、終端での未処理範囲保持を確認する。

### worktree所有境界

- S1: `docs/design/orchestration/knowledge-accumulation-workflow.md`
- S2: `docs/design/orchestration/knowledge-accumulation/acquisition-result-branching.md`
- S3: `docs/design/orchestration/knowledge-accumulation/update-work-items.md`
- S4: `docs/design/orchestration/knowledge-accumulation/update-result-resumption.md`
- S5: `docs/design/orchestration/knowledge-accumulation/workflow-termination.md`

共通起点はL5-M1と必要なL3契約を含むCommitとし、Merge順は`S1 → S2 → S3 → S4 → S5 → L5 Gate`とする。各Taskは単独Fileを所有するため別worktreeで先行ドラフト可能だが、確定とMergeは依存元の統合後に行う。L5-M2とは別Directoryで並行可能である。

### 既存構造の修正候補

#### 対象

`L5-M3`の直接依存

#### 現状

`L5-M1、L3-M1-S3、L3-M2-S3・S4`

#### 問題

最終出力契約だけでは、Knowledge AcquisitionのEpisode入力・Source Reference境界と、Knowledge UpdateのCandidate入力受入境界に隠れ依存が残る。

#### 修正案

現行依存へ`L3-M1-S1、L3-M2-S1`を追加する。

#### 理由・影響範囲・再承認単位

いずれもread-onlyの入力契約で循環は生じない。L5-M3台帳・DAG・Gate入力だけを変更し、L3責務と階層数は変更しない。L5-M3小分類5件と一括して再承認する。

### 承認対象

1. L5-M3小分類5件
2. 依存DAGとL5 Design Freeze Gateへの接続
3. worktree所有境界とMerge順
4. L5-M3直接依存の修正案

### 承認結果（2026-08-11）

ユーザーは、L5-M3小分類5件、依存DAG、L5 Design Freeze Gateへの接続、worktree所有境界・Merge順、直接依存へ`L3-M1-S1、L3-M2-S1`を追加する修正案を承認した。

## 決定 026: Parent Orchestration Skill実装・実Component統合の小分類

### 状態

承認済み（2026-08-11）

### 分割・監査方法

分割担当がIssue #1、承認記録、Task Mapを照合し、L5-M4を共通Parent実装、記事Workflow実装、Knowledge蓄積Workflow実装、Mock契約確認、L3実Component統合、L4実Component統合の6件へ分割した。別の妥当性確認担当が、未実装成果物の有無、S1とS2・S3のFile所有、Mock確認の独立完了性、L6との検証境界、L3・L4 Release依存、worktree並行性を独立監査した。

監査の結果、6件を維持し、全件`keep_unimplemented`と判定する。Gateでroot SkillとWorkflow referenceの呼出し境界を固定するため、S1〜S3は同一Gate Commitから別worktreeで並行実装できる。S5とS6も別の実統合Directoryを所有し、それぞれ対応するRelease後に独立して開始できる。

### L5-M4-S1

#### タスク名

`Parent entrypoint・Registry・共通実行制御を実装する`

#### 分割した理由

L5-M1で確定した外部受付、子Skill Registry、共通状態、相関、Budget・副作用安全性を、実行可能なroot Skillと共通referencesへ変換する成果物が未実装である。Workflow固有の処理順とFile所有が異なるためS2・S3から分離する。

主成果物は`.agents/skills/parent-orchestration/SKILL.md`、`references/common/**`、`references/child-registry.md`とする。S1だけがroot `SKILL.md`を編集し、Workflow識別、Registry dispatch、共通Envelope、相関ID、Budget判断、Mock差替え境界を実装する。

### L5-M4-S2

#### タスク名

`記事価値判定WorkflowのParent instructionsを実装する`

#### 分割した理由

L5-M2で確定するClaim fan-out／fan-in、追加調査Routing、最終Assessment受入・終端を、Parentから実行可能なinstructionsへ変換する成果物が未実装である。Knowledge蓄積Workflowと別referenceを所有できるため分離する。

主成果物は`.agents/skills/parent-orchestration/references/workflows/article-reading.md`とする。Article Analysis、Knowledge Search、Reading ValueをRegistry経由で制御し、canonical Claim、Knowledge Assessment、Search Traceを改変せず、追加調査要求と最終Assessmentを排他的に扱う。

### L5-M4-S3

#### タスク名

`Knowledge蓄積WorkflowのParent instructionsを実装する`

#### 分割した理由

L5-M3で確定するAcquisition結果3分岐、Update Work Item、部分成功後の安全な再開、終端を実行instructionsへ変換する成果物が未実装である。記事Workflowから自動接続しない別referenceとして分離する。

主成果物は`.agents/skills/parent-orchestration/references/workflows/knowledge-accumulation.md`とする。Candidate別結果と未処理集合を保持し、結果不明時の盲目的再実行とExposureからの自動起動を禁止する。

### L5-M4-S4

#### タスク名

`Mock ComponentsでParent固有Orchestration契約を確認する`

#### 分割した理由

実Component Release前に、S1〜S3を横断したRegistry dispatch、状態遷移、相関、Budget、再調査、部分成功、終端、Exposure非更新を決定論的に確認する成果物が必要である。実Component統合とは依存時期と失敗原因が異なるため分離する。

主成果物は`tests/contract/parent-orchestration/**`と局所Mockとする。production Skillを編集せず、正常、入力不足、部分失敗、Budget停止、再調査、結果不明を確認する。不適合はS1〜S3のOwnerへ戻す。

### L5-M4-S5

#### タスク名

`Parent Orchestration SkillをL3実Componentsへ接続し、Knowledge蓄積連携を確認する`

#### 分割した理由

Knowledge Acquisition／Updateの実Package接続はL3 Releaseを必要とする一方、L4 Releaseは不要である。Knowledge蓄積側の実成果物引渡しだけを独立して確認できるため分離する。

主成果物は`tests/component/parent-orchestration/knowledge-accumulation/**`とする。RegistryからL3実Skillsを起動し、Episode、Candidate、`no-candidate`、`input_insufficient`、更新結果を契約どおり受け渡す。L3 packageを変更せず、L3内部の意味評価やCLI Black-boxを再試験しない。

### L5-M4-S6

#### タスク名

`Parent Orchestration SkillをL4実Componentsへ接続し、記事価値判定連携を確認する`

#### 分割した理由

Article Analysis／Knowledge Search／Reading Valueの実Package接続はL4 Releaseを必要とするが、L3 Releaseとは独立して確認できる。記事側の実成果物引渡しと再調査・最終Assessment受入だけを確認する単位として分離する。

主成果物は`tests/component/parent-orchestration/article-reading/**`とする。L4 packageとKnowledge CLIを変更せず、各Skillの意味品質を再評価しない。

### 妥当性監査結果

- Issue 7、8.1、9、37〜39、46、51〜52章でParent Skillと2 Workflowの責務は定義済みだが、Skill package、Mock契約確認、実Component接続資産が未作成であるため6件すべて`keep_unimplemented`である。
- S1はroot Skillと共通制御、S2・S3はWorkflow別reference、S4はParent固有の横断Mock契約、S5・S6は対応する実Component引渡しを所有し、完了条件が重複しない。
- S4は共有の意味品質評価ではなく決定論的なParent固有構造・遷移・handoffだけを確認する。共有Fixture、Dataset、Harness、Agent評価、Acceptance、E2E、回帰はL6へ残す。
- 独立技術選定、追加の統合・移動・削除は不要である。

### 依存関係

```text
L5-M1〜M3
    ↓
L5 Design Freeze Gate
    ├─→ S1 共通実装 ─────────┐
    ├─→ S2 記事Workflow実装 ─┼─→ S1〜S3統合 → S4 Mock契約確認
    └─→ S3 蓄積Workflow実装 ─┘
                    ↓
          ┌─────────┴─────────┐
          ↓                   ↓
   L3 Release + S5      L4 Release + S6
          └─────────┬─────────┘
                    ↓
                L5-M4完了
                    ↓
          L6 Acceptance・E2E
```

着手上はS1〜S3をGate通過Commitから並行可能とする。Mergeはroot SkillのOwnerであるS1を先にし、S2・S3を順不同で統合してからS4へ進む。S5とS6はS4後、それぞれ対応するReleaseだけを待って並行可能である。

### L5 Design Freeze Gateとの接続

Gateでは全状態遷移、Registry、ID、Budget、Workflow reference名と呼出し境界、Mock差替え境界、Path Ownerを確定する。Gateに実装成果物を含めず、通過CommitをS1〜S3の共通起点とする。

### worktree所有境界

- S1: `.agents/skills/parent-orchestration/SKILL.md`、`references/common/**`、`references/child-registry.md`
- S2: `.agents/skills/parent-orchestration/references/workflows/article-reading.md`
- S3: `.agents/skills/parent-orchestration/references/workflows/knowledge-accumulation.md`
- S4: `tests/contract/parent-orchestration/**`と局所Mock
- S5: `tests/component/parent-orchestration/knowledge-accumulation/**`
- S6: `tests/component/parent-orchestration/article-reading/**`

L3・L4 Skill packageとL2 CLIはread-onlyで利用する。S4の共有supportを変更する必要があればS4を単一Ownerとして先にMergeする。共有評価資産はL6だけが変更する。Merge順は`S1 → S2・S3 → S4 → S5・S6 → L6`とする。

### L5-M4完了条件

実行可能なroot Skillと2 Workflow references、両WorkflowのMock契約合格、L3・L4それぞれの実Component連携合格、Registry／Version／状態／相関／Budget／終端の設計適合、記事評価からKnowledge更新を自動起動しないこと、未解決の契約差異がないことを必要とする。共有Acceptance・E2Eは後続L6の完了条件である。

### 既存構造の修正候補

#### 対象

決定022およびTask MapのL5-M4実Component統合開始条件

#### 現状

`L3 Release + L4 Release`後に実Component統合を一括開始する。

#### 問題

L3系とL4系は別Workflowであり、別component test Directoryを所有できるため、片方のReleaseを他方が待つ必要がなく過剰な直列化になる。

#### 修正案

`S4 + L3 Release → S5`、`S4 + L4 Release → S6`へ分岐し、`S5 + S6 → L5-M4完了`とする。

#### 理由・影響範囲・再承認単位

各実統合の開始を早めても、L5-M4とL5親の最終完了は両方を待つ。L5内DAG、M4台帳、実統合worktree OwnerとMerge順だけを変更し、L3・L4責務やRelease条件は変更しない。L5-M4小分類6件と一括して再承認する。

### 承認対象

1. L5-M4小分類6件
2. 依存DAG、L5 Gate、L5-M4完了条件
3. worktree所有境界とMerge順
4. 実Component統合開始条件の修正案

### 承認結果（2026-08-11）

ユーザーは、L5-M4小分類6件、依存DAG、L5 Design Freeze Gate、L5-M4完了条件、worktree所有境界・Merge順、L3系とL4系の実Component統合を各Release後に独立開始する修正案を承認した。

これによりL5-M2〜M4の小分類とL5内の依存修正を確定し、次の議論を`L6`「評価設計・横断品質保証を実装する」の中分類から小分類までの連続分割へ移す。

## 決定 027: 評価設計・横断品質保証の中分類

- 状態: 承認済み（決定033）
- 対象: `L6 評価設計・横断品質保証を実装する`
- 日付: 2026-08-11

### 調査方法

分割担当がIssue #1、承認記録、Task Mapを全文確認し、L6を評価規範、共有入力、共通実行基盤、Component品質評価、Workflow品質評価へ分割した。別の妥当性確認担当が、Issue 43〜45、49、51〜52章、各実装固有テストとの重複、設計先行Gate、対象別の部分依存、worktree所有を独立監査した。

### L6-M1

#### タスク名

`横断評価方針・Coverage／判定基準・Search Trace診断／Report契約を詳細設計する`

#### 分割した理由

共有DatasetやHarnessを実装する前に、何をどの基準で合否判定し、評価不能やAgent評価の揺らぎをどう扱い、既存Search Traceから失敗原因をどう診断して報告するかを一つの規範として固定する必要があるため分離する。productionのraw Search Traceと相関契約は再定義せず、L4-M2-S5とL5-M1-S4を固定入力として利用する。

主成果物は要件―Scenario―評価レイヤのCoverage matrix、合否・評価不能基準、決定論的／Agent評価規則、診断写像、Report schemaからなる`docs/design/evaluation/**`の規範仕様とする。

### L6-M2

#### タスク名

`Scenario A〜J・Invariantを検証する共有Dataset／Fixtureを設計・実装する`

#### 分割した理由

IssueはScenarioとInvariantの期待意味を定めているが、再現可能なStore snapshot、Article／Episode／Conversation入力、Evidence由来oracle、期待Assessment／Recommendation／更新履歴、version・provenance・隔離／reset方法は未作成である。評価対象に依存しない共有入力資産として単一Ownerにするため分離する。

### L6-M3

#### タスク名

`CLI・Skill・Workflow共通評価Harnessを設計・実装する`

#### 分割した理由

CLI、Component Skill、Workflowを同一の評価規範とDatasetで実行するには、対象別Adapter、隔離Storeのseed／reset、再現性、Agent試行、Trace収集、結果正規化、Report出力を共通化する必要がある。productionのRegistryやParent制御を再実装せず、評価実行基盤だけを単一Ownerにするため分離する。

Harness技術とrunnerの選定はこの中分類の詳細設計に含め、小分類で設計と実装をGate前後に分離する。

### L6-M4

#### タスク名

`CLI横断品質・L3／L4各Component Skillのレイヤ別／Agent評価を実行し、診断Reportを生成する`

#### 分割した理由

Issue 44章はCLIの検索だけでなく取得・更新、およびKnowledge Search、Acquisition、Update、Reading Valueの評価を要求している。各実装固有テストを繰り返さず、共有Scenarioに対する横断品質・意味品質と失敗原因を対象別に評価する責務として分離する。Parent Workflowの品質はL6-M5へ残す。

### L6-M5

#### タスク名

`記事価値判定・Knowledge蓄積WorkflowのAcceptance／E2E・回帰評価を実行する`

#### 分割した理由

実URLまたはEpisodeから最終Assessmentまたは永続化結果までのユーザー観測可能な振る舞いは、個別CLI・Skill評価やL5のMock／Component引渡し確認だけでは保証できない。2 Workflowの到達状態、Invariant、回帰を実Component経由で確認する最終品質単位として分離する。

### 妥当性監査結果

- 5件はいずれもIssue内で責務の存在は定義済みだが、評価規範・共有資産・実行基盤・評価実行成果物が未定義または未実装であり、タスクとして必要である。
- M1はproduction Traceを変更せず、評価時の診断・Reportだけを所有する。
- M4はCLI検索に限定せず、検索・取得・更新を含む。L3／L4のComponent Skill品質を所有し、Parent Workflow品質はM5が所有する。
- 共有資産契約や最終Reportを第6中分類にせず、M1の規範、M4の診断Report、M5のAcceptance Reportへ配置する。
- 統合・移動・削除は不要である。

### 依存関係

```text
M1 評価規範
  ↓
M2 Dataset設計 + M3 Harness設計
  ↓
L6 Evaluation Design Freeze Gate
  ├─→ M2 Dataset実装
  └─→ M3 Harness実装
          ↓
  ┌───────┼────────┐
  ↓       ↓        ↓
L2実装  L3 Release L4 Release
  ↓       ↓        ↓
M4 CLI  M4 L3     M4 L4 の対象別評価
  │       │        │
  │       └─→ L5-M4-S5 + Knowledge蓄積Acceptance
  │                └─→ L5-M4-S6 + 記事価値Acceptance
  └──────────────────────────────┐
                                 ↓
                  M4全評価 + M5両Workflow回帰
                                 ↓
                       Final Quality Gate
```

M5をM4全件とL5-M4全体の完了後まで一律に待たせず、対応するM4評価sliceとL5-M4-S5またはS6が揃ったWorkflowから開始する。小分類でこの部分完了依存を具体化する。

### L6 Evaluation Design Freeze Gate

M1完了、M2／M3の設計成果物、Issue 44・45・49・52章のCoverage、oracle、Dataset version／provenance、隔離／reset、Agent試行・合否規則、Trace診断写像、Report schema、実Path Ownerを確定する。独立成果物を持たないMilestoneなのでTask IDを付けない。

Final Quality Gateも、M4診断ReportとM5 Acceptance Reportが閾値を満たし、未解決blockerがないことを確認するMilestoneとし、Task IDを付けない。

### worktree所有境界

- M1: `docs/design/evaluation/spec/**`
- M2: `docs/design/evaluation/datasets/**`、`tests/evaluation/datasets/**`
- M3: `docs/design/evaluation/harness/**`、`docs/design/evaluation/adr/**`、`tests/evaluation/harness/**`
- M4: `tests/evaluation/suites/{cli,l3,l4}/**`、`tests/evaluation/reports/{cli,l3,l4}/**`
- M5: `tests/evaluation/workflows/**`、`tests/evaluation/reports/workflows/**`、`tests/evaluation/reports/final/**`

M2とM3はGate通過Commitから別worktreeで並行実装する。M4は対象別、M5はWorkflow別のworktreeへ分離する。M4・M5は共有Fixture、Harness、Report schemaをread-onlyで利用し、変更が必要ならM2・M3のOwnerへ戻す。

### 既存構造の修正候補

Task MapのL6依存にある旧表現`最終Acceptance・E2E: L5`を、`各Workflow Suite着手: 対応するL5-M4-S5またはS6＋Dataset／Harness readiness。Report確定: 対応M4評価。L6完了: L6-M5-S3＋Final Quality Gate`へ具体化する。

2 Workflowを不要に直列化しないための依存修正であり、L1〜L5の階層・責務・テスト所有境界は変更しない。L6中分類5件と全小分類と一括して承認対象にする。

## 決定 028: 横断評価仕様の小分類

- 状態: 承認済み（決定033）
- 対象: `L6-M1 横断評価方針・Coverage／判定基準・Search Trace診断／Report契約を詳細設計する`
- 日付: 2026-08-11

### L6-M1-S1

#### タスク名

`評価レイヤ・対象範囲・L1〜L6テスト所有境界を定義する`

#### 分割した理由

Issue 44章は評価レイヤを定めているが、各実装固有のUnit／Contract／ComponentテストとL6横断評価の境界、対象Artifactの評価可能時点は未定義である。重複実装を防ぐ共通方針として分離する。主成果物は`docs/design/evaluation/spec/evaluation-layer-policy.md`とする。

### L6-M1-S2

#### タスク名

`Scenario A〜J・Invariant・Done DefinitionのCoverage／検証責務Matrixを定義する`

#### 分割した理由

期待結論はIssueで確定しているが、各要件をどのレイヤ、観測種別、Owner、抽象的Oracle種別で検証するかは未定義である。意味を選び直さず検証責務へ写像し、具体入力・期待値はL6-M2-S2だけが所有する。主成果物は`docs/design/evaluation/spec/coverage-responsibility-matrix.md`とする。

### L6-M1-S3

#### タスク名

`評価指標・合否／評価不能・反復判定基準を定義する`

#### 分割した理由

Agent出力の許容差、反復回数、合格閾値、評価不能、重大度、停止条件はIssueで未確定である。決定論的評価とAgent評価の差を機械適用可能な基準へ落とすため分離する。主成果物は`docs/design/evaluation/spec/verdict-criteria.md`とする。

### L6-M1-S4

#### タスク名

`既存Search Trace・相関情報を用いた失敗原因診断手順を定義する`

#### 分割した理由

Issue 43章は診断すべき失敗原因を定めているが、raw Search TraceとParent相関情報から原因を分類する手順は未定義である。L4-M2-S5とL5-M1-S4の契約を変更しないread-only評価手順として分離する。主成果物は`docs/design/evaluation/spec/search-trace-diagnostics.md`とする。

### L6-M1-S5

#### タスク名

`評価結果Envelope・診断／集約Report契約を定義する`

#### 分割した理由

Dataset Version、対象Commit、反復結果、Trace参照、未評価範囲、最終判定をCaseから最終Reportまで同じ系譜で集約する共通表現が未定義であるため分離する。主成果物は`docs/design/evaluation/spec/evaluation-report-contract.md`とする。

### 依存関係・worktree

`S1 → S2 → S3`、`S1 + L4-M2-S5 + L5-M1-S4 → S4`、`S3・S4 → S5`とする。各Taskは`docs/design/evaluation/spec/**`内の対応Fileを単独所有し、Merge順は`S1 → S2・S4 → S3 → S5`とする。S1〜S5はL6 Evaluation Design Freeze Gateの入力である。

## 決定 029: 共有Dataset／Fixtureの小分類

- 状態: 承認済み（決定033）
- 対象: `L6-M2 Scenario A〜J・Invariantを検証する共有Dataset／Fixtureを設計・実装する`
- 日付: 2026-08-11

### L6-M2-S1

#### タスク名

`共有Dataset／Fixture Schema・Version・Provenance契約を設計する`

#### 分割した理由

Case ID、入力、初期Knowledge State、期待値、適用Invariant、Version、出典、互換性の共通形式が未定義である。Scenarioを再現可能な共有資産へ変換する設計Ownerとして分離する。主成果物は`docs/design/evaluation/datasets/dataset-contract.md`とする。

### L6-M2-S2

#### タスク名

`Scenario A〜J・Invariantの具体評価Case Catalog・期待Oracleを設計する`

#### 分割した理由

L6-M1-S2の検証責務を受け、Scenarioの結論を選び直さず、Case ID、具体入力、期待・禁止結果、利用Suite、実装Groupへ落とす設計成果物が必要であるため分離する。主成果物は`docs/design/evaluation/datasets/scenario-catalog.md`とする。

### L6-M2-S3

#### タスク名

`共有Dataset Manifest・Fixture静的検証基盤を実装する`

#### 分割した理由

複数worktreeで作るFixtureを同じSchema、参照整合、Coverage規則で検証する単一Owner実装が必要である。実対象を起動するHarnessとは責務が異なるため分離する。主成果物は`tests/evaluation/datasets/schema/**`と`tests/evaluation/datasets/tools/**`とする。

### L6-M2-S4

#### タスク名

`Scenario A〜E・検索／Knowledge State Invariantの共有Fixtureを実装する`

#### 分割した理由

空Store、完全一致、構成知識、矛盾、古いKnowledgeと検索／Knowledge State領域のInvariantは同じ資産群であり、他のEpisode／記事Fixtureと独立所有できるため分離する。主成果物は`tests/evaluation/datasets/scenarios/a-e/**`とする。

### L6-M2-S5

#### タスク名

`Scenario F〜H・Knowledge Acquisition／Update Invariantの共有Fixtureを実装する`

#### 分割した理由

質問、AI説明、ユーザー訂正とAcquisition／Update領域のInvariantはEpisode、Evidence、更新前後Stateを必要とし、検索Fixtureとは異なるライフサイクル資産であるため分離する。主成果物は`tests/evaluation/datasets/scenarios/f-h/**`とする。

### L6-M2-S6

#### タスク名

`Scenario I〜J・Reading Value Invariantの共有Fixtureを実装する`

#### 分割した理由

記事内容とKnowledge Stateの組合せ、Recognition Gain、Reliability、Attention Cost、推奨範囲とReading Value領域のInvariantを評価する資産は検索・更新Fixtureと異なるため分離する。`残余Invariant`という受け皿を設けず、Fixture化できない責務境界等はL6-M2-S2で適切な検証先へ配分する。主成果物は`tests/evaluation/datasets/scenarios/i-j-reading-value/**`とする。

### 依存関係・worktree

`M1-S1・S2 → S1`、`S1 + M1-S3 → S2`、`L6 Evaluation Design Freeze Gate → S3 → S4・S5・S6`とする。S1・S2は`docs/design/evaluation/datasets/**`、S3はSchema／Validator root、S4〜S6はScenario別Directoryを単独所有する。S4〜S6はS3統合Commitから別worktreeで並行実装する。M2は静的資産と検証を所有し、実行時投入はM3が所有する。

## 決定 030: 共通評価Harnessの小分類

- 状態: 承認済み（決定033）
- 対象: `L6-M3 CLI・Skill・Workflow共通評価Harnessを設計・実装する`
- 日付: 2026-08-11

### 技術選定判断

Harnessの実装言語、Test Runner、Codex Skill起動方式、Process分離、反復実行、依存・Lockfile、Report生成方式はIssueで未定義であり、全Adapterへ変更が波及する。したがって技術選定を独立leafとして設計前に完了させる。

### L6-M3-S1

#### タスク名

`共通評価Harnessの実行・再現性・隔離要件を定義する`

#### 分割した理由

技術選定前に、対象能力、実行環境、State初期化、並行性、時間制御、反復性、失敗隔離を技術非依存で固定する必要があるため分離する。主成果物は`docs/design/evaluation/harness/harness-requirements.md`とする。

### L6-M3-S2

#### タスク名

`評価Harnessの技術スタック・Runner・依存管理方式を選定しADRを確定する`

#### 分割した理由

Harness全体へ影響する未確定技術判断を実装Taskへ分散させると、Adapterごとに前提が割れて手戻りが発生するため分離する。主成果物は`docs/design/evaluation/adr/evaluation-harness-stack.md`とする。

### L6-M3-S3

#### タスク名

`Dataset Loader・Target Adapter・Oracle・Trace Collector・Reporter Architectureを詳細設計する`

#### 分割した理由

共通CoreとCLI／Skill／Workflow Adapterの境界、M2資産の読取、M1診断・Reportへの変換を実装前に固定する必要があるため分離する。主成果物は`docs/design/evaluation/harness/harness-architecture.md`とする。

### L6-M3-S4

#### タスク名

`評価Harness Core・Dataset Loader・Oracle／Report Pipelineを実装する`

#### 分割した理由

全Adapterが共有する反復実行、Case管理、判定、Trace収集、Report出力は単一Ownerが必要であり、Adapter Fakeで独立完了できるため分離する。主成果物は`tests/evaluation/harness/core/**`とHarness固有設定・依存定義とする。

### L6-M3-S5

#### タスク名

`検索・取得・更新CLI Target Adapterを実装する`

#### 分割した理由

Process起動、JSON、終了コード、作業Store隔離を扱うCLI境界はCodex Skill起動と異なるため分離する。主成果物は`tests/evaluation/harness/adapters/cli/**`とする。

### L6-M3-S6

#### タスク名

`L3／L4 Component Skill Target Adapterを実装する`

#### 分割した理由

Skill package起動とMarkdown成果物回収はCLI ProcessやParent Workflowと異なる共通境界であるため分離する。主成果物は`tests/evaluation/harness/adapters/skills/**`とする。

### L6-M3-S7

#### タスク名

`Parent Workflow Target Adapterを実装する`

#### 分割した理由

Workflow entrypoint、Run／Cycle、再調査、Budget、複数成果物の終端回収は個別Skill Adapterと異なるため分離する。主成果物は`tests/evaluation/harness/adapters/workflows/**`とする。

### 依存関係・worktree

`M1-S1・S3 → S1 → S2`、`S2 + M2-S1 + M1-S4・S5 → S3`、`L6 Evaluation Design Freeze Gate → S4 → S5・S6・S7`とする。S1・S3は`docs/design/evaluation/harness/**`、S2は`docs/design/evaluation/adr/**`を所有する。S4はHarness root、共通設定、依存・Lockfileの単独Ownerであり、共通Lockfile変更が不可避な場合もS4だけが変更する。S5〜S7はAdapter別Directoryを所有し、S4統合後に別worktreeで並行実装する。実Componentの合否評価はM4・M5へ残す。

## 決定 031: CLI・L3／L4 Component評価の小分類

- 状態: 承認済み（決定033）
- 対象: `L6-M4 CLI横断品質・L3／L4各Component Skillのレイヤ別／Agent評価を実行し、診断Reportを生成する`
- 日付: 2026-08-11

### L6-M4-S1

#### タスク名

`共有FixtureでCLI検索・取得・更新の横断品質を評価する`

#### 分割した理由

L2固有テストは公開契約と内部実装を確認するが、共有Knowledge State上で検索・取得・更新が横断的に成立するかは未評価であるため分離する。主成果物は`tests/evaluation/suites/cli/**`と`tests/evaluation/reports/cli/**`とする。

### L6-M4-S2

#### タスク名

`Knowledge Acquisition SkillのAgent評価を実行する`

#### 分割した理由

質問、AI説明、ユーザーEvidenceからKnowledge-worthy Candidateを抽出する意味品質はSkill構造確認では判定できず、Updateとは入力・判断・失敗原因が異なるため分離する。主成果物は`tests/evaluation/suites/l3/knowledge-acquisition/**`と`tests/evaluation/reports/l3/knowledge-acquisition/**`とする。

### L6-M4-S3

#### タスク名

`Knowledge Update SkillのAgent評価を実行する`

#### 分割した理由

重複抑止、Evidence追加、訂正、Supersede、非永続化の選択はCLI更新成功とは別の意味判断であり、Acquisitionとも完了条件が異なるため分離する。主成果物は`tests/evaluation/suites/l3/knowledge-update/**`と`tests/evaluation/reports/l3/knowledge-update/**`とする。

### L6-M4-S4

#### タスク名

`Article Analysis SkillのAgent評価を実行する`

#### 分割した理由

記事をClaim、Role、Importance、Location、Supportへ分解する意味品質はMarkdown構造確認とは別であり、他のL4 Skillsと失敗原因が異なるため分離する。主成果物は`tests/evaluation/suites/l4/article-analysis/**`と`tests/evaluation/reports/l4/article-analysis/**`とする。

### L6-M4-S5

#### タスク名

`Knowledge Search SkillのAgentic Search評価を実行する`

#### 分割した理由

直接・関連知識への到達、7 Assessment状態、探索停止、Search Trace診断はCLI検索結果だけでは評価できず、Article Analysis／Reading Valueとも判断責務が異なるため分離する。主成果物は`tests/evaluation/suites/l4/knowledge-search/**`と`tests/evaluation/reports/l4/knowledge-search/**`とする。

### L6-M4-S6

#### タスク名

`Reading Value SkillのAgent評価を実行する`

#### 分割した理由

Knowledge AssessmentからRecognition Gain、Reliability、Attention Cost、読書範囲を統合する意味品質は独立して評価できるため分離する。主成果物は`tests/evaluation/suites/l4/reading-value/**`と`tests/evaluation/reports/l4/reading-value/**`とする。

### 依存関係・worktree

全leafがM1-S2・S3・S5、M2-S3、M3-S4をread-onlyで利用する。その上で、S1は`M2-S4・S5 + M3-S5 + L2-M3-S9`、S2は`M2-S5 + M3-S6 + L3-M1-S4`、S3は`M2-S5 + M3-S5・S6 + L3-M2-S6`、S4は`M2-S6 + M3-S6 + L4-M1-S5`、S5は`M2-S4 + M3-S5・S6 + L4-M2-S7`、S6は`M2-S4・S6 + M3-S6 + L4-M3-S6`へ依存する。一律のL2〜L4 Release待ちは設けず、各対象の実装完了をTarget Readiness Gateとする。

各leafは対象別SuiteとReport Directoryを単独所有し、共有Harness・Fixtureをread-onlyで利用する。S2・S4・S6はS1を待たず並行可能である。S3・S5もSuite準備・試行開始はS1と並行できるが、Report確定は基盤となるS1のCLI横断品質合格を必要とする。M4全体完了は6件すべてを必要とするが、M5は必要な部分結果だけで進められる。

## 決定 032: Workflow Acceptance／E2E・回帰評価の小分類

- 状態: 承認済み（決定033）
- 対象: `L6-M5 記事価値判定・Knowledge蓄積WorkflowのAcceptance／E2E・回帰評価を実行する`
- 日付: 2026-08-11

### L6-M5-S1

#### タスク名

`Knowledge蓄積WorkflowのAcceptance／E2E評価を実行する`

#### 分割した理由

EpisodeからAcquisition、Update、CLI、Storeまでの連携はL3各SkillとCLIの個別評価だけでは保証できず、記事Workflowを待たず独立評価できるため分離する。主成果物は`tests/evaluation/workflows/knowledge-accumulation/**`と`tests/evaluation/reports/workflows/knowledge-accumulation/**`とする。

### L6-M5-S2

#### タスク名

`記事価値判定WorkflowのAcceptance／E2E評価を実行する`

#### 分割した理由

URLからArticle Analysis、Claim別Search、再調査、Reading Value、終了までの全体挙動は個別評価だけでは保証できず、Knowledge蓄積Workflowを待たず独立評価できるため分離する。主成果物は`tests/evaluation/workflows/article-reading/**`と`tests/evaluation/reports/workflows/article-reading/**`とする。

### L6-M5-S3

#### タスク名

`全Scenario回帰Suiteを実行し、最終Acceptance／診断Reportを確定する`

#### 分割した理由

両Workflowの個別合格後に、全Scenario・Invariant・Done DefinitionのCoverage、反復安定性、未解決失敗を一つのシステム判定へ集約する実行成果物が必要である。単なるReviewではなく全回帰実行を伴うためFinal Quality Gateから分離する。主成果物は`tests/evaluation/workflows/regression/**`と`tests/evaluation/reports/final/**`とする。

### 依存関係・worktree

S1のSuite準備・試行着手は`L6-M2-S5 + L6-M3-S5・S7 + L5-M4-S5`、Report確定は`L6-M4-S1〜S3`合格へ依存する。S2のSuite準備・試行着手は`L6-M2-S4・S6 + L6-M3-S5・S7 + L5-M4-S6`、Report確定は`L6-M4-S1・S4〜S6`合格へ依存する。S1とS2は対応するWorkflowごとに別worktreeで並行できる。`S1・S2 → S3 → Final Quality Gate`とし、S1・S2はWorkflow別Directory、S3は回帰Suiteと最終Report Directoryを単独所有する。

## L6小分類案の既存構造修正候補

1. L6-M4をL2・L3・L4全完了待ちにせず、各Target leafを対応Artifactの完成後に開始する。
2. L6-M5をM4・L5全完了待ちにせず、各Workflow Suiteは対応するDataset／Harness AdapterとL5-M4-S5またはS6で着手し、Report確定だけを対応するM4評価slice合格へ依存させる。

階層変更ではなく、部分完了を利用して過剰直列化を除く依存補正である。L6中分類5件、小分類27件、L6 Gate、Final Quality Gate、worktree所有境界と一括して承認対象にする。

## L6全小分類の妥当性監査結果

- 最終件数は27件で、追加、統合、移動、削除はない。
- 設計leaf 10件は`keep_missing_detail`、実装・評価leaf 17件は`keep_unimplemented`である。
- L6-M1-S2は検証責務の抽象Matrix、L6-M2-S2は具体Caseと期待Oracleを所有し、重複しない。
- L6-M1-S4は既存Traceと相関をread-onlyで利用し、production契約を再定義しない。
- InvariantはL6-M2-S4〜S6の各領域と、静的検証・M4・M5へ明示配分し、`残余Invariant`という受け皿を作らない。
- L6-M3-S2の独立技術選定は、全Adapterへ波及するHarnessの技術前提を実装前に固定するため必要である。
- M4の各leafはM3-S5またはS6の対象Adapterを直接依存に持つ。M5の両Workflowは到達に必要なM3-S7を直接依存に持つ。
- M4-S3とS5はSuite準備・試行をM4-S1と並行可能とするが、Report確定はCLI横断品質合格を必要とする。
- 各実装固有Unit／Contract／Component testはL2〜L5、共有意味品質・横断品質・Acceptance／E2E・回帰はL6という既存境界を維持する。
- L1〜L5の階層・責務に新たな修正はない。修正候補はTask MapのL6依存表現1件だけである。

### L6 Evaluation Design Freeze Gateの確定条件

`M1-S1〜S5 + M2-S1・S2 + M3-S1〜S3`を入力に、Issue 44・45・49・52章の全Requirement IDへOwner、観測、Oracle種別が割り当てられ、具体Caseに欠落がないこと、Dataset schema／version／provenance、Harness要件／ADR／Architecture、Agent反復／評価不能、既存Trace診断、Report schema、実Path Owner、共通起点Commit、Merge順が確定していることを確認する。実装や評価合格は含めず、独立Task IDを付けない。

### Final Quality Gateの確定条件

L6-M5-S3の全回帰・最終Reportが閾値を満たし、未解決blockerがないことだけを確認する。実行成果物はS3が所有するため、Gateへ独立Task IDを付けない。

### Path Ownerの確定案

- L6-M1: `docs/design/evaluation/spec/**`
- L6-M2設計: `docs/design/evaluation/datasets/**`
- L6-M2実装: `tests/evaluation/datasets/**`
- L6-M3要件・Architecture: `docs/design/evaluation/harness/**`
- L6-M3 ADR: `docs/design/evaluation/adr/evaluation-harness-stack.md`
- L6-M3実装: `tests/evaluation/harness/**`
- L6-M4 Suite／固定baseline Report: `tests/evaluation/suites/{cli,l3,l4}/**`、`tests/evaluation/reports/{cli,l3,l4}/**`
- L6-M5 Suite／固定baseline Report: `tests/evaluation/workflows/**`、`tests/evaluation/reports/workflows/**`、S3のみ`tests/evaluation/reports/final/**`

CIの一時Reportはcommit対象にせずartifact出力とする。共有Fixture、Harness、Report schemaの変更は対応Ownerへ戻し、M4・M5から変更しない。

### 承認対象

1. L6中分類5件
2. L6小分類27件
3. L6 Evaluation Design Freeze GateとFinal Quality Gate
4. 対象別Target Readiness、Workflow別着手／Report確定依存
5. worktree所有境界とMerge順
6. Task MapのL6依存表現を具体化する修正案

## 決定 033: L6中分類・小分類・依存関係の承認

- 状態: 承認済み
- 対象: `L6 評価設計・横断品質保証を実装する`
- 日付: 2026-08-11

### 承認結果

ユーザーは、決定027〜032で提示した以下をすべて承認した。

1. L6中分類5件
2. L6小分類27件
3. L6 Evaluation Design Freeze GateとFinal Quality Gate
4. 対象別Target Readiness、Workflow別の着手依存とReport確定依存
5. worktreeのPath Owner、共通起点Commit、Merge順
6. L6-M4・L6-M5を親単位で一括待機させず、対象別の部分完了を利用する依存修正

承認後の依頼に基づき、次の議論はIssue #1を原典とするL1〜L6全タスクの最終重複監査とする。過去の承認を維持理由にせず、同一成果物、同一完了条件、意味判断、検証責務、GateとTaskの重複を再確認する。

## 決定 034: Issue #1からの全タスク最終重複監査

- 状態: 調査完了
- 対象: 承認済みのL1〜L6全タスク、Gate、依存DAG、worktree所有境界
- 日付: 2026-08-11

### 調査方法

Issue #1を再取得して全文を原典とし、Task Mapと全承認記録を分割担当と妥当性確認担当が独立して照合した。名称の類似ではなく、主成果物、完了条件、意味判断、副作用、評価Oracle、Path Ownerが同じかを基準に、`keep_existing / merge / move / remove_defined / revise_existing`を判定した。

### 結論

- `merge`: 0件
- `move`: 0件
- `remove_defined`: 0件
- 漏れ: 0件
- 依存循環: 0件
- Path Owner競合: 0件
- 全重点疑義: `keep_existing`
- 追加の`revise_existing`: 0件

Issueは原則、責務、期待結論を定義しているが、各leafはSchema、適用手順、技術選定、具体Case、Fixture、Harness、合否基準、実装、Reportという未定義または未実装の成果物を持つ。したがって、Issue内容を再議論するだけのタスクは存在しない。

### 重点境界

1. L2固有テストとL2-M3-S9はUnit／IntegrationおよびL1公開CLI契約適合、L6-M4-S1は共有Fixture上の検索・取得・更新の横断品質と診断Reportを所有する。
2. L3-M2はCandidateの同一性、更新操作、副作用を所有し、L4-M2はArticle Claimのread-only探索、7状態Assessment、raw Search Traceを所有する。
3. L4は意味成果物を生成し、L5は起動、fan-out／fan-in、相関、Routing、Budget、終了を制御する。L5はL4の意味判断を再判定しない。
4. L5-M4はMock Contractと実Component引渡しを確認し、L6-M4／M5は共有Scenarioによる意味品質、Acceptance、E2E、回帰を評価する。
5. L6-M1-S2は抽象的な検証責務・Oracle種別、L6-M2-S2は具体Case・入力・期待値・禁止結果を所有する。
6. L6-M3-S5〜S7は再利用可能な対象起動・回収Adapter、L6-M4／M5は具体Scenarioの実行・判定・Reportを所有する。
7. 各Gateは成果物を作らないMilestoneであり、直前Taskの成果物を確認するだけなのでTaskと重複しない。
8. 各親の成果物欄は子成果物のroll-upであり、親が同じ成果物を別途作成する計画ではない。

### Issue章との対応

- 10〜14、25、48、51章: L1契約、L2物理実装・固有テスト
- 23〜26、39、47章: L3 Knowledge蓄積判断、L5 Workflow制御、L6評価
- 17〜22、27〜38、47章: L4記事・検索・Reading Value判断、L5記事Workflow、L6評価
- 7、8.1、16、24、26、37〜39章: L5 Parent Orchestration
- 43〜45、49、51〜52章: L6診断、共有評価資産、Acceptance、回帰

### 記録同期

決定033の承認結果に合わせて、決定027〜032、Task MapのL6中分類5件・小分類27件、L6依存表現、現在位置を承認済みへ同期した。タスク名、階層、成果物、依存DAG、Gate、worktree所有境界の修正はない。

## 決定 035: 全タスクの接続記載完全性監査

- 状態: 議論中
- 対象: 承認済みのL1〜L6全タスク、依存DAG、Gate／Release、worktree実行規則
- 日付: 2026-08-11

### 調査方法

Issue #1を再取得して原典とし、Task Mapと承認記録を分割担当・妥当性確認担当が独立して照合した。各leafについて、入力、着手依存、完了・Merge依存、後続利用者、Gate／Release、worktreeの共通起点・Path Owner・Merge順を確認した。親行と詳細DAGの表現差、非ID参照、部分完了依存の親単位への丸め込みも監査した。

### 結論

- タスク追加・統合・移動・削除: 0件
- 循環・到達不能・無効ID: 0件
- タスク構造の変更: 不要
- 依存記録の修正候補: 8件
- 現行記載で十分: L6親のWorkflow別着手・Report確定・完了条件

依存DAGの骨格は成立している。一方、leaf台帳では正しく記載されていても親行・全体図で省略されている接続、着手と完了・Mergeが図上で混同される接続、抽象名のため変更影響をIDから逆引きできない接続、後続利用者と設計worktree所有が明記されていない範囲が残っている。

### 既存構造の修正候補

#### 1. L2-M2親行と全体DAGへKnowledge Update側の探索要求を追加する

- 対象: `L2-M2`親行、全体DAG
- 現状: `L3-M2-S2 → L2-M2-S1`はleaf台帳とL2詳細DAGにあるが、L2-M2親行と全体DAGには`L4-M2-S1`だけが記載されている。
- 修正案: L2-M2の開始入力を`L1、L3-M2-S2、L4-M2-S1`とし、全体DAGにも`L3-M2-S2 → L2-M2-S1`を追加する。L3-M2全体や実CLI統合を開始条件にはしない。
- 理由: Knowledge Update固有の既存Knowledge特定要件も検索・Index設計の正式入力であり、親行から欠落すると更新経路の要求が追跡不能になるため。
- 影響範囲: L2-M2親行、全体DAGのみ。leafの責務と既存DAGは変更しない。

#### 2. L2親・L2-M1・L2-M2で設計着手と本番実装着手を分離する

- 対象: `L2`、`L2-M1`、`L2-M2`の親行
- 現状: 詳細節ではL2 Design Freeze Gate後に本番実装すると確定しているが、親行だけではGate前の詳細設計とGate後の本番実装を区別できない。
- 修正案: `詳細設計TaskはGate前に実施可能。本番実装Taskの着手はL2 Design Freeze Gate通過後`を親行へ明記する。親全体をGate後まで着手不可にはしない。
- 理由: 設計確定前の本番実装を防ぎつつ、技術選定・詳細設計の先行を妨げないため。
- 影響範囲: L2系親行の依存表現のみ。

#### 3. L2-M2-S5の着手依存と完了・Merge依存を図上でも分離する

- 対象: L2横断DAG、L2-M2内DAG
- 現状: leaf台帳は`着手: L2-M2-S3＋L2 Design Freeze Gate、完了・Merge: L2-M1-S10`と正しいが、図はS10完了後にしかS5へ着手できないように読める。
- 修正案: S5への辺を`着手`と`完了・Merge`の2種類で表示する。
- 理由: Port／Mockによる先行実装を維持し、正本照合の完了だけをS10へ依存させるため。
- 影響範囲: 図の表現のみ。

#### 4. L6-M3-S5の契約着手と実CLI Artifact接続完了を分離する

- 対象: `L6-M3-S5`
- 現状: `L6-M3-S4、L1-M2`だけが記載され、公開契約による先行実装と実Artifactへの接続完了が区別されていない。
- 修正案: `着手: L6-M3-S4、L1-M2-S6。完了・Merge: L2-M3-S8`とする。品質評価の確定は従来どおりL2-M3-S9後とする。
- 理由: Adapterを契約から先行実装しつつ、実Processへ接続していない状態を完成扱いしないため。
- 影響範囲: L6-M3-S5、L6 Gate後DAG、Target Readiness表現。

#### 5. L6-M3-S6のproducer契約を具体ID化する

- 対象: `L6-M3-S6`
- 現状: `L3／L4 producer契約`という抽象名であり、変更影響をIDから逆引きできない。
- 修正案: `L3-M1-S3、L3-M2-S3・S4、L4-M1-S4、L4-M2-S5、L4-M3-S5`を着手入力として明記する。実Skill ReleaseはM4の対象別Readiness条件に残す。
- 理由: Adapterの入力・出力Envelopeへ影響する契約Ownerを一意にするため。
- 影響範囲: L6-M3-S6とL6 Gate後DAGの表現のみ。

#### 6. L6-M3-S7の設計着手とParent Mock接続完了を分離する

- 対象: `L6-M3-S7`
- 現状: `L5-M1〜M3設計`という親単位参照だけで、設計確定と実行可能なParent入口の存在を区別できない。
- 修正案: `着手: L6-M3-S4、L5 Design Freeze Gate。完了・Merge: L5-M4-S4`とする。Gateの入力であるL5-M1〜M3各leafを重複して直接依存にせず、L3／L4実Component ReleaseはM5のWorkflow別Readiness条件に残す。
- 理由: Parent契約からAdapterを先行実装し、Mock Parentで入口・Run／Cycle・終端回収を確認してから完成させるため。
- 影響範囲: L6-M3-S7、L6 Gate後DAG、Workflow Readiness表現。

#### 7. 全leafの後続利用者と依存種別を逆引き可能にする

- 対象: Task Mapの全leaf
- 現状: 台帳は前提側の`直接依存`を中心とし、成果物の後続利用者、依存種別、Gate／Releaseへの接続を全leafから一意に逆引きできない。
- 修正案: 既存の直接依存欄を依存辺の正本として維持し、そこから逆引きする接続台帳を追加する。各辺について`Producer leaf、直接の後続利用者、利用成果物、接続種別（着手／確定／完了・Merge／Report確定／read-only）、入力となるGate`を記録し、利用先がない終端Taskは`なし`とする。推移的利用者は列挙せず、GateをTask化しない。
- 理由: consumer、変更影響、Merge順を前後両方向から監査可能にし、親名・抽象名・文章だけに隠れた依存を防ぐため。
- 影響範囲: Task Mapの表示・更新規則。タスク階層と責務は変更しない。

#### 8. L1とL2 Gate前設計Taskのworktree接続を明記する

- 対象: L1全leaf、L2-M1-S1〜S3、L2-M2-S1〜S4、L2-M3-S1・S2
- 現状: 論理DAGはあるが、設計文書の単一Path Owner、worktree共通起点Commit、Merge順が明記されていない。L2はGate後実装の所有規則だけが詳細である。
- 修正案: 設計成果物ごとのPath Owner、設計開始Commit、承認済みDAGに従うMerge順を追記する。L2実装Pathの実値をGate通過記録で確定する現行方針は維持する。
- 理由: 並行設計時の文書競合と、未Mergeの設計をGate入力へ混入させることを防ぐため。
- 影響範囲: L1・L2のworktree記録のみ。

### 維持する接続

L6親行のWorkflow別接続は、Knowledge蓄積と記事価値判定について`対応するL5-M4-S5またはS6＋Dataset／Harness readiness`、`対応M4評価後にReport確定`、`M5-S3＋Final Quality Gateで完了`を既に区別している。具体leafはL6-M5-S1・S2に列挙済みであり、親行へ同じ全IDを重複記載する必要はない。現行を維持する。

### 再承認単位（初回案。下記の再監査で訂正）

上記8件を「タスク構造を変えない依存記録の具体化」として一括再承認対象にする。承認後にTask Mapへ反映し、接続台帳と既存DAGの整合を再検証する。

### 記録完全性の再監査・追補

- 状態: 議論中
- 日付: 2026-08-11
- 結論: 上記8件の方向性は妥当だが、再承認可能な記録としては具体性が不足している。Task Mapの確定構造はまだ変更しない。

Issue #1、Task Map、決定035を分割担当と妥当性確認担当が独立して再照合した。その結果、各候補の`問題`、変更しない依存、直接の後続利用者、依存時点、利用方法、Gate／worktreeとの関係、および候補7・8で作る表の実体が不足していることを確認した。

#### 原典根拠

| 候補 | 原典・規則上の根拠 |
| --- | --- |
| 1 | Issue 6.7、24、39〜40章。検索は取得要求から逆算し、Knowledge Updateも既存Knowledge探索を必要とする |
| 2・3 | Issue 6.7、48、51章。物理方式は詳細設計で確定し、検索要件から逆算して設計する |
| 4〜6 | Issue 9、40、44、47、51章。JSON／Markdown境界、Codex／CLI責務、評価対象、次フェーズ成果物を接続する |
| 7・8 | Issue 51〜52章、および`decompose-tasks`の依存DAG・worktree所有規則 |

#### 候補別の不足と補正

##### 1. L2-M2親行と全体DAG

- 問題: 更新用探索要求の直接辺が親表示から欠落している。
- 着手依存の補正: `L3-M2-S2 → L2-M2-S1`を追加する。L3-M2全体は待たない。
- 完了・Merge依存: 変更なし。
- 直接の後続利用者: `L2-M2-S1`。
- Gate／worktree: L2 Design Freeze Gateの入力へ推移するが、worktree規則は変更しない。
- 再承認単位: 候補1〜6の依存表現。

##### 2. L2系親の設計着手・本番実装・完了条件

- 問題: 親行が子Taskの開始時点を一括表現し、Gate前設計とGate後実装を区別できない。
- `L2`: `設計着手: L1。本番実装: L2 Design Freeze Gate。完了: L2-M1〜M3`。
- `L2-M1`: `設計着手: L1およびL2-M2-S1の必要入力。本番実装: L2 Design Freeze Gate。完了: 全子Task（最終到達点L2-M1-S8〜S10）`。
- `L2-M2`: `設計着手: L1、L3-M2-S2、L4-M2-S1。本番実装: L2 Design Freeze Gate。完了・Merge: 全子Task、L2-M1`。
- 直接の後続利用者: L2-M1はL2-M2・L2-M3、L2-M2はL2-M3・L3／L4実CLI統合、L2はL3〜L6の実Artifact利用Task。
- worktree: Gate前設計の所有は候補8、Gate後実装の所有は既存L2 Gate記録を正とする。
- 再承認単位: 候補1〜6の依存表現。

##### 3. L2-M2-S5の着手辺と完了辺

- 問題: L2横断DAGでは`L2-M1-S10`完了後にしかS5へ着手できないように読める。
- 着手依存: `L2-M2-S3、L2 Design Freeze Gate`。
- 完了・Merge依存: `L2-M1-S10`。
- 直接の後続利用者: `L2-M2-S6〜S9、L2-M3-S5`。
- 図の補正: L2横断DAGは着手辺と完了辺へ分離する。L2-M2内DAGは既存の`Gate → S5着手`を維持し、`L2-M1-S10 → S5完了・Merge`だけを明記する。
- worktree: 既存L2 Gate記録を参照し、S5が共通検索Coreの単一Ownerとなる。
- 再承認単位: 候補1〜6の依存表現。

##### 4. L6-M3-S5 CLI Target Adapter

- 問題: 公開契約を用いた先行実装と、実CLI Processへの接続完了が混同されている。
- 着手依存: `L6-M3-S4、L1-M2-S6`。
- 完了・Merge依存: `L2-M3-S8`。
- 直接の後続利用者: `L6-M4-S1・S3・S5、L6-M5-S1・S2`。
- Report確定: CLI品質を利用する後続Reportは`L2-M3-S9`合格を対象別Readiness条件とする。
- worktree: `tests/evaluation/harness/**`内のCLI Adapter固有PathをS5が単一所有し、M3-S4統合後に分岐してMergeする。実PathはL6 Gate記録で確定する。
- 再承認単位: L6-M3-S5と上記直接後続の接続を、候補1〜6の依存表現として承認する。

##### 5. L6-M3-S6 Component Skill Target Adapter

- 問題: `L3／L4 producer契約`が抽象的であり、旧修正案は出力契約だけを列挙し、Skill起動に必要な入力契約を欠いていた。また`L3-M2-S4`はCLI Mappingでありproducer契約ではない。
- 着手依存: 共通で`L6-M3-S4`。Acquisition=`L3-M1-S1・S3`、Update=`L3-M2-S1・S3`、Article Analysis=`L4-M1-S1・S4`、Knowledge Search=`L4-M2-S2・S5`、Reading Value=`L4-M3-S1・S5`。
- 除外: `L3-M2-S4`と`L4-M2-S3`のCLI連携契約はL6-M3-S5側で扱う。
- 完了・Merge依存: 実Skill Releaseは待たず、契約とMock TargetでAdapterを完成可能とする。実Target readinessはL6-M4側で判定する。
- 直接の後続利用者: `L6-M4-S2〜S6`。
- Gate／worktree: L6 Gate通過後、`tests/evaluation/harness/**`内のComponent Skill Adapter固有PathをS6が単一所有し、M3-S4統合後にMergeする。実PathはL6 Gate記録で確定する。
- 再承認単位: 候補1〜6の依存表現。旧修正案の具体ID列挙は本項で置換する。

##### 6. L6-M3-S7 Parent Workflow Target Adapter

- 問題: Parent設計の確定と、実行可能なMock Parent入口への接続完了が混同されている。
- 着手依存: `L6-M3-S4、L5 Design Freeze Gate`。
- 完了・Merge依存: `L5-M4-S4`。
- 直接の後続利用者: `L6-M5-S1・S2`。
- Gate／worktree: L6 Gate通過後、`tests/evaluation/harness/**`内のParent Adapter固有PathをS7が単一所有し、M3-S4統合後にMergeする。L3／L4実Component ReleaseはM5側のReadinessに残す。実PathはL6 Gate記録で確定する。
- 再承認単位: L6-M3-S7、L5-M4-S4、L6-M5-S1・S2の接続を、候補1〜6の依存表現として承認する。

##### 7. 全leaf接続逆引き表

- 問題: 提案した列で`read-only`を依存時点と同列に扱っており、依存の時期と利用方法を区別できない。また、全leaf分の実レコードをまだ提示していない。
- 正本: 既存leaf台帳の直接依存を依存辺の正本とし、逆引き表はその派生表示とする。
- 必須列: `Producer leaf`、`直接の後続利用者`、`利用成果物`、`依存時点（着手／確定／完了・Merge／Report確定／Gate入力／Release条件）`、`利用方法（read-only／変更あり）`、`Gate名とGate入力成果物`、`worktree所有記録への参照`。
- 除外: 推移的利用者と親roll-upは列挙しない。終端Taskは`なし`、GateはMilestoneとして記録しTask IDを付けない。
- 未充足: 全leaf分の実レコードは未作成である。実表を提示するまでは候補7を承認可能と扱わない。
- 再承認単位: 全leafの接続逆引き表を独立した承認単位とする。

##### 8. L1・L2 Gate前設計Taskのworktree所有表

- 問題: 対象Taskは列挙済みだが、Path／Glob、単一Owner、共通起点Commit、Merge前提・順の実値がない。
- 必須対象: `L1-M1-S1〜S5`、`L1-M2-S1〜S6`、`L2-M1-S1〜S3`、`L2-M2-S1〜S4`、`L2-M3-S1・S2`。
- 必須列: `Task ID`、`設計成果物Path／Glob`、`単一Owner`、`共通起点Commit`、`Merge前提ID`、`Merge順`、`共有ADR索引／Schema参照／生成物の更新Owner`。
- 起点Commit: 実作業開始前に、Issue #1と承認済みTask Map／決定記録を含むplanning baseline Commitを作り、そのSHAを実行記録へ固定する。未Commit変更をworktree起点へ含めない。
- 共有物: Task Map・決定記録は計画Ownerだけが統合後に更新する。各設計worktreeから同時編集しない。共有ADR索引と生成物は表で指定する単一Ownerだけが変更する。
- 境界: 今回確定するのはGate前設計文書のPathである。Schema、Migration、Repository、Index、CLI Runtime等の本番実装Pathは、既存方針どおりL2 Gate通過記録で実値を確定する。
- 未充足: 対象20 leafの所有表は未作成である。実表を提示するまでは候補8を承認可能と扱わない。
- 再承認単位: L1・L2 Gate前設計leafのworktree所有表を独立した承認単位とする。

#### 再承認単位の訂正

上記8件を現時点で一括承認対象にはしない。次の3単位に分ける。

1. 候補1〜6の具体的な依存表現
2. 全leafの接続逆引き表
3. L1・L2 Gate前設計leafのworktree所有表

2と3は実表を提示してから承認対象に含める。タスク数、階層、責務、Gate数は変更しない。

---

## 決定 036: Gate前設計成果物の配置・単一Owner・接続台帳の提案

### 状態

承認済み（2026-08-11）

### 背景

GitHub Issue化とworktree分割の前に、各leafの成果物をどこへ出力し、誰が共有Pathを更新するかを確定する必要がある。分割担当と妥当性確認担当がIssue #1、Task Map、決定035を独立して再照合した結果、新しいleafの追加は不要であり、L1の11 leafとL2 Design Freeze Gate前の9 leafに具体Pathと単一Ownerを割り当てればよいと判断した。

今回確定候補とするのは設計文書のPathである。技術選定に依存するSchema、Migration、Repository、Index、CLI Runtime等の本番実装Pathは、L2 Design Freeze Gate通過記録で確定する。

### L1設計成果物の配置案

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
| L1-M2-S4 | `docs/design/cli-contract/schemas/**`（Catalogは`schemas/schema-catalog.md`） | L1-M2-S4 | L1-M2-S1〜S3 |
| L1-M2-S5 | `docs/design/cli-contract/versioning-and-compatibility.md` | L1-M2-S5 | L1-M2-S4 |
| L1-M2-S6 | `docs/design/cli-contract/requirements-traceability.md`、`docs/design/cli-contract/README.md` | L1-M2-S6 | L1-M2-S1〜S5 |

#### 配置理由

- `knowledge-model`と`cli-contract`を分け、論理データ契約と公開CLI契約のPath所有を交差させない。
- JSON SchemaとCatalogを`schemas/**`の単一Ownerへ集約し、並行worktree間の共有Path競合を防ぐ。
- Traceability確認後に次の契約を確定するため、`L1-M1-S5 → L1-M2-S1`を直接依存として具体化する。
- 各`README.md`は成果物一覧・依存・参照リンクだけを持ち、設計本文を複製しない。

#### L1 Merge DAG

```text
L1-M1: S1 → {S2, S3} → S4 → S5
L1-M2: L1-M1-S5 → S1 → {S2, S3} → S4 → S5 → S6
```

### L2 Design Freeze Gate前の設計成果物配置案

| Task ID | 設計成果物Path／Glob | 単一Owner | Merge前提 |
| --- | --- | --- | --- |
| L2-M2-S1 | `docs/design/search-infrastructure/search-and-index-requirements.md` | L2-M2-S1 | L1、L3-M2-S2、L4-M2-S1 |
| L2-M1-S1 | `docs/design/knowledge-store/adr/0001-persistence-stack.md` | L2-M1-S1 | L2-M2-S1 |
| L2-M1-S2 | `docs/design/knowledge-store/physical-model-history-schema-evolution.md` | L2-M1-S2 | L2-M1-S1 |
| L2-M1-S3 | `docs/design/knowledge-store/repository-transaction-boundary.md`、`docs/design/knowledge-store/README.md` | L2-M1-S3 | L2-M1-S2 |
| L2-M2-S2 | `docs/design/search-infrastructure/adr/0001-search-index-stack.md` | L2-M2-S2 | L2-M2-S1、L2-M1-S1〜S3 |
| L2-M2-S3 | `docs/design/search-infrastructure/candidate-search-architecture.md` | L2-M2-S3 | L2-M2-S2 |
| L2-M2-S4 | `docs/design/search-infrastructure/index-lifecycle-recovery.md`、`docs/design/search-infrastructure/README.md` | L2-M2-S4 | L2-M2-S3、L2-M1-S3 |
| L2-M3-S1 | `docs/design/cli-runtime/adr/0001-runtime-validation-distribution.md` | L2-M3-S1 | 調査着手=L1-M2-S6、ADR確定・Merge=L2-M1-S1＋L2-M2-S2 |
| L2-M3-S2 | `docs/design/cli-runtime/command-to-port-architecture.md`、`docs/design/cli-runtime/README.md` | L2-M3-S2 | L2-M3-S1、L2-M1-S3、L2-M2-S3・S4 |

#### 配置理由

- `knowledge-store`、`search-infrastructure`、`cli-runtime`を分け、物理正本、検索基盤、CLI内部実行境界のOwnerを明確にする。
- 検索要求を先に固定し、その要求に適合する永続化技術、検索技術、CLI Runtimeを順に確定することで実装後の設計変更を避ける。
- ADR番号はDirectory-localとし、全体連番やglobal ADR indexを作らない。これにより別技術領域の並行編集を避ける。
- `L2-M3-S1`は公開CLI契約後に調査を開始できるが、ADR確定はStore/Searchの技術選定後とし、不必要な待機と早過ぎる確定を分ける。

#### L2 Gate前 Merge DAG

```text
L2-M2-S1 → L2-M1-S1 → L2-M1-S2 → L2-M1-S3
L2-M1-S1〜S3 → L2-M2-S2 → L2-M2-S3 → L2-M2-S4
L1-M2-S6 → L2-M3-S1 調査着手
L2-M1-S1 + L2-M2-S2 → L2-M3-S1 ADR確定・Merge
L2-M3-S1 + L2-M1-S3 + L2-M2-S3・S4 → L2-M3-S2
上記全設計成果物 → L2 Design Freeze Gate
```

### worktree起点と共有物の所有規則

- 最初のTaskだけ、Issue #1、承認済みTask Map、決定記録を含むplanning baseline Commitから分岐する。
- 依存Taskは、必要な前提TaskがMergeされた統合Commitから分岐する。並行兄弟だけが同じ前提Commitを共有する。
- 実作業開始時にplanning baselineのSHAを実行記録へ固定する。承認前のため、現時点ではSHAを決めない。
- 設計worktreeは`docs/task-map.md`、`docs/task-decomposition-decisions.md`、接続逆引き台帳を編集しない。これらはplanning integration ownerだけが更新する。
- Traceability leafが不整合を発見した場合、上流文書を直接編集せず、元Ownerへ差し戻して再統合する。
- 各Directoryの`README.md`は表で指定した終端leafだけが更新し、成果物一覧・依存・参照リンク以外の設計本文を再記述しない。
- L1公開JSON SchemaはL1-M2-S4だけが更新し、後続Taskはread-onlyで利用する。
- Gate前に本番コードのPathを仮決めしない。本番Pathは技術選定と詳細設計を反映したL2 Design Freeze Gate記録の単一表で確定する。

### 接続逆引き台帳の配置案

- 新規派生文書: `docs/task-connections.md`
- 唯一の正本: `docs/task-map.md`のleaf直接依存
- Owner: planning integration ownerのみ
- Task Mapから`docs/task-connections.md`へ明示リンクする。
- 必須列: `Producer leaf`、`直接の後続利用者`、`利用成果物`、`依存時点`、`利用方法`、`Gate／Release`、`worktree所有記録への参照`。
- 同一Producer／Consumerでも成果物または依存時点が異なる場合は別行にする。
- 推移的依存と親roll-upは含めず、終端leafは後続利用者を`なし`とする。
- Gate／ReleaseはMilestoneとして記録し、Task IDを付けない。
- 生成元Task Map revisionを記録し、全直接辺が過不足なく一度だけ存在することを検証規則とする。
- 逆引き台帳から新しい依存判断を追加しない。差異があればTask Mapを承認後に修正し、その後に台帳を再生成する。

### 妥当性確認結果

- `revise_existing`: Pathをすべて`docs/design/...`の完全形で記録すること、`schema-catalog.md`を`schemas/**`配下へ置くこと、`L1-M2-S1`の直接依存を`L1-M1-S5`へ具体化すること。
- `keep_existing`: 20 leafという対象数、各leafの粒度、終端leafによるREADME所有、`docs/task-connections.md`を派生表示として分離する方針。
- `leaf追加・merge・move・remove_defined`: なし。

### 承認後に行う記録更新

1. `docs/task-map.md`へ上記Path／Owner／Merge前提と`docs/task-connections.md`へのリンクを反映する。
2. `docs/task-connections.md`を新規作成し、全leafの直接接続を実レコードとして記録する。
3. GitHub Issueテンプレートの`出力先・所有範囲`と`依存関係`へ、Taskごとの実値を転記する。
4. 実作業開始前にplanning baseline Commitを作り、そのSHAを実行記録へ固定する。

### 再承認単位

本決定の20 leaf配置表、Merge DAG、worktree所有規則、接続逆引き台帳の配置・生成規則を一括して承認対象とする。承認前はTask Mapと新規台帳を変更しない。

### 承認記録

- 2026-08-11: ユーザーが、20 leafの設計成果物Path、単一Owner、Merge DAG、worktree所有規則、`docs/task-connections.md`の配置・生成規則を一括承認した。
- 承認後の反映は、GitHub Issueテンプレート整備と接続逆引き台帳の実レコード作成を含めて行う。

---

## 決定 037: Leaf Task用GitHub Issueテンプレート

### 状態

作成済み（2026-08-11）

### タスク名

承認済みleafを実行するGitHub Issueテンプレートを作成する

### 分割した理由

Task MapのleafをGitHub Issueへ変換する際に、目的と完了条件だけでは、着手可能時点、完了・Merge条件、Gate、成果物Pathの単一Owner、worktree起点とMerge順を復元できない。Issue本文から新しい依存判断を作らず、承認済みplanning snapshotを実行単位へ転記する共通形式が必要なため、1種類のleafテンプレートとして独立させた。

### 配置

- `.github/ISSUE_TEMPLATE/leaf-task.md`
- 設計、実装、評価はテンプレートを分けず、`タスク種別`と理由付きの`該当なし`で扱う。
- `.github/ISSUE_TEMPLATE/config.yml`は、blank Issueを許可するか未決定で実益のある設定がないため、現時点では作成しない。

### 確定した記載項目

- Task ID、親ID、タスク名、種別
- Planning snapshot commit SHAと、そのSHAに固定したTask Mapリンク
- 原典の固定入力、未実施成果物、このleafで決める未確定事項、選び直さない事項
- 目的、実施内容、対象外、観測可能な完了条件
- 主成果物、書込みPath、単一Owner、read-only入力、共有資産Owner、Gate記録
- 着手依存、完了・Merge依存、Gateへの入力、Gate通過依存、Release条件
- 前提TaskのMerge commit、Gate通過記録、worktree起点、Path非競合を含む着手判定
- worktreeの起点、Branch、所有Path、並行／直列、Merge前提・順序、統合先
- Task固有の検証方法、合格条件、Evidence、後続評価への参照
- 正本との差異を発見した場合の停止・再承認手順

### 正本・派生表示との境界

- 直接依存の正本は`docs/task-map.md`とし、Issueから新しい依存辺を作らない。
- 後続利用者はIssueへ重複転記せず、planning snapshot SHAに固定した`docs/task-connections.md`へリンクする。
- GitHubの`Blocked by`表示は補助情報であり、依存関係の正本にしない。
- IssueのCloseだけでは後続Taskを解放せず、完了条件、Merge、Gate／Releaseを確認する。

### 妥当性確認結果

- 判定: `revise_existing`後に作成可能。
- 補正: `Task Map revision`をPlanning snapshot SHA＋固定リンクへ変更し、Gate入力とGate通過依存を分離した。
- 補正: 後続利用者の完全転記と変更禁止Pathの全列挙を廃し、派生台帳リンクとdefault-denyの書込み規則へ変更した。
- 補正: L6共有評価を各実装leafで再実施せず、Task固有検証と後続評価参照を分離した。
- テンプレート追加によるタスク数、依存DAG、Gate、Release、Path Ownerの変更はない。

---

## 決定 038: 承認済みTask MapのGitHub Issue移行

### 状態

実行承認済み（2026-08-11）

### タスク名

L1〜L6の親tracking Issueと全leaf Issueを作成する

### 分割した理由

Task Mapは計画の正本だが、実行時の担当、進捗、PR、検証EvidenceをTask単位で追跡できない。大分類・中分類を進捗集約用のtracking Issue、小分類を実行IssueとしてGitHubへ写像し、親子関係と依存関係を混同せず運用する必要があるため。

### 作成対象

- 大分類tracking Issue: 6件
- 中分類tracking Issue: 19件
- leaf実行Issue: 116件
- 新規Issue合計: 141件
- 原典Issue #1はrootとして維持し、直下の大分類6件への管理チェックリストだけを追記する。
- Design Freeze Gate、Target Readiness Gate、Release、Final Quality GateはMilestoneでありIssue化しない。

### 親子関係

- Issue #1 → L: 6辺
- L → M: 19辺
- M → S: 116辺
- Task Issue間の親子辺は合計135辺とする。
- 親Issueには直下の子だけを掲載し、孫や推移的依存を追加しない。
- 親Issueは子成果物のroll-upと進捗追跡だけを所有し、子と重複する成果物を作成しない。

### 作成・同期方式

- 同期実装: `scripts/task_issue_sync.rb`
- Task IDの一意Marker: `<!-- knowledge-task-id: ... -->`
- Planning snapshot Marker: `<!-- planning-snapshot: <40桁SHA> -->`
- 作成順: L → M → S。全Issue番号の確定後に親チェックリストと依存リンクを同期する。
- open／closed両方のIssueを照合し、同一Task IDの重複、Markerなしの同名タイトル、Task Map外の管理Issueがあれば書込みを停止する。
- 途中失敗時は作成済みIssueを削除せず、再実行時にMarkerから復元して未作成分を継続する。
- 既存本文と人手追記は管理Marker外に保持し、自動Close・自動Reopen・自動削除を行わない。

### Planning snapshot

Issue作成前に、Task Map、決定記録、完成済み接続台帳、Issueテンプレート、同期スクリプト、タスク分割Skillを同一planning baseline Commitへ固定する。全141 Issueは同じ40桁SHAと、そのSHAに固定したTask Map・接続台帳リンクを持つ。

### 接続記録

- `docs/task-connections.md`は全116 leafの直接接続、終端leaf、明示Gate／Release接続を記録する。
- `L6-M3-S5〜S7`は決定035の再監査結果に従い、着手依存と完了・Merge依存を具体leafへ展開してTask Mapへ同期した。
- GitHubの`Blocked by`やチェックリストは補助表示であり、依存関係の正本はPlanning snapshotのTask Mapとする。

### 実行後の監査条件

- Task IDとIssueが141件で一対一である。
- L=6、M=19、S=116、親子辺=135である。
- 全Issueのタイトル、主成果物、直接依存、Planning snapshot SHAが正本と一致する。
- Gate／Releaseを表す管理Issueが0件である。
- leaf 116件が必須Section、成果物Owner、固定リンクを持つ。
- 2回目の同期が作成・更新とも0件になる。

### 承認記録

- 2026-08-11: ユーザーが親tracking Issueを含む全TaskのGitHub Issue作成を承認し、最後まで作成するよう指示した。

---

## 決定 039: GitHub Issue作成前の独立監査と同期安全性修正

- 状態: 修正完了
- 対象: 全141 Task Issue、原典Issue #1の大分類接続、同期スクリプト
- 日付: 2026-08-11

### 監査結果と反映

- Planning snapshotを作成前に固定し、そのCommitにTask Map、決定記録、接続台帳、テンプレート、同期スクリプト、タスク分割Skillを含める。
- 接続台帳の`L3-M2-S2 → L2-M2-S1`重複を除き、Gate入力の1辺へ統一した。leaf間の直接接続は306辺である。
- `L1〜L5`、`Lx-M1〜M3`等の親範囲を同期時に展開し、依存Taskへのリンク欠落を防止した。
- leaf本文の依存を、着手、完了・Merge、Gate入力、Gate通過、Releaseの5区分で表示する。
- L3〜L6は承認記録のTask固有成果物Pathへ具体化した。L2 Gate後実装だけは、既定どおりTaskごとの専用PathをGate通過記録で実値化する。
- 自動生成領域と手動の進捗・Evidence領域を分離した。再同期は自動生成領域だけを置換し、Marker欠損・重複時は本文を上書きせず停止する。
- 親チェックリストの完了状態は子Issueのopen／closedから再生成し、人手の進捗記録とは分離する。

### 作成開始条件

1. 全116 leafの成果物Path文字列に意図しない同一Owner競合がない。
2. 接続台帳が306直接辺で再生成可能である。
3. Ruby構文確認とPlanning snapshotからの141件抽出が成功する。
4. GitHub上に同一Task ID MarkerまたはMarkerなし同名Issueが存在しない。
5. 作成後の再同期が作成・更新0件となり、全件検証が成功する。

---

## 決定 040: 全TaskのGitHub Issue作成と同期検証

- 状態: 完了
- 対象: 大分類6件、中分類19件、leaf 116件、原典Issue #1
- 日付: 2026-08-11

### 実行結果

- 専用ブランチ`planning/task-issues`を作成し、Planning snapshotを`9bf911b9122bf4ec51ec48312cc310de29bfcdff`へ固定した。
- `L1`〜`L6`をIssue #3〜#8、中分類19件をIssue #9〜#27、leaf 116件をIssue #28〜#143として作成した。
- 原典Issue #1へ6件の大分類Issueを接続した。既存Issue #2は変更していない。
- 各IssueへTask ID、親Issue、子Issue、直接依存、Gate／Release区分、成果物Path／Owner、固定Snapshotリンクを同期した。
- Gate／Releaseは成果物を持たないMilestoneとして扱い、独立Issueは作成していない。
- Task IDとGitHub Issue番号の対応を`docs/github-issue-map.md`へ記録した。

### 検証結果

- 初回同期後: `VERIFY OK L=6 M=19 S=116 TOTAL=141`
- 2回目の同期: 作成0件、更新0件、同じ全件検証に合格
- 独立した`--verify`: `VERIFY OK L=6 M=19 S=116 TOTAL=141`
- GitHub上のPlanning snapshot参照、親子関係、Task ID一意性、固定SHA、必須Sectionの整合を確認済み

### 承認記録

- 2026-08-11: ユーザーが専用ブランチの作成・push、`.codex/hooks.json`の追加コミット、および全141件のIssue作成を承認した。

---

## 決定 041: 本人による明示的な不知を表す8番目のKnowledge Assessment状態

- 状態: 承認済み・Commit前反映中
- 対象: Issue #1、L1-M1-S1、L4-M2、L4-M3、L4・L6のKnowledge Assessment参照、Issue #28〜#32
- 日付: 2026-08-11

### 前提知識

`no_evidence`は、必要な探索を行っても、利用者本人が対象命題を知っているとも知らないとも判断できる根拠を観測していない状態である。利用者本人が「知らない」と直接述べた場合は根拠が存在するため、`no_evidence`ではない。`uncertain`は根拠の不足や競合で一つの状態へ確定できない場合であり、明確な不知申告とは異なる。

### 問題

Issue #1は「知っている／知らない」という本人申告を根拠候補に含めながら、Knowledge Assessmentを7状態に固定していた。そのため、利用者本人が対象命題を知らないと明示し、その申告を採用済みの根拠にした場合を、既存7状態の意味を壊さず保存できなかった。

### 承認された修正

- `reported_unknown`（利用者本人が知らないと明示した状態）を8番目のKnowledge Assessment状態として追加する。
- `reported_unknown`は、利用者本人による明示的な不知申告を採用済みの根拠にした場合だけ使用する。
- 質問、説明依頼、AIによる説明、情報への露出、根拠の不在から`reported_unknown`を設定しない。
- `no_evidence`は未観測、`uncertain`は判定不能、`contradicted`は命題と矛盾する理解として、意味を変更しない。
- Knowledge Assessmentの状態数を参照するTask Map、接続台帳、Issue #1、Issue #28〜#32、後続設計は8状態へ同期する。

### 理由

未観測と本人による明示的な不知を区別し、状態から採用済みの根拠まで正確に説明できるようにするためである。既存7状態の意味を変更して異なる状況を一つへまとめることを避ける。

### 影響範囲と変更しない事項

- Knowledge Assessmentの状態数と、`reported_unknown`を扱う後続設計・評価だけを変更する。
- Recognition Gainの7種類、Reading Recommendationの3種類、Task数、階層、依存DAG、Gate数、成果物Ownerは変更しない。
- 過去の決定記録にある「7状態」は当時の承認内容として保持し、本決定が後から変更したことを追跡可能にする。

### Owner境界とCommit分離

- `docs/design/knowledge-model/logical-schema.md`は`L1-M1-S1`の単一Owner成果物として扱う。
- `docs/task-map.md`、`docs/task-connections.md`、`docs/task-decomposition-decisions.md`は、利用者が承認した8状態を計画へ同期するための明示的なplanning integration変更として扱う。
- `scripts/task_issue_sync.sh`は、親Taskをleafの着手依存へ誤って追加する問題を再発させず、Ruby実行環境への暗黙依存をなくすためのIssue同期ツール変更として扱う。
- `.agents/skills/explain-with-context/`は、利用者が別途依頼したプロジェクトSkillであり、Issue #28の成果物へ含めない。
- `.agents/skills/conduct-task-discussion/`と、同Skillをworktree引継ぎで指定するための`.agents/skills/manage-task-worktrees/SKILL.md`および`references/lifecycle.md`、`references/session-handoff-template.md`、`scripts/manage_worktree.sh`は、利用者が別途依頼した議論・引継ぎ手順の変更であり、Issue #28の成果物へ含めない。
- 利用者がCommitを承認した後も、上記のOwnerごとにCommitを分け、Issue #28の成果物Commitへ別Ownerの変更を混在させない。

### 承認記録

- 2026-08-11: ユーザーが選択肢Aを採用し、`reported_unknown`を8番目の状態として追加することを承認した。
- 2026-08-11: ユーザーが、発見済みの矛盾をCommit直前まで修正・同期するよう指示した。Commitは別途明示承認を得るまで行わない。

---

## 決定 042: Task Issue同期ツールをBashへ置換する

- 状態: 承認済み・Commit前反映中
- 対象: `scripts/task_issue_sync.sh`、現行の再生成手順、Task Issue同期の検証手順
- 日付: 2026-08-11

### 前提知識

`task_issue_sync`は製品本体のKnowledge CLIではない。承認済みTask Mapを読み、Task数・親子関係・依存関係・成果物Ownerを検査し、接続台帳を生成し、必要な場合だけGitHub Task Issueを同期する計画管理用ツールである。

Bashはコマンドをどの順番で実行し、途中で失敗した場合にどこで停止するかを記述する実行環境である。`jq`はJSONを項目単位で読み書きするコマンドであり、GitHub Issue本文に含まれる日本語、改行、Markdown記号を通常の文字列分割で壊さないために必要である。対応版はBash 3.2以上と`jq` 1.6以上とし、スクリプトが実行前に版を検査する。`rtk`はこのリポジトリでコマンドを実行する共通の入口、`git`は固定CommitからTask Mapを読む手段、`gh`はGitHub Issueを読み書きする手段である。

### 問題

従来の`scripts/task_issue_sync.rb`はRubyで記述されているが、このリポジトリはRubyの版や導入方法を固定していない。OS付属Rubyへの暗黙依存を残すと、将来同じ手順を実行できるかをリポジトリだけでは判断できない。

### 承認された修正

- 同期ツールをBashで記述した`scripts/task_issue_sync.sh`へ置換し、同値性確認後にRuby版を削除する。
- JSONの解析には`jq`を使用する。外部ライブラリや製品CLIの実装言語は追加で決定しない。
- Bash 3.2以上と`jq` 1.6以上を必須とする。これは、使用する正規表現、`scan`、`capture`、JSON変換の機能条件を固定し、実行環境の違いを開始時に検出するためである。
- `--check`は固定Planning snapshotだけを検査し、`--render-connections`は現在のTask Mapから接続台帳を生成する。どちらもGitHubへ接続しない。
- `--verify`は固定Planning snapshotとGitHub Issueを読み取って照合するが、書込みを行わない。
- `--apply`だけがGitHub Issueの作成・自動生成領域の更新・Issue #1の子一覧更新を行う。自動Close、自動Reopen、自動削除は行わない。
- 自動生成領域の外にある人手の進捗・Evidenceは保持し、管理Markerが欠落または重複している場合は書込み前に停止する。

### 同値性の合格条件

- 固定Planning snapshotから大分類6件、中分類19件、leaf 116件、合計141件を同じ条件で検査できる。
- 現在のTask Mapから生成する直接接続306件、終端leaf 3件、明示Gate／Release 19件の全文が従来版と一致する。
- 全141件のタイトル、本文、親子関係、依存区分、固定リンク、成果物Path、子Issue順を従来版と比較する。
- leaf Task IDを親Task IDとして重複解釈しない。
- GitHubへ書き込む`--apply`は、利用者が別途明示承認するまで実行しない。

### 履歴の扱い

決定038〜040にあるRuby版の記載は、当時そのツールでIssue #3〜#143を作成・検証した事実であるため変更しない。本決定以降の現行手順だけをBash版へ切り替える。

### Owner境界とCommit

本変更はplanning integration用ツールのOwner変更として扱い、Issue #28の論理スキーマ成果物およびプロジェクトSkillのCommitへ混在させない。Commit、push、PR、Issue更新、Mergeは、それぞれ利用者の明示承認を得た後に行う。

### 承認記録

- 2026-08-11: ユーザーがRuby版を残さず、シェルスクリプトで書くよう指示した。
