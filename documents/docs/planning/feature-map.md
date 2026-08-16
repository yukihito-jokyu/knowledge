# Feature Map

## FEAT-001: Evidence 起点の個人知識基盤

- **Goal:** ユーザー知識を Assertion、Evidence、Concept、Scope、Relation、Temporal Metadata として保存・再評価可能にする。
- **Primary Actor / Trigger:** Codex が Knowledge CLI を通じて知識を検索または更新する。
- **Covered Requirements:** REQ-009, REQ-010, REQ-011, REQ-012, REQ-013, REQ-014
- **Related Business Rules:** BR-002, BR-004, BR-007, BR-008, BR-010
- **High-level Acceptance:** Assertion を Evidence と分離して保持し、論理的検索・取得・更新プリミティブを JSON 境界で利用できる。履歴を失わずに訂正・置換を表現できる。
- **Dependencies:** なし
- **Status:** planned

## FEAT-002: Claim ごとの Agentic Knowledge Assessment

- **Goal:** Target Claim に対し、Evidence を反復探索してユーザー知識の状態、Known、Knowledge Gap、矛盾・時点差分を判断できるようにする。
- **Primary Actor / Trigger:** Article Claim を受け取った Codex が Knowledge Assessment を実行する。
- **Covered Requirements:** REQ-004, REQ-005, REQ-019, REQ-020, REQ-021
- **Related Business Rules:** BR-002, BR-003, BR-006, BR-007, BR-010
- **High-level Acceptance:** 直接一致がなくても分解・Concept・Relation・Evidence・矛盾・時点差分の探索を制御し、7つの状態と根拠を含む Assessment を返す。停止条件と Search Trace を適用できる。
- **Dependencies:** FEAT-001
- **Status:** planned

## FEAT-003: URL 記事の読書価値評価

- **Goal:** 技術記事を Claim に分解し、個人知識との差分と注意コストから読書推奨および対象箇所を説明する。
- **Primary Actor / Trigger:** ユーザーが技術記事 URL の評価を依頼する。
- **Covered Requirements:** REQ-001, REQ-002, REQ-003, REQ-006, REQ-007, REQ-008, REQ-019, REQ-020
- **Related Business Rules:** BR-001, BR-003, BR-009, BR-010
- **High-level Acceptance:** 記事を役割・重要度・位置・根拠付き Claim に分解し、`read_full`、`read_selected`、`skip` を理由、読む箇所、飛ばす箇所、認識利得、信頼性とともに返す。
- **Dependencies:** FEAT-002
- **Status:** planned

## FEAT-004: 会話・作業からの知識 Evidence 更新

- **Goal:** ユーザーの会話・作業エピソードから知識候補を抽出し、既存知識に対する適切な更新を行う。
- **Primary Actor / Trigger:** 意味のある Conversation / Task Episode が終了する。
- **Covered Requirements:** REQ-015, REQ-016, REQ-017, REQ-018, REQ-019
- **Related Business Rules:** BR-004, BR-005, BR-006, BR-008, BR-010
- **High-level Acceptance:** Evidence 候補を抽出して既存知識を照合し、create、attach-evidence、revise、supersede、skip を選べる。AI の説明・記事閲覧・質問だけから不正な更新をしない。
- **Dependencies:** FEAT-001
- **Status:** planned

## FEAT-005: 知識評価の品質保証

- **Goal:** Knowledge Search、Knowledge Update、Reading Value の正確性と原因追跡可能性を、受入シナリオと層別評価で確保する。
- **Primary Actor / Trigger:** 開発者が変更の妥当性を検証する。
- **Covered Requirements:** REQ-021
- **Related Business Rules:** BR-002, BR-003, BR-005, BR-008, BR-009
- **High-level Acceptance:** Issue #175 の Scenario A〜J を含む Fixture で、検索・更新・評価の各層を検証でき、Search Trace により失敗箇所を切り分けられる。
- **Dependencies:** FEAT-001, FEAT-002, FEAT-003, FEAT-004
- **Status:** planned

## FEAT-006: ローカル Semantic Search 移行

- **Goal:** 初期の字句検索基盤へ、ローカル Embedding とローカル Vector Index による Semantic Search を追加し、既存の個人知識・履歴を破壊せずに検索能力を拡張する。
- **Primary Actor / Trigger:** FEAT-001 が提供する既存 Knowledge Store に対し、Semantic Search の利用を開始する。
- **Covered Requirements:** REQ-011
- **Related Business Rules:** BR-002, BR-004, BR-010
- **High-level Acceptance:** 既存 Assertion から Embedding と Index を構築し、意味的候補を返す。既存 Evidence・Relation・履歴・字句検索の契約を維持し、類似候補を既知判定へ変換しない。
- **Dependencies:** FEAT-001
- **Status:** planned（初期提供後）

## FEAT-007: 隔離Knowledge Storeの明示選択

- **Goal:** Codex workspace sandboxなどでも、既定Storeを変更せずに利用者が書込み可能なローカルSQLite Storeを選択してKnowledge CLIとworkflowを実行できるようにする。
- **Primary Actor / Trigger:** 利用者がStoreパスを明示してKnowledge CLIまたはKnowledge Update workflowを実行する。
- **Covered Requirements:** REQ-022
- **Related Business Rules:** BR-010
- **High-level Acceptance:** 明示Storeで全Store利用operationが一貫して実行され、指定なしの既存操作は既定Storeを維持する。不正な指定は構造化CLI errorとして扱い、Codexの検証手順から同じStoreを利用できる。
- **Dependencies:** FEAT-001
- **Status:** planned

## 境界レビュー

- Feature は利用価値・システム能力で分けており、DB・API・UI などのコード層では分割していない。
- FEAT-003 は評価対象の記事を扱い、FEAT-004 はユーザー知識の観測・更新を扱うため、独立した状態変化と受入条件を持つ。
- UI-001 と UI-002 は Feature の目的を変えず、提供インターフェースと入力境界を決める横断設計事項として残す。
- FEAT-006 はコード層の追加ではなく、利用可能な検索能力を段階的に拡張する独立したシステム能力として分離する。
- FEAT-007 は保存先をOS設定領域から移す作業ではなく、sandboxを含む利用環境で明示的に隔離Storeを選ぶ能力として分離する。
