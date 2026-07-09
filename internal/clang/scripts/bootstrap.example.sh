#!/usr/bin/env bash
set -euo pipefail

# Do not change these variables
TMP_DIR="tmp"
OUTPUT_FILE="generated_cdbconfig.go"

# Change these variables to match the local environment
# Verify that this URL matches the compiler's version
BASE_URL="https://raw.githubusercontent.com/llvm/llvm-project/release/17.x"
TBLGEN="${TBLGEN:-clang-tblgen}"

# Refer to Clang's Driver CMakeLists.txt and Options.td for dependencies
TD_FILES=(
	"clang/include/clang/Driver/Options.td"
	"clang/include/clang/Driver/OptionDocEmitter.td"
	"llvm/include/llvm/Option/OptParser.td"
	"clang/include/clang/Basic/DiagnosticOptions.td"
	"clang/include/clang/Basic/DiagnosticGroups.td"
)

# Display environment configurations for verification
CLANG_PATH="${CLANG:-unset}"
CLANG_VERSION="unknown"
if [ "$CLANG_PATH" != "unset" ] && command -v "$CLANG_PATH" >/dev/null 2>&1; then
	CLANG_VERSION=$("$CLANG_PATH" --version | head -n 1)
fi

TBLGEN_PATH="${TBLGEN}"
TBLGEN_VERSION="unknown"
if command -v "$TBLGEN_PATH" >/dev/null 2>&1; then
	TBLGEN_VERSION=$("$TBLGEN_PATH" --version | head -n 1)
fi

echo "Clang path:       $CLANG_PATH"
echo "Clang version:    $CLANG_VERSION"
echo "TableGen path:    $TBLGEN_PATH"
echo "TableGen version: $TBLGEN_VERSION"
echo "Upstream URL:     $BASE_URL"
echo "TableGen files:   ${TD_FILES[*]}"

echo "Configure the variables above and delete this safety check to run." && exit 1

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
	"${TMP_DIR}/clang/include/clang/Driver/Options.td" \
	-o "${TMP_DIR}/options.json"

echo "Generating configuration..."
go run ./cmd/cdbconfiggen -o "$OUTPUT_FILE" "${TMP_DIR}/options.json"
