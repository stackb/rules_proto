load("//:compile.bzl", "proto_compile")

def android_grpc_compile(**kwargs):
    # Append the android plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//android:javalite"),
        Label("//android:grpc_javalite"),
    ]
    proto_compile(**kwargs)
