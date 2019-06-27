package main

import "fmt"

var gogoLibraryRuleTemplateString = `load("//{{ .Lang.Dir }}:{{ .Rule.Base }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Rule.Base }}_{{ .Rule.Kind }}_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

def {{ .Rule.Name }}(deps, **kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Rule.Base }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        deps = deps, # Forward only deps
        prefix_path = kwargs.get("importpath", ""),
    )
`

var gogoProtoLibraryRuleTemplate = mustTemplate(gogoLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    go_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = kwargs.get("go_deps", []) + [
            "@com_github_gogo_protobuf//proto:go_default_library",
        ],
        importpath = kwargs.get("importpath"),
        visibility = kwargs.get("visibility"),
    )`)

var gogoGrpcLibraryRuleTemplate = mustTemplate(gogoLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    go_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = kwargs.get("go_deps", []) + [
            "@com_github_gogo_protobuf//proto:go_default_library",
            "@org_golang_google_grpc//:go_default_library",
            "@org_golang_x_net//context:go_default_library",
        ],
        importpath = kwargs.get("importpath"),
        visibility = kwargs.get("visibility"),
    )`)

var gogoProtoLibraryExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "person_{{ .Lang.Name }}_library",
    go_deps = [
        "@com_github_gogo_protobuf//types:go_default_library",
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
    importpath = "github.com/stackb/rules_proto/{{ .Lang.Name }}/example/{{ .Rule.Name }}/person",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)`)

func makeGogo() *Language {
	lang := &Language{
		Dir:     "github.com/gogo/protobuf",
		Name:    "gogo",
		Plugins: make(map[string]*Plugin),
		Flags: commonLangFlags,
	}

	addGogoRules(lang, "gogo")
	addGogoRules(lang, "gogofast")
	addGogoRules(lang, "gogofaster")

	// These are of dubious utility here and hence excluded
	//
	// addGogoRules(lang, "gogotypes")
	// addGogoRules(lang, "gogoslick")

	return lang
}

func addGogoRules(language *Language, base string) {
	protoPlugin := "//github.com/gogo/protobuf:" + base
	grpcPlugin := "//github.com/gogo/protobuf:grpc_" + base

	language.Plugins[protoPlugin] = &Plugin{
		Tool: "@com_github_gogo_protobuf//protoc-gen-" + base,
	}

	language.Plugins[grpcPlugin] = &Plugin{
		Tool:    "@com_github_gogo_protobuf//protoc-gen-" + base,
		Options: []string{"plugins=grpc"},
	}

	language.Rules = append(language.Rules,
		&Rule{
			Name:             base + "_proto_compile",
			Base:             base,
			Kind:             "proto",
			Implementation:   aspectRuleTemplate,
			Plugins:          []string{protoPlugin},
			WorkspaceExample: goWorkspaceTemplate,
			BuildExample:     protoCompileExampleTemplate,
			Doc:              fmt.Sprintf("Generates %s protobuf artifacts", base),
			Attrs:            aspectProtoCompileAttrs,
		},
		&Rule{
			Name:             base + "_grpc_compile",
			Base:             base,
			Kind:             "grpc",
			Implementation:   aspectRuleTemplate,
			Plugins:          []string{grpcPlugin},
			WorkspaceExample: goWorkspaceTemplate,
			BuildExample:     grpcCompileExampleTemplate,
			Doc:              fmt.Sprintf("Generates %s protobuf+gRPC artifacts", base),
			Attrs:            aspectProtoCompileAttrs,
		},

		&Rule{
			Name:             base + "_proto_library",
			Base:             base,
			Kind:             "proto",
			Implementation:   gogoProtoLibraryRuleTemplate,
			WorkspaceExample: goWorkspaceTemplate,
			BuildExample:     gogoProtoLibraryExampleTemplate,
			Doc:              fmt.Sprintf("Generates %s protobuf library", base),
			Attrs:            append(aspectProtoCompileAttrs, goProtoAttrs...),
		},
		&Rule{
			Name:             base + "_grpc_library",
			Base:             base,
			Kind:             "grpc",
			Implementation:   gogoGrpcLibraryRuleTemplate,
			WorkspaceExample: goWorkspaceTemplate,
			BuildExample:     goGrpcLibraryExampleTemplate,
			Doc:              fmt.Sprintf("Generates %s protobuf+gRPC library", base),
			Attrs:            append(aspectProtoCompileAttrs, goProtoAttrs...),
		},
	)
}
