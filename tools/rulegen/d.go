package main

var dWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Lang.Name }}_deps")

{{ .Lang.Name }}_deps()

load("@io_bazel_rules_d//d:d.bzl", "d_repositories")

d_repositories()`)

var dProtoCompileExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:defs.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "person_{{ .Lang.Name }}_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)`)

var dProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir}}:d_proto_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@io_bazel_rules_d//d:d.bzl", "d_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create {{ .Lang.Name }} library
    d_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            "@com_github_dcarp_protobuf_d//:protosrc",
            "@com_github_dcarp_protobuf_d//:protobuf",
		],
        imports = ["external/com_github_dcarp_protobuf_d/src", name_pb],
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
		Notes: mustTemplate(`These rules use the protoc-gen-d plugin, which only supports proto3 .proto files.`),
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
				BazelCIExcludePlatforms: []string{"windows", "macos"},
			},
// 			&Rule{
// 				Name:             "d_grpc_compile",
// 				Kind:             "grpc",
// 				Implementation:   aspectRuleTemplate,
// 				Plugins:          []string{"//d:grpc_d"}, # TODO: Try https://github.com/huntlabs/grpc-dlang
// 				WorkspaceExample: dWorkspaceTemplate,
// 				BuildExample:     grpcCompileExampleTemplate,
// 				Doc:              "Generates d protobuf+gRPC artifacts",
// 				Attrs:            aspectProtoCompileAttrs,
// 				BazelCIExcludePlatforms: []string{"windows", "macos"},
// 			},
			&Rule{
				Name:             "d_proto_library",
				Kind:             "proto",
				Implementation:   dProtoLibraryRuleTemplate,
				WorkspaceExample: dWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates d protobuf library",
				Attrs:            aspectProtoCompileAttrs,
				BazelCIExcludePlatforms: []string{"windows", "macos"},
			},
// 			&Rule{
// 				Name:             "d_grpc_library",
// 				Kind:             "grpc",
// 				Implementation:   dGrpcLibraryRuleTemplate,
// 				WorkspaceExample: dWorkspaceTemplate,
// 				BuildExample:     grpcLibraryExampleTemplate,
// 				Doc:              "Generates d protobuf+gRPC library",
// 				Attrs:            aspectProtoCompileAttrs,
// 				BazelCIExcludePlatforms: []string{"windows", "macos"},
// 			},
		},
	}
}
