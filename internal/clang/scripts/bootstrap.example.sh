#!/usr/bin/env bash
set -euo pipefail

# If you do not have TableGen (or Clang) installed and do not plan to install
# them, replace the download and TableGen compilation steps below with:
#
#   - If Clang is installed but TableGen is missing:
#     Download a pre-compiled options.json dump online and save it to
#     tmp/options.json.
#
#   - If neither Clang nor TableGen are installed:
#     echo "{}" > tmp/options.json

# TODO: Remove this line only after adapting this script to your environment
echo "Read internal/clang/scripts/bootstrap.sh before continuing." && exit 1

# Constants (DO NOT CHANGE)
TMP_DIR="tmp"
OUTPUT_FILE="generated_cdbconfig.go"

# TODO: Configure all the variables below to match your local environment.

# Verify that this URL matches the compiler's version (e.g. release/17.x for Clang 17).
# Note: For Apple Clang, target Apple's Swift fork (e.g. https://raw.githubusercontent.com/swiftlang/llvm-project/).
BASE_URL="https://raw.githubusercontent.com/llvm/llvm-project/release/17.x"

# Path to the TableGen executable.
TBLGEN="/opt/homebrew/opt/llvm/bin/clang-tblgen"

# Locate Options.td within clang/ (release/^18.x uses clang/include/clang/Options/Options.td;
# release/<18.x uses clang/include/clang/Driver/Options.td).
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
"$TBLGEN" \
	-I "${TMP_DIR}/llvm/include" \
	-I "${TMP_DIR}/clang/include" \
	--dump-json \
	"${TMP_DIR}/${OPTIONS_TD_FILE}" \
	-o "${TMP_DIR}/options.json"

echo "Generating configuration..."
go run ./cmd/cdbconfiggen -o "$OUTPUT_FILE" "${TMP_DIR}/options.json"
