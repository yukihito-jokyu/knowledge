# FEAT-001 実装タスク

> **前提:** 利用者による詳細設計の明示承認済み。設計準備監査の結果は pass であり、実装時に公開契約や業務規則を再決定しない。

## Design Readiness Audit

| 監査項目 | 結果 | 根拠 |
| --- | --- | --- |
| 利用者による詳細設計承認 | pass | `.ai/workflow/state.yaml` の `features.FEAT-001.human_design_review: approved` |
| 独立設計レビュー | pass | [design-review.md](design-review.md) の総合判定 `pass` |
| 永続データ・履歴・migration | pass | [design.md](design.md)「Contract Completeness」、[database-schema.md](design/database-schema.md) |
| CLI option・JSON・error・exit契約 | pass | [design.md](design.md)「Interfaces / Data」、[CLI入力規約](design/cli-input-conventions.md)、全11操作資料 |
| operation別I/O・validation・DB read/write・transaction・図 | pass | [command-catalog.md](design/command-catalog.md)、[database-access-map.md](design/database-access-map.md)、`design/commands/` の全11資料 |
| 矛盾候補検索の公開契約 | pass | [DEC-FEAT-004](decisions/DEC-FEAT-004.md)、[search-contradictions.md](design/commands/search-contradictions.md) |
| Fixture・受入観点 | pass | [design.md](design.md)「Acceptance / Test Design」、各操作資料 |
| Go実装基盤・責務境界 | pass | [architecture.md](../../design/architecture.md)「Go module と配置規約」、[DEC-ARCH-002](../../design/decisions/DEC-ARCH-002.md) |

**結論:** pass。実装者が振る舞い、データ、公開契約、または実装基盤を新たに決める未確定事項は検出しなかった。

## 固定済み実装基盤

以下はTaskが新規に決める事項ではなく、利用者承認済みのInitial Designを実装時に守る制約である。

- Knowledge CLIはGo 1.26以上の単一Go moduleとする。SQLite driverはCGOを必要としない`modernc.org/sqlite`を用いる。
- CLI process entry、application、domain、SQLite persistence、migration asset、fixture、integration testの責務と依存方向は[architecture.md](../../design/architecture.md)「Go module と配置規約」に従う。
- migration assetはGo標準の`embed`で実行バイナリへ同梱し、実行時の外部migrationディレクトリ、保存先option、設定ファイルを追加しない。
- 共通チェックは`gofmt`、`go test ./...`、`go vet ./...`とする。依存管理では`go.mod`／`go.sum`と`go mod tidy`を用い、第三者lintは必須化しない。

## タスク一覧

| ID | タイトル | 論理領域 | 依存 |
| --- | --- | --- | --- |
| TASK-001-01 | SQLite v1 schema と migration 契約を実現する | persistence / migration | なし |
| TASK-001-02 | 共通 CLI 入出力・エラー境界を実現する | interface | TASK-001-01 |
| TASK-001-03 | Assertion と Evidence の取得・字句／Concept検索を実現する | retrieval | TASK-001-01, TASK-001-02 |
| TASK-001-04 | Relation検索と矛盾候補検索を実現する | retrieval | TASK-001-01, TASK-001-02 |
| TASK-001-05 | 時点差分検索を実現する | retrieval | TASK-001-01, TASK-001-02 |
| TASK-001-06 | Assertion作成と派生字句Index維持を実現する | mutation | TASK-001-01, TASK-001-02 |
| TASK-001-07 | Evidence追加・revision・置換履歴を実現する | mutation | TASK-001-03, TASK-001-06 |
| TASK-001-08 | 共通fixture・CLI process検証基盤を整備する | verification | TASK-001-02 |
| TASK-001-09 | 全操作の契約・履歴・失敗原子性を横断検証する | verification | TASK-001-03, TASK-001-04, TASK-001-05, TASK-001-06, TASK-001-07, TASK-001-08 |

## TASK-001-01: SQLite v1 schema と migration 契約を実現する

