load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("//internal:common.bzl", "check_bazel_minimum_version")

# Versions
MINIMUM_BAZEL_VERSION = "1.0.0"

# ENABLE_VERSION_NAGS = False
# VERSIONS = {
#     # Core
#     "rules_proto": {
#         "type": "github",
#         "org": "bazelbuild",
#         "repo": "rules_proto",
#         "ref": "97d8af4dc474595af3900dd85cb3a29ad28cc313",
#         "sha256": "602e7161d9195e50246177e7c55b2f39950a9cf7366f74ed5f22fd45750cd208",
#     },
#     "com_google_protobuf": { # When updating, also update Node.js requirements, Ruby requirements and C# requirements
#         "type": "github",
#         "org": "protocolbuffers",
#         "repo": "protobuf",
#         "ref": "v3.13.0",
#         "sha256": "9b4ee22c250fe31b16f1a24d61467e40780a3fbb9b91c3b65be2a376ed913a1a",
#         "binds": [
#             {
#                 "name": "protobuf_clib",
#                 "actual": "@com_google_protobuf//:protoc_lib",
#             },
#             {
#                 "name": "protobuf_headers",
#                 "actual": "@com_google_protobuf//:protobuf_headers",
#             },
#         ],
#     },
#     "com_github_grpc_grpc": { # When updating, also update Node.js requirements, Ruby requirements and C# requirements
#         "type": "github",
#         "org": "grpc",
#         "repo": "grpc",
#         "ref": "v1.32.0",
#         "sha256": "f880ebeb2ccf0e47721526c10dd97469200e40b5f101a0d9774eb69efa0bd07a",
#     },
#     "zlib": {
#         "type": "http",
#         "urls": [
#             "https://mirror.bazel.build/zlib.net/zlib-1.2.11.tar.gz",
#             "https://zlib.net/zlib-1.2.11.tar.gz",
#         ],
#         "sha256": "c3e5e9fdd5004dcb542feda5ee4f0ff0744628baf8ed2dd5d66f8ca1197cb1a1",
#         "strip_prefix": "zlib-1.2.11",
#         "build_file": "@rules_proto_grpc//third_party:BUILD.bazel.zlib",
#     },
#     "rules_python": {
#         "type": "github",
#         "org": "bazelbuild",
#         "repo": "rules_python",
#         "ref": "c064f7008a30f307ea7516cf52358a653011f82b",
#         "sha256": "b9cf39396181e8d4434625a3533240469ca21242442745bd0b672731555823b8",
#     },
#     "build_bazel_rules_swift": {
#         "type": "github",
#         "org": "bazelbuild",
#         "repo": "rules_swift",
#         "ref": "0.14.0",
#         "sha256": "fa746a50f442ea4bcce78b747182107b4f0041f868b285714364ce4508d19979",
#     },
#     "build_bazel_apple_support": {
#         "type": "github",
#         "org": "bazelbuild",
#         "repo": "apple_support",
#         "ref": "0.7.2",
#         "sha256": "519a3bc32132f7b5780e82c2fc6ad2a78d4b28b81561e6fd7b7e0b14ea110074",
#     },
#     "bazel_skylib": {
#         "type": "github",
#         "org": "bazelbuild",
#         "repo": "bazel-skylib",
#         "ref": "1.0.3",
#         "sha256": "7ac0fa88c0c4ad6f5b9ffb5e09ef81e235492c873659e6bb99efb89d11246bcb",
#     },


#     # Android
#     "build_bazel_rules_android": {
#         "type": "github",
#         "org": "bazelbuild",
#         "repo": "rules_android",
#         "ref": "9ab1134546364c6de84fc6c80b4202fdbebbbb35",
#         "sha256": "f329928c62ade05ceda72c4e145fd300722e6e592627d43580dd0a8211c14612",
#     },

