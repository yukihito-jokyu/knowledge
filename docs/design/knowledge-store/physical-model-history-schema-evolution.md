# 正本データの物理モデル・履歴・スキーマ変更設計

## 1. 目的とこの文書の位置

この文書は、利用者本人が示した知識の根拠を失わずに保存し、後から「どの根拠・どの命題・どの更新によって現在の判断になったか」を復元できる SQLite 3 の物理構造を定める。ここでいう**物理モデル**は、論理モデルが定めた情報を、SQLite の表、列、外部キー、一意制約としてどのように保存するかという設計である。表や列を先に決める理由は、実装時に履歴、参照、検索再構築に必要な情報を取り落とさないためである。

この設計の正本は、検索順位、全文検索表、ベクトル、キャッシュではない。正本は、利用者本人による根拠、知識判定用の命題、命題間又は索引語との関係、更新操作とその変更前後、計算済み状態である。検索用の派生索引は正本から作り直せなければならず、本書の表を唯一の復元経路とする。

| 固定入力 | この設計で守ること | 必要な理由 |
| --- | --- | --- |
| `adr/0001-persistence-stack.md` | SQLite 3、SQL制約、明示的トランザクション、版付きSQL移行を使う | 採用済みの保存方式と矛盾せず、参照整合性と原子的な更新をデータベース側でも検査するため |
| `docs/design/knowledge-model/` の論理設計 | 命題、根拠、状態、関係、更新履歴の意味・許可値・不変条件を変えない | 物理表の都合で利用者の知識判断の意味を変更しないため |
| L2-M2-S1 の検索・再構築要件 | 正本変更境界までの安定した列挙と、派生索引を再構築できる入力を残す | 索引障害後にも候補の意味と根拠を説明できるようにするため |

この文書は `L2-M1-S2` が所有する詳細設計である。SQLの初回定義は `L2-M1-S4`、移行の適用履歴と実行器は `L2-M1-S5`、トランザクションの開始・競合エラー・Repository（保存処理を業務処理から分ける内部操作窓口）は `L2-M1-S3` が所有する。本書はそれらが実装すべき表構造と不変条件を渡すが、個別SQL文、ドライバ、再試行回数は固定しない。

## 2. 用語、保存形式、共通規則

| 用語・識別名 | この文書での意味 | 必要な理由 |
| --- | --- | --- |
| `assertion`（命題） | 利用者本人が知っているかを一件ずつ判定する、正規化済みの具体的な文 | 単なる語句や根拠原文と、知識判定の対象を混同しないため |
| `evidence`（根拠） | 利用者本人が実際に示し、命題の判断へ採用した加工前の発言・コード・訂正など | AIの説明や検索語だけから利用者の知識を誤判定しないため |
| revision（版） | 同じ業務上のIDについて、ある時点の完全な値を表す不変の行 | 記録誤りを直しても旧値を上書きせず、変更前後を比較できるようにするため |
| commit（正本変更境界） | 一つの更新操作で確定した変更集合に、単調増加する順序を付けた記録。GitのCommitとは別物である | 指定時点までの列挙と二時点間の差分を、表示順や時計だけに依存せず再現するため |
| snapshot（完全な記録値） | 変更時点のレコード全属性をJSONとして固定した値。差分や要約ではない | 後から表構造が変わっても、当時の変更前後の意味を復元するため |

### 2.1 SQLiteの型と時点

- 識別子はすべて `TEXT` とする。識別子の発行方式は `L2-M1-S6` が決めるため、本書は値の見た目を仮定しない。空文字列は禁止する。
- 時点は、UTCのISO 8601文字列を `TEXT` で保存する。文字列比較と時系列順を一致させ、SQLite固有の日付解釈に依存しないためである。保存時にUTC・秒未満の精度を統一する方法は実装側が定める。
- 真偽値は `INTEGER` の `0` 又は `1`、列挙値は `TEXT` とし、`CHECK` 制約で許可値を限定する。SQLiteには専用の列挙型がないため、未定義値を保存しないために必要である。
- 複数値又は入れ子の値は、検索・参照・一意性が必要なものを子表へ正規化し、当時の完全値を残す履歴だけをJSONにする。JSONだけにすると外部キーと重複検査ができず、すべてを列に展開すると変更前後の完全復元が難しくなるためである。
- 全接続で `PRAGMA foreign_keys = ON` を有効にする。外部キーを定義しても接続側で無効なら、存在しない根拠や命題を参照できてしまうためである。

### 2.2 削除と更新

