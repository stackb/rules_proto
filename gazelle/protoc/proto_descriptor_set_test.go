package protoc

func exampleProtoDescriptorSet() *ProtoDescriptorSet {
	return NewProtoDescriptorSet(exampleProtoLibrary())
}

func ExampleProtoDescriptorSetRule() {
	printRule(exampleProtoDescriptorSet().Rule())
	// Output:
	// proto_descriptor_set(
	//     name = "foo_proto_descriptor",
	//     visibility = ["//visibility:public"],
	//     deps = [":foo_proto"],
	// )
}
