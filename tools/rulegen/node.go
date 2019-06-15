package main

var nodeGrpcCompileUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()`)

var nodeProtoLibraryUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@org_pubref_rules_node//node:rules.bzl", "node_repositories", "yarn_modules")

node_repositories()

yarn_modules(
    name = "proto_node_modules",
    deps = {
        "google-protobuf": "3.6.1",
    },
)`)

var nodeGrpcLibraryUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@org_pubref_rules_node//node:rules.bzl", "node_repositories", "yarn_modules")

node_repositories()

yarn_modules(
    name = "proto_node_modules",
    deps = {
        "google-protobuf": "3.6.1",
    },
)

yarn_modules(
    name = "grpc_node_modules",
    deps = {
        "grpc": "1.15.1",
    },
)`)

var nodeLibraryRuleTemplateString = `load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("//node:node_module_index.bzl", "node_module_index")
load("@org_pubref_rules_node//node:rules.bzl", "node_module")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    name_index = kwargs.get("name") + "_index"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create index
    node_module_index(
        name = name_index,
        compilation = name_pb,
    )
`

var nodeProtoLibraryRuleTemplate = mustTemplate(nodeLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    node_module(
        name = kwargs.get("name"),
        srcs = [name_pb],
        index = name_index,
        deps = [
            "@proto_node_modules//:_all_",
        ],
        visibility = kwargs.get("visibility"),
    )`)

var nodeGrpcLibraryRuleTemplate = mustTemplate(nodeLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    node_module(
        name = kwargs.get("name"),
        srcs = [name_pb],
        index = name_index,
        deps = [
            "@proto_node_modules//:_all_",
            "@grpc_node_modules//:_all_",
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
				Name:           "node_proto_compile",
				Kind:           "proto",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//node:js"},
				Usage:          usageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates node *.js protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "node_grpc_compile",
				Kind:           "grpc",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//node:js", "//node:grpc_js"},
				Usage:          nodeGrpcCompileUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates node *.js protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "node_proto_library",
				Kind:           "proto",
				Implementation: nodeProtoLibraryRuleTemplate,
				Usage:          nodeProtoLibraryUsageTemplate,
				Example:        protoLibraryExampleTemplate,
				Doc:            "Generates node *.js protobuf library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "node_grpc_library",
				Kind:           "grpc",
				Implementation: nodeGrpcLibraryRuleTemplate,
				Usage:          nodeGrpcLibraryUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates node *.js protobuf+gRPC library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
