# Task Session Handoff

このファイルはローカルセッション用であり、Gitへコミットしない。

## Task

- Issue: [#{{ISSUE_NUMBER}}]({{ISSUE_URL}})
- Task ID: `{{TASK_ID}}`
- Task名: {{ISSUE_TITLE}}
- Branch: `{{BRANCH}}`
- Worktree: `{{WORKTREE_PATH}}`
- Base ref: `{{BASE_REF}}`
- Worktree起点SHA: `{{BASE_SHA}}`
- Gate通過Commit: `{{GATE_COMMIT}}`
- Planning snapshot: `{{PLANNING_SNAPSHOT}}`
- 書込み可能なPath／Glob: {{OWNER_PATHS}}

## 開始手順

1. このファイルを読む。
2. GitHub Issue本文と固定Planning snapshotを読む。
3. 現在branch、HEAD、依存Evidence、Gate Evidenceを確認する。
4. 書込み可能なPath／Glob以外を変更しない。
5. Task固有の完了条件と検証方法を確認してから作業を始める。

## Codex開始プロンプト

```text
$manage-task-worktrees を使って .codex/task-session.local.md とTask Issueを読み、
記載されたleaf Taskだけを実行してください。
依存関係、Gate、成果物Ownerに差異があれば作業を止めて報告してください。
```

## 終了時

- Task固有テストまたはReviewを実行する。
- 変更PathがOwner境界内であることを確認する。
- Commit、PR、検証結果、統合CommitをIssueのhuman-progress領域へ記録する。
- Merge前にworktreeを削除しない。

