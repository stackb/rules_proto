// for some bizarre reason, naming this file 'protoc_js.go' makes it be ignored
// by the compiler?
package plugin

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocJsCommonPlugin{})
}

// ProtocJsCommonPlugin implements Plugin for the built-in protoc js/library plugin.
type ProtocJsCommonPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocJsCommonPlugin) Name() string {
	return "protoc:js:common"
}

// Configure implements part of the Plugin interface.
func (p *ProtocJsCommonPlugin) Configure(ctx *protoc.PluginContext, cfg *protoc.PluginConfiguration) {
	basename := strings.ToLower(ctx.ProtoLibrary.BaseName())
	library := basename + "_pb.js"
	if ctx.Rel != "" {
		library = path.Join(ctx.Rel, library)
	}

	cfg.Label = label.New("build_stack_rules_proto", "plugin/protoc", "commonjs")
	cfg.Outputs = []string{library}
	cfg.Options = []string{"import_style=commonjs"}
}
