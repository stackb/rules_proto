
load("//go:utils.bzl", "get_importmappings")
load("//github.com/gogo/protobuf:compile.bzl", "gogo_proto_compile", "gogo_grpc_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")


def gogo_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    gogo_proto_compile(
        name = name_pb,
        deps = deps,
        plugin_options = get_importmappings(kwargs.pop("importmappings", {})),
        visibility = visibility,
    )

    go_deps = kwargs.pop("go_deps", [])

    kwargs["deps"] = depset([
        "@com_github_gogo_protobuf//proto:go_default_library",
    ] + go_deps).to_list()

    print("deps: %r", kwargs["deps"])
    go_library(
        srcs = [name_pb],
        **kwargs
    )


def gogo_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    gogo_grpc_compile(
        name = name_pb,
        deps = deps,
        plugin_options = get_importmappings(kwargs.pop("importmappings")),
        visibility = visibility,
    )

    kwargs["deps"] = depset([
        "@com_github_gogo_protobuf//proto:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_x_net//context:go_default_library",
    ] + kwargs.get("go_deps", [])).to_list()

    go_library(
        srcs = [name_pb],
        **kwargs
    )
