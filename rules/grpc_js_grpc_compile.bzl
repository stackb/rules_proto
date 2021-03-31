load(
    "@build_stack_rules_proto//rules:proto_rules.bzl",
    "proto_compile_aspect",
    "proto_compile_aspect_rule_macro",
    "proto_compile_aspect_rule",
)

_default_plugins = [
    str(Label("//plugins/closure/proto:proto")),
    str(Label("//plugins/grpc_js/grpc:grpc")),
]

_grpc_js_grpc_compile_aspect = proto_compile_aspect(_default_plugins, "grpc_js_grpc_compile_aspect")

_grpc_js_grpc_compile_rule = proto_compile_aspect_rule(_grpc_js_grpc_compile_aspect)

def grpc_js_grpc_compile(**kwargs):
    proto_compile_aspect_rule_macro(_grpc_js_grpc_compile_rule, **kwargs)
