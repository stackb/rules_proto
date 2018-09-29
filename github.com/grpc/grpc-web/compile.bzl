load("//:compile.bzl", "proto_compile")

def closure_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/grpc/grpc-web:closure"))],
        **kwargs
    )

def commonjs_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/grpc/grpc-web:commonjs"))],
        **kwargs
    )

def commonjs_dts_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/grpc/grpc-web:commonjs_dts"))],
        **kwargs
    )

def ts_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/grpc/grpc-web:ts"))],
        **kwargs
    )
