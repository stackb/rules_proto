package testtools

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/testtools"
	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"

	"github.com/stackb/rules_proto/tools/protogen"
)

type ProtoRuleTestSuite struct {
	t                    *testing.T
	Rule                 *protogen.ProtoRule
	BzlFile              string
	WorkspaceExampleFile string
	BuildExampleFile     string
}

// LoadProtoRuleTestSuite creates new test suite based on the filename
// conventions of the output files of the proto_rule starlark rule.
func LoadProtoRuleTestSuite(t *testing.T, name string) *ProtoRuleTestSuite {
	rule, err := protogen.FromJSONFile(name + ".json")
	if err != nil {
		t.Fatal(err)
	}
	return &ProtoRuleTestSuite{
		t:                    t,
		Rule:                 rule,
		BzlFile:              mustRead(t, name+".bzl"),
		WorkspaceExampleFile: mustRead(t, name+".WORKSPACE"),
		BuildExampleFile:     mustRead(t, name+".BUILD"),
	}
}

func (r *ProtoRuleTestSuite) Run() {
	t := r.t
	ListFiles(".")

	gazellePath, ok := bazel.FindBinary("internal/gazellebinarytest", "gazelle_go_x")
	if !ok {
		t.Fatal("could not find gazelle binary")
	}

	files := []testtools.FileSpec{
		{Path: "WORKSPACE", Content: r.WorkspaceExampleFile},
		{Path: "BUILD.bazel", Content: r.BuildExampleFile},
		{Path: r.Rule.Name + ".bzl", Content: r.BzlFile},
		{Path: "foo.proto", Content: `syntax = "proto3";`},
	}
	dir, cleanup := testtools.CreateFiles(t, files)
	defer cleanup()

	bazel_testing.
	cmd := exec.Command(fmt.Sprintf("bazel build //:%s", r.Rule.Name))
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

// ListFiles - convenience debugging function to log the files under a given dir
func ListFiles(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("%v\n", err)
			return err
		}
		if info.Mode()&os.ModeSymlink > 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			log.Printf("%s -> %s", path, link)
			return nil
		}

		log.Println(path)
		return nil
	})
}

func mustRead(t *testing.T, filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("could not read %s: %v", filename, err)
	}
	return strings.TrimSpace(string(bytes))
}
