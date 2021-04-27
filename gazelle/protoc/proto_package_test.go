package protoc

import (
	"fmt"

	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/emicklei/proto"
)

const exampleDir = "proto/test"

func exampleProtoFile() *ProtoFile {
	file := NewProtoFile(exampleDir, "test.proto")
	file.pkg = proto.Package{
		Name: "proto.test",
	}
	file.messages = append(file.messages, proto.Message{
		Name: "Foo",
	})
	file.imports = append(file.imports,
		proto.Import{Filename: "google/protobuf/any.proto"},
		proto.Import{Filename: "foo/foo.proto"},
	)
	return file
}

func exampleProtoLibraryRule() *rule.Rule {
	rule := rule.NewRule("proto_library", "test_proto")
	rule.SetAttr("deps", []string{"//foo:foo_proto"})
	return rule
}

func exampleProtoLibrary() ProtoLibrary {
	return &OtherProtoLibrary{
		rule:  exampleProtoLibraryRule(),
		files: []*ProtoFile{exampleProtoFile()}}
}

func exampleProtoPackageConfig() *ProtoPackageConfig {
	c := NewProtoPackageConfig()
	if err := c.ParseDirectives(exampleDir, withDirectives(
		"proto_plugin", "py_proto label @fake//proto/plugin",
		"proto_plugin", "py_proto enabled true",
		"proto_lang", "py plugin py_proto",
		"proto_lang", "py enabled true",
	)); err != nil {
		panic("bad config: " + err.Error())
	}
	return c
}

func exampleProtoPackage() *protoPackage {
	return newProtoPackage(
		exampleDir,
		exampleProtoPackageConfig(),
		exampleProtoLibrary(),
	)
}

func ExampleProtoPackageRules() {
	printRules(exampleProtoPackage().Rules())
	// Output:
	// proto_compile(
	//     name = "test_py_compile",
	//     genfiles = ["test_pb2.py"],
	//     plugins = ["@fake//proto/plugin"],
	//     proto = "test_proto",
	// )
}

func printRules(rules []*rule.Rule) {
	file := rule.EmptyFile(exampleDir, "")
	for _, r := range rules {
		r.Insert(file)
	}
	fmt.Println(string(file.Format()))
}

func printRule(r *rule.Rule) {
	file := rule.EmptyFile(exampleDir, "")
	r.Insert(file)
	fmt.Println(string(file.Format()))
}

// func exampleGrpcLibraryRule() *rule.Rule {
// 	rule := rule.NewRule("proto_library", "greeter_proto")
// 	rule.SetAttr("deps", []string{
// 		"//rosetta/rosetta:common_proto",
// 		"@com_google_protobuf//:any_proto",
// 	})
// 	return rule
// }

// func exampleGrpcFile() *ProtoFile {
// 	file := NewProtoFile("test", "greeter.proto")
// 	file.imports = append(file.imports,
// 		&proto.Import{Filename: "google/protobuf/any.proto"},
// 		&proto.Import{Filename: "rosetta/rosetta/common.proto"},
// 	)
// 	file.services = append(file.services,
// 		&proto.Service{Name: "Greeter"},
// 	)
// 	return file
// }

// func exampleGrpcLibrary() ProtoLibrary {
// 	return &OtherProtoLibrary{rule: exampleGrpcLibraryRule(), files: []*ProtoFile{exampleGrpcFile()}}
// }
