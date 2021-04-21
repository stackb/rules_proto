load("@build_stack_rules_proto//rules:proto_plugin.bzl", "ProtoPluginInfo")
load("@rules_proto//proto:defs.bzl", "ProtoInfo")
load("proto_utils.bzl", "get_protoc_executable", "uniq", "descriptor_proto_path")

def _protoc_impl(ctx):

    ###
    ### Part 1: setup variables used in scope
    ###

    # <ProtoInfo> proto provider
    proto_info = ctx.attr.proto[ProtoInfo]
    # <list<ProtoPluginInfo>> plugins to be applied
    plugins = [plugin[ProtoPluginInfo] for plugin in ctx.attr.plugins]
    # <File> the protoc file
    protoc = get_protoc_executable(ctx)
    # <list<File>> tools for the compile action
    tools = [protoc]
    # <File> the descriptor_set_out file
    out = ctx.actions.declare_file(ctx.label.name + "_proto-descriptor-set.proto.bin")
    # <list<File>> files we expect to be generated
    genfiles = [ctx.actions.declare_file(rel, sibling = out) for rel in ctx.attr.generated_srcs]
    # <list<File>> outputs for the compile action
    outputs = [out] + genfiles
    # <list<File>> set of descriptors for the compile action
    descriptors = proto_info.transitive_descriptor_sets.to_list()
    # <list<File>> inputs for the compile action
    inputs = []
    # <list<string>> argument list for protoc execution
    args = []
    # <list<File>> The (filtered) set of .proto files to compile
    protos = []
    # <list<opaque>> Plugin input manifests
    input_manifests = []
    # <int> verbosity level
    verbose = ctx.attr.verbose

    ###
    ### Part 2: iterate over plugins
    ###
    for plugin in plugins:
        # Add in additional proto deps if specified by the compiler
        for info in plugin.supplementary_proto_deps:
            descriptors += info.transitive_descriptor_sets.to_list()
        # Add extra plugin data files
        inputs += plugin.data

        ###
        ### Part 2.1: build --plugin argument
        ###

        # Get plugin name
        plugin_name = plugin.protoc_plugin_name if plugin.protoc_plugin_name else plugin.name

        # Add plugin executable if not a built-in plugin
        plugin_tool = None
        if plugin.tool_executable:
            plugin_tool = plugin.tool_executable

        # Add plugin runfiles if plugin has a tool
        if plugin_tool:
            tools.append(plugin_tool)  
            plugin_runfiles, plugin_input_manifests = ctx.resolve_tools(tools = [plugin.tool])
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

        ###
        ### Part 2.2: build --{name}_out=OPTIONS argument
        ###
        out_arg = out.dirname
        if plugin.options:
            opts_str = ",".join(
                [option.replace("{name}", ctx.label.name) for option in plugin.options],
            )
            if plugin.separate_options_flag:
                args.append("--{}_opt={}".format(plugin_name, opts_str))
            else:
                out_arg = "{}:{}".format(opts_str, out_arg)
        args.append("--{}_out={}".format(plugin_name, out_arg))

        ###
        ### Part 2.3: build protos list
        ###
        for proto in proto_info.direct_sources:
            # Check for exclusion
            if any([
                proto.dirname.endswith(exclusion) or proto.path.endswith(exclusion)
                for exclusion in plugin.exclusions
            ]) or proto in protos:  # TODO: When using import_prefix, the ProtoInfo.direct_sources list appears to contain duplicate records, this line removes these. https://github.com/bazelbuild/bazel/issues/9127
                continue

            # Proto not excluded
            protos.append(proto)

        # Add in extra proto deps if attached to the plugin
        for info in plugin.supplementary_proto_deps:
            for srcs in info.direct_sources:
                protos += srcs.to_list()

    ###
    ### Part 3.1: build --descriptor_set_in args
    ###
    descriptor_set_list = uniq(descriptors)
    for d in descriptor_set_list:
        inputs.append(d)
    
    pathsep = ctx.configuration.host_path_separator
    args.append("--descriptor_set_in={}".format(pathsep.join(
        [d.path for d in descriptor_set_list],
    )))

    args.append("--descriptor_set_out="+out.path)

    # Add source proto files as descriptor paths
    for proto in uniq(protos):
        args.append(descriptor_proto_path(proto, proto_info))

    ###
    ### Step 3: declare action
    ###
    mnemonic = "Protoc"
    command = "mkdir -p {} && {} $@".format(ctx.label.package, protoc.path) # $@ is replaced with args list
    # for f in ctx.outputs.generated_srcs:
    #     command += "&& cp {} {}".format(f.short_path, f.path)

    if verbose > 0:
        command += " && echo '\n##### SANDBOX AFTER RUNNING PROTOC' && find . -type f "
    if verbose > 0:
        command = "echo '\n##### SANDBOX BEFORE RUNNING PROTOC' && find . -type l && " + command
    if verbose > 0:
        command = "env && " + command
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
        print("FINAL ARGS:", args)

    ctx_args = ctx.actions.args()
    ctx_args.add_all(args)
        
    ctx.actions.run_shell(
        progress_message = "Compiling protoc outputs for %r" % [f.basename for f in protos],
        # executable = protoc,
        command = command,
        arguments = [ctx_args],
        inputs = inputs,
        tools = tools,
        outputs = outputs,
        # input_manifests = input_manifests,
    )

    ###
    ### Step 4: generate providers
    ###
    return [
        DefaultInfo(files = depset(genfiles)),
    ]


protoc = rule(
    implementation = _protoc_impl,
    attrs = {
        "verbose": attr.int(
            doc = "The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*",
        ),
        "package_path": attr.string(
            doc = "The package_path option value",
        ),
        "proto": attr.label(
            doc = "The single ProtoInfo provider",
            mandatory = True,
            providers = [ProtoInfo],
        ),
        "plugins": attr.label_list(
            doc = "List of ProtoPluginInfo providers",
            mandatory = True,
            providers = [ProtoPluginInfo],
        ),
        "generated_srcs": attr.string_list(
            doc = "List of source files we expect to be generated (relative to package)",
            mandatory = True,
        ),
    },
    toolchains = [str(Label("@build_stack_rules_proto//toolchains:protoc_toolchain_type"))],
)
