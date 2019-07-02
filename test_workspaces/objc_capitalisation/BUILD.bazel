package(default_visibility = ["//visibility:private"])

load("@build_stack_rules_proto//objc:defs.bzl", "objc_proto_compile")

# Test that Objective-C files are correctly capitalised

proto_library(
    name = "proto_lib",
    srcs = [
        "demo.proto",
        "folder/nested_demo.proto",
    ],
)

objc_proto_compile(
    name = "objc_lib",
    deps = ["proto_lib"],
)

py_test(
    name = "test",
    srcs = ["main.py"],
    main = "main.py",
    data = ["objc_lib"],
)
