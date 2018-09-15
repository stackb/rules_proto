
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

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
    existing = native.existing_rules()  
    
    if "six" not in existing:
        native.new_http_archive(
            name = "six",
            build_file_content = SIX_BUILD,
            sha256 = "105f8d68616f8248e24bf0e9372ef04d3cc10104f1980f54d57b2ce73a5ad56a",
            urls = ["https://pypi.python.org/packages/source/s/six/six-1.10.0.tar.gz#md5=34eed507548117b2ab523ab14b2f8b55"],
        )
    if "io_bazel_rules_python" not in existing:
        http_archive(
            name = "io_bazel_rules_python",
            urls = ["https://github.com/bazelbuild/rules_python/archive/8b5d0683a7d878b28fffe464779c8a53659fc645.tar.gz"],
            strip_prefix = "rules_python-8b5d0683a7d878b28fffe464779c8a53659fc645",
            sha256 = "8b32d2dbb0b0dca02e0410da81499eef8ff051dad167d6931a92579e3b2a1d48",
        )
