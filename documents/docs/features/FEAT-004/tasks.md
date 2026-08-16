# FEAT-004 実装タスク

> **前提:** 利用者の「タスク分割に進んでください」（2026-08-15）を、独立レビュー`pass`後の詳細設計の明示承認およびTask分解開始指示として記録した。

## Design Readiness Audit

| 監査項目 | 結果 | 根拠 |
| --- | --- | --- |
| 利用者による詳細設計承認 | pass | `.ai/workflow/state.yaml` の `features.FEAT-004.human_design_review: approved` |
| 独立設計レビュー | pass | [design-review.md](design-review.md) の総合判定 `pass` |
| URL評価への組込み順序と回答不変 | pass | [DEC-FEAT-015](decisions/DEC-FEAT-015.md)、[design.md](design.md)「Behavioral Scenarios」 |
| Candidate / Update Result / 検索操作列 | pass | [workflow-contracts.md](design/workflow-contracts.md)「Candidate Knowledge」「Update Decision」「Update Result」 |
| 既存CLI・永続化・公開JSON境界 | pass | [design.md](design.md)「Scope / Out of Scope」「Contract Completeness」、FEAT-001 Command Catalog |
| 新規DB・migration・公開CLI/API | not_applicable | 新設・変更せず、既存Knowledge CLIをWorkflowから消費する。 |
| Fixture・受入観点 | pass | [design.md](design.md)「Acceptance / Test Design」 |

**結論:** pass。実装者に、取り込み対象、検索語列、Evidence分類、更新順、部分適用、回答返却、公開境界を再決定させる未確定事項はない。

## 固定済みの実装対象と配置

実装本体は親リポジトリの`skills/`配下に置く。`documents/docs/features/FEAT-004/`は正規の設計・handoff資料であり、実行するWorkflow本体ではない。

| 配置先（親リポジトリ基準） | 作成・変更 | 責務 |
| --- | --- | --- |
| `skills/reading-value/SKILL.md` | 変更 | Assessment本文を完成後・返却前にAcquisitionとUpdateを同期実行し、結果を本文へ反映せず返す。 |
| `skills/knowledge-acquisition/SKILL.md` | 作成 | URL評価Episodeから許可されたユーザー寄与だけをCandidate Knowledgeへ抽出する。 |
| `skills/knowledge-acquisition/references/artifact-contract.md` | 作成 | Episode、Candidate、Evidence分類、検索入力の成果物契約を定義する。 |
| `skills/knowledge-acquisition/references/verification.md` | 作成 | 候補化・除外・Evidence強度の検証契約を定義する。 |
| `skills/knowledge-update/SKILL.md` | 作成 | 候補検索、操作選択、既存CLI実行、停止・部分適用・結果不明を制御する。 |
| `skills/knowledge-update/references/artifact-contract.md` | 作成 | Update DecisionとUpdate Resultの成果物契約を定義する。 |
| `skills/knowledge-update/references/cli-operations.md` | 作成 | 消費する既存CLI operationとJSON・error・exit code契約を参照可能にする。 |
| `skills/knowledge-update/references/verification.md` | 作成 | 更新操作、競合、protocol error、中断の検証契約を定義する。 |

`cmd/knowledge/`、`internal/`、SQLite migration、Knowledge CLIの公開operation・JSONは変更しない。

## タスク一覧

| ID | タイトル | 論理領域 | 依存 |
| --- | --- | --- | --- |
| TASK-004-01 | URL評価EpisodeからCandidate Knowledgeを取得する | Codex workflow: acquisition | なし |
| TASK-004-02 | Candidateを既存Knowledgeへ照合・更新する | Codex workflow: update / CLI integration | TASK-004-01、FEAT-001 implementation handoff |
| TASK-004-03 | Reading Value Workflowへ同期更新を統合する | Codex workflow: orchestration | TASK-004-01、TASK-004-02、FEAT-003 implementation handoff |
| TASK-004-04 | 知識取得・更新Workflowの受入境界を検証する | verification | TASK-004-01〜TASK-004-03 |

## TASK-004-01: URL評価EpisodeからCandidate Knowledgeを取得する

