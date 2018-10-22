load("//:compile.bzl", "proto_compile")

def cpp_grpc_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//cpp:cpp")),
            str(Label("//cpp:grpc_cpp")),
        ],
        **kwargs
    )
