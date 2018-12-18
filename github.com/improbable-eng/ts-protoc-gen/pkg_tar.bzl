# Copyright 2015 The Bazel Authors. All rights reserved.
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
"""Rules for manipulation of various packaging."""

load("@bazel_tools//tools/build_defs/pkg:path.bzl", "compute_data_path", "dest_path")

# Filetype to restrict inputs
tar_filetype = [".tar", ".tar.gz", ".tgz", ".tar.xz", ".tar.bz2"]

def _pkg_tar_impl(ctx):
    """Implementation of the pkg_tar rule."""

    # Compute the relative path
    data_path = compute_data_path(ctx.outputs.out, ctx.attr.strip_prefix)

    build_tar = ctx.executable.build_tar
    args = [
        "--output=" + ctx.outputs.out.path,
        "--directory=" + ctx.attr.package_dir,
        "--mode=" + ctx.attr.mode,
        "--owner=" + ctx.attr.owner,
        "--owner_name=" + ctx.attr.ownername,
    ]

    file_inputs = ctx.files.srcs[:]

    # Add runfiles if requested
    if ctx.attr.include_runfiles:
        for f in ctx.attr.srcs:
            if hasattr(f, "default_runfiles"):
                run_files = f.default_runfiles.files.to_list()
                file_inputs += run_files

    # Fixes problem with external repository runfiles in the original pkg_tar
    # implementation.
    for f in file_inputs:
        src = f.path
        dst = dest_path(f, data_path)
        if dst.startswith("../"):
            dst = dst[3:]
        args.append("--file=%s=%s" % (src, dst))

    for target, f_dest_path in ctx.attr.files.items():
        target_files = target.files.to_list()
        if len(target_files) != 1:
            fail("Inputs to pkg_tar.files_map must describe exactly one file.")
        file_inputs += [target_files[0]]
        args += ["--file=%s=%s" % (target_files[0].path, f_dest_path)]
    if ctx.attr.modes:
        args += ["--modes=%s=%s" % (key, ctx.attr.modes[key]) for key in ctx.attr.modes]
    if ctx.attr.owners:
        args += ["--owners=%s=%s" % (key, ctx.attr.owners[key]) for key in ctx.attr.owners]
    if ctx.attr.ownernames:
        args += [
            "--owner_names=%s=%s" % (key, ctx.attr.ownernames[key])
            for key in ctx.attr.ownernames
        ]
    if ctx.attr.empty_files:
        args += ["--empty_file=%s" % empty_file for empty_file in ctx.attr.empty_files]
    if ctx.attr.empty_dirs:
        args += ["--empty_dir=%s" % empty_dir for empty_dir in ctx.attr.empty_dirs]
    if ctx.attr.extension:
        dotPos = ctx.attr.extension.find(".")
        if dotPos > 0:
            dotPos += 1
            args += ["--compression=%s" % ctx.attr.extension[dotPos:]]
    args += ["--tar=" + f.path for f in ctx.files.deps]
    args += [
        "--link=%s:%s" % (k, ctx.attr.symlinks[k])
        for k in ctx.attr.symlinks
    ]
    arg_file = ctx.actions.declare_file(ctx.label.name + ".args")
    ctx.actions.write(arg_file, "\n".join(args))

    ctx.actions.run_shell(
        command = "%s --flagfile=%s" % (build_tar.path, arg_file.path),
        inputs = file_inputs + ctx.files.deps + [arg_file],
        tools = [build_tar],
        outputs = [ctx.outputs.out],
        mnemonic = "PackageTar",
        use_default_shell_env = True,
    )

# A rule for creating a tar file, see README.md
_real_pkg_tar = rule(
    implementation = _pkg_tar_impl,
    attrs = {
        "strip_prefix": attr.string(),
        "package_dir": attr.string(default = "/"),
        "deps": attr.label_list(allow_files = tar_filetype),
        "srcs": attr.label_list(allow_files = True),
        "files": attr.label_keyed_string_dict(allow_files = True),
        "mode": attr.string(default = "0555"),
        "modes": attr.string_dict(),
        "owner": attr.string(default = "0.0"),
        "ownername": attr.string(default = "."),
        "owners": attr.string_dict(),
        "ownernames": attr.string_dict(),
        "extension": attr.string(default = "tar"),
        "symlinks": attr.string_dict(),
        "empty_files": attr.string_list(),
        "include_runfiles": attr.bool(default = False, mandatory = False),
        "empty_dirs": attr.string_list(),
        # Implicit dependencies.
        "build_tar": attr.label(
            default = Label("@bazel_tools//tools/build_defs/pkg:build_tar"),
            cfg = "host",
            executable = True,
            allow_files = True,
        ),
    },
    outputs = {
        "out": "%{name}.%{extension}",
    },
    executable = False,
)

def pkg_tar(**kwargs):
    # Compatibility with older versions of pkg_tar that define files as
    # a flat list of labels.
    if "srcs" not in kwargs:
        if "files" in kwargs:
            if not hasattr(kwargs["files"], "items"):
                label = "%s//%s:%s" % (native.repository_name(), native.package_name(), kwargs["name"])
                print("%s: you provided a non dictionary to the pkg_tar `files` attribute. " % (label,) +
                      "This attribute was renamed to `srcs`. " +
                      "Consider renaming it in your BUILD file.")
                kwargs["srcs"] = kwargs.pop("files")
    _real_pkg_tar(**kwargs)