正本表の行を物理削除しない。`revise`（同じ評価対象の記録誤りを訂正する更新）は新revisionを追加し、`supersede`（別の評価対象として旧命題を残したまま新命題へ置換する更新）は新しい命題と新→旧の関係を追加する。古い根拠、状態、参照を新しい命題へ自動で付け替えない。これにより、過去に何を観測し、どの判断をしたかを失わない。

`current`（現在版）を示す列は更新対象の本表へ置かず、`*_revision` 表で業務IDごとに `is_current = 1` を一件だけ許す部分一意索引で表す。現在値を本表へ複写すると履歴値とずれるためである。SQLiteの部分一意索引は `WHERE is_current = 1` で実現する。

## 3. 表の全体像

```text
knowledge_episode ─┐
                    ├─ evidence ─ evidence_revision ─ evidence_source
assertion ─ assertion_revision ─ assertion_scope_facet / assertion_temporal_version
     │  ├─ assertion_concept ─ concept ─ concept_alias / concept_scope_facet
     │  ├─ derived_knowledge_state ─ state_evidence / state_aspect / state_gap
     │  └─ state_publication ─ state_publication_revision
     └─ knowledge_relation ─ knowledge_relation_revision

knowledge_update_operation ─ knowledge_update_attempt ─ record_change
                            └─ state_publication_change
store_commit
```

矢印は外部キーによる参照を表す。`store_commit` は変更の確定順序を表し、各revision、操作試行、変更事実が参照する。検索索引の表、埋め込み値、索引世代の表はこの図に含めない。これらは正本を使う側の `L2-M2` が所有し、正本の唯一の復元経路になってはならない。

## 4. 正本の内容を保存する表

### 4.1 正本変更境界 `store_commit`

| 列 | 型・制約 | 意味と必要な理由 |
| --- | --- | --- |
| `commit_id` | `TEXT PRIMARY KEY`、空文字列禁止 | 一つの確定境界を指すID。外部のGit Commitと混同せず、各変更を同じ確定単位へ結び付けるため |
| `commit_sequence` | `INTEGER UNIQUE NOT NULL`、1以上 | 確定した順序。時計のずれや同時刻でも安定して列挙・差分取得するため |
| `operation_id` | `TEXT`、操作が由来なら外部キー | この境界を生んだ更新要求。初期データ又は保守移行はNULLを許し、利用者の知識更新との誤った対応付けを避けるため |
| `committed_at` | UTC時点、`NOT NULL` | 境界を確定した時点。各revisionの監査時点との順序を検査するため |
| `commit_kind` | `TEXT CHECK ('knowledge_update','state_recalculation','schema_maintenance')` | 知識変更、状態再計算、物理保守を区別する。派生索引側が再構築の対象範囲を判断できるようにするため |

`commit_sequence` を手で書き換えず、確定トランザクション内で次の値を一度だけ割り当てる。並行書込みの取得方法は `L2-M1-S3` が定める。

### 4.2 命題と適用範囲

`assertion` は業務上の同一性だけを保持し、表示・検索に使う内容は `assertion_revision` に保存する。本文、適用範囲、時点を直しても同じ評価対象なら同じ `assertion_id` を使い、新revisionを追加するためである。

| 表・列 | 型・制約 | 意味と必要な理由 |
| --- | --- | --- |
| `assertion.assertion_id` | `TEXT PRIMARY KEY` | 命題の業務上のID。根拠、状態、関係が本文変更後も同じ対象を参照するため |
| `assertion_revision.assertion_revision_id` | `TEXT PRIMARY KEY` | 命題の一版を指すID。変更前後を別行として保存するため |
| `assertion_revision.assertion_id` | `TEXT NOT NULL REFERENCES assertion` | この版が属する命題。別命題のrevisionを混ぜないため |
| `assertion_revision.previous_revision_id` | 同じ命題のrevisionを参照、初版のみNULL | 直接の変更前版。複数の並行更新を黙って併合せず、revision鎖を追えるようにするため |
| `assertion_revision.statement` | `TEXT NOT NULL`、空白だけ禁止 | 一件で知識判定できる正規化済み命題。根拠原文や概念名では代用できないため |
| `assertion_revision.language_tag` | `TEXT NOT NULL` | 命題本文の自然言語。`scope`のプログラミング言語とは別に、表示・比較方法を選ぶため |
| `assertion_revision.created_at` / `created_by_kind` / `created_by_id` | 時点、主体種別、条件付き主体ID | 初版の成立時点と内容決定者。AIや自動処理が利用者本人の根拠を作ったように見せないため |
| `assertion_revision.updated_at` / `updated_by_kind` / `updated_by_id` | 時点、主体種別、条件付き主体ID | 当該版の内容を最後に決めた時点と主体。作成時とは別に監査するため |
| `assertion_revision.store_commit_id` | `NOT NULL REFERENCES store_commit` | この版を確定した正本変更境界。差分列挙と変更履歴を結ぶため |
| `assertion_revision.is_current` | `INTEGER CHECK (0,1)` | 現在取得で選ぶ版。命題ごとに1件だけにして、旧版を通常取得へ混入させないため |

