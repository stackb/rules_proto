load("@bazel_tools//tools/build_defs/repo:jvm.bzl", "jvm_maven_import_external")

load(
    "//:deps.bzl",
    "com_google_guava_guava",
    "com_google_protobuf",
    "io_grpc_grpc_java",
)

load(
    "//protobuf:deps.bzl",
    "protobuf",
)


# From https://github.com/grpc/grpc-java/blob/master/repositories.bzl
def javax_annotation_javax_annotation_api(**kwargs):
    # Use //stub:javax_annotation for neverlink=1 support.
    if "javax_annotation_javax_annotation_api" not in native.existing_rules():
        jvm_maven_import_external(
            name = "javax_annotation_javax_annotation_api",
            artifact = "javax.annotation:javax.annotation-api:1.2",
            server_urls = ["http://central.maven.org/maven2"],
            artifact_sha256 = "5909b396ca3a2be10d0eea32c74ef78d816e1b4ead21de1d78de1f890d033e04",
            licenses = ["reciprocal"],  # CDDL License
        )

def java_proto_compile(**kwargs):
    protobuf(**kwargs)

def java_grpc_compile(**kwargs):
    java_proto_compile(**kwargs)
    io_grpc_grpc_java(**kwargs)

def java_proto_library(**kwargs):
    java_proto_compile(**kwargs)
    javax_annotation_javax_annotation_api(**kwargs)
    com_google_guava_guava(**kwargs)

def java_grpc_library(**kwargs):
    java_grpc_compile(**kwargs)
    java_proto_library(**kwargs)