- **目的:** FEAT-001 が所有する永続データ、制約、Index、schema version 1 を、再実行可能で原子的な migration として利用可能にする。
- **関連要件:** REQ-009、REQ-013、REQ-014、BR-004、BR-008、NFR-004、NFR-005、CON-002、CON-005。
- **論理領域:** persistence / migration。
- **作業内容:** Assertion、immutable revision、Scope、Evidence、Concept・Alias、Relation、Temporal Metadata、現行 Assertion の字句派生Index、および migration version の永続表現を実現する。初期化済み・未初期化・部分的または不整合なschemaの判定、単一transactionによる作成・version記録、rollbackと再実行可能性を設計どおりに扱う。
- **受入条件:**
  - schema v1 が全table・制約・参照・検索Index・字句派生Indexを持ち、Assertionとその履歴を失わずに表現できる。
  - `contradicts` と `supersedes` は Assertion→Assertion にだけ保存可能であり、他のRelation種別は設計どおり Assertion／Concept endpointを扱える。
  - migration は空Storeで全schemaとversion 1を同一transactionで作成し、適用済みの場合は変更せず成功する。
  - 部分schema・不整合なmigration履歴・DDL失敗では `storage_error` 相当となり、version記録や部分schemaを残さない。
  - v1の正規データは、将来のSemantic Indexを Assertion ID・現行本文・Scope・Concept から再構築できる。
  - Go 1.26以上の単一Go moduleで、CGOを必要としない`modernc.org/sqlite`を用いてmigrationを実行できる。
  - migration assetはGo標準の`embed`で同梱され、隔離したSQLite DBへ適用して初回適用・再実行・失敗時の原子性を確認できる。
- **依存:** なし。
- **対象外／注記:** 保存先、設定、バックアップ、同期、Semantic Indexの実装は対象外。DDLの正規根拠は [database-schema.md](design/database-schema.md)。SQLite adapter、連番migration asset、および実行バイナリへのasset同梱はInitial Designで固定済みの責務境界に従う。

## TASK-001-02: 共通 CLI 入出力・エラー境界を実現する

- **目的:** 全11操作が名前付きoptionを受け、決定論的なJSON response・stderr／stdout・終了コードを共通契約どおりに提供できるようにする。
- **関連要件:** REQ-014、NFR-005、CON-002、CON-003、DEC-FEAT-002、DEC-FEAT-003。
- **論理領域:** interface。
- **作業内容:** `knowledge <operation> --<option> <value>` の入力規約、繰返しoptionと複数値グループの順序規則、success/error envelope、error code、終了コード、stdout／stderrの分離を実現する。未知option、値不足、単一値optionの重複、グループ不完全、型・enum・空文字の入力不正を各操作のvalidationへ一貫して接続する。
- **受入条件:**
  - 全11操作が名前付きoption入力のみを受け、標準入力request JSONを公開契約として要求しない。
  - 成功時はstdoutに`ok: true`のJSON objectを1件だけ、失敗時はstdoutなし・stderrに`ok: false`のerror objectを1件だけ出力する。
  - `validation_error`=2、`not_found`=3、`conflict`=4、`storage_error`／`internal_error`=1、成功=0の終了コードを守る。
  - `null`、field省略、空配列、固定値の意味は各操作の項目定義と矛盾しない。
- **依存:** TASK-001-01（migration失敗を共通のstorage errorとして扱うため）。
- **対象外／注記:** 検索語生成、操作選択、知識状態判定はCodexの責務である。正規根拠は [CLI入力規約](design/cli-input-conventions.md) と [design.md](design.md)「Interfaces / Data」。

## TASK-001-03: Assertion と Evidence の取得・字句／Concept検索を実現する

