# DEC-FEAT-007: `search-temporal` の時点・有効期間照会契約

- **Status:** decided
- **Level:** L3
- **論点:** `search-temporal` を、時点情報を持つ現行 Assertion の候補取得に留めるか、指定時点または有効期間へ機械的に照合する公開検索として定めるか。

## 現状と、なぜ今決める必要があるか

実装済みの `search-temporal` は、`--concept` と Scope を selector とし、Temporal Metadata を持つ現行 Assertion を返す。これは [操作契約](../design/commands/search-temporal.md) と [UC-04](../design/use-cases/uc-04-temporal.md) に一致する。

一方、[TASK-001-05](../tasks.md) と Issue #181 は、指定した時点または有効期間に合致する Assertion を返し、不正な時刻を `validation_error` とすることを要求する。SQLite schema には `valid_from`、`valid_until` と時点・有効期間検索用の `idx_temporal_window` があるが、現在の公開optionとSQLには照合値・照合規則がない。

時点／期間の入力追加は公開CLI契約を変更し、null境界・期間一致・観測時刻の意味を固定する。そのため、実装へ進む前にL3決定を要する。

## 確定済みの制約

- CLI は候補の新しさ・正しさ・優先度、または知識状態を判断しない（BR-010、UC-04）。
- `valid_from` と `valid_until` は任意の Temporal Metadataであり、両方指定される場合だけ前後関係を持つ（[database schema](../design/database-schema.md)）。
- 現在の `search-temporal` は Concept または Scope の少なくとも一方を selector として必須にする。
- 公開CLI optionの追加・変更はL3の人間承認対象である。

## 選択肢

### 1. 現行契約を維持し、TaskとIssueの時点／期間照合要件を訂正する

`search-temporal` は、Temporal Metadataを持つ現行 Assertion の候補取得だけを行う。

- **利点:** 承認済みの公開CLI契約を変えず、現在の実装・operation契約・UC-04と一致する。
- **欠点:** 時点差分検索はメタデータの有無による候補取得に留まり、TaskとIssueが要求する指定時点／有効期間の照合を提供しない。

### 2. 時点・有効期間のフィルターを `search-temporal` に追加する（推奨）

現在のselectorを維持し、RFC 3339 UTCの時点または有効期間を公開optionで受け、`valid_from`／`valid_until` だけを機械的に照合する。`observed_at` と `last_verified` は結果表示に留め、有効性を推論しない。

- **利点:** TASK-001-05、Issue #181、時点・有効期間用Indexの目的と整合する。時刻比較に限定するため、CLIの意味判断禁止にも抵触しない。
- **欠点:** option名、時点と期間の入力規則、境界包含、片側・両側nullの扱いを公開契約として確定する必要がある。

### 3. 時点・有効期間検索を別operationまたは後続Featureへ分離する

`search-temporal` は候補列挙に固定し、時点照会を新たに設計する。

- **利点:** 現行v1契約を保ったまま、より広い時系列検索を独立して設計できる。
- **欠点:** Issue #181の目的を満たさず、初期提供の論理検索を追加で分割する根拠がない。

## 決定

利用者が選択肢2を承認した。`search-temporal` は既存の Concept／Scope selector に加え、次の任意の時点条件を受ける。

- `--at` は0または1回指定できるRFC 3339 UTC時刻であり、保存済み有効期間がその時点を**包含**するとき一致する。
- `--valid-from` と `--valid-until` は、両方を1回ずつ指定するRFC 3339 UTCの検索期間である。二つは `--at` と同時に指定できず、開始が終了以下でなければならない。
- 検索期間は保存済み有効期間と1点以上で**重複**するとき一致する。開始・終了はいずれも包含する。
- 保存済みの片側 `null` は開放境界とする。両側 `null` は有効期間を表さないため、`--at` または検索期間を指定した照会では一致しない。時点条件を省略した従来の候補取得では、引き続き返す。
- `observed_at` と `last_verified` は時点照合に使わず、結果表示だけに使う。
- 入力のRFC 3339 UTC時刻は、検証後にUTCの固定幅 `YYYY-MM-DDTHH:mm:ss.fffffffffZ`（小数秒9桁）へ正規化して保存・SQL bind・出力する。これによりSQLiteのTEXT比較と時系列順を一致させる。時刻として同値な入力表記は同じ値になる。既存StoreはSQLite adapter内のschema v2 Go migration handlerで`temporal_metadata`の全非null時刻を同じ形式へ変換してから照会に使う。

これにより、既存selectorだけの候補取得との後方互換性を保ちながら、Issue #181が要求する指定時点・有効期間の機械照合を提供する。CLIは時刻以外から有効性、新しさ、正しさを推論しない。

## 承認後の反映先

- `search-temporal` のoperation契約、UC-04、database schemaのIndex説明、schema v2 migration、設計レビュー。
- TASK-001-05 とimplementation handoff。
- RFC 3339 UTC（小数秒ありを含む）の正規化、不正時刻、時点・期間の同時指定、期間の逆順、上下限境界、null境界、期間重複、Concept Alias、Scope AND、空結果、read-only性を確認するCLI process境界テスト。
