BUILD_FILE = """load(
    "@build_bazel_rules_swift//swift/internal:swift_toolchain.bzl",
    "swift_toolchain",
)
swift_toolchain(
    name = "toolchain",
    arch = "{arch}",
    clang_executable = "{clang_executable}",
    os = "{os}",
    root = "{root}",
    visibility = ["//visibility:public"],
)
"""

def _swift_toolchain_impl(repository_ctx):
    """Creates BUILD targets for the Swift toolchain on Linux (see Dockerfile in this repo).
    Args:
      repository_ctx: The repository rule context.
    """
    repository_ctx.file("BUILD.bazel", BUILD_FILE.format(
        arch = repository_ctx.attr.arch,
        clang_executable = repository_ctx.attr.clang_executable,
        os = repository_ctx.attr.os,
        root = repository_ctx.attr.root,
    ))


_swift_toolchain = repository_rule(
    implementation = _swift_toolchain_impl,
    attrs = {
        "clang_executable": attr.string(
            doc = "Path to clang",
            default = "/usr/bin/clang",
        ),
        "arch": attr.string(
            doc = "Host architecture",
            default = "x86_64",
        ),
        "os": attr.string(
            doc = "Host os",
            default = "linux",
        ),
        "root": attr.string(
            doc = "Root swit dir (typically path/to/swiftc.dirname.dirname)",
            default = "/usr",
        ),
    },
)

def swift_toolchain(**kwargs):
    """
    swift_toolchain - configure toolchain for swift specifically for use within
    a docker sandbox.
    """

    # This should probably always be 'build_bazel_rules_swift_local_config', but
    # make it configurable just in case.
    kwargs["name"] = kwargs.pop("name", "build_bazel_rules_swift_local_config")
    _swift_toolchain(**kwargs)
