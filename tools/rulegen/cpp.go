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
            "@com_google_protobuf//:protoc_lib",
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
            "@com_google_protobuf//:protoc_lib",
            "@com_github_grpc_grpc//:grpc++",
            #"@com_github_grpc_grpc//:grpc++_reflection", # Disabled until fixed upstream
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
				Name:             "cpp_proto_compile",
				Kind:             "proto",
				Implementation:   compileRuleTemplate,
				Plugins:          []string{"//cpp:cpp"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates *.h,*.cc protobuf artifacts",
				Attrs:            append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:             "cpp_grpc_compile",
				Kind:             "grpc",
				Implementation:   compileRuleTemplate,
				Plugins:          []string{"//cpp:cpp", "//cpp:grpc_cpp"},
				WorkspaceExample: grpcWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates *.h,*.cc protobuf+gRPC artifacts",
				Attrs:            append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:             "cpp_proto_library",
				Kind:             "proto",
				Implementation:   cppProtoLibraryRuleTemplate,
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates *.h,*.cc protobuf library",
				Attrs:            append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:             "cpp_grpc_library",
				Kind:             "grpc",
				Implementation:   cppGrpcLibraryRuleTemplate,
				WorkspaceExample: grpcWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates *.h,*.cc protobuf+gRPC library",
				Attrs:            append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
