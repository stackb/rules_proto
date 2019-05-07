package main

var rubyProtoLibraryUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@com_github_yugui_rules_ruby//ruby:def.bzl", "ruby_register_toolchains")

ruby_register_toolchains()

load("@com_github_yugui_rules_ruby//ruby/private:bundle.bzl", "bundle_install")

bundle_install(
    name = "routeguide_gems_bundle",
    gemfile = "//ruby:Gemfile",
    gemfile_lock = "//ruby:Gemfile.lock",
)`)

var rubyGrpcLibraryUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

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

var rubyProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:ruby_proto_compile.bzl", "ruby_proto_compile")
load("@com_github_yugui_rules_ruby//ruby:def.bzl", "ruby_library")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    ruby_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    ruby_library(
        name = name,
        srcs = [name_pb],
        includes = ["{package}/%s" % name_pb],
        visibility = visibility,
    )`)

var rubyGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:ruby_grpc_compile.bzl", "ruby_grpc_compile")
load("@com_github_yugui_rules_ruby//ruby:def.bzl", "ruby_library")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    ruby_grpc_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    ruby_library(
        name = name,
        srcs = [name_pb],
        includes = ["{package}/%s" % name_pb],
        visibility = visibility,
    )`)

func makeRuby() *Language {
	return &Language{
		Dir:   "ruby",
		Name:  "ruby",
		Flags: commonLangFlags,
		Notes: aspectLangNotes,
		Rules: []*Rule{
			&Rule{
				Name:           "ruby_proto_compile",
				Implementation: aspectRuleTemplate,
				Plugins:        []string{"//ruby:ruby"},
				Usage:          usageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates *.ruby protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "ruby_grpc_compile",
				Implementation: aspectRuleTemplate,
				Plugins:        []string{"//ruby:ruby", "//ruby:grpc_ruby"},
				Usage:          grpcUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates *.ruby protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "ruby_proto_library",
				Implementation: rubyProtoLibraryRuleTemplate,
				Usage:          rubyProtoLibraryUsageTemplate,
				Example:        protoLibraryExampleTemplate,
				Doc:            "Generates *.rb protobuf library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "ruby_grpc_library",
				Implementation: rubyGrpcLibraryRuleTemplate,
				Usage:          rubyGrpcLibraryUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates *.rb protobuf+gRPC library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
