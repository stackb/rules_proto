#!/usr/bin/env bash

for LABEL in $(bazel query --config=quiet 'kind(".*_run", //example/...)' --output label)
do
    bazel run --config=quiet "$LABEL"
done
