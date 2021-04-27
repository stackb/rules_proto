package protoc

import "path"

func init() {
	MustRegisterPlugin("common_js_proto", &CommonJsPlugin{})
}

// CommonJsPlugin implements Plugin for the built-in js plugin with
// commonjs option.
type CommonJsPlugin struct{}

// ShouldApply implements part of the Plugin interface.
func (p *CommonJsPlugin) ShouldApply(rel string, cfg PackageConfig, lib ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// Outputs implements part of the Plugin interface.
func (p *CommonJsPlugin) Outputs(rel string, cfg PackageConfig, lib ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		base := f.Name
		pkg := f.Package()
		if pkg.Name != "" {
			base = path.Join(GoPackagePath(pkg.Name), base)
		}
		if f.HasMessages() || f.HasEnums() {
			srcs = append(srcs, base+"_pb.js")
		}
	}
	return srcs
}
