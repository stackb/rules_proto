---
layout: default
title: starlark_java
permalink: examples/starlark_java
parent: Examples
---


# starlark_java example

`bazel test //example/golden:starlark_java_test`


## `BUILD.bazel` (after gazelle)

~~~python
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")
load("@rules_java//java:defs.bzl", "java_library")

# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile
# gazelle:proto_rule java_library implementation lib/rules.star%java_library
# gazelle:proto_rule java_library deps @com_google_protobuf//:protobuf_java
# gazelle:proto_rule java_library deps @com_google_protobuf//java/core
# gazelle:proto_rule java_library visibility //visibility:public
# gazelle:proto_plugin java implementation lib/plugins.star%java
# gazelle:proto_language java rule proto_compile
# gazelle:proto_language java rule java_library
# gazelle:proto_language java plugin java

proto_library(
    name = "example_proto",
    srcs = ["example.proto"],
    visibility = ["//visibility:public"],
)

java_library(
    name = "example_java_library",
    srcs = ["example.srcjar"],
    visibility = ["//visibility:public"],
    deps = [
        "@com_google_protobuf//:protobuf_java",
        "@com_google_protobuf//java/core",
    ],
)

proto_compile(
    name = "example_java_compile",
    outputs = ["example.srcjar"],
    plugins = ["@build_stack_rules_proto//plugin/builtin:java"],
    proto = "example_proto",
)
~~~


## `BUILD.bazel` (before gazelle)

~~~python
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile
# gazelle:proto_rule java_library implementation lib/rules.star%java_library
# gazelle:proto_rule java_library deps @com_google_protobuf//:protobuf_java
# gazelle:proto_rule java_library deps @com_google_protobuf//java/core
# gazelle:proto_rule java_library visibility //visibility:public
# gazelle:proto_plugin java implementation lib/plugins.star%java
# gazelle:proto_language java rule proto_compile
# gazelle:proto_language java rule java_library
# gazelle:proto_language java plugin java
~~~


## `WORKSPACE`

~~~python
# Using scala_deps here because it provides rules_jvm_external, there is nothing
# scala-specific about this example.
load("@build_stack_rules_proto//deps:scala_deps.bzl", "scala_deps")

scala_deps()
~~~

