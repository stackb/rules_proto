---
layout: default
title: go_grpc
permalink: examples/go_grpc
parent: Examples
---


# go_grpc example

`bazel test //example/golden:go_grpc_test`


## `BUILD.bazel` (after gazelle)

~~~python
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@build_stack_rules_proto//rules/go:proto_go_library.bzl", "proto_go_library")
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")

# gazelle:proto_plugin protoc-gen-go implementation golang:protobuf:protoc-gen-go
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile
# gazelle:proto_rule proto_go_library implementation stackb:rules_proto:proto_go_library
# gazelle:proto_rule proto_go_library deps @org_golang_google_protobuf//reflect/protoreflect
# gazelle:proto_rule proto_go_library deps @org_golang_google_protobuf//runtime/protoimpl
# gazelle:proto_rule proto_go_library visibility //visibility:public
# gazelle:proto_language go plugin protoc-gen-go
# gazelle:proto_language go rule proto_compile
# gazelle:proto_language go rule proto_go_library

# gazelle:proto_plugin protoc-gen-go-grpc implementation grpc:grpc-go:protoc-gen-go-grpc
# gazelle:proto_plugin protoc-gen-go-grpc dep @org_golang_google_grpc//:go_default_library
# gazelle:proto_plugin protoc-gen-go-grpc dep @org_golang_google_grpc//codes
# gazelle:proto_plugin protoc-gen-go-grpc dep @org_golang_google_grpc//status
# gazelle:proto_language go plugin protoc-gen-go-grpc

proto_library(
    name = "pb_proto",
    srcs = ["example.proto"],
    visibility = ["//visibility:public"],
)

proto_compile(
    name = "pb_go_compile",
    outputs = [
        "example.pb.go",
        "example_grpc.pb.go",
    ],
    plugins = [
        "@build_stack_rules_proto//plugin/golang/protobuf:protoc-gen-go",
        "@build_stack_rules_proto//plugin/grpc/grpc-go:protoc-gen-go-grpc",
    ],
    proto = "pb_proto",
)

proto_go_library(
    name = "pb_go_proto",
    srcs = [
        "example.pb.go",
        "example_grpc.pb.go",
    ],
    importpath = "./",
    visibility = ["//visibility:public"],
    deps = [
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
# gazelle:proto_plugin protoc-gen-go implementation golang:protobuf:protoc-gen-go
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile
# gazelle:proto_rule proto_go_library implementation stackb:rules_proto:proto_go_library
# gazelle:proto_rule proto_go_library deps @org_golang_google_protobuf//reflect/protoreflect
# gazelle:proto_rule proto_go_library deps @org_golang_google_protobuf//runtime/protoimpl
# gazelle:proto_rule proto_go_library visibility //visibility:public
# gazelle:proto_language go plugin protoc-gen-go
# gazelle:proto_language go rule proto_compile
# gazelle:proto_language go rule proto_go_library

# gazelle:proto_plugin protoc-gen-go-grpc implementation grpc:grpc-go:protoc-gen-go-grpc
# gazelle:proto_plugin protoc-gen-go-grpc dep @org_golang_google_grpc//:go_default_library
# gazelle:proto_plugin protoc-gen-go-grpc dep @org_golang_google_grpc//codes
# gazelle:proto_plugin protoc-gen-go-grpc dep @org_golang_google_grpc//status
# gazelle:proto_language go plugin protoc-gen-go-grpc
~~~


## `WORKSPACE`

~~~python
~~~

