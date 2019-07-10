package main

var closureLibraryWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Lang.Name }}_deps")

{{ .Lang.Name }}_deps()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories()`)

var closureProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create {{ .Lang.Name }} library
    closure_js_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = PROTO_DEPS,
        visibility = kwargs.get("visibility"),
        suppress = [
            "JSC_LATE_PROVIDE_ERROR",
            "JSC_UNDEFINED_VARIABLE",
            "JSC_IMPLICITLY_NULLABLE_JSDOC",
            "JSC_STRICT_INEXISTENT_PROPERTY",
            "JSC_POSSIBLE_INEXISTENT_PROPERTY",
            "JSC_UNRECOGNIZED_TYPE_ERROR",
        ],
    )

PROTO_DEPS = [
    "@io_bazel_rules_closure//closure/protobuf:jspb"
]`)

func makeClosure() *Language {
	return &Language{
		Dir:   "closure",
		Name:  "closure",
		DisplayName: "Closure",
		Notes: mustTemplate("Rules for generating Closure protobuf `.js` files and libraries using standard Protocol Buffers. Libraries are created with `closure_js_library` from [rules_closure](https://github.com/bazelbuild/rules_closure)"),
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:             "closure_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//closure:js"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates Closure protobuf `.js` files",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "closure_proto_library",
				Kind:             "proto",
				Implementation:   closureProtoLibraryRuleTemplate,
				WorkspaceExample: closureLibraryWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates a Closure library with compiled protobuf `.js` files using `closure_js_library` from `rules_closure`",
				Attrs:            aspectProtoCompileAttrs,
			},
		},
	}
}
