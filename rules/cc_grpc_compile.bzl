load(
    "@build_stack_rules_proto//rules:proto_rules.bzl",
    "proto_compile_aspect",
    "proto_compile_aspect_rule_macro",
    "proto_compile_aspect_rule",
    "proto_compile_rule",
)

_default_plugins = [
    str(Label("//plugins/cc/proto:proto")),
    str(Label("//plugins/cc/grpc:grpc")),
]

_cc_grpc_compile_aspect = proto_compile_aspect(_default_plugins, "cc_grpc_compile_aspect")

_cc_grpc_compile_aspect_rule = proto_compile_aspect_rule(_cc_grpc_compile_aspect)

_cc_grpc_compile_rule = proto_compile_rule(_default_plugins)

def cc_grpc_compile(**kwargs):
    if kwargs.pop("transitive", False):
        proto_compile_aspect_rule_macro(_cc_grpc_compile_aspect_rule, **kwargs)
    else:
        _cc_grpc_compile_rule(**kwargs)
