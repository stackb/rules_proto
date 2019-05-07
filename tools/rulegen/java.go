package main

var javaGrpcUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//:deps.bzl", "io_grpc_grpc_java")

io_grpc_grpc_java()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")

grpc_java_repositories(omit_com_google_protobuf = True)

load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()`)

var javaProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:java_proto_compile.bzl", "java_proto_compile")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    java_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    native.java_library(
        name = name,
        srcs = [name_pb],
        deps = [str(Label("//java:proto_deps"))],
        exports = [
            str(Label("//java:proto_deps")),
        ],
        visibility = visibility,
    )`)

var javaGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:java_grpc_compile.bzl", "java_grpc_compile")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    java_grpc_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    native.java_library(
        name = name,
        srcs = [name_pb],
        deps = [str(Label("//java:grpc_deps"))],
        exports = [
            str(Label("//java:grpc_deps")),
        ],
        visibility = visibility,
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
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//java:java"},
				Usage:          usageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates a srcjar with protobuf *.java files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "java_grpc_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//java:java", "//java:grpc_java"},
				Usage:          usageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates a srcjar with protobuf+gRPC *.java files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "java_proto_library",
				Implementation: javaProtoLibraryRuleTemplate,
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
				},
			},
			&Rule{
				Name:           "java_grpc_library",
				Implementation: javaGrpcLibraryRuleTemplate,
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
