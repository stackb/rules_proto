package(default_visibility = ["//visibility:private"])

load("@build_stack_rules_proto//go:defs.bzl", "go_proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_test")

# Test that importpath on go_proto_library works

proto_library(
    name = "proto_lib_a",
    srcs = ["demo.proto"],
)

proto_library(
    name = "proto_lib_b",
    srcs = ["folder/nested_demo.proto"],
)

go_proto_library(
    name = "go_lib_a",
    deps = ["proto_lib_a"],
    importpath = "test_import_path",
)

go_proto_library(
    name = "go_lib_b",
    deps = ["proto_lib_b"],
    importpath = "test_import_path/folder",
)

go_test(
    name = "test",
    srcs = ["main.go"],
    deps = ["go_lib_a", "go_lib_b"],
    visibility = ["//visibility:public"],
)
