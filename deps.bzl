load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:maven_rules.bzl", "maven_jar")

# https://raw.githubusercontent.com/grpc/grpc/master/third_party/zlib.BUILD
ZLIB_BUILD = """
cc_library(
    name = "z",
    srcs = [
        "adler32.c",
        "compress.c",
        "crc32.c",
        "deflate.c",
        "infback.c",
        "inffast.c",
        "inflate.c",
        "inftrees.c",
        "trees.c",
        "uncompr.c",
        "zutil.c",
    ],
    hdrs = [
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
    ],
    includes = [
        ".",
    ],
    linkstatic = 1,
    visibility = [
        "//visibility:public",
    ],
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

def jar(name, artifact, sha1):
    """Declare a maven_jar
    """
    if name not in native.existing_rules():
        maven_jar(
            name = name,
            artifact = artifact,
            sha1 = sha1,
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

def com_github_madler_zlib(**kwargs):
    if "com_github_madler_zlib" not in native.existing_rules():
        http_archive(
            name = "com_github_madler_zlib",
            build_file_content = ZLIB_BUILD,
            strip_prefix = "zlib-cacf7f1d4e3d44d871b605da3b647f07d718623f",
            url = "https://github.com/madler/zlib/archive/cacf7f1d4e3d44d871b605da3b647f07d718623f.tar.gz",
            sha256 = "6d4d6640ca3121620995ee255945161821218752b551a1a180f4215f7d124d45",
        )

def external_zlib(**kwargs):
    com_github_madler_zlib(**kwargs)
    name = "zlib"
    if name not in native.existing_rules():
        native.bind(
            name = name,
            actual = "@com_github_madler_zlib//:z",
        )

def com_github_bazelbuild_bazel_gazelle(**kwargs):
    if "com_github_bazelbuild_bazel_gazelle" not in native.existing_rules():
        sha1 = "4bee5cae22da3b948d90293aff01928dd3b9f41a"
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
    ref = get_ref(name, "v1.20.1", kwargs)
    sha256 = get_sha256(name, "ba8b08a697b66e14af35da07753583cf32ff3d14dcd768f91b1bbe2e6c07c349", kwargs)
    github_archive(name, "grpc", "grpc", ref, sha256)

def io_bazel_rules_dotnet(**kwargs):
    name = "io_bazel_rules_dotnet"
    ref = get_ref(name, "7e907e130943d4c9391df6ad3b569e3e9b2efa9d", kwargs)  # PR#122
    sha256 = get_sha256(name, "17f6e070bb940441efadf516a0274db1db0e306130c279d09d400fc0b3c71899", kwargs)
    github_archive(name, "bazelbuild", "rules_dotnet", ref, sha256)

def io_bazel_rules_scala(**kwargs):
    name = "io_bazel_rules_scala"
    ref = get_ref(name, "14d9742496859faaf860b1adfc8126f3ed077921", kwargs)  # May 3, 2019
    sha256 = get_sha256(name, "72fc4357b29ec93951d472ee22a4cc3f30e170234a4ec73ff678f43f7e276bd4", kwargs)
    github_archive(name, "bazelbuild", "rules_scala", ref, sha256)

def io_bazel_rules_rust(**kwargs):
    name = "io_bazel_rules_rust"
    ref = get_ref(name, "2215277a2be52263ca5cd4e547cc4a50e320b828", kwargs)
    sha256 = get_sha256(name, "55d2ff891c25ebf589aff604c8f1b41afa3fe88dbc3b6f912cd44974111b413e", kwargs)
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
    ref = get_ref(name, "4fec67d1fcefe7c80d6f4cc2ae0841c9d90e429a", kwargs)  # post-18.3
    sha256 = get_sha256(name, "2ab9320c583b05b805a7f8f4005fd081606505a64308051008d32148d7f98e1f", kwargs)
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
        """,
            sha256 = "105f8d68616f8248e24bf0e9372ef04d3cc10104f1980f54d57b2ce73a5ad56a",
            urls = ["https://pypi.python.org/packages/source/s/six/six-1.10.0.tar.gz#md5=34eed507548117b2ab523ab14b2f8b55"],
        )

