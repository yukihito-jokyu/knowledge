# FEAT-002 要件

## 目的

Article Claim を入力として、Codex が既存 Knowledge CLI を反復利用し、ユーザーの知識状態を根拠付きで評価する。結果は後続の Reading Value が再判定せず利用できる Knowledge Assessment と、診断用の Search Trace に分離する。

## 含む要件

- REQ-004、REQ-005、REQ-019、REQ-020、REQ-021
- BR-002、BR-003、BR-006、BR-007、BR-010
- NFR-001、NFR-002、NFR-003、NFR-005、NFR-006
- CON-002、CON-003、CON-006
- DEC-REQ-001（初期提供では Semantic Search を使わない）

## Feature Acceptance

1. Claim 本文の直接探索だけで終えず、必要に応じてClaimの構成要素、Concept、Relation、Evidence、矛盾候補、時点差分を探索する。
2. `known`、`partially_known`、`inferable`、`contradicted`、`outdated`、`no_evidence`、`uncertain` のいずれか一つを、根拠・確信度・Known・Knowledge Gap とともに返す。
3. Evidence未観測を未知と断定せず、Relation の保存だけを理解の根拠にしない。
4. 検索は強い結論、飽和、増分なし、強い矛盾、または固定 Budget で停止し、操作過程を Search Trace として追跡できる。
5. CLI の技術エラーは知識状態へ変換せず、評価を失敗として呼出側へ返す。

## 範囲外

- Article URLの取得、Claim分解、記事内の重要度・根拠・位置の判断（FEAT-003）。
- Reading Recommendation と Attention Cost の評価（FEAT-003）。
- Conversation / Task Episode からのEvidence候補抽出・知識更新（FEAT-004）。
- Semantic Search、Embedding、Vector Index（FEAT-006）。
- Knowledge CLI の操作、JSON wire schema、SQLite schema、公開設定の追加・変更。
