package main

var testHeader = `
package main

import (
	"testing"
	"os"

	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
	"github.com/google/go-cmp/cmp"
)

var (
	// allow use of os package in other tests
	_os_Remove = os.Remove
	// allow use of cmp package in other tests
	_cmp_Diff = cmp.Diff
)

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: txtar,
	})
}
`
