package java

import (
	"log"
	"path"
	"strconv"

	"github.com/emicklei/proto"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const JavaProtoName = "java_proto"

func init() {
	protoc.Plugins().MustRegisterPlugin(JavaProtoName, &JavaPlugin{})
}

// JavaPlugin implements Plugin for the built-in protoc java plugin.
type JavaPlugin struct{}

// ShouldApply implements part of the Plugin interface.
func (p *JavaPlugin) ShouldApply(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// Outputs implements part of the Plugin interface.
func (p *JavaPlugin) Outputs(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) []string {
	// srcs := make([]string, 0)
	// for _, f := range lib.Files() {
	// 	outputs := make([]string, 0)

	// 	multipleFiles := javaMultipleFiles(f.GetOptions())
	// 	if multipleFiles {
	// 		for _, m := range f.messages {
	// 			outputs = append(outputs, m.Name+".java", m.Name+"OrBuilder.java")
	// 		}
	// 		for _, e := range f.enums {
	// 			outputs = append(outputs, e.Name+".java")
	// 		}
	// 	}

	// 	outerClassname := javaOuterClassname(f.GetOptions())
	// 	if outerClassname == "" {
	// 		outerClassname = strcase.ToCamel(f.Name)
	// 	}

	// 	if hasMatchingMessageName(f.messages, outerClassname) {
	// 		outerClassname = outerClassname + "OuterClass"
	// 	}
	// 	outputs = append(outputs, outerClassname+".java")

	// 	pkg := javaPackage(f.GetOptions())
	// 	if pkg != "" {
	// 		prefix := strings.ReplaceAll(pkg, ".", "/")
	// 		// If we are going to generate code outside of the package, bazel
	// 		// will not be happy.  Use a srcjar instead.
	// 		if prefix != rel {
	// 			outputs = []string{f.Name + ".srcjar"}
	// 		}
	// 	}

	// 	srcs = append(srcs, outputs...)
	// }
	return []string{srcjarFile(rel, lib.BaseName())}
}

// Out implements part the optional PluginOutProvider interface.
func (p *JavaPlugin) Out(rel string, cfg *protoc.PackageConfig, lib protoc.ProtoLibrary) string {
	return srcjarFile(rel, lib.BaseName())
}

func srcjarFile(dir, name string) string {
	return path.Join(dir, name+".srcjar")
}

// javaMultipleFiles is a utility function to seek for the java_outer_classname
// option.  If the option was not found the bool return argument is *true*.
func javaMultipleFiles(options []*proto.Option) bool {
	for _, opt := range options {
		if opt.Name != "java_multiple_files" {
			continue
		}
		value, err := strconv.ParseBool(opt.Constant.Source)
		if err != nil {
			log.Println("could not parse java_files_option value: %v", err)
			return false // since we did not parse, fallback to default
		}
		return value
	}
	return false
}

// javaOuterClassname is a utility function to seek for the java_outer_classname
// option.  If the option was not found the return argument is the empty string.
func javaOuterClassname(options []*proto.Option) string {
	for _, opt := range options {
		if opt.Name != "java_outer_classname" {
			continue
		}
		return opt.Constant.Source
	}
	return ""
}

// javaPackage is a utility function to seek for the java_package option.  If
// the option was not found the return argument is the empty string.
func javaPackage(options []*proto.Option) string {
	for _, opt := range options {
		if opt.Name != "java_package" {
			continue
		}
		return opt.Constant.Source
	}
	return ""
}

// hasMatchingMessageName is a utility function that searches for a top-level
// messsage matching the given name.
func hasMatchingMessageName(messages []*proto.Message, name string) bool {
	for _, message := range messages {
		if message.Name == name {
			return true
		}
	}
	return false
}
