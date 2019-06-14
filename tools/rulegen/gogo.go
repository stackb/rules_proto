package main

import "fmt"

var gogoLibraryRuleTemplateString = `load("//{{ .Lang.Dir }}:{{ .Rule.Base }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Rule.Base }}_{{ .Rule.Kind }}_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("//go:utils.bzl", "get_importmappings")

wkt_mappings = get_importmappings({
    "google/protobuf/any.proto": "github.com/gogo/protobuf/types",
    "google/protobuf/duration.proto": "github.com/gogo/protobuf/types",
    "google/protobuf/struct.proto": "github.com/gogo/protobuf/types",
    "google/protobuf/timestamp.proto": "github.com/gogo/protobuf/types",
    "google/protobuf/wrappers.proto": "github.com/gogo/protobuf/types",
})

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    importpath = kwargs.get("importpath")
    visibility = kwargs.get("visibility")
    go_deps = kwargs.get("go_deps", [])

    name_pb = name + "_pb"

    {{ .Rule.Base }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        deps = deps,
        plugin_options = get_importmappings(kwargs.pop("importmap", {})) + wkt_mappings,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )
`

var gogoProtoLibraryRuleTemplate = mustTemplate(gogoLibraryRuleTemplateString + `
    go_library(
        name = name,
        srcs = [name_pb],
        deps = go_deps + [
            "@com_github_gogo_protobuf//proto:go_default_library",
        ],
        importpath = importpath,
        visibility = visibility,
    )`)

var gogoGrpcLibraryRuleTemplate = mustTemplate(gogoLibraryRuleTemplateString + `
    go_library(
        name = name,
        srcs = [name_pb],
        deps = go_deps + [
            "@com_github_gogo_protobuf//proto:go_default_library",
            "@org_golang_google_grpc//:go_default_library",
            "@org_golang_x_net//context:go_default_library",
        ],
        importpath = importpath,
        visibility = visibility,
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
		Flags: append(commonLangFlags, &Flag{
			Category: "build",
			Name:     "incompatible_require_ctx_in_configure_features",
			Value:    "false",
		}),
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
			Name:           base + "_proto_compile",
			Base:           base,
			Kind:           "proto",
			Implementation: goCompileRuleTemplate,
			Plugins:        []string{protoPlugin},
			Usage:          goUsageTemplate,
			Example:        protoCompileExampleTemplate,
			Doc:            fmt.Sprintf("Generates %s protobuf artifacts", base),
			Attrs:          append(protoCompileAttrs, goProtoAttrs...),
		},
		&Rule{
			Name:           base + "_grpc_compile",
			Base:           base,
			Kind:           "grpc",
			Implementation: goCompileRuleTemplate,
			Plugins:        []string{grpcPlugin},
			Usage:          goUsageTemplate,
			Example:        grpcCompileExampleTemplate,
			Doc:            fmt.Sprintf("Generates %s protobuf+gRPC artifacts", base),
			Attrs:          append(protoCompileAttrs, goProtoAttrs...),
		},

		&Rule{
			Name:           base + "_proto_library",
			Base:           base,
			Kind:           "proto",
			Implementation: gogoProtoLibraryRuleTemplate,
			Usage:          goUsageTemplate,
			Example:        gogoProtoLibraryExampleTemplate,
			Doc:            fmt.Sprintf("Generates %s protobuf library", base),
			Attrs:          append(protoCompileAttrs, goProtoAttrs...),
		},
		&Rule{
			Name:           base + "_grpc_library",
			Base:           base,
			Kind:           "grpc",
			Implementation: gogoGrpcLibraryRuleTemplate,
			Usage:          goUsageTemplate,
			Example:        goGrpcLibraryExampleTemplate,
			Doc:            fmt.Sprintf("Generates %s protobuf+gRPC library", base),
			Attrs:          append(protoCompileAttrs, goProtoAttrs...),
		},
	)
}
