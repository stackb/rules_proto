load(
    "@build_stack_rules_proto//language/rules:proto_language.bzl",
    "ProtoLanguageInfo",
    "proto_language_info_to_struct",
)
load(
    "@build_stack_rules_proto//language/rules:provider_test.bzl",
    "provider_test_implementation",
    "provider_test_macro",
    "provider_test_rule_pair",
)

def _proto_language_info_provider_test_impl(ctx):
    return provider_test_implementation(ctx, ProtoLanguageInfo, proto_language_info_to_struct)

_proto_language_info_provider_test, _proto_language_info_provider_run = provider_test_rule_pair(
    _proto_language_info_provider_test_impl,
    ProtoLanguageInfo,
)

def proto_language_info_provider_test(**kwargs):
    provider_test_macro(
        _proto_language_info_provider_test,
        _proto_language_info_provider_run,
        **kwargs
    )
