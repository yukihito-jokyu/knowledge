# Backlog

| Feature | Status | Priority | Dependency | Planning Notes |
|---|---|---|---|---|
| FEAT-001 Evidence 起点の個人知識基盤 | planned | proposed: highest | なし | 初期は SQLite と字句検索を提供し、Semantic Search は DEC-REQ-001 に従う後続移行として設計境界を保持する。 |
| FEAT-002 Claim ごとの Agentic Knowledge Assessment | planned | proposed: high | FEAT-001 | 検索 Budget・Search Trace・7状態の Assessment 契約を設計対象とする。 |
| FEAT-003 URL 記事の読書価値評価 | planned | proposed: high | FEAT-002 | Primary Use Case。UI-001 の提供形態は未決定。 |
| FEAT-004 会話・作業からの知識 Evidence 更新 | planned | proposed: high | FEAT-001 | UI-002 の取り込み境界は未決定。 |
| FEAT-005 知識評価の品質保証 | planned | proposed: high | FEAT-001〜004 | Scenario A〜J を層別評価へ対応付ける。 |
| FEAT-006 ローカル Semantic Search 移行 | planned（初期提供後） | proposed: high | FEAT-001 | DEC-REQ-001 に従い、ローカル Embedding と Vector Index を追加する。既存正規データから Index を再構築する。 |
| FEAT-007 隔離Knowledge Storeの明示選択 | planned | proposed: high | FEAT-001 | Issue #233。Codex sandbox向けに明示Storeを選択可能にし、既定Storeは変更しない。 |

## 優先順位に関する注意

Priority は依存関係と Issue #175 の Primary Use Case に基づく提案であり、MVP／リリース範囲の決定ではない。MVP で FEAT-004 を含めるかは人間の優先順位判断が必要である。
