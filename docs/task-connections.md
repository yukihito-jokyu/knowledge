# Task Connections

## 目的

`docs/task-map.md`に記録されたleafの直接依存を、成果物を作るTaskから後続利用者へ逆引きする。

## 状態

完成。Task Map主台帳の全116 leafを抽出し、直接接続・終端leaf・明示Milestone接続を照合済み。

## Source of Truth

- 直接依存の唯一の正本は[Task Map](task-map.md)である。
- 本文書は派生表示であり、新しい依存判断を追加しない。
- GitHub Issue作成時はplanning baseline CommitのSHAへ固定する。

## 生成・検証規則

- `## タスク台帳`から次のレベル2見出しまでだけを解析する。
- leaf IDの明示・省略・範囲表現だけを展開し、推移的依存と親roll-upは含めない。
- 1行を`Producer leaf × Consumer leaf × 依存時点`の直接接続とする。
- 利用成果物はProducerの主成果物を転記し、consumerはread-onlyで利用する。
- Gate／ReleaseはMilestoneのまま保持し、Task IDを付けない。
- `L3-M2-S2 → L2-M2-S1`はTask Mapに明記された順方向のGate入力として記録する。
- 後続利用者がないleafは終端leaf表へ明記する。
- planning integration ownerだけが更新し、各worktreeはread-onlyで利用する。
- 再生成: `ruby scripts/task_issue_sync.rb --render-connections`

## 直接接続台帳

