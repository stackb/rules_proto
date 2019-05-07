package main

var goUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

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
        name = kwargs.get("name")
        name_plugin = name + "_plugin"
        proto_plugin(
            name = name_plugin,
            outputs = ["{package}/%s/{basename}.pb.go" % importpath],
            tool = "{{ with (index .Lang.Plugins (index .Rule.Plugins 0)) }}{{ .Tool }}{{ end }}",
        )
        kwargs["plugins"] = [name_plugin]
        kwargs.pop("importpath")

    # Define the default plugin if still not defined
    if not kwargs.get("plugins"):
        kwargs["plugins"] = [str(Label("{{ index .Rule.Plugins 0 }}"))]

    proto_compile(
        **kwargs
    )`)

var goProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Rule.Base}}_{{ .Rule.Kind }}_compile.bzl", "{{ .Rule.Base }}_{{ .Rule.Kind }}_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("//go:utils.bzl", "get_importmappings")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    importpath = kwargs.get("importpath")
    visibility = kwargs.get("visibility")
    go_deps = kwargs.get("go_deps", [])

    name_pb = name + "_pb"

    {{ .Rule.Base}}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        deps = deps,
        plugin_options = get_importmappings(kwargs.pop("importmap", {})),
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    go_library(
        name = name,
        srcs = [name_pb],
        deps = go_deps + [
            "@com_github_golang_protobuf//proto:go_default_library",
        ],
        importpath = importpath,
        visibility = visibility,
    )`)

var goGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:go_grpc_compile.bzl", "go_grpc_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("//go:utils.bzl", "get_importmappings")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    importpath = kwargs.get("importpath")
    visibility = kwargs.get("visibility")
    go_deps = kwargs.get("go_deps", [])

    name_pb = name + "_pb"

    go_grpc_compile(
        name = name_pb,
        deps = deps,
        plugin_options = get_importmappings(kwargs.pop("importmap", {})),
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    go_library(
        name = name,
        srcs = [name_pb],
        deps = go_deps + [
            "@com_github_golang_protobuf//proto:go_default_library",
            "@org_golang_google_grpc//:go_default_library",
            "@org_golang_x_net//context:go_default_library",
        ],
        importpath = importpath,
        visibility = visibility,
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
		// incompatible_disable_legacy_flags_cc_toolchain_api
		// incompatible_disable_tools_defaults_package
		// incompatible_enable_cc_toolchain_resolution
		//
		// incompatible_enable_legacy_cpp_toolchain_skylark_api
		Flags: append(commonLangFlags, &Flag{
			Category:    "build",
			Name:        "incompatible_enable_cc_toolchain_resolution",
			Value:       "false",
			Description: "In order to use find_cpp_toolchain, you must include the '@bazel_tools//tools/cpp:toolchain_type' in the toolchains argument to your rule.",
		}, &Flag{
			Category:    "build",
			Name:        "incompatible_require_ctx_in_configure_features",
			Value:       "false",
			Description: `/external/io_bazel_rules_rust/rust/private/rustc.bzl", line 143, in _get_linker_and_args cc_common.configure_features(cc_toolchain = cc_toolchain, reque..., ...) Incompatible flag --incompatible_require_ctx_in_configure_features has been flipped, and the mandatory parameter 'ctx' of cc_common.configure_features is missing`,
		}),
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
				Name:           "go_proto_compile",
				Base:           "go",
				Kind:           "proto",
				Implementation: goCompileRuleTemplate,
				Plugins:        []string{"//go:go"},
				Usage:          goUsageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates *.go protobuf artifacts",
				Attrs:          append(protoCompileAttrs, goProtoAttrs...),
			},
			&Rule{
				Name:           "go_grpc_compile",
				Base:           "go",
				Kind:           "grpc",
				Implementation: goCompileRuleTemplate,
				Plugins:        []string{"//go:grpc_go"},
				Usage:          goUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates *.go protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, goProtoAttrs...),
			},
			&Rule{
				Name:           "go_proto_library",
				Base:           "go",
				Kind:           "proto",
				Implementation: goProtoLibraryRuleTemplate,
				Usage:          goUsageTemplate,
				Example:        goProtoLibraryExampleTemplate,
				Doc:            "Generates *.go protobuf library",
				Attrs:          append(protoCompileAttrs, goProtoAttrs...),
			},
			&Rule{
				Name:           "go_grpc_library",
				Base:           "go",
				Kind:           "grpc",
				Implementation: goGrpcLibraryRuleTemplate,
				Usage:          goUsageTemplate,
				Example:        goGrpcLibraryExampleTemplate,
				Doc:            "Generates *.go protobuf+gRPC library",
				Attrs:          append(protoCompileAttrs, goProtoAttrs...),
			},
		},
	}
}
