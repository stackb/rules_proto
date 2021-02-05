package main

var nodeWorkspaceTemplate = mustTemplate(`load("@rules_proto_grpc//{{ .Lang.Dir }}:repositories.bzl", rules_proto_grpc_{{ .Lang.Name }}_repos="{{ .Lang.Name }}_repos")

rules_proto_grpc_{{ .Lang.Name }}_repos()

load("@build_bazel_rules_nodejs//:index.bzl", "yarn_install")

yarn_install(
    name = "nodejs_modules",
    package_json = "@rules_proto_grpc//nodejs:requirements/package.json",
    yarn_lock = "@rules_proto_grpc//nodejs:requirements/yarn.lock",
)`)

var nodeLibraryRuleTemplateString = `load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@build_bazel_rules_nodejs//:index.bzl", "js_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )
`

var nodeProtoLibraryRuleTemplate = mustTemplate(nodeLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    js_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = PROTO_DEPS,
        package_name = kwargs.get("name"),
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags"),
    )

PROTO_DEPS = [
    "@nodejs_modules//google-protobuf",
]`)

var nodeGrpcLibraryRuleTemplate = mustTemplate(nodeLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    js_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = GRPC_DEPS,
        package_name = kwargs.get("name"),
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags"),
    )

GRPC_DEPS = [
    "@nodejs_modules//google-protobuf",
    "@nodejs_modules//@grpc/grpc-js",
]`)

func makeNode() *Language {
	return &Language{
		Dir:   "nodejs",
		Name:  "nodejs",
		DisplayName: "Node.js",
		Notes: mustTemplate("Rules for generating Node.js protobuf and gRPC `.js` files using standard Protocol Buffers and gRPC."),
		Flags: commonLangFlags,
		SkipTestPlatforms: []string{},
		Rules: []*Rule{
			&Rule{
				Name:             "nodejs_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//nodejs:nodejs_plugin"},
				WorkspaceExample: nodeWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates Node.js protobuf `.js` artifacts",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{},
			},
			&Rule{
				Name:             "nodejs_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//nodejs:nodejs_plugin", "//nodejs:grpc_nodejs_plugin"},
				WorkspaceExample: nodeWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates Node.js protobuf+gRPC `.js` artifacts",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{},
			},
			&Rule{
				Name:             "nodejs_proto_library",
				Kind:             "proto",
				Implementation:   nodeProtoLibraryRuleTemplate,
				WorkspaceExample: nodeWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates a Node.js protobuf library using `js_library` from `rules_nodejs`",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{},
				Experimental:     true,
			},
			&Rule{
				Name:             "nodejs_grpc_library",
				Kind:             "grpc",
				Implementation:   nodeGrpcLibraryRuleTemplate,
				WorkspaceExample: nodeWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates a Node.js protobuf+gRPC library using `js_library` from `rules_nodejs`",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{},
				Experimental:     true,
			},
		},
	}
}
