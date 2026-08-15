# Assessment / Search Trace 契約

このSkillの正常実行は、永続化しない二つのMarkdown成果物を呼出側へ返す。公開JSON API、CLI operation、DB、設定を追加しない。

## 正常結果

一実行内で一意な `trace:<run-id>` を作る。Assessmentにはその参照だけを載せ、Trace本文を重複させない。`<run-id>` はその実行内でのみ有効であり、保存先やURLを表してはならない。

### Assessment

```markdown
# Knowledge Assessment

## Target Claim

<入力Claimをそのまま記載>

## Status

`known | partially_known | inferable | contradicted | outdated | no_evidence | uncertain`

## Confidence

`high | medium | low` — <Evidenceの直接性・完全性・競合の有無>

## Known

<支持された内容。ない場合は「なし」>

## Knowledge Gap

<未解決部分。ない場合は「なし」>

## Supporting Evidence

| Evidence ID | Assertion ID | 対応 | 要約 |
| --- | --- | --- | --- |
| `<id>` | `<id>` | direct / component / inference-path | <必要最小限の要約> |

## Related Knowledge

<明示Relationと導出経路。ない場合は「なし」>

## Contradictions

<矛盾候補と評価。ない場合は「なし」>

## Temporal Issues

<Scope・時点差分。ない場合は「なし」>

## Conclusion

<状態候補を選んだEvidenceベースの結論。未観測をユーザーの未知へ変換しない。>

## Trace Reference

`trace:<run-id>`
```

`Supporting Evidence` はEvidence IDと必要最小限の意味要約を置く。Evidence原文をAssessmentまたはTraceへ丸ごと複製しない。

### Search Trace

```markdown
# Search Trace

## Target Claim

<入力Claimをそのまま記載>

## Budget

| 区分 | 使用 / 上限 |
| --- | --- |
| 直接・構成要素・Concept探索 | <n> / 4 |
| 候補詳細・Relation・矛盾・時点探索 | <n> / 4 |
| Evidence hydration | <n> / 4 |
| 総CLI呼出し | <n> / 12 |
| Relation depth | <n> / 1 |
| 矛盾・時点探索 | <n> / 2 |

## Steps

| # | Operation | 最小入力要約 | 結果ID / Evidence ID | 増分 | 次の判断 |
| ---: | --- | --- | --- | --- | --- |
| 1 | `search-text` | Claim全文 | `<id>` またはなし | <増分> | <判断> |

## Stop Reason

<strong conclusion / saturation / no new information / no viable path / budget exhausted>
```

Traceにはoperation、入力の最小要約、ID、Budget、停止理由だけを残す。Evidenceの `raw_text`、CLI成功JSON全文、不要な個人情報を複製しない。

## technical_failure

`validation_error`、`not_found`、`storage_error`、`internal_error`、または成功JSONの契約不整合では、Assessmentを作らず、呼出側へ `technical_failure` と次の部分Traceだけを返す。

```markdown
# Search Trace

## Target Claim

<入力Claim>

## Budget

<正常Traceと同じ表>

## Steps

<失敗前までの実行済みStep>

## Technical Failure

| Operation | exit code | error code | 処置 |
| --- | ---: | --- | --- |
| `<operation>` | <n> | `<code>` | 停止。再試行しない。 |
```

この場合 `Trace Reference` を含むAssessment、知識状態、推測のerror codeを返してはならない。

## canceled

response開始前のexit 130・stdout/stderr無出力では、Assessmentを作らず、error codeなしで `canceled` と次の部分Traceを返す。

```markdown
# Search Trace

## Target Claim

<入力Claim>

## Budget

<正常Traceと同じ表>

## Steps

<中断前までの実行済みStep>

## Cancellation

| Operation | exit code | stdout / stderr | 処置 |
| --- | ---: | --- | --- |
| `<operation>` | 130 | responseなし | 中断を呼出側へ伝播。 |
```

Parent Orchestration自身がすでに中断済みなら、Traceも生成・保存せず中断をそのまま伝播する。
