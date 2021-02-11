load("@rules_proto//proto:defs.bzl", "ProtoInfo")
load(
    "@build_stack_rules_proto//rules/proto:provider_test.bzl",
    "provider_test_implementation",
    "provider_test_macro",
    "provider_test_rule_pair",
    "redact_host_configuration",
)

def proto_info_to_struct(info):
    return struct(
        check_deps_sources = [f.short_path for f in info.check_deps_sources.to_list()],
        direct_descriptor_set = info.direct_descriptor_set.short_path,
        direct_sources = [f.short_path for f in info.direct_sources],
        proto_source_root = info.proto_source_root,
        transitive_descriptor_sets = [f.short_path for f in info.transitive_descriptor_sets.to_list()],
        transitive_proto_path = [redact_host_configuration(f) for f in info.transitive_proto_path.to_list()],
        transitive_sources = [f.short_path for f in info.transitive_sources.to_list()],
    )

def _proto_info_provider_test_impl(ctx):
    return provider_test_implementation(ctx, ProtoInfo, proto_info_to_struct)

_proto_info_provider_test, _proto_info_provider_run = provider_test_rule_pair(
    _proto_info_provider_test_impl,
    ProtoInfo,
)

def proto_info_provider_test(**kwargs):
    provider_test_macro(
        _proto_info_provider_test,
        _proto_info_provider_run,
        **kwargs
    )