`assertion_scope_facet` は `(assertion_revision_id, dimension, normalized_value)` を一意にし、`dimension` は `domain`（技術・業務分野）、`language`（プログラミング言語）、`subject`（製品・ライブラリ等）だけを許す。`value` は人が読む原表記、`normalized_value` は比較用表記であり、片方だけでは表示又は表記揺れの照合を失う。空集合は「全範囲」ではなく「追加の適用範囲を観測していない」を表す。

`assertion_temporal` は一版につき0又は1行とし、`valid_from`、`valid_until`、`observed_at`、`last_verified` を保存する。`valid_from < valid_until` を満たす場合だけ両値を許し、`valid_until` 自体は有効期間に含めない。`assertion_temporal_version` は `(assertion_revision_id, version_value)` を一意にして、適用製品版を重複なく保存する。時点、製品版、保存時点を一列へ混ぜると、古い内容と古い記録を区別できないためである。

### 4.3 検索用索引語 `concept`

`concept` は命題そのものではなく、語句から関連命題を探すための正本の語彙である。概念の存在だけを、利用者本人が知っている証拠として扱わない。

| 列 | 型・制約 | 意味と必要な理由 |
| --- | --- | --- |
| `concept.concept_id` | `TEXT PRIMARY KEY` | 概念の業務ID。別名や関係から一貫して参照するため |
| `concept.canonical_label` / `normalized_label` | `TEXT NOT NULL` | 代表表記と比較表記。人が読める名称と表記揺れを抑える照合値を分けるため |
| `concept.language_tag` | `TEXT`、任意 | 代表表記の自然言語。言語により同じ文字列が別概念になる誤照合を抑えるため |
| `concept.created_at`、`created_by_*`、`updated_at`、`updated_by_*` | 命題と同じ監査構造 | 内容を誰がいつ決めたかを追跡し、保存実行者と混同しないため |
| `concept.store_commit_id` | `NOT NULL REFERENCES store_commit` | 作成又は更新を確定した境界。索引再構築の入力時点を特定するため |

`concept_alias` は原表記、比較表記、別名種別、任意の自然言語を持ち、`(concept_id, normalized_text, COALESCE(language_tag,''), alias_kind)` を一意にする。別名種別は `synonym`（同義語）、`abbreviation`（略称）、`identifier`（内部識別名）、`api_name`（公開操作名）、`spelling_variant`（綴り・表記揺れ）に限定する。代表表記と同じ比較表記を別名として登録しない。`concept_scope_facet` は命題の適用範囲と同じ三分類・一意性を使う。

### 4.4 根拠の文脈、内容、版

`knowledge_episode` は、根拠候補を評価した会話又は作業のまとまりであり、一断片だけでは失われる前後文脈を保存する。`evidence` は採用済み根拠の業務ID、`evidence_revision` はその完全な一版である。

