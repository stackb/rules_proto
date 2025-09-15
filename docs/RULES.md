
# Preconfigured `Rule` implementations

The rule name is an opaque string, but by convention they are maven-esqe
artifact identifiers that follow a GitHub org/repo/rule_name convention.

| Rule                                                                                               |
|----------------------------------------------------------------------------------------------------|
| [stackb:rules_proto:grpc_cc_library](/pkg/rule/rules_cc/grpc_cc_library.go)                        |
| [stackb:rules_proto:grpc_closure_js_library](/pkg/rule/rules_closure/grpc_closure_js_library.go)   |
| [stackb:rules_proto:grpc_java_library](/pkg/rule/rules_java/grpc_java_library.go)                  |
| [stackb:rules_proto:grpc_nodejs_library](/pkg/rule/rules_nodejs/grpc_nodejs_library.go)            |
| [stackb:rules_proto:grpc_web_js_library](/pkg/rule/rules_nodejs/grpc_web_js_library.go)            |
| [stackb:rules_proto:grpc_py_library](/pkg/rule/rules_python/grpc_py_library.go)                    |
| [stackb:rules_proto:proto_cc_library](/pkg/rule/rules_cc/proto_cc_library.go)                      |
| [stackb:rules_proto:proto_closure_js_library](/pkg/rule/rules_closure/proto_closure_js_library.go) |
| [stackb:rules_proto:proto_compile](/pkg/protoc/proto_compile.go)                                   |
| [stackb:rules_proto:proto_compiled_sources](/pkg/protoc/proto_compiled_sources.go)                 |
| [stackb:rules_proto:proto_descriptor_set](/pkg/protoc/proto_descriptor_set.go)                     |
| [stackb:rules_proto:proto_go_library](/pkg/rule/rules_go/go_library.go)                            |
| [stackb:rules_proto:proto_java_library](/pkg/rule/rules_java/proto_java_library.go)                |
| [stackb:rules_proto:proto_nodejs_library](/pkg/rule/rules_nodejs/proto_nodejs_library.go)          |
| [stackb:rules_proto:proto_py_library](/pkg/rule/rules_python/proto_py_library.go)                  |
| [bazelbuild:rules_scala:scala_proto_library](/pkg/rule/rules_scala/scala_proto_library.go)         |
