package main

var cppLibraryRuleTemplateString = `load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )
`

var cppProtoLibraryRuleTemplate = mustTemplate(cppLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    native.cc_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = PROTO_DEPS,
        includes = [name_pb],
        visibility = kwargs.get("visibility"),
    )

PROTO_DEPS = [
    "@com_google_protobuf//:protoc_lib",
]`)

var cppGrpcLibraryRuleTemplate = mustTemplate(cppLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    native.cc_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = GRPC_DEPS,
        includes = [name_pb],
        visibility = kwargs.get("visibility"),
    )

GRPC_DEPS = [
    "@com_google_protobuf//:protoc_lib",
    "@com_github_grpc_grpc//:grpc++",
    #"@com_github_grpc_grpc//:grpc++_reflection", # TODO: Disabled until fixed upstream
]`)

func makeCpp() *Language {
	return &Language{
		Dir:   "cpp",
		Name:  "cpp",
		DisplayName: "C++",
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:             "cpp_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//cpp:cpp"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates *.h,*.cc protobuf artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "cpp_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//cpp:cpp", "//cpp:grpc_cpp"},
				WorkspaceExample: grpcWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates *.h,*.cc protobuf+gRPC artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "cpp_proto_library",
				Kind:             "proto",
				Implementation:   cppProtoLibraryRuleTemplate,
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates *.h,*.cc protobuf library",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "cpp_grpc_library",
				Kind:             "grpc",
				Implementation:   cppGrpcLibraryRuleTemplate,
				WorkspaceExample: grpcWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates *.h,*.cc protobuf+gRPC library",
				Attrs:            aspectProtoCompileAttrs,
			},
		},
	}
}
