# データ所有権

| データ | Owner | 正規ソース | 利用者 |
|---|---|---|---|
| Knowledge Assertion、Evidence、Concept、Scope、Relation、Temporal Metadata | Knowledge CLI | Knowledge Store | Codex の Knowledge Search / Update |
| Derived User Knowledge State | Knowledge CLI | Evidence から導出 | Codex の Knowledge Search |
| Search Trace | Orchestration（永続化する場合は Knowledge CLI） | 実行時の検索過程 | 品質保証・診断 |
| Article Claim | Article Analysis | URL から取得した記事とその解析結果 | Knowledge Search、Reading Value |
| Knowledge Assessment | Knowledge Search | Claim と Knowledge CLI の検索・Evidence結果 | Reading Value |
| Reading Value Assessment | Reading Value | Article Claim と Knowledge Assessment | 利用者 |
| Candidate Knowledge | Knowledge Acquisition | Conversation / Task Episode | Knowledge Update |
| 更新判断 | Knowledge Update | Candidate と既存 Knowledge の照合 | Knowledge CLI |

## 更新原則

- Knowledge Store への変更は Knowledge Update が決定し、Knowledge CLI が実行する。
- 表の **Owner** は、当該データを保存・更新する唯一の経路を表す。Derived User Knowledge State の意味的な導出・知識評価は Codex の Knowledge Search が担い、Knowledge CLI は指定された結果を保存する場合にも決定論的な実行だけを担う。CLI が評価規則を決めたり、Evidence から状態を自律判定したりしてはならない。
- 訂正・置換は過去 Evidence または旧 Assertion を削除せず、履歴を維持する。
- 記事の閲覧・評価結果・AIの説明は Knowledge Evidence の正規ソースにならない。
