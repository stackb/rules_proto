package scalapb

import (
	"path"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const ScalaPBPluginName = "scalapb:scalapb:protoc-gen-scala"

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocGenScalaPlugin{})
}

// ProtocGenScalaPlugin implements Plugin for the scala plugin.
type ProtocGenScalaPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocGenScalaPlugin) Name() string {
	return ScalaPBPluginName
}

// Configure implements part of the Plugin interface.
func (p *ProtocGenScalaPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	srcjar := ctx.ProtoLibrary.BaseName() + "_scala.srcjar"
	if ctx.Rel != "" {
		srcjar = path.Join(ctx.Rel, srcjar)
	}
	options := ctx.PluginConfig.GetOptions()

	// if the plugin has 'grpc' statically configured, but the proto_library
	// does not contain services, remove it.
	if filtered, ok := removeAll(options, "grpc"); ok {
		if !protoc.HasServices(ctx.ProtoLibrary.Files()...) {
			options = filtered
		}
	}

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/scalapb/scalapb", "protoc-gen-scala"),
		Outputs: []string{srcjar},
		Options: options,
	}
}

// removeAll returns a copy of the slice will all elements matching 'key'
// removed.  If at least item was removed, return true.
func removeAll(src []string, key string) ([]string, bool) {
	dst := make([]string, 0)
	found := false
	for _, item := range src {
		if item == key {
			found = true
			continue
		}
		dst = append(dst, item)
	}
	return dst, found
}
