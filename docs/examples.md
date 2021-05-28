---
layout: default
title: examples
permalink: examples
has_children: true
---

See [examples/golden/testdata](examples/golden/testdata) for a complete list of example test cases.

Each one is tested as both a golden file test (`BUILD.in`, `BUILD.out`) as well
as a
[go_bazel_test](https://github.com/bazelbuild/rules_go/blob/4cd45a2ac59bd00ba54d23ebbdb7e5e2aed69007/go/tools/bazel_testing/def.bzl#L17)
to ensure the configuration is functional.