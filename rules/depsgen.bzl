"""depsgen.bzl
"""

load("@build_stack_rules_proto//rules:providers.bzl", "ProtoDependencyInfo")

def _depsgen_impl(ctx):
    config_json = ctx.outputs.json
    output_deps = ctx.outputs.deps

    config = struct(
        out = output_deps.path,
        name = ctx.label.name,
        deps = [dep[ProtoDependencyInfo] for dep in ctx.attr.deps],
    )

    ctx.actions.write(
        output = config_json,
        content = config.to_json(),
    )

    ctx.actions.run(
        mnemonic = "DepsGenerate",
        progress_message = "Generating %s deps" % ctx.attr.name,
        executable = ctx.file._depsgen,
        arguments = ["--config_json=%s" % config_json.path],
        inputs = [config_json],
        outputs = [output_deps],
    )

    return [DefaultInfo(
        files = depset([config_json, output_deps]),
    )]

depsgen = rule(
    implementation = _depsgen_impl,
    attrs = {
        "deps": attr.label_list(
            doc = "Top level dependencies to compute",
            providers = [ProtoDependencyInfo],
        ),
        "_depsgen": attr.label(
            doc = "The depsgen tool",
            default = "//cmd/depsgen",
            allow_single_file = True,
            executable = True,
            cfg = "host",
        ),
    },
    outputs = {
        "json": "%{name}.json",
        "deps": "%{name}_deps.bzl",
    },
)
