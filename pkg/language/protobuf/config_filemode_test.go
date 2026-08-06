package protobuf

import (
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
