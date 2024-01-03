package standard

import (
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: txtar,
	})
}

func TestRun(t *testing.T) {
	if err := bazel_testing.RunBazel("run", "@com_google_protobuf//:protoc"); err != nil {
		t.Fatal(err)
	}
}

var txtar = `
-- WORKSPACE --
local_repository(
    name = "build_stack_rules_proto",
    path = "../build_stack_rules_proto",
)

register_toolchains("@build_stack_rules_proto//toolchain:standard")

load("@build_stack_rules_proto//deps:protobuf_core_deps.bzl", "protobuf_core_deps")

protobuf_core_deps()

load("@rules_python//python:repositories.bzl", "py_repositories")

py_repositories()

-- BUILD.bazel --
# empty file
`
