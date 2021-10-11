"""proto_repostitory.bzl provides the proto_repository rule."""

# buildifier: disable=bzl-visibility
load("//rules/private:execution.bzl", "env_execute", "executable_extension")
load("@bazel_gazelle//internal:go_repository_cache.bzl", "read_cache_env")
load("@bazel_tools//tools/build_defs/repo:utils.bzl", "read_netrc", "use_netrc")

# We can't disable timeouts on Bazel, but we can set them to large values.
_GO_REPOSITORY_TIMEOUT = 86400

# Inspired by @bazel_tools//tools/build_defs/repo:http.bzl
def _get_auth(ctx, urls):
    """Given the list of URLs obtain the correct auth dict."""
    netrcfile = ""
    if "NETRC" in ctx.os.environ:
        netrcfile = ctx.os.environ["NETRC"]
    elif ctx.os.name.startswith("windows"):
        if "USERPROFILE" in ctx.os.environ:
            netrcfile = "%s/_netrc" % (ctx.os.environ["USERPROFILE"])
    elif "HOME" in ctx.os.environ:
        netrcfile = "%s/.netrc" % (ctx.os.environ["HOME"])

    if netrcfile and ctx.path(netrcfile).exists:
        netrc = read_netrc(ctx, netrcfile)
        return use_netrc(netrc, urls, {})

    return {}

def _proto_repository_impl(ctx):
    # stay
    if ctx.attr.urls:
        # HTTP mode
        ctx.download_and_extract(
            url = ctx.attr.urls,
            sha256 = ctx.attr.sha256,
            stripPrefix = ctx.attr.strip_prefix,
            type = ctx.attr.type,
            auth = _get_auth(ctx, ctx.attr.urls),
        )

    env = read_cache_env(ctx, str(ctx.path(Label("@bazel_gazelle_go_repository_cache//:go.env"))))
    env_keys = [
        # PATH is needed to locate git and other vcs tools.
        "PATH",

        # HOME is needed to locate vcs configuration files (.gitconfig).
        "HOME",

        # Settings below are used by vcs tools.
        "SSH_AUTH_SOCK",
        "SSL_CERT_FILE",
        "SSL_CERT_DIR",
        "HTTP_PROXY",
        "HTTPS_PROXY",
        "NO_PROXY",
        "http_proxy",
        "https_proxy",
        "no_proxy",
        "GIT_SSL_CAINFO",
        "GIT_SSH",
        "GIT_SSH_COMMAND",
    ]
    env.update({k: ctx.os.environ[k] for k in env_keys if k in ctx.os.environ})

    # Repositories are fetched. Determine if build file generation is needed.
    build_file_names = ctx.attr.build_file_name.split(",")
    existing_build_file = ""
    for name in build_file_names:
        path = ctx.path(name)
        if path.exists and not env_execute(ctx, ["test", "-f", path]).return_code:
            existing_build_file = name
            break

    generate = (ctx.attr.build_file_generation == "on" or ctx.attr.build_file_generation == "clean" or (not existing_build_file and ctx.attr.build_file_generation == "auto"))

    # remove any existing build files
    if ctx.attr.build_file_generation == "clean":
        cmd = ["find", ".", "-type", "f", "("]
        for i, name in enumerate(build_file_names):
            cmd += ["-name", name]
            if i + 1 < len(build_file_names):
                cmd.append("-o")
        cmd += [")", "-print", "-exec", "rm", "{}", "+"]
        result = env_execute(ctx, cmd, environment = env)
        if result.return_code:
            fail("failed to clean build files: " + result.stderr)

    # remove any other files
    if ctx.attr.deleted_files:
        result = env_execute(ctx, ["rm"] + [str(ctx.path(f)) for f in ctx.attr.deleted_files], environment = env)
        if result.return_code:
            fail("failed to remove deleted files: " + result.stderr)

    if generate:
        # Build file generation is needed.
        build_file_name = existing_build_file or build_file_names[0]

        # write an empty file to make sure the filegroup is satisfied
        ctx.file(ctx.attr.proto_repository_index_output_file, "")

        # Populate Gazelle directive at root build file and
        lines = ["# " + d for d in ctx.attr.build_directives] + [
            "",
            'exports_files(["%s"])' % ctx.attr.proto_repository_index_output_file,
        ]
        ctx.file(
            build_file_name,
            "\n".join(lines),
        )

        # Run Gazelle
        _gazelle = "@proto_repository_tools//:bin/gazelle{}".format(executable_extension(ctx))
        gazelle = ctx.path(Label(_gazelle))
        cmd = [
            gazelle,
            "-lang",
            ",".join(ctx.attr.languages),
            "-proto_language_config_file",
            ctx.path(ctx.attr.proto_language_config_file),
            "-proto_repository_index_output_file",
            ctx.path(ctx.attr.proto_repository_index_output_file),
            "-proto_repo_name",
            ctx.name,
            "-repo_root",
            ctx.path(""),
        ]
        if ctx.attr.build_file_name:
            cmd.extend(["-build_file_name", ctx.attr.build_file_name])
        if ctx.attr.build_tags:
            cmd.extend(["-build_tags", ",".join(ctx.attr.build_tags)])
        if ctx.attr.build_external:
            cmd.extend(["-external", ctx.attr.build_external])
        if ctx.attr.build_file_proto_mode:
            cmd.extend(["-proto", ctx.attr.build_file_proto_mode])
        if ctx.attr.proto_repository_index_files:
            cmd.extend(["-proto_repository_index_files", ",".join([str(ctx.path(i).realpath) for i in ctx.attr.proto_repository_index_files])])
        cmd.extend(ctx.attr.build_extra_args)
        cmd.append(ctx.path(""))
        result = env_execute(ctx, cmd, environment = env, timeout = _GO_REPOSITORY_TIMEOUT)
        if result.return_code:
            fail("failed to generate BUILD files: %s" % (
                result.stderr,
            ))
        if result.stderr:
            # buildifier: disable=print
            print("%s: %s" % (ctx.name, result.stderr))

    # Apply patches if necessary.
    patch(ctx)

