load("//github.com/gogo/protobuf:gogofaster_proto_compile.bzl", "gogofaster_proto_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("//go:utils.bzl", "get_importmappings")

wkt_mappings = get_importmappings({
    "google/protobuf/any.proto": "github.com/gogo/protobuf/types",
    "google/protobuf/duration.proto": "github.com/gogo/protobuf/types",
    "google/protobuf/struct.proto": "github.com/gogo/protobuf/types",
    "google/protobuf/timestamp.proto": "github.com/gogo/protobuf/types",
    "google/protobuf/wrappers.proto": "github.com/gogo/protobuf/types",
})

def gogofaster_proto_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    kwargs["plugin_options"] = kwargs.get("plugin_options", []) + get_importmappings(kwargs.get("importmap", {})) + wkt_mappings
    gogofaster_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k not in ("name", "importpath", "importmap", "go_deps")} # Forward args except name, importpath, importmap and go_deps
    )

    # Create gogo library
    go_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = kwargs.get("go_deps", []) + [
            "@com_github_gogo_protobuf//proto:go_default_library",
        ],
        importpath = kwargs.get("importpath"),
        visibility = kwargs.get("visibility"),
    )
