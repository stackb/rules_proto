package(default_visibility = ["//visibility:private"])

load("@build_stack_rules_proto//python:python_proto_library.bzl", "python_proto_library")

# Test that exclusions are correctly being applied

proto_library(
    name = "proto_lib",
    srcs = ["demo.proto"],
    deps = ["@com_google_protobuf//:timestamp_proto"],
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
