load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:maven_rules.bzl", "maven_jar")

# Special thing to get around maven jar issues
load("//closure:buildozer_http_archive.bzl", "buildozer_http_archive")

# Special dart_sdk_repository
load("//dart:sdk.bzl", "dart_sdk_repository")

load("//dart:dart_pub_deps.bzl", "dart_pub_deps")


def github_archive(name, org, repo, ref, sha256): 
    """Declare an http_archive from github
    """
    if name not in native.existing_rules():
        http_archive(
            name = name,
            strip_prefix = repo + "-" + ref,
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

def external_zlib(**kwargs):
    com_github_madler_zlib(**kwargs)
    name = "zlib"
    if name not in native.existing_rules():
        native.bind(
            name = name,
            actual = "@com_github_madler_zlib//:z",
        )


def com_github_madler_zlib(**kwargs):
    if "com_github_madler_zlib" not in native.existing_rules():
       http_archive(
           name = "com_github_madler_zlib",
           build_file = "@com_github_grpc_grpc//third_party:zlib.BUILD",
           strip_prefix = "zlib-cacf7f1d4e3d44d871b605da3b647f07d718623f",
           url = "https://github.com/madler/zlib/archive/cacf7f1d4e3d44d871b605da3b647f07d718623f.tar.gz",
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
    sha256 = get_sha256(name, "7aa0ab179eff56241b6cded9cd07324af2395ad56d5478e2f7dabdb42b65d3fb", kwargs)
    github_archive(name, "nanopb", "nanopb", ref, sha256)


def com_google_protobuf(**kwargs):
    name = "com_google_protobuf"
    ref = get_ref(name, "b829ff2a4614ff25048944b2cdc8e43b6488fda0", kwargs) 
    sha256 = get_sha256(name, "5d82a718e271e7fda626f983628e4b4601221788c2244763a9e57eda4cc667dd", kwargs)
    github_archive(name, "google", "protobuf", ref, sha256)


def com_github_grpc_grpc(**kwargs):
    name = "com_github_grpc_grpc"
    ref = get_ref(name, "f0fbee59326402e6c1f0a7c6a1f8899fd2b0025a", kwargs) # Bazel 0.20.0 workspace fixes.
    sha256 = get_sha256(name, "6cea8d3dbe289872343aafe67736539b0a28ecd062d3ff4965f157f5d9d77d1d", kwargs)
    github_archive(name, "grpc", "grpc", ref, sha256)


# NOTE(pcj): Using a different version of dotnet here that seems to have a bad assembly reference.
# Create an issue for this.
def io_bazel_rules_dotnet(**kwargs):
    name = "io_bazel_rules_dotnet"
    ref = get_ref(name, "e1254bf789f5008cf0f3d26fb6031695413a2d50", kwargs) 
    sha256 = get_sha256(name, "4fdda2aeeef61a2ebbeac7d5885fe361adbc73cc93b395494a84b41865065b4d", kwargs)
    github_archive(name, "stackb", "rules_dotnet", ref, sha256)


def io_bazel_rules_scala(**kwargs):
    name = "io_bazel_rules_scala"
    ref = get_ref(name, "f33c6a659e3af540de35df1413f57f31d36d70c7", kwargs) 
    sha256 = get_sha256(name, "fc5c25ff314d53ed895a4b98960650daa5e55c9e5e7e57bb822d813059a2947d", kwargs)
    github_archive(name, "bazelbuild", "rules_scala", ref, sha256)


def io_bazel_rules_rust(**kwargs):
    name = "io_bazel_rules_rust"
    ref = get_ref(name, "88022d175adb48aa5f8904f95dfc716c543b3f1e", kwargs) 
    sha256 = get_sha256(name, "d9832945f0fa7097ee548bd6fecfc814bd19759561dd7b06723e1c6a1879aa71", kwargs)
    github_archive(name, "bazelbuild", "rules_rust", ref, sha256)


def com_github_yugui_rules_ruby(**kwargs):
    name = "com_github_yugui_rules_ruby"
    ref = get_ref(name, "5976385c9c4b94647bc95e8bf9d9989f1dee4ee3", kwargs) # PR#8, 
    sha256 = get_sha256(name, "7991ded3b902aba4c13fa7bdd67132edfcc279930b356737c1a3d3b2686d08c8", kwargs)
    github_archive(name, "yugui", "rules_ruby", ref, sha256)


def com_github_grpc_ecosystem_grpc_gateway(**kwargs):
    name = "com_github_grpc_ecosystem_grpc_gateway"
    ref = get_ref(name, "8aa3d3f00fbaea619d864e688cd045497aa30fe8", kwargs) # Oct 2, 2018, 
    sha256 = get_sha256(name, "e21e3a3ac5e2be3869474f1c703db5d4f72e459e07084e9228eaaf2484e0dd48", kwargs)
    github_archive(name, "grpc-ecosystem", "grpc-gateway", ref, sha256)


def org_pubref_rules_node(**kwargs):
    name = "org_pubref_rules_node"
    ref = get_ref(name, "144b4b705e5e56b3ea5a64f47b086388be939b1d", kwargs)  
    sha256 = get_sha256(name, "c5d8f923bdf1ba905d7e46ce718d564814c2fad13673bc7dc4a2fce20c946315", kwargs)
    github_archive(name, "pubref", "rules_node", ref, sha256)


def build_bazel_rules_android(**kwargs):
    """Android Rules
    """
    name = "build_bazel_rules_android"
    ref = get_ref(name, "60f03a20cefbe1e110ae0ac7f25359822e9ea24a", kwargs) 
    sha256 = get_sha256(name, "4305b6cf6b098752a19fdb1abdc9ae2e069f5ff61359bfc3c752e4b4c862d18e", kwargs)
    github_archive(name, "bazelbuild", "rules_android", ref, sha256)


def build_bazel_rules_swift(**kwargs):
    """swift Rules
    """
    name = "build_bazel_rules_swift"
    ref = get_ref(name, "8f594d9a9b39ce471064cc13d35c07ea77a24628", kwargs) 
    sha256 = get_sha256(name, "0b6c009c18f4a1235db65acf841bbe85043d290f39ebb142436c046476a823e5", kwargs)
    github_archive(name, "bazelbuild", "rules_swift", ref, sha256)


def com_github_apple_swift_swift_protobuf(**kwargs):
    if "com_github_apple_swift_swift_protobuf" not in native.existing_rules():
        version = "86c579d280d629416c7bd9d32a5dfacab8e0b0b4" # v1.2.0
        http_archive(
            name = "com_github_apple_swift_swift_protobuf",
            url = "https://github.com/apple/swift-protobuf/archive/%s.tar.gz" % version,
            sha256 = "edd677ff3ad01a4090902c80e7a28671a2f48fa6ce06726f0616678465f575f1",
            strip_prefix = "swift-protobuf-" + version,
            build_file = "@build_bazel_rules_swift//third_party:com_github_apple_swift_swift_protobuf/BUILD.overlay",
        )    

def io_bazel_rules_go(**kwargs):
    """Go Rules
    """
    name = "io_bazel_rules_go"
    ref = get_ref(name, "0f0d007c89dc67a5a34490acafc5195b191f5045", kwargs) # 0.15.3 
    sha256 = get_sha256(name, "75a187b761dd3437c0722e3ab9a5c0835afc0acdd2cd1dc08f5d4810f409d57d", kwargs)
    github_archive(name, "bazelbuild", "rules_go", ref, sha256)


def io_bazel_rules_python(**kwargs):
    """python Rules
    """
    name = "io_bazel_rules_python"
    ref = get_ref(name, "8b5d0683a7d878b28fffe464779c8a53659fc645", kwargs)  
    sha256 = get_sha256(name, "8b32d2dbb0b0dca02e0410da81499eef8ff051dad167d6931a92579e3b2a1d48", kwargs)
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
    ref = get_ref(name, "78a4e1ba257bbe9a9d7a064c8cde8c5317059e17", kwargs) # ~11/27/2018 
    sha256 = get_sha256(name, "7e699c457f45151e5c61dce6bdcaf14e4cb840d95af86c13a9e7eecc49fd39a3", kwargs)
    github_archive(name, "dart-lang", "rules_dart", ref, sha256)


def dart_sdk(**kwargs):
    """The dart sdk
    """
    name = "dart_sdk"
    if name not in native.existing_rules():
        dart_sdk_repository(
            name = name,
        )


def dart_pub_deps_protoc_plugin(**kwargs):
    """Dart pub dependencies for the dart protoc plugin
    """
    name = "dart_pub_deps_protoc_plugin"
    if name not in native.existing_rules():
        dart_pub_deps(
            name = name,
            spec = str(Label("//dart:pubspec.yaml")),

            # these overrides were determined by manually browsing pub.dartlang.org, 
            # starting at protoc_plugin and going through all transitive dependencies,
            # pinning them to the version specified there (basically, seems like latest)
            
            override = {
                "analyzer": "0.34.0",
                "args": "1.5.1",
                "async": "2.0.8",
                "charcode": "1.1.2",
                "collection": "1.14.11",
                "convert": "2.0.2",
                "crypto": "2.0.6",
                "csslib": "0.14.6",
                "dart_style": "1.1.2",
                "front_end": "0.1.7",
                "fixnum": "0.10.9",
                "html": "0.13.3+3",
                "kernel": "0.3.7",
                "logging": "0.11.3+2",
                "meta": "1.1.6",
                "package_config": "1.0.5",
                "path": "1.6.2",
                "plugin": "0.2.0+3",
                "pub_semver": "1.4.2",
                "source_span": "1.4.1",
                "string_scanner": "1.0.4",
                "typed_data": "1.1.6",
                "watcher": "0.9.7+10",
                "utf": "0.9.0+5",
                "yaml": "2.1.15",
            },
        )

def dart_pub_deps_grpc(**kwargs):
    """Dart pub dependencies for grpc
    """
    name = "dart_pub_deps_grpc"
    if name not in native.existing_rules():
        dart_pub_deps(
            name = name,
            spec = str(Label("//dart:pubspec-grpc.yaml")),

            # these overrides were determined by manually browsing pub.dartlang.org, 
            # starting at protoc_plugin and going through all transitive dependencies,
            # pinning them to the version specified there (basically, seems like latest)
            
            override = {
                "async": "2.0.8",
                "collection": "1.14.11",
                "googlapis_auth": "0.2.6",
                "crypto": "2.0.6",
                "convert": "2.0.2",
                "charcode": "1.1.2",
                "typed_data": "1.1.6",
                "http": "0.12.0",
                "http_parser": "3.1.3",
                "source_span": "1.4.1",
                "string_scanner": "1.0.4",
                "path": "1.6.2",
                "http2": "0.1.9",
                "meta": "1.1.6",
            },
        )


def com_google_protobuf_lite(**kwargs):
    """A different branch of google/protobuf that contains the protobuf_lite plugin
    """
    name = "com_google_protobuf_lite"
    ref = get_ref(name, "5e8916e881c573c5d83980197a6f783c132d4276", kwargs) 
    sha256 = get_sha256(name, "d35902fb3cbe9afa67aad4e615a8224d0a531b8c06d32e100bdb235244748a3d", kwargs)
    github_archive(name, "protocolbuffers", "protobuf", ref, sha256)


def gmaven_rules(**kwargs):
    """A catalog of maven & android jars on google maven server
    """
    name = "gmaven_rules"
    ref = get_ref(name, "20180927-1", kwargs) 
    sha256 = get_sha256(name, "ddaa0f5811253e82f67ee637dc8caf3989e4517bac0368355215b0dcfa9844d6", kwargs)
    github_archive(name, "bazelbuild", "gmaven_rules", ref, sha256)


def io_grpc_grpc_java(**kwargs):
    """grpc java plugin and jars
    """
    name = "io_grpc_grpc_java"
    ref = get_ref(name, "3134daf471f90b8f00c037518fc64988a1cdc8f7", kwargs) # v1.15.0 
    sha256 = get_sha256(name, "a7d7def13fd019255ba6ef7499aa91dac38d0ec0f5d9c1262a75ae82f4d67174", kwargs)
    github_archive(name, "grpc", "grpc-java", ref, sha256)


def com_google_guava_guava(**kwargs):
    """grpc java plugin and jars
    """
    name = "com_google_guava_guava"
    artifact = get_artifact(name, "com.google.guava:guava:20.0", kwargs)
    sha1 = get_sha1(name, "89507701249388e1ed5ddcf8c41f4ce1be7831ef", kwargs)
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
            """
        )    


def io_bazel_rules_closure(**kwargs):
    name = "io_bazel_rules_closure"
    ref = get_ref(name, "1e12aa5612d758daf2df339991c8d187223a7ee6", kwargs) 
    sha256 = get_sha256(name, "663424d34fd067a8d066308eb2887fcaba36d73b354669ec1467498726a6b82c", kwargs)

    if "io_bazel_rules_closure" not in native.existing_rules():
        buildozer_http_archive(
            name = "io_bazel_rules_closure",
            urls = ["https://github.com/bazelbuild/rules_closure/archive/%s.tar.gz" % ref],
            sha256 = sha256,
            strip_prefix = "rules_closure-" + ref,
            label_list = ["//...:%java_binary", "//...:%java_library"],
            replace_deps = {
              "@com_google_code_findbugs_jsr305": "@com_google_code_findbugs_jsr305_3_0_0",
              "@com_google_errorprone_error_prone_annotations": "@com_google_errorprone_error_prone_annotations_2_1_3",
            },
            sed_replacements = {
              "closure/repositories.bzl": [
                "s|com_google_code_findbugs_jsr305|com_google_code_findbugs_jsr305_3_0_0|g",
                "s|com_google_errorprone_error_prone_annotations|com_google_errorprone_error_prone_annotations_2_1_3|g",
              ],
            },
        )


def com_github_stackb_grpc_js(**kwargs):
    """Grpc-web implementation (closure)
    """
    name = "com_github_stackb_grpc_js"
    ref = get_ref(name, "fb3d7dbd8bfc8e1d4fb259f76d75f59cd4b67c31", kwargs)  
    sha256 = get_sha256(name, "019c2fbaf3958fff3a8a68d775dde606f0eb77163f1256e1d700e78c9ed2db85", kwargs)
    github_archive(name, "stackb", "grpc.js", ref, sha256)


def build_bazel_rules_nodejs(**kwargs):
    """Rule node.js 
    """
    name = "build_bazel_rules_nodejs"
    ref = get_ref(name, "d334fd8e2274fb939cf447106dced97472534e80", kwargs)  
    sha256 = get_sha256(name, "5c69bae6545c5c335c834d4a7d04b888607993027513282a5139dbbea7166571", kwargs)
    github_archive(name, "bazelbuild", "rules_nodejs", ref, sha256)


def build_bazel_rules_typescript(**kwargs):
    """Rule for typescript 
    """
    name = "build_bazel_rules_typescript"
    ref = get_ref(name, "3488d4fb89c6a02d79875d217d1029182fbcd797", kwargs)  
    sha256 = get_sha256(name, "22ebe19999ce34de2f0329d29c7cac1cccd449cd61d0813aa0e633ac8dfaef80", kwargs)
    github_archive(name, "bazelbuild", "rules_typescript", ref, sha256)


def io_bazel_rules_webtesting(**kwargs):
    """Rule for browser testing 
    """
    name = "io_bazel_rules_webtesting"
    ref = get_ref(name, "e417b122a3d1e8f8a4cc09b1b05e2a5f52c8ecbb", kwargs)  
    sha256 = get_sha256(name, "76af36aac2aed3ca6984e0c45d1204582243fa6b81165bf0e42eacfffed9c769", kwargs)
    github_archive(name, "bazelbuild", "rules_webtesting", ref, sha256)


def ts_protoc_gen(**kwargs):
    """ d.ts generator 
    """
    name = "ts_protoc_gen"
    ref = get_ref(name, "67e0c93fa4e29539ec57c97d23116f366ffb94e7", kwargs)  
    sha256 = get_sha256(name, "03b87365116b3e829b648c57bf17b8d252b682c33264de5663e7f410c9ad0e55", kwargs)
    github_archive(name, "improbable-eng", "ts-protoc-gen", ref, sha256)


def bazel_gazelle(**kwargs):
    """ build file generator for go 
    """
    name = "bazel_gazelle"
    ref = get_ref(name, "6a1b93cc9b1c7e55e7d05a6d324bcf9d87ea3ab1", kwargs)  
    sha256 = get_sha256(name, "bc493cce447c02b361393a79e562a5f48f456705417ee76009a761a159540dd7", kwargs)
    github_archive(name, "bazelbuild", "bazel-gazelle", ref, sha256)


def bazel_skylib(**kwargs):
    """ bazel utils 
    """
    name = "bazel_skylib"
    ref = get_ref(name, "8cecf885c8bf4c51e82fd6b50b9dd68d2c98f757", kwargs)  
    sha256 = get_sha256(name, "68ef2998919a92c2c9553f7a6b00a1d0615b57720a13239c0e51d0ded5aa452a", kwargs)
    github_archive(name, "bazelbuild", "bazel-skylib", ref, sha256)


def com_github_grpc_grpc_web(**kwargs):
    """Rule for grpc-web 
    """
    name = "com_github_grpc_grpc_web"
    ref = get_ref(name, "92aa9f8fc8e7af4aadede52ea075dd5790a63b62", kwargs)  
    sha256 = get_sha256(name, "f4996205e6d1d72e2be46f1bda4d26f8586998ed42021161322d490537d8c9b9", kwargs)
    github_archive(name, "grpc", "grpc-web", ref, sha256)
