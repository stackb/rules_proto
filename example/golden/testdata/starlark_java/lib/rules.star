"""starlark rules definitions"""

def _make_java_library_rule(rctx, pctx):
    r = gazelle.Rule(
        kind = "java_library",
        name = pctx.proto_library.base_name + "_java_library",
        attrs = {
            "srcs": [pctx.proto_library.base_name + ".srcjar"],
            "deps": rctx.deps,
            "visibility": rctx.visibility,
        },
    )
    return r

def _provide_java_library(rctx, pctx):
    return struct(
        name = "java_library",
        rule = lambda: _make_java_library_rule(rctx, pctx),
        experimental_resolve_attr = "deps",
    )

protoc.Rule(
    name = "java_library",
    load_info = lambda: gazelle.LoadInfo(name = "@rules_java//java:defs.bzl", symbols = ["java_library"]),
    kind_info = lambda: gazelle.KindInfo(mergeable_attrs = {"srcs": True}, resolve_attrs = {"deps": True}),
    provide_rule = _provide_java_library,
)

# --------------------------------------------------

def _make_java_wrapper_rule(_rctx, pctx):
    r = gazelle.Rule(
        kind = "java_wrapper",
        name = pctx.proto_library.base_name + "_java_wrap",
        attrs = {
            "javalib": pctx.proto_library.base_name + "_java_library",
            "deps": [],
        },
    )
    return r

def _provide_java_wrapper(rctx, pctx):
    return struct(
        name = "java_wrapper",
        rule = lambda: _make_java_wrapper_rule(rctx, pctx),
        experimental_resolve_attr = "deps",
    )

protoc.Rule(
    name = "java_wrapper",
    load_info = lambda: gazelle.LoadInfo(name = "//:defs.bzl", symbols = ["java_wrapper"]),
    kind_info = lambda: gazelle.KindInfo(mergeable_attrs = {"javalib": True}, resolve_attrs = {"deps": True}),
    provide_rule = _provide_java_wrapper,
)
