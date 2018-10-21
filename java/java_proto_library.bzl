load("//java:java_proto_compile.bzl", "java_proto_compile")

def java_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    java_proto_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )
    native.java_library(
        name = name,
        srcs = [name_pb],
        deps = [str(Label("//java:proto_deps"))],
        exports = [
            str(Label("//java:proto_deps")),
        ],
        visibility = visibility,
    )

