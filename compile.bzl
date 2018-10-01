load("//:plugin.bzl", "ProtoPluginInfo")

ProtoCompileInfo = provider(fields = {
    "label": "label object",
    "plugins": "ProtoPluginInfo object",
    "descriptor": "descriptor set file",
    "outputs": "generated protoc outputs",
    "files": "final generated files",
    "protos": "generated protos (copies)",
    "args": "proto arguments",
    "tools": "proto tools",
    "verbose": "verbose level",
})


# Hack - providers indexing is by int, but I have not idea how to get the actual
# provider object here.
ProtoInfoProvider = 0

def _capitalize(s):
  """Capitalize a string - only first letter
  Args:
    s (string): The input string to be capitalized.
  Returns:
    (string): The capitalized string.
  """
  return s[0:1].upper() + s[1:]


def _pascal_case(s):
    """Convert pascal_case -> PascalCase
    Args:
        s (string): The input string to be capitalized.
    Returns:
        (string): The capitalized string.
    """
    return "".join([_capitalize(part) for part in s.split("_")])

def _get_output_sibling_file(pattern, proto, descriptor):
    if pattern.startswith("@package/"):
        return descriptor
    return proto

def _get_plugin_out(ctx, plugin):
    if not plugin.out:
        return None
    filename = plugin.out
    filename = filename.replace("%{name}", ctx.label.name)    
    return filename

def _copy_jar_to_srcjar(ctx, jar):
    srcjar = ctx.actions.declare_file("%s/%s.srcjar" % (ctx.label.name, ctx.label.name))
    ctx.actions.run_shell(
        mnemonic = "CopySrcjar",
        inputs = [jar],
        outputs = [srcjar],
        command = "cp %s %s" % (jar.path, srcjar.path),
    )
    return srcjar

def _get_output_filename(src, plugin, pattern):
    # If output to srcjar, don't emit a per-proto output file.
    if plugin.out:
        return None
    # Slice off this prefix if it exists, we don't use it here.
    if pattern.startswith("@package/"):
        pattern = pattern[len("@package/"):]
    basename = src.basename
    if basename.endswith(".proto"):
        basename = basename[:-6]
    elif basename.endswith(".protodevel"):
        basename = basename[:-11]

    filename = basename
   
    if pattern.find("{basename}") != -1:
        filename = pattern.replace("{basename}", basename)
    elif pattern.find("{basename|pascal}") != -1:
        filename = pattern.replace("{basename|pascal}", _pascal_case(basename))
    else:
        filename = basename + pattern

    return filename

def _get_proto_filename(src):
    parts = src.short_path.split("/")
    if len(parts) > 1 and parts[0] == "..":
        return "/".join(parts[2:])
    return src.short_path

def copy_proto(ctx, descriptor, src):
    proto = ctx.actions.declare_file(_get_proto_filename(src), sibling = descriptor)
    ctx.actions.run_shell(
        mnemonic = "CopyProto",
        inputs = [src],
        outputs = [proto],
        command = "cp %s %s" % (src.path, proto.path),
    )
    return proto

def _get_plugin_option(ctx, option):
    return option.replace("%{name}", ctx.label.name)

def _get_plugin_options(ctx, options):
    return [_get_plugin_option(ctx, option) for option in options]

def get_plugin_out_arg(ctx, outdir, plugin, plugin_outfiles):
    arg = outdir
    if plugin.outdir:
        arg = plugin.outdir.replace("%{name}", outdir)
    elif plugin.out:
        outfile = plugin_outfiles[plugin.name]
        #arg = "%s" % (outdir)
        #arg = "%s/%s" % (outdir, outfile.short_path)
        arg = outfile.path
    if plugin.options:
        arg = "%s:%s" % (",".join(_get_plugin_options(ctx, plugin.options)), arg) 
    return "--%s_out=%s" % (plugin.name, arg)  

def _get_plugin_outputs(ctx, descriptor, outputs, src, proto, plugin):
    for output in plugin.outputs:
        filename = _get_output_filename(src, plugin, output)
        if not filename:
            continue
        sibling = _get_output_sibling_file(output, proto, descriptor)
        outputs.append(ctx.actions.declare_file(filename, sibling = sibling))
    return outputs

proto_compile_attrs = {
    "deps": attr.label_list(
        mandatory = True,
        providers = ["proto"],
    ),
    "plugins": attr.label_list(
        providers = [ProtoPluginInfo],
        mandatory = True,
    ),
    "outputs": attr.output_list(
    ),
    "has_services": attr.bool(
        doc = "If the proto files(s) have a service rpc, generate grpc outputs",
    ),
    "protoc": attr.label(
        default = "@com_google_protobuf//:protoc",
        cfg = "host",
        executable = True,
    ),
    "verbose": attr.int(
        doc = "Increase verbose level for more debugging",
    ),
    "include_imports": attr.bool(
        doc = "Pass the --include_imports argument to the protoc_plugin",
        default = True,
    ),
    "include_source_info": attr.bool(
        doc = "Pass the --include_source_info argument to the protoc_plugin",
        default = True,
    ),
}

proto_compile_outputs = {
    "descriptor": "%{name}/descriptor.source.bin",
}

