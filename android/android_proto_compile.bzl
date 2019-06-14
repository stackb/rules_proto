load("//:compile.bzl", "proto_compile")

def android_proto_compile(**kwargs):
    # Append the android plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//android:javalite"),
    ]
    proto_compile(**kwargs)
