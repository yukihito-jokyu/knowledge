#!/usr/bin/env python3
"""Recompute coverage for production packages changed from a review baseline."""

from __future__ import annotations

import argparse
import json
import re
import subprocess
import sys
import tempfile
from decimal import Decimal
from pathlib import Path


PRODUCTION_ROOTS = (
    Path("cmd/knowledge"),
    Path("internal/application"),
    Path("internal/domain"),
    Path("internal/persistence/sqlite"),
)
TOTAL_PATTERN = re.compile(r"^total:\s+\(statements\)\s+([0-9]+(?:\.[0-9]+)?)%$")


def run(*args: str, cwd: Path) -> subprocess.CompletedProcess[str]:
    return subprocess.run(args, cwd=cwd, check=False, text=True, capture_output=True)


def repository_root() -> Path:
    result = run("git", "rev-parse", "--show-toplevel", cwd=Path.cwd())
    if result.returncode != 0:
        raise RuntimeError(result.stderr.strip() or "repository rootを取得できません")
    return Path(result.stdout.strip()).resolve()


def is_production_go_file(path: Path) -> bool:
    if path.suffix != ".go" or path.name.endswith("_test.go"):
        return False
    return any(path == root / path.name or root in path.parents for root in PRODUCTION_ROOTS)


def changed_production_files(repo: Path, base: str) -> list[Path]:
    verify = run("git", "rev-parse", "--verify", f"{base}^{{commit}}", cwd=repo)
    if verify.returncode != 0:
        raise RuntimeError(f"baseline commitを解決できません: {base}")

    diff = run("git", "diff", "--name-only", "--diff-filter=ACDMRTUXB", base, "--", cwd=repo)
    if diff.returncode != 0:
        raise RuntimeError(diff.stderr.strip() or "candidate diffを取得できません")
    untracked = run("git", "ls-files", "--others", "--exclude-standard", cwd=repo)
    if untracked.returncode != 0:
        raise RuntimeError(untracked.stderr.strip() or "未追跡fileを取得できません")

    paths = {
        Path(line)
        for line in (*diff.stdout.splitlines(), *untracked.stdout.splitlines())
        if line
    }
    return sorted(path for path in paths if is_production_go_file(path))


def package_for_directory(repo: Path, directory: Path) -> str | None:
    absolute = repo / directory
    if not absolute.is_dir() or not any(
        path.suffix == ".go" and not path.name.endswith("_test.go")
        for path in absolute.iterdir()
        if path.is_file()
    ):
        return None
    result = run(
        "go",
        "list",
        "-f",
        "{{if .GoFiles}}{{.ImportPath}}{{end}}",
        f"./{directory.as_posix()}",
        cwd=repo,
    )
    if result.returncode != 0:
        raise RuntimeError(result.stderr.strip() or f"packageを解決できません: {directory}")
    return result.stdout.strip() or None


def measure_package(repo: Path, package: str) -> dict[str, object]:
    with tempfile.TemporaryDirectory(prefix="knowledge-review-coverage-") as temp_dir:
        profile = Path(temp_dir) / "coverage.out"
        test = run("go", "test", package, f"-coverprofile={profile}", cwd=repo)
        if test.returncode != 0:
            return {
                "package": package,
                "status": "test_failed",
                "coverage_percent": None,
                "stderr": test.stderr.strip(),
            }

        blocks = [line for line in profile.read_text(encoding="utf-8").splitlines()[1:] if line]
        if not blocks:
            return {
                "package": package,
                "status": "not_applicable",
                "coverage_percent": None,
            }

        cover = run("go", "tool", "cover", f"-func={profile}", cwd=repo)
        if cover.returncode != 0:
            return {
                "package": package,
                "status": "coverage_failed",
                "coverage_percent": None,
                "stderr": cover.stderr.strip(),
            }
        total = next(
            (
                match.group(1)
                for line in cover.stdout.splitlines()
                if (match := TOTAL_PATTERN.match(line))
            ),
            None,
        )
        if total is None:
            return {
                "package": package,
                "status": "coverage_failed",
                "coverage_percent": None,
                "stderr": "coverage合計を取得できません",
            }
        return {
            "package": package,
            "status": "pass" if Decimal(total) == Decimal("100") else "below_threshold",
            "coverage_percent": total,
        }


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="candidateで変更されたproduction packageのcoverageを独立再測定します。"
    )
    parser.add_argument("--base", default="HEAD", help="Baseline Snapshotのcommit。default: HEAD")
    parser.add_argument("--json", action="store_true", help="結果をJSONで出力します。")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    try:
        repo = repository_root()
        changed_files = changed_production_files(repo, args.base)
        directories = sorted({path.parent for path in changed_files})
        removed_directories: list[str] = []
        packages: list[str] = []
        for directory in directories:
            package = package_for_directory(repo, directory)
            if package is None:
                removed_directories.append(directory.as_posix())
            else:
                packages.append(package)
        results = [measure_package(repo, package) for package in sorted(set(packages))]
    except RuntimeError as exc:
        print(str(exc), file=sys.stderr)
        return 2

    passed = all(result["status"] in {"pass", "not_applicable"} for result in results)
    report = {
        "base": args.base,
        "changed_production_files": [path.as_posix() for path in changed_files],
        "removed_package_directories": removed_directories,
        "packages": results,
        "verdict": "PASS" if passed else "FAIL",
    }
    if args.json:
        print(json.dumps(report, ensure_ascii=False, sort_keys=True))
    else:
        for result in results:
            coverage = result["coverage_percent"]
            suffix = "N/A" if coverage is None else f"{coverage}%"
            print(f"{result['package']}: {result['status']} ({suffix})")
        for directory in removed_directories:
            print(f"{directory}: package removed (coverage N/A)")
        print(f"verdict: {report['verdict']}")
    return 0 if passed else 1


if __name__ == "__main__":
    raise SystemExit(main())
