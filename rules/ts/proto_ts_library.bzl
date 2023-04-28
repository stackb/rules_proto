"proto_ts_library.bzl provides the proto_ts_library rule"

load("@aspect_rules_ts//ts:defs.bzl", "ts_project")

def proto_ts_library(
        name,
        srcs = None,
        tsconfig = None,
        args = [],
        data = [],
        deps = [],
        extends = None,
        source_map = False,
        declaration_map = False,
        resolve_json_module = None,
        composite = False,
        incremental = False,
        emit_declaration_only = False,
        transpiler = None,
        ts_build_info_file = None,
        declaration_dir = None,
        out_dir = None,
        root_dir = None,
        **kwargs):
    """Compiles a Typescript library from generated code.

    This is a thin wrapper around aspect_rules_ts's ts_project macro. Note that the Gazelle rule only supports customizing
    the args and tsconfig attributes.

    Additionally, it is necessary to pass
    ```
    deps = [
            "//:node_modules/long",
            "//:node_modules/protobufjs",
    ],
    ```
    meaning one must call npm_link_all_packages in the BUILD file at the *root* of the workspace (or figure out how to set rootDirs
    appropriately).

    See https://github.com/aspect-build/rules_ts/blob/007d1ab8f168879cad4ba6b741e8c6def20aa262/docs/rules.md#ts_project
    """
    ts_project(
        name = name,
        srcs = srcs,
        tsconfig = tsconfig or {},
        args = args,
        data = data,
        deps = deps,
        extends = extends,
        declaration = True,
        source_map = source_map,
        declaration_map = declaration_map,
        resolve_json_module = resolve_json_module,
        composite = composite,
        incremental = incremental,
        emit_declaration_only = emit_declaration_only,
        transpiler = transpiler,
        ts_build_info_file = ts_build_info_file,
        declaration_dir = declaration_dir,
        out_dir = out_dir,
        root_dir = root_dir,
        **kwargs
    )
