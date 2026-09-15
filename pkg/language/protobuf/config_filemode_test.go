package protobuf

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	gproto "github.com/bazelbuild/bazel-gazelle/language/proto"
)

// Regression: getOrCreatePackageConfig must point the cloned PackageConfig's
// embedded *config.Config at the CURRENT directory's config, not the parent's.
// Otherwise IsProtoFileMode reads the parent's proto-language config and
// misses a child directory that opted into `# gazelle:proto file`.
func TestGetOrCreatePackageConfig_RebindsConfig(t *testing.T) {
	pl := &protobufLang{name: "protobuf"}

	// Parent: default mode.
	parent := &config.Config{Exts: map[string]interface{}{}}
	parent.Exts["proto"] = &gproto.ProtoConfig{Mode: gproto.DefaultMode}
	pl.getOrCreatePackageConfig(parent)

	// Child: clone via gazelle's Config.Clone, then proto-lang flips mode to FileMode
	// for this dir only.
	child := parent.Clone()
	childProto := &gproto.ProtoConfig{Mode: gproto.FileMode}
	child.Exts["proto"] = childProto

	cfg := pl.getOrCreatePackageConfig(child)
	if cfg.Config != child {
		t.Fatalf("cloned PackageConfig.Config not rebound to child: got %p, want %p", cfg.Config, child)
	}

	// IsProtoFileMode-equivalent check.
	if gproto.GetProtoConfig(cfg.Config).Mode != gproto.FileMode {
		t.Errorf("expected FileMode via cfg.Config, got %v", gproto.GetProtoConfig(cfg.Config).Mode)
	}
}

// The manifest must use the paths Gazelle actually loads, including the source
// root fallback used when Gazelle runs inside a fetched repository.
func TestConfigInputsManifest(t *testing.T) {
	for _, sourceRootFallback := range []bool{false, true} {
		t.Run(fmt.Sprint(sourceRootFallback), func(t *testing.T) {
			dir := t.TempDir()
			workDir := dir
			if sourceRootFallback {
				workDir = filepath.Join(dir, "external")
				if err := os.Mkdir(workDir, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(workDir, "DO_NOT_BUILD_HERE"), []byte(dir), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			// Unique registry keys allow repeated tests in the same process.
			name := filepath.Base(filepath.Dir(dir)) + fmt.Sprint(sourceRootFallback) + ".star"
			filename := filepath.Join(dir, name)
			code := `
protoc.Plugin(name = "yaml_plugin", configure = lambda ctx: None)
protoc.Plugin(name = "flag_plugin", configure = lambda ctx: None)
protoc.Rule(
    name = "yaml_rule",
    load_info = lambda: None,
    kind_info = lambda: None,
    provide_rule = lambda rctx, pctx: None,
)
`
			if err := os.WriteFile(filename, []byte(code), 0o644); err != nil {
				t.Fatal(err)
			}
			yaml := fmt.Sprintf("starlarkPlugins:\n  - %s%%yaml_plugin\nstarlarkRules:\n  - %s%%yaml_rule\n", name, name)
			cfgFile := filepath.Join(dir, "config.yaml")
			if err := os.WriteFile(cfgFile, []byte(yaml), 0o644); err != nil {
				t.Fatal(err)
			}
			out := filepath.Join(dir, "inputs.json")
			c := config.New()
			c.WorkDir = workDir
			pl := NewProtobufLang("protobuf")
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			pl.RegisterFlags(fs, "update", c)
			if err := fs.Parse([]string{"-proto_configs", cfgFile, "-proto_plugin", name + "%flag_plugin", "-proto_config_inputs_out", out}); err != nil {
				t.Fatal(err)
			}
			if err := pl.CheckFlags(fs, c); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(out)
			if err != nil {
				t.Fatal(err)
			}
			var got []string
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if want := []string{filename}; !reflect.DeepEqual(got, want) {
				t.Fatalf("manifest = %v, want %v", got, want)
			}
		})
	}
}

func TestConfigInputsManifestEmpty(t *testing.T) {
	c := config.New()
	pl := NewProtobufLang("protobuf")
	out := filepath.Join(t.TempDir(), "inputs.json")
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	pl.RegisterFlags(fs, "update", c)
	if err := fs.Parse([]string{"-proto_config_inputs_out", out}); err != nil {
		t.Fatal(err)
	}
	if err := pl.CheckFlags(fs, c); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "[]" {
		t.Fatalf("manifest = %s, want []", data)
	}
	pl.configInputsOutFile = filepath.Join(out, "cannot-write.json")
	if err := pl.CheckFlags(fs, c); err == nil {
		t.Fatal("expected an error writing the manifest")
	}
}
