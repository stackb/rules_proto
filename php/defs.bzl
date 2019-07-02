# Aggregate all `php` rules to one loadable file
load(":php_proto_compile.bzl", _php_proto_compile="php_proto_compile")
load(":php_grpc_compile.bzl", _php_grpc_compile="php_grpc_compile")

php_proto_compile = _php_proto_compile
php_grpc_compile = _php_grpc_compile
