---
name: knowledge-update
description: Candidate Knowledgeを既存Knowledge CLIと照合し、履歴を保つ更新判断とUpdate Resultを一時Markdownで返す。新しいCLI、DB、ledgerは作らない。
---

# Knowledge Update

Candidate Knowledgeを一件ずつ既存Knowledgeと照合し、`create`、`attach-evidence`、`revise`、`supersede`、`skip`のいずれかを選んで既存CLIを呼ぶ。Knowledge Storeの知識状態を推測で断定せず、操作結果と停止位置を[成果物契約](references/artifact-contract.md)どおりに返す。

## 開始時の参照

最初のCLI呼出し前に、同梱の[成果物契約](references/artifact-contract.md)、[CLI操作リファレンス](references/cli-operations.md)、[検証契約](references/verification.md)をすべて読む。これらと呼出側から渡された承認済み契約に不一致または欠落があれば、CLIを呼ばず `failed` または `not_started` の停止結果を返す。リポジトリ内の設計資料、ソース、別履歴、記事本文、Codex応答を実行時の入力として探索しない。

## 入力境界

入力は完成済みの一つのCandidate Knowledge成果物だけである。Candidateの必須field、`search_queries`、Evidence原文、`evidence_kind`、`strength`、`scope`、`temporal`、`source_ordinal`、`candidate_ordinal`をそのまま扱う。`strength`は`strong`、`moderate`、`weak`のいずれかで、Knowledge StoreやCLIへ保存するfieldではなく、CandidateのEvidenceに基づく判断用派生値である。Candidateは`source_ordinal`昇順、同じ発話内は`candidate_ordinal`昇順で処理し、`candidate_id`の文字列順に並べ替えない。

Candidateゼロは正常な空入力であり、CLIを呼ばず、Decision一覧が空で全体状態が`completed`のUpdate Resultを返す。入力不成立、未完成Episode、記事・AI・質問だけの内容をCandidateやDecisionへ変換しない。

## 手順

1. 同梱参照とCandidateの構造を検証する。未知のfield、必須fieldの欠落、順序違反、queryの空文字、Evidenceの不整合、`strength`の欠落・不正値を補完・修正せず入力不成立として停止する。`strength`は`evidence_kind`とEvidence原文から契約どおり導出された値であることを確認し、CLIへ渡さない。
2. 各Candidateの`search_queries`を先頭から順に、各一回だけ `knowledge search-text --query '<query>'` へ渡す。空の`data.results`でも次のqueryを続ける。検索結果の`assertion_id`はquery間・結果内を通じて初出順に重複除去する。
3. 各queryで得た一意Assertion IDについて、意味照合に必要なので`get`を初出順に一回呼ぶ。`get-evidence`は、そのAssertionが意味的候補（attach、revise、supersede、duplicate-skip、conflict確認）になった場合だけ一回呼ぶ。不一致と判断したIDには`get-evidence`を呼ばない。IDを生成したり、未発見のIDを推測したりしない。create候補に既存Assertionがない場合は`get-evidence`不要である。createがconflictした場合は、検索で得たIDだけを候補にして`get`結果のnormalized text、scope、temporal、identityを照合し、一意に特定できた候補だけへ`get-evidence`を一回呼ぶ。検索外のIDやmetadataを推測して確認してはならず、一意に確認できなければ失敗する。
4. Candidateごとに次の一つを選ぶ。既存Assertionと同じ命題で新Evidenceなら`attach_evidence`、同一identityの本文・Scope・Temporalの訂正なら`revise`、別identityの理解が旧Assertionを置き換えるなら`create`後に`supersede`、既存に対応せず新しい命題なら`create`、許可済みCandidateが不十分・根拠不十分・同じEvidence記録済みなら`skip`とする。意味的な同一性、Evidenceの強度、訂正、置換はCodexが判断する。
5. `create`にはCandidateから`normalized-text`、明示されたscope/temporal、Evidence groupだけを渡す。Candidateには構造化Concept、Alias、Identifier、Relation fieldがないため、`search_queries`、原文、scope、その他metadataからConcept/Alias/Identifier/Relation optionを推測して追加しない。これらの構造化fieldが将来追加された場合は、別の承認済み契約によるspec recheckが必要である。`attach-evidence`、`revise`、`supersede`にも承認済みのnamed optionだけを渡し、成功stdout、失敗stderr、exit codeを同じOperation Resultへ記録する。
5. `create`にはCandidateから`normalized-text`、明示されたscope/temporal、Evidence groupだけを渡す。Candidateには構造化Concept、Alias、Identifier、Relation fieldがないため、`search_queries`、原文、scope、その他metadataからConcept/Alias/Identifier/Relation optionを推測して追加しない。これらの構造化fieldが将来追加された場合は、別の承認済み契約によるspec recheckが必要である。`create`がconflictしたとき、既存AssertionとEvidenceのkind/raw_text/observed_atが完全一致すれば、新しいmutationを呼ばず、`action:create`、`execution_status:applied`、rationaleに`already applied`を明記して、create conflict/get/get-evidenceの結果を記録する。Assertionだけが一意に一致しEvidenceが未付与なら、そのIDへ`attach-evidence`を実行し、`action:attach_evidence`、`execution_status:applied`とする。候補が曖昧、identity不一致、またはConcept/Alias等を安全に特定できない場合はattachへ変換せず、`failure_reason:cli_error`の`failed`で停止する。`attach-evidence`、`revise`、`supersede`にも承認済みのnamed optionだけを渡し、成功stdout、失敗stderr、exit codeを同じOperation Resultへ記録する。
6. `revise`が成功した後に`attach-evidence`が失敗・中断した場合は、revisionを削除せず`partially_applied`とする。`create`が成功した後に`supersede`が未適用と分かる失敗・中断なら、新Assertionを削除せず`partially_applied`とする。`supersede`の`conflict`はRelationの適用を確認できないため、`failed`かつ`failure_reason: outcome_unknown`とする。補償削除、自動再試行、再開ledgerは作らない。
7. `skip`はCLI操作なしで記録する。成功またはskipのCandidate後は次へ進む。検索・取得・mutationが失敗、protocol error、exit 130、部分適用、結果不明になった時点で停止し、後続CandidateはCLIを呼ばず`not_started`で順序どおり記録する。
8. 全Candidateを含むUpdate Resultを一時Markdownで返す。失敗・中断・部分適用を成功扱いにせず、Knowledge Store、Candidate、元のAssessment、別ledgerへ保存しない。