- **目的:** 保存済みの現行 Assertion とEvidenceを、文字列またはConceptから候補化・詳細取得できるようにする。
- **関連要件:** REQ-009、REQ-010、REQ-012、REQ-014、BR-002、BR-004、BR-010、CON-004、CON-005。
- **論理領域:** retrieval。
- **作業内容:** `search-text`、`search-concept`、`get`、`get-evidence` を設計どおりに実現する。字句検索では本文、Concept、Alias、Scope、API名、Identifierを対象に現行 Assertionを返す。Concept検索、Assertion詳細、Evidence履歴の取得では、JSON項目・順序・空結果・not-foundの契約を守る。
- **受入条件:**
  - `search-text` は指定文字列の照合fieldを示す候補を返し、空Storeまたは一致なしは成功の空配列となる。
  - `search-concept` はConcept名またはAliasからConceptと紐付く現行Assertionを返し、一致なしの表現は項目定義どおりである。
  - `get` は現行本文、Scope、Concept・Alias、Temporal Metadata、revisionを正しく返し、指定Assertionが不在なら`not_found`となる。
  - `get-evidence` は指定AssertionのEvidence履歴を返し、Evidenceが0件であることを未知・`known`等の意味判断に変換しない。
  - read操作はStoreを書き換えず、全field、固定値、`null`／省略／空配列の意味は各操作資料に一致する。
- **依存:** TASK-001-01、TASK-001-02。
- **対象外／注記:** Semantic Search、検索継続の判断、Evidenceの強度・知識状態の評価は対象外。操作別の正規契約は [command-catalog.md](design/command-catalog.md) の各資料。

## TASK-001-04: Relation検索と矛盾候補検索を実現する

- **目的:** 保存済みRelationを、検索起点と保存方向を取り違えずに検索し、矛盾候補から相手AssertionのEvidence確認へ接続できるようにする。
- **関連要件:** REQ-009、REQ-012、REQ-014、BR-007、BR-010、CON-004、CON-005、DEC-FEAT-004。
- **論理領域:** retrieval。
- **作業内容:** `search-related` と `search-contradictions` を実現する。前者は Assertion／Conceptのselector、Relation種別、検索起点から見た相手と保存方向を返す。後者は`contradicts`に限定し、selectorに一致するAssertionを`seed`、常に反対側のAssertionを`target`、保存時の向きを`direction`として返す。Concept selectorは所属Assertionへ展開し、両端が一致するRelationはseedごとに別結果とする。
- **受入条件:**
  - `search-related` は全許可Relation種別で、検索起点から見た相手と`incoming`／`outgoing`を返し、保存上のsource／targetとの比較を利用側へ強制しない。
  - `search-contradictions` は Assertion→Assertion の`contradicts`だけを対象にし、`target.id`を常に`get`および`get-evidence`のAssertion IDとして利用できる。
  - 検索起点が保存上target側でも、結果の`target`が検索起点自身にならない。
  - Concept selectorでRelation両端が対象となる場合、relation IDだけで重複排除せず、seedごとの2結果を返す。
  - Relationの存在を、CLIがユーザー理解や矛盾の意味の判定へ変換しない。空結果は成功の空配列である。
- **依存:** TASK-001-01、TASK-001-02。
- **対象外／注記:** Relationの意味推論、矛盾の解消・知識状態評価は対象外。`search-contradictions`の固定契約は [DEC-FEAT-004](decisions/DEC-FEAT-004.md)。

## TASK-001-05: 時点差分検索を実現する

- **目的:** 保存済みTemporal Metadataと現行Assertionを、指定した時点・有効期間・Scopeで決定論的に検索できるようにする。
- **関連要件:** REQ-009、REQ-012、REQ-014、BR-010、CON-004、CON-005。
- **論理領域:** retrieval。
- **作業内容:** `search-temporal` のoption、RFC 3339 UTC時刻、Scope集合照合、現行AssertionとTemporal Metadataの応答を設計どおりに実現する。
- **受入条件:**
  - 指定時点または有効期間に合致する保存済みTemporal Metadataを持つ現行Assertionだけを返す。
  - 指定Scopeがある場合は設計したScope集合照合に従い、未指定の場合はScopeで絞り込まない。
  - 不正な時刻・空文字・不完全なScope groupは`validation_error`となり、該当データがなければ成功の空配列となる。
  - 検索はTemporal Metadata以外から時点差分や古さを推論せず、Storeを書き換えない。
