load("@build_stack_rules_proto//rules:proto_plugin.bzl", "ProtoPluginInfo")
load("@rules_proto//proto:defs.bzl", "ProtoInfo")

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

def get_protoc_executable(ctx):
    protoc_toolchain_info = ctx.toolchains[str(Label("//protoc:toolchain_type"))]
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


def _proto_compile_impl(ctx):

    ###
    ### Part 1: setup variables used in scope
    ###

    # const <int> verbosity level
    verbose = ctx.attr.verbose
    # const <File> the protoc file from the toolchain
    protoc = get_protoc_executable(ctx)
    # const <ProtoInfo> proto provider
    proto_info = ctx.attr.proto[ProtoInfo]
    # const <list<ProtoPluginInfo>> plugins to be applied
    plugins = [plugin[ProtoPluginInfo] for plugin in ctx.attr.plugins]
    # const <dict<string,list<string>>> 
    plugin_options = ctx.attr.plugin_options
    # const <dict<string,string>> 
    plugin_outs = ctx.attr.plugin_out
    # const <list<File>> files we expect to be generated
    genfiles = ctx.outputs.generated_srcs
    # const <list<File>> outputs for the compile action
    outputs = [] + genfiles

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
            for srcs in info.direct_sources:
                protos += srcs.to_list()
            descriptors += info.transitive_descriptor_sets.to_list()

        # Include extra plugin data files
        inputs += plugin.data

        ### Part 2.2: build --plugin argument

        # const <string> The name of the plugin
        plugin_name = plugin.protoc_plugin_name if plugin.protoc_plugin_name else plugin.name
        # const <?File> Add plugin executable if not a built-in plugin
        plugin_tool = plugin.tool if plugin.tool else None

        # Add plugin runfiles if plugin has a tool
        if plugin_tool:
            tools.append(plugin_tool)
            # const <depset<File>, <list<opaque>>
            plugin_runfiles, plugin_input_manifests = ctx.resolve_tools(tools = [plugin.tool_target])
            if plugin_input_manifests:
                input_manifests.append(plugin_input_manifests) # TODO: check this
            inputs += plugin_runfiles.to_list()
            # If Windows, mangle the path. It's done a bit awkwardly with
            # `host_path_seprator` as there is no simple way to figure out what's
            # the current OS.
            plugin_tool_path = plugin_tool.path
            if ctx.configuration.host_path_separator == ";":
                plugin_tool_path = plugin_tool.path.replace("/", "\\")

            args.append("--plugin=protoc-gen-{}={}".format(plugin_name, plugin_tool_path))

        ### Part 2.3: build --{name}_out=OPTIONS argument

        options = [opt for opt in plugin_options.get(str(plugin.label), [])]
        options += plugin.options

        # mut <string>
        out = plugin.out
        if len(options) > 0:
            if plugin.separate_options_flag:
                args.append("--{}_opt={}".format(plugin_name, ",".join(options)))
            else:
                out = "{}:{}".format(",".join(options), out)
        # override with the out configured on the rule if specified
        out = plugin_outs.get(str(plugin.label), out)

        args.append("--{}_out={}".format(plugin_name, out))

    ###
    ### Part 3: trailing args
    ###

    ### Part 3.1: add descriptor sets

    descriptors = _uniq(descriptors)
    inputs += descriptors

    args.append("--descriptor_set_in={}".format(ctx.configuration.host_path_separator.join(
        [d.path for d in descriptors],
    )))

    ### Part 3.3: arg symbol replacements

    args = [arg.replace("{BIN_DIR}", ctx.bin_dir.path) for arg in args]
    args = [arg.replace("{PACKAGE}", ctx.label.package) for arg in args]
    args = [arg.replace("{NAME}", ctx.label.name) for arg in args]

    ### Part 3.4: add proto file args

    protos = _uniq(protos)
    for proto in protos:
        args.append(_descriptor_proto_path(proto, proto_info))

    ### Step 3.5: build args object

    final_args = ctx.actions.args()
    final_args.add_all(args)

    ###
    ### Step 4: command action
    ###

    command = "mkdir -p {} && {} $@".format(ctx.label.package, protoc.path) # $@ is replaced with args list
    if verbose > 2:
        command += " && echo '\n##### SANDBOX AFTER RUNNING PROTOC' && find . -type f "
        command = "env && echo '\n##### SANDBOX BEFORE RUNNING PROTOC' && find . -type l && " + command

    ### Step 3.6: declare action

    ctx.actions.run_shell(
        arguments = [final_args],
        command = command,
        inputs = inputs,
        # input_manifests = input_manifests, TODO
        mnemonic = "Protoc",
        outputs = outputs,
        progress_message = "Compiling protoc outputs for %r" % [f.basename for f in protos],
        tools = tools,
    )

    ### Step 3.7: debugging

    if verbose > 1:
        for f in inputs:
            print("INPUT:", f.path)
        for f in tools:
            print("TOOL:", f.path)
        for f in protos:
            print("PROTO:", f.path)
        for f in outputs:
            print("EXPECTED OUTPUT:", f.path)
        for a in args:
            print("ARG:", a)

    return [DefaultInfo(files = depset(genfiles))]


proto_compile = rule(
    implementation = _proto_compile_impl,
    attrs = {
        "args": attr.string_list(
            doc = "List of additional protoc args",
        ),
        "generated_srcs": attr.output_list(
            doc = "List of source files we expect to be generated (relative to package)",
            mandatory = True,
        ),
        "plugins": attr.label_list(
            doc = "List of ProtoPluginInfo providers",
            mandatory = True,
            providers = [ProtoPluginInfo],
        ),
        "plugin_options": attr.string_list_dict(
            doc = "List of additional options, keyed by proto_plugin label",
        ),
        "plugin_out": attr.string_dict(
            doc = "Output location, keyed by proto_plugin label",
        ),
        "proto": attr.label(
            doc = "The single ProtoInfo provider",
            mandatory = True,
            providers = [ProtoInfo],
        ),
        "verbose": attr.int(
            doc = "The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox before+after running protoc*",
        ),
    },
    toolchains = ["@build_stack_rules_proto//protoc:toolchain_type"],
)