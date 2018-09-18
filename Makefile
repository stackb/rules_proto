routeguide:
	bazel build closure/example/proto:routeguide \
	 cpp/example/proto:routeguide \
	 csharp/example/proto:routeguide \
	 dart/example/proto:routeguide \
	 go/example/proto:routeguide \
	 java/example/proto:routeguide \
	 node/example/proto:routeguide \
	 objc/example/proto:routeguide \
	 php/example/proto:routeguide \
	 python/example/proto:routeguide \
	 ruby/example/proto:routeguide

routeguide_test:
	bazel test go/example/routeguide/server:go_default_test 