## 失敗・中断の分類

- CLIの`validation_error`、`not_found`、`storage_error`、`internal_error`は、既存error JSONとexit codeを記録して停止する。状態や既適用を推測しない。
- stdout/stderrが契約と違う、成功JSONのdata形が違う、複数行JSONや成功・失敗の混在がある場合は`protocol_error`として停止する。
- response開始前の無出力exit 130は`canceled`。先行mutationがなければ全体状態も`canceled`、二段操作の後段なら`partially_applied`とし、いずれも後続を処理しない。
- `create`の`conflict`は、searchで得たIDに限定してgetでnormalized text/scope/temporal/identityを照合する。一意なAssertionだけへget-evidenceを行い、AssertionとEvidenceのkind/raw_text/observed_atが完全一致する場合は`action:create`の`applied`（rationaleに`already applied`）としてmutationなしで記録する。Assertionが一致してEvidenceが未付与ならattach-evidenceを行う。候補が曖昧、identity不一致、またはmetadataを安全に特定できない場合は`failure_reason:cli_error`の`failed`で停止し、attachしない。create conflict後のattachが失敗・中断した場合は、既存の二段操作規則に従い`partially_applied`または`failed`とし、後続を`not_started`にする。`attach-evidence`の重複conflictも完全一致を確認できる場合だけ既適用として扱う。確認できなければ技術失敗で停止する。`revise`の`conflict`は自動再実行しない。

## 非変更境界

このSkillはFEAT-001の既存CLIだけを呼ぶ。`search-semantic`、新CLI、公開JSON field、公開設定、SQLite schema、migration、workflow ledger、補償削除を追加・変更しない。CLI内部の一操作transactionを消費するが、二つのCLI操作を跨ぐ新しいtransactionは作らない。失敗時も過去Assertion、Evidence、revision、Relationを削除しない。

## 検証チェックリスト

- [ ] 同梱3参照を開始時に読んだ。
- [ ] Candidateの必須field、候補順、query列、Evidence原文を追跡できる。
- [ ] `strength`の存在・enum・`evidence_kind`からの導出を検証し、CLI保存fieldとして扱っていない。
- [ ] queryを一回ずつ順に実行し、空結果後も続行し、Assertion IDを初出順で重複除去した。
- [ ] 各一意Assertion IDへ`get`を一回、意味的候補だけへ`get-evidence`を一回行い、不一致IDをhydrationしていない。
- [ ] createでConcept/Alias/Identifier/Relationをmetadataから推測していない。
- [ ] create conflictではsearch由来のIDだけをgetし、normalized text/scope/temporal/identityを比較した。
- [ ] create conflictの完全一致を`action:create`の`already applied`として記録し、Evidence未付与だけを`attach_evidence`へ進め、曖昧・不一致・metadata特定不能は`cli_error`で停止した。
- [ ] create / attach-evidence / revise / supersede / skipの根拠、option、結果、IDを記録した。
- [ ] conflict、protocol error、既存error code、exit 130を契約どおり分類した。
- [ ] 二段操作の後段失敗・中断・結果不明で履歴削除や自動再試行を行わず、後続を`not_started`にした。
- [ ] Candidateゼロを空Decision一覧・全体`completed`として返した。
- [ ] Update Resultを保存せず、元のAssessment本文を変更しなかった。
