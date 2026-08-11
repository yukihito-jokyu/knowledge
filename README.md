# knowledge

## Task Issue同期スクリプト

`scripts/task_issue_sync.sh`は、製品本体ではなく、`docs/task-map.md`とGitHubの作業項目（Task Issue）を同じ計画内容に保つための管理用スクリプトです。Task Mapを変更したとき、接続台帳を再生成するとき、またはGitHub Issueとの不一致を確認・反映するときに使います。

### 必要な環境

- Bash 3.2以上: 処理順と、異常を検出した場合の停止を制御するために必要です。
- `jq` 1.6以上: GitHubが返す構造化データ形式（JSON）を、改行や文書記号を壊さず処理するために必要です。
- `rtk`と`git`: プロジェクト共通の実行入口を通して、Gitの変更履歴へ固定したTask Mapを読み取るために必要です。
- `gh`: GitHubをコマンドで操作し、Issueを照合または更新する場合だけ必要です。

### 使う場面とコマンド

#### 1. 固定したTask Mapを検査する

Planning snapshot SHAは、計画を固定したGitの変更履歴を表す40桁の識別値です。Task数、承認状態、親子関係、成果物の保存先（Path）の重複を検査します。GitHubへは接続しません。

```bash
rtk bash scripts/task_issue_sync.sh --check --snapshot <40桁のPlanning snapshot SHA>
```

#### 2. 接続台帳を再生成する

現在の`docs/task-map.md`から、成果物を作るTaskと、その成果物を直接利用する後続Taskの一覧を生成します。GitHubへは接続しません。

```bash
rtk bash scripts/task_issue_sync.sh --render-connections
```

出力内容を確認したうえで、`docs/task-connections.md`の「直接接続台帳」以降へ反映します。

#### 3. GitHub Issueとの一致を確認する

固定したTask Mapと既存のGitHub Task Issueを照合します。GitHubから読み取りますが、Issueは変更しません。

```bash
rtk bash scripts/task_issue_sync.sh --verify \
  --snapshot <40桁のPlanning snapshot SHA> \
  --repo <OWNER/REPOSITORY>
```

`OWNER`にはGitHubの利用者名または組織名、`REPOSITORY`にはリポジトリ名を指定します。`--repo`を省略した場合は`yukihito-jokyu/knowledge`を使用します。

#### 4. GitHub Issueへ反映する

不足するTask Issueの作成と、既存Issueの自動生成領域の更新を行います。人手で記録した進捗欄は保持し、Issueの自動Close・自動Reopen・自動削除は行いません。この操作だけはGitHubへ書き込むため、利用者が明示的に承認した場合に限り実行します。

```bash
rtk bash scripts/task_issue_sync.sh --apply \
  --snapshot <40桁のPlanning snapshot SHA> \
  --repo <OWNER/REPOSITORY>
```

実行順は`--check`、`--render-connections`、`--verify`、承認後の`--apply`です。事前検査で異常を検出した場合は、GitHubへ書き込まずに停止するため、原因を確認してください。
