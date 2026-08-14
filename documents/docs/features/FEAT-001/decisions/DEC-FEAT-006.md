# DEC-FEAT-006: Knowledge CLI の要求キャンセルと割込み終了

- **Status:** decided
- **Level:** L3
- **論点:** 対話端末からの割込みを、Knowledge CLI のSQLite初期化・照会・更新まで伝播する要求キャンセルとして扱うか。また、そのときの公開出力と終了状態をどう定めるか。

## 現状と、なぜ今決める必要があるか

現在の取得実装はSQLiteの照会へ`context.Context`を渡せるが、CLI入口とcomposition rootで`context.Background()`を新たに作る。そのため、呼出元の要求を途中で取り消す経路がない。`Ctrl-C`はGoランタイムの既定動作でプロセスを終了させるだけであり、照会停止、migration transactionの終了、Storeのclose、出力をCLI契約として制御できない。

この振る舞いは全11操作、既定Storeの初期化、SQLite migration、stdout/stderr、終了状態に共通する。既存のerror JSONと終了コードの契約へ例外または追加を設けるため、L3とする。

## 共通の設計条件

- `main`で要求Contextを1つ作り、CLI境界、composition root、application、domain port、SQLite adapterへ同じContextを渡す。下位層で`context.Background()`を作り直さない。
- 対話端末の割込みは`os.Interrupt`で受け、要求Contextをcancelする。公開CLI option、Store指定、環境変数、設定ファイルは追加しない。
- CLI境界はresponseの先頭byteを書き始める前にContext終了を確認する。終了を観測した時点がresponse開始前なら、実行結果・validation・Store初期化errorの種別にかかわらず中断を優先し、JSONを書かず終了コード130を返す。response開始後に到達した割込みは既出力を取り消さず、選択済みの通常responseを完了する。
- `os.UserConfigDir`や親ディレクトリ作成などContextを受けない初期化APIでは、呼出しの前後でContext終了を確認する。SQLiteの読取・migration・更新はContextを受けるAPIを用いる。
- 更新操作でContextがtransaction中に終了した場合は、成功応答を返さず、transactionはcommitしない。既存のrollback・履歴不変条件を維持する。
- 割込みの公開I/Oは、実バイナリへ割込みを送るprocess境界testで検証する。同期には`integrationtest` build tagだけで含める非公開gateを使い、migration、read、mutationが開始したことを親testが確認してから割込みを送る。このgateは通常ビルドに含めず、公開CLI option・公開環境変数・設定・Store overrideを追加しない。ContextがapplicationとSQLite adapterへ同一のまま到達することはunit testでも検証する。

## 選択肢

### 1. 割込みを協調的にcancelし、標準的な終了状態130で終える（推奨）

最初の`Ctrl-C`で要求Contextをcancelする。通常の成功・error JSONは出力せず、stdout/stderrを空にして終了コード`130`で終了する。

- **利点:** シェルの割込み慣例に従い、`storage_error`へ誤分類しない。JSONの機械可読な成功・失敗契約を、完了した操作に限定できる。SQLite照会とtransactionへキャンセルを伝播できる。
- **欠点:** 既存の共通error envelopeと終了コード1〜4に、割込み時だけの明示的な例外を追加する。cancelの処理中に再度割込みを受けた場合の即時終了方針を実装で定める必要がある。

### 2. 割込みを協調的にcancelし、stderrへ新しいJSON errorを出力する

最初の`Ctrl-C`で要求Contextをcancelし、例えば`{"ok":false,"error":{"code":"canceled",...}}`をstderrへ出力する。終了コードは130または新しい専用コードとする。

- **利点:** すべての異常終了でJSONを受け取れる。
- **欠点:** error codeと終了コードを追加し、呼出側の分岐が増える。端末の割込みを通常のCLI業務エラーと扱うため、シェルの慣例から外れる。

### 3. 現行のOS既定終了を維持する

Contextを入口から渡さず、`Ctrl-C`はプロセス終了に任せる。

- **利点:** 公開I/Oの追加設計が不要で、実装変更がない。
- **欠点:** SQLite処理への協調的なキャンセル、transactionの制御、終了時の観測可能な契約を提供できない。利用者の割込み要求を満たさないため採用しない。

## 決定

利用者承認により、選択肢1を採用する。最初の`Ctrl-C`は要求Contextをcancelし、通常のsuccess／error JSONを出力せず、stdout／stderrを空にして終了コード`130`で終了する。

割込みは保存・照会の失敗ではなく利用者が要求した中断であるため、`storage_error`や新規の業務error JSONに混ぜない。終了コード130だけを中断の観測点とし、既存の操作成功・error JSONの意味を保つ。response開始前にContext終了を観測した場合だけを中断の公開保証の線形化点とし、response開始後の割込みで既出力を巻き戻す契約は置かない。二度目以降の割込みの即時終了を追加の公開契約にはせず、シグナル購読を解除した後はOS既定の扱いへ復帰する。

## 反映する成果物

- [詳細設計](../design.md) の共通CLI error/exit規則、Store初期化、Acceptance / Test Design。
- [CLI入力規約](../design/cli-input-conventions.md) と、共通出力契約を参照する操作資料。
- [architecture](../../../design/architecture.md) と [cross-cutting concerns](../../../design/cross-cutting-concerns.md) の要求Context・signal処理・transaction規約。
- [tasks](../tasks.md) と [implementation handoff](../implementation-handoff.yaml) の共通実装／検証前提。
- Issue #179の受入条件、および通常compositionを使うCLI process境界test。
