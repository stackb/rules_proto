"toolchain.bzl provides the protoc toolchain rule"

def _protoc_impl(ctx):
    return [platform_common.ToolchainInfo(
        protoc_target = ctx.attr.tool,
        protoc_executable = ctx.executable.tool,
    )]

protoc = rule(
    implementation = _protoc_impl,
    attrs = {
        "tool": attr.label(
            doc = "The protocol compiler tool",
            allow_single_file = True,
            executable = True,
            cfg = "exec",
        ),
    },
    provides = [platform_common.ToolchainInfo],
)
