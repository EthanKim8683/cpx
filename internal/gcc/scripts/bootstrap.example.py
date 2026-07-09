import sys
import urllib.request
import subprocess
from pathlib import Path

if True:  # TODO: Set to False only after adapting this script to your environment.
    print("Read internal/gcc/scripts/bootstrap.py before continuing.")
    sys.exit(1)

# If you do not have GCC installed and do not plan to install it, set OPT_FILES=()
# to skip downloading and generate an empty option configuration.

# TODO: Configure the variables below to match your local environment.

# Constants (DO NOT CHANGE)
PKG_DIR = Path(__file__).resolve().parent.parent
TMP_DIR = PKG_DIR / "tmp"

# Verify that this URL matches your compiler's version (e.g. releases/gcc-16 for GCC 16).
BASE_URL = "https://raw.githubusercontent.com/gcc-mirror/gcc/releases/gcc-16/gcc"

# Refer to the upstream Makefile.in to find options.cc (options.c for older versions) source dependencies.
# Check the target architecture (via $GCC -dumpmachine) to include target-specific files (e.g. config/aarch64/aarch64.opt).
OPT_FILES = [
    "c-family/c.opt",
    "common.opt",
    "params.opt",
    "analyzer/analyzer.opt",
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
        "go",
        "run",
        str(PKG_DIR / "cmd" / "cdbconfiggen"),
        "-o",
        str(PKG_DIR / "generated_cdbconfig.go"),
    ]
    + file_args,
    check=True,
)

print("Bootstrap complete")
