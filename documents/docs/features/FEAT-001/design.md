# FEAT-001 詳細設計（v1 公開契約）

> **状態:** review_ready — DEC-FEAT-003 により、v1 の入力をrequest JSONから名前付きoptionへ変更した。全11操作へ、バリデーションとは別の入出力項目定義（型・出現回数・固定値・意味）を追加した。DEC-FEAT-004 により `search-contradictions` は検索起点ごとの相手 Assertion を返し、保存方向を `direction` で保持する。独立 Design Review の pass 後は利用者による詳細設計レビューを待つ。

## Feature Summary

Evidence を正規根拠とする Assertion、Concept、Scope、Relation、Temporal Metadata をローカル SQLite に保存し、Codex が決定論的 JSON CLI で検索・取得・更新できるようにする。初期提供の検索は字句・Concept・Relation・Evidence・矛盾候補・時点差分であり、Semantic Search は FEAT-006 へ延期する。

## Scope / Out of Scope

- **含む:** REQ-009〜014 の初期提供、履歴保持、SQLite schema v1、名前付きoption入力とJSON出力を持つCLIの検索・取得・更新操作、字句 Index の維持、DEC-FEAT-005で固定するユーザー固有の既定Store。
- **含まない:** `search-semantic`、Embedding／Vector Index、Codex による意味判断・検索戦略・知識状態判定、保存先の選択・設定・同期・共有、記事・会話の入力境界。

## Related Requirements / Business Rules

REQ-009〜014、BR-002、BR-004、BR-007、BR-008、BR-010、NFR-004〜006、CON-001〜005、DEC-ARCH-001、DEC-REQ-001、DEC-FEAT-001、DEC-FEAT-002、DEC-FEAT-003 を根拠とする。Issue #175 §10〜13、§40、§44、§48 の論理操作・責務境界を継承する。

## Behavioral Scenarios

| Trigger | Main flow | Result |
| --- | --- | --- |
| Codex が検索語または Concept を指定 | CLI が Index と関連を決定論的に照会 | 候補 Assertion と検索根拠を返す。0件は成功の空配列である。 |
| Codex が Assertion ID を指定 | CLI が現行 Assertion または Evidence 履歴を取得 | 詳細または Evidence 一覧を返す。ID 不在は not-found である。 |
| Knowledge Update が作成・Evidence追加・修正・置換を決定済みとして渡す | CLI が入力を検証し、履歴を保持する SQLite transaction を実行 | 成功時だけ更新後の識別子・revision を返す。失敗時は rollback する。 |

CLI は結果から `known`、`no_evidence`、矛盾の意味、同一性、置換の妥当性を判断しない。空の検索結果は `no_evidence` を導出する材料になり得るだけであり、知識状態ではない。

目的別の操作列、各操作を選ぶ責務、および操作の過不足・矛盾の監査は、ユースケース別資料を参照する。

- [UC-01: 文字列と根拠の確認](design/use-cases/uc-01-text-and-evidence.md)
- [UC-02: Concept と Relation の確認](design/use-cases/uc-02-concept-and-relation.md)
- [UC-03: 矛盾候補の確認](design/use-cases/uc-03-contradictions.md)
- [UC-04: 時点情報の確認](design/use-cases/uc-04-temporal.md)
- [UC-05: 新しい Assertion の登録](design/use-cases/uc-05-create.md)
- [UC-06: 既存 Assertion の更新履歴](design/use-cases/uc-06-update-history.md)

## Selected Design Analyses

- **採用:** 正常・入力不正・not-found・競合の flow、CLI／SQLite 境界の sequence、論理データモデル、wire schema、transaction／idempotency、Fixture 観点。
- **不採用:** 状態遷移図。長期ワークフローや CLI 内の知識状態機械はなく、履歴は immutable record と current revision の参照で表現するため。
- **Semantic Search:** 初期提供では not_applicable。FEAT-006 が正規 Assertion ID・現行正規化本文・Scope・Concept から再構築する。

## Responsibilities

| Logical responsibility | Responsibility |
| --- | --- |
| Codex / Knowledge Update | 意味的同一性、Evidence の知識価値・強度、操作選択、検索継続、知識状態の評価を判断する。 |
| Knowledge CLI | 指定入力の validation、ID管理、SQLite read/write、字句／Relation Index 照会、履歴保全、JSON出力を決定論的に実行する。 |
| Knowledge Store | Assertion、revision、Evidence、Concept、Scope、Relation、Temporal Metadata と検索用派生 Index を保存する。 |

