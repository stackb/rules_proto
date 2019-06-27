package main

var swiftWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load(
    "@build_bazel_rules_swift//swift:repositories.bzl",
    "swift_rules_dependencies",
)

swift_rules_dependencies()

load(
    "@build_bazel_apple_support//lib:repositories.bzl",
    "apple_support_dependencies",
)

apple_support_dependencies()`)

var swiftLibraryRuleTemplate = mustTemplate(`load("@build_bazel_rules_swift//swift:swift.bzl", _{{ .Lang.Name }}_{{ .Rule.Kind }}_library = "{{ .Lang.Name }}_{{ .Rule.Kind }}_library")

{{ .Lang.Name }}_{{ .Rule.Kind }}_library = _{{ .Lang.Name }}_{{ .Rule.Kind }}_library`)

var swiftGrpcLibraryExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "person_{{ .Lang.Name }}_library",
    flavor = "client",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)`)

func makeSwift() *Language {
	return &Language{
		Dir:  "swift",
		Name: "swift",
		// TravisExclusionReason: "travis incompatible",
		PresubmitEnvVars: map[string]string{
			"CC": "clang",
		},
		Flags: append(commonLangFlags, &Flag{
			Category: "build",
			Name:     "strategy=SwiftCompile",
			Value:    "standalone",
		}),
		Rules: []*Rule{
			&Rule{
				Name:             "swift_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//swift:swift"},
				WorkspaceExample: swiftWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates swift protobuf artifacts",
				Attrs:            aspectProtoCompileAttrs,
				BazelCIExclusionReason: "experimental",
				Experimental:     true,
			},
			&Rule{
				Name:             "swift_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//swift:grpc_swift"},
				WorkspaceExample: swiftWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates swift protobuf+gRPC artifacts",
				Attrs:            aspectProtoCompileAttrs,
				BazelCIExclusionReason: "experimental",
				Experimental:     true,
			},
			&Rule{
				Name:             "swift_proto_library",
				Kind:             "proto",
				Implementation:   swiftLibraryRuleTemplate,
				WorkspaceExample: swiftWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates swift protobuf library",
				Attrs:            aspectProtoCompileAttrs,
				BazelCIExclusionReason: "experimental",
				Experimental:     true,
			},
			&Rule{
				Name:             "swift_grpc_library",
				Kind:             "grpc",
				Implementation:   swiftLibraryRuleTemplate,
				WorkspaceExample: swiftWorkspaceTemplate,
				BuildExample:     swiftGrpcLibraryExampleTemplate,
				Doc:              "Generates swift protobuf+gRPC library",
				Attrs:            aspectProtoCompileAttrs,
				BazelCIExclusionReason: "experimental",
				Experimental:     true,
			},
		},
	}
}
