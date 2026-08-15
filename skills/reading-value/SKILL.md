---
name: reading-value
description: 利用者がCodex会話で指定した公開技術記事URL一件を、安全なArticle Analysis、ClaimごとのKnowledge Assessment、根拠追跡可能な読書推奨へ統合する。公開技術記事の読書価値を評価したいときに使う。
---

# Reading Value

Codex会話で受けた技術記事URL一件を、Article Analysis、ClaimごとのKnowledge Assessment、Reading Value Assessmentへ統合する。専用CLI、API、Web UI、永続保存を作らず、記事・利用者の知識・Knowledge Storeも変更しない。

## 開始

開始前に、必ず[Article Analysis契約](references/article-analysis.md)、[Reading Value統合契約](references/reading-value.md)、既存`skills/knowledge-search/`の`SKILL.md`とAssessment / Search Trace契約を読む。利用者入力は一件の絶対HTTP(S) URLだけとし、fragmentは同一記事内の位置補助として扱う。空、複数、相対URL、HTTP(S)以外は接続せず、入力を一件の絶対HTTP(S) URLとして指定し直すよう求める。

## 手順

1. 初期URLを契約どおり接続前に検査する。安全な接続先固定と接続時のDNS再解決禁止を保証できる取得能力がなければ、接続せず`article_unavailable`で停止する。
2. Cookie、認証情報、セッション、フォーム送信、その他の状態変更を使わない読み取り取得だけを行う。redirectを受けたら、次の接続前にそのURLを同じ手順で検査する。
3. 中心本文とClaimの根拠・読者が到達できる位置を確認し、契約のMarkdown形式でArticle Analysisを作成する。中心本文またはClaim根拠・位置が欠ける場合は、部分Claimを作らず`article_unavailable`で停止する。技術記事でない、または評価可能なClaimがゼロなら`article_not_evaluable`で停止する。Claimと無関係な要素だけの欠落は`content_limitations`へ記録できる。
4. 各Article Claimを、role、importance、location、support、scope、target_version、temporal_contextを改変せず一件ずつKnowledge Searchへ渡す。各正常結果だけをClaim ID、Assessment reference、Trace referenceの一対一のAssessment Mapへ記録する。
5. 一件でもKnowledge Searchが`technical_failure`または`canceled`なら、Assessment、知識状態、推奨へ変換せず、その失敗または中断を呼出側へ伝播して停止する。
6. [Reading Value統合契約](references/reading-value.md)に従い、Mapにある正常AssessmentとArticle Claimだけを統合し、再調査が必要かを判定する。
7. 重要な不足だけは、対象Claim、不足情報、必要理由を明記してArticle AnalysisまたはKnowledge Searchへ限定再調査を依頼できる。初回後の再調査は合計二回までとし、上限後は不確実性を`Why`または`Limitations`へ残す。Article再調査ではClaim Reconciliationを作り、同一四属性のretained ClaimだけMapを維持し、changed／addedは各ClaimをKnowledge Searchへ再度渡して新しい正常AssessmentとTraceをMapへ記録し、removedはMapから除外する。再調査後は、更新したArticle AnalysisとAssessment Mapで手順5と6を再実行する。
8. 再調査が不要、または上限に達した後、全Claimが正常Assessmentへ対応する最新Mapだけから、`read_full`、`read_selected`、`skip`のいずれか一つを決め、Reading Value Assessmentを会話内へ返す。

## 引き渡し境界

正常時は一件のReading Value Assessmentを会話内へ返す。Article Analysis、Knowledge Assessment、Search Trace、Assessment Mapは実行中のMarkdown成果物だけとして扱う。Map外のClaim、Assessment、Traceを最終評価の根拠に使わない。記事全体の品質順位、真偽断定、利用者の熟達度判定は行わない。

## 検証

- 初期URLと全redirectが、接続前に検査されていることを確認する。
- 接続が検査済みIPに固定され、接続時にhostnameを再名前解決しないことを確認する。
- 各Claimに必須fieldと本文内の根拠・位置があること、または設計済みの停止結果になっていることを確認する。
- 各正常ClaimがAssessment Mapから一つのAssessmentとTraceへ追跡でき、失敗または中断時にReading Value Assessmentを作らないことを確認する。
- 再調査が二回を超えず、四属性が完全一致するretained ClaimだけAssessmentを再利用することを確認する。
- CLI、API、Web UI、DB、公開設定、記事本文の永続保存を追加していないことを確認する。
