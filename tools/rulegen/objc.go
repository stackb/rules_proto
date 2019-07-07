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
        deps = [
            "@com_google_protobuf//:protobuf_objc",
        ],
        includes = [name_pb],
        visibility = kwargs.get("visibility"),
    )
`)

var objcGrpcLibraryRuleTemplate = mustTemplate(objcLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    native.objc_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            "@com_google_protobuf//:protobuf_objc",
            "@com_github_grpc_grpc//:grpc++",
            "@build_stack_rules_proto//objc:grpc_lib",
        ],
        includes = [name_pb],
        visibility = kwargs.get("visibility"),
    )
`)

func makeObjc() *Language {
	return &Language{
		Dir:   "objc",
		Name:  "objc",
		Flags: commonLangFlags,
		BazelCIExcludePlatforms: []string{"ubuntu1604", "ubuntu1804", "windows"},
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
			&Rule{
				Name:             "objc_proto_library",
				Kind:             "proto",
				Implementation:   objcProtoLibraryRuleTemplate,
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates objc protobuf library",
				Attrs:            aspectProtoCompileAttrs,
			},
// 			&Rule{ // Disabled due to issues fetching gRPC dependencies
// 				Name:             "objc_grpc_library",
// 				Kind:             "grpc",
// 				Implementation:   objcGrpcLibraryRuleTemplate,
// 				WorkspaceExample: grpcWorkspaceTemplate,
// 				BuildExample:     grpcLibraryExampleTemplate,
// 				Doc:              "Generates objc protobuf+gRPC library",
// 				Attrs:            aspectProtoCompileAttrs,
// 			},
		},
	}
}
