load("//:compile.bzl", "proto_compile")

def grpc_js_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/stackb/grpc.js:grpc.js"))],
        **kwargs
    )
