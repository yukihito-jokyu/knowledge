# DEC-FEAT-015: URL評価の回答とKnowledge Updateの実行順

> **状態:** `decided`  
> **Decision Level:** L3 — Codex workflowの実行・利用者体験契約  
> **対象Feature:** FEAT-004

## 判断が必要なこと

Reading Value Assessmentの内容を完成させた後、Knowledge AcquisitionとKnowledge Updateをどの順序で実行し、いつ会話へ回答を返すかを決める。

ここでいう**Knowledge Update**は、URL評価中のユーザー由来の説明・コード・訂正だけを候補として既存Knowledgeと照合し、必要なときだけKnowledge CLIへ保存を依頼する後続処理である。記事本文やCodexの回答を保存する処理ではない。

## 現在の状態と事実

- DEC-FEAT-013の当初表現は、URL評価の回答完了をKnowledge Acquisitionの開始条件としていた。今回のDecisionで「Assessment本文の完成後・会話への返却前」へ明確化する。
- 現行の`skills/reading-value/SKILL.md`は、Reading Value Assessmentを会話へ返すと実行を終了する。回答を返した後に同じ実行を継続する仕組みはない。
- job queue、scheduler、callback、保存済み会話を再開する仕組みは、既存の設計・実装配置にない。
- Knowledge CLIは既存の同期CLIとして利用できる。CLI更新失敗は回答内容を変更しないことがFEAT-004の要件である。

したがって、「回答を返した後に自動実行する」という結果を保つには、新しい後続実行機構が必要になる。一方で、回答の内容を完成させた後・会話へ返す前に更新を実行するなら、既存の一回のWorkflow実行だけで完結する。

**重要な影響:** Reading Value Assessmentを完成させた後に更新しても、そのAssessmentは更新前のKnowledge Storeに基づく。更新した知識は次回以降のURL評価には使えるが、今回の推薦・理由には反映されない。今回の回答にも反映するには、Knowledge AcquisitionとKnowledge UpdateをKnowledge Searchより前に完了させるか、更新後にAssessmentを再実行する必要がある。

## 選択肢

### A. 回答を会話へ返す前に、後続処理を同じWorkflowで実行する（推奨）

Reading Value Assessmentの内容を完成させた後、Knowledge AcquisitionとKnowledge Updateを実行し、その結果にかかわらず同じ回答本文を会話へ返す。

- 利用者への影響: Knowledge CLI処理の分だけ回答表示が遅れる可能性がある。更新失敗でも回答本文は変わらない。
- 実装への影響: `skills/reading-value/`から新設Workflowを同期で呼ぶ。job queue・再開用保存・公開APIは不要。
- 後続作業への影響: 一回のURL評価実行内でEpisode、Candidate、Update Resultを追跡できる。失敗・中断は既存の同期的な結果契約へ収められる。

### B. 回答を会話へ返した後、別の後続実行として自動起動する

回答表示の完了を契機に、別実行でKnowledge AcquisitionとKnowledge Updateを開始する。

- 利用者への影響: 回答表示は遅れない。後続更新が失敗・中断しても、別途その結果を確認する方法が必要になる。
- 実装への影響: 後続実行を起動・追跡するlifecycle、job queue、callback、または同等の仕組みと、Episodeの引渡し・失敗時の扱いを新たに定義する必要がある。
- 後続作業への影響: FEAT-004の範囲を超える実行基盤・運用・保存／再開契約のDecisionと設計変更が必要になる。

### C. URL評価の開始時にKnowledge AcquisitionとKnowledge Updateを実行してから評価する

URLとともに受け取った当該エピソード内のユーザー由来寄与を先に候補化・更新し、その更新後のKnowledge StoreでKnowledge SearchとReading Value Assessmentを実行する。

- 利用者への影響: 今回のAssessmentは更新済みの知識を反映する。更新失敗・中断時は、評価を続けるか停止するかの失敗方針を別途固定する必要がある。
- 実装への影響: DEC-FEAT-013の「回答完了後に取り込む」開始条件を変更する。URL評価の入力境界、失敗時の評価可否、候補に使えるユーザー寄与の範囲を再設計する必要がある。
- 後続作業への影響: 現行のDEC-FEAT-013を置き換えるL3 Decisionと、FEAT-003 / FEAT-004間の順序・失敗契約の再レビューが必要になる。

## 推奨

**Aを推奨する。** Knowledge Acquisitionは評価中にユーザー自身が示した知識候補を抽出し、Knowledge Updateはそれを次回以降の評価に使えるKnowledge Storeへ反映する処理である。今回のAssessmentを再計算する処理ではないため、同期処理後にも同じ回答本文を返す。Bは未承認の実行基盤を増やし、Cはユーザー由来寄与が評価の終了前には確定しないという入力境界を崩す。

## 決定

2026-08-15に人間が **A** を確定した。Reading Value Assessmentの**内容を完成**させた後、まだ会話へ返す前に、同じURL評価Workflow内でKnowledge Acquisition、続けてKnowledge Updateを同期実行する。その成功、skip、失敗、中断、部分適用のいずれでも、完成済みのReading Value Assessment本文は変更せず会話へ返す。

この順序は、今回の読書推奨を更新結果で変えるためではない。URL評価中のユーザー由来寄与を、次回以降のKnowledge SearchとReading Value Assessmentに利用できるよう蓄積するためである。更新が失敗・中断・部分適用なら、Knowledge Storeが完全には反映されないことをUpdate Resultへ残すが、回答本文の成功条件にはしない。

## 承認後に更新する資料

- `docs/features/FEAT-004/design.md`: 起動条件、正常・中断flow、Implementation Deliverables / Placement、Assumptions。
- `docs/features/FEAT-004/design/workflow-contracts.md`: EpisodeとUpdate Resultの開始・受渡し・中断規約。
- `docs/features/FEAT-004/design-review.md`: 差し戻し事項の再レビュー。
- Cを選ぶ場合は、`docs/features/FEAT-004/decisions/DEC-FEAT-013.md`、`docs/features/FEAT-003/design.md`、関連するhandoffを含む開始順・失敗契約の再設計。
