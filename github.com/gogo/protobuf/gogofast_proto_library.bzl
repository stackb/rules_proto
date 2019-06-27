load("//github.com/gogo/protobuf:gogofast_proto_compile.bzl", "gogofast_proto_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

def gogofast_proto_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    gogofast_proto_compile(
        name = name_pb,
        deps = deps, # Forward only deps
        prefix_path = kwargs.get("importpath", ""),
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
