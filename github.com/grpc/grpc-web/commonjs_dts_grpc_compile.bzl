load("//:compile.bzl", "proto_compile")

def commonjs_dts_grpc_compile(**kwargs):
    proto_compile(
		plugins = [
			str(Label("//github.com/grpc/grpc-web:commonjs_dts")),
		],
        **kwargs
	)
