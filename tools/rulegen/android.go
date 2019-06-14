package main

var androidLibraryUsageTemplateString = `load("@build_stack_rules_proto//:deps.bzl", "io_grpc_grpc_java")

io_grpc_grpc_java()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")

grpc_java_repositories(
    omit_com_google_protobuf = True,
)

load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@build_bazel_rules_android//android:sdk_repository.bzl", "android_sdk_repository")

android_sdk_repository(name = "androidsdk")`

var androidGrpcLibraryUsageTemplate = mustTemplate(androidLibraryUsageTemplateString)

var androidProtoLibraryUsageTemplate = mustTemplate("# The set of dependencies loaded here is excessive for android proto alone\n# (but simplifies our setup)\n" + androidLibraryUsageTemplateString)

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
            str(Label("//android:{{ .Rule.Kind }}_deps")),
        ],
        exports = [
            str(Label("//android:{{ .Rule.Kind }}_deps")),
        ],
        visibility = kwargs.get("visibility"),
    )`)

func makeAndroid() *Language {
	return &Language{
		Dir:  "android",
		Name: "android",
		Flags: append([]*Flag{}, &Flag{
			Category: "build",
			Name:     "incompatible_disable_deprecated_attr_params",
			Value:    "false",
		}),
		Rules: []*Rule{
			&Rule{
				Name:           "android_proto_compile",
				Kind:           "proto",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//android:javalite"},
				Usage:          usageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates android protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "android_grpc_compile",
				Kind:           "grpc",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//android:javalite", "//android:grpc_javalite"},
				Usage:          javaGrpcUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates android protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "android_proto_library",
				Kind:           "proto",
				Implementation: androidLibraryRuleTemplate,
				Usage:          androidProtoLibraryUsageTemplate,
				Example:        protoLibraryExampleTemplate,
				Doc:            "Generates android protobuf library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
				Flags: []*Flag{
					{
						Category: "build",
						Name:     "incompatible_remove_native_maven_jar",
						Value:    "false",
					},
					{
						Category: "build",
						Name:     "incompatible_disallow_struct_provider_syntax",
						Value:    "false",
					},
					{
						Category: "build",
						Name:     "incompatible_use_toolchain_resolution_for_java_rules",
						Value:    "false",
					},
				},
			},
			&Rule{
				Name:           "android_grpc_library",
				Kind:           "grpc",
				Implementation: androidLibraryRuleTemplate,
				Usage:          androidGrpcLibraryUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates android protobuf+gRPC library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
				Flags: []*Flag{
					{
						Category: "build",
						Name:     "incompatible_remove_native_maven_jar",
						Value:    "false",
					},
					{
						Category: "build",
						Name:     "incompatible_disallow_struct_provider_syntax",
						Value:    "false",
					},
					{
						Category: "build",
						Name:     "incompatible_use_toolchain_resolution_for_java_rules",
						Value:    "false",
					},
				},
			},
		},
	}
}
