# UC-04: 時点情報を持つ Assertion を確認する

**目的:** Concept または Scope を手掛かりに、時点情報を持つ現行 Assertion とその根拠を確認する。

**根拠:** REQ-009、REQ-012、REQ-014、BR-004、BR-010、CON-004。

| 順序 | 操作 | 使い方 | 次の判断材料 |
| --- | --- | --- | --- |
| 1 | [`search-temporal`](../commands/search-temporal.md) | Concept または Scope を指定し、必要に応じて時点または有効期間で、時点情報を持つ現行 Assertion を調べる。 | 時点情報を持つ候補。 |
| 2 | [`get`](../commands/get.md) | 確認する候補の Assertion ID を指定する。 | 本文、Scope、Temporal Metadata。 |
| 3 | [`get-evidence`](../commands/get-evidence.md) | 必要な候補の根拠を取得する。 | Evidence 履歴。 |

**終了条件:** `search-temporal` の空結果は成功である。CLI は候補の新しさ・正しさ・優先度、またはユーザーの知識状態を結論づけない。

**利用者・前提:** Codex。候補の新しさ・正しさ・優先度は CLI が判断しない。

## 操作列とデータの受渡し例

### 1. 時点情報を持つ候補を検索する

```text
knowledge search-temporal --concept channel \
  --scope-key language --scope-value Go \
  --at 2026-08-14T00:00:00Z
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
        "temporal": {
          "valid_from": "2025-01-01T00:00:00Z",
          "valid_until": null,
          "version_scope": "Go 1.24",
          "observed_at": "2026-08-13T00:00:00Z",
          "last_verified": null
        }
      }
    ]
  }
}
```

`results[0].assertion_id` の **`asrt_01` を次へ渡す値** とする。空の `results` は成功終了する。

### 2. 候補の全 revision を確認する

```text
knowledge get --assertion-id asrt_01
```

標準出力:

```json
{
  "ok": true,
  "data": {
    "assertion_id": "asrt_01",
    "current_revision": 2,
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
      },
      {
        "revision": 2,
        "normalized_text": "unbuffered channelへのsendは対応するreceiverがreadyになるまでblockする",
        "scope": [
          {
            "key": "language",
            "value": "Go"
          }
        ],
        "temporal": {
          "valid_from": "2025-01-01T00:00:00Z",
          "valid_until": null,
          "version_scope": "Go 1.24",
          "observed_at": "2026-08-13T00:00:00Z",
          "last_verified": null
        }
      }
    ],
    "concepts": [
      {
        "concept_id": "cpt_01",
        "name": "channel",
        "aliases": []
      }
    ],
    "aliases": []
  }
}
```

`data.assertion_id` の **`asrt_01` を次へ渡す値** とする。

### 3. 根拠を確認する

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

- CLI は候補の新しさ・正しさ・優先度、知識状態を結論づけない。
- このユースケースに `search-contradictions` は含めない。矛盾候補の探索は独立した [UC-03](uc-03-contradictions.md) である。
- 詳細契約: [`search-temporal`](../commands/search-temporal.md)、[`get`](../commands/get.md)、[`get-evidence`](../commands/get-evidence.md)。
