package main

var closureLibraryUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories(
    omit_com_google_protobuf = True,
)`)

var closureProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:closure_proto_compile.bzl", "closure_proto_compile")
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")
    transitive = kwargs.pop("transitive", True)
    transitivity = kwargs.get("transitivity")

    name_pb = name + "_pb"

    closure_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        transitive = transitive,
        transitivity = transitivity,
    )

    closure_js_library(
        name = name,
        srcs = [name_pb],
        deps = ["@io_bazel_rules_closure//closure/protobuf:jspb"],
        visibility = visibility,
        internal_descriptors = [name_pb + "/descriptor.source.bin"],
        lenient = True,
    )
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")
`)

func makeClosure() *Language {
	return &Language{
		Dir:  "closure",
		Name: "closure",
		Rules: []*Rule{
			&Rule{
				Name:           "closure_proto_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//closure:js"},
				Usage:          usageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates closure *.js protobuf+gRPC files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "closure_proto_library",
				Implementation: closureProtoLibraryRuleTemplate,
				Usage:          closureLibraryUsageTemplate,
				Example:        protoLibraryExampleTemplate,
				Doc:            "Generates a closure_library with compiled protobuf *.js files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