- **目的:** 完成したURL評価Episode内の許可されたユーザー寄与だけを、追跡・再評価・検索に使えるCandidate Knowledgeへ変換する。
- **関連要件:** REQ-015、REQ-017、REQ-019、BR-004〜006、BR-010、NFR-001、NFR-003、NFR-006、DEC-FEAT-013、DEC-FEAT-015。
- **論理領域:** Codex workflow: acquisition。
- **作業内容:** `knowledge-acquisition` Workflowと成果物・検証参照資料を作成する。Episode、ユーザー寄与、Candidateの必須項目、発話順、Evidence原文、観測時刻、種類・強度の導出、`search_queries`を契約どおりに生成する。質問、AI説明、記事本文・閲覧・評価・要約、理由のない技術判断はCandidateにしない。
- **受入条件:**
  - URL評価のAssessment本文完成後・返却前に渡された当該Episodeだけを入力とし、失敗・中断した評価、通常会話、別作業、保存済み履歴は読まない。
  - Candidateは`episode_id`、順序、完全なEvidence原文、`observed_at`、Evidence kind、派生strength、正規化Assertion、Scope/Temporal、抽出理由、検索入力を保持する。
  - 説明・推論・コード・訂正・理由を伴う技術判断はstrong、自己申告はmoderate、概念認識はweakとして導出し、複合寄与は根拠種類ごとに分割する。
  - 除外入力にはCandidateもUpdate Decisionも作らず、許可される根拠がゼロなら後続へ空Candidate一覧を渡せる。
  - 各Candidateの検索入力は候補Assertionを先頭に、訂正Evidenceで引用された完全な旧命題（ある場合のみ）、原文に明示されたConcept・Alias・Identifierを出現順・重複なしで続け、Scopeや推測した同義語を加えない。
- **依存:** なし。
- **対象外／注記:** Knowledge Storeの照合・更新と、回答返却の組込みは後続Taskで扱う。Knowledge CLIの変更はしない。

## TASK-004-02: Candidateを既存Knowledgeへ照合・更新する

- **目的:** Candidateごとに既存Knowledgeを決められた字句検索列で確認し、履歴を失わない既存CLI操作と結果を一意に記録する。
- **関連要件:** REQ-016、REQ-018〜019、BR-004、BR-006、BR-008、BR-010、NFR-003〜004、NFR-006、DEC-FEAT-014〜015。
- **論理領域:** Codex workflow: update / CLI integration。
- **作業内容:** `knowledge-update` Workflowと成果物・CLI参照・検証資料を作成する。Candidateの`search_queries`を一語ずつ既存`search-text`へ渡し、Assertion IDを初出順に集約して必要時だけ取得する。Codexが`create`、`attach-evidence`、`revise`、`supersede`、`skip`を選び、既存CLIのJSON、error、exit codeを消費してDecisionとUpdate Resultを作る。
- **受入条件:**
  - 各検索入力を一回ずつ順に実行し、空結果でも次へ進み、返却Assertion IDを初出順で重複除去する。検索失敗・中断時は残りqueryと後続Candidateを実行しない。
  - `create`、`attach-evidence`、`revise`、`supersede`、`skip`の選択条件と、対象・置換ID、検索根拠、操作結果、失敗理由をCandidateごとに追跡できる。
  - `revise`後のEvidence追加、または`create`後の`supersede`が失敗・中断したとき、成功済みの履歴を削除せず、契約どおり`partially_applied`または`failed/outcome_unknown`を残す。
  - 最初の`failed`、`canceled`、`partially_applied`以降のCandidateは`not_started`として全件記録し、CLIを呼ばない。候補ゼロは空Decision一覧・`completed`となる。
  - protocol error、CLI error、exit 130、既存Evidence重複・conflictの扱いが、既存CLI公開契約およびUpdate Result契約と一致する。Episodeの自動再開・自動再実行はしない。
- **依存:** TASK-004-01、FEAT-001 implementation handoff（既存Knowledge CLI）。
- **対象外／注記:** 意味検索、新CLI operation、SQLite schema・migration、補償削除、再開ledgerは追加しない。

## TASK-004-03: Reading Value Workflowへ同期更新を統合する