| 表・列 | 型・制約 | 意味と必要な理由 |
| --- | --- | --- |
| `knowledge_episode.episode_id` | `TEXT PRIMARY KEY` | 会話・作業範囲を指すID。根拠を前後文脈へ戻すため |
| `knowledge_episode.episode_kind` | `TEXT CHECK ('conversation','task')` | 会話か作業か。説明内容の分類ではなく、範囲の再取得方法を選ぶため |
| `knowledge_episode.source_locator` | `TEXT NOT NULL UNIQUE` | 範囲全体への参照先。断片だけでは評価の前提が検証できず、同じ範囲を二重登録しないため |
| `knowledge_episode.started_at` / `ended_at` | UTC時点、両方`NOT NULL`、`started_at <= ended_at` | 評価を始めてよい完了済み文脈の時間範囲。未完了範囲を完了済み根拠として使わないため |
| `knowledge_episode.summary` | `TEXT`、任意、存在時は空白だけ禁止 | 範囲を人が判別できる短い説明。参照先だけでは作業目的を理解できないため |
| `knowledge_episode.created_at` / `created_by_kind` / `created_by_id` / `updated_at` / `updated_by_kind` / `updated_by_id` | 命題と同じ監査構造 | 会話・作業範囲の記録を誰がいつ作成・更新したかを追跡するため |
| `evidence.evidence_id` | `TEXT PRIMARY KEY` | 根拠の業務ID。訂正や記録誤りの修正後も同じ観測を指すため |
| `evidence_revision.evidence_revision_id` / `previous_revision_id` / `is_current` | 命題revisionと同じ規則 | 根拠の記録誤りだけを訂正しても、採用済み原文と旧版を消さないため |
| `evidence_revision.assertion_id` | `NOT NULL REFERENCES assertion` | 根拠が支持・反証・自己申告する命題。別の命題へ自動付替えしないため |
| `evidence_revision.evidence_kind` | `TEXT CHECK ('explanation','reasoning','code','correction','knowledge_self_report','technical_decision')` | 根拠の内容分類。`reasoning`は理由付け、`knowledge_self_report`は利用者本人の知識申告、`technical_decision`は技術選択の記録を表す。出典取得形式と混同しないため |
| `evidence_revision.raw_content` | `TEXT NOT NULL`、空白だけ禁止 | 利用者本人が示した加工前の内容。要約や状態理由で置換すると再評価できないため |
| `evidence_revision.stance` | 非自己申告では `TEXT CHECK ('supports','challenges','qualifies')`、自己申告ではNULL | 命題を支持する、反対・反証候補を示す、又は適用条件を限定するという根拠の向き。自己申告の内容は別列で表し、根拠の向きと状態導出を混同しないため |
| `evidence_revision.strength` | `TEXT CHECK ('strong','moderate','weak')` | 根拠の強さ。件数ではなく内容の強さを状態判断で説明するため |
| `evidence_revision.knowledge_report` | `knowledge_self_report` のときだけ `TEXT NOT NULL CHECK ('reports_known','reports_unknown','recognizes','reports_experience')`、他種別ではNULL | 利用者本人の自己申告内容。`reported_unknown`を推測だけで作らず、申告の種類も失わないため |
| `evidence_revision.store_commit_id` | `NOT NULL REFERENCES store_commit` | 根拠版の確定境界。どの索引再構築へ含めるかを判定するため |

`evidence_source` は一つの `evidence_revision` に一件だけ置く。`source_kind` は `message`（会話の一発言）又は `file`（ファイル）だけを許し、`source_locator`、`episode_id`、`originator_kind`、`captured_at` を必須にする。`originator_kind` は必ず `user` とし、Codex又は自動処理だけが作った出典を利用者本人の根拠へ保存しない。任意の `content_digest` は元データ取り違え検出用、`fragment_locator` は発言番号・行範囲等で採用箇所を特定するために置く。出典の取得形式と `evidence_kind` の意味分類を一列へ混ぜない。

根拠の適用範囲と時点は `evidence_scope_facet`、`evidence_temporal`、`evidence_temporal_version` に、命題と同じ規則で保存する。根拠の観測時点と、データベースへ保存した時点は異なるため、前者を `observed_at`、後者を `store_commit.committed_at` として分ける。

### 4.5 状態と現在公開の選択

`derived_knowledge_state` は根拠から計算した結果であり、根拠の代替正本ではない。状態を直接改版せず、再計算のたびに新しい行を追加する。これにより、同じ根拠集合から当時どの状態を導いたかを残す。

