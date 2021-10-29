"ts_proto_library.bzl provides the ts_proto_library rule"

load("@build_bazel_rules_nodejs//:providers.bzl", "DeclarationInfo")

def _ts_proto_library_impl(ctx):
    """
    Implementation for ts_proto_library rule
    """
    dts_outputs = []
    js_outputs = []
    es5_outputs = []

    for f in ctx.files.srcs:
        dts_outputs.append(f)

    direct_outputs = js_outputs + dts_outputs
    outputs = depset(direct = direct_outputs)
    transitive_declarations = depset()  # TODO: actually gather this

    return [
        DefaultInfo(files = outputs),
        DeclarationInfo(
            declarations = dts_outputs,
            transitive_declarations = depset(dts_outputs),
            type_blocklisted_declarations = depset(),
        ),
    ]

ts_proto_library = rule(
    implementation = _ts_proto_library_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
        ),
        "deps": attr.label_list(
            providers = [DeclarationInfo],
        ),
        "_tsc": attr.label(
            allow_files = True,
            executable = True,
            cfg = "host",
            default = Label("@npm//typescript/bin:tsc"),
        ),
    },
)
