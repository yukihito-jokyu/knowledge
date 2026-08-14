# CLI入力規約

## 端末割込み

全操作で最初の`Ctrl-C`は実行中の要求を中断する。CLI境界がresponseの先頭byteを書き始める前に中断を観測した場合、success／error JSONを出力せず、stdout／stderrを空にして終了コード130で終了する。response開始後の割込みは既出力を取り消さない。これは名前付きoptionのvalidationや通常のerror responseとは別の共通終了規約である。詳細は[設計本文](../design.md#cli入力とjson出力)と[DEC-FEAT-006](../decisions/DEC-FEAT-006.md)に従う。

## 目的と根拠

この資料は、11操作へ値を渡す方法を共通化する。公開契約の根拠は [DEC-FEAT-003](../decisions/DEC-FEAT-003.md) である。出力 JSON の共通形式は [設計本文](../design.md) を参照する。

## 基本形

```text
knowledge <operation> --<option> <value>
```

- 入力は名前付きoptionだけで渡し、標準入力のJSONや、JSON文字列をoption値へ埋め込む形式は使わない。
- 値に空白を含むときは、呼び出し側がシェルの規則で一つの値として渡す。例: `--normalized-text '本文を一つの値として渡す'`。
- option名は小文字とハイフンで表す。未知のoption、値を持たないoption、値を二重に指定するoptionは `validation_error`（exit 2）である。
- 単一値optionは一度だけ指定する。繰返し可能かどうかは各操作のI/O表で定める。

## 繰返しデータの表し方

複数の値を必要とする入力は、値を一つずつ受けるoptionを順番に置く。区切り文字付き文字列やJSON文字列を解釈しない。

| 種類 | 開始optionと後続option | 規則 |
| --- | --- | --- |
| Scope | `--scope-key`、`--scope-value` | `--scope-key` が1件を開始し、その後に `--scope-value` をちょうど1回置く。同じScope keyは同一操作内で重複不可。 |
| Concept | `--concept`、`--concept-alias` | `--concept` が1件を開始する。後続する `--concept-alias` は次の `--concept` まで、直前のConceptに属する。Conceptなしの `--concept-alias` は不正。 |
| Assertion Alias | `--alias-kind`、`--alias-value` | `--alias-kind` が1件を開始し、次の `--alias-kind` の前に `--alias-value` をちょうど1回置く。 |
| Relation | `--relation-type`、`--relation-target-kind`、`--relation-target-id` | `--relation-type` が1件を開始し、次の `--relation-type` の前にtargetのkindとIDを各1回置く。 |
| Evidence | `--evidence-kind`、`--evidence-text`、`--evidence-observed-at` | `--evidence-kind` が1件を開始し、次の `--evidence-kind` の前にtextとobserved-atを各1回置く。 |

この順序規則を満たさない、途中のレコード、重複する必須構成要素は `validation_error` とする。操作が許すレコード数と各値のenum・空値・重複・参照先は、各操作のバリデーション設計で定める。

## 時点情報

`--valid-from`、`--valid-until`、`--version-scope`、`--observed-at`、`--last-verified` のいずれかを指定すると時点情報を作る。指定しない時点fieldは `null` である。これらを一つも指定しなければ、時点情報は `null` である。各optionは一度だけ指定できる。

## 仮定と範囲

- shellの引用・エスケープはCLI外の呼び出し規則であり、Knowledge CLIの設定や保存先を追加しない。
- これは入力の表現だけを変える。内部で複数の値をSQL bind用の表現へ変換することは公開入力形式ではない。
