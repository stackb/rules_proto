load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

def node_proto_compile_deps():
    pass

def node_proto_library_deps(
  rules_node_version = "1c60708c599e6ebd5213f0987207a1d854f13e23",
  rules_node_sha256 = "248efb149bfa86d9d778b43949351015b23a8339405a9878467a1583ff6df348",
):
    node_proto_compile_deps()

    existing = native.existing_rules()

    if "org_pubref_rules_node" not in existing:
        http_archive(
            name = "org_pubref_rules_node",
            urls = ["https://github.com/pubref/rules_node/archive/%s.tar.gz" % rules_node_version],
            sha256 = rules_node_sha256,
            strip_prefix = "rules_node-" + rules_node_version,
        )
