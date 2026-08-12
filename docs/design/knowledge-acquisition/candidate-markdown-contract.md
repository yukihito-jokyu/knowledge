# Candidate Markdown・候補なし・入力不足の引渡し契約

## 1. 目的と担当境界

この文書は、Knowledge Acquisition（利用者本人が会話や作業で示した内容から、知識を判断する候補を取り出す処理）が、Knowledge Candidate（既存の知識記録との比較・更新判断の前段階にある知識候補）を後続へ渡すMarkdown（見出しと本文から成る受渡し文書）の形式を定める。

同じ受渡し文書で、候補が一件以上ある結果、十分に評価したが候補がない正常結果、入力が不足して安全に評価できない結果を区別する。区別が必要なのは、候補なしを利用者本人の不知や障害と誤認せず、入力不足を通常終了として失わないためである。

この文書の単一Owner（変更責任を持つTask）は`L3-M1-S3`である。`L3-M1-S1`の入力・Evidence候補採否境界、`L3-M1-S2`のCandidate正規化とL1 Field写像は読取り入力として使い、それらの意味を変更しない。既存Knowledgeとの比較、保存・更新操作の選択は`L3-M2`、起動、再試行、停止、Knowledge Updateを起動しない制御は`L5`、実行用Fixture・Dataset・評価Harnessは`L6`の担当であり、本書では決めない。

## 2. 成果物の共通規則

### 2.1 文字・見出し・順序

成果物はUTF-8のMarkdown一文書とし、機械識別に必要な値はコード表記で記す。見出しの順序、項目名、結果種別の識別値は、この契約に定めたものから変更・省略・翻訳してはならない。後続が結果種別とCandidate境界を安定して読み取り、版ごとの推測をしないためである。説明本文の自然言語は日本語以外でもよいが、原文、命題、出典の内容を翻訳・要約して置き換えてはならない。

候補ありの文書では第3章の共通Sectionを一度だけ置き、その後に第4章のCandidate Sectionを一件ずつ出現順で置く。候補なしと入力不足の文書ではCandidate Sectionを置かず、第5章または第6章の結果Sectionだけを置く。空のCandidate Section、未使用項目の空文字列、`TBD`、推測値を出力してはならない。存在しない値を空欄と見せると、後続が未観測・不要・取得失敗を区別できないためである。

### 2.2 共通Section

すべての結果は、先頭から次のSectionをこの順に持つ。

```markdown
# Knowledge Acquisition Result

## Contract
- contract_name: `knowledge-acquisition-candidate-markdown`
- contract_version: `1.0`

## Episode
- episode_id: `<評価済みEpisodeの識別子>`
- ended_at: `<評価範囲を確定した終了時点>`

## Result
- result_kind: `<candidate | no-candidate | input_insufficient>`
```

`contract_name`は、この文書が何の受渡し形式かを示す固定名であり、別のMarkdownをCandidate結果と取り違えないために必須である。`contract_version`は構造互換性を判定する版であり、内容の知識状態やEpisodeの版ではない。`1.0`以外を受信した側は、同じ意味と構造を推測せず互換性確認へ渡す。

`episode_id`は、候補、診断、後続の更新判断を同じ完了済みEpisodeへ結び直す識別子であり、空値や別Episodeの値を禁止する。`ended_at`は、訂正を含む評価範囲が閉じた時点であり、空値、推定時刻、開始時点より前の時刻を禁止する。`candidate`と`no-candidate`では両方とも実値を必須とする。一方、`input_insufficient`で当該値そのものが取得不能な場合だけ、該当する値を固定リテラル`unavailable`（値を推測・補完せず取得不能と明示する印）にできる。この印は実際のEpisode識別子・時刻ではなく、対応する`diagnostics`に`episode_identity_missing`又は`episode_end_missing`を必ず置く。これにより、必須入力の欠落を空欄や偽の値で隠さず、候補あり又は正常な候補なしと誤認しない。`result_kind`は評価可否と候補有無だけを表す分類であり、`known`や`reported_unknown`のような利用者本人の知識状態、保存・更新の成否、再試行可否を表してはならない。

## 3. `candidate` の全体構造

