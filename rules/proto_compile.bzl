"""Compilation functions
"""

load(
    ":proto_utils.bzl",
    "copy_file",
    "descriptor_proto_path",
    "flatten",
    "get_output_filename",
    "homogenize",
    "strip_path_prefix",
    "uniq",
)
load(
    ":proto_providers.bzl",
    "ProtoCompileInfo",
    "ProtoCompileInputInfo",
    "ProtoCompileOutputInfo",
)

def proto_compile(compilation):
    """Run protoc against a set of plugins.

    Args:
        compilation: <ProtoCompilationInfo>
    Returns:
        <list<ProtoCompileOutputInfo>> the compilation output infos
    """

    ###
    ### Part 1: setup variables used in scope
    ###

    verbose = compilation.verbose
    single_action = compilation.single_action
    transitive = compilation.transitive
    proto_info = compilation.proto_info
    plugins = compilation.plugins
    rel_outdir = compilation.rel_outdir
    full_outdir = compilation.full_outdir

    # <list<File>> The list of generated artifacts like 'foo_pb2.py' that we
    # expect to be produced.
    output_files = []

    # <list<File>> The list of generated artifact directories that we
    # expect to be produced.
    output_dirs_list = []

    # <list<ProtoCompileInputInfo>> The list of compilation structs that are to be run
    # together as a single action.  This list is empty unless the
    # 'single_action' bool is True.
    actions = []

    ###
    ### Part 2: iterate over plugins
    ###
    for plugin in plugins:
        ###
        ### Part 2.1: fetch plugin tool and runfiles
        ###

        # <list<File>> Files required for running the plugins
        plugin_runfiles = []

        # <list<opaque>> Plugin input manifests
        plugin_input_manifests = None

        # Get plugin name
        plugin_name = plugin.name
        if plugin.protoc_plugin_name:
            plugin_name = plugin.protoc_plugin_name

        # Add plugin executable if not a built-in plugin
        plugin_tool = None
        if plugin.tool_executable:
            plugin_tool = plugin.tool_executable

        # Add plugin runfiles if plugin has a tool
        if plugin.tool:
            plugin_runfiles, plugin_input_manifests = compilation.resolve_tools(tools = [plugin.tool])
            plugin_runfiles = plugin_runfiles.to_list()

        # Add extra plugin data files
        plugin_runfiles += plugin.data

        # Check plugin outputs
        if plugin.output_directory and (plugin.out or plugin.outputs):
            fail("Proto plugin {} cannot use output_directory in conjunction with outputs or out".format(plugin.name))

        ###
        ### Part 2.2: gather proto files and filter by exclusions
        ###

        # <list<File>> The filtered set of .proto files to compile
        protos = []

        for proto in proto_info.direct_sources:
            # Check for exclusion
            if any([
                proto.dirname.endswith(exclusion) or proto.path.endswith(exclusion)
                for exclusion in plugin.exclusions
            ]) or proto in protos:  # TODO: When using import_prefix, the ProtoInfo.direct_sources list appears to contain duplicate records, this line removes these. https://github.com/bazelbuild/bazel/issues/9127
                continue

            # Proto not excluded
            protos.append(proto)

        # Skip plugin if all proto files have now been excluded
        if len(protos) == 0:
            if verbose > 2:
                print('Skipping plugin "{}" for "{}" as all proto files have been excluded'.format(plugin.name, compilation.label))
            continue

        ###
        ### Part 2.3: declare per-proto generated outputs from plugin
        ###

        # <list<File>> The list of generated artifacts like 'foo_pb2.py' that we
        # expect to be produced by this plugin only
        plugin_outputs = []

        for proto in protos:
            for pattern in plugin.outputs:
                output_filename = get_output_filename(proto, pattern, proto_info)
                filename = "{}/{}".format(rel_outdir, output_filename)
                file = compilation.actions.declare_file(filename)
                plugin_outputs.append(file)

        # Append current plugin outputs to global outputs before looking at
        # per-plugin outputs; these are manually added globally as there may
        # be srcjar outputs.
        output_files.extend(plugin_outputs)

        ###
        ### Part 2.4: declare per-plugin artifacts
        ###

        # Some protoc plugins generate a set of output files (like python) while
        # others generate a single 'archive' file that contains the individual
        # outputs (like java). Jar outputs are gathered as a special case as we need to
        # post-process them to have a 'srcjar' extension (java_library rules don't
        # accept source jars with a 'jar' extension)

        out_file = None
        if plugin.out:
            # Define out file
            out_file = compilation.actions.declare_file("{}/{}".format(
                rel_outdir,
                plugin.out.replace("{name}", compilation.label.name),
            ))
            plugin_outputs.append(out_file)

            if not out_file.path.endswith(".jar"):
                # Add output direct to global outputs
                output_files.append(out_file)
            else:
                # Create .srcjar from .jar for global outputs
                output_files.append(copy_file(
                    compilation.actions,
                    out_file,
                    "{}.srcjar".format(out_file.basename.rpartition(".")[0]),
                    sibling = out_file,
                ))

        ###
        ### Part 2.5: declare plugin output directory
        ###

        # Some plugins outputs a structure that cannot be predicted from the
        # input file paths alone. For these plugins, we simply declare the
        # directory.

        if plugin.output_directory:
            out_file = compilation.actions.declare_directory(rel_outdir + "/" + "_plugin_" + plugin.name)
            plugin_outputs.append(out_file)
            output_dirs_list.append(out_file)

        ###
        ### Part 2.6: build args
        ###

        # <list<string> argument list for protoc execution
        args = []

        # Add descriptors
        pathsep = compilation.host_path_separator
        args.append("--descriptor_set_in={}".format(pathsep.join(
            [f.path for f in proto_info.transitive_descriptor_sets.to_list()],
        )))

        # Add plugin if not built-in
        if plugin_tool:
            # If Windows, mangle the path. It's done a bit awkwardly with
            # `host_path_seprator` as there is no simple way to figure out what's
            # the current OS.
            plugin_tool_path = None
            if compilation.host_path_separator == ";":
                plugin_tool_path = plugin_tool.path.replace("/", "\\")
            else:
                plugin_tool_path = plugin_tool.path

            args.append("--plugin=protoc-gen-{}={}".format(plugin_name, plugin_tool_path))

        # Add plugin out arg
        out_arg = out_file.path if out_file else full_outdir

        if plugin.options:
            opts_str = ",".join(
                [option.replace("{name}", compilation.label.name) for option in plugin.options],
            )
            if plugin.separate_options_flag:
                args.append("--{}_opt={}".format(plugin_name, opts_str))
            else:
                out_arg = "{}:{}".format(opts_str, out_arg)
        args.append("--{}_out={}".format(plugin_name, out_arg))

        # Add source proto files as descriptor paths
        for proto in protos:
            args.append(descriptor_proto_path(proto, proto_info))

        ###
        ### Part 2.7: schedule command
        ###

        action = ProtoCompileInputInfo(
            actions = compilation.actions,
            args = args,
            input_manifests = plugin_input_manifests if plugin_input_manifests else [],
            inputs = proto_info.transitive_descriptor_sets.to_list() + plugin_runfiles,  # Proto files are not inputs, as they come via the descriptor sets
            output_directory = full_outdir,
            outputs = plugin_outputs,
            protoc = compilation.protoc,
            tools = [plugin_tool] if plugin_tool else [],
            use_default_shell_env = plugin.use_built_in_shell_environment,
            verbose = verbose,
        )

        if single_action:
            actions.append(action)
        else:
            _proto_compile_action(action)

    if len(actions) > 0:
        _proto_compile_action(merge_proto_compile_input_infos(actions))

    ###
    ### Step 3: generate providers
    ###

    output_files_dict = {}
    if output_files:
        output_files_dict[full_outdir] = output_files

    # Gather transitive info
    transitive_output_dirs_list = []
    for transitive_info in compilation.transitive_outs:
        output_files_dict.update(**transitive_info.output_files)
        transitive_output_dirs_list.append(transitive_info.output_dirs)

    return [
        ProtoCompileOutputInfo(
            files = output_files_dict,
            dirs = depset(direct = output_dirs_list, transitive = transitive_output_dirs_list),
        ),
    ]

