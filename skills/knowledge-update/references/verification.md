# Knowledge Update 検証契約

fixtureごとに隔離したKnowledge Storeを用意し、実際の`knowledge` binaryをプロセス境界で起動する。各caseでargv、stdout、stderr、exit code、実行順、Update Result、DBの履歴を観測する。成功JSONとerror JSONを推測で合成しない。

## 共通oracle

- Candidate処理順は`source_ordinal`→`candidate_ordinal`。
- `search_queries`は列の順に一回ずつ実行し、空結果後も次を呼ぶ。
- Assertion IDは全検索結果の初出順で重複除去し、その集合だけをhydrationする。
- 各一意Assertion IDへ`get`を初出順に一回行い、意味的候補だけへ`get-evidence`を一回行う。不一致IDへのget-evidence、search結果にないIDの推測、metadataからのConcept/Alias/Identifier/Relation推測を許可しない。
- 失敗・protocol error・exit 130・部分適用後は後続CandidateのCLI呼出しがなく、Update Resultに後続が`not_started`で残る。
- 成功時はstdout一つ、失敗時はstderr一つ、exit codeと既存error codeが契約に一致する。
- DBでは旧Assertion、旧Evidence、旧revision、既存Relationが削除されず、新規workflow ledger・migration・schema差分がない。

## 再現ケース

