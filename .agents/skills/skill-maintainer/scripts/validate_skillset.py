#!/usr/bin/env python3
"""Lightweight structural validation for this skillset.

Checks only deterministic structure. It does not evaluate model behavior.
"""
from __future__ import annotations

import sys
from pathlib import Path

try:
    import yaml
except ImportError:
    print("ERROR: PyYAML is required to run this validator.", file=sys.stderr)
    raise SystemExit(2)


def load_frontmatter(path: Path) -> dict:
    text = path.read_text(encoding="utf-8")
    if not text.startswith("---\n"):
        raise ValueError("SKILL.md must start with YAML frontmatter")
    parts = text.split("---", 2)
    if len(parts) != 3:
        raise ValueError("invalid YAML frontmatter delimiters")
    data = yaml.safe_load(parts[1]) or {}
    return data


def resolve_project_paths(repo: Path) -> tuple[list[Path], Path | None]:
    """Resolve the split implementation and documentation roots of this project."""
    direct_workflow = repo / ".ai" / "workflow"
    documents_dir = repo / "documents"
    documents_workflow = documents_dir / ".ai" / "workflow"

    if direct_workflow.exists():
        skill_roots = [repo / ".agents" / "skills"]
        parent_skill_root = repo.parent / ".agents" / "skills"
        if repo.name == "documents" and parent_skill_root.exists():
            skill_roots.append(parent_skill_root)
        return skill_roots, direct_workflow

    if documents_workflow.exists():
        return [repo / ".agents" / "skills", documents_dir / ".agents" / "skills"], documents_workflow

    return [repo / ".agents" / "skills"], None


def direct_child_skill_files(skill_roots: list[Path]) -> list[Path]:
    """Collect direct-child skills once even if another root links to the same skill."""
    files: list[Path] = []
    seen: set[Path] = set()
    for root in skill_roots:
        for path in sorted(root.glob("*/SKILL.md")):
            resolved = path.resolve()
            if resolved not in seen:
                seen.add(resolved)
                files.append(path)
    return files


def main() -> int:
    repo = Path.cwd()
    skill_roots, workflow_root = resolve_project_paths(repo)
    errors: list[str] = []
    warnings: list[str] = []

    skill_files = direct_child_skill_files(skill_roots)
    if not skill_files:
        errors.append("No direct-child skills found under the configured .agents/skills roots")

    names: dict[str, Path] = {}
    for path in skill_files:
        try:
            fm = load_frontmatter(path)
        except Exception as exc:  # noqa: BLE001
            errors.append(f"{path}: {exc}")
            continue

        name = fm.get("name")
        description = fm.get("description")
        if not isinstance(name, str) or not name.strip():
            errors.append(f"{path}: missing non-empty frontmatter name")
        elif name != path.parent.name:
            errors.append(f"{path}: frontmatter name '{name}' != directory '{path.parent.name}'")
        elif name in names:
            errors.append(f"duplicate skill name '{name}': {names[name]} and {path}")
        else:
            names[name] = path

        if not isinstance(description, str) or not description.strip():
            errors.append(f"{path}: missing non-empty frontmatter description")
        elif len(description) > 1024:
            errors.append(f"{path}: description exceeds 1024 characters")

    direct_resolved = {p.resolve() for p in skill_files}
    nested: list[Path] = []
    seen_nested: set[Path] = set()
    for root in skill_roots:
        for path in root.rglob("SKILL.md"):
            resolved = path.resolve()
            if resolved not in direct_resolved and resolved not in seen_nested:
                seen_nested.add(resolved)
                nested.append(path)
    for path in nested:
        warnings.append(
            f"Nested SKILL.md detected: {path}. Confirm it is intended to be independently discovered."
        )

    artifact_map_path = workflow_root / "artifact-map.yaml" if workflow_root else None
    if artifact_map_path and artifact_map_path.exists():
        try:
            artifact_map = yaml.safe_load(artifact_map_path.read_text(encoding="utf-8")) or {}
            for artifact_name, spec in (artifact_map.get("artifacts") or {}).items():
                owner = (spec or {}).get("owner")
                if owner and owner not in names and owner != "implementation-skill-builder":
                    errors.append(
                        f"artifact '{artifact_name}' references unknown owner skill '{owner}'"
                    )
        except Exception as exc:  # noqa: BLE001
            errors.append(f"{artifact_map_path}: invalid YAML: {exc}")
    else:
        errors.append("Missing workflow artifact-map.yaml at this repository entry point")

    if workflow_root:
        for path in workflow_root.glob("*.yaml"):
            try:
                yaml.safe_load(path.read_text(encoding="utf-8"))
            except Exception as exc:  # noqa: BLE001
                errors.append(f"{path}: invalid YAML: {exc}")

    for warning in warnings:
        print(f"WARN: {warning}")
    for error in errors:
        print(f"ERROR: {error}")

    if errors:
        print(f"FAILED: {len(errors)} error(s), {len(warnings)} warning(s)")
        return 1

    print(f"OK: {len(skill_files)} skills, {len(warnings)} warning(s)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
