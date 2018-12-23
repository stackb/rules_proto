package main

var tsProtocGenUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@org_pubref_rules_node//node:rules.bzl", "node_repositories", "yarn_modules")

node_repositories()

load("@build_bazel_rules_nodejs//:defs.bzl", "node_repositories")

node_repositories(package_json = ["@ts_protoc_gen//:package.json"])

load("@build_bazel_rules_typescript//:defs.bzl", "ts_setup_workspace")

ts_setup_workspace()

load("@io_bazel_rules_webtesting//web:repositories.bzl", "browser_repositories", "web_test_repositories")

web_test_repositories()

load("@build_bazel_rules_nodejs//:defs.bzl", "npm_install")

npm_install(
    name = "deps",
    package_json = "@ts_protoc_gen//:package.json",
    package_lock_json = "@ts_protoc_gen//:package-lock.json",
)`)

func makeGithubComImprobableTsProtocGen() *Language {
	return &Language{
		Dir:  "github.com/improbable-eng/ts-protoc-gen",
		Name: "ts-protoc-gen",
		Rules: []*Rule{
			&Rule{
				Name:           "ts_proto_compile",
				Experimental:   true,
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//github.com/improbable-eng/ts-protoc-gen:ts"},
				Usage:          tsProtocGenUsageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates typescript protobuf t.ds files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
