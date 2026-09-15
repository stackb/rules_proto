#!/usr/bin/env python3
"""Exercise proto_repository invalidation with real Bazel and Gazelle, offline.

Pass --gazelle (the built //cmd/gazelle binary), --gazelle-repo (the fetched
bazel-gazelle source directory), and optionally --bazel. All generated files
and Bazel state live in a temporary directory.
"""

import argparse
import hashlib
import io
import pathlib
import shutil
import subprocess
import sys
import tarfile
import tempfile
import textwrap


def write(path: pathlib.Path, content: str, executable: bool = False) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content)
    if executable:
        path.chmod(0o755)


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--bazel", default="bazel")
    parser.add_argument("--gazelle", type=pathlib.Path, required=True)
    parser.add_argument("--gazelle-repo", type=pathlib.Path, required=True)
    args = parser.parse_args()
    root = pathlib.Path(__file__).resolve().parents[1]
    gazelle = args.gazelle.resolve()
    with tempfile.TemporaryDirectory(prefix="proto-repository-inputs-") as tmp:
        base = pathlib.Path(tmp)
        workspace = base / "workspace"
        workspace.mkdir()
        write(workspace / "WORKSPACE", "")
        write(workspace / "BUILD.bazel", 'exports_files(["config.yaml", "seed.csv"])\n')
        shutil.copy(root / "rules/proto/proto_repository.bzl", workspace / "repo.bzl")
        module = ['module(name = "test")']
        for name in [
            "bazel_gazelle",
            "bazel_gazelle_go_repository_cache",
            "bazel_gazelle_go_repository_tools",
            "proto_repository_tools",
            # Bazel's built-in module declares these even for a fetch-only test.
            "rules_license",
            "buildozer",
            "platforms",
            "zlib",
            "rules_proto",
            "bazel_features",
            "protobuf",
            "rules_java",
            "rules_cc",
            "rules_python",
            "rules_shell",
            "apple_support",
        ]:
            directory = base / name
            write(directory / "MODULE.bazel", f'module(name = "{name}")\n')
            write(directory / "BUILD.bazel", 'exports_files(glob(["**"]))\n')
            module += [
                f'bazel_dep(name = "{name}")',
                f'local_path_override(module_name = "{name}", path = "{directory}")',
            ]
        for name in ["common.bzl", "go_repository_cache.bzl"]:
            dest = base / "bazel_gazelle/internal" / name
            dest.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy(args.gazelle_repo / "internal" / name, dest)
        write(
            base / "bazel_gazelle/internal/BUILD.bazel",
            'exports_files(glob(["*.bzl"]))\n',
        )
        write(base / "bazel_gazelle_go_repository_cache/go.env", "")
        # The archive is local and already extracted; fetch_repo has no work.
        write(
            base / "bazel_gazelle_go_repository_tools/bin/fetch_repo",
            "#!/bin/sh\nexit 0\n",
            True,
        )
        # A changing comment makes upstream regeneration observable downstream,
        # even when a config edit leaves Gazelle's semantic index unchanged.
        wrapper = f"#!{sys.executable}\n" + textwrap.dedent(f"""\
            import pathlib
            import subprocess
            import sys
            import uuid
            subprocess.run([{str(gazelle)!r}, *sys.argv[1:]], check=True)
            with pathlib.Path('imports.csv').open('a') as out:
                out.write('# generation ' + str(uuid.uuid4()) + '\\n')
            """)
        write(base / "proto_repository_tools/bin/gazelle", wrapper, True)
        archive = base / "source.tar.gz"
        with tarfile.open(archive, "w:gz") as tar:
            for name, content in {
                "BUILD.bazel": 'exports_files(["imports.csv"])\n',
                "REPO.bazel": "",
                "WORKSPACE": "",
                "example.proto": 'syntax = "proto3"; package example; message Example {}\n',
                # Also exercise a Starlark file inside the fetched repository:
                # it must not be passed to ctx.watch, which forbids own paths.
                "internal.star": 'protoc.Plugin(name = "internal", configure = lambda ctx: None)\n',
            }.items():
                data = content.encode()
                info = tarfile.TarInfo(name)
                info.size = len(data)
                tar.addfile(info, io.BytesIO(data))
        digest = hashlib.sha256(archive.read_bytes()).hexdigest()
        module.append(
            'proto_repository = use_repo_rule("//:repo.bzl", "protobuf_go_repository")'
        )
        common = f'urls = ["{archive.as_uri()}"], sha256 = "{digest}", build_config = "//:WORKSPACE", build_file_generation = "on"'
        module += [
            f'proto_repository(name = "producer", apparent_name = "producer", cfgs = ["//:config.yaml"], imports = ["//:seed.csv"], {common})',
            f'proto_repository(name = "consumer", apparent_name = "consumer", imports = ["@producer//:imports.csv"], {common})',
        ]
        write(workspace / "MODULE.bazel", "\n".join(module) + "\n")
        config = workspace / "config.yaml"
        write(
            config,
            "starlarkPlugins:\n  - plugin.star%test\n  - internal.star%internal\n",
        )
        plugin = workspace / "plugin.star"
        write(plugin, 'protoc.Plugin(name = "test", configure = lambda ctx: None)\n')
        seed = workspace / "seed.csv"
        write(seed, "# seed\n")
        output = base / "output"
        # Match the source-root marker used by Gazelle in external repositories.
        write(output / "DO_NOT_BUILD_HERE", str(workspace))
        command = [
            args.bazel,
            "--batch",
            f"--output_user_root={base / 'bazel'}",
            f"--output_base={output}",
            "--ignore_all_rc_files",
            "fetch",
            "--repo=@consumer",
            "--incompatible_autoload_externally=",
            "--lockfile_mode=off",
            "--noshow_progress",
        ]

        def fetch() -> tuple[str, str]:
            result = subprocess.run(
                command, cwd=workspace, capture_output=True, text=True, timeout=60
            )
            if result.returncode:
                raise RuntimeError(result.stdout + result.stderr)
            indexes = []
            for name in ["producer", "consumer"]:
                files = list((output / "external").glob(f"*+{name}/imports.csv"))
                if len(files) != 1:
                    raise AssertionError(f"expected one {name} index, found {files}")
                indexes.append(files[0].read_text())
            return indexes[0], indexes[1]

        previous = fetch()
        assert fetch() == previous, "unchanged repositories should be reused"
        for path in [config, seed, plugin]:
            with path.open("a") as out:
                out.write("\n# changed\n")
            current = fetch()
            assert all(a != b for a, b in zip(previous, current)), (
                f"{path.name} did not invalidate both repositories"
            )
            assert fetch() == current, f"{path.name} caused unnecessary regeneration"
            previous = current
            print(
                f"PASS: {path.name} invalidates producer and consumer; unchanged rerun is stable",
                flush=True,
            )

        for name in ["first.star", "second.star"]:
            target = workspace / name
            write(target, plugin.read_text() + f"\n# {name}\n")
            plugin.unlink()
            plugin.symlink_to(target)
            current = fetch()
            assert all(a != b for a, b in zip(previous, current)), (
                "Starlark symlink change did not invalidate both repositories"
            )
            assert fetch() == current, (
                "Starlark symlink caused unnecessary regeneration"
            )
            previous = current
        print(
            "PASS: Starlark symlink replacement and retargeting invalidate both repositories",
            flush=True,
        )


if __name__ == "__main__":
    main()
