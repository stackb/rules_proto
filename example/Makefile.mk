android_android_proto_compile_example:
	cd example/android/android_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

android_android_grpc_compile_example:
	cd example/android/android_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

android_android_proto_library_example:
	cd example/android/android_proto_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

android_android_grpc_library_example:
	cd example/android/android_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

android_examples: android_android_proto_compile_example android_android_grpc_compile_example android_android_proto_library_example android_android_grpc_library_example

closure_closure_proto_compile_example:
	cd example/closure/closure_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

closure_closure_proto_library_example:
	cd example/closure/closure_proto_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

closure_examples: closure_closure_proto_compile_example closure_closure_proto_library_example

cpp_cpp_proto_compile_example:
	cd example/cpp/cpp_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

cpp_cpp_grpc_compile_example:
	cd example/cpp/cpp_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

cpp_cpp_proto_library_example:
	cd example/cpp/cpp_proto_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

cpp_cpp_grpc_library_example:
	cd example/cpp/cpp_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

cpp_examples: cpp_cpp_proto_compile_example cpp_cpp_grpc_compile_example cpp_cpp_proto_library_example cpp_cpp_grpc_library_example

csharp_csharp_proto_compile_example:
	cd example/csharp/csharp_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

csharp_csharp_grpc_compile_example:
	cd example/csharp/csharp_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

csharp_csharp_proto_library_example:
	cd example/csharp/csharp_proto_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

csharp_csharp_grpc_library_example:
	cd example/csharp/csharp_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

csharp_examples: csharp_csharp_proto_compile_example csharp_csharp_grpc_compile_example csharp_csharp_proto_library_example csharp_csharp_grpc_library_example

d_d_proto_compile_example:
	cd example/d/d_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

d_examples: d_d_proto_compile_example

dart_dart_proto_compile_example:
	cd example/dart/dart_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

dart_dart_grpc_compile_example:
	cd example/dart/dart_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

dart_dart_proto_library_example:
	cd example/dart/dart_proto_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

dart_dart_grpc_library_example:
	cd example/dart/dart_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

dart_examples: dart_dart_proto_compile_example dart_dart_grpc_compile_example dart_dart_proto_library_example dart_dart_grpc_library_example

go_go_proto_compile_example:
	cd example/go/go_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

go_go_grpc_compile_example:
	cd example/go/go_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

go_go_proto_library_example:
	cd example/go/go_proto_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

go_go_grpc_library_example:
	cd example/go/go_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

go_examples: go_go_proto_compile_example go_go_grpc_compile_example go_go_proto_library_example go_go_grpc_library_example

java_java_proto_compile_example:
	cd example/java/java_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

java_java_grpc_compile_example:
	cd example/java/java_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

java_java_proto_library_example:
	cd example/java/java_proto_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

java_java_grpc_library_example:
	cd example/java/java_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

java_examples: java_java_proto_compile_example java_java_grpc_compile_example java_java_proto_library_example java_java_grpc_library_example

node_node_proto_compile_example:
	cd example/node/node_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

node_node_grpc_compile_example:
	cd example/node/node_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

node_node_proto_library_example:
	cd example/node/node_proto_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

node_node_grpc_library_example:
	cd example/node/node_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

node_examples: node_node_proto_compile_example node_node_grpc_compile_example node_node_proto_library_example node_node_grpc_library_example

objc_objc_proto_compile_example:
	cd example/objc/objc_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

objc_objc_grpc_compile_example:
	cd example/objc/objc_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

objc_examples: objc_objc_proto_compile_example objc_objc_grpc_compile_example

php_php_proto_compile_example:
	cd example/php/php_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

php_php_grpc_compile_example:
	cd example/php/php_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

php_examples: php_php_proto_compile_example php_php_grpc_compile_example

python_python_proto_compile_example:
	cd example/python/python_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

python_python_grpc_compile_example:
	cd example/python/python_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

python_python_proto_aspect_compile_example:
	cd example/python/python_proto_aspect_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

python_python_grpc_aspect_compile_example:
	cd example/python/python_grpc_aspect_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

python_python_proto_library_example:
	cd example/python/python_proto_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

python_python_grpc_library_example:
	cd example/python/python_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

python_examples: python_python_proto_compile_example python_python_grpc_compile_example python_python_proto_aspect_compile_example python_python_grpc_aspect_compile_example python_python_proto_library_example python_python_grpc_library_example

ruby_ruby_proto_compile_example:
	cd example/ruby/ruby_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

ruby_ruby_grpc_compile_example:
	cd example/ruby/ruby_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

ruby_ruby_proto_library_example:
	cd example/ruby/ruby_proto_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

ruby_ruby_grpc_library_example:
	cd example/ruby/ruby_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

ruby_examples: ruby_ruby_proto_compile_example ruby_ruby_grpc_compile_example ruby_ruby_proto_library_example ruby_ruby_grpc_library_example

rust_rust_proto_compile_example:
	cd example/rust/rust_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

rust_rust_grpc_compile_example:
	cd example/rust/rust_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

rust_rust_proto_library_example:
	cd example/rust/rust_proto_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

