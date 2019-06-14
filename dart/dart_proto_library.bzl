load("//dart:dart_proto_compile.bzl", "dart_proto_compile")
load("@io_bazel_rules_dart//dart/build_rules:core.bzl", "dart_library")

def dart_proto_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    dart_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create dart library
    dart_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            str(Label("@vendor_protobuf//:protobuf")),
        ],
        pub_pkg_name = kwargs.get("name"),
        visibility = kwargs.get("visibility"),
    )