上位のデータ所有権表における **Derived User Knowledge State の Owner=Knowledge CLI** は、保存・更新の唯一の経路を指す。BR-010 に従い、Evidenceからの意味的な導出・知識状態の評価は Codex の Knowledge Search が担い、Knowledge CLI は指定された保存を決定論的に実行するだけである。FEAT-001 は導出済み Knowledge State を保存・返却せず、その入力となる Evidence と検索結果だけを扱う。

## Interfaces / Data

### 識別子と共通表現

- `assertion_id`、`concept_id`、`evidence_id`、`relation_id` は CLI が生成する不透明な文字列であり、利用者は意味を解釈しない。
- `scope` は `key` と `value` から成る配列である。同一 assertion revision 内で `key` は一意とする。`aliases` は Assertion に付与する `api_name` または `identifier` の配列であり、現行 Assertion 内で同じ `kind` と `value` を重複させない。
- `temporal` は任意。`valid_from`、`valid_until`、`version_scope`、`observed_at`、`last_verified` を保持でき、時刻はRFC 3339 UTC入力を固定幅 `YYYY-MM-DDTHH:mm:ss.fffffffffZ` へ正規化して保持・出力し、未設定は `null` とする。
- Evidence は `raw_text` と `kind` を必須とする。`kind` は Issue #175 の根拠例を機械的に区別する enum であり、値は command 資料に従う。
- optionの未知名、値不足、単一値optionの重複、複数値グループの不完全な並びは `validation_error` とする。JSON出力の省略可能 field は `null` と区別し、配列は空配列で表す。

### 共通結果 envelope

成功時は stdout に次の JSON object を1件だけ出力する。失敗時は stdout を出力せず、stderr に error object を1件だけ出力する。

```json
{
  "ok": true,
  "data": {}
}
```

```json
{
  "ok": false,
  "error": {
    "code": "validation_error",
    "message": "query must not be empty",
    "field": "query"
  }
}
```

error code は `validation_error`、`not_found`、`conflict`、`storage_error`、`internal_error` である。終了コードは success=0、validation=2、not-found=3、conflict=4、storage/internal=1 とする。端末からの`Ctrl-C`による要求中断だけは、responseの先頭byteを書き始める前にcancelを観測した場合、通常のerror envelopeを出力せず、stdout／stderrを空にして終了コード130で終える（[DEC-FEAT-006](decisions/DEC-FEAT-006.md)）。response開始後の割込みは既出力を取り消さず、選択済みの通常responseを完了する。

### CLI入力とJSON出力

各論理操作を `knowledge <operation> --<option> <value>` として呼び出す。入力の名前付きoption、繰返しデータの順序規則は [CLI入力規約](design/cli-input-conventions.md) と各操作資料に従う。response JSON object は標準出力または標準エラーへ出力する。標準入力のrequest JSON、保存先を選ぶoption、設定は定義しない。

端末の最初の`Ctrl-C`は、実行中の要求を中断する。CLI境界がresponse開始前に中断を観測した場合、success／error JSONを出力せず、stdout／stderrを空にして終了コード130で終了する。これは全11操作とStore初期化に共通する例外であり、入力不正や保存失敗を表すerror codeではない。response開始後は既出力を取り消さない。

### Store の既定保存先と初期化

通常ビルドの `knowledge` は、`os.UserConfigDir()` 配下の `knowledge-cli/knowledge.db` を唯一の既定Storeとする。OSごとの解決結果は、macOSでは `~/Library/Application Support/knowledge-cli/knowledge.db`、Unixでは `$XDG_CONFIG_HOME/knowledge-cli/knowledge.db` または `~/.config/knowledge-cli/knowledge.db`、Windowsでは `%AppData%\\knowledge-cli\\knowledge.db` となる。

