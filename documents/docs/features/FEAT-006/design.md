# FEAT-006 詳細設計: ローカル Semantic Search 移行

> **状態:** blocked — [DEC-FEAT-023](decisions/DEC-FEAT-023.md) の L3 判断待ち。ローカル実行は既決だが、モデルと Index の提供方式が未確定のため、SQLite schema と `search-semantic` の物理／wire 契約は未完成である。

## Feature Summary

正規 Assertion を再作成せず、現行の正規化本文を入力にローカル Embedding を生成し、ローカル Vector Index で意味的候補を取得する。Semantic Search は候補検索だけを決定論的に実行し、候補からユーザーが知っているか、Evidence が十分か、意味的に同一かを判断しない。

## Scope / Out of Scope

含む要件、受入結果、対象外は [requirements.md](requirements.md) を正とする。既存11操作とその公開 JSON、既定 Store、Evidence・Relation・履歴・字句検索は変更しない。

## Related Requirements / Business Rules

REQ-011、BR-002、BR-004、BR-010、CON-002〜006、NFR-004〜006、DEC-REQ-001、DEC-FEAT-001、および [DEC-FEAT-023](decisions/DEC-FEAT-023.md) を根拠とする。

## Behavioral Scenarios

| Trigger | Preconditions | Main flow | Result |
| --- | --- | --- | --- |
| Codex が意味検索を要求する | Store の migration と Semantic Index が利用可能 | query をローカルで Embedding 化し、Vector Index から候補 Assertion ID を得て、現行 Assertion を取得する | 類似候補と検索根拠を返す。候補なしは成功の空配列。 |
| 既存 Store を初めて Semantic Search 対応へ開く | 正規 Assertion と既存 migration が整合している | migration が現行 Assertion から派生 Embedding / Index を構築する | 正規データを変更せず、検索可能な Index を得る。 |
| Assertion を create / revise する | 既存 mutation が入力を検証済み | 正規データの更新と派生 Embedding / Index 更新を同一の整合単位で完了する | 成功時だけ新しい現行 Assertion が Semantic Search の対象になる。 |
| migration・Embedding・Index 更新が失敗または要求が中断される | response 開始前 | Index 更新を確定せず既存の正規データ・既存派生データを保つ | `storage_error` または既存の中断規約で終了し、部分更新を残さない。 |

## Selected Design Analyses

| 分析 | 判定 | 理由 |
| --- | --- | --- |
| シナリオ、正常・異常 flow | 採用 | 検索、migration、mutation と失敗復帰を定めるため。 |
| 責務・interface・データ・永続化・error | 採用 | CLI、Embedding runtime、Vector Index、SQLite の境界を跨ぐため。 |
| transaction / idempotency | 採用 | 正規データと派生 Index の整合、再構築、中断を扱うため。 |
| フローチャート、シーケンス図 | 採用予定 | 実行順・rollback・失敗の戻り先が重要。ただし Index 基盤未決のため契約確定後に operation 資料へ記載する。 |
| 状態遷移図 | 不採用 | 長期業務状態は増やさず、migration version と派生 Index の整合は永続化契約で表現するため。 |
| UI state / security認可 | 不採用 | CLI とローカル Store だけが対象で、新しい UI・認可境界はない。 |

## Responsibilities

| 論理責務 | 担当 |
| --- | --- |
| Query の生成、候補への後続検索、semantic equivalence・Evidence・知識状態の判断 | Codex / Knowledge Search |
| query / 現行 Assertion の Embedding 生成、Vector Index の照会・更新・再構築、候補の決定論的な返却 | Knowledge CLI |
| Assertion、Evidence、revision、Relation の正規保管 | Knowledge Store（既存） |
| Embedding と Vector Index の派生保管 | Semantic Search の派生 Store。正規データではなく、現行 Assertion から再構築可能であることが必須。 |

## Implementation Deliverables / Placement

| root-relative 配置先 | 実行成果物の責務 | 依存 | 変更しない領域 |
| --- | --- | --- | --- |
| `internal/domain/` | Semantic Search の論理入力・候補結果・永続化 port | 既存 retrieval の意味中立な domain 境界 | Codex の意味判断を実装しない。 |
| `internal/application/` | Semantic Search の application operation | domain port | 検索戦略・知識状態判断を実装しない。 |
| `internal/persistence/sqlite/` | migration、正規 Assertion を入力とする派生 Index の再構築・照会・mutation同期 | SQLite Store、選定済み local embedding / index capability | Assertion/Evidence/Relation/revision の正規表を再作成・削除しない。 |
| `internal/persistence/sqlite/migrations/` | 既存連番の次の embedded schema migration asset | 既存 migration 規約 | 適用済み migration を編集しない。 |
| `cmd/knowledge/` | `search-semantic` の名前付き option 入力、既存 JSON envelope / error / exit 境界への接続 | application | 既存11操作の I/O と既定 Store 契約を変えない。 |
| `internal/**` と `cmd/knowledge/**` の既存 test 配置 | unit・SQLite integration・CLI process boundary の検証 | 既存 Fixture と migration test 規約 | 通常利用の Codex workflow Skill を test asset として変更しない。 |

