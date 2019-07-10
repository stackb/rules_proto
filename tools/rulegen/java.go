package main

var javaProtoWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Lang.Name }}_deps")

{{ .Lang.Name }}_deps()

load("@io_grpc_grpc_java//:repositories.bzl", "com_google_guava")

com_google_guava()`)

var javaGrpcWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Lang.Name }}_deps")

{{ .Lang.Name }}_deps()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")

grpc_java_repositories(
    omit_com_google_protobuf = True,
    omit_net_zlib = True
)`)

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
    )

PROTO_DEPS = [
    "@com_google_guava_guava//jar",
    "@com_google_protobuf//:protobuf_java",
    "@javax_annotation_javax_annotation_api//jar",
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
    )

GRPC_DEPS = [
    "@com_google_guava_guava//jar",
    "@com_google_protobuf//:protobuf_java",
    "@com_google_protobuf//:protobuf_java_util",
    "@javax_annotation_javax_annotation_api//jar",
    "@io_grpc_grpc_java//core",
    "@io_grpc_grpc_java//protobuf",
    "@io_grpc_grpc_java//stub",
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
				Plugins:          []string{"//java:java"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates a Java protobuf srcjar artifact",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "java_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//java:java", "//java:grpc_java"},
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
