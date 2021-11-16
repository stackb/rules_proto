package protoc

import (
	"fmt"

	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bazelbuild/bazel-gazelle/config"
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
	return NewOtherProtoLibrary(nil, exampleProtoLibraryRule(), exampleFile())
}

func examplePackageConfig() *PackageConfig {
	emptyConfig := &config.Config{}
	c := NewPackageConfig(emptyConfig)
	if err := c.ParseDirectives(exampleDir, withDirectives(
		"proto_rule", "proto_compile implementation stackb:rules_proto:proto_compile",
		"proto_plugin", "fake_proto implementation protoc:fake",
		"proto_plugin", "fake_proto enabled true",
		"proto_language", "fake plugin fake_proto",
		"proto_language", "fake enabled true",
		"proto_language", "fake rule proto_compile",
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

func ExamplePackage() {
	printRules(examplePackage().Rules())
	// Output:
	// proto_compile(
	//     name = "test_fake_compile",
	//     outputs = ["test_fake.pb.go"],
	//     plugins = ["@build_stack_rules_proto//plugin/builtin:fake"],
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
