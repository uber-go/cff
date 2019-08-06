#!/bin/bash

set -euo pipefail

PROJECT_ROOT=${PROJECT_ROOT:-$(git rev-parse --show-toplevel)}
CURRENT_PATH=$(pwd)
BAZEL_CMD="${PROJECT_ROOT}/tools/bazel"
export BAZEL_CMD

echo "Running CFFv2 tool over all cff rules under $CURRENT_PATH..."

bazel query 'kind(cff, ...)' --output xml |
  bazel run //src/go.uber.org/cff/cmd/cff-bazelxml-to-cffcmd:bin -- --project-root="$PROJECT_ROOT" "$@"
