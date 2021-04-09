load(
    "@build_stack_rules_proto//rules:proto_rules.bzl",
    "proto_compile_aspect",
    "proto_compile_aspect_rule_macro",
    "proto_compile_aspect_rule",
    "proto_compile_rule",
)

_default_plugins = [
    str(Label("//plugins/gogo:gogoslick")),
]

_gogoslick_proto_compile_aspect = proto_compile_aspect(_default_plugins, "gogoslick_proto_compile_aspect")

_gogoslick_proto_compile_aspect_rule = proto_compile_aspect_rule(_gogoslick_proto_compile_aspect)

_gogoslick_proto_compile_rule = proto_compile_rule(_default_plugins)

def gogoslick_proto_compile(**kwargs):
    if kwargs.pop("transitive", False):
        proto_compile_aspect_rule_macro(_gogoslick_proto_compile_aspect_rule, **kwargs)
    else:
        _gogoslick_proto_compile_rule(**kwargs)
