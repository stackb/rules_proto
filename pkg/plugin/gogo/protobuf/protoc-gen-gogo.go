package protobuf

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
		protoc.Plugins().MustRegisterPlugin(&GogoPlugin{variant})
	}
}

// GogoPlugin implements Plugin for the the gogo_* family of plugins.
type GogoPlugin struct {
	variant string
}

// Name implements part of the Plugin interface.
func (p *GogoPlugin) Name() string {
	return "gogo:protobuf:protoc-gen-" + p.variant
}

// Configure implements part of the Plugin interface.
func (p *GogoPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !p.shouldApply(ctx.ProtoLibrary) {
		return nil
	}

	grpcOptions := p.grpcOptions(ctx.Rel, ctx.PluginConfig, ctx.ProtoLibrary)
	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/gogo/protobuf", "protoc-gen-"+p.variant),
		Outputs: p.outputs(ctx.ProtoLibrary),
		Options: append(grpcOptions, ctx.PluginConfig.GetOptions()...),
	}
}

func (p *GogoPlugin) shouldApply(lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() || f.HasServices() {
			return true
		}
	}
	return false
}

func (p *GogoPlugin) outputs(lib protoc.ProtoLibrary) []string {
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

func (p *GogoPlugin) grpcOptions(rel string, cfg protoc.LanguagePluginConfig, lib protoc.ProtoLibrary) []string {
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
