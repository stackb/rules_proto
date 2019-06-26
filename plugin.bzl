ProtoPluginInfo = provider(fields = {
    "name": "The proto plugin name",
    "options": "A list of options to pass to the compiler for this plugin",
    "outputs": "Output filenames generated on a per-proto basis. Example: '{basename}_pb2.py",
    "out": "Output filename generated on a per-plugin basis; to be used in the value for --NAME-out=OUT",
    "tool": "The plugin binary. If absent, it is assumed the plugin is built-in to protoc itself",
    "tool_executable": "The plugin binary executable. If absent, it is assumed the plugin is built-in to protoc itself",
    "exclusions": "Exclusion filters to apply when generating outputs with this plugin. Used to prevent generating files that are included in the protobuf library, for example. Can exclude either by proto name prefix or by proto folder prefix",
    "data": "Additional files required for running the plugin",
})


def _proto_plugin_impl(ctx):
    # Build ProtoPluginInfo provider
    return [
        ProtoPluginInfo(
            name = ctx.attr.name,
            options = ctx.attr.options,
            outputs = ctx.attr.outputs,
            out = ctx.attr.out,
            tool = ctx.attr.tool,
            tool_executable = ctx.executable.tool,
            exclusions = ctx.attr.exclusions,
            data = ctx.files.data,
        )
    ]


proto_plugin = rule(
    implementation = _proto_plugin_impl,
    attrs = {
        "options": attr.string_list(
            doc = "A list of options to pass to the compiler for this plugin",
        ),
        "outputs": attr.string_list(
            doc = "Output filenames generated on a per-proto basis. Example: '{basename}_pb2.py'",
        ),
        "out": attr.string(
            doc = "Output filename generated on a per-plugin basis; to be used in the value for --NAME-out=OUT",
        ),
        "tool": attr.label(
            doc = "The plugin binary. If absent, it is assumed the plugin is built-in to protoc itself",
            cfg = "host",
            allow_files = True,
            executable = True,
        ),
        "exclusions": attr.string_list(
            doc = "Exclusion filters to apply when generating outputs with this plugin. Used to prevent generating files that are included in the protobuf library, for example. Can exclude either by proto name prefix or by proto folder prefix",
        ),
        "data": attr.label_list(
            doc = "Additional files required for running the plugin",
            allow_files = True,
        ),
    },
)
