load("//:compile.bzl", "proto_compile")

def cpp_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//cpp:cpp"))],
        **kwargs
    )

def grpc_cpp_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//cpp:cpp")), str(Label("//cpp:grpc_cpp"))],
        **kwargs
    )
