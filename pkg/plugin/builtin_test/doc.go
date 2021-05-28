// plugin builtin_test is for testing of the pkg/builtin plugins.
package builtin_test

// what I want: create testdata examples that get their gazelle run, then take
// the BUILD.out and other files and put that in a txtar, and run that as a
// go_bazel_test.  So here Is what I think we need.

// 1. Write a rule that takes a filegroup(testdata/**) glob and produces a list
// of corresponding {DIRNAME}_bazel_test.go files, one foreach dir in the
// testdata dir.
//
// 2. Execute this as a go_bazel_test, as a macro.
