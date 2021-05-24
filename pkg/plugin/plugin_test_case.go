package plugin

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

func PluginTestCases(t *testing.T, subject protoc.Plugin, cases map[string]PluginTestCase) {
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tc.Run(t, subject)
		})
	}
}

type PluginTestCase struct {
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

type PluginConfigurationOption func(c *protoc.PluginConfiguration)

func WithConfiguration(options ...PluginConfigurationOption) *protoc.PluginConfiguration {
	c := &protoc.PluginConfiguration{}
	for _, opt := range options {
		opt(c)
	}
	return c
}

func WithSkip(skip bool) PluginConfigurationOption {
	return func(c *protoc.PluginConfiguration) {
		c.Skip = skip
	}
}

func WithName(name string) PluginConfigurationOption {
	return func(c *protoc.PluginConfiguration) {
		c.Name = name
	}
}

func WithOutputs(outputs ...string) PluginConfigurationOption {
	return func(c *protoc.PluginConfiguration) {
		c.Outputs = outputs
	}
}

func WithDirectives(items ...string) (d []rule.Directive) {
	if len(items)%2 != 0 {
		panic("directive list must be a sequence of key/value pairs")
	}
	if len(items) < 2 {
		return
	}
	for i := 1; i < len(items); i = i + 2 {
		d = append(d, rule.Directive{Key: items[i-1], Value: items[i]})
	}
	return
}

func (tc *PluginTestCase) Run(t *testing.T, subject protoc.Plugin) {
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

	got := &protoc.PluginConfiguration{}
	subject.Configure(ctx, got)

	if got.Skip != tc.Configuration.Skip {
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

func mustExecProtoc(t *testing.T, protoc, dir string, args ...string) {
	cmd := exec.Command(protoc, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("protoc exec error: %v\n\n%s", err, out)
	}
}

// mustListFiles - convenience debugging function to log the files under a given
// dir, excluding proto files and the extra binaries here.
func mustListFiles(t *testing.T, dir string) []string {
	files := make([]string, 0)

	if err := filepath.Walk(dir, func(relname string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatal(err)
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(relname) == ".proto" {
			return nil
		}
		files = append(files, relname)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	return files
}

// fileExists checks if a file exists and is not a directory before we try using
// it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if info == nil {
		return false
	}
	return !info.IsDir()
}

// listFiles - convenience debugging function to log the files under a given dir
func listFiles(dir string) error {
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
