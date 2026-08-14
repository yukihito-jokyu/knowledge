# UC-05: 新しい Assertion を根拠とともに登録する

**目的:** Codex が既存 Assertion との照合を終え、新しい Assertion と判断した場合に、初期 Evidence と関連情報を一貫して保存する。

**根拠:** REQ-009、REQ-013、REQ-014、BR-004、BR-008。

| 順序 | 操作 | 使い方 | 結果 |
| --- | --- | --- | --- |
| 1 | [`search-text`](../commands/search-text.md) または [`search-concept`](../commands/search-concept.md) | Codex が既存候補を調べる。必要に応じて UC-02 の Relation 検索も行う。 | 更新操作の選択に必要な候補。 |
| 2 | [`create`](../commands/create.md) | Codex が「新規」と判断した Assertion、初期 Evidence、任意の Concept・Scope・Relation・時点情報を一つのoption列で指定する。 | 新しい Assertion と初期 revision が同一 transaction で保存される。 |
| 3 | [`get`](../commands/get.md) または [`get-evidence`](../commands/get-evidence.md) | 必要な場合だけ、成功 response の Assertion ID で保存結果を確認する。 | 保存済みの Assertion または Evidence。 |

**終了条件:** 同一の正規化本文・Scope の現行 Assertion は `conflict`、指定 Relation target がない場合は `not_found` で終了する。Codex は失敗後に自動で `attach-evidence` 等へ切り替えず、再照合して次の操作を判断する。

**利用者・前提:** Codex が `create` を選択する。CLI は意味的な同一性や、失敗後に別 mutation へ切り替える判断をしない。

## 操作列とデータの受渡し例

### 1. 既存候補を検索する

```text
knowledge search-text --query select
```

標準出力（新規と判断できる空結果の例）:

```json
{
  "ok": true,
  "data": {
    "results": []
  }
}
```

空結果そのものは新規登録の許可ではない。Codex が根拠に基づいて新規と判断した場合だけ次へ進む。

### 2. Assertion と初期 Evidence を作成する

```text
knowledge create \
  --normalized-text 'selectで複数の通信操作を待てる' \
  --scope-key language --scope-value Go \
  --concept select \
  --evidence-kind user_explanation \
  --evidence-text 'selectを使うと複数channelを待てます' \
  --evidence-observed-at 2026-08-13T00:00:00Z
```

標準出力:

```json
{
  "ok": true,
  "data": {
    "assertion_id": "asrt_03",
    "revision": 1,
    "evidence_ids": ["evd_03"],
    "concepts": [
      {
        "concept_id": "cpt_03",
        "name": "select"
      }
    ],
    "relation_ids": []
  }
}
```

`data.assertion_id` の **`asrt_03` を次へ渡す値** とする。

### 3. 保存結果を確認する

```text
knowledge get-evidence --assertion-id asrt_03
```

標準出力:

```json
{
  "ok": true,
  "data": {
    "assertion_id": "asrt_03",
    "evidence": [
      {
        "evidence_id": "evd_03",
        "kind": "user_explanation",
        "raw_text": "selectを使うと複数channelを待てます",
        "observed_at": "2026-08-13T00:00:00Z"
      }
    ]
  }
}
```

`data.evidence[0].evidence_id` は作成時に返った `evd_03` と照合できる。必要なら同じ `asrt_03` を `knowledge get` に渡して Assertion の内容も確認する。

## 終了条件と境界

- 同一の正規化本文・Scope の現行 Assertion は `conflict`、指定 Relation target がない場合は `not_found` で終了する。
- 詳細契約: [`search-text`](../commands/search-text.md)、[`create`](../commands/create.md)、[`get`](../commands/get.md)、[`get-evidence`](../commands/get-evidence.md)。
