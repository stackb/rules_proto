load("//:plugin.bzl", "ProtoPluginInfo")
load(
    "//:common.bzl",
    "ProtoCompileInfo",
    "apply_plugin_transitivity_rules",
    "copy_jar_to_srcjar",
    "copy_proto",
    "get_output_filename",
    "get_output_sibling_file",
    "get_plugin_options",
    "get_plugin_out",
    "get_plugin_runfiles",
)


def get_plugin_out_arg(ctx, outdir, plugin, plugin_outfiles):
    """Build the --<plugin>_out argument

    Args:
      ctx: the <ctx> object
      outdir: the package output directory <string>
      plugin: the <PluginInfo> object.
      plugin_outfiles: The <dict<string,<File>>.  For example, {closure: "library.js"}

    Returns
      <string> for the protoc arg list.
    """

    arg = outdir
    if plugin.outdir:
        arg = plugin.outdir.replace("{name}", outdir)
    elif plugin.out:
        outfile = plugin_outfiles[plugin.name]
        arg = outfile.path

    # Collate a list of options from the plugin itself PLUS options from the
    # global plugin_options list (if they exist)
    options = []
    options += getattr(plugin, "options", [])
    options += getattr(ctx.attr, "plugin_options", [])

    if options:
        arg = "%s:%s" % (",".join(get_plugin_options(ctx.label.name, options)), arg)
    return "--%s_out=%s" % (plugin.name, arg)


def get_plugin_outputs(ctx, descriptor, outputs, src, proto, plugin):
    """Get the predicted generated outputs for a given plugin

    Args:
      ctx: the <ctx> object
      descriptor: the descriptor <Generated File>
      outputs: the list of outputs.
      src: the orginal .proto source <Source File>.
      proto: the copied .proto <Generated File> (the one in the package 'staging area')
      plugin: the <PluginInfo> object.

    Returns:
      <list<Generated File>> the augmented list of files that will be generated
    """
    for output in plugin.outputs:
        filename = get_output_filename(src, plugin, output)
        if not filename:
            continue
        sibling = get_output_sibling_file(output, proto, descriptor)
        outputs.append(ctx.actions.declare_file(filename, sibling = sibling))
    return outputs


