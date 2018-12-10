# NOTE: this is a slightly modified version of
# https://github.com/dart-lang/rules_dart/blob/master/dart/build_rules/internal/sdk.bzl
# to experiment with different versions.  A better solution would be a PR on the
# rules_dart repository to make it customizable.

# Copyright 2016 The Bazel Authors. All rights reserved.
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
"""A repository rule to download the Dart SDK."""

SDK_BUILD_FILE = """
package(default_visibility = [ "//visibility:public" ])

filegroup(
  name = "dart_vm",
  srcs = ["bin/dart"],
)

filegroup(
  name = "analyzer",
  srcs = ["bin/dartanalyzer"],
)

filegroup(
  name = "dart2js",
  srcs = ["bin/dart2js"],
)

filegroup(
  name = "dart2js_support",
  srcs = glob([
      "bin/dart",
      "bin/snapshots/dart2js.dart.snapshot",
      "lib/**",
  ]),
)

filegroup(
  name = "dev_compiler",
  srcs = ["bin/dartdevc"],
)

filegroup(
    name = "ddc_support",
    srcs = glob(["lib/dev_compiler/legacy/*.*"]),
)

filegroup(
  name = "pub",
  srcs = ["bin/pub"],
)

filegroup(
  name = "pub_support",
  srcs = glob([
      "version",
      "bin/dart",
      "bin/snapshots/pub.dart.snapshot",
  ]),
)

filegroup(
  name = "sdk_summary",
  srcs = ["lib/_internal/strong.sum"],
)

filegroup(
  name = "sdk_summary_dill",
  srcs = ["lib/_internal/vm_platform_strong.dill"],
)

filegroup(
  name = "kernel_worker_snapshot",
  srcs = ["bin/snapshots/kernel_worker.dart.snapshot"],
)

filegroup(
    name = "lib_files_no_summaries",
    srcs = glob([
        "lib/**/*.dart",
        "lib/dart_client.platform",
        "lib/dart_shared.platform",
        "version",
    ]),
)

filegroup(
    name = "lib_files",
    srcs = glob([
        ":lib_files_no_summaries",
        ":sdk_summaries",
    ]),
)

"""

_hosted_prefix = "https://storage.googleapis.com/dart-archive/channels/dev/release"
# _hosted_prefix = "https://storage.googleapis.com/dart-archive/channels/stable/release"
_linux_file = "dartsdk-linux-x64-release.zip"
_mac_file = "dartsdk-macos-x64-release.zip"

_version = "2.2.0-dev.0.0"
_linux_sha = "fe2087156d5215cf5f333cf6cd72d9a461fb1e475bcc60e48e53d348306ca0d9"
_mac_sha = "3e6953d30c1c1427b689da9f7effb86380dde1475d0291634c3868651926a6b9"

def _sdk_repository_impl(repository_ctx):
    """Downloads the appropriate SDK for the current OS."""
    os_name = repository_ctx.os.name

    file_name = False

    if "linux" in os_name:
        file_name = _linux_file
        sha = _linux_sha
    elif "mac os" in os_name:
        file_name = _mac_file
        sha = _mac_sha

    if not file_name:
        fail("Cannot find SDK for OS: %s" % os_name)

    repository_ctx.download_and_extract(
        url = "%s/%s/sdk/%s" % (_hosted_prefix, _version, file_name),
        sha256 = sha,
        stripPrefix = "dart-sdk/",
    )
    repository_ctx.file("BUILD", SDK_BUILD_FILE)

dart_sdk_repository = repository_rule(
    implementation = _sdk_repository_impl,
)
