package plugintest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/stackb/rules_proto/pkg/protoc"
)

// Cases is a utility function that runs a mapping of test cases.
func Cases(t *testing.T, subject protoc.Plugin, cases map[string]Case) {
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tc.Run(t, subject)
		})
	}
}

// Case holds the inputs and expected outputs for black-box testing of
// a plugin implementation.
type Case struct {
	// The base name of the proto file to mock parse.  If not set, defaults to 'test' ('test.proto')
	Basename string
	// The relative package path
	Rel string
	// The Configuration
	// Optional directives for the package config
	Directives []rule.Directive
	// The input proto file source.  "syntax = proto3" will be automatically prepended.
	Input string
	// The expected value for the final configuration state
	Configuration *protoc.PluginConfiguration
}

func (tc *Case) Run(t *testing.T, subject protoc.Plugin) {
	execrootDir := os.Getenv("TEST_TMPDIR")
	defer os.RemoveAll(execrootDir)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	protocPath := filepath.Join(cwd, "protoc")

	basename := tc.Basename
	if basename == "" {
		basename = "test"
	}
	filename := basename + ".proto"

	f := protoc.NewFile(tc.Rel, filename)
	in := "syntax = \"proto3\";\n\n" + tc.Input
	if err := f.ParseReader(strings.NewReader(in)); err != nil {
		t.Fatalf("unparseable proto file: %s: %v", tc.Input, err)
	}
	c := protoc.NewPackageConfig()
	if err := c.ParseDirectives(tc.Rel, tc.Directives); err != nil {
		t.Fatalf("bad directives: %v", err)
	}
	r := rule.NewRule("proto_library", basename+"_proto")
	pluginConfig, ok := c.Plugin(tc.Configuration.Name)
	if !ok {
		t.Fatalf("unregistered plugin configuration %q (%+v)", subject.Name(), c)
	}
	lib := protoc.NewOtherProtoLibrary(r, f)
	ctx := &protoc.PluginContext{
		Rel:           tc.Rel,
		ProtoLibrary:  lib,
		PackageConfig: *c,
		PluginConfig:  pluginConfig,
	}

	got := subject.Configure(ctx)

	if got == nil != tc.Configuration.Skip {
		t.Errorf("%T.Skip: want %t, got %t", subject, tc.Configuration.Skip, got.Skip)
	}
	outputs := got.Outputs
	if len(tc.Configuration.Outputs) != len(outputs) {
		t.Fatalf("%T.Outputs: want %d, got %d (%v)", subject, len(tc.Configuration.Outputs), len(outputs), outputs)
	}
	for i, got := range outputs {
		want := tc.Configuration.Outputs[i]
		if want != got {
			t.Errorf("%T.Outputs[%d]: want %q, got %q", subject, i, want, got)
		}
	}

	// relDir is the location where the proto files are written.  A BUILD.bazel
	// file containing the proto_library would normally be here.
	relDir := filepath.Join(".", tc.Rel)
	if err := os.MkdirAll(filepath.Join(execrootDir, relDir), os.ModePerm); err != nil {
		t.Fatalf("relDir: %v", err)
	}
	if err := ioutil.WriteFile(filepath.Join(execrootDir, relDir, filename), []byte(in), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// gendir is the root location where we expect generated files to be
	// written.  Within a bazel action, this is the execroot unless the "Out"
	// setting is configured.
	outDir := "."
	if got.Out != "" {
		outDir = filepath.Join(outDir, got.Out)
		if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
			t.Fatalf("outDir: %v", err)
		}
	}

	args := []string{
		"--proto_path=.", // this is the default (just a reminder)  The execroot is '.'
		fmt.Sprintf("--%s_out=%s:%s", tc.Configuration.Name, strings.Join(got.Options, ","), outDir),
		filepath.Join(tc.Rel, filename),
	}

	t.Log("protoc args:", args)

	mustExecProtoc(t, protocPath, execrootDir, args...)

	actuals := mustListFiles(t, execrootDir)
	if len(tc.Configuration.Outputs) != len(actuals) {
		t.Fatalf("%T.Actuals: want %d, got %d: %v", subject, len(tc.Configuration.Outputs), len(actuals), actuals)
	}

	for _, want := range outputs {
		realpath := filepath.Join(execrootDir, want)
		if !fileExists(realpath) {
			t.Errorf("expected file %q was not produced: (got %v)", want, actuals)
		}
	}

}
