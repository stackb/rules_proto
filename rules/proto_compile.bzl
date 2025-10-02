"""proto_compile.bzl provides the proto_compile rule.

This runs the protoc tool and generates output source files.
"""

load("@rules_proto//proto:defs.bzl", "ProtoInfo")
load(":providers.bzl", "ProtoCompileInfo", "ProtoPluginInfo")

def _uniq(iterable):
    """Returns a list of unique elements in `iterable`.

    Requires all the elements to be hashable.
    Args:
      iterable: An iterable to filter.
    Returns:
      A new list with all unique elements from `iterable`.
    """
    unique_elements = {element: None for element in iterable}

    return list(unique_elements.keys())

def _ctx_replace_args(ctx, args):
    return [_ctx_replace_arg(ctx, arg) for arg in args]

def _ctx_replace_arg(ctx, arg):
    arg = arg.replace("{BIN_DIR}", ctx.bin_dir.path)
    arg = arg.replace("{PACKAGE}", ctx.label.package)
    arg = arg.replace("{NAME}", ctx.label.name)

    if arg.find("{PROTO_LIBRARY_BASENAME}") != -1:
        basename = ctx.attr.proto.label.name
        if basename.endswith("_proto"):
            basename = basename[:len(basename) - len("_proto")]
        arg = arg.replace("{PROTO_LIBRARY_BASENAME}", basename)
    return arg

def _plugin_label_key(label):
    """_plugin_label_key converts a label into a string.  

    This is needed due to an edge case about how Labels are parsed and
    represented. Consider the label
    "@build_stack_rules_proto//plugin/scalapb/scalapb:protoc-gen-scala". If this
    string is the value for an attr.label in the same workspace
    build_stack_rules_proto, the workspace name is actually ommitted and becomes
    the empty string.  However, if is is the value for an attr.string and then
    parsed into a label in Starlark, the workspace name is preserved.  To resolve
    this issue, we just ignore the workspace name altogether, hoping that no-one
    tries to use two different plugins having a different workspace_name but
    otherwise identical package and name.
    """
    key = "%s:%s" % (label.package, label.name)

    return key

def get_protoc_executable(ctx):
    if ctx.file.protoc:
        return ctx.file.protoc
    protoc_toolchain_info = ctx.toolchains[str(Label("//toolchain:protoc"))]
    return protoc_toolchain_info.protoc_executable

def _descriptor_proto_path(proto, proto_info):
    """Convert a proto File to the path within the descriptor file.

    Adapted from https://github.com/bazelbuild/rules_go
    """

    # Strip proto_source_root
    path = _strip_path_prefix(proto.path, proto_info.proto_source_root)

    # Strip root
    path = _strip_path_prefix(path, proto.root.path)

    # Strip workspace root
    path = _strip_path_prefix(path, proto.owner.workspace_root)

    return path

def _strip_path_prefix(path, prefix):
    """Strip a prefix from a path if it exists and any remaining prefix slashes

    Args:
        path: <string>
        prefix: <string>
    Returns:
        <string>
    """
    if path.startswith(prefix):
        path = path[len(prefix):]
    if path.startswith("/"):
        path = path[1:]
    return path

def is_windows(ctx):
    return ctx.configuration.host_path_separator == ";"