| Producer leaf | 直接の後続利用者 | 利用成果物 | 依存時点 | 利用方法 | Gate／Release | worktree所有記録への参照 |
| --- | --- | --- | --- | --- | --- | --- |
| L1-M1-S1 | L1-M1-S2 | 全論理レコードのデータ辞書 | 着手 | read-only | — | Task Map「L1-M1内」 |
| L1-M1-S1 | L1-M1-S3 | 全論理レコードのデータ辞書 | 着手 | read-only | — | Task Map「L1-M1内」 |
| L1-M1-S1 | L1-M1-S5 | 全論理レコードのデータ辞書 | 着手 | read-only | — | Task Map「L1-M1内」 |
| L1-M1-S2 | L1-M1-S4 | 関連種別・方向・多重度・参照整合性の制約表 | 着手 | read-only | — | Task Map「L1-M1内」 |
| L1-M1-S2 | L1-M1-S5 | 関連種別・方向・多重度・参照整合性の制約表 | 着手 | read-only | — | Task Map「L1-M1内」 |
| L1-M1-S3 | L1-M1-S4 | Evidence、導出結果、根拠追跡、再計算・無効化の契約 | 着手 | read-only | — | Task Map「L1-M1内」 |
| L1-M1-S3 | L1-M1-S5 | Evidence、導出結果、根拠追跡、再計算・無効化の契約 | 着手 | read-only | — | Task Map「L1-M1内」 |
| L1-M1-S4 | L1-M1-S5 | 更新操作の遷移表と履歴系譜規則 | 着手 | read-only | — | Task Map「L1-M1内」 |
| L1-M1-S4 | L3-M2-S3 | 更新操作の遷移表と履歴系譜規則 | 着手 | read-only | — | Task Map「L1-M1内」 |
| L1-M1-S4 | L5-M1-S5 | 更新操作の遷移表と履歴系譜規則 | 着手 | read-only | — | Task Map「L1-M1内」 |
| L1-M1-S5 | L1-M2-S1 | Issue要件と設計要素の双方向トレーサビリティ表 | 着手 | read-only | — | Task Map「L1-M1内」 |
| L1-M2-S1 | L1-M2-S2 | 公開コマンド台帳 | 着手 | read-only | — | Task Map「L1-M2内」 |
| L1-M2-S1 | L1-M2-S3 | 公開コマンド台帳 | 着手 | read-only | — | Task Map「L1-M2内」 |
| L1-M2-S1 | L1-M2-S4 | 公開コマンド台帳 | 着手 | read-only | — | Task Map「L1-M2内」 |
| L1-M2-S1 | L1-M2-S6 | 公開コマンド台帳 | 着手 | read-only | — | Task Map「L1-M2内」 |
| L1-M2-S1 | L2-M1-S3 | 公開コマンド台帳 | 着手 | read-only | — | Task Map「L1-M2内」 |
| L1-M2-S2 | L1-M2-S4 | 集合取得規則とコマンド別適用表 | 着手 | read-only | — | Task Map「L1-M2内」 |
| L1-M2-S2 | L1-M2-S6 | 集合取得規則とコマンド別適用表 | 着手 | read-only | — | Task Map「L1-M2内」 |
| L1-M2-S3 | L1-M2-S4 | Error台帳と終了コード対応表 | 着手 | read-only | — | Task Map「L1-M2内」 |
| L1-M2-S3 | L1-M2-S6 | Error台帳と終了コード対応表 | 着手 | read-only | — | Task Map「L1-M2内」 |
| L1-M2-S3 | L2-M1-S3 | Error台帳と終了コード対応表 | 着手 | read-only | — | Task Map「L1-M2内」 |
| L1-M2-S4 | L1-M2-S5 | 公開操作のJSON Schema一式 | 着手 | read-only | — | Task Map「L1-M2内」 |
| L1-M2-S4 | L1-M2-S6 | 公開操作のJSON Schema一式 | 着手 | read-only | — | Task Map「L1-M2内」 |
| L1-M2-S5 | L1-M2-S6 | 契約Versioning・互換性方針 | 着手 | read-only | — | Task Map「L1-M2内」 |
| L1-M2-S6 | L2-M3-S9 | Issue要件とCLI契約の双方向対応表 | 着手 | read-only | — | Task Map「L1-M2内」 |
| L1-M2-S6 | L6-M3-S5 | Issue要件とCLI契約の双方向対応表 | 着手 | read-only | — | Task Map「L1-M2内」 |
| L2-M1-S1 | L2-M1-S2 | L2-M2-S1のReview Gate基準を満たす永続化技術ADR | 着手 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S1 | L2-M2-S2 | L2-M2-S1のReview Gate基準を満たす永続化技術ADR | 着手 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S1 | L2-M3-S1 | L2-M2-S1のReview Gate基準を満たす永続化技術ADR | 確定 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S10 | L2-M2-S10 | 論理Aggregateの復元、Evidence・履歴追跡、正本Snapshot列挙・Commit済み差分取得Port | 完了・Merge | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S10 | L2-M2-S5 | 論理Aggregateの復元、Evidence・履歴追跡、正本Snapshot列挙・Commit済み差分取得Port | 完了・Merge | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S10 | L2-M3-S5 | 論理Aggregateの復元、Evidence・履歴追跡、正本Snapshot列挙・Commit済み差分取得Port | 完了・Merge | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S2 | L2-M1-S3 | L2-M2-S1の検索・再構築要件を満たす物理ストレージ詳細設計書 | 着手 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S2 | L2-M1-S4 | L2-M2-S1の検索・再構築要件を満たす物理ストレージ詳細設計書 | 着手 | read-only | L2-M1-S2、L2 Design Freeze Gate | Task Map「L2-M1内」 |
| L2-M1-S2 | L2-M2-S2 | L2-M2-S1の検索・再構築要件を満たす物理ストレージ詳細設計書 | 着手 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S3 | L2-M1-S6 | Commit済み正本変更境界を含む実行時永続化設計書 | 着手 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S3 | L2-M2-S2 | Commit済み正本変更境界を含む実行時永続化設計書 | 着手 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S3 | L2-M2-S4 | Commit済み正本変更境界を含む実行時永続化設計書 | 着手 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S3 | L2-M3-S2 | Commit済み正本変更境界を含む実行時永続化設計書 | 着手 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S4 | L2-M1-S5 | 初回DDL／Migration | 着手 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S5 | L2-M1-S6 | Version管理されたMigration Runner | 着手 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S6 | L2-M1-S10 | 取得・更新が共有する低水準永続化層 | 着手 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S6 | L2-M1-S7 | 取得・更新が共有する低水準永続化層 | 着手 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S6 | L2-M1-S8 | 取得・更新が共有する低水準永続化層 | 着手 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S6 | L2-M1-S9 | 取得・更新が共有する低水準永続化層 | 着手 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S7 | L2-M1-S10 | EvidenceとStateの実行時整合経路 | 着手 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S7 | L2-M1-S8 | EvidenceとStateの実行時整合経路 | 着手 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S7 | L2-M1-S9 | EvidenceとStateの実行時整合経路 | 着手 | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S8 | L2-M2-S10 | 原子的な追加型更新Repository | 完了・Merge | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S8 | L2-M3-S6 | 原子的な追加型更新Repository | 完了・Merge | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S9 | L2-M2-S10 | 非破壊更新Repositoryと履歴系譜 | 完了・Merge | read-only | — | Task Map「L2-M1内」 |
| L2-M1-S9 | L2-M3-S6 | 非破壊更新Repositoryと履歴系譜 | 完了・Merge | read-only | — | Task Map「L2-M1内」 |
| L2-M2-S1 | L2-M1-S1 | 技術非依存の検索・鮮度・整合・再構築要件とReview Gate基準 | 着手 | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S1 | L2-M2-S2 | 技術非依存の検索・鮮度・整合・再構築要件とReview Gate基準 | 着手 | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S10 | L2-M3-S7 | 全検索Providerを統合するIndex Lifecycle実装 | 完了・Merge | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S2 | L2-M2-S3 | 検索・Embedding・同期方式の技術選定ADR | 着手 | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S2 | L2-M3-S1 | 検索・Embedding・同期方式の技術選定ADR | 確定 | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S3 | L2-M2-S4 | 内部Component、Index物理表現、共通Candidateの詳細設計 | 着手 | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S3 | L2-M2-S5 | 内部Component、Index物理表現、共通Candidateの詳細設計 | 着手 | read-only | 着手: L2-M2-S3、L2 Design Freeze Gate | Task Map「L2-M2内」 |
| L2-M2-S3 | L2-M3-S2 | 内部Component、Index物理表現、共通Candidateの詳細設計 | 着手 | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S4 | L2-M2-S10 | Commit境界から世代切替・復旧までのIndex Lifecycle設計 | 着手 | read-only | 着手: L2-M2-S4、L2 Design Freeze Gate | Task Map「L2-M2内」 |
| L2-M2-S4 | L2-M2-S6 | Commit境界から世代切替・復旧までのIndex Lifecycle設計 | 着手 | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S4 | L2-M2-S7 | Commit境界から世代切替・復旧までのIndex Lifecycle設計 | 着手 | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S4 | L2-M2-S8 | Commit境界から世代切替・復旧までのIndex Lifecycle設計 | 着手 | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S4 | L2-M2-S9 | Commit境界から世代切替・復旧までのIndex Lifecycle設計 | 着手 | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S4 | L2-M3-S2 | Commit境界から世代切替・復旧までのIndex Lifecycle設計 | 着手 | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S5 | L2-M2-S6 | Provider共通契約、Filter・順序・Pagination、正本照合Core | 着手 | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S5 | L2-M2-S7 | Provider共通契約、Filter・順序・Pagination、正本照合Core | 着手 | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S5 | L2-M2-S8 | Provider共通契約、Filter・順序・Pagination、正本照合Core | 着手 | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S5 | L2-M2-S9 | Provider共通契約、Filter・順序・Pagination、正本照合Core | 着手 | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S5 | L2-M3-S5 | Provider共通契約、Filter・順序・Pagination、正本照合Core | 完了・Merge | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S6 | L2-M2-S10 | Lexical Indexと文字列候補検索Provider | 完了・Merge | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S6 | L2-M3-S5 | Lexical Indexと文字列候補検索Provider | 完了・Merge | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S7 | L2-M2-S10 | Version追跡可能なEmbedding・Semantic検索Provider | 完了・Merge | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S7 | L2-M3-S5 | Version追跡可能なEmbedding・Semantic検索Provider | 完了・Merge | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S8 | L2-M2-S10 | Relation系構造Indexと候補検索Provider | 完了・Merge | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S8 | L2-M3-S5 | Relation系構造Indexと候補検索Provider | 完了・Merge | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S9 | L2-M2-S10 | 独立物理Indexを必須としないTemporal検索Provider | 完了・Merge | read-only | — | Task Map「L2-M2内」 |
| L2-M2-S9 | L2-M3-S5 | 独立物理Indexを必須としないTemporal検索Provider | 完了・Merge | read-only | — | Task Map「L2-M2内」 |
| L2-M3-S1 | L2-M3-S2 | CLI Runtime、Schema Validator、Build・配布方式の技術選定ADR | 着手 | read-only | — | Task Map「L2-M3内」 |
| L2-M3-S2 | L2-M3-S3 | 全公開CommandのHandler・Port所有表と内部CLI詳細設計 | 着手 | read-only | L2-M3-S2、L2 Design Freeze Gate | Task Map「L2-M3内」 |
| L2-M3-S3 | L2-M3-S4 | Process Lifecycle、Command Interface、Dispatch、help・version入口 | 着手 | read-only | — | Task Map「L2-M3内」 |
| L2-M3-S3 | L2-M3-S8 | Process Lifecycle、Command Interface、Dispatch、help・version入口 | 着手 | read-only | — | Task Map「L2-M3内」 |
| L2-M3-S4 | L2-M3-S5 | JSON入出力、stdout・stderr、内部Error変換、終了コード境界 | 着手 | read-only | — | Task Map「L2-M3内」 |
| L2-M3-S4 | L2-M3-S6 | JSON入出力、stdout・stderr、内部Error変換、終了コード境界 | 着手 | read-only | — | Task Map「L2-M3内」 |
| L2-M3-S4 | L2-M3-S7 | JSON入出力、stdout・stderr、内部Error変換、終了コード境界 | 着手 | read-only | — | Task Map「L2-M3内」 |
| L2-M3-S4 | L2-M3-S8 | JSON入出力、stdout・stderr、内部Error変換、終了コード境界 | 着手 | read-only | — | Task Map「L2-M3内」 |
| L2-M3-S5 | L2-M3-S8 | 全取得・検索Commandの読取りAdapter | 着手 | read-only | — | Task Map「L2-M3内」 |
| L2-M3-S5 | L3-M2-S6 | 全取得・検索Commandの読取りAdapter | 着手 | read-only | — | Task Map「L2-M3内」 |
| L2-M3-S5 | L4-M2-S7 | 全取得・検索Commandの読取りAdapter | 着手 | read-only | — | Task Map「L2-M3内」 |
| L2-M3-S6 | L2-M3-S8 | create・attach-evidence・revise・supersedeの更新Adapter | 着手 | read-only | — | Task Map「L2-M3内」 |
| L2-M3-S6 | L3-M2-S6 | create・attach-evidence・revise・supersedeの更新Adapter | 着手 | read-only | — | Task Map「L2-M3内」 |
| L2-M3-S7 | L2-M3-S8 | Index Lifecycle操作のCLI Adapter | 着手 | read-only | — | Task Map「L2-M3内」 |
| L2-M3-S8 | L2-M3-S9 | 全Adapterを配線したVersion付き実行Artifact | 着手 | read-only | — | Task Map「L2-M3内」 |
| L2-M3-S8 | L3-M2-S6 | 全Adapterを配線したVersion付き実行Artifact | 着手 | read-only | — | Task Map「L2-M3内」 |
| L2-M3-S8 | L4-M2-S7 | 全Adapterを配線したVersion付き実行Artifact | 着手 | read-only | — | Task Map「L2-M3内」 |
| L2-M3-S8 | L6-M3-S5 | 全Adapterを配線したVersion付き実行Artifact | 完了・Merge | read-only | — | Task Map「L2-M3内」 |
| L2-M3-S9 | L3-M2-S6 | 配布Artifactに対するCLI単体Contract Suiteと適合結果 | Release条件 | read-only | L3 ReleaseはL2-M3-S9後 | Task Map「L2-M3内」 |
| L2-M3-S9 | L4-M2-S7 | 配布Artifactに対するCLI単体Contract Suiteと適合結果 | Release条件 | read-only | L4 ReleaseはL2-M3-S9後 | Task Map「L2-M3内」 |
| L2-M3-S9 | L6-M4-S1 | 配布Artifactに対するCLI単体Contract Suiteと適合結果 | 着手 | read-only | — | Task Map「L2-M3内」 |
| L3-M1-S1 | L3-M1-S2 | Episode入力・Evidence採否境界設計書 | 確定 | read-only | — | Task Map「L3-M1内」 |
| L3-M1-S1 | L3-M1-S3 | Episode入力・Evidence採否境界設計書 | 着手 | read-only | — | Task Map「L3-M1内」 |
| L3-M1-S1 | L5-M3-S1 | Episode入力・Evidence採否境界設計書 | 着手 | read-only | — | Task Map「L3-M1内」 |
| L3-M1-S1 | L6-M3-S6 | Episode入力・Evidence採否境界設計書 | 着手 | read-only | — | Task Map「L3-M1内」 |
| L3-M1-S2 | L3-M1-S3 | Candidate正規化・L1 Field写像仕様 | 着手 | read-only | — | Task Map「L3-M1内」 |
| L3-M1-S3 | L3-M2-S1 | Candidate Markdown契約とcanonical例 | 着手 | read-only | — | Task Map「L3-M1内」 |
| L3-M1-S3 | L3-M2-S2 | Candidate Markdown契約とcanonical例 | 着手 | read-only | — | Task Map「L3-M1内」 |
| L3-M1-S3 | L5-M1-S2 | Candidate Markdown契約とcanonical例 | 着手 | read-only | — | Task Map「L3-M1内」 |
| L3-M1-S3 | L5-M3-S1 | Candidate Markdown契約とcanonical例 | 着手 | read-only | — | Task Map「L3-M1内」 |
| L3-M1-S3 | L5-M3-S2 | Candidate Markdown契約とcanonical例 | 着手 | read-only | — | Task Map「L3-M1内」 |
| L3-M1-S3 | L6-M3-S6 | Candidate Markdown契約とcanonical例 | 着手 | read-only | — | Task Map「L3-M1内」 |
| L3-M1-S4 | L6-M4-S2 | 固有構造契約に適合するKnowledge Acquisition Skill package | 着手 | read-only | — | Task Map「L3-M1内」 |
| L3-M2-S1 | L3-M2-S3 | Candidate受入・評価可否境界設計書 | 着手 | read-only | — | Task Map「L3-M2内」 |
| L3-M2-S1 | L5-M3-S3 | Candidate受入・評価可否境界設計書 | 着手 | read-only | — | Task Map「L3-M2内」 |
| L3-M2-S1 | L6-M3-S6 | Candidate受入・評価可否境界設計書 | 着手 | read-only | — | Task Map「L3-M2内」 |
| L3-M2-S2 | L2-M2-S1 | 更新対象探索・意味比較設計書 | Gate入力 | read-only | L2 Design Freeze Gate入力 | Task Map「L3-M2内」 |
| L3-M2-S2 | L3-M2-S3 | 更新対象探索・意味比較設計書 | 着手 | read-only | — | Task Map「L3-M2内」 |
| L3-M2-S2 | L3-M2-S4 | 更新対象探索・意味比較設計書 | 確定 | read-only | — | Task Map「L3-M2内」 |
| L3-M2-S3 | L3-M2-S4 | Update Decision Markdown契約 | 確定 | read-only | — | Task Map「L3-M2内」 |
| L3-M2-S3 | L5-M1-S2 | Update Decision Markdown契約 | 着手 | read-only | — | Task Map「L3-M2内」 |
| L3-M2-S3 | L5-M3-S3 | Update Decision Markdown契約 | 着手 | read-only | — | Task Map「L3-M2内」 |
| L3-M2-S3 | L5-M3-S4 | Update Decision Markdown契約 | 着手 | read-only | — | Task Map「L3-M2内」 |
| L3-M2-S3 | L6-M3-S6 | Update Decision Markdown契約 | 着手 | read-only | — | Task Map「L3-M2内」 |
| L3-M2-S4 | L5-M1-S2 | CLI利用・結果引渡し契約 | 着手 | read-only | — | Task Map「L3-M2内」 |
| L3-M2-S4 | L5-M3-S4 | CLI利用・結果引渡し契約 | 着手 | read-only | — | Task Map「L3-M2内」 |
| L3-M2-S5 | L3-M2-S6 | Mockで固有構造契約を確認済みのKnowledge Update Skill package | 着手 | read-only | — | Task Map「L3-M2内」 |
| L3-M2-S6 | L6-M4-S3 | M2固有の実CLI Component連携 | 着手 | read-only | — | Task Map「L3-M2内」 |
| L4-M1-S1 | L4-M1-S3 | 記事入力・取得方式・Source同一性・解析可否境界設計 | 着手 | read-only | — | Task Map「L4-M1内」 |
| L4-M1-S1 | L4-M1-S4 | 記事入力・取得方式・Source同一性・解析可否境界設計 | 着手 | read-only | — | Task Map「L4-M1内」 |
| L4-M1-S1 | L5-M2-S1 | 記事入力・取得方式・Source同一性・解析可否境界設計 | 着手 | read-only | — | Task Map「L4-M1内」 |
| L4-M1-S1 | L6-M3-S6 | 記事入力・取得方式・Source同一性・解析可否境界設計 | 着手 | read-only | — | Task Map「L4-M1内」 |
| L4-M1-S2 | L4-M1-S3 | Overview・Claim分解／正規化設計 | 着手 | read-only | — | Task Map「L4-M1内」 |
| L4-M1-S2 | L4-M1-S4 | Overview・Claim分解／正規化設計 | 着手 | read-only | — | Task Map「L4-M1内」 |
| L4-M1-S3 | L4-M1-S4 | Location・Support根拠追跡設計 | 着手 | read-only | — | Task Map「L4-M1内」 |
| L4-M1-S4 | L4-M2-S2 | Article Analysis Markdown・局所再分析結果契約 | 着手 | read-only | — | Task Map「L4-M1内」 |
| L4-M1-S4 | L4-M3-S1 | Article Analysis Markdown・局所再分析結果契約 | 着手 | read-only | — | Task Map「L4-M1内」 |
| L4-M1-S4 | L5-M1-S2 | Article Analysis Markdown・局所再分析結果契約 | 着手 | read-only | — | Task Map「L4-M1内」 |
| L4-M1-S4 | L5-M2-S1 | Article Analysis Markdown・局所再分析結果契約 | 着手 | read-only | — | Task Map「L4-M1内」 |
| L4-M1-S4 | L5-M2-S4 | Article Analysis Markdown・局所再分析結果契約 | 着手 | read-only | — | Task Map「L4-M1内」 |
| L4-M1-S4 | L6-M3-S6 | Article Analysis Markdown・局所再分析結果契約 | 着手 | read-only | — | Task Map「L4-M1内」 |
| L4-M1-S5 | L6-M4-S4 | 実行可能なArticle Analysis Skill package | 着手 | read-only | — | Task Map「L4-M1内」 |
| L4-M2-S1 | L2-M2-S1 | 技術非依存な探索能力・遷移・必要取得情報・局所停止理由 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S1 | L4-M2-S2 | 技術非依存な探索能力・遷移・必要取得情報・局所停止理由 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S1 | L4-M2-S3 | 技術非依存な探索能力・遷移・必要取得情報・局所停止理由 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S1 | L4-M2-S4 | 技術非依存な探索能力・遷移・必要取得情報・局所停止理由 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S1 | L4-M2-S5 | 技術非依存な探索能力・遷移・必要取得情報・局所停止理由 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S2 | L4-M2-S3 | Claim受入・検索用部分Claim／variant生成手順 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S2 | L4-M2-S4 | Claim受入・検索用部分Claim／variant生成手順 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S2 | L4-M2-S5 | Claim受入・検索用部分Claim／variant生成手順 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S2 | L5-M2-S2 | Claim受入・検索用部分Claim／variant生成手順 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S2 | L6-M3-S6 | Claim受入・検索用部分Claim／variant生成手順 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S3 | L4-M2-S5 | Knowledge Search固有CLI Mapping・失敗引渡し契約 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S4 | L4-M2-S5 | Evidenceから7状態・Known・Gap・Confidenceを導く手順 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S5 | L4-M3-S1 | Assessment・Trace・実行状態を分離した出力契約 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S5 | L5-M1-S2 | Assessment・Trace・実行状態を分離した出力契約 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S5 | L5-M1-S4 | Assessment・Trace・実行状態を分離した出力契約 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S5 | L5-M2-S3 | Assessment・Trace・実行状態を分離した出力契約 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S5 | L5-M2-S4 | Assessment・Trace・実行状態を分離した出力契約 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S5 | L6-M1-S4 | Assessment・Trace・実行状態を分離した出力契約 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S5 | L6-M3-S6 | Assessment・Trace・実行状態を分離した出力契約 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S6 | L4-M2-S7 | Mockで固有構造契約を確認済みのKnowledge Search Skill package | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M2-S7 | L6-M4-S5 | Knowledge Search固有の実CLI Component連携 | 着手 | read-only | — | Task Map「L4-M2内」 |
| L4-M3-S1 | L4-M3-S2 | M1・M2入力のClaim整合・評価可否境界 | 着手 | read-only | — | Task Map「L4-M3内」 |
| L4-M3-S1 | L4-M3-S3 | M1・M2入力のClaim整合・評価可否境界 | 着手 | read-only | — | Task Map「L4-M3内」 |
| L4-M3-S1 | L4-M3-S5 | M1・M2入力のClaim整合・評価可否境界 | 着手 | read-only | — | Task Map「L4-M3内」 |
| L4-M3-S1 | L5-M2-S3 | M1・M2入力のClaim整合・評価可否境界 | 着手 | read-only | — | Task Map「L4-M3内」 |
| L4-M3-S1 | L6-M3-S6 | M1・M2入力のClaim整合・評価可否境界 | 着手 | read-only | — | Task Map「L4-M3内」 |
| L4-M3-S2 | L4-M3-S4 | 7状態を再判定しないRecognition Gain適用手順 | 着手 | read-only | — | Task Map「L4-M3内」 |
| L4-M3-S3 | L4-M3-S4 | Supportに追跡可能なReliability・Applicability判断手順 | 着手 | read-only | — | Task Map「L4-M3内」 |
| L4-M3-S4 | L4-M3-S5 | 3推奨と読書範囲のClaim横断統合手順 | 着手 | read-only | — | Task Map「L4-M3内」 |
| L4-M3-S5 | L5-M1-S2 | 最終Assessment／追加調査要求の排他的出力契約 | 着手 | read-only | — | Task Map「L4-M3内」 |
| L4-M3-S5 | L5-M2-S4 | 最終Assessment／追加調査要求の排他的出力契約 | 着手 | read-only | — | Task Map「L4-M3内」 |
| L4-M3-S5 | L5-M2-S5 | 最終Assessment／追加調査要求の排他的出力契約 | 着手 | read-only | — | Task Map「L4-M3内」 |
| L4-M3-S5 | L6-M3-S6 | 最終Assessment／追加調査要求の排他的出力契約 | 着手 | read-only | — | Task Map「L4-M3内」 |
| L4-M3-S6 | L6-M4-S6 | 実行可能なReading Value Skill package | 着手 | read-only | — | Task Map「L4-M3内」 |
| L5-M1-S1 | L5-M1-S3 | 外部受付、Workflow識別、Run初期化要求のumbrella契約 | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M1-S1 | L5-M1-S4 | 外部受付、Workflow識別、Run初期化要求のumbrella契約 | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M1-S2 | L5-M1-S3 | 子Skill識別・Version・Invocation・成果物引渡し契約 | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M1-S2 | L5-M1-S4 | 子Skill識別・Version・Invocation・成果物引渡し契約 | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M1-S3 | L5-M1-S5 | producer意味状態を保持する共通実行状態・診断Envelope | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M1-S3 | L5-M2-S5 | producer意味状態を保持する共通実行状態・診断Envelope | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M1-S3 | L5-M3-S2 | producer意味状態を保持する共通実行状態・診断Envelope | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M1-S3 | L5-M3-S5 | producer意味状態を保持する共通実行状態・診断Envelope | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M1-S4 | L5-M1-S5 | ID Owner、成果物系譜、Search Trace相関契約 | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M1-S4 | L5-M2-S4 | ID Owner、成果物系譜、Search Trace相関契約 | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M1-S4 | L5-M3-S3 | ID Owner、成果物系譜、Search Trace相関契約 | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M1-S4 | L5-M3-S4 | ID Owner、成果物系譜、Search Trace相関契約 | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M1-S4 | L6-M1-S4 | ID Owner、成果物系譜、Search Trace相関契約 | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M1-S5 | L5-M2-S4 | Budget、再実行判定、停止理由、副作用安全性契約 | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M1-S5 | L5-M2-S5 | Budget、再実行判定、停止理由、副作用安全性契約 | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M1-S5 | L5-M3-S2 | Budget、再実行判定、停止理由、副作用安全性契約 | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M1-S5 | L5-M3-S4 | Budget、再実行判定、停止理由、副作用安全性契約 | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M1-S5 | L5-M3-S5 | Budget、再実行判定、停止理由、副作用安全性契約 | 着手 | read-only | — | Task Map「L5-M1内」 |
| L5-M2-S1 | L5-M2-S2 | 記事Workflow開始・Article Analysis起動／成果物受入 | 着手 | read-only | — | Task Map「L5-M2内」 |
| L5-M2-S1 | L5-M2-S5 | 記事Workflow開始・Article Analysis起動／成果物受入 | 着手 | read-only | — | Task Map「L5-M2内」 |
| L5-M2-S2 | L5-M2-S3 | canonical Claimから検索Work Itemへの一対多展開契約 | 着手 | read-only | — | Task Map「L5-M2内」 |
| L5-M2-S2 | L5-M2-S5 | canonical Claimから検索Work Itemへの一対多展開契約 | 着手 | read-only | — | Task Map「L5-M2内」 |
| L5-M2-S3 | L5-M2-S4 | Claim別結果の機械的fan-in・Reading Value入力Bundle | 着手 | read-only | — | Task Map「L5-M2内」 |
| L5-M2-S3 | L5-M2-S5 | Claim別結果の機械的fan-in・Reading Value入力Bundle | 着手 | read-only | — | Task Map「L5-M2内」 |
| L5-M2-S4 | L5-M2-S5 | 許可済み追加調査要求の固有Routing・Cycle反映契約 | 着手 | read-only | — | Task Map「L5-M2内」 |
| L5-M3-S1 | L5-M3-S2 | Knowledge蓄積Workflow開始・Acquisition起動／成果物受入 | 着手 | read-only | — | Task Map「L5-M3内」 |
| L5-M3-S1 | L5-M3-S5 | Knowledge蓄積Workflow開始・Acquisition起動／成果物受入 | 着手 | read-only | — | Task Map「L5-M3内」 |
| L5-M3-S2 | L5-M3-S3 | Acquisition結果を意味変更しない3分岐契約 | 着手 | read-only | — | Task Map「L5-M3内」 |
| L5-M3-S2 | L5-M3-S5 | Acquisition結果を意味変更しない3分岐契約 | 着手 | read-only | — | Task Map「L5-M3内」 |
| L5-M3-S3 | L5-M3-S4 | Candidate追跡単位・batch／直列既定のUpdate起動契約 | 着手 | read-only | — | Task Map「L5-M3内」 |
| L5-M3-S3 | L5-M3-S5 | Candidate追跡単位・batch／直列既定のUpdate起動契約 | 着手 | read-only | — | Task Map「L5-M3内」 |
| L5-M3-S4 | L5-M3-S5 | Candidate別結果集約・pending集合・安全な再開順序 | 着手 | read-only | — | Task Map「L5-M3内」 |
| L5-M4-S1 | L5-M4-S2 | root Parent Skill、共通制御、Registry実体 | 完了・Merge | read-only | — | Task Map「L5-M4内」 |
| L5-M4-S1 | L5-M4-S3 | root Parent Skill、共通制御、Registry実体 | 完了・Merge | read-only | — | Task Map「L5-M4内」 |
| L5-M4-S1 | L5-M4-S4 | root Parent Skill、共通制御、Registry実体 | 着手 | read-only | — | Task Map「L5-M4内」 |
| L5-M4-S2 | L5-M4-S4 | 記事Workflow実行reference | 着手 | read-only | — | Task Map「L5-M4内」 |
| L5-M4-S3 | L5-M4-S4 | Knowledge蓄積Workflow実行reference | 着手 | read-only | — | Task Map「L5-M4内」 |
| L5-M4-S4 | L5-M4-S5 | Parent固有Mock・決定論的契約確認 | 着手 | read-only | L5-M4-S4、L3 Release | Task Map「L5-M4内」 |
| L5-M4-S4 | L5-M4-S6 | Parent固有Mock・決定論的契約確認 | 着手 | read-only | L5-M4-S4、L4 Release | Task Map「L5-M4内」 |
| L5-M4-S4 | L6-M3-S7 | Parent固有Mock・決定論的契約確認 | 完了・Merge | read-only | — | Task Map「L5-M4内」 |
| L5-M4-S5 | L6-M5-S1 | L3実SkillsとのKnowledge蓄積Component連携 | 着手 | read-only | — | Task Map「L5-M4内」 |
| L5-M4-S6 | L6-M5-S2 | L4実Skillsとの記事価値判定Component連携 | 着手 | read-only | — | Task Map「L5-M4内」 |
| L6-M1-S1 | L6-M1-S2 | 評価レイヤ・対象・非対象・Owner方針 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S1 | L6-M1-S4 | 評価レイヤ・対象・非対象・Owner方針 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S1 | L6-M2-S1 | 評価レイヤ・対象・非対象・Owner方針 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S1 | L6-M3-S1 | 評価レイヤ・対象・非対象・Owner方針 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S2 | L6-M1-S3 | Requirement―評価レイヤ―観測―Owner―Oracle種別の追跡表 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S2 | L6-M2-S1 | Requirement―評価レイヤ―観測―Owner―Oracle種別の追跡表 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S2 | L6-M4-S1 | Requirement―評価レイヤ―観測―Owner―Oracle種別の追跡表 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S2 | L6-M4-S2 | Requirement―評価レイヤ―観測―Owner―Oracle種別の追跡表 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S2 | L6-M4-S3 | Requirement―評価レイヤ―観測―Owner―Oracle種別の追跡表 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S2 | L6-M4-S4 | Requirement―評価レイヤ―観測―Owner―Oracle種別の追跡表 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S2 | L6-M4-S5 | Requirement―評価レイヤ―観測―Owner―Oracle種別の追跡表 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S2 | L6-M4-S6 | Requirement―評価レイヤ―観測―Owner―Oracle種別の追跡表 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S3 | L6-M1-S5 | 決定論的／Agent評価の判定・反復基準 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S3 | L6-M2-S2 | 決定論的／Agent評価の判定・反復基準 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S3 | L6-M3-S1 | 決定論的／Agent評価の判定・反復基準 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S3 | L6-M4-S1 | 決定論的／Agent評価の判定・反復基準 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S3 | L6-M4-S2 | 決定論的／Agent評価の判定・反復基準 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S3 | L6-M4-S3 | 決定論的／Agent評価の判定・反復基準 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S3 | L6-M4-S4 | 決定論的／Agent評価の判定・反復基準 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S3 | L6-M4-S5 | 決定論的／Agent評価の判定・反復基準 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S3 | L6-M4-S6 | 決定論的／Agent評価の判定・反復基準 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S4 | L6-M1-S5 | 既存Trace・相関情報からの失敗原因診断手順 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S4 | L6-M3-S3 | 既存Trace・相関情報からの失敗原因診断手順 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S5 | L6-M3-S3 | Caseから最終Reportまでの共通Envelope・集約契約 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S5 | L6-M4-S1 | Caseから最終Reportまでの共通Envelope・集約契約 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S5 | L6-M4-S2 | Caseから最終Reportまでの共通Envelope・集約契約 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S5 | L6-M4-S3 | Caseから最終Reportまでの共通Envelope・集約契約 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S5 | L6-M4-S4 | Caseから最終Reportまでの共通Envelope・集約契約 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S5 | L6-M4-S5 | Caseから最終Reportまでの共通Envelope・集約契約 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M1-S5 | L6-M4-S6 | Caseから最終Reportまでの共通Envelope・集約契約 | 着手 | read-only | — | Task Map「L6-M1内」 |
| L6-M2-S1 | L6-M2-S2 | Dataset／Fixture共通Schema・互換性契約 | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S1 | L6-M2-S3 | Dataset／Fixture共通Schema・互換性契約 | 着手 | read-only | L6 Evaluation Design Freeze Gate、L6-M2-S1・S2 | Task Map「L6-M2内」 |
| L6-M2-S1 | L6-M3-S3 | Dataset／Fixture共通Schema・互換性契約 | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S2 | L6-M2-S3 | 全Scenario・Invariantの具体Case Catalog・期待Oracle | 着手 | read-only | L6 Evaluation Design Freeze Gate、L6-M2-S1・S2 | Task Map「L6-M2内」 |
| L6-M2-S3 | L6-M2-S4 | Manifest・Schema・参照・Coverage Validator | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S3 | L6-M2-S5 | Manifest・Schema・参照・Coverage Validator | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S3 | L6-M2-S6 | Manifest・Schema・参照・Coverage Validator | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S3 | L6-M4-S1 | Manifest・Schema・参照・Coverage Validator | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S3 | L6-M4-S2 | Manifest・Schema・参照・Coverage Validator | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S3 | L6-M4-S3 | Manifest・Schema・参照・Coverage Validator | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S3 | L6-M4-S4 | Manifest・Schema・参照・Coverage Validator | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S3 | L6-M4-S5 | Manifest・Schema・参照・Coverage Validator | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S3 | L6-M4-S6 | Manifest・Schema・参照・Coverage Validator | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S4 | L6-M4-S1 | A〜E・検索／State領域の初期State・Evidence・期待Assessment Fixture | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S4 | L6-M4-S5 | A〜E・検索／State領域の初期State・Evidence・期待Assessment Fixture | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S4 | L6-M4-S6 | A〜E・検索／State領域の初期State・Evidence・期待Assessment Fixture | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S4 | L6-M5-S2 | A〜E・検索／State領域の初期State・Evidence・期待Assessment Fixture | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S5 | L6-M4-S1 | F〜H・蓄積領域のEpisode・Evidence・更新前後State Fixture | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S5 | L6-M4-S2 | F〜H・蓄積領域のEpisode・Evidence・更新前後State Fixture | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S5 | L6-M4-S3 | F〜H・蓄積領域のEpisode・Evidence・更新前後State Fixture | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S5 | L6-M5-S1 | F〜H・蓄積領域のEpisode・Evidence・更新前後State Fixture | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S6 | L6-M4-S4 | I〜J・Reading Value領域の記事・Knowledge・期待推奨Fixture | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S6 | L6-M4-S6 | I〜J・Reading Value領域の記事・Knowledge・期待推奨Fixture | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M2-S6 | L6-M5-S2 | I〜J・Reading Value領域の記事・Knowledge・期待推奨Fixture | 着手 | read-only | — | Task Map「L6-M2内」 |
| L6-M3-S1 | L6-M3-S2 | 技術非依存なHarness能力・再現性・隔離要件 | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S2 | L6-M3-S3 | Harness技術選定ADR | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S3 | L6-M3-S4 | Harness Interface・データフロー・失敗境界設計 | 着手 | read-only | L6 Evaluation Design Freeze Gate、L6-M3-S3 | Task Map「L6-M3内」 |
| L6-M3-S4 | L6-M3-S5 | 共通Harness Core・Loader・判定・Report Pipeline | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S4 | L6-M3-S6 | 共通Harness Core・Loader・判定・Report Pipeline | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S4 | L6-M3-S7 | 共通Harness Core・Loader・判定・Report Pipeline | 着手 | read-only | 着手: L6-M3-S4、L5 Design Freeze Gate | Task Map「L6-M3内」 |
| L6-M3-S4 | L6-M4-S1 | 共通Harness Core・Loader・判定・Report Pipeline | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S4 | L6-M4-S2 | 共通Harness Core・Loader・判定・Report Pipeline | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S4 | L6-M4-S3 | 共通Harness Core・Loader・判定・Report Pipeline | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S4 | L6-M4-S4 | 共通Harness Core・Loader・判定・Report Pipeline | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S4 | L6-M4-S5 | 共通Harness Core・Loader・判定・Report Pipeline | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S4 | L6-M4-S6 | 共通Harness Core・Loader・判定・Report Pipeline | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S5 | L6-M4-S1 | CLI Process・JSON・Store隔離Adapter | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S5 | L6-M4-S3 | CLI Process・JSON・Store隔離Adapter | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S5 | L6-M4-S5 | CLI Process・JSON・Store隔離Adapter | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S5 | L6-M5-S1 | CLI Process・JSON・Store隔離Adapter | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S5 | L6-M5-S2 | CLI Process・JSON・Store隔離Adapter | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S6 | L6-M4-S2 | Component Skill起動・Markdown回収Adapter | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S6 | L6-M4-S3 | Component Skill起動・Markdown回収Adapter | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S6 | L6-M4-S4 | Component Skill起動・Markdown回収Adapter | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S6 | L6-M4-S5 | Component Skill起動・Markdown回収Adapter | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S6 | L6-M4-S6 | Component Skill起動・Markdown回収Adapter | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S7 | L6-M5-S1 | Parent entrypoint・Run／Cycle・終端回収Adapter | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M3-S7 | L6-M5-S2 | Parent entrypoint・Run／Cycle・終端回収Adapter | 着手 | read-only | — | Task Map「L6-M3内」 |
| L6-M4-S1 | L6-M4-S3 | CLI横断評価Suite・診断Report | Report確定 | read-only | — | Task Map「L6-M4内」 |
| L6-M4-S1 | L6-M4-S5 | CLI横断評価Suite・診断Report | Report確定 | read-only | — | Task Map「L6-M4内」 |
| L6-M4-S1 | L6-M5-S1 | CLI横断評価Suite・診断Report | Report確定 | read-only | — | Task Map「L6-M4内」 |
| L6-M4-S1 | L6-M5-S2 | CLI横断評価Suite・診断Report | Report確定 | read-only | — | Task Map「L6-M4内」 |
| L6-M4-S2 | L6-M5-S1 | Acquisition Agent評価Suite・診断Report | Report確定 | read-only | — | Task Map「L6-M4内」 |
| L6-M4-S3 | L6-M5-S1 | Update Agent評価Suite・診断Report | Report確定 | read-only | — | Task Map「L6-M4内」 |
| L6-M4-S4 | L6-M5-S2 | Article Analysis Agent評価Suite・診断Report | Report確定 | read-only | — | Task Map「L6-M4内」 |
| L6-M4-S5 | L6-M5-S2 | Knowledge Search Agentic評価Suite・診断Report | Report確定 | read-only | — | Task Map「L6-M4内」 |
| L6-M4-S6 | L6-M5-S2 | Reading Value Agent評価Suite・診断Report | Report確定 | read-only | — | Task Map「L6-M4内」 |
| L6-M5-S1 | L6-M5-S3 | Knowledge蓄積Workflow E2E Suite・Report | 着手 | read-only | — | Task Map「L6-M5内」 |
| L6-M5-S2 | L6-M5-S3 | 記事価値判定Workflow E2E Suite・Report | 着手 | read-only | — | Task Map「L6-M5内」 |

