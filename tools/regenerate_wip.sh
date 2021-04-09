#!/usr/bin/env bash

set euo pipefail

for LABEL in $(bazel query --config=quiet 'kind(".*_run", //rules/proto/gogo/...)' --output label)
do
    bazel run --config=quiet "$LABEL"
done
