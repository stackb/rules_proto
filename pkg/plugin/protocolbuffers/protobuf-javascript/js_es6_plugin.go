package protobuf_javascript

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&JsES6Plugin{})
}

// JsES6Plugin implements Plugin for the built-in protoc js/library plugin.
type JsES6Plugin struct{}

// Name implements part of the Plugin interface.
func (p *JsES6Plugin) Name() string {
	return "protocolbuffers:protobuf-javascript:es6"
}

// Configure implements part of the Plugin interface.
func (p *JsES6Plugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	basename := strings.ToLower(ctx.ProtoLibrary.BaseName())
	library := basename + "_pb.js"
	if ctx.Rel != "" {
		library = path.Join(ctx.Rel, library)
	}

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "protocolbuffers/protobuf-javascript", "es6"),
		Outputs: []string{library},
	}
}
