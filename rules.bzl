ProtoPluginInfo = provider(fields = {
    "name": "proto plugin name",
    "outputs": "outputs to be generated",
    "tool": "plugin tool",
    "options": "proto options",
    "out": "aggregate proto output",
})

ProtoCompileInfo = provider(fields = {
    "plugins": "ProtoPluginInfo object",
    "descriptor": "descriptor set file",
    "outputs": "generated protoc outputs",
    "files": "final generated files",
    "protos": "generated protos (copies)",
    "args": "proto arguments",
    "tools": "proto tools",
    "verbose": "verbose level",
})


def _proto_plugin_impl(ctx):
    return [ProtoPluginInfo(
        name = ctx.label.name,
        options = ctx.attr.options,
        out = ctx.attr.out,
        outputs = ctx.attr.outputs,
        tool = ctx.executable.tool,
    )]

proto_plugin = rule(
    implementation = _proto_plugin_impl,
    attrs = {
        "options": attr.string_list(
            doc = "An list of options to pass to the compiler.",
        ),
        "outputs": attr.string_list(
            doc = "Output filenames generated on a per-proto basis.  Example: '{basename}_pb2.py'",
        ),
        "out": attr.string(
            doc = "Output filename generated on a per-plugin basis; to be used in the value for --NAME-out=OUT",
        ),
        "tool": attr.label(
            doc = "The plugin binary.  If absent, assume the plugin is a built-in to protoc itself",
            cfg = "host",
            executable = True,
        ),
    }
)


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
    if plugin.out:
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
}

proto_compile_outputs = {
    "descriptor": "%{name}/descriptor.bin",
}

def proto_compile_impl(ctx):
    protoc = ctx.executable.protoc
    has_services = ctx.attr.has_services
    descriptor = ctx.outputs.descriptor
    outdir = descriptor.dirname
    deps = [dep.proto for dep in ctx.attr.deps]
    plugins = [plugin[ProtoPluginInfo] for plugin in ctx.attr.plugins]
    tools = {}
    protos = []
    outputs = []
    args = []
    directs = {}
    srcjars = []
    plugin_outfiles = {}

    for dep in deps:
        for plugin in plugins:
            print("checking plugin %s" % plugin.name)
            filename = _get_plugin_out(ctx, plugin)
            if not filename:
                continue
            out = ctx.actions.declare_file(filename, sibling = descriptor)
            outputs.append(out)
            plugin_outfiles[plugin.name] = out
            print("registered proto out %s: %s" % (plugin.name, out.path))
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
            #print("transitive source: %r" % src)
            proto = copy_proto(ctx, descriptor, src)
            protos.append(proto)
        for e in dep.transitive_proto_path:
            print("proto_path: %r" % e)
        # for e in dep.transitive_descriptor_sets:
        #     print("descriptor_set: %r" % e)

    args += ["--descriptor_set_out=%s" % descriptor.path]
    args += ["--proto_path=%s" % outdir]        
    for plugin in plugins:
        args += [get_plugin_out_arg(ctx, outdir, plugin, plugin_outfiles)]        
        if plugin.tool:    
            tools[plugin.name] = plugin.tool

    args += ["--plugin=protoc-gen-%s=%s" % (k, v.path) for k, v in tools.items()]        
    args += [proto.path for proto in protos]

    command = " ".join([protoc.path] + args)
    if ctx.attr.verbose > 0:
        print("PROTOC COMMAND: %s" % command)
    if ctx.attr.verbose > 1:
        command += "&& echo '\n##### SANDBOX AFTER RUNNING PROTOC' && find ."
    if ctx.attr.verbose > 2:
        command = "echo '\n##### SANDBOX BEFORE RUNNING PROTOC' && find . && " + command

    ctx.actions.run_shell(
        command = command,
        inputs = [protoc] + tools.values() + protos,
        outputs = outputs + [descriptor] + ctx.outputs.outputs,
    )

    files = [] + ctx.outputs.outputs
    files += outputs
    if len(plugin_outfiles) > 0:
        files += plugin_outfiles.values()

    return [ProtoCompileInfo(
        plugins = plugins,
        protos = protos,
        outputs = outputs,
        files = files,
        tools = tools,
        args = args,
        descriptor = descriptor,
    ), DefaultInfo(files = depset(files))]

proto_compile = rule(
    implementation = proto_compile_impl,
    attrs = proto_compile_attrs,
    outputs = proto_compile_outputs,
    output_to_genfiles = True,
)