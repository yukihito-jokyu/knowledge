# UC-03: 矛盾候補を確認する

**目的:** 保存済みの `contradicts` Relation から、矛盾候補として記録された Assertion を確認する。

**根拠:** REQ-009、REQ-012、REQ-014、BR-004、BR-010、CON-004。

| 順序 | 操作 | 使い方 | 次の判断材料 |
| --- | --- | --- | --- |
| 1 | [`search-contradictions`](../commands/search-contradictions.md) | Codex が決めた selector を一つ指定し、保存済みの `contradicts` Relation を調べる。 | 矛盾候補の Relation と対象。 |
| 2 | [`get`](../commands/get.md) | 比較する候補を個別に取得する。 | 本文、Scope、Temporal Metadata。 |
| 3 | [`get-evidence`](../commands/get-evidence.md) | 比較する候補の Evidence を取得する。 | 根拠の履歴。 |

**終了条件:** `search-contradictions` の空結果は成功である。CLI は矛盾の意味、どちらが正しいか、またはユーザーの知識状態を結論づけない。

**利用者・前提:** Codex。`target` は常に、その結果の `seed` と反対側の Assertion である。`direction` は保存済み Relation の向きを残す。Concept selector に一つの Relation の両端が一致する場合は、`seed` ごとに別の結果として返る。矛盾の正誤やどちらが正しいかは、CLI でなく Codex が判断する。

## 操作列とデータの受渡し例

### 1. 矛盾候補の Relation を取得する

```text
knowledge search-contradictions --assertion-id asrt_01
```

標準出力:

```json
{
  "ok": true,
  "data": {
    "results": [
      {
        "relation_id": "rel_02",
        "direction": "outgoing",
        "seed": {
          "kind": "assertion",
          "id": "asrt_01"
        },
        "target": {
          "kind": "assertion",
          "id": "asrt_09",
          "normalized_text": "unbuffered channelへのsendはblockしない"
        }
      }
    ]
  }
}
```

`results[0].target.id` の **`asrt_09` を次へ渡す値** とする。`target` は常にこの結果の `seed` の相手なので、保存上のRelation方向を判定する必要はない。空の `results` は矛盾候補が未登録という成功結果である。

Conceptをselectorにした場合、同じRelationの両端がそのConceptに属していれば、同じ `relation_id` が異なる `seed` と `target` の組で2件返る。これは重複ではなく、どちらのAssertionから相手を確認するかを区別する結果である。

### 2. 矛盾候補の内容を取得する

```text
knowledge get --assertion-id asrt_09
```

標準出力:

```json
{
  "ok": true,
  "data": {
    "assertion_id": "asrt_09",
    "current_revision": 1,
    "revisions": [
      {
        "revision": 1,
        "normalized_text": "unbuffered channelへのsendはblockしない",
        "scope": [
          {
            "key": "language",
            "value": "Go"
          }
        ],
        "temporal": null
      }
    ],
    "concepts": [],
    "aliases": []
  }
}
```

`data.assertion_id` の **`asrt_09` を次へ渡す値** とする。

### 3. 根拠を取得する

```text
knowledge get-evidence --assertion-id asrt_09
```

標準出力:

```json
{
  "ok": true,
  "data": {
    "assertion_id": "asrt_09",
    "evidence": [
      {
        "evidence_id": "evd_09",
        "kind": "user_explanation",
        "raw_text": "sendは待たないと思う",
        "observed_at": "2026-08-12T00:00:00Z"
      }
    ]
  }
}
```

`data.evidence` は候補を支える Evidence 履歴であり、CLI は矛盾の結論を返さない。

## 終了条件と境界

- `search-contradictions` の空結果は成功である。
- このユースケースに `search-temporal` は含めない。時点情報の探索は独立した [UC-04](uc-04-temporal.md) である。
- 詳細契約: [`search-contradictions`](../commands/search-contradictions.md)、[`get`](../commands/get.md)、[`get-evidence`](../commands/get-evidence.md)。
