---
name: knowledge-search
description: 一つの独立評価可能なClaimを既存Knowledge CLIで有限に反復探索し、Evidenceに基づくKnowledge Assessmentの状態候補を決定する。記事Claimの知識評価、Evidence探索、known・inferable・contradicted等の判定を依頼されたときに使う。CLIの公開契約、DB、設定を変更する実装には使わない。
---

# Knowledge Search

一つのTarget Claimを既存Knowledge CLIの読取操作で同期探索する。Claimの意味対応、次の探索経路、Evidence強度、評価状態はCodexが判断し、CLIへ委ねない。

## 開始

最初のCLI呼出し前に、同梱の[CLI操作リファレンス](references/cli-operations.md)と[Assessment / Search Trace契約](references/artifact-contract.md)を読む。実行時にリポジトリ内外の設計資料・ソースコード・URLを正規契約として読まない。必要な操作、option、成功・失敗JSON、終了コード、成果物形式はこのSkill本文と同梱リファレンスだけで判断する。

## Input Boundary

- 一つの独立評価可能なTarget Claimを受け取る。Claimを勝手に再分解しない。
- Claimに明示されたScope、version、時点だけを使う。未指定のfilterやselectorを推測しない。現在日時、評価日時、一般知識から`--at`、`--valid-from`、`--valid-until`、Scopeを補完しない。
- `search-semantic`、Embedding、Vector Index、CLIの書込み操作は使わない。
- 新しいCLI operation、JSON field、保存、公開設定を作らない。

## Exploration Procedure

1. 利用するCLI operationの必須option、enum、selector、値の組を正規資料で確認する。不正な入力を発見したらCLIを実行せずtechnical failureとして止める。
2. `search-text`を一度だけ実行し、Claim全体への直接Evidenceがあるか評価する。
3. 未解決部分がある場合だけ、Claimの構成要素または明示Conceptを必要な範囲で選び、`search-text`または`search-concept`を使う。
4. 結果から得たIDだけに対して、必要時に`get`、`get-evidence`、`search-related`、`search-contradictions`、`search-temporal`を呼ぶ。`search-temporal`にはConceptまたはScope selectorが必要であり、時点条件だけでは呼ばない。時点filterを付けるならClaimに明示された時点だけを使う。
5. 各呼出しの前にBudgetと同一論理queryの既実行有無を確認する。呼出し直後に、結果ID、Evidence ID、増分、次の判断を実行中の探索記録へ残す。
6. 強い結論、飽和、増分なし、探索可能な未解決経路なし、またはBudget到達で止める。未解決部分はKnowledge Gapへ残す。

## Fixed Limits

- CLI呼出しは最大12回。直接／構成要素／Concept探索は合計4回、候補詳細／Relation／矛盾／時点探索は合計4回、Evidence hydrationは最大4回。
- Relation深さは1。`search-contradictions`と`search-temporal`は合計2回まで。
- 同じoperationと同じ論理queryを再実行しない。上限を超える探索で根拠を補強しない。

## Evaluation Rules

Evidence原文、Scope、時点、Claimとの意味対応に基づき、次の優先順で一つだけ状態候補を決める。

- `uncertain`: 相反するEvidenceまたは複数の意味対応に優位性がない。
- `contradicted`: Claimと両立しない理解を直接支えるEvidenceが優位である。
- `outdated`: 現行のScopeに対し古い理解だけがあり、新しい訂正または時点根拠が差分を示す。
- `known`: Claim全体を直接または意味的に同等なAssertionとして支えるstrong Evidenceがあり、矛盾・時点差分がない。
- `partially_known`: Claimの一部または構成知識だけにEvidenceがある。
- `inferable`: Claim自体の直接Evidenceはないが、必要な構成知識と明示Relationから導出でき、その経路もEvidenceで支持される。
- `no_evidence`: 探索した経路に支持Evidenceがない。

Concept一致、Relationの存在、候補の類似、EvidenceなしAssertion、質問、AI説明だけを`known`の根拠にしない。`no_evidence`をユーザーの未知や熟達度の結論へ変換しない。

## Output Boundary

- 正常探索だけで、同梱の成果物契約どおりのMarkdown `Knowledge Assessment` と対応する一時的なMarkdown `Search Trace` を返す。Assessmentには実行内だけで有効な`trace:<run-id>`参照を載せる。
- Assessmentにはstatus、confidence、Known、Knowledge Gap、Supporting Evidence ID、明示Relation、矛盾・時点差分、結論を含める。Traceには使用Budget、実行済みoperation／論理query、結果ID、停止理由を含める。
- `technical_failure`または`canceled`ではAssessmentを作らない。部分Traceと状態だけを呼出側へ返す。Parent Orchestration自身が中断済みならTraceも作らず中断を伝播する。
- いずれもファイル保存、公開JSON API、新CLI operation、DB更新、設定追加を行わない。

## Failure Boundary

- `validation_error`、`not_found`、`storage_error`、`internal_error`、JSONプロトコル不整合では評価状態を作らない。同じ入力で再試行しない。
- response開始前のexit 130・無出力はerror JSONを期待せず、`canceled`として扱う。
- 空Store、空検索結果、未登録Conceptは成功結果であり、technical failureと混同しない。

## Verification

- 直接探索が最初のCLI呼出しであることを確認する。
- 7状態の優先規則、Evidenceの直接性、Scope／時点を説明できることを確認する。
- Budget、同一query禁止、Relation深さ、矛盾／時点探索上限を探索記録で確認する。
- 新しいCLI公開契約、DB変更、永続化、設定がないことを確認する。
