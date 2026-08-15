# FEAT-003 要件

## 目的

利用者がCodex会話で指定した技術記事URLを、Codex側のParent OrchestrationがArticle Claimへ分解し、FEAT-002のKnowledge Assessmentと統合して、利用者固有の読書推奨と読むべき箇所を自然言語で返す。

## 実装対象と配置

実装はKnowledge CLIのGoコードではなく、親リポジトリのCodex Workflow Skillとして行う。新設する配置先は `skills/reading-value/` とし、`documents/` 配下はこの機能の設計資料だけを置く。

| 作成物 | 配置先 | 実装する責務 |
| --- | --- | --- |
| Reading Value Workflow | `skills/reading-value/SKILL.md` | Codex会話のURL入力を受け、記事取得、Article Claim作成、FEAT-002呼出し、Assessment Map管理、最大2回の再調査、最終Markdown返却を一つのワークフローとして制御する。 |
| 記事解析・取得リファレンス | `skills/reading-value/references/article-analysis.md` | URL公開性検査、検査済みIPへの接続先固定、リダイレクト、本文取得可否、Claim分解と再照合を規定する。 |
| 推奨統合リファレンス | `skills/reading-value/references/reading-value.md` | Assessment Map、`read_full`／`read_selected`／`skip`の決定、再調査、Reading Value Assessmentの出力を規定する。 |
| 検証リファレンス | `skills/reading-value/references/verification.md` | 正常系、再調査、取得失敗、公開外URL、DNS再解決、部分本文、Knowledge Search失敗の確認手順を規定する。 |

既存の `skills/knowledge-search/` はClaimごとのKnowledge Assessmentを提供する依存先として呼び出す。Knowledge CLIの `cmd/`、`internal/`、SQLite migration、公開CLI JSONには変更を加えない。

## 含む要件

- REQ-001、REQ-002、REQ-003、REQ-006、REQ-007、REQ-008、REQ-019、REQ-020。
- BR-001、BR-003、BR-009、BR-010。
- NFR-001、NFR-002、NFR-003、NFR-005、NFR-006。
- CON-002、CON-003、CON-006、DEC-REQ-001、DEC-FEAT-010、DEC-FEAT-011。

## Feature Acceptance

1. 利用者がCodex会話で技術記事のHTTPまたはHTTPS URLを一件指定すると、評価を開始できる。
2. 記事本文をClaim単位に分解し、各Claimに本文、`core`／`supporting`／`example`／`background`／`incidental`、`high`／`medium`／`low`、記事内位置、Supportを記録する。
3. 各ClaimをFEAT-002へ渡し、Knowledge Assessmentを再判定せずに使用する。
4. `read_full`、`read_selected`、`skip`の一つを、数値スコアではなく理由を伴う自然言語で返す。
5. `read_selected`では、読む箇所、読む必要が薄い箇所、得られる認識利得、および根拠の信頼性を説明する。
6. Reading Valueが重要な不足を検出した場合、Article AnalysisまたはKnowledge Searchへ最大2回まで限定再調査を依頼できる。
7. 記事取得・解析・Knowledge Searchの技術失敗または中断を知識状態や読書推奨へ変換しない。

## 範囲外

- 専用CLI、JSON API、Web UI、記事評価の永続化、記事からのKnowledge Evidence自動登録。
- Knowledge CLIのoperation、JSON wire schema、SQLite schema、公開設定の追加・変更。
- Knowledge Assessmentの検索戦略、Evidence強度、7状態の再判定（FEAT-002）。
- 書籍、動画、Release Note、公式ドキュメント、論文、ニュース、カンファレンス資料の評価（ASM-001）。
- 一般的な記事品質、AI生成かどうか、利用者の熟達度の評価。
