# FEAT-004 ワークフロー成果物契約

> Codex workflow間のMarkdown成果物に適用する論理契約。Knowledge CLIのJSON契約を変更しない。

## URL Evaluation Episode

| Field | 型・必須性 | 意味 |
| --- | --- | --- |
| `episode_id` | opaque string、必須 | 一回のURL評価実行を表す安定ID。当該実行内の候補・判断・結果を対応付ける。中断後の再送・再開には使わない。 |
| `article_url` | URL string、必須 | 評価を依頼された技術記事URL。 |
| `completed_at` | RFC 3339 UTC string、必須 | Reading Value Assessmentの内容を完成させた時刻。候補抽出を開始できる時点を表す。会話への返却時刻ではない。 |
| `user_contributions` | list、必須。空可 | 当該URL評価内でユーザーが発した内容。下表の要素を発話順で持つ。Codex応答と記事本文は含めない。 |

`completed_at`がない、またはReading Value Assessmentの内容が完成していないEpisodeはKnowledge Acquisitionへ渡さない。

`user_contributions` の各要素は次のfieldを持つ。

| Field | 型・必須性 | 意味 |
| --- | --- | --- |
| `ordinal` | 1以上のinteger、必須 | Episode内のユーザー発話順。重複不可。 |
| `source_text` | string、必須 | ユーザー発話またはユーザーが提示したコードの原文。空白だけは不可。 |
| `source_type` | `message` / `code`、必須 | 自然言語のユーザー発話か、ユーザーが提示したコードかを示す。 |
| `observed_at` | RFC 3339 UTC string、必須 | 当該ユーザー寄与を観測した時刻。 |

## Candidate Knowledge

| Field | 型・必須性 | 意味 |
| --- | --- | --- |
| `episode_id` | opaque string、必須 | 入力Episodeへの参照。 |
| `candidate_id` | opaque string、必須 | 同じEpisode内で一意の候補ID。 |
| `source_ordinal` | integer、必須 | 根拠にした`user_contributions.ordinal`。複数の根拠を一候補へ混在させない。 |
| `candidate_ordinal` | 1以上のinteger、必須 | 同じ`source_ordinal`から分けた候補の出現順。発話内の先頭から独立命題を抽出した順に連番とし、重複不可。 |
| `source_excerpt` | string、必須 | ユーザー由来Evidenceの必要最小限の原文。空白だけは不可。 |
| `evidence_raw_text` | string、必須 | `source_ordinal`の原文のうち、候補を直接支える完全なEvidence原文。省略・要約・意味を変える編集は不可。CLIのEvidence `raw_text`にはこの値を渡す。 |
| `evidence_kind` | enum、必須 | `user_explanation`、`user_reasoning`、`user_code`、`self_report`、`concept_recognition`、`correction`のいずれか。理由を伴う技術判断は`user_reasoning`に分類する。 |
| `strength` | enum、必須 | `strong`、`moderate`、`weak`。保存fieldではなく、`evidence_kind`と`evidence_raw_text`を根拠にCodexがこの候補時点で再評価した派生表示。 |
| `observed_at` | RFC 3339 UTC string、必須 | 参照する`user_contributions.observed_at`と同じ値。 |
| `proposed_assertion` | string、必須 | ユーザー知識として独立評価可能な正規化候補。 |
| `search_queries` | 1件以上のstring list、必須 | 既存`search-text --query`へ順番に渡す検索入力。先頭は`proposed_assertion`そのもの。訂正Evidenceに引用符で囲まれた旧命題がある場合だけ、その完全な引用文を二番目に置ける。残りは原文で明示されたConcept・Alias・Identifierだけを出現順に置く。同じ文字列は最初の一件だけにし、Scopeや推測した同義語は加えない。 |
| `scope` | key/value list、必須。空可 | 各要素は空でない`key`と`value`を各1つ持つ。発言で明示された適用範囲だけであり、同じkeyを重複させず、推測で補わない。 |
| `temporal` | objectまたはnull、必須 | 発言で明示された`valid_from`、`valid_until`、`version_scope`、`observed_at`、`last_verified`だけを持つ。時刻はRFC 3339 UTC、未指定fieldは`null`、情報がなければobject全体が`null`。 |
| `extraction_rationale` | string、必須 | 候補化したユーザー由来の理由。 |

Candidateは一つの`evidence_kind`だけを持つ。一つの原文が複数の独立命題または複数種類の根拠を直接支える場合は、命題・根拠種類ごとに候補を分ける。`candidate_id`は追跡だけに使うopaque IDであり、処理順には使わない。原文を切り詰めて意味を変えてはならない。Knowledge Updateは各`search_queries`を先頭から一回ずつ`search-text --query`へ渡し、空結果でも次のqueryへ進む。各queryで返ったAssertion IDは最初に現れた順で重複を除いて集め、必要時の`get`/`get-evidence`はこの集合だけに行う。`search-text`の技術失敗または中断では残りのquery・Candidateを実行しない。

### Evidence強度の導出

