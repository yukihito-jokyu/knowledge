# FEAT-001 詳細設計レビュー（矛盾候補・既定Store・要求キャンセル）

- **対象:** `search-contradictions` の C-01対応、DEC-FEAT-005の通常ビルド既定Store、およびDEC-FEAT-006の要求キャンセル。公開 I/O、SQLite schema、process composition、結合テスト、上位要件・Decision・現行実装を独立に照合した。
- **総合判定:** **pass**
- **人間 Decision:** 不要。`DEC-FEAT-006` の公開契約は利用者承認済みであり、R-01／R-02の実装・検証上の未決定事項も公開契約を変えずに解消された。
- **handoff 可否:** 可。実装は固定済みのresponse開始境界と非公開同期gateに従う。

## 正規根拠台帳

| 根拠 | 確定状態 | 適用範囲 |
| --- | --- | --- |
| REQ-009、REQ-012、REQ-014、BR-010、CON-004 | confirmed | Relation と矛盾候補の論理検索、Codex と CLI の責務分離、JSON の機械境界 |
| DEC-FEAT-002、DEC-FEAT-003 | decided (L3) | 11操作、JSON response、名前付きoption入力 |
| DEC-FEAT-004 | decided (L3) | `seed` / 相手 `target` / 保存方向 `direction`、`contradicts` を Assertion→Assertion に限定 |
| DEC-FEAT-005 | decided (L3) | `os.UserConfigDir()/knowledge-cli/knowledge.db`、初回migration、公開Store overrideなし |
| DEC-FEAT-006 | decided (L3) | 最初の`Ctrl-C`で同一要求Contextを全層へ伝播し、JSONなし・stdout／stderr空・exit 130、mutationはcommitしない |

## 契約根拠トレーサビリティ

| 境界 | 設計上の契約 | 根拠判定 | 監査結果 |
| --- | --- | --- | --- |
| 検索結果 | `seed` は selector に一致した Assertion、`target` は常にその反対側 Assertion | explicit | `search-contradictions` の項目表・JSON例・UC-03で一致する。`target.id`を `get` / `get-evidence` へそのまま渡せる。 |
| 保存方向 | `direction` は `outgoing` / `incoming` | explicit | SQL-01の2分岐、項目表、用語、UC-03で一致する。相手の選択と保存方向を混同しない。 |
| Concept selector | 所属 Assertion を `seed` 候補へ展開し、両端が該当するRelationは2結果を返す | explicit | `selected_seed` と両方向の `UNION ALL` が、各endpointを別の `seed` として返す。relation IDだけでは統合しない。 |
| `contradicts` endpoint | Assertion→Assertion のみ | explicit | DEC-FEAT-004、`relations` table のCHECK、`create` の入力項目・validation、search SQLのいずれも一致する。 |
| 既定Store | OS標準のユーザー設定ディレクトリ配下の1箇所 | explicit | DEC-FEAT-005、詳細設計、architecture、handoffが `os.UserConfigDir()/knowledge-cli/knowledge.db` で一致する。 |
| process境界検証 | 通常compositionを、一時化したOS設定ディレクトリ環境で実行 | explicit / derived | 公開Store overrideを追加せず、初回migration・空Store・保存先初期化失敗を実バイナリで観測する。 |
| 要求キャンセル | response先頭byte前に観測した最初の`Ctrl-C`は同一ContextをCLI入口からSQLiteへ渡し、通常errorへ写像せずexit 130で終える | explicit | DEC-FEAT-006、設計本文、architecture、cross-cutting concerns、tasks、handoffは、response開始前の優先順位、非Context API前後の確認、開始後に既出力を取り消さない境界を一致して記載する。 |
| 割込みのprocess検証 | `integrationtest`限定の非公開gateでmigration／read／mutation開始を確認してからSIGINTを送る | explicit | Decision、acceptance、architecture、cross-cutting concerns、TASK-001-08／09、handoffが、通常composition・既定Store・公開option／環境変数／設定を変更しない同一方式を定める。 |
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
- DEC-FEAT-006、`design.md`、CLI入力規約、command catalog、architecture、cross-cutting concerns、tasks、handoffは、最初の割込みの公開値（JSONなし、stdout／stderr空、exit 130）とContext伝播先を矛盾なく参照する。

## DEC-FEAT-006 の反証結果

### R-01: 取消しと通常responseの優先順位（解消）

DEC-FEAT-006:17-18はresponse先頭byte前を線形化点として固定し、結果・validation・Store初期化errorより130と無出力を優先する。`os.UserConfigDir`と親ディレクトリ作成の前後確認も明記された。response開始後の割込みは既出力を取り消さず、選択済みresponseを完了するため、実現不能な巻戻し契約を含まない。`design.md`、architecture、CLI入力規約、command catalog、cross-cutting concerns、tasks、handoffは同じ境界を参照する。

### R-02: process境界の割込み検証が決定論的に実施できない（解消）

DEC-FEAT-006:20は、`integrationtest` build tagだけに含める非公開gateでmigration、read、mutationの開始を親testが確認してから割込みを送る方式を固定した。TASK-001-08／09とAcceptanceは、stdout／stderr空・exit 130、Context伝播、mutationの非commitとDB事後状態の不変を検証対象とする。このgateは通常composition・既定Store解決・公開CLI option／環境変数／設定へ露出しないため、DEC-FEAT-005および公開I/O契約とも整合する。

## コマンド体系・不要契約の再確認

`search-contradictions` は `search-related --relation-type contradicts` と低層のRelation読取りを共有するが、REQ-012、CON-004、DEC-FEAT-002が矛盾候補を独立した論理検索として要求・承認している。削除または統合は根拠がない。

今回追加された `seed` と `direction` は、既決のDEC-FEAT-004により、`target`を常に後続取得可能な相手として固定しつつ保存方向を失わないために必要な最小情報である。DEC-FEAT-005の既定Storeは通常ビルドを実行可能にする最小の運用契約であり、保存先を選ぶ入力、設定、運用仕様、AIによる意味判断は追加していない。

## 総合判定と次のゲート

**判定:** `pass`

旧C-01は解消された。`contradicts` の Assertion→Assertion 制約は、Concept endpointを一般に否定する新規規約ではなく、明示承認済みDEC-FEAT-004に基づく矛盾候補検索の型安全な範囲限定である。DEC-FEAT-005の既定Storeは通常compositionだけで解決し、結合テストはOS設定ディレクトリ環境を隔離するため、公開CLIの保存先指定や利用者Storeへの接続を生じない。DEC-FEAT-006はresponse先頭byte前の取消し優先順位と、tag限定の非公開同期gateによる決定論的なprocess検証を固定し、通常compositionと公開契約を変えない。残る実装上の差分は、承認済み設計を実現する作業であり、仕様未決定ではない。

**次のゲート:** 実装で、同一Contextの全層伝播、非Context初期化API前後とresponse直前の取消し確認、commit前の取消し確認、`integrationtest`実バイナリによるmigration／read／mutationのSIGINT検証を完了する。通常ビルド・既定Store・公開CLI I/Oの回帰を横断確認する。
