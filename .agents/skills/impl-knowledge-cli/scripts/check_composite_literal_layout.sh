#!/bin/sh

set -eu

go run "$(dirname "$0")/check_composite_literal_layout.go"
