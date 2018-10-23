package main

var javaGrpcUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//:deps.bzl", "io_grpc_grpc_java")

io_grpc_grpc_java()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")
grpc_java_repositories(
    omit_com_google_protobuf = True,
)

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
        transitive = True,
        visibility = visibility,
    )
    native.java_library(
        name = name,
        srcs = [name_pb],
        deps = [str(Label("//java:proto_deps"))],
        exports = [
            str(Label("//java:proto_deps")),
        ],
        visibility = visibility,
    )
`)

var javaGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:java_grpc_compile.bzl", "java_grpc_compile")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    java_grpc_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )
    native.java_library(
        name = name,
        srcs = [name_pb],
        deps = [str(Label("//java:grpc_deps"))],
        exports = [
            str(Label("//java:grpc_deps")),
        ],
        visibility = visibility,
    )
`)

func makeJava() *Language {
	return &Language{
		Dir:  "java",
		Name: "java",
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
			},
			&Rule{
				Name:           "java_grpc_library",
				Implementation: javaGrpcLibraryRuleTemplate,
				Usage:          javaGrpcUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates a jar with compiled protobuf+gRPC *.class files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
