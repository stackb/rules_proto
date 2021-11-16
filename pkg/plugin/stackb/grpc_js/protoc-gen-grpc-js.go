package grpc_js

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocGenGrpcJs{})
}

// ProtocGenGrpcJs implements Plugin for the https://github.com/stackb/grpc.js/tree/master/protoc-gen-grpc-js tool.
type ProtocGenGrpcJs struct{}

// Name implements part of the Plugin interface.
func (p *ProtocGenGrpcJs) Name() string {
	return "stackb:grpc.js:protoc-gen-grpc-js"
}

// Configure implements part of the Plugin interface.
func (p *ProtocGenGrpcJs) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !protoc.HasServices(ctx.ProtoLibrary.Files()...) {
		return nil
	}
	basename := strings.ToLower(ctx.ProtoLibrary.BaseName())
	jsFile := basename + ".grpc.js"
	if ctx.Rel != "" {
		jsFile = path.Join(ctx.Rel, jsFile)
	}
	out := "out=" + jsFile
	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/stackb/grpc_js", "protoc-gen-grpc-js"),
		Outputs: []string{jsFile},
		Options: append(ctx.PluginConfig.GetOptions(), out),
	}
}
