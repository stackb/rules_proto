package protoc

import (
	"path"
	"strings"
)

func init() {
	MustRegisterPlugin("cc_proto", &CcPlugin{})
}

// CcPlugin implements Plugin for the built-in protoc python plugin.
type CcPlugin struct{}

// ShouldApply implements part of the Plugin interface.
func (p *CcPlugin) ShouldApply(rel string, cfg PackageConfig, lib ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// Outputs implements part of the Plugin interface.
func (p *CcPlugin) Outputs(rel string, cfg PackageConfig, lib ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		base := f.Name
		pkg := f.Package()
		if pkg.Name != "" {
			base = path.Join(GoPackagePath(pkg.Name), base)
		}
		if f.HasMessages() || f.HasEnums() {
			srcs = append(srcs, base+".pb.cc", base+".pb.h")
		}
	}
	return srcs
}
