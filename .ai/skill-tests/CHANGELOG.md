# Skill Tests Changelog

## 2026-08-14

- `impl-knowledge-cli`: Goおよびmigration SQLの説明コメントを日本語で端的に書く規約を、英語の説明語を残さない形へ明確化。複数要素の配列・sliceリテラルを1要素1行にする規約と回帰scenarioを追加。
- `impl-knowledge-cli`: テーブル駆動unit testとstatement coverage 100%をSkillスクリプトで強制する規約を追加。回帰scenarioに検査失敗条件を追加。
- `impl-knowledge-cli`: SQLite adapterのSQL文を使う関数の直前に個別の複数行定数として置く規約を追加。
- `impl-knowledge-cli`: 固定個数のIN条件は`?`プレースホルダーへ束縛する規約を追加。
