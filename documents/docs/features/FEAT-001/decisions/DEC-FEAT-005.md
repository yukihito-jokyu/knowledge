# DEC-FEAT-005: Knowledge Store の既定保存先と初回初期化

- **Status:** decided
- **Level:** L3
- **論点:** 公開CLIの保存済みAssertion／Evidenceを、通常ビルドの `knowledge` コマンドがどのローカルSQLite Storeへ接続するか。

## なぜ今決める必要があるか

TASK-001-03 の取得操作は、SQLite Storeへ接続できなければ通常ビルドで実行できない。既存設計は保存先を定義しない一方、`cmd/knowledge` を composition root として定めているため、実行可能なCLIを提供するには既定保存先と初回初期化の契約が必要である。

この決定は全11操作が共有する永続化・運用上の振る舞いを変更するため、L3とする。

## 選択肢

### 1. OSのユーザー設定ディレクトリ配下を既定保存先にする（推奨）

`os.UserConfigDir()` が返すディレクトリ配下の `knowledge/knowledge.db` をStoreとする。

- 初回実行時に親ディレクトリ、SQLite DB、未適用migrationを作成・適用する。
- `os.UserConfigDir()` の解決、ディレクトリ作成、DB open、migrationの失敗は `storage_error` とする。
- 空Storeでは検索は成功の空配列、Assertion指定の取得は `not_found` とする。
- `--store` option、環境変数、設定ファイル、実行時の外部migrationディレクトリは追加しない。

**利点:** OSごとのユーザー単位の標準配置に従い、実行ディレクトリに左右されず、既存の公開CLI入力契約を増やさない。

**欠点:** 利用者が任意のStoreを選ぶ機能は初期提供に含まれず、保存先はOSごとに異なる。

### 2. 実行ディレクトリの固定相対パスを使う

例として `./knowledge.db` をStoreとする。

**利点:** パスを説明しやすい。

**欠点:** 実行場所ごとにデータが分散し、同じ利用者の知識Storeとして安定しないため採用しない。

### 3. 保存先を公開optionまたは設定で指定する

`--store`、環境変数、設定ファイル等でStoreを指定する。

**利点:** 利用者が保存先を選べる。

**欠点:** 既承認のCLI・設定契約を拡張し、運用・共有・バックアップ方針も追加で必要になるため、初期提供では採用しない。

## 決定

利用者の明示承認により、選択肢1を採用する。

- 通常ビルドの `cmd/knowledge` は、各操作の実行前に上記の既定Storeをcomposition rootで解決する。
- 読み取り操作は、Storeの初回作成・migration以後、Assertion、Evidence、Concept、Scope、Relation、Temporal Metadata、派生Indexを変更しない。
- 結合テストは公開CLIへ保存先指定を追加せず、OS標準のユーザー設定ディレクトリ環境だけをテストごとの一時領域へ隔離して、通常のcompositionを実行する。

## 影響範囲

- [詳細設計](../design.md) のScope、CLI境界、Security/NFR、Migration、Open Issues
- [実装タスク](../tasks.md) のTASK-001-01、TASK-001-02、TASK-001-03、TASK-001-08およびFeature全体の対象外
- [implementation handoff](../implementation-handoff.yaml)
- Initial Designのarchitectureおよびcross-cutting concerns（既定保存先を横断責務として明記する場合）
- Issue #179 の実装前提および受入条件

## 後続ゲート

- このDecisionを反映した詳細設計・タスク・handoff・Initial Designの独立レビュー。
- 更新後のFEAT-001詳細設計に対する利用者レビュー。
