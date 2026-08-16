# DEC-FEAT-023: 明示的なKnowledge Store指定

- **Status:** decided
- **Level:** L3
- **Decision source:** 利用者の「承認します」（2026-08-16、Issue #233作成前の公開`--store` option提案への承認）
- **論点:** Codex sandboxを含む書込み制約のある実行環境で、Knowledge CLIのStoreをどのように選択可能にするか。

## 現在の状態と必要性

DEC-FEAT-005は指定なしの通常利用にOSのユーザー設定領域を採用し、保存先optionを初期提供の対象外とした。この既定StoreはmacOSではユーザーのLibrary配下となるため、Codex workspace sandboxから開けず`storage_error`となる。`~/.knowledge`への既定変更もhome領域にあるためsandboxの書込み制約を解消せず、実行場所ごとに既定データを分散させる。

Issue #233は、既定Storeを維持したまま、利用者が書込み可能なStoreだけを明示選択できることを求める。これは全11operationが共有する公開CLI・保存先・運用契約を変更するためL3である。

## 選択肢

### A. invocationごとの公開`--store` optionで絶対パスを指定する（採用）

```text
knowledge --store /absolute/path/to/knowledge.db search-text --query channel
```

- **利点:** Codex sandboxではworkspace内のファイルを選べる。既定Storeを変えず、呼出しごとに対象Storeが明確である。
- **欠点:** workflow実行時も同じoptionを明示して渡す必要がある。

### B. 既定Storeを`~/.knowledge`またはworkspaceへ移す

- **利点:** invocationの指定は不要。
- **欠点:** `~/.knowledge`はsandboxから書込不能なhome領域であり、workspace既定化は実行場所ごとに知識を分散させる。既存のOS標準配置の契約も破るため採用しない。

### C. 環境変数または設定ファイルでStoreを指定する

- **利点:** command lineを短くできる。
- **欠点:** 見えない設定がStoreの選択を左右し、workflow・CI・通常利用の再現性を下げる。設定の優先順位・漏洩・永続管理も追加で必要になるため採用しない。

## 決定

利用者承認により選択肢Aを採用する。これはDEC-FEAT-005の「保存先optionを追加しない」という制約を、**明示Store選択に関してのみ置き換える**。DEC-FEAT-005が定めた指定なしの既定Storeと初回初期化は維持する。

- 公開構文は `knowledge --store <absolute-path> <operation> [operation options]` とする。`--store`はoperationより前に一度だけ置く。
- `<absolute-path>`は空でない絶対SQLiteファイルパスとする。相対パス、重複、値なしは`validation_error`（stderr JSON、exit 2）。
- 指定Storeはfallbackなしで唯一の接続先とする。親ディレクトリ作成、DB open、migration失敗は`storage_error`（stderr JSON、exit 1）。
- 指定なし時だけ、既定Store `os.UserConfigDir()/knowledge/knowledge.db` を解決する。
- 環境変数、設定ファイル、外部migrationディレクトリ、権限昇格は追加しない。

## 影響範囲

- 全11CLI operationの共通入力境界とStore composition。
- process境界integration testのStore隔離方式。
- `verification/`内でSkillsが使用するCLI invocation。
- DEC-FEAT-005、Initial Design、FEAT-001共通CLI資料の後続整合更新。

## 承認後に更新する資料

- `docs/design/architecture.md` と `cross-cutting-concerns.md`
- FEAT-001のCLI入力規約
- FEAT-007の詳細設計、task、handoff
- `verification/README.md` と検証用SkillのStore指定手順
