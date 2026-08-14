#!/usr/bin/env python3
"""Fingerprint the explicit normative sources listed in a Verification Packet."""

from __future__ import annotations

import hashlib
import subprocess
import sys
from pathlib import Path


def main() -> int:
    if len(sys.argv) < 2:
        print("usage: source_fingerprint.py <repo-relative-source>...", file=sys.stderr)
        return 2

    repo = Path(
        subprocess.check_output(("git", "rev-parse", "--show-toplevel"), stderr=subprocess.DEVNULL)
        .decode()
        .strip()
    ).resolve()
    digest = hashlib.sha256()

    inventory: set[tuple[str, str]] = set()
    source_files: set[Path] = set()
    for raw_path in sorted(set(sys.argv[1:])):
        source = repo / raw_path
        resolved = source.resolve()
        try:
            resolved.relative_to(repo)
        except ValueError:
            print(f"repository外のpathは指定できません: {raw_path}", file=sys.stderr)
            return 2
        relative_root = resolved.relative_to(repo).as_posix()
        if resolved.is_dir():
            inventory.add(("ROOT", relative_root))
            directory_files = {path.resolve() for path in resolved.rglob("*") if path.is_file()}
            if not directory_files:
                print(f"正規source directoryが空です: {raw_path}", file=sys.stderr)
                return 2
            source_files.update(directory_files)
        elif resolved.is_file():
            inventory.add(("FILE", relative_root))
            source_files.add(resolved)
        else:
            print(f"正規source pathが存在しません: {raw_path}", file=sys.stderr)
            return 2

    for source_type, relative in sorted(inventory):
        digest.update(source_type.encode() + b"\0" + relative.encode() + b"\0")

    for source in sorted(source_files):
        try:
            source.resolve().relative_to(repo)
        except ValueError:
            print(f"repository外を参照するsourceです: {source}", file=sys.stderr)
            return 2
        relative = source.relative_to(repo).as_posix()
        digest.update(b"PATH\0" + relative.encode() + b"\0")
        digest.update(b"CONTENT\0" + source.read_bytes() + b"\0")

    print(digest.hexdigest())
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
