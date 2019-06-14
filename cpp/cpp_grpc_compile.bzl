load("//:compile.bzl", "proto_compile")

def cpp_grpc_compile(**kwargs):
    # Append the cpp plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//cpp:cpp"),
        Label("//cpp:grpc_cpp"),
    ]
    proto_compile(**kwargs)
