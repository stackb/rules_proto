package main

var androidLibraryWorkspaceTemplateString = `load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Lang.Name }}_deps")

{{ .Lang.Name }}_deps()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")

grpc_java_repositories(
    omit_com_google_protobuf = True,
    omit_com_google_protobuf_javalite = True,
    omit_net_zlib = True
)

load("@build_bazel_rules_android//android:sdk_repository.bzl", "android_sdk_repository")

android_sdk_repository(name = "androidsdk")`

var androidGrpcLibraryWorkspaceTemplate = mustTemplate(androidLibraryWorkspaceTemplateString)

var androidProtoLibraryWorkspaceTemplate = mustTemplate("# The set of dependencies loaded here is excessive for android proto alone\n# (but simplifies our setup)\n" + androidLibraryWorkspaceTemplateString)

var androidLibraryRuleTemplateString = `load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@build_bazel_rules_android//android:rules.bzl", "android_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )
`

var androidProtoLibraryRuleTemplate = mustTemplate(androidLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    android_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = PROTO_DEPS,
        exports = PROTO_DEPS,
        visibility = kwargs.get("visibility"),
    )

PROTO_DEPS = [
    "@com_google_guava_guava_android//jar",
    "@com_google_protobuf//:protobuf_javalite",
    "@javax_annotation_javax_annotation_api//jar"
]`)

var androidGrpcLibraryRuleTemplate = mustTemplate(androidLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    android_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = GRPC_DEPS,
        exports = GRPC_DEPS,
        visibility = kwargs.get("visibility"),
    )

GRPC_DEPS = [
    "@com_google_guava_guava_android//jar",
    "@com_google_protobuf//:protobuf_javalite",
    "@com_google_protobuf//:protobuf_java_util",
    "@javax_annotation_javax_annotation_api//jar",
    "@io_grpc_grpc_java//core",
    "@io_grpc_grpc_java//protobuf-lite",
    "@io_grpc_grpc_java//stub",
]`)

func makeAndroid() *Language {
	return &Language{
		Dir:  "android",
		Name: "android",
		DisplayName: "Android",
		Flags: commonLangFlags,
		SkipDirectoriesMerge: true,
		SkipTestPlatforms: []string{"all"},
		Rules: []*Rule{
			&Rule{
				Name:             "android_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//android:java"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates android protobuf artifacts",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{"none"},
			},
			&Rule{
				Name:             "android_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//android:java", "//android:grpc_javalite"},
				WorkspaceExample: javaGrpcWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates android protobuf+gRPC artifacts",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{"none"},
			},
			&Rule{
				Name:             "android_proto_library",
				Kind:             "proto",
				Implementation:   androidProtoLibraryRuleTemplate,
				WorkspaceExample: androidProtoLibraryWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates android protobuf library",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{"none"},
			},
			&Rule{
				Name:             "android_grpc_library",
				Kind:             "grpc",
				Implementation:   androidGrpcLibraryRuleTemplate,
				WorkspaceExample: androidGrpcLibraryWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates android protobuf+gRPC library",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{"none"},
			},
		},
	}
}
