package testtools

import (
	"os"
	"os/exec"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/testtools"
	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

type RuleTestSuite struct {
	Name string
}

func (r *RuleTestSuite) Run(t *testing.T) {
	gazellePath, ok := bazel.FindBinary("internal/gazellebinarytest", "gazelle_go_x")
	if !ok {
		t.Fatal("could not find gazelle binary")
	}

	files := []testtools.FileSpec{
		{Path: "WORKSPACE"},
		{Path: "BUILD.bazel", Content: "# gazelle:prefix example.com/test"},
		{Path: "foo.go", Content: "package foo"},
		{Path: "foo.proto", Content: `syntax = "proto3";`},
	}
	dir, cleanup := testtools.CreateFiles(t, files)
	defer cleanup()

	cmd := exec.Command(gazellePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	testtools.CheckFiles(t, dir, []testtools.FileSpec{{
		Path: "BUILD.bazel",
		Content: `
load("@io_bazel_rules_go//go:def.bzl", "go_library")
# gazelle:prefix example.com/test
go_library(
    name = "test",
    srcs = ["foo.go"],
    importpath = "example.com/test",
    visibility = ["//visibility:public"],
)
x_library(name = "x_default_library")
`,
	}})
}
