# FEAT-007 実装タスク

> **前提:** 利用者による詳細設計の明示承認済み。独立設計レビューは`pass`であり、実装者は公開契約を再決定しない。

## Design Readiness Audit

| 監査項目 | 結果 | 根拠 |
| --- | --- | --- |
| 利用者による詳細設計承認 | pass | `.ai/workflow/state.yaml` の `features.FEAT-007.human_design_review: approved` |
| 独立設計レビュー | pass | [design-review.md](design-review.md) |
| CLI input・JSON error・exit契約 | pass | [design.md](design.md)「Interfaces / Data」、[DEC-FEAT-023](decisions/DEC-FEAT-023.md) |
| Store選択・migration・fallback禁止 | pass | [design.md](design.md)「Store解決契約」、`State / Interaction` |
| 全11operationへの適用範囲 | pass | [design.md](design.md)「Operation Documentation」およびFEAT-001の全operation資料 |
| Go実装基盤・検証基盤 | pass | [architecture.md](../../design/architecture.md)「Go module と配置規約」、既存Go moduleとtest/integration |
| 検証用workflowの入力・Store伝播 | pass | [design.md](design.md)「検証用workflowのStore伝播入口」 |

**結論:** pass。実装時に振る舞い、公開CLI、Store、error、migration、workflow入力を再設計する未決定事項はない。

## 固定済み実装基盤

- Go単一module、CGO不要のSQLite driver、埋込みmigration、責務分離はInitial Designと既存実装に従う。
- 共通チェックは`gofmt`、`go test ./...`、`go vet ./...`とし、公開CLIは実プロセス境界で検証する。
- `--store`はoperation前に一度だけ置く任意global optionである。指定なしの既定Store、success JSON、既存operation固有I/Oは維持する。
- `verification/`の変更は、review/audit済みのclean source commitを隔離worktreeでbuildしてから、専用ブランチ`chore/verification-environment`だけに追従コミットし、Feature実装ブランチへ混在させない。

## タスク一覧

| ID | タイトル | 論理領域 | 依存 |
| --- | --- | --- | --- |
| TASK-007-01 | 明示Store選択を共通CLI境界へ実現する | interface / persistence | なし |
| TASK-007-02 | Store選択の公開資料と通常Skill参照を整合する | documentation / workflow | TASK-007-01 |
| TASK-007-03 | 隔離検証workflowを専用ブランチへ追従する | workflow verification | TASK-007-01, TASK-007-02 |
| TASK-007-04 | Store選択の互換性と全operation適用を横断検証する | verification | TASK-007-01, TASK-007-02, TASK-007-03 |

## TASK-007-01: 明示Store選択を共通CLI境界へ実現する

- **目的:** 全11operationが、明示した絶対Storeまたは従来の既定Storeを一意に選び、同じ初期化・migration・error境界を利用できるようにする。
- **関連要件:** REQ-022、NFR-007、CON-002、CON-003、CON-007、BR-010、DEC-FEAT-023。
- **論理領域:** interface / persistence。
- **作業内容:** global `--store`の構文・位置・一回性・絶対パス検証、指定/既定Storeの解決、既存の親作成・open・migrationと全operationへの選択済みStore伝達を実現する。
- **受入条件:**
  - `knowledge --store <absolute-path> <operation> [operation options]`が全11operationで同一の指定Storeだけを利用する。
  - 指定なしは`os.UserConfigDir()/knowledge/knowledge.db`を従来どおり利用する。
  - 値なし、重複、相対パス、不正位置はStoreをopenせず、stderr JSONの`validation_error`とexit 2になる。
  - 親作成、DB open、migrationの失敗は指定Store以外へfallbackせず、stderr JSONの`storage_error`とexit 1になる。
  - success JSON、既存error envelope、stdout/stderr、exit code、Ctrl-C規約、domain I/O、schema/migration assetを変更しない。
- **依存:** なし。
- **対象外／注記:** 環境変数、設定ファイル、既定Storeの移動、権限昇格、Store共有は対象外。

## TASK-007-02: Store選択の公開資料と通常Skill参照を整合する

- **目的:** 新旧の正規資料と通常workflowが矛盾なく、呼出側が必要時にglobal optionを扱えるようにする。
- **関連要件:** REQ-022、NFR-007、CON-007、DEC-FEAT-023。
- **論理領域:** documentation / workflow。
- **作業内容:** FEAT-001のCLI入力、Store、Decision、task、handoffに残る保存先option禁止の限定範囲をDEC-FEAT-023に整合し、通常SkillのCLI operation参照へ任意global optionの位置を反映する。
- **受入条件:**
  - 正規資料に`--store`を禁止する記述が残らず、指定なし既定Storeの契約は維持される。
  - global optionの位置、不正入力、error/exit契約がFEAT-007設計と一致する。
  - 通常のroot `skills/`はStoreを自動指定せず、呼出側の明示Store実行コンテキストだけを正しいargvで扱える。
- **依存:** TASK-007-01。
- **対象外／注記:** Reading ValueのURL-only入力契約と通常Skillの既定Store挙動は変えない。

## TASK-007-03: 隔離検証workflowを専用ブランチへ追従する

