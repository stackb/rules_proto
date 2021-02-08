#!/usr/bin/env bash

set euo pipefail

for LABEL in $(bazel query --config=quiet 'kind(".*_run", //example/... union //language/... union //docs/...)' --output label)
do
    bazel run --config=quiet "$LABEL"
done
