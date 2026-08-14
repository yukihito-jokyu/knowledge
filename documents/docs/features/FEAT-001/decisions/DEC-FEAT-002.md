# DEC-FEAT-002: 初期 Knowledge CLI の公開 JSON 契約

- **Status:** superseded_in_part
- **Level:** L3
- **論点:** Issue #175 が列挙する初期 Knowledge CLI の論理操作について、JSON出力契約と Index 管理の公開範囲を誰がどの根拠で確定するか。入力形式は DEC-FEAT-003 により更新された。

## 確定済みの根拠

- CLI は Codex との機械的境界で JSON を入出力し、保存・検索・取得・更新・Index 管理を決定論的に担う（REQ-014、CON-002、CON-003、Issue #175 §8.7、§9、§40）。
- 論理操作として `search-text`、`search-concept`、`search-related`、`get`、`get-evidence`、`search-contradictions`、`search-temporal`、`create`、`attach-evidence`、`revise`、`supersede` が示される（Issue #175 §12、§13）。
- `search-semantic` は初期提供では公開しない（DEC-REQ-001、DEC-FEAT-001）。

## 決定

人間判断として **A** を採用する。`docs/features/FEAT-001/design.md`、`design/command-catalog.md`、`design/commands/**`、`design/database-schema.md` および `design/database-access-map.md` に記載した v1 契約を承認する。

- 11 の論理操作は `knowledge <operation>` として提供し、成功 response を標準出力、error response を標準エラーへ出力する。入力は [DEC-FEAT-003](DEC-FEAT-003.md) の名前付きoption契約に従う。
- 共通の success / error envelope、error code、終了コード、未知 field・省略・`null` の規則は `design.md` の「Interfaces / Data」に従う。
- 操作別の入力option / response field、必須性、Evidence / Relation の enum、重複・競合時の挙動は各 command 資料に従う。
- Index 管理は CLI の責務とし、mutation または migration と同じ transaction で内部的に維持する。独立した公開 Index 管理操作、保存先、設定、運用仕様は v1 契約に含めない。
- SQLite schema は公開 JSON 契約を実現する内部表現であり、保存先を規定しない。将来の Semantic Search は FEAT-006 が安定 Assertion ID と現行データから追加する。

## 決定前に未確定だった事項

- 各操作の request / success / error JSON の field、必須性、enum、ネスト、代表値。
- JSON の受渡し方法、stdout / stderr の割当、終了コード。
- Index 管理を公開操作にするか、公開する場合の操作・I/O、公開しない場合のライフサイクル責任。
- 公開契約のバージョニングおよび互換性表現。

## なぜ今決める必要があるか

本 Feature は JSON CLI と SQLite Store を持ち、利用者から操作別 I/O・SQL・図を作成するよう求められている。未決定のままでは、Operation Documentation Coverage Gate と Contract Completeness Gate を満たせず、Task Breakdown へ進めない。

## 選択肢

### A. 承認済み根拠に基づく詳細設計案を、公開 v1 契約として人間が承認する

Feature Design が各 field と transport を提案し、人間がレビューして v1 契約として承認する。

- **利点:** 現行要件を保ったまま操作別資料を完成できる。
- **欠点:** 承認前の提案を設計資料としてレビューする必要がある。

### B. 人間が既存または希望する v1 契約を提示する

提示された field、transport、Index 管理方針を根拠として Feature Design に反映する。

- **利点:** 公開契約の意思決定を明確に所有できる。
- **欠点:** 詳細設計の開始に必要な入力が増える。

### C. 初期提供の公開 CLI 範囲を縮小または後続 Feature へ移す

CLI の論理能力は保持しつつ、公開 I/O を要する操作の提供時期を要件・Feature Planning から再判断する。

- **利点:** v1 公開契約を今確定しない選択肢を取れる。
- **欠点:** REQ-014 と FEAT-001 の受入条件に影響し、上流成果物の変更が必要になる。

## 採用理由

Issue #175 が論理操作と JSON 境界を明示しているため、根拠にない保存先・設定・運用仕様を追加せず、操作別 wire schema だけを最小の v1 契約として承認する。

## 影響範囲

- FEAT-001 の操作別設計、SQLite schema、Migration、テスト、実装 handoff。
- FEAT-002、FEAT-004、FEAT-006 が利用する Knowledge CLI の互換性境界。

## 解消した停止条件

- `docs/features/FEAT-001/design.md` と `design/commands/**` の operation I/O。
- Design Review、Task Breakdown、`implementation-handoff.yaml`。
