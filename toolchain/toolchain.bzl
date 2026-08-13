"""toolchain.bzl provides the protoc toolchain rule"""

load("@rules_proto//proto:proto_common.bzl", "proto_common")

# PROTOC_TOOLCHAIN_TYPE is the toolchain type provided by this ruleset.
PROTOC_TOOLCHAIN_TYPE = Label("//toolchain:protoc")

# PROTO_TOOLCHAIN_TYPE is the toolchain type used by the `proto_toolchain` rule
# from protobuf/rules_proto.  This is the toolchain type that
# --incompatible_enable_proto_toolchain_resolution resolves `protoc` from, and
# the one registered by rulesets that supply a prebuilt protoc such as
# https://github.com/aspect-build/toolchains_protoc.
PROTO_TOOLCHAIN_TYPE = Label("@com_google_protobuf//bazel/private:proto_toolchain_type")

# INCOMPATIBLE_ENABLE_PROTO_TOOLCHAIN_RESOLUTION reports whether the
# --incompatible_enable_proto_toolchain_resolution flag is enabled.
INCOMPATIBLE_ENABLE_PROTO_TOOLCHAIN_RESOLUTION = getattr(
    proto_common,
    "INCOMPATIBLE_ENABLE_PROTO_TOOLCHAIN_RESOLUTION",
    False,
)

def use_protoc_toolchains():
    """Returns the list of toolchains that provide protoc.

    When --incompatible_enable_proto_toolchain_resolution is enabled, protoc is
    additionally sourced from the proto toolchain, and //toolchain:protoc
    becomes optional (a build that resolves protoc via the proto toolchain does
    not need to register one of ours at all).

    Returns:
        list: of toolchain types for the `toolchains` argument of a rule.
    """
    if not INCOMPATIBLE_ENABLE_PROTO_TOOLCHAIN_RESOLUTION:
        return [config_common.toolchain_type(PROTOC_TOOLCHAIN_TYPE, mandatory = True)]
    return [
        config_common.toolchain_type(PROTO_TOOLCHAIN_TYPE, mandatory = False),
        config_common.toolchain_type(PROTOC_TOOLCHAIN_TYPE, mandatory = False),
    ]

def find_protoc(ctx, override = None):
    """Resolves the protoc tool for a rule that uses `use_protoc_toolchains`.

    Resolution order:

    1. the `override` file, if given (typically the `protoc` rule attribute).
    2. the proto toolchain, if
       --incompatible_enable_proto_toolchain_resolution is enabled and such a
       toolchain is registered.
    3. the //toolchain:protoc toolchain.

    Args:
        ctx: the rule context.
        override: optional <File> that takes precedence over the toolchains.

    Returns:
        struct: having an `executable` <File> and a `tool` field suitable for
        the `tools` argument of a ctx.actions method.
    """
    if override:
        return struct(executable = override, tool = override)

    if INCOMPATIBLE_ENABLE_PROTO_TOOLCHAIN_RESOLUTION:
        proto_toolchain = ctx.toolchains[PROTO_TOOLCHAIN_TYPE]
        if proto_toolchain:
            proto_compiler = proto_toolchain.proto.proto_compiler
            return struct(executable = proto_compiler.executable, tool = proto_compiler)

    protoc_toolchain = ctx.toolchains[PROTOC_TOOLCHAIN_TYPE]
    if not protoc_toolchain:
        fail("no protoc toolchain was resolved for %s: register one of type '%s' " % (ctx.label, PROTOC_TOOLCHAIN_TYPE) +
             "(for example 'register_toolchains(\"@build_stack_rules_proto//toolchain:standard\")') " +
             "or, with --incompatible_enable_proto_toolchain_resolution, one of type '%s'" % PROTO_TOOLCHAIN_TYPE)

    return struct(
        executable = protoc_toolchain.protoc_executable,
        tool = protoc_toolchain.protoc_executable,
    )

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
