"""providers.bzl
"""

ProtoPluginInfo = provider(
    "ProtoPluginInfo provides metadata about how a protoc plugin should be run",
    fields = {
        "name": "The proto plugin name",
        "label": "The proto plugin label",
        "options": "A list of options to pass to the compiler for this plugin",
        "tool": "The plugin binary executable",
        "tool_target": "The plugin tool target attr",
        "use_built_in_shell_environment": "Whether the tool should use the built in shell environment or not",
        "protoc_plugin_name": "The name used for the plugin binary on the protoc command line. Useful for targeting built-in plugins. Uses plugin name when not set",
        "exclusions": "Exclusion filters to apply when generating outputs with this plugin. Used to prevent generating files that are included in the protobuf library, for example. Can exclude either by proto name prefix or by proto folder prefix",
        "mods": "awk expressions to apply to the output files",
        "data": "Additional files required for running the plugin",
        "out": "The format for the --x_out argument.  Defaults to to {BIN_DIR}",
        "supplementary_proto_deps": "Additional proto dependencies whose descriptors/files should be included in all protoc invocations",
        "deps": "The list of workspace dependencies for this plugin",
    },
)

ProtoCompileInfo = provider(
    "ProtoCompileInfo provides downstream rules with the outputs of proto_compile",
    fields = {
        "label": "The proto_compile rule label (type Label)",
        "outputs": "The output files from the rule (type List[File])",
        "output_files_by_rel_path": "The output files from the rule (type Dict[String,List[File]]).  The keys are the package-relative paths of the file and values are the file objects.",
    },
)
