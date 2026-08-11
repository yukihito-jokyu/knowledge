---
name: manage-task-worktrees
description: GitHub tracking Issueから実行可能なleaf Taskを展開し、依存DAGと書込みPathから並行waveを提案したうえで、Taskごとに隔離したGit worktree、ブランチ、VS Code、Codexセッションを作成・再開・完了・削除まで管理する。親Issueを起点に作業を始める場合、依存Task・Gate・基準Commit・Path競合を確認して並行作業を計画する場合に使用する。
---

# Task Worktree管理

Task Issueを1つのブランチ、worktree、VS Codeウィンドウ、Codexセッションへ対応付ける。

## 必須手順

1. 変更操作の前に[ライフサイクル](references/lifecycle.md)を全文読む。
2. 指定Issueがtracking Taskなら、`plan`や`start`を実行せず、Issue本文と固定Planning snapshotから子孫leafを再帰的に展開する。
3. leafの依存DAG、Gate、基準ref、Owner Path、既存worktreeを照合し、現在のready frontierと将来waveをユーザーへ提案する。
4. ユーザーが選んだ現在waveの各leafについてPrimary checkoutから `plan` を実行する。機械判定に加え、複雑Globと共有資産のPath競合を意味的に監査する。
5. wave全体の作成内容、並行・直列理由、Merge順、外部操作を示し、必要な承認を得てから各leafの `start` を実行する。後続waveへ承認を持ち越さない。
6. VS Codeが開いたら、ユーザーへ各ウィンドウでCodexの開始と引継ぎファイルの読込みを依頼する。
7. 作業完了時は `finish` とIssue固有の検証を実行し、Commit・PR・IssueへEvidenceを残す。
8. Merge後に限り、明示承認を得て `remove` を実行する。

GitHub Issue、ネットワーク、GUI、push、PR、削除を伴う操作では、それぞれ既存の承認規則に従う。

## 親Issueからの開始

`Lx`と`Lx-My`はtracking Task、`Lx-My-Sz`はworktreeを持つleaf Taskとして扱う。tracking Issueを指定された場合は[ライフサイクルの親Issue展開](references/lifecycle.md#親issueからleafへの展開)を実行し、次を表で提示する。

- 現在開始できるready frontierと、依存統合後に解放される将来wave
- 各leafのIssue、Task ID、依存Evidence、Gate、Owner Path、base ref／SHA
- 並行可能な組、直列化する理由、Merge順、次waveの解放条件

将来waveのbase SHAは先行Taskの統合前に確定しない。tracking Issue自身へ `start` を実行せず、ユーザーが承認した現在waveのleafだけを `plan`、`start` する。

## コマンド

Skillディレクトリを `<skill>` として、すべてPrimary checkoutから実行する。

```shell
rtk bash <skill>/scripts/manage_worktree.sh plan <issue-number>
rtk bash <skill>/scripts/manage_worktree.sh start <issue-number>
rtk bash <skill>/scripts/manage_worktree.sh status [issue-number]
rtk bash <skill>/scripts/manage_worktree.sh open <issue-number>
rtk bash <skill>/scripts/manage_worktree.sh finish <issue-number>
rtk bash <skill>/scripts/manage_worktree.sh remove <issue-number> --merged-into <ref> --confirm
```

通常の基準refは `origin/main` とする。Gate後のTaskはGate通過Commitを明示する。

```shell
rtk bash <skill>/scripts/manage_worktree.sh plan 46 --base <gate-ref> --gate-commit <sha>
rtk bash <skill>/scripts/manage_worktree.sh start 46 --base <gate-ref> --gate-commit <sha>
```

`plan`はworktreeを作成しない。単純なPath重複は検査するが、複雑Glob、除外規則、Schema、Migration、Interface、Registry、DI、Lockfile、生成物、共有FixtureはCodexがTask MapとIssueを意味的に監査する。`start --no-open`はVS Codeを開かず、worktreeと引継ぎファイルだけを作成する。

## 配置規則

Issue #28、Task ID `L1-M1-S1`の例:

```text
worktree: <repo>/.worktrees/issue-28-l1-m1-s1
branch:   task/issue-28-l1-m1-s1
handoff: <worktree>/.codex/task-session.local.md
```

`/.worktrees/`と`/.codex/task-session.local.md`はGit管理対象外にする。worktree内から別worktreeを作成しない。

## Codexへの引継ぎ

`start`はVS Codeを新規ウィンドウで開き、ローカル引継ぎファイルを表示する。Codex UIの開始はユーザーが手動で行う。新しいCodexは次を守る。

- 引継ぎファイル、Task Issue、固定Planning snapshotを最初に読む。
- Issueの単一Owner Pathだけを変更する。
- 意味判断、依存関係、Gate条件を作業セッション内で追加しない。
- 契約差異や未解決TBDを見つけた場合は実装を止め、元セッションへ戻す。
- 親Taskではなくleaf Issueだけを実行単位にする。

## 安全境界

- dirtyなworktreeを削除しない。
- `--merged-into`へHEADが含まれないworktreeを削除しない。
- `--force`、`git reset --hard`、未確認のbranch削除を使わない。
- 依存Issueが未完了、統合Commitが未記録、または基準refに含まれない場合は開始しない。
- Gate依存がある場合、`--gate-commit`なしでは開始しない。
- 並行Taskの書込みPathが重なる場合は開始せず、単一Ownerへ直列化する。
- 自動Path検査の合格だけで並行可能と断定しない。
- 未完了の将来wave用worktreeを先に作成して依存確認を迂回しない。
- 自動Merge、自動push、自動Issue closeは行わない。
