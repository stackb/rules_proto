package protoc

import "log"

func init() {
	MustRegisterProtoPlugin("cc_proto", &CcProtoPlugin{})
}

// CcProtoPlugin implements ProtoPlugin for the built-in protoc python plugin.
type CcProtoPlugin struct{}

// ShouldApply implements part of the ProtoPlugin interface.
func (p *CcProtoPlugin) ShouldApply(rel string, cfg *ProtoPackageConfig, lib ProtoLibrary) bool {
	log.Printf("checking %d files...", len(lib.Files()))
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
		if f.HasMessages() || f.HasEnums() {
			srcs = append(srcs, f.Name+"_pb.cc")
		}
	}
	return srcs
}
