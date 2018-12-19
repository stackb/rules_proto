tag := gcr.io/stack-build/rules_proto/bazel:0.20.0

build:
	docker build . -t ${tag}

push: build
	docker push ${tag}

