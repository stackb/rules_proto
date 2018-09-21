load("//:compile.bzl", "proto_compile")

def objc_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//objc:objc"))],
        **kwargs
    )

def objc_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//objc:objc")), str(Label("//objc:grpc_objc"))],
        **kwargs
    )
