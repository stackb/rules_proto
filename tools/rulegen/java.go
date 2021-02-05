package main

var javaProtoWorkspaceTemplate = mustTemplate(`load("@rules_proto_grpc//{{ .Lang.Dir }}:repositories.bzl", rules_proto_grpc_{{ .Lang.Name }}_repos="{{ .Lang.Name }}_repos")

rules_proto_grpc_{{ .Lang.Name }}_repos()`)

var javaGrpcWorkspaceTemplate = mustTemplate(`load("@rules_proto_grpc//{{ .Lang.Dir }}:repositories.bzl", rules_proto_grpc_{{ .Lang.Name }}_repos="{{ .Lang.Name }}_repos")

rules_proto_grpc_{{ .Lang.Name }}_repos()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")

grpc_java_repositories()`)

var javaLibraryRuleTemplateString = `load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )
`

var javaProtoLibraryRuleTemplate = mustTemplate(javaLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    native.java_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = PROTO_DEPS,
        exports = PROTO_DEPS,
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags"),
    )

PROTO_DEPS = [
    "@com_google_protobuf//:protobuf_java",
]`)

var javaGrpcLibraryRuleTemplate = mustTemplate(javaLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    native.java_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = GRPC_DEPS,
        runtime_deps = ["@io_grpc_grpc_java//netty"],
        exports = GRPC_DEPS,
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags"),
    )

GRPC_DEPS = [  # From https://github.com/grpc/grpc-java/blob/3ce5df3f78e8fd4a619a791914087dd4c0562835/compiler/BUILD.bazel#L21-L27
    "@io_grpc_grpc_java//api",
    "@io_grpc_grpc_java//protobuf",
    "@io_grpc_grpc_java//stub",
    "@io_grpc_grpc_java//stub:javax_annotation",
    "@com_google_code_findbugs_jsr305//jar",
    "@com_google_guava_guava//jar",
    "@com_google_protobuf//:protobuf_java",
    "@com_google_protobuf//:protobuf_java_util",
]`)

func makeJava() *Language {
	return &Language{
		Dir:              "java",
		Name:             "java",
		DisplayName:      "Java",
		Notes: mustTemplate("Rules for generating Java protobuf and gRPC `.jar` files and libraries using standard Protocol Buffers and [gRPC-Java](https://github.com/grpc/grpc-java). Libraries are created with the Bazel native `java_library`"),
		Flags:            commonLangFlags,
		SkipDirectoriesMerge: true,
		Rules: []*Rule{
			&Rule{
				Name:             "java_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//java:java_plugin"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates a Java protobuf srcjar artifact",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "java_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//java:java_plugin", "//java:grpc_java_plugin"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates a Java protobuf+gRPC srcjar artifact",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "java_proto_library",
				Kind:             "proto",
				Implementation:   javaProtoLibraryRuleTemplate,
				WorkspaceExample: javaProtoWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates a Java protobuf library using `java_library`",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "java_grpc_library",
				Kind:             "grpc",
				Implementation:   javaGrpcLibraryRuleTemplate,
				WorkspaceExample: javaGrpcWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates a Java protobuf+gRPC library using `java_library`",
				Attrs:            aspectProtoCompileAttrs,
			},
		},
	}
}
