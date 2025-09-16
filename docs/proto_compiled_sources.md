---
layout: default
title: proto_compiled_sources
permalink: examples/proto_compiled_sources
parent: Examples
---


# proto_compiled_sources example

[`testdata files`](/example/golden/testdata/proto_compiled_sources)


## `Integration Test`

`bazel test @@//example/golden:proto_compiled_sources_test`)


## `BUILD.bazel` (before gazelle)

~~~python
# gazelle:proto_strip_import_prefix /src
~~~


## `BUILD.bazel` (after gazelle)

~~~python
load("@build_stack_rules_proto//rules:proto_compiled_sources.bzl", "proto_compiled_sources")
load("@rules_proto//proto:defs.bzl", "proto_library")

# gazelle:proto_strip_import_prefix /src

proto_library(
    name = "svc_proto",
    srcs = ["svc.proto"],
    strip_import_prefix = "/src",
    visibility = ["//visibility:public"],
)

proto_compiled_sources(
    name = "svc_python_compiled_sources",
    srcs = ["svc_pb2.py"],
    output_mappings = ["svc_pb2.py=/idl/svc_pb2.py"],
    plugins = ["@build_stack_rules_proto//plugin/builtin:python"],
    proto = "svc_proto",
    visibility = ["//visibility:public"],
)~~~


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

