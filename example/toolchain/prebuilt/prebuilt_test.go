package prebuilt

import (
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: `
-- WORKSPACE --
local_repository(
	name = "build_stack_rules_proto",
	path = "../build_stack_rules_proto",
)

register_toolchains("@build_stack_rules_proto//toolchain:prebuilt")

load("@build_stack_rules_proto//deps:prebuilt_protoc_deps.bzl", "prebuilt_protoc_deps")

prebuilt_protoc_deps()

-- BUILD.bazel --
# empty file
`,
	})
}

func TestRun(t *testing.T) {
	if err := bazel_testing.RunBazel("run", "@build_stack_rules_proto//toolchain:protoc.exe"); err != nil {
		t.Fatal(err)
	}
}
