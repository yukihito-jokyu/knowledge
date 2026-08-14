# ドメイン境界

## Content Assessment Domain

記事を Article Claim へ変換し、Knowledge Assessment を統合してユーザー固有の Reading Value を決定する。

- 所有する概念: Article、Article Claim、Claim Role、Importance、Support、Knowledge Assessment、Recognition Gain、Reading Recommendation、Attention Cost。
- 所有しない概念: ユーザー知識の保存・更新の物理的管理。

## Personal Knowledge Domain

ユーザー知識の Evidence を根拠として Assertion、Concept、Scope、Relation、Temporal 情報および Derived Knowledge State を管理する。

- 所有する概念: Knowledge Assertion、Evidence、Concept、Scope、Relation、Temporal Metadata、Derived User Knowledge State、Search Trace。
- 所有しない概念: 記事 Claim の意味解釈、検索 Query の生成、読書価値の決定。

## Orchestration Domain

Content Assessment と Personal Knowledge をまたぐ処理の開始、再調査、停止、成果物の引き渡しを管理する。

- 所有する概念: Workflow 実行、検索 Budget、再調査制御、終了判定。
- 所有しない概念: 永続データの直接更新、CLI 内部の検索アルゴリズム。

## 境界不変条件

- Evidence が知識状態の正規根拠であり、Derived State は Evidence に還元できる。
- Relation の保存とユーザーがその Relation を理解していることを同一視しない。
- Reading Value は記事品質の一般評価ではない。
