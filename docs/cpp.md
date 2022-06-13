---
layout: default
title: cpp
permalink: examples/cpp
parent: Examples
---


# cpp example

`bazel test //example/golden:cpp_test`


## `BUILD.bazel` (after gazelle)

~~~python
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")

# "proto_rule" instantiates the proto_compile rule
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile

# "proto_plugin" instantiates the builtin cpp plugin
# gazelle:proto_plugin cpp implementation builtin:cpp

# "proto_language" binds the rule(s) and plugin(s) together
# gazelle:proto_language cpp rule proto_compile
# gazelle:proto_language cpp plugin cpp

proto_library(
    name = "example_proto",
    srcs = ["example.proto"],
    visibility = ["//visibility:public"],
)

proto_compile(
    name = "example_cpp_compile",
    outputs = [
        "example.pb.cc",
        "example.pb.h",
    ],
    plugins = ["@build_stack_rules_proto//plugin/builtin:cpp"],
    proto = "example_proto",
)
~~~


## `BUILD.bazel` (before gazelle)

~~~python
# "proto_rule" instantiates the proto_compile rule
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile

# "proto_plugin" instantiates the builtin cpp plugin
# gazelle:proto_plugin cpp implementation builtin:cpp

# "proto_language" binds the rule(s) and plugin(s) together
# gazelle:proto_language cpp rule proto_compile
# gazelle:proto_language cpp plugin cpp
~~~


## `WORKSPACE`

~~~python
~~~

