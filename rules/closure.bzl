# rules for closure
load("@build_stack_rules_proto//rules:closure_proto_compile.bzl", _closure_proto_compile = "closure_proto_compile")
load("@build_stack_rules_proto//rules:closure_proto_library.bzl", _closure_proto_library = "closure_proto_library")

closure_proto_compile = _closure_proto_compile
closure_proto_library = _closure_proto_library
