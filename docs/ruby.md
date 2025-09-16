---
layout: default
title: ruby
permalink: examples/ruby
parent: Examples
---


# ruby example

[`testdata files`](/example/golden/testdata/ruby)


## `Integration Test`

`bazel test @@//example/golden:ruby_test`)


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


## `BUILD.bazel` (after gazelle)

~~~python
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")
load("@rules_proto//proto:defs.bzl", "proto_library")

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
    outputs = ["example_pb.rb"],
    plugins = ["@build_stack_rules_proto//plugin/builtin:ruby"],
    proto = "example_proto",
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

~~~