#     # C#
#     "io_bazel_rules_dotnet": {
#         "type": "github",
#         "org": "bazelbuild",
#         "repo": "rules_dotnet",
#         "ref": "98cc58708e0ea150a8737e7f83a74f0f41ececf8",
#         "sha256": "1b61f931391cd449fa60bb823c511db30d0567ecc2c6ef9d393bfba391c9a2da",
#     },

#     # Closure
#     "io_bazel_rules_closure": {
#         "type": "github",
#         "org": "bazelbuild",
#         "repo": "rules_closure",
#         "ref": "0.11.0",
#         "sha256": "d66deed38a0bb20581c15664f0ab62270af5940786855c7adc3087b27168b529",
#     },

#     # D
#     "io_bazel_rules_d": {
#         "type": "github",
#         "org": "bazelbuild",
#         "repo": "rules_d",
#         "ref": "064206527421aa6d1a4c11fbee5c93fd1619d99e",
#         "sha256": "e514c8be2d029563fec81eb1ecf3a931491cc5a02960cb3ab4331ac5c9124893",
#     },
#     "com_github_dcarp_protobuf_d": {
#         "type": "http",
#         "urls": ["https://github.com/dcarp/protobuf-d/archive/v0.6.2.tar.gz"],
#         "sha256": "5509883fa042aa2e1c8c0e072e52c695fb01466f572bd828bcde06347b82d465",
#         "strip_prefix": "protobuf-d-0.6.2",
#         "build_file": "@rules_proto_grpc//third_party:BUILD.bazel.com_github_dcarp_protobuf_d",
#     },

#     # Go
#     "io_bazel_rules_go": {
#         "type": "github",
#         "org": "bazelbuild",
#         "repo": "rules_go",
#         "ref": "v0.24.3",
#         "sha256": "e37e7937141a1deea40ee2f06a7850fc520e2272de7aacd85ad8a2ace11d2e83",
#     },
#     "bazel_gazelle": {
#         "type": "github",
#         "org": "bazelbuild",
#         "repo": "bazel-gazelle",
#         "ref": "v0.21.1",
#         "sha256": "2423201f91471ea87925b81962258e27a22cd8ebb4fe355bf033dcf2ad668541",
#     },

#     # grpc-gateway
#     "grpc_ecosystem_grpc_gateway": {
#         "type": "github",
#         "org": "grpc-ecosystem",
#         "repo": "grpc-gateway",
#         "ref": "v1.15.0",
#         "sha256": "0630c364e47aa7f813dd92f1874c778e496251304719c65e959675b15f7c7f15",
#     },

#     # gRPC web
#     "com_github_grpc_grpc_web": {
#         "type": "github",
#         "org": "grpc",
#         "repo": "grpc-web",
#         "ref": "1.2.1",
#         "sha256": "23cf98fbcb69743b8ba036728b56dfafb9e16b887a9735c12eafa7669862ec7b",
#     },

#     # Java
#     "io_grpc_grpc_java": {
#         "type": "github",
#         "org": "grpc",
#         "repo": "grpc-java",
#         "ref": "v1.32.1",  # Bug in 1.32.0 release means 1.32.1 should be used
#         "sha256": "e5d691f80e7388035c34616a17830ec2687fb2ef5c5d9c9b79c605a7addb78ab",
#     },

#     # NodeJS
#     # Use .tar.gz in release assets, not the Github generated source .tar.gz
#     "build_bazel_rules_nodejs": {
#         "type": "http",
#         "urls": ["https://github.com/bazelbuild/rules_nodejs/releases/download/2.2.0/rules_nodejs-2.2.0.tar.gz"],
#         "sha256": "4952ef879704ab4ad6729a29007e7094aef213ea79e9f2e94cbe1c9a753e63ef",
#     },

