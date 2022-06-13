---
layout: default
title: gomerge
permalink: examples/gomerge
parent: Examples
---


# gomerge example

`bazel test //example/golden:gomerge_test`


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


## `WORKSPACE`

~~~python
~~~

