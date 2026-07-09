#!/usr/bin/env bash
set -euo pipefail

# If you are missing dependencies and do not plan on installing them:
#
# 1. Clang is installed, but TableGen is not:
#    Obtain a pre-compiled options.json dump for your version and place it under tmp/options.json and run:
#    go run ./cmd/cdbconfiggen -o "$OUTPUT_FILE" "${TMP_DIR}/options.json"
#
# 2. Neither Clang nor TableGen are installed:
#    You can generate an empty option configuration by running:
#    echo "{}" > ${TMP_DIR}/options.json
#    go run ./cmd/cdbconfiggen -o "$OUTPUT_FILE" "${TMP_DIR}/options.json"

# TODO: Remove this line only after adapting this script to your environment
echo "Read internal/clang/scripts/bootstrap.sh before continuing." && exit 1

# Constants (DO NOT CHANGE)
TMP_DIR="tmp"
OUTPUT_FILE="generated_cdbconfig.go"

# TODO: Configure these variables to match the local environment
# Verify that this URL matches the compiler's version (e.g. release/17.x for Clang 17)
# Note: For Apple Clang, target Apple's Swift fork (e.g. https://raw.githubusercontent.com/swiftlang/llvm-project/)
BASE_URL="https://raw.githubusercontent.com/llvm/llvm-project/release/17.x"
TBLGEN="/opt/homebrew/opt/llvm/bin/clang-tblgen"
# Locate Options.td within clang/ (paths may vary in forks like Apple Clang):
# - LLVM >= 18: clang/include/clang/Options/Options.td
# - LLVM <= 17: clang/include/clang/Driver/Options.td
OPTIONS_TD_FILE="clang/include/clang/Driver/Options.td"
# Recursively locate dependencies by following 'include "..."' directives in Options.td,
# resolving them using clang/include and llvm/include as search paths.
TD_FILES=(
  "$OPTIONS_TD_FILE"
  "llvm/include/llvm/Option/OptParser.td"
)

mkdir -p "$TMP_DIR"
for file in "${TD_FILES[@]}"; do
	src="${BASE_URL}/${file}"
	dest="${TMP_DIR}/${file}"
	echo "Downloading $file..."
	mkdir -p "$(dirname "$dest")"
	curl -fsSL "$src" -o "$dest"
done

echo "Generating options.json dump..."
"$TBLGEN_PATH" \
	-I "${TMP_DIR}/llvm/include" \
	-I "${TMP_DIR}/clang/include" \
	--dump-json \
	"${TMP_DIR}/${OPTIONS_TD_FILE}" \
	-o "${TMP_DIR}/options.json"

echo "Generating configuration..."
go run ./cmd/cdbconfiggen -o "$OUTPUT_FILE" "${TMP_DIR}/options.json"
