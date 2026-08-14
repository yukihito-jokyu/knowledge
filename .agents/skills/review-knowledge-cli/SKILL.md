---
name: review-knowledge-cli
description: Knowledge CLIの候補diffを独立かつ読み取り専用でレビューし、責務分離、公開CLI境界、SQLite migration・transaction、失敗経路、test oracle、変更package別coverageの欠陥をfinding ID付きで判定するSkill。実装証拠が揃った後、最終整合性監査の前に使う。
---

# Knowledge CLI Reviewer

## Read-only Boundary

ファイル、Issue、Planning成果物、外部状態を変更しない。実装者の説明を正しい前提にせず、実際のdiff、正規資料、test codeから反証する。[orchestration-contract.md](../impl-knowledge-cli/references/orchestration-contract.md)を読み、Review Reportを返す。

## Required Inputs

- Verification Packet=`READY`
- Baseline Snapshot
- candidate ID、candidate diff、変更file一覧
- Implementation Report=`READY_FOR_REVIEW`
- 追加・移動fileのfile配置matrix
- 実行済みcommandと受入条件別test対応

candidate IDまたはsource IDを再計算できない、候補diffが証拠作成後に変わっている、必須commandが失敗・未実行、または既存変更との境界が不明なら`BLOCKED`として再採取を求める。

## Procedure

1. Baselineとcandidate diffを比較し、今回変更、既存変更、未追跡fileを分ける。Planning成果物の無断変更とscope外変更を探す。
2. 追加・移動fileの配置matrixを実diffと照合し、owner、呼出元、実行時期、更新周期が配置先の責務と一致するか確認する。孤立script、一回限りのscript、rootへの仮置き、既存Task・helperとの重複を探す。移動・rename・削除がある場合は、hidden fileを含むrepository全体から旧pathと旧実行commandを検索し、CI、Taskfile、Skill本文を含む全callerが更新済みか反証する。
3. Verification Packetの各受入条件を、production code、unit test、fixture、process testへ追跡する。
4. CLI / application / domain / persistenceの依存方向と責務を確認する。
5. 公開I/Oが承認済みoption、JSON、stdout/stderr、error、exit codeとbyte-levelで一致するか確認する。
6. SQLite migration、transaction、SQL、data互換性、失敗・cancel経路を確認する。
7. testが実装と同じ誤りを複製していないか、境界値、否定例、DB事後状態をassertしているか確認する。
8. `python3 .agents/skills/review-knowledge-cli/scripts/review_test_coverage.py --base <Baseline SnapshotのHEAD> --json`でcandidateが変更したproduction packageを抽出してcoverageを独立再測定し、Implementation Reportの対象package一覧と実測値を照合する。100%でも弱いassertionや未実行process経路を指摘する。
9. 証拠が現在candidateと一致し疑義がなければfull gateを無条件に再実行しない。疑義、再現、欠落証拠に絞って読み取り専用commandを実行し、実行内容を報告する。

## Review Checklist

- validation前後で副作用が起きない
- successはcommit後だけ返り、rollback後のDBが不変
- migrationの新規DB、upgrade、再実行、失敗時version/data rollback
- null、空結果、重複、順序、開放境界、最大/最小値
- SQL placeholder、row close/error、transaction lifecycle、context/cancel伝播
- default Storeとprocess testのOS設定directory隔離
- fixtureが期待I/Oの正本で、test codeに契約が重複していない
- review用coverage scriptが抽出した全変更production packageに実測結果があり、100%未満、test失敗、実装報告からのpackage欠落をfindingにしている
- 追加・移動fileが責務ownerのdirectoryにあり、呼出元と更新周期に反する逆依存や孤立がない
- 一回限りの調査・変換scriptがcandidateへ残らず、継続運用scriptにはcallerとtestがある
- Skill補助scriptが所有Skillの`scripts/`にあり、Pythonへ統一され、別拡張子のhelperや同じ検査の重複実装がない
- 移動・rename・削除したfileの旧pathや旧実行commandを参照するcallerが、hidden file、CI、Taskfile、Skill本文に残っていない
- 変更した共通層を使う近接操作に明白なregressionがない

## Findings and Verdict

findingは`REV-001`のような安定IDを付け、severity、candidate ID、根拠file/symbol、影響する受入条件、再現または不足test、最小の修正方向を示す。

- `BLOCKER` / `HIGH`: data loss、公開契約違反、誤結果、atomicity、必須検証欠落。`BLOCKED`。
- `MEDIUM`: 現実的なedge case、test oracle、保守性が正しさへ影響する欠陥。原則`BLOCKED`。
- `LOW`: 要件・正しさへ影響しない改善。非blocking riskとして残せる。
- 公開契約または正規資料の再判断が必要: `NEEDS_SPEC_RECHECK`。

`PASS`はblocking findingがなく、受入条件別の検証証拠を確認できた場合だけ返す。指摘がない場合も、確認範囲、実行command、未検証riskを明示する。
