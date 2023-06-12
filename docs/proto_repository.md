---
layout: default
title: proto_repository
permalink: examples/proto_repository
parent: Examples
---


# proto_repository example

`bazel test //example/golden:proto_repository_test`


## `BUILD.bazel` (after gazelle)

~~~python
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@build_stack_rules_proto//rules/go:proto_go_library.bzl", "proto_go_library")
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")

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


## `BUILD.bazel` (before gazelle)

~~~python
~~~


## `WORKSPACE`

~~~python
# ----------------------------------------------------
# proto_repository
# ----------------------------------------------------

load("@build_stack_rules_proto//rules/proto:proto_repository.bzl", "proto_repository")

proto_repository(
    name = "googleapis",
    build_directives = [
        "gazelle:proto_language go enabled true",
    ],
    build_file_generation = "on",
    build_file_proto_mode = "file",
    cfgs = ["//:config.yaml"],
    reresolve_known_proto_imports = True,
    sha256 = "b9dbc65ebc738a486265ef7b708e9449bf361541890091983e946557ee0a4bfc",
    strip_prefix = "googleapis-66759bdf6a5ebb898c2a51c8649aefd1ee0b7ffe",
    type = "zip",
    urls = ["https://codeload.github.com/googleapis/googleapis/zip/66759bdf6a5ebb898c2a51c8649aefd1ee0b7ffe"],
)
~~~

