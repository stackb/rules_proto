#!/bin/bash

#
# Script is written to be run from the WORKSPACE root
#
set -eu
set -o pipefail

PROTOBUF_VERSION="3.6.1"
GRPC_VERSION="1.15.0"

TOOL="./bazel-bin/external/io_bazel_rules_dotnet/tools/nuget2bazel/nuget2bazel"

# Compile the tool, cannot use sandboxing
bazel build --spawn_strategy=standalone @io_bazel_rules_dotnet//tools/nuget2bazel

# Make temp directory
TMPDIR=$(mktemp -d ${TMPDIR:-/tmp}/nuget.XXXXXXXXXX)

# Tool only works on WORKSPACE files so we have to do these shenanigans
touch "${TMPDIR}/WORKSPACE"

# Gather protobuf deps
"${TOOL}" add --path "/${TMPDIR}" Google.Protobuf "${PROTOBUF_VERSION}"

# Create the nuget.bzl file, indent 4 spaces
echo 'load("@io_bazel_rules_dotnet//dotnet:defs.bzl", "nuget_package")' > ./csharp/nuget/nuget.bzl
echo "def nuget_protobuf_packages():" >> csharp/nuget/nuget.bzl
sed -e 's/^/    /' "${TMPDIR}/WORKSPACE" >> csharp/nuget/nuget.bzl

# Reset the workspace file
echo "" > "${TMPDIR}/WORKSPACE"

# Rerun it for grpc deps
"${TOOL}" add --path "/${TMPDIR}" Grpc "${GRPC_VERSION}"

# Similarly, write the gprc.bzl file
echo "def nuget_grpc_packages():" >> ./csharp/nuget/nuget.bzl
sed -e 's/^/    /' "${TMPDIR}/WORKSPACE" >> ./csharp/nuget/nuget.bzl

# Cleanup
rm -rf "${TMPDIR}"
