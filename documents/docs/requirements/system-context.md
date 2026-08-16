# システムコンテキスト

## 目的

継続的に技術を学ぶエンジニアが URL を指定したコンテンツについて、自身の既存知識との差分と注意コストを踏まえ、読む価値と読むべき箇所を判断できるようにする。一般的な記事品質や AI 生成の有無ではなく、ユーザー固有の有意味な認識利得を扱う。

## アクター

- 主ユーザー: 継続的に技術キャッチアップを行うエンジニア
- Codex: 意味理解、検索方針、知識評価、読書価値判定および各サブワークフローの制御を担う AI Runtime

## 外部システム

- 対象コンテンツの URL とその記事本文
- Knowledge CLI: ユーザー知識を決定論的に保存・検索・更新する永続知識インターフェース
- Knowledge Store: Knowledge CLI が管理する Evidence と知識関連情報の保存先
- Codex workspace sandbox: OSのユーザー設定領域に書込みできない場合がある、利用者が選んだworkspace内Storeを利用できる実行環境

## 対象範囲

- URL 起点の技術記事を Claim へ分解し、ユーザー知識との差分に基づいて `read_full`、`read_selected`、`skip` を推奨する。
- 会話・作業エピソードから、ユーザー知識の根拠候補を抽出し、Knowledge Store を継続的に更新する。
- Evidence 起点の Knowledge State、関連・矛盾・時点差分を含む探索、評価可能な検索・更新を提供する。
- 利用者が明示したローカルSQLite Storeを、CLI invocation単位で選択できる。指定がない場合の既定Storeは維持する。

## 対象外

- AI 生成コンテンツの判定
- 一般向けの記事ランキングや記事単体の品質採点
- ユーザーを初心者・中級者・上級者へ固定的に分類すること
- 閲覧・要約・AI による説明だけからの知識獲得判定
- Knowledge CLI 内での AI 判断・Agent 実行

## 主要な用語

- **Knowledge Assertion**: ユーザーが持つ可能性のある、独立して評価可能な正規化済みの命題。
- **Evidence**: Assertion をユーザーが持つと判断する根拠となる原文、コード、訂正など。
- **Derived User Knowledge State**: Evidence から導出される、時点を持つ知識状態。
- **Article Claim**: 記事を知識差分と比較するために分解した、役割・重要度・位置・根拠を持つ主張。
- **Recognition Gain**: Fact、Understanding、Structural、Decision、Practical、Correction、Exploration の各利得。

## 正規ソース

- Issue #175「要件定義」: 承認済み統合設計（2026-08-12取得）
