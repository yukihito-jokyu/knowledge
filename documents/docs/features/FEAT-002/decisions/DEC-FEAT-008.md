# DEC-FEAT-008: 初期 Agentic Search の固定Budgetと探索順序

- **Status:** decided
- **Level:** L2
- **Decision:** FEAT-002 は1 Claimごとに最大12回の Knowledge CLI 呼出し、最大4件の Evidence hydration、Relation の深さ1を固定上限とする。公開option、設定ファイル、環境変数は追加しない。

## 根拠

REQ-021、NFR-002、ASM-002、および横断設定方針は、停止可能な検索を要求しつつ検索 Budget 等を公開設定として追加しない。初期提供のCLIはSemantic Searchを持たず、字句・Concept・Relation・Evidence・矛盾・時点検索を提供するため、直接探索から限定的な拡張へ進む上限をFeature内の制御契約として固定する。

## 内訳

- 直接字句探索と構成要素またはConcept探索: 合計最大4回（最初の直接探索1回を含む）。
- 候補の詳細取得、Relation探索、矛盾候補または時点差分探索: 合計最大4回。Relation深さは1、矛盾・時点探索は合計2回までとする。
- 候補のEvidence取得: 最大4回。

同じ操作・同じ論理queryを再実行しない。上限に達する前でも、強い結論、候補・Evidenceの増分停止、または探索可能な未解決経路がないとき停止する。

## 影響

この値は Codex 側のFeature内ワークフローに限る。CLIの公開契約、Store、後続のSemantic Search移行を変更しない。Fixture評価で過小探索が繰り返し確認された場合は、L2 Decisionを更新して上限を見直す。
