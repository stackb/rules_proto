load("//:compile.bzl", "proto_compile")

def ts_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/improbable-eng/ts-protoc-gen:ts"))],
        **kwargs
    )

def ts_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/improbable-eng/ts-protoc-gen:ts_grpc"))],
        **kwargs
    )
