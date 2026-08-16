# Knowledge CLI 操作リファレンス

Knowledge Updateが消費するFEAT-001 v1の既存CLI契約を参照する。全操作は `knowledge <operation> <named-options>`。成功はstdoutに一つの`{ "ok": true, "data": ... }`、失敗はstderrに一つの`{ "ok": false, "error": ... }`、成功exitは0である。JSON、option、保存先をこのSkillで拡張しない。

## 共通error / exit

| error.code | exit | Updateでの扱い |
| --- | ---: | --- |
| `validation_error` | 2 | CLI errorで停止。再試行しない。 |
| `not_found` | 3 | CLI errorで停止。IDを推測しない。 |
| `conflict` | 4 | 既存契約どおり扱う。attach/createは完全一致を読取確認できる時だけ既適用、supersedeはoutcome_unknown。 |
| `storage_error` / `internal_error` | 1 | CLI errorで停止。 |
| response前の無出力exit 130 | 130 | error JSONなしのcanceled。 |

成功JSONの形、stdout/stderrの混在、data形、複数行を検査する。不一致は`protocol_error`、error_code nullである。

## 読取操作

```text
knowledge search-text --query '<query>'
knowledge get --assertion-id '<assertion-id>'
knowledge get-evidence --assertion-id '<assertion-id>'
```

`search-text`はCandidateの各`search_queries`を一回ずつ実行する。`data.results[]`の`assertion_id`を初出順に重複除去し、各一意Assertion IDへ意味照合用の`get`を初出順に一回行う。`get-evidence`は、そのAssertionがattach、revise、supersede、duplicate-skip、conflict確認の意味的候補になった場合だけ一回行い、意味的不一致と判断したIDには行わない。空results、Evidenceなしはexit 0の正常結果であり、失敗としない。IDは検索結果以外から作らない。create候補に既存Assertionがない場合はget-evidence不要であり、create conflict時もsearchで得たID以外を推測して確認しない。

## `create`

```text
knowledge create --normalized-text '<text>' \
  [--scope-key '<key>' --scope-value '<value>']... \
  [--concept '<name>' [--concept-alias '<alias>']...]... \
  [--alias-kind '<api_name|identifier>' --alias-value '<value>']... \
  [--relation-type '<type>' --relation-target-kind '<assertion|concept>' \
   --relation-target-id '<id>']... \
  [--valid-from '<RFC3339 UTC>'] [--valid-until '<RFC3339 UTC>'] \
  [--version-scope '<scope>'] [--observed-at '<RFC3339 UTC>'] \
  [--last-verified '<RFC3339 UTC>'] \
  --evidence-kind '<kind>' --evidence-text '<raw text>' \
  --evidence-observed-at '<RFC3339 UTC>'
```

UpdateのCandidateには構造化Concept、Alias、Identifier、Relation fieldがないため、createは`--normalized-text`、明示されたscope/temporal、Evidence groupだけを渡す。`search_queries`、Evidence原文、scope、strength、その他metadataからConcept/Alias/Identifier/Relation optionを推測してはならない。これらの構造化fieldが将来追加された場合は、別の承認済み契約とspec recheckが必要である。既存CLI自体が許可するoption一覧はFEAT-001契約を参照するが、今回のCandidateから渡すoptionを拡張しない。成功dataは`assertion_id`、`revision`、`evidence_ids`、`relation_ids`、`concepts`を持つ。`conflict`で既存AssertionとEvidenceの完全一致を確認できない場合は失敗で停止する。

### create conflictの確認と分岐

createが`conflict`（stderrのerror JSON、exit 4）を返したら、search-textの結果に含まれるAssertion IDだけを候補にする。各候補へ`get --assertion-id`を一回行い、normalized text、scope、temporal、identityを比較する。候補が一意に一致した場合だけ、そのIDへ`get-evidence --assertion-id`を一回行う。検索結果にないID、Concept/Alias/Identifier/Relationを原文やqueryから推測してはならない。

- 既存AssertionとEvidenceのkind/raw_text/observed_atが完全一致する場合は、追加mutationなしで、既存enumの`action:create`、`execution_status:applied`、rationaleの`already applied`として記録する。`cli_operations`の順序はcreate(conflict)、get、get-evidenceであり、確認IDは各result_idsとsearch_evidenceに残す。
- Assertionが一意に一致し、対象Evidenceが未付与の場合だけ、そのAssertion IDへ上記`attach-evidence`を呼ぶ。Decisionは`action:attach_evidence`、`execution_status:applied`とし、create(conflict)、get、get-evidence、attach-evidenceの結果を順に残す。
- 一意に一致しない、identityが不一致、またはCandidateに構造化fieldがないため対象Concept/Alias等を安全に特定できない場合はattachを呼ばない。`failure_reason:cli_error`の`failed`として停止し、後続CandidateのCLIを呼ばない。新しい`already_applied` actionは公開しない。

## `attach-evidence`

```text
knowledge attach-evidence --assertion-id '<assertion-id>' \
  --evidence-kind '<kind>' --evidence-text '<raw text>' \
  --evidence-observed-at '<RFC3339 UTC>'
```

Evidence kindは`user_explanation`、`user_reasoning`、`user_code`、`self_report`、`concept_recognition`、`correction`。成功dataは`assertion_id`、`evidence_id`。同一Assertion内のkind・raw_text・observed_at重複は`conflict`である。完全一致を`get-evidence`で確認できる場合だけ既適用として記録し、それ以外はCLI errorで停止する。

## `revise`

```text
knowledge revise --assertion-id '<assertion-id>' \
  --normalized-text '<new text>' \
  [--scope-key '<key>' --scope-value '<value>']... \
  [--valid-from '<RFC3339 UTC>'] [--valid-until '<RFC3339 UTC>'] \
  [--version-scope '<scope>'] [--observed-at '<RFC3339 UTC>'] \
  [--last-verified '<RFC3339 UTC>']
```

成功dataは`assertion_id`、`previous_revision`、`revision`。旧revisionとEvidenceは削除されない。訂正Evidenceを残す必要がある場合は、成功後に上記`attach-evidence`を別操作として実行し、後段失敗・中断を`partially_applied`にする。`revise`のconflictは自動再実行しない。

## `supersede`

```text
knowledge supersede \
  --superseded-assertion-id '<old-assertion-id>' \
  --replacement-assertion-id '<new-assertion-id>'
```

`create`成功で得た新Assertion IDを`--replacement-assertion-id`へ渡し、旧IDを`--superseded-assertion-id`へ渡す。成功dataは`relation_id`、`relation_type: supersedes`、両Assertion ID。旧Assertionを削除しない。create成功後のsupersede失敗・中断は、未適用が明らかなら`partially_applied`、`conflict`なら`failed/outcome_unknown`とする。

## 禁止事項

- `search-semantic`、Embedding、Vector Index、新しい検索・更新operationを呼ばない。
- 未定義option、位置引数、外部生成ID、公開JSON fieldを使わない。
- 二段操作を一つの新規transactionへ包まない。
- 補償削除、履歴削除、自動再試行、再開ledgerを作らない。
