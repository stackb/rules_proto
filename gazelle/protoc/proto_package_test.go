package protoc

import (
	"fmt"

	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/emicklei/proto"
)

func exampleProtoLibraryRule() *rule.Rule {
	rule := rule.NewRule("proto_library", "foo_proto")
	rule.SetAttr("deps", []string{"//rosetta/rosetta:common_proto"})
	return rule
}

func exampleGrpcLibraryRule() *rule.Rule {
	rule := rule.NewRule("proto_library", "greeter_proto")
	rule.SetAttr("deps", []string{
		"//rosetta/rosetta:common_proto",
		"@com_google_protobuf//:any_proto",
	})
	return rule
}

func exampleProtoFile() *ProtoFile {
	file := NewProtoFile("test", "foo.proto")
	file.messages = append(file.messages, &proto.Message{
		Name: "Foo",
	})
	file.imports = append(file.imports,
		&proto.Import{Filename: "google/protobuf/any.proto"},
		&proto.Import{Filename: "rosetta/rosetta/common.proto"},
	)
	return file
}

func exampleGrpcFile() *ProtoFile {
	file := NewProtoFile("test", "greeter.proto")
	file.imports = append(file.imports,
		&proto.Import{Filename: "google/protobuf/any.proto"},
		&proto.Import{Filename: "rosetta/rosetta/common.proto"},
	)
	file.services = append(file.services,
		&proto.Service{Name: "Greeter"},
	)
	return file
}

func exampleProtoLibraryRules() []*rule.Rule {
	return []*rule.Rule{exampleProtoLibraryRule()}
}

func exampleProtoLibrary() ProtoLibrary {
	return &OtherProtoLibrary{rule: exampleProtoLibraryRule(), files: []*ProtoFile{exampleProtoFile()}}
}

func exampleGrpcLibrary() ProtoLibrary {
	return &OtherProtoLibrary{rule: exampleGrpcLibraryRule(), files: []*ProtoFile{exampleGrpcFile()}}
}

func exampleProtoPackageConfig() *protoPackageConfig {
	c := newProtoPackageConfig()
	c.languages[ProtoDescriptorSetLanguageName] = MustLookupProtoLanguage(ProtoDescriptorSetLanguageName)
	return c
}

func exampleProtoPackage() *ProtoPackage {
	return NewProtoPackage(nil,
		"rosetta/rosetta/foo",
		exampleProtoPackageConfig(),
		exampleProtoLibraryRules(),
		exampleProtoFile())
}

func exampleRstubsProtoPackageConfig() *protoPackageConfig {
	c := newProtoPackageConfig()
	c.languages[ProtoDescriptorSetLanguageName] = MustLookupProtoLanguage(ProtoDescriptorSetLanguageName)
	return c
}

func exampleRstubsProtoPackage() *ProtoPackage {
	return NewProtoPackage(
		nil,
		"rosetta/rosetta/foo",
		exampleRstubsProtoPackageConfig(),
		exampleProtoLibraryRules(),
		exampleProtoFile())
}

func ExampleProtoPackageRules() {
	printRules(exampleProtoPackage().Rules())
	// Output:
	// proto_descriptor_set(
	//     name = "foo_proto_descriptor",
	//     visibility = ["//visibility:public"],
	//     deps = [":foo_proto"],
	// )
	//
	// py_proto_compile(
	//     name = "foo_py_proto_compile",
	//     deps = [":foo_proto"],
	// )
	//
	// proto_compile_test(
	//     name = "foo_py_proto_compile_test",
	//     srcs = ["foo_pb2.py"],
	//     rule = ":foo_py_proto_compile",
	//     visibility = ["//visibility:private"],
	// )
	//
	// py_library(
	//     name = "foo_py_library",
	//     srcs = ["foo_pb2.py"],
	//     imports = ["../../.."],
	//     visibility = ["//rosetta:__subpackages__"],
	// )
}

func ExampleProtoPackageImports() {
	for _, i := range exampleProtoPackage().Imports() {
		fmt.Println(i)
	}
	// Output:
	// proto_descriptor_set
	// py_proto_compile
	// proto_compile_test
	// py_library
}

func ExampleRstubsProtoPackageRules() {
	printRules(exampleRstubsProtoPackage().Rules())
	// Output:
	// proto_descriptor_set(
	//     name = "foo_proto_descriptor",
	//     visibility = ["//visibility:public"],
	//     deps = [":foo_proto"],
	// )
	//
	// py_proto_compile(
	//     name = "foo_py_proto_compile",
	//     deps = [":foo_proto"],
	// )
	//
	// proto_compile_test(
	//     name = "foo_py_proto_compile_test",
	//     srcs = ["foo_pb2.py"],
	//     rule = ":foo_py_proto_compile",
	//     visibility = ["//visibility:private"],
	// )
	//
	// py_rstubs_proto_compile(
	//     name = "foo_py_rstubs_proto_compile",
	//     deps = [":foo_proto"],
	// )
	//
	// proto_compile_test(
	//     name = "foo_py_rstubs_proto_compile_test",
	//     srcs = ["foo.py"],
	//     rule = ":foo_py_rstubs_proto_compile",
	//     visibility = ["//visibility:private"],
	// )
	//
	// py_library(
	//     name = "foo_py_library",
	//     srcs = [
	//         "foo.py",
	//         "foo_pb2.py",
	//     ],
	//     imports = ["../../.."],
	//     visibility = ["//rosetta:__subpackages__"],
	// )
}

func printRule(r *rule.Rule) {
	file := rule.EmptyFile("test", "")
	r.Insert(file)
	fmt.Println(string(file.Format()))
}

func printRules(rules []*rule.Rule) {
	file := rule.EmptyFile("test", "")
	for _, r := range rules {
		r.Insert(file)
	}
	fmt.Println(string(file.Format()))
}
