load(
    "@build_stack_rules_proto//language/rules:proto_rule.bzl",
    "ProtoRuleInfo",
    "proto_rule_info_to_struct",
)
load(
    "@build_stack_rules_proto//language/rules:provider_test.bzl",
    "provider_test_implementation",
    "provider_test_macro",
    "provider_test_rule_pair",
)

def _proto_rule_info_provider_test_impl(ctx):
    return provider_test_implementation(ctx, ProtoRuleInfo, proto_rule_info_to_struct)

_proto_rule_info_provider_test, _proto_rule_info_provider_run = provider_test_rule_pair(
    _proto_rule_info_provider_test_impl,
    ProtoRuleInfo,
)

def proto_rule_info_provider_test(**kwargs):
    provider_test_macro(
        _proto_rule_info_provider_test,
        _proto_rule_info_provider_run,
        **kwargs
    )
