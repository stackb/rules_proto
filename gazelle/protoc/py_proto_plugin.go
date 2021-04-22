package protoc

// PyProtoPlugin implements ProtoPlugin for the built-in protoc python plugin.
type PyProtoPlugin struct {
}

func (p *PyProtoPlugin) ShouldApply(rel string, cfg *protoPackageConfig, lib ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// GeneratedSrcs implements part of the ProtoPlugin interface
func (p *PyProtoPlugin) GeneratedSrcs(rel string, cfg *protoPackageConfig, lib ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			srcs = append(srcs, f.Basename+"_pb2.py")
		}
	}
	return srcs
}
