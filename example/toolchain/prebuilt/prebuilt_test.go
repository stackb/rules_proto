package prebuilt

import (
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: `
-- BUILD.bazel --
# empty file
`,
		ModuleFileSuffix: `
# TODO(pcj): why do we need a bazel_dep for rules_go?  It would seem to be setup by go_bazel_test already
bazel_dep(name = "rules_go", version = "0.57.0", repo_name = "io_bazel_rules_go")
go_sdk = use_extension("@io_bazel_rules_go//go:extensions.bzl", "go_sdk")
go_sdk.download(version = "1.23.1")`,
	})
}

func TestRun(t *testing.T) {
	if err := bazel_testing.RunBazel("run", "@build_stack_rules_proto//toolchain:protoc.exe"); err != nil {
		t.Fatal(err)
	}
}
