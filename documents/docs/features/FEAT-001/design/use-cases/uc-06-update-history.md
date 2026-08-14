# UC-06: 既存 Assertion の根拠を追加・内容を訂正・旧知識を置換する

**目的:** 既存データを物理削除せず、更新の意味に応じて Evidence 追加、revision 追加、または置換 Relation の保存を行う。

**根拠:** REQ-013、REQ-014、BR-004、BR-010、NFR-006。

| 更新の意味 | 操作 | 使い方 | 履歴への影響 |
| --- | --- | --- | --- |
| 同じ Assertion を支える新しい根拠 | [`attach-evidence`](../commands/attach-evidence.md) | 対象 Assertion ID と新しい Evidence を指定する。 | Assertion は変えず Evidence を追加する。 |
| 正規化本文、Scope、または時点情報の訂正 | [`revise`](../commands/revise.md) | 対象 Assertion ID と訂正後の本文・Scope・時点情報を指定する。 | 現行 revision を更新し、過去 revision を保持する。 |
| 旧 Assertion を削除せず新 Assertion へ置き換える関係 | [`supersede`](../commands/supersede.md) | 旧 Assertion と新 Assertion を指定する。 | `supersedes` Relation を追加し、両 Assertion と Evidence を保持する。 |
| 保存済み内容・根拠の再確認 | [`get`](../commands/get.md)、[`get-evidence`](../commands/get-evidence.md) | 更新対象または更新後の Assertion ID を指定する。 | 履歴と結果を確認するだけで変更しない。 |

**終了条件:** 対象 ID の不存在は `not_found`、重複 Evidence・同一 revision・自己参照／循環／重複 Relation は `conflict`、保存失敗は rollback 後の `storage_error` で終了する。

**利用者・前提:** Codex が更新の意味に応じて三つの mutation の一つを選ぶ。CLI は選択しない。

## 操作列とデータの受渡し例

### A. 同じ Assertion を支える Evidence を追加する

```text
knowledge attach-evidence --assertion-id asrt_01 \
  --evidence-kind correction \
  --evidence-text 'receiverがreadyになるまでblockする理解です' \
  --evidence-observed-at 2026-08-13T00:00:00Z
```

標準出力:

```json
{
  "ok": true,
  "data": {
    "assertion_id": "asrt_01",
    "evidence_id": "evd_02"
  }
}
```

`data.assertion_id` の **`asrt_01` を次へ渡し**、`knowledge get-evidence` で Evidence 履歴に `evd_02` が含まれることを確認できる。

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
      },
      {
        "evidence_id": "evd_02",
        "kind": "correction",
        "raw_text": "receiverがreadyになるまでblockする理解です",
        "observed_at": "2026-08-13T00:00:00Z"
      }
    ]
  }
}
```

`data.assertion_id` は、確認した Assertion の ID `asrt_01` である。

### B. 本文・Scope・時点情報を訂正する

```text
knowledge revise --assertion-id asrt_01 \
  --normalized-text 'unbuffered channelへのsendは対応するreceiverがreadyになるまでblockする' \
  --scope-key language --scope-value Go \
  --version-scope 'Go 1.24' --observed-at 2026-08-13T00:00:00Z
```

標準出力:

```json
{
  "ok": true,
  "data": {
    "assertion_id": "asrt_01",
    "previous_revision": 1,
    "revision": 2
  }
}
```

`data.assertion_id` の **`asrt_01` を `knowledge get` へ渡す**。応答の `current_revision` は `2`、`revisions` は revision 1 と 2 の両方を含む。

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
          "valid_from": null,
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
    "aliases": [
      {
        "kind": "identifier",
        "value": "chan<-"
      }
    ]
  }
}
```

`data.assertion_id` は、確認した Assertion の ID `asrt_01` である。

### C. 新しい Assertion が古い Assertion を置換する関係を保存する

```text
knowledge supersede --superseded-assertion-id asrt_01 \
  --replacement-assertion-id asrt_03
```

標準出力:

```json
{
  "ok": true,
  "data": {
    "relation_id": "rel_03",
    "relation_type": "supersedes",
    "superseded_assertion_id": "asrt_01",
    "replacement_assertion_id": "asrt_03"
  }
}
```

`data.superseded_assertion_id` の **`asrt_01`** と `data.replacement_assertion_id` の **`asrt_03` を次へ渡す値** とする。

### D. 旧・新の Assertion がともに残っていることを確認する

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
          "valid_from": null,
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
    "aliases": [
      {
        "kind": "identifier",
        "value": "chan<-"
      }
    ]
  }
}
```

```text
knowledge get --assertion-id asrt_03
```

標準出力:

```json
{
  "ok": true,
  "data": {
    "assertion_id": "asrt_03",
    "current_revision": 1,
    "revisions": [
      {
        "revision": 1,
        "normalized_text": "selectで複数の通信操作を待てる",
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
        "concept_id": "cpt_03",
        "name": "select",
        "aliases": []
      }
    ],
    "aliases": []
  }
}
```

二つの成功応答があるため、`supersede` は旧 Assertion を削除せず、置換 Relation を追加する操作だと確認できる。

## 終了条件と境界

- 対象 ID の不存在は `not_found`、重複 Evidence・同一 revision・自己参照／循環／重複 Relation は `conflict`、保存失敗は rollback 後の `storage_error` で終了する。
- 各成功応答を受けた後も、別の mutation を選ぶかどうかは Codex が判断する。
- 詳細契約: [`attach-evidence`](../commands/attach-evidence.md)、[`revise`](../commands/revise.md)、[`supersede`](../commands/supersede.md)、[`get`](../commands/get.md)、[`get-evidence`](../commands/get-evidence.md)。
