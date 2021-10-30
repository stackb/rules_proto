"ts_proto_library.bzl provides the ts_proto_library rule"

load("@build_bazel_rules_nodejs//:providers.bzl", "DeclarationInfo", "JSModuleInfo")

def _ts_proto_library_impl(ctx):
    """
    Implementation for ts_proto_library rule
    """

    # .d.ts files that will be created by the tsc action
    dts_outputs = []

    # .js files that will be created by the tsc action
    js_outputs = []

    for f in ctx.files.srcs:
        base = f.basename[0:-len(f.extension)]
        js_outputs.append(ctx.actions.declare_file(base + "js", sibling = f))
        dts_outputs.append(ctx.actions.declare_file(base + "d.ts", sibling = f))

    # all outputs (.d.ts + .js)
    outputs = js_outputs + dts_outputs

    # tools we need for the action
    tools = [ctx.executable._tsc]

    # the action command
    command = [ctx.executable._tsc.path, " --declaration", "--esModuleInterop"] + [f.path for f in ctx.files.srcs]

    ctx.actions.run_shell(
        mnemonic = "TscProtoLibrary",
        inputs = depset(
            direct = ctx.files.srcs,
            transitive = [depset(ctx.files._ts_proto_deps)],
        ),
        tools = depset(tools),
        outputs = outputs,
        progress_message = "tsc %s" % ctx.label,
        command = " ".join(command),
    )

    files = depset(direct = outputs)

    transitive_declarations = depset(transitive = [])

    return [
        DefaultInfo(files = files),
        DeclarationInfo(
            declarations = dts_outputs,
            transitive_declarations = depset(transitive = [depset(dts_outputs)]),
            type_blocklisted_declarations = depset(),
        ),
        JSModuleInfo(
            direct_sources = depset(js_outputs),
            sources = depset(js_outputs),
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
        "_ts_proto_deps": attr.label_list(
            allow_files = True,
            default = [
                # "@npm//@nestjs/common",
                # "@npm//@nestjs/core",
                # "@npm//@nestjs/microservices",
                # "@npm//@types/bytebuffer",
                # "@npm//@types/node",
                # "@npm//axios",
                # "@npm//grpc",
                "@npm//long",
                "@npm//protobufjs",
                # "@npm//rxjs",
            ],
        ),
        "_tsc": attr.label(
            allow_files = True,
            executable = True,
            cfg = "host",
            default = Label("@npm//typescript/bin:tsc"),
        ),
    },
)