`result_kind`が`candidate`である結果は、必須入力が揃った完了済みEpisodeを評価し、一件以上のCandidateを作れた場合だけ使う。Candidateの件数は一件以上とし、件数ゼロ、入力不備、評価途中には使わない。Candidateの存在は保存済みKnowledgeの作成、利用者本人の`known`、既存Knowledgeとの同一性を意味しない。

各Candidateは次の見出しから開始する。`<n>`は1から連続する10進整数であり、Candidateの並び順を固定するために必要である。同じEpisode内では、Evidence候補が現れた順を基本とし、同一の根拠から複数Candidateを分けたときは原文中の出現順にする。意味的な重要度、保存優先度、更新操作の順で並べ替えてはならない。これらは後続の`L3-M2`が判断する事項である。

```markdown
## Candidate <n>

### Assertion
- statement: `<独立評価可能な正規化済み一命題>`
- language: `<statementの自然言語>`

### Concept <k>
- canonical_label: `<検索用の代表表記>`
- normalized_label: `<検索・比較用に表記揺れを抑えた表記>`
- aliases: [`<同じConceptを指す別表記>`, ...]

### Scope
- facets:
  - kind: `<domain | language | subject>`
    original_label: `<原表記>`
    normalized_label: `<正規化表記>`

### Temporal
- observed_at: `<原文を観測した時点>`
- valid_from: `<原文で確認できる場合だけ>`
- valid_until: `<原文で確認できる場合だけ>`
- version_scope: `<原文で確認できる場合だけ>`

### Evidence <m>
- evidence_kind: `<explanation | reasoning | code | correction | knowledge_self_report | technical_decision>`
- stance: `<supports | challenges | qualifies>`
- knowledge_report: `<reports_known | reports_unknown | recognizes | reports_experience>`
- strength: `<strong | moderate | weak>`
- raw_content: |
    <利用者本人の加工前の原文>
- content_language: `<raw_contentの自然言語>`
- source_kind: `<message | file>`
- source_locator: `<元の会話又はファイルへ再到達できる参照先>`
- fragment_locator: `<出典内の具体的位置。条件付き>`
- originator:
  - actor_kind: `user`
- episode_id: `<共通Sectionと同じ識別子>`
- captured_at: `<出典を取得した時点>`
- content_digest: `<内容照合値。条件付き>`
- observed_at: `<この根拠を観測した時点>`

### Normalization rationale
<原文からこの一命題、Concept、Scope、時点、根拠を分離・統合した理由と、未観測として残した事項>
```

### 3.1 Candidateごとの必須性と禁止事項

`Assertion.statement`は、利用者本人が知っているかを単独で評価できる平叙文一件である。複数命題の束、能力ラベル、Evidence原文、Codexの説明を禁止する。`language`は命題の自然言語であり、プログラミング言語を入れない。原文と正規化済み命題を分離しないと、後続がどの内容を評価・更新するかを再確認できない。

`Concept <k>`は一件以上置き、`<k>`は1から連続する10進整数とする。複数のConceptは、命題本文に現れる順、同順位なら`canonical_label`のUnicodeコードポイント順で並べる。各Conceptは`canonical_label`と`normalized_label`を必須とし、`aliases`は同じ概念を指すと説明できる別表記がある場合だけ置く。別名がなければ`aliases: []`と明示する。Conceptへ命題全文、知識状態、根拠の強さ、Codexの推測を入れてはならない。Conceptは「利用者本人がその語を知っている」という主張ではなく、命題を検索・比較する入口だからである。

`Scope.facets`は原文から観測できた場合だけ置き、空の`facets`配列は置かない。`kind`の`domain`は技術・業務分野、`language`はプログラミング言語、`subject`は製品・ライブラリ・機能を表す。同じ`kind`の複数値はそのいずれか、異なる`kind`の値はすべて満たす条件として扱う。未記載は全範囲で有効という意味ではなく、条件を観測していないことを表す。似た名称の自然言語とプログラミング言語を混同して検索・比較することを防ぐため、分類を分ける。

