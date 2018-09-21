load("//:compile.bzl", "proto_compile")

def csharp_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//csharp:csharp"))],
        **kwargs
    )

def csharp_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//csharp:csharp")), str(Label("//csharp:grpc_csharp"))],
        **kwargs
    )
