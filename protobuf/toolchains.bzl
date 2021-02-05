def _protoc_toolchain_impl(ctx):
    return [platform_common.ToolchainInfo(
        protoc_target = ctx.attr.protoc,
        protoc_executable = ctx.executable.protoc,
    )]

protoc_toolchain = rule(
    implementation = _protoc_toolchain_impl,
    attrs = {
        "protoc": attr.label(
            doc = "The protocol compiler tool",
            default = "@com_google_protobuf//:protoc",
            executable = True,
            cfg = "exec",
        ),
    },
    provides = [platform_common.ToolchainInfo],
)
