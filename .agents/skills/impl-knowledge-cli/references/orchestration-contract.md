# Knowledge CLI Orchestration Contract

このcontractは、仕様照合、実装、コードレビュー、最終整合性監査の受け渡しを固定する。各roleは該当sectionを満たし、根拠のない完了宣言を返さない。

## Source Authority

| 種別 | 権限 |
|---|---|
| `AGENTS.md`、workflow policy、artifact map | 作業方法、Decision level、Artifact Ownerを規定する |
| 要件、Business Rule、人間承認済みDecision | product intent、公開契約、意味を規定する |
| 詳細設計、operation資料 | 承認済み契約を具体化する |
| Task、handoff | 今回の実装範囲と完了条件を選ぶ |
| Issue | 進行と受入条件を追跡する。単独では契約を承認しない |
| 現行code、test、fixture | 実態の証拠。仕様逸脱を正当化しない |

矛盾は、明示的なDecision、承認記録、supersedes関係でだけ解消する。解消根拠がなければ隠さず判定へ反映する。

## Candidate Identity

すべての実装後roleは同じcandidate IDを参照する。candidate IDは、現在の`HEAD`と、ignore対象を除くworking treeのtracked・untracked変更をスクリプト自身が列挙し、path、内容、mode、削除状態から算出する。`python3 .agents/skills/impl-knowledge-cli/scripts/candidate_fingerprint.py --json`でmanifestとIDを生成し、candidate変更後は再生成する。呼出側が変更pathを指定して集合を狭めてはならない。

承認済みCLI candidateから別branchへbinary・fixture・コピーSkillなどの派生成果物を同期するhandoffでは、review/audit前にcandidateを対象実装branchへcommitし、`git status --porcelain`が空であることを確認してからcandidate fingerprintを算出する。このclean candidateの`HEAD` commit、candidate ID、`candidate_fingerprint.py --json`の完全出力、source IDをImplementation Report、Review Report、Audit Reportへ同じ値で残す。gate後にこのcommitを変更・amendしない。派生writerは、reportに記録されたcommitを隔離worktreeでcheckoutし、binaryなどCLI由来の派生成果物をそのworktreeからだけ生成する。未commitまたはdirty candidateへのPASSを派生成物の同期根拠にしてはならない。

CLI由来binaryとworkflow Skillコピーのsourceは同一commitである必要はない。workflowコピーは、workflow writerが自身の完了条件を満たしたclean workflow source commitから取得する。同期writerはorchestratorから受け取ったPASS済みReview ReportとAudit Reportの全文を改変せず、`verification/evidence/<cli-source-commit>/review-report.md`と`audit-report.md`へ保存する。各reportは同じCLI source commitとcandidate IDを明記し、専用検証branchへ同期するまで変更しない。同期manifestはJSON objectとし、`cli_candidate`へCLI source commit、candidate ID、`candidate_fingerprint.py --json`の完全な出力、normative source inventoryの順序付き配列、同じ配列への`source_fingerprint.py`出力を保存する。`workflow_source`へworkflow source commit、コピー対象rootの順序付き配列、同じ配列への`source_fingerprint.py`出力を保存する。`gate_evidence`へReview/Audit Report各々のrepo相対path、SHA-256、verdict、CLI source commit、candidate IDを保存し、report本文・manifest・CLI fingerprintの値が一致しない場合は停止する。`artifacts`へbinaryとコピー対象ファイルをrepo相対path順に列挙し、各pathのSHA-256を保存する。manifest自身のhashと同期commitは自己参照になるため保存せず、専用branchの同期commitとmanifest内容をGit履歴で照合する。これにより、異なる責務のwriterを混在させず、同期後に各成果物のsourceとgateを反証できる。

handoffがこの派生同期または人間専用Runtime受入を完了条件に含む場合、Audit PASS後のIssue Finalizationはdeferredとする。orchestratorはworkflow writerの同期reportとRuntime受入の実施/未実施状態を受け、全Taskの受入条件が満たされた時だけimplementerへIssue finalizationを委譲する。

Verification Packetは、正規source inventoryとしてFeature directoryなどのroot directoryと個別fileを列挙し、`python3 .agents/skills/impl-knowledge-cli/scripts/source_fingerprint.py <repo-relative-source>...`でdirectory配下を再帰列挙してsource IDを生成する。これにより新しいDecision等の追加も検知する。implementer、reviewer、auditorは各gate直前にsource IDを再計算し、仕様変更時はverifierへ戻す。

