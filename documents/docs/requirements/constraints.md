# 制約

### CON-001: AI Runtime

Codex を AI Runtime とし、AI Agent は Codex 側にのみ置く。

- **Source:** Issue #175 ヘッダー, §7, §40

### CON-002: 永続化インターフェース

Persistent Knowledge Interface は自作 Knowledge CLI とし、CLI は AI 判断を内包しない。

- **Source:** Issue #175 ヘッダー, §6.8, §8.7, §40

### CON-003: データ境界

Knowledge CLI と Codex の機械的境界は JSON、Codex の各サブワークフロー間の成果物は必須セクションを備えた Markdown を基本とする。

- **Source:** Issue #175 §9, §17, §36, §47

### CON-004: 論理的検索プリミティブ

Knowledge CLI は字句、意味、Concept、Relation、Evidence、矛盾候補、時点差分の検索・取得を論理的に分離して提供する。初期提供では DEC-REQ-001 により Semantic Search を延期するが、後続追加を妨げない論理契約とデータ表現を維持する。

- **Source:** Issue #175 §11–§13

### CON-005: 検索起点

保存形式・物理実装は、Article Claim から関連 Knowledge に到達する検索要件を先に満たすよう設計する。

- **Source:** Issue #175 §6.7, §48, §49

### CON-006: 未確定の物理実装

物理 DB、Embedding Engine、Index Library、Migration 方式、具体的な検索パラメータは本要件では確定しない。

- **Source:** Issue #175 §16, §48

### CON-007: 明示的なStore選択

Storeの選択は、各CLI invocationの公開optionで絶対SQLiteファイルパスを明示して行う。環境変数・設定ファイル・実行ディレクトリからの暗黙解決は追加しない。

- **Source:** Issue #233
