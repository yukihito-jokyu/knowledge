# Reading Value 検証契約

この文書は、Reading Value Workflowの通常系、停止結果、再調査、公開URL取得境界を、実装者以外が同じ入力と観測点で確認するための受入検証契約である。外部記事サイトの可用性、任意URL取得の汎用セキュリティ基盤、Knowledge Searchの探索品質は検証対象に含めない。

## 共通観測規則

- 各ケースは、会話入力、Article Analysis、Knowledge Searchの応答、Reading Value Assessmentまたは停止結果を分けて記録する。
- 正常ケースでは、最終評価の各Claim IDをAssessment MapのAssessment reference、Trace reference、Article Analysisの`location`まで追跡する。自由文による対応付けは合格根拠にしない。
- 停止ケースでは、Article Analysis、Assessment Map、Reading Value Assessmentのいずれを作成しないかを明示して確認する。停止理由を`no_evidence`、`skip`、利用者の未知へ変換してはならない。
- URL公開性検査の記録と、実際の接続先が検査済みIPに固定された記録は別々に観測する。一方だけでは合格にしない。
- 全ケースを、Cookie、認証情報、既存セッション、フォーム送信、POSTその他の状態変更なしで実行する。

## 正常記事と根拠追跡

### V-001: Claimから最終推奨まで追跡できる

複数の評価可能Claim、各Claimの本文内`location`と`support`、正常なKnowledge AssessmentおよびSearch Traceを持つ公開技術記事を入力にする。

次を確認する。

- Article Analysisの各Claimに、Claim ID、`claim`、`role`、`importance`、`location`、`support`がある。
- 各Claimの`scope`、`target_version`、`temporal_context`は記事に明示がなければ`なし`であり、それ以外では記事の値を記録する。
- 各Claimをこれらの属性を改変せず一件ずつKnowledge Searchへ渡す。記事全文を一括で渡さない。
- Assessment Mapに、各Claim IDごとに正常なAssessment referenceとTrace referenceが一対一である。
- `Why`、`Read`、`Skip`、`Knowledge Gains`で参照するClaimが、Assessment Mapを経由してAssessment、Trace、`location`へ追跡できる。
- Map外のClaim、Assessment、Traceを最終推奨の根拠に使わない。
- 正常な各fixtureのReading Value Assessmentに、`Article`、`Recommendation`、`Assessment Map`、`Why`、`Read`、`Skip`、`Knowledge Gains`、`Reliability`、`Limitations`、`Conclusion`がある。`skip`時の`Read`など、空表または「なし」を許すSectionはReading Value成果物契約に従う。

### V-002: 三つの推奨を根拠と禁止断定で確認する

次の独立した記事・Assessmentのfixtureを用意する。

| 推奨 | 入力条件 | 必須観測 | 禁止される断定 |
| --- | --- | --- | --- |
| `read_full` | 高い認識利得が記事全体に分散する、またはClaim間の構造・因果に全体文脈が必要 | 全文が必要な理由、中心Gain、Reliability、低優先箇所の扱い | 未観測Claimを利用者が未知だと断定すること、記事全体の真偽・品質を断定すること |
| `read_selected` | 重要かつ位置特定可能なGainがあり、残りは既知・低重要度・根拠不足・再包装などで相対的に低価値 | 読むClaim IDと位置、飛ばす位置、各Gain、Why、Reliability | 位置のないGainを読む対象にすること、`no_evidence`を未知へ変換すること |
| `skip` | 主な内容が既知、またはGain候補が低重要度・低適用性・根拠不足でAttention Costを正当化しない | 読まない理由、重要な例外の位置、未知を断定していない説明 | 記事一般の品質順位・真偽を断定すること、未観測を未知とみなすこと |

いずれのfixtureでも`Recommendation`は一つだけであり、根拠はAssessment Map内の正常AssessmentとArticle Claimだけであることを確認する。

さらに、`inferable`を`known`へ、`no_evidence`を利用者の未知へ、`uncertain`を確定Gainへ変換しない状態を含むfixtureを用意する。Assessmentの状態・意味を再判定せず、`uncertain`は必要に応じて限定再調査候補または`Limitations`に残ることを確認する。

## 再調査

### V-003: Claim再利用と再評価

初回Article Analysisと、再調査後のAnalysisに次のClaimを含める。

| 照合結果 | 四属性（`claim`、`scope`、`target_version`、`temporal_context`） | 期待結果 |
| --- | --- | --- |
| retained | すべて同一 | 既存Claim IDとAssessment Mapを再利用する |
| changed | 一つ以上が異なる | 旧Assessmentを無効にし、Knowledge Searchを再実行する |
| added | 新Analysisだけにある | Knowledge Searchを実行し、新しいMap項目を作る |
| removed | 旧Analysisだけにある | Mapと最終推奨から除外する |

