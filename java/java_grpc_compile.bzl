load("//:compile.bzl", "proto_compile")

def java_grpc_compile(**kwargs):
    # Append the java plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//java:java"),
        Label("//java:grpc_java"),
    ]
    proto_compile(**kwargs)
