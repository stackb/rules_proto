#!/usr/bin/env bash

for LABEL in $(bazel query --config=quiet 'kind(".*_test", //example/...)' --output label | grep golden)
do
    bazel run --config=quiet "$LABEL"
done
