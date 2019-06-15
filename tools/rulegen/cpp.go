package main

var cppLibraryRuleTemplateString = `load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )
`

var cppProtoLibraryRuleTemplate = mustTemplate(cppLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    native.cc_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            "//external:protobuf_clib",
        ],
        includes = [name_pb],
        visibility = kwargs.get("visibility"),
    )`)

var cppGrpcLibraryRuleTemplate = mustTemplate(cppLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    native.cc_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            "//external:protobuf_clib",
            "@com_github_grpc_grpc//:grpc++",
            "@com_github_grpc_grpc//:grpc++_reflection",
        ],
        # This seems magical to me.
        includes = [name_pb],
        visibility = kwargs.get("visibility"),
    )`)

func makeCpp() *Language {
	return &Language{
		Dir:   "cpp",
		Name:  "cpp",
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:           "cpp_proto_compile",
				Kind:           "proto",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//cpp:cpp"},
				Usage:          usageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates *.h,*.cc protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "cpp_grpc_compile",
				Kind:           "grpc",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//cpp:cpp", "//cpp:grpc_cpp"},
				Usage:          grpcUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates *.h,*.cc protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "cpp_proto_library",
				Kind:           "proto",
				Implementation: cppProtoLibraryRuleTemplate,
				Usage:          usageTemplate,
				Example:        protoLibraryExampleTemplate,
				Doc:            "Generates *.h,*.cc protobuf library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "cpp_grpc_library",
				Kind:           "grpc",
				Implementation: cppGrpcLibraryRuleTemplate,
				Usage:          grpcUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates *.h,*.cc protobuf+gRPC library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
