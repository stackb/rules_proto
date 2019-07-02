package(default_visibility = ["//visibility:private"])

load("@build_stack_rules_proto//python:defs.bzl", "python_proto_library")

# Use the proto_source_root attribute

proto_library(
    name = "proto_lib",
    srcs = ["demo.proto"],
    proto_source_root = "package",
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
