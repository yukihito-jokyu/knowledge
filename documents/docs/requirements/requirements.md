# 機能要件

## URL 評価

### REQ-001

- **Statement:** ユーザーは技術記事の URL を起点に、記事を読む価値の評価を依頼できる。
- **Source:** Issue #175 §1, §38, §50, §52
- **Rationale:** 読書前の注意コストを Codex が肩代わりするため。
- **Status:** confirmed

### REQ-002

- **Statement:** システムは記事を、ユーザー知識と独立に比較可能な複数の Article Claim へ分解する。
- **Source:** Issue #175 §8.2, §27, §29, §52
- **Rationale:** 記事全体の一括検索では差分と読むべき箇所を判定できないため。
- **Status:** confirmed

### REQ-003

- **Statement:** 各 Article Claim は本文、記事内の役割、重要度、位置、記事内の根拠を保持する。
- **Source:** Issue #175 §27, §28, §38, §47
- **Rationale:** 未知 Claim の数だけで記事の価値を決めないため。
- **Status:** confirmed

### REQ-004

- **Statement:** システムは Claim ごとに Knowledge CLI を反復検索し、直接・部分・関連・前提・矛盾・時点差分の Evidence を評価可能にする。
- **Source:** Issue #175 §8.3, §14, §15, §42, §52
- **Rationale:** 一回の検索結果だけではユーザー知識との差分を適切に判定できないため。
- **Status:** confirmed

### REQ-005

- **Statement:** システムは Claim ごとに `known`、`partially_known`、`inferable`、`contradicted`、`outdated`、`no_evidence`、`uncertain` のいずれかを、判定確信度とともに示せる。
- **Source:** Issue #175 §17–§20
- **Rationale:** 既知・未知を二値化せず、Evidence に基づく知識差分を明示するため。
- **Status:** confirmed

### REQ-006

- **Statement:** システムは Knowledge Assessment を統合し、`read_full`、`read_selected`、`skip` のいずれかと、その理由を返す。
- **Source:** Issue #175 §30–§36, §52
- **Rationale:** ユーザーの読書判断と注意コストを支援するため。
- **Status:** confirmed

### REQ-007

- **Statement:** `read_selected` の場合、システムは読む価値がある箇所、読む必要が薄い箇所、得られる認識利得および根拠の信頼性を説明する。
- **Source:** Issue #175 §34, §36, §38, §52
- **Rationale:** 記事全体を読む注意コストを避けつつ価値を保つため。
- **Status:** confirmed

### REQ-008

- **Statement:** Reading Value の最終出力は、単一の数値スコアではなく理由を含む自然言語の評価とする。
- **Source:** Issue #175 §35
- **Rationale:** 評価理由のブラックボックス化を防ぐため。
- **Status:** confirmed

## Knowledge Store と検索

### REQ-009

- **Statement:** システムは Knowledge Assertion を中心に、Concept、Evidence、Scope、Relation、Temporal Metadata を扱える。
- **Source:** Issue #175 §10
- **Rationale:** Claim との差分を意味単位で検索・再評価するため。
- **Status:** confirmed

### REQ-010

- **Statement:** システムは Assertion 本文、Concept、Alias、Scope、API 名および Identifier を対象とする字句検索を提供する。
- **Source:** Issue #175 §11, §12
- **Rationale:** 固有名詞や API 名を含む知識へ到達するため。
- **Status:** confirmed

### REQ-011

- **Statement:** システムは正規化済み Assertion を中心に意味検索を提供し、意味検索結果のみでユーザーが知っているとは判定しない。初期提供では DEC-REQ-001 により実装を延期し、後続移行で提供する。
- **Source:** Issue #175 §11, §12, §41
- **Rationale:** 類似度と知識保有を混同しないため。
- **Status:** confirmed（初期提供は deferred）

### REQ-012

- **Statement:** システムは Concept、Relation、Evidence、矛盾候補、時点差分を入口とする検索・取得を提供する。
- **Source:** Issue #175 §11–§13
- **Rationale:** 直接一致しない部分知識や古い理解を調査するため。
- **Status:** confirmed

### REQ-013

- **Statement:** システムは Knowledge Assertion の作成、Evidence 追加、正規化表現・Scope の修正、および旧知識を削除しない置換関係の保存を提供する。
- **Source:** Issue #175 §13, §25
- **Rationale:** 知識の訂正と履歴保持を両立するため。
- **Status:** confirmed

### REQ-014

- **Statement:** Knowledge CLI は Codex との機械的境界で JSON を入出力し、保存・検索・取得・更新・Index 管理を決定論的に実行する。
- **Source:** Issue #175 §8.7, §9, §40, §48
- **Rationale:** AI の意味判断と永続化基盤の責務を分離するため。
- **Status:** confirmed

## 知識の継続更新

### REQ-015

- **Statement:** システムは意味のある会話または作業エピソードの終了後に、ユーザー発言・推論・コード・技術的訂正・明示的自己申告・技術判断から知識候補を抽出できる。
- **Source:** Issue #175 §22–§24
- **Rationale:** ユーザー知識を継続的に観測するため。
- **Status:** confirmed

### REQ-016

- **Statement:** システムは知識候補を既存 Knowledge と照合し、create、attach-evidence、revise、supersede、または更新しない、を判断できる。
- **Source:** Issue #175 §8.6, §13, §24, §26
- **Rationale:** 重複 Assertion の作成を避け、根拠を蓄積するため。
- **Status:** confirmed

### REQ-017

- **Statement:** システムは評価または AI の説明だけでユーザー知識を更新せず、質問だけを未知の Evidence として登録しない。
- **Source:** Issue #175 §6.4, §6.5, §23, §45
- **Rationale:** 誤った知識状態の推定を防ぐため。
- **Status:** confirmed

### REQ-018

- **Statement:** システムはユーザーによる後日の訂正を強い Evidence として扱い、過去 Evidence を物理削除せずに現在の知識状態を再評価できる。
- **Source:** Issue #175 §25, §45
- **Rationale:** 知識推定の訂正可能性と履歴を保つため。
- **Status:** confirmed

## ワークフロー制御

### REQ-019

- **Statement:** Codex は Article Analysis、Knowledge Search、Reading Value、Knowledge Acquisition、Knowledge Update のワークフローを制御し、成果物を引き渡す。
- **Source:** Issue #175 §7, §8.1, §46
- **Rationale:** 各責務を分離した評価と継続更新を実現するため。
- **Status:** confirmed

### REQ-020

- **Statement:** Reading Value は必要に応じて Article Analysis または Knowledge Search へ再調査を依頼でき、Codex は再調査回数を制御する。
- **Source:** Issue #175 §37, §38
- **Rationale:** 重要 Claim の根拠不足や Evidence 不足を補うため。
- **Status:** confirmed

### REQ-021

- **Statement:** システムは検索停止条件と検索 Budget を適用し、検索過程を Search Trace として記録可能にする。
- **Source:** Issue #175 §16, §43
- **Rationale:** Agentic Search の無限化を防ぎ、誤判定を分析可能にするため。
- **Status:** confirmed

### REQ-022

- **Statement:** 利用者は、Knowledge CLIの起動ごとに絶対パスで指定したローカルSQLite Storeを選択できる。指定がない場合は、既存のOSユーザー設定領域にある既定Storeを使う。
- **Source:** Issue #233
- **Rationale:** Codex workspace sandboxなど、OSのユーザー設定領域へ書き込めない実行環境でも、利用者が書込み可能なStoreを明示して検索・更新を実行できるようにするため。
- **Status:** confirmed
