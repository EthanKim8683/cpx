#!/usr/bin/env bash
# bootstrap.sh - Parameterized Clang CDB Option Parser Bootstrapper
# This script fetches Clang TableGen files, compiles them into a JSON dump, and runs cdbconfiggen.
set -euo pipefail

# -- Configuration --
# These values should be adapted for the local environment.
CLANG_BIN="${CLANG:-}"
TBLGEN="${TBLGEN:-clang-tblgen}"
BASE_URL="${BASE_URL:-https://raw.githubusercontent.com/llvm/llvm-project/release/17.x}"

# Fixed paths in the cpx package structure
TMP_DIR="internal/clang/tmp"
OUTPUT_FILE="internal/clang/generated_cdbconfig.go"

echo "=== Clang Option Config Bootstrapper ==="

# Validate that CLANG environment variable is set
if [ -z "$CLANG_BIN" ]; then
    echo "Error: CLANG environment variable is not set." >&2
    echo "Please set it in your .env file (e.g. CLANG=/usr/bin/clang) so direnv loads it automatically." >&2
    exit 1
fi

if ! command -v "$CLANG_BIN" >/dev/null 2>&1; then
    echo "Error: '$CLANG_BIN' is not executable or not in PATH." >&2
    exit 1
fi

if ! command -v "$TBLGEN" >/dev/null 2>&1; then
    echo "Error: '$TBLGEN' is not executable or not in PATH." >&2
    echo "Please install llvm/clang development packages or specify TBLGEN path." >&2
    exit 1
fi

# Display compiler version info
echo "--- Local Compiler Version ---"
"$CLANG_BIN" --version
echo "------------------------------"

# Display tblgen version info
echo "--- Local TableGen Version ---"
"$TBLGEN" --version
echo "------------------------------"

echo "Base URL: $BASE_URL"
echo "======================================"

mkdir -p "$TMP_DIR"

TD_FILES=(
    "clang/include/clang/Driver/Options.td"
    "clang/include/clang/Driver/OptionDocEmitter.td"
    "llvm/include/llvm/Option/OptParser.td"
    "clang/include/clang/Basic/DiagnosticOptions.td"
    "clang/include/clang/Basic/DiagnosticGroups.td"
)

echo "Downloading TableGen source files..."
for file in "${TD_FILES[@]}"; do
    target_path="${TMP_DIR}/${file}"
    mkdir -p "$(dirname "$target_path")"
    
    echo "  - $file"
    curl -fsSL -o "$target_path" "${BASE_URL}/${file}"
done

echo "Generating JSON dump using tblgen..."
"$TBLGEN" \
    -I "${TMP_DIR}/llvm/include" \
    -I "${TMP_DIR}/clang/include" \
    --dump-json \
    "${TMP_DIR}/clang/include/clang/Driver/Options.td" \
    -o "${TMP_DIR}/options.json"

echo "Running cdbconfiggen..."
go run ./internal/clang/cmd/cdbconfiggen -o "$OUTPUT_FILE" "${TMP_DIR}/options.json"

echo "Success! Config generated at $OUTPUT_FILE"
