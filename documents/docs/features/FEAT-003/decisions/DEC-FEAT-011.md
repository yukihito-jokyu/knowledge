# DEC-FEAT-011: URL評価の再調査上限

- **Status:** decided
- **Level:** L2
- **Decision:** Parent Orchestrationは、初回のArticle Analysis・Knowledge Search・Reading Valueの後に、Reading Valueから起票された再調査を合計2回まで実行する。一回の再調査は、Article AnalysisまたはKnowledge Searchのどちらか一方に対する、対象Claim・不足情報・必要理由を明記した限定依頼とする。

## 根拠

REQ-020は必要な再調査とその回数制御を求め、Issue #175 §37は一方向Pipelineを禁じる。一方、再調査の無制限な往復は利用者の注意コストを肩代わりする目的とNFR-002に反する。初回処理と最大2回の限定再調査により、重要な不足を補いつつ終了を保証する。

## 詳細

- Article Analysisへの再調査は、重要ClaimのSupport、位置、役割、全体文脈の不足を確認する場合だけ行う。結果は変更・追加されたClaimだけを置換し、変更されないClaimを再分解しない。
- Knowledge Searchへの再調査は、重要Claimの`no_evidence`、`uncertain`、Scopeまたは時点の不足を補う場合だけ行う。対象Claimごとの実行はFEAT-002の12回CLI呼出しBudgetを新しい一実行として守る。
- Article再調査後、Claim本文および明示されたScope、対象version、時点文脈が同一のClaimだけは既存Assessmentを再利用する。これらのいずれかが変わったClaimと追加Claimは、旧Assessmentを無効にしてKnowledge Searchを再実行する。削除ClaimはReading Valueの入力から除外する。この追随検索は、Article Analysisへの一回の再調査に含まれ、別の再調査枠を消費しない。
- 再調査後に同じ不足が残る、または上限に達した場合、Reading Valueは未確定性をReliabilityまたはWhyに明記して最終判断する。追加の往復は行わない。
- Article取得失敗やKnowledge Searchの`technical_failure`／`canceled`は再調査枠を消費せず、評価を失敗または中断として呼出側へ返す。
