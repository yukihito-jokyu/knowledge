# FEAT-003 実装タスク

> **前提:** 利用者の「ではタスク分割に進んでください」（2026-08-15）を、更新済み詳細設計の明示承認およびTask分解開始指示として記録した。独立設計レビューは `pass` である。

## Design Readiness Audit

| 監査項目 | 結果 | 根拠 |
| --- | --- | --- |
| 利用者による詳細設計承認 | pass | `.ai/workflow/state.yaml` の `features.FEAT-003.human_design_review: approved` |
| 独立設計レビュー | pass | [design-review.md](design-review.md) の総合判定 `pass` |
| 利用者I/O、推奨、再調査規則 | pass | [design.md](design.md)「Behavioral Scenarios」「Re-evaluation」および DEC-FEAT-010〜011 |
| 記事取得の公開ネットワーク境界 | pass | [article-analysis-template.md](design/article-analysis-template.md) および [DEC-FEAT-012](decisions/DEC-FEAT-012.md) |
| Claim、Assessment Map、最終出力の形式 | pass | [article-analysis-template.md](design/article-analysis-template.md)、[reading-value-template.md](design/reading-value-template.md) |
| CLI・JSON・SQLite migration・公開設定 | not_applicable | 本FeatureはCodex workflow Skillを新設する。Knowledge CLIの公開契約・DB・設定は変更しない。 |
| Fixture・受入観点 | pass | [design.md](design.md)「Acceptance / Test Design」および [DEC-FEAT-012](decisions/DEC-FEAT-012.md) |

**結論:** pass。実装者が推奨基準、再調査上限、公開URL取得境界、公開CLI/JSON/永続化契約を再決定する未確定事項はない。

## 固定済みの実装対象と配置

実装は親リポジトリの `skills/reading-value/` に新設する。`documents/docs/features/FEAT-003/` は正規の設計・handoff資料であり、実行する機能本体ではない。

| 配置先（親リポジトリ基準） | 作成物 | 責務 |
| --- | --- | --- |
| `skills/reading-value/SKILL.md` | Reading Value Workflow | URL一件の受付からArticle Analysis、Claim単位のKnowledge Assessment、Markdown推奨までを統合する。 |
| `skills/reading-value/references/article-analysis.md` | 記事取得・分析契約 | 公開URL検証、検査済みIPへの接続固定、redirect、本文成立性、Article Claim、再調査時の照合を扱う。 |
| `skills/reading-value/references/reading-value.md` | 評価・推奨契約 | Assessment Map、再利用／無効化、推奨と出力を扱う。 |
| `skills/reading-value/references/verification.md` | 検証契約 | 通常・失敗・安全境界の再現可能な確認を定義する。 |

既存の `skills/knowledge-search/` はClaimごとのAssessmentとSearch Traceを返す依存先として利用する。記事取得、推奨判断、Reading Value Workflowを同SkillやKnowledge CLIの `cmd/`・`internal/`・SQLite migration・公開CLI JSONへ追加しない。

## タスク一覧

| ID | タイトル | 論理領域 | 依存 |
| --- | --- | --- | --- |
| TASK-003-01 | 安全な記事取得とArticle Analysisを実現する | Codex workflow: article acquisition / analysis | なし |
| TASK-003-02 | Claim評価統合・再調査・読書推奨を実現する | Codex workflow: assessment orchestration / recommendation | TASK-003-01、FEAT-002 implementation handoff |
| TASK-003-03 | Reading Value Workflowの受入・安全境界を検証する | verification | TASK-003-01、TASK-003-02 |

## TASK-003-01: 安全な記事取得とArticle Analysisを実現する

