load("//:aspect.bzl", "ProtoLibraryAspectNodeInfo")
load("//:plugin.bzl", "ProtoPluginInfo")

ProtoCompileInfo = provider(fields = {
    
})

def proto_compile_impl(ctx):
    files = [] 

    for dep in ctx.attr.deps:
        aspect = dep[ProtoLibraryAspectNodeInfo]
        files += aspect.outputs

    # print("final files: %r" % files)
    return [DefaultInfo(files = depset(files))]

# def proto_compile_impl(ctx):

#     ###
#     ### Part 1: setup variables used in scope
#     ###

#     # <int> verbose level
#     verbose = ctx.attr.verbose

#     # <File> the protoc tool
#     protoc = ctx.executable.protoc

#     # <bool> if this is a gRPC compilation
#     has_services = ctx.attr.has_services

#     # <File> for the output descriptor.  Often used as the sibling in
#     # 'declare_file' actions.
#     descriptor = ctx.outputs.descriptor

#     # <string> The directory where that generated descriptor is.
#     outdir = descriptor.dirname

#     # <list<ProtoInfo>> A list of ProtoInfo
#     deps = [dep.proto for dep in ctx.attr.deps]

#     # <list<PluginInfo>> A list of PluginInfo 
#     plugins = [plugin[ProtoPluginInfo] for plugin in ctx.attr.plugins]

#     # <list<File>> The list of .proto files that will exist in the 'staging
#     # area'.  We copy them from their source location into place such that a
#     # single '-I.' at the package root will satisfy all import paths.
#     protos = []

#     # <dict<string,File>> The set of .proto files to compile, used as the final
#     # list of arguments to protoc.  This is a subset of the 'protos' list that
#     # are directly specified in the proto_library deps, but excluding other
#     # transitive .protos.  For example, even though we might transitively depend
#     # on 'google/protobuf/any.proto', we don't necessarily want to actually
#     # generate artifacts for it when compiling 'foo.proto'. Maintained as a dict
#     # for set semantics.  The key is the value from File.path.  
#     targets = {}

#     # <dict<string,File>> A mapping from plugin name to the plugin tool. Used to
#     # generate the --plugin=protoc-gen-KEY=VALUE args
#     plugin_tools = {}

#     # <dict<string,<File> A mapping from PluginInfo.name to File.  In the case
#     # of plugins that specify a single output 'archive' (like java), we gather
#     # them in this dict.  It is used to generate args like
#     # '--java_out=libjava.jar'.  
#     plugin_outfiles = {}

#     # <list<File>> The list of srcjars that we're generating (like
#     # 'foo.srcjar').
#     srcjars = []

#     # <list<File>> The list of generated artifacts like 'foo_pb2.py' that we
#     # expect to be produced.
#     outputs = []
    
#     # Additional data files from plugin.data needed by plugin tools that are not
#     # single binaries. 
#     data = []

#     ###
#     ### Part 2: gather plugin.out artifacts
#     ###

#     # Some protoc plugins generate a set of output files (like python) while
#     # others generate a single 'archive' file that contains the individual
#     # outputs (like java).  This first loop is for the latter type.  In this
#     # scenario, the PluginInfo.out attribute will exist; the predicted file
#     # output location is relative to the package root, marked by the descriptor
#     # file. Jar outputs are gathered as a special case as we need to
#     # post-process them to have a 'srcjar' extension (java_library rules don't
#     # accept source jars with a 'jar' extension)
#     for plugin in plugins:
#         if plugin.executable:    
#             plugin_tools[plugin.name] = plugin.executable
#         data += plugin.data + get_plugin_runfiles(plugin.tool)

#         filename = _get_plugin_out(ctx, plugin)
#         if not filename:
#             continue
#         out = ctx.actions.declare_file(filename, sibling = descriptor)
#         outputs.append(out)
#         plugin_outfiles[plugin.name] = out
#         if out.path.endswith(".jar"):
#             srcjar = _copy_jar_to_srcjar(ctx, out)
#             srcjars.append(srcjar)

#     ###
#     ### Part 3a: Gather generated artifacts for each dependency .proto source file.
#     ###

#     for dep in deps:

#         # Iterate all the directly specified .proto files.  If we have already
#         # processed this one, skip it to avoid declaring duplicate outputs.
#         # Create an action to copy the proto into our staging area.  Consult the
#         # plugin to assemble the actual list of predicted generated artifacts
#         # and save these in the 'outputs' list.  
#         for src in dep.direct_sources:
#             if targets.get(src.path):
#                 continue
#             proto = copy_proto(ctx, descriptor, src)
#             targets[src] = proto
#             protos.append(proto)

#         # Iterate all transitive .proto files.  If we already processed in the
#         # loop above, skip it. Otherwise add a copy action to get it into the
#         # 'staging area'
#         for src in dep.transitive_sources:
#             if targets.get(src):
#                 continue
#             if verbose > 2:
#                 print("transitive source: %r" % src)
#             proto = copy_proto(ctx, descriptor, src)
#             protos.append(proto)
#             if ctx.attr.transitive:
#                 targets[src] = proto


