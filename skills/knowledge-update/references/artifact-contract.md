# Update成果物契約

この契約は、Candidate Knowledgeを既存Knowledge CLIへ適用する一時Markdown成果物の論理契約である。FEAT-001のCLI JSON、保存先、DB、migrationを再定義しない。

## Candidate入力

入力は`knowledge-acquisition`が返す一つのCandidate Knowledgeである。各Candidateは次を持つ。

| Field | 型・必須性 | 契約 |
| --- | --- | --- |
| `episode_id` / `candidate_id` | opaque string、必須 | Candidate Knowledgeとの追跡参照。再開・再実行の鍵にしない。 |
| `source_ordinal` / `candidate_ordinal` | integer、必須 | 処理順。source昇順、同一source内candidate昇順。 |
| `evidence_raw_text` | string、必須 | CLIへ渡す完全なユーザーEvidence原文。省略・要約しない。 |
| `evidence_kind` | enum、必須 | `user_explanation`、`user_reasoning`、`user_code`、`self_report`、`concept_recognition`、`correction`。 |
| `strength` | enum、必須 | `strong`、`moderate`、`weak`。Evidenceから導出する判断用派生値であり、CLI保存fieldではない。 |
| `observed_at` | RFC 3339 UTC、必須 | Evidenceの観測時刻。 |
| `proposed_assertion` | string、必須 | 独立評価可能な候補Assertion。 |
| `search_queries` | 1件以上のstring list、必須 | 先頭は`proposed_assertion`そのもの。以後は原文に明示されたConcept・Alias・Identifierのみを出現順で持ち、完全一致文字列を重複させない。 |
| `scope` | key/value list、必須 | 明示されたScopeだけ。空可、key重複不可。 |
| `temporal` | objectまたはnull、必須 | 明示された時点情報だけ。 |
| `source_excerpt` / `extraction_rationale` | string、必須 | Candidateの確認と抽出理由。 |

### Evidence kindからstrengthへの導出

| `evidence_kind` | `strength` |
| --- | --- |
| `user_explanation`、`user_reasoning`、`user_code`、`correction` | `strong` |
| `self_report` | `moderate` |
| `concept_recognition` | `weak` |

`strength`が欠落、上記enum外、または`evidence_kind`と導出表が一致しないCandidateは入力不成立である。修正・推測・CLI呼出しをせず停止する。`strength`はCandidate時点の説明用派生値であり、`create`、`attach-evidence`、`revise`、`supersede`のoptionやKnowledge Storeの保存fieldではない。

## Update Decision

各Candidateは、成功・skip・停止のいずれでも一件のDecisionを持つ。除外入力にはDecisionを作らない。

| Field | 型・必須性 | 契約 |
| --- | --- | --- |
| `episode_id` / `candidate_id` | opaque string、必須 | 入力への追跡参照。 |
| `action` | enum、必須 | `create`、`attach_evidence`、`revise`、`supersede`、`skip`、`not_started`。 |
| `target_assertion_id` | opaque stringまたはnull、必須 | attach/reviseの対象、supersedeの旧Assertion。それ以外はnull。 |
| `replacement_assertion_id` | opaque stringまたはnull、必須 | supersedeのcreate成功後の新Assertion。それ以外はnull。 |
| `rationale` | string、必須 | 意味対応、更新、skip、停止の根拠。 |
| `search_evidence` | Assertion/Evidence ID list、必須 | 判断に使った既存結果。空可。 |
| `execution_status` | enum、必須 | `not_started`、`applied`、`skipped`、`partially_applied`、`failed`、`canceled`。 |
| `failure_reason` | enumまたはnull、必須 | failed時のみ`cli_error`、`protocol_error`、`outcome_unknown`。その他はnull。 |
| `cli_operations` | ordered list、必須 | 実行した既存CLI operation result。空可。 |

`skip`は`skipped`かつ操作一覧空である。選択前の検索・取得失敗は`not_started` actionかつ`failed`/`canceled`、先行Candidate停止による未処理は`not_started`/`not_started`である。後段が未適用と分かる二段操作は`partially_applied`、`supersede`の`conflict`は`failed`/`outcome_unknown`とする。

## CLI operation result

各要素は実行順に次を持つ。

