package protoc

import (
	"path"
	"strings"

	"github.com/emicklei/proto"
)

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
		base := f.Name
		// see https://github.com/gogo/protobuf/blob/master/protoc-gen-gogo/generator/generator.go#L347
		if goPackage, _, ok := goPackageOption(f.GetOptions()); ok {
			base = path.Join(goPackage, base)
		}
		if f.HasMessages() || f.HasEnums() {
			srcs = append(srcs, base+".pb.go")
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

// goPackageOption is a utility function to seek for the go_package option and
// split it.  If present the return values will be populated with the importpath
// and alias (e.g. github.com/foo/bar/v1;bar -> "github.com/foo/bar/v1", "bar").
// If the option was not found the bool return argument is false.
func goPackageOption(options []*proto.Option) (string, string, bool) {
	for _, opt := range options {
		if opt.Name != "go_package" {
			continue
		}
		parts := strings.SplitN(opt.Constant.Source, ";", 2)
		switch len(parts) {
		case 0:
			return "", "", true
		case 1:
			return parts[0], "", true
		case 2:
			return parts[0], parts[1], true
		default:
			return parts[0], strings.Join(parts[1:], ";"), true
		}
	}

	return "", "", false
}
