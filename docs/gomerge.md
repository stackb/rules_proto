---
layout: default
title: gomerge
permalink: examples/gomerge
parent: Examples
---


# gomerge example

[`testdata files`](/example/golden/testdata/gomerge)


## `Integration Test`

`bazel test @@//example/golden:gomerge_test`)


## `BUILD.bazel` (before gazelle)

~~~python
# gazelle:proto file

# gazelle:proto_plugin protoc-gen-go implementation golang:protobuf:protoc-gen-go
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile
# gazelle:proto_rule proto_go_library implementation stackb:rules_proto:proto_go_library
# gazelle:proto_rule proto_go_library deps @org_golang_google_protobuf//reflect/protoreflect
# gazelle:proto_rule proto_go_library deps @org_golang_google_protobuf//runtime/protoimpl
# gazelle:proto_rule proto_go_library resolve google/protobuf/([a-z]+).proto @org_golang_google_protobuf//types/known/${1}pb
# gazelle:proto_rule proto_go_library visibility //visibility:public
# gazelle:proto_language go plugin protoc-gen-go
# gazelle:proto_language go rule proto_compile
# gazelle:proto_language go rule proto_go_library
~~~


## `BUILD.bazel` (after gazelle)

~~~python
# gazelle:proto file

# gazelle:proto_plugin protoc-gen-go implementation golang:protobuf:protoc-gen-go
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile
# gazelle:proto_rule proto_go_library implementation stackb:rules_proto:proto_go_library
# gazelle:proto_rule proto_go_library deps @org_golang_google_protobuf//reflect/protoreflect
# gazelle:proto_rule proto_go_library deps @org_golang_google_protobuf//runtime/protoimpl
# gazelle:proto_rule proto_go_library resolve google/protobuf/([a-z]+).proto @org_golang_google_protobuf//types/known/${1}pb
# gazelle:proto_rule proto_go_library visibility //visibility:public
# gazelle:proto_language go plugin protoc-gen-go
# gazelle:proto_language go rule proto_compile
# gazelle:proto_language go rule proto_go_library
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

~~~

