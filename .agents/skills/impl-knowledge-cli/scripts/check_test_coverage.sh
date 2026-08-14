#!/bin/sh

set -eu

sh "$(dirname "$0")/check_composite_literal_layout.sh"

coverage_file=$(mktemp "${TMPDIR:-/tmp}/knowledge-coverage.XXXXXX")
trap 'rm -f "$coverage_file"' EXIT

go test ./... -coverpkg=./... -coverprofile="$coverage_file"

coverage=$(go tool cover -func="$coverage_file" | awk '/^total:/ { sub(/%/, "", $3); print $3 }')
if [ "$coverage" != "100.0" ]; then
	printf 'テストカバレッジは100%%である必要があります。実測値: %s%%\n' "$coverage" >&2
	exit 1
fi