`Temporal.observed_at`は必須であり、Candidateを作成した時点、Markdownを出力した時点、任意に選んだEvidenceの時点を入れてはならない。Candidateが持つ全`Evidence <m>.observed_at`のうち最も早い時点を、Candidateの`Temporal.observed_at`へ必ず写す。Candidateが一件のEvidenceだけを持つ場合も、二つの値は同一になる。各Evidenceの`observed_at`は、その原文断片を観測した時点であり、Candidateの最初の観測時点より前であってはならない。これにより、複数の原文を一つのCandidateへまとめても、その命題を支える観測が始まった時点を一意に追跡でき、後から追加した根拠や文書作成時刻を知識内容の観測時点と取り違えない。`valid_from`、`valid_until`、`version_scope`は原文で確認できる場合だけ置き、未観測の値を現在有効、無期限、特定版と推測しない。時点と版を分けるのは、古い内容や対象版の違いを現在の矛盾と誤認しないためである。

`Evidence <m>`はCandidateごとに一件以上必要であり、`<m>`は1から連続する10進整数とする。`raw_content`は利用者本人の加工前の該当原文であり、要約、翻訳、Codexの説明で代用してはならない。`originator.actor_kind`は必ず`user`である。AI、system、作成者不明、第三者の原文だけではこのSectionを作れない。`source_kind`は会話内発言なら`message`、ファイルから取得した原文なら`file`だけを許可し、両方や内容分類を禁止する。`source_locator`と`captured_at`は必須で、`fragment_locator`は長い出典・複数断片など、出典内の範囲を一意に示す必要がある場合に必須、出典全体が単一断片の場合のみ省略できる。`content_digest`は内容照合値を得られる場合だけ置く。これらがなければ、後続が原文、位置、作成者、取得時点を追跡できず、同じ断片の重複採用やAI内容の誤採用を防げない。

`evidence_kind`は観測した利用者本人の行為の分類、`stance`は命題の正しさへの向き、`knowledge_report`は利用者本人の知識についての明示申告、`strength`はその一件の根拠がどれほど明確かを表す。分類軸が違うため互いの代用にしてはならない。`knowledge_self_report`の場合は`knowledge_report`を必須とし、`stance`を置かない。それ以外の`evidence_kind`では`stance`を必須とし、`knowledge_report`を置かない。質問だけを`reports_unknown`へ変換してはならない。`strength`は全Evidenceの件数、利用者の熟達度、最終知識状態、Candidateの優先度ではない。

`Normalization rationale`は必須であり、Codexが何を原文から分けたか、何を同一Candidateへ統合したか、何を未観測として残したかを説明する。これはEvidenceではなく、`raw_content`の代用にも保存根拠にもならない。後続が欠損値を推測で補ったり、正規化の境界を読み違えたりしないために必要である。

## 4. `no-candidate` の結果

`result_kind`が`no-candidate`である結果は、完了済みEpisode、原文、出典情報がそろい、Episode全体を採否規則に従って評価した結果、利用者本人由来で採用条件を満たすEvidence候補が一件もなかった場合だけ使う。利用者本人が何も知らない、知識状態が`reported_unknown`である、保存・更新を試みて失敗した、という意味ではない。入力不足・未終了・出典解決不能・前後関係の切断が一つでもある場合は使わず、第5章の`input_insufficient`を使う。

```markdown
## No candidate
- evaluation_completed: `true`
- candidate_count: `0`
- exclusion_summary:
  - `<採用しなかった原文の分類と、その原文が候補にならない理由>`
```

`evaluation_completed`は必ず`true`、`candidate_count`は必ず`0`とする。`exclusion_summary`は一件以上とし、AIのみの内容、質問・依頼のみ、露出のみ、本人以外が作成した原文、非知識的な会話・操作のうち該当した分類と理由を記す。個々の不採用原文をすべて複写する必要はないが、入力不足を不採用として隠してはならない。この概要が必要なのは、候補がない理由を「評価済みだが採用対象なし」として検証可能にし、L5が入力回復の制御へ誤分岐しないためである。

## 5. `input_insufficient` の結果

