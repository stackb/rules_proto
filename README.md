# `rules_proto`

Bazel skylark rules for building protocol buffers +/- gRPC :sparkles:.

<table border="0"><tr>
<td><img src="https://bazel.build/images/bazel-icon.svg" height="180"/></td>
<td><img src="https://github.com/pubref/rules_protobuf/blob/master/images/wtfcat.png" height="180"/></td>
<td><img src="https://avatars2.githubusercontent.com/u/7802525?v=4&s=400" height="180"/></td>
</tr><tr>
<td>Bazel</td>
<td>rules_proto</td>
<td>gRPC</td>
</tr></table>



## Rules

| Status | Lang | Rule | Description
| ---    | ---: | :--- | :--- |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [android](/android) | [android_proto_compile](/android#android_proto_compile) | Generates android protobuf artifacts ([example](/example/android/android_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [android](/android) | [android_grpc_compile](/android#android_grpc_compile) | Generates android protobuf+gRPC artifacts ([example](/example/android/android_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [android](/android) | [android_proto_library](/android#android_proto_library) | Generates android protobuf library ([example](/example/android/android_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [android](/android) | [android_grpc_library](/android#android_grpc_library) | Generates android protobuf+gRPC library ([example](/example/android/android_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [closure](/closure) | [closure_proto_compile](/closure#closure_proto_compile) | Generates closure *.js protobuf+gRPC files ([example](/example/closure/closure_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [closure](/closure) | [closure_proto_library](/closure#closure_proto_library) | Generates a closure_library with compiled protobuf *.js files ([example](/example/closure/closure_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [cpp](/cpp) | [cpp_proto_compile](/cpp#cpp_proto_compile) | Generates *.h,*.cc protobuf artifacts ([example](/example/cpp/cpp_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [cpp](/cpp) | [cpp_grpc_compile](/cpp#cpp_grpc_compile) | Generates *.h,*.cc protobuf+gRPC artifacts ([example](/example/cpp/cpp_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [cpp](/cpp) | [cpp_proto_library](/cpp#cpp_proto_library) | Generates *.h,*.cc protobuf library ([example](/example/cpp/cpp_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [cpp](/cpp) | [cpp_grpc_library](/cpp#cpp_grpc_library) | Generates *.h,*.cc protobuf+gRPC library ([example](/example/cpp/cpp_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [csharp](/csharp) | [csharp_proto_compile](/csharp#csharp_proto_compile) | Generates csharp protobuf artifacts ([example](/example/csharp/csharp_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [csharp](/csharp) | [csharp_grpc_compile](/csharp#csharp_grpc_compile) | Generates csharp protobuf+gRPC artifacts ([example](/example/csharp/csharp_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [csharp](/csharp) | [csharp_proto_library](/csharp#csharp_proto_library) | Generates csharp protobuf library ([example](/example/csharp/csharp_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [csharp](/csharp) | [csharp_grpc_library](/csharp#csharp_grpc_library) | Generates csharp protobuf+gRPC library ([example](/example/csharp/csharp_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [d](/d) | [d_proto_compile](/d#d_proto_compile) | Generates d protobuf artifacts ([example](/example/d/d_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [dart](/dart) | [dart_proto_compile](/dart#dart_proto_compile) | Generates dart protobuf artifacts ([example](/example/dart/dart_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [dart](/dart) | [dart_grpc_compile](/dart#dart_grpc_compile) | Generates dart protobuf+gRPC artifacts ([example](/example/dart/dart_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [dart](/dart) | [dart_proto_library](/dart#dart_proto_library) | Generates dart protobuf library ([example](/example/dart/dart_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [dart](/dart) | [dart_grpc_library](/dart#dart_grpc_library) | Generates dart protobuf+gRPC library ([example](/example/dart/dart_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [go](/go) | [go_proto_compile](/go#go_proto_compile) | Generates *.go protobuf artifacts ([example](/example/go/go_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [go](/go) | [go_grpc_compile](/go#go_grpc_compile) | Generates *.go protobuf+gRPC artifacts ([example](/example/go/go_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [go](/go) | [go_proto_library](/go#go_proto_library) | Generates *.go protobuf library ([example](/example/go/go_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [go](/go) | [go_grpc_library](/go#go_grpc_library) | Generates *.go protobuf+gRPC library ([example](/example/go/go_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [java](/java) | [java_proto_compile](/java#java_proto_compile) | Generates a srcjar with protobuf *.java files ([example](/example/java/java_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [java](/java) | [java_grpc_compile](/java#java_grpc_compile) | Generates a srcjar with protobuf+gRPC *.java files ([example](/example/java/java_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [java](/java) | [java_proto_library](/java#java_proto_library) | Generates a jar with compiled protobuf *.class files ([example](/example/java/java_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [java](/java) | [java_grpc_library](/java#java_grpc_library) | Generates a jar with compiled protobuf+gRPC *.class files ([example](/example/java/java_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [node](/node) | [node_proto_compile](/node#node_proto_compile) | Generates node *.js protobuf artifacts ([example](/example/node/node_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [node](/node) | [node_grpc_compile](/node#node_grpc_compile) | Generates node *.js protobuf+gRPC artifacts ([example](/example/node/node_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [node](/node) | [node_proto_library](/node#node_proto_library) | Generates node *.js protobuf library ([example](/example/node/node_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [node](/node) | [node_grpc_library](/node#node_grpc_library) | Generates node *.js protobuf+gRPC library ([example](/example/node/node_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [objc](/objc) | [objc_proto_compile](/objc#objc_proto_compile) | Generates objc protobuf artifacts ([example](/example/objc/objc_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [objc](/objc) | [objc_grpc_compile](/objc#objc_grpc_compile) | Generates objc protobuf+gRPC artifacts ([example](/example/objc/objc_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [php](/php) | [php_proto_compile](/php#php_proto_compile) | Generates php protobuf artifacts ([example](/example/php/php_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [php](/php) | [php_grpc_compile](/php#php_grpc_compile) | Generates php protobuf+gRPC artifacts ([example](/example/php/php_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [python](/python) | [python_proto_compile](/python#python_proto_compile) | Generates *.py protobuf artifacts ([example](/example/python/python_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [python](/python) | [python_grpc_compile](/python#python_grpc_compile) | Generates *.py protobuf+gRPC artifacts ([example](/example/python/python_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [python](/python) | [python_proto_library](/python#python_proto_library) | Generates *.py protobuf library ([example](/example/python/python_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [python](/python) | [python_grpc_library](/python#python_grpc_library) | Generates *.py protobuf+gRPC library ([example](/example/python/python_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [ruby](/ruby) | [ruby_proto_compile](/ruby#ruby_proto_compile) | Generates *.ruby protobuf artifacts ([example](/example/ruby/ruby_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [ruby](/ruby) | [ruby_grpc_compile](/ruby#ruby_grpc_compile) | Generates *.ruby protobuf+gRPC artifacts ([example](/example/ruby/ruby_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [ruby](/ruby) | [ruby_proto_library](/ruby#ruby_proto_library) | Generates *.rb protobuf library ([example](/example/ruby/ruby_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [ruby](/ruby) | [ruby_grpc_library](/ruby#ruby_grpc_library) | Generates *.rb protobuf+gRPC library ([example](/example/ruby/ruby_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [rust](/rust) | [rust_proto_compile](/rust#rust_proto_compile) | Generates rust protobuf artifacts ([example](/example/rust/rust_proto_compile)) |
| experimental | [rust](/rust) | [rust_grpc_compile](/rust#rust_grpc_compile) | Generates rust protobuf+gRPC artifacts ([example](/example/rust/rust_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [rust](/rust) | [rust_proto_library](/rust#rust_proto_library) | Generates rust protobuf library ([example](/example/rust/rust_proto_library)) |
| experimental | [rust](/rust) | [rust_grpc_library](/rust#rust_grpc_library) | Generates rust protobuf+gRPC library ([example](/example/rust/rust_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [scala](/scala) | [scala_proto_compile](/scala#scala_proto_compile) | Generates *.scala protobuf artifacts ([example](/example/scala/scala_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [scala](/scala) | [scala_grpc_compile](/scala#scala_grpc_compile) | Generates *.scala protobuf+gRPC artifacts ([example](/example/scala/scala_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [scala](/scala) | [scala_proto_library](/scala#scala_proto_library) | Generates *.scala protobuf library ([example](/example/scala/scala_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [scala](/scala) | [scala_grpc_library](/scala#scala_grpc_library) | Generates *.scala protobuf+gRPC library ([example](/example/scala/scala_grpc_library)) |
| experimental | [swift](/swift) | [swift_proto_compile](/swift#swift_proto_compile) | Generates swift protobuf artifacts ([example](/example/swift/swift_proto_compile)) |
| experimental | [swift](/swift) | [swift_grpc_compile](/swift#swift_grpc_compile) | Generates swift protobuf+gRPC artifacts ([example](/example/swift/swift_grpc_compile)) |
| experimental | [swift](/swift) | [swift_proto_library](/swift#swift_proto_library) | Generates swift protobuf library ([example](/example/swift/swift_proto_library)) |
| experimental | [swift](/swift) | [swift_grpc_library](/swift#swift_grpc_library) | Generates swift protobuf+gRPC library ([example](/example/swift/swift_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [gogo](/github.com/gogo/protobuf) | [gogo_proto_compile](/github.com/gogo/protobuf#gogo_proto_compile) | Generates gogo protobuf artifacts ([example](/example/github.com/gogo/protobuf/gogo_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [gogo](/github.com/gogo/protobuf) | [gogo_grpc_compile](/github.com/gogo/protobuf#gogo_grpc_compile) | Generates gogo protobuf+gRPC artifacts ([example](/example/github.com/gogo/protobuf/gogo_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [gogo](/github.com/gogo/protobuf) | [gogo_proto_library](/github.com/gogo/protobuf#gogo_proto_library) | Generates gogo protobuf library ([example](/example/github.com/gogo/protobuf/gogo_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [gogo](/github.com/gogo/protobuf) | [gogo_grpc_library](/github.com/gogo/protobuf#gogo_grpc_library) | Generates gogo protobuf+gRPC library ([example](/example/github.com/gogo/protobuf/gogo_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [gogo](/github.com/gogo/protobuf) | [gogofast_proto_compile](/github.com/gogo/protobuf#gogofast_proto_compile) | Generates gogofast protobuf artifacts ([example](/example/github.com/gogo/protobuf/gogofast_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [gogo](/github.com/gogo/protobuf) | [gogofast_grpc_compile](/github.com/gogo/protobuf#gogofast_grpc_compile) | Generates gogofast protobuf+gRPC artifacts ([example](/example/github.com/gogo/protobuf/gogofast_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [gogo](/github.com/gogo/protobuf) | [gogofast_proto_library](/github.com/gogo/protobuf#gogofast_proto_library) | Generates gogofast protobuf library ([example](/example/github.com/gogo/protobuf/gogofast_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [gogo](/github.com/gogo/protobuf) | [gogofast_grpc_library](/github.com/gogo/protobuf#gogofast_grpc_library) | Generates gogofast protobuf+gRPC library ([example](/example/github.com/gogo/protobuf/gogofast_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [gogo](/github.com/gogo/protobuf) | [gogofaster_proto_compile](/github.com/gogo/protobuf#gogofaster_proto_compile) | Generates gogofaster protobuf artifacts ([example](/example/github.com/gogo/protobuf/gogofaster_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [gogo](/github.com/gogo/protobuf) | [gogofaster_grpc_compile](/github.com/gogo/protobuf#gogofaster_grpc_compile) | Generates gogofaster protobuf+gRPC artifacts ([example](/example/github.com/gogo/protobuf/gogofaster_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [gogo](/github.com/gogo/protobuf) | [gogofaster_proto_library](/github.com/gogo/protobuf#gogofaster_proto_library) | Generates gogofaster protobuf library ([example](/example/github.com/gogo/protobuf/gogofaster_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [gogo](/github.com/gogo/protobuf) | [gogofaster_grpc_library](/github.com/gogo/protobuf#gogofaster_grpc_library) | Generates gogofaster protobuf+gRPC library ([example](/example/github.com/gogo/protobuf/gogofaster_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [grpc-gateway](/github.com/grpc-ecosystem/grpc-gateway) | [gateway_grpc_compile](/github.com/grpc-ecosystem/grpc-gateway#gateway_grpc_compile) | Generates grpc-gateway *.go files ([example](/example/github.com/grpc-ecosystem/grpc-gateway/gateway_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [grpc-gateway](/github.com/grpc-ecosystem/grpc-gateway) | [gateway_swagger_compile](/github.com/grpc-ecosystem/grpc-gateway#gateway_swagger_compile) | Generates grpc-gateway swagger *.json files ([example](/example/github.com/grpc-ecosystem/grpc-gateway/gateway_swagger_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [grpc-gateway](/github.com/grpc-ecosystem/grpc-gateway) | [gateway_grpc_library](/github.com/grpc-ecosystem/grpc-gateway#gateway_grpc_library) | Generates grpc-gateway library files ([example](/example/github.com/grpc-ecosystem/grpc-gateway/gateway_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [grpc.js](/github.com/stackb/grpc.js) | [closure_grpc_compile](/github.com/stackb/grpc.js#closure_grpc_compile) | Generates protobuf closure grpc *.js files ([example](/example/github.com/stackb/grpc.js/closure_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [grpc.js](/github.com/stackb/grpc.js) | [closure_grpc_library](/github.com/stackb/grpc.js#closure_grpc_library) | Generates protobuf closure library *.js files ([example](/example/github.com/stackb/grpc.js/closure_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [grpc-web](/github.com/grpc/grpc-web) | [closure_grpc_compile](/github.com/grpc/grpc-web#closure_grpc_compile) | Generates a closure *.js protobuf+gRPC files ([example](/example/github.com/grpc/grpc-web/closure_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [grpc-web](/github.com/grpc/grpc-web) | [commonjs_grpc_compile](/github.com/grpc/grpc-web#commonjs_grpc_compile) | Generates a commonjs *.js protobuf+gRPC files ([example](/example/github.com/grpc/grpc-web/commonjs_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [grpc-web](/github.com/grpc/grpc-web) | [commonjs_dts_grpc_compile](/github.com/grpc/grpc-web#commonjs_dts_grpc_compile) | Generates a commonjs_dts *.js protobuf+gRPC files ([example](/example/github.com/grpc/grpc-web/commonjs_dts_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [grpc-web](/github.com/grpc/grpc-web) | [ts_grpc_compile](/github.com/grpc/grpc-web#ts_grpc_compile) | Generates a commonjs *.ts protobuf+gRPC files ([example](/example/github.com/grpc/grpc-web/ts_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [grpc-web](/github.com/grpc/grpc-web) | [closure_grpc_library](/github.com/grpc/grpc-web#closure_grpc_library) | Generates protobuf closure library *.js files ([example](/example/github.com/grpc/grpc-web/closure_grpc_library)) |

## Overview

These rules are the successor to <https://github.com/pubref/rules_protobuf> and
are in a pre-release status.  The primary goals are:

1. Interoperate with the native `proto_library` rules and other proto support in
   the bazel ecosystem as much as possible.
2. Provide a `proto_plugin` rule to support custom protoc plugins.
3. Minimal dependency loading.  Proto rules should not pull in more dependencies
   than they absolutely need.

> NOTE: in general, try to use the native proto library rules when possible to
minimize external dependencies in your project.  Add `rules_proto` when you have
more complex proto requirements such as when dealing with multiple output
languages, gRPC, unsupported (native) language support, or custom proto plugins.

## Installation

Add rules_proto your `WORKSPACE`:

```python
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "build_stack_rules_proto",
    urls = ["https://github.com/stackb/rules_proto/archive/{GIT_COMMIT_ID}.tar.gz"],
    sha256 = "{ARCHIVE_TAR_GZ_SHA256}",
    strip_prefix = "rules_proto-{GIT_COMMIT_ID}",
)
```

> Important: Follow instructions in the language-specific `README.md` for
additional workspace dependencies that may be required.

## Usage

**Step 1**: write a protocol buffer file (example: `thing.proto`):

```proto
syntax = "proto3";

package example;

import "google/protobuf/any.proto";

message Thing {
    string name = 1;
    google.protobuf.Any payload = 2;
}
```

**Step 2**: write a `BAZEL.build` file with a native `proto_library` rule:

```python
proto_library(
    name = "thing_proto",
    srcs = ["thing.proto"],
    deps = ["@com_google_protobuf//:any_proto"],
)
```

In this example we have a dependency on a well-known type `any.proto` hance the
`proto_library` to `proto_library` dependency.

**Step 3**: add a `cpp_proto_compile` rule (substitute `cpp_` for the language
of your choice).

> NOTE: In this example `thing.proto` does not include service definitions
(gRPC).  For protos with services, use the `cpp_grpc_compile` rule instead.

```python
# BUILD.bazel
load("@build_stack_rules_proto//cpp:cpp_proto_compile.bzl", "cpp_proto_compile")

cpp_proto_compile(
    name = "cpp_thing_proto",
    deps = [":thing_proto"],
)
```

But wait, before we can build this, we need to load the dependencies necessary
for this rule (from [cpp/README.md](/cpp/README.md)):

**Step 4**: load the workspace macro corresponding to the build rule.

```python
# WORKSPACE
load("@build_stack_rules_proto//cpp:deps.bzl", "cpp_proto_compile")

cpp_proto_compile()
```

Note that the workspace macro has the same name of the build rule.  Refer to the
[cpp/deps.bzl](/cpp/deps.bzl) for details on what other dependencies are loaded.

We're now ready to build the rule:

**Step 5**: build it.

```sh
$ bazel build //example/proto:cpp_thing_proto
Target //example/proto:cpp_thing_proto up-to-date:
  bazel-genfiles/example/proto/cpp_thing_proto/example/proto/thing.pb.h
  bazel-genfiles/example/proto/cpp_thing_proto/example/proto/thing.pb.cc
```

If we were only interested for the generated file artifacts, the
`cpp_grpc_compile` rule would be fine.  However, for convenience we'd rather
have the outputs compiled into an `*.so` file.  To do that, let's change the
rule from `cpp_proto_compile` to `cpp_proto_library`:

```python
# WORKSPACE
load("@build_stack_rules_proto//cpp:deps.bzl", "cpp_proto_library")

cpp_proto_library()
```

```python
# BUILD.bazel
load("@build_stack_rules_proto//cpp:cpp_proto_library.bzl", "cpp_proto_library")

cpp_proto_library(
    name = "cpp_thing_proto",
    deps = [":thing_proto"],
)
```

```sh
$ bazel build //example/proto:cpp_thing_proto
Target //example/proto:cpp_thing_proto up-to-date:
  bazel-bin/example/proto/libcpp_thing_proto.a
  bazel-bin/example/proto/libcpp_thing_proto.so  bazel-genfiles/example/proto/cpp_thing_proto/example/proto/thing.pb.h
  bazel-genfiles/example/proto/cpp_thing_proto/example/proto/thing.pb.cc
```

This way, we can use `//example/proto:cpp_thing_proto` as a dependency of any
other `cc_library` or `cc_binary` rule as per normal.

> NOTE: the `cpp_proto_library` implicitly calls `cpp_proto_compile`, and we can
access that rule by adding `_pb` at the end of the rule name, like `bazel build
//example/proto:cpp_thing_proto_pb`.

### Summary

* There are generally four rule flavors for any language `{lang}`:
`{lang}_proto_compile`, `{lang}_proto_library`, `{lang}_grpc_compile`, and
`{lang}_grpc_library`.

* If you are solely interested in the source code artifacts, use the
  `{lang}_{proto|grpc}_compile` rule.

* If your proto file has services, use the `{lang}_grpc_{compile|library}`
  rule instead.

* Load any external dependencies needed for the rule via the
  `load("@build_stack_rules_proto//{lang}:deps.bzl",
  "{lang}_{proto|grpc}_{compile|library}")`.

## Developers

### Code Layout

Each language `{lang}` has a top-level subdirectory that contains:

1. `{lang}/README.md`: generated documentation for the rule(s).

1. `{lang}/deps.bzl`: contains macro functions that declare repository rule
   dependencies for that language.  The name of the macro corresponds to the
   name of the build rule you'd like to use.

2. `{lang}/{rule}.bzl`: rule implementation(s) of the form
`{lang}_{kind}_{type}`, where `kind` is one of `proto|grpc` and `type` is one of
`compile|library`.

3. `{lang}/BUILD.bazel`: contains `proto_plugin()` declarations for the available
   plugins for that language.

4. `{lang}/example/{rule}/`: contains a generated `WORKSPACE` and `BUILD.bazel`
demonstrating standalone usage.

5. `{lang}/example/routeguide/`: routeguide example implementation, if possible.

The root directory contains the base rule defintions:

* `plugin.bzl`: A build rule that defines the name, tool binary, and options for
  a particular proto plugin.

* `compile.bzl`: A build rule that contains the `proto_compile` rule.  This rule
  calls `protoc` with a given list of plugins and generates output files.

Additional protoc plugins and their rules are scoped to the github repository
name where the plugin resides.  For example, there are 3 grpc-web
implementations in this repo:
`[github.com/improbable-eng/grpc-web](./github.com/improbable-eng/grpc-web),
[github.com/grpc/grpc-web](./github.com/grpc/grpc-web), and
[github.com/stackb/grpc.js](./github.com/stackb/grpc.js).

### Rule Generation

To help maintain consistency of the rule implementations and documentation, all
of the rule implementations are generated by the tool `//tools/rulegen`. Changes
in the main `README.md` should be placed in `tools/rulegen/README.header.md` or
`tools/rulegen/README.footer.md`.  Changes to generated rules should be put in
the source files (example: `tools/rulegen/java.go`).

### Transitivity

Briefly, here's how the rules work:

1. Using the `proto_library` graph, collect all the `*.proto` files directly and
transitively required for the protobuf compilation.

2. Copy the `*.proto` files into a "staging" area in the bazel sandbox such that a
single `-I.` will satisfy all imports.

3. Call `protoc OPTIONS FILELIST` and generate outputs.

The concept of *transitivity* (as defined here) affects which files in the set
`*.proto` files are named in the `FILELIST`.  If we only list direct
dependencies then we say `transitive = False`.  If all files are named, then
`transitive = True`.  The set of files that can be included or excluded from the
`FILELIST` are called *transitivity rules*, which can be defined on a per-rule
or per-plugin basis.  Please grep the codebase for examples of their usage.

## Developing Custom Plugins

Follow the pattern seen in the multiple examples in this repository.  The basic idea is:

1. Load the plugin rule: `load("@build_stack_rules_proto//:plugin.bzl", "proto_plugin")`.
2. Define the rule, giving it a `name`, `options` (not mandatory), `tool`, and
   `outputs`.
3. `tool` is a label that refers to the binary executable for the plugin itself.
4. `outputs` is a list of strings that predicts the pattern of files generated
   by the plugin.  Specifying outputs is the only attribute that requires much
   mental effort.

## Contributing

Contributions welcome; please create Issues or GitHub pull requests.

## ProtoInfo provider

A few notes about the 'proto' provider.  Bazel initially implemented
[providers](https://docs.bazel.build/versions/master/skylark/lib/Provider.html)
using simple strings.  For example, one could write a rule that depends on
`proto_library` rules as follows:

```python
my_custom_rule = rule(
    implementation = my_custom_rule_impl,
    attrs = {
        "deps": attr.label_list(
            doc = "proto_library dependencies",
            mandatory = True,
            providers = ["proto"],
        ),
    }
) 
```

We can then collect a list of the provided "info" objects as follows:

```python
def my_custom_rule_impl(ctx):
    # list<ProtoInfo> A list of ProtoInfo values
    deps = [ dep.proto for dep in ctx.attr.deps ]
```

This style of provider usage is now considered "legacy".  As of bazel 0.22.0
(Jan 2019)
[4817df](https://github.com/bazelbuild/bazel/commit/4817df1e862fe2de3b17905dbbc0c7badff6bb4b),
the
[ProtoInfo](https://docs.bazel.build/versions/master/skylark/lib/ProtoInfo.html)
provider was introduced.  Rules using this provider should be changed to:

```python
my_custom_rule = rule(
    implementation = my_custom_rule_impl,
    attrs = {
        "deps": attr.label_list(
            doc = "proto_library dependencies",
            mandatory = True,
            providers = [ProtoInfo],
        ),
 
```

```python
def my_custom_rule_impl(ctx):
    # <list<ProtoInfo>> A list of ProtoInfo
    deps = [ dep[ProtoInfo] for dep in ctx.attr.deps ]
```

A flag to toggle availability of the "legacy" provider was introduced in
[e3f4ce](https://github.com/bazelbuild/bazel/commit/e3f4ce739f0dbe0fea2b5580655fc96012d162cd)
as `--incompatible_disable_legacy_proto_provider`.  Therefore, if you still want
the legacy behavior specify
`--incompatible_disable_legacy_proto_provider=false`.

Additional reading:

- https://github.com/bazelbuild/bazel/issues/3701
- https://github.com/bazelbuild/bazel/issues/6901
