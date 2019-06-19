
def _protoc_toolchain_impl(ctx):
    return [platform_common.ToolchainInfo(
        protoc = ctx.executable.protoc,
    )]

protoc_toolchain = rule(
    implementation = _protoc_toolchain_impl,
    attrs = {
        "protoc": attr.label(
            doc = "The protocol compiler tool",
            default = "@com_google_protobuf//:protoc",
            executable = True,
            cfg = "host", # TODO: Change to exec when available
        ),
    },
    provides = [platform_common.ToolchainInfo],
)
