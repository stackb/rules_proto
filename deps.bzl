load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:jvm.bzl", "jvm_maven_import_external")
load("//:utils.bzl", "github_archive", "get_ref", "get_sha256")

# Versions
VERSIONS = {
    # Core
    "com_google_protobuf": {
        "type": "github",
        "org": "protocolbuffers",
        "repo": "protobuf",
        "ref": "2d0fdcb791dd59599cce812841149ff779b7c371", # 3.9.0 rc1 + PR 6310
        "sha256": "a6d34f9f32421425f6e2971de121356cf61efa576100dd37762b83869f05b17a",
        "binds": [
            {
                "name": "protobuf_clib",
                "actual": "@com_google_protobuf//:protoc_lib",
            },
            {
                "name": "protobuf_headers",
                "actual": "@com_google_protobuf//:protobuf_headers",
            },
        ],
    },
    "com_github_grpc_grpc": {
        "type": "github",
        "org": "grpc",
        "repo": "grpc",
        "ref": "v1.21.0", # 3.9.0 rc1 + PR 6310
        "sha256": "8da7f32cc8978010d2060d740362748441b81a34e5425e108596d3fcd63a97f2",
    },
    "zlib": {
        "type": "http",
        "urls": ["https://zlib.net/zlib-1.2.11.tar.gz"],
        "sha256": "c3e5e9fdd5004dcb542feda5ee4f0ff0744628baf8ed2dd5d66f8ca1197cb1a1",
        "strip_prefix": "zlib-1.2.11",
        "build_file": "@com_google_protobuf//:third_party/zlib.BUILD",
    },

    # Misc
    "bazel_skylib": {
        "type": "github",
        "org": "bazelbuild",
        "repo": "bazel-skylib",
        "ref": "be3b1fc838386bdbea39d9750ea4411294870575", # Apr 13, 2019
        "sha256": "6128dd2af9830430e0ae404cb6fdce754fb80ed88942e1a0865a7f376bb68c4e",
    },
    "com_github_bazelbuild_buildtools": {
        "type": "github",
        "org": "bazelbuild",
        "repo": "buildtools",
        "ref": "6415663945d3248207da955aafa1fa2af1a0f2ed",
        "sha256": "d1e28237d1f4c2255c504246b4f3fd36f74d590f2974491b4399a84c58b495ed",
    },

    # Android
    "build_bazel_rules_android": {
        "type": "github",
        "org": "bazelbuild",
        "repo": "rules_android",
        "ref": "9ab1134546364c6de84fc6c80b4202fdbebbbb35", # 2019-06-19
        "sha256": "f329928c62ade05ceda72c4e145fd300722e6e592627d43580dd0a8211c14612",
    },
    "com_google_guava_guava_android": {
        "type": "jvm_maven_import_external",
        "artifact": "com.google.guava:guava:27.0.1-android",
        "server_urls": ["http://central.maven.org/maven2"],
        "artifact_sha256": "caf0955aed29a1e6d149f85cfb625a89161b5cf88e0e246552b7ffa358204e28",
    },
    "com_google_protobuf_javalite": {
        "type": "github",
        "org": "protocolbuffers",
        "repo": "protobuf",
        "ref": "fa08222434bc58d743e8c2cc716bc219c3d0f44e",
        "sha256": "b04b08d31208be32aafdf5842d1b6073d527a67ff8d2cf4b17ee8f22a5273758",
    },

    # csharp
    "io_bazel_rules_dotnet": {
        "type": "github",
        "org": "bazelbuild",
        "repo": "rules_dotnet",
        "ref": "e9537b4a545528b11b270dfa124f3193bdb2d78e", # June 26, 2019
        "sha256": "9ee5429417190f00b2c970ba628db833e7ce71323efb646b9ce6b3aaaf56f125",
    },

    # Closure
    "io_bazel_rules_closure": {
        "type": "github",
        "org": "bazelbuild",
        "repo": "rules_closure",
        "ref": "03110588392d8c6c05b99c08a6f1c2121604ca27",
        "sha256": "12d17f241fc73f562825b0bba11a21040580d82da00be4689f5942c05508d8da",
    },

    # D
    "io_bazel_rules_d": {
        "type": "github",
        "org": "bazelbuild",
        "repo": "rules_d",
        "ref": "c4af62269c85dd5dcab0be119196baa5da4662b6", # June 28, 2019 + PR 30
        "sha256": "ef380076035d42bfc8b9a5547092779792de4b0cf718b9623a7c1923b0cd23e6",
    },
    "com_github_dcarp_protobuf_d": {
        "type": "http",
        "urls": ["https://github.com/dcarp/protobuf-d/archive/v0.5.0.tar.gz"],
        "sha256": "67a037dc29242f0d2f099746da67f40afff27c07f9ab48dda53d5847620db421",
        "strip_prefix": "protobuf-d-0.5.0",
        "build_file": Label("//d:com_github_dcarp_protobuf_d.BUILD.bazel"),
    },

    # Go
    "io_bazel_rules_go": {
        "type": "github",
        "org": "bazelbuild",
        "repo": "rules_go",
        "ref": "6fc21c78143ff1d4ea98100e8fd7a928d45abd00", # 0.18.6
        "sha256": "d9f58122d7cece7c73ddd4408e90ba0ac48bf45b58de74550cc446319ad61617",
    },
    "bazel_gazelle": {
        "type": "github",
        "org": "bazelbuild",
        "repo": "bazel-gazelle",
        "ref": "63ddd72aa315d020456f1a96bc6fcca9405810cb",
        "sha256": "bff502b74b6ad77d0b9b558ebe99b030d6ba9ab0e3a8b4cb396448bf7fe88ab4",
    },

    # gRPC gateway
    "grpc_ecosystem_grpc_gateway": {
        "type": "github",
        "org": "grpc-ecosystem",
        "repo": "grpc-gateway",
        "ref": "79ff520b46091f8148bafeafd6e798826d6d47c2", # Apr 2019
        "sha256": "a8d283391d1e37b2bea798082f198187dd1edfed03da00f5be96edc6dadfde44",
    },

    # gRPC web
    "com_github_grpc_grpc_web": {
        "type": "github",
        "org": "grpc",
        "repo": "grpc-web",
        "ref": "ffe8e9c9036f4ec7d5b55da75b1758b1f57fbf8d",
        "sha256": "936ca06fe7a9b55c1e334e4869e1d153fec68d92d750d2b550e41e1c5580b4dd",
    },

    # gRPC.js
    "com_github_stackb_grpc_js": {
        "type": "github",
        "org": "stackb",
        "repo": "grpc.js",
        "ref": "d075960a9e62846ce92ae1029a777c141809f489",
        "sha256": "c0f422823486986ea965fd36a0f5d3380151516421a6de8b69b72778cf3798a4",
    },

    # Java
    "rules_jvm_external": {
        "type": "github",
        "org": "bazelbuild",
        "repo": "rules_jvm_external",
        "ref": "e359007344bc53133e1e54c891670d08453d4827", # May 7 2019
        "sha256": "150c8cd5a3abe8b2da09235ebe5aedd0a379440d9f6a15d1c99c2b1e560a09f9",
    },
    "io_grpc_grpc_java": {
        "type": "github",
        "org": "grpc",
        "repo": "grpc-java",
        "ref": "46ef51576f7128874feba8067e13fc7b7e131f6f", # v1.21.0 + plugin fix
        "sha256": "c86c32a06a3db8d88839c682dc18e5acee130e29e4c009d64826fb580729fd91",
    },
    "javax_annotation_javax_annotation_api": {
        "type": "jvm_maven_import_external",
        "artifact": "javax.annotation:javax.annotation-api:1.2",
        "server_urls": ["http://central.maven.org/maven2"],
        "artifact_sha256": "5909b396ca3a2be10d0eea32c74ef78d816e1b4ead21de1d78de1f890d033e04",
        "licenses": ["reciprocal"], # CDDL License
    },
    "com_google_errorprone_error_prone_annotations": {
        "type": "jvm_maven_import_external",
        "artifact": "com.google.errorprone:error_prone_annotations:2.3.2",
        "server_urls": ["http://central.maven.org/maven2"],
        "artifact_sha256": "357cd6cfb067c969226c442451502aee13800a24e950fdfde77bcdb4565a668d",
        "licenses": ["notice"], # Apache 2.0
        "binds": [
            {
                "name": "error_prone_annotations",
                "actual": "@com_google_errorprone_error_prone_annotations//jar",
            },
        ]
    },

    # NodeJS
    "build_bazel_rules_nodejs": {
        "type": "github",
        "org": "bazelbuild",
        "repo": "rules_nodejs",
        "ref": "fb1133d1259992c196df31e2da6540fa6d54f191", # 0.32.2
        "sha256": "4e91b794f01cc8b1c690ae4d7ae8815ac6a4d4e9d2a5ac5988233ad0be564355",
    },

    # Python
    "com_apt_itude_rules_pip": {
        "type": "github",
        "org": "apt-itude",
        "repo": "rules_pip",
        "ref": "ce667087818553cdc4b1a2258fc53df917c4f87c", # 2019-07-07
        "sha256": "5cabd6bfb9cef095d0d076faf5e7acd5698f7172e803059c21c4e700a07b131b",
    },
    "six": {
        "type": "http",
        "urls": ["https://pypi.python.org/packages/source/s/six/six-1.12.0.tar.gz"],
        "sha256": "d16a0141ec1a18405cd4ce8b4613101da75da0e9a7aec5bdd4fa804d0e0eba73",
        "strip_prefix": "six-1.12.0",
        "build_file_content": """
genrule(
  name = "copy_six",
  srcs = ["six.py"],
  outs = ["__init__.py"],
  cmd = "cp $< $(@)",
)
py_library(
  name = "six",
  srcs = ["__init__.py"],
  srcs_version = "PY2AND3",
  visibility = ["//visibility:public"],
)
""",
    },

    # Ruby
    "com_github_yugui_rules_ruby": {
        "type": "github",
        "org": "yugui",
        "repo": "rules_ruby",
        "ref": "73479cdc6a34a8d940cc3c904badf7a2ae6bdc6d", # PR#8
        "sha256": "bd88b1aa144f70bb3f069ff3ddc5ddba032311ce27fb40b7276db694dcb63490",
    },

    # Rust
    "io_bazel_rules_rust": {
        "type": "github",
        "org": "bazelbuild",
        "repo": "rules_rust",
        "ref": "3fac9fe0001d2a829d8ddaf3033b5171c049abdb", # 2019-07-02
        "sha256": "299108772020c103eefacb4de30873d45224e8e0e6c11df7b56ffd11d959e212",
    },

    # Scala
    "io_bazel_rules_scala": {
        "type": "github",
        "org": "bazelbuild",
        "repo": "rules_scala",
        "ref": "14d9742496859faaf860b1adfc8126f3ed077921", # May 3, 2019
        "sha256": "72fc4357b29ec93951d472ee22a4cc3f30e170234a4ec73ff678f43f7e276bd4",
    },
    "com_github_scalapb_scalapb": {
        "type": "http",
        "urls": ["https://github.com/scalapb/ScalaPB/releases/download/v0.8.0/scalapbc-0.8.0.zip"],
        "sha256": "bda0b44b50f0a816342a52c34e6a341b1a792f2a6d26f4f060852f8f10f5d854",
        "strip_prefix": "scalapbc-0.8.0/lib",
        "build_file_content": """
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
    },

    # Swift
    "build_bazel_rules_swift": {
        "type": "github",
        "org": "bazelbuild",
        "repo": "rules_swift",
        "ref": "c935de3d04a8d24feb09a57df3b33a328be5d863", # 0.11.1
        "sha256": "797593aef1401c3fedfe0762ec073bfe7619ab8e4e26558614a1daa491e501a4",
    },
    "com_github_apple_swift_swift_protobuf": {
        "type": "http",
        "urls": ["https://github.com/apple/swift-protobuf/archive/1.4.0.tar.gz"],
        "sha256": "efa256d572d19fc23756a30089129af523173ad29a84ee87800fa88f056efaac",
        "strip_prefix": "swift-protobuf-1.4.0",
        "build_file": "@build_bazel_rules_swift//third_party:com_github_apple_swift_swift_protobuf/BUILD.overlay",
    },
}


def _generic_dependency(name, **kwargs):
    if name not in VERSIONS:
        fail("Name {} not in VERSIONS".format(name))
    dep = VERSIONS[name]

    if dep["type"] == "github":
        if name not in native.existing_rules():
            ref = get_ref(name, dep["ref"], kwargs)
            sha256 = get_sha256(name, dep["sha256"], kwargs)
            github_archive(name, dep["org"], dep["repo"], ref, sha256)
        else:
            pass
            #print("Dependency '{}' has already been declared, skipping".format(name))

    elif dep["type"] == "http":
        if name not in native.existing_rules():
            args = {k: v for k, v in dep.items() if k in ["urls", "sha256", "strip_prefix", "build_file", "build_file_content"]}
            http_archive(name = name, **args)
        else:
            pass
            #print("Dependency '{}' has already been declared, skipping".format(name))

    elif dep["type"] == "jvm_maven_import_external":
        if name not in native.existing_rules():
            args = {k: v for k, v in dep.items() if k in ["artifact", "server_urls", "artifact_sha256"]}
            jvm_maven_import_external(name = name, **args)
        else:
            pass
            #print("Dependency '{}' has already been declared, skipping".format(name))

    else:
        fail("Unknown dependency type {}".format(dep))

    if "binds" in dep:
        for bind in dep["binds"]:
            if bind["name"] not in native.existing_rules():
                native.bind(
                    name = bind["name"],
                    actual = bind["actual"],
                )


#
# Core
#
def external_zlib(**kwargs):
    _generic_dependency("zlib", **kwargs)

def com_google_protobuf(**kwargs):
    _generic_dependency("com_google_protobuf", **kwargs)

def com_github_grpc_grpc(**kwargs):
    _generic_dependency("com_github_grpc_grpc", **kwargs)


#
# Misc
#
def bazel_skylib(**kwargs):
    _generic_dependency("bazel_skylib", **kwargs)

def com_github_bazelbuild_buildtools(**kwargs):
    _generic_dependency("com_github_bazelbuild_buildtools", **kwargs)


#
# Android
#
def build_bazel_rules_android(**kwargs):
    _generic_dependency("build_bazel_rules_android", **kwargs)

def com_google_guava_guava_android(**kwargs):
    _generic_dependency("com_google_guava_guava_android", **kwargs)

def com_google_protobuf_javalite(**kwargs):
    _generic_dependency("com_google_protobuf_javalite", **kwargs)


#
# Closure
#
def io_bazel_rules_closure(**kwargs):
    _generic_dependency("io_bazel_rules_closure", **kwargs)


#
# csharp
#
def io_bazel_rules_dotnet(**kwargs):
    _generic_dependency("io_bazel_rules_dotnet", **kwargs)


#
# D
#
def io_bazel_rules_d(**kwargs):
    _generic_dependency("io_bazel_rules_d", **kwargs)

def com_github_dcarp_protobuf_d(**kwargs):
    _generic_dependency("com_github_dcarp_protobuf_d", **kwargs)


#
# Go
#
def io_bazel_rules_go(**kwargs):
    _generic_dependency("io_bazel_rules_go", **kwargs)

def bazel_gazelle(**kwargs):
    _generic_dependency("bazel_gazelle", **kwargs)


#
# gRPC gateway
#
def grpc_ecosystem_grpc_gateway(**kwargs):
    _generic_dependency("grpc_ecosystem_grpc_gateway", **kwargs)


#
# gRPC web
#
def com_github_grpc_grpc_web(**kwargs):
    _generic_dependency("com_github_grpc_grpc_web", **kwargs)


#
# gRPC.js
#
def com_github_stackb_grpc_js(**kwargs):
    _generic_dependency("com_github_stackb_grpc_js", **kwargs)


#
# Java
#
def rules_jvm_external(**kwargs):
    _generic_dependency("rules_jvm_external", **kwargs)

def io_grpc_grpc_java(**kwargs):
    _generic_dependency("io_grpc_grpc_java", **kwargs)

def javax_annotation_javax_annotation_api(**kwargs):
    # Use //stub:javax_annotation for neverlink=1 support.
    _generic_dependency("javax_annotation_javax_annotation_api", **kwargs)

def com_google_errorprone_error_prone_annotations(**kwargs):
    # Use //stub:javax_annotation for neverlink=1 support.
    _generic_dependency("javax_annotation_javax_annotation_api", **kwargs)


#
# NodeJS
#
def build_bazel_rules_nodejs(**kwargs):
    _generic_dependency("build_bazel_rules_nodejs", **kwargs)


#
# Python
#
def com_apt_itude_rules_pip(**kwargs):
    _generic_dependency("com_apt_itude_rules_pip", **kwargs)

def six(**kwargs):
    _generic_dependency("six", **kwargs)


#
# Ruby
#
def com_github_yugui_rules_ruby(**kwargs):
    _generic_dependency("com_github_yugui_rules_ruby", **kwargs)


#
# Rust
#
def io_bazel_rules_rust(**kwargs):
    _generic_dependency("io_bazel_rules_rust", **kwargs)


#
# Scala
#
def io_bazel_rules_scala(**kwargs):
    _generic_dependency("io_bazel_rules_scala", **kwargs)

def com_github_scalapb_scalapb(**kwargs):
    _generic_dependency("com_github_scalapb_scalapb", **kwargs)


#
# Swift
#
def build_bazel_rules_swift(**kwargs):
    _generic_dependency("build_bazel_rules_swift", **kwargs)

def com_github_apple_swift_swift_protobuf(**kwargs):
    _generic_dependency("com_github_apple_swift_swift_protobuf", **kwargs)
