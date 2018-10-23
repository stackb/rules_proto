load("//:compile.bzl", "proto_compile")

def csharp_proto_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//csharp:csharp")),
        ],
        **kwargs
    )
