# FEAT-004 要件再構成: 会話・作業からの知識 Evidence 更新

> **状態:** `completed` — エピソード取り込み境界と回答・更新の実行順はDEC-FEAT-013、DEC-FEAT-015で決定済み。独立設計レビューと人間承認を完了し、実装へ引き渡す。

## Goal

意味のある会話または作業エピソードから、ユーザー知識の根拠となり得る候補を抽出し、既存Knowledgeとの照合に基づいて、履歴を保つ更新を行えるようにする。

## Included Requirements

| ID | 要求 | Featureで扱う結果 |
| --- | --- | --- |
| REQ-015 | エピソード終了後に、ユーザー発言・推論・コード・技術的訂正・明示的自己申告・技術判断から候補を抽出する。 | Evidence候補とその原文・種類・候補Assertionを作る。 |
| REQ-016 | 候補を既存Knowledgeと照合し、create、attach-evidence、revise、supersede、更新しない、を判断する。 | Codexが更新操作を選び、Knowledge CLIへ決定論的操作を渡す。 |
| REQ-017 | 評価、AI説明、質問だけでは更新しない。 | 根拠外入力を候補化・保存しない。 |
| REQ-018 | 後日の訂正を強いEvidenceとして扱い、過去Evidenceを物理削除せず再評価できる。 | `correction` Evidenceと履歴保持操作を選ぶ。 |
| REQ-019 | CodexがKnowledge AcquisitionとKnowledge Updateを制御し、成果物を引き渡す。 | Acquisition、Update、CLI実行結果を分離して受け渡す。 |

## Related Business Rules

- BR-004: Evidenceを知識状態の正規根拠とし、再評価可能にする。
- BR-005: 露出、AI説明、閲覧、質問だけを知識獲得にしない。
- BR-006: Evidence強度を区別する。
- BR-008: 訂正・更新で物理削除せず履歴を保持する。
- BR-010: 意味判断はCodex、保存・検索・更新の決定論的実行はKnowledge CLIが担う。

## Scope

- Conversation / Task Episodeから許可されるユーザー由来の根拠を抽出する。
- 候補を既存Assertionと意味的に照合し、既存のKnowledge CLI操作を選ぶ。
- `create`、`attach-evidence`、`revise`、`supersede`の既決CLI契約を使って履歴を保つ。
- 更新しない判断と理由を成果物として残す。

## Out of Scope

- Knowledge CLIの公開operation、JSON、SQLite schema、migrationの変更。
- AI説明、記事閲覧、要約、質問だけをEvidenceとして保存すること。
- 記事Claim分解、Knowledge Assessment、Reading Valueの再設計。
- エピソード入力の保存方式、Codex連携方式、利用者UIの実装詳細。これらは下記Decisionの承認後に境界だけを定義する。

## Preconditions

- FEAT-001が提供するKnowledge CLIの `search-text`、`get`、`get-evidence`、`create`、`attach-evidence`、`revise`、`supersede` 契約を変更せず利用できる。
- ユーザーが技術記事URLを渡し、CodexがReading Value Assessmentの内容を完成させている（まだ会話へ返していない）。

## Acceptance Conditions

1. 許可されたユーザー由来の発言・推論・コード・訂正・自己申告・技術判断だけから候補を抽出できる。
2. 各候補について、根拠、Evidence kind、強度、候補Assertion、更新判断またはskip理由を追跡できる。
3. 字句検索で発見・確認できた既存Assertionを重複作成せず、既存AssertionへのEvidence追加・revision・置換を適切に選べる。
4. 訂正では過去Evidenceと旧Assertionを物理削除しない。
5. AI説明、記事の閲覧・評価・要約、質問のみは更新されない。
6. CLIのvalidation、not-found、conflict、storage failureでは、更新判断・実行結果を混同せず失敗を呼出側へ返す。

## 入力境界の決定

[DEC-FEAT-013](decisions/DEC-FEAT-013.md)および[DEC-FEAT-015](decisions/DEC-FEAT-015.md)により、CodexがURL評価のAssessment本文を完成させた後、会話へ返す前に、そのURL評価ワークフローだけを候補抽出の対象とする。任意会話・他作業・保存済み履歴は対象外である。
