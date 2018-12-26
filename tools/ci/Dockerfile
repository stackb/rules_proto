FROM l.gcr.io/google/bazel:0.20.0

RUN apt-get update && apt-get install -y \
  xz-utils \
  bzip2 \
  ruby \
  make 

# ERROR: Analysis of target '//scala/example/routeguide:client' failed; build aborted: no such package '@org_json4s_json4s_jackson_2_12//jar': Traceback (most recent call last):
# 	File "/root/.cache/bazel/_bazel_root/adcbc75d8a72a02f3d3df29ee04882f7/external/bazel_tools/tools/build_defs/repo/maven_rules.bzl", line 280
# 		_maven_artifact_impl(ctx, "jar", _maven_jar_build_file_te...)
# 	File "/root/.cache/bazel/_bazel_root/adcbc75d8a72a02f3d3df29ee04882f7/external/bazel_tools/tools/build_defs/repo/maven_rules.bzl", line 230, in _maven_artifact_impl
# 		_check_dependencies(ctx)
# 	File "/root/.cache/bazel/_bazel_root/adcbc75d8a72a02f3d3df29ee04882f7/external/bazel_tools/tools/build_defs/repo/maven_rules.bzl", line 56, in _check_dependencies
# 		fail(("%s requires %s as a dependency...)))
# org_json4s_json4s_jackson_2_12 requires mvn as a dependency. Please check your PATH.
RUN apt-get install -y \
  maven

# https://github.com/bazelbuild/bazel/issues/6834
RUN ln -s /usr/bin/gold /usr/local/bin/gold

# xz required by:
# ERROR: /root/.cache/bazel/_bazel_root/adcbc75d8a72a02f3d3df29ee04882f7/external/io_bazel_rules_closure/third_party/expat/BUILD:7:1: Executing genrule @io_bazel_rules_closure//third_party/expat:libexpat_k8 failed (Exit 2)
# tar (child): xz: Cannot exec: No such file or directory

# ruby required for yugui rules_ruby

