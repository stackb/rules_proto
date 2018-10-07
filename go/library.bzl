load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", _go_proto_library = "go_proto_library")

load("//go:compile.bzl", "go_proto_compile", "go_grpc_compile")

def go_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    importpath = kwargs.get("importpath")
    visibility = kwargs.get("visibility")
    go_deps = kwargs.get("go_deps", [])

    name_pb = name + "_pb"

    go_proto_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
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


def go_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    importpath = kwargs.get("importpath")
    visibility = kwargs.get("visibility")
    go_deps = kwargs.get("go_deps", [])

    name_pb = name + "_pb"

    go_grpc_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )

    go_library(
        name = name,
        srcs = [name_pb],
        deps = go_deps + [
            "@com_github_golang_protobuf//proto:go_default_library",
            "@org_golang_google_grpc//:go_default_library",
            "@org_golang_x_net//context:go_default_library",
        ],
        importpath = importpath,
        visibility = visibility,
    )



def golang_proto_library(**kwargs):
    name = kwargs.get("name")
    importpath = kwargs.get("importpath")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")
    compilers = kwargs.get("compilers")
    if not compilers:
        compilers = ["@io_bazel_rules_go//proto:go_proto"]
    _go_proto_library(
        name = name,
        compilers = compilers,
        importpath = importpath,
        proto = deps[0],
        visibility = visibility,
    )

def golang_grpc_library(**kwargs):
    name = kwargs.get("name")
    importpath = kwargs.get("importpath")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")
    compilers = kwargs.get("compilers")
    if not compilers:
        compilers = ["@io_bazel_rules_go//proto:go_grpc"]
    _go_proto_library(
        name = name,
        compilers = compilers,
        importpath = importpath,
        proto = deps[0],
        visibility = visibility,
    )
