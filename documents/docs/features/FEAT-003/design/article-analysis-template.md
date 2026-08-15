# Article Analysis 成果物契約

Article Analysisは、取得済みの記事をユーザー知識に依存しない比較可能なClaim群へ変換する、URL評価実行内のMarkdown成果物である。記事全文の要約やKnowledge Assessmentではない。初期提供では永続化しない。

```markdown
# Article Analysis

## Article

- source_url: <利用者が指定したURL>
- resolved_url: <取得後の最終URL>
- title: <記事タイトル。取得できなければ「不明」>
- retrieved_at: <評価実行時刻>
- content_limitations: <欠落した本文、ログイン要求、動的表示など。なければ「なし」>

## Overview

<記事の主題と、Claim群を読むために必要な最小の文脈。>

## Claims

### CLM-001

- claim: <ユーザー知識と独立に評価できる一つの命題>
- role: <core | supporting | example | background | incidental>
- importance: <high | medium | low>
- location: <節見出し、順序、利用可能ならアンカー。読者が到達できる位置>
- support: <本文で命題を支える説明、実装例、失敗例、測定、一次資料参照など>
- scope: <任意。Claimが成り立つ対象・条件。なければ「なし」>
- target_version: <任意。対象製品・仕様・実装のversion。なければ「なし」>
- temporal_context: <任意。Claimの時点・有効期間。なければ「なし」>

## Coverage Notes

<Claimにしなかった反復、装飾、独立評価不能な断片と、その理由。>

## Claim Reconciliation

<再調査時だけ、旧Claim ID、新Claim ID、retained | changed | added | removed、理由を記す。初回は「初回解析」。>
```

## 記録規則

- `CLM-###` は単一実行内だけで有効な順序付き識別子であり、Article Claimの永続IDではない。
- `claim` は一つの事実、因果、構造、判断、実践手順、訂正、または探索先として評価できる粒度にする。複数の独立命題は分割するが、単独で意味を失う因果・トレードオフは必要な最小単位として一つに保つ。
- `role` は記事における役割であり、利用者が理解している度合いではない。`core`は中心結論、`supporting`は中心結論を支える知識、`example`は具体例、`background`は前提、`incidental`は補足である。
- `importance` は記事の主題と読書判断への寄与を表す。`high`は中心結論・安全性・重要な設計判断など、`medium`は理解または適用を実質的に支えるもの、`low`は補助的詳細とする。役割だけから機械的に決めない。
- `location` は少なくとも見出しまたは段落上の位置を示す。URL fragmentが利用できる場合だけ補助的に添える。生成した引用位置や存在しないアンカーを捏造しない。
- `support` は記事内部の根拠の種類と対象を説明する。Supportが弱い・主張だけである場合はその事実を記録する。
- 記事の取得・解析から、利用者が読んだ・理解した・知っているというEvidenceを生成してはならない。
- `content_limitations` は記事の評価可能性を示す。ナビゲーション等、Claimの根拠にならない要素だけの欠落は記録してよい。一方、中心本文またはClaimの根拠となる節・位置情報が欠ける場合は `article_unavailable` として終了し、部分的な本文からClaimを評価しない。
- 再調査では、`claim`、`scope`、`target_version`、`temporal_context` がすべて同一のClaimだけが `retained` となり、既存のClaim IDを維持する。それ以外は `changed`、`added`、`removed` として照合表へ記録する。
