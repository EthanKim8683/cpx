#!/usr/bin/env bash
set -euo pipefail

# Do not change these variables
TMP_DIR="tmp"
OUTPUT_FILE="generated_cdbconfig.go"

# Change these variables to match the local environment
# Verify that this URL matches the compiler's version
BASE_URL="https://raw.githubusercontent.com/gcc-mirror/gcc/releases/gcc-16/gcc"
# Refer to the upstream GCC Makefile for option source dependencies
OPT_FILES=(
	"c-family/c.opt"
	"common.opt"
	"params.opt"
	"analyzer/analyzer.opt"
)

# Display environment configurations for verification
GCC_PATH="unset"
GCC_VERSION="unset"
if [ -n "${GCC:-}" ]; then
	GCC_PATH="$GCC"
	GCC_VERSION=$("$GCC_PATH" --version | head -n 1)
fi

echo "GCC path:     $GCC_PATH"
echo "GCC version:  $GCC_VERSION"
echo "Upstream URL: $BASE_URL"
echo "Option files: ${OPT_FILES[*]}"

echo "Please verify the configuration above and delete this safety check line to continue bootstrapping." && exit 1

mkdir -p "$TMP_DIR"
for file in "${OPT_FILES[@]}"; do
	src="${BASE_URL}/${file}"
	dest="${TMP_DIR}/${file}"
	echo "Downloading $file..."
	mkdir -p "$(dirname "$dest")"
	curl -fsSL "$src" -o "$dest"
done

echo "Generating configuration..."
go run ./cmd/cdbconfiggen -o "$OUTPUT_FILE" "${OPT_FILES[@]/#/$TMP_DIR/}"