- **目的:** URL評価の回答本文を完成した後・返却前に、AcquisitionとUpdateを同一Workflowで実行し、更新結果によらず同じ本文を返す。
- **関連要件:** REQ-019、BR-010、NFR-003、NFR-006、DEC-FEAT-013、DEC-FEAT-015。
- **論理領域:** Codex workflow: orchestration。
- **作業内容:** 既存Reading Value Workflowを、URL評価Episodeの生成、Assessment本文完成後の同期的なAcquisition→Updateの呼出し、Update Resultの一時保持、同一本文の返却へ更新する。URL評価が本文完成前に失敗・中断した場合はAcquisition・Updateを開始しない。
- **受入条件:**
  - URL評価一回ごとに不透明な`episode_id`を割り当て、各ユーザー寄与の観測時刻をAcquisitionへ渡す。
  - Assessment本文を完成した後、会話へ返す前にAcquisition、続けてUpdateを同期実行する。
  - Updateがcompleted、skipped、failed、canceled、partially applied、outcome unknownのいずれでも、完成済みのAssessment本文・推奨・理由を変更せず返す。
  - Assessment本文を完成できないURL評価では、Candidate抽出・CLI更新・Update Resultを作らない。
  - Update Resultは実行中の呼出側成果物に限り、Knowledge Store・別ledger・公開UI/APIへ保存しない。
- **依存:** TASK-004-01、TASK-004-02、FEAT-003 implementation handoff（Reading Value Workflow）。
- **対象外／注記:** 非同期job、callback、scheduler、保存済み会話の再開、回答後の別実行は追加しない。

## TASK-004-04: 知識取得・更新Workflowの受入境界を検証する

- **目的:** Evidenceの許可境界、検索列、更新履歴、失敗・中断、回答不変を、実装者以外も再現可能に確認する。
- **関連要件:** REQ-015〜019、BR-004〜006、BR-008、BR-010、NFR-001、NFR-003〜004、NFR-006、DEC-FEAT-013〜015。
- **論理領域:** verification。
- **作業内容:** AcquisitionとUpdateの検証参照資料および既存検証手順を用い、候補抽出、検索、CLI操作、統合順序の観測oracleを整備する。
- **受入条件:**
  - 説明・推論・コード・訂正・理由を伴う技術判断、自己申告、概念認識、複合寄与、質問のみ、AI説明のみを区別して、候補化・除外・強度を確認できる。
  - 検索入力順、検索結果IDの重複除去、空候補、複数Candidate、skip、途中停止後の`not_started`を確認できる。
  - create、attach、revise、supersede、重複Evidence、conflict、protocol error、exit 130、部分適用、結果不明がUpdate Decision / Resultへ正しく表れることを確認できる。
  - URL評価本文の完成後・返却前に更新が実行され、更新結果にかかわらず同じ本文が返り、本文未完成時は更新を開始しないことを確認できる。
  - Knowledge CLIの`cmd/knowledge/`、`internal/`、SQLite migration、公開CLI JSONを変更していないことを確認できる。
- **依存:** TASK-004-01、TASK-004-02、TASK-004-03。
- **対象外／注記:** Scenario A〜Jをまたぐ横断Fixture評価基盤はFEAT-005の所有とする。

## Dependency Notes

- TASK-004-01がCandidateの入力・順序・Evidence境界を固定し、TASK-004-02はその契約を消費する。
- TASK-004-03は両Workflowの完成成果物をReading Valueの返却前へ接続するため、TASK-004-01とTASK-004-02に依存する。
- TASK-004-04は取得・更新・統合の観測を行うため、三Taskに依存する。

## Implementation Readiness Review

| 項目 | 記録 |
| --- | --- |
| 判定 | **pass** |
| reviewer | `/root/feat004_implementation_readiness`（独立・読み取り専用） |
| 親リポジトリ俯瞰の範囲 | `AGENTS.md`、`skills/reading-value/`、`skills/knowledge-search/`、`cmd/knowledge/`、`internal/`、`test/integration/`、`Taskfile.yml`、FEAT-001 CLI契約、FEAT-003 handoff。 |
| 根拠 | 実装先は親rootの`skills/`に固定され、CLI・SQLite・公開JSONを変更しない。利用するEvidence kindと7 operationは既存契約・実装に存在する。Candidate契約→CLI更新→Reading Value統合→検証の順序、および成功・skip・失敗・中断・部分適用・結果不明・`not_started`・本文不変の受入oracleが揃う。Go共通検証とCLIプロセス境界テストの既存手順も確認済み。 |
