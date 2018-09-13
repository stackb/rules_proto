
#load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

SIX_BUILD = """
genrule(
  name = "copy_six",
  srcs = ["six-1.10.0/six.py"],
  outs = ["six.py"],
  cmd = "cp $< $(@)",
)

py_library(
  name = "six",
  srcs = ["six.py"],
  srcs_version = "PY2AND3",
  visibility = ["//visibility:public"],
)
"""

def py_proto_deps():
    if "six" not in native.existing_rules():
        native.new_http_archive(
            name = "six",
            build_file_content = SIX_BUILD,
            sha256 = "105f8d68616f8248e24bf0e9372ef04d3cc10104f1980f54d57b2ce73a5ad56a",
            urls = ["https://pypi.python.org/packages/source/s/six/six-1.10.0.tar.gz#md5=34eed507548117b2ab523ab14b2f8b55"],
        )