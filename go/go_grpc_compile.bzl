load("//:compile.bzl", "proto_compile")
load("//:plugin.bzl", "proto_plugin")

def go_grpc_compile(**kwargs):
    # If importpath specified, declare a custom plugin that should correctly
    # predict the output location.
    importpath = kwargs.get("importpath")
    if importpath and not kwargs.get("plugins"):
        name = kwargs.get("name")
        name_plugin = name + "_plugin"
        proto_plugin(
            name = name_plugin,
            outputs = ["{package}/%s/{basename}.pb.go" % importpath],
            tool = "@com_github_golang_protobuf//protoc-gen-go",
        )
        kwargs["plugins"] = [name_plugin]
        kwargs.pop("importpath")
    # Define the default plugin if still not defined
    if not kwargs.get("plugins"):
        kwargs["plugins"] = [str(Label("//go:grpc_go"))]

    proto_compile(
        **kwargs
    )
