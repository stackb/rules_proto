load("@rules_proto//proto:defs.bzl", "ProtoInfo")
load("//tools/gencopy:gencopy.bzl", "gencopy_action", "gencopy_attrs", "gencopy_config")

def _proto_info_provider_test_impl(ctx):
    outputs = []

    config = gencopy_config(ctx)

    n = 0
    for info in [dep[ProtoInfo] for dep in ctx.attr.deps]:
        n += 1
        info_proto = ctx.actions.declare_file("%s.%d.prototext" % (ctx.label.name, n))
        outputs.append(info_proto)
        ctx.actions.write(
            output = info_proto,
            content = proto_info_to_struct(info).to_proto(),
        )

    script, runfiles = gencopy_action(ctx, config, outputs)

    return [
        DefaultInfo(
            files = depset(outputs),
            runfiles = runfiles,
            executable = script,
        ),
    ]

_proto_info_provider_test = rule(
    implementation = _proto_info_provider_test_impl,
    attrs = dict(
        gencopy_attrs,
        deps = attr.label_list(
            doc = "The proto_library rule",
            providers = [ProtoInfo],
        ),
    ),
    executable = True,
    test = True,
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

# transform:
# from "bazel-out/darwin-fastbuild/bin/external/com_google_protobuf/_virtual_imports/any_proto"
# to "bazel-out/{HOST_CONFIGURATION}/bin/external/com_google_protobuf/_virtual_imports/any_proto}
def redact_host_configuration(filename):
    if not filename.startswith("bazel-out/"):
        return filename
    parts = filename.split('/')
    parts[1] = "{HOST_CONFIGURATION}"
    return "/".join(parts)

def proto_info_provider_test(**kwargs):
    srcs = kwargs.pop("srcs", [])
    deps = kwargs.pop("deps", [])
    name = kwargs.pop("name")

    update_name = name + ".checkin"

    _proto_info_provider_test(
        name = name,
        deps = deps,
        srcs = srcs,
        mode = "check",
    )

    _proto_info_provider_test(
        name = update_name,
        deps = deps,
        mode = "update",
    )
