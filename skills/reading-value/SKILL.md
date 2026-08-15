---
name: reading-value
description: 利用者がCodex会話で指定した公開技術記事URL一件を、安全に取得して根拠位置付きのArticle Analysisへ分解する。記事取得・Article Claimの抽出または再調査時に使う。Knowledge Assessmentと読書推奨の判断には使わない。
---

# Reading Value

この段階では、Codex会話で受けた技術記事URL一件をArticle Analysisへ変換する。専用CLI、API、Web UI、永続保存を作らず、記事・利用者の知識・Knowledge Storeも変更しない。

## 開始

取得前に、必ず[Article Analysis契約](references/article-analysis.md)を読む。利用者入力は一件の絶対HTTP(S) URLだけとし、fragmentは同一記事内の位置補助として扱う。空、複数、相対URL、HTTP(S)以外は接続せず、入力を一件の絶対HTTP(S) URLとして指定し直すよう求める。

## 手順

1. 初期URLを契約どおり接続前に検査する。安全な接続先固定と接続時のDNS再解決禁止を保証できる取得能力がなければ、接続せず`article_unavailable`で停止する。
2. Cookie、認証情報、セッション、フォーム送信、その他の状態変更を使わない読み取り取得だけを行う。redirectを受けたら、次の接続前にそのURLを同じ手順で検査する。
3. 中心本文とClaimの根拠・読者が到達できる位置を確認し、契約のMarkdown形式でArticle Analysisを作成する。
4. 中心本文またはClaim根拠・位置が欠ける場合は、部分Claimを作らず`article_unavailable`で停止する。技術記事でない、または評価可能なClaimがゼロなら`article_not_evaluable`で停止する。Claimと無関係な要素だけの欠落は`content_limitations`へ記録できる。
5. 再調査では、前回のArticle Analysisと照合してClaim Reconciliationを作成する。

## 引き渡し境界

成功時はArticle Analysisだけを親オーケストレーションへ返す。ClaimごとのKnowledge Assessment、Assessment Map、再調査回数の制御、`read_full`／`read_selected`／`skip`、記事品質の評価はこの段階の対象外であり、推測して返さない。

## 検証

- 初期URLと全redirectが、接続前に検査されていることを確認する。
- 接続が検査済みIPに固定され、接続時にhostnameを再名前解決しないことを確認する。
- 各Claimに必須fieldと本文内の根拠・位置があること、または設計済みの停止結果になっていることを確認する。
- CLI、API、Web UI、DB、公開設定、記事本文の永続保存を追加していないことを確認する。
