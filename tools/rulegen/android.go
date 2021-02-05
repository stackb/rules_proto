package main

var androidLibraryWorkspaceTemplateString = `load("@rules_proto_grpc//{{ .Lang.Dir }}:repositories.bzl", rules_proto_grpc_{{ .Lang.Name }}_repos="{{ .Lang.Name }}_repos")

rules_proto_grpc_{{ .Lang.Name }}_repos()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")

grpc_java_repositories()

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
        tags = kwargs.get("tags"),
    )

PROTO_DEPS = [
    "@com_google_protobuf//:protobuf_javalite",
]`)

var androidGrpcLibraryRuleTemplate = mustTemplate(androidLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    android_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = GRPC_DEPS,
        exports = GRPC_DEPS,
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags"),
    )

GRPC_DEPS = [
    "@io_grpc_grpc_java//api",
    "@io_grpc_grpc_java//protobuf-lite",
    "@io_grpc_grpc_java//stub",
    "@io_grpc_grpc_java//stub:javax_annotation",
    "@com_google_code_findbugs_jsr305//jar",
    "@com_google_guava_guava//jar",
    "@com_google_protobuf//:protobuf_javalite",
    "@com_google_protobuf//:protobuf_java_util",
]`)

func makeAndroid() *Language {
	return &Language{
		Dir:  "android",
		Name: "android",
		DisplayName: "Android",
		Notes: mustTemplate("Rules for generating Android protobuf and gRPC `.jar` files and libraries using standard Protocol Buffers and [gRPC-Java](https://github.com/grpc/grpc-java). Libraries are created with `android_library` from [rules_android](https://github.com/bazelbuild/rules_android)"),
		Flags: commonLangFlags,
		SkipDirectoriesMerge: true, // Jar files are not needed to be in merged directory
		SkipTestPlatforms: []string{"all"},
		Rules: []*Rule{
			&Rule{
				Name:             "android_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//android:javalite_plugin"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates an Android protobuf `.jar` artifact",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{"none"},
			},
			&Rule{
				Name:             "android_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//android:javalite_plugin", "//android:grpc_javalite_plugin"},
				WorkspaceExample: javaGrpcWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates Android protobuf+gRPC `.jar` artifacts",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{"none"},
			},
			&Rule{
				Name:             "android_proto_library",
				Kind:             "proto",
				Implementation:   androidProtoLibraryRuleTemplate,
				WorkspaceExample: androidProtoLibraryWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates an Android protobuf library using `android_library` from `rules_android`",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{"none"},
			},
			&Rule{
				Name:             "android_grpc_library",
				Kind:             "grpc",
				Implementation:   androidGrpcLibraryRuleTemplate,
				WorkspaceExample: androidGrpcLibraryWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates Android protobuf+gRPC library using `android_library` from `rules_android`",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{"none"},
			},
		},
	}
}
