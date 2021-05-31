---
layout: default
title: php
permalink: examples/php
parent: Examples
---


# php example

`bazel test //example/golden:php_test`


## `BUILD.bazel` (after gazelle)

~~~python
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")

# "proto_rule" instantiates the proto_compile rule
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile

# "proto_plugin" instantiates the builtin php plugin
# gazelle:proto_plugin php implementation builtin:php

# "proto_language" binds the rule(s) and plugin(s) together
# gazelle:proto_language php rule proto_compile
# gazelle:proto_language php plugin php

proto_library(
    name = "example_proto",
    srcs = ["example.proto"],
    visibility = ["//visibility:public"],
)

proto_compile(
    name = "example_php_compile",
    outputs = [
        "E.php",
        "M.php",
        # "S.php",
    ],
    plugins = ["@build_stack_rules_proto//plugin/builtin:php"],
    proto = "example_proto",
    verbose = True,
)
~~~


## `BUILD.bazel` (before gazelle)

~~~python
# "proto_rule" instantiates the proto_compile rule
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile

# "proto_plugin" instantiates the builtin php plugin
# gazelle:proto_plugin php implementation builtin:php

# "proto_language" binds the rule(s) and plugin(s) together
# gazelle:proto_language php rule proto_compile
# gazelle:proto_language php plugin php
~~~


## `WORKSPACE`

~~~python
~~~

