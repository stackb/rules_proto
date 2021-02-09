
load("@build_stack_rules_proto//LANG:py_proto_compile.bzl", _py_proto_compile = "py_proto_compile")

py_proto_compile = _py_proto_compile
load("@build_stack_rules_proto//LANG:py_proto_library.bzl", _py_proto_library = "py_proto_library")

py_proto_library = _py_proto_library