- **目的:** Codex会話から受ける技術記事URL一件を安全に取得・解析し、根拠位置を持つArticle Claimへ分解できるようにする。
- **関連要件:** REQ-001〜003、REQ-008、BR-001、BR-003、BR-010、NFR-001〜003、NFR-005、CON-002〜003、DEC-FEAT-010、DEC-FEAT-012。
- **論理領域:** Codex workflow: article acquisition / analysis。
- **作業内容:** `skills/reading-value/SKILL.md` の会話起点と `skills/reading-value/references/article-analysis.md` を新設する。記事URL、全redirect、名前解決結果を接続前に検査し、検査済みの公開到達可能IPだけへ接続先を固定して、接続時に再名前解決しない。保証できない取得能力は使用せず`article_unavailable`で停止する。本文からClaim ID、本文、役割、重要度、位置、support、および存在する場合のscope、target version、temporal contextを作成し、本文・Claim根拠の不足と非Claim要素の欠落を区別する。再調査時はClaim Reconciliationを作成する。
- **受入条件:**
  - Codex会話で技術記事の絶対HTTP(S) URLを一件だけ受け、専用CLI、API、Web UI、永続保存を追加しない。
  - 初期URLとredirect先は、HTTP/HTTPS、既定port、userinfoなし、公開到達可能な解決先だけを満たす場合にだけ接続を試みる。内部・ローカル・予約済み宛先、非既定port、資格情報付きURL、未検査redirectには接続しない。
  - 接続は検査済みIPへ固定し、接続時にhostnameを再名前解決しない。固定を保証できない場合は`article_unavailable`で停止する。
  - Cookie、認証情報、セッション、フォーム送信、状態変更を伴う取得を行わない。
  - Article ClaimはID、Claim、Role、Importance、Location、Supportを必須とし、scope、target version、temporal contextは該当時に明示する。
  - 本文またはClaim根拠の位置が取れない場合は部分Claim評価を開始せず`article_unavailable`とする。一方、Claim根拠に無関係な非本文要素の欠落は`content_limitations`として記録できる。
  - 技術記事でない、または評価可能なArticle Claimを一件も抽出できない場合は`article_not_evaluable`とし、空出力、要約、一般的な品質評価、読書推奨で代替しない。
  - 再調査時に、既存Claimをretained、changed、added、removedへ照合できる。
- **依存:** なし。
- **対象外／注記:** Claimごとの知識状態判定、Assessmentの生成、読書推奨はTASK-003-02で扱う。Knowledge CLIおよび既存`skills/knowledge-search/`への記事取得機能の追加は対象外。

## TASK-003-02: Claim評価統合・再調査・読書推奨を実現する

- **目的:** Article Claimごとに既存Knowledge Searchを利用してAssessment／Traceを対応付け、利用者に固有の読書価値を`read_full`、`read_selected`、`skip`として返せるようにする。
- **関連要件:** REQ-001、REQ-003、REQ-006〜008、REQ-019〜020、BR-001、BR-003、BR-009〜010、NFR-001〜003、NFR-005〜006、CON-002〜003、DEC-FEAT-010、DEC-FEAT-011、DEC-FEAT-012、DEC-FEAT-008、DEC-FEAT-009。
- **論理領域:** Codex workflow: assessment orchestration / recommendation。
- **作業内容:** `skills/reading-value/references/reading-value.md` を新設し、Workflowから各Article Claimを既存`skills/knowledge-search/`へ渡す。正常なAssessmentとSearch TraceをClaim IDごとに一対一でAssessment Mapへ対応付け、最終出力の評価根拠をMapのみに限定する。記事再調査は合計二回までとし、Claim、scope、target version、temporal contextがすべて同一のretained Claimだけ既存Assessmentを再利用する。changed／added Claimは再評価し、removed Claimは最終評価から除外する。再調査で続くKnowledge Assessmentはこの上限を消費しない。全Claimの評価が正常終了した場合だけ推奨を返す。
- **受入条件:**
  - 各Claimは既存Knowledge Searchの正常なAssessment一件と、そのTrace一件へAssessment Mapで対応付けられ、最終出力からClaim ID、Assessment、Traceを追跡できる。
  - `technical_failure`または`canceled`となったClaimが一件でもあれば、知識状態・推奨へ変換せず、Reading Value Assessmentを作成しない。
  - 利用者の既知／未知、推論可能性、Claimの重要度と支援根拠だけから、`read_full`、`read_selected`、`skip`を決定する。未観測を利用者の未知と断定しない。
  - `read_selected`では読むべきClaimと根拠位置を明示し、`skip`でも記事全体の品質評価や真偽断定をしない。
  - 記事再調査は初回後に合計二回までとし、再調査の理由とClaim Reconciliationを保持する。
  - 再利用判定はClaim本文、scope、target version、temporal contextの四属性すべてが同一の場合だけ行う。changed／addedは再評価、removedは除外する。
  - 評価、再調査、推奨の一時Markdown成果物は実行中の受渡しに限り、DB、公開JSON、公開設定、新しいCLI operationを追加しない。