| 表・列 | 型・制約 | 意味と必要な理由 |
| --- | --- | --- |
| `derived_knowledge_state.state_id` | `TEXT PRIMARY KEY` | 一回の計算結果を指すID。旧結果を上書きせず履歴参照するため |
| `derived_knowledge_state.assertion_id` | `NOT NULL REFERENCES assertion` | 評価対象命題。別命題の根拠を現在状態として返さないため |
| `derived_knowledge_state.status` | 8値の `CHECK` | `known`、`partially_known`、`inferable`、`contradicted`、`outdated`、`reported_unknown`、`no_evidence`、`uncertain`。再計算待ちはここへ追加しないため |
| `derived_knowledge_state.confidence` | `REAL CHECK (0.0 <= value AND value <= 1.0)` | 選んだ分類が確かである度合い。熟達度や命題の真偽ではないため |
| `derived_knowledge_state.rationale` | `TEXT NOT NULL` | 状態、使った根拠の役割、適用範囲・時点比較、競合解消を説明する理由文。原文の複写で代用しないため |
| `derived_knowledge_state.calculated_at` | UTC時点、`NOT NULL` | 計算時点。根拠の観測時点・公開時点と区別するため |
| `derived_knowledge_state.valid_from` / `valid_until` | UTC時点、各任意。両方ある場合は `valid_from < valid_until`、`valid_until`自体は有効期間に含めない | この計算結果を適用できる開始・終了境界。`calculated_at`は計算を実行した時点であり、内容が有効な期間を表さないため分離する |
| `state_temporal_version(state_id, version_value)` | `state_id REFERENCES derived_knowledge_state`、複合主キー、`version_value`は空白だけ禁止 | 適用する製品・仕様・APIの版を0件以上、重複なく保存する。製品版に依存しない場合は行を作らず、空文字列や空配列で「不明」を推測しないため |
| `derived_knowledge_state.observed_at` | UTC時点、`NOT NULL`、`calculated_at` と同値 | 計算済み状態を観測した時点。論理契約の `temporal.observed_at=evaluated_at` を物理的に強制し、根拠の観測時点や正本変更境界の時点と混同しないため |
| `derived_knowledge_state.last_verified` | UTC時点、任意 | この状態を最後に検証した時点。計算時点・有効期間・保存時点とは異なり、再検証の必要性を判断するための情報である |
| `derived_knowledge_state.created_at` / `created_by_kind` / `created_by_id` | 時点と主体種別は`NOT NULL`、主体IDは条件付き。主体種別は `user`、`codex`、`system` | 初期の状態内容を誰がいつ決めたか。`user`ではIDを禁止し、`codex`・`system`ではIDを必須にする。保存処理だけを行った主体を内容決定者へ置き換えないため |
| `derived_knowledge_state.updated_at` / `updated_by_kind` / `updated_by_id` | 時点と主体種別は`NOT NULL`、主体IDは条件付き。`updated_at >= created_at` | 最後に状態内容を決めた主体と時点。新規状態で変更がなければ作成監査値と同値にし、`store_commit.committed_at`（確定境界）や`calculated_at`（計算実行）と異なる役割を失わないため |
| `derived_knowledge_state.store_commit_id` | `NOT NULL REFERENCES store_commit` | 計算結果を確定した境界。索引側が正本時点を照合するため |

`state_evidence` は、`state_id REFERENCES derived_knowledge_state(state_id)` と `evidence_id REFERENCES evidence(evidence_id)` を個別の外部キーにし、`(state_id, evidence_id)` を複合主キーにする。このため、状態導出に実際に使った根拠だけを保存でき、存在しない状態・根拠への参照と、同じ状態への同じ根拠の重複を禁止できる。`status='no_evidence'` では0件、それ以外では1件以上を要求する規則は、複数行をまたぐため `L2-M1-S4` がトリガー又は確定前検証として実装する。

`state_aspect(state_id, aspect_text)` と `state_gap(state_id, gap_text)` は、それぞれ `state_id REFERENCES derived_knowledge_state` と本文を複合主キーにする。本文は前後空白を除いて1文字以上とし、順序は意味にしない。同じ状態内の重複、空文字列、別状態の行を許さない。前者は観測できた既知部分、後者は観測できない又は不足する部分であり、どちらも根拠ID・状態本文・`rationale`で代用しない。`status='partially_known'` では両表に1件以上を要求し、他の状態で0件を禁止しない。これにより、部分既知である理由と記事で補うべき箇所を、重複なく具体的に取得できる。

`state_publication` と `state_publication_revision` は、状態の値ではなく「現在状態を返してよいか」を表す選択情報である。`publication_status` は `available`（返せる）又は `recalculation_required`（再計算が必要）の二値だけを許す。前者では同じ命題の `current_state_id` を必須にし、後者ではNULLにする。`recalculation_required` を `no_evidence` 又は `uncertain` に読み替えることは禁止する。前者は根拠未観測、後者は根拠評価後の競合であり、再計算待ちとは別だからである。

`state_publication` は `assertion_id` を主キーとして一命題一行の業務上の選択を保持する。`state_publication_revision` は `publication_revision_id` を主キーとし、`assertion_id`、`publication_status`、`current_state_id`、`operation_id`、`changed_at`、`previous_publication_revision_id`、`is_current`、`store_commit_id` を持つ。`is_current=1` は命題ごとに一件だけ許し、`available` の `current_state_id` が同じ命題の状態を参照することは複合外部キー又は確定前検証で確認する。これにより、現在選択を上書きせず、再計算開始時に旧状態を返さないことと、完了後にどの状態を公開したかを両立する。

### 4.6 関係

`knowledge_relation` は命題又は概念のつながりの業務ID、`knowledge_relation_revision` はその一版である。関係は利用者本人がその関係を知っている証明ではない。

