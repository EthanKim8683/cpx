#!/usr/bin/env bash
set -euo pipefail

# If you do not have GCC and do not plan on installing it, set OPT_FILES=()
# to bypass downloading and generate an empty option configuration.

# TODO: Remove this line only after adapting this script to your environment
echo "Read internal/gcc/scripts/bootstrap.sh before continuing." && exit 1

# Constants (DO NOT CHANGE)
TMP_DIR="tmp"
OUTPUT_FILE="generated_cdbconfig.go"

# TODO: Configure these variables to match the local environment
# Verify that this URL matches the compiler's version (e.g. releases/gcc-14 for GCC 14)
BASE_URL="https://raw.githubusercontent.com/gcc-mirror/gcc/releases/gcc-16/gcc"
# Refer to the upstream gcc/gcc/Makefile.in to find options.cc (options.c for older versions) source dependencies.
# Check the target architecture (via $GCC_PATH -dumpmachine) to include target-specific files (e.g. gcc/config/aarch64/aarch64.opt).
OPT_FILES=(
	"c-family/c.opt"
	"common.opt"
	"params.opt"
	"analyzer/analyzer.opt"
)

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