# Knowledge Assessment 成果物契約

Knowledge Assessment は ClaimごとのMarkdown成果物であり、Reading Value が知識状態を再判定せず利用する正規入力である。Search Traceは [search-trace.md](search-trace.md) に分離する。

```markdown
# Knowledge Assessment

## Target Claim

<Article Analysisから受け取った一つのClaim。本Featureでは再分解しない。>

## Status

<known | partially_known | inferable | contradicted | outdated | no_evidence | uncertain>

## Confidence

<high | medium | low>

## Known

- <Evidence IDを伴う、観測済みの知識>

## Knowledge Gap

- <Claimのうち直接Evidenceまたは導出根拠が不足する部分。なければ「なし」。>

## Supporting Evidence

| Assertion ID | Evidence ID | Strength | Relevance | Why |
| --- | --- | --- | --- | --- |

## Related Knowledge

<探索したRelationと、Relation単独では理解の根拠にしない旨>

## Contradictions

<候補、対応Evidence、または「なし」>

## Temporal Issues

<対象version・時点、対応Evidence、または「なし」>

## Conclusion

<状態、Confidence、Known/GAP/矛盾・時点差分を結ぶ説明>

## Trace Reference

<対応するSearch Traceへの相対参照または実行内識別子>
```

`Strength` は `strong`、`moderate`、`weak` のいずれかとし、Evidence種別を機械的な序列へ置換しない。ユーザー自身の正しい説明・推論・コードは通常strong、明示的自己申告は通常moderate、Concept認識だけはweakとする。後日の訂正は、過去Evidenceとの比較で強い反証となり得る。最終的な関連性、正しさ、時間的な優先はCodexがEvidence原文とClaimの意味から判断する。
