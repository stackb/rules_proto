package protoc

import (
	"path"
	"strings"

	"github.com/iancoleman/strcase"
)

func init() {
	MustRegisterProtoPlugin("java_grpc", &JavaGrpcPlugin{})
}

// JavaGrpcPlugin implements ProtoPlugin for the built-in protoc java plugin.
type JavaGrpcPlugin struct{}

// ShouldApply implements part of the ProtoPlugin interface.
func (p *JavaGrpcPlugin) ShouldApply(rel string, cfg *ProtoPackageConfig, lib ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// GeneratedSrcs implements part of the ProtoPlugin interface.
func (p *JavaGrpcPlugin) GeneratedSrcs(rel string, cfg *ProtoPackageConfig, lib ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		genfiles := make([]string, 0)

		multipleFiles := javaMultipleFiles(f.GetOptions())
		if multipleFiles {
			for _, m := range f.messages {
				genfiles = append(genfiles, m.Name+".java", m.Name+"OrBuilder.java")
			}
			for _, e := range f.enums {
				genfiles = append(genfiles, e.Name+".java")
			}
		}

		outerClassname := javaOuterClassname(f.GetOptions())
		if outerClassname == "" {
			outerClassname = strcase.ToCamel(f.Name)
		}

		if hasMatchingMessageName(f.messages, outerClassname) {
			outerClassname = outerClassname + "OuterClass"
		}
		genfiles = append(genfiles, outerClassname+".java")

		pkg := javaPackage(f.GetOptions())
		if pkg != "" {
			prefix := strings.ReplaceAll(pkg, ".", "/")
			// If we are going to generate code outside of the package, bazel
			// will not be happy.  Use a srcjar instead.
			if prefix != rel {
				genfiles = []string{f.Name + "_grpc.srcjar"}
			}
		}

		srcs = append(srcs, genfiles...)
	}
	return srcs
}

// GeneratedOut implements part the optional PluginOutProvider interface.
func (p *JavaGrpcPlugin) GeneratedOut(rel string, cfg *ProtoPackageConfig, lib ProtoLibrary) string {
	return path.Join("{BIN_DIR}", "{PACKAGE}", lib.BaseName()+"_grpc.srcjar")
}
