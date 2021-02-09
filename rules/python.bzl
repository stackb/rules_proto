# rules for python
load("@build_stack_rules_proto//rules:py_grpc_compile.bzl", _py_grpc_compile = "py_grpc_compile")
load("@build_stack_rules_proto//rules:py_proto_compile.bzl", _py_proto_compile = "py_proto_compile")
load("@build_stack_rules_proto//rules:py_proto_library.bzl", _py_proto_library = "py_proto_library")

py_grpc_compile = _py_grpc_compile
py_proto_compile = _py_proto_compile
py_proto_library = _py_proto_library
