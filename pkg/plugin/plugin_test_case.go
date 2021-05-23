package plugin

import (
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

type PluginTestCase struct {
	// The name of the proto file to mock parse.  If not set, defaults to 'test.proto'
	Filename string
	// The relative package path
	Rel string
	// Optional directives for the package config
	Directives []rule.Directive
	// The input proto file source.  "syntax = proto3" will be automatically prepended.
	Input string
	// The expected value for "ShouldApply"
	ShouldApply bool
	// The expected outputs for "Outputs"
	Outputs []string
}

func (tc *PluginTestCase) Run(t *testing.T, subject protoc.Plugin) {
	filename := tc.Filename
	if filename == "" {
		filename = "test.proto"
	}

	f := protoc.NewFile(tc.Rel, filename)
	in := "syntax = \"proto3\";\n\n" + tc.Input
	if err := f.ParseReader(strings.NewReader(in)); err != nil {
		t.Fatalf("unparseable proto file: %s: %v", tc.Input, err)
	}
	c := protoc.NewPackageConfig()
	if err := c.ParseDirectives(tc.Rel, tc.Directives); err != nil {
		t.Fatalf("bad directives: %v", err)
	}
	r := rule.NewRule("proto_library", "test_proto")
	lib := protoc.NewOtherProtoLibrary(r, f)
	// p := protoc.NewPackage(tc.Rel, c, lib)

	apply := subject.ShouldApply(tc.Rel, *c, lib)
	if apply != tc.ShouldApply {
		t.Errorf("%T.ShouldApply: want %t, got %t", subject, apply, tc.ShouldApply)
	}

	outputs := subject.Outputs(tc.Rel, *c, lib)
	if len(tc.Outputs) != len(outputs) {
		t.Fatalf("%T.Outputs: want %d, got %d", subject, len(tc.Outputs), len(outputs))
	}

	for i, got := range outputs {
		want := tc.Outputs[i]
		if want != got {
			t.Errorf("%T.Outputs[%d]: want %q, got %q", subject, i, want, got)
		}
	}
}

func PluginTestCases(t *testing.T, subject protoc.Plugin, cases map[string]PluginTestCase) {
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tc.Run(t, subject)
		})
	}
}
