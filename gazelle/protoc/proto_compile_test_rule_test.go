package protoc

func exampleProtoCompileTestRule() *ProtoCompileTest {
	return &ProtoCompileTest{NewProtoRule(exampleGrpcLibrary(), "py", "grpc", "compile")}
}

func ExampleProtoCompileTestRule() {
	printRule(exampleProtoCompileTestRule().Rule())
	// Output:
	// proto_compile_test(
	//     name = "greeter_py_grpc_compile_test",
	//     srcs = ["greeter_pb2_grpc.py"],
	//     rule = ":greeter_py_grpc_compile",
	//     visibility = ["//visibility:private"],
	// )
}
