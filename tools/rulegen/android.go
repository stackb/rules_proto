package main

var androidLibraryWorkspaceTemplateString = `load("@build_stack_rules_proto//:deps.bzl", "io_grpc_grpc_java")

io_grpc_grpc_java()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")

grpc_java_repositories(
    omit_com_google_protobuf = True,
)

load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@build_bazel_rules_android//android:sdk_repository.bzl", "android_sdk_repository")

android_sdk_repository(name = "androidsdk")`

var androidGrpcLibraryWorkspaceTemplate = mustTemplate(androidLibraryWorkspaceTemplateString)

var androidProtoLibraryWorkspaceTemplate = mustTemplate("# The set of dependencies loaded here is excessive for android proto alone\n# (but simplifies our setup)\n" + androidLibraryWorkspaceTemplateString)

var androidLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@build_bazel_rules_android//android:rules.bzl", "android_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create {{ .Lang.Name }} library
    android_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            Label("//android:{{ .Rule.Kind }}_deps"),
        ],
        exports = [
            Label("//android:{{ .Rule.Kind }}_deps"),
        ],
        visibility = kwargs.get("visibility"),
    )`)

func makeAndroid() *Language {
	return &Language{
		Dir:  "android",
		Name: "android",
		Flags: commonLangFlags,
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
			},
			&Rule{
				Name:             "android_proto_library",
				Kind:             "proto",
				Implementation:   androidLibraryRuleTemplate,
				WorkspaceExample: androidProtoLibraryWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates android protobuf library",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "android_grpc_library",
				Kind:             "grpc",
				Implementation:   androidLibraryRuleTemplate,
				WorkspaceExample: androidGrpcLibraryWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates android protobuf+gRPC library",
				Attrs:            aspectProtoCompileAttrs,
			},
		},
	}
}
