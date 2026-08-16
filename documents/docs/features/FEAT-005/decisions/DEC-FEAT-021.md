# DEC-FEAT-021: テスト専用Runtime受入launcher

- **Status:** decided
- **Level:** L2（FEAT-005のテスト実行方式）
- **Date:** 2026-08-16
- **Approval:** 利用者承認（2026-08-16）

## Decision

`test/integration/` に、明示的に起動した場合だけCodex Runtimeを呼ぶ受入launcherを置く。launcherは一つの`case_id`ごとに、Fixtureからseedした隔離Store、既存Skill、既存`knowledge` binaryだけを渡す。Runtimeは一時Markdown成果物と`execution_mode: runtime_acceptance`のCase Resultを呼出し元へ返す。

通常の`go test ./...`、CI、およびAIが実行する検証はRuntimeを起動しない。Runtime受入は人間専用の`task test:runtime-acceptance`だけで明示起動し、AIはそのtaskまたは`codex exec`を実行せず、利用者が提示した結果だけを評価する。taskは一意な実行ログをGit管理外の`.runtime-acceptance-logs/`へ保存する。launcherは環境変数で明示された実行だけを受け付け、結果をrepository、Skill本文、公開CLI、通常Storeへ保存しない。Reading Valueは起動しない。

## Consequences

- A〜J／XのRuntime受入は、Codex Runtimeが利用可能な開発環境でだけ実行する。
- 実行不能時は`not_run`を記録し、Feature合格へ読み替えない。
- 既存Workflowの本文・通常入力契約、公開CLI/JSON、SQLite schemaは変更しない。

## Affected Artifacts

- `docs/features/FEAT-005/design.md`
- `docs/features/FEAT-005/tasks.md`
- `docs/features/FEAT-005/implementation-handoff.yaml`
- `testdata/fixtures/acceptance/knowledge-quality/`
- `test/integration/`
