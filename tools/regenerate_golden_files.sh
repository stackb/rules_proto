#!/usr/bin/env bash

set euo pipefail

for LABEL in $(bazel query --config=quiet 'kind(".*_run", //example/... union //rules/proto/... union //docs/...)' --output label)
do
    bazel run --config=quiet "$LABEL"
done
