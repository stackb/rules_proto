---
layout: default
title: scala
permalink: examples/scala
parent: Examples
---


# scala example

`bazel test //example/golden:scala_test`


## `BUILD.bazel` (after gazelle)

~~~python
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@build_stack_rules_proto//rules/scala:grpc_scala_library.bzl", "grpc_scala_library")
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")

proto_library(
    name = "syntax_proto",
    srcs = [
        "noun.proto",
        "package.proto",
        "pronoun.proto",
        "service.proto",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//proto:proto_proto",
        "@scalaapis//scalapb:scalapb_proto",
    ],
)

grpc_scala_library(
    name = "syntax_scala_library",
    srcs = [
        "syntax_akka_grpc.srcjar",
        "syntax_scala.srcjar",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//lib:scala",
        "//proto:proto_scala_library",
        "@com_google_protobuf//:protobuf_java",
        "@maven_akka//:com_lightbend_akka_grpc_akka_grpc_runtime_2_12",
        "@maven_akka//:com_typesafe_akka_akka_actor_2_12",
        "@maven_akka//:com_typesafe_akka_akka_http_core_2_12",
        "@maven_akka//:com_typesafe_akka_akka_stream_2_12",
        "@maven_scala//:com_thesamet_scalapb_lenses_2_12",
        "@maven_scala//:com_thesamet_scalapb_scalapb_runtime_2_12",
        "@maven_scala//:com_thesamet_scalapb_scalapb_runtime_grpc_2_12",
        "@maven_scala//:io_grpc_grpc_api",
        "@maven_scala//:io_grpc_grpc_protobuf",
        "@maven_scala//:io_grpc_grpc_stub",
        "@scalaapis//scalapb:scalapb_scala_library",
    ],
)

proto_compile(
    name = "syntax_scala_compile",
    options = {"@build_stack_rules_proto//plugin/scalapb/scalapb:protoc-gen-scala": ["grpc"]},
    outputs = [
        "syntax_akka_grpc.srcjar",
        "syntax_scala.srcjar",
    ],
    plugins = [
        "@build_stack_rules_proto//plugin/akka/akka-grpc:protoc-gen-akka-grpc",
        "@build_stack_rules_proto//plugin/scalapb/scalapb:protoc-gen-scala",
    ],
    proto = "syntax_proto",
)
~~~


## `BUILD.bazel` (before gazelle)

~~~python
~~~


## `WORKSPACE`

~~~python
# ----------------------------------------------------
# scala
# ----------------------------------------------------

load("@build_stack_rules_proto//deps:scala_deps.bzl", "scala_deps")

scala_deps()

load("@io_bazel_rules_scala//:scala_config.bzl", "scala_config")

scala_config(scala_version = "2.12.11")

load("@io_bazel_rules_scala//scala:scala.bzl", "scala_repositories")

scala_repositories()

load("@io_bazel_rules_scala//scala:toolchains.bzl", "scala_register_toolchains")

scala_register_toolchains()

# ----------------------------------------------------
# maven
# ----------------------------------------------------

load(
    "@rules_jvm_external//:defs.bzl",
    "maven_install",
)

maven_install(
    name = "maven_scala",
    artifacts = [
        "com.thesamet.scalapb:lenses_2.12:0.11.10",
        "com.thesamet.scalapb:scalapb-json4s_2.12:0.12.0",
        "com.thesamet.scalapb:scalapb-runtime_2.12:0.11.10",
        "com.thesamet.scalapb:scalapb-runtime-grpc_2.12:0.11.10",
        "com.thesamet.scalapb:scalapbc_2.12:0.11.10",
        "org.json4s:json4s-core_2.12:4.0.3",
    ],
    fetch_sources = True,
    repositories = ["https://repo1.maven.org/maven2"],
)

# ----------------------------------------------------
# akka
# ----------------------------------------------------

maven_install(
    name = "maven_akka",
    artifacts = [
        "com.lightbend.akka.grpc:akka-grpc-codegen_2.12:2.1.3",
        "com.lightbend.akka.grpc:akka-grpc-runtime_2.12:2.1.3",
    ],
    fetch_sources = True,
    repositories = ["https://repo1.maven.org/maven2"],
)

# ----------------------------------------------------
# proto_repository
# ----------------------------------------------------

load("@build_stack_rules_proto//rules/proto:proto_repository.bzl", "proto_repository")

proto_repository(
    name = "scalaapis",
    build_directives = ["gazelle:proto_language scala enabled true"],
    build_file_generation = "on",
    build_file_proto_mode = "file",
    cfgs = ["//:config.yaml"],
    sha256 = "1ac039f79b0825fe2e7e5ddf24e330632d63b70a7a42bfd39ded5bb1fb648811",
    # the typical importpath is 'scalapb/scalapb.proto', so strip the prefix up
    # to that directory.
    strip_prefix = "ScalaPB-a4e0e02c0f5b160877d5f97f6902dbec4c633afe/protobuf",
    type = "zip",
    urls = ["https://codeload.github.com/scalapb/ScalaPB/zip/a4e0e02c0f5b160877d5f97f6902dbec4c633afe"],
)
~~~

