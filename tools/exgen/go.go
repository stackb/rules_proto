package main

var goUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()`)

var goProtoCompileRuleTemplate = mustTemplate(`load("//:compile.bzl", "proto_compile")
load("//:plugin.bzl", "proto_plugin")

def {{ .Rule.Name }}(**kwargs):
    # If importpath specified, declare a custom plugin that should correctly
    # predict the output location.
    importpath = kwargs.get("importpath")
    if importpath and not kwargs.get("plugins"):
        name = kwargs.get("name")
        name_plugin = name + "_plugin"
        proto_plugin(
            name = name_plugin,
            outputs = ["{package}/%s/{basename}.pb.go" % importpath],
            tool = "@com_github_golang_protobuf//protoc-gen-go",
        )
        kwargs["plugins"] = [name_plugin]
        kwargs.pop("importpath")
    # Define the default plugin if still not defined
    if not kwargs.get("plugins"):
        kwargs["plugins"] = [str(Label("//go:go"))]

    proto_compile(
        **kwargs
    )
`)

var goGrpcCompileRuleTemplate = mustTemplate(`load("//:compile.bzl", "proto_compile")
load("//:plugin.bzl", "proto_plugin")

def {{ .Rule.Name }}(**kwargs):
    # If importpath specified, declare a custom plugin that should correctly
    # predict the output location.
    importpath = kwargs.get("importpath")
    if importpath and not kwargs.get("plugins"):
        name = kwargs.get("name")
        name_plugin = name + "_plugin"
        proto_plugin(
            name = name_plugin,
            options = ["plugins=grpc"],
            outputs = ["{package}/%s/{basename}.pb.go" % importpath],
            tool = "@com_github_golang_protobuf//protoc-gen-go",
        )
        kwargs["plugins"] = [name_plugin]
        kwargs.pop("importpath")
    # Define the default plugin if still not defined
    if not kwargs.get("plugins"):
        kwargs["plugins"] = [str(Label("//go:grpc_go"))]

    proto_compile(
        **kwargs
    )
`)

var goProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Name}}:go_proto_compile.bzl", "go_proto_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    importpath = kwargs.get("importpath")
    visibility = kwargs.get("visibility")
    go_deps = kwargs.get("go_deps", [])

    name_pb = name + "_pb"

    go_proto_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )

    go_library(
        name = name,
        srcs = [name_pb],
        deps = go_deps + [
            "@com_github_golang_protobuf//proto:go_default_library",
        ],
        importpath = importpath,
        visibility = visibility,
    )
`)

var goGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Name}}:go_grpc_compile.bzl", "go_grpc_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    importpath = kwargs.get("importpath")
    visibility = kwargs.get("visibility")
    go_deps = kwargs.get("go_deps", [])

    name_pb = name + "_pb"

    go_grpc_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )

    go_library(
        name = name,
        srcs = [name_pb],
        deps = go_deps + [
            "@com_github_golang_protobuf//proto:go_default_library",
            "@org_golang_google_grpc//:go_default_library",
            "@org_golang_x_net//context:go_default_library",
        ],
        importpath = importpath,
        visibility = visibility,
    )
`)

var goProtoLibraryExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
	name = "person_{{ .Lang.Name }}_library",
    importpath = "github.com/stackb/rules_proto/{{ .Lang.Name }}/example/{{ .Rule.Name }}/person",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
    go_deps = [
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
)`)

var goGrpcLibraryExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
	name = "greeter_{{ .Lang.Name }}_library",
    importpath = "github.com/stackb/rules_proto/{{ .Lang.Name }}/example/{{ .Rule.Name }}/greeter",
	deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
    go_deps = [
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
)`)

func makeGo() *Language {
	return &Language{
		Dir:  "go",
		Name: "go",
		Rules: []*Rule{
			&Rule{
				Name:           "go_proto_compile",
				Implementation: goProtoCompileRuleTemplate,
				Plugins:        []string{"//go:go"},
				Usage:          goUsageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates *.go protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "go_grpc_compile",
				Implementation: goGrpcCompileRuleTemplate,
				Plugins:        []string{"//go:go", "//go:grpc_go"},
				Usage:          goUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates *.go protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "go_proto_library",
				Implementation: goProtoLibraryRuleTemplate,
				Usage:          goUsageTemplate,
				Example:        goProtoLibraryExampleTemplate,
				Doc:            "Generates *.go protobuf library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "go_grpc_library",
				Implementation: goGrpcLibraryRuleTemplate,
				Usage:          goUsageTemplate,
				Example:        goGrpcLibraryExampleTemplate,
				Doc:            "Generates *.go protobuf+gRPC library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
