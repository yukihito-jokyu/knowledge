# FEAT-002 詳細設計レビュー

- **対象:** `FEAT-002 Claim ごとの Agentic Knowledge Assessment` の詳細設計一式
- **レビュー範囲:** Claim単位の反復探索、7状態のAssessment、Search Trace、既存Knowledge CLI境界、Budget、技術失敗・中断の伝播
- **総合判定:** **pass**
- **人間Decision:** 不要。公開契約の追加・変更、L3/L4の未承認Decisionは検出しなかった。

## 正規根拠台帳

| 根拠 | 確定状態 | 適用範囲 |
| --- | --- | --- |
| GitHub Issue #175（要件定義）§6.2、§6.6〜6.8、§8、§12、§14〜21、§42〜45、§49、§52 | confirmed | Evidence起点、Codex/CLI境界、7状態、有限探索、Trace、未観測と未知の分離 |
| REQ-004、REQ-005、REQ-019〜REQ-021 | confirmed | 探索経路、Assessment、成果物の引渡し、再調査、Budget・Trace |
| BR-002、BR-003、BR-006、BR-007、BR-010 | confirmed | `no_evidence`、`inferable`、Evidence強度、Relation、責務境界 |
| NFR-001〜NFR-003、NFR-005〜NFR-006 | confirmed | 説明可能性、有限性、診断可能性、論理契約互換性、受入評価 |
| CON-002、CON-003、CON-006 | confirmed | 自作CLI、JSON/Markdown境界、物理実装を追加確定しないこと |
| DEC-REQ-001 | decided (L3) | 初期提供ではSemantic Searchを利用しない |
| DEC-FEAT-008 | decided (L2) | 1 Claim当たり最大12呼出し、Evidence最大4件、Relation深さ1、矛盾・時点探索の合計最大2回、公開設定を追加しない |
| FEAT-001 implementation handoff、Command Catalog、各operation設計、DEC-FEAT-006 | approved / ready | 既存11操作、JSON/error/exit code、`Ctrl-C`の無JSON・exit 130契約 |

## 契約根拠トレーサビリティ

| 境界 | 設計上の契約 | 根拠判定 | レビュー結果 |
| --- | --- | --- | --- |
| 入力 | 独立評価可能な単一Claimと、明示されたScope・時点だけを受け取る | explicit | 要件・Issueと整合。未指定filterを捏造しない。 |
| 検索戦略 | 直接字句探索から、構成要素/Concept、Relation、Evidence、矛盾、時点差分へ限定拡張する | explicit / derived | REQ-004、Issue #175 §42、DEC-FEAT-008と整合。`search-semantic`を呼ばない。 |
| CLI境界 | CLIは既存の検索・取得JSONだけを返し、意味対応、Evidence強度、状態、次queryはCodexが判断する | explicit | BR-010、CON-002/003、FEAT-001 handoffと整合。新規operation・JSON schema・migrationを追加しない。 |
| Assessment | 7状態の一つ、Confidence、Known、Gap、Evidence、矛盾・時点差分、Trace参照をMarkdownで返す | explicit | REQ-005、NFR-001、Issue #175 §17〜21・§47と整合。`no_evidence`を未知にしていない。 |
| 状態分類 | `known`と`inferable`を直接性で分離し、Relation単独を根拠にしない | explicit | BR-003、BR-007、Issue #175 §18・§49と整合。 |
| 通常CLI失敗 | `validation_error`、`not_found`、`storage_error`、`internal_error`、JSONプロトコル不整合は`technical_failure`で停止し、AssessmentなしでTraceとParentへ伝播する | explicit / derived | 本文、Trace、flowchart、sequence diagram、受入設計の記述が同じ分岐で一致する。空配列・未登録Conceptは成功結果であり混同しない。 |
| CLI中断 | response開始前の`Ctrl-C`によるexit 130・無出力はerror JSON/error codeではなく`canceled`で停止し、Assessmentなしで中断を伝播する | explicit | DEC-FEAT-006、本文、Trace、flowchart、sequence diagram、受入設計のすべてで`technical_failure`から分離されている。Parent自身も中断済みならTraceを出力・保存しない。 |
| Trace | 操作、入力要約、候補/Evidence ID、増分、Budget、停止理由をAssessmentと分離する | explicit | NFR-003、Issue #175 §43、CON-003と整合。Evidence原文を複製しない。 |
| Budget | 全CLI呼出し最大12、探索最大4、補助探索最大4、Evidence取得最大4、Relation深さ1、矛盾・時点探索合計最大2 | explicit | DEC-FEAT-008の全上限を本文、Search Trace、受入設計が同じ意味で定める。 |

## 資料間整合監査

| 照合対象 | 結果 |
| --- | --- |
| 要件 → Feature Map / Traceability | REQ-004、REQ-005、REQ-019〜021とBR-002、BR-003、BR-006、BR-007、BR-010のFEAT-002割当ては一致する。 |
| FEAT-002要件 → 設計 / Assessment / Trace | 反復探索、7状態、Evidence根拠、有限Budget、技術失敗の非状態化、Markdown成果物分離は一致する。 |
| DEC-REQ-001 → CLI利用表 | `search-semantic`を呼ばず、字句・Concept・Relation・Evidence・矛盾・時点差分だけを使うため一致する。 |
| FEAT-001公開CLI → 操作列 | `search-text`、`search-concept`、`get`、`get-evidence`、`search-related`、`search-contradictions`、`search-temporal`だけを消費する。`search-temporal`はConceptまたはScope selectorが必須のため、時点だけでは実行しない。 |
| FEAT-001公開CLI → 通常エラーの伝播 | `validation_error`（exit 2）、`not_found`（exit 3）、`storage_error` / `internal_error`（exit 1）は、Assessmentなし→`technical_failure` Trace→Parentへ評価失敗、という一意の経路で本文・Trace・両図・受入設計が一致する。 |
| FEAT-001公開CLI → 中断の伝播 | response開始前のexit 130・無出力は、error codeなし→`canceled`部分Trace→Assessmentなしの中断伝播で一意に一致する。 |
| DEC-FEAT-008 → 12回上限 | 4回探索 + 4回補助探索 + 4回Evidence取得は最大12回であり、13回目を禁止する受入設計もある。 |
| DEC-FEAT-008 → 矛盾・時点探索 | 本文のMain FlowとBudget / Stop Conditionsは両operationを合計2回までとし、到達後はGapへ残す。Traceは`contradiction_temporal_calls: <0..2>`を必須記録し、2到達後の不実行を規定する。受入設計も合計3回目を禁止する。したがって前回R-002は解消済み。 |
| Scenario A〜J / NFR-006 | 空Store、完全一致、構成知識、矛盾、古い知識、質問/AI説明、競合Evidenceを個別に扱う。横断FixtureセットをFEAT-005が所有する分離もFeature Mapと整合する。 |

## 差し戻し事項

なし。前回R-002（`search-contradictions`と`search-temporal`の合計上限2回）は、本文・Trace・受入設計へ一貫して反映された。

## 総合判定と次のゲート

**判定:** `pass`

承認済みL2 Decisionを含むBudget制約は実装可能な形で資料間一致している。公開CLI/JSON/永続化/設定の追加・変更、未承認のL3/L4 Decision、実装者へ再設計を委ねる欠落は検出しなかった。

**次のゲート:** 利用者によるFEAT-002詳細設計レビューと明示承認。承認前にtask breakdownまたはimplementation handoffへ進めない。
