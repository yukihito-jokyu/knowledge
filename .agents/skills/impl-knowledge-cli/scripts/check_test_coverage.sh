#!/bin/sh

set -eu

sh "$(dirname "$0")/check_composite_literal_layout.sh"

go list -f '{{if and .GoFiles .TestGoFiles}}{{.ImportPath}}{{end}}' ./... |
while IFS= read -r package; do
	[ -n "$package" ] || continue
	coverage_file=$(mktemp "${TMPDIR:-/tmp}/knowledge-coverage.XXXXXX")
	trap 'rm -f "$coverage_file"' EXIT HUP INT TERM
	go test "$package" -coverprofile="$coverage_file"
	coverage=$(go tool cover -func="$coverage_file" | awk '/^total:/ { sub(/%/, "", $3); print $3 }')
	rm -f "$coverage_file"
	trap - EXIT HUP INT TERM
	if [ "$coverage" != "100.0" ]; then
		printf '%s: テストカバレッジは100%%である必要があります。実測値: %s%%\n' "$package" "$coverage" >&2
		exit 1
	fi
done
