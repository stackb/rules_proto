load(":proto_compile.bzl", "coalesce_outputs", "get_output_directories", "proto_compile")
load(":proto_providers.bzl", "ProtoCompilationInfo", "ProtoCompileOutputInfo")
load("@build_stack_rules_proto//rules:proto_plugin.bzl", "ProtoPluginInfo")
load("@rules_proto//proto:defs.bzl", "ProtoInfo")
load("proto_utils.bzl", "get_bool_attr", "get_int_attr", "get_protoc_executable")

def _proto_compile_aspect_rule_impl(ctx):
    """The implementation function for the aspect rule.

    When this function is evaluated, bazel has visited all the deps via
    the aspect mechanism and populated the ctx.attr.deps with ProtoCompileOutputInfo.

    Args:
        ctx: <ctx> the bazel ctx object
    Returns:
        <list<providers>> the list of providers
    """
    return coalesce_outputs(ctx, [dep[ProtoCompileOutputInfo] for dep in ctx.attr.deps])

def _proto_compile_impl(ctx):
    """The implementation function for the non-aspect rule.

    When this function is evaluated, a compilation object is constructed
    from the given ctx.  The compilation object is then executed and the
    ProtoCompileOutputInfo are coalesced into the final list of providers.

    Args:
        ctx: <ctx> the bazel ctx object
    Returns:
        <list<providers>> the list of providers
    """

    rel_outdir, full_outdir = get_output_directories(ctx.bin_dir, ctx.label, "")
    proto_infos = [dep[ProtoInfo] for dep in ctx.attr.deps]
    plugins = [plugin[ProtoPluginInfo] for plugin in ctx.attr.plugins]

    compilation = ProtoCompilationInfo(
        actions = ctx.actions,
        bin_dir = ctx.bin_dir,
        full_outdir = full_outdir,
        host_path_separator = ctx.configuration.host_path_separator,
        label = ctx.label,
        output_directory_prefix = "",
        plugins = plugins,
        proto_info = proto_infos[0],
        protoc = get_protoc_executable(ctx),
        rel_outdir = rel_outdir,
        resolve_tools = ctx.resolve_tools,
        single_action = ctx.attr.single_action,
        transitive = ctx.attr.transitive,
        transitive_outs = [],
        verbose = ctx.attr.verbose,
    )

    return coalesce_outputs(ctx, proto_compile(compilation))

def _proto_compile_aspect_impl(target, ctx):
    """The implementation function for the aspect function.

    When this function is evaluated, bazel is visiting one of the deps via
    the aspect mechanism.

    Args:
        target: <target> the bazel target object being visited
        ctx: <ctx> the bazel ctx (aspect variant) object
    Returns:
        <list<ProtoCompileOutputInfo>> the list of providers
    """

    rel_outdir, full_outdir = get_output_directories(ctx.bin_dir, ctx.label, ctx.attr._prefix)
    proto_info = target[ProtoInfo]
    plugins = [plugin[ProtoPluginInfo] for plugin in ctx.attr._plugins]
    transitive_outs = [dep[ProtoCompileOutputInfo] for dep in ctx.rule.attr.deps]

    return proto_compile(ProtoCompilationInfo(
        actions = ctx.actions,
        bin_dir = ctx.bin_dir,
        full_outdir = full_outdir,
        host_path_separator = ctx.configuration.host_path_separator,
        label = ctx.label,
        output_directory_prefix = ctx.attr._prefix,
        plugins = plugins,
        proto_info = proto_info,
        protoc = get_protoc_executable(ctx),
        rel_outdir = rel_outdir,
        resolve_tools = ctx.resolve_tools,
        single_action = get_bool_attr(ctx.attr, "single_action_string"),
        transitive = get_bool_attr(ctx.attr, "transitive_string"),
        transitive_outs = transitive_outs,
        verbose = get_int_attr(ctx.attr, "verbose_string"),
    ))

