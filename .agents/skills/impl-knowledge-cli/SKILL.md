---
name: impl-knowledge-cli
description: Knowledge CLIのGo実装を、実装前仕様照合・単一writerによる実装・独立コードレビュー・実装後の全体整合性監査へ分担し、証拠付きで完了させるオーケストレーター。FEAT-001または後続FeatureのCLI、SQLite migration、公開I/O、fixture、プロセス結合テストを変更する依頼で使う。
---

# Knowledge CLI Implementation Orchestrator

## Goal

承認済みのKnowledge CLI設計だけを、短い実装フィードバックと独立した品質ゲートで安全に完了する。オーケストレーター自身は製品コード・Planning成果物・Issueを変更しない。

## Required Roles and Contract

- `verify-knowledge-cli-spec`: 実装前の仕様照合
- `impl-knowledge-cli-implementation`: 唯一のwriter
- `review-knowledge-cli`: 実装diffの独立コードレビュー
- `audit-knowledge-cli-conformance`: コードレビュー合格後の全体整合性監査
- 全roleは[orchestration-contract.md](references/orchestration-contract.md)のauthority、baseline、candidate ID、入出力packet、verdictを使用する。

## Orchestration Procedure

0. `git status`と対象diffからbaselineと既存変更を記録し、製品コードのwriterを実装role一つに限定する。
1. verifierへIssue、Feature、baselineを渡す。`READY`だけを実装へ進め、`BLOCKED`は根拠を解消し、`NEEDS_HUMAN_DECISION`は人間へ返す。Decision変更が必要なら実装roleに書かせず、`planning-orchestrator`からartifact ownerの`feature-design`へ依頼し、承認後に再照合する。
2. implementerへVerification Packetを渡す。実装中は変更近傍のtestを優先し、candidate完成時に必須full gateを一度実行する。この段階ではIssueを更新しない。Implementation Reportが`READY_FOR_REVIEW`でなければreviewを開始しない。
3. 実装diff、candidate ID、source ID、Implementation Reportをreviewerへ渡す。`PASS`なら次へ進み、`BLOCKED`はfinding ID単位で同じimplementerへ戻す。修正後はcandidate IDと証拠を更新し、同じreviewerが再確認する。`NEEDS_SPEC_RECHECK`はverifierへ戻す。
4. **コードレビューが同一candidateへ`PASS`した後だけ**、新しい読み取り専用subagentでauditorを起動する。先入観を抑えるため、reviewの詳細findingではなくPASS事実と非blocking riskだけを渡す。`BLOCKED`は実装修正後にreviewから再開し、`NEEDS_SPEC_RECHECK`はverifierから再開する。
5. 同じ原因のfindingが2回続いた場合、局所修正を繰り返さず、前提・責務境界・test oracleを再評価する。未承認の公開契約で実装を修正しない。
6. 同一candidateへreviewとauditがともに`PASS`した後だけ、implementerへIssue Finalizationを依頼する。更新後にIssueを再読し、受入条件・検証結果・未完了項目が証拠と一致することを確認する。Issue更新失敗は実装成功と分けて報告する。

## Role Boundaries

- オーケストレーターだけがsubagentを起動し、packetとverdictを統合する。
- 事前調査・コードレビュー・最終整合性監査subagentはファイル、Issue、外部状態を変更しない。
- 製品コード、test、fixtureと最終Issue更新は実装subagentだけが行う。
- Planning成果物はPlanning workflowのartifact ownerだけが変更する。実装系roleはDecision案を含めて直接書かない。

## Use When

- FEAT-001または後続Featureで、Knowledge CLIのGoコード、SQLite persistence、migration、派生Index、CLI process integration testを実装・変更する。
- `cmd/knowledge`、`internal/application`、`internal/domain`、`internal/persistence/sqlite`、`testdata/fixtures`、`test/integration`の責務境界に関わる。

## Do Not Use When

- Planning成果物、要件、公開CLI契約、業務規則を新規に決定または変更する。
- Semantic Search、ローカルEmbedding、Vector Indexを実装する。これらはFEAT-006の承認済み設計が必要である。
- Codexによる検索戦略、Evidence価値、Knowledge Assessmentなど、CLIが持たない意味判断を実装する。

## Required Inputs

- `documents/AGENTS.md`
- `documents/docs/features/{feature_id}/implementation-handoff.yaml`
- 対象Featureの承認済み詳細設計とDecision
- `documents/docs/design/architecture.md`
- `documents/docs/design/cross-cutting-concerns.md`
- `documents/docs/design/technology-constraints.md`
- `documents/docs/design/decisions/DEC-ARCH-002.md`
- 実装対象に近い既存コードとtest

## Safety and Speed

- authority、入力packet、出力packet、verdictは共通契約へ一元化し、role間で同じ説明を再解釈しない。
- 編集中は変更近傍のtestを先に回し、失敗を早く局所化する。candidate完成前に毎回full suiteを回さない。
- candidate完成時の`gofmt`、`go test ./...`、`go vet ./...`、lint、対象coverage、必要なprocess integration testは省略しない。有効な最新証拠はreviewerが重複実行せず、欠落・陳腐化・疑義がある範囲だけ再実行する。
- test成功は観測した範囲の証拠であり、受入条件・公開I/O・migration・transaction・全体利用との対応付けを別に確認する。
- reviewerとauditorは独立させるが、最終監査はユーザー指定どおりコードレビュー合格後に直列実行する。

## Completion Criteria

- Verification Packetが`READY`である。
- 同一candidate IDへCode Review ReportとConformance Audit Reportがともに`PASS`である。
- 必須検証が同一candidateの最新証拠として成功している。
- 未承認の公開契約またはArchitecture/Product Decisionを導入していない。
- Issue Finalization後の再読結果が、受入条件・検証証拠・未完了項目と一致している。
