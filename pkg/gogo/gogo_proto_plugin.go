package protoc

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const gogoGrpcPluginOption = "plugins=grpc"

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

// Name implements part of the Plugin interface.
func (p *GogoPlugin) Name() string {
	return "gogo:protobuf:" + p.Variant
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
		if goPackage, _, ok := protoc.GoPackageOption(f.Options()); ok {
			base = path.Join(goPackage, base)
		} else if pkg.Name != "" {
			base = path.Join(strings.ReplaceAll(pkg.Name, ".", "/"), base)
		}
		srcs = append(srcs, base+".pb.go")
	}
	return srcs
}

// Options implements part of the optional PluginOptionsProvider
// interface.  If the library contains services, apply the grpc plugin.
func (p *GogoPlugin) Options(rel string, c protoc.PackageConfig, lib protoc.ProtoLibrary) []string {
	cfg, ok := c.Plugin(p.Variant)
	if !ok {
		panic("unable to access the plugin config: " + p.Variant)
	}

	// if the configuration specifically states that we don't want grpc, return
	// early
	if want, ok := cfg.Options[gogoGrpcPluginOption]; ok && !want {
		return nil
	}

	for _, f := range lib.Files() {
		if f.HasServices() {
			return []string{gogoGrpcPluginOption}
		}
	}

	return nil
}
