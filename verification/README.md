# Knowledge CLI / Skills 検証環境

`bin/knowledge` は、このリポジトリの現在のKnowledge CLIをビルドした実行バイナリです。Knowledge Storeの保存先は既存CLIの契約どおりOSのユーザー設定領域です。検証用のStoreを使う場合は、既存CLIが参照するユーザー設定領域を事前に隔離してください。

## URLを評価してKnowledgeを更新するSkill

使うSkillは **`reading-value`** です。これは次のSkillを順に呼びます。

1. `reading-value`: URLの記事を読む価値を評価する。
2. `knowledge-acquisition`: 同じURL評価中のユーザー由来の技術的回答だけをCandidate Knowledgeにする。
3. `knowledge-update`: Candidateを既存Knowledgeと照合し、必要な更新をこの環境のCLIで行う。

Codexへは `reading-value` を選択してから、**一件の絶対HTTP(S) URLだけ**を渡してください。例:

```text
$reading-value https://example.com/technical-article
```

`reading-value` の入力契約はURL一件だけです。Skill名を本文へ混ぜたり、複数URLを同時に渡したりしません。技術的な説明、訂正、コード、自己申告など、同じURL評価Episode内でユーザーが示した内容だけが更新候補になります。記事本文、AIの説明、質問だけの発話は保存されません。

## CLIの確認

```sh
./bin/knowledge search-text --query '確認したい知識'
```

Skillのコピーは `.agents/skills/` にあります。Knowledge Updateが呼ぶ `knowledge` をこのディレクトリの実行バイナリへ解決するには、実行環境のPATH先頭へ `verification/bin` を追加してください。

```sh
export PATH="$PWD/bin:$PATH"
```

この状態で `knowledge-update` が呼ぶ `knowledge` は `bin/knowledge` を使います。
