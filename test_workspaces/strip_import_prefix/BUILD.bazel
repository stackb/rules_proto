package(default_visibility = ["//visibility:private"])

load("@build_stack_rules_proto//python:defs.bzl", "python_proto_library")

# Use the strip_import_prefix attribute

proto_library(
    name = "proto_lib",
    srcs = ["a/b/demo.proto"],
    strip_import_prefix = "a/b",
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
