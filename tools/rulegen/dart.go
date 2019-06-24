package main

var dartWorkspaceTemplateString = `load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

# rules_go used here to compile a wrapper around the protoc-gen-grpc plugin
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@io_bazel_rules_dart//dart/build_rules:repositories.bzl", "dart_repositories")

dart_repositories()

load("@dart_pub_deps_protoc_plugin//:deps.bzl", dart_protoc_plugin_deps = "pub_deps")

dart_protoc_plugin_deps()`

var dartProtoWorkspaceTemplate = mustTemplate(dartWorkspaceTemplateString)

var dartGrpcWorkspaceTemplate = mustTemplate(dartWorkspaceTemplateString + `

load("@dart_pub_deps_grpc//:deps.bzl", dart_grpc_deps = "pub_deps")

dart_grpc_deps()`)

var dartLibraryRuleTemplateString = `load("//{{ .Lang.Dir}}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@io_bazel_rules_dart//dart/build_rules:core.bzl", "dart_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )
`

var dartProtoLibraryRuleTemplate = mustTemplate(dartLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    dart_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            Label("@vendor_protobuf//:protobuf"),
        ],
        pub_pkg_name = kwargs.get("name"),
        visibility = kwargs.get("visibility"),
    )`)

var dartGrpcLibraryRuleTemplate = mustTemplate(dartLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    dart_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            Label("@vendor_protobuf//:protobuf"),
            Label("@vendor_grpc//:grpc"),
        ],
        pub_pkg_name = kwargs.get("name"),
        visibility = kwargs.get("visibility"),
    )`)

func makeDart() *Language {
	return &Language{
		Dir:   "dart",
		Name:  "dart",
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:             "dart_proto_compile",
				Kind:             "proto",
				Implementation:   compileRuleTemplate,
				Plugins:          []string{"//dart:dart"},
				WorkspaceExample: dartProtoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates dart protobuf artifacts",
				Attrs:            append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:             "dart_grpc_compile",
				Kind:             "grpc",
				Implementation:   compileRuleTemplate,
				Plugins:          []string{"//dart:grpc_dart"},
				WorkspaceExample: dartProtoWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates dart protobuf+gRPC artifacts",
				Attrs:            append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:             "dart_proto_library",
				Kind:             "proto",
				Implementation: dartProtoLibraryRuleTemplate,
				WorkspaceExample: dartProtoWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates dart protobuf library",
				Attrs:            append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:             "dart_grpc_library",
				Kind:             "grpc",
				Implementation:   dartGrpcLibraryRuleTemplate,
				WorkspaceExample: dartGrpcWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates dart protobuf+gRPC library",
				Attrs:            append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
