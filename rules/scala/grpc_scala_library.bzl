"grpc_scala_library.bzl provides a scala_library for grpc files."

load("@io_bazel_rules_scala//scala:scala.bzl", "scala_library")

def grpc_scala_library(**kwargs):
    deps = kwargs.get("deps", [])
    kwargs.setdefault("exports", deps)
    scala_library(**kwargs)
