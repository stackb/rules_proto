load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# https://raw.githubusercontent.com/grpc/grpc/master/third_party/zlib.BUILD
ZLIB_BUILD = """
package(default_visibility = ["//visibility:public"])

licenses(["notice"])  # BSD/MIT-like license (for zlib)

_ZLIB_HEADERS = [
    "crc32.h",
    "deflate.h",
    "gzguts.h",
    "inffast.h",
    "inffixed.h",
    "inflate.h",
    "inftrees.h",
    "trees.h",
    "zconf.h",
    "zlib.h",
    "zutil.h",
]

_ZLIB_PREFIXED_HEADERS = ["zlib/include/" + hdr for hdr in _ZLIB_HEADERS]

# In order to limit the damage from the `includes` propagation
# via `:zlib`, copy the public headers to a subdirectory and
# expose those.
genrule(
    name = "copy_public_headers",
    srcs = _ZLIB_HEADERS,
    outs = _ZLIB_PREFIXED_HEADERS,
    cmd = "cp $(SRCS) $(@D)/zlib/include/",
    visibility = ["//visibility:private"],
)

cc_library(
    name = "zlib",
    srcs = [
        "adler32.c",
        "compress.c",
        "crc32.c",
        "deflate.c",
        "gzclose.c",
        "gzlib.c",
        "gzread.c",
        "gzwrite.c",
        "infback.c",
        "inffast.c",
        "inflate.c",
        "inftrees.c",
        "trees.c",
        "uncompr.c",
        "zutil.c",
        # Include the un-prefixed headers in srcs to work
        # around the fact that zlib isn't consistent in its
        # choice of <> or "" delimiter when including itself.
    ] + _ZLIB_HEADERS,
    hdrs = _ZLIB_PREFIXED_HEADERS,
    copts = select({
        "@bazel_tools//src/conditions:windows": [],
        "//conditions:default": [
            "-Wno-unused-variable",
            "-Wno-implicit-function-declaration",
        ],
    }),
    includes = ["zlib/include/"],
)
"""

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

def get_artifact(name, default, kwargs):
    key = name + "_artifact"
    return kwargs.get(key, default)

def get_sha256(name, default, kwargs):
    key = name + "_sha256"
    return kwargs.get(key, default)

def get_sha1(name, default, kwargs):
    key = name + "_sha1"
    return kwargs.get(key, default)

def external_protobuf_clib(**kwargs):
    name = "protobuf_clib"
    if name not in native.existing_rules():
        native.bind(
            name = name,
            actual = "@com_google_protobuf//:protoc_lib",
        )

def external_protobuf_headers(**kwargs):
    name = "protobuf_headers"
    if name not in native.existing_rules():
        native.bind(
            name = name,
            actual = "@com_google_protobuf//:protobuf_headers",
        )

def external_protocol_compiler(**kwargs):
    name = "protocol_compiler"
    if name not in native.existing_rules():
        native.bind(
            name = name,
            actual = "@com_google_protobuf//:protoc",
        )

def external_nanopb(**kwargs):
    com_github_nanopb_nanopb(**kwargs)
    name = "nanopb"
    if name not in native.existing_rules():
        native.bind(
            name = name,
            actual = "@com_github_nanopb_nanopb//:nanopb",
        )

def external_libssl(**kwargs):
    boringssl(**kwargs)
    name = "libssl"
    if name not in native.existing_rules():
        native.bind(
            name = name,
            actual = "@boringssl//:ssl",
        )

def zlib(**kwargs):
    if "zlib" not in native.existing_rules():
        http_archive(
            name = "zlib",
            build_file_content = ZLIB_BUILD,
            strip_prefix = "zlib-1.2.11",
            url = "https://github.com/madler/zlib/archive/1.2.11.tar.gz",
            sha256 = "c3e5e9fdd5004dcb542feda5ee4f0ff0744628baf8ed2dd5d66f8ca1197cb1a1",
        )

def external_zlib(**kwargs):
    zlib(**kwargs)
    name = "zlib"
    if name not in native.existing_rules():
        native.bind(
            name = name,
            actual = "@zlib//:z",
        )

# grpc also wants it named "//external:madler_zlib"
def external_madler_zlib(**kwargs):
    zlib(**kwargs)
    name = "madler_zlib"
    if name not in native.existing_rules():
        native.bind(
            name = name,
            actual = "@zlib//:z",
        )

