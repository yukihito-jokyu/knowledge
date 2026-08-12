# Candidate正規化とL1 Field写像設計

## 1. 目的

本書は、完了した`KnowledgeEpisode`（評価対象の会話または作業のまとまり）から、利用者本人が示した内容を`Candidate`（まだ保存済みKnowledgeではない、更新評価へ渡す候補）として整理し、L1の論理データモデルが必要とする項目へ対応付ける設計である。

正規化とは、意味を変えずに表記揺れを抑え、検索・比較・後続の判断を可能にするよう整理することである。原文を短く要約して捨てることではない。候補の命題本文、検索用索引語、根拠原文、適用範囲を一つの文章や一つの項目へ混在させると、利用者本人が何を知っているかを再評価できなくなる。そのため本書は、それぞれを別のField（保存する情報の一項目）へ写像する。

本書の目的はCandidateを新規保存することではない。Candidateは後続のKnowledge Updateが既存Knowledgeと意味比較し、`create`、`attach-evidence`、`revise`、`supersede`または更新しない、のいずれを選ぶための入力である。従って、候補生成だけで新しいKnowledge Assertion（知識判定用の命題）やEvidence（採用済みの根拠）を永続化してはならない。

## 2. 入力、出力、責務境界

### 2.1 入力

入力は、L3-M1-S1が採用可能と判定した利用者本人由来のEvidence候補と、それを解釈するために必要な完了済みEpisodeの前後関係である。Candidateを作る前に、各候補について原文の提示者、出典位置、Episodeの完了境界を確認する。

AIの説明、利用者が情報を見た事実、質問だけは、単独では利用者本人の知識を示すEvidence候補にしない。これは「露出」と「知識獲得」を混同しないためである。利用者本人が説明、推論、コード、訂正、自己申告、または技術判断として示した内容だけを、出典とともに候補化できる。

### 2.2 出力

出力は、独立して評価できる一つ以上のCandidateである。一つのCandidateは一つの正規化済み命題と、その命題を判断できる根拠・文脈・検索情報の組である。複数の独立命題が一つの出典断片にある場合は、命題ごとにCandidateを分ける。同じ命題に複数の根拠断片がある場合は、一Candidateへ複数の根拠候補を関連付けてよい。

Candidateが作れない場合の`no-candidate`、入力が不足して判定できない場合の`input_insufficient`のMarkdown表現はL3-M1-S3の責務である。本書は、それらを選ぶための最終出力形式を決めない。

### 2.3 Codex、CLI、後続Taskの担当

| 担当 | 行うこと | 行わないこと | 分ける理由 |
| --- | --- | --- | --- |
| Codex | 原文の意味、命題の粒度、同一性候補、根拠種別・向き・強さ、Scope、時点情報を判断する | 文字列一致や件数だけで知識状態を決めない | 意味は文脈を読まなければ判断できないため |
| CLI（コマンドライン処理） | 指定されたFieldの必須性、型、列挙値、参照形式を検証し、保存・取得する | Candidateから意味を推測し、既存Knowledgeとの同一性を決めない | 同じ入力から同じ検証結果を返す決定的処理へ意味判断を混入させないため |
| L3-M1-S3 | Candidate Markdown、`no-candidate`、`input_insufficient`の引渡し形式を定義する | 本書のField意味を変更しない | 他Taskが読める出力契約と、候補の意味的整理を分離するため |
| L3-M2 | Candidateを既存Knowledgeと比較し、更新操作または非更新を決める | Candidate生成時に既存Knowledgeを更新しない | 候補の抽出と更新判断を分離し、重複作成や履歴破壊を防ぐため |

## 3. Candidateの最小構造

CandidateはL1の保存レコードではなく、L1へ渡す値を欠落なく揃える作業中の構造である。識別子、監査情報、保存操作IDは後続の更新処理が付与する。候補段階で勝手にIDを確定すると、既存レコードとの同一性を未確認のまま新規作成を前提にしてしまうためである。