| `evidence_kind` | `strength` | 導出規則 |
| --- | --- | --- |
| `user_explanation` / `user_reasoning` / `user_code` | `strong` | ユーザー自身が正しい技術的説明、理由付け、またはコードを直接示す。 |
| `correction` | `strong` | 以前の理解・Assertionを明示的に訂正する。 |
| `self_report` | `moderate` | 経験・利用・知識を明示的に自己申告するが、説明・推論・コードで直接は裏付けない。 |
| `concept_recognition` | `weak` | 用語や概念を知っていることだけを示し、説明・推論・コード・自己申告を含まない。 |

理由を伴う技術判断は`user_reasoning`であり`strong`とする。理由を伴わない結論・方針表明だけの技術判断は、正しい説明・推論・コード・明示的自己申告・概念認識のいずれにも当たらないためCandidateを作らない。同一寄与に複数の根拠がある場合も、一候補へ混在させず上表に従う候補へ分ける。質問だけ、または表のどれにも該当しない内容はCandidateを作らない。

## Update Decision

| Field | 型・必須性 | 意味 |
| --- | --- | --- |
| `episode_id` / `candidate_id` | opaque string、必須 | 入力候補への追跡参照。 |
| `action` | enum、必須 | `create`、`attach_evidence`、`revise`、`supersede`、`skip`、`not_started`。`not_started`は選択前の技術失敗・中断、または前のCandidate停止により未処理となったCandidateだけに使う。 |
| `target_assertion_id` | opaque stringまたはnull、必須 | `attach_evidence`・`revise`では対象Assertion。`supersede`では旧Assertion。`create`・`skip`・`not_started`は`null`。 |
| `replacement_assertion_id` | opaque stringまたはnull、必須 | `supersede`の新Assertion。`create`成功後に得る。他actionは`null`。 |
| `rationale` | string、必須 | 意味的同一性、訂正、置換、または許可済みCandidateをskipする根拠。入力境界から除外された発言にはDecisionを作らない。 |
| `search_evidence` | Assertion/Evidence ID list、必須。空可 | 判断に使った既存Knowledgeへの参照。 |
| `execution_status` | enum、必須 | `not_started`、`applied`、`skipped`、`partially_applied`、`failed`、`canceled`。`failed`は成功・不成功を判定できない結果不明を含む。 |
| `failure_reason` | enumまたはnull、必須 | `not_started`・`applied`・`skipped`・`canceled`・`partially_applied`では`null`。`failed`時は`cli_error`、`protocol_error`、`outcome_unknown`のいずれか。 |
| `cli_operations` | ordered list、必須。空可 | 下表のoperation resultを実行順に持つ。 |

`skip`は`execution_status: skipped`かつ`cli_operations: []`である。`skip`は許可済みCandidateが命題として不十分、知識根拠にならない、または同じEvidenceが既に記録済みの場合だけに使い、質問・AI・記事本文など入力境界から除外された発言には作らない。選択前の検索・取得が失敗または中断したCandidateは`action: not_started`、`execution_status: failed`または`canceled`とし、rationaleへ停止理由、`cli_operations`へ失敗operationを残す。前のCandidate停止で処理しなかったCandidateは`action: not_started`、`execution_status: not_started`、`target_assertion_id`・`replacement_assertion_id`・`failure_reason`は`null`、`search_evidence`・`cli_operations`は空、rationaleは先行停止を示す。`revise`だけ、または`supersede`の`create`だけが成功し、後段が未適用と分かる失敗・中断なら、`execution_status: partially_applied`で成功IDと失敗した後段operationを残す。`supersede`の`conflict`はRelationを読む既存operationがないため、後段が未適用とは断定できない例外として`execution_status: failed`、`failure_reason: outcome_unknown`とする。

`cli_operations` の各要素は、`operation`（既存operation名）、`query`（`search-text`ではその回の`search_queries`要素、他operationでは`null`）、`status`（`applied`、`failed`、`canceled`）、`result_ids`（成功時に返ったAssertion / Evidence / Relation IDのlist。失敗・中断時は空）、`error_code`（既存CLI error codeまたは`null`）、`exit_code`（integer）を持つ。`failed`で`failure_reason: cli_error`なら`error_code`が必須、`failure_reason: protocol_error`なら`error_code: null`を許可する。`canceled`では`error_code: null`かつ`exit_code: 130`である。

## Update Result

一つのEpisodeの処理結果であり、`episode_id`、`completed_at`、**全Candidateを順番どおり含む**Update Decision一覧、全体状態を含む。候補がゼロの場合は、Decision一覧を空、全体状態を`completed`とし、Candidateに紐付かない`skip` Decisionは作らない。候補がある場合の全体状態は、`partially_applied`があれば`partially_applied`、先行mutationなしのexit 130なら`canceled`、技術失敗または結果不明なら`failed`、すべて`applied`または`skipped`なら`completed`とする。最初の`failed`、`canceled`、または`partially_applied`のDecision以降のCandidateは上記`not_started`で記録し、CLIを呼ばない。候補は`source_ordinal`昇順、同じ発話内では`candidate_ordinal`昇順で処理する。

Update Resultは呼出側へ返す一時的なMarkdown成果物であり、Knowledge Storeまたは別ledgerへ保存しない。中断・応答喪失後に同じ`episode_id`を再開・自動再実行しない。

URL評価への回答本文・記事本文の複製、知識状態ラベル、読書推奨は含めない。
