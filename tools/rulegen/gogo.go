package main

import "fmt"

var gogoLibraryRuleTemplateString = `load("//{{ .Lang.Dir }}:{{ .Rule.Base }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Rule.Base }}_{{ .Rule.Kind }}_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Rule.Base }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        prefix_path = kwargs.get("importpath", ""),
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )
`

var gogoProtoLibraryRuleTemplate = mustTemplate(gogoLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    go_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = kwargs.get("go_deps", []) + PROTO_DEPS,
        importpath = kwargs.get("importpath"),
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags"),
    )

PROTO_DEPS = [
    "@com_github_gogo_protobuf//proto:go_default_library",
]`)

var gogoGrpcLibraryRuleTemplate = mustTemplate(gogoLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    go_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = kwargs.get("go_deps", []) + GRPC_DEPS,
        importpath = kwargs.get("importpath"),
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags"),
    )

GRPC_DEPS = [
    "@com_github_gogo_protobuf//proto:go_default_library",
    "@org_golang_google_grpc//:go_default_library",
    "@org_golang_google_grpc//codes:go_default_library",
    "@org_golang_google_grpc//status:go_default_library",
    "@org_golang_x_net//context:go_default_library",
]`)

var gogoProtoLibraryExampleTemplate = mustTemplate(`load("@rules_proto_grpc//{{ .Lang.Dir }}:defs.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "person_{{ .Lang.Name }}_library",
    go_deps = [
        "@com_github_gogo_protobuf//types:go_default_library",
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
    importpath = "github.com/rules-proto-grpc/rules_proto_grpc/{{ .Lang.Name }}/example/{{ .Rule.Name }}/person",
    deps = ["@rules_proto_grpc//example/proto:person_proto"],
)`)

func makeGogo() *Language {
	lang := &Language{
		Dir:     "github.com/gogo/protobuf",
		Name:    "gogo",
		DisplayName: "Go (gogoprotobuf)",
		Notes: mustTemplate("Rules for generating Go protobuf and gRPC `.go` files and libraries using [gogo/protobuf](https://github.com/gogo/protobuf). Libraries are created with `go_library` from [rules_go](https://github.com/bazelbuild/rules_go)"),
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
	protoPlugin := "//github.com/gogo/protobuf:" + base + "_plugin"
	grpcPlugin := "//github.com/gogo/protobuf:grpc_" + base + "_plugin"

	language.Rules = append(language.Rules,
		&Rule{
			Name:             base + "_proto_compile",
			Base:             base,
			Kind:             "proto",
			Implementation:   aspectRuleTemplate,
			Plugins:          []string{protoPlugin},
			WorkspaceExample: goWorkspaceTemplate,
			BuildExample:     protoCompileExampleTemplate,
			Doc:              fmt.Sprintf("Generates %s protobuf `.go` artifacts", base),
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
			Doc:              fmt.Sprintf("Generates %s protobuf+gRPC `.go` artifacts", base),
			Attrs:            aspectProtoCompileAttrs,
		},

		&Rule{
			Name:             base + "_proto_library",
			Base:             base,
			Kind:             "proto",
			Implementation:   gogoProtoLibraryRuleTemplate,
			WorkspaceExample: goWorkspaceTemplate,
			BuildExample:     gogoProtoLibraryExampleTemplate,
			Doc:              fmt.Sprintf("Generates a Go %s protobuf library using `go_library` from `rules_go`", base),
			Attrs:            append(aspectProtoCompileAttrs, goProtoAttrs...),
		},
		&Rule{
			Name:             base + "_grpc_library",
			Base:             base,
			Kind:             "grpc",
			Implementation:   gogoGrpcLibraryRuleTemplate,
			WorkspaceExample: goWorkspaceTemplate,
			BuildExample:     goGrpcLibraryExampleTemplate,
			Doc:              fmt.Sprintf("Generates a Go %s protobuf+gRPC library using `go_library` from `rules_go`", base),
			Attrs:            append(aspectProtoCompileAttrs, goProtoAttrs...),
			SkipTestPlatforms: []string{"windows"}, // gRPC go lib rules fail on windows due to bad path
		},
	)
}
