package protoc

import "path"

func init() {
	MustRegisterPlugin("cc_grpc", &CcGrpcPlugin{})
}

// CcGrpcPlugin implements Plugin for the built-in protoc python plugin.
type CcGrpcPlugin struct{}

// ShouldApply implements part of the Plugin interface.
func (p *CcGrpcPlugin) ShouldApply(rel string, cfg PackageConfig, lib ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasServices() {
			return true
		}
	}
	return false
}

// Outputs implements part of the Plugin interface.
func (p *CcGrpcPlugin) Outputs(rel string, cfg PackageConfig, lib ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		base := f.Name
		pkg := f.Package()
		if pkg.Name != "" {
			base = path.Join(GoPackagePath(pkg.Name), base)
		}
		if f.HasServices() {
			srcs = append(srcs, base+".grpc.pb.cc", base+".grpc.pb.h")
		}
	}
	return srcs
}
