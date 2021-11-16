"""proto_plugin.bzl provides the "proto_plugin" rule.

A "proto_plugin" rule wraps metadata about a proto compiler plugin.
"""

load("@rules_proto//proto:defs.bzl", "ProtoInfo")
load("@build_stack_rules_proto//rules:providers.bzl", "ProtoDependencyInfo", "ProtoPluginInfo")

def _proto_plugin_impl(ctx):
    return [
        ProtoPluginInfo(
            name = ctx.attr.name,
            label = ctx.label,
            options = ctx.attr.options,
            out = ctx.attr.out,
            tool = ctx.executable.tool,
            tool_target = ctx.attr.tool,
            use_built_in_shell_environment = ctx.attr.use_built_in_shell_environment,
            protoc_plugin_name = ctx.attr.protoc_plugin_name,
            exclusions = ctx.attr.exclusions,
            mods = ctx.attr.mods,
            data = ctx.files.data,
            supplementary_proto_deps = [dep[ProtoInfo] for dep in ctx.attr.supplementary_proto_deps],
            deps = [dep[ProtoDependencyInfo] for dep in ctx.attr.deps],
        ),
    ]

proto_plugin = rule(
    implementation = _proto_plugin_impl,
    attrs = {
        "tool": attr.label(
            doc = "The plugin binary. If absent, it is assumed the plugin is built-in to protoc itself and builtin_plugin_name will be used if available, otherwise the plugin name",
            cfg = "exec",
            allow_files = True,
            executable = True,
        ),
        "use_built_in_shell_environment": attr.bool(
            doc = "Whether the tool should use the built in shell environment or not",
            default = False,
        ),
        "protoc_plugin_name": attr.string(
            doc = "The name used for the plugin binary on the protoc command line. Useful for targeting built-in plugins. Uses plugin name when not set",
        ),
        "options": attr.string_list(
            doc = "A list of options to pass to the compiler for this plugin",
        ),
        "out": attr.string(
            doc = "The output scheme for the plugin.  Can be a string like '.' or a symbol such as {BIN_DIR} or {PACKAGE}.",
            default = "{BIN_DIR}",
        ),
        "exclusions": attr.string_list(
            doc = "Exclusion filters to apply when generating outputs with this plugin. Used to prevent generating files that are included in the protobuf library, for example. Can exclude either by proto name prefix or by proto folder prefix",
        ),
        "data": attr.label_list(
            doc = "Additional files required for running the plugin",
            allow_files = True,
        ),
        "supplementary_proto_deps": attr.label_list(
            doc = "Additional proto files/descriptors to be placed on the argument list",
            allow_files = True,
            providers = [ProtoInfo],
        ),
        "deps": attr.label_list(
            doc = "Workspace dependencies the plugin tool requires",
            providers = [ProtoDependencyInfo],
        ),
        "mods": attr.string_dict(
            doc = "content modifications to apply to the output files",
        ),
    },
)
