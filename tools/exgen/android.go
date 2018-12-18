package main

var androidProtoLibraryUsageTemplate = mustTemplate(`
# The set of dependencies loaded here is excessive for android proto alone
# (but simplifies our setup)
load("@build_stack_rules_proto//:deps.bzl", "io_grpc_grpc_java")
io_grpc_grpc_java()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")
grpc_java_repositories(
    omit_com_google_protobuf = True,
)

load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")
{{ .Rule.Name }}()

load("@build_bazel_rules_android//android:sdk_repository.bzl", "android_sdk_repository")
android_sdk_repository(name = "androidsdk")

load("@gmaven_rules//:gmaven.bzl", "gmaven_rules")
gmaven_rules()`)

var androidGrpcLibraryUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//:deps.bzl", "io_grpc_grpc_java")
io_grpc_grpc_java()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")
grpc_java_repositories(
    omit_com_google_protobuf = True,
)

load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")
{{ .Rule.Name }}()

#load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")
#grpc_deps()

load("@build_bazel_rules_android//android:sdk_repository.bzl", "android_sdk_repository")
android_sdk_repository(name = "androidsdk")

load("@gmaven_rules//:gmaven.bzl", "gmaven_rules")
gmaven_rules()`)

var androidProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Rule.Base}}_{{ .Rule.Kind }}_compile.bzl", "{{ .Rule.Base }}_{{ .Rule.Kind }}_compile")
load("@build_bazel_rules_android//android:rules.bzl", "android_library")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    {{ .Rule.Base}}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )

    android_library(
        name = name,
        srcs = [name_pb],
        deps = [
            str(Label("//android:proto_deps")),
        ],
        exports = [
            str(Label("//android:proto_deps")),
        ],
        visibility = visibility,
    )
`)

var androidGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Rule.Base}}_{{ .Rule.Kind }}_compile.bzl", "{{ .Rule.Base }}_{{ .Rule.Kind }}_compile")
load("@build_bazel_rules_android//android:rules.bzl", "android_library")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    {{ .Rule.Base}}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )

    android_library(
        name = name,
        srcs = [name_pb],
        deps = [
            str(Label("//android:grpc_deps")),
        ],
        exports = [
            str(Label("//android:grpc_deps")),
        ],
        visibility = visibility,
	)
`)

func makeAndroid() *Language {
	return &Language{
		Dir:  "android",
		Name: "android",
		Rules: []*Rule{
			&Rule{
				Name:           "android_proto_compile",
				Base:           "android",
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
				Base:           "android",
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
				Base:           "android",
				Kind:           "proto",
				Usage:          androidProtoLibraryUsageTemplate,
				Example:        protoLibraryExampleTemplate,
				Implementation: androidProtoLibraryRuleTemplate,
				Doc:            "Generates android protobuf library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "android_grpc_library",
				Base:           "android",
				Kind:           "grpc",
				Implementation: androidGrpcLibraryRuleTemplate,
				Usage:          androidGrpcLibraryUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates android protobuf+gRPC library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