| Field | 型・必須性 | 契約 |
| --- | --- | --- |
| `operation` | 既存operation名、必須 | `search-text`、`get`、`get-evidence`、`create`、`attach-evidence`、`revise`、`supersede`のいずれか。 |
| `query` | stringまたはnull、必須 | search-textならその一回のquery、その他はnull。 |
| `status` | enum、必須 | `applied`、`failed`、`canceled`。読取成功も`applied`と記録する。 |
| `result_ids` | string list、必須 | 成功時のAssertion/Evidence/Relation ID。失敗・中断は空。 |
| `error_code` | enumまたはnull、必須 | 既存CLI error code。protocol errorとcanceledはnull。 |
| `exit_code` | integer、必須 | 実測exit code。canceledは130。 |

## Update Result全体

```markdown
# Update Result

## Episode
- episode_id: <opaque id>
- completed_at: <RFC 3339 UTC>

## Overall Status
`completed | failed | canceled | partially_applied`

## Decisions
### Candidate 1
- episode_id: <opaque id>
- candidate_id: <opaque id>
- source_ordinal: <integer>
- candidate_ordinal: <integer>
- action: <enum>
- target_assertion_id: <id or null>
- replacement_assertion_id: <id or null>
- rationale: <根拠>
- search_evidence: [<id>]
- execution_status: <enum>
- failure_reason: <enum or null>
- cli_operations:
  - operation: <existing operation>
    query: <query or null>
    status: <enum>
    result_ids: [<id>]
    error_code: <code or null>
    exit_code: <integer>

## Persistence
- update_result: none
- knowledge_store_compensation_delete: none
- resume_ledger: none
```

全Candidateを入力順に含める。Candidateゼロは`Decisions: []`と`Overall Status: completed`であり、Candidateに紐付かないskipを作らない。最初の失敗・中断・部分適用以降は、全後続Candidateを`action: not_started`、`execution_status: not_started`、ID null、failure_reason null、検索・操作空で記録する。Candidate入力のstrength欠落・不正値はこの処理へ入る前に停止し、Update Decisionを捏造しない。

全体状態は、部分適用があれば`partially_applied`、先行mutationなしのexit 130なら`canceled`、技術失敗または結果不明なら`failed`、全Decisionがapplied/skippedなら`completed`である。

## 操作分類

| 状況 | Decision |
| --- | --- |
| 新しい命題、既存対応なし | `create` |
| 同一Assertion、未記録Evidence | `attach_evidence` |
| 同一identityの本文・Scope・Temporal訂正 | `revise`、訂正Evidenceが必要なら続けてattach |
| 別identityが旧Assertionを置換 | `create`、続けて`supersede` |
| 許可済みCandidateが不十分、根拠不十分、同じEvidence済み | `skip` |

### create conflictのDecision

`create`のconflictは新しいaction enumを追加せず、次の既存Decisionで表現する。

- create conflict後、searchで得たIDだけを`get`し、normalized text、scope、temporal、identityが一意に一致したAssertionだけへ`get-evidence`を一回行う。Assertionのkind/raw_text/observed_atとCandidate Evidenceが完全一致なら既適用であり、`action: create`、`execution_status: applied`、`failure_reason: null`とする。`rationale`には`already applied`を明記し、新しいmutationは呼ばない。操作順と結果はcreate conflict、get、get-evidenceとして`cli_operations`へ残す。
- Assertionが一意に一致するがEvidenceが未付与なら、`action: attach_evidence`、`target_assertion_id: <一致したID>`、`execution_status: applied`とする。`cli_operations`にはcreate conflict、get、get-evidence、attach-evidenceを実行順に残す。attachの失敗・中断は既存の`partially_applied`または`failed`とし、後続Candidateを`not_started`にする。
- 一意に特定できない、identityが不一致、またはCandidateに存在しないConcept/Alias/Identifier/Relationを安全に特定できない場合はattachへ変換しない。`action: create`、`execution_status: failed`、`failure_reason: cli_error`とし、確認に使った操作結果だけを残して停止する。後続Candidateは`not_started`である。

`already_applied` actionは追加しない。create Decisionの`target_assertion_id`と`replacement_assertion_id`は既存field契約どおりnullとし、確認したIDは`search_evidence`と`cli_operations.result_ids`、既適用の根拠は`rationale`で追跡する。
