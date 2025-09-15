

# Preconfigured `Plugin` implementations

The plugin name is an opaque string, but by convention they are maven-esqe
artifact identifiers that follow a GitHub org/repo/plugin_name convention.

| Plugin                                                                                                                  |
|-------------------------------------------------------------------------------------------------------------------------|
| [builtin:cpp](/pkg/plugin/builtin/cpp_plugin.go)                                                                        |
| [builtin:csharp](/pkg/plugin/builtin/csharp_plugin.go)                                                                  |
| [builtin:java](/pkg/plugin/builtin/java_plugin.go)                                                                      |
| [builtin:js:closure](/pkg/plugin/builtin/js_closure_plugin.go)                                                          |
| [builtin:js:common](/pkg/plugin/builtin/js_common_plugin.go)                                                            |
| [builtin:objc](/pkg/plugin/builtin/objc_plugin.go)                                                                      |
| [builtin:php](/pkg/plugin/builtin/php_plugin.go)                                                                        |
| [builtin:python](/pkg/plugin/builtin/python_plugin.go)                                                                  |
| [builtin:pyi](/pkg/plugin/builtin/pyi_plugin.go)                                                                        |
| [builtin:ruby](/pkg/plugin/builtin/ruby_plugin.go)                                                                      |
| [grpc:grpc:cpp](/pkg/plugin/builtin/grpc_grpc_cpp.go)                                                                   |
| [grpc:grpc:protoc-gen-grpc-python](/pkg/plugin/grpc/grpc/protoc-gen-grpc-python.go)                                     |
| [golang:protobuf:protoc-gen-go](/pkg/plugin/golang/protobuf/protoc-gen-go.go)                                           |
| [grpc:grpc-go:protoc-gen-go-grpc](/pkg/plugin/grpc/grpcgo/protoc-gen-go-grpc.go)                                        |
| [grpc:grpc-java:protoc-gen-grpc-java](/pkg/plugin/grpc/grpcjava/protoc-gen-grpc-java.go)                                |
| [grpc:grpc-node:protoc-gen-grpc-node](/pkg/plugin/grpc/grpcnode/protoc-gen-grpc-node.go)                                |
| [grpc:grpc-web:protoc-gen-grpc-web](/pkg/plugin/grpc/grpcweb/protoc-gen-grpc-web.go)                                    |
| [gogo:protobuf:protoc-gen-combo](/pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                          |
| [gogo:protobuf:protoc-gen-gogo](/pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                           |
| [gogo:protobuf:protoc-gen-gogofast](/pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                       |
| [gogo:protobuf:protoc-gen-gogofaster](/pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                     |
| [gogo:protobuf:protoc-gen-gogoslick](/pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                      |
| [gogo:protobuf:protoc-gen-gogotypes](/pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                      |
| [gogo:protobuf:protoc-gen-gostring](/pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                       |
| [grpc-ecosystem:grpc-gateway:protoc-gen-grpc-gateway](/pkg/plugin/grpcecosystem/grpcgateway/protoc-gen-grpc-gateway.go) |
| [scalapb:scalapb:protoc-gen-scala](/pkg/plugin/scalapb/scalapb/protoc_gen_scala.go)                                     |
| [stackb:grpc.js:protoc-gen-grpc-js](/pkg/plugin/stackb/grpc_js/protoc-gen-grpc-js.go)                                   |
| [stephenh:ts-proto:protoc-gen-ts-proto](/pkg/plugin/stephenh/ts-proto/protoc-gen-ts-proto.go)                           |

## Preconfigured `proto_plugin` rules

```sh
$ bazel query 'kind(proto_plugin,@build_stack_rules_proto//...)' --keep_going
```

```py
@build_stack_rules_proto//plugin/akka/akka-grpc:protoc-gen-akka-grpc
@build_stack_rules_proto//plugin/bufbuild:connect-es
@build_stack_rules_proto//plugin/bufbuild:es
@build_stack_rules_proto//plugin/builtin:closurejs
@build_stack_rules_proto//plugin/builtin:commonjs
@build_stack_rules_proto//plugin/builtin:cpp
@build_stack_rules_proto//plugin/builtin:csharp
@build_stack_rules_proto//plugin/builtin:java
@build_stack_rules_proto//plugin/builtin:objc
@build_stack_rules_proto//plugin/builtin:php
@build_stack_rules_proto//plugin/builtin:pyi
@build_stack_rules_proto//plugin/builtin:python
@build_stack_rules_proto//plugin/builtin:ruby
@build_stack_rules_proto//plugin/gogo/protobuf:protoc-gen-combo
@build_stack_rules_proto//plugin/gogo/protobuf:protoc-gen-gogo
@build_stack_rules_proto//plugin/gogo/protobuf:protoc-gen-gogofast
@build_stack_rules_proto//plugin/gogo/protobuf:protoc-gen-gogofaster
@build_stack_rules_proto//plugin/gogo/protobuf:protoc-gen-gogoslick
@build_stack_rules_proto//plugin/gogo/protobuf:protoc-gen-gogotypes
@build_stack_rules_proto//plugin/gogo/protobuf:protoc-gen-gostring
@build_stack_rules_proto//plugin/golang/protobuf:protoc-gen-go
@build_stack_rules_proto//plugin/grpc/grpc:protoc-gen-grpc-cpp
@build_stack_rules_proto//plugin/grpc/grpc:protoc-gen-grpc-python
@build_stack_rules_proto//plugin/grpc/grpc-go:protoc-gen-go-grpc
@build_stack_rules_proto//plugin/grpc/grpc-java:protoc-gen-grpc-java
@build_stack_rules_proto//plugin/grpc/grpc-node:protoc-gen-grpc-node
@build_stack_rules_proto//plugin/grpc/grpc-web:protoc-gen-grpc-web
@build_stack_rules_proto//plugin/scalapb/scalapb:protoc-gen-scala
@build_stack_rules_proto//plugin/scalapb/scalapb:protoc-gen-scala-grpc
@build_stack_rules_proto//plugin/scalapb/zio-grpc:protoc-gen-zio-grpc
@build_stack_rules_proto//plugin/stackb/grpc_js:protoc-gen-grpc-js
@build_stack_rules_proto//plugin/stephenh/ts-proto:protoc-gen-ts-proto
```
