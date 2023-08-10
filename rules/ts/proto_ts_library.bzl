"proto_ts_library.bzl provides the proto_ts_library rule"

load("@aspect_rules_ts//ts:defs.bzl", "ts_project")

def proto_ts_library(
        name,
        declaration = True,
        tsconfig = None,
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
        tsconfig = tsconfig or {
            "compilerOptions": {
                "declaration": declaration,
            },
        },
        **kwargs
    )
