#!/usr/bin/env python3
"""Fingerprint an in-scope working-tree candidate without changing the repository."""

from __future__ import annotations

import hashlib
import json
import os
import subprocess
import sys
from pathlib import Path


def git(*args: str) -> bytes:
    return subprocess.check_output(("git", *args), stderr=subprocess.DEVNULL)


def main() -> int:
    if len(sys.argv) > 2 or len(sys.argv) == 2 and sys.argv[1] != "--json":
        print("usage: candidate_fingerprint.py [--json]", file=sys.stderr)
        return 2

    repo = Path(git("rev-parse", "--show-toplevel").decode().strip()).resolve()
    head = git("rev-parse", "HEAD").decode().strip()
    tracked_changes = set(
        filter(
            None,
            git("diff", "--name-only", "-z", "--no-renames", "HEAD", "--")
            .decode()
            .split("\0"),
        )
    )
    untracked_changes = set(
        filter(
            None,
            git("ls-files", "--others", "--exclude-standard", "-z").decode().split("\0"),
        )
    )
    changed_paths = sorted(tracked_changes | untracked_changes)
    digest = hashlib.sha256()
    digest.update(b"HEAD\0" + head.encode() + b"\0")
    manifest: list[dict[str, str | int | None]] = []

    for raw_path in changed_paths:
        candidate = (repo / raw_path).resolve()
        try:
            relative = candidate.relative_to(repo)
        except ValueError:
            print(f"repository外のpathは指定できません: {raw_path}", file=sys.stderr)
            return 2

        relative_text = relative.as_posix()
        tracked = subprocess.run(
            ("git", "ls-files", "--error-unmatch", "--", relative_text),
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
            check=False,
        ).returncode == 0

        digest.update(b"PATH\0" + relative_text.encode() + b"\0")
        if not candidate.exists():
            if not tracked:
                print(f"存在しない未追跡pathです: {relative_text}", file=sys.stderr)
                return 2
            digest.update(b"DELETED\0")
            manifest.append({"path": relative_text, "state": "deleted", "mode": None, "sha256": None})
            continue
        if not candidate.is_file():
            print(f"file以外は指定できません: {relative_text}", file=sys.stderr)
            return 2

        mode = os.stat(candidate).st_mode & 0o777
        content = candidate.read_bytes()
        content_hash = hashlib.sha256(content).hexdigest()
        digest.update(f"MODE\0{mode:o}\0".encode())
        digest.update(b"CONTENT\0" + content + b"\0")
        manifest.append(
            {
                "path": relative_text,
                "state": "untracked" if relative_text in untracked_changes else "tracked_change",
                "mode": mode,
                "sha256": content_hash,
            }
        )

    candidate_id = digest.hexdigest()
    if len(sys.argv) == 2:
        print(
            json.dumps(
                {"candidate_id": candidate_id, "head": head, "files": manifest},
                ensure_ascii=False,
                sort_keys=True,
            )
        )
    else:
        print(candidate_id)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