#     # Python
#     "subpar": {
#         "type": "github",
#         "org": "google",
#         "repo": "subpar",
#         "ref": "2.0.0",
#         "sha256": "b80297a1b8d38027a86836dbadc22f55dc3ecad56728175381aa6330705ac10f",
#     },
#     "six": {
#         "type": "http",
#         "urls": ["https://pypi.python.org/packages/source/s/six/six-1.13.0.tar.gz"],
#         "sha256": "30f610279e8b2578cab6db20741130331735c781b56053c59c4076da27f06b66",
#         "strip_prefix": "six-1.13.0",
#         "build_file": "@rules_proto_grpc//third_party:BUILD.bazel.six",
#     },

#     # Ruby
#     "com_github_yugui_rules_ruby": {
#         "type": "github",
#         "org": "yugui",
#         "repo": "rules_ruby",
#         "ref": "73479cdc6a34a8d940cc3c904badf7a2ae6bdc6d", # PR#8
#         "sha256": "bd88b1aa144f70bb3f069ff3ddc5ddba032311ce27fb40b7276db694dcb63490",
#     },

#     # Rust
#     "io_bazel_rules_rust": {
#         "type": "github",
#         "org": "bazelbuild",
#         "repo": "rules_rust",
#         "ref": "e64700dc9b8b3869bce4f77b78c33cb9d088cc4b",
#         "sha256": "eb384450d3b89332b386173233daa66a71e13cf63fe6d9ee51bd09fba0eb41f2",
#     },

#     # Scala
#     "io_bazel_rules_scala": {
#         "type": "github",
#         "org": "bazelbuild",
#         "repo": "rules_scala",
#         "ref": "6280cdbdb03bbace36e5458ca73745b80a9fe467",
#         "sha256": "723ac4c2eda86c6a5d9cbe64bde36f17185e7205acf8064a2b8bb1aea2fbf831",
#     },
#     "com_github_scalapb_scalapb": {
#         "type": "http",
#         "urls": ["https://github.com/scalapb/ScalaPB/releases/download/v0.9.7/scalapbc-0.9.7.zip"],  # Matches version in https://github.com/bazelbuild/rules_scala/blob/master/scala_proto/private/scala_proto_default_repositories.bzl
#         "sha256": "623f626e97cca119b2a12c4e1d9a3c85aab9f9fd6dcb8dc22b4f704b824da94e",
#         "strip_prefix": "scalapbc-0.9.7",
#         "build_file": "@rules_proto_grpc//third_party:BUILD.bazel.com_github_scalapb_scalapb",
#     },

#     # Swift
#     "com_github_grpc_grpc_swift_patched": {
#         "type": "github",
#         "org": "grpc",
#         "repo": "grpc-swift",
#         "ref": "0.11.0",
#         "sha256": "82e0a3d8fe2b9ee813b918e1a674f5a7c6dc024abe08109a347b686db6e57432",
#         "build_file": "@build_bazel_rules_swift//third_party:com_github_grpc_grpc_swift/BUILD.overlay",
#         "patch_cmds": [
#             "sed -i.bak -e 's/if fileDescriptor.services.count > 0/if true/g' Sources/protoc-gen-swiftgrpc/main.swift",  # Make grpc plugin always output files
#             "rm Sources/protoc-gen-swiftgrpc/main.swift.bak",  # Remove .bak file we had to create above due to mac OS sed weirdness
#         ],
#     },
# }


# def _generic_dependency(name, **kwargs):
#     if name not in VERSIONS:
#         fail("Name {} not in VERSIONS".format(name))
#     dep = VERSIONS[name]

#     existing_rules = native.existing_rules()
#     if dep["type"] == "github":
#         # Resolve ref and sha256
#         ref = kwargs.get(name + "_ref", dep["ref"])
#         sha256 = kwargs.get(name + "_sha256", dep["sha256"])

#         # Fix GitHub naming quirk in path
#         strippedRef = ref
#         if strippedRef.startswith("v"):
#             strippedRef = ref[1:]

#         # Generate URLs
#         urls = [
#             "https://github.com/{}/{}/archive/{}.tar.gz".format(dep["org"], dep["repo"], ref),
#         ]

