load("//:compile.bzl", "proto_compile")

def android_grpc_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//android:java")),
            str(Label("//android:grpc_javalite")),
        ],
        **kwargs
    )