def _proto_compile_impl(ctx):
    ###
    ### Part 1: setup variables used in scope
    ###

    # out_dir is used in conjunction with file.short_path to determine root
    # output file paths
    out_dir = ctx.bin_dir.path
    if ctx.label.workspace_root:
        out_dir = "/".join([out_dir, ctx.label.workspace_root])

    if len(ctx.attr.srcs) > 0 and len(ctx.outputs.outputs) > 0:
        fail("rule must provide 'srcs' or 'outputs' (but not both)")

    # <dict<string,File>: output files mapped by their package-relative path.
    # This struct is given to the provider.
    output_files_by_rel_path = {}

    # const <dict<string,string>.  The key is the file basename, value is the
    # short_path of the output file.
    output_short_paths_by_basename = {}

    # renames is a mapping from the output filename that was produced by the
    # plugin to the actual name we want to output.
    renames = {}

    if len(ctx.attr.srcs):
        # assume filenames in srcs are already package-relative
        for name in ctx.attr.srcs:
            rel = "/".join([ctx.label.package, name])
            actual_name = name + ctx.attr.output_file_suffix
            if actual_name != name:
                renames[rel] = "/".join([ctx.label.package, actual_name])
            f = ctx.actions.declare_file(actual_name)
            output_files_by_rel_path[rel] = f
            output_short_paths_by_basename[name] = rel
    else:
        for f in ctx.outputs.outputs:
            # rel = _get_package_relative_path(ctx.label, f.short_path)
            rel = f.short_path
            output_files_by_rel_path[rel] = f
            output_short_paths_by_basename[f.basename] = rel

    # const <bool> verbosity flag
    verbose = ctx.attr.verbose

    # const <File> the protoc file from the toolchain
    protoc = get_protoc_executable(ctx)

    # const <ProtoInfo> proto provider
    proto_info = ctx.attr.proto[ProtoInfo]

    # const <list<ProtoPluginInfo>> plugins to be applied
    plugins = [plugin[ProtoPluginInfo] for plugin in ctx.attr.plugins]

    # const <dict<string,string>>
    outs = {_plugin_label_key(Label(k)): v for k, v in ctx.attr.outs.items()}

    # mut <list<File>> set of descriptors for the compile action
    descriptors = proto_info.transitive_descriptor_sets.to_list()

    # mut <list<File>> tools for the compile action
    tools = [protoc]

    # mut <list<string>> argument list for protoc execution
    args = [] + ctx.attr.args

    # mut <list<File>> inputs for the compile action
    inputs = []

    # mut <list<File>> The (filtered) set of .proto files to compile
    protos = []

    # mut <list<opaque>> Plugin input manifests
    input_manifests = []

    # mut <dict<string,string>> post-processing modifications for the compile action
    mods = dict()

    ###
    ### Part 2: per-plugin args
    ###

    for plugin in plugins:
        ### Part 2.1: build protos list

        # add all protos unless excluded
        for proto in proto_info.direct_sources:
            if any([
                proto.dirname.endswith(exclusion) or proto.path.endswith(exclusion)
                for exclusion in plugin.exclusions
            ]) or proto in protos:  # TODO: When using import_prefix, the ProtoInfo.direct_sources list appears to contain duplicate records, this line removes these. https://github.com/bazelbuild/bazel/issues/9127
                continue

            # Proto not excluded
            protos.append(proto)

        # augment proto list with those attached to plugin
        for info in plugin.supplementary_proto_deps:
            for src in info.direct_sources:
                protos.append(src)
            descriptors += info.transitive_descriptor_sets.to_list()

        # Include extra plugin data files
        inputs += plugin.data

        ### Part 2.2: build --plugin argument

        # const <string> The name of the plugin
        plugin_name = plugin.protoc_plugin_name if plugin.protoc_plugin_name else plugin.name

        # const <?File> Add plugin executable if not a built-in plugin
        plugin_tool = plugin.tool if plugin.tool else None
        is_builtin = plugin.tool == None

        # Add plugin runfiles if plugin has a tool
        if plugin_tool:
            tools.append(plugin_tool)
            tools.append(plugin.tool_target[DefaultInfo].files_to_run)

            # If Windows, mangle the path.
            plugin_tool_path = plugin_tool.path
            if is_windows(ctx):
                plugin_tool_path = plugin_tool.path.replace("/", "\\")

            args.append("--plugin=protoc-gen-{}={}".format(plugin_name, plugin_tool_path))

        ### Part 2.3: build --{name}_out=OPTIONS argument

        # mut <string>
        out = plugin.out
        if ctx.label.workspace_root:
            # special handling for "{BIN_DIR}".  If we are dealing with a
            # formatted output string (like for a .srcjar), cannot just append
            # "external/repo" to the string.
            if out.find("{BIN_DIR}") != -1:
                out = out.replace("{BIN_DIR}", "{BIN_DIR}/" + ctx.label.workspace_root)
            else:
                out = "/".join([out, ctx.label.workspace_root])

        # dict<key=label.package+label.name,value=list<string>>
        options = {_plugin_label_key(Label(k)): v for k, v in ctx.attr.options.items()}

        # const <list<string>>
        opts = plugin.options + [opt for opt in options.get(_plugin_label_key(plugin.label), [])]
        if is_builtin and opts:
            # builtins can't use the --opt flags
            out = "{}:{}".format(",".join(opts), out)
        else:
            for opt in opts:
                args.append("--{}_opt={}".format(plugin_name, opt))

        # override with the out configured on the rule if specified
        plugin_out = outs.get(_plugin_label_key(plugin.label), None)
        if plugin_out:
            # bin-dir relative is implied for plugin_out overrides.  Workspace
            # root might be empty, so filter empty strings via this list
            # comprehension.
            out = "/".join([e for e in [ctx.bin_dir.path, ctx.label.workspace_root, plugin_out] if e])
        args.append("--{}_out={}".format(plugin_name, out))

        ### Part 2.4: setup awk modifications if any
        for k, v in plugin.mods.items():
            mods[k] = v

    ###
    ### Part 3: trailing args
    ###

    ### Part 3.1: add descriptor sets

    descriptors = _uniq(descriptors)
    inputs += descriptors

    args.append("--descriptor_set_in={}".format(ctx.configuration.host_path_separator.join(
        [d.path for d in descriptors],
    )))

    ### Part 3.2: add proto file args

    protos = _uniq(protos)
    for proto in protos:
        args.append(_descriptor_proto_path(proto, proto_info))

    ### Step 3.3: build args object

    replaced_args = _ctx_replace_args(ctx, _uniq(args))
    final_args = ctx.actions.args()
    final_args.use_param_file("@%s", use_always = False)
    final_args.add_all(replaced_args)

    ###
    ### Step 4: command action
    ###
    commands = [
        "set -euo pipefail",
        "mkdir -p ./" + ctx.label.package,
        protoc.path + " $@",  # $@ is replaced with args list
    ]

    # if the rule declares any mappings, setup copy file commands to move them
    # into place
    if len(ctx.attr.output_mappings) > 0:
        copy_commands = []
        for mapping in ctx.attr.output_mappings:
            basename, _, intermediate_filename = mapping.partition("=")
            output_short_path = output_short_paths_by_basename.get(basename)
            if not output_short_path:
                fail("the mapped file '%s' was not listed in outputs" % basename)
            copy_commands.append("cp '{dir}/{src}' '{dir}/{dst}'".format(
                dir = out_dir,
                src = intermediate_filename,
                dst = output_short_path,
            ))
        copy_script = ctx.actions.declare_file(ctx.label.name + "_copy.sh")
        ctx.actions.write(copy_script, "\n".join(copy_commands), is_executable = True)
        inputs.append(copy_script)
        commands.append(copy_script.path)

    # if there are any mods to apply, set those up now
    if len(mods):
        mv_commands = []
        for suffix, action in mods.items():
            for output_short_path in output_short_paths_by_basename.values():
                if output_short_path.endswith(suffix):
                    mv_commands.append("awk '{action}' {dir}/{short_path} > {dir}/{short_path}.tmp".format(
                        action = action,
                        dir = out_dir,
                        short_path = output_short_path,
                    ))
                    mv_commands.append("mv {dir}/{short_path}.tmp {dir}/{short_path}".format(
                        dir = out_dir,
                        short_path = output_short_path,
                    ))
        mv_script = ctx.actions.declare_file(ctx.label.name + "_mv.sh")
        ctx.actions.write(mv_script, "\n".join(mv_commands), is_executable = True)
        inputs.append(mv_script)
        commands.append(mv_script.path)

    # if the ctx.attr.output_file_suffix was set in conjunction with
    # ctx.attr.srcs, we want to rename all the output files to a different
    # suffix (e.g. foo.ts -> foo.ts.gen).  The relocates the files that were
    # generated by protoc plugins to a different name.  This is used by the
    # 'proto_compiled_sources' rule.  The reason is that if we also have a
    # `foo.ts` source file sitting in the workspace (checked into git), rules
    # like `ts_project` will perform a 'copy_to_bin' action on the file.  If we
    # didn't do this rename, the ts_project rule and the proto_compile rule
    # would attempt to create the same output file in bazel-bin (foo.ts),
    # causing an error.
    #
    # In the case of proto_compiled_sources, executing `bazel run
    # //proto:foo_ts.update` would generate the file
    # `bazel-bin/proto/foo.ts.gen` and the gencopy operation will copy that file
    # to `WORKSPACE/proto/foo.ts`, essentially making the `.gen` a
    # temporary-like file.
    if len(renames):
        rename_commands = []
        for src, dst in renames.items():
            rename_commands.append("mv {dir}/{src} {dir}/{dst}".format(
                dir = out_dir,
                src = src,
                dst = dst,
            ))
        rename_script = ctx.actions.declare_file(ctx.label.name + "_rename.sh")
        ctx.actions.write(rename_script, "\n".join(rename_commands), is_executable = True)
        inputs.append(rename_script)
        commands.append(rename_script.path)

    if verbose:
        before = ["env", "pwd", "ls -al .", "echo '\n##### SANDBOX BEFORE RUNNING PROTOC'", "find * -type l | grep -v node_modules"]
        after = ["echo '\n##### SANDBOX AFTER RUNNING PROTOC'", "find * -type f"]
        commands = before + commands + after

        for c in commands:
            # buildifier: disable=print
            print("COMMAND:", c)
        for a in replaced_args:
            # buildifier: disable=print
            print("ARG:", a)
        for f in protos:
            # buildifier: disable=print
            print("PROTO:", f.path)
        for f in inputs:
            # buildifier: disable=print
            print("INPUT:", f.path)
        for f in output_files_by_rel_path.values():
            # buildifier: disable=print
            print("EXPECTED OUTPUT:", f.path)

    ctx.actions.run_shell(
        arguments = [final_args],
        command = "\n".join(commands),
        inputs = inputs,
        mnemonic = "Protoc",
        outputs = output_files_by_rel_path.values(),
        progress_message = "Compiling protoc outputs for %r" % [f.basename for f in protos],
        tools = tools,
        input_manifests = input_manifests,
        env = {"BAZEL_BINDIR": ctx.bin_dir.path},
    )

    outputs = output_files_by_rel_path.values()

    providers = [
        ProtoCompileInfo(
            label = ctx.label,
            outputs = outputs,
            output_files_by_rel_path = output_files_by_rel_path,
        ),
    ]
    if ctx.attr.default_info:
        providers.append(DefaultInfo(files = depset(outputs)))

    return providers

