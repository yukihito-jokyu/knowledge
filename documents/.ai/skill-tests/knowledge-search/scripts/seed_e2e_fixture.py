#!/usr/bin/env python3
"""Seed one isolated real Knowledge CLI store for a knowledge-search E2E case."""

import argparse
import json
import os
import subprocess
import sys


def invoke(knowledge, home, *args):
    env = os.environ.copy()
    env["HOME"] = home
    result = subprocess.run([knowledge, *args], capture_output=True, text=True, env=env)
    if result.returncode != 0:
        raise RuntimeError(f"{args[0]} failed ({result.returncode}): {result.stderr.strip()}")
    payload = json.loads(result.stdout)
    if not payload.get("ok"):
        raise RuntimeError(f"{args[0]} returned error: {result.stdout.strip()}")
    return payload["data"]


def create(knowledge, home, text, evidence_kind, evidence_text, *extra):
    data = invoke(knowledge, home, "create", "--normalized-text", text,
                  "--evidence-kind", evidence_kind, "--evidence-text", evidence_text,
                  "--evidence-observed-at", "2026-08-15T00:00:00Z", *extra)
    return data["assertion_id"]


def seed_known(knowledge, home):
    return {"claim": "Goのcontext cancellationは長時間実行を中断できる", "assertion_id": create(knowledge, home, "Goのcontext cancellationは長時間実行を中断できる", "user_explanation", "context cancellationを監視する長時間実行は中断できる", "--concept", "Go", "--concept", "context")}


def seed_partially_known(knowledge, home):
    return {"claim": "cancellation stops long task and releases resources", "assertion_id": create(knowledge, home, "cancellation stops long task", "user_explanation", "Cancellation stops a running long task", "--concept", "cancellation", "--concept", "task")}


def seed_inferable(knowledge, home):
    consequence = create(knowledge, home, "stopped task becomes complete", "user_reasoning", "A stopped task is treated as complete", "--concept", "task")
    premise = create(knowledge, home, "cancel signal stops task", "user_reasoning", "A task that receives a cancel signal stops", "--concept", "cancel", "--concept", "task", "--relation-type", "causes", "--relation-target-kind", "assertion", "--relation-target-id", consequence)
    return {"claim": "cancel signal causes task to become complete", "assertion_id": premise}


def seed_contradicted(knowledge, home):
    assertion_id = create(knowledge, home, "cancellation cannot stop long task", "user_explanation", "Cancellation cannot stop the long task", "--concept", "cancellation")
    invoke(knowledge, home, "attach-evidence", "--assertion-id", assertion_id, "--evidence-kind", "correction", "--evidence-text", "Correction: cancellation stops a long task when the task observes it", "--evidence-observed-at", "2026-08-15T00:00:00Z")
    return {"claim": "cancellation cannot stop long task", "assertion_id": assertion_id}


def seed_outdated(knowledge, home):
    assertion_id = create(knowledge, home, "Tool X is beta", "self_report", "Tool X was beta on 2025-01-01", "--concept", "Tool X", "--version-scope", "2025", "--valid-from", "2025-01-01T00:00:00Z", "--valid-until", "2025-12-31T00:00:00Z")
    invoke(knowledge, home, "attach-evidence", "--assertion-id", assertion_id, "--evidence-kind", "correction", "--evidence-text", "Correction: Tool X became stable on 2026-08-01", "--evidence-observed-at", "2026-08-01T00:00:00Z")
    return {"claim": "Tool X is beta", "assertion_id": assertion_id}


def seed_uncertain(knowledge, home):
    assertion_id = create(knowledge, home, "feature Y is enabled", "user_explanation", "Observation A: feature Y is enabled", "--concept", "feature Y")
    invoke(knowledge, home, "attach-evidence", "--assertion-id", assertion_id, "--evidence-kind", "correction", "--evidence-text", "Correction: observation B says feature Y is not enabled", "--evidence-observed-at", "2026-08-15T00:00:00Z")
    return {"claim": "feature Y is enabled", "assertion_id": assertion_id}


CASES = {"known": seed_known, "partially_known": seed_partially_known, "inferable": seed_inferable, "contradicted": seed_contradicted, "outdated": seed_outdated, "uncertain": seed_uncertain, "no_evidence": lambda _knowledge, _home: {"claim": "根拠のない独立した検証用Claim"}}


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--knowledge", required=True)
    parser.add_argument("--home", required=True)
    parser.add_argument("--case", choices=sorted(CASES), required=True)
    args = parser.parse_args()
    os.makedirs(args.home, exist_ok=True)
    print(json.dumps(CASES[args.case](args.knowledge, args.home), ensure_ascii=False))


if __name__ == "__main__":
    try:
        main()
    except Exception as error:
        print(str(error), file=sys.stderr)
        sys.exit(1)
