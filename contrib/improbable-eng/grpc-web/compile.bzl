load("//:compile.bzl", "proto_compile")

def ts_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//ts:ts"))],
        **kwargs
    )

def ts_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//ts:grpc_ts"))],
        **kwargs
    )
