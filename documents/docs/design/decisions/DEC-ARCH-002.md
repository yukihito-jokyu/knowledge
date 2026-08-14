# DEC-ARCH-002: Knowledge CLI の Go 実装基盤

- **Status:** decided
- **Level:** L3
- **論点:** Knowledge CLIをGoで実装すること、およびGo実装で必要となるSQLite driverと最低対応Go versionをどう固定するか。

## 現在の状態と、なぜ今決める必要があるか

Knowledge CLIは、Codexから受け取った名前付きoptionを処理し、SQLiteへ保存・検索・更新した結果をJSONで返す単一のローカルCLIである。人間承認により実装言語はGoとなった。一方、SQLite接続には標準ライブラリ外のdriverが必要であり、最低Go versionはmodule・CI・配布可能な実行バイナリの互換性を決める。

現在のリポジトリにはGo module、既存のSQLite driver、CI実行環境、最低対応Go versionの根拠がない。開発環境にGo 1.26.1が存在することは確認できたが、これは利用者・CIの互換性保証ではない。

Go採用は確定済みである。SQLite driverと最低Go versionは、決定後に変えると依存、ビルド環境、SQLiteの挙動確認を横断して見直すため、FEAT-001実装開始前に固定する。

## 決定済み

利用者承認により、以下を採用する。

- Knowledge CLI v1と後続Featureの実装言語は **Go** とする。
- 依存は単一のGo moduleで管理する。
- 実行入口、非公開のapplication/domain/persistence、migration asset、fixture／integration testを分離する。
- 共通チェックはGo標準toolchainの `gofmt`、`go test ./...`、`go vet ./...` とする。

Goは単一バイナリとして配布しやすく、決定論的なCLI、SQLite transaction、JSON入出力、後続のローカル検索を同じtoolchainで扱える。これは公開CLI契約を変更せず、実装の責務境界だけを確定する。

## 追加決定

利用者は選択肢 **A** を承認した。SQLite driverには **`modernc.org/sqlite`** を採用し、`go.mod` の最低対応Go versionは **1.26** とする。

- `modernc.org/sqlite` は純GoのSQLite driverとしてmodule依存に追加する。CGOやC compilerをbuild前提にしない。
- 開発・CIではGo 1.26以上を使用し、`gofmt`、`go test ./...`、`go vet ./...` を実行する。
- driverの具体的な接続初期化は `internal/persistence/sqlite` の責務であり、driver固有の挙動や設定を公開CLI契約へ露出しない。

## 選択肢（SQLite driverと最低対応 Go version）

### A. `modernc.org/sqlite` と Go 1.26 以上を採用する（推奨）

- **何が変わるか:** 純GoのSQLite driverをmodule依存として追加し、Go 1.26以上でbuild・test・静的検査を保証する。
- **利点:** C compilerを必要とせず、ローカルCLIをGo toolchainだけでbuildしやすい。現在確認できる開発環境のGo 1.26.1とも整合する。
- **欠点:** driver実装をmodule依存として取得・更新する。Go 1.26未満の環境は対象外となる。

### B. `github.com/mattn/go-sqlite3` と Go 1.26 以上を採用する

- **何が変わるか:** CGOを用いるSQLite driverをmodule依存として追加し、build環境にC compilerとCGO利用可能な環境を要求する。
- **利点:** SQLite C libraryを利用する一般的なdriver選択肢である。
- **欠点:** OS・cross compile・CIにnative toolchainの条件が加わり、単一ローカルCLIの導入・検証が重くなる。

### C. driverと最低Go versionを実装開始後まで保留する

- **何が変わるか:** Goの責務境界とmigration asset運用だけを確定し、SQLite接続・migration・transactionを伴う実装とintegration testは保留する。
- **利点:** 現時点で依存を固定しない。
- **欠点:** FEAT-001のSQLite migration、transaction、integration testを実装できず、今回の実装準備の目的を満たさない。

## 採用理由

初期提供はローカル専用・単一CLIであるため、CGOという追加のbuild前提を持たないAを採用する。これによりSQLiteのtransactionとmigrationをGoの共通test手順で検証しやすい。native SQLite固有の要件は現行要件にないためBは採らず、CはFEAT-001の実装準備を止めるため採らない。

## 決定の影響

- `go.mod` のGo directiveとSQLite driver依存。
- 開発・CIで実行する `gofmt`、`go test ./...`、`go vet ./...` の対象環境。
- `internal/persistence/sqlite` の接続初期化とintegration test。

公開CLIのoption、JSON、終了コード、SQLite schema、保存先、設定、運用仕様は変更しない。
