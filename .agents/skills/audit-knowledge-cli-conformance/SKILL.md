---
name: audit-knowledge-cli-conformance
description: Knowledge CLIのコードレビューPASS後に、現行実装と全体利用経路を要件、詳細仕様、承認済みDecision、handoff、Issue、既存DB、実プロセス挙動へ読み取り専用で再対応付けし、差分レビューでは見落とす矛盾・仕様逸脱・既存機能破壊を判定する最終整合性監査Skill。
---

# Knowledge CLI Conformance Auditor

## Read-only Boundary

ファイル、Issue、Planning成果物、外部状態を変更しない。コード品質レビューを繰り返すのではなく、利用者が複数操作を組み合わせたときに要件どおり使えるかを独立して反証する。[orchestration-contract.md](../impl-knowledge-cli/references/orchestration-contract.md)を読み、Audit Reportを返す。

## Preconditions

- Verification Packet=`READY`
- 現在candidate IDに対するReview Report=`PASS`
- Implementation Report=`READY_FOR_REVIEW`とcandidate IDが一致する
- Verification Packetのsource IDを再計算し、一致する

reviewにblocking findingが残る場合は監査せず`BLOCKED`として返す。source IDが変わった場合は監査せず`NEEDS_SPEC_RECHECK`として返す。anchoringを避けるため、詳細なreview結論は入力にせず、`PASS`の事実と非blocking riskだけを受け取る。

## Procedure

1. 要件、Business Rule、承認済みDecision、詳細設計、operation資料、use case、handoff、Issue受入条件を読み、Verification Packetのmatrixを原典から再確認する。
2. 変更した公開operation、shared CLI boundary、domain port、table/index、migrationを起点にimpact mapを作る。
3. diffだけでなく、変更層を共有する既存operation、直前schemaから移行したDB、通常compositionの利用経路を調べる。
4. 各受入条件を、仕様、production path、fixture、実binaryの観測結果へ対応付ける。test名やcoverage率だけを証拠にしない。
5. 変更に応じて、少なくとも次のうち該当する経路を確認する。
   - 新規DBと直前version DBでの起動・migration・再起動
   - create/revise後にsearch/get/get-evidenceなどへ値を受け渡す操作列
   - selector、Alias、Scope AND、null、期間境界、順序、空結果
   - validation/storage/conflict/cancel時のstdout、stderr、exit code、DB不変性
   - 変更していない近接operationを1つ以上通すregression経路
6. process testが通常compositionを使うか、fixtureの期待値が仕様から導かれ、実装結果の写しになっていないか確認する。
7. code reviewのcommandを無条件に繰り返さず、全体利用の不足証拠に絞って安全なprocess testを実行する。実行できない経路は未検証riskとして隠さない。

## Distinction from Code Review

- reviewerはcandidate diffの実装品質と局所的な正しさを判定する。
- auditorは要件から利用経路へ縦断し、既存操作と文書間を横断する。
- auditorはreviewerの`PASS`を要件適合の証明として流用しない。

## Findings and Verdict

findingは`AUD-001`のような安定IDを付け、candidate ID、要件根拠、影響する利用経路、再現手順、観測結果、不足oracle、修正先を示す。

- 実装・test・fixtureで直せる矛盾または破壊: `BLOCKED`
- 公開契約、Decision、Planning資料の再判断が必要: `NEEDS_SPEC_RECHECK`
- 全受入条件と影響経路にblocking矛盾がない: `PASS`

再監査では以前の`AUD-nnn`を引き継ぎ、同じ論点はIDを維持して共通contractのfinding statusを更新する。

`PASS`でも、監査した既存operation、DB version、process caseと、実行不能だった残存riskを明示する。Issue checkboxは監査中に変更しない。
