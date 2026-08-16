# FEAT-006 要件再構成: ローカル Semantic Search 移行

## 目的

既存の Knowledge Store に保存された現行 Assertion を正規ソースとして、ローカルで Embedding を生成し、ローカル Vector Index から意味的に近い候補を検索できるようにする。既存の Evidence、Relation、revision 履歴および字句検索を破壊・意味変更しない。

## 含む要件

| ID | 要求 |
| --- | --- |
| REQ-011 | 正規化済み Assertion 中心の意味検索を提供し、意味検索結果だけでユーザーの知識状態を判定しない。 |
| BR-002 | 候補がないことをユーザーが知らないことと同一視しない。 |
| BR-004 | Evidence を知識状態の正規根拠として維持する。 |
| BR-010 | Codex は検索戦略と意味判断、CLI は決定論的な検索・保存処理を担う。 |
| DEC-REQ-001 / DEC-FEAT-001 | Embedding と Vector Index はローカルで実行し、既存 Assertion ID と正規データから再構築する。 |

## 受入結果

- 既存 Assertion の安定 ID と現行正規化本文から Embedding と Vector Index を構築・再構築できる。
- `search-semantic` は意味的候補と検索に用いた Index / model の識別情報を決定論的に返す。
- Semantic Search の候補、順位、類似度は Evidence や Derived User Knowledge State ではなく、Codex が後続の検索・評価を決めるための候補である。
- 既存の11操作、JSON envelope、stdout/stderr、終了コード、既定 Store、履歴および字句検索の契約を維持する。
- 失敗・中断した migration / Index 構築では、正規 Assertion、Evidence、Relation、履歴、既存字句 Index を変更しない。

## 対象外

- Codex による query 生成、検索継続・停止、semantic equivalence、ユーザー知識状態・Evidence 強度の判断。
- リモート Embedding API、個人知識の外部送信、同期・共有。
- 既存 Assertion / Evidence / Relation / revision の物理削除・再作成。
- 既存11操作の入力・出力互換性を壊す変更。

## 未解決の前提

ローカル Embedding モデルの配布・更新方法と Vector Index の実装基盤は未決である。この選択はモデル識別子・ベクトル次元・Index の永続化 schema・migration・初回構築時の失敗契約を決めるため、[DEC-FEAT-023](decisions/DEC-FEAT-023.md) の人間判断を要する。
