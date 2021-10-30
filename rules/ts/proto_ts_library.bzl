"proto_ts_library.bzl provides the proto_ts_library rule"

load("@build_bazel_rules_nodejs//:providers.bzl", "DeclarationInfo", "JSModuleInfo")

def _proto_ts_library_impl(ctx):
    """
    Implementation for proto_ts_library rule

    Args:
        ctx: the rule context object
    Returns:
        list of providers
    """

    # list<depset<File>>
    transitive_dts = [depset(ctx.files.deps)]

    # list<depset<File>>>
    transitive_js_modules = []

    # gather transitive .d.ts files
    for dep in ctx.attr.deps:
        if DeclarationInfo in dep:
            transitive_dts.append(dep[DeclarationInfo].transitive_declarations)
        if JSModuleInfo in dep:
            transitive_js_modules.append(dep[JSModuleInfo].sources)

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

    # for the JSModuleInfo provider
    js_module = depset(js_outputs)
    transitive_js_modules.append(js_module)

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
            direct_sources = js_module,
            sources = depset(transitive = transitive_js_modules),
        ),
    ]

proto_ts_library = rule(
    implementation = _proto_ts_library_impl,
    attrs = {
        "deps": attr.label_list(
            doc = "dependencies that provide .d.ts files (typically other proto_ts_library rules)",
            providers = [DeclarationInfo],
        ),
        "srcs": attr.label_list(
            allow_files = True,
            doc = "source .ts files",
        ),
        "tsc": attr.label(
            allow_files = True,
            cfg = "host",
            mandatory = True,
            doc = "typescript compiler executable",
            executable = True,
        ),
        # TODO(pcj): just use args here
        "tsc_options": attr.string_list(
            default = ["--esModuleInterop"],
            doc = "additional options for the tsc compile action",
        ),
    },
    # toolchains = ["@build_stack_rules_proto//rules/ts:ts_compiler"],
)
