load("//:compile.bzl", "proto_compile")

def php_grpc_compile(**kwargs):
    proto_compile(
		plugins = [
			str(Label("//php:php")),
			str(Label("//php:grpc_php")),
		],
        **kwargs
	)
