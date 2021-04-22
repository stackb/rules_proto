package protoc

import (
	"fmt"
	"strings"
)

func init() {
	MustRegisterProtoPlugin("closure_js_proto", &JsProtoPlugin{})
}

// JsProtoPlugin implements ProtoPlugin for the built-in protoc python plugin.
type JsProtoPlugin struct{}

// ShouldApply implements part of the ProtoPlugin interface.
func (p *JsProtoPlugin) ShouldApply(rel string, cfg *ProtoPackageConfig, lib ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// GeneratedSrcs implements part of the ProtoPlugin interface.
func (p *JsProtoPlugin) GeneratedSrcs(rel string, cfg *ProtoPackageConfig, lib ProtoLibrary) []string {
	base := strings.ToLower(lib.BaseName())
	return []string{base + ".js"}
	// for _, f := range lib.Files() {
	// 	genfiles := make([]string, 0)
	// 	for _, m := range f.messages {
	// 		base := strings.ToLower(m.Name)
	// 		genfiles = append(genfiles, base+".js")
	// 	}
	// 	for _, e := range f.enums {
	// 		base := strings.ToLower(e.Name)
	// 		genfiles = append(genfiles, base+".js")
	// 	}

	// 	srcs = append(srcs, genfiles...)
	// }
	// return srcs
}

// GeneratedOptions implements part of the optional PluginOptionsProvider
// interface.
func (p *JsProtoPlugin) GeneratedOptions(rel string, c *ProtoPackageConfig, lib ProtoLibrary) []string {
	library := fmt.Sprintf("library=%s/%s", rel, lib.BaseName())
	return []string{library}
}
