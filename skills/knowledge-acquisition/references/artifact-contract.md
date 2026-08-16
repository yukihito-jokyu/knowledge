# Candidate Knowledge 成果物契約

この契約は、承認済みのEpisode/Candidate論理契約を、Codex workflow間で受け渡すための一時Markdown成果物として具体化する。Knowledge CLIの公開JSON、保存先、DB、migration、Update Decisionの契約ではない。

このSkillは、実行場所や`workflow-contracts.md`の物理配置を前提にしない。呼出側は正規の要件・handoff・workflow契約を必ず渡し、Skillはそれを根拠として同梱契約と照合する。正規契約が欠落している場合、契約が不一致の場合、または実行時にリポジトリ内の特定パスを探索・読込する必要がある場合は、Candidateを生成せず停止する。

## 入力: URL Evaluation Episode

一回のURL評価実行を表す一つのEpisodeだけを受け取る。Assessment本文の内容が完成した後、会話へ返す前に渡される。`completed_at`がない、完成保証がない、または評価が失敗・中断した場合は入力不成立である。

| Field | 必須・型 | 契約 |
| --- | --- | --- |
| `episode_id` | 必須、opaque string | Episode内の候補を対応付ける不透明ID。意味や順序を埋め込まない。 |
| `article_url` | 必須、URL string | 評価対象URL。記事本文をCandidateへ複製するためには使わない。 |
| `completed_at` | 必須、RFC 3339 UTC string | Assessment本文の内容を完成させた時刻。回答返却時刻ではない。 |
| `user_contributions` | 必須、list（空可） | 当該URL評価内のユーザー寄与だけを発話順に含む。Codex応答、記事本文、閲覧・評価・要約結果を含めない。 |

`user_contributions`の各要素は次を持つ。

| Field | 必須・型 | 契約 |
| --- | --- | --- |
| `ordinal` | 必須、1以上のinteger | Episode内のユーザー寄与順。重複不可。入力リスト上も昇順であること。 |
| `source_text` | 必須、空白だけでないstring | ユーザー発話またはユーザー提示コードの原文。省略・要約・意味変更をしない。 |
| `source_type` | 必須、`message` / `code` | 自然言語のユーザー発話か、ユーザー提示コードか。 |
| `observed_at` | 必須、RFC 3339 UTC string | 当該寄与を観測した時刻。Candidateにも同じ値を渡す。 |

## 入力検証と停止

次のどれかに該当する場合、`invalid_input`の停止結果を返し、Candidate一覧とUpdate Decisionを作らない。

- 必須fieldの欠落、空の`episode_id`、空の`article_url`、空の`source_text`。
- `completed_at`または`observed_at`がRFC 3339 UTCでない。
- `ordinal`が1未満、重複、または入力リスト上で降順・同順位になっている。
- `source_type`が`message` / `code`以外。

Assessment本文が未完成、URL評価が失敗・中断、または完了時刻を確定できない場合は`assessment_incomplete`として停止する。完成済みEpisodeから許可されるユーザーEvidenceが一件も抽出されない場合だけは、正常なCandidateゼロとする。

## 出力: Candidate Knowledge

正常時のMarkdown成果物は次の外枠を持つ。`article_url`と`completed_at`はEpisodeの追跡メタデータであり、Candidate各件には未定義の追加fieldを持ち込まない。

```markdown
# Candidate Knowledge

## Episode

- episode_id: <opaque id>
- article_url: <url>
- completed_at: <RFC 3339 UTC>

## Candidates

### Candidate 1

- episode_id: <opaque id>
- candidate_id: <opaque id>
- source_ordinal: <integer>
- candidate_ordinal: <integer>
- source_excerpt: <必要最小限のユーザー原文>
- evidence_raw_text: <候補を直接支える完全なユーザーEvidence原文>
- evidence_kind: <enum>
- strength: <enum>
- observed_at: <RFC 3339 UTC>
- proposed_assertion: <独立評価可能な候補Assertion>
- search_queries:
  - <proposed_assertionそのもの>
  - <訂正Evidenceで引用された旧命題（ある場合だけ）>
  - <原文に明示されたConcept/Alias/Identifier>
- scope: <明示されたkey/value list、なければ []>
- temporal: <明示情報object、なければ null>
- extraction_rationale: <候補化したユーザー由来の理由>

## Handoff

- destination: knowledge-update
- persistence: none
- update_decision: none
```

候補がない場合は`## Candidates`の直下を`なし`とし、`Handoff`は空一覧を渡すことと`update_decision: none`を示す。停止結果では`Candidate Knowledge`の候補外枠を作らず、停止理由だけを返す。

### Candidate field契約

