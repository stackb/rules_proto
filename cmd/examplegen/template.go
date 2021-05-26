package main

var header = `
package main

import (
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: txtar,
	})
}

func TestBuild(t *testing.T) {
	if err := bazel_testing.RunBazel("build", "..."); err != nil {
		t.Fatal(err)
	}
}

`

var workspace = `
local_repository(
    name = "build_stack_rules_proto",
    path = "../build_stack_rules_proto",
)

# ====================================================
# Toolchains
# ====================================================
register_toolchains("@build_stack_rules_proto//protoc:toolchain")

# ====================================================
# External Dependencies
# ====================================================

load(
    "@build_stack_rules_proto//:deps.bzl",
    "bazel_gazelle",
    "bazel_skylib",
    "com_github_grpc_grpc",
    "com_google_protobuf",
    "io_bazel_rules_go",
    "build_bazel_rules_swift",
    "rules_proto",
    "rules_jvm_external",
    "io_grpc_grpc_java",
    "rules_python",
    "zlib",
)

io_bazel_rules_go()

bazel_gazelle()

bazel_skylib()

com_github_grpc_grpc()

build_bazel_rules_swift()

io_grpc_grpc_java()

rules_jvm_external()

# ==================================================
# Go
# ==================================================

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

# ==================================================
# Go Tool Deps
# ==================================================

load("@build_stack_rules_proto//:go_deps.bzl", "go_deps")

go_deps()

# ==================================================
# Protobuf Core
# ==================================================

rules_proto()

rules_python()

com_google_protobuf()

zlib()

`
