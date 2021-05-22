package golang

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	// PbjsName is the name the plugin is registered under.
	PbjsName = "pbjs"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(PbjsName, &PbjsPlugin{})
}

// PbjsPlugin implements Plugin for the the protobufjs ;pbjs' plugin.
type PbjsPlugin struct{}

// Label implements part of the Plugin interface.
func (p *PbjsPlugin) Label() label.Label {
	return label.New("build_stack_rules_proto", "protobufjs/protobufjs", "pbjs")
}

func (p *PbjsPlugin) ShouldApply(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// Outputs implements part of the Plugin interface
func (p *PbjsPlugin) Outputs(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		if !(f.HasMessages() || f.HasEnums()) {
			continue
		}
		base := f.Name
		// pkg := f.Package()
		// // see https://github.com/gogo/protobuf/blob/master/protoc-gen-gogo/generator/generator.go#L347
		// if goPackage, _, ok := protoc.GoPackageOption(f.Options()); ok {
		// 	base = path.Join(goPackage, base)
		// } else if pkg.Name != "" {
		// base = path.Join(strings.ReplaceAll(pkg.Name, ".", "/"), base)
		srcs = append(srcs, base+".pb.js")
	}
	return srcs
}
