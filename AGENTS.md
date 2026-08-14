# AI Development Workflow

このリポジトリでは、Planning領域とImplementation領域を分離する。`documents/` 配下の詳細なPlanning規約は、同ディレクトリ以下で引き続き適用される。

## 実装の起点

- 製品コードの実装は、このリポジトリ直下で行う。
- 設計・handoff・要件の正規資料は `documents/docs/` にある。実装時の参照パスは `documents/` を先頭にする。
- Implementation Skill は `.agents/skills/impl-<scope>/SKILL.md` に置く。実装に必要な再利用可能Skillがなければ `implementation-skill-builder` を使う。1 TaskごとのSkill生成は禁止する。
- 承認済みの公開CLI/JSON、保存先・設定・運用仕様を、実装の都合で変更または追加しない。設計変更が必要なら `documents/` 配下のDecision Policyに従う。

## 検証

- Go実装では、承認済み設計に従い `gofmt`、`go test ./...`、`go vet ./...` を実行する。
- SQLite migration、公開CLI入出力、stdout/stderr、exit codeは、実行プロセス境界で検証する。

## 人間の指摘からの継続的改善

人間から成果物、説明、質問、レビュー、進行、または判断の仕方について指摘を受けて対応した後は、毎回、その原因が個別成果物だけにあるのか、AGENTS.md・workflow policy・Skillの規約不足にもあるのかを判定する。

- 再発可能な規約不足なら、`skill-maintainer` を使い、最小の正規ルールと regression scenario を修正・追加する。
- プロダクト固有または一回限りの内容なら、Skillを不必要に一般化せず、対象成果物だけを修正する。
- SkillやAGENTS.mdを修正する必要がないと判断した場合も、その理由を最終報告で簡潔に示す。
