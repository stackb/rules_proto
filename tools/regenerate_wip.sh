#!/usr/bin/env bash

set euo pipefail

for LABEL in $(bazel query --config=quiet 'kind(".*_run", //language/nodejs/nodejs_grpc_library/...)' --output label)
do
    bazel run --config=quiet "$LABEL"
done
