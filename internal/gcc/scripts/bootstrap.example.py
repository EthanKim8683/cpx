#!/usr/bin/env python3
import os
import sys
import urllib.request
import subprocess
from pathlib import Path

# TODO: Remove this block only after adapting this script to your environment.
print("Read internal/gcc/scripts/bootstrap.py before continuing.")
sys.exit(1)

# If you do not have GCC installed and do not plan to install it, set OPT_FILES=()
# to skip downloading and generate an empty option configuration.

# Constants (DO NOT CHANGE)
TMP_DIR = Path("./internal/gcc/tmp")
OUTPUT_FILE = Path("./internal/gcc/generated_cdbconfig.go")

# TODO: Configure the variables below to match your local environment.

# Verify that this URL matches your compiler's version (e.g. releases/gcc-14 for GCC 14).
BASE_URL = "https://raw.githubusercontent.com/gcc-mirror/gcc/releases/gcc-16/gcc"

# Refer to the upstream gcc/gcc/Makefile.in to find options.cc (options.c for older versions) source dependencies.
# Check the target architecture (via $GCC_PATH -dumpmachine) to include target-specific files (e.g. gcc/config/aarch64/aarch64.opt).
OPT_FILES = [
    "c-family/c.opt",
    "common.opt",
    "params.opt",
    "analyzer/analyzer.opt",
]

TMP_DIR.mkdir(parents=True, exist_ok=True)

# Download option files
for file in OPT_FILES:
    src_url = f"{BASE_URL}/{file}"
    dest_path = TMP_DIR / file
    dest_path.parent.mkdir(parents=True, exist_ok=True)
    print(f"Downloading {file}...")
    try:
        urllib.request.urlretrieve(src_url, dest_path)
    except Exception as e:
        print(f"Error downloading {file}: {e}", file=sys.stderr)
        sys.exit(1)

print("Generating configuration...")
file_args = [str(TMP_DIR / file) for file in OPT_FILES]
try:
    subprocess.run(
        ["go", "run", "./internal/gcc/cmd/cdbconfiggen", "-o", str(OUTPUT_FILE)]
        + file_args,
        check=True,
    )
except subprocess.CalledProcessError as e:
    print(f"Error running cdbconfiggen: {e}", file=sys.stderr)
    sys.exit(1)

print("Bootstrap complete")
