# DEC-FEAT-018: Reading Valueは既存検証契約への参照に限定する

- **Status:** decided
- **Level:** L4（Feature受入範囲の調整）
- **Date:** 2026-08-15
- **Decision by:** 利用者承認

## Context

既存Reading Value Workflowは公開HTTP(S) URLだけを通常入力として受け、内部でArticle Analysisを作る。固定Article Analysis／Assessment Mapを直接入力する承認済みのテスト入口はない。通常Skillを変更せず、外部記事にも依存しないFEAT-005のFixtureから、Reading Valueを決定論的に実行することはできない。

## Decision

FEAT-005は、Knowledge Search、Knowledge Acquisition、Knowledge UpdateおよびCLI／Store境界を固定Fixtureで実行・観測する。Scenario A、G、I、JのReading Value期待値は、既存FEAT-003検証契約の必須観測節へ`reading_value_reference`として対応付けるだけとし、Reading Value Workflowを起動しない。

## Consequences

- `skills/reading-value/`を変更せず、通常利用時のコンテキストを増やさない。
- I/Jの推奨の正しさはFEAT-003検証契約が所有し、FEAT-005はその必須観測節への参照完全性だけを確認する。
- FEAT-005のCase ResultとRuntime受入評価はReading Value実行結果を持たない。

## Affected Artifacts

- `docs/features/FEAT-005/requirements.md`
- `docs/features/FEAT-005/design.md`
- `docs/features/FEAT-005/design/scenario-catalog.md`
