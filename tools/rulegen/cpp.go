package main

var cppProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:cpp_proto_compile.bzl", "cpp_proto_compile")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    cpp_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    native.cc_library(
        name = name,
        srcs = [name_pb],
        deps = [
            "//external:protobuf_clib",
        ],
        includes = [name_pb],
        visibility = visibility,
    )`)

var cppGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:cpp_grpc_compile.bzl", "cpp_grpc_compile")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    cpp_grpc_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    native.cc_library(
        name = name,
        srcs = [name_pb],
        deps = [
            "//external:protobuf_clib",
            "@com_github_grpc_grpc//:grpc++",
            "@com_github_grpc_grpc//:grpc++_reflection",
        ],
        # This seems magical to me.
        includes = [name_pb],
        visibility = visibility,
    )`)

var commonLangFlags = []*Flag{
	{
		Category: "build",
		Name:     "incompatible_enable_cc_toolchain_resolution",
		Value:    "false",
	},
}

func makeCpp() *Language {
	return &Language{
		Dir:   "cpp",
		Name:  "cpp",
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:           "cpp_proto_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//cpp:cpp"},
				Usage:          usageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates *.h,*.cc protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "cpp_grpc_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//cpp:cpp", "//cpp:grpc_cpp"},
				Usage:          grpcUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates *.h,*.cc protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "cpp_proto_library",
				Implementation: cppProtoLibraryRuleTemplate,
				Usage:          usageTemplate,
				Example:        protoLibraryExampleTemplate,
				Doc:            "Generates *.h,*.cc protobuf library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "cpp_grpc_library",
				Implementation: cppGrpcLibraryRuleTemplate,
				Usage:          grpcUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates *.h,*.cc protobuf+gRPC library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
				Flags: []*Flag{
					{
						Category:    "build",
						Name:        "incompatible_disable_legacy_proto_provider",
						Value:       "false",
						Description: "grpc/gprc has not migrated to ProtoInfo provider",
					},
					{
						Category:    "build",
						Name:        "incompatible_depset_is_not_iterable",
						Value:       "false",
						Description: "com_github_grpc_grpc/bazel/generate_cc.bzl line 10, in generate_cc_impl",
					},
					{
						Category:    "build",
						Name:        "incompatible_disallow_struct_provider_syntax",
						Value:       "false",
						Description: "com_github_grpc_grpc/bazel/generate_cc.bzl: 81",
					},
				},
			},
		},
	}
}
