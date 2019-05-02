load(
    "//:deps.bzl",
    "com_google_protobuf",
    "io_bazel_rules_dart",
    "io_bazel_rules_go",
)
load(
    "//protobuf:deps.bzl",
    "protobuf",
)

# Special dart_sdk_repository
load("//dart:sdk.bzl", "dart_sdk_repository")
load("//dart:dart_pub_deps.bzl", "dart_pub_deps")

def dart_sdk(**kwargs):
    """The dart sdk
    """
    name = "dart_sdk"
    if name not in native.existing_rules():
        dart_sdk_repository(
            name = name,
        )

def dart_pub_deps_protoc_plugin(**kwargs):
    """Dart pub dependencies for the dart protoc plugin
    """
    name = "dart_pub_deps_protoc_plugin"
    if name not in native.existing_rules():
        dart_pub_deps(
            name = name,
            spec = str(Label("//dart:pubspec.yaml")),

            # these overrides were determined by manually browsing
            # pub.dartlang.org, starting at protoc_plugin and going through all
            # transitive dependencies, pinning them to the version specified
            override = {
                "dart_style": "1.2.7",
                "analyzer": "0.36.1",
                "args": "1.5.1",
                "charcode": "1.1.2",
                "collection": "1.14.11",
                "convert": "2.1.1",
                "typed_data": "1.1.6",
                "crypto": "2.0.6",
                "front_end": "0.1.16",
                "kernel": "0.3.16",
                "package_config": "1.0.5",
                "path": "1.6.2",
                "yaml": "2.1.15",
                "source_span": "1.5.5",
                "term_glyph": "1.1.0",
                "string_scanner": "1.0.4",
                "meta": "1.1.7",
                "glob": "1.1.7",
                "async": "2.2.0",
                "html": "0.14.0+1",
                "csslib": "0.15.0",
                "pub_semver": "1.4.2",
                "watcher": "0.9.7+10",
                "fixnum": "0.10.9",
                "protobuf": "0.13.11",
            },
        )

def dart_pub_deps_grpc(**kwargs):
    """Dart pub dependencies for grpc
    """
    name = "dart_pub_deps_grpc"
    if name not in native.existing_rules():
        dart_pub_deps(
            name = name,
            spec = str(Label("//dart:pubspec-grpc.yaml")),

            # these overrides were determined by manually browsing
            # pub.dartlang.org, starting at grpc and going through all
            # transitive dependencies, pinning them to the version specified
            # there
            override = {
                "async": "2.2.0",
                "googlapis_auth": "0.2.7",
                "http": "0.12.0+2",
                "http_parser": "3.1.3",
                "pedantic": "1.5.0",
            },
        )

def dart_proto_compile(**kwargs):
    protobuf(**kwargs)
    io_bazel_rules_go(**kwargs)
    io_bazel_rules_dart(**kwargs)
    dart_sdk(**kwargs)
    dart_pub_deps_protoc_plugin(**kwargs)

def dart_grpc_compile(**kwargs):
    dart_proto_compile(**kwargs)

def dart_proto_library(**kwargs):
    dart_proto_compile(**kwargs)

def dart_grpc_library(**kwargs):
    dart_grpc_compile(**kwargs)
    dart_proto_library(**kwargs)
    dart_pub_deps_grpc(**kwargs)