proto_repository = repository_rule(
    implementation = _proto_repository_impl,
    attrs = {

        # Attributes for a repository that should be downloaded via HTTP.
        "urls": attr.string_list(),
        "strip_prefix": attr.string(),
        "type": attr.string(),
        "sha256": attr.string(),

        # Attributes for a repository that needs automatic build file generation
        "build_external": attr.string(
            values = [
                "",
                "external",
                "vendored",
            ],
        ),
        "build_file_name": attr.string(default = "BUILD.bazel,BUILD"),
        "build_file_generation": attr.string(
            default = "on",
            values = [
                "auto",
                "clean",
                "off",
                "on",
            ],
        ),
        "build_naming_convention": attr.string(
            values = [
                "go_default_library",
                "import",
                "import_alias",
            ],
            default = "import_alias",
        ),
        "build_tags": attr.string_list(),
        "build_file_proto_mode": attr.string(
            values = [
                "",
                "file",
                "default",
                "package",
                "disable",
                "disable_global",
                "legacy",
            ],
        ),
        "build_extra_args": attr.string_list(),
        # "build_config": attr.label(default = "@bazel_gazelle_go_repository_config//:WORKSPACE"),
        "build_directives": attr.string_list(default = []),

        # language specific configuration
        "proto_language_config_file": attr.label(),

        # language specific configuration
        "proto_repository_index_output_file": attr.string(default = "protoresolve.csv"),
        "proto_repository_index_files": attr.label_list(
            allow_files = True,
        ),
        "deleted_files": attr.string_list(),
        # Patches to apply after running gazelle.
        "patches": attr.label_list(),
        "patch_tool": attr.string(default = "patch"),
        "patch_args": attr.string_list(default = ["-p0"]),
        "patch_cmds": attr.string_list(default = []),
        "languages": attr.string_list(default = ["proto", "protobuf", "protoresolve"]),
    },
)

# Copied from @bazel_tools//tools/build_defs/repo:utils.bzl
def patch(ctx):
    """Implementation of patching an already extracted repository

    Args:
        ctx: the context object
    """
    bash_exe = ctx.os.environ["BAZEL_SH"] if "BAZEL_SH" in ctx.os.environ else "bash"
    for patchfile in ctx.attr.patches:
        command = "{patchtool} {patch_args} < {patchfile}".format(
            patchtool = ctx.attr.patch_tool,
            patchfile = ctx.path(patchfile),
            patch_args = " ".join([
                "'%s'" % arg
                for arg in ctx.attr.patch_args
            ]),
        )
        st = ctx.execute([bash_exe, "-c", command])
        if st.return_code:
            fail("Error applying patch %s:\n%s%s" %
                 (str(patchfile), st.stderr, st.stdout))
    for cmd in ctx.attr.patch_cmds:
        st = ctx.execute([bash_exe, "-c", cmd])
        if st.return_code:
            fail("Error applying patch command %s:\n%s%s" %
                 (cmd, st.stdout, st.stderr))
