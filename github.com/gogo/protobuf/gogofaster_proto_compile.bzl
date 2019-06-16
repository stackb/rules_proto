load("//:compile.bzl", "proto_compile")
load("//:plugin.bzl", "proto_plugin")

def gogofaster_proto_compile(**kwargs):
    # If importpath specified, declare a custom plugin that should correctly
    # predict the output location.
    importpath = kwargs.get("importpath")
    if importpath and not kwargs.get("plugins"):
        name_plugin = kwargs.get("name") + "_plugin"
        proto_plugin(
            name = name_plugin,
            outputs = ["{package}/%s/{basename}.pb.go" % importpath],
            tool = "@com_github_gogo_protobuf//protoc-gen-gogofaster",
        )
        kwargs["plugins"] = [name_plugin]
        kwargs.pop("importpath")

    # Define the default plugin if still not defined
    if not kwargs.get("plugins"):
        kwargs["plugins"] = [Label("//github.com/gogo/protobuf:gogofaster")]

    proto_compile(
        **kwargs
    )
