load("@build_stack_rules_proto//:plugin.bzl", "ProtoPluginInfo")
load(
    "@build_stack_rules_proto//:provider_test.bzl",
    "provider_test_implementation",
    "provider_test_macro",
    "provider_test_rule_pair",
)

def proto_plugin_info_to_struct(info):
    return struct(
        name = info.name,
        options = info.options,
        outputs = info.outputs,
        output_directory = info.output_directory,
        tool = info.tool if info.tool else "",
        tool_executable = info.tool_executable.short_path if info.tool_executable else "",
        use_built_in_shell_environment = info.use_built_in_shell_environment,
        protoc_plugin_name = info.protoc_plugin_name,
        exclusions = info.exclusions,
        data = [f.short_path for f in info.data],
        separate_options_flag = info.separate_options_flag,
    )

def _proto_plugin_info_provider_test_impl(ctx):
    return provider_test_implementation(ctx, ProtoPluginInfo, proto_plugin_info_to_struct)

_proto_plugin_info_provider_test, _proto_plugin_info_provider_run = provider_test_rule_pair(
    _proto_plugin_info_provider_test_impl,
    ProtoPluginInfo,
)

def proto_plugin_info_provider_test(**kwargs):
    provider_test_macro(
        _proto_plugin_info_provider_test,
        _proto_plugin_info_provider_run,
        **kwargs
    )
