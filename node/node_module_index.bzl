load("//:compile.bzl", "ProtoCompileInfo")

NodeModuleIndexInfo = provider(fields = {
    "name": "rule name",
    "index": "index.js file",
})

def _get_js_variable_name(file):
    name = file.basename.rstrip(".js")

    # Deal with special characters here?
    return name

def _get_js_output_file_name(ctx, file):
    filename = file.short_path
    filename = filename[len(ctx.label.package) + 1:]
    return filename

def _node_module_index_impl(ctx):
    compilation = ctx.attr.compilation[ProtoCompileInfo]

    index_js = ctx.actions.declare_file("%s/index.js" % (compilation.label.name))

    exports = {}

    for output in compilation.outputs:
        if output.path.endswith("_pb.js"):
            name = _get_js_variable_name(output)
            exports[name] = _get_js_output_file_name(ctx, output)
        elif output.path.endswith("_grpc_pb.js"):
            name = _get_js_variable_name(output)
            exports[name] = _get_js_output_file_name(ctx, output)

    content = []
    content.append("module.exports = {")
    for name, path in exports.items():
        content.append("    '%s': require('./%s')," % (name, path))
    content.append("}")

    ctx.actions.write(
        output = index_js,
        content = "\n".join(content),
    )

    return [NodeModuleIndexInfo(
        name = ctx.label.name,
        index = index_js,
    ), DefaultInfo(
        files = depset([index_js]),
    )]

node_module_index = rule(
    implementation = _node_module_index_impl,
    attrs = {
        "compilation": attr.label(
            providers = [ProtoCompileInfo],
            mandatory = True,
        ),
    },
    output_to_genfiles = True,
)
