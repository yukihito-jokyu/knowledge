# DEC-FEAT-023: ローカル Embedding モデルと Vector Index の提供方式

- **Status:** proposed
- **Level:** L3
- **対象 Feature:** FEAT-006

## 判断が必要な機能と現在の状態

FEAT-006 は、既存 Assertion の現行正規化本文をローカルで Embedding 化し、ローカル Vector Index により `search-semantic` の候補を返す機能である。

既に確定しているのは「外部 API ではなく、ローカル Embedding モデルとローカル Vector Index を使う」ことである（DEC-REQ-001、DEC-FEAT-001）。一方、モデルをどのように利用者の端末へ提供・更新するか、またどの Vector Index 実装を同じ配布物で管理するかは未確定である。

## 既知の事実

- Knowledge CLI は単一 Go module、純 Go SQLite driver、埋込み migration、OS標準の既定 Store を使う。外部 migration ディレクトリ、保存先 option、設定ファイルは追加しない。
- 既存 Store は安定した Assertion ID、現行正規化本文、Scope、Concept、immutable な Evidence・revision・Relation を保持している。
- 現在の Go module に Embedding runtime、モデル asset、Vector Index dependency は存在しない。
- Semantic Search の追加は既存11操作の公開契約を破壊せず、正規 Assertion から派生データを再構築する必要がある。

## 未確定事項と AI だけで確定できない理由

モデルの配布サイズ、初回利用時のネットワーク要件、更新責任、端末要件、利用を許容するライセンスは、利用者の運用方針とリリース方針に依存する。これらは SQLite schema、model version の保持、Index 互換性および障害時の利用体験を変えるため、実装上の局所判断にはできない。

## 選択肢

### A. モデルと Vector Index を Knowledge CLI の配布物に同梱する

特定バージョンのローカル Embedding モデルと対応 runtime / Vector Index を CLI として配布し、初回利用でもネットワークを必要としない。

- **利用者への影響:** オフラインで利用でき、追加セットアップは不要。ただし配布物が大きくなり、モデル更新は CLI 更新と同時になる。
- **公開契約・後続作業への影響:** model ID・version・ベクトル次元・Index format を migration として固定できる。モデル更新時は互換性を判定し、Index を再構築する後続 migration が必要になる。

### B. 初回利用時に固定モデルをダウンロードしてローカルキャッシュする

CLI が固定の配布元からモデルを取得し、既定 Store の近傍または別のローカルキャッシュへ保存する。

- **利用者への影響:** 配布物は小さくなるが、初回利用時にネットワークとダウンロード失敗時の復旧が必要になる。
- **公開契約・後続作業への影響:** 配布元、完全性検証、キャッシュ保存先、再試行・更新・削除方針という新しい運用契約が必要になる。現在の「設定・運用仕様を追加しない」方針との整合確認も必要になる。

### C. 利用者が管理するローカル Embedding runtime / model を必須にする

CLI は別途導入済みのローカル runtime を呼び出し、モデルの取得・更新を利用者側に委ねる。

- **利用者への影響:** CLI 配布物を小さく保てるが、事前セットアップ、runtime の起動、モデル選択、障害対応が必要になる。
- **公開契約・後続作業への影響:** runtime の発見方法、接続・バージョン・エラー、設定または環境要件を定める必要がある。既存の公開設定を追加しない方針からの変更になり得る。

## 推奨

**A（配布物への同梱）を推奨する。** ローカル専用・外部送信なし・オフライン利用という既決方針に最も直接的に適合し、モデルと Index の互換性を CLI が一元管理できるためである。B と C はモデル運用を外へ出せる一方、未承認のネットワーク・保存先・設定・障害対応契約を増やす。

承認後、Feature Design は選択された提供方式に合う具体的なモデル識別子、ライセンス確認、runtime / Vector Index、ベクトル次元、永続化 schema、migration、`search-semantic` wire schema、再構築と失敗復旧を確定する。A を承認する場合も、同梱する具体的なモデルと Index library は、その選定根拠を添えてこの Decision の追記または後続 L2 Decision として記録する。

## 承認後に更新する設計資料

- `docs/features/FEAT-006/design.md`
- `docs/features/FEAT-006/design/**`
- 必要なら `docs/features/FEAT-006/requirements.md`

## 解消しない場合の停止条件

モデル識別子・ベクトル次元・Index 永続化・初期化失敗の物理／wire 契約を確定できないため、Contract Completeness Gate を満たす詳細設計および Task Breakdown へ進めない。
