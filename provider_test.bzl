load("@build_stack_rules_proto//tools/gencopy:gencopy.bzl", "gencopy_action", "gencopy_attrs", "gencopy_config")

def provider_test_implementation(ctx, provider, provider_to_struct):
    outputs = []

    config = gencopy_config(ctx)

    n = 0
    for info in [dep[provider] for dep in ctx.attr.deps]:
        n += 1
        info_proto = ctx.actions.declare_file("%s.%d.prototext" % (ctx.label.name, n))
        outputs.append(info_proto)
        ctx.actions.write(
            output = info_proto,
            content = provider_to_struct(info).to_proto(),
        )

    script, runfiles = gencopy_action(ctx, config, outputs)

    return [
        DefaultInfo(
            files = depset(outputs),
            runfiles = runfiles,
            executable = script,
        ),
    ]

def provider_test_rule(implementation, provider):
    return rule(
        implementation = implementation,
        attrs = dict(
            gencopy_attrs,
            deps = attr.label_list(
                doc = "The upstream rule that provides the info to test",
                providers = [provider],
            ),
        ),
        executable = True,
        test = True,
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


def provider_test_macro(provider_test_rule, **kwargs):
    srcs = kwargs.pop("srcs", [])
    deps = kwargs.pop("deps", [])
    name = kwargs.pop("name")

    update_target_label_name = 'golden'
    update_name = "%s.%s" % (name, update_target_label_name)

    provider_test_rule(
        name = name,
        deps = deps,
        srcs = srcs,
        mode = "check",
        update_target_label_name = update_target_label_name,
    )

    provider_test_rule(
        name = update_name,
        deps = deps,
        mode = "update",
        update_target_label_name = update_target_label_name,
    )
