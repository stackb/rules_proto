package protoc

import (
	"fmt"
	"path"
	"strings"
)

func init() {
	MustRegisterPlugin("closure_js_proto", &ClosureJsPlugin{})
}

// ClosureJsPlugin implements Plugin for the built-in protoc python plugin.
type ClosureJsPlugin struct{}

// ShouldApply implements part of the Plugin interface.
func (p *ClosureJsPlugin) ShouldApply(rel string, cfg PackageConfig, lib ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// Outputs implements part of the Plugin interface.
func (p *ClosureJsPlugin) Outputs(rel string, cfg PackageConfig, lib ProtoLibrary) []string {
	base := strings.ToLower(lib.BaseName())
	return []string{path.Join(rel, base+".js")}
}

// Options implements part of the optional PluginOptionsProvider
// interface.
func (p *ClosureJsPlugin) Options(rel string, cfg PackageConfig, lib ProtoLibrary) []string {
	library := fmt.Sprintf("library=%s/%s", rel, lib.BaseName())
	return []string{library}
}
