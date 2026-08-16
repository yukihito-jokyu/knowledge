# 横断的関心事

## Error と不確実性

- 検索結果がない場合は失敗または未知ではなく `no_evidence` として表現する。
- Evidence が矛盾して収束しない場合は `uncertain` として表現する。
- 外部記事取得・CLI 実行・保存などの技術エラーは、知識評価状態と区別して呼出側へ伝播する。
- Knowledge CLIの端末割込みは技術errorではなく利用者要求による中断として扱う。最初の`Ctrl-C`をresponse開始前に観測した場合は要求Contextをcancelし、JSONを出力せずstdout／stderrを空にして終了コード130で終了する。response開始後の既出力は取り消さない（DEC-FEAT-006）。

## 観測性

- 検索戦略の過程は Search Trace として記録可能にする。
- Trace は正式な Knowledge Assessment 本文に原則含めず、診断・評価用に分離する。
- 検索、評価、更新の失敗を層別に切り分けられる情報を残す。

## 検証

- 決定論的な CLI の検索・取得・更新を単体で検証する。
- Knowledge Search、Knowledge Acquisition、Knowledge Update、Reading Value は Fixture を用いた層別評価を行う。
- Scenario A〜J は少なくとも一つの検証に対応付ける。
- Go実装の共通チェックは `gofmt`、`go test ./...`、`go vet ./...` とする。`gofmt` の結果が差分なしであること、testとvetが成功することを変更完了の最低条件とする。
- `gofmt`、`go test`、`go vet` はGo標準toolchainに含まれる。任意の第三者lint toolを必須化しない。追加する場合は、対象ルール・導入理由・CI実行環境を別Decisionで固定する。
- domainとapplicationのunit testは対応する内部packageに置く。公開option、JSON、stdout/stderr、exit code、migrationの再実行・rollbackは `test/integration/` と `testdata/fixtures/` を用いてprocess境界で検証する。
- SQLiteを使うtestは、OS標準のユーザー設定ディレクトリ環境をテストごとの一時領域へ隔離し、通常のcompositionから独立した一時DBを接続する。固定の保存先・共有DB・ユーザーの既存Storeに依存してはならない。DB指定を公開CLI option、設定、環境変数として追加しない。
- 通常ビルドのStoreだけはDEC-FEAT-005により `os.UserConfigDir()/knowledge-cli/knowledge.db` に固定する。実プロセス境界の検証では、既定Storeの解決・初回migration・`storage_error`を観測する。
- CLI process entryからSQLite adapterまで同一の要求Contextを渡すことをunit testで検証する。実バイナリへの`Ctrl-C`は、`integrationtest` build tagにだけ含める非公開同期gateでmigration・read・mutationの開始を確認してから送る。通常compositionと隔離したOS設定ディレクトリを用いるprocess境界testで、JSONなし・stdout／stderrなし・終了コード130として検証し、mutation中の中断ではDB事後状態が不変であることを確認する。このgateは通常ビルドと公開CLI契約に含めない。

## 依存管理とmigration

- Go依存は `go.mod` と `go.sum` を正規記録とし、追加・更新後は `go mod tidy` により整合させる。実装成果物にvendor directoryを含めない。
- migration SQL assetは連番を再利用・編集しない。不変の既適用versionを変更する必要がある場合は、後続versionを追加する。
- migration runnerは起動時に、FEAT-001のmigration契約に従って未適用versionだけを順序どおりに適用する。`modernc.org/sqlite` のdriver固有の接続・migration機構を公開契約に露出しない。

## セキュリティとプライバシー

- ユーザーの会話・作業 Evidence は個人知識データとして扱う。
- 初期提供はローカル専用であり、リモート共有・同期・共同利用を扱わない。
- ローカル端末上の保存データを保護する OS 標準のアクセス制御を前提とし、追加の暗号化要件は現時点では定義しない。

## 設定

- 検索 Budget、Top-K、Relation Depth 等の値を公開設定として追加しない。必要性は、該当Featureの要件と公開契約を伴うDecisionで扱う。
- リモート接続用の秘密情報は初期提供では扱わない。
