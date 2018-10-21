load("//:compile.bzl", "proto_compile")

def cpp_proto_compile(**kwargs):
    proto_compile(
		plugins = [
			str(Label("//cpp:cpp")),
		],
        **kwargs
	)
