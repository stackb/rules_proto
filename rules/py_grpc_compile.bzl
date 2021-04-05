load(
    "@build_stack_rules_proto//rules:proto_rules.bzl",
    "proto_compile_aspect",
    "proto_compile_aspect_rule_macro",
    "proto_compile_aspect_rule",
    "proto_compile_rule",
)

_default_plugins = [
    str(Label("//plugins/python/grpc:grpc")),
]

_py_grpc_compile_aspect = proto_compile_aspect(_default_plugins, "py_grpc_compile_aspect")

_py_grpc_compile_aspect_rule = proto_compile_aspect_rule(_py_grpc_compile_aspect)

_py_grpc_compile_rule = proto_compile_rule(_default_plugins)

def py_grpc_compile(**kwargs):
    if kwargs.pop("transitive", False):
        proto_compile_aspect_rule_macro(_py_grpc_compile_aspect_rule, **kwargs)
    else:
        _py_grpc_compile_rule(**kwargs)
