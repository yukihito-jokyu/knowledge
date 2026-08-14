# DEC-FEAT-003: 初期 Knowledge CLI の名前付きオプション入力契約

- **Status:** decided
- **Level:** L3
- **論点:** Knowledge CLI の11操作に値を渡す公開入力形式。

## 決定

人間が選択した **1** に従い、入力は操作名の後ろに続く名前付きoptionとする。

```text
knowledge <operation> --<option> <value>
```

- 標準入力へ request JSON を渡さない。
- 配列または複数レコードは、同じ名前付きoptionまたはそのグループを繰り返して表す。
- 具体的なoption名、必須性、繰返し規則、グループ規則は [CLI入力規約](../design/cli-input-conventions.md) と各操作設計で定める。
- 成功 JSON は stdout、error JSON は stderr に出力する既存契約を維持する。

## 根拠と採用理由

- 人間が選択した公開 v1 契約である。
- 出力は機械可読な JSON を維持しつつ、呼び出し側は値をコマンドライン上で読み書きできる。
- JSON文字列をoption値に埋め込まないため、複合入力も入力形式が一貫する。

## 影響範囲

- [DEC-FEAT-002](DEC-FEAT-002.md) の「request JSON を標準入力に渡す」部分を置き換える。
- データモデル、SQLite、検索方式、成功・失敗 JSON のschema、終了コードは変更しない。

## 未決定事項

なし。保存先、設定、運用仕様は決めない。
