
load("@io_bazel_rules_go//proto:def.bzl", _go_proto_library = "go_proto_library")

go_proto_library = _go_proto_library

def go_grpc_library(**kwargs):
    kwargs.setdefault("compilers", ["@io_bazel_rules_go//proto:go_grpc"])
    go_proto_library(**kwargs)