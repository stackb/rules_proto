package main

var swiftUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")
{{ .Rule.Name }}()

# rules_go used here to compile a wrapper around the protoc-gen-swift plugin
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

load("@build_stack_rules_proto//{{ .Lang.Dir }}:repositories.bzl", "swift_toolchain")
# Default values work with linux, x86_64, /usr/local/bin/clang. 
swift_toolchain(
	#root = "/home/pcj/.local/share/umake/swift/swift-lang/usr",
)

# Uncomment for ocal development with swift installed on your machine
# load("@build_bazel_rules_swift//swift:repositories.bzl", "swift_rules_dependencies")
# swift_rules_dependencies()

`)

var swiftProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Rule.Base}}_{{ .Rule.Kind }}_compile.bzl", "{{ .Rule.Base }}_{{ .Rule.Kind }}_compile")
load("@build_bazel_rules_swift//swift:swift.bzl", _swift_proto_library = "swift_proto_library")

def {{ .Rule.Name }}(**kwargs):
    _swift_proto_library(**kwargs)
`)

var swiftGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Rule.Base}}_{{ .Rule.Kind }}_compile.bzl", "{{ .Rule.Base }}_{{ .Rule.Kind }}_compile")
load("@build_bazel_rules_swift//swift:swift.bzl", _swift_proto_library = "swift_proto_library")

def {{ .Rule.Name }}(**kwargs):
    _swift_proto_library(**kwargs)
`)

func makeSwift() *Language {
	return &Language{
		Dir:  "swift",
		Name: "swift",
		Notes: mustTemplate(`
**NOTE**: The swift rules are essentially non-functional.  The protoc-plugin core dumps despite all efforts thus far on linux.`),
		Rules: []*Rule{
			&Rule{
				Experimental:   true,
				Name:           "swift_proto_compile",
				Base:           "swift",
				Kind:           "proto",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//swift:swift"},
				Usage:          swiftUsageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates swift protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Experimental:   true,
				Name:           "swift_grpc_compile",
				Base:           "swift",
				Kind:           "grpc",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//swift:grpc_swift"},
				Usage:          swiftUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates swift protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Experimental:   true,
				Name:           "swift_proto_library",
				Base:           "swift",
				Kind:           "proto",
				Usage:          swiftUsageTemplate,
				Example:        protoLibraryExampleTemplate,
				Implementation: swiftProtoLibraryRuleTemplate,
				Doc:            "Generates swift protobuf library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Experimental:   true,
				Name:           "swift_grpc_library",
				Base:           "swift",
				Kind:           "grpc",
				Implementation: swiftGrpcLibraryRuleTemplate,
				Usage:          swiftUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates swift protobuf+gRPC library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