| Candidateの要素 | 意味 | 必要な理由 | L1への主な写像 |
| --- | --- | --- | --- |
| `assertion` | 利用者本人が知っているかを単独で評価できる、正規化済みの一命題 | 粗い能力評価や複数主張を比較単位にしないため | `KnowledgeAssertion.statement`、`language`、`scope`、`temporal` |
| `concepts` | 命題へ到達するための検索用索引語と別名 | 命題本文と異なる語句、API名、略称から検索できるようにするため | `Concept.canonical_label`、`normalized_label`、`aliases`、`scope` |
| `evidence_candidates` | 利用者本人の加工前の発言・コード・訂正・申告・判断と、その分類 | 命題への整理誤りや状態判断を、原文から再確認できるようにするため | `Evidence`の各Field |
| `source_context` | 出典断片、原文の提示者、Episode、取得時点 | 原文だけでは不足する前後関係と出典追跡を保持するため | `SourceReference`、`KnowledgeEpisode` |
| `scope` | 技術分野、プログラミング言語、対象製品などの適用条件 | 同じ表現でも対象が違う内容を混同しないため | `Scope.facets` |
| `temporal` | 有効期間、製品版、観測時点、確認時点 | 版差・時点差を現在も有効な内容や矛盾と取り違えないため | `TemporalMetadata` |
| `normalization_rationale` | Codexがどの原文から何を分離・統合し、何を未観測として残したかの説明 | 後続が推測で値を補ったり、正規化の妥当性を検証不能にしたりしないため | Candidateの引渡し説明。L1の`Evidence.raw_content`の代用にはしない |

## 4. 正規化手順

### 4.1 一命題へ分解する

1. Evidence候補の原文とEpisodeの前後関係を読み、利用者本人が実際に示した内容を確認する。
2. 「利用者本人がこの内容を知っているか」を他の内容と独立して評価できる最小の文へ分ける。
3. 原文の言い回しをそのまま複写せず、意味を保つ平叙文として`assertion.statement`へ整理する。
4. 原文にない条件、因果、例外、製品版、時点を追加しない。観測できない値は空欄または未観測として残す。

たとえば「Goのsliceにappendすると、容量が足りなければ新しい配列になる。だから元の配列を共有する場合がある」という説明は、再配列の条件と共有の可能性が別々に評価できるなら二つのCandidateへ分ける。反対に、単独では意味を成さない限定条件は、その命題のScopeまたは本文に残す。接続詞で複数命題を束ねると、既知部分と未知部分を区別できないためである。

### 4.2 原文と正規化済み命題を分離する

`Evidence.raw_content`には利用者本人の加工前の発言、コード、訂正、判断を該当範囲のまま入れる。`KnowledgeAssertion.statement`には、その原文から整理した一命題を入れる。Codexの説明文、要約、判断理由で`raw_content`を置き換えてはならない。

この分離が必要なのは、後から「どの発言をどの命題へ整理したか」を検証し、根拠同士の競合や訂正を再評価するためである。`normalization_rationale`は整理の理由であり、本人の原文でも根拠そのものでもない。

### 4.3 検索用索引語を抽出する

`Concept`は「利用者本人がその語を知っている」という命題ではない。命題の検索入口である。命題本文、原文、または明確な技術用語から、関連命題へ到達するために必要な語だけを抽出する。

- 代表表記を`canonical_label`、検索・比較用の表記を`normalized_label`へ写像する。
- 略称、API名、綴り違い、別言語表記は、同じ意味を指すとCodexが説明できる場合だけ`aliases`へ入れる。
- 同じ文字列でも技術分野や対象が異なる場合は、Scopeを分けて別のConcept候補とする。
- 命題本文、根拠の強さ、知識状態、Codexの推測をConceptへ入れない。

### 4.4 根拠候補を分類する

一つの原文断片ごとに、以下の`evidence_kind`（利用者本人の何を観測したかを表す分類）を一つ選ぶ。根拠として使う役割と、知識状態を決めることは別である。

| `evidence_kind` | 選ぶ条件 | 必要な理由 |
| --- | --- | --- |
| `explanation` | 仕組み、意味、条件を利用者本人が説明した | 対象命題への直接理解を確認するため |
| `reasoning` | 前提から結論へ至る推論を示した | 結論だけでなく導出過程を分けて評価するため |
| `code` | 利用者本人がコードを作成または提示した | 知識を実際に適用した観測を残すため |
| `correction` | 以前の説明・判断の何を誤りとし何へ直すかを明示した | 訂正前の根拠を消さず、変更の意味を追跡するため |
| `knowledge_self_report` | 知っている、知らない、認識した、経験したと本人が申告した | 命題の正しさへの向きと本人の知識申告を混同しないため |
| `technical_decision` | 複数案を比較し、採否と理由を示した | 知識を意思決定へ適用した行為を識別するため |

`stance`は命題の正しさに対する向きで、`supports`（支持）、`challenges`（反証）、`qualifies`（条件・例外による限定）のいずれかである。`knowledge_self_report`以外で必須とする。`knowledge_report`は本人の知識に関する申告で、`reports_known`、`reports_unknown`、`recognizes`、`reports_experience`のいずれかであり、`knowledge_self_report`でのみ必須とする。質問したことだけを`reports_unknown`にしてはならない。