- **依存:** TASK-001-01、TASK-001-02。
- **対象外／注記:** `outdated`等の知識状態判定は対象外。正規契約は [search-temporal.md](design/commands/search-temporal.md)。

## TASK-001-06: Assertion作成と派生字句Index維持を実現する

- **目的:** 新しいAssertionと付随データを、重複・参照整合性・履歴保全を守りつつ一貫して作成する。
- **関連要件:** REQ-009、REQ-010、REQ-013、REQ-014、BR-004、BR-008、BR-010、NFR-004、NFR-005、CON-002、CON-004。
- **論理領域:** mutation。
- **作業内容:** `create` でAssertion初版、Scope、Evidence、Concept・Alias、Assertion Alias、Temporal Metadata、Relationを設計どおりに保存する。Concept名・Aliasの横断一意性、同一本文・Scopeの現行Assertionの衝突、Relation endpointの実在・種別制約を検証し、現行Assertionの字句派生Indexを同一transactionで構築する。
- **受入条件:**
  - `create` は成功時だけ、作成された識別子・revisionを設計したJSONで返す。
  - 必須項目、型、空文字、enum、重複Scope／Concept／Alias、同一正規化本文・Scope、Relation endpoint不在・種別不一致を、操作資料のerror分類で扱う。
  - `contradicts` Relationを作成する場合は、sourceとtargetがAssertionであることを検証する。
  - Assertion、付随データ、Relation、字句派生Indexの書込みは一操作・一transactionで行い、いずれの失敗でも全てrollbackする。
  - 作成後の字句検索は現行本文、Concept、Alias、Scope、API名、IdentifierからそのAssertionを検出できる。
- **依存:** TASK-001-01、TASK-001-02。
- **対象外／注記:** 同一性の意味判断、createと他の更新操作の選択、公開Index管理操作は対象外。正規契約は [create.md](design/commands/create.md)。

## TASK-001-07: Evidence追加・revision・置換履歴を実現する

- **目的:** Evidenceの追記、Assertionの現行表現更新、旧知識を残す置換関係を、履歴を失わずに実現する。
- **関連要件:** REQ-013、REQ-014、BR-004、BR-008、BR-010、NFR-004、NFR-005、CON-002、CON-004。
- **論理領域:** mutation。
- **作業内容:** `attach-evidence`、`revise`、`supersede` を実現する。Evidenceはimmutableに追加し、revisionは既存revisionを更新せず新しい現行revisionを作成してScope・Temporal Metadata・字句派生Indexを同期する。置換は旧Assertionを削除せず、Assertion間の`supersedes` Relationを保存し、自己参照・循環・重複を競合として扱う。
- **受入条件:**
  - `attach-evidence` は既存AssertionにEvidenceだけを追記し、対象不在は`not_found`、入力不正は`validation_error`、保存失敗はrollbackとなる。
  - `revise` は新revisionと現行revision参照を原子的に更新し、過去の本文・Scope・Temporal Metadataを消さず、検索用派生Indexを現行revisionへ更新する。
  - `supersede` は旧・新Assertionを保持したままAssertion→Assertionの置換Relationを一度だけ保存し、自己参照・循環・既存同一関係を`conflict`として扱う。
  - 三操作とも成功responseはcommit後だけに返し、失敗時はstdoutなしで部分更新を残さない。
  - `get`／`get-evidence`および検索で、更新後の現行情報と過去Evidence・旧Assertionを設計どおりに確認できる。
- **依存:** TASK-001-03（更新後の現行情報・Evidence履歴・検索到達性を確認するため）、TASK-001-06（作成済みの基礎データと字句Index維持が必要なため）。
- **対象外／注記:** Evidenceの価値・強度、更新操作の選択、知識状態の再評価はCodexの責務。操作別の正規契約は各更新操作資料。

## TASK-001-08: 共通fixture・CLI process検証基盤を整備する