def proto_compile_impl(ctx):
    verbose = ctx.attr.verbose
    protoc = ctx.executable.protoc
    has_services = ctx.attr.has_services
    descriptor = ctx.outputs.descriptor
    outdir = descriptor.dirname
    deps = [dep.proto for dep in ctx.attr.deps]
    #datadeps = [dep[DefaultInfo].data_runfiles for dep in ctx.attr.deps]
    #datadeps = [dep[DefaultInfo].default_runfiles for dep in ctx.attr.deps]
    # datadeps = [dep[ProtoSupportDataInfo].default_runfiles for dep in ctx.attr.deps]
    # print("datadeps: %r" % datadeps)
    # for dat in datadeps:
    #     print("dat: %r" % dat)
    #     # files, symlinks, empty_filenames
    #     for f in dat.symlinks:
    #         print("datfile: %r" % f)

    if verbose:
        print("Starting proto compile...")

    plugins = [plugin[ProtoPluginInfo] for plugin in ctx.attr.plugins]
    tools = {}
    protos = []
    outputs = []
    args = []
    directs = {}
    srcjars = []
    plugin_outfiles = {}
    # Aggregate files from plugin.data.  For example for some reason the
    # dart_plugin executable does not pull in the dart_sdk dart binary.
    data = []

    for dep in deps:
        for plugin in plugins:
            filename = _get_plugin_out(ctx, plugin)
            if not filename:
                continue
            out = ctx.actions.declare_file(filename, sibling = descriptor)
            outputs.append(out)
            plugin_outfiles[plugin.name] = out
            # Special handling for output jars
            if out.path.endswith(".jar"):
                srcjar = _copy_jar_to_srcjar(ctx, out)
                srcjars.append(srcjar)

        for src in dep.direct_sources:
            if directs.get(src.path):
                continue
            directs[src.path] = src
            proto = copy_proto(ctx, descriptor, src)
            protos.append(proto)

            for plugin in plugins:
                outputs = _get_plugin_outputs(ctx, descriptor, outputs, src, proto, plugin)

        for src in dep.transitive_sources:
            if directs.get(src.path):
                continue
            if verbose > 2:
                print("transitive source: %r" % src)
            proto = copy_proto(ctx, descriptor, src)
            protos.append(proto)
        for e in dep.transitive_proto_path:
            if verbose > 2:
                print("proto_path: %r" % e)
        for e in dep.transitive_descriptor_sets:
            if verbose > 2:
                print("descriptor_set: %r" % e)

    args += ["--descriptor_set_out=%s" % descriptor.path]
    if ctx.attr.include_imports:
        args += ["--include_imports"]
    if ctx.attr.include_source_info:
        args += ["--include_source_info"]

    args += ["--proto_path=%s" % outdir]        
    for plugin in plugins:
        data += plugin.data
        args += [get_plugin_out_arg(ctx, outdir, plugin, plugin_outfiles)]        
        if plugin.executable:    
            tools[plugin.name] = plugin.executable

    args += ["--plugin=protoc-gen-%s=%s" % (k, v.path) for k, v in tools.items()]        
    args += [proto.path for proto in directs.values()]

    command = " ".join([protoc.path] + args)
    if verbose > 0:
        print("PROTOC COMMAND: %s" % command)
    if verbose > 1:
        command += "&& echo '\n##### SANDBOX AFTER RUNNING PROTOC' && find ."
    if verbose > 2:
        command = "echo '\n##### SANDBOX BEFORE RUNNING PROTOC' && find . && " + command
    if verbose > 3:
        command = "env && " + command

    for plugin in plugins:
        data += get_tool_files(plugin.tool)

    ctx.actions.run_shell(
        command = command,
        inputs = [protoc] + tools.values() + protos + data,
        outputs = outputs + [descriptor] + ctx.outputs.outputs,
    )

    files = [] + ctx.outputs.outputs

    if len(srcjars) > 0:
        files += srcjars
    else:
        files += outputs
        if len(plugin_outfiles) > 0:
            files += plugin_outfiles.values()

    return [ProtoCompileInfo(
        label = ctx.label,
        plugins = plugins,
        protos = protos,
        outputs = outputs,
        files = files,
        tools = tools,
        args = args,
        descriptor = descriptor,
    ), DefaultInfo(files = depset(files))]


def get_tool_files(tool):
    files = []
    if not tool:
        return files
    info = tool[DefaultInfo]
    if not info:
        return files
    if info.files:
        files += info.files.to_list()
    if info.default_runfiles:
        runfiles = info.default_runfiles
        if runfiles.files:
            files += runfiles.files.to_list()
    return files


proto_compile = rule(
    implementation = proto_compile_impl,
    attrs = proto_compile_attrs,
    outputs = proto_compile_outputs,
    output_to_genfiles = True,
)


def invoke(proto_compile_rule, name_suffix, kwargs):
    """Invoke a proto_compile rule using kwargs

    Invoke is a convenience function for library rules that call proto_compile
    rules.  Rather than having to do the same boilerplate across many different
    files, this function centralizes the logic of calling proto_compile rules
    using kwargs.

    Args:
      proto_compile_rule: the rule function to invoke
      name_suffix: a suffix for the kwargs.name to use for the rule
      kwargs: the **kwargs dict, passed directly (not decontucted)

    Returns:
      The name of the invoked rule. This can be used in the srcs label of a library rule.
    """ 

    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")
    verbose = kwargs.get("verbose")
    rule_name = name + name_suffix

    proto_compile_rule(
        name = rule_name,
        deps = deps,
        visibility = visibility,
        verbose = verbose,
    )

    return rule_name

