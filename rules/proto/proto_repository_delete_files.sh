#!/bin/bash
set -euo pipefail

if [ $# -lt 2 ]; then
    echo "Usage: $0 <base-path> <pattern1> [pattern2] [pattern3] ..."
    echo "Examples:"
    echo "  $0 . 'unittest*.proto'"
    echo "  $0 /src 'unittest*.proto' '*test*.proto' 'map_*unittest.proto'"
    exit 1
fi

BASE_PATH="$1"
shift

# Collect all files matching the glob patterns
TEMP_FILE=$(mktemp)
for pattern in "$@"; do
    find "$BASE_PATH" -path "*$pattern" -type f >> "$TEMP_FILE"
done

# Remove duplicates
sort -u "$TEMP_FILE" > "${TEMP_FILE}.sorted"
mv "${TEMP_FILE}.sorted" "$TEMP_FILE"

if [ ! -s "$TEMP_FILE" ]; then
    echo "No files found matching the patterns."
    rm -f "$TEMP_FILE"
    exit 0
fi

# Delete the files
xargs rm -v < "$TEMP_FILE"
rm -f "$TEMP_FILE"

exit 0