def io_bazel_rules_dart(**kwargs):
    """Dart Rules
    """
    name = "io_bazel_rules_dart"
    ref = get_ref(name, "07aa5a42827f74d707ad3abcd3edbc14c7cad837", kwargs)  # Mar 11 (fork of dart-lang/rules_dart)
    sha256 = get_sha256(name, "836aa1908fda2c5f25f5a8dc298399d60252006c2953d170b95a33bfc5b5de14", kwargs)
    github_archive(name, "FKint", "rules_dart", ref, sha256)

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

def gmaven_rules(**kwargs):
    """A catalog of maven & android jars on google maven server
    """
    name = "gmaven_rules"
    ref = get_ref(name, "20180927-1", kwargs)
    sha256 = get_sha256(name, "ddaa0f5811253e82f67ee637dc8caf3989e4517bac0368355215b0dcfa9844d6", kwargs)
    github_archive(name, "bazelbuild", "gmaven_rules", ref, sha256)

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
    ref = get_ref(name, "3c24dc6fe1b8f3e5c89b919c38a4eefe216397d3", kwargs)  # v1.19.0 and changes up to PR #5456
    sha256 = get_sha256(name, "1eeb136874a58a0a311a0701016aced96919f501ced0372013eb1708724ab046", kwargs)
    github_archive(name, "grpc", "grpc-java", ref, sha256)

def com_google_guava_guava(**kwargs):
    """grpc java plugin and jars
    """
    name = "com_google_guava_guava"
    artifact = get_artifact(name, "com.google.guava:guava:20.0", kwargs)
    sha1 = get_sha1(name, "89507701249388e1ed5ddcf8c41f4ce1be7831ef", kwargs)
    jar(name, artifact, sha1)

def com_google_guava_guava_android(**kwargs):
    """android-specific guava 
    """
    name = "com_google_guava_guava_android"
    artifact = get_artifact(name, "com.google.guava:guava:27.0.1-android", kwargs)
    sha1 = get_sha1(name, "b7e1c37f66ef193796ccd7ea6e80c2b05426182d", kwargs)
    jar(name, artifact, sha1)

def com_thesamet_scalapb_scalapb_json4s(**kwargs):
    """json parsing library for scala
    """
    name = "com_thesamet_scalapb_scalapb_json4s"
    artifact = get_artifact(name, "com.thesamet.scalapb:scalapb-json4s_2.12:0.7.1", kwargs)
    sha1 = get_sha1(name, "808eeb6cfaa359a6c6a3dd2ea2c0374caee30c28", kwargs)
    jar(name, artifact, sha1)

def org_json4s_json4s_jackson_2_12(**kwargs):
    """json parsing library for scala
    """
    name = "org_json4s_json4s_jackson_2_12"
    artifact = get_artifact(name, "org.json4s:json4s-jackson_2.12:3.6.1", kwargs)
    sha1 = get_sha1(name, "864cf214dcd5686929f1c7f8d61344195c828b35", kwargs)
    jar(name, artifact, sha1)

def org_json4s_json4s_core_2_12(**kwargs):
    """json parsing library for scala - core
    """
    name = "org_json4s_json4s_core_2_12"
    artifact = get_artifact(name, "org.json4s:json4s-core_2.12:3.6.1", kwargs)
    sha1 = get_sha1(name, "7a619365089281c6015b80c499ff3b3cb196572f", kwargs)
    jar(name, artifact, sha1)

def org_json4s_json4s_ast_2_12(**kwargs):
    """json parsing library for scala
    """
    name = "org_json4s_json4s_ast_2_12"
    artifact = get_artifact(name, "org.json4s:json4s-ast_2.12:3.6.1", kwargs)
    sha1 = get_sha1(name, "cf937592788dfa654acb9679b97eb1e691bf69f8", kwargs)
    jar(name, artifact, sha1)

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
