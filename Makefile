routeguide_compile:
	bazel build \
		//closure/example/proto:routeguide \
		//cpp/example/proto:routeguide \
		//csharp/example/proto:routeguide \
		//go/example/proto:routeguide \
		//java/example/proto:routeguide \
		//node/example/proto:routeguide \
		//objc/example/proto:routeguide \
		//php/example/proto:routeguide \
		//python/example/proto:routeguide \
		//ruby/example/proto:routeguide \
		//contrib/improbable-eng/ts-protoc-gen/example/proto:routeguide \
		//contrib/grpc/grpc-web/example/proto:routeguide \

#		//dart/example/proto:routeguide \

routeguide_clients:
	bazel build \
		//closure/example/routeguide/client \
		//cpp/example/routeguide:client \
		//python/example/routeguide:client \
		//java/example/routeguide:client \
		//go/example/routeguide/client \
		//contrib/grpc/grpc-web/example/routeguide/client:bundle \
		//contrib/stackb/grpc.js/example/routeguide/client:bundle 

routeguide_servers:
	bazel build \
		//python/example/routeguide:server \
		//cpp/example/routeguide:server \
		//java/example/routeguide:server \
		//go/example/routeguide/server \

routeguide_tests:
	bazel test \
		//closure/example/routeguide/... \
		//contrib/stackb/grpc.js/example/routeguide/... \
		//contrib/grpc/grpc-web/example/routeguide/... \
		//cpp/example/routeguide/... \
		//java/example/routeguide/... \
		//go/example/routeguide/... \
