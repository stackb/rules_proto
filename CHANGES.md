Assorted bugfixes

* Bump version of rules_dotnet and add .bazelrc 
  (see https://github.com/bazelbuild/rules_dotnet/issues/72)
* Don't use invoke_transitive.
* Makefile cleanup.
* Test google/protobuf 3.6.2.1
* Fix grpc.js and update sha1 to that repo.
* Apply buildifier
* Relieve :deps.bzl of dart deps

passed: 

* Use remote cache
* bazel build --config=remote --remote_instance_name=main //cpp/... //java/... //objc/... //php/... //python/... //go/... 
