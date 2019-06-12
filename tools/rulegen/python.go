package main

var pythonGrpcCompileUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()`)

var pythonProtoLibraryUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@io_bazel_rules_python//python:pip.bzl", "pip_import", "pip_repositories")

pip_repositories()

pip_import(
    name = "protobuf_py_deps",
    requirements = "@build_stack_rules_proto//python/requirements:protobuf.txt",
)

load("@protobuf_py_deps//:requirements.bzl", protobuf_pip_install = "pip_install")

protobuf_pip_install()`)

var pythonGrpcLibraryUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@io_bazel_rules_python//python:pip.bzl", "pip_import", "pip_repositories")

pip_repositories()

pip_import(
    name = "protobuf_py_deps",
    requirements = "@build_stack_rules_proto//python/requirements:protobuf.txt",
)

load("@protobuf_py_deps//:requirements.bzl", protobuf_pip_install = "pip_install")

protobuf_pip_install()

pip_import(
    name = "grpc_py_deps",
    requirements = "@build_stack_rules_proto//python/requirements:grpc.txt",
)

load("@grpc_py_deps//:requirements.bzl", grpc_pip_install = "pip_install")

grpc_pip_install()`)

var pythonProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:python_proto_compile.bzl", "python_proto_compile")
load("@protobuf_py_deps//:requirements.bzl", protobuf_requirements = "all_requirements")

def python_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    python_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    native.py_library(
        name = name,
        srcs = [name_pb],
        deps = protobuf_requirements,
        # This magically adds REPOSITORY_NAME/PACKAGE_NAME/{name_pb} to PYTHONPATH
        imports = [name_pb],
        visibility = visibility,
    )

# Alias
py_proto_library = python_proto_library`)

var pythonGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:python_grpc_compile.bzl", "python_grpc_compile")
load("@protobuf_py_deps//:requirements.bzl", protobuf_requirements = "all_requirements")
load("@grpc_py_deps//:requirements.bzl", grpc_requirements = "all_requirements")

def python_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    python_grpc_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    native.py_library(
        name = name,
        srcs = [name_pb],
        deps = depset(protobuf_requirements + grpc_requirements).to_list(),
        # This magically adds REPOSITORY_NAME/PACKAGE_NAME/{name_pb} to PYTHONPATH
        imports = [name_pb],
        visibility = visibility,
    )

# Alias
py_grpc_library = python_grpc_library`)

func makePython() *Language {
	return &Language{
		Dir:   "python",
		Name:  "python",
		Flags: commonLangFlags,
		Notes: aspectLangNotes,
		Rules: []*Rule{
			&Rule{
				Name:           "python_proto_compile",
				Implementation: aspectRuleTemplate,
				Plugins:        []string{"//python:python"},
				Usage:          usageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates *.py protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "python_grpc_compile",
				Implementation: aspectRuleTemplate,
				Plugins:        []string{"//python:python", "//python:grpc_python"},
				Usage:          pythonGrpcCompileUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates *.py protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
				Flags: []*Flag{
					{
						Name:     "incompatible_enable_cc_toolchain_resolution",
						Value:    "false",
						Category: "build",
					},
				},
			},
			&Rule{
				Name:           "python_proto_library",
				Implementation: pythonProtoLibraryRuleTemplate,
				Usage:          pythonProtoLibraryUsageTemplate,
				Example:        protoLibraryExampleTemplate,
				Doc:            "Generates *.py protobuf library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "python_grpc_library",
				Implementation: pythonGrpcLibraryRuleTemplate,
				Usage:          pythonGrpcLibraryUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates *.py protobuf+gRPC library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
