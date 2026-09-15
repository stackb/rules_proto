"""proto_repository.bzl provides the proto_repository rule."""

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

load("@build_stack_rules_proto//rules/proto:proto_repository.bzl", "protobuf_go_repository", _proto_repository_attrs = "proto_repository_attrs")

# TODO: narrow the set of available attrs
starlark_repository_attrs = _proto_repository_attrs

def starlark_repository(**kwargs):
    """starlark_repository wraps proto_repository and sets the language to starlark_bundle

    Args:
      **kwargs: the kwargs dict passed to protobuf_go_repository
    """
    name = kwargs.get("name")
    kwargs.setdefault("apparent_name", name)

    protobuf_go_repository(**kwargs)
