# Knowledge Acquisition 検証契約

この契約は、TASK-004-01のCandidate抽出を実装者以外が同じ入力観測点で確認するためのもの。対象はCodex workflowの一時Markdown成果物であり、Knowledge CLIの検索・更新品質やReading Valueの統合は対象外とする。

## 共通観測規則

- 各ケースは、入力Episode、Assessment完成状態、抽出結果、停止理由、後続handoffを別々に記録する。
- 正常ケースでは、各Candidateの`source_ordinal`、`candidate_ordinal`、`candidate_id`、`evidence_raw_text`、`observed_at`を入力寄与へ機械的に追跡する。
- `evidence_raw_text`と`source_excerpt`は入力のユーザー原文と文字単位で比較する。要約、`...`による省略、記事本文・AI本文の混入を合格にしない。
- Candidateの並びは`source_ordinal`昇順、同一寄与では`candidate_ordinal`昇順で確認し、`candidate_id`の文字列順をoracleにしない。
- 成果物全体を検索し、`Update Decision`、CLI operation、Knowledge Store、DB、migration、保存ledgerが存在しないことを確認する。

## AC-01: 完成Episodeだけを受け付ける

### V-001: 完成済みURL評価Episode

完成したAssessment、非空のopaque `episode_id`、URL、RFC 3339 UTCの`completed_at`、ordinal順のユーザー寄与を入力する。

期待結果:

- `Candidate Knowledge`成果物が返る。
- 対象は当該Episodeの`user_contributions`だけであり、他の会話、別作業、保存済み履歴を参照しない。
- CandidateからKnowledge Updateへ一時的にhandoffされ、保存・CLI・Update Decisionはない。

### V-002: 未完成・失敗・中断の境界

`completed_at`欠落、Assessment未完成、URL評価失敗、URL評価中断をそれぞれ入力する。

期待結果:

- `assessment_incomplete`または同等の停止理由で停止する。
- Candidate、空Candidate一覧、Update Decision、Update Resultを作らない。
- 同一`episode_id`を自動再開・再実行しない。

### V-003: 無関係Episodeの境界

通常会話、別作業、保存済み履歴、完成していない別URL評価を、完成済み対象Episodeと同じ入力源に置く。

期待結果:

- 対象Episode以外の入力からCandidateを作らない。
- 対象Episode側に許可Evidenceがなければ、正常な`Candidates: なし`を返す。
- 無関係入力を理由にCandidateやUpdate Decisionを作らない。

## AC-02: Candidateの追跡可能性と原文

### V-004: 必須fieldと順序

複数の寄与を異なる`ordinal`と`observed_at`で入力し、各寄与から複数Candidateを抽出する。

期待結果:

- 各Candidateが契約の14 field（`episode_id`、`candidate_id`、`source_ordinal`、`candidate_ordinal`、`source_excerpt`、`evidence_raw_text`、`evidence_kind`、`strength`、`observed_at`、`proposed_assertion`、`search_queries`、`scope`、`temporal`、`extraction_rationale`）を持つ。
- `source_ordinal`が入力寄与へ、`observed_at`が同じ寄与へ追跡できる。
- `candidate_id`はEpisode内で一意だが、意味や処理順を示さないopaque IDである。
- `scope`と`temporal`は原文で明示された値だけを持ち、未指定ならそれぞれ空list・`null`である。

### V-005: 原文保持とプライバシー

長いユーザー説明、ユーザー提示コード、記事本文の引用、Codex応答の説明を同じ入力に含める。

期待結果:

- 許可されたCandidateの`evidence_raw_text`は、候補を直接支えるユーザー原文を省略・要約せず保持する。
- `source_excerpt`は必要最小限のユーザー原文である。
- 記事本文・Codex応答・閲覧・評価・要約はCandidateへ複製されず除外される。

## AC-03: Evidence kind / strength

### V-006: 全種類と強度

次の入力を一件ずつ用意する: 技術説明、理由付き推論、ユーザー提示コード、明示的訂正、説明なし自己申告、概念名だけの認識。

期待結果:

| 入力 | 期待`evidence_kind` | 期待`strength` |
| --- | --- | --- |
| 技術説明 | `user_explanation` | `strong` |
| 理由付き推論・技術判断 | `user_reasoning` | `strong` |
| ユーザーコード | `user_code` | `strong` |
| 明示的訂正 | `correction` | `strong` |
| 自己申告 | `self_report` | `moderate` |
| 概念認識 | `concept_recognition` | `weak` |