rust_rust_grpc_library_example:
	cd example/rust/rust_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

rust_examples: rust_rust_proto_compile_example rust_rust_grpc_compile_example rust_rust_proto_library_example rust_rust_grpc_library_example

scala_scala_proto_compile_example:
	cd example/scala/scala_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

scala_scala_grpc_compile_example:
	cd example/scala/scala_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

scala_scala_proto_library_example:
	cd example/scala/scala_proto_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

scala_scala_grpc_library_example:
	cd example/scala/scala_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

scala_examples: scala_scala_proto_compile_example scala_scala_grpc_compile_example scala_scala_proto_library_example scala_scala_grpc_library_example

swift_swift_proto_compile_example:
	cd example/swift/swift_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

swift_swift_grpc_compile_example:
	cd example/swift/swift_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

swift_swift_proto_library_example:
	cd example/swift/swift_proto_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

swift_swift_grpc_library_example:
	cd example/swift/swift_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

swift_examples: swift_swift_proto_compile_example swift_swift_grpc_compile_example swift_swift_proto_library_example swift_swift_grpc_library_example

gogo_gogo_proto_compile_example:
	cd example/github.com/gogo/protobuf/gogo_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

gogo_gogo_grpc_compile_example:
	cd example/github.com/gogo/protobuf/gogo_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

gogo_gogo_proto_library_example:
	cd example/github.com/gogo/protobuf/gogo_proto_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

gogo_gogo_grpc_library_example:
	cd example/github.com/gogo/protobuf/gogo_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

gogo_gogofast_proto_compile_example:
	cd example/github.com/gogo/protobuf/gogofast_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

gogo_gogofast_grpc_compile_example:
	cd example/github.com/gogo/protobuf/gogofast_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

gogo_gogofast_proto_library_example:
	cd example/github.com/gogo/protobuf/gogofast_proto_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

gogo_gogofast_grpc_library_example:
	cd example/github.com/gogo/protobuf/gogofast_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

gogo_gogofaster_proto_compile_example:
	cd example/github.com/gogo/protobuf/gogofaster_proto_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

gogo_gogofaster_grpc_compile_example:
	cd example/github.com/gogo/protobuf/gogofaster_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

gogo_gogofaster_proto_library_example:
	cd example/github.com/gogo/protobuf/gogofaster_proto_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

gogo_gogofaster_grpc_library_example:
	cd example/github.com/gogo/protobuf/gogofaster_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

gogo_examples: gogo_gogo_proto_compile_example gogo_gogo_grpc_compile_example gogo_gogo_proto_library_example gogo_gogo_grpc_library_example gogo_gogofast_proto_compile_example gogo_gogofast_grpc_compile_example gogo_gogofast_proto_library_example gogo_gogofast_grpc_library_example gogo_gogofaster_proto_compile_example gogo_gogofaster_grpc_compile_example gogo_gogofaster_proto_library_example gogo_gogofaster_grpc_library_example

grpc-gateway_gateway_grpc_compile_example:
	cd example/github.com/grpc-ecosystem/grpc-gateway/gateway_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

grpc-gateway_gateway_swagger_compile_example:
	cd example/github.com/grpc-ecosystem/grpc-gateway/gateway_swagger_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

grpc-gateway_gateway_grpc_library_example:
	cd example/github.com/grpc-ecosystem/grpc-gateway/gateway_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

grpc-gateway_examples: grpc-gateway_gateway_grpc_compile_example grpc-gateway_gateway_swagger_compile_example grpc-gateway_gateway_grpc_library_example

grpc.js_closure_grpc_compile_example:
	cd example/github.com/stackb/grpc.js/closure_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

grpc.js_closure_grpc_library_example:
	cd example/github.com/stackb/grpc.js/closure_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

grpc.js_examples: grpc.js_closure_grpc_compile_example grpc.js_closure_grpc_library_example

grpc-web_closure_grpc_compile_example:
	cd example/github.com/grpc/grpc-web/closure_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

grpc-web_commonjs_grpc_compile_example:
	cd example/github.com/grpc/grpc-web/commonjs_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

grpc-web_commonjs_dts_grpc_compile_example:
	cd example/github.com/grpc/grpc-web/commonjs_dts_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

grpc-web_ts_grpc_compile_example:
	cd example/github.com/grpc/grpc-web/ts_grpc_compile; \
	bazel build --disk_cache=../../bazel-disk-cache //...

grpc-web_closure_grpc_library_example:
	cd example/github.com/grpc/grpc-web/closure_grpc_library; \
	bazel build --disk_cache=../../bazel-disk-cache //...

grpc-web_examples: grpc-web_closure_grpc_compile_example grpc-web_commonjs_grpc_compile_example grpc-web_commonjs_dts_grpc_compile_example grpc-web_ts_grpc_compile_example grpc-web_closure_grpc_library_example

all_examples: grpc-web_closure_grpc_compile_example grpc-web_commonjs_grpc_compile_example grpc-web_commonjs_dts_grpc_compile_example grpc-web_ts_grpc_compile_example grpc-web_closure_grpc_library_example
