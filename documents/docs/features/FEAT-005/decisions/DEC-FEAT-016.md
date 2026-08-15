# DEC-FEAT-016: Scenario A〜Jを固定Fixtureと層別oracleで評価する

- **Status:** decided
- **Level:** L2（FEAT-005内の検証構成）
- **Date:** 2026-08-15

## Context

Issue #175はScenario A〜Jによる受入評価を求める。一方、Knowledge Search、Update、Reading Valueは責務と成果物が分かれており、最終推奨だけを比較すると、検索漏れ・Evidence解釈・更新境界・統合判断のどこで失敗したかを判断できない。

## Facts

- FEAT-001は隔離SQLite Storeと実CLIプロセス境界のFixtureを提供する。
- FEAT-002はAssessmentとSearch Traceを分離し、Traceに検索過程と停止理由を残す。
- FEAT-003とFEAT-004は、Article／EpisodeからReading Valueと更新までの成果物・検証契約を提供する。
- 公開CLI、永続schema、各Workflowの業務判断は、それぞれのFeatureで既に固定されている。

## Decision

Scenario A〜Jごとに、固定入力・隔離Store・期待成果物・期待Store差分を持つ受入Fixtureを用意する。各Fixtureは、必要な層だけを単独でも実行でき、End-to-End実行では同じケース識別子と入力を利用する。結果は最終推奨だけでなく、CLI process記録、Search Trace、Candidate／Update Result、Assessment Map、Store差分を層別oracleとして照合する。

Fixtureはテスト専用であり、製品の公開入力・出力・設定・永続化契約にしない。外部サイトの可用性に依存せず、記事本文取得後のArticle Analysis入力とEpisode入力を固定する。

## Considered Alternatives

### 最終Reading Valueだけを比較する

短いが、同じ不正な結論に至った原因をSearch、Update、統合のどこにも帰属できないため採用しない。

### 各Featureの既存単体Fixtureだけを使い、Scenario横断Fixtureを作らない

個別契約の検証には有効だが、同じ根拠がEnd-to-Endで維持されること、Scenario A〜Jの網羅性を証明できないため採用しない。

### 実際の外部記事URLをFixtureにする

記事変更・取得不能・ネットワーク条件で結果が揺れ、受入評価の再現性を失うため採用しない。

## Consequences

- Scenarioごとの入力とoracleを保守する必要がある。
- 既存Featureの契約変更は不要で、各契約をテスト根拠として参照する。
- 新しい製品公開契約、DB、migration、ledgerは作らない。

## Affected Artifacts

- `docs/features/FEAT-005/requirements.md`
- `docs/features/FEAT-005/design.md`
- FEAT-005の実装時に作るテストFixtureおよび検証結果。