def merge_proto_compile_input_infos(infos):
    """merge_proto_compile_input_infos merges a list of infos into one.

    Args:
        infos: <list<ProtoCompileInputInfo>> A list of infos to merge.
    Returns:
        <ProtoCompileInputInfo> the single merged action.
    """
    return ProtoCompileInputInfo(
        actions = homogenize("actions", [info.actions for info in infos]),
        args = uniq(flatten([info.args for info in infos])),
        input_manifests = flatten([info.input_manifests for info in infos]),  # don't apply 'uniq' here to avoid 'Error: unhashable type: 'RunfilesSupplierImpl'
        inputs = uniq(flatten([info.inputs for info in infos])),
        output_directory = homogenize("output_directory", [info.output_directory for info in infos]),
        outputs = uniq(flatten([info.outputs for info in infos])),
        protoc = homogenize("protoc", [info.protoc for info in infos]),
        tools = uniq(flatten([info.tools for info in infos])),
        use_default_shell_env = homogenize("use_default_shell_env", [info.use_default_shell_env for info in infos]),
        verbose = homogenize("verbose", [info.verbose for info in infos]),
    )

def _proto_compile_action(info):
    """Declare a run_shell action for the given compilation.

    Args:
        info: <ProtoCompileInputInfo>
    """
    mnemonic = "ProtoCompile"
    full_outdir = info.output_directory

    command = ("mkdir -p '{}' && ".format(full_outdir)) + info.protoc.path + " $@"  # $@ is replaced with args list

    # Amend command with debug options
    if info.verbose > 0:
        print("{}:".format(mnemonic), info.protoc.path, info.args)

    if info.verbose > 1:
        command += " && echo '\n##### SANDBOX AFTER RUNNING PROTOC' && find . -type f "

    if info.verbose > 2:
        command = "echo '\n##### SANDBOX BEFORE RUNNING PROTOC' && find . -type l && " + command

    if info.verbose > 3:
        command = "env && " + command
        for f in info.inputs:
            print("INPUT:", f.path)
        for f in info.tools:
            print("TOOL:", f.path)
        for f in info.outputs:
            print("EXPECTED OUTPUT:", f.path)

    args = info.actions.args()
    args.add_all(info.args)

    info.actions.run_shell(
        progress_message = "Compiling protoc outputs...",
        command = command,
        arguments = [args],
        inputs = info.inputs,
        tools = [info.protoc] + info.tools,
        outputs = info.outputs,
        use_default_shell_env = info.use_default_shell_env,
        input_manifests = info.input_manifests,
    )

