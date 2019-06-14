package main

var objcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create {{ .Lang.Name }} library
    native.objc_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        includes = [name_pb],
        visibility = kwargs.get("visibility"),
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
