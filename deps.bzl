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
    name = "zlib"
    if name not in native.existing_rules():
        http_archive(
            name = name,
            build_file = "@com_google_protobuf//:third_party/zlib.BUILD",
            sha256 = "629380c90a77b964d896ed37163f5c3a34f6e6d897311f1df2a7016355c45eff",
            strip_prefix = "zlib-1.2.11",
            urls = ["https://github.com/madler/zlib/archive/v1.2.11.tar.gz"],
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
    ref = get_ref(name, "0.29.0", kwargs)
    sha256 = get_sha256(name, "f3ef44916e6be705ae862c0520bac6834dd2ff1d4ac7e5abc61fe9f12ce7a865", kwargs)
    github_archive(name, "bazelbuild", "buildtools", ref, sha256)

def boringssl(**kwargs):
    if "boringssl" not in native.existing_rules():
        http_archive(
            name = "boringssl",
            # on the chromium-stable-with-bazel branch
            url = "https://boringssl.googlesource.com/boringssl/+archive/dcd3e6e6ecddf059adb48fca45bc7346a108bdd9.tar.gz",
	    sha256 = "e051159ca173513a44588eb82a3af476f2819b225494fad26d2b825f1f6a668b",
        )

def com_github_nanopb_nanopb(**kwargs):
    name = "com_github_nanopb_nanopb"
    ref = get_ref(name, "5e5280de13b89c4a71ba2e270980a6435bfd17f1", kwargs)
    sha256 = get_sha256(name, "1eb1aafea696ba3744dde2238e6237a18a104d07eb38bc1db93392de210f5422", kwargs)
    github_archive(name, "nanopb", "nanopb", ref, sha256)

def com_google_protobuf(**kwargs):
    name = "com_google_protobuf"
    ref = get_ref(name, "v3.10.0", kwargs)
    sha256 = get_sha256(name, "758249b537abba2f21ebc2d02555bf080917f0f2f88f4cbe2903e0e28c4187ed", kwargs)
    github_archive(name, "protocolbuffers", "protobuf", ref, sha256)

def com_github_grpc_grpc(**kwargs):
    name = "com_github_grpc_grpc"
    ref = get_ref(name, "v1.24.1", kwargs)
    sha256 = get_sha256(name, "ffadb8c6bcd725b60c370484062363c4c476335fbd5f377dcc66ac9c91aeae03", kwargs)
    github_archive(name, "grpc", "grpc", ref, sha256)

def io_bazel_rules_dotnet(**kwargs):
    name = "io_bazel_rules_dotnet"
    ref = get_ref(name, "9a99d60543807196026000fb590f462331a91d19", kwargs)
    sha256 = get_sha256(name, "623932c09b8f90483fc3aa3cdef0f86c9c1270cdcebcd317841e903e0fa3d95a", kwargs)
    github_archive(name, "bazelbuild", "rules_dotnet", ref, sha256)

def io_bazel_rules_scala(**kwargs):
    name = "io_bazel_rules_scala"
    ref = get_ref(name, "f985e5e0d6364970be8d6f15d262c8b0e0973d1b", kwargs)
    sha256 = get_sha256(name, "4276b2ab877d6e1271825933eea00869248d32948d42770bfe4fedd491b2824c", kwargs)
    github_archive(name, "bazelbuild", "rules_scala", ref, sha256)

def io_bazel_rules_rust(**kwargs):
    name = "io_bazel_rules_rust"
    ref = get_ref(name, "d28b121396974a628b9cdb29b6ed7f4e370edb4e", kwargs)
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
    ref = get_ref(name, "b60d84e5635233d2d8fc905041576bd44cefbb94", kwargs)
    sha256 = get_sha256(name, "446e633c679260f971da74555369bf45f0165475354ad59a23416448e2571f5f", kwargs)
    github_archive(name, "bazelbuild", "rules_android", ref, sha256)

def build_bazel_rules_swift(**kwargs):
    """swift Rules
    """
    name = "build_bazel_rules_swift"
    ref = get_ref(name, "0.13.0", kwargs)
    sha256 = get_sha256(name, "617e568aa8263c454f63362f5ab837038da710d646510b8f4a6760ff6361f714", kwargs)
    github_archive(name, "bazelbuild", "rules_swift", ref, sha256)

def com_github_apple_swift_swift_protobuf(**kwargs):
    if "com_github_apple_swift_swift_protobuf" not in native.existing_rules():
        version = "1.7.0"
        http_archive(
            name = "com_github_apple_swift_swift_protobuf",
            url = "https://github.com/apple/swift-protobuf/archive/%s.tar.gz" % version,
            sha256 = "3654f75d1de8806678ea7c942903a6fcdaba477e0fc0a53439cdc381a5f3e4c0",
            strip_prefix = "swift-protobuf-" + version,
            build_file = "@build_bazel_rules_swift//third_party:com_github_apple_swift_swift_protobuf/BUILD.overlay",
        )

def io_bazel_rules_go(**kwargs):
    """Go Rules
    """
    name = "io_bazel_rules_go"
    ref = get_ref(name, "v0.19.5", kwargs)
    sha256 = get_sha256(name, "b2493537eab9f65715d9236223a38f40553e11d7ec27499f294792d1f10dcfc3", kwargs)
    github_archive(name, "bazelbuild", "rules_go", ref, sha256)

def rules_cc(**kwargs):
    name = "rules_cc"
    ref = get_ref(name, "a508235df92e71d537fcbae0c7c952ea6957a912", kwargs)
    sha256 = get_sha256(name, "", kwargs)
    github_archive(name, "bazelbuild", "rules_cc", ref, sha256)

def rules_java(**kwargs):
    name = "rules_java"
    ref = get_ref(name, "db9b3c7f7e7e2f4f66589259f0e2d332b00ae631", kwargs)
    sha256 = get_sha256(name, "", kwargs)
    github_archive(name, "bazelbuild", "rules_java", ref, sha256)

def rules_proto(**kwargs):
    name = "rules_proto"
    ref = get_ref(name, "97d8af4dc474595af3900dd85cb3a29ad28cc313", kwargs)
    sha256 = get_sha256(name, "602e7161d9195e50246177e7c55b2f39950a9cf7366f74ed5f22fd45750cd208", kwargs)
    github_archive(name, "bazelbuild", "rules_proto", ref, sha256)

def rules_python(**kwargs):
    """python Rules
    """
    name = "rules_python"
    ref = get_ref(name, "5aa465d5d91f1d9d90cac10624e3d2faf2057bd5", kwargs)
    sha256 = get_sha256(name, "e220053c4454664c09628ffbb33f245e65f5fe92eb285fbd0bc3a26f173f99d0", kwargs)
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
    ref = get_ref(name, "0579d30b7667a04b252489ab130b449882a7bdba", kwargs)
    sha256 = get_sha256(name, "bf8d7e7d76f4abef5a732614ac06c0ccffbe5aa5fdc983ea4fa3a81ec68e1f8c", kwargs)
    github_archive(name, "bazelbuild", "rules_d", ref, sha256)

def rules_jvm_external(**kwargs):
    """Fetch maven artifacts
    """
    name = "rules_jvm_external"
    ref = get_ref(name, "05ba43aa5b671269cf0dfe2f91ec8f26dcea998e", kwargs)
    sha256 = get_sha256(name, "02e33287aa6fa129be0a3569ddba0c84ef8eb8b1e5f6f5348373bee559447642", kwargs)
    github_archive(name, "bazelbuild", "rules_jvm_external", ref, sha256)

def io_grpc_grpc_java(**kwargs):
    """grpc java plugin and jars
    """
    name = "io_grpc_grpc_java"
    ref = get_ref(name, "v1.24.0", kwargs)
    sha256 = get_sha256(name, "8b495f58aaf75138b24775600a062bbdaa754d85f7ab2a47b2c9ecb432836dd1", kwargs)
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
    ref = get_ref(name, "1.0.2", kwargs)
    sha256 = get_sha256(name, "e5d90f0ec952883d56747b7604e2a15ee36e288bb556c3d0ed33e818a4d971f2", kwargs)
    github_archive(name, "bazelbuild", "bazel-skylib", ref, sha256)

def com_github_grpc_grpc_web(**kwargs):
    """Rule for grpc-web
    """
    name = "com_github_grpc_grpc_web"
    ref = get_ref(name, "ffe8e9c9036f4ec7d5b55da75b1758b1f57fbf8d", kwargs)
    sha256 = get_sha256(name, "936ca06fe7a9b55c1e334e4869e1d153fec68d92d750d2b550e41e1c5580b4dd", kwargs)
    github_archive(name, "grpc", "grpc-web", ref, sha256)
