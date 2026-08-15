# FEAT-004 詳細設計・独立レビュー

> **レビュー日:** 2026-08-15  
> **対象:** `docs/features/FEAT-004/requirements.md`、`design.md`、`design/workflow-contracts.md`、`decisions/DEC-FEAT-013.md`〜`DEC-FEAT-015.md`  
> **方式:** 設計者および過去レビューの結論を前提にせず、正規要件、承認済みDecision、初期設計、FEAT-001の既存CLI契約と現行設計を読み取り専用で照合した。

## 総合判定

**`pass`**

取り込み境界、Assessment本文の完成後・返却前という実行順、回答本文を更新結果で変えない規約、候補の抽出・順序・停止・結果記録、既存CLIの利用範囲は、正規要件と承認済みDecisionに整合する。実装者へ公開契約、保存形式、または業務規則の再設計を委ねる未解決事項はない。

この判定はhandoff可を意味しない。次のゲートは、利用者による詳細設計の明示承認の確認後のTask Breakdownである。

## 正規ソース台帳

| ソース | 確定性 | 適用範囲 |
| --- | --- | --- |
| `docs/requirements/requirements.md` REQ-015〜019 | confirmed | 候補抽出、更新判断、除外、訂正、Workflow制御 |
| `docs/requirements/business-rules.md` BR-004〜006、BR-008、BR-010 | confirmed | Evidence根拠・強度、履歴保全、責務分離 |
| `docs/requirements/nonfunctional-requirements.md` NFR-001、NFR-003、NFR-004、NFR-006 | confirmed | 追跡可能性、失敗の切分け、履歴、Fixture |
| `docs/design/architecture.md`、`docs/design/domain-boundaries.md` | initial design | Acquisition / Update / CLIの責務・Markdown成果物境界 |
| `docs/features/FEAT-001/design/command-catalog.md` および `commands/{search-text,get,get-evidence,create,attach-evidence,revise,supersede}.md` | review-ready public contract | 既存CLIのoperation、単一`--query`、JSON、error/exit code、mutation原子性 |
| DEC-FEAT-013 | decided / L3 / 人間決定 | URL評価Workflowに限定した自動取り込み境界 |
| DEC-FEAT-014 | decided / L2 | 複数CLI操作の部分適用、補償削除・自動再開なし |
| DEC-FEAT-015 | decided / L3 / 人間決定 | 同一Workflowで同期更新後、回答本文を変えず返却 |

## 根拠トレーサビリティ

| 契約・境界 | 設計上の扱い | 正規根拠 | 判定 | レビュー結果 |
| --- | --- | --- | --- | --- |
| 取り込み対象・開始時点 | Assessment本文完成後、返却前に当該URL評価Workflowだけを対象にする | REQ-015、REQ-019、DEC-FEAT-013、DEC-FEAT-015 | explicit | 整合 |
| 回答の返却順 | Acquisition→Updateを同一Workflowで同期実行し、結果によらず完成済み本文を返す | DEC-FEAT-015 | explicit | 整合 |
| 候補化・除外 | ユーザー由来の説明・推論・コード・訂正・自己申告・技術判断を対象にし、AI説明・記事・閲覧・評価・要約・質問のみを除外する | REQ-015、REQ-017、BR-005、DEC-FEAT-013 | explicit | 整合。除外入力にはCandidateもDecisionも作らず、候補ゼロは空一覧・`completed`である。 |
| Evidence強度 | 説明・推論・コード・訂正=`strong`、自己申告=`moderate`、概念認識=`weak`。理由を伴う技術判断は推論として扱う | BR-006、REQ-018 | explicit / derived | 整合。理由を伴わない判断の除外は、根拠のない知識化を避けるFeature内L2規則として明記済み。 |
| 候補の検索入力・順序 | `proposed_assertion`を先頭に、原文で明示されたConcept・Alias・Identifierだけを発話順・候補順で一回ずつ`search-text --query`へ渡す。IDは初出順で重複除去する | REQ-016、BR-010、FEAT-001 `search-text`契約 | derived | 整合。Candidate契約が検索列・停止・取得対象を一意に定義する。 |
| 更新選択・責務 | Codexが意味照合と操作選択を担い、CLIが検索・保存を決定論的に実行する | REQ-016、REQ-019、BR-010、初期設計 | explicit | 整合 |
| CLI / 永続化境界 | 新operation・公開JSON・schema・migrationを追加せず既存7 operationを消費する | REQ-014、初期設計、FEAT-004対象外 | explicit | 整合 |
| 訂正・置換・履歴 | 物理削除せず、二操作の後段失敗・中断は部分適用または結果不明として残す | REQ-018、BR-008、NFR-004、DEC-FEAT-014 | explicit | 整合 |
| Candidateゼロ・複数Candidate・停止後 | ゼロはDecision一覧`[]`と`completed`、複数は定義順で処理し、最初の停止後の残りは`not_started`で記録する | REQ-015、REQ-019、NFR-003、workflow-contracts | derived | 整合 |

## 資料間整合

- `requirements.md`、`design.md`、`workflow-contracts.md`、DEC-FEAT-013、DEC-FEAT-015は、Assessmentを**完成後、会話への返却前**に取り込む点で一致する。
- 候補化しない質問・AI説明・記事本文等はDecisionを持たず、候補がないEpisodeの結果は空Decision一覧と`completed`である。`skip`は許可済みCandidateだけに限定され、入力境界と矛盾しない。
- Candidateは`search_queries`を保持し、既存の単一query `search-text`を呼ぶ列、重複除去、検索失敗・中断時の停止、必要時の`get`/`get-evidence`対象を、本文・成果物契約・受入条件で一致して定義する。
- `source_ordinal`と`candidate_ordinal`による全Candidateの順序、成功・`skip`後の継続、失敗・中断・部分適用・結果不明後の`not_started`記録が、正常フロー、状態図、フローチャート、シーケンス図、成果物契約、受入条件で整合する。
- `revise`成功後の`attach-evidence`失敗・exit 130は`partially_applied`、`supersede`の`conflict`は`failed/outcome_unknown`とする扱いが、本文、図、成果物契約、DEC-FEAT-014で一致する。
- `create`、`attach-evidence`、`revise`、`supersede`に渡すEvidence kind・原文・観測時刻、および既存CLIの重複・conflict・exit 130契約は、FEAT-001の公開契約を変更せず消費している。

## 人間Decisionの境界

新たなL3/L4の未承認Decisionは検出しなかった。URL評価Workflowへ限定する取り込み境界と、Assessment本文完成後・返却前に同期更新する順序は、人間決定済みのDEC-FEAT-013、DEC-FEAT-015により確定している。候補の詳細化、順序、失敗分類は、その境界と既存CLI契約から導かれるFeature内L2規則として記録済みである。

## 次のゲート

利用者による詳細設計の明示承認を確認し、`task-breakdown`で論理TaskとImplementation Handoffを作成する。