def com_github_bazelbuild_bazel_gazelle(**kwargs):
    if "com_github_bazelbuild_bazel_gazelle" not in native.existing_rules():
        sha1 = "a79ae21dcb2e1f4d36c2b99bb14e27816c5f4100"
        http_archive(
            name = "com_github_bazelbuild_bazel_gazelle",
            strip_prefix = "bazel-gazelle-" + sha1,
            url = "https://github.com/bazelbuild/bazel-gazelle/archive/%s.tar.gz" % sha1,
            patch_cmds = [
                # Expose the go_library targets so we can use it!
                "sed -i 's|//:__subpackages__|//visibility:public|g' internal/rule/BUILD.bazel",
            ],
        )

def com_github_bazelbuild_buildtools(**kwargs):
    name = "com_github_bazelbuild_buildtools"
    ref = get_ref(name, "6415663945d3248207da955aafa1fa2af1a0f2ed", kwargs)
    sha256 = get_sha256(name, "d1e28237d1f4c2255c504246b4f3fd36f74d590f2974491b4399a84c58b495ed", kwargs)
    github_archive(name, "bazelbuild", "buildtools", ref, sha256)

def boringssl(**kwargs):
    if "boringssl" not in native.existing_rules():
        http_archive(
            name = "boringssl",
            # on the chromium-stable-with-bazel branch
            url = "https://boringssl.googlesource.com/boringssl/+archive/dcd3e6e6ecddf059adb48fca45bc7346a108bdd9.tar.gz",
        )

def com_github_nanopb_nanopb(**kwargs):
    name = "com_github_nanopb_nanopb"
    ref = get_ref(name, "ae9901f2a31500e8fdc93fa9804d24851c58bb1e", kwargs)

    # should be 8bbbb1e78d4ddb0a1919276924ab10d11b631df48b657d960e0c795a25515735?
    sha256 = get_sha256(name, "7aa0ab179eff56241b6cded9cd07324af2395ad56d5478e2f7dabdb42b65d3fb", kwargs)
    github_archive(name, "nanopb", "nanopb", ref, sha256)

def com_google_protobuf(**kwargs):
    name = "com_google_protobuf"
    ref = get_ref(name, "v3.7.1", kwargs)
    sha256 = get_sha256(name, "f1748989842b46fa208b2a6e4e2785133cfcc3e4d43c17fecb023733f0f5443f", kwargs)

    github_archive(name, "google", "protobuf", ref, sha256)

def com_github_grpc_grpc(**kwargs):
    name = "com_github_grpc_grpc"
    ref = get_ref(name, "acd79569bff44550c87e09768e2184e91c7eb610", kwargs)
    sha256 = get_sha256(name, "f5bb5ee99fb877ecb915773ac7f5a6e62cee9ce89731cccc779b0b5eb2cacfbf", kwargs)
    github_archive(name, "grpc", "grpc", ref, sha256)

def io_bazel_rules_dotnet(**kwargs):
    name = "io_bazel_rules_dotnet"
    ref = get_ref(name, "0.0.4", kwargs) 
    sha256 = get_sha256(name, "8bd425654e142739b0da9ff182dbf735b7560ebd50b000627a02dba5fb2a759f", kwargs)
    github_archive(name, "bazelbuild", "rules_dotnet", ref, sha256)

def io_bazel_rules_scala(**kwargs):
    name = "io_bazel_rules_scala"
    ref = get_ref(name, "30b80b03a410994a8abb93d5a3f81b0d1f5cb96f", kwargs)  # Feb 3, 2020
    sha256 = get_sha256(name, "342ec47a670ecf5167857c83ac3281be9590c27ce8d4b04ce46fa63e16de7a58", kwargs)
    github_archive(name, "bazelbuild", "rules_scala", ref, sha256)

def io_bazel_rules_rust(**kwargs):
    name = "io_bazel_rules_rust"
    ref = get_ref(name, "d28b121396974a628b9cdb29b6ed7f4e370edb4e", kwargs)  # May 8, 2019
    sha256 = get_sha256(name, "58b8786e00b3489ce127e001670fd991547bb7db315e8a214915a2fa0b83743f", kwargs)
    github_archive(name, "bazelbuild", "rules_rust", ref, sha256)

def com_github_yugui_rules_ruby(**kwargs):
    name = "com_github_yugui_rules_ruby"
    ref = get_ref(name, "73479cdc6a34a8d940cc3c904badf7a2ae6bdc6d", kwargs)  # PR#8,
    sha256 = get_sha256(name, "bd88b1aa144f70bb3f069ff3ddc5ddba032311ce27fb40b7276db694dcb63490", kwargs)
    github_archive(name, "yugui", "rules_ruby", ref, sha256)

def grpc_ecosystem_grpc_gateway(**kwargs):
    name = "grpc_ecosystem_grpc_gateway"
    ref = get_ref(name, "79ff520b46091f8148bafeafd6e798826d6d47c2", kwargs)  # Apr 2019
    sha256 = get_sha256(name, "a8d283391d1e37b2bea798082f198187dd1edfed03da00f5be96edc6dadfde44", kwargs)
    github_archive(name, "grpc-ecosystem", "grpc-gateway", ref, sha256)

