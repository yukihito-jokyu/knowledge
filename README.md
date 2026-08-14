# knowledge

技術記事を「自分にとって今読む価値があるか」という観点で評価するための、個人知識基盤です。

記事の主張を、利用者自身の知識を裏付ける Evidence（発言、コード、訂正など）と照合し、読むべきか・どこを読むべきかを説明できる状態を目指します。一般的な記事の品質評価や、AI 生成コンテンツの判定は目的に含みません。

> 現在は設計・実装計画の段階です。実行可能な Knowledge CLI やアプリケーションはまだリポジトリに含まれていません。

## 目指す体験

利用者が技術記事の URL を渡すと、Codex が記事を複数の Claim（検証可能な主張）へ分解します。各 Claim を個人知識と照合し、次のいずれかを理由とともに提案します。

- `read_full`: 記事全体を読む価値が高い
- `read_selected`: 価値のある箇所に絞って読む
- `skip`: 現時点では読む必要が薄い

評価では、既知・部分的に既知・推測可能・矛盾・古い可能性・根拠なし・不確実を区別します。検索結果がないことだけで「知らない」と断定せず、Evidence を根拠に扱うことを重視します。

## 仕組み

```text
技術記事 URL / 会話・作業エピソード
                │
                ▼
           Codex Orchestration
  ┌─────────────┼──────────────────┐
  │ Article     │ Knowledge Search │ Reading Value
  │ Analysis    │ / Update         │
  └─────────────┴──────────────────┘
                │ JSON
                ▼
          Knowledge CLI（予定）
                │
                ▼
   ローカル SQLite の Knowledge Store
```

Knowledge CLI は保存・検索・取得・更新を決定論的に実行する境界です。Claim の意味解釈、検索戦略、知識状態の判定、読書価値の判断は Codex 側が担います。

## 計画中の機能

| Feature | 内容 | 状態 |
| --- | --- | --- |
| FEAT-001 | Assertion・Evidence・Concept・Relation・時点情報を保存し、JSON CLI で検索・更新する | 計画済み |
| FEAT-002 | Claim ごとに根拠を探索して知識状態を評価する | 計画済み |
| FEAT-003 | URL 記事の読書価値と推奨箇所を説明する | 計画済み |
| FEAT-004 | 会話・作業から Evidence 候補を抽出し、知識を更新する | 計画済み |
| FEAT-005 | 検索・更新・評価をシナリオで検証する | 計画済み |
| FEAT-006 | ローカル Semantic Search を追加する | 初期提供後 |

初期の検索は字句検索を中心とし、Semantic Search は後続の機能として扱います。

## 設計上の原則

- Evidence と、そこから導く知識状態を混同しない
- 過去の知識や訂正を削除せず、revision と置換関係で履歴を保つ
- 個人知識はローカル SQLite にのみ保存し、初期提供では同期・共有・リモート接続を扱わない
- CLI の公開境界は名前付き option と JSON 出力に固定する
- AI による意味判断と、CLI による決定論的な永続化を分離する

## 実装予定の Knowledge CLI

FEAT-001 では Go 1.26 以上、`modernc.org/sqlite`、埋め込み migration を採用する予定です。CLI は以下の 11 操作を提供する設計です。

```text
search-text            search-concept          search-related
get                    get-evidence            search-contradictions
search-temporal        create                  attach-evidence
revise                 supersede
```

成功時は stdout に JSON を 1 件、失敗時は stderr に JSON のエラーを 1 件出力します。具体的な option、JSON schema、終了コードは実装前にすでに詳細設計で固定されています。

## 資料の読み方

- [システムコンテキスト](documents/docs/requirements/system-context.md): プロダクトの目的、範囲、用語
- [機能要件](documents/docs/requirements/requirements.md): 正規の要件一覧
- [Feature Map](documents/docs/planning/feature-map.md): 機能ごとの目的、依存関係、進行状況
- [初期アーキテクチャ](documents/docs/design/architecture.md): Codex、CLI、SQLite の責務境界
- [FEAT-001 詳細設計](documents/docs/features/FEAT-001/design.md): 最初の Knowledge CLI の公開契約
- [コマンド一覧](documents/docs/features/FEAT-001/design/command-catalog.md): 各 CLI 操作の設計
- [実装タスク](documents/docs/features/FEAT-001/tasks.md): FEAT-001 の実装順序と受入条件

`documents/docs/` が要件・設計・handoff の正規資料です。製品コードの実装はリポジトリ直下で行い、承認済みの公開 CLI・JSON・保存仕様は、実装の都合で変更しません。

## 開発に参加する場合

実装を始める際は、対象 Feature の詳細設計と実装タスクを先に確認してください。Go 実装の共通検証は次のとおりです。

```bash
gofmt -w <対象ファイル>
go test ./...
go vet ./...
```

公開 CLI の JSON、標準出力・標準エラー、終了コード、SQLite migration は、実行プロセス境界で検証します。設計の変更が必要な場合は、先に [Decision Policy](documents/.ai/workflow/decision-policy.md) に従って判断を記録します。
