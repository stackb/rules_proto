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
				Name:             "objc_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//objc:objc"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates objc protobuf artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "objc_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//objc:objc", "//objc:grpc_objc"},
				WorkspaceExample: grpcWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates objc protobuf+gRPC artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			// &Rule{
			// 	Name:             "objc_proto_library",
			//  Kind:             "proto",
			// 	Implementation:   objcLibraryRuleTemplate,
			// 	WorkspaceExample: protoWorkspaceTemplate,
			// 	BuildExample:     protoLibraryExampleTemplate,
			// 	Doc:              "Generates objc protobuf library",
			// 	Attrs:            aspectProtoCompileAttrs,
			// },
			// &Rule{
			// 	Name:             "objc_grpc_library",
			//  Kind:             "grpc",
			// 	Implementation:   objcLibraryRuleTemplate,
			// 	WorkspaceExample: grpcWorkspaceTemplate,
			// 	BuildExample:     grpcLibraryExampleTemplate,
			// 	Doc:              "Generates objc protobuf+gRPC library",
			// 	Attrs:            aspectProtoCompileAttrs,
			// },
		},
	}
}