`strength`（一件の根拠が知識状態をどれほど明確に示すか）は、`strong`、`moderate`、`weak`から選ぶ。これは利用者の熟達度、根拠件数、最終状態の確かさではない。明示的な不知申告は`knowledge_self_report`、`reports_unknown`、`moderate`であり、`stance`を持たない。意味確認済みの明示的訂正は`correction`、`strong`とするが、訂正内容を推測できるだけでは`strong`にしない。

### 4.5 Scopeと時点情報を付ける

Scopeは、`domain`（技術・業務分野）、`language`（プログラミング言語）、`subject`（製品、ライブラリ、機能）の`ScopeFacet`からなる。異なる種類のFacetはすべて満たすAND条件、同じ種類の複数値はそのいずれかのOR条件である。原文の自然言語はScopeではなく、命題の`language`または根拠の`content_language`へ入れる。

時点情報は、観測した時点をEvidenceの`temporal.observed_at`へ必ず写像する。製品版や有効期間が原文で確認できる場合だけ、`valid_from`、`valid_until`、`version_scope`を付ける。空のScopeや省略した時点情報は、全範囲・全期間で有効という意味ではなく、追加情報を観測していないという意味である。

## 5. CandidateからL1 Fieldへの写像

### 5.1 `KnowledgeAssertion`への写像

| Candidateの情報 | L1 Field | 写像規則 | 禁止・制約 |
| --- | --- | --- | --- |
| 正規化済みの一命題 | `KnowledgeAssertion.statement` | 一つの独立評価可能な平叙文へ整える | 粗い能力ラベル、複数命題の束、根拠原文の複写を入れない |
| 命題の自然言語 | `language` | 原文ではなく`statement`の言語を設定する | プログラミング言語名を入れない |
| 適用条件 | `scope.facets` | 観測できた`domain`、`language`、`subject`を原表記と正規化表記で設定する | 未観測の条件を補わない。空は全範囲を意味しない |
| 有効期間・製品版 | `temporal.valid_from`、`valid_until`、`version_scope` | 原文に根拠がある条件だけを設定する | 省略を「現在有効」と断定しない |
| 保存時の識別・監査 | `assertion_id`、`audit` | Candidate段階では未確定。後続の更新操作で付与する | 新規IDや作成日時を候補生成側で確定しない |

### 5.2 `Concept`への写像

| Candidateの情報 | L1 Field | 写像規則 | 禁止・制約 |
| --- | --- | --- | --- |
| 人が読める代表語 | `canonical_label` | 命題を検索する中心となる用語を選ぶ | 命題全文や知識状態を入れない |
| 検索・比較用の語 | `normalized_label` | 表記揺れを抑えた表記を作る | 意味が違う語を同一表記へ統合しない |
| 略称・API名・綴り違い | `aliases` | 同じConceptを指すと説明できる表記だけ追加する | 代表表記のない別名だけを作らない |
| 適用条件 | `scope` | 命題と同じ条件を機械的に複写せず、当該用語の意味を区別する条件を設定する | 同名異義をScopeなしで統合しない |
| 保存時の識別・監査 | `concept_id`、`audit` | Candidate段階では未確定 | Candidateだけで永続化しない |

### 5.3 `Evidence`と出典への写像

| Candidateの情報 | L1 Field | 写像規則 | 禁止・制約 |
| --- | --- | --- | --- |
| 対象命題 | `assertion_id` | 後続が確定した命題IDへ関連付ける | Candidate段階で仮IDを保存しない |
| 観測行為 | `evidence_kind`、`stance`、`knowledge_report`、`strength` | 第4.4節の分類を使う。条件付きFieldは該当しないとき省略する | 質問を不知申告へ、反証を不知申告へ読み替えない |
| 加工前の本人表現 | `raw_content`、`content_language` | 出典断片の該当範囲をそのまま保持する | Codexの要約・説明で置き換えない |
| 出典の取得形式・位置 | `source.source_kind`、`source_locator`、`fragment_locator` | 会話発言なら`message`、ファイルなら`file`を設定し、該当位置を可能な限り示す | `message`と`file`を同時に設定しない |
| 原文の提示者と文脈 | `source.originator`、`source.episode_id`、`source.captured_at` | 提示者が`user`であること、完了Episode、取得時点を記録する | Codex・Systemだけの出典をEvidenceにしない |
| 内容照合 | `source.content_digest` | 出典内容を機械照合できるときだけ設定する | 照合値がないことを内容の否定と扱わない |
| 観測・有効情報 | `temporal.observed_at`、必要に応じて他の`temporal`属性 | `observed_at`は必須。出典取得時点以前とする | 観測していない製品版・期間を推測しない |
| 保存時の識別・監査 | `evidence_id`、`audit` | Candidate段階では未確定 | 候補化と同時にEvidence正本を作らない |

