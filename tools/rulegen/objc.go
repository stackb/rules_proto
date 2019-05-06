package main

var objcProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:objc_proto_compile.bzl", "objc_proto_compile")
def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

	objc_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    native.objc_library(
        name = name,
        srcs = [name_pb],
        includes = [name_pb],
        visibility = visibility,
    )
`)

var objcGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:objc_grpc_compile.bzl", "objc_grpc_compile")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    objc_grpc_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    native.objc_library(
        name = name,
        srcs = [name_pb],
        includes = [name_pb],
        visibility = visibility,
    )
`)

func makeObjc() *Language {
	return &Language{
		Dir:   "objc",
		Name:  "objc",
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:           "objc_proto_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//objc:objc"},
				Usage:          usageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates objc protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "objc_grpc_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//objc:objc", "//objc:grpc_objc"},
				Usage:          grpcUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates objc protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			// &Rule{
			// 	Name:           "objc_proto_library",
			// 	Implementation: objcProtoLibraryRuleTemplate,
			// 	Usage:          usageTemplate,
			// 	Example:        protoLibraryExampleTemplate,
			// 	Doc:            "Generates objc protobuf library",
			// 	Attrs:          append(protoCompileAttrs, []*Attr{}...),
			// },
			// &Rule{
			// 	Name:           "objc_grpc_library",
			// 	Implementation: objcGrpcLibraryRuleTemplate,
			// 	Usage:          grpcUsageTemplate,
			// 	Example:        grpcLibraryExampleTemplate,
			// 	Doc:            "Generates objc protobuf+gRPC library",
			// 	Attrs:          append(protoCompileAttrs, []*Attr{}...),
			// },
		},
	}
}
