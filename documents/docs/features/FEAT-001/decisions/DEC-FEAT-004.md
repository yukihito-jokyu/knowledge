# DEC-FEAT-004: 矛盾候補検索の検索起点・相手・保存方向

- **Status:** decided
- **Level:** L3
- **論点:** `search-contradictions` が保存済み Relation の source / target をそのまま返すと、検索起点が target 側の場合に、後続の `get` / `get-evidence` へ検索起点自身を渡してしまう。

## 確定済みの根拠

- Knowledge CLI は、矛盾候補を含む検索・取得を論理的に分離して提供する（REQ-012、CON-004）。
- Codex は検索結果を用いて Evidence を評価し、CLI は意味判断をしない（REQ-004、BR-010）。
- Relation は Assertion または Concept endpoint を持つが、`search-contradictions` の利用目的は矛盾候補 Assertion とそのEvidenceを確認することである（database-schema、UC-03）。
- `search-related` は、検索起点から見た相手を `target`、保存方向を `direction` として返す既存 v1 原則である。

## 決定

利用者の明示的な承認により、`search-contradictions` に同じ利用原則を適用する。

- 結果ごとに、selector に一致した Assertion を `seed`、その反対側の Assertion を `target` として返す。
- `target` は常に `knowledge get --assertion-id` および `knowledge get-evidence --assertion-id` に渡せる Assertion ID とする。
- 保存時の Relation 向きは失わず、`direction` に `outgoing` / `incoming` として返す。
- Concept selector は所属 Assertion を検索起点候補へ展開する。同じ Relation の両端が一致する場合は、`seed` ごとに別結果を返し、relation_idだけで統合しない。
- この用途の endpoint を Assertion に閉じるため、`contradicts` Relation は Assertion→Assertion に限定する。

## 採用理由

利用側は常に `target.id` を読めば相手 Assertion を次の取得操作へ渡せ、保存方向の比較処理を繰り返さない。`direction` を残すため、保存上の向きも確認できる。Concept selectorで両端が一致する場合に一方を除外すると、確認可能な経路を隠すため、`seed` を含めて別結果にする。

## 影響範囲

- `search-contradictions` の成功 JSON、SQL、ユースケース例。
- `contradicts` Relation の入力検証とSQLite schema v1制約。
- `search-related`、他のRelation種別、保存先・設定・運用仕様には影響しない。
