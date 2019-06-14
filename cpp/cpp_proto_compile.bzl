load("//:compile.bzl", "proto_compile")

def cpp_proto_compile(**kwargs):
    # Append the cpp plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//cpp:cpp"),
    ]
    proto_compile(**kwargs)
