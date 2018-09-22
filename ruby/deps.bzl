load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def ruby_proto_library_deps(
  rules_ruby_version = "5976385c9c4b94647bc95e8bf9d9989f1dee4ee3", # PR#8,
  rules_ruby_sha256 = "7991ded3b902aba4c13fa7bdd67132edfcc279930b356737c1a3d3b2686d08c8", 
):
    existing = native.existing_rules()

    if "com_github_yugui_rules_ruby" not in existing:
        http_archive(
            name = "com_github_yugui_rules_ruby",
            urls = ["https://github.com/yugui/rules_ruby/archive/%s.tar.gz" % rules_ruby_version],
            sha256 = rules_ruby_sha256,
            strip_prefix = "rules_ruby-" + rules_ruby_version,
        )
        # native.local_repository(
        #     name = "com_github_yugui_rules_ruby",
        #     path = "/home/pcj/github/yugui/rules_ruby",
        # )
