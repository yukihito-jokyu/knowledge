#!/usr/bin/env ruby
# frozen_string_literal: true

require "json"
require "open3"
require "optparse"

ROOT = File.expand_path("..", __dir__)
TASK_MAP = File.join(ROOT, "docs", "task-map.md")
EXPECTED_COUNTS = { "L" => 6, "M" => 19, "S" => 116 }.freeze
REPO_DEFAULT = "yukihito-jokyu/knowledge"
ROOT_ISSUE = 1

Task = Struct.new(:id, :name, :plan_state, :execution_state, :deliverable, :dependency, keyword_init: true) do
  def kind
    return "L" if id.match?(/^L\d+$/)
    return "M" if id.match?(/^L\d+-M\d+$/)
    "S"
  end

  def parent_id
    case kind
    when "L" then nil
    when "M" then id.sub(/-M\d+$/, "")
    else id.sub(/-S\d+$/, "")
    end
  end
end

def abort_with(message)
  warn "ERROR: #{message}"
  exit 1
end

def parse_tasks(text)
  start_marker = "## タスク台帳"
  start_index = text.index(start_marker)
  abort_with("タスク台帳が見つかりません") unless start_index
  section = text[(start_index + start_marker.length)..]
  section = section.split(/^## /, 2).first
  tasks = []
  section.each_line do |line|
    next unless line.start_with?("| L")
    cells = line.split("|").map(&:strip)
    next unless cells.length >= 8
    id, name, plan_state, execution_state, deliverable, dependency = cells[1, 6]
    next unless id.match?(/^L[1-6](?:-M\d+)?(?:-S\d+)?$/)
    tasks << Task.new(id: id, name: name, plan_state: plan_state,
                      execution_state: execution_state, deliverable: deliverable,
                      dependency: dependency)
  end
  tasks
end

def validate_tasks!(tasks)
  duplicates = tasks.group_by(&:id).select { |_id, rows| rows.length > 1 }.keys
  abort_with("Task ID重複: #{duplicates.join(', ')}") unless duplicates.empty?
  counts = tasks.group_by(&:kind).transform_values(&:length)
  EXPECTED_COUNTS.each do |kind, expected|
    abort_with("#{kind}件数が#{counts.fetch(kind, 0)}件です（期待: #{expected}）") unless counts.fetch(kind, 0) == expected
  end
  abort_with("承認済みでないTaskがあります") unless tasks.all? { |task| task.plan_state == "承認済み" }
  ids = tasks.map(&:id)
  tasks.reject { |task| task.kind == "L" }.each do |task|
    abort_with("親Task欠落: #{task.id} -> #{task.parent_id}") unless ids.include?(task.parent_id)
  end
  tasks
end

def read_tasks_from_file
  validate_tasks!(parse_tasks(File.read(TASK_MAP)))
end

def read_tasks_from_snapshot(snapshot)
  out, status = Open3.capture2("rtk", "git", "-C", ROOT, "show", "#{snapshot}:docs/task-map.md")
  abort_with("snapshotからTask Mapを読めません: #{snapshot}") unless status.success?
  validate_tasks!(parse_tasks(out))
end

def expand_leaf_refs(text)
  refs = []
  dependency_part = text.split("。出力を", 2).first
  dependency_part.scan(/L(\d+)-M(\d+)-S(\d+)((?:[・〜]S\d+)*)/) do |l, m, first, tail|
    values = [first.to_i]
    current = first.to_i
    tail.scan(/([・〜])S(\d+)/) do |operator, value_text|
      value = value_text.to_i
      if operator == "〜"
        abort_with("逆順のleaf範囲: #{text}") if value < current
        values.concat(((current + 1)..value).to_a)
      else
        values << value
      end
      current = value
    end
    values.each { |value| refs << "L#{l}-M#{m}-S#{value}" }
  end
  refs.uniq
end

def expand_parent_refs(text)
  refs = []
  text.scan(/L(\d+)-M(\d+)((?:[・〜]M\d+)*)/) do |l, first, tail|
    values = [first.to_i]
    current = first.to_i
    tail.scan(/([・〜])M(\d+)/) do |operator, value_text|
      value = value_text.to_i
      values.concat(operator == "〜" ? ((current + 1)..value).to_a : [value])
      current = value
    end
    values.each { |value| refs << "L#{l}-M#{value}" }
  end
  text.scan(/L(\d+)〜L(\d+)/) do |first, last|
    refs.concat((first.to_i..last.to_i).map { |value| "L#{value}" })
  end
  refs
end

def dependency_clauses(text)
  clauses = []
  text.split("。").each do |raw_clause|
    clause = raw_clause.strip
    next if clause.empty? || clause.start_with?("出力を")
    label = nil
    value = clause
    if clause =~ /\A([^:：]+)[:：]\s*(.*)\z/
      label = Regexp.last_match(1)
      value = Regexp.last_match(2)
    end
    phase = case label
            when "完了・Merge" then "完了・Merge"
            when "Report確定" then "Report確定"
            when "確定", "ADR確定" then "確定"
            when "Gateへの入力" then "Gate入力"
            else
              clause.start_with?("Mergeは") ? "完了・Merge" :
                (clause.include?("Releaseは") ? "Release条件" : "着手")
            end
    clauses << [phase, value, clause]
  end
  clauses
end

def task_refs(text, tasks)
  ids = tasks.map(&:id)
  refs = expand_leaf_refs(text) + expand_parent_refs(text)
  without_leaf_refs = text.gsub(/L\d+-M\d+-S\d+(?:[・〜]S\d+)*/, "")
  without_leaf_refs.scan(/L\d+(?:-M\d+)?/) { |id| refs << id if ids.include?(id) }
  refs.select { |id| ids.include?(id) }.uniq
end

def decision_for(task)
  map = {
    "L1-M1" => "004", "L1-M2" => "005",
    "L2-M1" => "007", "L2-M2" => "010", "L2-M3" => "014",
    "L3-M1" => "016", "L3-M2" => "017",
    "L4-M1" => "019", "L4-M2" => "020", "L4-M3" => "021",
    "L5-M1" => "023", "L5-M2" => "024", "L5-M3" => "025", "L5-M4" => "026",
    "L6-M1" => "028", "L6-M2" => "029", "L6-M3" => "030", "L6-M4" => "031", "L6-M5" => "032"
  }
  return map.fetch(task.parent_id) if task.kind == "S"
  return map.fetch(task.id) if task.kind == "M"
  { "L1" => "001・002", "L2" => "006", "L3" => "015", "L4" => "018",
    "L5" => "022", "L6" => "027・033" }.fetch(task.id)
end

def task_type(task)
  return "evaluation" if task.id.match?(/^L6-M[45]-/)
  return "implementation" if task.name.match?(/実装|接続|実行Artifact|Black-box検証/)
  "design"
end

def path_for(task)
  id = task.id
  exact = {
    "L1-M1-S1" => "docs/design/knowledge-model/logical-schema.md", "L1-M1-S2" => "docs/design/knowledge-model/relations-and-integrity.md",
    "L1-M1-S3" => "docs/design/knowledge-model/evidence-derived-state.md", "L1-M1-S4" => "docs/design/knowledge-model/update-operations-and-lineage.md",
    "L1-M1-S5" => "docs/design/knowledge-model/requirements-traceability.md、docs/design/knowledge-model/README.md",
    "L1-M2-S1" => "docs/design/cli-contract/command-catalog.md", "L1-M2-S2" => "docs/design/cli-contract/collection-contract.md",
    "L1-M2-S3" => "docs/design/cli-contract/errors-and-exit-codes.md", "L1-M2-S4" => "docs/design/cli-contract/schemas/**",
    "L1-M2-S5" => "docs/design/cli-contract/versioning-and-compatibility.md", "L1-M2-S6" => "docs/design/cli-contract/requirements-traceability.md、docs/design/cli-contract/README.md",
    "L2-M1-S1" => "docs/design/knowledge-store/adr/0001-persistence-stack.md", "L2-M1-S2" => "docs/design/knowledge-store/physical-model-history-schema-evolution.md",
    "L2-M1-S3" => "docs/design/knowledge-store/repository-transaction-boundary.md、docs/design/knowledge-store/README.md",
    "L2-M2-S1" => "docs/design/search-infrastructure/search-and-index-requirements.md", "L2-M2-S2" => "docs/design/search-infrastructure/adr/0001-search-index-stack.md",
    "L2-M2-S3" => "docs/design/search-infrastructure/candidate-search-architecture.md", "L2-M2-S4" => "docs/design/search-infrastructure/index-lifecycle-recovery.md、docs/design/search-infrastructure/README.md",
    "L2-M3-S1" => "docs/design/cli-runtime/adr/0001-runtime-validation-distribution.md", "L2-M3-S2" => "docs/design/cli-runtime/command-to-port-architecture.md、docs/design/cli-runtime/README.md",
    "L3-M1-S1" => "docs/design/knowledge-acquisition/episode-input-evidence-boundary.md", "L3-M1-S2" => "docs/design/knowledge-acquisition/candidate-normalization-field-mapping.md",
    "L3-M1-S3" => ".agents/skills/knowledge-acquisition/references/candidate-markdown-contract.md", "L3-M1-S4" => ".agents/skills/knowledge-acquisition/SKILL.md、.agents/skills/knowledge-acquisition/references/**（S3所有Fileを除く）",
    "L3-M2-S1" => "docs/design/knowledge-update/candidate-acceptance.md", "L3-M2-S2" => "docs/design/knowledge-update/existing-knowledge-comparison.md",
    "L3-M2-S3" => ".agents/skills/knowledge-update/references/update-decision-contract.md", "L3-M2-S4" => ".agents/skills/knowledge-update/references/cli-command-mapping.md",
    "L3-M2-S5" => ".agents/skills/knowledge-update/SKILL.md、.agents/skills/knowledge-update/references/**（S3・S4所有Fileを除く）", "L3-M2-S6" => "tests/component/knowledge-update/**",
    "L4-M1-S1" => "docs/design/article-analysis/article-input-boundary.md", "L4-M1-S2" => "docs/design/article-analysis/overview-claim-decomposition.md",
    "L4-M1-S3" => "docs/design/article-analysis/location-support-traceability.md", "L4-M1-S4" => ".agents/skills/article-analysis/references/article-analysis-markdown-contract.md",
    "L4-M1-S5" => ".agents/skills/article-analysis/SKILL.md、.agents/skills/article-analysis/references/**（S4所有Fileを除く）",
    "L4-M2-S1" => "docs/design/knowledge-search/agentic-search-requirements.md", "L4-M2-S2" => "docs/design/knowledge-search/target-claim-query-reconstruction.md",
    "L4-M2-S3" => ".agents/skills/knowledge-search/references/cli-command-mapping.md", "L4-M2-S4" => "docs/design/knowledge-search/assessment-decision.md",
    "L4-M2-S5" => ".agents/skills/knowledge-search/references/knowledge-search-output-contract.md", "L4-M2-S6" => ".agents/skills/knowledge-search/SKILL.md、.agents/skills/knowledge-search/references/**（S3・S5所有Fileを除く）",
    "L4-M2-S7" => "tests/component/knowledge-search/**",
    "L4-M3-S1" => "docs/design/reading-value/input-alignment-and-evaluability.md", "L4-M3-S2" => "docs/design/reading-value/recognition-gain-application.md",
    "L4-M3-S3" => "docs/design/reading-value/reliability-applicability.md", "L4-M3-S4" => "docs/design/reading-value/recommendation-range-aggregation.md",
    "L4-M3-S5" => ".agents/skills/reading-value/references/reading-value-output-contract.md", "L4-M3-S6" => ".agents/skills/reading-value/SKILL.md、.agents/skills/reading-value/references/**（S5所有Fileを除く）",
    "L5-M1-S1" => "docs/design/orchestration/common-control-contract.md", "L5-M1-S2" => "docs/design/orchestration/child-skill-registry-contract.md",
    "L5-M1-S3" => "docs/design/orchestration/execution-state-envelope.md", "L5-M1-S4" => "docs/design/orchestration/correlation-trace-contract.md", "L5-M1-S5" => "docs/design/orchestration/execution-control-policy.md",
    "L5-M2-S1" => "docs/design/orchestration/article-reading-workflow.md", "L5-M2-S2" => "docs/design/orchestration/article-reading/article-search-fanout.md",
    "L5-M2-S3" => "docs/design/orchestration/article-reading/article-search-fanin-handoff.md", "L5-M2-S4" => "docs/design/orchestration/article-reading/article-research-cycle-routing.md", "L5-M2-S5" => "docs/design/orchestration/article-reading/article-workflow-termination.md",
    "L5-M3-S1" => "docs/design/orchestration/knowledge-accumulation-workflow.md", "L5-M3-S2" => "docs/design/orchestration/knowledge-accumulation/acquisition-result-branching.md",
    "L5-M3-S3" => "docs/design/orchestration/knowledge-accumulation/update-work-items.md", "L5-M3-S4" => "docs/design/orchestration/knowledge-accumulation/update-result-resumption.md", "L5-M3-S5" => "docs/design/orchestration/knowledge-accumulation/workflow-termination.md",
    "L5-M4-S1" => ".agents/skills/parent-orchestration/SKILL.md、.agents/skills/parent-orchestration/references/common/**、.agents/skills/parent-orchestration/references/child-registry.md",
    "L5-M4-S2" => ".agents/skills/parent-orchestration/references/workflows/article-reading.md", "L5-M4-S3" => ".agents/skills/parent-orchestration/references/workflows/knowledge-accumulation.md",
    "L5-M4-S4" => "tests/contract/parent-orchestration/**", "L5-M4-S5" => "tests/component/parent-orchestration/knowledge-accumulation/**", "L5-M4-S6" => "tests/component/parent-orchestration/article-reading/**",
    "L6-M1-S1" => "docs/design/evaluation/spec/evaluation-layer-policy.md", "L6-M1-S2" => "docs/design/evaluation/spec/coverage-responsibility-matrix.md",
    "L6-M1-S3" => "docs/design/evaluation/spec/verdict-criteria.md", "L6-M1-S4" => "docs/design/evaluation/spec/search-trace-diagnostics.md", "L6-M1-S5" => "docs/design/evaluation/spec/evaluation-report-contract.md",
    "L6-M2-S1" => "docs/design/evaluation/datasets/dataset-contract.md", "L6-M2-S2" => "docs/design/evaluation/datasets/scenario-catalog.md",
    "L6-M2-S3" => "tests/evaluation/datasets/schema/**、tests/evaluation/datasets/tools/**", "L6-M2-S4" => "tests/evaluation/datasets/scenarios/a-e/**",
    "L6-M2-S5" => "tests/evaluation/datasets/scenarios/f-h/**", "L6-M2-S6" => "tests/evaluation/datasets/scenarios/i-j-reading-value/**",
    "L6-M3-S1" => "docs/design/evaluation/harness/harness-requirements.md", "L6-M3-S2" => "docs/design/evaluation/adr/evaluation-harness-stack.md",
    "L6-M3-S3" => "docs/design/evaluation/harness/harness-architecture.md", "L6-M3-S4" => "tests/evaluation/harness/core/**",
    "L6-M3-S5" => "tests/evaluation/harness/adapters/cli/**", "L6-M3-S6" => "tests/evaluation/harness/adapters/skills/**", "L6-M3-S7" => "tests/evaluation/harness/adapters/workflows/**",
    "L6-M5-S1" => "tests/evaluation/workflows/knowledge-accumulation/**、tests/evaluation/reports/workflows/knowledge-accumulation/**",
    "L6-M5-S2" => "tests/evaluation/workflows/article-reading/**、tests/evaluation/reports/workflows/article-reading/**",
    "L6-M5-S3" => "tests/evaluation/workflows/regression/**、tests/evaluation/reports/final/**"
  }
  return exact[id] if exact.key?(id)
  return "L2 Design Freeze Gate通過記録の#{id}専用Path（Issue作成時TBD、Ready前に実値化）" if id.start_with?("L2-")
  if id.start_with?("L6-M4-")
    target = { "S1" => "cli", "S2" => "l3/knowledge-acquisition", "S3" => "l3/knowledge-update",
               "S4" => "l4/article-analysis", "S5" => "l4/knowledge-search", "S6" => "l4/reading-value" }.fetch(id[/S\d+$/])
    return "tests/evaluation/suites/#{target}/**、tests/evaluation/reports/#{target}/**"
  end
  abort_with("Path規則がありません: #{id}")
end

def validate_paths!(tasks)
  paths = tasks.select { |task| task.kind == "S" }.group_by { |task| path_for(task) }
  collisions = paths.select { |_path, owners| owners.length > 1 }
  return if collisions.empty?
  details = collisions.map { |path, owners| "#{path}=#{owners.map(&:id).join(',')}" }
  abort_with("leaf間で成果物Pathが重複しています: #{details.join('; ')}")
end

def issue_link(number, task)
  "##{number} `[#{task.id}] #{task.name}`"
end

def linked_dependencies(task, tasks, mapping)
  refs = task_refs(task.dependency, tasks)
  return "なし" if refs.empty?
  refs.map do |id|
    dependency_task = tasks.find { |candidate| candidate.id == id }
    number = mapping[id]
    number ? issue_link(number, dependency_task) : "`#{id}`"
  end.join("、")
end

DEPENDENCY_KINDS = ["着手依存", "完了・Merge依存", "Gateへの入力", "Gate通過依存", "Release条件"].freeze

def dependency_kind(label, clause)
  return "完了・Merge依存" if ["完了", "完了・Merge", "Report確定", "確定", "ADR確定"].include?(label)
  return "Gateへの入力" if label == "Gateへの入力"
  return "Release条件" if ["完了・Release", "Release条件"].include?(label) || clause.include?("Releaseは")
  return "Gate通過依存" if ["本番実装", "Mock実装", "本番Skill実装"].include?(label)
  "着手依存"
end

def dependency_rows(task, tasks, mapping)
  grouped = Hash.new { |hash, key| hash[key] = [] }
  task.dependency.split("。").each do |raw|
    clause = raw.strip
    next if clause.empty? || clause.start_with?("出力を")
    label = nil
    value = clause
    if clause =~ /\A([^:：]+)[:：]\s*(.*)\z/
      label = Regexp.last_match(1)
      value = Regexp.last_match(2)
    end
    refs = task_refs(value, tasks).map do |id|
      dependency_task = tasks.find { |candidate| candidate.id == id }
      mapping[id] ? issue_link(mapping[id], dependency_task) : "`#{id}`（Issue作成待ち）"
    end
    milestone = value.scan(/(?:L\d+(?:-M\d+)?\s+)?(?:Design Freeze Gate|Target Readiness Gate|Final Quality Gate)|L\d+ Release/).uniq
    rendered = (refs + milestone.map { |item| "`#{item}`" }).uniq
    rendered << value if rendered.empty?
    grouped[dependency_kind(label, clause)].concat(rendered)
  end
  DEPENDENCY_KINDS.map do |kind|
    values = grouped[kind].uniq
    "| #{kind} | #{values.empty? ? '該当なし' : values.join('、')} | #{values.empty? ? 'Task Mapに直接条件なし' : 'Task Map原文の依存時点を満たす'} |"
  end.join("\n")
end

def editable_tracking_body
  <<~MD

    <!-- human-progress:start -->
    ## 実行記録（手動更新欄）

    - [ ] 直下の子Issueがすべて完了している
    - [ ] 親の到達状態が子成果物の総和として成立している
    - [ ] Gate／Release条件を満たしている

    - 統合Commit／Release evidence:
    - 未解決事項:
    <!-- human-progress:end -->
  MD
end

def editable_leaf_body
  <<~MD

    <!-- human-progress:start -->
    ## 実行記録（手動更新欄）

    ### 実施チェック

    - [ ] 主成果物を作成した
    - [ ] 直接依存とGate条件を確認した
    - [ ] Owner／Path境界を守った
    - [ ] Task固有の検証を完了した
    - [ ] 未解決の契約差異とReady後のTBDがない

    ### 着手・検証Evidence

    - worktree起点SHA:
    - Branch:
    - 依存TaskのMerge commit:
    - Gate通過記録:
    - 検証Command／Review:
    - 結果:
    - 関連PR:
    <!-- human-progress:end -->
  MD
end

def fixed_url(repo, snapshot, path)
  "https://github.com/#{repo}/blob/#{snapshot}/#{path}"
end

def generated_wrap(task, snapshot, content)
  <<~BODY
    <!-- knowledge-task-id: #{task.id} -->
    <!-- generated-content:start -->
    <!-- planning-snapshot: #{snapshot} -->
    #{content.rstrip}
    <!-- generated-content:end -->
  BODY
end

def tracking_body(task, tasks, mapping, repo, snapshot, issue_states = {})
  parent = task.kind == "L" ? "原典Issue ##{ROOT_ISSUE}" : issue_link(mapping.fetch(task.parent_id), tasks.find { |item| item.id == task.parent_id })
  children = tasks.select { |item| item.parent_id == task.id }
  child_lines = children.map do |child|
    number = mapping[child.id]
    checked = issue_states[child.id] == "CLOSED" ? "x" : " "
    number ? "- [#{checked}] #{issue_link(number, child)}" : "- [ ] `[#{child.id}] #{child.name}`（Issue作成待ち）"
  end.join("\n")
  content = <<~MD
    ## タスク情報

    - Task ID: `#{task.id}`
    - 親Task: #{parent}
    - 区分: #{task.kind == 'L' ? '大分類tracking' : '中分類tracking'}
    - Planning snapshot commit SHA: `#{snapshot}`
    - Task Map固定リンク: #{fixed_url(repo, snapshot, 'docs/task-map.md')}
    - 原典Issue: ##{ROOT_ISSUE}
    - 関連する決定ID: 決定#{decision_for(task)}

    ## 目的・主成果物

    #{task.name}。到達状態は「#{task.deliverable}」。

    ## 直接依存

    - Task Map原文: #{task.dependency}
    - 参照Task: #{linked_dependencies(task, tasks, mapping)}

    ## 直下の子Issue

    #{child_lines}

    ## 完了条件

    - 直下の子Issueがすべて完了している
    - `#{task.deliverable}`が子成果物の総和として成立している
    - Task Mapに記載されたGate／Release条件を満たしている
    - 未解決の契約差異がない

    ## 境界

    このIssueは進捗と統合条件を追跡する。子Issueと重複する設計・実装・評価成果物を独自に作成しない。Gate／ReleaseはMilestoneであり、別Issueを作成しない。
  MD
  generated_wrap(task, snapshot, content) + editable_tracking_body
end

def leaf_body(task, tasks, mapping, repo, snapshot)
  parent_task = tasks.find { |item| item.id == task.parent_id }
  parent = issue_link(mapping.fetch(task.parent_id), parent_task)
  path = path_for(task)
  downstream_url = fixed_url(repo, snapshot, "docs/task-connections.md")
  content = <<~MD
    ## タスク情報

    - Task ID: `#{task.id}`
    - 親Task ID: `#{task.parent_id}`
    - 親Issue: #{parent}
    - タスク名: #{task.name}
    - タスク種別: `#{task_type(task)}`
    - Planning snapshot commit SHA: `#{snapshot}`
    - Task Map固定リンク: #{fixed_url(repo, snapshot, 'docs/task-map.md')}
    - 原典Issue: ##{ROOT_ISSUE}
    - 関連する原典章: Issue #1、および決定#{decision_for(task)}に記録された対応章
    - 関連する決定ID: 決定#{decision_for(task)}・034〜037

    ## 目的

    #{task.name}ことで、`#{task.deliverable}`を完成させる。

    ## 原典との差分

    - 固定入力: Issue #1、Planning snapshotのTask Map、決定#{decision_for(task)}、直接依存の承認済み成果物
    - 原典で決定済みだが未実施の成果物: #{task.deliverable}
    - このleafで決める未確定事項: Task名が要求する詳細化または実装に限定する
    - 選び直さない事項: Issue #1とTask Mapで確定した責務境界、依存DAG、Gate／Release、Codex／CLI境界

    ## 実施内容

    - [ ] `#{task.deliverable}`を指定Owner境界で作成する
    - [ ] 直接依存とGate条件を確認し、Task固有のReviewまたは検証を記録する

    ## 成果物と所有

    | 項目 | 内容 |
    | --- | --- |
    | 主成果物 | #{task.deliverable} |
    | 書込み可能なPath／Glob | `#{path}` |
    | 単一Owner | `#{task.id}` |
    | read-only入力 | #{linked_dependencies(task, tasks, mapping)} |
    | 共有資産と単一Owner | Task Mapのworktree所有境界に従う。共有物は記載Ownerへ変更を戻す |
    | Gate通過記録 | Task Map原文のGate条件に従う。Gate非対象なら該当なし |

    書込み可能なのは上表の所有Pathと承認済みの明示例外だけとし、それ以外は変更しない。

    ## 完了条件

    - [ ] `#{task.deliverable}`が作成され、Task名の到達状態を満たす
    - [ ] Task Mapの直接依存、Gate、Releaseの必要条件を満たす
    - [ ] 単一Ownerと書込みPathの境界を守る
    - [ ] Task固有の検証結果またはReview evidenceをIssueへ記録する
    - [ ] 未解決の契約差異と、Ready後のTBDがない

    ## 対象外

    - Issue #1で既に決定した原則・責務の選び直し
    - 別Taskが単一Ownerである成果物、共有評価、推移的依存の再実装

    ## 依存関係

    | 種別 | Task／Milestone | 必要な成果物・状態 |
    | --- | --- | --- |
    #{dependency_rows(task, tasks, mapping)}

    - 後続接続の固定リンク: #{downstream_url}

    ## 着手判定

    | 確認項目 | 結果・Evidence |
    | --- | --- |
    | Planning snapshot SHAとTask Map固定リンクが一致する | 未確認 |
    | 着手依存TaskのMerge commit | 着手時に記録 |
    | 必要なGateの通過記録、またはGate前Taskであること | 着手時に記録 |
    | worktree起点SHA | 着手時に記録 |
    | 必須値に未解決TBDがない | 着手時に確認 |
    | 並行Taskと書込みPathが競合しない | 着手時に確認 |

    - [ ] 上表を確認し、このTaskは着手可能である

    ## worktree・Merge

    | 項目 | 内容 |
    | --- | --- |
    | planning baseline SHA | `#{snapshot}` |
    | worktree起点SHA | 着手時に確定 |
    | Branch | 着手時に確定 |
    | 所有Path／Glob | `#{path}` |
    | 共有物と単一Owner | Task Mapのworktree所有境界に従う |
    | 並行可能Task | Pathが交差せず、必要な依存を満たすTask |
    | 直列化するTaskと理由 | Task Mapの直接依存・Merge順に従う |
    | Merge前提 | 直接依存の必要状態、Gate通過、Task固有検証合格 |
    | Merge順 | Task Mapの承認済みDAGに従う |
    | 統合先 | planning baselineを含む統合branch |

    ## 検証

    | 種別 | 方法・Command | 合格条件 | Evidence |
    | --- | --- | --- | --- |
    | 静的確認 | 成果物・参照・Owner境界をReview | 不足・重複・境界違反がない | 完了時に記録 |
    | Task固有テスト／Review | 成果物種別に応じた固有検証 | Task固有の完了条件を満たす | 完了時に記録 |
    | 契約・統合確認 | Task Mapに明記された場合のみ実施 | 対象契約に適合する | 完了時に記録 |
    | 後続評価 | L6等の後続Issueで実施 | このIssueでは重複実施しない | 後続Issueへ記録 |

    ## 差異を発見した場合

    - [ ] 作業を停止する
    - [ ] Issue内で新しい依存、Path、設計判断を決めない
    - [ ] Task MapまたはGate記録の修正案を議論記録へ残す
    - [ ] 必要な再承認後、Planning snapshotとIssueを同期する

    ## 関連Issue・PR

    - 原典Issue: ##{ROOT_ISSUE}
    - 親Issue: #{parent}
    - Blocked by: #{linked_dependencies(task, tasks, mapping)}
    - 関連PR: 未作成
  MD
  content = content.gsub("- [ ]", "- 要件:")
  generated_wrap(task, snapshot, content) + editable_leaf_body
end

def body_for(task, tasks, mapping, repo, snapshot, issue_states = {})
  task.kind == "S" ? leaf_body(task, tasks, mapping, repo, snapshot) : tracking_body(task, tasks, mapping, repo, snapshot, issue_states)
end

def gh_json(repo)
  out, err, status = Open3.capture3("rtk", "gh", "issue", "list", "--repo", repo, "--state", "all",
                                   "--limit", "500", "--json", "number,title,body,url,state")
  abort_with("gh issue list失敗: #{err}") unless status.success?
  JSON.parse(out)
end

def managed_issues(issues)
  by_id = Hash.new { |hash, key| hash[key] = [] }
  issues.each do |issue|
    marker = issue.fetch("body", "").match(/<!-- knowledge-task-id: (L\d+(?:-M\d+)?(?:-S\d+)?) -->/)
    by_id[marker[1]] << issue if marker
  end
  duplicates = by_id.select { |_id, matches| matches.length > 1 }
  abort_with("GitHub上のTask ID重複: #{duplicates.keys.join(', ')}") unless duplicates.empty?
  by_id.transform_values(&:first)
end

def title_conflicts!(issues, tasks, managed)
  tasks.each do |task|
    prefix = "[#{task.id}]"
    conflicts = issues.select { |issue| issue["title"].start_with?(prefix) && managed[task.id] != issue }
    abort_with("Markerなし／別Markerのタイトル競合: #{task.id} -> #{conflicts.map { |i| i['number'] }.join(',')}") unless conflicts.empty?
  end
end

def run_gh_with_body(args, body)
  out, err, status = Open3.capture3("rtk", "gh", *args, stdin_data: body)
  abort_with("gh #{args.join(' ')} 失敗: #{err}") unless status.success?
  out.strip
end

def create_issue(repo, task, body)
  title = "[#{task.id}] #{task.name}"
  output = run_gh_with_body(["issue", "create", "--repo", repo, "--title", title, "--body-file", "-"], body)
  number = output[/\/(\d+)\z/, 1]
  abort_with("Issue番号を取得できません: #{output}") unless number
  puts "CREATE #{task.id} -> ##{number}"
  number.to_i
end

def replace_generated(existing_body, desired_body)
  start_marker = "<!-- generated-content:start -->"
  end_marker = "<!-- generated-content:end -->"
  starts = existing_body.scan(start_marker).length
  ends = existing_body.scan(end_marker).length
  abort_with("管理Marker不整合のため本文を更新できません") unless starts == 1 && ends == 1
  desired_generated = desired_body[/#{Regexp.escape(start_marker)}.*?#{Regexp.escape(end_marker)}/m]
  existing_body.sub(/#{Regexp.escape(start_marker)}.*?#{Regexp.escape(end_marker)}/m, desired_generated)
end

def update_issue(repo, issue, task, body)
  title = "[#{task.id}] #{task.name}"
  merged_body = replace_generated(issue.fetch("body", ""), body)
  return false if issue["title"] == title && issue.fetch("body", "") == merged_body
  run_gh_with_body(["issue", "edit", issue["number"].to_s, "--repo", repo,
                    "--title", title, "--body-file", "-"], merged_body)
  puts "UPDATE #{task.id} -> ##{issue['number']}"
  true
end

def root_children_block(tasks, mapping, issue_states = {})
  lines = tasks.select { |task| task.kind == "L" }.map do |task|
    checked = issue_states[task.id] == "CLOSED" ? "x" : " "
    "- [#{checked}] #{issue_link(mapping.fetch(task.id), task)}"
  end.join("\n")
  <<~MD.rstrip
    <!-- task-map-children:start -->
    ## Task Map: 大分類Issue

    Planning snapshotに基づく直下の大分類tracking Issueです。Gate／ReleaseはIssue化していません。

    #{lines}
    <!-- task-map-children:end -->
  MD
end

def update_root_issue(repo, root_issue, tasks, mapping, issue_states = {})
  body = root_issue.fetch("body", "")
  block = root_children_block(tasks, mapping, issue_states)
  start_marker = "<!-- task-map-children:start -->"
  end_marker = "<!-- task-map-children:end -->"
  desired = if body.include?(start_marker) && body.include?(end_marker)
              body.sub(/#{Regexp.escape(start_marker)}.*?#{Regexp.escape(end_marker)}/m, block)
            else
              "#{body.rstrip}\n\n#{block}\n"
            end
  return false if desired == body
  run_gh_with_body(["issue", "edit", root_issue["number"].to_s, "--repo", repo, "--body-file", "-"], desired)
  puts "UPDATE ROOT -> ##{root_issue['number']}"
  true
end

def verify!(tasks, issues, snapshot, repo)
  managed = managed_issues(issues)
  expected_ids = tasks.map(&:id)
  extra = managed.keys - expected_ids
  missing = expected_ids - managed.keys
  abort_with("余分な管理Issue: #{extra.join(', ')}") unless extra.empty?
  abort_with("不足する管理Issue: #{missing.join(', ')}") unless missing.empty?
  counts = managed.keys.group_by do |id|
    id.match?(/-S\d+$/) ? "S" : (id.match?(/-M\d+$/) ? "M" : "L")
  end.transform_values(&:length)
  abort_with("GitHub件数不一致: #{counts}") unless EXPECTED_COUNTS.all? { |kind, value| counts.fetch(kind, 0) == value }
  managed.each do |id, issue|
    task = tasks.find { |candidate| candidate.id == id }
    abort_with("タイトル不一致: #{id}") unless issue["title"] == "[#{id}] #{task.name}"
    abort_with("snapshot不一致: #{id}") unless issue.fetch("body", "").include?("<!-- planning-snapshot: #{snapshot} -->")
    abort_with("固定リンク不一致: #{id}") unless issue.fetch("body", "").include?(fixed_url(repo, snapshot, "docs/task-map.md"))
    if task.kind == "S"
      %w[タスク情報 目的 原典との差分 実施内容 成果物と所有 完了条件 対象外 依存関係 着手判定 worktree・Merge 検証 差異を発見した場合 関連Issue・PR].each do |heading|
        abort_with("leaf必須Section欠落: #{id} #{heading}") unless issue.fetch("body", "").include?("## #{heading}")
      end
    else
      child_section = issue.fetch("body", "")[/## 直下の子Issue\n\n(.*?)\n\n## 完了条件/m, 1].to_s
      expected_children = tasks.select { |candidate| candidate.parent_id == id }
      expected_children.each do |child|
        count = child_section.scan("[#{child.id}]").length
        abort_with("親子チェックリスト不一致: #{id} -> #{child.id} count=#{count}") unless count == 1
      end
      abort_with("親子チェックリスト件数不一致: #{id}") unless child_section.scan(/^- \[[ x]\] #\d+ /).length == expected_children.length
    end
  end
  root_issue = issues.find { |issue| issue["number"] == ROOT_ISSUE }
  abort_with("原典Issue ##{ROOT_ISSUE}がありません") unless root_issue
  root_section = root_issue.fetch("body", "")[/<!-- task-map-children:start -->(.*?)<!-- task-map-children:end -->/m, 1].to_s
  tasks.select { |task| task.kind == "L" }.each do |task|
    abort_with("原典Issueの大分類接続不一致: #{task.id}") unless root_section.scan("[#{task.id}]").length == 1
  end
  puts "VERIFY OK L=#{counts['L']} M=#{counts['M']} S=#{counts['S']} TOTAL=#{managed.length}"
end

options = { repo: REPO_DEFAULT, mode: "check", snapshot: nil }
OptionParser.new do |parser|
  parser.banner = "Usage: task_issue_sync.rb [--check|--apply|--verify] --snapshot SHA [--repo OWNER/REPO]"
  parser.on("--check") { options[:mode] = "check" }
  parser.on("--apply") { options[:mode] = "apply" }
  parser.on("--verify") { options[:mode] = "verify" }
  parser.on("--render-connections") { options[:mode] = "render_connections" }
  parser.on("--repo REPO") { |value| options[:repo] = value }
  parser.on("--snapshot SHA") { |value| options[:snapshot] = value }
end.parse!

if options[:mode] == "render_connections"
  tasks = read_tasks_from_file
  validate_paths!(tasks)
  by_id = tasks.each_with_object({}) { |task, result| result[task.id] = task }
  edges = []
  milestones = []
  tasks.select { |task| task.kind == "S" }.each do |consumer|
    dependency_clauses(consumer.dependency).each do |phase, value, raw_clause|
      expand_leaf_refs(value).each do |producer|
        abort_with("接続先に存在しないleaf: #{producer} -> #{consumer.id}") unless by_id[producer]&.kind == "S"
        gate_release = raw_clause.match?(/Gate|Release/) ? raw_clause : "—"
        edges << [producer, consumer.id, by_id.fetch(producer).deliverable, phase, "read-only", gate_release,
                  "Task Map「#{by_id.fetch(producer).parent_id}内」"]
      end
      next unless raw_clause.match?(/Gate|Release/)
      milestone = raw_clause.gsub(/L\d+-M\d+-S\d+(?:[・〜]S\d+)*[、]?/, "").strip
      milestones << [milestone, consumer.id, phase] unless milestone.empty?
    end
  end
  # L3-M2-S2の出力は、同Taskの依存ではなくL2-M2-S1への承認済みGate入力である。
  producer = by_id.fetch("L3-M2-S2")
  edges.reject! { |row| row[0] == producer.id && row[1] == "L2-M2-S1" }
  edges << [producer.id, "L2-M2-S1", producer.deliverable, "Gate入力", "read-only",
            "L2 Design Freeze Gate入力", "Task Map「L3-M2内」"]
  edges.uniq!
  edges.sort_by! { |row| [row[0], row[1], row[3]] }
  puts "| Producer leaf | 直接の後続利用者 | 利用成果物 | 依存時点 | 利用方法 | Gate／Release | worktree所有記録への参照 |"
  puts "| --- | --- | --- | --- | --- | --- | --- |"
  edges.each { |row| puts "| #{row.join(' | ')} |" }
  producers = edges.map(&:first).uniq
  terminals = tasks.select { |task| task.kind == "S" && !producers.include?(task.id) }
  puts
  puts "### 終端leaf"
  puts
  puts "| Leaf Task | 後続利用者 | 主成果物 |"
  puts "| --- | --- | --- |"
  terminals.each { |task| puts "| #{task.id} | なし（終端Task） | #{task.deliverable} |" }
  puts
  puts "### 明示されたGate／Release接続"
  puts
  puts "| Milestone表現 | Consumer leaf | 依存時点 |"
  puts "| --- | --- | --- |"
  milestones.uniq.each { |row| puts "| #{row.join(' | ')} |" }
  if milestones.empty?
    puts "| なし | — | — |"
  end
  exit
end

snapshot = options[:snapshot]
abort_with("--snapshotに40桁SHAが必要です") unless snapshot&.match?(/\A[0-9a-f]{40}\z/)
tasks = read_tasks_from_snapshot(snapshot)
validate_paths!(tasks)

if options[:mode] == "check"
  puts "CHECK OK L=6 M=19 S=116 TOTAL=141 snapshot=#{snapshot}"
  exit
end

issues = gh_json(options[:repo])
managed = managed_issues(issues)
title_conflicts!(issues, tasks, managed)

if options[:mode] == "verify"
  verify!(tasks, issues, snapshot, options[:repo])
  exit
end

abort_with("未知のmode: #{options[:mode]}") unless options[:mode] == "apply"

mapping = managed.transform_values { |issue| issue["number"] }
%w[L M S].each do |kind|
  tasks.select { |task| task.kind == kind }.each do |task|
    next if mapping.key?(task.id)
    mapping[task.id] = create_issue(options[:repo], task, body_for(task, tasks, mapping, options[:repo], snapshot))
  end
end

issues = gh_json(options[:repo])
managed = managed_issues(issues)
mapping = managed.transform_values { |issue| issue["number"] }
issue_states = managed.transform_values { |issue| issue["state"] }
tasks.each do |task|
  issue = managed.fetch(task.id)
  update_issue(options[:repo], issue, task, body_for(task, tasks, mapping, options[:repo], snapshot, issue_states))
end

issues = gh_json(options[:repo])
root_issue = issues.find { |issue| issue["number"] == ROOT_ISSUE }
abort_with("原典Issue ##{ROOT_ISSUE}がありません") unless root_issue
update_root_issue(options[:repo], root_issue, tasks, mapping, issue_states)

verify!(tasks, gh_json(options[:repo]), snapshot, options[:repo])
