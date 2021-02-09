load(
    "@build_stack_rules_proto//rules:proto_aspect.bzl",
    "proto_compile_aspect",
    "proto_compile_rule_macro",
    "proto_compile_rule",
)

_default_plugins = [
    str(Label("//python:python_plugin")),
    str(Label("//python:grpc_python_plugin")),
]

_py_grpc_compile_aspect = proto_compile_aspect(_default_plugins, "py_grpc_compile_aspect")

_py_grpc_compile_rule = proto_compile_rule(_py_grpc_compile_aspect)

def py_grpc_compile(**kwargs):
    proto_compile_rule_macro(_py_grpc_compile_rule, **kwargs)
