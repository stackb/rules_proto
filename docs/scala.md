---
layout: default
title: scala
permalink: examples/scala
parent: Examples
---


# scala example

[`testdata files`](/example/golden/testdata/scala)


## `Integration Test`

`bazel test @@//example/golden:scala_test`)


## `BUILD.bazel` (before gazelle)

~~~python
~~~


## `BUILD.bazel` (after gazelle)

~~~python
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")
load("@build_stack_rules_proto//rules/scala:grpc_scala_library.bzl", "grpc_scala_library")
load("@build_stack_rules_proto//rules/scala:proto_scala_library.bzl", "proto_scala_library")
load("@rules_proto//proto:defs.bzl", "proto_library")

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
    name = "syntax_grpc_scala_library",
    srcs = ["syntax_akka_grpc.srcjar"],
    visibility = ["//visibility:public"],
    deps = [
        "//lib:scala",
        "//proto:proto_proto_scala_library",
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
        "@scalaapis//scalapb:scalapb_proto_scala_library",
    ],
)

proto_compile(
    name = "syntax_scala_compile",
    options = {"@build_stack_rules_proto//plugin/scalapb/scalapb:protoc-gen-scala-grpc": ["grpc"]},
    outputs = [
        "syntax_akka_grpc.srcjar",
        "syntax_scala.srcjar",
        "syntax_scala_grpc.srcjar",
    ],
    plugins = [
        "@build_stack_rules_proto//plugin/akka/akka-grpc:protoc-gen-akka-grpc",
        "@build_stack_rules_proto//plugin/scalapb/scalapb:protoc-gen-scala",
        "@build_stack_rules_proto//plugin/scalapb/scalapb:protoc-gen-scala-grpc",
    ],
    proto = "syntax_proto",
)

proto_scala_library(
    name = "syntax_proto_scala_library",
    srcs = ["syntax_scala.srcjar"],
    visibility = ["//visibility:public"],
    exports = ["@com_google_protobuf//:protobuf_java"],
    deps = [
        "//lib:scala",
        "//proto:proto_proto_scala_library",
        "@com_google_protobuf//:protobuf_java",
        "@maven_scala//:com_thesamet_scalapb_lenses_2_12",
        "@maven_scala//:com_thesamet_scalapb_scalapb_runtime_2_12",
        "@scalaapis//scalapb:scalapb_proto_scala_library",
    ],
)
~~~


## `MODULE.bazel (snippet)`

~~~python
# ----------------------------------------------------
# scala
# ----------------------------------------------------

load("@build_stack_rules_proto//deps:scala_deps.bzl", "scala_deps")
load("@io_bazel_rules_scala//:scala_config.bzl", "scala_config")
load("@io_bazel_rules_scala//scala:scala.bzl", "scala_repositories")
load("@io_bazel_rules_scala//scala:toolchains.bzl", "scala_register_toolchains")
load(
    "@rules_jvm_external//:defs.bzl",
    "maven_install",
)
load("@build_stack_rules_proto//rules/proto:proto_repository.bzl", "proto_repository")

scala_deps()

scala_config(scala_version = "2.12.18")

scala_repositories()

scala_register_toolchains()

# ----------------------------------------------------
# maven
# ----------------------------------------------------

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

bazel_dep(name = "rules_go", version = "0.57.0", repo_name = "io_bazel_rules_go")

# -------------------------------------------------------------------
# Configuration: Go
# -------------------------------------------------------------------

go_sdk = use_extension("@io_bazel_rules_go//go:extensions.bzl", "go_sdk")
go_sdk.download(version = "1.23.1")

# -------------------------------------------------------------------
# Configuration: protobuf
# -------------------------------------------------------------------

register_toolchains("@build_stack_rules_proto//toolchain:prebuilt")

bazel_dep(name = "gazelle", version = "0.45.0", repo_name = "bazel_gazelle")
bazel_dep(name = "protobuf", version = "32.0", repo_name = "com_google_protobuf")
bazel_dep(name = "rules_jvm_external", version = "6.8")
bazel_dep(name = "rules_scala", version = "7.0.0", repo_name = "io_bazel_rules_scala")

# -------------------------------------------------------------------
# Configuration: Gazelle
# -------------------------------------------------------------------


maven = use_extension("@rules_jvm_external//:extensions.bzl", "maven")
maven.install(
    name = "maven",
    # these artifacts are specified only to disambiguate versions between
    # multiple contributing modules (bazel mod warnings)
    artifacts = [
        "com.google.code.findbugs:jsr305:3.0.2",
        "com.google.code.gson:gson:2.11.0",
        "com.google.errorprone:error_prone_annotations:2.30.0",
        "com.google.guava:guava:32.0.1-jre",
    ],
    known_contributing_modules = [
        "",
        "build_stack_rules_proto",
        "grpc-java",
        "protobuf",
    ],
    # lock_file = "maven_install.json",
)
maven.install(
    name = "maven_scala",
    artifacts = [
        "com.thesamet.scalapb:compilerplugin_2.12:0.11.17",
        "com.thesamet.scalapb:lenses_2.12:0.11.5",
        "com.thesamet.scalapb:scalapb-json4s_2.12:0.12.0",
        "com.thesamet.scalapb:scalapb-runtime_2.12:0.11.5",
        "com.thesamet.scalapb:scalapb-runtime-grpc_2.12:0.11.5",
    ],
    known_contributing_modules = ["", "build_stack_rules_proto"],    
    # lock_file = "rules_jvm_external~~maven~maven_scala_install.json",
    repositories = ["https://repo1.maven.org/maven2"],
)
maven.install(
    name = "maven_akka",
    artifacts = [
        "com.lightbend.akka.grpc:akka-grpc-codegen_2.12:2.1.3",
        "com.lightbend.akka.grpc:akka-grpc-runtime_2.12:2.1.3",
    ],
    fetch_sources = True,
    known_contributing_modules = ["", "build_stack_rules_proto"],    
    # lock_file = "//:maven_akka_install.json",
    repositories = ["https://repo1.maven.org/maven2"],
)
use_repo(
    maven,
    "com_google_protobuf_protobuf_java_3_19_6",
    "com_thesamet_scalapb_compilerplugin_2_12_0_11_13",
    "maven",
    "maven_akka",
    "maven_scala",
)

proto_repository = use_extension("@build_stack_rules_proto//extensions:proto_repository.bzl", "proto_repository", dev_dependency = True)
proto_repository.archive(
    name = "scalaapis",
    build_directives = ["gazelle:proto_language scala enabled true"],
    build_file_generation = "clean",
    build_file_proto_mode = "file",
    cfgs = ["//:config.yaml"],
    sha256 = "1ac039f79b0825fe2e7e5ddf24e330632d63b70a7a42bfd39ded5bb1fb648811",
    # the typical importpath is 'scalapb/scalapb.proto', so strip the prefix up
    # to that directory.
    strip_prefix = "ScalaPB-a4e0e02c0f5b160877d5f97f6902dbec4c633afe/protobuf",
    type = "zip",
    urls = ["https://codeload.github.com/scalapb/ScalaPB/zip/a4e0e02c0f5b160877d5f97f6902dbec4c633afe"],
)
use_repo(
    proto_repository,
    "scalaapis",
)

~~~