`result_kind`が`input_insufficient`である結果は、Candidateの有無を安全に判断する必須入力が欠ける、または評価範囲を確定できない場合だけ使う。代表例はEpisode未終了、`episode_id`又は`ended_at`の欠落、原文を再取得できない、必須の`SourceReference`がない、`source_locator`を解決できない、前後関係が切断されて訂正を確認できない場合である。評価途中、AIのみの原文、質問のみの原文は、必要な入力がそろっていれば入力不足ではなく、採否評価の対象である。

```markdown
## Input insufficient
- evaluation_completed: `false`
- candidate_count: `0`
- diagnostics:
  - code: `<episode_not_completed | episode_identity_missing | episode_end_missing | source_reference_missing | source_unresolvable | raw_content_unavailable | context_truncated>`
    affected_input: `<不足又は解決不能な入力要素>`
    explanation: `<Candidate有無を安全に判断できない理由>`
```

`evaluation_completed`は必ず`false`、`candidate_count`は必ず`0`とする。`diagnostics`は一件以上とし、各診断は`code`、`affected_input`、`explanation`をすべて持つ。`code`は不足・解決不能の種類、`affected_input`は影響したEpisode又は出典の項目、`explanation`は候補有無を確定できない理由を表す。再試行回数、再取得方法、停止、終了、Knowledge Updateを起動するかどうかを診断へ追加してはならない。これらはL5の制御判断であり、M1が決めるとproducerの意味結果とWorkflow制御を混同するためである。

## 6. 規範例

以下は構造と値の選び方を示す規範例であり、実行Fixture、Dataset、テスト入力ではない。例の内容から利用者本人の実際の知識状態を推測してはならない。

### 6.1 Candidateが二件ある例

```markdown
# Knowledge Acquisition Result

## Contract
- contract_name: `knowledge-acquisition-candidate-markdown`
- contract_version: `1.0`

## Episode
- episode_id: `episode-example-001`
- ended_at: `2026-08-12T09:30:00Z`

## Result
- result_kind: `candidate`

## Candidate 1

### Assertion
- statement: `Goのsliceへ要素を追加すると、容量不足時に新しい配列が確保される。`
- language: `ja`

### Concept 1
- canonical_label: `Go slice`
- normalized_label: `go slice`
- aliases: [`slice`]

### Scope
- facets:
  - kind: `language`
    original_label: `Go`
    normalized_label: `go`

### Temporal
- observed_at: `2026-08-12T09:20:00Z`

### Evidence 1
- evidence_kind: `explanation`
- stance: `supports`
- strength: `moderate`
- raw_content: |
    Goのsliceへ要素を追加すると、容量が足りない場合は新しい配列が確保される。
- content_language: `ja`
- source_kind: `message`
- source_locator: `conversation://example/episode-001`
- fragment_locator: `message-12`
- originator:
  - actor_kind: `user`
- episode_id: `episode-example-001`
- captured_at: `2026-08-12T09:20:00Z`
- observed_at: `2026-08-12T09:20:00Z`

### Normalization rationale
原文の容量不足時の再配列だけを独立命題にした。配列共有の説明は別に評価できるため、Candidate 2へ分けた。製品版と有効期間は原文で観測していない。

## Candidate 2

### Assertion
- statement: `Goのsliceは元の配列を共有する場合がある。`
- language: `ja`

### Concept 1
- canonical_label: `Go slice`
- normalized_label: `go slice`
- aliases: [`slice`, `配列共有`]

### Scope
- facets:
  - kind: `language`
    original_label: `Go`
    normalized_label: `go`

### Temporal
- observed_at: `2026-08-12T09:20:00Z`

### Evidence 1
- evidence_kind: `explanation`
- stance: `supports`
- strength: `weak`
- raw_content: |
    元の配列を共有する場合がある。
- content_language: `ja`
- source_kind: `message`
- source_locator: `conversation://example/episode-001`
- fragment_locator: `message-12:second-sentence`
- originator:
  - actor_kind: `user`
- episode_id: `episode-example-001`
- captured_at: `2026-08-12T09:20:00Z`
- observed_at: `2026-08-12T09:20:00Z`

### Normalization rationale
同じ発言内の配列共有についての説明を、Candidate 1の再配列条件とは別に評価できる命題として分けた。共有する条件は原文から特定できないため追加していない。
```

### 6.2 十分な評価後の候補なしの例

```markdown
# Knowledge Acquisition Result

