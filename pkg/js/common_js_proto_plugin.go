package protoc

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocJsPlugin{})
}

// ProtocJsPlugin implements Plugin for the built-in js plugin.
type ProtocJsPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocJsPlugin) Name() string {
	return "protoc:js"
}

// Label implements part of the Plugin interface.
func (p *ProtocJsPlugin) Label() label.Label {
	return label.New("build_stack_rules_proto", "plugin/protoc", "js")
}

// ShouldApply implements part of the Plugin interface.
func (p *ProtocJsPlugin) ShouldApply(ctx *protoc.PluginContext) bool {
	for _, f := range ctx.ProtoLibrary.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// Outputs implements part of the Plugin interface.
func (p *ProtocJsPlugin) Outputs(ctx *protoc.PluginContext) []string {
	srcs := make([]string, 0)
	for _, f := range ctx.ProtoLibrary.Files() {
		base := f.Name
		pkg := f.Package()
		if pkg.Name != "" {
			base = path.Join(strings.ReplaceAll(pkg.Name, ".", "/"), base)
		}
		if f.HasMessages() || f.HasEnums() {
			srcs = append(srcs, base+"_pb.js")
		}
	}
	return srcs
}
