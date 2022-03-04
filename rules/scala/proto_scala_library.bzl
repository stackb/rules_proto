"proto_scala_library.bzl provides a scala_library for proto files."

load("@io_bazel_rules_scala//scala:scala.bzl", "scala_library")

def proto_scala_library(**kwargs):
    deps = kwargs.get("deps", [])
    kwargs.setdefault("exports", deps)
    scala_library(deps = deps, exports = exports, **kwargs)
