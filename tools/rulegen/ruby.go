package main

var rubyProtoLibraryWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@com_github_yugui_rules_ruby//ruby:def.bzl", "ruby_register_toolchains")

ruby_register_toolchains()

load("@com_github_yugui_rules_ruby//ruby/private:bundle.bzl", "bundle_install")

bundle_install(
    name = "routeguide_gems_bundle",
    gemfile = "//ruby:Gemfile",
    gemfile_lock = "//ruby:Gemfile.lock",
)`)

var rubyGrpcLibraryWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@com_github_yugui_rules_ruby//ruby:def.bzl", "ruby_register_toolchains")

ruby_register_toolchains()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@com_github_yugui_rules_ruby//ruby/private:bundle.bzl", "bundle_install")

bundle_install(
    name = "routeguide_gems_bundle",
    gemfile = "//ruby:Gemfile",
    gemfile_lock = "//ruby:Gemfile.lock",
)`)

var rubyLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@com_github_yugui_rules_ruby//ruby:def.bzl", "ruby_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create {{ .Lang.Name }} library
    ruby_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        includes = ["{package}/%s" % name_pb],
        visibility = kwargs.get("visibility"),
    )`)

func makeRuby() *Language {
	return &Language{
		Dir:   "ruby",
		Name:  "ruby",
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:             "ruby_proto_compile",
				Kind:             "proto",
				Implementation:   compileRuleTemplate,
				Plugins:          []string{"//ruby:ruby"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates *.ruby protobuf artifacts",
				Attrs:            append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:             "ruby_grpc_compile",
				Kind:             "grpc",
				Implementation:   compileRuleTemplate,
				Plugins:          []string{"//ruby:ruby", "//ruby:grpc_ruby"},
				WorkspaceExample: grpcWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates *.ruby protobuf+gRPC artifacts",
				Attrs:            append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:             "ruby_proto_aspect_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Experimental:     true,
				Plugins:          []string{"//ruby:ruby"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates *.ruby protobuf artifacts using aspect based compilation",
				Attrs:            append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:             "ruby_grpc_aspect_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Experimental:     true,
				Plugins:          []string{"//ruby:ruby", "//ruby:grpc_ruby"},
				WorkspaceExample: grpcWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates *.ruby protobuf+gRPC artifacts using aspect based compilation",
				Attrs:            append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:             "ruby_proto_library",
				Kind:             "proto",
				Implementation:   rubyLibraryRuleTemplate,
				WorkspaceExample: rubyProtoLibraryWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates *.rb protobuf library",
				Attrs:            append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:             "ruby_grpc_library",
				Kind:             "grpc",
				Implementation:   rubyLibraryRuleTemplate,
				WorkspaceExample: rubyGrpcLibraryWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates *.rb protobuf+gRPC library",
				Attrs:            append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
