# FEAT-005 詳細設計レビュー

## 対象と判定

- **対象:** `requirements.md`、`design.md`、`design/scenario-catalog.md`、DEC-FEAT-018〜020
- **総合判定:** **pass**
- **次のゲート:** 利用者による詳細設計レビュー。Task分割・implementation handoffは、この承認後にのみ進められる。

今回の独立レビューでは、Reading Valueを固定Fixtureから実行せず、A/G/I/Jでは既存FEAT-003検証契約への参照完全性だけを確認するというDEC-FEAT-018の境界を反証的に確認した。固定Article Analysis、Assessment Map、Reading Value Assessment、推奨・Read／Skip／GainをFEAT-005のFixture入力、期待oracle、Case Result、Runtime入力へ要求する記述は残っていない。DEC-FEAT-019反映後の再レビューでは、`search-text`のSIGINTを既存の`integrationtest`限定同期gateで固定する。DEC-FEAT-020はEvidence単位Temporalの非破壊migrationと既存`get-evidence`の任意`temporal`表示だけを例外として許可し、既存値なしEvidenceのJSON形状、通常binary、既定挙動を維持する。

## 正規ソース

| ソース | 状態 | この設計への適用 |
| --- | --- | --- |
| `docs/features/FEAT-005/requirements.md` | confirmed | A〜Jの固定受入Fixture、層別診断、A/G/I/Jの既存FEAT-003検証契約への対応付け、Workflow非変更を求める。 |
| DEC-FEAT-016 | decided / L2 | 固定Fixture、隔離Store、層別oracleを定める。Reading Value実行については後続L4 Decisionに従う。 |
| DEC-FEAT-017 | decided / L2 | CLI境界とRuntime受入評価を分離し、既存Skillを変更しない。Runtime不在時は`not_run`とする。 |
| DEC-FEAT-018 | decided / L4、利用者承認済み | Reading Value Workflowを起動せず、A/G/I/Jを既存検証契約への`reading_value_reference`だけに限定する。 |
| DEC-FEAT-019 | decided / L3、利用者承認済み | `search-text`中断の再現性のため、`internal/persistence/sqlite/`の`integrationtest`限定同期gateだけを拡張できる。 |
| DEC-FEAT-020 | decided / L3、利用者承認済み | Evidence単位Temporalを非破壊migrationで保存し、既存`get-evidence`のEvidenceへ任意`temporal`を表示できる。 |
| `docs/features/FEAT-003/requirements.md`、`../skills/reading-value/references/verification.md` | approved / existing | V-002は三つの推奨の観測、V-004は記事・AI成果物をStoreへ保存しないことを、通常URL入力Workflowの検証として所有する。 |

## 契約根拠と照合

| 境界 | FEAT-005設計 | 根拠 | 判定 |
| --- | --- | --- | --- |
| Fixture配置 | `testdata/fixtures/`と`test/integration/`のテスト専用領域。通常Skillを変更しない。 | 要件の非変更境界、DEC-FEAT-017/018 | explicit・整合 |
| CLI／Store | 既存CLI、隔離Store、stdout/stderr/exit code、履歴差分を観測する。新operationは追加せず、DEC-FEAT-020のEvidence Temporal migrationと既存`get-evidence`任意表示だけを追加する。 | 要件7、DEC-FEAT-016/017/020、FEAT-001 | explicit・整合 |
| Runtime評価 | Search、Acquisition、Updateを固定Claim／Episodeで観測し、Runtime不在は`not_run`としてFeature受入に算入しない。 | DEC-FEAT-017、design実行モード | explicit・整合 |
| A/G/I/JのReading Value | `reading_value_reference`はV-002またはV-004の節IDと参照先の存在だけを確認する。Runtime入力・`expected`の実行oracle・Case Resultではない。 | DEC-FEAT-018、design Interfaces / Data、catalog | explicit・整合 |
| FEAT-003の所有範囲 | 推奨、Assessment Map、Read／Skip／Gain、固定Article AnalysisはFEAT-003通常Workflow検証が所有する。 | FEAT-003 requirements、verification V-001/V-002/V-004、DEC-FEAT-018 | explicit・整合 |
| `search-text`中断 | response開始前にSIGINTを送る同期を、`integrationtest` build tag時だけのSQLite read gateで固定する。 | DEC-FEAT-019、design Error Semantics、TASK005-02/04 | explicit・整合 |

## 資料間整合

- 要件のA/G/I/J対応はcatalogでA/I/J→V-002、G→V-004として一意に表され、各caseのSearch／Acquisition／Update期待値とは分離されている。
- catalogのA/G/I/Jには、固定Article Analysis、Assessment Map、Reading Value Assessment、推奨、Read／Skip／Gainを入力・期待oracleとして置いていない。I/JはSearchのAssessment／Traceだけを観測する。
- `design.md`の`expected`はAssessment、Trace、Candidate、Decision／Result、禁止操作、Store差分、停止結果に限定される。`reading_value_reference`は別論理項目であり、Case Resultの`first_mismatch_layer`にも混在しない。
- Assumptionsは固定Episodeだけを合成データとし、Article Analysisは本FeatureのFixtureへ含めない。これはFEAT-003の通常URL入力境界およびDEC-FEAT-018と整合する。
- DEC-FEAT-019の例外は`internal/persistence/sqlite/`の`integrationtest`限定同期gateである。DEC-FEAT-020はEvidence Temporal migration、Store read、既存`get-evidence`の任意`temporal`表示を追加する。tasks.mdとimplementation-handoff.yamlの変更可能範囲・禁止範囲はこの二つのDecisionに一致する。
- `git diff --check`は成功した。レビュー対象資料に未解消の形式不備は検出されなかった。

## 差し戻し事項

なし。
