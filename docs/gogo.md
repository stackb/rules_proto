---
layout: default
title: gogo
permalink: examples/gogo
parent: Examples
---


# gogo example

`bazel test //example/golden:gogo_test`


## `BUILD.bazel` (after gazelle)

~~~python
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")
load("@rules_proto//proto:defs.bzl", "proto_library")

# "proto_rule" instantiates the proto_compile rule
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile

# "proto_plugin" instantiates the builtin gogo plugin
# gazelle:proto_plugin gogo implementation gogo:protobuf:protoc-gen-gogo

# "proto_language" binds the rule(s) and plugin(s) together
# gazelle:proto_language gogo rule proto_compile
# gazelle:proto_language gogo plugin gogo

proto_library(
    name = "example_proto",
    srcs = ["example.proto"],
    visibility = ["//visibility:public"],
)

proto_compile(
    name = "example_gogo_compile",
    options = {"@build_stack_rules_proto//plugin/gogo/protobuf:protoc-gen-gogo": ["plugins=grpc"]},
    outputs = ["example.pb.go"],
    plugins = ["@build_stack_rules_proto//plugin/gogo/protobuf:protoc-gen-gogo"],
    proto = "example_proto",
)
~~~


## `BUILD.bazel` (before gazelle)

~~~python
# "proto_rule" instantiates the proto_compile rule
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile

# "proto_plugin" instantiates the builtin gogo plugin
# gazelle:proto_plugin gogo implementation gogo:protobuf:protoc-gen-gogo

# "proto_language" binds the rule(s) and plugin(s) together
# gazelle:proto_language gogo rule proto_compile
# gazelle:proto_language gogo plugin gogo
~~~


## `WORKSPACE`

~~~python
~~~

