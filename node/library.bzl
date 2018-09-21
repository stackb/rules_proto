load("//node:compile.bzl", "node_proto_compile", "node_grpc_compile")
load("//node:node_module_index.bzl", "node_module_index")
load("@org_pubref_rules_node//node:rules.bzl", "node_module")

def node_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    name_index = name + "_index"

    node_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
    )
    node_module_index(
        name = name_index,
        compilation = name_pb,
    )
    node_module(
        name = name,
        srcs = [name_pb],
        index = name_index,
        deps = [
            "@proto_node_modules//:_all_",
        ],
        visibility = visibility,
    )


def node_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    name_index = name + "_index"

    node_grpc_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
    )
    node_module_index(
        name = name_index,
        compilation = name_pb,
    )
    node_module(
        name = name,
        srcs = [name_pb],
        index = name_index,
        deps = [
            "@grpc_node_modules//:_all_",
        ],
        visibility = visibility,
    )
