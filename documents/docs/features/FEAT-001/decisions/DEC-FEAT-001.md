# DEC-FEAT-001: Semantic Search の Embedding 実行方式と Vector Index

- **Status:** decided
- **Level:** L3
- **論点:** ローカル SQLite Store と連携する `search-semantic` を、どこで Embedding を生成し、どの方式で類似検索するか。

## 決定

**C を初期提供として採用し、後続移行では A（ローカル埋め込みモデル + ローカル Vector Index）を採用する。**

初期の SQLite データ表現には、Assertion を再作成せず Embedding と Vector Index を後付けできる安定した識別子、正規化 Assertion 本文、Scope、Concept を保持する。移行は正規データから Index を再構築し、Evidence・Relation・履歴を変更しない。

## なぜ今決める必要があったか

REQ-011 と CON-004 は Semantic Search を必須とする。Embedding の実行場所は、個人 Evidence の外部送信、ネットワーク可用性、費用、依存関係、検索品質、初期の SQLite 境界に影響する。実装タスクと受入評価を確定するために必要である。

## 選択肢

### A. ローカル埋め込みモデル + ローカル Vector Index

- **利点:** Evidence を外部送信せず、オフラインでも検索できる。SQLite ローカル専用方針と整合する。
- **欠点:** モデル配布・実行環境・Index 容量の管理が必要で、初期セットアップが重い。
- **影響:** CLI はローカルのモデル実行と Index を管理する。

### B. 外部 Embedding API + ローカル Vector Index

- **利点:** 高品質なモデルを容易に利用でき、モデル運用を自前で持たない。
- **欠点:** Assertion と Scope を外部送信し、API Key、費用、ネットワーク障害、データ取り扱いを管理する必要がある。Evidence 本文を送るかは別途限定が必要。
- **影響:** CLI またはその外側に外部 AI 接続・設定・秘密情報管理の境界が加わる。

### C. 初期は SQLite の字句検索のみを提供し、Semantic Search を後続 Feature に延期

- **利点:** 初期実装を軽くできる。
- **欠点:** REQ-011、CON-004、Scenario C の探索要件を満たさず、承認済み要件から外れる。
- **影響:** 要件変更の承認が必要になる。

## 推奨

初期提供の複雑さを抑えるため C を選択し、最終的なプライバシー境界と必須論理能力を満たす移行先として A を選択した。

## 影響

- FEAT-001 は字句・Concept・Relation・Evidence・時点差分の検索を初期提供する。
- Semantic Search の実装・評価は後続移行の Feature として backlog 化する。
- 初期のデータ・JSON 契約は、Semantic Search の追加によって既存の検索・更新結果を破壊しない。
