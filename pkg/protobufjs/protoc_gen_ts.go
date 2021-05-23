// ts_proto.go implements a protoc.Plugin for https://github.com/stephenh/ts-proto.
package golang

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	// ProtocGenTsName is the name the plugin is registered under.
	ProtocGenTsName = "protoc_gen_ts"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(ProtocGenTsName, &ProtocGenTsPlugin{})
}

// ProtocGenTsPlugin implements Plugin for the the protobufjs ;pbts' plugin.
type ProtocGenTsPlugin struct{}

// Label implements part of the Plugin interface.
func (p *ProtocGenTsPlugin) Label() label.Label {
	return label.New("build_stack_rules_proto", "protobufjs/protobufjs", "protoc_gen_ts")
}

func (p *ProtocGenTsPlugin) ShouldApply(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// Outputs implements part of the Plugin interface
func (p *ProtocGenTsPlugin) Outputs(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		if !(f.HasMessages() || f.HasEnums()) {
			continue
		}
		base := f.Name
		pkg := f.Package()
		if pkg.Name != "" {
			base = path.Join(strings.ReplaceAll(pkg.Name, ".", "/"), base)
		}
		// // see https://github.com/gogo/protobuf/blob/master/protoc-gen-gogo/generator/generator.go#L347
		// if goPackage, _, ok := protoc.GoPackageOption(f.Options()); ok {
		// 	base = path.Join(goPackage, base)
		// } else if pkg.Name != "" {
		// base = path.Join(strings.ReplaceAll(pkg.Name, ".", "/"), base)
		srcs = append(srcs, base+".ts")
	}
	return srcs
}