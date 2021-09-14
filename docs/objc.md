---
layout: default
title: objc
permalink: examples/objc
parent: Examples
---


# objc example

`bazel test //example/golden:objc_test`


## `BUILD.bazel` (after gazelle)

~~~python
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")
load("@rules_proto//proto:defs.bzl", "proto_library")

# "proto_rule" instantiates the proto_compile rule
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile

# "proto_plugin" instantiates the builtin objc plugin
# gazelle:proto_plugin objc implementation builtin:objc

# "proto_language" binds the rule(s) and plugin(s) together
# gazelle:proto_language objc rule proto_compile
# gazelle:proto_language objc plugin objc

proto_library(
    name = "example_proto",
    srcs = ["example.proto"],
    visibility = ["//visibility:public"],
)

proto_compile(
    name = "example_objc_compile",
    outputs = [
        "Example.pbobjc.h",
        "Example.pbobjc.m",
    ],
    plugins = ["@build_stack_rules_proto//plugin/builtin:objc"],
    proto = "example_proto",
)
~~~


## `BUILD.bazel` (before gazelle)

~~~python
# "proto_rule" instantiates the proto_compile rule
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile

# "proto_plugin" instantiates the builtin objc plugin
# gazelle:proto_plugin objc implementation builtin:objc

# "proto_language" binds the rule(s) and plugin(s) together
# gazelle:proto_language objc rule proto_compile
# gazelle:proto_language objc plugin objc
~~~


## `WORKSPACE`

~~~python
~~~