- 各操作の実行前に、composition rootが親ディレクトリを作成し、SQLite Storeをopenして未適用migrationを適用する。
- 初回実行時は空Storeを作成する。初期化後の検索は成功の空配列、存在しないAssertionの取得は `not_found` を返す。
- ユーザー設定ディレクトリの解決、親ディレクトリ作成、Store open、migrationの失敗は `storage_error` とする。ただし各Context非対応APIの前後およびresponse開始前に割込みを観測した場合は、`storage_error`より中断（終了コード130・無出力）を優先する。
- 割込みで要求Contextが終了した場合は、初期化・migration・照会・更新のいずれも、response開始前に観測した限り成功JSONを返さない。更新中のtransactionはcommit前にContextを確認してcommitせず、既存のrollbackと履歴不変条件を維持する。
- この既定保存先以外を指定する公開option、環境変数、設定ファイル、実行時の外部migrationディレクトリは提供しない。
- read操作は、上記の初回初期化後、Assertion、Evidence、Concept、Scope、Relation、Temporal Metadata、派生Indexを変更しない。

### SQL の可変長値規約（内部）

操作資料中の `:name[]` は、CLI がその操作内だけで繰返しoptionから組み立てた入力集合または前段 query の行集合を、同数の scalar prepared-statement parameter へ展開することを表す内部表記である。これは公開optionや設定ではない。

- `IN (:name[])` は要素数ぶんの scalar bind に展開する。空集合なら、その `IN` を実行する後続 query を省略して空行集合とする。
- 繰返しoptionは入力順を保ってvalidationするが、SQL結果の順序は各操作資料の `ORDER BY` に従う。
- 前段queryの候補ID集合は、同一read operation内でCLIが保持する。候補が0件なら hydration queryは実行せず、空の `results` を返す。
- 固定個数の値は従来どおり単一の `:name` bind を使う。
- **Scope 集合の例外:** `create` の SQL-02 と `search-temporal` の SQL-01 だけは、検証済みの Scope group を CLI 内部で `[{"key":"…","value":"…"}]` に直列化した単一 bind `:scope_json` を使う。これは入力 JSON でも公開 field でもなく、SQL の組集合照合専用である。Scope 未指定時は必ず `[]` とし、`json_each(:scope_json)` は0行となる。`create` では既存 Scope も0件の Assertion だけが同一 Scope と一致し、`search-temporal` ではScopeによる絞込みをしない。この二つ以外の可変長値は前記の scalar bind 展開を使う。

## Physical / Wire Schema

- SQLite v1 DDL、migration の前後条件、関係・Index は [database-schema.md](design/database-schema.md) を参照。
- operation ごとのoption入力、success／error JSON、exit code、図、SQL は [command-catalog.md](design/command-catalog.md) を入口とする。
- 利用目的ごとの操作列とユースケース整合監査は、上記のユースケース別資料を参照する。
- DB の operation 別 read/write は [database-access-map.md](design/database-access-map.md) を参照。

## Operation Documentation

Operation Documentation Coverage Gate は該当する。初期提供の全11操作を [command-catalog.md](design/command-catalog.md) に列挙し、各資料に I/O、入出力項目定義（型・出現回数・固定値・意味）、バリデーション設計、フローチャート、シーケンス図、DB接続を記載する。`search-semantic` と独立した公開 Index 管理操作は初期提供では not_applicable とする。

## Error / Edge Cases

- 空検索は `ok: true` と空の `results` を返す。CLI は知識状態を返さない。
- 不存在 ID、selector の不足・複数指定、空文字、無効な時刻／enum は `validation_error` または `not_found` で更新しない。
- `supersede` の自己参照・循環、既に同じ置換関係がある場合は `conflict` とする。
- mutation は一操作を一 transaction とし、失敗時に Assertion、Evidence、Relation、Index の部分更新を残さない。
- `Ctrl-C`による中断をresponse開始前に観測した場合は`storage_error`へ写像せず、JSONを出力しない終了コード130とする。ContextはCLI入口からapplication、domain port、SQLite adapterまで同一のまま伝播し、下位層で新たに作り直さない。
- 同一操作の再送を識別するrequest-id／冪等性キーは根拠がないため初期公開契約に含めない。再送時の重複回避は Codex 側の操作選択責務とする。

## Security / NFR Considerations

