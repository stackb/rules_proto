load("//:compile.bzl", "proto_compile")

def swift_grpc_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//swift:grpc_swift")),
        ],
        **kwargs
    )
