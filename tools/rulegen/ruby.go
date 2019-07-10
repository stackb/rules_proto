package main

var rubyProtoLibraryWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Lang.Name }}_deps")

{{ .Lang.Name }}_deps()

load("@com_github_yugui_rules_ruby//ruby:def.bzl", "ruby_register_toolchains")

ruby_register_toolchains()

load("@com_github_yugui_rules_ruby//ruby/private:bundle.bzl", "bundle_install")

bundle_install(
    name = "routeguide_gems_bundle",
    gemfile = "//ruby:Gemfile",
    gemfile_lock = "//ruby:Gemfile.lock",
)`)

var rubyGrpcLibraryWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Lang.Name }}_deps")

{{ .Lang.Name }}_deps()

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
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create {{ .Lang.Name }} library
    ruby_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        includes = [name_pb],
        visibility = kwargs.get("visibility"),
    )`)

func makeRuby() *Language {
	return &Language{
		Dir:   "ruby",
		Name:  "ruby",
		DisplayName: "Ruby",
		Notes: mustTemplate("Rules for generating Ruby protobuf and gRPC `.rb` files and libraries using standard Protocol Buffers and gRPC. Libraries are created with `ruby_library` from [rules_ruby](https://github.com/yugui/rules_ruby)"),
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:             "ruby_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//ruby:ruby"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates Ruby protobuf `.rb` artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "ruby_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//ruby:ruby", "//ruby:grpc_ruby"},
				WorkspaceExample: grpcWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates Ruby protobuf+gRPC `.rb` artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "ruby_proto_library",
				Kind:             "proto",
				Implementation:   rubyLibraryRuleTemplate,
				WorkspaceExample: rubyProtoLibraryWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates a Ruby protobuf library using `ruby_library` from `rules_ruby`",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "ruby_grpc_library",
				Kind:             "grpc",
				Implementation:   rubyLibraryRuleTemplate,
				WorkspaceExample: rubyGrpcLibraryWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates a Ruby protobuf+gRPC library using `ruby_library` from `rules_ruby`",
				Attrs:            aspectProtoCompileAttrs,
			},
		},
	}
}
