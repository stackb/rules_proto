package main

var testHeader = `
package main

import (
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
	"github.com/google/go-cmp/cmp"
)

func TestMain(m *testing.M) {
	// allow use of cmp in other tests, justify import here
	cmp.Diff("", "")
	
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: txtar,
	})
}
`
