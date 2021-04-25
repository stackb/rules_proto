package protoc

import (
	"path"
	"strings"
)

func init() {
	MustRegisterProtoPlugin("cc_proto", &CcProtoPlugin{})
}

// CcProtoPlugin implements ProtoPlugin for the built-in protoc python plugin.
type CcProtoPlugin struct{}

// ShouldApply implements part of the ProtoPlugin interface.
func (p *CcProtoPlugin) ShouldApply(rel string, cfg *ProtoPackageConfig, lib ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// GeneratedSrcs implements part of the ProtoPlugin interface.
func (p *CcProtoPlugin) GeneratedSrcs(rel string, cfg *ProtoPackageConfig, lib ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		base := f.Name
		pkg := f.Package()
		if pkg.Name != "" {
			base = path.Join(protoPackagePath(pkg.Name), base)
		}
		if f.HasMessages() || f.HasEnums() {
			srcs = append(srcs, base+".pb.cc", base+".pb.h")
		}
	}
	return srcs
}

func protoPackagePath(pkg string) string {
	return strings.ReplaceAll(pkg, ".", "/")
}
