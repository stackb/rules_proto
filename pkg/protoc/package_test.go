package protoc

import (
	"fmt"

	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/emicklei/proto"
)

const exampleDir = "proto/test"

func exampleFile() *File {
	file := NewFile(exampleDir, "test.proto")
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
	return NewOtherProtoLibrary(exampleProtoLibraryRule(), exampleFile())
}

func examplePackageConfig() *PackageConfig {
	c := NewPackageConfig()
	if err := c.ParseDirectives(exampleDir, withDirectives(
		"proto_plugin", "fake_proto label @fake//proto/plugin",
		"proto_plugin", "fake_proto enabled true",
		"proto_language", "fake plugin fake_proto",
		"proto_language", "fake enabled true",
	)); err != nil {
		panic("bad config: " + err.Error())
	}
	return c
}

func examplePackage() *Package {
	return NewPackage(
		exampleDir,
		examplePackageConfig(),
		exampleProtoLibrary(),
	)
}

func ExamplePackageRules() {
	printRules(examplePackage().Rules())
	// Output:
	// proto_compile(
	//     name = "test_fake_compile",
	//     outputs = ["test_fake.pb.go"],
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
