package main

var objcLibraryRuleTemplateString = `load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )
`

var objcProtoLibraryRuleTemplate = mustTemplate(objcLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    native.objc_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = PROTO_DEPS,
        includes = [name_pb],
        copts = kwargs.get("copts"),
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags"),
    )

PROTO_DEPS = [
    "@com_google_protobuf//:protobuf_objc",
]`)

var objcGrpcLibraryRuleTemplate = mustTemplate(objcLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    native.objc_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = GRPC_DEPS,
        includes = [name_pb],
        copts = kwargs.get("copts"),
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags"),
    )

GRPC_DEPS = [
    "@com_google_protobuf//:protobuf_objc",
    "@com_github_grpc_grpc//:grpc++",
    "@rules_proto_grpc//objc:grpc_lib",
]`)

func makeObjc() *Language {
	return &Language{
		Dir:   "objc",
		Name:  "objc",
		DisplayName: "Objective-C",
		Notes: mustTemplate("Rules for generating Objective-C protobuf and gRPC `.m` & `.h` files and libraries using standard Protocol Buffers and gRPC. Libraries are created with the Bazel native `objc_library`"),
		Flags: commonLangFlags,
		SkipTestPlatforms: []string{"linux", "windows"},
		Rules: []*Rule{
			&Rule{
				Name:             "objc_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//objc:objc_plugin"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates Objective-C protobuf `.m` & `.h` artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "objc_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//objc:objc_plugin", "//objc:grpc_objc_plugin"},
				WorkspaceExample: grpcWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates Objective-C protobuf+gRPC `.m` & `.h` artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "objc_proto_library",
				Kind:             "proto",
				Implementation:   objcProtoLibraryRuleTemplate,
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates an Objective-C protobuf library using `objc_library`",
				Attrs:            aspectProtoCompileAttrs,
			},
// 			&Rule{ // Disabled due to issues fetching gRPC dependencies
// 				Name:             "objc_grpc_library",
// 				Kind:             "grpc",
// 				Implementation:   objcGrpcLibraryRuleTemplate,
// 				WorkspaceExample: grpcWorkspaceTemplate,
// 				BuildExample:     grpcLibraryExampleTemplate,
// 				Doc:              "Generates an Objective-C protobuf+gRPC library using `objc_library`",
// 				Attrs:            aspectProtoCompileAttrs,
// 			},
		},
	}
}
