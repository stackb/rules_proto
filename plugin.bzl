ProtoPluginInfo = provider(fields = {
    "name": "proto plugin name",
    "outputs": "outputs to be generated",
    "tool": "plugin tool",
    "executable": "plugin tool executable",
    "options": "proto options",
    "out": "aggregate proto output",
    "outdir": "whether to use the package output dir",
    "data": "additional data"
})

def _proto_plugin_impl(ctx):
    return [ProtoPluginInfo(
        name = ctx.label.name,
        options = ctx.attr.options,
        out = ctx.attr.out,
        outdir = ctx.attr.outdir,
        outputs = ctx.attr.outputs,
        tool = ctx.attr.tool,
        executable = ctx.executable.tool,
        data = ctx.files.data,
    )]

proto_plugin = rule(
    implementation = _proto_plugin_impl,
    attrs = {
        "options": attr.string_list(
            doc = "An list of options to pass to the compiler.",
        ),
        "outputs": attr.string_list(
            doc = "Output filenames generated on a per-proto basis.  Example: '{basename}_pb2.py'",
        ),
        "out": attr.string(
            doc = "Output filename generated on a per-plugin basis; to be used in the value for --NAME-out=OUT",
        ),
        "outdir": attr.string(
            doc = "If present, overrides the file.path from out; to be used in the value for --NAME-out=OUT",
        ),
        "tool": attr.label(
            doc = "The plugin binary.  If absent, assume the plugin is a built-in to protoc itself",
            cfg = "host",
            allow_files = True,
            executable = True,
        ),
        "data": attr.label_list(
            doc = "Additional files that should travel with the plugin",
            allow_files = True,
        ),
    }
)

