package(default_visibility = ["//visibility:private"])

load("@build_stack_rules_proto//python:python_proto_library.bzl", "python_proto_library")

# Set the transitive attribute to false and check that the transitive protos are not built (as expected)

proto_library(
    name = "base_proto_lib",
    srcs = ["base.proto"],
)

proto_library(
    name = "layer_proto_lib",
    srcs = ["layer.proto"],
    deps = ["base_proto_lib"],
)

python_proto_library(
    name = "py_lib",
    transitive = False,
    deps = ["layer_proto_lib"],
)

py_test(
    name = "test",
    srcs = ["main.py"],
    main = "main.py",
    deps = ["py_lib"],
)
