package protoc

func exampleProtoRule() *ProtoRule {
	return NewProtoRule(exampleProtoLibrary(), "py", "proto", "library")
}

func ExampleProtoRuleRule() {
	printRule(exampleProtoRule().Rule())
	// Output:
	// py_proto_library(
	//     name = "foo_py_proto_library",
	//     deps = [":foo_proto"],
	// )
}
