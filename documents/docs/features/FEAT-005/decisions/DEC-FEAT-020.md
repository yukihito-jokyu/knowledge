# DEC-FEAT-020: Evidence単位のTemporal情報の保存方法

- **Status:** decided
- **Level:** L3（共有SQLite永続化形式）
- **Date:** 2026-08-16
- **Decision by:** 利用者承認（選択肢A）

## Context

`FEAT005-E-SEARCH-OUTDATED` は、旧理解のEvidence `ev-e-old` を
`version_scope=go1.21`、訂正Evidence `ev-e-correction` を
`version_scope=go1.22+` として固定Storeへ再現する。

現行SQLite schemaでは、`temporal_metadata` はAssertion revision単位であり、
`evidence` はEvidence単位の`version_scope`を持たない。このため、承認済み
Scenario catalogにある二つのEvidenceの異なるTemporal情報を、実Storeへ忠実に
保存・検索・検証できない。

## Decision

選択肢Aを採用する。Evidenceに紐付くTemporal metadataをSQLite schemaと
domain/store契約へ追加する。

以下は実装時の不変条件とする。

- EvidenceのTemporalは任意であり、既存Evidenceと既存CLI呼出しは互換に読む。
- `version_scope`を含むTemporal値はEvidence IDに一意に紐付ける。
- migrationは既存Storeを破壊せず、再起動可能である。
- `get-evidence`の公開JSONには、Evidence Temporalを省略可能な`temporal`として
  表示する。既存の値なしEvidenceのresponse形状は維持する。

## Considered Alternatives

次のいずれかを人間が決定する。

### A. Evidence Temporalを永続化する（採用）

Evidenceに紐付くTemporal metadataをSQLite schemaとdomain/store契約へ追加する。

- 利点: catalogの固定seedとoracleを変更せず、Evidence単位の時点判断を将来の
  Searchにも利用できる。
- 影響: migration、domain/store、CLI read response、fixture、複数Featureの
  永続化契約。互換性を伴うL3変更となる。

### B. TemporalはAssertion revision単位へ正規化する

`FEAT005-E` のcatalogを変更し、Evidence個別の`version_scope`をAssertion
revisionのTemporalへ移す。

- 利点: schema変更が不要。
- 影響: Evidence単位の時点根拠を保持しないという意味変更になるため、Eの
  Assessment/Trace oracleと関連Feature設計の再承認が必要。

### C. Evidence TemporalをFEAT-005受入対象から外す

catalogからEvidence個別のTemporal要求を削除する。

- 利点: 最小の実装範囲でFixtureを完成できる。
- 影響: `outdated`判定の根拠が弱くなり、受入Scenarioの意味を縮小する。

## Consequences

`FEAT005-E` の固定Store seedと期待oracleを実DBで再現できる。migration、
domain/store、CLI read response、fixture、実プロセス検証の実装が必要となる。

## Affected Artifacts

- `docs/features/FEAT-005/design/scenario-catalog.md`
- `docs/features/FEAT-005/design.md`
- `docs/features/FEAT-005/tasks.md`
- SQLite schema / migration / domain Store contract（選択肢Aの場合）
