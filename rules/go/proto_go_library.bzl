"proto_go_library.bzl provides a go_library for proto files."

load("@io_bazel_rules_go//go:def.bzl", "go_library")

def proto_go_library(**kwargs):
    go_library(**kwargs)
