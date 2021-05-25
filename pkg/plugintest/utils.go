package plugintest

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/stackb/rules_proto/pkg/protoc"
)

// PluginConfigurationOption modifies a configuration in-place
type PluginConfigurationOption func(c *protoc.PluginConfiguration)

// WithConfiguration creates a new PluginConfiguration and applies all the given
// options.
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

// WithSkip assigns th Skip field.
func WithName(name string) PluginConfigurationOption {
	return func(c *protoc.PluginConfiguration) {
		c.Name = name
	}
}

// WithSkip assigns th Skip field.
func WithOutputs(outputs ...string) PluginConfigurationOption {
	return func(c *protoc.PluginConfiguration) {
		c.Outputs = outputs
	}
}

// WithSkip assigns th Skip field.
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
