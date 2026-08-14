# 技術制約と保留事項

## 採用済み

- AI Runtime は Codex とする。
- 永続知識インターフェースは自作 Knowledge CLI とする。
- Codex と CLI の機械的境界は JSON とする。
- Codex のサブワークフロー間成果物は Markdown を基本とする。
- CLI は決定論的な保存・検索・取得・更新を担い、AI Agent を内包しない。
- Knowledge Store はローカル専用の埋め込み SQLite とし、複数端末同期・共同利用・リモート共有を初期提供に含めない。
- Knowledge CLI 実装言語は Go とし、単一のGo moduleで依存を管理する（[DEC-ARCH-002](decisions/DEC-ARCH-002.md) の人間承認）。
- SQLite driverは純Goの `modernc.org/sqlite` とし、`go.mod` の最低対応Go versionは1.26とする（[DEC-ARCH-002](decisions/DEC-ARCH-002.md)）。CGOおよびC compilerをbuild前提にしない。
- 実行入口、非公開のapplication/domain/persistence、SQLite migration asset、fixture／integration testを分離する。具体的な責務と依存方向は [architecture.md](architecture.md#go-module-と配置規約) に従う。
- migration SQL assetは実行バイナリに同梱し、Go標準toolchainの `gofmt`、`go test ./...`、`go vet ./...` を共通のformat・test・静的検査とする。

## 必須の論理能力

- Assertion 中心、Evidence 起点、履歴保持。
- Lexical、Semantic、Relation、Temporal の探索。初期提供では DEC-REQ-001 により Semantic Search を延期する。
- Concept、Evidence、矛盾候補の検索・取得。
- 物理実装を変えても論理契約を維持すること。

## Deferred

- Embedding Engine、Vector Index は Semantic Search を追加するFEAT-006で扱う。初期は字句検索のみであり、SQLite driverの選定はEmbedding／Vector Indexの選定を意味しない。
- URL 評価の利用者インターフェース。
- Conversation / Task Episode の入力境界。

FEAT-001のmigration論理契約は詳細設計で決定済みであり、assetはGo標準の`embed`で同梱する。Embedding Engine と Vector Index は、Semantic Search を追加する後続移行で A（ローカル実行）として決定済みである。URL 評価の利用者インターフェースとエピソード入力境界は、各 Feature の詳細設計まで保留する。
