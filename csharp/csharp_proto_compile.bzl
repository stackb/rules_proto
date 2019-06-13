load("//:compile.bzl", "proto_compile")

def csharp_proto_compile(**kwargs):
    # Prepend the csharp plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//csharp:csharp"),
    ]
    proto_compile(**kwargs)
