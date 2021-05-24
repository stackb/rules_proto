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

	protoc.Plugins().MustRegisterPlugin(&StarlarkPlugin{config})
}

// Name implements part of the Plugin interface.
func (p *StarlarkPlugin) Name() string {
	return "starlark"
}

// Label implements part of the Plugin interface.
func (p *StarlarkPlugin) Label() label.Label {
	var l label.Label
	return l
}

// ShouldApply implements part of the Plugin interface.
func (p *StarlarkPlugin) ShouldApply(ctx *protoc.PluginContext) bool {
	return false
}

// Outputs implements part of the Plugin interface.
func (p *StarlarkPlugin) Outputs(ctx *protoc.PluginContext) []string {
	srcs := make([]string, 0)
	return srcs
}