| 列 | 型・制約 | 意味と必要な理由 |
| --- | --- | --- |
| `relation_id` / `relation_revision_id` | `TEXT` の主キー | 関係とその版を別々に識別し、記録誤りの訂正履歴を残すため |
| `source_entity_kind` / `source_entity_id` | `assertion` 又は `concept` と対応ID | 関係の始点。向きを失うと前提・原因・置換を逆向きに辿るため |
| `target_entity_kind` / `target_entity_id` | `assertion` 又は `concept` と対応ID | 関係の終点。始点と同じ形式で、接続先の種別を明示するため |
| `relation_kind` | `TEXT CHECK ('related_to','prerequisite','causes','contributes_to','contradicts','supersedes')` | 関連、前提、原因、寄与、矛盾、置換の六つの意味だけを保存する。単なる関連という曖昧な値にしないため |
| `is_current` / `previous_revision_id` / `store_commit_id` | 他revisionと同じ規則 | 関係の記録訂正を追跡し、正本時点を確定するため |

関係の向きは `source_entity_*`、`target_entity_*` と `relation_kind` そのもので表すため、独立した `direction` 列は作らない。`prerequisite`、`causes`、`contributes_to`、`supersedes` は始点から終点へ読む有向関係であり、逆向き検索は既存行を終点側から検索する。`related_to` と `contradicts` は対称関係なので、`entity_kind` の固定順（`assertion`、`concept`）とID昇順で始点・終点を正規化し、左右を入れ替えた二重保存を禁止する。

始点・終点の外部キーは異種表へ向くため、SQLiteの単一外部キーでは表せない。`source_entity_kind` と `target_entity_kind` に応じて該当表のIDが存在すること、種別ごとの許可組合せ、自己参照、適用範囲・時点も含む同じ論理関係の重複、`prerequisite` と `supersedes` の全経路循環がないことを `L2-M1-S4` のトリガー又は確定前検証で強制する。`supersedes` は新命題→旧命題だけを許し、命題IDが同じなら禁止する。`causes` と `contributes_to` の循環は論理設計上許可されるため、この二種を一律に拒否しない。

関係自身に成立条件を持たせるため、`knowledge_relation_scope_facet`、`knowledge_relation_temporal`、`knowledge_relation_temporal_version` を `relation_revision_id` 配下へ置く。列と一意性は命題の適用範囲・時点表と同じであり、関係の条件を始点又は終点から自動複写しない。接続先と関係の適用条件は別の意味を持つためである。

## 5. 操作履歴と変更前後の完全保存

`KnowledgeUpdateOperation`、`KnowledgeUpdateAttempt`、`RecordChange`、`StatePublicationChange` は、業務データではなく更新を追跡する正本である。更新要求の再試行、競合、拒否と、実際に確定した変更を区別するために必要である。

| 表 | 主な列と制約 | 必要な理由 |
| --- | --- | --- |
| `knowledge_update_operation` | `operation_id`、`operation_kind`（`create`、`attach-evidence`、`revise`、`supersede`）、`request_fingerprint`、`decided_by_*`、`decision_rationale`、`requested_at`、`result`（`committed`、`no_change`、`rejected`）、`completed_at` | 同じ操作IDの再試行で二重更新しないため。`request_fingerprint` は意味要素を正規化して得る照合値であり、同じIDへの別要求を検出する |
| `knowledge_update_attempt` | `attempt_id`、`operation_id`、`attempt_sequence`、`request_fingerprint`、`executed_by_*`、`received_at`、`completed_at`、`result`（`committed`、`replayed`、`recalculation_required`、`rejected`、`failed`）、`diagnostic`、`resume_cursor`、`terminal_operation_id` | 一つの論理要求の各保存試行を残すため。操作全体の終端結果と、一回の途中到達・失敗を混同しない |
| `record_change` | `change_id`、`originating_attempt_id`、対象種別・ID、`change_kind`、`before_revision_id`、`before_snapshot_json`、`after_revision_id`、`after_snapshot_json` | 変更前後の完全値を確定後に上書きしないため。`created`では前値を禁止し、`revised`では前値を必須にする |
| `state_publication_change` | `publication_change_id`、`originating_attempt_id`、命題ID、前後の公開選択revisionと完全snapshot | 計算結果そのものと「現在値として返す選択」の切替を別々に監査するため |

`record_change.before_snapshot_json` と `after_snapshot_json` は、対象revisionの列と子表の値を含む完全な正規JSONでなければならない。読取り高速化のために操作表へ同じJSONを複写しても、変更表を正本とし、不一致は補正せず整合性違反として拒否する。

論理モデルで順序付き一覧又は対象一覧である属性を、JSON配列だけへ閉じ込めない。次の中間表を置き、外部キー、一意制約、順序を検査できる形にする。これらはSnapshotの代替ではなく、操作経路を辿るための正規化された参照である。

