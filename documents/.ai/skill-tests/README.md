# Skill regression scenarios

`skill-maintainer` がSkillを修正するときに、再発防止用のシナリオをここへ追加する。

推奨形式:

```yaml
id: SKILLTEST-001
target_skill: feature-design
scenario: 単純な読み取り専用Feature
input_summary: 状態を持たず、既存Architectureの範囲内で完結する
expected:
  - 人間判断を要求しない
  - 不要な状態遷移図を必須化しない
  - feature-design.mdのみを更新する
forbidden:
  - architecture.mdを直接変更する
```

これは自動評価基盤が存在しない場合でも、Skill修正時の手動・AI再実行チェックとして利用する。
