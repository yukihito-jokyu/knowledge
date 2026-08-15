# DEC-FEAT-014: 複数CLI操作の失敗時に履歴を補償削除しない

> **状態:** `decided`  
> **Decision Level:** L2 — Feature Scoped  
> **対象Feature:** FEAT-004

## 決定

Knowledge Updateが一候補に対して複数の既存CLI操作を必要とする場合、先行操作の成功後に後段操作が失敗または中断しても、成功済みのAssertion、revision、Evidence、Relationを補償削除しない。

対象となる操作列は、訂正時の`revise`→`attach-evidence`と、置換時の`create`→`supersede`である。後段失敗時は、先行操作が成功していれば`partially_applied`を返し、成功IDと失敗operationを残す。`supersede`の`conflict`は既存Relationを読むoperationがないため、既適用と推測せず`failed`かつ`outcome_unknown`とする。

同じ`episode_id`の自動再実行・再開は初期提供で行わない。Update Resultは呼出側へ返す一時的Markdown成果物であり、再開用ledgerを導入しない。

## 根拠

- FEAT-001はEvidenceと旧Assertionの物理削除を提供せず、BR-008とNFR-004は履歴保全を求める。
- `revise`、`create`、`supersede`はそれぞれ単独ではtransactionを持つが、二CLI操作を横断する原子性は既存契約にない。
- 補償削除、Relation読取、再開ledgerを追加すると、既存CLIまたは横断的な永続化契約を変更する必要がある。

## 影響

部分適用・結果不明を隠さず返すが、URL評価への回答は改変しない。利用者表示や後日の解決操作はUI-001または後続Featureで扱う。
