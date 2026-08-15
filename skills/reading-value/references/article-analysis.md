# Article Analysis契約

## 目的と境界

Article Analysisは、取得済みの公開技術記事を、利用者の知識と独立に評価できるArticle Claim群へ変換する実行内Markdown成果物である。記事全文の要約、Knowledge Assessment、読書推奨、一般的な品質評価ではない。記事URL・本文・Claimを永続化せず、記事から利用者が読んだ、理解した、または知っているというEvidenceも作らない。

成功時の成果物は次の形式にする。

```markdown
# Article Analysis

## Article

- source_url: <利用者が指定したURL>
- resolved_url: <取得後の最終URL>
- title: <記事タイトル。取得できなければ「不明」>
- retrieved_at: <評価実行時刻>
- content_limitations: <Claimに無関係な欠落。なければ「なし」>

## Overview

<主題とClaim群を読むために必要な最小の文脈。>

## Claims

### CLM-001

- claim: <独立に評価できる一つの命題>
- role: <core | supporting | example | background | incidental>
- importance: <high | medium | low>
- location: <節見出しと段落上の位置。利用可能なら既存anchorを補助的に付記>
- support: <命題を支える本文内の説明、実装例、失敗例、測定、一次資料参照など>
- scope: <該当する対象・条件。なければ「なし」>
- target_version: <該当する製品・仕様・実装version。なければ「なし」>
- temporal_context: <該当する時点・有効期間。なければ「なし」>

## Coverage Notes

<Claimにしなかった反復、装飾、独立評価不能な断片と理由。>

## Claim Reconciliation

<初回は「初回解析」。再調査時は旧Claim ID、新Claim ID、分類、理由。>
```

`CLM-###`は当該実行内だけで有効な順序付きIDであり、永続IDではない。

## URL入力と公開性検査

一件の絶対URLだけを受け入れる。fragmentは記事内位置の補助であり、別入力にはしない。各接続候補（初期URLと全redirect先）は、**その接続の前**に次をすべて満たすときだけ接続を試みる。

- schemeが`http`または`https`である。
- authorityにuserinfoがない。
- portが省略されているか、`http`では80、`https`では443である。
- hostnameを解決した接続候補が一件以上あり、**すべて**がグローバル到達可能IPである。IP literalも同じ判定をする。
- hostnameがlocalhostその他のローカル名ではなく、候補にloopback、private、link-local、unique-local、multicast、unspecified、予約済みその他の公開外IPが一つもない。

解決失敗、候補なし、候補の一つでも公開外、検査不能、または不適格なURLがあれば、その対象へ接続せず`article_unavailable`で停止する。redirectの`Location`は現在URLに対して解決した後、次の接続前に同じ検査を行う。redirect loop、Location欠落・解決不能、または取得能力がredirect先を接続前に検査できない場合も`article_unavailable`で停止する。

## 接続固定と取得制限

接続は、直前の公開性検査で承認したIP候補のみに固定する。接続開始時にhostnameをDNS再解決してはならない。HTTP HostまたはTLSの名前情報が必要な場合も、接続先IPは検査済みIPに固定されたままでなければならない。

この固定と再名前解決禁止を保証できないブラウザ、HTTPクライアント、プロキシその他の取得能力は使用しない。この保証が確認できない時点で、ネットワーク接続を行わず`article_unavailable`で停止する。検査後のDNS回答変化、接続失敗、TLS/HTTP取得失敗を別IPへの再解決・再接続で回避しない。

取得は匿名の読み取りだけに限定する。Cookie、認証情報、既存セッション、フォーム送信、POST等の状態変更、または状態を変える操作を送らない。ログイン、同意、動的操作などが中心本文やClaim根拠に必要なら`article_unavailable`とする。

## 本文成立性と停止結果

中心本文、Claimを支える節、または各Claimの読者が到達できる位置が欠ける場合は、部分的な本文からClaimを作らず`article_unavailable`で停止する。停止時はArticle Analysis、Knowledge Assessment、読書推奨を返さず、平易な取得・解析不能理由だけを返す。

中心本文と根拠・位置が揃っていても、技術記事でない、または評価可能なClaimを一件も抽出できない場合は`article_not_evaluable`で停止する。空のClaim一覧、記事要約、一般的品質評価で代替してはならない。

ナビゲーション、装飾、関連リンクなど、Claim根拠にならない非本文要素だけの欠落は、具体的な内容を`content_limitations`に記録して正常に続行できる。

## Claim抽出規則

- `claim`は事実、因果、構造、判断、実践手順、訂正、または探索先として単独評価できる一命題にする。独立命題は分割するが、単独で意味を失う因果・トレードオフは最小の一単位に保つ。
- `role`は記事内の役割である。`core`は中心結論、`supporting`はその根拠、`example`は具体例、`background`は前提、`incidental`は補足であり、利用者の理解度を表さない。
- `importance`は記事主題と読書判断への寄与で決める。中心結論・安全性・重要設計判断は`high`、理解や適用を実質的に支えるものは`medium`、補助的詳細は`low`とする。roleだけから機械的に決めない。
- `location`は少なくとも節見出しまたは段落上の位置を示す。実在するURL fragmentだけを補助に使い、引用位置・anchorを捏造しない。
- `support`には命題を支える記事内根拠の種類と対象を記し、主張だけで弱い場合はその制限も記す。
- `scope`、`target_version`、`temporal_context`は記事が明示するときだけ記録し、明示されない場合はそれぞれ`なし`とする。

## 再調査時のClaim Reconciliation

初回は`初回解析`と記す。再調査では旧Analysisと新Analysisの各Claimを、`claim`、`scope`、`target_version`、`temporal_context`の四属性で照合し、各結果を次のいずれかと理由付きで記録する。

| 分類 | 意味 |
| --- | --- |
| `retained` | 四属性がすべて同一。既存Claim IDを維持する。 |
| `changed` | 対応するClaimがあるが、四属性の少なくとも一つが異なる。 |
| `added` | 新AnalysisだけにあるClaim。 |
| `removed` | 旧AnalysisだけにあるClaim。 |

同じ意味に見えるという自由な推測で`retained`にしてはならない。照合表は旧Claim ID、新Claim ID、分類、理由を明記する。`changed`、`added`、`removed`の後続評価・推奨上の扱いはこの契約では決めない。
