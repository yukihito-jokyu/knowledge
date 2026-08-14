# FEAT-001 詳細設計レビュー（矛盾候補の検索起点・相手・保存方向）

- **対象:** `search-contradictions` の C-01 対応。公開 I/O、SQL、SQLite schema、Relation 作成時の検証、図、UC-03、および上位要件・Decision を独立に照合した。
- **総合判定:** **pass**
- **人間 Decision:** 不要。`DEC-FEAT-004` は利用者が明示承認した L3 Decision であり、今回の設計はその実装詳細への反映である。
- **handoff 可否:** 利用者による FEAT-001 詳細設計の明示承認前のため、不可。

## 正規根拠台帳

| 根拠 | 確定状態 | 適用範囲 |
| --- | --- | --- |
| REQ-009、REQ-012、REQ-014、BR-010、CON-004 | confirmed | Relation と矛盾候補の論理検索、Codex と CLI の責務分離、JSON の機械境界 |
| DEC-FEAT-002、DEC-FEAT-003 | decided (L3) | 11操作、JSON response、名前付きoption入力 |
| DEC-FEAT-004 | decided (L3) | `seed` / 相手 `target` / 保存方向 `direction`、`contradicts` を Assertion→Assertion に限定 |

## 契約根拠トレーサビリティ

| 境界 | 設計上の契約 | 根拠判定 | 監査結果 |
| --- | --- | --- | --- |
| 検索結果 | `seed` は selector に一致した Assertion、`target` は常にその反対側 Assertion | explicit | `search-contradictions` の項目表・JSON例・UC-03で一致する。`target.id`を `get` / `get-evidence` へそのまま渡せる。 |
| 保存方向 | `direction` は `outgoing` / `incoming` | explicit | SQL-01の2分岐、項目表、用語、UC-03で一致する。相手の選択と保存方向を混同しない。 |
| Concept selector | 所属 Assertion を `seed` 候補へ展開し、両端が該当するRelationは2結果を返す | explicit | `selected_seed` と両方向の `UNION ALL` が、各endpointを別の `seed` として返す。relation IDだけでは統合しない。 |
| `contradicts` endpoint | Assertion→Assertion のみ | explicit | DEC-FEAT-004、`relations` table のCHECK、`create` の入力項目・validation、search SQLのいずれも一致する。 |
| 既存Relation検索 | `search-related` は全Relation種別・Assertion/Concept endpointを扱う | explicit / derived | `contradicts` の型別制約はこの種別だけであり、他RelationのConcept endpoint契約を狭めていない。 |

## C-01 の反証確認

| 反証ケース | 確認結果 |
| --- | --- |
| `asrt_09 → asrt_01` を保存し、`--assertion-id asrt_01` を指定する | incoming分岐で `seed=asrt_01`、`target=asrt_09`、`direction=incoming` となる。旧C-01の「起点自身を返す」状態は生じない。 |
| `asrt_01 → asrt_09` を保存し、`--assertion-id asrt_01` を指定する | outgoing分岐で `seed=asrt_01`、`target=asrt_09`、`direction=outgoing` となる。 |
| 同一Conceptに `asrt_01` と `asrt_09` が属し、その間にRelationがある | 両Assertionが`selected_seed`へ入り、outgoing/incomingの各1行を返す。`seed`ごとに相手が明確で、2行を誤って重複排除しない。 |
| `contradicts` の片端がConceptである既存データ | v1 schema CHECK と `create` validation が保存を拒否するため、検索応答でAssertionとして扱うデータ喪失・型不一致は生じない。 |

## 資料間整合

- `search-contradictions` の I/O、項目定義、validation、フローチャート、シーケンス図、SQL-01は、`seed`、相手`target`、`direction`の意味を一致して記載する。
- `database-schema.md` は `contradicts` と `supersedes` だけを Assertion→Assertion にCHECKし、`create` は新規Assertionをsourceにしつつ、`contradicts` のtarget kindをAssertionだけに制限する。両者は整合する。
- UC-03は `results[0].target.id` を次の `get` / `get-evidence` の `--assertion-id` として具体的に受け渡す。これは検索結果契約と一致する。
- `design.md`、`database-access-map.md`、DEC-FEAT-004に、検索対象・保存方向・影響範囲について矛盾する記述は検出しなかった。

## コマンド体系・不要契約の再確認

`search-contradictions` は `search-related --relation-type contradicts` と低層のRelation読取りを共有するが、REQ-012、CON-004、DEC-FEAT-002が矛盾候補を独立した論理検索として要求・承認している。削除または統合は根拠がない。

今回追加された `seed` と `direction` は、既決のDEC-FEAT-004により、`target`を常に後続取得可能な相手として固定しつつ保存方向を失わないために必要な最小情報である。不要な入力、保存先、設定、運用仕様、AIによる意味判断は追加されていない。

## 総合判定と次のゲート

**判定:** `pass`

旧C-01は解消された。`contradicts` の Assertion→Assertion 制約は、Concept endpointを一般に否定する新規規約ではなく、明示承認済みDEC-FEAT-004に基づく矛盾候補検索の型安全な範囲限定である。

**次のゲート:** 利用者による FEAT-001 詳細設計の明示レビュー・承認。承認後にのみ task breakdown と implementation handoff を作成できる。