| 表 | 主キー・主な列 | 制約と必要な理由 |
| --- | --- | --- |
| `operation_expected_revision` | `(operation_id, target_kind, target_id)`、`expected_revision_id` | 操作が比較した現在版を一つずつ保存する。対象をJSONだけにすると、競合検出の対象と版を参照整合性で確認できないため |
| `operation_affected_assertion` | `(operation_id, assertion_id)` | 状態再計算の影響を受ける命題を重複なく保存する。再計算対象の漏れを検査するため |
| `operation_record_change` | `(operation_id, sequence_no)`、`change_id UNIQUE` | 同じ操作に属する全`RecordChange`を成立順に一度だけ集約する。別操作・別試行の変更を混ぜないため |
| `operation_publication_change` | `(operation_id, sequence_no)`、`publication_change_id UNIQUE` | 同じ操作に属する全`StatePublicationChange`を成立順に一度だけ集約する。公開停止・再開の履歴を欠かさないため |
| `operation_result_reference` | `(operation_id, entity_kind, entity_id)` | `committed`又は`no_change`の業務結果として返す命題・根拠・状態・関係を重複なく保存する。変更事実のIDやSnapshotを返却対象へ混ぜないため |
| `attempt_record_change` | `(attempt_id, sequence_no)`、`change_id UNIQUE` | 変更を実際に確定した一試行を固定する。部分成功を後続の再試行へ誤って帰属させないため |
| `attempt_publication_change` | `(attempt_id, sequence_no)`、`publication_change_id UNIQUE` | 公開選択を切り替えた一試行を固定する。再計算待ちの安全な途中状態を追跡するため |

`record_change` と `state_publication_change` には、各一覧から導ける `operation_id` も `NOT NULL` で持たせ、対応する`originating_attempt_id`の操作IDと必ず一致させる。`knowledge_update_attempt` には、論理モデルの `request_fingerprint`、`decision_rationale`、`diagnostic`、`resume_cursor`、`terminal_operation_id` を、その存在条件とともに列として保存する。`failure_reason` のような一語では、拒否規則、対象、再開可否を説明できないためである。

`operation_id` が構文上有効で照合値まで計算できた後は、同じID・同じ照合値で確定済みなら `replayed` として既存結果を返し、新しい変更を作らない。同じID・異なる照合値は衝突として拒否する。意味上の重複を判断するのはCodexであり、SQLiteは照合値、一意制約、想定revisionの一致だけを機械検査する。

## 6. 参照整合性・一意性・実装時の検査

次のうち通常の外部キー・一意制約で表せる規則はDDLに置く。複数表の件数、異種参照、全経路循環の規則は、トリガー又は同一トランザクション内の確定前検証として実装する。アプリケーションの記憶だけに委ねると、別実装経路が規則を破れるためである。

| 規則 | 強制方法 | 必要な理由 |
| --- | --- | --- |
| 根拠、状態、revision、変更履歴が存在する親を参照する | 外部キー、削除禁止 | 孤立行から根拠や履歴を失わないため |
| 一つの業務IDに現在revisionがちょうど一件 | 部分一意索引＋確定前検証 | 通常取得と次更新がどの版を使うか一意にするため |
| 根拠の出典作成者は `user` | `CHECK` と出典表の制約 | Codex・自動処理だけの情報を利用者知識の根拠にしないため |
| `available` の公開選択は同じ命題の状態だけを指す | 複合外部キー又はトリガー | 別命題又は古い状態を現在値として返さないため |
| `no_evidence` の根拠数は0、他状態は1以上 | 確定前検証 | 未観測と根拠評価済みの状態を混同しないため |
| 状態の `observed_at` は `calculated_at` と同値、適用期間は開始より終了が後 | `CHECK` と確定前検証 | 計算結果の観測時点を計算時点から乖離させず、終了境界を含まない有効期間を一貫して判定するため |
| 状態の適用製品版、既知部分、不足部分を重複・空値なく親状態へ結ぶ | 子表の複合主キー、外部キー、空白禁止 `CHECK` | 一つの状態に同じ製品版又は同じ部分説明を重ねず、孤立した一覧行を作らないため |
| `partially_known` の状態は既知部分と不足部分を各1件以上持つ | 確定前検証 | 「一部だけを知っている」という分類の内訳を、どちらか片方だけにして曖昧にしないため |
| 状態の監査主体・時点は有効な `ActorReference` と順序を満たす | `CHECK` と確定前検証 | 内容決定者を保存実行者と混同せず、作成後に更新されたという時間順を検証するため |
| `reported_unknown` は明示的不知根拠を使う | 確定前検証 | 利用者が知らないことを推測だけで保存しないため |
| revision鎖、`supersedes`関係は循環しない | 再帰問合せを使う確定前検証 | 新旧の前後関係を無限に辿る壊れた履歴を作らないため |
| `expected_revision` が現在値と一致する | 更新トランザクションの条件 | 後から到着した更新が先行更新を黙って上書きしないため |

