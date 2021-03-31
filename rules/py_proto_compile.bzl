load(
    "@build_stack_rules_proto//rules:proto_rules.bzl",
    "proto_compile_aspect",
    "proto_compile_aspect_rule",
    "proto_compile_aspect_rule_macro",
    "proto_compile_rule",
)

_default_plugins = [
    str(Label("//plugins/python/proto:proto")),
]

_py_proto_compile_direct = proto_compile_rule(_default_plugins)
_py_proto_compile_aspect = proto_compile_aspect(_default_plugins, "py_proto_compile_aspect")

_py_proto_compile_aspect_rule = proto_compile_aspect_rule(_py_proto_compile_aspect)

def py_proto_compile_transitive(**kwargs):
    proto_compile_aspect_rule_macro(_py_proto_compile_aspect_rule, **kwargs)

def py_proto_compile_with_aspect(**kwargs):
    proto_compile_aspect_rule_macro(_py_proto_compile_aspect_rule, **kwargs)

def py_proto_compile(**kwargs):
    _py_proto_compile_direct(**kwargs)
