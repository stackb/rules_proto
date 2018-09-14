load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

PLUGIN_VERSION = "1.9.0"

def java_proto_deps():
    existing = native.existing_rules()

    if "protoc_gen_grpc_java_linux_x86_64" not in existing:
        native.http_file(
            name = "protoc_gen_grpc_java_linux_x86_64",
            url = "https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/{plugin_version}/protoc-gen-grpc-java-{plugin_version}-linux-x86_64.exe".format(
                plugin_version = PLUGIN_VERSION,
            ),
            sha256 = "f20cc8c052eea904c5a979c140237696e3f187f35deac49cd70b16dc0635f463",
        )

    if "protoc_gen_grpc_java_linux_macosx" not in existing:
        native.http_file(
            name = "protoc_gen_grpc_java_linux_macosx",
            url = "https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/{plugin_version}/protoc-gen-grpc-java-{plugin_version}-macosx.exe".format(
                plugin_version = PLUGIN_VERSION,
            ),
            sha256 = "593937361f99e8b145fe29c78c71cdd00e8327ae88de010729479eb2acdc1de9",
        )

    if "io_grpc_grpc_java" not in existing:
        http_archive(
            name = "io_grpc_grpc_java",
            urls = ["https://github.com/grpc/grpc-java/archive/v1.15.0.tar.gz"],
            strip_prefix = "grpc-java-1.15.0",
            sha256 = "8a131e773b1c9c0442e606b7fc85d7fc6739659281589d01bd917ceda218a1c7",
        )