- **目的:** 利用者が検証workspaceでURLだけを渡し、Codex sandboxから同一の隔離Storeを使って評価・更新できるようにする。
- **関連要件:** REQ-022、NFR-007、CON-007、DEC-FEAT-023。
- **論理領域:** workflow verification。
- **作業内容:** CLI実装branchへreview/audit前にcommitされたクリーンcandidateについて、独立コードレビューと最終整合性監査がともにPASSしたCLI source commitとcandidate IDをbinary同期元として固定する。そのsource commitの隔離worktreeだけでbinaryをbuildする。orchestratorから渡されたPASS済みReview/Audit Report全文を改変せず`verification/evidence/<cli-source-commit>/`へ置き、report本文のsource commit/candidate IDが同期元と一致することを確認する。root Skillコピーは、`reading-value`、`knowledge-acquisition`、`knowledge-search`、`knowledge-update`の順序付き4 rootだけを、`impl-knowledge-workflow`が自身の検証を完了してcommitしたclean workflow source commitから得る。各コピーは、利用者が`verification/`をworkspaceとして開いた物理パスを`pwd -P`で一度取得し、`.agents/skills/reading-value/`と`bin/knowledge`のmarkerを確認してから`knowledge.db`の絶対パスを決める。`source-manifest.json`へ、CLI candidateの`candidate_fingerprint.py --json`完全出力、source inventoryの順序付きroot一覧と`source_fingerprint.py`出力、両source commit、workflowのコピー対象root一覧と`source_fingerprint.py`出力、Review/Audit Report各path/SHA-256/verdict/source commit/candidate ID、binary/コピーSkillのpathごとのSHA-256を保存する。manifest自身のhashおよび同期commitは記録せず、同一の専用branch commitとGit履歴で確定する。
- **受入条件:**
  - 利用者は`verification/`をworkspaceとして開き、`$reading-value <一件の絶対HTTP(S) URL>`だけを渡せる。
- 検証用Skill群は、利用者が開いた`verification/` workspaceの`pwd -P`とmarker確認だけを使い、環境変数、設定、URL以外の利用者入力、任意CWDに依存せず、全CLI呼出しへ同じ`verification/knowledge.db`を渡す。
  - binary不在または固定配置からStoreを解決できない場合は、CLIを起動せず技術的停止となり、既定Storeへのfallbackをしない。
- review/audit PASS済みのCLI source commit、clean candidate IDとcandidate fingerprint JSON、source inventory/fingerprint、workflow source commit/root/fingerprint、Review/Audit Report全文と各hash、生成binary、コピーSkillの全pathが`source-manifest.json`と同期Git commitで反証可能に対応付けられ、変更は`chore/verification-environment`にのみcommit/pushされる。
- **依存:** TASK-007-01、TASK-007-02。
- **対象外／注記:** 検証用Skillを別配置へコピーして使うこと、環境変数・設定によるStore指定は対象外。

## TASK-007-04: Store選択の互換性と全operation適用を横断検証する

- **目的:** 公開CLI、SQLite初期化、全operation、workflow伝播が既定利用を壊さず隔離Storeで成立することを観測する。
- **関連要件:** REQ-022、NFR-007、CON-002、CON-003、CON-007、DEC-FEAT-023。
- **論理領域:** verification。
- **作業内容:** unitと実プロセスのintegration test、通常Go検査、利用者実行可能な検証記録を用い、指定/既定Store、migration、error I/O、11operation、verification workflowを横断確認する。AIの静的Skill照合・CLI process検証と、利用者だけが実行するCodex Runtime受入を分離する。
- **受入条件:**
  - workspace内の一時絶対パスで初回migration、create後のread、再実行を実プロセスで観測できる。
  - 指定なし既定Store互換、validation errorとstorage errorのstdout/stderr/exit codeを観測できる。
  - 全11operationが指定Storeで少なくとも一度起動し、同じStoreを利用する。
  - `gofmt`、`go test ./...`、`go vet ./...`が成功する。
  - 利用者が実行する`$reading-value <一件の絶対HTTP(S) URL>`について、追加入力を要求しない会話記録、全CLI invocationの同一絶対`--store` argv記録、`verification/knowledge.db`の生成または更新、既定Storeへfallbackしなかった確認を受領できる。利用者実行がない場合、このRuntime受入は未実行のまま完了扱いにしない。
- **依存:** TASK-007-01、TASK-007-02、TASK-007-03。
- **対象外／注記:** sandbox外の既定Storeへ書き込む実機検証は要求しない。

## Feature 全体の対象外

- 既定Storeの移動、環境変数・設定ファイル、権限昇格、同期・共有・remote Store。
- SQLite schema、migration asset、domain I/O、既存11operationの固有JSON schemaの変更。

## Implementation Readiness Review

- **reviewer:** `feat007_readiness_review`（Task作成者と異なる独立subagent）
- **親リポジトリ俯瞰範囲:** `cmd/knowledge`の共有Store composition、`test/integration`の実binary oracle、`impl-knowledge-cli`、`impl-knowledge-workflow`、root `skills/`、専用`chore/verification-environment`ブランチ、適用される`AGENTS.md`。
- **初回finding:** IR-007-001（Planning Owner外の旧正規資料更新）、IR-007-002（専用検証branchのwriter/source同期）、IR-007-003（人間専用Runtime受入oracle）。
- **継続finding:** IR-007-002は、candidate IDがdirty worktree fingerprintでありsource revisionではないため、同期元commitと永続証跡が未定義として再オープンした。CLIのclean committed candidate規則とCLI/workflow source分離を追加後、補正後の独立再レビューを要する。
- **finalization:** Issue #233はTASK-007-03の同期reportとTASK-007-04の人間専用Runtime受入が完了または未実行として明確になるまでfinalizeしない。
