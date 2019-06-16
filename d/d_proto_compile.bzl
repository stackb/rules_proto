load("//:compile.bzl", "proto_compile")
load("//:plugin.bzl", "proto_plugin")

def d_proto_compile(**kwargs):
    # If package specified, declare a custom plugin that should correctly
    # predict the output location.
    package = kwargs.get("package")
    if package and not kwargs.get("plugins"):
        name_plugin = kwargs.get("name") + "_plugin"
        proto_plugin(
            name = name_plugin,
            outputs = ["{package}/%s/{basename}.d" % package],
            tool = "@com_github_dcarp_protobuf_d//:protoc-gen-d",
        )
        kwargs["plugins"] = [name_plugin]
        kwargs.pop("package")

    # Define the default plugin if still not defined
    if not kwargs.get("plugins"):
        kwargs["plugins"] = [Label("//d:d")]

    proto_compile(
        **kwargs
    )
