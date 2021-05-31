---
layout: default
title: ruby
permalink: examples/ruby
parent: Examples
---


# ruby example

`bazel test //example/golden:ruby_test`


## `BUILD.bazel` (after gazelle)

~~~python
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")

# "proto_rule" instantiates the proto_compile rule
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile

# "proto_plugin" instantiates the builtin ruby plugin
# gazelle:proto_plugin ruby implementation builtin:ruby

# "proto_language" binds the rule(s) and plugin(s) together
# gazelle:proto_language ruby rule proto_compile
# gazelle:proto_language ruby plugin ruby

proto_library(
    name = "example_proto",
    srcs = ["example.proto"],
    visibility = ["//visibility:public"],
)

proto_compile(
    name = "example_ruby_compile",
    outputs = [
        "example_pb.rb",
    ],
    plugins = ["@build_stack_rules_proto//plugin/builtin:ruby"],
    proto = "example_proto",
)
~~~


## `BUILD.bazel` (before gazelle)

~~~python
# "proto_rule" instantiates the proto_compile rule
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile

# "proto_plugin" instantiates the builtin ruby plugin
# gazelle:proto_plugin ruby implementation builtin:ruby

# "proto_language" binds the rule(s) and plugin(s) together
# gazelle:proto_language ruby rule proto_compile
# gazelle:proto_language ruby plugin ruby
~~~


## `WORKSPACE`

~~~python
~~~

