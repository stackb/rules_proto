---
layout: default
title: builtins
permalink: examples/builtins
parent: examples
---


# builtins example

`bazel test //example/golden:builtins_test`


## `BUILD.bazel` (after gazelle)

~~~python
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")

# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile
# gazelle:proto_plugin cpp implementation builtin:cpp
# gazelle:proto_plugin java implementation builtin:java
# gazelle:proto_plugin closurejs implementation builtin:js:closure
# gazelle:proto_plugin commonjs implementation builtin:js:common
# gazelle:proto_plugin python implementation builtin:python
# gazelle:proto_plugin ruby implementation builtin:ruby
# gazelle:proto_plugin objc implementation builtin:objc
# gazelle:proto_language builtins rule proto_compile
# gazelle:proto_language builtins plugin cpp
# gazelle:proto_language builtins plugin java
# gazelle:proto_language builtins plugin closurejs
# gazelle:proto_language builtins plugin commonjs
# gazelle:proto_language builtins plugin python
# gazelle:proto_language builtins plugin ruby
# gazelle:proto_language builtins plugin objc

proto_library(
    name = "test_proto",
    srcs = ["test.proto"],
    visibility = ["//visibility:public"],
)

proto_compile(
    name = "test_builtins_compile",
    outs = {"@build_stack_rules_proto//plugin/builtin:java": "test.srcjar"},
    options = {
        "@build_stack_rules_proto//plugin/builtin:closurejs": [
            "import_style=closure",
            "library=test_closure",
        ],
        "@build_stack_rules_proto//plugin/builtin:commonjs": ["import_style=commonjs"],
    },
    outputs = [
        "Test.pbobjc.h",
        "Test.pbobjc.m",
        "test.pb.cc",
        "test.pb.h",
        "test.srcjar",
        "test_closure.js",
        "test_pb.js",
        "test_pb.rb",
        "test_pb2.py",
    ],
    plugins = [
        "@build_stack_rules_proto//plugin/builtin:closurejs",
        "@build_stack_rules_proto//plugin/builtin:commonjs",
        "@build_stack_rules_proto//plugin/builtin:cpp",
        "@build_stack_rules_proto//plugin/builtin:java",
        "@build_stack_rules_proto//plugin/builtin:objc",
        "@build_stack_rules_proto//plugin/builtin:python",
        "@build_stack_rules_proto//plugin/builtin:ruby",
    ],
    proto = "test_proto",
)
~~~


## `BUILD.bazel` (before gazelle)

~~~python
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile
# gazelle:proto_plugin cpp implementation builtin:cpp
# gazelle:proto_plugin java implementation builtin:java
# gazelle:proto_plugin closurejs implementation builtin:js:closure
# gazelle:proto_plugin commonjs implementation builtin:js:common
# gazelle:proto_plugin python implementation builtin:python
# gazelle:proto_plugin ruby implementation builtin:ruby
# gazelle:proto_plugin objc implementation builtin:objc
# gazelle:proto_language builtins rule proto_compile
# gazelle:proto_language builtins plugin cpp
# gazelle:proto_language builtins plugin java
# gazelle:proto_language builtins plugin closurejs
# gazelle:proto_language builtins plugin commonjs
# gazelle:proto_language builtins plugin python
# gazelle:proto_language builtins plugin ruby
# gazelle:proto_language builtins plugin objc
~~~


## `WORKSPACE`

~~~python
local_repository(
    name = "build_stack_rules_proto",
    path = "../build_stack_rules_proto",
)

register_toolchains("@build_stack_rules_proto//toolchain:standard")

# == Externals ==

load("@build_stack_rules_proto//deps:core_deps.bzl", "core_deps")

core_deps()

# == Go ==

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.18.2")

# == Gazelle ==

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

# == Protobuf ==

load("@build_stack_rules_proto//deps:protobuf_core_deps.bzl", "protobuf_core_deps")

protobuf_core_deps()
~~~

