build:
	docker build . -t gcr.io/stack-build/rules_proto/rbe-ubuntu16-04-plus:latest

push: build
	docker push gcr.io/stack-build/rules_proto/rbe-ubuntu16-04-plus:latest
	# gcr.io/stack-build/rules_proto/rbe-ubuntu16-04-plus@sha256:4564ad5704ee1b31000dd0728ee5271af20434000a503065e7ad76e80b6ac4e5