#!/usr/bin/env bash

#
# this test only asserts that the generated bundle file exists.
#

set -euo pipefail

# find .

readonly bundle_js="./example/routeguide/closure/bundle.js"

if [[ ! -f "${bundle_js}" ]]; then
    echo "missing file: ${bundle_js}"
fi