#         # Check for existing rule
#         if name not in existing_rules:
#             http_archive(
#                 name = name,
#                 strip_prefix = dep["repo"] + "-" + strippedRef,
#                 urls = urls,
#                 sha256 = sha256,
#                 **{k: v for k, v in dep.items() if k in ["build_file", "patch_cmds"]}
#             )
#         elif existing_rules[name]["kind"] != "http_archive":
#             if ENABLE_VERSION_NAGS:
#                 print("Dependency '{}' has already been declared with a different rule kind. Found {}, expected http_archive".format(
#                     name, existing_rules[name]["kind"],
#                 ))
#         elif existing_rules[name]["urls"] != tuple(urls):
#             if ENABLE_VERSION_NAGS:
#                 print("Dependency '{}' has already been declared with a different version. Found urls={}, expected {}".format(
#                     name, existing_rules[name]["urls"], tuple(urls)
#                 ))

#     elif dep["type"] == "http":
#         if name not in existing_rules:
#             args = {k: v for k, v in dep.items() if k in ["urls", "sha256", "strip_prefix", "build_file", "build_file_content"]}
#             http_archive(name = name, **args)
#         elif existing_rules[name]["kind"] != "http_archive":
#             if ENABLE_VERSION_NAGS:
#                 print("Dependency '{}' has already been declared with a different rule kind. Found {}, expected http_archive".format(
#                     name, existing_rules[name]["kind"],
#                 ))
#         elif existing_rules[name]["urls"] != tuple(dep["urls"]):
#             if ENABLE_VERSION_NAGS:
#                 print("Dependency '{}' has already been declared with a different version. Found urls={}, expected {}".format(
#                     name, existing_rules[name]["urls"], tuple(dep["urls"])
#                 ))

#     elif dep["type"] == "local":
#         if name not in existing_rules:
#             args = {k: v for k, v in dep.items() if k in ["path"]}
#             native.local_repository(name = name, **args)
#         elif existing_rules[name]["kind"] != "local_repository":
#             if ENABLE_VERSION_NAGS:
#                 print("Dependency '{}' has already been declared with a different rule kind. Found {}, expected local_repository".format(
#                     name, existing_rules[name]["kind"],
#                 ))
#         elif existing_rules[name]["path"] != dep["path"]:
#             if ENABLE_VERSION_NAGS:
#                 print("Dependency '{}' has already been declared with a different version. Found path={}, expected {}".format(
#                     name, existing_rules[name]["path"], dep["urls"]
#                 ))

#     else:
#         fail("Unknown dependency type {}".format(dep))

#     if "binds" in dep:
#         for bind in dep["binds"]:
#             if bind["name"] not in native.existing_rules():
#                 native.bind(
#                     name = bind["name"],
#                     actual = bind["actual"],
#                 )


# Write version data. Required for both upb and rules_rust
def _store_bazel_version(repository_ctx):
    repository_ctx.file("BUILD", "exports_files(['def.bzl'])")
    repository_ctx.file("bazel_version.bzl", "bazel_version = \"{}\"".format(native.bazel_version))
    repository_ctx.file("def.bzl", "BAZEL_VERSION='{}'".format(native.bazel_version))

bazel_version_repository = repository_rule(
    implementation = _store_bazel_version,
)


#
# Toolchains
#
def rules_proto_grpc_toolchains():
    check_bazel_minimum_version(MINIMUM_BAZEL_VERSION)
    native.register_toolchains(str(Label("//protobuf:protoc_toolchain")))


# #
# # Core
# #
# def rules_proto_grpc_repos(**kwargs):
#     check_bazel_minimum_version(MINIMUM_BAZEL_VERSION)

#     bazel_version_repository(name = "bazel_version")

