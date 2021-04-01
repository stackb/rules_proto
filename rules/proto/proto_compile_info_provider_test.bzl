load("@build_stack_rules_proto//rules:proto_providers.bzl", "ProtoCompileInfo")
load(
    "@build_stack_rules_proto//rules/proto:provider_test.bzl",
    "provider_test_implementation",
    "provider_test_macro",
    "provider_test_rule_pair",
)

def proto_compile_info_to_struct(info):
    return struct(
        label = str(info.label),
        output_files = output_files_to_dict(info.output_files),
        output_dirs = [f.short_path for f in info.output_dirs.to_list()],
    )

def output_files_to_dict(d):
    return struct(**{k: [f.short_path for f in files] for k, files in d.items()})

def _proto_compile_provider_test_impl(ctx):
    return provider_test_implementation(ctx, ProtoCompileInfo, proto_compile_info_to_struct)

_proto_compile_provider_test, _proto_compile_provider_run = provider_test_rule_pair(
    _proto_compile_provider_test_impl,
    ProtoCompileInfo,
)

def proto_compile_info_provider_test(**kwargs):
    provider_test_macro(
        _proto_compile_provider_test,
        _proto_compile_provider_run,
        **kwargs
    )
