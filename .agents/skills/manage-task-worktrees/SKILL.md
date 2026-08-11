---
name: manage-task-worktrees
description: GitHub Task Issueごとに隔離したGit worktreeとブランチを作成し、VS Codeで別Codexセッションへ引き継ぎ、状態確認、完了準備、安全な削除まで管理する。依存Task・Gate・基準Commit・書込みPathを確認しながら並行作業を始める場合や、worktreeを再開・整理する場合に使用する。
---

# Task Worktree管理

Task Issueを1つのブランチ、worktree、VS Codeウィンドウ、Codexセッションへ対応付ける。

## 必須手順

1. 変更操作の前に[ライフサイクル](references/lifecycle.md)を全文読む。
2. Primary checkoutから `plan` を実行し、Issue、依存Task、Gate、基準ref、Path競合を確認する。
3. 作成内容と外部操作をユーザーへ示し、必要な承認を得てから `start` を実行する。
4. VS Codeが開いたら、ユーザーへCodexの開始と引継ぎファイルの読込みを依頼する。
5. 作業完了時は `finish` とIssue固有の検証を実行し、Commit・PR・IssueへEvidenceを残す。
6. Merge後に限り、明示承認を得て `remove` を実行する。

GitHub Issue、ネットワーク、GUI、push、PR、削除を伴う操作では、それぞれ既存の承認規則に従う。

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

`plan`はworktreeを作成しない。`start --no-open`はVS Codeを開かず、worktreeと引継ぎファイルだけを作成する。

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
- 自動Merge、自動push、自動Issue closeは行わない。

