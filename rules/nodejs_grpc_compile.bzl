load(
    "@build_stack_rules_proto//rules:proto_aspect.bzl",
    "proto_compile_aspect",
    "proto_compile_rule",
    "proto_compile_rule_macro",
)

_default_plugins = [
    str(Label("//plugins/nodejs/proto:proto")),
    str(Label("//plugins/nodejs/grpc:grpc")),
]

_nodejs_grpc_compile_aspect = proto_compile_aspect(_default_plugins, "nodejs_grpc_compile_aspect")

_nodejs_grpc_compile_rule = proto_compile_rule(_nodejs_grpc_compile_aspect)

def nodejs_grpc_compile(**kwargs):
    proto_compile_rule_macro(_nodejs_grpc_compile_rule, **kwargs)
