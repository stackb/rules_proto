# rules for grpc_js
load("@build_stack_rules_proto//rules:grpc_js_grpc_compile.bzl", _grpc_js_grpc_compile = "grpc_js_grpc_compile")
load("@build_stack_rules_proto//rules:grpc_js_grpc_library.bzl", _grpc_js_grpc_library = "grpc_js_grpc_library")

grpc_js_grpc_compile = _grpc_js_grpc_compile
grpc_js_grpc_library = _grpc_js_grpc_library
