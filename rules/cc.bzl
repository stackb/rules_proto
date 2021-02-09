# rules for cc
load("@build_stack_rules_proto//rules:cc_grpc_compile.bzl", _cc_grpc_compile = "cc_grpc_compile")
load("@build_stack_rules_proto//rules:cc_grpc_library.bzl", _cc_grpc_library = "cc_grpc_library")
load("@build_stack_rules_proto//rules:cc_proto_compile.bzl", _cc_proto_compile = "cc_proto_compile")
load("@build_stack_rules_proto//rules:cc_proto_library.bzl", _cc_proto_library = "cc_proto_library")

cc_grpc_compile = _cc_grpc_compile
cc_grpc_library = _cc_grpc_library
cc_proto_compile = _cc_proto_compile
cc_proto_library = _cc_proto_library
