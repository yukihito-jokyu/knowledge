# FEAT-001 要件

## 目的

Evidence を正規根拠とする個人知識をローカル SQLite に保存し、Codex が JSON 経由で決定論的に検索・取得・更新できる Knowledge CLI を提供する。

## 含む要件

- REQ-009〜REQ-014
- BR-002、BR-004、BR-007、BR-008、BR-010
- NFR-001、NFR-004、NFR-005、NFR-006
- CON-001〜CON-005
- DEC-ARCH-001

## Feature Acceptance

1. Assertion、Evidence、Concept、Scope、Relation、Temporal Metadata をそれぞれの意味で保存・取得できる。
2. Evidence が存在しないことをユーザーが知らないことへ変換しない情報を維持できる。
3. 初期提供では、字句・Concept・Relation・Evidence・矛盾候補・時点差分の論理検索を JSON で呼び出せる。Semantic Search は DEC-REQ-001 および DEC-FEAT-001 に従い後続移行で提供する。
4. create、attach-evidence、revise、supersede は履歴を失わず、SQLite の一貫した更新として実行される。
5. CLI は AI による意味判断や知識状態判定を実行しない。

## 範囲外

- 検索 Query の生成、検索の継続判断、`known` 等の Assessment 判定。
- 記事の Claim 分解、読書価値判断、会話からの Evidence 候補抽出。
- 複数端末同期、共同利用、リモート共有。
- Semantic Search の実装（DEC-REQ-001 に従う後続移行で扱う）。
