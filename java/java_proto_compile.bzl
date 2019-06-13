load("//:compile.bzl", "proto_compile")

def java_proto_compile(**kwargs):
    # Prepend the java plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//java:java"),
    ]
    proto_compile(**kwargs)
