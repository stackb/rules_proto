load("//:compile.bzl", "proto_compile")

load("@build_bazel_rules_swift//swift:swift.bzl", "swift_proto_library")

def swift_proto_compile(**kwargs):
    swift_proto_library(**kwargs)    
    # proto_compile(
    #     plugins = [str(Label("//swift:swift"))],
    #     **kwargs
    # )

def swift_grpc_compile(**kwargs):
    # proto_compile(
    #     plugins = [str(Label("//swift:grpc_swift"))],
    #     **kwargs
    # )