#     rules_proto(**kwargs)
#     rules_python(**kwargs)
#     build_bazel_rules_swift(**kwargs)
#     build_bazel_apple_support(**kwargs)
#     bazel_skylib(**kwargs)

#     six(**kwargs)
#     com_google_protobuf(**kwargs)
#     com_github_grpc_grpc(**kwargs)
#     external_zlib(**kwargs)

# def rules_proto(**kwargs):
#     _generic_dependency("rules_proto", **kwargs)

# def rules_python(**kwargs):
#     _generic_dependency("rules_python", **kwargs)

# def build_bazel_rules_swift(**kwargs):
#     _generic_dependency("build_bazel_rules_swift", **kwargs)

# def build_bazel_apple_support(**kwargs):
#     _generic_dependency("build_bazel_apple_support", **kwargs)

# def com_google_protobuf(**kwargs):
#     _generic_dependency("com_google_protobuf", **kwargs)

# def com_github_grpc_grpc(**kwargs):
#     _generic_dependency("com_github_grpc_grpc", **kwargs)

# def external_zlib(**kwargs):
#     _generic_dependency("zlib", **kwargs)



# #
# # Misc
# #
# def bazel_skylib(**kwargs):
#     _generic_dependency("bazel_skylib", **kwargs)


# #
# # Android
# #
# def build_bazel_rules_android(**kwargs):
#     _generic_dependency("build_bazel_rules_android", **kwargs)


# #
# # Closure
# #
# def io_bazel_rules_closure(**kwargs):
#     _generic_dependency("io_bazel_rules_closure", **kwargs)


# #
# # C#
# #
# def io_bazel_rules_dotnet(**kwargs):
#     _generic_dependency("io_bazel_rules_dotnet", **kwargs)


# #
# # D
# #
# def io_bazel_rules_d(**kwargs):
#     _generic_dependency("io_bazel_rules_d", **kwargs)

# def com_github_dcarp_protobuf_d(**kwargs):
#     _generic_dependency("com_github_dcarp_protobuf_d", **kwargs)


# #
# # Go
# #
# def io_bazel_rules_go(**kwargs):
#     _generic_dependency("io_bazel_rules_go", **kwargs)

# def bazel_gazelle(**kwargs):
#     _generic_dependency("bazel_gazelle", **kwargs)


# #
# # gRPC gateway
# #
# def grpc_ecosystem_grpc_gateway(**kwargs):
#     _generic_dependency("grpc_ecosystem_grpc_gateway", **kwargs)


# #
# # gRPC web
# #
# def com_github_grpc_grpc_web(**kwargs):
#     _generic_dependency("com_github_grpc_grpc_web", **kwargs)


# #
# # Java
# #
# def io_grpc_grpc_java(**kwargs):
#     _generic_dependency("io_grpc_grpc_java", **kwargs)


# #
# # NodeJS
# #
# def build_bazel_rules_nodejs(**kwargs):
#     _generic_dependency("build_bazel_rules_nodejs", **kwargs)


# #
# # Python
# #
# def subpar(**kwargs):
#     _generic_dependency("subpar", **kwargs)

# def six(**kwargs):
#     _generic_dependency("six", **kwargs)


# #
# # Ruby
# #
# def com_github_yugui_rules_ruby(**kwargs):
#     _generic_dependency("com_github_yugui_rules_ruby", **kwargs)


# #
# # Rust
# #
# def io_bazel_rules_rust(**kwargs):
#     _generic_dependency("io_bazel_rules_rust", **kwargs)


# #
# # Scala
# #
# def io_bazel_rules_scala(**kwargs):
#     _generic_dependency("io_bazel_rules_scala", **kwargs)

# def com_github_scalapb_scalapb(**kwargs):
#     _generic_dependency("com_github_scalapb_scalapb", **kwargs)


# #
# # Swift
# #
# def com_github_grpc_grpc_swift_patched(**kwargs):
#     _generic_dependency("com_github_grpc_grpc_swift_patched", **kwargs)
