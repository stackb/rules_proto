load(
    "@build_stack_rules_proto//rules:proto_rules.bzl",
    "proto_compile_aspect",
    "proto_compile_aspect_rule_macro",
    "proto_compile_aspect_rule",
    "proto_compile_rule",
)

_default_plugins = [
    str(Label("//plugins/gogo:gogofast")),
]

_gogofast_proto_compile_aspect = proto_compile_aspect(_default_plugins, "gogofast_proto_compile_aspect")

_gogofast_proto_compile_aspect_rule = proto_compile_aspect_rule(_gogofast_proto_compile_aspect)

_gogofast_proto_compile_rule = proto_compile_rule(_default_plugins)

def gogofast_proto_compile(**kwargs):
    if kwargs.pop("transitive", False):
        proto_compile_aspect_rule_macro(_gogofast_proto_compile_aspect_rule, **kwargs)
    else:
        _gogofast_proto_compile_rule(**kwargs)
