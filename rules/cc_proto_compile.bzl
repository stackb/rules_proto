load(
    "@build_stack_rules_proto//rules:proto_aspect.bzl",
    "proto_compile_aspect",
    "proto_compile_rule_macro",
    "proto_compile_rule",
)

_default_plugins = [
    str(Label("//cc:cc_plugin")),
]

_cc_proto_compile_aspect = proto_compile_aspect(_default_plugins, "cc_proto_compile_aspect")

_cc_proto_compile_rule = proto_compile_rule(_cc_proto_compile_aspect)

def cc_proto_compile(**kwargs):
    proto_compile_rule_macro(_cc_proto_compile_rule, **kwargs)
