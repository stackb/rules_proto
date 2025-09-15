---
layout: default
title: proto_repository
permalink: examples/proto_repository
parent: Examples
---


# proto_repository example

[`testdata files`](/example/golden/testdata/proto_repository)


## `Integration Test`

`bazel test @@//example/golden:proto_repository_test`)


## `BUILD.bazel` (before gazelle)

~~~python
~~~


## `BUILD.bazel` (after gazelle)

~~~python
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")
load("@build_stack_rules_proto//rules/go:proto_go_library.bzl", "proto_go_library")
load("@rules_proto//proto:defs.bzl", "proto_library")

proto_library(
    name = "app_proto",
    srcs = ["app.proto"],
    visibility = ["//visibility:public"],
    deps = [
        "@googleapis//google/api:annotations_proto",
        "@googleapis//google/api:field_behavior_proto",
    ],
)

proto_compile(
    name = "app_go_compile",
    output_mappings = [
        "app.pb.go=github.com/example/app/app.pb.go",
        "app_grpc.pb.go=github.com/example/app/app_grpc.pb.go",
    ],
    outputs = [
        "app.pb.go",
        "app_grpc.pb.go",
    ],
    plugins = [
        "@build_stack_rules_proto//plugin/golang/protobuf:protoc-gen-go",
        "@build_stack_rules_proto//plugin/grpc/grpc-go:protoc-gen-go-grpc",
    ],
    proto = "app_proto",
    visibility = ["//visibility:public"],
)

proto_go_library(
    name = "app_go_proto",
    srcs = [
        "app.pb.go",
        "app_grpc.pb.go",
    ],
    importpath = "github.com/example/app",
    visibility = ["//visibility:public"],
    deps = [
        "@googleapis//google/api:annotations_go_proto",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//runtime/protoimpl",
    ],
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

register_toolchains("@build_stack_rules_proto//toolchain:prebuilt")

bazel_dep(name = "gazelle", version = "0.45.0", repo_name = "bazel_gazelle")

# -------------------------------------------------------------------
# Configuration: Go
# -------------------------------------------------------------------

go_deps = use_extension("@bazel_gazelle//:extensions.bzl", "go_deps")
go_deps.module(
    path = "google.golang.org/protobuf",
    sum = "h1:OgPcDAFKHnH8X3O4WcO4XUc8GRDeKsKReqbQtiCj7N8=",
    version = "v1.36.6",
)
go_deps.module(
    path = "google.golang.org/grpc",
    sum = "h1:OgPcDAFKHnH8X3O4WcO4XUc8GRDeKsKReqbQtiCj7N8=",
    version = "v1.67.3",
)
use_repo(
    go_deps,
    "org_golang_google_grpc",
    "org_golang_google_protobuf",
)

# -------------------------------------------------------------------
# Configuration: Protobuf
# -------------------------------------------------------------------

proto_repository = use_extension("@build_stack_rules_proto//extensions:proto_repository.bzl", "proto_repository", dev_dependency = True)
proto_repository.archive(
    name = "protoapis",
    build_directives = [
        "gazelle:exclude testdata",
        "gazelle:exclude google/protobuf/compiler/ruby",
        "gazelle:exclude google/protobuf/util/internal/testdata",
        "gazelle:proto_language go enable true",
    ],
    build_file_generation = "clean",
    build_file_proto_mode = "file",
    cfgs = ["//:config.yaml"],
    deleted_files = [
        "google/protobuf/*test*.proto",
        "google/protobuf/*unittest*.proto",
        "google/protobuf/compiler/cpp/*test*.proto",
        "google/protobuf/util/*test*.proto",
        "google/protobuf/util/*unittest*.proto",
        "google/protobuf/util/json_format*.proto",
    ],
    sha256 = "4514213c25a5b87e1948aeeb4c40effc55d11d60871ca5b903a2779005fc48ce",
    strip_prefix = "protobuf-9650e9fe8f737efcad485c2a8e6e696186ae3862/src",
    type = "zip",
    urls = [
        "https://codeload.github.com/protocolbuffers/protobuf/zip/9650e9fe8f737efcad485c2a8e6e696186ae3862",
    ],
)
proto_repository.archive(
    name = "googleapis",
    build_directives = [
        "gazelle:exclude google/example",
        "gazelle:exclude google/ads/googleads/v7/services",
        "gazelle:exclude google/ads/googleads/v8/services",
        "gazelle:exclude google/cloud/recommendationengine/v1beta1",
        "gazelle:proto_language go enabled true",
    ],
    build_file_generation = "clean",
    build_file_proto_mode = "file",
    cfgs = ["//:config.yaml"],
    imports = ["@protoapis//:imports.csv"],
    reresolve_known_proto_imports = True,
    sha256 = "95da12951c7d570980d5152f6cca9e1cb795ddc6b6dd7e9423bdffde28290f7a",
    strip_prefix = "googleapis-02710fa0ea5312d79d7fb986c9c9823fb41049a9",
    type = "zip",
    urls = [
        "https://codeload.github.com/googleapis/googleapis/zip/02710fa0ea5312d79d7fb986c9c9823fb41049a9",
    ],
)
use_repo(
    proto_repository,
    "googleapis",
    "protoapis",
)

~~~

