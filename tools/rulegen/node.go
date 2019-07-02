package main

var nodeGrpcCompileWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()`)

var nodeProtoLibraryWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@build_bazel_rules_nodejs//:defs.bzl", "node_repositories", "yarn_install")

node_repositories()

yarn_install(
    name = "node_modules",
    package_json = "//node:requirements/package.json",
    yarn_lock = "//node:requirements/yarn.lock",
)`)

var nodeGrpcLibraryWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@build_bazel_rules_nodejs//:defs.bzl", "node_repositories")

yarn_install(
    name = "node_modules",
    package_json = "@build_stack_rules_proto//node:requirements/package.json",
    yarn_lock = "@build_stack_rules_proto//node:requirements/yarn.lock",
)`)

var nodeLibraryRuleTemplateString = `load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@build_bazel_rules_nodejs//:defs.bzl", "npm_package")

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
    npm_package(
        name = kwargs.get("name"),
        deps = [name_pb],
        packages = [
            "@node_modules//google-protobuf",
        ],
        visibility = kwargs.get("visibility"),
    )`)

var nodeGrpcLibraryRuleTemplate = mustTemplate(nodeLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    npm_package(
        name = kwargs.get("name"),
        deps = [name_pb],
        packages = [
            "@node_modules//google-protobuf",
            "@node_modules//grpc",
        ],
        visibility = kwargs.get("visibility"),
    )`)

func makeNode() *Language {
	return &Language{
		Dir:   "node",
		Name:  "node",
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:             "node_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//node:js"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates node *.js protobuf artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "node_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//node:js", "//node:grpc_js"},
				WorkspaceExample: nodeGrpcCompileWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates node *.js protobuf+gRPC artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
// 			&Rule{
// 				Name:             "node_proto_library",
// 				Kind:             "proto",
// 				Implementation:   nodeProtoLibraryRuleTemplate,
// 				WorkspaceExample: nodeProtoLibraryWorkspaceTemplate,
// 				BuildExample:     protoLibraryExampleTemplate,
// 				Doc:              "Generates node *.js protobuf library",
// 				Attrs:            aspectProtoCompileAttrs,
// 			},
// 			&Rule{
// 				Name:             "node_grpc_library",
// 				Kind:             "grpc",
// 				Implementation:   nodeGrpcLibraryRuleTemplate,
// 				WorkspaceExample: nodeGrpcLibraryWorkspaceTemplate,
// 				BuildExample:     grpcLibraryExampleTemplate,
// 				Doc:              "Generates node *.js protobuf+gRPC library",
// 				Attrs:            aspectProtoCompileAttrs,
// 			},
		},
	}
}
