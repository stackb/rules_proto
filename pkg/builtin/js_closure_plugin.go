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
	protoc.Plugins().MustRegisterPlugin(&JsClosurePlugin{})
}

// JsClosurePlugin implements Plugin for the built-in protoc js/library plugin.
type JsClosurePlugin struct{}

// Name implements part of the Plugin interface.
func (p *JsClosurePlugin) Name() string {
	return "builtin:js:closure"
}

// Configure implements part of the Plugin interface.
func (p *JsClosurePlugin) Configure(ctx *protoc.PluginContext, cfg *protoc.PluginConfiguration) {
	basename := strings.ToLower(ctx.ProtoLibrary.BaseName())
	library := path.Join(ctx.Rel, basename+".js")

	cfg.Label = label.New("build_stack_rules_proto", "plugin/builtin", "closurejs")
	cfg.Outputs = []string{library}
	// cfg.Options = []string{"import_style=closure", fmt.Sprintf("library=%s", strings.TrimSuffix(library, filepath.Ext(library)))}
}
