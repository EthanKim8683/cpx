import os
import sys
import urllib.request
import subprocess
from pathlib import Path

if True:  # TODO: Set to False only after adapting this script to your environment.
    print("Read internal/gcc/scripts/bootstrap.py before continuing.")
    sys.exit(1)

# If you do not have GCC installed and do not plan to install it, set OPT_FILES=()
# to skip downloading and generate an empty option configuration.

# Constants (DO NOT MODIFY)
GO = os.environ.get("GO", "go")
PKG_DIR = Path(__file__).resolve().parent.parent
TMP_DIR = PKG_DIR / "tmp"

# TODO: Configure the variables below to match your local environment.

# Verify that this URL matches your compiler's version (e.g. releases/gcc-16 for GCC 16).
BASE_URL = "https://raw.githubusercontent.com/gcc-mirror/gcc/releases/gcc-16"

# Refer to the upstream gcc/Makefile.in to find the .opt files options.cc (options.c for older versions) depends on.
# Check the target architecture (via $GCC -dumpmachine) to identify target-specific files (e.g. gcc/config/aarch64/aarch64.opt).
OPT_FILES = [
    "gcc/c-family/c.opt",
    "gcc/common.opt",
    "gcc/params.opt",
    "gcc/analyzer/analyzer.opt",
]

for file in OPT_FILES:
    src_url = f"{BASE_URL}/{file}"
    dest_path = TMP_DIR / file
    dest_path.parent.mkdir(parents=True, exist_ok=True)
    print(f"Downloading {file}...")
    urllib.request.urlretrieve(src_url, dest_path)

print("Generating configuration...")
file_args = [str(TMP_DIR / file) for file in OPT_FILES]
subprocess.run(
    [
        GO,
        "run",
        str(PKG_DIR / "cmd" / "cdbconfiggen"),
        "-o",
        str(PKG_DIR / "generated_cdbconfig.go"),
    ]
    + file_args,
    check=True,
)

print("Bootstrap complete")
