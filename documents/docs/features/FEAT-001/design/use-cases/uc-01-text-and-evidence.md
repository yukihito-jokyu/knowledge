# UC-01: 特定の語や API 名から Assertion と根拠を確認する

**目的:** 固有名詞、API 名、Identifier、本文、Scope など既知の文字列から候補を見つけ、候補と Evidence を確認する。

**根拠:** REQ-010、REQ-012、REQ-014、BR-002、CON-004。

| 順序 | 操作 | 使い方 | 次の判断材料 |
| --- | --- | --- | --- |
| 1 | [`search-text`](../commands/search-text.md) | Codex が決めた `query` で字句候補を得る。 | `results` の Assertion ID と、どの照合 field が一致したか。 |
| 2 | [`get`](../commands/get.md) | 確認する候補の Assertion ID を指定し、現行内容・関連データを得る。 | 正規化本文、Scope、Concept、Alias、時点情報。 |
| 3 | [`get-evidence`](../commands/get-evidence.md) | 根拠を確認する候補の Assertion ID を指定する。 | Evidence 履歴。 |

**終了条件:** `search-text` が空なら、このユースケースは成功して終了する。Codex が不足状態をどう扱うかは FEAT-001 の範囲外である。候補 ID が不正または消失していれば `get` または `get-evidence` は `not_found` で終了する。

**利用者・前提:** Codex。`query` は Codex が決める。初期提供は字句照合であり、意味検索はしない。

## 操作列とデータの受渡し例

### 1. 文字列で候補を検索する

```text
knowledge search-text --query channel
```

標準出力:

```json
{
  "ok": true,
  "data": {
    "results": [
      {
        "assertion_id": "asrt_01",
        "normalized_text": "unbuffered channelへのsendはreceiverがreadyになるまでblockする",
        "revision": 1,
        "concepts": [
          {
            "concept_id": "cpt_01",
            "name": "channel"
          }
        ],
        "scope": [
          {
            "key": "language",
            "value": "Go"
          }
        ],
        "matched_fields": ["assertion_text", "concept"]
      }
    ]
  }
}
```

`results[0].assertion_id` の **`asrt_01` を次へ渡す値** とする。空の `results` は成功で、このユースケースを終了する。

### 2. 候補の内容と履歴を取得する

```text
knowledge get --assertion-id asrt_01
```

標準出力では `current_revision`、本文、Scope、時点情報、Concept、Alias が返る。Evidence は含まれない。

```json
{
  "ok": true,
  "data": {
    "assertion_id": "asrt_01",
    "current_revision": 1,
    "revisions": [
      {
        "revision": 1,
        "normalized_text": "unbuffered channelへのsendはreceiverがreadyになるまでblockする",
        "scope": [
          {
            "key": "language",
            "value": "Go"
          }
        ],
        "temporal": null
      }
    ],
    "concepts": [
      {
        "concept_id": "cpt_01",
        "name": "channel",
        "aliases": []
      }
    ],
    "aliases": [
      {
        "kind": "identifier",
        "value": "chan<-"
      }
    ]
  }
}
```

`data.assertion_id` の **`asrt_01` を次へ渡す値** とする。

### 3. 根拠の履歴を取得する

```text
knowledge get-evidence --assertion-id asrt_01
```

標準出力:

```json
{
  "ok": true,
  "data": {
    "assertion_id": "asrt_01",
    "evidence": [
      {
        "evidence_id": "evd_01",
        "kind": "user_explanation",
        "raw_text": "channelって受け側いないとsend側止まりますよね",
        "observed_at": "2026-08-13T00:00:00Z"
      }
    ]
  }
}
```

## 終了条件と境界

- 手順1が空結果なら成功終了する。CLI は不足状態や知識状態を結論づけない。
- ID が不正または消失していれば手順2・3は `not_found` で終了する。
- 詳細契約: [`search-text`](../commands/search-text.md)、[`get`](../commands/get.md)、[`get-evidence`](../commands/get-evidence.md)。
