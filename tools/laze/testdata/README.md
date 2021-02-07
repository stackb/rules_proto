# Gazelle test cases

This directory contains a suite of test cases for the laze language plugin for
Gazelle.

It's called "laze" because I'm lazy, and it sort of sounds like "gazelle", or
"glaze".

Please note that there are no `BUILD` or `BUILD.bazel` files in subdirs, insted
there are `BUILD.in` and `BUILD.out` describing what the `BUILD` should look
like initially and what the `BUILD` file should look like after the run. These
names are special because they are not recognized by Bazel as a proper `BUILD`
file, and therefore are included in the data dependency by the recursive data
glob in `//gazelle:go_default_test`. If you would like to include any extremely
complicated tests that contain proper `BUILD` files you will need to manually
add them to the `//gazelle:go_default_test` target's `data` section.

## `proto_plugin`

For every `proto_plugin` rule, generate a corresponding
`proto_plugin_info_provider_test`.

## `proto_rule`

For every `proto_rule` rule, generate a corresponding
`proto_rule_info_provider_test` and a `proto_rule_test`.
