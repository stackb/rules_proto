"""proto_repostitory.bzl provides the proto_repository rule."""

# Copyright 2014 The Bazel Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@bazel_features//:features.bzl", "bazel_features")
load("@build_stack_rules_proto//rules/proto:proto_repository.bzl", "proto_repository_attrs", proto_repository_repo_rule = "proto_repository")

def _extension_metadata(
        module_ctx,
        *,
        root_module_direct_deps = None,
        root_module_direct_dev_deps = None,
        reproducible = False):
    """returns the module_ctx.extension_metadata in a bazel-version-aware way

    This function was copied from the bazel-gazelle repository.
    """

    if not hasattr(module_ctx, "extension_metadata"):
        return None
    metadata_kwargs = {}
    if bazel_features.external_deps.extension_metadata_has_reproducible:
        metadata_kwargs["reproducible"] = reproducible
    return module_ctx.extension_metadata(
        root_module_direct_deps = root_module_direct_deps,
        root_module_direct_dev_deps = root_module_direct_dev_deps,
        **metadata_kwargs
    )

def _proto_repository_impl(module_ctx):
    # named_repos is a dict<K,V> where V is the kwargs for the actual
    # "proto_repository" repo rule and K is the tag.name (the name given by the
    # MODULE.bazel author)
    named_archives = {}

    # iterate all the module tags and gather a list of named_archives.
    #
    # TODO(pcj): what is the best practice for version selection here? Do I need
    # to check if module.is_root and handle that differently?
    #
    for module in module_ctx.modules:
        for tag in module.tags.archive:
            kwargs = {
                attr: getattr(tag, attr)
                for attr in _proto_repository_archive_attrs.keys()
                if hasattr(tag, attr)
            }
            named_archives[tag.name] = kwargs

    # declare a repository rule foreach one
    for _, kwargs in named_archives.items():
        proto_repository_repo_rule(**kwargs)

    return _extension_metadata(
        module_ctx,
        reproducible = True,
    )

_proto_repository_archive_attrs = proto_repository_attrs | {
    "name": attr.string(
        doc = "The repo name.",
        mandatory = True,
    ),
}

proto_repository = module_extension(
    implementation = _proto_repository_impl,
    tag_classes = dict(
        archive = tag_class(
            doc = "declare an http_archive repository that is post-processed by a custom version of gazelle that includes the 'protobuf' language",
            attrs = _proto_repository_archive_attrs,
        ),
    ),
)
