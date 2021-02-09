load(
    "@build_stack_rules_proto//rules:proto_aspect.bzl",
    "proto_compile_aspect",
    "proto_compile_rule_macro",
    "proto_compile_rule",
)

_default_plugins = [
    str(Label("//python:python_plugin")),
]

_py_proto_compile_aspect = proto_compile_aspect(_default_plugins, "py_proto_compile_aspect")

_py_proto_compile_rule = proto_compile_rule(_py_proto_compile_aspect)

def py_proto_compile(**kwargs):
    proto_compile_rule_macro(_py_proto_compile_rule, **kwargs)
