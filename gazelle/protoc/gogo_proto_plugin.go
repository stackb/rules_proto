package protoc

func init() {
	for _, variant := range []string{
		"combo",
		"gogo",
		"gogofast",
		"gogofaster",
		"gogoslick",
		"gogotypes",
		"gostring",
	} {
		MustRegisterProtoPlugin(variant, &GogoProtoPlugin{Variant: variant})
	}
}

// GogoProtoPlugin implements ProtoPlugin for the the gogo_* family of plugins.
type GogoProtoPlugin struct {
	Variant string
}

func (p *GogoProtoPlugin) ShouldApply(rel string, cfg *ProtoPackageConfig, lib ProtoLibrary) bool {
	return true
}

// GeneratedSrcs implements part of the ProtoPlugin interface
func (p *GogoProtoPlugin) GeneratedSrcs(rel string, cfg *ProtoPackageConfig, lib ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			srcs = append(srcs, f.Name+".pb.go")
		}
	}
	return srcs
}
