"""rules_proto based language rules."""
# Ref: https://github.com/stackb/rules_proto/blob/master/pkg/protoc/starlark_rule.go
def _make_proto_python_rule(rule_cfg, protoc_cfg):
    return gazelle.Rule(
        kind = "mypy_stubs",
        name = protoc_cfg.proto_library.base_name + "_pyi_library",
        attrs = {
           "srcs": [x.replace(".proto", "_pb2.pyi") for x in protoc_cfg.proto_library.srcs],
           "visibility": rule_cfg.visibility,
        },
     )

def _provide_python(rctx, pctx):
    return struct(
        name = "proto_python_stubs_library",
        rule = lambda: _make_proto_python_rule(rctx, pctx),
        experimental_resolve_attr = "deps",
    )

protoc.Rule(
    name = "proto_python_stubs_library",
    load_info = lambda: gazelle.LoadInfo(name = "//example/starlark_proto/lib:rules.bzl", symbols = ["mypy_stubs"]),
    kind_info = lambda: gazelle.KindInfo(mergeable_attrs = {"srcs": True}, resolve_attrs = {"deps": True}),
    provide_rule = _provide_python,
)

# vim: set ft=python
