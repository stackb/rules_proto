load("//:compile.bzl", "proto_compile")

def csharp_grpc_compile(**kwargs):
    # Append the csharp plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//csharp:csharp"),
        Label("//csharp:grpc_csharp"),
    ]
    proto_compile(**kwargs)
