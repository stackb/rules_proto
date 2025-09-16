---
layout: default
title: goembed
permalink: examples/goembed
parent: Examples
---


# goembed example

[`testdata files`](/example/golden/testdata/goembed)


## `Integration Test`

`bazel test @@//example/golden:goembed_test`)


## `BUILD.bazel` (before gazelle)

~~~python
~~~


## `BUILD.bazel` (after gazelle)

~~~python
~~~


## `MODULE.bazel (snippet)`

~~~python

bazel_dep(name = "rules_go", version = "0.57.0", repo_name = "io_bazel_rules_go")

# -------------------------------------------------------------------
# Configuration: Go
# -------------------------------------------------------------------

go_sdk = use_extension("@io_bazel_rules_go//go:extensions.bzl", "go_sdk")
go_sdk.download(version = "1.23.1")

# -------------------------------------------------------------------
# Configuration: protobuf
# -------------------------------------------------------------------

register_toolchains("@build_stack_rules_proto//toolchain:standard")

bazel_dep(name = "protobuf", version = "32.0", repo_name = "com_google_protobuf")

~~~