def org_pubref_rules_node(**kwargs):
    name = "org_pubref_rules_node"
    ref = get_ref(name, "9ebfa90ca3283bb0f92ae5f337173a5a5a4d98aa", kwargs)
    sha256 = get_sha256(name, "cb1bf3d64c0b323bc515748902df9fef9ecfcc37c7aa84253d7e99d876f1196a", kwargs)
    github_archive(name, "pubref", "rules_node", ref, sha256)

def build_bazel_rules_android(**kwargs):
    """Android Rules
    """
    name = "build_bazel_rules_android"
    ref = get_ref(name, "ac9d2df31b2b9c77c1c43d4a3cfce789758320c5", kwargs)  # Apr 2019
    sha256 = get_sha256(name, "5bbead25993489d50290d8361f16fd08958c6b7b78e260091a0fd4f691518fb9", kwargs)
    github_archive(name, "bazelbuild", "rules_android", ref, sha256)

def build_bazel_rules_swift(**kwargs):
    """swift Rules
    """
    name = "build_bazel_rules_swift"
    ref = get_ref(name, "004597eeb9b3e0a2bda6c9f6232f8687a01e73b0", kwargs)  # Apr 12, 2019
    sha256 = get_sha256(name, "1fa65113903c03d0c7ae0fd79cfcf5cda1695d96411afa58a8022569db53e3c2", kwargs)
    github_archive(name, "bazelbuild", "rules_swift", ref, sha256)

def com_github_apple_swift_swift_protobuf(**kwargs):
    if "com_github_apple_swift_swift_protobuf" not in native.existing_rules():
        version = "1.4.0"
        http_archive(
            name = "com_github_apple_swift_swift_protobuf",
            url = "https://github.com/apple/swift-protobuf/archive/%s.tar.gz" % version,
            sha256 = "efa256d572d19fc23756a30089129af523173ad29a84ee87800fa88f056efaac",
            strip_prefix = "swift-protobuf-" + version,
            build_file = "@build_bazel_rules_swift//third_party:com_github_apple_swift_swift_protobuf/BUILD.overlay",
        )

def io_bazel_rules_go(**kwargs):
    """Go Rules
    """
    name = "io_bazel_rules_go"
    ref = get_ref(name, "fabf03c1cd31bcf15fb945d932cef322b242be3a", kwargs)  # post-18.6
    sha256 = get_sha256(name, "d40144fcb282e167a1836c4d1a9aabdd457051717a4648153a89aa34bd9f8e6a", kwargs)
    github_archive(name, "bazelbuild", "rules_go", ref, sha256)

def io_bazel_rules_python(**kwargs):
    """python Rules
    """
    name = "io_bazel_rules_python"
    ref = get_ref(name, "965d4b4a63e6462204ae671d7c3f02b25da37941", kwargs)  # 2019-03-07
    sha256 = get_sha256(name, "3e55ec4f7e151b048e950965f956c1e0633fc76449905f40dba671574eac574c", kwargs)
    github_archive(name, "bazelbuild", "rules_python", ref, sha256)

def six(**kwargs):
    name = "six"
    if name not in native.existing_rules():
        http_archive(
            name = name,
            build_file_content = """
genrule(
  name = "copy_six",
  srcs = ["six-1.12.0/six.py"],
  outs = ["six.py"],
  cmd = "cp $< $(@)",
)

py_library(
  name = "six",
  srcs = ["six.py"],
  srcs_version = "PY2AND3",
  visibility = ["//visibility:public"],
)
        """,
            sha256 = "d16a0141ec1a18405cd4ce8b4613101da75da0e9a7aec5bdd4fa804d0e0eba73",
            urls = ["https://pypi.python.org/packages/source/s/six/six-1.12.0.tar.gz"],
        )

def io_bazel_rules_d(**kwargs):
    """d Rules
    """
    name = "io_bazel_rules_d"
    ref = get_ref(name, "99c22ceeac4b883f97b1a420f98d4540e47978ca", kwargs)
    sha256 = get_sha256(name, "ba8eb23c5753de0ba6e743e27e40f0eef1c3b08b3eaabd1bf782f87bca1ada2c", kwargs)
    github_archive(name, "bazelbuild", "rules_d", ref, sha256)

def com_google_protobuf_lite(**kwargs):
    """A different branch of google/protobuf that contains the protobuf_lite plugin
    """
    name = "com_google_protobuf_lite"
    ref = get_ref(name, "3cf3be9959928bf8a7133d323eaf6a5a8d5afdd7", kwargs)  # latest as of Apr 9, 2019
    sha256 = get_sha256(name, "9f28fd96ccd1f87e0b2d23c622db2a87c87ff91dc30dd13a6dd3bff11738e608", kwargs)
    github_archive(name, "protocolbuffers", "protobuf", ref, sha256)

