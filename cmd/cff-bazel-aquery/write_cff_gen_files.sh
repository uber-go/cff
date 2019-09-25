#!/bin/bash

set -euo pipefail
set -x

PROJECT_ROOT=${PROJECT_ROOT:-$(git rev-parse --show-toplevel)}

echo "Running CFFv2 tool over all cff rules under $(pwd)..."

for rule in $(bazel query 'kind(cff, ...)'); do
  bazel build "$rule"
  bazel aquery "$rule" --output=proto | bazel run //src/go.uber.org/cff/cmd/cff-bazel-aquery:bin -- "$@" "$PROJECT_ROOT"
done

