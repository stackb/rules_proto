load("//go:go_proto_compile.bzl", "go_proto_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("//go:utils.bzl", "get_importmappings")

def go_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    plugins = kwargs.get("plugins", [])
    importpath = kwargs.get("importpath")
    visibility = kwargs.get("visibility")
    go_deps = kwargs.get("go_deps", [])

    name_pb = name + "_pb"

    go_proto_compile(
        name = name_pb,
        deps = deps,
        plugins = plugins,
        plugin_options = get_importmappings(kwargs.pop("importmap", {})),
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    go_library(
        name = name,
        srcs = [name_pb],
        deps = go_deps + [
            "@com_github_golang_protobuf//proto:go_default_library",
        ],
        importpath = importpath,
        visibility = visibility,
    )
