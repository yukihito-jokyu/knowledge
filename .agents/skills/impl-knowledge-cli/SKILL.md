---
name: impl-knowledge-cli
description: Knowledge CLIのGo実装におけるCLI境界、application/domain/persistenceの責務分離、SQLite migration・transaction、process境界の統合検証を実装する際に使う。Semantic Searchや、承認済み公開契約の変更には使わない。
---

# Knowledge CLI Go Implementation

## Goal

承認済みのKnowledge CLI設計を、Go 1.26以上、`modernc.org/sqlite`、埋込みmigration、隔離SQLite DBを使うCLI process検証で、一貫して実装・検証する。

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

## Project Conventions

- Knowledge CLIはGo 1.26以上の単一moduleであり、SQLite driverはCGOを必要としない`modernc.org/sqlite`である。
- `cmd/knowledge`は公開CLI境界と依存の組立てだけを担い、SQL・migration・ドメイン規則を持たない。
- `internal/application`はtransport・SQLite driverから独立し、`internal/domain`は外部I/Oに依存しない。
- `internal/persistence/sqlite`がSQLite接続、SQL、migration、派生字句Indexを所有し、applicationまたはCLI entryをimportしない。
- migration SQL assetはGo標準`embed`で同梱する。既適用migrationを編集せず、必要な変更には後続versionを追加する。
- 公開CLIのoption、JSON、stdout/stderr、exit code、保存先、設定、運用仕様を追加・変更しない。
- Goおよびmigration SQLの説明コメントは、識別子・コンパイラ指示を除いて日本語で端的に記述し、英語の説明語を残さない。
- Goの複数要素の配列・sliceリテラルは、各要素を1行ずつに記述する。
- Goの複合リテラルは、各フィールドを1行ずつに記述する。`scripts/check_composite_literal_layout.sh`で違反を検査する。
- SQLite adapterのSQL文は、使う関数の直前に個別のパッケージレベル定数として置き、複数行で整形する。関数内へ直接記述しない。
- SQLへ値を渡すときは`?`プレースホルダーを使う。固定個数の`IN`条件を文字列置換で組み立てない。
- unit testはテーブル駆動で記述し、Skill配下の`scripts/check_test_coverage.sh`で全対象packageのstatement coverage 100%を確認する。
- 実行環境上再現不能な分岐がある場合は、対象コードへ日本語の理由をコメントし、テスト可能な内部依存を注入してカバーする。

## File Placement Rules

- CLI process entry: `cmd/knowledge/`
- application: `internal/application/`
- domainモデル・不変条件・永続化port: `internal/domain/`
- SQLite adapter、SQL、migration runner、派生Index: `internal/persistence/sqlite/`
- 不変migration asset: `internal/persistence/sqlite/migrations/`
- 再現可能なfixture: `testdata/fixtures/`
- CLI／SQLiteのprocess境界integration test: `test/integration/`

## Autonomous Decisions

- 承認済み責務境界の内側での非公開型、関数、package分割、test helper、fixture構成を選択できる。
- 既存コードと設計に従い、最小の内部実装を選択できる。
- 不明な公開契約は推測で補わず、設計資料を確認する。

## Escalate When

- 新しい公開CLI option、JSON field、error code、exit code、設定、保存先、運用契約が必要になる。
- 新しいlibrary、framework、SQLite driver、Go version、またはArchitecture境界の変更が必要になる。
- 詳細設計・Decision・実行可能な既存コードが矛盾し、局所的な実装判断で解決できない。

## Procedure

1. 対象Issueと対応する`implementation-handoff.yaml`を読み、依存Taskの完了を確認する。
2. 対象操作の詳細設計から、option、JSON、validation、SQL、transaction、fixture・受入条件を抽出する。
3. 既存の近接コードを確認し、責務境界を越えない最小の実装を行う。
4. unit testと、必要なCLI process境界のintegration testを追加または更新する。testごとに隔離SQLite DBを使い、公開DB指定optionや設定を追加しない。
   - 利用者が確認すべきCLIコマンドの棚卸しと自動検証を求めた場合、入力引数、期待stdout、期待stderr JSON、exit codeを`testdata/fixtures/`へ置き、`test/integration/`から実バイナリを起動して検証する。fixtureを公開運用コマンドや別の利用者向け実行ファイルとして扱わない。
   - Taskfileは開発用の実行入口であり、fixtureや期待結果の正本にしない。Taskfileへのtest task追加・維持はユーザーの指定に従う。
   - 未実装の後続操作で成功結果を検証できない場合は、そのケースの現時点の期待結果と、成功契約を検証する後続Taskを明確に分ける。
5. migrationは初回適用、再実行、失敗時rollbackを確認する。mutationはcommit後だけ成功responseを返し、失敗時は部分更新を残さない。
6. 検証を実行し、設計との不一致はPlanning成果物を勝手に変えず報告する。

## Validation

- `gofmt` を変更したGoファイルに実行する。
- `go test ./...`
- `go vet ./...`
- `task lint`
- `.agents/skills/impl-knowledge-cli/scripts/check_test_coverage.sh`
- 変更内容に応じて、対象のCLI process integration testとmigrationの初回・再実行・rollback検証を実行する。

## Completion Criteria

- 対象Taskの受入条件と操作別設計を満たす。
- application/domain/persistence/CLIの依存方向を守る。
- migration、transaction、公開I/Oを該当するtestで検証する。
- 利用者確認コマンドを求められた場合、fixtureとprocess境界testが引数、stdout、stderr JSON、exit codeを対応付けている。
- 未承認の公開契約またはArchitecture/Product Decisionを導入しない。
