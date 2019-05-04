load("//github.com/stackb/grpc.js:closure_grpc_compile.bzl", "closure_grpc_compile")
load("//closure:closure_proto_compile.bzl", "closure_proto_compile")
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def closure_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    name_pb_grpc = name + "_pb_grpc"

    closure_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    closure_grpc_compile(
        name = name_pb_grpc,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    closure_deps = kwargs.get("closure_deps", [])

    closure_js_library(
        name = name,
        srcs = [name_pb, name_pb_grpc],
        deps = [
            "@io_bazel_rules_closure//closure/library",
            "@io_bazel_rules_closure//closure/protobuf:jspb",
            "@com_github_stackb_grpc_js//js/grpc/stream:observer",
            "@com_github_stackb_grpc_js//js/grpc/stream/observer:call",
            "@com_github_stackb_grpc_js//js/grpc",
            "@com_github_stackb_grpc_js//js/grpc:api",
            "@com_github_stackb_grpc_js//js/grpc:options",
        ] + closure_deps,
        internal_descriptors = [
            name_pb + "/descriptor.source.bin",
            name_pb_grpc + "/descriptor.source.bin",
        ],
        suppress = [
            "JSC_IMPLICITLY_NULLABLE_JSDOC",
        ],
        library_level_checks = False,
        visibility = visibility,
    )
