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

// GeneratedOptions implements part of the optional PluginOptionsProvider
// interface.  If the library contains services, apply the grpc plugin.
func (p *GogoProtoPlugin) GeneratedOptions(rel string, c *ProtoPackageConfig, lib ProtoLibrary) []string {
	for _, f := range lib.Files() {
		if f.HasServices() {
			return []string{"plugins=grpc"}
		}
	}
	return nil
}