def proto_compile_impl(ctx):
    ###
    ### Part 1: setup variables used in scope
    ###

    # <struct> The resolved protoc toolchain
    protoc_toolchain_info = ctx.toolchains["@build_stack_rules_proto//protobuf:toolchain_type"]

    # <Target> The resolved protoc compiler from the protoc toolchain
    protoc = protoc_toolchain_info.protoc

    # <int> verbose level
    verbose = ctx.attr.verbose

    # <File> for the output descriptor.  Often used as the sibling in
    # 'declare_file' actions.
    descriptor = ctx.outputs.descriptor

    # <string> The directory where that generated descriptor is.
    outdir = descriptor.dirname

    # <list<ProtoInfo>> A list of ProtoInfo
    deps = [dep[ProtoInfo] for dep in ctx.attr.deps]

    # <list<PluginInfo>> A list of PluginInfo
    plugins = [plugin[ProtoPluginInfo] for plugin in ctx.attr.plugins]

    # <list<File>> The list of .proto files that will exist in the 'staging
    # area'.  We copy them from their source location into place such that a
    # single '-I.' at the package root will satisfy all import paths.
    protos = []

    # <dict<string,File>> The set of .proto files to compile, used as the final
    # list of arguments to protoc.  This is a subset of the 'protos' list that
    # are directly specified in the proto_library deps, but excluding other
    # transitive .protos.  For example, even though we might transitively depend
    # on 'google/protobuf/any.proto', we don't necessarily want to actually
    # generate artifacts for it when compiling 'foo.proto'. Maintained as a dict
    # for set semantics.  The key is the value from File.path.
    targets = {}

    # <dict<string,File>> A mapping from plugin name to the plugin tool. Used to
    # generate the --plugin=protoc-gen-KEY=VALUE args
    plugin_tools = {}

    # <dict<string,<File> A mapping from PluginInfo.name to File.  In the case
    # of plugins that specify a single output 'archive' (like java), we gather
    # them in this dict.  It is used to generate args like
    # '--java_out=libjava.jar'.
    plugin_outfiles = {}

    # <list<File>> The list of srcjars that we're generating (like
    # 'foo.srcjar').
    srcjars = []

    # <list<File>> The list of generated artifacts like 'foo_pb2.py' that we
    # expect to be produced.
    outputs = []

    # Additional data files from plugin.data needed by plugin tools that are not
    # single binaries.
    data = []

    ###
    ### Part 2: gather plugin.out artifacts
    ###

    # Some protoc plugins generate a set of output files (like python) while
    # others generate a single 'archive' file that contains the individual
    # outputs (like java).  This first loop is for the latter type.  In this
    # scenario, the PluginInfo.out attribute will exist; the predicted file
    # output location is relative to the package root, marked by the descriptor
    # file. Jar outputs are gathered as a special case as we need to
    # post-process them to have a 'srcjar' extension (java_library rules don't
    # accept source jars with a 'jar' extension)
    for plugin in plugins:
        if plugin.executable:
            plugin_tools[plugin.name] = plugin.executable
        data += plugin.data + get_plugin_runfiles(plugin.tool)

        filename = get_plugin_out(ctx.label.name, plugin)
        if not filename:
            continue
        out = ctx.actions.declare_file(filename, sibling = descriptor)
        outputs.append(out)
        plugin_outfiles[plugin.name] = out
        if out.path.endswith(".jar"):
            srcjar = copy_jar_to_srcjar(ctx, out)
            srcjars.append(srcjar)

    ###
    ### Part 3a: Gather generated artifacts for each dependency .proto source file.
    ###

    for dep in deps:
        # Iterate all the directly specified .proto files.  If we have already
        # processed this one, skip it to avoid declaring duplicate outputs.
        # Create an action to copy the proto into our staging area.  Consult the
        # plugin to assemble the actual list of predicted generated artifacts
        # and save these in the 'outputs' list.
        for src in dep.direct_sources:
            if targets.get(src.path):
                continue
            proto = copy_proto(ctx, descriptor, src)
            targets[src] = proto
            protos.append(proto)

        # Iterate all transitive .proto files.  If we already processed in the
        # loop above, skip it. Otherwise add a copy action to get it into the
        # 'staging area'
        for src in dep.transitive_sources.to_list():
            if targets.get(src):
                continue
            if verbose > 2:
                print("transitive source: %r" % src)
            proto = copy_proto(ctx, descriptor, src)
            protos.append(proto)
            if ctx.attr.transitive:
                targets[src] = proto

    ###
    ### Part 3b: apply transitivity rules
    ###

    # If the 'transitive = true' was enabled, we collected all the protos into
    # the 'targets' list.
    # At this point we want to post-process that list and remove any protos that
    # might be incompatible with the plugin transitivity rules.
    if ctx.attr.transitive:
        for plugin in plugins:
            targets = apply_plugin_transitivity_rules(ctx, targets, plugin)

    ###
    ### Part 3c: collect generated artifacts for all in the target list of protos to compile
    ###
    for src, proto in targets.items():
        for plugin in plugins:
            outputs = get_plugin_outputs(ctx, descriptor, outputs, src, proto, plugin)

    ###
    ### Part 4: build list of arguments for protoc
    ###

    args = ["--descriptor_set_out=%s" % descriptor.path]

    # By default we have a single 'proto_path' argument at the 'staging area'
    # root.
    args += ["--proto_path=%s" % outdir]

    if ctx.attr.include_imports:
        args += ["--include_imports"]

    if ctx.attr.include_source_info:
        args += ["--include_source_info"]

    for plugin in plugins:
        args += [get_plugin_out_arg(ctx, outdir, plugin, plugin_outfiles)]

    args += ["--plugin=protoc-gen-%s=%s" % (k, v.path) for k, v in plugin_tools.items()]
    args += [proto.path for proto in targets.values()]

    ###
    ### Part 5: build the final protoc command and declare the action
    ###

    mnemonic = "ProtoCompile"
    command = " ".join([protoc.path] + args)

    if verbose > 0:
        print("%s: %s" % (mnemonic, command))
    if verbose > 1:
        command += " && echo '\n##### SANDBOX AFTER RUNNING PROTOC' && find ."
    if verbose > 2:
        command = "echo '\n##### SANDBOX BEFORE RUNNING PROTOC' && find . && " + command
    if verbose > 3:
        command = "env && " + command
        for f in outputs:
            print("expected output: %q", f.path)

    ctx.actions.run_shell(
        mnemonic = mnemonic,
        command = command,
        inputs = protos + data,
        outputs = outputs + [descriptor] + ctx.outputs.outputs,
        tools = [protoc] + plugin_tools.values(),
    )

    ###
    ### Part 6: assemble output providers
    ###

    # The files for 'DefaultInfo' include any explicit outputs for the rule.  If
    # we are generating srcjars, use those as the final outputs rather than
    # their '.jar' intermediates.  Otherwise include all the file outputs.
    # NOTE: this looks a little wonky here.  It probably works in simple cases
    # where there list of plugins has length 1 OR all outputting to jars OR all
    # not outputting to jars.  Probably would break here if they were mixed.
    files = [] + ctx.outputs.outputs

    if len(srcjars) > 0:
        files += srcjars
    else:
        files += outputs
        if len(plugin_outfiles) > 0:
            files += plugin_outfiles.values()

    return [
        ProtoCompileInfo(
            label = ctx.label,
            plugins = plugins,
            protos = protos,
            outputs = outputs,
            files = files,
            tools = plugin_tools,
            args = args,
            descriptor = descriptor,
        ),
        DefaultInfo(files = depset(files)),
    ]


