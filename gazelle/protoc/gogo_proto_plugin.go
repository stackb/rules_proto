package protoc

func init() {
	MustRegisterProtoPlugin("gogo", &GogoProtoPlugin{Name: "gogo"})
	MustRegisterProtoPlugin("gogofast", &GogoProtoPlugin{Name: "gogofast"})
	MustRegisterProtoPlugin("gogofaster", &GogoProtoPlugin{Name: "gogofaster"})
}

// GogoProtoPlugin implements ProtoPlugin for the the gogo_* family of plugins.
type GogoProtoPlugin struct {
	Name string
}

func (p *GogoProtoPlugin) ShouldApply(rel string, cfg *protoPackageConfig, lib ProtoLibrary) bool {
	return true
}

// GeneratedSrcs implements part of the ProtoPlugin interface
func (p *GogoProtoPlugin) GeneratedSrcs(rel string, cfg *protoPackageConfig, lib ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			srcs = append(srcs, f.Name+".pb.go")
		}
	}
	return srcs
}
