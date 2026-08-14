# 業務規則・不変条件

### BR-001: 読書価値はユーザー固有である

記事の一般的な品質・情報量・AI生成の有無ではなく、特定ユーザーの既存知識に対する有意味な認識利得と Attention Cost により評価する。

- **Source:** Issue #175 §1–§5, §41, §49

### BR-002: 未観測は未知ではない

Knowledge Store に Evidence がない場合、ユーザーが知らないと断定せず `no_evidence` として扱う。初期状態を初心者とみなしてはならない。

- **Source:** Issue #175 §6.2, §18, §45, §49

### BR-003: 推論可能は既知ではない

前提知識から Claim を導出できても、Claim 自体を認識している Evidence がなければ `inferable` とし `known` と同一視しない。

- **Source:** Issue #175 §18, §49

### BR-004: Evidence が知識状態の正規根拠である

Derived User Knowledge State は Evidence から導出する。キャッシュまたは永続化した状態は再評価可能でなければならない。

- **Source:** Issue #175 §6.6, §10.5, §49

### BR-005: 露出は知識獲得ではない

AI の説明、記事評価、要約提示、閲覧のみをユーザー知識の Evidence として扱ってはならない。質問だけも未知の Evidence にしてはならない。

- **Source:** Issue #175 §6.4–§6.5, §23, §45, §49

### BR-006: Evidence の強度を区別する

ユーザー自身の正しい説明・推論・コードは strong、明示的な自己申告は moderate、概念認識のみは weak と扱う。確信度は熟達度ラベルではない。

- **Source:** Issue #175 §19, §21, §23

### BR-007: Relation の存在は理解の根拠ではない

Knowledge Store 内に Relation が保存されていても、ユーザーがその Relation を理解しているとは判定しない。

- **Source:** Issue #175 §10.6

### BR-008: 削除ではなく履歴を保持する

訂正・更新時に過去 Evidence および旧 Knowledge を物理削除せず、現在の知識状態を再評価する。

- **Source:** Issue #175 §13, §25

### BR-009: 読書推奨は根拠を伴う

未知 Claim の数だけで `read_full` とせず、Novelty、Importance、Evidence Quality、Applicability、および Attention Cost を評価する。

- **Source:** Issue #175 §28, §31–§35, §45

### BR-010: Codex と CLI の責務を維持する

Claim の意味理解、検索戦略、知識評価、読書価値判断は Codex が担う。CLI は指定された決定論的な保存・検索・更新のみを担う。

- **Source:** Issue #175 §6.8, §8.7, §14, §40, §49
