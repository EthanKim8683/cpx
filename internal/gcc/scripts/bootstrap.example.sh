#!/usr/bin/env bash
set -euo pipefail

# Do not change these variables
TMP_DIR="tmp"
OUTPUT_FILE="generated_cdbconfig.go"

# Change these variables to match the local environment
# Verify that this URL matches the compiler's version (e.g. releases/gcc-14 for GCC 14)
BASE_URL="https://raw.githubusercontent.com/gcc-mirror/gcc/releases/gcc-16/gcc"
# Refer to the upstream GCC Makefile to find options.cc (options.c for older versions) source dependencies.
# You can check the target architecture (via $GCC_PATH -dumpmachine) to determine if target-specific
# option files are needed (e.g. config/aarch64/aarch64.opt, config/i386/i386.opt, or config/darwin.opt).
OPT_FILES=(
	"c-family/c.opt"
	"common.opt"
	"params.opt"
	"analyzer/analyzer.opt"
)

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

# Safety barrier: edit the configurations above, then delete this line to execute the script
echo "Configure the variables above and delete this safety check to run." && exit 1

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