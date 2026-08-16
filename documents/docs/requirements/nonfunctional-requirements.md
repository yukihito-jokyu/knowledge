# 非機能要件

### NFR-001: 判定の説明可能性

Knowledge Assessment と Reading Value Assessment は、結論だけでなく Evidence、Known、Knowledge Gap、矛盾・時点差分、推奨理由を追跡可能にする。

- **Source:** Issue #175 §17, §36, §43, §52
- **Status:** confirmed

### NFR-002: 検索の有限性

Agentic Search は強い Evidence、飽和、Evidence 増分の停止、強い矛盾または検索 Budget により停止可能でなければならない。

- **Source:** Issue #175 §16
- **Status:** confirmed

### NFR-003: 検索品質の検証可能性

検索・評価・更新の誤りを切り分けられる Search Trace を記録可能にし、層別の Fixture／評価で検証可能にする。

- **Source:** Issue #175 §43–§45
- **Status:** confirmed

### NFR-004: 履歴保全

Evidence と supersede 前の Knowledge を保持し、訂正後にも再評価できなければならない。

- **Source:** Issue #175 §13, §25, §44
- **Status:** confirmed

### NFR-005: 論理契約の互換性

物理 DB、Embedding Engine、Index Library、Migration 方式を後で変更しても、Codex と Knowledge CLI 間の論理契約を壊してはならない。

- **Source:** Issue #175 §48
- **Status:** confirmed

### NFR-006: 受入評価

空 Store、完全一致、構成知識、誤解、古い知識、質問のみ、AI説明のみ、訂正、一部のみ高価値、些末な未知の各シナリオを検証可能にする。

- **Source:** Issue #175 §44–§45
- **Status:** confirmed

### NFR-007: Store選択の互換性と隔離

Store指定がない既存CLI invocationは既定Storeの解決・初回migration・JSON/error/exit codeの契約を維持する。指定があるinvocationは、その一意なStoreだけを初期化・migration・読書きし、他のStoreへ接続しない。

- **Source:** Issue #233
- **Status:** confirmed
