# DEC-FEAT-009: AssessmentはCodex側Knowledge Searchワークフローとして実行する

- **Status:** decided
- **Level:** L2
- **Decision:** FEAT-002のClaim反復探索、Evidenceに基づく評価、7状態の判定、Assessment/Traceの組立ては、Codex側のKnowledge Searchワークフローとして実行する。Knowledge CLIには既存の決定論的検索・取得operationの呼出しだけを行い、Goコード、新規CLI operation、JSON wire schema、SQLite schema、公開設定を追加・変更しない。

## 根拠

人間Decision（2026-08-15）は、Issue #197 / TASK-002-01の実装先をKnowledge CLIではなく既存ArchitectureのCodex Knowledge Search側とした。これはBR-010、CON-002/003、Architectureの責務境界と一致する。実行形態の詳細は[Knowledge Search Workflow 契約](../design/knowledge-search-workflow.md)で定める。

## 影響

TASK-002-01はKnowledge CLIのapplication/domain実装ではなく、Codex側ワークフローの探索・評価責務へ再分解する必要がある。TASK-002-02とTASK-002-03は同ワークフローの成果物境界と検証として、再分解後のTASK-002-01へ依存する。Taskおよびimplementation handoffの更新はtask-breakdown ownerが行う。
