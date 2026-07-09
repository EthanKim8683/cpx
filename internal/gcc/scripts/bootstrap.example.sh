#!/usr/bin/env bash
# bootstrap.sh - Parameterized GCC CDB Option Parser Bootstrapper
# This script fetches GCC .opt files from upstream and runs cdbconfiggen.
set -euo pipefail

# -- Configuration --
# These values should be adapted for the local environment.
GCC_BIN="${GCC:-}"
BASE_URL="${BASE_URL:-https://raw.githubusercontent.com/gcc-mirror/gcc/releases/gcc-14/gcc}"

# Fixed paths in the cpx package structure
TMP_DIR="internal/gcc/tmp"
OUTPUT_FILE="internal/gcc/generated_cdbconfig.go"

echo "=== GCC Option Config Bootstrapper ==="

# Validate that GCC environment variable is set
if [ -z "$GCC_BIN" ]; then
    echo "Error: GCC environment variable is not set." >&2
    echo "Please set it in your .env file (e.g. GCC=/usr/bin/gcc) so direnv loads it automatically." >&2
    exit 1
fi

if ! command -v "$GCC_BIN" >/dev/null 2>&1; then
    echo "Error: '$GCC_BIN' is not executable or not in PATH." >&2
    exit 1
fi

# Display compiler version info
echo "--- Local Compiler Version ---"
"$GCC_BIN" --version
echo "------------------------------"

echo "Base URL: $BASE_URL"
echo "======================================"

mkdir -p "$TMP_DIR"

# List of .opt files to fetch (adapted for this environment)
OPT_FILES=(
    "common.opt"
    "params.opt"
    "c-family/c.opt"
    "analyzer/analyzer.opt"
)

echo "Downloading option source files..."
for file in "${OPT_FILES[@]}"; do
    target_path="${TMP_DIR}/${file}"
    mkdir -p "$(dirname "$target_path")"
    
    echo "  - $file"
    curl -fsSL -o "$target_path" "${BASE_URL}/${file}"
done

echo "Running cdbconfiggen..."
go run ./internal/gcc/cmd/cdbconfiggen -o "$OUTPUT_FILE" "$TMP_DIR"

echo "Success! Config generated at $OUTPUT_FILE"
