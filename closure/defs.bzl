# Aggregate all `closure` rules to one loadable file
load(":closure_proto_compile.bzl", _closure_proto_compile="closure_proto_compile")
load(":closure_proto_library.bzl", _closure_proto_library="closure_proto_library")

closure_proto_compile = _closure_proto_compile
closure_proto_library = _closure_proto_library
