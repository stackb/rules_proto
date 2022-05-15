"""starlark rules definitions"""

def _make_java_library_rule(rctx, pctx):
    r = gazelle.Rule(
        kind = "java_library",
        name = pctx.proto_library.base_name + "_java_library",
        attrs = {
            "srcs": [pctx.proto_library.base_name + ".srcjar"],
        },
    )

def _provide_java_library(rctx, pctx):
    return struct(
        name = "java_library",
        rule = lambda: _make_java_library_rule(rctx, pctx),
    )

protoc.Rule(
    name = "java_library",
    load_info = lambda: gazelle.LoadInfo(name = "@rules_java//java:defs.bzl", symbols = ["java_library"]),
    kind_info = lambda: gazelle.KindInfo(merge_attrs = {"srcs": True}, resolve_attrs = {"deps": True}),
    provide_rule = _provide_java_library,
)
