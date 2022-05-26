package grpcgatewayts

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&protocGenGrpcGatewayTsPlugin{})
}

// protocGenGrpcGatewayTsPlugin implements Plugin for protoc-gen-grpc-gateway-ts.
type protocGenGrpcGatewayTsPlugin struct{}

// Name implements part of the Plugin interface.
func (p *protocGenGrpcGatewayTsPlugin) Name() string {
	return "grpc-ecosystem:grpc-gateway-ts:protoc-gen-grpc-gateway-ts"
}

// Configure implements part of the Plugin interface.
func (p *protocGenGrpcGatewayTsPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !p.shouldApply(ctx.ProtoLibrary) {
		return nil
	}
	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/grpc-ecosystem/grpc-gateway-ts", "protoc-gen-grpc-gateway-ts"),
		Outputs: p.outputs(ctx.Rel, ctx.ProtoLibrary, p.fetchModuleFilename(ctx.PluginConfig)),
		Options: p.options(ctx.Rel, ctx.PluginConfig),
	}
}

func (p *protocGenGrpcGatewayTsPlugin) shouldApply(lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasServices() {
			return true
		}
	}
	return false
}

func (p *protocGenGrpcGatewayTsPlugin) options(rel string, cfg protoc.LanguagePluginConfig) []string {
	options := make([]string, 0)
	for _, o := range cfg.GetOptions() {
		if strings.HasPrefix(o, "fetch_module_directory=") {
			continue
		}
		options = append(options, o)
	}
	options = append(options, "fetch_module_directory="+rel)
	return options
}

func (p *protocGenGrpcGatewayTsPlugin) fetchModuleFilename(cfg protoc.LanguagePluginConfig) string {
	const filenameArg = "fetch_module_filename="
	const defaultName = "fetch.pb.ts"
	for _, o := range cfg.GetOptions() {
		if strings.HasPrefix(o, filenameArg) {
			return strings.TrimPrefix(o, filenameArg)
		}
	}
	return defaultName
}

func (p *protocGenGrpcGatewayTsPlugin) outputs(rel string, lib protoc.ProtoLibrary, fetch string) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		if !f.HasServices() {
			continue
		}
		base := f.Name
		if rel != "" {
			base = path.Join(rel, base)
			fetch = path.Join(rel, fetch)
		}
		srcs = append(srcs, base+".pb.ts", fetch)
	}
	return srcs
}
