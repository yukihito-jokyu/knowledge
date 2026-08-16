# Knowledge CLI 操作リファレンス

このファイルは `knowledge-search` Skillが使う既存の読み取りCLI契約である。コマンドはすべて `knowledge [--store <absolute-path>] <operation> <named-options>` の形で実行する。`--store`は呼出側が明示Store実行コンテキストを渡す場合だけoperation前に一度置き、通常Skillは自動追加しない。標準出力は成功JSONだけ、標準エラーは失敗JSONだけを読む。成功は exit 0 と `{ "ok": true, "data": ... }`、失敗は `{ "ok": false, "error": { "code": ..., "message": ..., "field": ... } }` である。

| error.code | exit code | 扱い |
| --- | ---: | --- |
| `validation_error` | 2 | technical failure。入力を直して再試行しない。 |
| `not_found` | 3 | technical failure。ID由来の呼出しなら探索経路を止める。 |
| `storage_error` / `internal_error` | 1 | technical failure。再試行しない。 |
| response開始前の exit 130・無出力 | 130 | canceled。Assessmentを作らない。 |

空配列、未登録Concept、Relationなし、Evidenceなしは exit 0 の正常結果であり、technical failureではない。

## 操作一覧

### `search-text`

```text
knowledge search-text --query '<query>'
```

`--query` は空白だけでない文字列を一回指定する。字句照合であり、意味検索ではない。最初の呼出しに必ず使う。`data.results` は常に配列で、各候補の `assertion_id`、`normalized_text`、`revision`、`concepts[]`（`concept_id`, `name`）、`scope[]`（`key`, `value`）、`matched_fields[]` を返す。`matched_fields` は `assertion_text`、`concept`、`alias`、`scope_key`、`scope_value` のいずれかである。

### `search-concept`

```text
knowledge search-concept --concept '<concept-or-alias>'
```

`--concept` は空白だけでない文字列を一回指定する。完全一致のConcept名またはAliasを検索する。`data.concept` は一致時に `{ concept_id, name }`、未登録時は `null`。`data.results[]` は `assertion_id`、`normalized_text`、`revision`、`scope[]` を返す。未登録は空の成功結果である。

### `get`

```text
knowledge get --assertion-id '<assertion-id>'
```

`--assertion-id` は空でないIDを一回指定する。必ず既出の候補またはRelation結果から得たIDだけを渡す。`data` は `assertion_id`、`current_revision`、全 `revisions[]`（`normalized_text`、`scope[]`、`temporal`）、`concepts[]`、`aliases[]` を返す。`temporal` は `valid_from`、`valid_until`、`version_scope`、`observed_at`、`last_verified` を持つか `null` である。

### `get-evidence`

```text
knowledge get-evidence --assertion-id '<assertion-id>'
```

`--assertion-id` は既出IDだけを一回指定する。`data.evidence[]` は `evidence_id`、`kind`、`raw_text`、`observed_at` を返す。Evidenceなしは空の成功結果である。`kind` は `user_explanation`、`user_reasoning`、`user_code`、`self_report`、`concept_recognition`、`correction` のいずれかである。

### `search-related`

```text
knowledge search-related --seed-kind assertion --seed-id '<assertion-id>' \
  --relation-type prerequisite
```

`--seed-kind` は `assertion` または `concept`、`--seed-id` は対応する既出IDを一回指定する。`--relation-type` は省略可能で、繰返す場合は `related_to`、`prerequisite`、`causes`、`contributes_to`、`contradicts`、`supersedes` のいずれかだけを使う。`data.results[]` は `relation_id`、`relation_type`、`direction`、`target` を返す。`target.kind` は `assertion` または `concept`、`target.id` は次の操作へ渡せるIDである。Relationなしは空の成功結果である。

### `search-contradictions`

```text
knowledge search-contradictions --assertion-id '<assertion-id>'
# または
knowledge search-contradictions --concept '<concept-or-alias>'
```

`--assertion-id` と `--concept` のどちらか一方だけを指定する。`data.results[]` は `relation_id`、`direction`、`seed`、`target`を返し、`target.kind` は常に `assertion` である。これは保存済みの矛盾候補であり、正誤を確定しない。Assertion不存在またはConcept未登録は空の成功結果である。

### `search-temporal`

```text
knowledge search-temporal --concept '<concept-or-alias>' \
  --scope-key '<key>' --scope-value '<value>' \
  --at '2026-08-14T00:00:00Z'
```

`--concept` または一組以上の `--scope-key` / `--scope-value` が必須である。Scopeは同じ数だけ、同じkeyを重複させずに対で指定する。`--at` はRFC 3339 UTCで、`--valid-from` と `--valid-until` の組とは併用しない。期間で探すときは `--valid-from` と `--valid-until` をRFC 3339 UTCで対にし、開始を終了以前にする。時点条件はTarget Claimに明示された値だけを使い、現在日時・評価日時・一般知識から作らない。`data.results[]` は `assertion_id`、`normalized_text`、`temporal`を返す。時点情報を持つ候補だけを機械的に絞り込み、新しさや正しさは判断しない。

## 呼出し時の共通規則

- named optionだけを使い、位置引数・未定義option・重複した単一値optionを渡さない。
- IDを外部から作らない。検索・取得結果から得たIDだけを後続操作へ渡す。
- 成功JSONでも `data` の形がこのリファレンスと違う場合はJSONプロトコル不整合のtechnical failureとして停止する。
- stderrの失敗JSON、exit code、stdoutを同じ呼出しの結果として扱う。成功・失敗の両方を出力したり、複数行JSONを連結したりしない。
