package(default_visibility = ["//visibility:private"])

load("@build_stack_rules_proto//python:defs.bzl", "python_proto_library")

# Test that outputs are correctly predicted for python plugin when using proto
# files in paths with dashes. These are converted to underscores by the python
# protoc plugin as dashes cannot be used in an import name.
# See https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/compiler/python/python_generator.cc#L89-L95

proto_library(
    name = "proto_lib",
    srcs = ["path-with-dashes/demo.proto"],
)

python_proto_library(
    name = "py_lib",
    deps = ["proto_lib"],
)

py_test(
    name = "test",
    srcs = ["main.py"],
    main = "main.py",
    deps = ["py_lib"],
)
