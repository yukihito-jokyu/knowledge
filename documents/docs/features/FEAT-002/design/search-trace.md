# Search Trace 契約

Search Trace は Knowledge Assessment 本文から分離する診断・評価用Markdown成果物である。永続保存、CLIへの書込み、新しいJSON APIは初期提供の範囲外とする。

## 必須項目

```markdown
# Search Trace

## Target Claim

<評価対象の一つのClaim>

## Budget

- used: <0..12>
- limit: 12
- evidence_hydrated: <0..4>
- relation_depth: <0..1>
- contradiction_temporal_calls: <0..2>

## Steps

| # | Intent | CLI operation | Input summary | Result IDs | Evidence IDs | New finding | Decision |
| --- | --- | --- | --- | --- | --- | --- | --- |

## Stop Reason

<strong_conclusion | saturation | no_increment | no_expandable_path | budget_exhausted | technical_failure | canceled>

## Technical Failure

<該当時だけ、operation、error code、再試行しない理由。知識状態は記載しない。>

## Cancellation

<該当時だけ、operation、exit_code: 130、error codeなし、部分Traceを返せたか。>
```

## 記録規則

- `Intent` は直接一致、構成要素、Concept、Relation、Evidence、矛盾、時点差分のいずれを確かめるかを記す。
- `Input summary` はClaim本文、ID、Concept名、Scope、時点条件を必要最小限で記し、Evidence原文を複製しない。
- `Result IDs` と `Evidence IDs` は、Assessment の主張を遡れる不透明IDとして残す。候補なしは空配列と明記する。
- `New finding` は候補、Evidence、矛盾、時点情報のいずれが新規に得られたかを記し、単にコマンドが成功したことは増分に数えない。
- `contradiction_temporal_calls` は `search-contradictions` と `search-temporal` の実行回数の合計である。2に達した後は、この二操作を実行せず未解決部分をKnowledge Gapへ残す。
- 呼出し前に、Codexは既存operation資料の必須option、enum、selector、値の組を検証する。CLIがなお`validation_error`を返した場合は、内部の操作生成不整合として`technical_failure`で停止する。Traceにはoperation、既存error code `validation_error`、その入力要約、再試行しない理由を残し、呼出側へ同じ失敗を渡す。
- CLIの`storage_error`、`internal_error`、`not_found`（直前の成功JSONから得たIDの場合）、またはJSONプロトコル不整合は `technical_failure` で停止する。Traceには既存error codeを残し、Assessmentは出力しない。空結果は技術失敗ではない。
- response開始前のCLI `Ctrl-C` はerror JSONを持たずexit 130で終了するため、`technical_failure`ではなく利用者中断の`canceled`で停止する。Knowledge Searchはoperation、`exit_code: 130`、error codeなし、実行済みstepを部分Traceへ記録し、Assessmentなしで中断を呼出側へ伝播する。Parent Orchestration自体も中断済みなら部分Traceの出力・永続化を試みず、その中断を優先する。
