load("//swift:swift_grpc_compile.bzl", "swift_grpc_compile")
load("@build_bazel_rules_swift//swift:swift.bzl", _swift_proto_library = "swift_proto_library")

def swift_grpc_library(**kwargs):
    _swift_proto_library(**kwargs)
