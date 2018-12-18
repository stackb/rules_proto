
load("//:deps.bzl", 
    "com_google_protobuf",
    # "dart_pub_deps_protoc_plugin",
    # "dart_pub_deps_grpc",
    # "dart_sdk",
    "io_bazel_rules_dart",
    "io_bazel_rules_go",
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

            # these overrides were determined by manually browsing pub.dartlang.org, 
            # starting at protoc_plugin and going through all transitive dependencies,
            # pinning them to the version specified there (basically, seems like latest)
            
            override = {
                "analyzer": "0.34.0",
                "args": "1.5.1",
                "async": "2.0.8",
                "charcode": "1.1.2",
                "collection": "1.14.11",
                "convert": "2.0.2",
                "crypto": "2.0.6",
                "csslib": "0.14.6",
                "dart_style": "1.1.2",
                "front_end": "0.1.7",
                "fixnum": "0.10.9",
                "html": "0.13.3+3",
                "kernel": "0.3.7",
                "logging": "0.11.3+2",
                "meta": "1.1.6",
                "package_config": "1.0.5",
                "path": "1.6.2",
                "plugin": "0.2.0+3",
                "pub_semver": "1.4.2",
                "source_span": "1.4.1",
                "string_scanner": "1.0.4",
                "typed_data": "1.1.6",
                "watcher": "0.9.7+10",
                "utf": "0.9.0+5",
                "yaml": "2.1.15",
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

            # these overrides were determined by manually browsing pub.dartlang.org, 
            # starting at protoc_plugin and going through all transitive dependencies,
            # pinning them to the version specified there (basically, seems like latest)
            
            override = {
                "async": "2.0.8",
                "collection": "1.14.11",
                "googlapis_auth": "0.2.6",
                "crypto": "2.0.6",
                "convert": "2.0.2",
                "charcode": "1.1.2",
                "typed_data": "1.1.6",
                "http": "0.12.0",
                "http_parser": "3.1.3",
                "source_span": "1.4.1",
                "string_scanner": "1.0.4",
                "path": "1.6.2",
                "http2": "0.1.9",
                "meta": "1.1.6",
            },
        )



def dart_proto_compile(**kwargs):
    com_google_protobuf(**kwargs)
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