## 7. スキーマ変更（Schema Evolution）の規則

ここでいうスキーマ変更は、SQLite内の表・列・制約を変え、既存の正本を新しいプログラムでも安全に読めるようにする作業である。利用者の根拠や履歴を失わず、異なる版のデータベースを同じ構造だと誤認しないために必要である。

1. 物理スキーマの変更は、順序番号と内容が固定された新しいSQL移行ファイルだけで追加する。既適用ファイルを編集・削除・並べ替えない。
2. 移行ファイルの識別子、内容照合値、適用時点、適用結果を保存する実行履歴の表とRunnerは `L2-M1-S5` が所有する。本書は、その履歴が物理スキーマの適用状態を検査可能にしなければならないという要求だけを渡す。
3. 後方互換な追加（任意列、既定値を持つ列、新表、読取り専用索引）は、旧行の意味を変えない移行として追加できる。必須列、許可値の縮小、分割・統合、既存値の意味変更は、旧値の変換・検証・復元方法を同じ移行設計で明記しなければならない。
4. 正本の列・表を廃止する場合も、直ちに削除しない。まず新旧の読取り経路、変換済み件数、参照整合性、履歴snapshotの復元を検証し、後方互換の終了を別の承認済み設計で決める。履歴JSONは当時の列名を含み得るため、表削除だけで過去の意味を失わないことを確認する必要がある。
5. 派生索引の再構築に影響する正本表・列・正規化規則が変わる場合は、正本スキーマ版と再構築対象の `store_commit.commit_sequence` を取り出せるようにする。これは再構築へ含める正本変更境界の単調増加する順序番号であり、その番号以下で確定した正本を入力にする境界を指す。表示順や時計ではなく、§4.1で定義した一意の確定順を使うためである。索引定義、モデル、索引世代との対応と互換性判定は `L2-M2-S4` が所有する。
6. 移行は、失敗しても未適用又は完全適用のどちらかになるトランザクション境界で実行する。SQLiteで一部のDDLが持つ制約への具体的対応はRunner設計で検証し、途中で適用済みと記録してはならない。

## 8. 後続Taskへの受渡し

| 後続Task | 渡す物理契約 | 成果物を渡せるか | 当該Taskが今すぐ着手できるか |
| --- | --- | --- | --- |
| `L2-M1-S3`（#41） | `store_commit` を正本変更境界とし、操作・試行・変更前後snapshot・現在revision・公開選択を同一トランザクションで整合させる必要がある | はい。この文書を読取り入力として渡せる | いいえ。Task Map上、この文書の統合と、同Taskが持つ他の直接依存が必要である |
| `L2-M1-S4`（#42） | 本書の表、列、外部キー、一意制約、トリガー又は確定前検証の必須規則をBaseline SQLへ具体化する | はい | いいえ。`L2 Design Freeze Gate` と本書の統合が着手条件である |
| `L2-M2-S2`（#50） | 検索索引は正本でなく、`store_commit.commit_sequence` 以下で確定した命題・根拠・関係・時点・適用範囲から再構築する。この列は索引入力を固定する正本変更境界の順序であり、時刻や索引世代そのものではない | はい | いいえ。L2-M1-S1〜S3など、同Task自身の全依存が未充足である |

「渡せる」はこの設計書を後続が入力として読めること、「今すぐ着手できる」は後続Taskの全依存・統合・Gateを満たすかを表す別の判定である。未Commitの本書を、後続の確定済み入力として扱ってはならない。

## 9. 今回の確認範囲と残る境界

- 本書は、論理モデルが定める根拠、命題、関係、計算済み状態、更新履歴をSQLiteの正本へ写し、物理削除・派生索引への依存を禁止する。
- 初期SQL、SQL構文のSQLite版差異、Migration適用履歴、接続設定、ID発行、トランザクション実装、Repositoryの操作名、検索索引製品と索引世代管理は、それぞれ後続Ownerが決める。未決定の漏れではなく、担当境界を保つために本書が固定しない事項である。
- GitHub APIへ接続できず、Issue #1、#40、#41、#42、#50の最新本文はこの作業時点で再取得できなかった。Task ID、依存、Owner、直接後続は、固定Planning snapshot、現在のTask Map、接続台帳、`docs/github-issue-map.md`、先行ADRから確認した。この接続不能は、本書に採用したSQLite物理設計の根拠を変更しないが、Issue本文固有の受入条件・指定決定記録は独立レビュー時に再照合が必要である。
