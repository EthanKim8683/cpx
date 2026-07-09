#!/usr/bin/env python3
import os
import sys
import urllib.request
import subprocess
from pathlib import Path

# TODO: Remove this block only after adapting this script to your environment.
print("Read internal/clang/scripts/bootstrap.py before continuing.")
sys.exit(1)

# If you do not have TableGen (or Clang) installed and do not plan to install
# them, replace the download and TableGen compilation steps below with:
#
#   - If Clang is installed but TableGen is missing:
#     Download a pre-compiled options.json dump online and save it to
#     tmp/options.json.
#
#   - If neither Clang nor TableGen are installed:
#     echo "{}" > tmp/options.json

# Constants (DO NOT CHANGE)
TMP_DIR = Path("./internal/clang/tmp")
OUTPUT_FILE = Path("./internal/clang/generated_cdbconfig.go")

# TODO: Configure the variables below to match your local environment.

# Verify that this URL matches your driver's version (e.g. release/17.x for Clang 17).
# Note: For Apple Clang, target Apple's Swift fork (e.g. https://raw.githubusercontent.com/swiftlang/llvm-project/).
BASE_URL = "https://raw.githubusercontent.com/llvm/llvm-project/release/17.x"

# Path to your TableGen executable.
TBLGEN = "/opt/homebrew/opt/llvm/bin/clang-tblgen"

# Locate Options.td within clang/.
# For LLVM ^18.x: clang/include/clang/Options/Options.td
# For LLVM <18.x: clang/include/clang/Driver/Options.td
OPTIONS_TD_FILE = "clang/include/clang/Driver/Options.td"

# Locate dependencies by recursively following 'include "..."' directives in .td files,
# resolving them using clang/include and llvm/include as search paths.
TD_FILES = [
    OPTIONS_TD_FILE,
    "llvm/include/llvm/Option/OptParser.td",
]

# Clang is configured in the .env file at the repository root and loaded via direnv
clang_bin = os.environ.get("CLANG", "")
clang_path = clang_bin if clang_bin else "unset"
clang_version = "unset"

if clang_bin:
    try:
        result = subprocess.run([clang_bin, "--version"], capture_output=True, text=True, check=True)
        clang_version = result.stdout.splitlines()[0]
    except Exception:
        clang_version = "unknown"

tblgen_version = "unset"
if TBLGEN:
    try:
        result = subprocess.run([TBLGEN, "--version"], capture_output=True, text=True, check=True)
        tblgen_version = result.stdout.splitlines()[0]
    except Exception:
        tblgen_version = "unknown"

print(f"Clang path:       {clang_path}")
print(f"Clang version:    {clang_version}")
print(f"TableGen path:    {TBLGEN}")
print(f"TableGen version: {tblgen_version}")
print(f"Upstream URL:     {BASE_URL}")
print(f"TableGen files:   {' '.join(TD_FILES)}")

TMP_DIR.mkdir(parents=True, exist_ok=True)

# Download option files
for file in TD_FILES:
    src_url = f"{BASE_URL}/{file}"
    dest_path = TMP_DIR / file
    dest_path.parent.mkdir(parents=True, exist_ok=True)
    print(f"Downloading {file}...")
    try:
        urllib.request.urlretrieve(src_url, dest_path)
    except Exception as e:
        print(f"Error downloading {file}: {e}", file=sys.stderr)
        sys.exit(1)

print("Generating options.json dump...")
try:
    subprocess.run([
        TBLGEN,
        "-I", str(TMP_DIR / "llvm/include"),
        "-I", str(TMP_DIR / "clang/include"),
        "--dump-json",
        str(TMP_DIR / OPTIONS_TD_FILE),
        "-o", str(TMP_DIR / "options.json")
    ], check=True)
except subprocess.CalledProcessError as e:
    print(f"Error running TableGen: {e}", file=sys.stderr)
    sys.exit(1)

print("Generating configuration...")
try:
    subprocess.run([
        "go", "run", "./internal/clang/cmd/cdbconfiggen",
        "-o", str(OUTPUT_FILE),
        str(TMP_DIR / "options.json")
    ], check=True)
except subprocess.CalledProcessError as e:
    print(f"Error running cdbconfiggen: {e}", file=sys.stderr)
    sys.exit(1)

print("Bootstrap complete")
