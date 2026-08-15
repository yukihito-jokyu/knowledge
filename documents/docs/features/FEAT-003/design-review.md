# FEAT-003 詳細設計レビュー

- **対象:** FEAT-003: URL記事の読書価値評価
- **レビュー種別:** 独立・読み取り専用レビュー
- **判定:** pass
- **次のゲート:** 利用者による詳細設計レビュー。承認後にのみTask分割へ進める。

## 正規ソース台帳

| ソース | 確定性 | 適用範囲 |
| --- | --- | --- |
| Issue #175（正規要件資料に転記された§1、§7〜8、§27〜38、§43、§50、§52） | 原典 | URL起点の評価、Claim、知識評価、読書推奨、再調査、検証可能性 |
| `docs/requirements/requirements.md` REQ-001〜003、REQ-006〜008、REQ-019〜020 | confirmed | FEAT-003の機能境界 |
| `docs/requirements/business-rules.md` BR-001、BR-003、BR-009〜010 | confirmed | 利用者固有の価値、推論可能・未観測の扱い、CLI責務 |
| `docs/requirements/nonfunctional-requirements.md` NFR-001〜003、NFR-005〜006 | confirmed | 追跡可能性、有限性、fixture評価、論理契約互換性 |
| `docs/requirements/constraints.md` CON-002〜003、CON-006 | confirmed | Codex Runtime、CLI JSON境界、Markdown成果物、未確定物理実装 |
| `docs/design/architecture.md`、`domain-boundaries.md`、`data-ownership.md`、`cross-cutting-concerns.md` | approved initial design | Codex／CLI責務、成果物所有、外部失敗の伝播 |
| DEC-FEAT-010 | L3 decided | Codex会話を入口、会話内Markdownを返す |
| DEC-FEAT-011、DEC-FEAT-012 | L2 decided | 再調査上限、公開記事取得の安全境界 |

## 契約根拠トレーサビリティ

| 境界 | 設計上の契約 | 根拠 | 判定 |
| --- | --- | --- | --- |
| 利用者I/O | Codex会話で一件の技術記事URLを受け、Reading Value AssessmentをMarkdownで返す。専用CLI/API/Web UIは作らない。 | REQ-001、REQ-008、CON-001〜003、DEC-FEAT-010 | explicit |
| 実装配置 | 親リポジトリの `skills/reading-value/` にWorkflow本体と3つのreferenceを新設する。 | 利用者の明示要求、親リポジトリ`AGENTS.md`のCodex workflow配置規約、FEAT-003設計 | explicit |
| 依存 | `skills/knowledge-search/` をClaim単位で利用し、同SkillのAssessment／Trace契約を消費する。 | REQ-019、FEAT-002契約、`skills/knowledge-search/SKILL.md` | derived |
| CLI・永続化 | `cmd/`、`internal/`、SQLite migration、公開CLI JSON、設定を変更しない。 | BR-010、CON-002〜003、NFR-005、初期設計 | explicit |
| 記事取得 | HTTP(S)・既定port・公開到達可能IPのみ。初期URLと全redirectを接続前に検査し、検査済みIPへ固定する。Cookie等を送らない。 | 初期設計がFEAT-003へ委譲、DEC-FEAT-012 | decision |
| Article Claim | 本文、役割、重要度、位置、Supportを必須として、実行内Markdownに保持する。 | REQ-002〜003、CON-003 | explicit |
| Assessment／Trace追跡 | Claimごとに正常AssessmentとTrace参照を一対一でAssessment Mapへ対応付け、最終出力から追跡可能にする。 | NFR-001〜003、FEAT-002成果物契約 | derived |
| 再調査 | 合計2回まで。記事再調査で同一のClaim本文・Scope・version・時点だけ既存Assessmentを再利用し、変更／追加は再評価、削除は除外する。 | REQ-020、NFR-002、DEC-FEAT-011 | decision |
| 失敗 | 記事取得・解析失敗、Knowledge Searchの`technical_failure`／`canceled`を知識状態や推奨に変換せず伝播する。 | BR-002〜003、初期設計、FEAT-002契約 | derived |

## 資料間整合の確認

### 実装対象と配置

親リポジトリを俯瞰して確認した。`skills/knowledge-search/SKILL.md`と同梱referenceが存在し、`skills/reading-value/`は未作成である。新設先、作成物、責務、既存依存先が[FEAT-003設計](design.md#implementation-deliverables--placement)と[FEAT-003要件](requirements.md#実装対象と配置)で一致する。

この配置は、親リポジトリ`AGENTS.md`の「Knowledgeプロダクトに付随するCodex workflow Skillは`skills/<skill-name>/`に置く」という規約に適合する。Knowledge CLI実装領域と分離されており、既存`knowledge-search`へ記事取得や推奨判断を混入させない責務境界も明確である。

### 成果物・失敗・再調査

- Article Analysis templateのClaim必須属性と、設計の入力・責務表がREQ-002〜003と一致する。
- Reading Value templateの`Assessment Map`は、FEAT-002の正常Assessmentが持つ`Trace Reference`と矛盾しない。`technical_failure`／`canceled`にはAssessmentを作らない既存契約とも整合する。
- `read_full`、`read_selected`、`skip`の決定条件は、未知Claim数による機械判定を禁じるBR-009、および`inferable`／`no_evidence`の扱いに関するBR-002〜003と整合する。
- 再調査後のClaim照合規則はArticle Analysis template、DEC-FEAT-011、設計のAssessment Map規則で同じ四属性を比較対象としており、再利用・無効化・除外に矛盾はない。
- 記事取得失敗時に部分本文から推奨しない規則と、Claim根拠と無関係な要素のみの欠落を記録して評価可能とする規則は、template、error flow、acceptanceで一致する。

### セキュリティ・公開境界

DEC-FEAT-012の全redirectに対する接続前検査、検査済みIP固定、再名前解決禁止、認証情報・Cookie不送信が、実装deliverable、edge case、acceptanceに反映されている。HTTP(S)だけでは内部接続を防げないという外部境界のリスクを、公開I/OやCLI設定の追加なしにFeature内L2決定として閉じている。

実装時に当該固定を保証できる記事取得能力がなければ`article_unavailable`として停止するため、安全性を未検証の能力に委ねて成功扱いにしない。これはDEC-FEAT-012の明示契約である。

## 差し戻し事項

なし。

## 留意事項（実装時の確認対象）

- `skills/reading-value/`は現時点で未作成であり、設計どおりの新設対象である。
- URL取得能力が検査済みIPへの接続先固定を実際に保証できるかは、実装時に検証する。保証できない能力を選択した場合は、DEC-FEAT-012に従って対象記事を`article_unavailable`にし、取得を続行してはならない。
- 上記は未決Decisionや設計不足ではなく、設計済みの失敗条件を実装・検証するための確認事項である。

## 総合判定

**pass**。L3 DecisionはDEC-FEAT-010で利用者承認済みであり、公開CLI/API、永続化、設定、Knowledge CLI契約を暗黙に増やしていない。実装者へ意味判断、再調査、外部安全境界を再設計させる未解決事項はない。次のゲートは利用者による詳細設計レビューであり、このレビューだけでimplementation handoffを作成してはならない。
