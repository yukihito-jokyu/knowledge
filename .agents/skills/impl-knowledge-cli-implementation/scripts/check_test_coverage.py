#!/usr/bin/env python3
"""Require 100% statement coverage for every package with production Go files."""

from __future__ import annotations

import re
import subprocess
import sys
import tempfile
from decimal import Decimal
from pathlib import Path


TOTAL_PATTERN = re.compile(r"^total:\s+\(statements\)\s+([0-9]+(?:\.[0-9]+)?)%$")


def run(*args: str, cwd: Path, capture_output: bool = False) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        args,
        cwd=cwd,
        check=False,
        text=True,
        capture_output=capture_output,
    )


def main() -> int:
    repo_result = run("git", "rev-parse", "--show-toplevel", cwd=Path.cwd(), capture_output=True)
    if repo_result.returncode != 0:
        sys.stderr.write(repo_result.stderr)
        return repo_result.returncode
    repo = Path(repo_result.stdout.strip()).resolve()

    packages_result = run(
        "go",
        "list",
        "-f",
        "{{if .GoFiles}}{{.ImportPath}}{{end}}",
        "./...",
        cwd=repo,
        capture_output=True,
    )
    if packages_result.returncode != 0:
        sys.stderr.write(packages_result.stderr)
        return packages_result.returncode

    packages = [line for line in packages_result.stdout.splitlines() if line]
    for package in packages:
        with tempfile.NamedTemporaryFile(prefix="knowledge-coverage-", suffix=".out") as profile:
            test_result = run(
                "go",
                "test",
                package,
                f"-coverprofile={profile.name}",
                cwd=repo,
            )
            if test_result.returncode != 0:
                return test_result.returncode

            statement_blocks = [
                line
                for line in Path(profile.name).read_text(encoding="utf-8").splitlines()[1:]
                if line
            ]
            if not statement_blocks:
                print(f"{package}: coverable statementがないためcoverage対象外です。")
                continue

            cover_result = run(
                "go",
                "tool",
                "cover",
                f"-func={profile.name}",
                cwd=repo,
                capture_output=True,
            )
            if cover_result.returncode != 0:
                sys.stderr.write(cover_result.stderr)
                return cover_result.returncode

            total = next(
                (
                    match.group(1)
                    for line in cover_result.stdout.splitlines()
                    if (match := TOTAL_PATTERN.match(line))
                ),
                None,
            )
            if total is None:
                print(f"{package}: coverage合計を取得できません。", file=sys.stderr)
                return 1
            if Decimal(total) != Decimal("100"):
                print(
                    f"{package}: テストカバレッジは100%である必要があります。実測値: {total}%",
                    file=sys.stderr,
                )
                return 1

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