## Baseline Snapshot

- 対象Feature / Task / Issue
- 開始時刻
- baseline candidate ID、HEAD、`candidate_fingerprint.py --json`完全出力、`git status --short`
- 既存の変更fileと未追跡file、および既存変更のdiff要約
- 今回変更してよい範囲
- 既存変更と重なるfile、そのbaseline diff、および所有権が不明なfile

## Verification Packet

- `verdict`: `READY` | `BLOCKED` | `NEEDS_HUMAN_DECISION`
- source IDと、その算出に使った正規source inventory root・個別file一覧
- 対象範囲 / 対象外 / 依存Task
- 正規source一覧と承認状態
- 受入条件matrix: ID / 観測可能な条件 / 正規source / 実装surface / 公開I/OまたはDB契約 / 最も強い検証oracle / 実装可否
- 未承認契約または矛盾: ID / 根拠 / Decision level / Artifact Owner / 停止する範囲
- 実装者、reviewer、auditorへのrisk引継ぎ

## Implementation Report

- `verdict`: `READY_FOR_REVIEW` | `BLOCKED` | `NEEDS_SPEC_RECHECK`
- candidate IDと変更file一覧
- 再計算したsource ID
- 既存変更と今回変更の区別
- file配置matrix: 追加・移動file / owner / 主な呼出元 / 実行時期 / 配置理由 / 重複確認
- 受入条件別証拠: ID / production file・symbol / test・fixture・process case / 観測結果
- 公開契約差分と承認Decision。差分がなければ`none`
- migration / transaction / rollback証拠
- command results: exact command / pass・fail・not_run / exit status / 対象範囲
- 変更production packageごとのcoverage実測値
- 未解決事項と残存risk
- Issue finalization: `deferred` | `completed` | `failed`

長い生logをそのまま渡さず、失敗箇所と再現に必要な行を要約する。成功commandも実行した事実と対象範囲を残す。

## Review Report

- `verdict`: `PASS` | `BLOCKED` | `NEEDS_SPEC_RECHECK`
- candidate ID / 再計算したsource ID / 確認したdiff / 正規source / file配置matrixの照合結果
- findings: `REV-nnn` / severity / 根拠 / 受入条件 / 再現または不足test / 修正方向 / status
- finding status: `open` | `resolved` | `rejected_with_evidence` | `accepted_nonblocking_risk`
- reviewerが独立抽出した変更production package一覧とcoverage再測定結果
- reviewerが実行したcommand
- 非blocking riskと未検証範囲

## Audit Report

- `verdict`: `PASS` | `BLOCKED` | `NEEDS_SPEC_RECHECK`
- candidate ID / 再計算したsource ID
- 要件 / 仕様 / 実装・利用経路 / 実プロセス検証matrix
- impact map: 変更境界と影響する既存operation / DB version
- findings: `AUD-nnn` / 要件根拠 / 利用経路 / 再現と観測 / 不足oracle / 修正先 / status
- 実行したprocess case
- 監査済み範囲と残存risk

## Finding Lifecycle

ReviewとAuditのfindingは共通して`open`、`resolved`、`rejected_with_evidence`、`accepted_nonblocking_risk`のいずれかを持つ。blocking findingが`open`の間はPASSにしない。再review・再auditでは以前のfinding IDと根拠を入力に含め、同じ論点はIDを維持して状態を更新する。新しいauditorへは直前Audit findingを渡すが、Code Reviewの詳細findingは渡さない。

## Gate Rules

- `READY`以外で実装を開始しない。
- 実装中はIssueを完了更新しない。
- Implementation Report=`READY_FOR_REVIEW`、全必須command成功、source ID一致以外でcode reviewを開始しない。
- Review Report=`PASS`以外で最終整合性監査を開始しない。
- Audit Report=`PASS`以外でIssueをfinalizeしない。
- candidate変更後はImplementation Reportを更新し、code reviewとauditを再実行する。
- 公開契約・Decision findingは実装へ直接戻さず、仕様照合とArtifact Ownerへ戻す。
- 同じ原因のfindingが2回続いたら、前提、責務境界、oracleを再評価してから次のpatchを決める。