- **目的:** 全11操作を実行プロセス境界で検証するための、隔離SQLite DB・再利用可能なfixture・JSON／終了コードの観測基盤を提供する。
- **関連要件:** REQ-014、BR-004、BR-008、NFR-004、NFR-005、CON-002、CON-003、DEC-ARCH-002。
- **論理領域:** verification。
- **作業内容:** Go標準toolchainで実行できる、隔離SQLite DBを使ったCLI process検証の共通基盤を整備する。migration適用済み空Store、操作ごとの初期データ、stdout／stderr、JSON、終了コード、Storeの事後状態を決定論的に観測するfixtureと共通手順を提供する。
- **受入条件:**
  - 隔離SQLite DBへ埋込みmigrationを適用し、適用済みDBへの再実行が成功することを共通手順で確認できる。
  - CLI processの入力option、stdout、stderr、終了コード、JSON response、Storeの事後状態を、各操作の契約へ対応付けて観測できる。
  - 空Storeと、Assertion、Evidence、Concept、Scope、Relation、Temporal Metadataを含む再利用可能な初期データを用意できる。
  - success、`validation_error`、`not_found`、`conflict`、`storage_error`／`internal_error`の観測方法を提供し、stdout／stderrと終了コードの共通契約を確認できる。
  - 共通基盤は承認済みの`gofmt`、`go test ./...`、`go vet ./...`で利用でき、公開CLI契約を追加しない。
- **依存:** TASK-001-02（CLIの共通I/O・error境界を観測するため）。
- **対象外／注記:** 個別操作の完全な正常系・失敗系、履歴保全、rollback、契約横断の判定はTASK-001-09で行う。Issue #175 の Scenario A〜Jを横断する統合・層別評価（NFR-006）はFEAT-005が担当する。

## TASK-001-09: 全操作の契約・履歴・失敗原子性を横断検証する

- **目的:** 全11操作の公開CLI契約、検索到達性、履歴保全、更新失敗時の無変更を、TASK-001-08の共通基盤で再現可能に確認する。
- **関連要件:** REQ-009〜REQ-014、BR-002、BR-004、BR-007、BR-008、BR-010、NFR-004、NFR-005、CON-001〜CON-005、DEC-FEAT-004。
- **論理領域:** verification。
- **作業内容:** 空Store、字句一致、Concept／Relation、Evidence、矛盾候補、Temporal Metadata、作成・追加・revision・置換、validation／not-found／conflict／storage failureを、操作別設計と共通fixture・CLI process検証基盤により横断確認する。JSON field、fixed value、終了コード、stdout／stderr、migration再実行、rollback、将来Semantic Index再構築に必要な正規データを検証する。
- **受入条件:**
  - 全11操作に正常、空結果または該当する失敗結果を含む検証があり、公開JSON・終了コード・stdout／stderrの契約を確認できる。
  - 空Storeの検索は空配列で成功し、`known`、`no_evidence`、矛盾や古さの意味判断を返さない。
  - 本文、Concept、Concept Alias、Scope、API名、Identifierからの字句検索と、Evidence履歴、Relation、Temporal検索を確認できる。
  - `search-contradictions` のincoming／outgoing、Concept selectorの両端一致、`target.id`の後続取得可能性、Assertion→Assertion制約を確認できる。
  - mutation失敗ではAssertion、Evidence、Relation、revision、字句派生Index、migration versionの部分更新が残らないことを確認できる。
  - migrationの初回適用・再実行、及びv1 catalogにSemantic Searchが現れず後続Indexを正規データから再構築できる前提を確認できる。
- **依存:** TASK-001-03、TASK-001-04、TASK-001-05、TASK-001-06、TASK-001-07、TASK-001-08。
- **対象外／注記:** Codexによる検索戦略、Search Trace、Knowledge Assessment、Article評価のfixtureは後続Featureの範囲。Issue #175 の Scenario A〜J を横断する統合・層別評価（NFR-006）はFEAT-005が担当し、本Taskはそれに利用可能なCLI／Storeの決定論的fixtureだけを担う。

## Feature 全体の対象外

- `search-semantic`、Embedding、Vector Index（FEAT-006）。
- Query生成、検索継続・停止、同一性・Evidence価値の意味判断、`known`等のKnowledge Assessment。
- 保存先、設定、暗号化、バックアップ、同期、共有、リモート接続。
