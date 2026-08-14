# 初期アーキテクチャ

## 設計ドライバー

- URL 評価では、記事 Claim とユーザー知識の差分を Evidence 付きで説明する必要がある。
- `no_evidence` と未知、`inferable` と既知、閲覧と知識獲得を厳密に分離する必要がある。
- Codex の意味判断を、再現可能な Knowledge CLI の保存・検索・更新から分離する必要がある。
- 継続更新・履歴保持・Search Trace と Scenario A〜J による検証可能性が必要である。

## 最小責務境界

```text
利用者
  │ URL評価依頼 / エピソード入力
  ▼
Codex Orchestration
  ├─ Article Analysis: 記事を Claim へ分解
  ├─ Knowledge Search: 検索戦略と Evidence に基づく Assessment
  ├─ Reading Value: Assessment を読書推奨へ統合
  ├─ Knowledge Acquisition: エピソードから Evidence 候補を抽出
  └─ Knowledge Update: 候補と既存知識を照合して更新を判断
             │ JSON
             ▼
Knowledge CLI
  ├─ `cmd/knowledge`: process entry（option受領、標準入出力、終了コード）
  ├─ `internal/application`: 操作の組立て、transaction境界、公開契約の写像
  ├─ `internal/domain`: 知識モデル・不変条件・永続化port
  ├─ `internal/persistence/sqlite`: SQLite adapter、migration、派生Index
  └─ 保存・検索・取得・更新・Index管理 → Knowledge Store
```

## 依存方向と禁止事項

- Codex 側の各ワークフローは Knowledge CLI の論理プリミティブへ依存する。
- Knowledge CLI は Codex のプロンプト、Claim の意味、検索継続、known 判定、読書価値判定へ依存してはならない。
- Knowledge Store は Knowledge CLI だけが直接変更する。
- `cmd/knowledge` は公開CLI境界だけを担い、最初の`os.Interrupt`から要求Contextを作成して下位層へ伝播する。ドメイン規則、SQL、migrationを持たない。
- `internal/application` は CLI transport から独立して操作を組み立てる。domainが定義する永続化portを通じて必要な永続化操作を行い、SQLite固有の型・SQL・driverを参照しない。
- `internal/domain` は Knowledge CLI の意味的に中立なモデルと不変条件を担う。`database/sql`、CLI引数、JSONの送受信、SQLite固有の実装へ依存してはならない。
- `internal/persistence/sqlite` は domain が定義する永続化portを実装し、SQLite接続、SQL、migration、派生Indexを所有する。applicationまたはCLI entryをimportしてはならない。
- composition root は `cmd/knowledge` とし、applicationとSQLite adapterの接続をここで行う。依存方向は `cmd/knowledge → application → domain`、`cmd/knowledge → persistence/sqlite → domain` とする。
- Article Analysis と Reading Value は Knowledge State を直接更新してはならない。
- Knowledge Acquisition は候補を抽出し、Knowledge Update が更新判断を所有する。

## Go module と配置規約

Knowledge CLI は単一の Go module として管理する。`go.mod` の最低Go versionは1.26とし、SQLite driverは `modernc.org/sqlite` とする（[DEC-ARCH-002](decisions/DEC-ARCH-002.md)）。module pathは実装開始時にリポジトリの配布先に合わせて定める内部識別子であり、公開CLI契約ではない。

| 配置 | 責務 | 境界 |
| --- | --- | --- |
| `cmd/knowledge/` | CLI process entry、依存の組立て、既定Storeの解決、要求Context、stdout/stderr・exit codeの受渡し | `os.Interrupt`から要求Contextを生成し、DEC-FEAT-005の既定Storeを解決してapplicationとSQLite adapterを接続する。response開始前にContext終了を観測した場合はJSONを出力せず終了コード130とする。SQL・ドメイン判断を持たない。 |
| `internal/application/` | 11操作の実行順、入力からdomain値への変換、transactionを使う操作の組立て | CLI transportとSQLite driverを参照しない。公開fieldの正規仕様はFEAT-001設計に従う。 |
| `internal/domain/` | Assertion、Evidence、Relation等の意味的に中立な値、不変条件、永続化port | 外部I/Oへ依存しない。Codexの意味判断・知識状態評価を実装しない。 |
| `internal/persistence/sqlite/` | SQLite adapter、SQL、接続初期化、migration実行、派生字句Index | domainの永続化portを実装する。公開CLIのoption／JSONを定義しない。 |
| `internal/persistence/sqlite/migrations/` | 連番付きの不変SQL migration asset | adapterだけが読み込む内部asset。保存先・外部操作・運用設定を新設しない。 |
| `testdata/fixtures/` | CLI入力、期待JSON、seedデータ等の再現可能なfixture | プロダクトコードから読み込まない。公開契約の検証データとして扱う。 |
| `test/integration/` | process境界を通るCLI／SQLite integration test | `go test ./...` で実行する。unit testを置き換えない。 |

migration assetは Go 標準の `embed` を用いて実行バイナリへ同梱する。実行時に外部migrationディレクトリ、保存先option、設定ファイルを要求しない。通常ビルドの既定StoreはDEC-FEAT-005に従い、`os.UserConfigDir()/knowledge/knowledge.db` とする。適用規則（連番、単一transaction、再実行、破損schemaの扱い）はFEAT-001の [database-schema.md](../features/FEAT-001/design/database-schema.md#migration) を唯一の正規仕様とする。

integration testは、OS標準のユーザー設定ディレクトリ環境をテストごとの一時領域へ隔離し、通常のcompositionでSQLite adapterを接続する。公開CLIのoption、設定、環境変数としてDB指定を追加せず、利用者の既存Storeにも接続しない。

要求Contextはprocess entryからcomposition root、application、domainの永続化port、SQLite adapterへ同一のまま渡す。下位層が`context.Background()`で要求を作り直すことは許可しない。CLI境界はparse、Context非対応の保存先解決／ディレクトリ作成、実行結果の各境界とresponse開始直前でContext終了を確認し、response開始前なら結果・error分類より終了コード130と無出力を優先する。Context終了時、SQLite adapterは読取・migration・transactionへContext対応APIを用い、mutation transactionをcommitしない。response開始後の既出力は取り消さない。公開上の割込み規約はDEC-FEAT-006に従う。

process境界testは、`integrationtest` build tagにだけ含める非公開同期gateでSQLiteのmigration、read、mutation処理開始を確認してから割込みを送る。このgateは通常ビルドに含めず、通常composition・既定Store解決を差し替えず、公開CLI option・公開環境変数・設定を追加しない。

## 外部境界

- 記事本文の取得・URL 解析の具体的方式は FEAT-003 の詳細設計で決定する。
- Codex と Knowledge CLI の境界は JSON とする。
- Codex ワークフロー間は必須セクションを持つ Markdown 成果物を基本とする。
- Knowledge Store はローカル専用の埋め込み SQLite とし、Knowledge CLI を唯一のアクセス経路とする。

## 保留中の横断判断

- UI-001: URL 評価の利用者インターフェース
- UI-002: Conversation / Task Episode の取り込み境界
