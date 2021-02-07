load(
    "@build_stack_rules_proto//:plugin.bzl",
    "ProtoPluginInfo",
    "proto_plugin_info_to_struct",
)
load(
    "@build_stack_rules_proto//:provider_test.bzl",
    "provider_test_implementation",
    "provider_test_macro",
    "provider_test_rule_pair",
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