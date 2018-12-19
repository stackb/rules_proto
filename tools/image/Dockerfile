FROM l.gcr.io/google/bazel:0.20.0

RUN apt-get update && apt-get install -y \
  xz-utils \
  bzip2 \
  ruby 

# xz required by:
# ERROR: /root/.cache/bazel/_bazel_root/adcbc75d8a72a02f3d3df29ee04882f7/external/io_bazel_rules_closure/third_party/expat/BUILD:7:1: Executing genrule @io_bazel_rules_closure//third_party/expat:libexpat_k8 failed (Exit 2)
# tar (child): xz: Cannot exec: No such file or directory

# ruby required for yugui rules_ruby

