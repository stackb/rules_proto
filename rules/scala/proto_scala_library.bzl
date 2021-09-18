"proto_scala_library.bzl provides a scala_library for proto files."

load("@io_bazel_rules_scala//scala:scala.bzl", "scala_library")
# load("@io_bazel_rules_scala//scala_proto:default_dep_sets.bzl", "DEFAULT_SCALAPB_COMPILE_DEPS")

def proto_scala_library(**kwargs):
    scala_library(**kwargs)
