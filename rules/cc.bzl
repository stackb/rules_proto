# rules for cc
load("@build_stack_rules_proto//rules:cc_proto_compile.bzl", _cc_proto_compile = "cc_proto_compile")
load("@build_stack_rules_proto//rules:cc_proto_library.bzl", _cc_proto_library = "cc_proto_library")

cc_proto_compile = _cc_proto_compile
cc_proto_library = _cc_proto_library
