load("//:plugin.bzl", "ProtoPluginInfo")
load(
    "//:aspect.bzl",
    "ProtoLibraryAspectNodeInfo",
    "proto_compile_aspect_attrs",
    "proto_compile_aspect_impl",
    "proto_compile_attrs",
    "proto_compile_impl",
)

# "Aspects should be top-level values in extension files that define them."

ruby_grpc_compile_aspect = aspect(
    implementation = proto_compile_aspect_impl,
    provides = ["proto_compile", ProtoLibraryAspectNodeInfo],
    attr_aspects = ["deps"],
    attrs = dict(
        proto_compile_aspect_attrs,
        _plugins = attr.label_list(
            doc = "List of protoc plugins to apply",
            providers = [ProtoPluginInfo],
            default = [
                str(Label("//ruby:ruby")),
                str(Label("//ruby:grpc_ruby")),
            ],
        ),
    ),
)

_rule = rule(
    implementation = proto_compile_impl,
    attrs = dict(
        proto_compile_attrs,
        deps = attr.label_list(
            mandatory = True,
            providers = ["proto", "proto_compile", ProtoLibraryAspectNodeInfo],
            aspects = [ruby_grpc_compile_aspect],
        ),
    ),
)

def ruby_grpc_compile(**kwargs):
    _rule(
        verbose_string = "%s" % kwargs.get("verbose", 0),
        plugin_options_string = ";".join(kwargs.get("plugin_options", [])),
        **kwargs
    )
