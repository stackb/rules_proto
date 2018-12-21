load("//:compile.bzl", "proto_compile_attrs", "proto_compile_impl")
load("//:aspect.bzl", "ProtoLibraryAspectNodeInfo", "proto_compile_aspect_attrs", "proto_compile_aspect_impl")
load("//:plugin.bzl", "ProtoPluginInfo")

# "Aspects should be top-level values in extension files that define them."

_aspect = aspect(
    implementation = proto_compile_aspect_impl,
    attr_aspects = ["deps"],
    attrs = proto_compile_aspect_attrs + {
        "_plugins": attr.label_list(
            doc = "List of protoc plugins to apply",
            providers = [ProtoPluginInfo],
            default = [
                str(Label("//github.com/grpc/grpc-web:commonjs_dts")),
            ],
        ),
    },
)

_rule = rule(
    implementation = proto_compile_impl,
    attrs = proto_compile_attrs + {
        "deps": attr.label_list(
            mandatory = True,
            providers = ["proto", ProtoLibraryAspectNodeInfo],
            aspects = [_aspect],
        ),    
    },
)

def commonjs_dts_grpc_compile(**kwargs):
    _rule(
        verbose_string = "%s" % kwargs.get("verbose", 0),
        plugin_options_string = ";".join(kwargs.get("plugin_options", [])),
        **kwargs)

