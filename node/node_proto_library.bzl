load("//node:node_proto_compile.bzl", "node_proto_compile")
load("//node:node_module_index.bzl", "node_module_index")
load("@org_pubref_rules_node//node:rules.bzl", "node_module")

def node_proto_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    name_index = kwargs.get("name") + "_index"
    node_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create index
    node_module_index(
        name = name_index,
        compilation = name_pb,
    )

    # Create node library
    node_module(
        name = kwargs.get("name"),
        srcs = [name_pb],
        index = name_index,
        deps = [
            "@proto_node_modules//:_all_",
        ],
        visibility = kwargs.get("visibility"),
    )