## Contract
- contract_name: `knowledge-acquisition-candidate-markdown`
- contract_version: `1.0`

## Episode
- episode_id: `episode-example-002`
- ended_at: `2026-08-12T10:00:00Z`

## Result
- result_kind: `no-candidate`

## No candidate
- evaluation_completed: `true`
- candidate_count: `0`
- exclusion_summary:
  - `質問・依頼のみ`: 利用者本人は説明を求めたが、知識内容又は明示的な不知を示していない。
  - `AIのみの内容`: AIが説明を提示したが、利用者本人による別の説明、推論、訂正、申告は観測されなかった。
```

### 6.3 入力不足の例

```markdown
# Knowledge Acquisition Result

## Contract
- contract_name: `knowledge-acquisition-candidate-markdown`
- contract_version: `1.0`

## Episode
- episode_id: `episode-example-003`
- ended_at: `2026-08-12T10:30:00Z`

## Result
- result_kind: `input_insufficient`

## Input insufficient
- evaluation_completed: `false`
- candidate_count: `0`
- diagnostics:
  - code: `context_truncated`
    affected_input: `conversation://example/episode-003, message-18の前後文脈`
    explanation: `後続の訂正を確認できず、原文が候補条件を満たすか安全に判断できない。`
```

## 7. 後続へ渡す確定事項

| 直接の後続Task | 今回渡すもの | この成果物を渡せるか | そのTaskが今すぐ開始できるか |
| --- | --- | --- | --- |
| `L3-M2-S1` | 三つの`result_kind`、Candidateの必須項目、入力不足診断と、Candidateだけを更新評価へ渡す境界 | はい。受入側は本書の構造を読取り入力にできる | 未確認。本Task以外に必要な`L1`入力・統合状況は、この契約だけでは確認できない |
| `L3-M2-S2` | 正規化済み命題、Concept、Scope、時点、原文・出典・正規化理由を含むCandidate境界 | はい。既存Knowledge探索・意味比較の読取り入力にできる | 未確認。本Task以外の着手依存と統合状況は、この契約だけでは確認できない |
| `L5-M1-S2` | 契約名・版、三結果、候補境界、入力不足診断 | はい。子Skillの成果物引渡し契約へ参照できる | 未確認。ほかのproducer契約が必要であるため、この契約だけで着手可否を決めない |
| `L5-M3-S1` | Episode識別子、終了時点、Acquisition成果物の三結果 | はい。Knowledge蓄積Workflowの起動・受入設計へ参照できる | 未確認。固定Task Mapは`L5-M1`も入力とするため、本書だけでは全依存を満たさない |
| `L5-M3-S2` | `candidate`、`no-candidate`、`input_insufficient`の意味を変えない分岐入力 | はい。分岐TaskはCandidate内容を再判定せずに受け取れる | 未確認。`L5-M3-S1`と共通制御契約も必要なため、本書だけでは全依存を満たさない |
| `L6-M3-S6` | producer契約としての文書構造、Candidate境界、三結果と診断 | はい。後続Adapterの読取り入力にできる | 未確認。ほかのL3・L4 producer契約とHarness Coreが必要なため、本書だけでは全依存を満たさない |

この表の「渡せる」は、今回のMarkdown契約を後続が入力として読めることを意味する。「今すぐ開始できる」は、後続自身のすべての着手依存とGateが満たされているかという別の判定である。二つを混同すると、一つの契約が完成しただけで後続全体を開始可能又は不可能と誤判定するため、分けて記録する。

## 8. 完全性確認

- Candidateあり、十分な評価後の候補なし、入力不足による評価不能を排他的に表し、`no-candidate`を利用者本人の不知や入力障害へ読み替えない。
- Candidateごとに一命題、Concept、条件、時点、利用者本人の加工前原文、出典、正規化理由を分離し、後続が原文へ戻れるようにした。
- 順序、複数Candidate境界、必須・条件付き項目、禁止値を定め、未観測値を推測で補わないようにした。
- L3-M2による比較・更新、L5による再試行・終了・非起動制御、L6のFixture・評価を本書から除外し、責務を重複させない。
