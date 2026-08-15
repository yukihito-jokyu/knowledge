# FEAT-002 詳細設計レビュー

- **対象:** `FEAT-002 Claim ごとの Agentic Knowledge Assessment` の詳細設計一式、およびそれを実装へ渡す `tasks.md` / `implementation-handoff.yaml`
- **レビュー原典:** Issue #197、REQ-004 / REQ-005 / REQ-019〜021、BR-002 / BR-003 / BR-006 / BR-007 / BR-010、NFR-001〜003 / 005 / 006、CON-002 / 003 / 006、承認済みDecision、Initial Design、FEAT-001の公開CLI契約
- **総合判定:** **pass**
- **人間Decision:** 不要。利用者は2026-08-15に実装先をKnowledge CLIではなくCodex側workflowとする方針を選択済みであり、L3/L4の未承認Decisionは検出しなかった。

## 正規根拠台帳

| 根拠 | 確定状態 | 適用範囲 |
| --- | --- | --- |
| Issue #197 | confirmed | Claimごとの直接探索、限定拡張、7状態、`inferable`、`search-temporal`のselector条件、有限停止、資料同期 |
| REQ-004、REQ-005、REQ-019〜REQ-021 | confirmed | 反復探索、Assessment、Codexワークフローの成果物引渡し、Budget・Trace |
| BR-002、BR-003、BR-006、BR-007、BR-010 | confirmed | 未観測と未知の分離、`inferable`、Evidence強度、Relation、Codex/CLI責務境界 |
| CON-002、CON-003、CON-006 | confirmed | CLIにAI判断を置かないこと、CLI JSON／Codex間Markdown境界、物理実装を要件だけで確定しないこと |
| DEC-REQ-001 | decided (L3) | 初期提供でSemantic Searchを利用しない |
| DEC-FEAT-008 | decided (L2) | Claimごとの固定Budget、探索順序、停止条件。公開設定を追加しない |
| DEC-FEAT-009 | decided (L2) | 評価とAssessment/TraceはCodex側Knowledge Search workflowで実行し、TASK／handoffを再分解する |
| FEAT-001 Command Catalog、各operation契約、DEC-FEAT-006 | approved / ready | 既存11操作、JSON/error/exit code、response開始前のexit 130・無出力 |

## 契約根拠トレーサビリティ

| 境界 | 設計上の契約 | 根拠判定 | レビュー結果 |
| --- | --- | --- | --- |
| 実行責務 | Claimの意味、query、探索継続、Evidence強度、7状態、Assessment/TraceはCodex Knowledge Searchが担う | explicit | BR-010、CON-002/003、architecture、DEC-FEAT-009と整合。Knowledge CLIへGo実装を追加しない。 |
| CLI公開I/O / DB / 設定 | 既存read operationだけを既存JSON、stdout/stderr、exit code契約で同期利用する。operation、wire schema、SQLite、migration、公開Budget設定は変えない | explicit | FEAT-001 handoff・Command Catalog・各operation資料と整合。`search-semantic`も呼ばない。 |
| 入力・通常返却 | 一つの独立評価可能なClaimと明示されたScope/version/時点だけを受け取り、正常時はAssessmentとTrace参照を返す | explicit | workflow契約、Assessment／Trace契約は一致。未指定filterを推測しない。 |
| 状態判定 | 7状態を優先規則とEvidence原文、Scope、時点、意味対応で区別する。`no_evidence`は未知断定でなく、Relation/候補類似だけは根拠にしない | explicit | Issue #197、REQ-005、BR-002/003/006/007と整合。 |
| 時点探索 | `search-temporal`はConceptまたはScope selectorがある場合だけ利用する | explicit | Issue #197およびFEAT-001 `search-temporal` validationと整合。 |
| Budget・停止 | 12呼出し、探索4、補助探索4、Evidence 4、Relation深さ1、矛盾・時点探索合計2。同一論理queryを再実行しない | explicit | DEC-FEAT-008、design、Trace契約、受入設計で一致。 |
| 技術失敗・中断 | 既存error JSON／プロトコル不整合は`technical_failure`、response開始前のexit 130・無出力は`canceled`としてAssessmentなしでParentへ伝播する | explicit / derived | DEC-FEAT-006と設計・workflow・Traceで一致。空結果は成功結果として分離される。 |
| 実装handoff | Codex側workflowの実装・検証対象と依存を、CLI Goのapplication/domainから区別して渡す | downstream action | DEC-FEAT-009に従い、利用者の詳細設計承認後のtask-breakdownがTASK／handoffを再分解する。現行Task／handoffは前回設計時点の下流成果物であり、このレビューの不整合ではない。 |

## 資料間整合監査

- `requirements.md`、`design.md`、`knowledge-search-workflow.md`、Assessment／Trace契約、DEC-FEAT-008／009の間では、Codex側実行、CLI境界、入出力、Budget、停止、技術失敗／中断の規則が整合する。未承認のCLI公開I/O、SQLite、migration、設定の追加は検出しなかった。
- FEAT-001の各operation契約と照合した。`search-temporal`のselector必須条件、空結果の成功扱い、既存error、exit 130の扱いに矛盾はない。
- 現行リポジトリの製品実装はGo CLI（`cmd/knowledge`、`internal/**`、`test/integration`）だけであり、Codex Knowledge Searchを実行する既存workflow artifact／実装先は存在しない。これはCLIを変更する根拠ではない。設計はCodex AI Runtime内の論理workflow、入出力、停止、失敗／中断、成果物境界を定義しており、具体的なartifact・ファイル配置はPlanning/Implementation境界に従い下流で選ぶ。
- `tasks.md` と `implementation-handoff.yaml` はDEC-FEAT-009前の下流成果物であり、現workflow stateでは更新済み詳細設計の人間承認およびtask-breakdownがまだ完了していない。そのため現時点で実装へ用いてはならない。承認後、task-breakdown ownerがDEC-FEAT-009に従ってTASK-002-01〜03、logical area、依存、handoffをCodex側workflow向けに更新する必要がある。

## 下流工程への注意

- DEC-FEAT-009の影響どおり、利用者がこの更新済み詳細設計を承認した後に、`task-breakdown` ownerは旧Task／handoffをCodex側workflow用へ再分解・更新する。
- その更新前の`implementation-handoff.yaml`はGo Knowledge CLI実装へ渡してはならない。これは本レビューの`pass`を妨げないが、implementation-readyへ進む前の必須更新である。

## 総合判定と次のゲート

**判定:** `pass`

Codex側へ移す方針、既存CLI契約の消費、実行境界、入出力、Budget／停止、Trace／Assessment、技術失敗／中断の設計は正規根拠に適合する。次のゲートは利用者による更新済み詳細設計の明示承認であり、その後にtask-breakdownを開始する。既存Task／handoffはこの設計の実装開始根拠に使わない。
