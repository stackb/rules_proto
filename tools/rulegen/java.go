package main

var javaGrpcUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//:deps.bzl", "io_grpc_grpc_java")

io_grpc_grpc_java()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")

grpc_java_repositories(omit_com_google_protobuf = True)

load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()`)

var javaLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create {{ .Lang.Name }} library
    native.java_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [str(Label("//java:{{ .Rule.Kind }}_deps"))],
        exports = [
            str(Label("//java:{{ .Rule.Kind }}_deps")),
        ],
        visibility = kwargs.get("visibility"),
    )`)

func makeJava() *Language {
	return &Language{
		Dir:              "java",
		Name:             "java",
		RouteGuideClient: "//java/example/routeguide:client",
		RouteGuideServer: "//java/example/routeguide:server",
		Flags:            commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:           "java_proto_compile",
				Kind:           "proto",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//java:java"},
				Usage:          usageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates a srcjar with protobuf *.java files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "java_grpc_compile",
				Kind:           "grpc",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//java:java", "//java:grpc_java"},
				Usage:          usageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates a srcjar with protobuf+gRPC *.java files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "java_proto_library",
				Kind:           "proto",
				Implementation: javaLibraryRuleTemplate,
				Usage:          usageTemplate,
				Example:        protoLibraryExampleTemplate,
				Doc:            "Generates a jar with compiled protobuf *.class files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
				Flags: []*Flag{
					{
						Category:    "build",
						Name:        "incompatible_disallow_struct_provider_syntax",
						Value:       "false",
						Description: "bazel_tools/tools/jdk/java_toolchain_alias.bzl",
					},
					{
						Category: "build",
						Name:     "incompatible_remove_native_maven_jar",
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
				Name:           "java_grpc_library",
				Kind:           "grpc",
				Implementation: javaLibraryRuleTemplate,
				Usage:          javaGrpcUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates a jar with compiled protobuf+gRPC *.class files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
				Flags: []*Flag{
					{
						Category:    "build",
						Name:        "incompatible_disallow_struct_provider_syntax",
						Value:       "false",
						Description: "bazel_tools/tools/jdk/java_toolchain_alias.bzl",
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