理由のない技術判断はCandidateにならない。`strength`は保存状態ではなくCandidateでの派生表示である。

### V-007: 複合寄与の分割

一つのユーザー寄与に、説明と自己申告、コードと訂正、または独立した二つの命題を含める。

期待結果:

- 独立命題・根拠種類ごとにCandidateが分かれる。
- 各Candidateは一つの`evidence_kind`だけを持ち、`candidate_ordinal`が寄与内の出現順になる。
- 各`evidence_raw_text`はそのCandidateを直接支える完全な原文で、別命題の意味を混ぜない。

## AC-04: 除外境界

### V-008: 除外入力のみ

質問、Codex回答、記事本文、記事の閲覧内容、評価・要約、理由のない結論、通常会話、別作業、履歴だけを、完成済みEpisodeの寄与として与える。

期待結果:

- Candidateは0件。
- 完成済みEpisodeなので、失敗ではなく`Candidates: なし`の空成果物を返す。
- CandidateもUpdate Decisionも作らない。`skip` Decisionで代替しない。

### V-009: 入力構造エラー

必須field欠落、空`source_text`、非RFC 3339時刻、0以下のordinal、重複ordinal、降順ordinal、未知の`source_type`をそれぞれ入力する。

期待結果:

- `invalid_input`で停止し、入力を修正・並べ替え・推測補完しない。
- Candidate、Update Decision、保存、CLI呼出しはない。

## AC-05: Assertionとsearch_queries

### V-010: Assertion先頭と明示語の順序

例として、ユーザー原文に`SQLite`、`WAL`、`busy_timeout`を明示し、`database`のようなScope・推測語と、`SQLite`の再出現を含める。

期待結果:

- `search_queries[0]`が`proposed_assertion`と完全一致する。
- 訂正Evidenceに引用された旧命題がある場合だけ、`search_queries[1]`はその完全な引用文である。残りの後続queryは原文に明示されたConcept・Alias・Identifierだけで、出現順が保持される。
- 同じ文字列は大小文字・記号・空白を含めて完全一致で一度だけ現れる。
- 推測した同義語、翻訳語、Scope由来語、原文にない関連語はない。

### V-011: 複数Candidateの順序

複数寄与と一つの複合寄与を入力し、各CandidateのIDを意図的に順序と異なる文字列にする。

期待結果:

- 出力順は`source_ordinal`、次に`candidate_ordinal`だけで決まる。
- `candidate_id`の辞書順、生成時刻、Assertionの文字列順では並べ替えない。

## AC-06: 空結果と非変更境界

### V-012: Candidate zero

完成済みEpisodeに空の`user_contributions`、または除外入力だけを与える。

期待結果:

- 正常なCandidate成果物の`Candidates`は空で、後続handoffも空一覧を示す。
- Update Decision / Update Resultは作らない。
- 空結果を「ユーザーが知らない」「検索失敗」「skip Decision」と解釈しない。

### V-013: 永続化・公開境界

実行前後の作業tree、Knowledge Store、CLI呼出し記録、DB schema、migration一覧、公開JSON契約を比較する。

期待結果:

- 通常のTASK-004-01実装で変更できるのは、`skills/knowledge-acquisition/SKILL.md`、`skills/knowledge-acquisition/references/artifact-contract.md`、`skills/knowledge-acquisition/references/verification.md`の3ファイルだけである。
- DEC-FEAT-022の承認済み横断契約同期に限り、`skills/knowledge-update/references/artifact-contract.md`および根拠となるDecision・設計・handoffを同時に更新できる。
- `cmd/knowledge/`、`internal/`、SQLite、migration、公開CLI/JSON、`skills/knowledge-update/SKILL.md`、`skills/reading-value/`に変更がない。
- Candidate成果物は一時的で、Store・DB・ledger・設定へ保存されない。

## 合格判定

V-001〜V-013の期待結果を満たし、成果物のfield・原文・順序・分類・検索語を入力と照合できることをTASK-004-01の合格条件とする。ケースを実行していない場合は未検証と記録し、成功扱いにしない。
