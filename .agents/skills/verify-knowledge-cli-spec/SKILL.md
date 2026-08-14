---
name: verify-knowledge-cli-spec
description: Knowledge CLIのIssue、要件、詳細設計、承認状態付きDecision、handoff、依存Task、現行コードを実装前に読み取り専用で照合し、受入条件別の実装可能範囲・検証oracle・未承認公開契約を判定する調査Skill。実装開始前と契約問題の検出後にオーケストレーターから使う。
---

# Knowledge CLI Specification Verifier

## Read-only Boundary

ファイル、Issue、workflow state、外部状態を変更しない。[orchestration-contract.md](../impl-knowledge-cli/references/orchestration-contract.md)を読み、Verification Packetを返す。

## Source Authority

1. 適用される`AGENTS.md`とworkflow policyは作業方法・Ownershipを規定する。
2. 要件・Business Ruleと、人間承認済みDecisionはproduct intentと契約を規定する。Decisionのstatus、承認記録、supersedes関係を確認する。
3. 詳細設計とoperation資料は承認済み契約を具体化する。上位根拠を黙って変更できない。
4. Taskとhandoffは今回の実装範囲を選ぶ。新しい公開契約を承認できない。
5. Issueは進行と受入条件のtrackerであり、単独では仕様の正規根拠にならない。
6. 現行コードとtestは実態の証拠であり、仕様との不一致を正当化しない。

文書の新しさだけで優先順位を決めない。矛盾を解消する明示的なDecisionまたはsupersedes参照がなければ、矛盾として報告する。

## Procedure

1. `documents/AGENTS.md`、Decision Policy、artifact map、workflow state、対象Featureの要件・設計・Decision・handoff・Task、Issueを読む。
2. handoffのdesign readiness、独立design review、人間承認、依存Taskを確認する。実装開始条件を満たさない場合は`BLOCKED`にする。
3. Issueの各受入条件へ安定したIDを付け、正規根拠、対象operation、実装surface、公開I/O、最も強い検証oracle、対象外を対応付ける。
4. 現行コード、近接test、migration version、fixtureを読み、`implemented`、`missing`、`contradictory`、`not_in_scope`を分ける。
5. 新規または変更されるoption、JSON field、error/exit code、保存先、設定、運用、schema互換性、Architecture境界を列挙し、承認済みDecisionへ紐付ける。
6. 根拠がない公開契約、L3/L4、Business Rule衝突は`NEEDS_HUMAN_DECISION`にする。資料欠落、未完了依存、機械的に解消可能な不整合は`BLOCKED`にする。
7. reviewerが重点確認すべきfailure modeと、auditorが確認すべき既存利用経路を引継ぎへ記録する。

## Required Checks

- CLIの引数、型、必須性、相互排他、stdout/stderr、JSON schema、error code、exit code
- SQLite schema、migration前後条件、rollback、既存DB互換性
- transaction、cancel、read-only性、履歴・派生Indexの不変条件
- fixtureとprocess testで受入条件を観測できるか
- out-of-scope、後続Feature、未実装操作との境界
- Issue checkboxと実際の証拠が一致しているか

## Output

共通contractのVerification Packetを使う。判定は`READY`、`BLOCKED`、`NEEDS_HUMAN_DECISION`のいずれか1つとし、各受入条件に正規sourceと検証oracleがなければ`READY`にしない。不明点を一般的な質問へ変換せず、調査済み根拠と判断が必要な一点を示す。
