## Plugin Implementations

The plugin name is an opaque string, but by convention they are maven-esqe
artifact identifiers that follow a GitHub org/repo/plugin_name convention.

| Plugin                                                                                                                 |
|------------------------------------------------------------------------------------------------------------------------|
| [builtin:cpp](pkg/plugin/builtin/cpp_plugin.go)                                                                        |
| [builtin:csharp](pkg/plugin/builtin/csharp_plugin.go)                                                                  |
| [builtin:java](pkg/plugin/builtin/java_plugin.go)                                                                      |
| [builtin:js:closure](pkg/plugin/builtin/js_closure_plugin.go)                                                          |
| [builtin:js:common](pkg/plugin/builtin/js_common_plugin.go)                                                            |
| [builtin:objc](pkg/plugin/builtin/objc_plugin.go)                                                                      |
| [builtin:php](pkg/plugin/builtin/php_plugin.go)                                                                        |
| [builtin:python](pkg/plugin/builtin/python_plugin.go)                                                                  |
| [builtin:pyi](pkg/plugin/builtin/pyi_plugin.go)                                                                        |
| [builtin:ruby](pkg/plugin/builtin/ruby_plugin.go)                                                                      |
| [grpc:grpc:cpp](pkg/plugin/builtin/grpc_grpc_cpp.go)                                                                   |
| [grpc:grpc:protoc-gen-grpc-python](pkg/plugin/grpc/grpc/protoc-gen-grpc-python.go)                                     |
| [golang:protobuf:protoc-gen-go](pkg/plugin/golang/protobuf/protoc-gen-go.go)                                           |
| [grpc:grpc-go:protoc-gen-go-grpc](pkg/plugin/grpc/grpcgo/protoc-gen-go-grpc.go)                                        |
| [grpc:grpc-java:protoc-gen-grpc-java](pkg/plugin/grpc/grpcjava/protoc-gen-grpc-java.go)                                |
| [grpc:grpc-node:protoc-gen-grpc-node](pkg/plugin/grpc/grpcnode/protoc-gen-grpc-node.go)                                |
| [grpc:grpc-web:protoc-gen-grpc-web](pkg/plugin/grpc/grpcweb/protoc-gen-grpc-web.go)                                    |
| [gogo:protobuf:protoc-gen-combo](pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                          |
| [gogo:protobuf:protoc-gen-gogo](pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                           |
| [gogo:protobuf:protoc-gen-gogofast](pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                       |
| [gogo:protobuf:protoc-gen-gogofaster](pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                     |
| [gogo:protobuf:protoc-gen-gogoslick](pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                      |
| [gogo:protobuf:protoc-gen-gogotypes](pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                      |
| [gogo:protobuf:protoc-gen-gostring](pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                       |
| [grpc-ecosystem:grpc-gateway:protoc-gen-grpc-gateway](pkg/plugin/grpcecosystem/grpcgateway/protoc-gen-grpc-gateway.go) |
| [scalapb:scalapb:protoc-gen-scala](pkg/plugin/scalapb/scalapb/protoc_gen_scala.go)                                     |
| [stackb:grpc.js:protoc-gen-grpc-js](pkg/plugin/stackb/grpc_js/protoc-gen-grpc-js.go)                                   |
| [stephenh:ts-proto:protoc-gen-ts-proto](pkg/plugin/stephenh/ts-proto/protoc-gen-ts-proto.go)                           |

## Rule Implementations

The rule name is an opaque string, but by convention they are maven-esqe
artifact identifiers that follow a GitHub org/repo/rule_name convention.

| Plugin                                                                                            |
|---------------------------------------------------------------------------------------------------|
| [stackb:rules_proto:grpc_cc_library](pkg/rule/rules_cc/grpc_cc_library.go)                        |
| [stackb:rules_proto:grpc_closure_js_library](pkg/rule/rules_closure/grpc_closure_js_library.go)   |
| [stackb:rules_proto:grpc_java_library](pkg/rule/rules_java/grpc_java_library.go)                  |
| [stackb:rules_proto:grpc_nodejs_library](pkg/rule/rules_nodejs/grpc_nodejs_library.go)            |
| [stackb:rules_proto:grpc_web_js_library](pkg/rule/rules_nodejs/grpc_web_js_library.go)            |
| [stackb:rules_proto:grpc_py_library](pkg/rule/rules_python/grpc_py_library.go)                    |
| [stackb:rules_proto:proto_cc_library](pkg/rule/rules_cc/proto_cc_library.go)                      |
| [stackb:rules_proto:proto_closure_js_library](pkg/rule/rules_closure/proto_closure_js_library.go) |
| [stackb:rules_proto:proto_compile](pkg/protoc/proto_compile.go)                                   |
| [stackb:rules_proto:proto_compiled_sources](pkg/protoc/proto_compiled_sources.go)                 |
| [stackb:rules_proto:proto_descriptor_set](pkg/protoc/proto_descriptor_set.go)                     |
| [stackb:rules_proto:proto_go_library](pkg/rule/rules_go/go_library.go)                            |
| [stackb:rules_proto:proto_java_library](pkg/rule/rules_java/proto_java_library.go)                |
| [stackb:rules_proto:proto_nodejs_library](pkg/rule/rules_nodejs/proto_nodejs_library.go)          |
| [stackb:rules_proto:proto_py_library](pkg/rule/rules_python/proto_py_library.go)                  |
| [bazelbuild:rules_scala:scala_proto_library](pkg/rule/rules_scala/scala_proto_library.go)         |

Please consult the `example/` directory and unit tests for more additional
detail.


# Writing Custom Plugins and Rules

Custom plugin implementations and rule implementations can be written in golang
or starlark.  Golang implementations are statically compiled into the final
`gazelle_binary` whereas starlark plugins are evaluated at gazelle runtime.

## +/- of golang implementations

- `+` Full power of a statically compiled language, the golang stdlib, and
  external dependencies.
- `+` Easier to test.
- `+` API not experimental.
- `-` Cannot be used in a `proto_repository` rule without forking
  stackb/rules_proto.
- `-` Initial setup harder, often housed within your own custom gazelle
  extension.

Until a dedicated tutorial is available, please consult the source code for
examples.

## +/- of starlark implementations

- `+` More familiar to developer with starlark experience but not golang.
- `+` Easier setup (*.star files in your gazelle repository)
- `+` Possible to use in conjunction with the `proto_repository` rule.
- `-` Limited API: can only reference state that has been already configured via gazelle directives.
- `-` Not possible to implement stateful design.
- `-` No standard library.

Until a dedicated tutorial is available, please consult the reference example in
`example/testdata/starlark_java`.
