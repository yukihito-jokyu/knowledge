# DEC-FEAT-017: 自動CLI境界とRuntime受入評価を分離する

- **Status:** decided
- **Level:** L2（FEAT-005内のテスト実行方式）
- **Date:** 2026-08-15

## Context

FEAT-005は既存Knowledge CLIのプロセス境界だけでなく、Codex Runtime内で生成されるAssessment、Search Trace、Candidate、Update Resultを層別に観測する。Reading Valueは通常URL入力の既存FEAT-003検証契約を参照する。既存Go integration testは実`knowledge` binaryと隔離SQLite Storeを実行できるが、Codex RuntimeのWorkflowを起動する公開CLIやCI runnerは存在しない。

通常利用時に読み込まれる`skills/`へFixtureや受入oracleを追加すると、実行コンテキストを不要に増やす。

## Decision

同じ固定Fixtureと`case_id`を使い、次の二つを分離する。

1. CI／Go integration testは、Fixture構造、実CLI、stdout/stderr、exit code、隔離Store、履歴・差分、技術失敗・中断を自動観測する。
2. 開発者が起動するテスト専用のCodex Runtime評価は、Fixture、既存Skill、隔離Store接続情報を一回だけ受け取り、Workflowの一時Markdown成果物とCase Resultを呼出しsessionへ返す。

Runtime評価の指示はFixture領域に置く。既存Skillは変更せず、成果物はStore、repository、公開UI/API、Skill本文へ保存しない。CIにCodex Runtimeがない場合、Runtime評価を自動実行したことにせず`not_run`として理由を残す。A〜J／XのCLI／Search／Acquisition／UpdateのFeature受入はRuntime評価が完了するまで成立しない。Reading Valueは既存FEAT-003検証契約への参照を確認する。

## Considered Alternatives

### 受入oracleを既存Skillへ追記する

通常利用時に不要なFixture・判断表が読み込まれ、実行コンテキストを増やすため採用しない。

### Go integration testだけで全層を合格とする

Codex Runtime成果物を生成できず、Search、Acquisition、Updateの受入oracleを観測できないため採用しない。Reading Valueは既存FEAT-003検証契約で観測する。

### RuntimeをCIに必須導入する

採用済みのRuntime runner、認証、CI実行契約がなく、FEAT-005内で新しい運用・公開境界を導入することになるため採用しない。

## Consequences

- Case Resultは`case_id`と`execution_mode`（`cli_boundary`／`runtime_acceptance`）の組で区別し、一方の結果で他方を上書きしない。
- Runtime未利用のCIはFeature全体を合格にできない。
- 通常のSkill実行コンテキストと既存Workflow契約は増えない。

## Affected Artifacts

- `docs/features/FEAT-005/design.md`
- `docs/features/FEAT-005/design/scenario-catalog.md`
- `docs/features/FEAT-005/tasks.md`
- `docs/features/FEAT-005/implementation-handoff.yaml`
