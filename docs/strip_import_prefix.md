---
layout: default
title: strip_import_prefix
permalink: examples/strip_import_prefix
parent: Examples
---


# strip_import_prefix example

[`testdata files`](/example/golden/testdata/strip_import_prefix)


## `Integration Test`

`bazel test @@//example/golden:strip_import_prefix_test`)


## `BUILD.bazel` (before gazelle)

~~~python
# gazelle:proto_strip_import_prefix /module_lib/util/nested
~~~


## `BUILD.bazel` (after gazelle)

~~~python
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")
load("@build_stack_rules_proto//rules/py:proto_py_library.bzl", "proto_py_library")
load("@rules_proto//proto:defs.bzl", "proto_library")

# gazelle:proto_strip_import_prefix /module_lib/util/nested

proto_library(
    name = "prefix_test_proto",
    srcs = ["test.proto"],
    strip_import_prefix = "/module_lib/util/nested",
    visibility = ["//visibility:public"],
)

proto_compile(
    name = "prefix_test_python_compile",
    output_mappings = ["test_pb2.py=/prefix/test_pb2.py"],
    outputs = ["test_pb2.py"],
    plugins = ["@build_stack_rules_proto//plugin/builtin:python"],
    proto = "prefix_test_proto",
)

proto_py_library(
    name = "prefix_test_py_library",
    srcs = ["test_pb2.py"],
    imports = [".."],
    visibility = ["//visibility:public"],
    deps = ["@com_google_protobuf//:protobuf_python"],
)
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

register_toolchains("@build_stack_rules_proto//toolchain:prebuilt")

bazel_dep(name = "protobuf", version = "32.0", repo_name = "com_google_protobuf")

~~~

