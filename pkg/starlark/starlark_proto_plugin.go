package starlark

import (
	"context"
	"fmt"
	"os"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
	"github.com/stripe/skycfg"
	"github.com/stripe/skycfg/gogocompat"
)

// StarlarkPlugin implements Plugin for a set of plugins that can be progratically driven by starlark.
type StarlarkPlugin struct {
	config *skycfg.Config
}

func init() {
	filename := ""
	config, err := skycfg.Load(context.Background(), filename, skycfg.WithProtoRegistry(gogocompat.ProtoRegistry()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading %q: %v\n", filename, err)
		os.Exit(1)
	}

	protoc.Plugins().MustRegisterPlugin("skycfg", &StarlarkPlugin{config})
}

// Label implements part of the Plugin interface.
func (p *StarlarkPlugin) Label() label.Label {
	var l label.Label
	return l
}

// ShouldApply implements part of the Plugin interface.
func (p *StarlarkPlugin) ShouldApply(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) bool {
	return false
}

// Outputs implements part of the Plugin interface.
func (p *StarlarkPlugin) Outputs(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) []string {
	srcs := make([]string, 0)
	return srcs
}
