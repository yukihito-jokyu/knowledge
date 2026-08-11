# Worktreeライフサイクル

## 前提

- Primary checkoutのリポジトリ直下から実行する。
- `git`、`gh`、`code`を利用でき、`gh auth status`が成功する。
- `/.worktrees/`と`/.codex/task-session.local.md`が`.gitignore`に登録されている。
- 実行対象は`Lx-My-Sz`形式のleaf Task Issueである。
- GitHub Issue本文のTask ID Marker、依存表、成果物Ownerを正本として読む。

## 状態遷移

```text
未作成
  └─ plan合格
       └─ start
            └─ 作業中
                 ├─ open/status
                 └─ finish確認
                      └─ Commit・push・PR・Merge
                           └─ remove
                                └─ worktree削除済み
```

`plan`と`status`はread-onlyである。`start`はbranch、worktree、ローカル引継ぎファイルを作成し、既定ではVS Codeを開く。`finish`は完了判定材料を表示するだけでCommitやpushを行わない。`remove`はMerge済みでcleanなworktreeだけを削除する。

## 基準ref

通常は`origin/main`を使用する。次の場合は`--base`を必須とする。

- Design Freeze Gate通過Commitから複数Taskを並行開始する。
- 依存Taskを統合した専用integration refから開始する。
- `origin/main`へ未統合だが、承認済みの共通起点Commitがある。

依存Issueはclosedだけでは不十分である。Issueのhuman-progress領域に統合Commitを記録し、そのCommitが基準refに含まれることを確認する。

Gate依存があるTaskは`--gate-commit <sha>`を指定する。スクリプトはGate Commitが基準refに含まれることを確認する。Gateの意味的な合格判定は行わない。

## 開始

最初にdry-runする。

```shell
rtk bash .agents/skills/manage-task-worktrees/scripts/manage_worktree.sh plan 28
```

表示内容を確認後、外部アクセスとGUI起動の承認を得て開始する。

```shell
rtk bash .agents/skills/manage-task-worktrees/scripts/manage_worktree.sh start 28
```

VS Codeを後から開く場合:

```shell
rtk bash .agents/skills/manage-task-worktrees/scripts/manage_worktree.sh start 28 --no-open
rtk bash .agents/skills/manage-task-worktrees/scripts/manage_worktree.sh open 28
```

同じTaskのbranchとworktreeが正しく存在する場合、`start`は破壊せず再開として扱う。別Pathで同じbranchがcheckoutされている場合は停止する。

## 新しいCodexセッション

VS Code内でCodexを手動開始し、次のように依頼する。

```text
$manage-task-worktrees を使って .codex/task-session.local.md を読み、
記載されたTask Issueだけを実行してください。
```

新しいセッションは、着手前に次を確認する。

1. Task IDとIssue番号が一致する。
2. 現在branchと引継ぎbranchが一致する。
3. HEADと引継ぎのworktree起点SHAが整合する。
4. 書込みPathがIssueの単一Owner境界内である。
5. 依存・GateのEvidenceが引継ぎとIssueに存在する。

## 並行作業

同じGate Commitから別worktreeを作る場合も、次を満たす必要がある。

- branch名とworktree PathがTaskごとに一意である。
- 書込みPath／Globが重ならない。
- shared Registry、Lockfile、生成物、共通Fixtureは承認済み単一Ownerだけが変更する。
- Merge順があるTaskは、worktree作成を並行化しても確定・Mergeを直列化する。
- 上流契約をconsumer側から変更しない。

Path競合はスクリプトだけで完全判定しない。CodexがIssueの「成果物と所有」とTask Mapを読み、意味的に確認する。

## 完了準備

```shell
rtk bash .agents/skills/manage-task-worktrees/scripts/manage_worktree.sh finish 28
```

出力されたbranch、HEAD、変更Path、Owner PathをIssueと照合する。その後、Task固有テストを実行し、ユーザーの依頼範囲に従ってCommit、push、PRを行う。

Issueへ少なくとも次を記録する。

- worktree起点SHA
- 実装または設計Commit
- 実行した検証Commandと結果
- PR
- 統合Commit
- Gate／Releaseに渡すEvidence
- 未解決事項の有無

## 削除

PRのMergeとIssue Evidence記録後に実行する。

```shell
rtk bash .agents/skills/manage-task-worktrees/scripts/manage_worktree.sh remove 28 \
  --merged-into origin/main \
  --confirm
```

ローカルbranchも削除する場合だけ`--delete-branch`を付ける。スクリプトはdirty状態、Merge未完了、Path不一致で停止する。`--force`は提供しない。

## 障害時

- Issue Markerがない: tracking Issueの可能性があるため開始しない。
- Task IDがleafでない: 親Taskは進捗管理用なので開始しない。
- 依存Commitがない: 依存IssueへEvidenceを記録する。
- Gate Commitがない: Gate確認セッションへ戻る。
- 基準refにCommitがない: integration順を修正し、refを更新する。
- Pathが既に存在する: 登録済みworktreeとbranchを確認し、推測で削除しない。
- `code`がない: `start --no-open`を使用し、利用可能なVS Code起動方法をユーザーと決める。

