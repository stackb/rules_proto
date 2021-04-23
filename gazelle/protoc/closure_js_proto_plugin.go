package protoc

import (
	"fmt"
	"path"
	"strings"
)

func init() {
	MustRegisterProtoPlugin("closure_js_proto", &ClosureJsProtoPlugin{})
}

// ClosureJsProtoPlugin implements ProtoPlugin for the built-in protoc python plugin.
type ClosureJsProtoPlugin struct{}

// ShouldApply implements part of the ProtoPlugin interface.
func (p *ClosureJsProtoPlugin) ShouldApply(rel string, cfg *ProtoPackageConfig, lib ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// GeneratedSrcs implements part of the ProtoPlugin interface.
func (p *ClosureJsProtoPlugin) GeneratedSrcs(rel string, cfg *ProtoPackageConfig, lib ProtoLibrary) []string {
	base := strings.ToLower(lib.BaseName())
	return []string{path.Join(rel, base+".js")}
}

// GeneratedOptions implements part of the optional PluginOptionsProvider
// interface.
func (p *ClosureJsProtoPlugin) GeneratedOptions(rel string, c *ProtoPackageConfig, lib ProtoLibrary) []string {
	library := fmt.Sprintf("library=%s/%s", rel, lib.BaseName())
	return []string{library}
}
