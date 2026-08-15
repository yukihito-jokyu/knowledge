# DEC-FEAT-019: `search-text` の中断をテスト専用同期gateで再現する

- **Status:** decided
- **Level:** L3（既存テスト専用基盤の利用範囲）
- **Date:** 2026-08-15
- **Decision by:** 利用者承認（選択肢1）

## Context

`FEAT005-X-SEARCH-CANCELED` は、最初の `search-text` をresponse開始前・無出力のexit 130で中断し、停止位置と後続未実行を再現可能に観測する。既存の実プロセス中断同期gateは`get`と`search-temporal`にあり、`search-text`にはない。外部から即時に割り込むだけでは、検索開始前または完了後とのraceになり固定Fixtureにならない。

## Decision

既存の`integrationtest` build tag限定の同期gateを、`search-text`のread段階にも適用できるようにする。これは`test/integration/`からSIGINTを送る時点を同期するテスト専用の内部変更である。

通常バイナリ、公開CLIのoption・JSON・exit code、SQLite schema・migration、通常の検索挙動、公開設定は変更しない。`internal/`の変更は、このtest-only gateの追加に限る。

## Considered Alternatives

### 即時SIGINTを送る

同期点がなく再現性を保証できないため採用しない。

### 中断対象を既存gateのあるoperationへ変える

`search-text`を最初のoperationとするScenario Xの意味とSearch Trace契約を変更するため採用しない。

## Consequences

- Xの中断caseを実プロセス境界で決定論的に検証できる。
- production runtimeの振る舞いと公開契約は増えない。
- FEAT-005の`internal/`非変更境界には、この限定的なtest-only例外を記録する。

## Affected Artifacts

- `docs/features/FEAT-005/design.md`
- `docs/features/FEAT-005/tasks.md`
- `docs/features/FEAT-005/implementation-handoff.yaml`
