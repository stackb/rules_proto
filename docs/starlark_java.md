---
layout: default
title: starlark_java
permalink: examples/starlark_java
parent: Examples
---


# starlark_java example

[`testdata files`](/example/golden/testdata/starlark_java)


## `Integration Test`

`bazel test @@//example/golden:starlark_java_test`)


## `BUILD.bazel` (before gazelle)

~~~python
~~~


## `BUILD.bazel` (after gazelle)

~~~python
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")
load("@rules_java//java:defs.bzl", "java_library")
load("@rules_proto//proto:defs.bzl", "proto_library")
load("//:defs.bzl", "java_wrapper")

proto_library(
    name = "order_proto",
    srcs = ["order.proto"],
    visibility = ["//visibility:public"],
    deps = ["//customer:customer_proto"],
)

java_library(
    name = "order_java_library",
    srcs = ["order.srcjar"],
    visibility = ["//visibility:public"],
    deps = [
        "//customer:customer_java_library",
        "@com_google_protobuf//:protobuf_java",
        "@com_google_protobuf//java/core",
    ],
)

java_wrapper(
    name = "order_java_wrap",
    javalib = "order_java_library",
    deps = ["//customer:customer_java_wrap"],
)

proto_compile(
    name = "order_java_compile",
    outs = {"@build_stack_rules_proto//plugin/builtin:java": "order/order.srcjar"},
    outputs = ["order.srcjar"],
    plugins = ["@build_stack_rules_proto//plugin/builtin:java"],
    proto = "order_proto",
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

register_toolchains("@build_stack_rules_proto//toolchain:standard")

bazel_dep(name = "protobuf", version = "32.0", repo_name = "com_google_protobuf")

~~~

