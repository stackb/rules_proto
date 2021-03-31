load(
    "@build_stack_rules_proto//rules:proto_rules.bzl",
    "proto_compile_aspect",
    "proto_compile_aspect_rule_macro",
    "proto_compile_aspect_rule",
)

_default_plugins = [
    str(Label("//plugins/closure/proto:proto")),
]

_closure_proto_compile_aspect = proto_compile_aspect(_default_plugins, "closure_proto_compile_aspect")

_closure_proto_compile_aspect_rule = proto_compile_aspect_rule(_closure_proto_compile_aspect)

def closure_proto_compile(**kwargs):
    proto_compile_aspect_rule_macro(_closure_proto_compile_aspect_rule, **kwargs)
