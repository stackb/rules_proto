package protoc

import (
	"fmt"
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ClosureJsPlugin{})
}

// ClosureJsPlugin implements Plugin for the built-in protoc python plugin.
type ClosureJsPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ClosureJsPlugin) Name() string {
	return "protoc:js:closure"
}

// Label implements part of the Plugin interface.
func (p *ClosureJsPlugin) Label() label.Label {
	return label.New("build_stack_rules_proto", "protocolbuffers/protobuf", "closure_js_plugin")
}

// ShouldApply implements part of the Plugin interface.
func (p *ClosureJsPlugin) ShouldApply(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// Outputs implements part of the Plugin interface.
func (p *ClosureJsPlugin) Outputs(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) []string {
	base := strings.ToLower(lib.BaseName())
	return []string{path.Join(rel, base+".js")}
}

// Options implements part of the optional PluginOptionsProvider
// interface.
func (p *ClosureJsPlugin) Options(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) []string {
	library := fmt.Sprintf("library=%s/%s", rel, lib.BaseName())
	return []string{library}
}
