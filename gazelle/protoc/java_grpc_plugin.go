package protoc

func init() {
	MustRegisterProtoPlugin("java_grpc", &JavaGrpcPlugin{})
}

// JavaGrpcPlugin implements ProtoPlugin for the built-in protoc java plugin.
type JavaGrpcPlugin struct{}

// ShouldApply implements part of the ProtoPlugin interface.
func (p *JavaGrpcPlugin) ShouldApply(rel string, cfg *ProtoPackageConfig, lib ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasServices() {
			return true
		}
	}
	return false
}

// GeneratedSrcs implements part of the ProtoPlugin interface.
func (p *JavaGrpcPlugin) GeneratedSrcs(rel string, cfg *ProtoPackageConfig, lib ProtoLibrary) []string {
	return []string{srcjarFile(rel, lib.BaseName()+"_grpc")}
}

// GeneratedOut implements part the optional PluginOutProvider interface.
func (p *JavaGrpcPlugin) GeneratedOut(rel string, cfg *ProtoPackageConfig, lib ProtoLibrary) string {
	return srcjarFile(rel, lib.BaseName()+"_grpc")
}