| ID | Fixture / process観測 | 期待oracle |
| --- | --- | --- |
| VU-001 | Candidate zeroを入力 | CLI呼出しなし、Decision一覧空、全体`completed`。 |
| VU-002 | query列を`assertion`、空結果、明示Concept、同じIDを返すqueryの順で入力 | search-textは各一回、空後も継続、IDは初出順一度だけ、hydrationは必要IDだけ。 |
| VU-003 | 複数Candidateをsource/candidate ordinalが異なる文字列IDで入力 | ordinal順に処理し、ID文字列順ではない。 |
| VU-004 | 既存対応なしのCandidateを隔離DBで実CLI process起動 | argvが`normalized-text`、scope/temporal、Evidence groupだけで、Concept/Alias/Identifier/Relationを推測していないこと、stdoutの成功JSON、stderr空、exit 0、`assertion_id`・`revision`・`evidence_ids`・`relation_ids`・`concepts`、新Assertion/初期revision/EvidenceがDBに残ることを記録。実binary未実行時は`not_run`。 |
| VU-005 | 同一Assertionと未記録Evidenceを隔離DBで実CLI process起動 | search後に対象IDへget一回、意味的候補へget-evidence一回、argvのkind・原文・観測時刻、stdout成功JSON、stderr空、exit 0、返却Assertion/Evidence ID、既存Assertionと新EvidenceがDBに残ることを記録。 |
| VU-006 | 本文・Scope・Temporalの訂正を隔離DBで実CLI process起動 | search→get一回→意味的候補ならget-evidence一回→reviseの必須option、stdout成功JSON、stderr空、exit 0、`assertion_id`・`previous_revision`・`revision`、旧revision/Evidenceと新revisionがDBに残ることを記録。訂正Evidenceが必要ならattachを二段目にする。 |
| VU-007 | 別identityの置換を隔離DBで実CLI process起動 | create argvにmetadataを推測したoptionがなく、stdout成功JSON/exit 0の`assertion_id`をsupersede後段argvへ渡す。両operationのstdout、stderr、exit code、relation ID、旧Assertion・新Assertion・旧revision/Evidence・supersedes RelationがDBに残ることを記録。 |
| VU-008 | 許可済みだが不十分、または同一Evidence済み | CLI呼出しなし、`skip`/`skipped`、操作一覧空。除外入力にはDecisionなし。 |
| VU-008a | strength欠落、不正enum、またはevidence_kindと導出不一致 | CLI呼出しなし、入力不成立で停止。`strength`はCLI保存fieldではなく、導出表どおりの派生値だけを受理。 |
| VU-009 | search-textのvalidation/not-found/storage/internal error、JSON形不整合 | 最初の失敗を記録して停止。残りquery・CandidateにCLIなし。`protocol_error`はerror_code null。 |
| VU-010 | attach/createの重複Evidence conflict | 完全一致をget-evidenceで確認できる場合だけ既適用。それ以外はfailed/cli_error、後続はnot_started。 |
| VU-010a | create conflict後、search由来IDのgetでAssertion（normalized text/scope/temporal/identity）を一意に特定し、get-evidenceでkind/raw_text/observed_atが完全一致 | argvのcreate必須option、stderrのconflict JSON、exit 4、続くget/get-evidenceのargv、stdout成功JSON、stderr空、exit 0、result IDsと実行順を観測する。新規mutationなし、Decisionは`action:create`、`execution_status:applied`、rationaleに`already applied`、新Evidence/Assertion/履歴差分なし。後続Candidateは通常どおり処理可能。 |
| VU-010b | create conflict後、getでAssertionを一意に特定するがEvidenceが未付与 | create conflict→get→get-evidence→attach-evidenceのargv、stdout/stderr、exit code、必須option、Assertion/Evidence result IDsを観測する。Decisionは`action:attach_evidence`、`execution_status:applied`。既存Assertion、旧revision、既存Evidenceを保持し、新EvidenceだけがDBに残る。 |
| VU-010c | create conflict後、候補が複数、identity不一致、またはConcept/Alias等を安全に特定不能 | searchで得たIDだけのget argvと各stdout/error、create conflict stderr/exit 4を観測し、unique match以外へget-evidenceを呼ばず、attachも呼ばない。Decisionは`action:create`、`execution_status:failed`、`failure_reason:cli_error`、後続CandidateはCLI呼出しなしで`not_started`。既存Assertion/Evidence/revision/Relationを変更・削除しない。 |
| VU-011 | revise成功後のattach失敗 | revisionをDBから削除せず、`partially_applied`、成功revision ID、失敗attach operation、後続not_started。 |
| VU-012 | revise成功後のattachがresponse前exit 130 | exit 130、error JSONなし、`partially_applied`、後続not_started。 |
| VU-013 | create成功後のsupersedeが未適用と分かる失敗またはexit 130 | 新Assertionを削除せず`partially_applied`、成功IDと後段operationを記録。 |
| VU-014 | create成功後のsupersede conflict | Relationを読取確認せず、`failed`、`failure_reason: outcome_unknown`、後続not_started。 |
| VU-015 | 全operationの成功・skipを混在 | 各Candidateのargv、stdout成功JSON、stderr、exit code、必須option、result IDs、create/attach/revise/supersede後の旧Assertion・旧revision・Evidence・Relationを含むDB履歴が順序どおりで、全体`completed`。実binary未実行時は`not_run`。 |
| VU-016 | worktree/DB/公開CLIの前後比較 | workflowの新規保存・ledger・migration・schema・公開operation・JSON追加なし。既存CLI一操作のtransactionだけが消費される。 |

## AC対応

| AC | 主な証拠 |
| --- | --- |
| 検索列・空結果・hydration | VU-001〜VU-003のprocess log、VU-002のID列。 |
| operation判断と追跡 | VU-004〜VU-008、VU-015のDecisionとOperation Result。 |
| 停止・失敗・中断・部分適用 | VU-009〜VU-014、VU-010a〜VU-010cのexit、stderr、Update Result、後続呼出し不在。 |
| 履歴保全と非変更境界 | VU-011〜VU-014、VU-016のDB/schema/diff比較。 |
| Candidate順序・空候補 | VU-001、VU-003、VU-015。 |

未実行caseは成功扱いにせず、`not_run`としてImplementation Reportへ記録する。今回のTaskはworkflow docsのみで実binary fixtureを変更・追加しないため、実CLI process観測は`not_applicable`（oracle定義済み）とする。Go/SQLite実装変更も対象外であり、Go gateは`not_applicable`とする。実行可能なMarkdown構造検査とscope検査は別途実行する。
