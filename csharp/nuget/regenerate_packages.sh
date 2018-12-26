#!/bin/bash

#
# Script is written to be run from the WORKSPACE root
#
set -eu
set -o pipefail

PROTOBUF_VERSION="3.6.1"
GRPC_VERSION="1.17.1"

OUTPUT_FILE="./csharp/nuget/nuget.bzl"

# Compile the tool
bazel build @io_bazel_rules_dotnet//tools/nuget2bazel

TOOL="$(pwd)/bazel-bin/external/io_bazel_rules_dotnet/tools/nuget2bazel/nuget2bazel.bash"

# Make temp directory
TMPDIR=$(mktemp -d ${TMPDIR:-/tmp}/nuget.XXXXXXXXXX)

# Tool only works on WORKSPACE files so we have to do these shenanigans
touch "${TMPDIR}/WORKSPACE"

# Gather protobuf deps
"${TOOL}" add --skipSha256 --path "/${TMPDIR}" Google.Protobuf "${PROTOBUF_VERSION}"

# Create the nuget.bzl file, indent 4 spaces
echo 'load("@io_bazel_rules_dotnet//dotnet:defs.bzl", "nuget_package")' > "${OUTPUT_FILE}"
echo "" >> "${OUTPUT_FILE}"
echo "def nuget_protobuf_packages():" >> "${OUTPUT_FILE}"
sed -e 's/^/    /' "${TMPDIR}/WORKSPACE" >> "${OUTPUT_FILE}"

echo "" >> "${OUTPUT_FILE}"
echo "def nuget_grpc_packages():" >> "${OUTPUT_FILE}"

# Reset the workspace file; rerun it for grpc deps
echo "" > "${TMPDIR}/WORKSPACE"
"${TOOL}" add --skipSha256 --path "/${TMPDIR}" Grpc "${GRPC_VERSION}"

# # Similarly, write the gprc deps
sed -e 's/^/    /' "${TMPDIR}/WORKSPACE" >> "${OUTPUT_FILE}"

# Cleanup
rm -rf "${TMPDIR}"
