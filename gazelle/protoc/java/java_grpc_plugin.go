package java

import "github.com/stackb/rules_proto/gazelle/protoc"

const JavaGrpcName = "java_grpc"

func init() {
	protoc.MustRegisterProtoPlugin(JavaGrpcName, &JavaGrpcPlugin{})
}

// JavaGrpcPlugin implements ProtoPlugin for the built-in protoc java plugin.
type JavaGrpcPlugin struct{}

// ShouldApply implements part of the ProtoPlugin interface.
func (p *JavaGrpcPlugin) ShouldApply(rel string, cfg protoc.ProtoPackageConfig, lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasServices() {
			return true
		}
	}
	return false
}

// GeneratedSrcs implements part of the ProtoPlugin interface.
func (p *JavaGrpcPlugin) GeneratedSrcs(rel string, cfg protoc.ProtoPackageConfig, lib protoc.ProtoLibrary) []string {
	return []string{srcjarFile(rel, lib.BaseName()+"_grpc")}
}

// GeneratedOut implements part the optional PluginOutProvider interface.
func (p *JavaGrpcPlugin) GeneratedOut(rel string, cfg protoc.ProtoPackageConfig, lib protoc.ProtoLibrary) string {
	return srcjarFile(rel, lib.BaseName()+"_grpc")
}
