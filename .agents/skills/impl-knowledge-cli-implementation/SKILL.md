---
name: impl-knowledge-cli-implementation
description: 承認済みのKnowledge CLI仕様をGoで実装する唯一のwriterとして、SQLite migration・transaction・公開CLI境界・fixture・プロセス結合テスト・package別coverageを証拠付きで完成させるSkill。オーケストレーターからREADY判定を受けた変更と、全gate通過後のIssue finalizationで使う。
---

# Knowledge CLI Implementer

## Preconditions

[orchestration-contract.md](../impl-knowledge-cli/references/orchestration-contract.md)を読み、`READY`のVerification Packet、Baseline Snapshot、対象Feature・Task・Issue・受入条件ID、承認済みの公開I/OとDecision参照を受け取る。

いずれかが欠ける、または正規資料がPacketと矛盾する場合は実装せず、`NEEDS_SPEC_RECHECK`として返す。

## Writable Scope

- `cmd/knowledge/`
- `internal/application/`
- `internal/domain/`
- `internal/persistence/sqlite/`
- `testdata/fixtures/`
- `test/integration/`
- 実装に必要な既存の開発用設定
- 全gate通過後に明示的に委譲された対象Issueのfinalization

要件、詳細設計、Decision、Task、handoffなどのPlanning成果物は変更しない。必要な変更はfindingとしてオーケストレーターへ返す。

## Implementation Procedure

1. Baseline Snapshotと現在のworktreeを照合し、今回の対象fileと既存変更を分離する。重なる既存変更を消去・巻き戻し・上書きしない。
2. 受入条件ごとに、変更する責務、先に追加するtest、必要なfixture、実プロセス検証を決める。追加・移動・rename・削除するfileはowner、呼出元、実行時期、配置理由を先に決める。移動・rename・削除前にhidden fileを含むrepository全体から旧pathと旧実行commandのcallerを検索する。仕様で固定されていない公開契約を補完しない。
3. 対象testを先に失敗させるか、既存の欠落を再現してから、責務境界内の最小変更を実装する。
4. 編集中は変更packageと近接integration testを実行し、短いfeedback loopを使う。
5. candidate diffが安定したら、全追加・移動fileの配置を再確認する。移動・rename・削除したfileはCI、Taskfile、Skill本文を含むrepository全体を再検索し、旧callerがなく全callerが新pathと実行commandを使うことを確認する。その後、format、literal layout検査、全体test、静的検査、coverage、該当process testを独立したgateとして実行する。
6. candidate ID、source ID、diff、file配置matrix、受入条件対応、command result、coverage、残存riskをImplementation Reportとして返す。この段階ではIssueを更新しない。

## Architecture Rules

- `cmd/knowledge`は公開CLI境界、validation、JSON/stdout/stderr/exit code、依存の組立てだけを担う。SQL、migration、domain ruleを置かない。
- `internal/application`はuse caseを調整し、transportやSQLite driverへ依存しない。
- `internal/domain`は外部I/Oへ依存せず、不変条件とportを所有する。
- `internal/persistence/sqlite`は接続、SQL、migration、派生字句Indexを所有し、applicationやCLI entryをimportしない。
- 公開option、JSON field、error code、exit code、保存先、設定、運用仕様は、Verification Packetが参照する承認済み契約以外を追加・変更しない。

## SQLite and Transaction Rules

- migration SQL assetはGo標準`embed`で同梱する。既適用migrationを編集せず、後続versionを追加する。
- migrationは新規DB、直前versionからのupgrade、再実行、途中失敗時rollbackを検証する。version記録とdata更新を同じatomic boundaryに置く。
- mutationはcommit成功後だけsuccessを返し、失敗・cancel時に部分更新を残さない。
- SQLは使用関数の直前に個別のpackage-level定数として複数行で置く。値は`?`placeholderで束縛し、固定個数の`IN`条件を文字列置換で組み立てない。
- null、空結果、境界値、重複、順序、複数filterの組合せを、該当する仕様に沿ってtestする。

## Test and Fixture Rules

