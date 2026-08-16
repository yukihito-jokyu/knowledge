---
name: impl-knowledge-workflow
description: KnowledgeプロダクトのCodex workflow Skillと承認済み検証用コピーを、承認済みhandoffに従って実装・検証するwriter。skills/配下のReading Value、Knowledge Acquisition、Knowledge Updateとその参照契約を変更するTaskで使う。
---

# Knowledge Workflow Implementer

## Preconditions

承認済みのimplementation handoff、対象Task、関連Decision、既存の対象Skillと参照契約を読む。公開CLI、JSON、SQLite、migration、またはPlanning成果物の変更が必要なら編集せず、仕様照合へ戻す。

## Writable Scope

- `skills/reading-value/`
- `skills/knowledge-acquisition/`
- `skills/knowledge-update/`
- 変更内容を検証するための `.ai/skill-tests/`
- 承認済みhandoffが明示する専用検証ブランチ上の `verification/`

要件、設計、Decision、Task、handoff、`cmd/knowledge/`、`internal/`、SQLite schema・migration、既存Knowledge CLIの公開I/Oは変更しない。

## Procedure

1. baseline、承認済みsource fingerprint、既存変更を確認する。対象外の既存変更を消去・上書きしない。
2. 受入条件ごとに、変更するWorkflow境界、呼出し順、停止経路、一時成果物、本文不変のoracleを対応付ける。
3. 対象Skillと同梱参照を最小限更新する。未承認のfield、保存先、再実行、非同期実行、CLI operationを補完しない。
4. handoffが専用検証ブランチの`verification/`を対象にするときは、CLI実装branchへreview/audit前にcommitされたクリーンcandidateについて、独立コードレビューと最終整合性監査がともにPASSしたCLI source commit、candidate ID、`candidate_fingerprint.py --json`の完全な出力、順序付きsource inventoryと`source_fingerprint.py`出力だけをbinary同期元にする。orchestratorから渡されたPASS済みReview/Audit Report全文を改変せず、`verification/evidence/<cli-source-commit>/review-report.md`と`audit-report.md`へ保存し、各reportのCLI source commitとcandidate IDが同期元と一致することを確認する。CLI source commitから新しい隔離worktreeを作り、そのworktreeだけでbinaryをbuildする。root workflow Skillのコピーは、workflow writer自身が検証してcommitした別のclean workflow source commitから取得し、コピー対象rootを順序付きで固定して`source_fingerprint.py`出力を得る。専用ブランチの`verification/source-manifest.json`へ、両source commit、CLI candidate出力・source inventory/fingerprint、workflow root/fingerprint、Review/Audit Report各path/SHA-256/verdict/source commit/candidate ID、binaryおよびコピー対象の各repo相対pathとSHA-256を保存する。report本文、manifest、CLI fingerprintの値に不一致があれば停止する。manifest、report、binary、Skillコピーを同一commitでpushし、その同期commitはGit履歴で照合する。Feature実装ブランチへ`verification/`を混在させず、いずれかのsource、clean candidate、PASS、report、manifestが未確定なら停止する。
5. 対象Skillを先頭から読み、正常・失敗・中断経路が受入条件と矛盾しないことを確認する。必要なら受入境界を再現するskill testを追加する。Codex Runtimeの会話受入はAIが実行せず、利用者が提示した会話記録、CLI argv記録、隔離DB事後状態だけを受け取って判定する。
6. `git diff --check`、対象Skillのfrontmatter検証、対象Skillと参照契約の相互参照を実行する。変更内容、同期元candidate ID、受入条件別証拠、実行command、利用者実行待ちの範囲、残存riskを報告する。

## Completion Criteria

- 承認済みの実装先だけを変更している。
- 同期順序、停止経路、本文不変、一時成果物境界を明示している。
- 公開CLI・DB・永続ledger・非同期実行を追加していない。
- 専用検証ブランチを更新した場合、review/audit済みsource candidateとの対応、専用branchだけへのcommit/push、利用者実行の未実施範囲を区別している。
- 実行した検証と未実行範囲を区別して報告している。