Store はローカル専用 SQLite、CLI が唯一のアクセス経路である。既定保存先はDEC-FEAT-005に従う。保存先の選択、暗号化、バックアップ、アクセス制御の追加仕様は定めない。OS 標準アクセス制御を前提とし、秘密情報・リモート接続は扱わない。隠しまたは通常非表示のOS標準ディレクトリは利便性のためであり、追加のセキュリティ境界ではない。

## Acceptance / Test Design

- 空 Store の各検索は空配列を返し、未知・`known` 判定を返さない。
- Assertion本文、Concept、Alias、Scope、API名、Identifier の字句検索 Fixture を用意する。
- Evidence追加、revision、supersede 後も過去 record を取得でき、失敗 mutation はDBを変更しない。
- Relation、矛盾候補、時点差分は保存済み relation／metadata のみから候補を返す。
- Semantic Search が初期 CLI catalog に現れず、将来 Index を Assertion ID から再構築できる現行データを保持する。
- 通常ビルドがOS標準のユーザー固有既定Storeを初期化して全操作を実行し、失敗時は公開error／exit契約どおり `storage_error` を返す。
- `integrationtest` build tagにだけ含める非公開同期gateでmigration・read・mutationの開始を確認してから実バイナリへ`Ctrl-C`を送る。通常compositionと隔離したOS設定ディレクトリを使い、stdout／stderrが空で終了コード130となること、mutationの事後DB状態が不変であることをprocess境界で確認する。このgateは通常ビルド、公開CLI option、公開環境変数、設定、Store overrideに含めない。Issue #179 の取得実装では migration／read を検証し、mutation executor を導入する TASK-001-06／07 で mutation の gate・事後DB状態検証を追加する。unit testでは、同一のcancel済みContextがapplicationとSQLite adapterに到達し、更新transactionがcommitしないことを確認する。

## Contract Completeness

| 項目 | 状態 | 根拠・注記 |
| --- | --- | --- |
| 永続データ・履歴 | complete | SQLite v1 と immutable history を DEC-FEAT-002 が承認。 |
| CLI option入力・JSON出力 wire schema | complete | DEC-FEAT-003に基づき、共通規約、全11操作資料、全ユースケース例を更新済み。 |
| 検索・取得規則 | complete | operation別資料に記載。 |
| 更新・重複・原子性 | complete | 操作選択はCodex、保存上の重複制約を明記。 |
| Migration と既定Store | complete | v1と将来互換のmigration、およびDEC-FEAT-005のユーザー固有既定保存先。 |
| 要求キャンセルと割込み終了 | complete | DEC-FEAT-006により、全操作のContext伝播、JSONなし、stdout／stderrなし、終了コード130を固定。 |
| Fixture | complete | 上記 Acceptance / Test Design。 |

## Assumptions / Open Issues

- **DEC-FEAT-002 (L3, superseded_in_part):** JSON出力・exit code・Indexを独立公開操作にしない方針は維持し、request JSON標準入力はDEC-FEAT-003で置き換えた。
- **DEC-FEAT-003 (L3, decided):** 入力は名前付きoptionとし、標準入力request JSONを使わない。
- **DEC-FEAT-004 (L3, decided):** `search-contradictions` の結果は `seed` と常にその相手である `target` を返し、保存方向は `direction` で保持する。矛盾候補を `get` / `get-evidence` へ安全に受け渡すため、`contradicts` Relation は Assertion 間に限定する。
- **DEC-FEAT-005 (L3, decided):** 通常ビルドは `os.UserConfigDir()/knowledge-cli/knowledge.db` を既定Storeとして初期化する。保存先の選択、設定、運用責任を追加しない。
- **DEC-FEAT-006 (L3, decided):** 最初の`Ctrl-C`は要求Contextをcancelし、success／error JSONおよびstdout／stderrを出さず終了コード130で終了する。
- **DEC-FEAT-007 (L3, decided):** `search-temporal` は任意のRFC 3339 UTC時点または閉区間で、保存済み有効期間を機械的に照合する。時点は包含、期間は重複で一致し、片側`null`は開放境界、両側`null`は時点条件がある照会では不一致とする。時刻は固定幅UTCへ正規化し、SQLiteのTEXT比較を時系列比較として用いる。
- **L2:** SQLite schema v1 の table／column／Index 名は、承認済み JSON 契約を満たす内部表現である。
- **既決:** Semantic Search は FEAT-006 まで延期する。
