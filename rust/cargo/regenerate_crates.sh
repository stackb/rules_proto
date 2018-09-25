#!/bin/bash
set -eu
set -o pipefail

if ! (which cargo &>/dev/null); then
    echo "Cannot find cargo in PATH, please install Cargo." >&2
    exit 1
fi

if ! (which cargo-raze &>/dev/null); then
    echo "Cannot find cargo-raze in PATH, please install cargo-raze (cargo install cargo-raze)." >&2
    exit 1
fi

cd "$(dirname "${BASH_SOURCE[0]}")"

rm -fr remote
# Now run cargo generate-lockfile && cargo raze
cargo generate-lockfile && cargo raze

# We need to update the build file path.
#sed -i.bak 's|//rust|@org_pubref_rules_proto//rust|g' "crates.bzl"
#rm crates.bzl.bak

# Remove these outdated rust_bench_test references
find remote -name '*BUILD' | xargs sed -i.bak '/rust_bench_test/d'