package protoc

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/emicklei/proto"
	"github.com/stackb/rules_proto/pkg/protoc"
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
		protoc.Plugins().MustRegisterPlugin(variant, &GogoPlugin{Variant: variant})
	}
}

// GogoPlugin implements Plugin for the the gogo_* family of plugins.
type GogoPlugin struct {
	Variant string
}

// Label implements part of the Plugin interface.
func (p *GogoPlugin) Label() label.Label {
	return label.New("build_stack_rules_proto", "gogo/protobuf", p.Variant+"_plugin")
}

func (p *GogoPlugin) ShouldApply(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() || f.HasServices() {
			return true
		}
	}
	return false
}

// Outputs implements part of the Plugin interface
func (p *GogoPlugin) Outputs(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		base := f.Name
		pkg := f.Package()
		// see https://github.com/gogo/protobuf/blob/master/protoc-gen-gogo/generator/generator.go#L347
		if goPackage, _, ok := goPackageOption(f.Options()); ok {
			base = path.Join(goPackage, base)
		} else if pkg.Name != "" {
			base = path.Join(strings.ReplaceAll(pkg.Name, ".", "/"), base)
		}
		if f.HasMessages() || f.HasEnums() {
			srcs = append(srcs, base+".pb.go")
		}
	}
	return srcs
}

// Options implements part of the optional PluginOptionsProvider
// interface.  If the library contains services, apply the grpc plugin.
func (p *GogoPlugin) Options(rel string, c *protoc.PackageConfig, lib protoc.ProtoLibrary) []string {
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
func goPackageOption(options []proto.Option) (string, string, bool) {
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
