load("//swift:swift_proto_compile.bzl", "swift_proto_compile")
load("@build_bazel_rules_swift//swift:swift.bzl", _swift_proto_library = "swift_proto_library")

def swift_proto_library(**kwargs):
    _swift_proto_library(**kwargs)

