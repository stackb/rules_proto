load("//:compile.bzl", "proto_compile")

def android_proto_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//android:java")),
        ],
        **kwargs
    )
