"grpc_zio_scala_library.bzl provides a scala_library for grpc files."

load("@io_bazel_rules_scala//scala:scala.bzl", "scala_library")

def grpc_zio_scala_library(**kwargs):
    scala_library(**kwargs)
