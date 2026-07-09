#!/usr/bin/env bash
set -euo pipefail

# If you do not have GCC and do not plan on installing it, set OPT_FILES=()
# to bypass downloading and generate an empty option configuration.

# Constants (DO NOT CHANGE)
TMP_DIR="tmp"
OUTPUT_FILE="generated_cdbconfig.go"

# TODO: Configure these variables to match the local environment
# Verify that this URL matches the compiler's version (e.g. releases/gcc-14 for GCC 14)
BASE_URL="https://raw.githubusercontent.com/gcc-mirror/gcc/releases/gcc-16/gcc"
# Refer to the upstream GCC Makefile to find options.cc (options.c for older versions) source dependencies.
# You can check the target architecture (via $GCC_PATH -dumpmachine) to determine if target-specific
# option files are needed (e.g. config/aarch64/aarch64.opt, config/i386/i386.opt, or config/darwin.opt).
# Note: Option files are version-dependent (e.g. analyzer/analyzer.opt was added in GCC 10; params.opt in GCC 5).
OPT_FILES=(
	"c-family/c.opt"
	"common.opt"
	"params.opt"
	"analyzer/analyzer.opt"
)

# TODO: If unset, set GCC in .env at the root of the project
GCC_PATH="unset"
GCC_VERSION="unset"
if [ -n "${GCC:-}" ]; then
	GCC_PATH="$GCC"
	GCC_VERSION=$("$GCC_PATH" --version | head -n 1)
fi

# If the path shows "unset" but GCC is installed, set the GCC path in the .env file at the repository root
echo "GCC path:     $GCC_PATH"
echo "GCC version:  $GCC_VERSION"
echo "Upstream URL: $BASE_URL"
echo "Option files: ${OPT_FILES[*]}"

# TODO: Remove this line after configuring the variables
echo "Read internal/gcc/scripts/bootstrap.sh before continuing." && exit 1

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