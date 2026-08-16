# FEAT-007 要件: 隔離Knowledge Storeの明示選択

## 目的

Codex workspace sandboxなど、OSユーザー設定領域へ書込みできない環境でも、利用者が書込み可能なローカルSQLite Storeを明示してKnowledge CLIとworkflowを実行できるようにする。

## 対象要件

- **REQ-022:** invocationごとの明示Store選択。
- **NFR-007:** 指定Storeの隔離と、指定なし既定Storeの互換性。
- **CON-002 / CON-003:** Knowledge CLIの決定論的JSON境界を維持する。
- **CON-007:** 絶対パスの公開optionだけでStoreを選び、暗黙設定を追加しない。
- **BR-010:** Codexは利用者のStore指定をCLIへ渡すだけで、CLIは保存・検索・更新を決定論的に実行する。

## 受入条件

1. `knowledge --store <absolute-path> <operation> ...` は、指定Storeの親ディレクトリ・DB・未適用migrationを必要時に作成し、そのoperationを実行する。
2. 11個の既存Store利用operationはすべて、同一invocationで指定されたStoreだけを使う。
3. `--store`を省略した既存invocationは、`os.UserConfigDir()/knowledge/knowledge.db` を既定Storeとして使い続ける。
4. `--store`の欠落値、重複、相対パス、未知operationとの組合せは、実行前に`validation_error`、stderrのJSON、終了コード2となる。
5. 絶対パスが有効でも、親ディレクトリ作成、DB open、migrationが失敗した場合は`storage_error`、stderrのJSON、終了コード1となる。
6. 指定Storeの初期化またはmigration失敗はtransactionの既存規則に従い、別のStoreへのfallbackを行わない。
7. 利用者が`verification/`をworkspaceとして開くCodex利用手順は、そのworkspace物理絶対パスから`knowledge.db`を導出して`--store`で渡し、4つの検証用Skillが同じbinaryとStoreを利用する方法を示す。

## 対象外

- 既定Storeを`~/.knowledge`、workspace、または別のOS設定領域へ移動すること。
- 環境変数、設定ファイル、または任意の実行ディレクトリからのStore自動選択。利用者が開いた固定`verification/` workspaceの物理パスを検証用Skillが利用する場合だけは除く。
- Storeの同期、共有、暗号化、sandbox権限の昇格または回避。
- SQLite schema、既存11operationのdomain data、JSON成功responseの変更。
