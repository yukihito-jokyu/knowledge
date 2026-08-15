# Reading Value Assessment 成果物契約

Reading Value Assessmentは、Article Analysisと正常終了したKnowledge Assessmentを統合し、利用者へ返す最終Markdown成果物である。主出力は数値スコアではなく、理由を伴う読書推奨とする。

```markdown
# Reading Value Assessment

## Article

- source_url: <利用者が指定したURL>
- title: <記事タイトルまたは「不明」>

## Recommendation

<read_full | read_selected | skip>

## Assessment Map

| Claim ID | Assessment reference | Trace reference |
| --- | --- | --- |

## Why

<推薦を、重要Claim、Knowledge Assessment、Support、Attention Costの観点で説明する。>

## Read

| Location | Claim IDs | Expected recognition gain | Why read |
| --- | --- | --- | --- |

## Skip

| Location | Claim IDs | Why lower priority |
| --- | --- | --- |

## Knowledge Gains

| Claim IDs | Assessment reference | Gain type | Expected gain |
| --- | --- | --- | --- |

## Reliability

<記事内Support、一次資料参照、実装例、失敗例、未解決の不足を、独立した品質スコアにせず説明する。>

## Limitations

<no_evidence、uncertain、取得欠落、再調査上限など、結論を限定する事実。なければ「なし」。>

## Conclusion

<利用者がどう読むべきかを平易に述べる。>
```

## 記録規則

- `Recommendation`は必ず一つ。表示文では`read_full`を「全文を読む」、`read_selected`を「選んだ箇所を読む」、`skip`を「今回は読まない」と平易に説明してよい。
- `Assessment Map` の各Claim IDは、正常完了したKnowledge Assessmentの実行時参照と、そのAssessmentが保持するSearch Traceの実行時参照へ一対一に対応付ける。`Read`、`Skip`、`Knowledge Gains` は、この表にあるClaim IDとAssessment referenceだけを参照する。
- `Read`は`read_selected`で必須とし、`read_full`では全体文脈が必要な理由を記す。`skip`では空表または「なし」とする。
- `Skip`は`read_selected`で必須とし、既知・低重要度・根拠の弱さ・再包装などを具体的に示す。`read_full`でも優先度が低い補足があれば記してよい。
- `Gain type`は`Fact`、`Understanding`、`Structural`、`Decision`、`Practical`、`Correction`、`Exploration`の一つ以上とする。Assessmentの状態をGain typeへ機械的に置換しない。
- `no_evidence`は利用者が知らない証拠ではなく、新規可能性を示す限定的な材料である。Reliability、Importance、Applicability、Attention Costと併せて扱い、単独で強い推薦理由にしない。
- `Reliability`は記事が認識更新の根拠としてどこまで採用可能かを説明する。記事一般のランキング、著者の人格評価、AI生成判定は含めない。
