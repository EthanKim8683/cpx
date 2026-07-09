import os
import sys
import urllib.request
import subprocess
from pathlib import Path

if True:  # TODO: Set to False only after adapting this script to your environment.
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

# Constants (DO NOT MODIFY)
GO = os.environ.get("GO", "go")
PKG_DIR = Path(__file__).resolve().parent.parent
TMP_DIR = PKG_DIR / "tmp"

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

for file in TD_FILES:
    src_url = f"{BASE_URL}/{file}"
    dest_path = TMP_DIR / file
    dest_path.parent.mkdir(parents=True, exist_ok=True)
    print(f"Downloading {file}...")
    urllib.request.urlretrieve(src_url, dest_path)

print("Generating options.json dump...")
subprocess.run(
    [
        TBLGEN,
        "-I",
        str(TMP_DIR / "llvm" / "include"),
        "-I",
        str(TMP_DIR / "clang" / "include"),
        "--dump-json",
        str(TMP_DIR / OPTIONS_TD_FILE),
        "-o",
        str(TMP_DIR / "options.json"),
    ],
    check=True,
)

print("Generating configuration...")
subprocess.run(
    [
        GO,
        "run",
        str(PKG_DIR / "cmd" / "cdbconfiggen"),
        "-o",
        str(PKG_DIR / "generated_cdbconfig.go"),
        str(TMP_DIR / "options.json"),
    ],
    check=True,
)

print("Bootstrap complete")
