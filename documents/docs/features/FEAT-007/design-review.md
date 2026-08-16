# FEAT-007 独立設計レビュー

- **対象:** FEAT-007「隔離Knowledge Storeの明示選択」
- **レビュー日:** 2026-08-16
- **レビュー方法:** 設計者と独立した読み取り専用照合
- **総合判定:** `pass`

## 正規ソース台帳

| ソース | 確定状態 | 適用範囲 |
| --- | --- | --- |
| [Issue #233](https://github.com/yukihito-jokyu/knowledge/issues/233) | confirmed | sandboxでの明示Store、既定Store互換、全operation一貫性、検証環境、実装前承認 |
| `requirements/requirements.md` REQ-022 | confirmed | invocation単位の絶対パス指定と指定なし既定Store |
| `requirements/nonfunctional-requirements.md` NFR-007 | confirmed | 指定Storeの隔離と既存契約の互換 |
| `requirements/constraints.md` CON-002、CON-003、CON-007 | confirmed | JSON CLI境界と、option以外の暗黙解決禁止 |
| [DEC-FEAT-023](decisions/DEC-FEAT-023.md) | L3 decided | `--store`の構文、検証、error、fallback禁止、DEC-FEAT-005の限定置換 |
| `design/architecture.md`、`cross-cutting-concerns.md` | confirmed | composition、migration、process境界検証、Context・設定境界 |
| `features/FEAT-001/decisions/DEC-FEAT-005.md` | decided（限定置換対象） | 指定なし既定Storeと初回初期化を継続 |

## 契約根拠トレーサビリティ

| 境界 | FEAT-007設計上の契約 | 根拠 | 判定 |
| --- | --- | --- | --- |
| 公開CLI入力 | `knowledge [--store <absolute-path>] <operation> [operation options]`。optionはoperation前に最大一回 | Issue #233、REQ-022、CON-007、DEC-FEAT-023 | explicit |
| 入力失敗 | 値なし・重複・相対パス・不正位置・未知operationは、Store open前に`validation_error`／stderr JSON／exit 2 | FEAT-007要件受入条件4、DEC-FEAT-023 | explicit |
| 永続化・migration | 有効な指定先だけを作成・open・migrationする | Issue #233、REQ-022、NFR-007、DEC-FEAT-023 | explicit |
| open失敗 | 親作成・SQLite open・migration失敗はfallbackなしの`storage_error`／stderr JSON／exit 1 | FEAT-007要件受入条件5〜6、DEC-FEAT-023 | explicit |
| 既定互換 | optionなし時のみ`os.UserConfigDir()/knowledge/knowledge.db`を使う | Issue #233、REQ-022、NFR-007、DEC-FEAT-005、DEC-FEAT-023 | explicit |
| 既存operation | 11operationは選択済みの同一Storeだけを使い、success JSONを変えない | Issue #233、FEAT-007要件受入条件2、NFR-007 | explicit |
| workflow伝播 | 検証用Skillが固定配置から`verification/knowledge.db`を絶対パスとして一意に解決し、全CLI invocationへargvで渡す | Issue #233、FEAT-007要件受入条件7、design.md「検証用workflowのStore伝播入口」 | derived |
| 設定・権限 | environment/config/CWDの暗黙解決、権限昇格、他Store探索をしない | Issue #233、CON-007、DEC-FEAT-023 | explicit |
| 検証環境の配置 | 現行検証環境は`chore/verification-environment`に隔離し、Feature実装と混在させない | Issue #233、design.md「Implementation Deliverables」 | explicit |

## 資料間整合

- `cmd/knowledge`の現行実装は、操作parse後に`openDefaultRetrievalStore`で`os.UserConfigDir()`配下を開き、作成・検索・履歴更新で同じ既定Store解決を使う。FEAT-007はこの選択境界だけを一般化し、domain I/O、schema、migration asset、operation固有I/Oを変更しないため、責務分離と整合する。
- 現行CLIは`validation_error`をexit 2、`storage_error`をexit 1としてstderr JSONへ出す。設計の失敗分類・stdout/stderr・exit codeは既存公開error境界と整合する。割込みは既存のexit 130・無出力規約を優先するため、通常errorの追加とも衝突しない。
- FEAT-001資料には、DEC-FEAT-023以前の「保存先optionを追加しない」という記述が残る。設計はその相反箇所をScope/CLI入力/Store/Security/Assumptions、CLI入力規約、DEC-FEAT-005、tasks、handoff、Skill operation referencesまで列挙し、同一変更で新旧資料を並存させないことを固定している。DEC-FEAT-023が置換する範囲を保存先option禁止だけに限定し、既定Store・初回初期化を維持するため、未解消の契約矛盾ではない。
- `reading-value`の入力を一件の絶対HTTP(S) URLだけに保つ一方、検証用コピーだけがSkill自身の固定配置を入口にStoreを決定し、`knowledge-search`と`knowledge-update`へ必須実行コンテキストとして渡す。したがって、利用者入力にStore設定を追加せず、各CLI呼出しで同一の明示argvを使える。通常の`skills/`がoptionを自動追加しないことも明記され、CON-007と整合する。
- unit、process integration、全11operation regression、URLだけのuser-run workflow verificationが、それぞれparse、Store選択、migration、I/O境界、Skill伝播を観測する。Issue #233の検証要件を網羅する。

## 前回差し戻し事項の再照合

| Finding | 結果 | 根拠 |
| --- | --- | --- |
| DR-007-001: URLだけのSkill入力から同一Storeを渡す入口が未確定 | resolved | 固定`verification/`配置からの正規絶対パス解決、`verification/knowledge.db`、全CLI invocationの明示argv、失敗時のfallback禁止が設計された。 |
| DR-007-002: FEAT-001資料との整合更新範囲が未確定 | resolved | 資料別の置換・維持内容と、同一変更で旧資料を残さない責務が設計された。 |

## 人間Decisionの境界

公開CLIと保存先のL3判断は、DEC-FEAT-023に利用者承認済みとして記録されている。検証用Skillの固定配置から同一argvを導く部分は、その承認済み公開契約を実装可能なworkflow境界へ落とす詳細であり、新しい公開設定・保存先・権限契約を追加しない。追加のL3/L4 Decisionは不要である。

## 次のゲート

`pass`。次は**利用者によるFEAT-007詳細設計の明示承認**である。承認前にtask breakdown、implementation handoff、実装へ進んではならない。