| Field | 必須・型 | 契約 |
| --- | --- | --- |
| `episode_id` | 必須、opaque string | 入力Episodeと同じ値。 |
| `candidate_id` | 必須、opaque string | Episode内で一意な追跡ID。意味・順序・Assertion IDを表さない。 |
| `source_ordinal` | 必須、integer | 根拠にした寄与の`ordinal`。一候補に複数寄与を混在させない。 |
| `candidate_ordinal` | 必須、1以上のinteger | 同じ`source_ordinal`内で独立命題を抽出した順。重複不可。 |
| `source_excerpt` | 必須、string | 表示・確認に必要な最小限のユーザー原文。空白だけ、要約、推測文は不可。 |
| `evidence_raw_text` | 必須、string | 候補を直接支える完全なEvidence原文。省略記号、要約、意味を変える編集、記事・AI本文の混入は不可。 |
| `evidence_kind` | 必須、enum | `user_explanation`、`user_reasoning`、`user_code`、`self_report`、`concept_recognition`、`correction`のいずれか。 |
| `strength` | 必須、enum | `strong`、`moderate`、`weak`。Evidenceから導出する表示値で、Knowledge Storeの保存fieldではない。 |
| `observed_at` | 必須、RFC 3339 UTC string | 参照寄与の`observed_at`と同じ値。 |
| `proposed_assertion` | 必須、string | ユーザー知識として独立評価できる候補。Evidenceの意味を越えて一般化しない。 |
| `search_queries` | 必須、1件以上のstring list | 先頭は`proposed_assertion`そのもの。訂正Evidenceに引用された旧命題がある場合だけ二番目にその完全な引用文を置ける。残りは原文に明示されたConcept・Alias・Identifierだけ。 |
| `scope` | 必須、key/value list（空可） | 発話で明示された適用範囲だけ。空でないkey/value、同一key重複なし、推測補完なし。 |
| `temporal` | 必須、objectまたはnull | 明示された`valid_from`、`valid_until`、`version_scope`、`observed_at`、`last_verified`だけ。情報がなければnull。 |
| `extraction_rationale` | 必須、string | なぜ許可されたユーザー由来Evidenceとして候補化したか。AIや記事の説明を理由にしない。 |

Candidateは一つの`evidence_kind`だけを持つ。複合寄与は同じ`source_ordinal`から複数Candidateへ分け、各`evidence_raw_text`がその候補の完全な直接根拠になるようにする。並びは`source_ordinal`昇順、同一寄与では`candidate_ordinal`昇順とし、`candidate_id`の文字列順で並べ替えない。

## Evidence分類と強度

| 入力の観測 | `evidence_kind` | `strength` |
| --- | --- | --- |
| ユーザー自身の正しい技術的説明 | `user_explanation` | `strong` |
| 理由・因果・比較を伴う技術的推論、理由付き技術判断 | `user_reasoning` | `strong` |
| ユーザーが提示した技術コード | `user_code` | `strong` |
| 以前の理解・Assertionを明示的に訂正 | `correction` | `strong` |
| 説明・推論・コードを伴わない経験・利用・知識の自己申告 | `self_report` | `moderate` |
| 用語・概念を知っていることだけを示す言及 | `concept_recognition` | `weak` |

理由のない結論・方針表明はCandidateにしない。単に技術語を含む質問も除外する。複数種類が同じ寄与にある場合は一件へ混在させず分割する。

## 除外・プライバシー境界

通常会話、別作業、保存済み履歴、失敗・中断した評価、Codex応答、記事本文、記事の閲覧・評価・要約、質問だけの発話、理由のない技術判断はCandidateにしない。URLはEpisodeを追跡するメタデータとしてのみ扱い、記事本文を引用・保存しない。Codexが生成した説明をユーザーEvidenceとして再利用しない。

除外入力にはCandidateもUpdate Decisionも作らない。Acquisitionは除外されたことを示すための`skip` Decisionすら作らない。`skip`を含むUpdate Decisionは後続Knowledge Updateの責務であり、本成果物には含めない。

## search_queries規則

1. 先頭要素は`proposed_assertion`と完全一致させる。
2. 訂正Evidenceに引用符で囲まれた旧命題がある場合だけ、その完全な引用文を二番目の要素へ置く。引用がない、または引用文が先頭要素と完全一致する場合は置かない。
3. Evidence原文を左から右へ走査し、ユーザーがConcept、Alias、Identifierとして明示した文字列だけを残りの後続要素へ置く。コード識別子、明示されたAPI・パッケージ・コマンド名、引用・ラベル付きの用語は対象になり得る。
4. 文字列は原文の綴り、大小文字、記号、空白を保つ。完全一致する文字列が再出現した場合は最初の一件だけ残す。
5. 推測した同義語、翻訳語、関連語、Scopeの値、一般化した表現を追加しない。Concept・Alias・Identifierか判断できない語は追加しない。

## Handoff境界

Candidate Knowledgeは呼出側へ返す一時的なMarkdown成果物であり、Knowledge Updateが後続で検索・判断に使う。Acquisitionは保存、CLI呼出し、Update Decision、Update Result、Knowledge stateの断定をしない。ファイル、DB、ledger、設定への永続化を行わず、応答喪失後に同じ`episode_id`を再開・自動再実行しない。