proto_compile = rule(
    implementation = proto_compile_impl,
    attrs = {
        "deps": attr.label_list(
            doc = "List of labels that provide a `ProtoInfo` (such as `native.proto_library`)",
            providers = [ProtoInfo],
            mandatory = True,
        ),
        "plugins": attr.label_list(
            doc = "List of labels that provide a `ProtoPluginInfo` used as plugins to protoc",
            providers = [ProtoPluginInfo],
            mandatory = True,
        ),
        "plugin_options": attr.string_list(
            doc = "List of additional 'global' plugin options (applies to all plugins). To apply plugin specific options, use the `options` attribute on `proto_plugin`",
        ),
        "outputs": attr.output_list(
            doc = "List of additional expected generated file outputs",
        ),
        "verbose": attr.int(
            doc = "The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*",
            default = 0,
        ),
        "include_imports": attr.bool(
            doc = "Pass the --include_imports argument to the protoc_plugin",
            default = True,
        ),
        "include_source_info": attr.bool(
            doc = "Pass the --include_source_info argument to the protoc_plugin",
            default = True,
        ),
        "transitive": attr.bool(
            doc = "Generate outputs for both *.proto directly named in `deps` AND all their transitive proto_library dependencies",
            default = True,
        ),
        "transitivity": attr.string_dict(
            doc = "Transitive filters to apply when the 'transitive' property is enabled. This string_dict can be used to exclude or explicitly include protos from the compilation list by using `exclude` or `include` respectively as the dict value",
            default = {},
        ),
    },
    outputs = {
        "descriptor": "%{name}/descriptor.source.bin",
    },
    output_to_genfiles = True,
    toolchains = ["@build_stack_rules_proto//protobuf:toolchain_type"],
)