## 6. 正規化時の不変条件

1. Candidateの存在は、利用者本人がその命題を`known`（知っている）であることを意味しない。知識状態はEvidenceから導出する別の結果である。
2. 根拠がないことは、利用者本人が知らないことを意味しない。未観測と明示的な不知申告を混同しない。
3. 似た文字列、同じConcept、検索結果の近さだけで、命題の意味的同一性を確定しない。同一性と更新操作は後続のL3-M2が判断する。
4. Evidenceの原文と出典を削除、上書き、Codexの説明への置換をしない。訂正があっても旧根拠を残す。
5. `strength`、`confidence`、利用者の熟達度を混同しない。Candidateでは根拠一件の明確さだけを`strength`として記録する。
6. Scope・時点情報が未観測であることを、無条件に有効であること、または無効であることへ変換しない。
7. Candidateは既存Knowledgeの更新を実行しない。更新対象の検索、意味比較、履歴操作、保存はL3-M2およびL1/L2の責務である。

## 7. 具体例

### 7.1 説明から二つのCandidateへ分ける例

利用者本人の原文が「Goのsliceへ要素を追加すると、容量が足りない場合は新しい配列が確保される。元の配列を共有する場合がある」であるとする。この一文には、再配列の条件と配列共有の可能性という別々に評価できる内容がある。

| 要素 | Candidate A | Candidate B |
| --- | --- | --- |
| `assertion.statement` | Goのsliceへ要素を追加すると、容量不足時に新しい配列が確保される。 | Goのsliceは元の配列を共有する場合がある。 |
| `concepts` | Go、slice、append | Go、slice、配列共有 |
| `Evidence.raw_content` | 利用者本人の原文のうち再配列を説明する範囲 | 利用者本人の原文のうち共有を説明する範囲 |
| `evidence_kind` / `stance` | `explanation` / `supports` | `explanation` / `supports` |
| `strength` | 原文が条件と結果を明確に説明しているかをCodexが判断する | 原文が共有の条件をどこまで説明しているかをCodexが判断する |
| Scope | `language=Go`。他の条件は観測できた場合のみ追加 | `language=Go`。他の条件は観測できた場合のみ追加 |

同じ出典断片を二Candidateが参照しても、原文を一つの複合命題へ潰さない。後続の更新処理が、既存Knowledgeとの同一性をCandidateごとに比較できるようにするためである。

### 7.2 不知申告を質問から分ける例

利用者本人が「`append`の再配置条件は分からない。どう確認すればよいか」と発言した場合、前半に明示的不知申告があり、後半は質問である。

- 不知申告のCandidate根拠は、`evidence_kind=knowledge_self_report`、`knowledge_report=reports_unknown`、`strength=moderate`、`stance`なしとする。
- 「どう確認すればよいか」は質問であり、単独で新しいEvidenceや追加の不知申告にはしない。
- このCandidateの存在だけで`reported_unknown`を保存しない。状態の十分条件と競合確認はL1-M1-S3および後続更新処理が扱う。

## 8. 後続Taskへの引継ぎ

L3-M1-S3は、各Candidateについて少なくとも本書の`assertion`、`concepts`、`evidence_candidates`、`source_context`、`scope`、`temporal`、`normalization_rationale`を読めるMarkdown契約を定義する必要がある。値が未観測の場合は、推測値で埋めず、その事実を読者が区別できる表現にする。

L3-M2は、Candidateを受け取った後に既存Knowledgeを探索し、意味的同一性、更新操作、保存可否を判断する。本書の正規化済み`statement`、Scope、時点、原文・出典が、その比較と履歴保持に必要な入力である。Candidate正規化の完了は、L3-M2の開始または保存実行の許可を意味しない。

## 9. 確認表

| 確認対象 | 本書での確認箇所 |
| --- | --- |
| Candidateが独立評価可能な一命題である | 第3章、第4.1節、第5.1節 |
| 命題、Concept、Evidence、状態を混同しない | 第3章、第4.2節から第4.4節、第6章 |
| 原文・出典・Episodeへ戻れる | 第3章、第5.3節 |
| Scope、時点、製品版の未観測を推測しない | 第4.5節、第5章、第6章 |
| CodexとCLI、候補化と更新の責務が分離される | 第2.3節、第6章、第8章 |
| 後続L3-M1-S3が必要な引渡し情報を得られる | 第8章 |