proto_compile = rule(
    implementation = _proto_compile_impl,
    attrs = {
        "args": attr.string_list(
            doc = "List of additional protoc args",
        ),
        "outputs": attr.output_list(
            doc = "List of source files we expect to be generated (relative to package)",
        ),
        "srcs": attr.string_list(
            doc = "List of source files we expect to be regenerated (relative to package)",
        ),
        "plugins": attr.label_list(
            doc = "List of ProtoPluginInfo providers",
            mandatory = True,
            providers = [ProtoPluginInfo],
        ),
        "options": attr.string_list_dict(
            doc = "List of additional options, keyed by proto_plugin label",
        ),
        "outs": attr.string_dict(
            doc = "Output location, keyed by proto_plugin label",
        ),
        "output_mappings": attr.string_list(
            doc = "strings of the form A=B where A is a file named in attr.outputs and B is the actual file generated in the execroot",
        ),
        "proto": attr.label(
            doc = "The single ProtoInfo provider",
            mandatory = True,
            providers = [ProtoInfo],
        ),
        "protoc": attr.label(
            doc = "Overrides the protoc from the toolchain",
            allow_single_file = True,
            executable = True,
            cfg = "exec",
        ),
        "verbose": attr.bool(
            doc = "The verbosity flag.",
            default = False,
        ),
        "default_info": attr.bool(
            doc = "If false, do not return the DefaultInfo provider",
            default = True,
        ),
        "output_file_suffix": attr.string(
            doc = "If set, copy the output files to a new set having this suffix",
        ),
    },
    toolchains = ["@build_stack_rules_proto//toolchain:protoc"],
)
