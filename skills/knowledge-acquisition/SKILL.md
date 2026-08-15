---
name: knowledge-acquisition
description: 完成したURL評価Episodeから、ユーザー由来の技術的EvidenceだけをCandidate Knowledgeへ抽出する。Knowledge Store、Update Decision、CLI、DBを変更しない。
---

# Knowledge Acquisition

完成した一回のURL評価Episodeから、ユーザー自身が示した技術的Evidenceだけを、Knowledge Updateへ渡す一時的なCandidate Knowledge成果物へ変換する。候補抽出は保存や知識保有の判定ではない。

## 開始と入力契約

開始前に同梱の[成果物契約](references/artifact-contract.md)と[検証契約](references/verification.md)を読む。実行時の入力は一つの`URL Evaluation Episode`だけとし、他の会話、保存済み履歴、別作業、記事取得結果、Codex応答を読み込まない。

このSkillは、実行場所や`workflow-contracts.md`の物理配置を前提にしない。呼出側から渡された正規契約と同梱の成果物契約を照合し、契約の不一致や正規契約の欠落があればCandidateを生成せず停止する。リポジトリ内の特定パスを実行時に探索しない。

入力には次を必須とする。

- `episode_id`: 不透明な非空文字列。
- `article_url`: 評価対象のURL文字列。
- `completed_at`: Assessment本文の内容を完成させた時刻を示すRFC 3339 UTC文字列。会話へ返した時刻ではない。
- `user_contributions`: 当該Episode内のユーザー寄与のリスト。空を許すが、各要素に`ordinal`、`source_text`、`source_type`、`observed_at`を持つ。

呼出側は、`completed_at`が存在する場合に限り、Reading Value Assessmentの内容が完成済みであることを保証する。この保証がない、評価が失敗・中断した、または必須入力が欠けるEpisodeは抽出対象にしない。

## 手順

1. 入力の必須field、型、非空条件、RFC 3339 UTC、`ordinal`の正数・重複なし・入力順を検証する。不正な重複や順序を修正・推測せず、入力エラーとして停止する。
2. `completed_at`とAssessment完成済みの入口条件を確認する。未完成、失敗、中断、再開されたEpisodeは候補を作らず停止する。
3. `user_contributions`だけを`ordinal`順に走査する。ユーザーが発した説明、推論、コード、訂正、自己申告、概念認識、理由付き技術判断だけを、独立評価可能な命題へ抽出する。
4. 一つの寄与に複数の独立命題、または複数のEvidence種類がある場合は、命題・種類ごとにCandidateを分ける。同じ寄与の候補は先頭から`candidate_ordinal`を1ずつ付ける。
5. 各Candidateに、意味を変えない完全な`evidence_raw_text`、必要最小限の`source_excerpt`、原文から明示された`scope`と`temporal`、抽出理由、`proposed_assertion`を設定する。原文の要約、黙った正規化、記事本文やAI本文の複製をしない。
6. `evidence_kind`と`strength`を[成果物契約](references/artifact-contract.md)の表で分類する。理由のない技術判断、質問、概念以外の単なる言及はCandidateにしない。
7. `search_queries`は`proposed_assertion`をそのまま先頭に置き、続けてEvidence原文に明示されたConcept・Alias・Identifier文字列だけを出現順に置く。同一文字列は最初の一件だけにし、推測した同義語、翻訳語、Scope由来語を加えない。
8. Candidateを`source_ordinal`昇順、同じ寄与内では`candidate_ordinal`昇順に並べ、成果物契約どおりの一時Markdown成果物として呼出側へ返す。`candidate_id`は追跡用の不透明IDであり、処理順や意味を表さない。

## 除外と停止

次の内容からはCandidateを作らない。入力に混在していた場合も、ユーザー由来であることを確認できる部分だけを扱い、除外部分を`evidence_raw_text`や`source_excerpt`へコピーしない。

- 通常会話、URL評価と無関係な別作業、保存済み履歴。
- 失敗・中断・未完成のURL評価。
- Codexの応答、記事本文、閲覧内容、記事の評価・要約・検索結果。
- 質問だけの発話、理由のない結論・方針表明、単なる非技術的な感想。

入力field不足、RFC 3339 UTC違反、空の`source_text`、重複`ordinal`、順序違反、または完了条件不成立では、`invalid_input`または`assessment_incomplete`の停止結果だけを返す。候補一覧やUpdate Decisionを捏造しない。同じ入力を自動再試行・再開しない。

完成済みEpisodeで許可されたEvidenceが一件もない場合は失敗ではなく、`Candidates: なし`の空Candidate成果物を返す。この場合もUpdate Decisionは作らず、Knowledge Updateへ空一覧を渡す。

## 出力形式と境界

正常時は、[成果物契約](references/artifact-contract.md)のMarkdown形式で`Candidate Knowledge`を一件返す。候補ゼロでも同じ成果物の空一覧とする。停止時は停止理由を含むMarkdownだけを返し、CandidateやUpdate Decisionを含めない。

このSkillは次を行わない。

- Candidateの保存、Knowledge Storeの読み書き、Update Decision / Update Resultの作成。
- `search-text`その他のKnowledge CLI操作、公開CLI/JSONの変更、SQLite、migration、ledgerの追加。
- `skills/knowledge-update/`、Reading Value、記事本文、Codex応答の保存。
- 同一`episode_id`の自動再開、自動再実行、補償削除。

## 検証チェックリスト

- 完成済みURL評価Episodeだけを受け、未完成・失敗・中断・無関係Episodeを停止または除外した。
- 必須field、寄与の順序、候補の順序、完全なEvidence原文、観測時刻を追跡できる。
- `user_explanation`、`user_reasoning`、`user_code`、`correction`、`self_report`、`concept_recognition`の分類と`strong` / `moderate` / `weak`が契約どおりである。
- 複合寄与を独立Candidateへ分け、質問・AI・記事・理由なし判断を候補化していない。
- `search_queries`の先頭、明示語の出現順、完全一致の重複除去、推測語の不在を確認した。
- 候補ゼロが有効な空結果であり、Update Decision、CLI、DB、migration、保存がないことを確認した。
- 詳細な観測手順と期待結果は[検証契約](references/verification.md)に従った。
