This directory contains a script to generate the nuget.bzl file, which declares
the nuget protobuf and grpc dependencies.

After running this script one must manually update the grpc.core `nuget_package`
rule `core_files` attribute with the runtime libraries.  

To inspect these files, `(cd $(bazel info output_base)/external/grpc.core &&
find runtimes/)` and copy those file paths into the attribute.