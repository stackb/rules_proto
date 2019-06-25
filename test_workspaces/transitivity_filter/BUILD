package(default_visibility = ["//visibility:private"])

load("@build_stack_rules_proto//python:python_proto_library.bzl", "python_proto_library")

# Use the transitivity attribute to exclude a transitive file

proto_library(
    name = "base_proto_lib",
    srcs = [
        "base.proto",
        "exclude_dir/exclude_base.proto",
    ],
)

proto_library(
    name = "layer_proto_lib",
    srcs = ["layer.proto"],
    deps = ["base_proto_lib"],
)

python_proto_library(
    name = "py_lib",
    transitive = True,
    transitivity = {
        "exclude_dir": "exclude",
    },
    deps = ["layer_proto_lib"],
)

py_test(
    name = "test",
    srcs = ["main.py"],
    main = "main.py",
    deps = ["py_lib"],
)
