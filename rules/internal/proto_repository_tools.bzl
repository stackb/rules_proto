# Copyright 2019 The Bazel Authors. All rights reserved.
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

load("//rules/internal:execution.bzl", "env_execute", "executable_extension")
load("@bazel_gazelle//internal:go_repository_cache.bzl", "read_cache_env")

_PROTO_REPOSITORY_TOOLS_BUILD_FILE = """
package(default_visibility = ["//visibility:public"])

filegroup(
    name = "gazelle",
    srcs = ["bin/gazelle{extension}"],
)
"""

def _proto_repository_tools_impl(ctx):
    # Create a link to the rules_proto repo. This will be our GOPATH.
    env = read_cache_env(ctx, str(ctx.path(ctx.attr.go_cache)))
    extension = executable_extension(ctx)
    go_tool = env["GOROOT"] + "/bin/go" + extension

    # symlink into GOROOT
    rules_proto_path = ctx.path(Label("@build_stack_rules_proto//:WORKSPACE"))
    ctx.symlink(
        rules_proto_path.dirname,
        "src/github.com/stackb/rules_proto",
    )

    env.update({
        "GOPATH": str(ctx.path(".")),
        # TODO(gravypod): make this more hermetic
        "GO111MODULE": "off",
        # workaround: avoid the Go SDK paths from leaking into the binary
        "GOROOT_FINAL": "GOROOT",
        # workaround: avoid cgo paths in /tmp leaking into binary
        "CGO_ENABLED": "0",
    })

    if "PATH" in ctx.os.environ:
        # workaround: to find gcc for go link tool on Arm platform
        env["PATH"] = ctx.os.environ["PATH"]
    if "GOPROXY" in ctx.os.environ:
        env["GOPROXY"] = ctx.os.environ["GOPROXY"]

    # Build the tools.
    args = [
        go_tool,
        "install",
        "-ldflags",
        "-w -s",
        "-gcflags",
        "all=-trimpath=" + env["GOPATH"],
        "-asmflags",
        "all=-trimpath=" + env["GOPATH"],
        "github.com/stackb/rules_proto/cmd/gazelle",
    ]
    result = env_execute(ctx, args, environment = env)
    if result.return_code:
        fail("failed to build tools: " + result.stderr)

    # add a build file to export the tools
    ctx.file(
        "BUILD.bazel",
        _PROTO_REPOSITORY_TOOLS_BUILD_FILE.format(extension = executable_extension(ctx)),
        False,
    )

proto_repository_tools = repository_rule(
    _proto_repository_tools_impl,
    attrs = {
        "go_cache": attr.label(
            mandatory = True,
            allow_single_file = True,
        ),
    },
    environ = [
        "GOCACHE",
        "GOPATH",
        "GO_REPOSITORY_USE_HOST_CACHE",  # TODO(gravypod): Find out why this is here in rules_go
    ],
)
"""proto_repository_tools is a synthetic repository used by proto_repository.
"""
