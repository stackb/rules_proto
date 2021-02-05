package main

var goWorkspaceTemplate = mustTemplate(`load("@rules_proto_grpc//:repositories.bzl", "bazel_gazelle", "io_bazel_rules_go")

io_bazel_rules_go()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

bazel_gazelle()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

load("@rules_proto_grpc//{{ .Lang.Dir }}:repositories.bzl", rules_proto_grpc_{{ .Lang.Name }}_repos="{{ .Lang.Name }}_repos")

rules_proto_grpc_{{ .Lang.Name }}_repos()`)


var goLibraryRuleTemplateString = `load("//{{ .Lang.Dir }}:{{ .Rule.Base}}_{{ .Rule.Kind }}_compile.bzl", "{{ .Rule.Base }}_{{ .Rule.Kind }}_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Rule.Base}}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        prefix_path = kwargs.get("importpath", ""),
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )
`

var goProtoLibraryRuleTemplate = mustTemplate(goLibraryRuleTemplateString + `
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
    "@com_github_golang_protobuf//proto:go_default_library",
    "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
    "@org_golang_google_protobuf//runtime/protoimpl:go_default_library",
]`)

var goGrpcLibraryRuleTemplate = mustTemplate(goLibraryRuleTemplateString + `
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
    "@com_github_golang_protobuf//proto:go_default_library",
    "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
    "@org_golang_google_protobuf//runtime/protoimpl:go_default_library",
    "@org_golang_google_grpc//:go_default_library",
    "@org_golang_google_grpc//codes:go_default_library",
    "@org_golang_google_grpc//status:go_default_library",
    "@org_golang_x_net//context:go_default_library",
]`)

var goProtoCompileExampleTemplate = mustTemplate(`load("@rules_proto_grpc//{{ .Lang.Dir }}:defs.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "person_{{ .Lang.Name }}_proto",
    importpath = "github.com/rules-proto-grpc/rules_proto_grpc/{{ .Lang.Name }}/example/{{ .Rule.Name }}/person",
    deps = ["@rules_proto_grpc//example/proto:person_proto"],
)`)

var goGrpcCompileExampleTemplate = mustTemplate(`load("@rules_proto_grpc//{{ .Lang.Dir }}:defs.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "greeter_{{ .Lang.Name }}_grpc",
    importpath = "github.com/rules-proto-grpc/rules_proto_grpc/{{ .Lang.Name }}/example/{{ .Rule.Name }}/greeter",
    deps = ["@rules_proto_grpc//example/proto:greeter_grpc"],
)`)

var goProtoLibraryExampleTemplate = mustTemplate(`load("@rules_proto_grpc//{{ .Lang.Dir }}:defs.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "person_{{ .Lang.Name }}_library",
    go_deps = [
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
    importpath = "github.com/rules-proto-grpc/rules_proto_grpc/{{ .Lang.Name }}/example/{{ .Rule.Name }}/person",
    deps = ["@rules_proto_grpc//example/proto:person_proto"],
)`)

var goGrpcLibraryExampleTemplate = mustTemplate(`load("@rules_proto_grpc//{{ .Lang.Dir }}:defs.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "greeter_{{ .Lang.Name }}_library",
    go_deps = [
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
    importpath = "github.com/rules-proto-grpc/rules_proto_grpc/{{ .Lang.Name }}/example/{{ .Rule.Name }}/greeter",
    deps = ["@rules_proto_grpc//example/proto:greeter_grpc"],
)`)

var goProtoAttrs = []*Attr{
	&Attr{
		Name:      "importpath",
		Type:      "string",
		Default:   "None",
		Doc:       "Importpath for the generated artifacts",
		Mandatory: false,
	},
}

func makeGo() *Language {
	return &Language{
		Dir:  "go",
		Name: "go",
		DisplayName: "Go",
		Notes: mustTemplate("Rules for generating Go protobuf and gRPC `.go` files and libraries using [golang/protobuf](https://github.com/golang/protobuf). Libraries are created with `go_library` from [rules_go](https://github.com/bazelbuild/rules_go)"),
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:             "go_proto_compile",
				Base:             "go",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//go:go_plugin"},
				WorkspaceExample: goWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates Go protobuf `.go` artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "go_grpc_compile",
				Base:             "go",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//go:grpc_go_plugin"},
				WorkspaceExample: goWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates Go protobuf+gRPC `.go` artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "go_proto_library",
				Base:             "go",
				Kind:             "proto",
				Implementation:   goProtoLibraryRuleTemplate,
				WorkspaceExample: goWorkspaceTemplate,
				BuildExample:     goProtoLibraryExampleTemplate,
				Doc:              "Generates a Go protobuf library using `go_library` from `rules_go`",
				Attrs:            append(aspectProtoCompileAttrs, goProtoAttrs...),
			},
			&Rule{
				Name:             "go_grpc_library",
				Base:             "go",
				Kind:             "grpc",
				Implementation:   goGrpcLibraryRuleTemplate,
				WorkspaceExample: goWorkspaceTemplate,
				BuildExample:     goGrpcLibraryExampleTemplate,
				Doc:              "Generates a Go protobuf+gRPC library using `go_library` from `rules_go`",
				Attrs:            append(aspectProtoCompileAttrs, goProtoAttrs...),
			},
		},
	}
}
