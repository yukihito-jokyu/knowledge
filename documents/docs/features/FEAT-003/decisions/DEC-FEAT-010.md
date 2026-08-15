# DEC-FEAT-010: URL評価の利用者入力境界

- **Status:** decided
- **Level:** L3
- **論点:** FEAT-003 のURL評価を、利用者がどの公開入口から開始し、結果をどの形で受け取るか。

## なぜ今決める必要があるか

REQ-001 は技術記事URLを起点に評価を依頼できることを求める。一方、初期設計は UI-001 として入力・結果提示の利用者インターフェースを未決定としている。入口を確定しなければ、入力バリデーション、取得失敗の返却、最終成果物の受渡し、および実装範囲を固定できない。

## 選択肢

### A. Codex会話からURLを直接受け付け、Markdown成果物を返す

利用者がCodexへURL評価を依頼し、Codex側のParent Orchestrationが実行する。最終結果は会話内で提示するReading Value Assessmentとする。

- **利点:** Issue #175のAI RuntimeがCodexであること、FEAT-002のCodex側ワークフロー、Codex間Markdown成果物契約と整合する。Knowledge CLIの公開契約を増やさない。
- **欠点:** 独立したCLI、HTTP API、専用GUIを提供しない。

### B. 専用CLIまたはJSON APIとして提供する

URLとオプションを機械的な公開I/Oで受け付け、構造化responseを返す。

- **利点:** 外部自動化や他の利用者インターフェースから利用しやすい。
- **欠点:** 新しい公開contract、error体系、設定・認証・実行環境の設計が必要であり、Initial Designと複数Featureへの影響を伴う。

### C. 専用Web UIとして提供する

URL入力画面と評価結果画面を提供する。

- **利点:** URL入力と読書箇所の視認性を専用に最適化できる。
- **欠点:** UI技術、配布、状態管理、セキュリティ境界を新設する必要があり、本Featureのワークフロー設計を超える。

## 推奨

**A** を推奨する。初期提供のCodex Runtimeという正規コンテキストと、既存のCodex側Knowledge Search契約に一致し、Knowledge CLIの公開JSON・SQLite契約を変更せずにREQ-001を満たすためである。BまたはCを選ぶ場合は、Initial Designへの変更要求と追加のL3設計が必要になる。

## 決定

2026-08-15に利用者が **A** を承認した。URL評価はCodex会話でURLを受け取り、Codex側ワークフローが実行したMarkdownのReading Value Assessmentを同じ会話で返す。専用CLI、JSON API、Web UIは初期提供に含めない。

## 影響

- A: FEAT-003内で、会話入力、記事取得、Markdown成果物、再調査、および最終提示を詳細化できる。
- B/C: `docs/design/**` の公開境界・技術制約を更新するChange Requestが必要となり、FEAT-003詳細設計はその承認まで停止する。
