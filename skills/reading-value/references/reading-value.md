# Reading Value統合契約

## 境界

Reading Valueは、Article AnalysisのClaimと正常終了したKnowledge Assessmentを統合し、会話内に一件のReading Value Assessmentを返す。Assessmentの状態、Confidence、Known、Knowledge Gap、Evidence、矛盾、時点差分、Knowledge Searchの探索戦略・Budgetを再判定または変更しない。Article本文、Claim、Assessment、Trace、最終評価は実行中のMarkdown成果物であり、永続化しない。

Article Analysisの各Claimを、属性を改変せず一件ずつKnowledge Searchへ渡す。正常なAssessmentだけを、Claim ID、Assessment reference、Assessment内のTrace referenceで構成するAssessment Mapへ一対一で記録する。`Read`、`Skip`、`Knowledge Gains`、`Why`はMapにあるClaim IDとAssessment referenceだけを根拠にし、自由文の推測でAssessmentやTraceを結び直さない。

一件でも`technical_failure`または`canceled`なら、Assessment、知識状態、推奨を作らず、失敗または中断を呼出側へ伝播する。部分成功を`no_evidence`、`skip`、利用者の未知へ変換してはならない。

## 読書推奨

利用者固有の認識利得を、Assessmentが示す既知・部分既知・推論可能性・矛盾または時点差分、Claimのrole・importance・support・位置・適用可能性、Attention Costから判断し、次のいずれか一つを選ぶ。

| 推奨 | 選ぶ条件 | 必須説明 |
| --- | --- | --- |
| `read_full` | 高い認識利得が記事全体に分散する、またはClaim間の因果・構造・判断に全体文脈が必要 | 全文が必要な理由、中心Gain、Reliability、低優先箇所の扱い |
| `read_selected` | 重要かつ十分にSupportされたGainが特定位置にあり、残りは既知・低重要度・根拠不足・再包装などで相対的に低価値 | 読むClaim IDと位置、飛ばす位置、各Gain、Why、Reliability |
| `skip` | 主な内容が既知、またはGain候補が低重要度・低適用性・根拠不足でAttention Costを正当化しない | 読まない理由、未知を断定していないこと、重要な例外の位置 |

`known`はFact Gainを小さくし得るが、新しい構造化はStructural Gainになり得る。`partially_known`、`inferable`、`contradicted`、`outdated`は、Assessmentの結論を変えずに、Gain候補として扱える。`no_evidence`は利用者の未知の証拠ではなく、新規可能性を示す限定材料にとどめる。`uncertain`は確定Gainにせず、必要なら再調査候補またはLimitationsに残す。記事一般の品質順位、著者の人格評価、AI生成判定、記事の真偽断定を行わない。

## 再調査

初回後に起票できる再調査は合計二回までとする。一回ごとに、Article AnalysisまたはKnowledge Searchのいずれか一方へ、対象Claim、不足情報、必要理由を明記した限定依頼を出す。取得失敗、`technical_failure`、`canceled`は再調査枠を消費せず、直ちに停止・伝播する。上限後も重要な不確実性が残る場合は、追加の往復をせず、正常Assessmentだけで保守的に判断し、`Why`または`Limitations`へ残す。

Article再調査ではClaim Reconciliationを保持する。`claim`、`scope`、`target_version`、`temporal_context`がすべて同一の`retained` Claimだけ、既存Claim IDとAssessment Mapを維持する。`changed`と`added`は旧Mapを無効にしてKnowledge Searchを再実行し、`removed`はMapと最終評価から除外する。この追随検索はArticle再調査に含まれ、追加の再調査枠を消費しない。Knowledge Searchだけへの再調査は、一回の再調査枠を消費する。

## 出力

全ClaimのKnowledge Searchが正常終了した場合だけ、次の形式で返す。

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

<重要Claim、Assessment、Support、Attention Costに基づく理由。>

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

<Support、一次資料参照、実装例、失敗例、未解決の不足を説明する。>

## Limitations

<不確実性または「なし」。>

## Conclusion

<利用者がどう読むべきかを平易に述べる。>
```

`Recommendation`は一つだけにする。`read_selected`では`Read`と`Skip`を必須とし、`Read`にはClaim IDと記事内位置を記す。`skip`では`Read`を空表または「なし」とする。`Gain type`は`Fact`、`Understanding`、`Structural`、`Decision`、`Practical`、`Correction`、`Exploration`の一つ以上とし、AssessmentのStatusを機械的に置換しない。
