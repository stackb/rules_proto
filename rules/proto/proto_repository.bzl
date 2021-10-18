"""proto_repostitory.bzl provides the proto_repository rule."""

# Copyright 2014 The Bazel Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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
    # TODO(#549): vcs repositories are not cached and still need to be fetched.
    # Download the repository or module.
    fetch_repo_args = None

    if ctx.attr.urls:
        # HTTP mode
        for key in ("commit", "tag", "vcs", "remote", "version", "sum", "replace"):
            if getattr(ctx.attr, key):
                fail("cannot specifiy both urls and %s" % key, key)
        ctx.download_and_extract(
            url = ctx.attr.urls,
            sha256 = ctx.attr.sha256,
            stripPrefix = ctx.attr.strip_prefix,
            type = ctx.attr.type,
            auth = _get_auth(ctx, ctx.attr.urls),
        )
    elif ctx.attr.commit or ctx.attr.tag:
        # repository mode
        rev = None
        if ctx.attr.commit:
            rev = ctx.attr.commit
            rev_key = "commit"
        elif ctx.attr.tag:
            rev = ctx.attr.tag
            rev_key = "tag"
        for key in ("urls", "strip_prefix", "type", "sha256", "version", "sum", "replace"):
            if getattr(ctx.attr, key):
                fail("cannot specify both %s and %s" % (rev_key, key), key)

        if ctx.attr.vcs and not ctx.attr.remote:
            fail("if vcs is specified, remote must also be")

        fetch_repo_args = ["-dest", ctx.path(""), "-importpath", ctx.attr.importpath]
        if ctx.attr.remote:
            fetch_repo_args.extend(["--remote", ctx.attr.remote])
        if rev:
            fetch_repo_args.extend(["--rev", rev])
        if ctx.attr.vcs:
            fetch_repo_args.extend(["--vcs", ctx.attr.vcs])
    elif ctx.attr.version:
        # module mode
        for key in ("urls", "strip_prefix", "type", "sha256", "commit", "tag", "vcs", "remote"):
            if getattr(ctx.attr, key):
                fail("cannot specify both version and %s" % key)
        if not ctx.attr.sum:
            fail("if version is specified, sum must also be")

        fetch_path = ctx.attr.replace if ctx.attr.replace else ctx.attr.importpath
        fetch_repo_args = [
            "-dest=" + str(ctx.path("")),
            "-importpath=" + fetch_path,
            "-version=" + ctx.attr.version,
            "-sum=" + ctx.attr.sum,
        ]
    else:
        fail("one of urls, commit, tag, or importpath must be specified")

    env = read_cache_env(ctx, str(ctx.path(Label("@bazel_gazelle_go_repository_cache//:go.env"))))
    env_keys = [
        # Respect user proxy and sumdb settings for privacy.
        # TODO(jayconrod): gazelle in go_repository mode should probably
        # not go out to the network at all. This means *the build*
        # goes out to the network. We tolerate this for downloading
        # archives, but finding module roots is a bit much.
        "GOPROXY",
        "GONOPROXY",
        "GOPRIVATE",
        "GOSUMDB",
        "GONOSUMDB",

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

    if fetch_repo_args:
        # Disable sumdb in fetch_repo. In module mode, the sum is a mandatory
        # attribute of go_repository, so we don't need to look it up.
        fetch_repo_env = dict(env)
        fetch_repo_env["GOSUMDB"] = "off"

        # Override external GO111MODULE, because it is needed by module mode, no-op in repository mode
        fetch_repo_env["GO111MODULE"] = "on"

        fetch_repo = str(ctx.path(Label("@bazel_gazelle_go_repository_tools//:bin/fetch_repo{}".format(executable_extension(ctx)))))
        result = env_execute(
            ctx,
            [fetch_repo] + fetch_repo_args,
            environment = fetch_repo_env,
            timeout = _GO_REPOSITORY_TIMEOUT,
        )
        if result.return_code:
            fail("failed to fetch %s: %s" % (ctx.name, result.stderr))
        if result.stderr:
            # buildifier: disable=print
            print("fetch_repo: " + result.stderr)

    # Repositories are fetched. Determine if build file generation is needed.
    build_file_names = ctx.attr.build_file_name.split(",")
    existing_build_file = ""
    for name in build_file_names:
        path = ctx.path(name)
        if path.exists and not env_execute(ctx, ["test", "-f", path]).return_code:
            existing_build_file = name
            break

    generate = (ctx.attr.build_file_generation == "on" or (not existing_build_file and ctx.attr.build_file_generation == "auto"))

    # remove any existing build files
    if ctx.attr.build_file_expunge:
        cmd = ["find", ".", "-type", "f", "("]
        for i, name in enumerate(build_file_names):
            cmd += ["-name", name]
            if i + 1 < len(build_file_names):
                cmd.append("-o")
        cmd += [")", "-print", "-exec", "rm", "{}", "+"]
        result = env_execute(ctx, cmd, environment = env)
        if result.return_code:
            fail("failed to expunge build files: " + result.stderr)

    # remove any other files
    if ctx.attr.deleted_files:
        result = env_execute(ctx, ["rm"] + [str(ctx.path(f)) for f in ctx.attr.deleted_files], environment = env)
        if result.return_code:
            fail("failed to remove deleted files: " + result.stderr)

    if generate:
        # Build file generation is needed.
        build_file_name = existing_build_file or build_file_names[0]

        # Populate Gazelle directive at root build file and
        lines = ["# " + d for d in ctx.attr.build_directives] + [
            "",
            'exports_files(["%s"])' % ctx.attr.imports_out,
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
            "-go_repository_mode",
            "-mode",
            "fix",
            "-repo_root",
            ctx.path(""),
            "-repo_config",
            ctx.path(ctx.attr.build_config),
            "-proto_repo_name",
            ctx.name,
        ]
        if ctx.attr.version:
            cmd.append("-go_repository_module_mode")
        if ctx.attr.importpath:
            cmd.extend(["-go_prefix", ctx.attr.importpath])
        if ctx.attr.build_file_name:
            cmd.extend(["-build_file_name", ctx.attr.build_file_name])
        if ctx.attr.build_tags:
            cmd.extend(["-build_tags", ",".join(ctx.attr.build_tags)])
        if ctx.attr.build_external:
            cmd.extend(["-external", ctx.attr.build_external])
        if ctx.attr.build_file_proto_mode:
            cmd.extend(["-proto", ctx.attr.build_file_proto_mode])
        if ctx.attr.build_naming_convention:
            cmd.extend(["-go_naming_convention", ctx.attr.build_naming_convention])

        # protobuf extension
        if ctx.attr.languages:
            cmd.extend(["-lang", ",".join(ctx.attr.languages)])
        if ctx.attr.cfgs:
            cfgs = ",".join([str(ctx.path(f).realpath) for f in ctx.attr.cfgs])
            cmd.extend(["-proto_configs", cfgs])
        if ctx.attr.imports_out:
            cmd.extend(["-proto_imports_out", ctx.path(ctx.attr.imports_out)])
        if ctx.attr.imports:
            protoimports = ",".join([str(ctx.path(lbl).realpath) for lbl in ctx.attr.imports])
            cmd.extend(["-proto_imports_in", protoimports])
        if ctx.attr.override_go_googleapis:
            cmd.extend(["-override_go_googleapis"])

        cmd.extend(ctx.attr.build_extra_args)
        cmd.append(ctx.path(""))

        # for arg in cmd:
        #     print("gazelle cmd arg", arg)
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

go_repository = repository_rule(
    implementation = _proto_repository_impl,
    attrs = {
        # Fundamental attributes of a go repository
        "importpath": attr.string(mandatory = False),  # True in go_repository

        # Attributes for a repository that should be checked out from VCS
        "commit": attr.string(),
        "tag": attr.string(),
        "vcs": attr.string(
            default = "",
            values = [
                "",
                "git",
                "hg",
                "svn",
                "bzr",
            ],
        ),
        "remote": attr.string(),

        # Attributes for a repository that should be downloaded via HTTP.
        "urls": attr.string_list(),
        "strip_prefix": attr.string(),
        "type": attr.string(),
        "sha256": attr.string(),

        # Attributes for a module that should be downloaded with the Go toolchain.
        "version": attr.string(),
        "sum": attr.string(),
        "replace": attr.string(),

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
            default = "auto",
            values = [
                "auto",
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
        "build_config": attr.label(default = "@bazel_gazelle_go_repository_config//:WORKSPACE"),
        "build_directives": attr.string_list(default = []),

        # Patches to apply after running gazelle.
        "patches": attr.label_list(),
        "patch_tool": attr.string(default = "patch"),
        "patch_args": attr.string_list(default = ["-p0"]),
        "patch_cmds": attr.string_list(default = []),

        # protobuf extension specific configuration
        "build_file_expunge": attr.bool(),
        "languages": attr.string_list(),
        "cfgs": attr.label_list(allow_files = True),
        "imports": attr.label_list(
            allow_files = True,
        ),
        "imports_out": attr.string(default = "imports.csv"),
        "deleted_files": attr.string_list(),
        "override_go_googleapis": attr.bool(),
    },
)

def proto_repository(**kwargs):
    kwargs.setdefault("languages", ["proto", "protobuf"])
    kwargs.setdefault("build_file_expunge", True)
    kwargs.setdefault("build_file_generation", "on")
    go_repository(**kwargs)

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
