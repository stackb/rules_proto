load("@build_stack_rules_proto//scala:scala_grpc_compile.bzl", "scala_grpc_compile")
load("@io_bazel_rules_scala//scala:scala.bzl", "scala_library")

def scala_grpc_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    scala_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create scala library
    scala_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [str(Label("//scala:grpc_deps"))],
        exports = [
            str(Label("//scala:grpc_deps")),
        ],
        visibility = kwargs.get("visibility"),
    )