def coalesce_outputs(ctx, outs):
    """Aggregate output files and dirs created by the aspect as it has walked the deps

    Args:
        ctx: <ctx>
        outs: <list<ProtoCompileOutputInfo>>
    Returns:
        <list<providers>>
    """
    output_files_dicts = [info.files for info in outs]
    output_dirs = depset(transitive = [info.dirs for info in outs])

    # Check merge_directories and prefix_path
    if not ctx.attr.merge_directories and ctx.attr.prefix_path:
        fail("Attribute prefix_path cannot be set when merge_directories is false")

    # Build outputs
    final_output_files = {}
    final_output_files_list = []
    final_output_dirs = depset()
    prefix_path = ctx.attr.prefix_path

    if not ctx.attr.merge_directories:
        # Pass on outputs directly when not merging
        for output_files_dict in output_files_dicts:
            final_output_files.update(**output_files_dict)
            final_output_files_list = [f for files in final_output_files.values() for f in files]
        final_output_dirs = output_dirs

    elif output_dirs:
        # If we have any output dirs specified, we declare a single output
        # directory and merge all files in one go. This is necessary to prevent
        # path prefix conflicts

        # Declare single output directory
        dir_name = ctx.label.name
        if prefix_path:
            dir_name = dir_name + "/" + prefix_path
        new_dir = ctx.actions.declare_directory(dir_name)
        final_output_dirs = depset(direct = [new_dir])

        # Build copy command for directory outputs
        # Use cp {}/. rather than {}/* to allow for empty output directories from a plugin (e.g when no service exists,
        # so no files generated)
        command_parts = ["cp -r {} '{}'".format(
            " ".join(["'" + d.path + "/.'" for d in output_dirs.to_list()]),
            new_dir.path,
        )]

        # Extend copy command with file outputs
        command_input_files = []
        for output_files_dict in output_files_dicts:
            for root, files in output_files_dict.items():
                for file in files:
                    # Strip root from file path
                    path = strip_path_prefix(file.path, root)

                    # Prefix path is contained in new_dir.path created above and
                    # used below

                    # Add command to copy file to output
                    command_input_files.append(file)
                    command_parts.append("cp '{}' '{}'".format(
                        file.path,
                        "{}/{}".format(new_dir.path, path),
                    ))

        # Add debug options
        if ctx.attr.verbose > 1:
            command_parts = command_parts + ["echo '\n##### SANDBOX AFTER MERGING DIRECTORIES'", "find . -type l"]
        if ctx.attr.verbose > 2:
            command_parts = ["echo '\n##### SANDBOX BEFORE MERGING DIRECTORIES'", "find . -type l"] + command_parts
        if ctx.attr.verbose > 0:
            print("Directory merge command: {}".format(" && ".join(command_parts)))

        # Copy directories and files to shared output directory in one action
        ctx.actions.run_shell(
            mnemonic = "CopyDirs",
            inputs = depset(direct = command_input_files, transitive = [output_dirs]),
            outputs = [new_dir],
            command = " && ".join(command_parts),
            progress_message = "copying directories and files to {}".format(new_dir.path),
        )

    else:
        # Otherwise, if we only have output files, build the output tree by
        # aggregating files created by aspect into one directory

        output_root = ctx.bin_dir.path + "/"

        if ctx.label.workspace_root:
            output_root += ctx.label.workspace_root + "/"

        if ctx.label.package:
            output_root += ctx.label.package + "/"

        output_root += ctx.label.name
        final_output_files[output_root] = []

        for output_files_dict in output_files_dicts:
            for root, files in output_files_dict.items():
                for file in files:
                    # Strip root from file path
                    base = file.basename if ctx.attr.output_in_package_root else file.path
                    path = strip_path_prefix(base, root)

                    # Prepend prefix path if given
                    if prefix_path:
                        path = prefix_path + "/" + path

                    # Copy file to output
                    final_output_files[output_root].append(copy_file(
                        ctx.actions,
                        file,
                        path,
                        # TODO(pcj): is this necessary?
                        # "{}/{}".format(ctx.label.name, path),
                    ))

        final_output_files_list = final_output_files[output_root]

    # Create depset containing all outputs
    if ctx.attr.merge_directories:
        # If we've merged directories, we have copied files/dirs that are now direct rather than
        # transitive dependencies
        all_outputs = depset(direct = final_output_files_list + final_output_dirs.to_list())
    else:
        # If we have not merged directories, all files/dirs are transitive
        all_outputs = depset(
            transitive = [depset(direct = final_output_files_list), final_output_dirs],
        )

    # Create default and proto compile providers
    return [
        ProtoCompileInfo(
            label = ctx.label,
            output_files = final_output_files,
            output_dirs = final_output_dirs,
        ),
        DefaultInfo(
            files = all_outputs,
            data_runfiles = ctx.runfiles(transitive_files = all_outputs),
        ),
    ]

def get_output_directories(bin_dir, label, prefix):
    """get_output_directories returns the relative and full dirs.

    Args:
        bin_dir: <string> the ctx.bin_dir argument
        label: <Label> the ctx.label argument
        prefix: <strig> an optional prefix
    Returns:
        The directory where the outputs will be generated, relative to
        the package. This contains the aspect _prefix attr to disambiguate
        different aspects that may share the same plugins and would otherwise
        try to touch the same file.
    """

    rel_outdir = "{}/{}_verb0".format(label.name, prefix)
    full_outdir = bin_dir.path + "/"
    if label.workspace_root:
        full_outdir += label.workspace_root + "/"
    if label.package:
        full_outdir += label.package + "/"
    full_outdir += rel_outdir
    return rel_outdir, full_outdir
