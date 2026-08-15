---
name: impl-knowledge-workflow
description: KnowledgeプロダクトのCodex workflow Skillを、承認済みhandoffに従って実装・検証するwriter。skills/配下のReading Value、Knowledge Acquisition、Knowledge Updateとその参照契約を変更するTaskで使う。
---

# Knowledge Workflow Implementer

## Preconditions

承認済みのimplementation handoff、対象Task、関連Decision、既存の対象Skillと参照契約を読む。公開CLI、JSON、SQLite、migration、またはPlanning成果物の変更が必要なら編集せず、仕様照合へ戻す。

## Writable Scope

- `skills/reading-value/`
- `skills/knowledge-acquisition/`
- `skills/knowledge-update/`
- 変更内容を検証するための `.ai/skill-tests/`

要件、設計、Decision、Task、handoff、`cmd/knowledge/`、`internal/`、SQLite schema・migration、既存Knowledge CLIの公開I/Oは変更しない。

## Procedure

1. baseline、承認済みsource fingerprint、既存変更を確認する。対象外の既存変更を消去・上書きしない。
2. 受入条件ごとに、変更するWorkflow境界、呼出し順、停止経路、一時成果物、本文不変のoracleを対応付ける。
3. 対象Skillと同梱参照を最小限更新する。未承認のfield、保存先、再実行、非同期実行、CLI operationを補完しない。
4. 対象Skillを先頭から読み、正常・失敗・中断経路が受入条件と矛盾しないことを確認する。必要なら受入境界を再現するskill testを追加する。
5. `git diff --check`、対象Skillのfrontmatter検証、対象Skillと参照契約の相互参照を実行する。変更内容、受入条件別証拠、実行command、残存riskを報告する。

## Completion Criteria

- 承認済みの実装先だけを変更している。
- 同期順序、停止経路、本文不変、一時成果物境界を明示している。
- 公開CLI・DB・永続ledger・非同期実行を追加していない。
- 実行した検証と未実行範囲を区別して報告している。
