---
name: reading-value
description: 利用者がCodex会話で指定した公開技術記事URL一件を、安全なArticle Analysis、ClaimごとのKnowledge Assessment、根拠追跡可能な読書推奨へ統合する。公開技術記事の読書価値を評価したいときに使う。
---

# Reading Value

Codex会話で受けた技術記事URL一件を、Article Analysis、ClaimごとのKnowledge Assessment、Reading Value Assessmentへ統合する。専用CLI、API、Web UI、永続保存を作らない。Assessment本文を完成した後だけ、同じWorkflow内でKnowledge AcquisitionとKnowledge Updateを同期実行できるが、記事・Assessment・Update Resultを永続保存しない。

## 開始

開始前に、必ず[Article Analysis契約](references/article-analysis.md)、[Reading Value統合契約](references/reading-value.md)、[検証契約](references/verification.md)、既存`skills/knowledge-search/`、`skills/knowledge-acquisition/`、`skills/knowledge-update/`の`SKILL.md`とそれぞれが指定する成果物契約を読む。利用者入力は一件の絶対HTTP(S) URLだけとし、fragmentは同一記事内の位置補助として扱う。空、複数、相対URL、HTTP(S)以外は接続せず、入力を一件の絶対HTTP(S) URLとして指定し直すよう求める。

## 手順

0. URL入力を受理して一回のURL評価を開始した時点で、不透明な`episode_id`を一つ割り当てる。当該評価中のユーザー寄与だけを受信順に記録し、それぞれの受信時刻をRFC 3339 UTCの`observed_at`として保持する。
1. 初期URLを契約どおり接続前に検査する。安全な接続先固定と接続時のDNS再解決禁止を保証できる取得能力がなければ、接続せず`article_unavailable`で停止する。
2. Cookie、認証情報、セッション、フォーム送信、その他の状態変更を使わない読み取り取得だけを行う。redirectを受けたら、次の接続前にそのURLを同じ手順で検査する。
3. 中心本文とClaimの根拠・読者が到達できる位置を確認し、契約のMarkdown形式でArticle Analysisを作成する。中心本文またはClaim根拠・位置が欠ける場合は、部分Claimを作らず`article_unavailable`で停止する。技術記事でない、または評価可能なClaimがゼロなら`article_not_evaluable`で停止する。Claimと無関係な要素だけの欠落は`content_limitations`へ記録できる。
4. 各Article Claimを、role、importance、location、support、scope、target_version、temporal_contextを改変せず一件ずつKnowledge Searchへ渡す。各正常結果だけをClaim ID、Assessment reference、Trace referenceの一対一のAssessment Mapへ記録する。
5. 一件でもKnowledge Searchが`technical_failure`または`canceled`なら、Assessment、知識状態、推奨へ変換せず、その失敗または中断を呼出側へ伝播して停止する。
6. [Reading Value統合契約](references/reading-value.md)に従い、Mapにある正常AssessmentとArticle Claimだけを統合し、再調査が必要かを判定する。
7. 重要な不足だけは、対象Claim、不足情報、必要理由を明記してArticle AnalysisまたはKnowledge Searchへ限定再調査を依頼できる。初回後の再調査は合計二回までとし、上限後は不確実性を`Why`または`Limitations`へ残す。Article再調査ではClaim Reconciliationを作り、同一四属性のretained ClaimだけMapを維持し、changed／addedは各ClaimをKnowledge Searchへ再度渡して新しい正常AssessmentとTraceをMapへ記録し、removedはMapから除外する。再調査後は、更新したArticle AnalysisとAssessment Mapで手順5と6を再実行する。
8. 再調査が不要、または上限に達した後、全Claimが正常Assessmentへ対応する最新Mapだけから、`read_full`、`read_selected`、`skip`のいずれか一つを決め、Reading Value Assessmentの本文・推奨・理由を完成する。
9. 完成済みAssessmentの`completed_at`をRFC 3339 UTCで確定し、開始時に割り当てた`episode_id`、対象URL、順序付き`user_contributions`と組み合わせてEpisodeを完成する。過去の会話、保存済み履歴、別作業、Article Analysis、Knowledge Assessment、Codex応答をEpisodeへ混ぜない。
10. 同じWorkflow実行内でKnowledge Acquisitionへ完成Episodeを渡す。正常なCandidate KnowledgeだけをKnowledge Updateへ順に渡し、Candidateが空でもUpdateを実行して空のUpdate Resultを受け取る。Acquisitionが停止結果を返した場合はUpdate・CLI更新・Update Resultを開始せず、それでも完成済みAssessmentを返す。Update Resultの全体状態が`completed`、`failed`、`canceled`、`partially_applied`のいずれでも、またDecisionが`skipped`または`failure_reason: outcome_unknown`を含んでも、Assessment本文・推奨・理由を変更せず、失敗として伝播させない。Update Resultは呼出側の実行中Markdown成果物だけとして保持し、会話出力、Knowledge Store、別ledger、公開UI/APIへ保存・表示しない。
11. 手順8でAssessmentを完成できない停止経路では、Candidate Knowledge、Update Result、Knowledge Acquisition、Knowledge Update、CLI更新を開始しない。非同期job、callback、scheduler、回答後の別実行、保存済み会話の再開、自動再実行を追加しない。

## 引き渡し境界

正常時は、同期更新の完了後に、一件の完成済みReading Value Assessmentを会話内へ返す。Article Analysis、Knowledge Assessment、Search Trace、Assessment Map、URL Evaluation Episode、Candidate Knowledge、Update Resultは実行中のMarkdown成果物だけとして扱う。Map外のClaim、Assessment、Traceを最終評価の根拠に使わない。記事全体の品質順位、真偽断定、利用者の熟達度判定は行わない。

## 検証

受入・安全境界の再現手順と観測oracleは、[検証契約](references/verification.md)に従う。

- 初期URLと全redirectが、接続前に検査されていることを確認する。
- 接続が検査済みIPに固定され、接続時にhostnameを再名前解決しないことを確認する。
- 各Claimに必須fieldと本文内の根拠・位置があること、または設計済みの停止結果になっていることを確認する。
- 各正常ClaimがAssessment Mapから一つのAssessmentとTraceへ追跡でき、失敗または中断時にReading Value Assessmentを作らないことを確認する。
- 再調査が二回を超えず、四属性が完全一致するretained ClaimだけAssessmentを再利用することを確認する。
- 正常Assessmentの完成後・会話返却前に、同一URL評価EpisodeからKnowledge Acquisition、続けてKnowledge Updateを同期実行し、全Update結果でもAssessment本文・推奨・理由が不変であることを確認する。
- Assessment未完成の停止経路では、Candidate、Update Result、Knowledge Acquisition、Knowledge Update、CLI更新を開始しないことを確認する。
- Update Resultが実行中成果物だけであり、非同期実行、再開、公開UI/API、別ledger、記事本文の永続保存を追加していないことを確認する。
- 統合検証では、`completed`、空Candidate、`skipped`、`failed`、`canceled`、`partially_applied`、`failure_reason: outcome_unknown`ごとに、Assessment完成→Acquisition→Update→同一Assessment返却の順序、本文・推奨・理由の完全一致、Update Resultの会話出力不在を観測する。ArticleまたはSearchの停止結果ではAssessment未完成後にAcquisition・Update・CLIが開始されないことを観測する。
