package main

var objcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

	{{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
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
				Kind:           "proto",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//objc:objc"},
				Usage:          usageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates objc protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "objc_grpc_compile",
				Kind:           "grpc",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//objc:objc", "//objc:grpc_objc"},
				Usage:          grpcUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates objc protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			// &Rule{
			// 	Name:           "objc_proto_library",
			//  Kind:           "proto",
			// 	Implementation: objcLibraryRuleTemplate,
			// 	Usage:          usageTemplate,
			// 	Example:        protoLibraryExampleTemplate,
			// 	Doc:            "Generates objc protobuf library",
			// 	Attrs:          append(protoCompileAttrs, []*Attr{}...),
			// },
			// &Rule{
			// 	Name:           "objc_grpc_library",
			//  Kind:           "grpc",
			// 	Implementation: objcLibraryRuleTemplate,
			// 	Usage:          grpcUsageTemplate,
			// 	Example:        grpcLibraryExampleTemplate,
			// 	Doc:            "Generates objc protobuf+gRPC library",
			// 	Attrs:          append(protoCompileAttrs, []*Attr{}...),
			// },
		},
	}
}
