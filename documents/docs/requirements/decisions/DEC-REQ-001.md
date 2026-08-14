# DEC-REQ-001: Semantic Search の段階的提供

- **Status:** decided
- **Level:** L3
- **Decision:** 初期提供では Semantic Search を実装せず、SQLite の字句検索を提供する。Semantic Search は後続の移行で A（ローカル埋め込みモデル + ローカル Vector Index）として追加する。

## 理由

初期実装のモデル配布・ローカル実行・Index 管理の複雑さを分離し、Evidence 起点の保存・字句検索・履歴保持を先に提供するため。

## 維持する要件

- Semantic Search は最終システムの論理要件として維持する。
- 初期のデータ表現は、既存 Assertion を再作成せずに Embedding と Vector Index を後付けできる必要がある。
- `no_evidence` と未知、類似候補と既知判定を混同しない不変条件は維持する。

## 初期提供の制限

- `search-semantic` は利用可能な検索プリミティブに含めない。
- Codex は Semantic Search を必要としない字句・Concept・Relation・Evidence・時点差分の探索を使用する。
- Semantic Search が必要な評価は、後続移行の受入範囲とする。
