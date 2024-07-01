"proto_ts_library.bzl provides the proto_ts_library rule"

load("@aspect_rules_ts//ts:defs.bzl", "ts_project")

def proto_ts_library(
        name,
        declaration = True,
        tsconfig = None,
        **kwargs):
    """Compiles a Typescript library from generated code.
    """
    ts_project(
        name = name,
        # validate = False,
        declaration = declaration,
        tsconfig = tsconfig or {
            "compilerOptions": {
                "declaration": declaration,
                "esModuleInterop": True,
            },
        },
        **kwargs
    )
