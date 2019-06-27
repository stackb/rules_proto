load("//github.com/gogo/protobuf:gogo_proto_compile.bzl", "gogo_proto_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

def gogo_proto_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    gogo_proto_compile(
        name = name_pb,
        prefix_path = kwargs.get("importpath", ""),
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create gogo library
    go_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = kwargs.get("go_deps", []) + [
            "@com_github_gogo_protobuf//proto:go_default_library",
        ],
        importpath = kwargs.get("importpath"),
        visibility = kwargs.get("visibility"),
    )