さらに、初回後のArticle AnalysisまたはKnowledge Searchへの限定再調査は合計二回までであることを確認する。二回後に不足が残る場合は、三回目を開始せず、不確実性を`Why`または`Limitations`に残す。Article再調査後にchangedまたはadded Claimを追随評価するKnowledge Searchは、そのArticle再調査に含まれ、追加の再調査回数を消費しない。

### V-004: Knowledge Storeを更新しない

正常記事、`article_unavailable`、`article_not_evaluable`、Knowledge Searchの`technical_failure`または`canceled`をそれぞれ実行する前後で、隔離したKnowledge Storeの内容と更新操作記録を比較する。記事URL・本文・Article Claim・Assessment・Trace・Reading Value Assessment・AI説明をEvidenceその他の永続データとして追加・更新せず、更新操作がゼロであることを確認する。

## 停止結果

### V-005: 入力形式エラー

空、複数、相対URL、非HTTP(S) URLを利用者入力としてそれぞれ与える。いずれも接続を開始せず、URL一件の指定を求める入力エラーとなり、Article Analysis、Assessment Map、Reading Value Assessmentを作成しないことを確認する。

### V-006: 記事本文・Claim根拠の不足

中心本文、Claimを支える節、または読者が到達できるClaimの位置のいずれかを欠く記事を入力にする。`article_unavailable`で停止し、Article Analysis、部分Claim、Assessment Map、Reading Value Assessmentを作成しないことを確認する。

### V-007: 評価可能Claimなし

技術記事でない、または評価可能なClaimを一件も抽出できない本文を入力にする。`article_not_evaluable`で停止し、Article Analysis、Assessment Map、Reading Value Assessmentを作成しない。空のClaim一覧、記事要約、一般的な品質評価で代替しないことを確認する。

### V-008: Knowledge Searchの失敗・中断

少なくとも一つのClaimに対し、Knowledge Searchが`technical_failure`または`canceled`を返すfixtureをそれぞれ実行する。各結果をAssessment、知識状態、`skip`その他の推奨へ変換せず、失敗または中断を呼出側へ伝播する。失敗したClaimのAssessmentを作らず、Assessment MapとReading Value Assessmentを作成しないことを確認する。

## 公開URL取得の安全境界

### V-009: 接続前に拒否する初期URLとredirect

絶対HTTP(S) URLとして形式は正しい初期URLとredirect先を、それぞれ次の不適格な条件にしたfixtureを用意する。

- localhost、内部、private、loopback、link-local、unique-local、multicast、unspecified、予約済みその他の公開外宛先
- 非既定port
- userinfo付きauthority
- 解決失敗、候補なし、または一つでも公開外IPを含む名前解決結果
- redirect先の非HTTP(S) scheme、redirect loop、`Location`欠落・解決不能、未検査redirect

各ケースで当該宛先へ接続せず、`article_unavailable`で停止することを確認する。

### V-010: IP固定とDNS再名前解決禁止

公開性検査で承認したIPと、接続時に使用したIPを観測できる取得能力で確認する。検査後に名前解決結果が変化するfixtureでも、接続先が検査済みIP以外にならず、hostnameを再名前解決しないことを確認する。

この固定を保証できないブラウザ、HTTPクライアント、プロキシその他の能力では、ネットワーク接続を始めず`article_unavailable`で停止することを確認する。接続失敗を別IPへの再解決・再接続で回避してはならない。

### V-011: 匿名読み取りと非Claim要素の欠落

正常な公開技術記事を、Cookie、認証、セッション、フォーム送信、POSTその他の状態変更なしで取得する。Claim根拠に無関係なナビゲーション、装飾、関連リンクなどだけが欠けるfixtureでは、具体的な欠落を`content_limitations`に記録して評価を続けられることを確認する。中心本文またはClaim根拠が欠ける場合はV-006を優先する。

## 非変更境界

最終確認ではcandidate diffを検査し、次の差分がないことを確認する。

- Knowledge CLIの公開CLI JSON、option、stdout、stderr、exit code
- SQLite schemaまたはmigration
- `cmd/`、`internal/`配下の実装
- 公開設定、保存先、記事本文の永続化
- `skills/knowledge-search/`の探索戦略、状態分類、Budget、記事取得、推奨機能

この検証契約は、上記の変更を伴わない会話内Markdown Workflowだけを対象にする。
