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


def com_google_protobuf(**kwargs):
    name = "com_google_protobuf"
    ref = get_ref(name, "48cb18e5c419ddd23d9badcfe4e9df7bde1979b2", kwargs) # ref referenced by com_github_grpc_grpc
    sha256 = get_sha256(name, "013cc34f3c51c0f87e059a12ea203087a7a15dca2e453295345e1d02e2b9634c", kwargs)
    github_archive(name, "google", "protobuf", ref, sha256)


def com_github_grpc_grpc(**kwargs):
    name = "com_github_grpc_grpc"
    ref = get_ref(name, "5f84445781ef29e50435c2eea661ca435a19b6bc", kwargs) # v1.15.0
    sha256 = get_sha256(name, "8c09d35806c857dcb3a383c81052cc7ba1c5855e827d0792a83bd57251821e29", kwargs)
    github_archive(name, "grpc", "grpc", ref, sha256)


def io_bazel_rules_dotnet(**kwargs):
    name = "io_bazel_rules_dotnet"
    ref = get_ref(name, "281e4eeb389a30b4fd5dcd99915b35c925eff1ba", kwargs) 
    sha256 = get_sha256(name, "6913f43275eb098d28b46e96215e1d7425a6d7301d6a2805920e8796f2734d8f", kwargs)
    github_archive(name, "bazelbuild", "rules_dotnet", ref, sha256)


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


def org_pubref_rules_node(**kwargs):
    name = "org_pubref_rules_node"
    ref = get_ref(name, "1c60708c599e6ebd5213f0987207a1d854f13e23", kwargs)  
    sha256 = get_sha256(name, "248efb149bfa86d9d778b43949351015b23a8339405a9878467a1583ff6df348", kwargs)
    github_archive(name, "pubref", "rules_node", ref, sha256)


def build_bazel_rules_android(**kwargs):
    """Android Rules
    """
    name = "build_bazel_rules_android"
    ref = get_ref(name, "60f03a20cefbe1e110ae0ac7f25359822e9ea24a", kwargs) 
    sha256 = get_sha256(name, "4305b6cf6b098752a19fdb1abdc9ae2e069f5ff61359bfc3c752e4b4c862d18e", kwargs)
    github_archive(name, "bazelbuild", "rules_android", ref, sha256)


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
    ref = get_ref(name, "32ec7eb8585ba618e1236ccce43d158b6bda2f7b", kwargs)  
    sha256 = get_sha256(name, "b8107ad02bb3f3667900f18a585bb490192d6a193e975e82b2766470da1c7ebd", kwargs)
    github_archive(name, "bazelbuild", "rules_dart", ref, sha256)


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
            spec = "//dart:pubspec.yaml",
            override = {
                "path": "1.6.2",
                "analyzer": "0.32.5",
                "crypto": "2.0.6",
                "async": "2.0.8",
                "fixnum": "0.10.8",
                "collection": "1.14.11",
                "dart_style": "1.1.3",
                "source_span": "1.4.1",
                "args": "1.5.0",
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
    ref = get_ref(name, "c94ef115b4e8eea526d5b54b829cfc7542f39bc5", kwargs)  
    sha256 = get_sha256(name, "bf3b7fca7803a9187e6d6780089cad593997c46d76c5d78ba3202ce8b5e424b2", kwargs)
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


def com_github_grpc_grpc_web(**kwargs):
    """Rule for grpc-web 
    """
    name = "com_github_grpc_grpc_web"
    ref = get_ref(name, "92aa9f8fc8e7af4aadede52ea075dd5790a63b62", kwargs)  
    sha256 = get_sha256(name, "f4996205e6d1d72e2be46f1bda4d26f8586998ed42021161322d490537d8c9b9", kwargs)
    github_archive(name, "grpc", "grpc-web", ref, sha256)
