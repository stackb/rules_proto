"ts_proto_library.bzl provides the ts_proto_library rule"

load("@build_bazel_rules_nodejs//:providers.bzl", "DeclarationInfo", "JSModuleInfo")

def _ts_proto_library_impl(ctx):
    """
    Implementation for ts_proto_library rule

    Args:
        ctx: the rule context object
    Returns:
        list of providers
    """

    # list<depset<File>>
    transitive_dts = [depset(ctx.files.npm_deps)]

    # gather transitive .d.ts files
    for info in [dep[DeclarationInfo] for dep in ctx.attr.deps]:
        transitive_dts.append(info.transitive_declarations)

    # list<File>: .d.ts files that will be created by the tsc action
    dts_outputs = []

    # .js files that will be created by the tsc action
    js_outputs = []

    # all srcs files are expected to be .ts files.
    for f in ctx.files.srcs:
        (base, _, _) = f.basename.rpartition(".")
        dts_outputs.append(
            ctx.actions.declare_file(base + ".d.ts", sibling = f),
        )
        js_outputs.append(
            ctx.actions.declare_file(base + ".js", sibling = f),
        )

    # all outputs (.d.ts + .js)
    outputs = js_outputs + dts_outputs

    # build the tsc command
    command = [ctx.executable.tsc.path, " --declaration"]
    command += ctx.attr.tsc_options
    command += [f.path for f in ctx.files.srcs]

    ctx.actions.run_shell(
        command = " ".join(command),
        inputs = depset(
            direct = ctx.files.srcs,
            transitive = transitive_dts,
        ),
        mnemonic = "TscProtoLibrary",
        outputs = outputs,
        progress_message = "tsc %s" % ctx.label,
        tools = [ctx.executable.tsc],
    )

    return [
        # primary provider needed for downstream rules such as ts_project.
        DeclarationInfo(
            declarations = dts_outputs,
            # ts_project wants direct /and/ transitive .d.ts in 'transitive'
            transitive_declarations = depset(
                transitive = [depset(dts_outputs)] + transitive_dts,
            ),
            type_blocklisted_declarations = depset(),
        ),
        # default info contains all files, in case someone wanted it for a
        # filegroup etc.
        DefaultInfo(files = depset(direct = outputs)),
        # js files, in the event someone would like this for js_library etc.
        JSModuleInfo(
            direct_sources = depset(js_outputs),
            sources = depset(js_outputs),
        ),
    ]

ts_proto_library = rule(
    implementation = _ts_proto_library_impl,
    attrs = {
        "deps": attr.label_list(
            doc = "dependencies that provide .d.ts files (typically other ts_proto_library rules)",
            providers = [DeclarationInfo],
        ),
        "npm_deps": attr.label_list(
            allow_files = True,
            default = [
                "@npm//long",
                "@npm//protobufjs",
            ],
            doc = "additional npm library dependencies",
        ),
        "srcs": attr.label_list(
            allow_files = True,
            doc = "source .ts files",
        ),
        "tsc": attr.label(
            allow_files = True,
            cfg = "host",
            default = Label("@npm//typescript/bin:tsc"),
            doc = "typescript compiler executable",
            executable = True,
        ),
        "tsc_options": attr.string_list(
            default = ["--esModuleInterop"],
            doc = "additional options for the tsc compile action",
        ),
    },
)
