# 未解決事項

## L2 Assumptions

### ASM-001: 初期リリースの対象コンテンツ

初期の URL 評価は Issue #175 が明示する技術記事を対象とし、書籍・動画・Release Note・公式ドキュメント・論文・ニュース・カンファレンス資料への拡張は将来範囲とする。

- **根拠:** Primary Use Case および §50
- **影響:** 要件・Feature Planning
- **状態:** assumed

### ASM-002: 検索 Budget の定義時期

検索回数、Top-K、Relation Depth 等の具体的な既定値は、実装可能な設計・評価時に決定する。現時点では停止可能で設定可能であることだけを要件とする。

- **根拠:** Issue #175 §16
- **影響:** 初期設計・Feature設計
- **状態:** assumed

## 解決済みDecision

### DEC-REQ-001: Semantic Search の段階的提供

初期提供は字句検索を中心とし、Semantic Search はローカル Embedding と Vector Index を採用する後続移行で追加する。詳細は [DEC-REQ-001](decisions/DEC-REQ-001.md) を参照。

## 未決定（後続設計で解決）

### UI-001: URL 評価の利用者インターフェース

URL の入力、評価結果、読書箇所の提示をどの利用者インターフェースで提供するかは未定。既定のシステム境界・論理契約を変更しない範囲で初期設計にて決定する。

- **Decision Level:** L3
- **影響:** 初期設計、URL 評価 Feature

### UI-002: 会話・作業エピソードの取り込み境界

Knowledge Acquisition の入力となる Conversation / Task Episode を、どの Codex 連携・明示操作・保存単位から受け取るかは未定。

- **Decision Level:** L3
- **影響:** 初期設計、Knowledge Acquisition Feature

### UI-003: Knowledge CLI の物理実装

データベース、Embedding Engine、Index Library、Migration 方式および検索パラメータは後続の CLI 詳細設計で決定する。

- **Decision Level:** L2
- **影響:** 初期設計、Knowledge Store Feature