def rules_jvm_external(**kwargs):
    """Fetch maven artifacts
    """
    name = "rules_jvm_external"
    ref = get_ref(name, "e359007344bc53133e1e54c891670d08453d4827", kwargs)  # May 7 2019
    sha256 = get_sha256(name, "150c8cd5a3abe8b2da09235ebe5aedd0a379440d9f6a15d1c99c2b1e560a09f9", kwargs)
    github_archive(name, "bazelbuild", "rules_jvm_external", ref, sha256)

def io_grpc_grpc_java(**kwargs):
    """grpc java plugin and jars
    """
    name = "io_grpc_grpc_java"
    ref = get_ref(name, "v1.27.0", kwargs)  
    sha256 = get_sha256(name, "a23970d15ee790c2bf36544976977eb45d3498c3efecc304717d6fbd8ba0fcc8", kwargs)
    github_archive(name, "grpc", "grpc-java", ref, sha256)


def com_github_scalapb_scalapb(**kwargs):
    """scala compiler plugin
    """
    if "com_github_scalapb_scalapb" not in native.existing_rules():
        http_archive(
            name = "com_github_scalapb_scalapb",
            url = "https://github.com/scalapb/ScalaPB/releases/download/v0.8.0/scalapbc-0.8.0.zip",
            sha256 = "bda0b44b50f0a816342a52c34e6a341b1a792f2a6d26f4f060852f8f10f5d854",
            strip_prefix = "scalapbc-0.8.0/lib",
            build_file_content = """
java_import(
    name = "compilerplugin",
    jars = ["com.thesamet.scalapb.compilerplugin-0.8.0.jar"],
    visibility = ["//visibility:public"],
)
java_import(
    name = "scala-library",
    jars = ["org.scala-lang.scala-library-2.11.12.jar"],
    visibility = ["//visibility:public"],
)
            """,
        )

def com_github_stackb_grpc_js(**kwargs):
    """Grpc-web implementation (closure)
    """
    name = "com_github_stackb_grpc_js"
    ref = get_ref(name, "d075960a9e62846ce92ae1029a777c141809f489", kwargs)
    sha256 = get_sha256(name, "c0f422823486986ea965fd36a0f5d3380151516421a6de8b69b72778cf3798a4", kwargs)
    github_archive(name, "stackb", "grpc.js", ref, sha256)

def bazel_gazelle(**kwargs):
    """ build file generator for go
    """
    name = "bazel_gazelle"
    ref = get_ref(name, "63ddd72aa315d020456f1a96bc6fcca9405810cb", kwargs)
    sha256 = get_sha256(name, "bff502b74b6ad77d0b9b558ebe99b030d6ba9ab0e3a8b4cb396448bf7fe88ab4", kwargs)
    github_archive(name, "bazelbuild", "bazel-gazelle", ref, sha256)

def bazel_skylib(**kwargs):
    """ bazel utils
    """
    name = "bazel_skylib"
    ref = get_ref(name, "be3b1fc838386bdbea39d9750ea4411294870575", kwargs)  # Apr 13, 2019
    sha256 = get_sha256(name, "6128dd2af9830430e0ae404cb6fdce754fb80ed88942e1a0865a7f376bb68c4e", kwargs)
    github_archive(name, "bazelbuild", "bazel-skylib", ref, sha256)

def com_github_grpc_grpc_web(**kwargs):
    """Rule for grpc-web
    """
    name = "com_github_grpc_grpc_web"
    ref = get_ref(name, "ffe8e9c9036f4ec7d5b55da75b1758b1f57fbf8d", kwargs)
    sha256 = get_sha256(name, "936ca06fe7a9b55c1e334e4869e1d153fec68d92d750d2b550e41e1c5580b4dd", kwargs)
    github_archive(name, "grpc", "grpc-web", ref, sha256)

def bazel_version():
    """Rule for setting the bazel version repository
    """
    bazel_version_repository(
        name = "upb_bazel_version",
    )

# Write version data. Required for both upb and rules_rust
def _store_bazel_version(repository_ctx):
    repository_ctx.file("BUILD", "exports_files(['def.bzl'])")
    repository_ctx.file("bazel_version.bzl", "bazel_version = \"{}\"".format(native.bazel_version))
    repository_ctx.file("def.bzl", "BAZEL_VERSION='{}'".format(native.bazel_version))

bazel_version_repository = repository_rule(
    implementation = _store_bazel_version,
)