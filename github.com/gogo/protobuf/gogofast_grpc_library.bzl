load("//github.com/gogo/protobuf:gogofast_grpc_compile.bzl", "gogofast_grpc_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("//go:utils.bzl", "get_importmappings")

wkt_mappings = get_importmappings({
	"google/protobuf/any.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/duration.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/struct.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/timestamp.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/wrappers.proto": "github.com/gogo/protobuf/types",
})

def gogofast_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    importpath = kwargs.get("importpath")
    visibility = kwargs.get("visibility")
    go_deps = kwargs.get("go_deps", [])

    name_pb = name + "_pb"

    gogofast_grpc_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        plugin_options = get_importmappings(kwargs.pop("importmap", {})) + wkt_mappings,
        visibility = visibility,
    )

    go_library(
        name = name,
        srcs = [name_pb],
        deps = go_deps + [
            "@com_github_gogo_protobuf//proto:go_default_library",
            "@org_golang_google_grpc//:go_default_library",
            "@org_golang_x_net//context:go_default_library",
        ],
        importpath = importpath,
        visibility = visibility,
    )

