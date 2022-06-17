---
layout: default
title: grpc_gateway
permalink: examples/grpc_gateway
parent: Examples
---


# grpc_gateway example

`bazel test //example/golden:grpc_gateway_test`


## `BUILD.bazel` (after gazelle)

~~~python
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@build_stack_rules_proto//rules/go:proto_go_library.bzl", "proto_go_library")
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")

# gazelle:proto_plugin protoc-gen-go implementation golang:protobuf:protoc-gen-go
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile
# gazelle:proto_rule proto_go_library implementation stackb:rules_proto:proto_go_library
# gazelle:proto_plugin protoc-gen-go dep @org_golang_google_protobuf//reflect/protoreflect
# gazelle:proto_plugin protoc-gen-go dep @org_golang_google_protobuf//runtime/protoimpl
# gazelle:proto_rule proto_go_library resolve google/protobuf/([a-z]+).proto @org_golang_google_protobuf//types/known/${1}pb
# gazelle:proto_rule proto_go_library resolve google/([a-z]+)/([a-z]+).proto @org_golang_google_genproto//googleapis/${1}/${2}
# gazelle:proto_rule proto_go_library visibility //visibility:public
# gazelle:proto_language go plugin protoc-gen-go
# gazelle:proto_language go rule proto_compile
# gazelle:proto_language go rule proto_go_library

# gazelle:proto_plugin protoc-gen-go-grpc implementation grpc:grpc-go:protoc-gen-go-grpc
# gazelle:proto_plugin protoc-gen-go-grpc dep @org_golang_google_grpc//:go_default_library
# gazelle:proto_plugin protoc-gen-go-grpc dep @org_golang_google_grpc//codes
# gazelle:proto_plugin protoc-gen-go-grpc dep @org_golang_google_grpc//status
# gazelle:proto_language go plugin protoc-gen-go-grpc

# gazelle:proto_plugin protoc-gen-grpc-gateway implementation grpc-ecosystem:grpc-gateway:protoc-gen-grpc-gateway
# gazelle:proto_plugin protoc-gen-grpc-gateway option logtostderr=true
# gazelle:proto_plugin protoc-gen-grpc-gateway option generate_unbound_methods=true
# gazelle:proto_plugin protoc-gen-grpc-gateway dep @org_golang_google_grpc//grpclog
# gazelle:proto_plugin protoc-gen-grpc-gateway dep @org_golang_google_grpc//metadata
# gazelle:proto_plugin protoc-gen-grpc-gateway dep @org_golang_google_protobuf//proto
# gazelle:proto_plugin protoc-gen-grpc-gateway dep @com_github_grpc_ecosystem_grpc_gateway_v2//runtime
# gazelle:proto_plugin protoc-gen-grpc-gateway dep @com_github_grpc_ecosystem_grpc_gateway_v2//utilities
# gazelle:proto_plugin protoc-gen-grpc-gateway dep @com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-openapiv2/options
# gazelle:proto_language go plugin protoc-gen-grpc-gateway

proto_library(
    name = "pb_proto",
    srcs = ["example.proto"],
    visibility = ["//visibility:public"],
    deps = ["@go_googleapis//google/api:annotations_proto"],
)

proto_compile(
    name = "pb_go_compile",
    options = {"@build_stack_rules_proto//plugin/grpc-ecosystem/grpc-gateway:protoc-gen-grpc-gateway": [
        "generate_unbound_methods=true",
        "logtostderr=true",
    ]},
    outputs = [
        "example.pb.go",
        "example.pb.gw.go",
        "example_grpc.pb.go",
    ],
    plugins = [
        "@build_stack_rules_proto//plugin/golang/protobuf:protoc-gen-go",
        "@build_stack_rules_proto//plugin/grpc-ecosystem/grpc-gateway:protoc-gen-grpc-gateway",
        "@build_stack_rules_proto//plugin/grpc/grpc-go:protoc-gen-go-grpc",
    ],
    proto = "pb_proto",
)

proto_go_library(
    name = "pb_go_proto",
    srcs = [
        "example.pb.go",
        "example.pb.gw.go",
        "example_grpc.pb.go",
    ],
    importpath = "./",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-openapiv2/options",
        "@com_github_grpc_ecosystem_grpc_gateway_v2//runtime",
        "@com_github_grpc_ecosystem_grpc_gateway_v2//utilities",
        "@org_golang_google_genproto//googleapis/api/annotations",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//grpclog",
        "@org_golang_google_grpc//metadata",
        "@org_golang_google_grpc//status",
        "@org_golang_google_protobuf//proto",
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
# gazelle:proto_plugin protoc-gen-go dep @org_golang_google_protobuf//reflect/protoreflect
# gazelle:proto_plugin protoc-gen-go dep @org_golang_google_protobuf//runtime/protoimpl
# gazelle:proto_rule proto_go_library resolve google/protobuf/([a-z]+).proto @org_golang_google_protobuf//types/known/${1}pb
# gazelle:proto_rule proto_go_library resolve google/([a-z]+)/([a-z]+).proto @org_golang_google_genproto//googleapis/${1}/${2}
# gazelle:proto_rule proto_go_library visibility //visibility:public
# gazelle:proto_language go plugin protoc-gen-go
# gazelle:proto_language go rule proto_compile
# gazelle:proto_language go rule proto_go_library

# gazelle:proto_plugin protoc-gen-go-grpc implementation grpc:grpc-go:protoc-gen-go-grpc
# gazelle:proto_plugin protoc-gen-go-grpc dep @org_golang_google_grpc//:go_default_library
# gazelle:proto_plugin protoc-gen-go-grpc dep @org_golang_google_grpc//codes
# gazelle:proto_plugin protoc-gen-go-grpc dep @org_golang_google_grpc//status
# gazelle:proto_language go plugin protoc-gen-go-grpc

# gazelle:proto_plugin protoc-gen-grpc-gateway implementation grpc-ecosystem:grpc-gateway:protoc-gen-grpc-gateway
# gazelle:proto_plugin protoc-gen-grpc-gateway option logtostderr=true
# gazelle:proto_plugin protoc-gen-grpc-gateway option generate_unbound_methods=true
# gazelle:proto_plugin protoc-gen-grpc-gateway dep @org_golang_google_grpc//grpclog
# gazelle:proto_plugin protoc-gen-grpc-gateway dep @org_golang_google_grpc//metadata
# gazelle:proto_plugin protoc-gen-grpc-gateway dep @org_golang_google_protobuf//proto
# gazelle:proto_plugin protoc-gen-grpc-gateway dep @com_github_grpc_ecosystem_grpc_gateway_v2//runtime
# gazelle:proto_plugin protoc-gen-grpc-gateway dep @com_github_grpc_ecosystem_grpc_gateway_v2//utilities
# gazelle:proto_plugin protoc-gen-grpc-gateway dep @com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-openapiv2/options
# gazelle:proto_language go plugin protoc-gen-grpc-gateway
~~~


## `WORKSPACE`

~~~python
~~~