### 終端leaf

| Leaf Task | 後続利用者 | 主成果物 |
| --- | --- | --- |
| L5-M2-S5 | なし（終端Task） | 最終Assessment受入、異常・部分・Budget終端契約 |
| L5-M3-S5 | なし（終端Task） | 正常・部分・異常・Budget終端と未処理範囲 |
| L6-M5-S3 | なし（終端Task） | 全Scenario回帰Suite・最終Acceptance／診断Report |

### 明示されたGate／Release接続

| Milestone表現 | Consumer leaf | 依存時点 |
| --- | --- | --- |
| L2 Design Freeze Gate | L2-M1-S4 | 着手 |
| 着手: L2 Design Freeze Gate | L2-M2-S5 | 着手 |
| 着手: L2 Design Freeze Gate | L2-M2-S10 | 着手 |
| L2 Design Freeze Gate | L2-M3-S3 | 着手 |
| L3 Skill Design Freeze Gate | L3-M1-S4 | 着手 |
| L3 Skill Design Freeze Gate | L3-M2-S5 | 着手 |
| L3 Releaseは後 | L3-M2-S6 | Release条件 |
| L4 Skill Design Freeze Gate | L4-M1-S5 | 着手 |
| L4 Skill Design Freeze Gate | L4-M2-S6 | 着手 |
| L4 Releaseは後 | L4-M2-S7 | Release条件 |
| L4 Skill Design Freeze Gate | L4-M3-S6 | 着手 |
| L5 Design Freeze Gate | L5-M4-S1 | 着手 |
| L5 Design Freeze Gate | L5-M4-S2 | 着手 |
| L5 Design Freeze Gate | L5-M4-S3 | 着手 |
| L3 Release | L5-M4-S5 | 着手 |
| L4 Release | L5-M4-S6 | 着手 |
| L6 Evaluation Design Freeze Gate、 | L6-M2-S3 | 着手 |
| L6 Evaluation Design Freeze Gate、 | L6-M3-S4 | 着手 |
| 着手: L5 Design Freeze Gate | L6-M3-S7 | 着手 |