_proto_compile_attrs = {
    "verbose": attr.int(
        doc = "The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*",
    ),
    "single_action": attr.bool(
        doc = "Whether to run all plugins in a single action",
        default = False,
    ),
    "transitive": attr.bool(
        doc = "Whether to compile transitive deps",
        default = False,
    ),
    "prefix_path": attr.string(
        doc = "Path to prefix to the generated files in the output directory. Cannot be set when merge_directories == False",
    ),
    "merge_directories": attr.bool(
        doc = "If true, all generated files are merged into a single directory with the name of current label and these new files returned as the outputs. If false, the original generated files are returned across multiple roots",
        default = True,
    ),
    "output_in_package_root": attr.bool(
        doc = "If true, copy the output files into the package output directory rather than a subdirectory thereof",
        default = True,
    ),
}

_proto_compile_aspect_attrs = {
    "verbose_string": attr.string(
        doc = "String version of the verbose string, used for aspect",
        values = ["", "None", "0", "1", "2", "3", "4"],
        default = "0",
    ),
    "single_action_string": attr.string(
        doc = "A boolean (as a string) on whether to run all plugins in a single action",
        values = ["True", "False"],
        default = "False",
    ),
    "transitive_string": attr.string(
        doc = "A boolean (as a string) on whether to compile transitive deps",
        values = ["True", "False"],
        default = "False",
    ),
}

def proto_compile_rule(default_plugins):
    """proto_compile_rule returns a new rule using the given default plugins.

    This rule does not use an aspect.

    Args:
        default_plugins: <list<string>> A list of labels the provide ProtoPluginInfo.
    Returns:
        <rule> the rule function
    """
    return rule(
        implementation = _proto_compile_impl,
        attrs = dict(
            _proto_compile_attrs,
            plugins = attr.label_list(
                doc = "List of protoc plugins to apply",
                providers = [ProtoPluginInfo],
                default = default_plugins,
            ),
            deps = attr.label_list(
                mandatory = True,
                providers = [ProtoInfo],
            ),
        ),
        toolchains = [str(Label("@build_stack_rules_proto//toolchains:protoc_toolchain_type"))],
    )

def proto_compile_aspect_rule(aspect):
    """proto_compile_aspect_rule returns a new aspect rule using the given aspect function.

    Args:
        aspect: <aspect> The aspect function
    Returns:
        <rule> the aspect rule function
    """
    return rule(
        implementation = _proto_compile_aspect_rule_impl,
        attrs = dict(
            _proto_compile_aspect_attrs.items() + _proto_compile_attrs.items(),
            deps = attr.label_list(
                mandatory = True,
                providers = [ProtoInfo, ProtoCompileOutputInfo],
                aspects = [aspect],
            ),
        ),
    )

def proto_compile_aspect_rule_macro(aspect_rule, **kwargs):
    """Wraps an aspect rule with a macro that pre-populates the *_string atttributes.

    Args:
        aspect_rule: <rule> the rule function to wrap.
        **kwargs: <*> items for the rule function
    """
    aspect_rule(
        single_action_string = "{}".format(kwargs.get("single_action", False)),
        verbose_string = "{}".format(kwargs.get("verbose", 0)),
        transitive_string = "{}".format(kwargs.get("transitive", False)),
        merge_directories = kwargs.pop("merge_directories", True),
        **kwargs
    )

def proto_compile_aspect(default_plugins, default_prefix):
    """proto_compile_aspect returns a new aspect having the given default plugins and prefix.

    Args:
        default_plugins: <list<string>> A list of labels the provide ProtoPluginInfo.
        default_prefix: <string> An optional prefix to disambiguate output directory paths.
    Returns:
        <aspect> the aspect function
    """
    return aspect(
        implementation = _proto_compile_aspect_impl,
        provides = [ProtoCompileOutputInfo],
        attr_aspects = ["deps"],
        attrs = dict(
            _proto_compile_aspect_attrs,
            _plugins = attr.label_list(
                doc = "List of protoc plugins to apply",
                providers = [ProtoPluginInfo],
                default = default_plugins,
            ),
            _prefix = attr.string(
                doc = "String used to disambiguate aspects when generating outputs",
                default = default_prefix,
            ),
        ),
        toolchains = [str(Label("@build_stack_rules_proto//toolchains:protoc_toolchain_type"))],
    )
