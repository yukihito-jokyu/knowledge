# Feature Dependencies

| From | To | 種別 | 理由 |
|---|---|---|---|
| FEAT-002 | FEAT-001 | data prerequisite | Assessment には Assertion、Evidence、Concept、Relation、Temporal 情報の検索・取得が必要。 |
| FEAT-003 | FEAT-002 | behavioral prerequisite | 記事の読書価値評価は Claim ごとの Knowledge Assessment を統合する。 |
| FEAT-004 | FEAT-001 | data prerequisite | Evidence の追加、訂正、履歴保持には知識基盤が必要。 |
| FEAT-006 | FEAT-001 | data prerequisite | Semantic Search 移行は、安定した Assertion 識別子と既存の正規データから Index を構築する。 |
| FEAT-005 | FEAT-001〜004 | optional sequencing | 各 Feature の完成に伴い評価 Fixture を追加できる。全機能完成後にのみ開始する必須依存ではない。 |

## 推奨順序

1. FEAT-001: Evidence 起点の個人知識基盤
2. FEAT-002: Claim ごとの Agentic Knowledge Assessment
3. FEAT-003: URL 記事の読書価値評価
4. FEAT-004: 会話・作業からの知識 Evidence 更新
5. FEAT-005: 知識評価の品質保証（各 Feature と並行して拡充）
6. FEAT-006: ローカル Semantic Search 移行

FEAT-003、FEAT-004、FEAT-006 の優先順位は、どの利用シナリオと検索能力を先に提供するかに関わるため、ここでは推奨であり確定ではない。
