# A collection of targets that build routeguide clients
clients:
	bazel build \
		//android/example/routeguide:client \
		//closure/example/routeguide/client \
		//cpp/example/routeguide:client \
		//csharp/example/routeguide:client \
		//dart/example/routeguide:client \
		//go/example/routeguide/client \
		//java/example/routeguide:client \
		//node/example/routeguide:client \
		//python/example/routeguide:client \
		//ruby/example/routeguide:client \
		//scala/example/routeguide:client \
		//github.com/grpc/grpc-web/example/routeguide/closure:bundle \
		//github.com/stackb/grpc.js/example/routeguide/client:bundle \

# 		//rust/example/routeguide:client \

# A collection of targets that build routeguide servers
servers:
	bazel build \
		//cpp/example/routeguide:server \
		//csharp/example/routeguide:server \
		//dart/example/routeguide:server \
		//go/example/routeguide/server \
		//java/example/routeguide:server \
		//node/example/routeguide:server \
		//python/example/routeguide:server \
		//ruby/example/routeguide:server \


#		//rust/example/routeguide:server \

# A collection of test targets
tests:
	bazel test \
		//closure/example/routeguide/... \
		//github.com/stackb/grpc.js/example/routeguide/... \
		//cpp/example/routeguide/... \
		//java/example/routeguide/... \
		//go/example/routeguide/... \

# grpc-web closure test seems to crash phantomjs.  Todo move to headless-chrome
# for rules_closure.
closure_test:
	bazel test \
		//github.com/grpc/grpc-web/example/routeguide/...

all: clients servers tests
