package main

var goWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()`)

var goCompileRuleTemplate = mustTemplate(`load("//:compile.bzl", "proto_compile")
load("//:plugin.bzl", "proto_plugin")

def {{ .Rule.Name }}(**kwargs):
    # If importpath specified, declare a custom plugin that should correctly
    # predict the output location.
    importpath = kwargs.get("importpath")
    if importpath and not kwargs.get("plugins"):
        name_plugin = kwargs.get("name") + "_plugin"
        proto_plugin(
            name = name_plugin,
            outputs = ["{package}/%s/{basename}.pb.go" % importpath],
            tool = "{{ with (index .Lang.Plugins (index .Rule.Plugins 0)) }}{{ .Tool }}{{ end }}",
        )
        kwargs["plugins"] = [name_plugin]
        kwargs.pop("importpath")

    # Define the default plugin if still not defined
    if not kwargs.get("plugins"):
        kwargs["plugins"] = [Label("{{ index .Rule.Plugins 0 }}")]

    proto_compile(
        **kwargs
    )`)

var goLibraryRuleTemplateString = `load("//{{ .Lang.Dir }}:{{ .Rule.Base}}_{{ .Rule.Kind }}_compile.bzl", "{{ .Rule.Base }}_{{ .Rule.Kind }}_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("//go:utils.bzl", "get_importmappings")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    kwargs["plugin_options"] = kwargs.get("plugin_options", []) + get_importmappings(kwargs.get("importmap", {}))
    {{ .Rule.Base}}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k not in ("name", "importpath", "importmap", "go_deps")} # Forward args except name, importpath, importmap and go_deps
    )
`

var goProtoLibraryRuleTemplate = mustTemplate(goLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    go_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = kwargs.get("go_deps", []) + [
            "@com_github_golang_protobuf//proto:go_default_library",
        ],
        importpath = kwargs.get("importpath"),
        visibility = kwargs.get("visibility"),
    )`)

var goGrpcLibraryRuleTemplate = mustTemplate(goLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    go_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = kwargs.get("go_deps", []) + [
            "@com_github_golang_protobuf//proto:go_default_library",
            "@org_golang_google_grpc//:go_default_library",
            "@org_golang_x_net//context:go_default_library",
        ],
        importpath = kwargs.get("importpath"),
        visibility = kwargs.get("visibility"),
    )`)

var goProtoCompileExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "person_{{ .Lang.Name }}_proto",
    importpath = "github.com/stackb/rules_proto/{{ .Lang.Name }}/example/{{ .Rule.Name }}/person",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)`)

var goGrpcCompileExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "greeter_{{ .Lang.Name }}_grpc",
    importpath = "github.com/stackb/rules_proto/{{ .Lang.Name }}/example/{{ .Rule.Name }}/greeter",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)`)

var goProtoLibraryExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "person_{{ .Lang.Name }}_library",
    go_deps = [
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
    importpath = "github.com/stackb/rules_proto/{{ .Lang.Name }}/example/{{ .Rule.Name }}/person",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)`)

var goGrpcLibraryExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "greeter_{{ .Lang.Name }}_library",
    go_deps = [
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
    importpath = "github.com/stackb/rules_proto/{{ .Lang.Name }}/example/{{ .Rule.Name }}/greeter",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)`)

var goProtoAttrs = []*Attr{
	&Attr{
		Name:      "importpath",
		Type:      "string",
		Default:   "None",
		Doc:       "Importpath for the generated artifacts",
		Mandatory: false,
	},
	&Attr{
		Name:      "importmap",
		Type:      "string_dict",
		Default:   "None",
		Doc:       "A dictionary of the form `{ K: V}` that dictates the importpath `V` for a matching imported proto file `K`",
		Mandatory: false,
	},
}

func makeGo() *Language {
	return &Language{
		Dir:  "go",
		Name: "go",
		Flags: commonLangFlags,
		Plugins: map[string]*Plugin{
			"//go:go": &Plugin{
				Tool: "@com_github_golang_protobuf//protoc-gen-go",
			},
			"//go:grpc_go": &Plugin{
				Tool:    "@com_github_golang_protobuf//protoc-gen-go",
				Options: []string{"plugins=grpc"},
			},
		},
		Rules: []*Rule{
			&Rule{
				Name:             "go_proto_compile",
				Base:             "go",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//go:go"},
				WorkspaceExample: goWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates *.go protobuf artifacts",
				Attrs:            append(aspectProtoCompileAttrs, goProtoAttrs...),
			},
			&Rule{
				Name:             "go_grpc_compile",
				Base:             "go",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//go:grpc_go"},
				WorkspaceExample: goWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates *.go protobuf+gRPC artifacts",
				Attrs:            append(aspectProtoCompileAttrs, goProtoAttrs...),
			},
			&Rule{
				Name:             "go_proto_library",
				Base:             "go",
				Kind:             "proto",
				Implementation:   goProtoLibraryRuleTemplate,
				WorkspaceExample: goWorkspaceTemplate,
				BuildExample:     goProtoLibraryExampleTemplate,
				Doc:              "Generates *.go protobuf library",
				Attrs:            append(aspectProtoCompileAttrs, goProtoAttrs...),
			},
			&Rule{
				Name:             "go_grpc_library",
				Base:             "go",
				Kind:             "grpc",
				Implementation:   goGrpcLibraryRuleTemplate,
				WorkspaceExample: goWorkspaceTemplate,
				BuildExample:     goGrpcLibraryExampleTemplate,
				Doc:              "Generates *.go protobuf+gRPC library",
				Attrs:            append(aspectProtoCompileAttrs, goProtoAttrs...),
			},
		},
	}
}
