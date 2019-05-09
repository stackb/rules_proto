# Copyright (C) 2017 The Dagger Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Origally taken from
# https://github.com/google/bazel-common/tree/master/tools/javadoc and modified to
# operate on srcjars.
#
# Example usage:
#
# javadoc_library(
#     name = "javadoc",
#     srcjars = ["//proto/api:java_pb"],
#     deps = ["//proto/api:java"],
#     doctitle = "My gRPC Java API",
#     root_packages = [
#         "com.example.api.v1",
#     ],
# )

def _check_non_empty(value, name):
    if not value:
        fail("%s must be non-empty" % name)

def _android_jar(android_api_level):
    if android_api_level == -1:
        return None
    return Label("@androidsdk//:platforms/android-%s/android.jar" % android_api_level)

def _javadoc_library(ctx):
    _check_non_empty(ctx.attr.root_packages, "root_packages")

    transitive_deps = [dep.java.transitive_deps for dep in ctx.attr.deps]
    if ctx.attr._android_jar:
        transitive_deps.append(ctx.attr._android_jar.files)

    classpath = depset([], transitive = transitive_deps).to_list()

    java_home = str(ctx.attr._jdk[java_common.JavaRuntimeInfo].java_home)

    javadoc_command = [
        java_home + "/bin/javadoc",
        '-sourcepath .',
        " ".join(ctx.attr.root_packages),
        "-use",
        "-subpackages", ":".join(ctx.attr.root_packages),
        "-encoding UTF8",
        "-classpath", ":".join([jar.path for jar in classpath]),
        "-notimestamp",
        "-d tmp",
        "-Xdoclint:-missing",
        "-quiet",
    ]

    if ctx.attr.doctitle:
        javadoc_command.append('-doctitle "%s"' % ctx.attr.doctitle)

    if ctx.attr.exclude_packages:
        javadoc_command.append("-exclude %s" % ":".join(ctx.attr.exclude_packages))

    for link in ctx.attr.external_javadoc_links:
        javadoc_command.append("-linkoffline {0} {0}".format(link))

    if ctx.attr.bottom_text:
        javadoc_command.append("-bottom '%s'" % ctx.attr.bottom_text)

    jar_command = "%s/bin/jar cf %s -C tmp ." % (java_home, ctx.outputs.jar.path)

    unjar_command = "%s/bin/jar xf %s" % (java_home, " ".join([f.path for f in ctx.files.srcjars]))

    srcs = depset(transitive = [src.files for src in ctx.attr.srcjars]).to_list()
    ctx.actions.run_shell(
        inputs = srcs + classpath + ctx.files._jdk,
        command = "%s && %s && %s" % (unjar_command, " ".join(javadoc_command), jar_command),
        outputs = [ctx.outputs.jar],
    )

javadoc_library = rule(
    attrs = {
        "srcjars": attr.label_list(allow_files = [".srcjar"]),
        "deps": attr.label_list(),
        "doctitle": attr.string(default = ""),
        "root_packages": attr.string_list(),
        "exclude_packages": attr.string_list(),
        "android_api_level": attr.int(default = -1),
        "bottom_text": attr.string(default = ""),
        "external_javadoc_links": attr.string_list(),
        "_android_jar": attr.label(
            default = _android_jar,
            allow_single_file = True,
        ),
        "_jdk": attr.label(
            default = Label("@bazel_tools//tools/jdk:current_java_runtime"),
            providers = [java_common.JavaRuntimeInfo],
        ),
    },
    outputs = {"jar": "%{name}.jar"},
    implementation = _javadoc_library,
)

"""
Generates a Javadoc jar path/to/target/<name>.jar.
Arguments:
  srcjars: sourcejar files to process
  deps: targets that contain references to other types referenced in Javadoc. This can be the
      java_library/android_library target(s) for the same sources
  root_packages: Java packages to include in generated Javadoc. Any subpackages not listed in
      exclude_packages will be included as well
  exclude_packages: Java packages to exclude from generated Javadoc
  android_api_level: If Android APIs are used, the API level to compile against to generate
      Javadoc
  doctitle: title for Javadoc's index.html. See javadoc -doctitle
  bottom_text: text passed to javadoc's `-bottom` flag
  external_javadoc_links: a list of URLs that are passed to Javadoc's `-linkoffline` flag
"""
