load("//closure:compile.bzl", "closure_proto_compile")
load("//contrib/grpc/grpc-web:compile.bzl", "grpc_web_proto_compile")
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def web_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    name_pb_grpc = name + "_pb_grpc"

    closure_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
    )
    
    grpc_web_proto_compile(
        name = name_pb_grpc,
        deps = deps,
        visibility = visibility,
        verbose = verbose,
    )

    closure_js_library(
        name = name,
        srcs = [name_pb, name_pb_grpc],
        deps = [
            "@com_github_grpc_grpc_web//javascript/net/grpc/web:abstractclientbase",
            "@com_github_grpc_grpc_web//javascript/net/grpc/web:clientreadablestream",
            "@com_github_grpc_grpc_web//javascript/net/grpc/web:grpcwebclientbase",
            "@com_github_grpc_grpc_web//javascript/net/grpc/web:error",
            "@io_bazel_rules_closure//closure/library",
            "@io_bazel_rules_closure//closure/protobuf:jspb",
        ],
        lenient = True,
        visibility = visibility,
    )