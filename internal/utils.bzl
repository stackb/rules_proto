load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def github_archive(name, org, repo, ref, sha256):
    """Declare an http_archive from github
    """

    # correct github quirk about removing the 'v' in front of tags
    stripRef = ref
    if stripRef.startswith("v"):
        stripRef = ref[1:]

    if name not in native.existing_rules():
        http_archive(
            name = name,
            strip_prefix = repo + "-" + stripRef,
            urls = [
                "https://mirror.bazel.build/github.com/%s/%s/archive/%s.tar.gz" % (org, repo, ref),
                "https://github.com/%s/%s/archive/%s.tar.gz" % (org, repo, ref),
            ],
            sha256 = sha256,
        )

def get_ref(name, default, kwargs):
    key = name + "_ref"
    return kwargs.get(key, default)

def get_sha256(name, default, kwargs):
    key = name + "_sha256"
    return kwargs.get(key, default)
