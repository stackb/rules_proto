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
        transitive = True,
    )

    native.cc_library(
        name = name,
        srcs = [name_pb],
        deps = [
            "//external:protobuf_clib",
        ],
        includes = [name_pb],
        visibility = visibility,
    )
`)

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
        transitive = True,
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
    )
`)

func makeCpp() *Language {
	return &Language{
		Dir:  "cpp",
		Name: "cpp",
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
			},
		},
	}
}
