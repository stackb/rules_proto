
load("@bazel_tools//tools/build_defs/repo:utils.bzl", "workspace_and_buildfile")

def _http_archive_impl(ctx):
    """Buildozer implementation of the http_archive rule."""
    if not ctx.attr.url and not ctx.attr.urls:
        fail("At least one of url and urls must be provided")
    if ctx.attr.build_file and ctx.attr.build_file_content:
        fail("Only one of build_file and build_file_content can be provided.")

    all_urls = []
    if ctx.attr.urls:
        all_urls = ctx.attr.urls
    if ctx.attr.url:
        all_urls = [ctx.attr.url] + all_urls

    ctx.download_and_extract(
        all_urls,
        "",
        ctx.attr.sha256,
        ctx.attr.type,
        ctx.attr.strip_prefix,
    )

    if ctx.os.name == "mac os x":
      buildozer_urls = [ctx.attr.buildozer_mac_url]
      buildozer_sha256 = ctx.attr.buildozer_mac_sha256
    else:
      buildozer_urls = [ctx.attr.buildozer_linux_url]
      buildozer_sha256 = ctx.attr.buildozer_linux_sha256

    ctx.download(
        buildozer_urls,
        output = "buildozer",
        sha256 = buildozer_sha256,
        executable = True,
    )


    if ctx.attr.label_list:
        args = ["./buildozer", "-root_dir", ctx.path(".")]
        args += ["replace deps %s %s" % (k, v) for k, v in ctx.attr.replace_deps.items()]
        args += ctx.attr.label_list

        result = ctx.execute(args, quiet = False)
        if result.return_code:
            fail("Buildozer failed: %s" % result.stderr)

    if ctx.attr.sed_replacements:
        sed = ctx.which("sed")
        if not sed:
            fail("sed utility not found")
        # For each file (dict key) in the target list...
        for filename, replacements in ctx.attr.sed_replacements.items():
            # And each sed replacement to make (dict value)...
            for replacement in replacements:
                args = [sed, "-i.bak", replacement, filename]
                # execute the replace on that file.
                result = ctx.execute(args, quiet = False)
                if result.return_code:
                    fail("Buildozer failed: %s" % result.stderr)

    workspace_and_buildfile(ctx)


_http_archive_attrs = {
    "url": attr.string(),
    "urls": attr.string_list(),
    "sha256": attr.string(),
    "strip_prefix": attr.string(),
    "type": attr.string(),
    "build_file": attr.label(allow_single_file = True),
    "build_file_content": attr.string(),
    "replace_deps": attr.string_dict(),
    "sed_replacements": attr.string_list_dict(),
    "label_list": attr.string_list(),
    "workspace_file": attr.label(allow_single_file = True),
    "workspace_file_content": attr.string(),
    "buildozer_linux_url": attr.string(
        default = "https://github.com/bazelbuild/buildtools/releases/download/0.15.0/buildozer",
    ),
    "buildozer_linux_sha256": attr.string(
        default = "be07a37307759c68696c989058b3446390dd6e8aa6fdca6f44f04ae3c37212c5",
    ),
    "buildozer_mac_url": attr.string(
        default = "https://github.com/bazelbuild/buildtools/releases/download/0.15.0/buildozer.osx",
    ),
    "buildozer_mac_sha256": attr.string(
        default = "294357ff92e7bb36c62f964ecb90e935312671f5a41a7a9f2d77d8d0d4bd217d",
    ),
}

buildozer_http_archive = repository_rule(
    implementation = _http_archive_impl,
    attrs = _http_archive_attrs,
)
"""
http_archive implementation that applies buildozer and sed replacements in the
downloaded archive.

Refer to documentation of the typical the http_archive rule in http.bzl.  This
rule lacks the patch functionality of the original.

Following download and extraction of the archive, this rule will:

1. Execute a single buildozer command.
2. Execute a list of sed commands. 

The buildozer command is constructed from the `replace_deps` and `label_list`
attributes.  For each A -> B mapping in the replace_deps dict, a command like
'replace deps A B' will be appended.  The list of labels to match are taken from
the label_list attribute.  Refer to buildozer documentation for an explanation
of the replace deps command.

The sed commands are constructed from the `sed_replacements` attribute.  These
sed commands might not be necessary if buildozer was capable of replacement
within *.bzl files, but currently it cannot.  This attribute is a
string_list_dict, meaning the dict keys are filename to modify (in place), and
each dict value is are list of sed commands to apply onto that file.  The value
typically looks something like 's|A|B|g'. 
"""
