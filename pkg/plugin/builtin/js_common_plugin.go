// for some bizarre reason, naming this file 'protoc_js.go' makes it be ignored
// by the compiler?
package builtin

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&JsCommonPlugin{})
}

// JsCommonPlugin implements Plugin for the built-in protoc js/library plugin.
type JsCommonPlugin struct{}

// Name implements part of the Plugin interface.
func (p *JsCommonPlugin) Name() string {
	return "builtin:js:common"
}

// Configure implements part of the Plugin interface.
func (p *JsCommonPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	basename := strings.ToLower(ctx.ProtoLibrary.BaseName())
	library := basename + "_pb.js"
	if ctx.Rel != "" {
		library = path.Join(ctx.Rel, library)
	}

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/builtin", "commonjs"),
		Outputs: []string{library},
		Options: append(ctx.PluginConfig.GetOptions(), "import_style=commonjs"),
	}
}