- **依存:** TASK-003-01、FEAT-002 implementation handoff（`skills/knowledge-search/`のClaim入力、Assessment、Search Trace、技術失敗・中断契約）。
- **対象外／注記:** Knowledge Searchの探索戦略・状態分類・CLI呼出しBudgetの再設計、Knowledge CLIの変更、記事本文の永続アーカイブ、一般的な記事品質ランキングは対象外。

## TASK-003-03: Reading Value Workflowの受入・安全境界を検証する

- **目的:** 新設Workflowが通常・失敗・再調査・外部URL安全境界を、実装者以外も再現できる形で検証できるようにする。
- **関連要件:** REQ-001〜003、REQ-006〜008、REQ-019〜020、BR-001、BR-003、BR-009〜010、NFR-001〜003、NFR-005〜006、CON-002〜003、DEC-FEAT-010〜012。
- **論理領域:** verification。
- **作業内容:** `skills/reading-value/references/verification.md` を新設し、設計済みfixtureと観測oracleを整備する。会話起点、Article Claim、Assessment Map、三つの推奨、再調査時の再利用／無効化、記事取得の拒否・停止、Knowledge Searchの失敗伝播を確認する。検証はKnowledge CLIの公開I/Oを変更しない実行境界で行い、URL検査と実際の接続先固定を別々に観測する。
- **受入条件:**
  - 記事のClaim、Assessment Map、推奨の根拠がClaim IDからAssessment、Trace、位置まで追跡できる。
  - `read_full`、`read_selected`、`skip`それぞれで、決定根拠と禁止される断定が設計どおりに確認できる。
  - retained Claimの四属性が全て同一ならAssessmentを再利用し、changed／addedは再評価、removedは除外する。記事再調査三回目を開始しない。
  - 本文・Claim根拠の不足は`article_unavailable`、評価可能なClaimを抽出できない場合は`article_not_evaluable`とし、どちらもReading Value Assessmentを出力しない。Knowledge Searchの`technical_failure`、`canceled`も同様に推奨へ変換しない。
  - ローカル／内部／予約済み宛先、非HTTP(S)、非既定port、userinfo付きURL、不適格redirect、接続時の再名前解決、IP固定を保証できない取得能力が`article_unavailable`で停止することを確認する。
  - 正常記事ではCookie、認証、状態変更を伴わずに取得し、非Claim要素だけの欠落は`content_limitations`として出力できる。
  - Knowledge CLIの公開CLI JSON、SQLite migration、`cmd/`、`internal/`、公開設定に変更がないことを確認する。
- **依存:** TASK-003-01、TASK-003-02。
- **対象外／注記:** 外部記事サイトの可用性保証、任意URL取得一般のセキュリティ基盤、Knowledge Searchそのものの探索品質評価は対象外。

## Dependency Notes

- TASK-003-01は、新設Skillの入口と記事取得・分析の安全境界を固定するため最初に行う。
- TASK-003-02はTASK-003-01のArticle Claimを入力とし、既存`skills/knowledge-search/`の確定済み契約を消費する。
- TASK-003-03は、Workflow全体と外部取得境界を観測するため、両方の実装Taskに依存する。