- production packageのunit testはtable-drivenを基本とし、正常系だけでなくvalidation、not-found、conflict、storage failure、rollbackを対象にする。
- 公開CLI変更は実binaryを起動する`test/integration/`で、引数、stdout、stderr JSON、exit code、DB事後状態を検証する。
- 再利用する期待値と入力dataは`testdata/fixtures/`を正本とし、test codeとTaskfileへ重複させない。
- testごとに隔離SQLite DBと隔離したOS設定directoryを使い、公開Store overrideを追加しない。
- 変更した各production packageを`python3 .agents/skills/impl-knowledge-cli-implementation/scripts/check_test_coverage.py`で個別測定し、statement coverage 100%を確認する。coverage率だけでassertionの妥当性を判断しない。

## File and Script Placement Rules

- fileは実行時責務、変更理由、主な呼出元、更新周期を所有する最小の既存directoryへ置く。便宜だけでrepository rootや汎用`utils`へ置かない。
- product runtime codeは対応する`cmd` / `application` / `domain` / `persistence`へ、process testは`test/integration`へ、再利用fixtureは`testdata/fixtures`へ置く。
- Skillの工程同一性を扱うscriptはオーケストレーターの`scripts/`、実装品質gateを扱うscriptは実装Skillの`scripts/`へ置く。Skill補助scriptはPythonへ統一し、`python3`で明示的に呼ぶ。
- 一回限りの調査・変換scriptはtemporary directoryで実行し、candidateへ残さない。継続運用するscriptだけを、明確なcaller、test、ownerがある配置へ追加する。
- 既存Taskfile、test helper、Skill scriptと同じ責務を重複実装しない。新しい配置区分や共有ownerが必要なら、実装で既成事実化せず`BLOCKED`としてオーケストレーターへ返す。

## Coding Conventions

- Goとmigration SQLの説明commentは、識別子・compiler指示を除き日本語で端的に書く。
- 複数要素のarray/slice literalとcomposite literalのfieldは1行ずつ記述する。自動検査は、既存gateとの互換性を保つため同一行に複数のkeyed fieldがあるcomposite literalを対象とする。
- `python3 .agents/skills/impl-knowledge-cli-implementation/scripts/check_composite_literal_layout.py`を実行する。

## Validation Strategy

編集loopでは変更packageのtestと該当integration testを使う。最終candidateでは`gofmt`後に、`go test ./...`、`go vet ./...`、`task lint`、package別coverage script、変更に対応するCLI process integration testをすべて実行する。

`codex exec`、またはそれを内部で起動する人間専用Runtime受入taskは、AIが実行してはならない。Runtime受入が必要な場合は、通常の自動gateと別taskへ分離し、目的、前提となる認証、実行コマンド、期待される観測を利用者へ示す。実行結果は利用者が提示した場合だけ受け取り、AIの自動検証結果として扱わない。

互いにrepository artifactを書かないcommandは並列実行してよい。全必須gate成功時だけ`READY_FOR_REVIEW`を返す。失敗時は`BLOCKED`とし、command、exit status、主要error、未実行gateを記録する。正規sourceのfingerprintが変わった場合は`NEEDS_SPEC_RECHECK`とする。

## Final Candidate Checklist

- 追加・移動した全fileについてowner、呼出元、実行時期、配置理由を説明でき、孤立script、一回限りのscript、責務重複、repository rootへの仮置きがない。
- 移動・rename・削除したfileの旧pathと旧実行commandをhidden fileを含むrepository全体で検索し、CI、Taskfile、Skill本文を含む全callerが更新されている。
- Skill補助scriptは所有Skillの`scripts/`配下にあり、Python以外の補助scriptや暗黙の実行方法を増やしていない。
- 公開I/O、migration、transaction、fixture、process testがVerification Packetの受入条件へ対応している。
- format、literal layout、`go test ./...`、`go vet ./...`、`task lint`、package別coverage、該当process integration testが同一candidateで成功している。

## Issue Finalization Mode

オーケストレーターから、同じcandidate IDに対するReview Report=`PASS`、Audit Report=`PASS`、受入条件別証拠を受け取った場合だけ実行する。

1. Issue本文を読み直す。
2. 証拠がある受入条件だけを完了へ更新する。
3. 未対応・未検証項目は未完了のまま理由を記録する。
4. 実装結果と検証commandを簡潔に追記する。
5. 更新後のIssueを読み戻し、変更内容を返す。

Issueの受入条件を削除・弱体化しない。Issueを閉じるのは、ユーザー指示または既存workflowで明示されている場合だけとする。

## Output

共通contractのImplementation Report形式で返す。実測していないcommandを`PASS`にせず、既存変更と今回変更を区別し、公開契約差分がない場合も`none`と明記する。
