# AI Development Skillset — Draft v0.1

要件定義からImplementation Handoffまでを横展開可能なPlanning領域として固定し、実装領域はプロジェクト・技術スタックに応じてSkillを生成するための叩き台です。

## Design Principles

1. **PlanningとImplementationを分離する**
   - Planning: What / Why / Responsibility / Acceptance / Dependency
   - Implementation: How / File / Symbol / Library / Framework-specific choice

2. **人間は重要なDecision Pointに集中する**
   - L1/L2: AI自走
   - L3: Cross-feature / Contract -> 人間承認
   - L4: Product / Business / Architecture -> 人間と議論

3. **Artifact Ownerを1つにする**
   - 別SkillがOwnerのファイルを直接書き換えない
   - 変更が必要ならOwner Skillへ戻す

4. **詳細設計はFeatureごとにJust-in-Timeで行う**
   - Initial Designで全Featureを詳細設計しない

5. **Implementation Skillは必要に応じて生成する**
   - 1 Task = 1 Skillは禁止
   - 再利用される技術依存の実装判断だけSkill化する

6. **Skill自身も保守対象にする**
   - 問題 -> Root Cause -> Minimal Fix -> Regression Scenario

## Skill Set

### Planning

- `planning-orchestrator`
- `requirements-structuring`
- `feature-planning`
- `initial-design`
- `feature-design`
- `task-breakdown`

### Implementation Factory

- `implementation-skill-builder`

### Maintenance

- `skill-maintainer`

## Workflow

```text
Raw Requirements
      |
      v
requirements-structuring
      |
      v
feature-planning
      |
      v
initial-design
      |
      +------------------------------+
      v                              |
feature-design (one feature)         |
      |                              |
      v                              |
task-breakdown                       |
      |                              |
      v                              |
implementation-handoff.yaml          |
=====================================|=== Planning / Implementation boundary
      |                              |
      v                              |
implementation-skill-builder         |
      |                              |
      v                              |
project-specific impl-* skills       |
      |                              |
      v                              |
Implementation ----------------------+  next feature

skill-maintainer -> repairs workflow/skills when friction or defects appear
```

## Repository Layout

```text
.agents/
└── skills/
    ├── planning-orchestrator/
    │   └── SKILL.md
    ├── requirements-structuring/
    │   └── SKILL.md
    ├── feature-planning/
    │   └── SKILL.md
    ├── initial-design/
    │   └── SKILL.md
    ├── feature-design/
    │   └── SKILL.md
    ├── task-breakdown/
    │   └── SKILL.md
    ├── implementation-skill-builder/
    │   └── SKILL.md
    ├── skill-maintainer/
    │   └── SKILL.md
    └── impl-*/
        └── SKILL.md   # generated later

.ai/
├── workflow/
│   ├── artifact-map.yaml
│   ├── decision-policy.md
│   ├── state.yaml
│   ├── implementation-handoff-schema.yaml
│   └── implementation-skills.yaml
└── skill-tests/
    └── ...

docs/
├── requirements/
├── planning/
├── design/
└── features/
    └── FEAT-XXX/
        ├── requirements.md
        ├── design.md
        ├── decisions/
        ├── design-change-requests/
        ├── tasks.md
        └── implementation-handoff.yaml
```

## Artifact Ownership Summary

| Artifact | Owner |
|---|---|
| `.ai/workflow/state.yaml` | planning-orchestrator |
| `docs/requirements/**` | requirements-structuring |
| `docs/planning/**` | feature-planning |
| `docs/design/**` | initial-design |
| `docs/features/{id}/requirements.md` | feature-design |
| `docs/features/{id}/design.md` | feature-design |
| `docs/features/{id}/tasks.md` | task-breakdown |
| `docs/features/{id}/implementation-handoff.yaml` | task-breakdown |
| `.agents/skills/impl-*` | implementation-skill-builder |
| `.ai/workflow/implementation-skills.yaml` | implementation-skill-builder |

`skill-maintainer` はSkill/Workflow定義の修正権限を持ちますが、Project Artifactを書き換えて不具合を隠す用途には使いません。

## Suggested First Trial

最初から巨大プロジェクト全体へ適用せず、既知の要件を持つ小規模プロジェクトで以下を1周させることを推奨します。

1. 生の要件資料を用意
2. `planning-orchestrator` を起点にする
3. 1 Featureだけ `implementation-handoff.yaml` まで進める
4. `implementation-skill-builder` で必要な `impl-*` Skillを生成
5. 実装して、詰まった箇所を記録
6. `skill-maintainer` でSkillを修正
7. 同じ条件のRegression Scenarioを追加

## Draft Status

これはv0.1の叩き台です。特に実運用で検証すべき点:

- Feature粒度判定の安定性
- L2/L3境界が人間の期待に合うか
- `task-breakdown` がImplementationへ情報を渡しすぎ/不足にならないか
- `implementation-skill-builder` のSkill粒度
- OrchestratorがOwner Skillの責務を侵食しないか
- Skill修正時のRegression Scenarioの実効性
