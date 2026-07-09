#!/usr/bin/env bash
set -euo pipefail

# If you are missing dependencies and do not plan on installing them:
#
# 1. Has Clang but no TableGen:
#    Since you are using Clang, you should look into obtaining a pre-compiled options.json dump
#    online for your version and placing it under tmp/options.json.
#
# 2. No Clang and no TableGen:
#    If you do not use Clang at all, you can generate an empty option configuration by running:
#    echo "{}" > tmp/options.json

# Constants (DO NOT CHANGE)
TMP_DIR="tmp"
OUTPUT_FILE="generated_cdbconfig.go"

# TODO: Configure these variables to match the local environment
# Verify that this URL matches the compiler's version (e.g. release/17.x for Clang 17)
# Note: For Apple Clang, target Apple's Swift fork (e.g. https://raw.githubusercontent.com/swiftlang/llvm-project/swift-5.10-RELEASE)
BASE_URL="https://raw.githubusercontent.com/llvm/llvm-project/release/17.x"
TBLGEN="${TBLGEN:-clang-tblgen}"

# Refer to Clang's Options.td and diagnostic TableGen files for dependencies.
# Note: TableGen files are version-dependent.
TD_FILES=(
  "clang/include/clang/Driver/Options.td"
  "llvm/include/llvm/Option/OptParser.td"
)

CLANG_PATH="unset"
CLANG_VERSION="unset"
if [ -n "${CLANG:-}" ]; then
	CLANG_PATH="$CLANG"
	CLANG_VERSION=$("$CLANG_PATH" --version | head -n 1)
fi

TBLGEN_PATH="${TBLGEN:-clang-tblgen}"
TBLGEN_VERSION="unset"
if command -v "$TBLGEN_PATH" >/dev/null 2>&1; then
	TBLGEN_VERSION=$("$TBLGEN_PATH" --version | head -n 1)
fi

# If the path shows "unset" but Clang is installed, set the CLANG path in the .env file at the repository root
echo "Clang path:       $CLANG_PATH"
echo "Clang version:    $CLANG_VERSION"
echo "TableGen path:    $TBLGEN_PATH"
echo "TableGen version: $TBLGEN_VERSION"
echo "Upstream URL:     $BASE_URL"
echo "TableGen files:   ${TD_FILES[*]}"

# TODO: Remove this line after configuring the variables
echo "Read internal/clang/scripts/bootstrap.sh before continuing." && exit 1

mkdir -p "$TMP_DIR"
for file in "${TD_FILES[@]}"; do
	src="${BASE_URL}/${file}"
	dest="${TMP_DIR}/${file}"
	echo "Downloading $file..."
	mkdir -p "$(dirname "$dest")"
	curl -fsSL "$src" -o "$dest"
done

echo "Generating options.json dump..."
# TODO: Adjust these TableGen include paths (-I) and compile target file if your compiler fork
# or targeted version uses a different options directory layout.
"$TBLGEN_PATH" \
	-I "${TMP_DIR}/llvm/include" \
	-I "${TMP_DIR}/clang/include" \
	--dump-json \
	"${TMP_DIR}/clang/include/clang/Driver/Options.td" \
	-o "${TMP_DIR}/options.json"

echo "Generating configuration..."
go run ./cmd/cdbconfiggen -o "$OUTPUT_FILE" "${TMP_DIR}/options.json"
