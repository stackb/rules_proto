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
)
`)

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
)

`)

var nodeProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:node_proto_compile.bzl", "node_proto_compile")
load("//node:node_module_index.bzl", "node_module_index")
load("@org_pubref_rules_node//node:rules.bzl", "node_module")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    name_index = name + "_index"

    node_proto_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )
    node_module_index(
        name = name_index,
        compilation = name_pb,
    )
    node_module(
        name = name,
        srcs = [name_pb],
        index = name_index,
        deps = [
            "@proto_node_modules//:_all_",
        ],
        visibility = visibility,
    )
`)

var nodeGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:node_grpc_compile.bzl", "node_grpc_compile")
load("//node:node_module_index.bzl", "node_module_index")
load("@org_pubref_rules_node//node:rules.bzl", "node_module")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    name_index = name + "_index"

    node_grpc_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )
    node_module_index(
        name = name_index,
        compilation = name_pb,
    )
    node_module(
        name = name,
        srcs = [name_pb],
        index = name_index,
        deps = [
            "@proto_node_modules//:_all_",
            "@grpc_node_modules//:_all_",
        ],
        visibility = visibility,
    )
`)

func makeNode() *Language {
	return &Language{
		Dir:  "node",
		Name: "node",
		Rules: []*Rule{
			&Rule{
				Name:           "node_proto_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//node:js"},
				Usage:          usageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates node *.js protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "node_grpc_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//node:js", "//node:grpc_js"},
				Usage:          nodeGrpcCompileUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates node *.js protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "node_proto_library",
				Implementation: nodeProtoLibraryRuleTemplate,
				Usage:          nodeProtoLibraryUsageTemplate,
				Example:        protoLibraryExampleTemplate,
				Doc:            "Generates node *.js protobuf library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "node_grpc_library",
				Implementation: nodeGrpcLibraryRuleTemplate,
				Usage:          nodeGrpcLibraryUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates node *.js protobuf+gRPC library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
