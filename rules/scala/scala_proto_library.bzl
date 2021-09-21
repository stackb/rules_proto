"scala_proto_library.bzl is a thin wrapper for @io_bazel_rules_scala//scala:scala_proto.bzl%scala_proto_library."

load("@io_bazel_rules_scala//scala_proto:scala_proto.bzl", _scala_proto_library = "scala_proto_library")

def scala_proto_library(**kwargs):
    _scala_proto_library(**kwargs)
