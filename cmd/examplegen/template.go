package main

var testHeader = `
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