#     ###
#     ### Part 3cb: apply transitivity rules
#     ###

#     # If the 'transitive = true' was enabled, we collected all the protos into
#     # the 'targets' list.  
#     # At this point we want to post-process that list and remove any protos that
#     # might be incompatible with the plugin transitivity rules.
#     if ctx.attr.transitive:
#         for plugin in plugins:
#             targets = _apply_plugin_transitivity_rules(ctx, targets, plugin)

#     ###
#     ### Part 3c: collect generated artifacts for all in the target list of protos to compile
#     ###
#     for proto in targets.values():
#         for plugin in plugins:
#             outputs = _get_plugin_outputs(ctx, descriptor, outputs, proto, plugin)

#     ###
#     ### Part 4: build list of arguments for protoc
#     ###

#     args = ["--descriptor_set_out=%s" % descriptor.path]

#     # By default we have a single 'proto_path' argument at the 'staging area'
#     # root.
#     args += ["--proto_path=%s" % outdir]        

#     if ctx.attr.include_imports:
#         args += ["--include_imports"]

#     if ctx.attr.include_source_info:
#         args += ["--include_source_info"]

#     for plugin in plugins:
#         args += [get_plugin_out_arg(ctx, outdir, plugin, plugin_outfiles)]        

#     args += ["--plugin=protoc-gen-%s=%s" % (k, v.path) for k, v in plugin_tools.items()]

#     args += [proto.path for proto in targets.values()]

#     ###
#     ### Part 5: build the final protoc command and declare the action
#     ###

#     mnemonic = "ProtoCompile"

#     command = " ".join([protoc.path] + args)

#     if verbose > 0:
#         print("%s: %s" % (mnemonic, command))
#     if verbose > 1:
#         command += " && echo '\n##### SANDBOX AFTER RUNNING PROTOC' && find ."
#     if verbose > 2:
#         command = "echo '\n##### SANDBOX BEFORE RUNNING PROTOC' && find . && " + command
#     if verbose > 3:
#         command = "env && " + command
#         for f in outputs:
#             print("expected output: %q", f.path)    

#     ctx.actions.run_shell(
#         mnemonic = mnemonic,
#         command = command,
#         inputs = [protoc] + plugin_tools.values() + protos + data,
#         outputs = outputs + [descriptor] + ctx.outputs.outputs,
#     )

#     ###
#     ### Part 6: assemble output providers
#     ###

#     # The files for 'DefaultInfo' include any explicit outputs for the rule.  If
#     # we are generating srcjars, use those as the final outputs rather than
#     # their '.jar' intermediates.  Otherwise include all the file outputs.  
#     # NOTE: this looks a little wonky here.  It probably works in simple cases
#     # where there list of plugins has length 1 OR all outputting to jars OR all
#     # not outputting to jars.  Probably would break here if they were mixed.
#     files = [] + ctx.outputs.outputs

#     if len(srcjars) > 0:
#         files += srcjars
#     else:
#         files += outputs
#         if len(plugin_outfiles) > 0:
#             files += plugin_outfiles.values()

#     return [DefaultInfo(files = depset(files))]


proto_compile_attrs = {
    # "plugins": attr.label_list(
    #     doc = "List of protoc plugins to apply",
    #     providers = [ProtoPluginInfo],
    #     mandatory = True,
    # ),
    "plugin_options": attr.string_list(
        doc = "List of additional 'global' options to add (applies to all plugins)",
    ),
    "plugin_options_string": attr.string(
        doc = "(internal) List of additional 'global' options to add (applies to all plugins)",
    ),
    "outputs": attr.output_list(
        doc = "Escape mechanism to explicitly declare files that will be generated",
    ),
    "has_services": attr.bool(
        doc = "If the proto files(s) have a service rpc, generate grpc outputs",
    ),
    # "protoc": attr.label(
    #     doc = "The protoc tool",
    #     default = "@com_google_protobuf//:protoc",
    #     cfg = "host",
    #     executable = True,
    # ),
    "verbose": attr.int(
        doc = "Increase verbose level for more debugging",
    ),
    "verbose_string": attr.string(
        doc = "Increase verbose level for more debugging",
    ),
    # "include_imports": attr.bool(
    #     doc = "Pass the --include_imports argument to the protoc_plugin",
    #     default = True,
    # ),
    # "include_source_info": attr.bool(
    #     doc = "Pass the --include_source_info argument to the protoc_plugin",
    #     default = True,
    # ),
    "transitive": attr.bool(
        doc = "Emit transitive artifacts",
    ),
    "transitivity": attr.string_dict(
        doc = "Transitive rules.  When the 'transitive' property is enabled, this string_dict can be used to exclude protos from the compilation list",
    ),
}

# proto_compile = rule(
#     implementation = _proto_compile_impl,
    
#     outputs = {
#         # "descriptor": "%{name}/descriptor.source.bin",
#     },
#     output_to_genfiles = True,
# )

