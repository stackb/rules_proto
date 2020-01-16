load("@bazel_tools//tools/build_defs/repo:jvm.bzl", "jvm_maven_import_external")

load(
    "//:deps.bzl",
    "com_google_protobuf",
    "io_grpc_grpc_java",
)

load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def com_google_guava_guava(**kwargs):
    if "com_google_guava_guava" not in native.existing_rules():
        jvm_maven_import_external(
            name = "com_google_guava_guava",
            artifact = "com.google.guava:guava:20.0",
            server_urls = ["https://central.maven.org/maven2"],
            artifact_sha256 = "36a666e3b71ae7f0f0dca23654b67e086e6c93d192f60ba5dfd5519db6c288c8",
            licenses = ["reciprocal"],  # CDDL License
        )


# From https://github.com/grpc/grpc-java/blob/master/repositories.bzl
def javax_annotation_javax_annotation_api(**kwargs):
    # Use //stub:javax_annotation for neverlink=1 support.
    if "javax_annotation_javax_annotation_api" not in native.existing_rules():
        jvm_maven_import_external(
            name = "javax_annotation_javax_annotation_api",
            artifact = "javax.annotation:javax.annotation-api:1.2",
            server_urls = ["https://central.maven.org/maven2"],
            artifact_sha256 = "5909b396ca3a2be10d0eea32c74ef78d816e1b4ead21de1d78de1f890d033e04",
            licenses = ["reciprocal"],  # CDDL License
        )

# From https://github.com/grpc/grpc-java/blob/master/repositories.bzl
def com_google_errorprone_error_prone_annotations(**kwargs):
    # Use //stub:javax_annotation for neverlink=1 support.
    if "com_google_errorprone_error_prone_annotations" not in native.existing_rules():
        jvm_maven_import_external(
            name = "com_google_errorprone_error_prone_annotations",
            artifact = "com.google.errorprone:error_prone_annotations:2.3.2",
            server_urls = ["https://central.maven.org/maven2"],
            artifact_sha256 = "357cd6cfb067c969226c442451502aee13800a24e950fdfde77bcdb4565a668d",
            licenses = ["notice"],  # Apache 2.0
        )
    if "error_prone_annotations" not in native.existing_rules():
        native.bind(
            name = "error_prone_annotations",
            actual = "@com_google_errorprone_error_prone_annotations//jar",
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
    com_google_errorprone_error_prone_annotations(**kwargs)
