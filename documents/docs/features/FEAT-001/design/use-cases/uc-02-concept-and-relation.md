# UC-02: Concept または Relation から関連する Assertion をたどる

**目的:** 直接の文字列一致だけでは足りないときに、Concept を入口として候補を探し、指定した検索起点との Relation を両方向から確認する。

**根拠:** REQ-009、REQ-012、REQ-014、CON-004。

| 順序 | 操作 | 使い方 | 次の判断材料 |
| --- | --- | --- | --- |
| 1 | [`search-concept`](../commands/search-concept.md) | Codex が決めた Concept 名または Alias を完全照合する。 | Concept に紐付く現行 Assertion。 |
| 2 | [`search-related`](../commands/search-related.md) | Assertion または Concept を検索起点として、必要な Relation 種別を指定する。 | 起点に対する相手と保存済みの向き。 |
| 3 | [`get`](../commands/get.md) | 調べる候補の Assertion ID を指定する。 | 個別 Assertion の詳細。 |
| 4 | [`get-evidence`](../commands/get-evidence.md) | 必要な候補の根拠を取得する。 | Evidence 履歴。 |

**終了条件:** `search-concept` の空結果は成功である。`search-related` は、指定した検索起点が存在しないときだけ `not_found` で終了し、Relation がない場合は成功の空結果となる。

**利用者・前提:** Codex。Concept 名または Alias と、Relation をたどる起点を Codex が選ぶ。`target` は常に検索起点の相手、`direction` は保存済み Relation の向きである。

## 操作列とデータの受渡し例

### 1. Concept から候補 Assertion を得る

```text
knowledge search-concept --concept channel
```

標準出力:

```json
{
  "ok": true,
  "data": {
    "concept": {
      "concept_id": "cpt_01",
      "name": "channel"
    },
    "results": [
      {
        "assertion_id": "asrt_01",
        "normalized_text": "unbuffered channelへのsendはreceiverがreadyになるまでblockする",
        "revision": 1,
        "scope": [
          {
            "key": "language",
            "value": "Go"
          }
        ]
      }
    ]
  }
}
```

Relation の起点にする `results[0].assertion_id`、すなわち **`asrt_01` を次へ渡す値** とする。Concept 未登録なら `concept: null` と空の `results` が返り成功終了する。

### 2. Assertion を起点に Relation をたどる

```text
knowledge search-related --seed-kind assertion --seed-id asrt_01 \
  --relation-type causes
```

標準出力:

```json
{
  "ok": true,
  "data": {
    "results": [
      {
        "relation_id": "rel_01",
        "relation_type": "causes",
        "direction": "outgoing",
        "target": {
          "kind": "assertion",
          "id": "asrt_02",
          "normalized_text": "goroutineがreceiverとして起動している"
        }
      }
    ]
  }
}
```

`results[0].target.id` の **`asrt_02` を次へ渡す値** とする。空の `results` は Relation がないという成功結果である。

### 3. 相手 Assertion と根拠を確認する

```text
knowledge get --assertion-id asrt_02
```

標準出力:

```json
{
  "ok": true,
  "data": {
    "assertion_id": "asrt_02",
    "current_revision": 1,
    "revisions": [
      {
        "revision": 1,
        "normalized_text": "goroutineがreceiverとして起動している",
        "scope": [],
        "temporal": null
      }
    ],
    "concepts": [],
    "aliases": []
  }
}
```

`data.assertion_id` の **`asrt_02` を次へ渡す値** とする。

### 4. 相手 Assertion の根拠を取得する

```text
knowledge get-evidence --assertion-id asrt_02
```

標準出力:

```json
{
  "ok": true,
  "data": {
    "assertion_id": "asrt_02",
    "evidence": [
      {
        "evidence_id": "evd_04",
        "kind": "user_code",
        "raw_text": "go func() { <-ch }()",
        "observed_at": "2026-08-13T00:00:00Z"
      }
    ]
  }
}
```

## 終了条件と境界

- `search-related` は検索起点がない場合のみ `not_found` で終了し、Relation がない場合は成功の空結果である。
- CLI は Relation の意味や検索継続を判断しない。
- 詳細契約: [`search-concept`](../commands/search-concept.md)、[`search-related`](../commands/search-related.md)、[`get`](../commands/get.md)、[`get-evidence`](../commands/get-evidence.md)。
