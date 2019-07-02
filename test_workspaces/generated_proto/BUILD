package(default_visibility = ["//visibility:private"])

load("@build_stack_rules_proto//python:defs.bzl", "python_proto_library")

# Test that generated proto files are handled correctly

genrule(
    name = "generate_proto",
    srcs = ["demo.proto.template"],
    outs = ["demo.proto"],
    cmd = "cp '$(location demo.proto.template)' '$@'",
)

proto_library(
    name = "proto_lib",
    srcs = ["demo.proto"],
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
