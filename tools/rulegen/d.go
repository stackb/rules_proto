package main

var dWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@io_bazel_rules_d//d:d.bzl", "d_repositories")

d_repositories()`)

var dProtoCompileExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "person_{{ .Lang.Name }}_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)`)

var dCompileRuleTemplate = mustTemplate(`load("//:compile.bzl", "proto_compile")
load("//:plugin.bzl", "proto_plugin")

def {{ .Rule.Name }}(**kwargs):
    # If package specified, declare a custom plugin that should correctly
    # predict the output location.
    package = kwargs.get("package")
    if package and not kwargs.get("plugins"):
        name_plugin = kwargs.get("name") + "_plugin"
        proto_plugin(
            name = name_plugin,
            outputs = ["%s/{basename}.d" % package],
            tool = "{{ with (index .Lang.Plugins (index .Rule.Plugins 0)) }}{{ .Tool }}{{ end }}",
        )
        kwargs["plugins"] = [name_plugin]
        kwargs.pop("package")

    # Define the default plugin if still not defined
    if not kwargs.get("plugins"):
        kwargs["plugins"] = [Label("{{ index .Rule.Plugins 0 }}")]

    proto_compile(
        **kwargs
    )`)

var dProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir}}:d_proto_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@io_bazel_rules_d//d:d.bzl", "d_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create {{ .Lang.Name }} library
    d_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            "@com_github_dcarp_protobuf_d//:protosrc",
            "@com_github_dcarp_protobuf_d//:protobuf",
		],
        imports = ["external/com_github_dcarp_protobuf_d/src"],
        visibility = kwargs.get("visibility"),
    )`)

var dGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir}}:d_grpc_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@io_bazel_rules_d//d:d.bzl", "d_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create {{ .Lang.Name }} library
    d_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
        ],
        visibility = kwargs.get("visibility"),
    )`)

func makeD() *Language {
	return &Language{
		Dir:   "d",
		Name:  "d",
		Flags: commonLangFlags,
		Plugins: map[string]*Plugin{
			"//d:d": &Plugin{
				Tool: "@com_github_dcarp_protobuf_d//:protoc-gen-d",
			},
		},
		Rules: []*Rule{
			&Rule{
				Name:             "d_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//d:d"},
				WorkspaceExample: dWorkspaceTemplate,
				BuildExample:     dProtoCompileExampleTemplate,
				Doc:              "Generates d protobuf artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			// &Rule{
			// 	Name:             "d_grpc_compile",
			//  Kind:             "grpc",
			// 	Implementation:   aspectRuleTemplate,
			// 	Plugins:          []string{"//d:grpc_d"},
			// 	WorkspaceExample: dWorkspaceTemplate,
			// 	BuildExample:     grpcCompileExampleTemplate,
			// 	Doc:              "Generates d protobuf+gRPC artifacts",
			// 	Attrs:            aspectProtoCompileAttrs,
			// },
			// &Rule{
			// 	Name:             "d_proto_library",
			//  Kind:             "proto",
			// 	Implementation:   dProtoLibraryRuleTemplate,
			// 	WorkspaceExample: dWorkspaceTemplate,
			// 	BuildExample:     protoLibraryExampleTemplate,
			// 	Doc:              "Generates d protobuf library",
			// 	Attrs:            aspectProtoCompileAttrs,
			// },
			// &Rule{
			// 	Name:             "d_grpc_library",
			//  Kind:             "grpc",
			// 	Implementation:   dGrpcLibraryRuleTemplate,
			// 	WorkspaceExample: dWorkspaceTemplate,
			// 	BuildExample:     grpcLibraryExampleTemplate,
			// 	Doc:              "Generates d protobuf+gRPC library",
			// 	Attrs:            aspectProtoCompileAttrs,
			// },
		},
	}
}