これらは既存の [Architecture](../../design/architecture.md#go-module-と配置規約)、FEAT-001 実装、親 `AGENTS.md` に基づく配置である。具体的なモデル runtime / Vector Index dependency は DEC-FEAT-023 解消後にこの Feature 内で定める。

## Interfaces / Data / Physical / Wire Schema

`search-semantic` は REQ-011 と CON-004 により必要な新規 CLI operation である。入力 option、正常 response、error response、候補順序、類似度の表現、件数上限、model / Index version の露出範囲は、選択した model と Index の互換性に依存する。

SQLite の派生 table / index 名、column、型、NULL性、key・constraint、migration version・前後条件、再構築方式、および operation 別 JSON schema は未確定であり、実装へ委任しない。DEC-FEAT-023 解消後に `design/command-catalog.md`、`design/commands/search-semantic.md`、database schema reference、relationship map、access map を作成して complete にする。

## Error / Edge Cases

- query が意味的候補を返さないことは成功の空配列であり、CLI は `no_evidence` や `unknown` を返さない。
- モデル未配備・モデル互換性不一致・Index 破損・Embedding / Index 構築失敗は、正規データを変更せず技術 error として扱う必要がある。error code と復旧可能性は未決 Decision 解消後に確定する。
- response開始前の最初の `Ctrl-C` は、既存の無出力・終了コード130・commitしない規約を維持する。
- revision により現行本文が変わった Assertion は、旧 revision の Embedding を候補に残さない。ただし旧 revision、Evidence、Relation 自体は削除しない。

## Security / NFR Considerations

外部への Assertion / Scope / Evidence 送信を行わず、ローカル専用の既決データ境界を維持する。モデル配布・更新のネットワーク要件、ライセンス、容量、端末要件は DEC-FEAT-023 の選択により確定する。

## Acceptance / Test Design

Decision 解消後、少なくとも以下を具体的な Fixture・操作・期待 JSON / DB状態へ落とす。

- 空 Store、既存 Store upgrade、意味的に近い候補、候補なし、字句検索と異なる表現を検証する。
- create / revise 後に現行 Assertion だけが検索対象となり、Evidence・Relation・旧 revision が保たれることを検証する。
- migration / Embedding / Index 更新の失敗・中断時に、正規表と既存派生 Index が部分更新されないことを SQLite と実プロセス境界で検証する。
- 既存11操作の JSON、stdout/stderr、終了コード、字句検索結果が回帰しないことを検証する。
- Semantic Search の結果だけでは知識状態を判定しないことを、CLI response に state field を含めないことと Codex側利用契約で検証する。

## Contract Completeness

| 項目 | 状態 | 根拠・残作業 |
| --- | --- | --- |
| 永続データ・履歴 | partial | 正規データ不変と再構築原則は確定。派生 schema は未決。 |
| CLI / JSON wire schema | partial | 新規 `search-semantic` は必要。option・response は未決。 |
| 検索規則 | partial | 候補検索のみ、空結果成功、知識状態を判定しない点は確定。順位・上限・score表現は未決。 |
| 更新・整合・原子性 | partial | 正規データを破壊せず、現行 Assertion と派生 Index の整合を保つ原則は確定。実現する transaction 境界は未決。 |
| Schema / Store / Index | partial | embedded migration・連番・再実行・rollback の既存規約を継承。派生 schema と model / Index compatibility は未決。 |
| Fixture / test | partial | 必須観点を定めた。具体 Fixture と oracle は未決。 |
| Operation Documentation Coverage | pending | JSON CLI と永続 Store を持つため該当。DEC-FEAT-023 解消後に operation資料と DB資料を作成する。 |

## Assumptions

- 選定される Embedding runtime と Vector Index は、外部への個人知識送信なしにローカルで実行できる。
- 既存 `assertion_id` と現行 `normalized_text` が再構築の入力として利用できる。この前提は FEAT-001 schema と現行実装で確認済みである。

## Decisions

- [DEC-FEAT-023](decisions/DEC-FEAT-023.md): ローカルモデルと Vector Index の提供方式（L3、proposed）。

## Open Issues

- DEC-FEAT-023 の人間判断待ち。解消まで `design_ready` ではない。
