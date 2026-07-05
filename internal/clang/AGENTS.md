# Clang CDB Configuration Generator — Agent Bootstrapping

Generate `internal/clang/generated_cdbconfig.go` — a Go source file
containing Clang's option configuration as a `*cdb.Config` value for the
compilation database.

The agent handles the full pipeline: detecting the installed Clang, fetching
upstream TableGen sources, producing a JSON dump, and running the code
generator. The upstream LLVM/Clang repository is the ground truth for every
step — each phase below mirrors a stage in the upstream build system.

Intermediate files go in `internal/clang/scratch/`. Create it if it does not
exist. It is gitignored.

## Steps

### 1. Detect the installed Clang

Read the `CLANG` environment variable (set in `.env` and loaded via direnv)
to find the Clang executable path. Run:

```sh
$CLANG --version
```

Parse the output to determine the version number and distribution (upstream
LLVM, Apple Clang, etc.). The version output typically identifies which one.

If the distribution cannot be determined, ask the user to clarify.

### 2. Locate the upstream source

Find the matching LLVM/Clang source repository on GitHub. This repository
contains the TableGen files and build system definitions that dictate how
option generation works — all subsequent steps derive from it.

| Distribution | Raw content base URL | Notes |
|---|---|---|
| **LLVM** (upstream) | `https://raw.githubusercontent.com/llvm/llvm-project/<tag>/` | Tags: `release/<major>.x` or `llvmorg-<version>` |
| **Apple Clang** | `https://raw.githubusercontent.com/swiftlang/llvm-project/<tag>/` | Apple's fork. Only publishes Swift tags, not clang version tags. |

Apple Clang is built from Apple's fork (`swiftlang/llvm-project`), but Apple
does not publish clang version tags — only Swift release and snapshot tags.
For option definitions, fall back to the matching upstream LLVM release
(`llvm/llvm-project` with `release/<major>.x`). The option definitions are
nearly identical. For example, Apple Clang 17 → `release/17.x`.

Append a relative file path to the base URL to fetch raw content. If the
distribution does not match one of these, search GitHub for the appropriate
repository and derive its raw content base URL.

### 3. Determine which `.td` files are required

Do not assume a fixed set of files. Reverse-engineer the required TableGen
sources from the upstream build system:

1. **Find the tablegen rule.** Locate the CMakeLists.txt that defines the
   `tablegen()` call for the Clang driver options target. It specifies the
   entry-point `.td` file and backend flags (e.g. `-gen-opt-parser-defs`).
2. **Follow `include` directives.** Trace `include` directives within the
   `.td` files to discover transitive dependencies. TableGen resolves includes
   relative to the `-I` flags, which in the LLVM build correspond to
   `llvm/include` and `clang/include`.
3. **Record what you find.** Note each required `.td` file, its repository
   path, and its role (entry point, schema, etc.).

### 4. Fetch the TableGen source files

Download the required `.td` files into `internal/clang/scratch/`, **preserving
the directory structure from the repository**. This is critical — tblgen
resolves `include` directives relative to `-I` paths, and a mirrored layout
ensures those paths resolve correctly.

For example, if the repository has:
```
llvm-project/clang/include/clang/Driver/Options.td
llvm-project/llvm/include/llvm/Option/OptParser.td
```

Then `scratch/` should contain:
```
internal/clang/scratch/clang/include/clang/Driver/Options.td
internal/clang/scratch/llvm/include/llvm/Option/OptParser.td
```

Use the raw content base URL from step 2. If network access is restricted,
inform the user and ask them to place the files manually.

### 5. Generate the JSON dump

Replicate the upstream tablegen invocation for standalone use. The
CMakeLists.txt rule from step 3 identifies the entry-point `.td` file and
backend flags. Substitute `--dump-json` for the original backend flag to
produce JSON instead of C++ source.

Either `clang-tblgen` or `llvm-tblgen` can be used — both support
`--dump-json` and produce identical output.

Infer the `-I` flags from the repository layout so that `include` directives
resolve. For example, if the CMakeLists.txt says:
```cmake
set(LLVM_TARGET_DEFINITIONS Options.td)
tablegen(LLVM Options.inc -gen-opt-parser-defs)
```

And the entry-point file is at `clang/include/clang/Driver/Options.td`, the
equivalent standalone invocation is:
```sh
clang-tblgen -I ./internal/clang/scratch/llvm/include \
  -I ./internal/clang/scratch/clang/include \
  ./internal/clang/scratch/clang/include/clang/Driver/Options.td \
  --dump-json > ./internal/clang/scratch/options.json
```

Verify the exact paths and flags against the upstream source for the detected
version.

### 6. Fallback: if tblgen is not available

If neither `clang-tblgen` nor `llvm-tblgen` can be found:

1. Ask the user if they would like to install one. Provide platform-specific
   instructions if possible.
2. If declined, check for a pre-built JSON dump online (CI artifacts, LLVM
   release assets, package repositories).
3. As a last resort, ask the user to provide the dump file directly.

Do not install system packages without user permission.

### 7. Verify the dump

Before running the generator, confirm the dump is valid:

- Non-empty and valid JSON.
- Contains `"!tablegen_json_version"` with value `1`.
- Contains defs with `Option` in their superclass list.

A quick way to check all three at once:
```sh
python3 -c "
import json
with open('./internal/clang/scratch/options.json') as f:
    data = json.load(f)
assert data.get('!tablegen_json_version') == 1, 'bad version'
opts = sum(1 for k,v in data.items()
    if not k.startswith('!')
    and isinstance(v, dict)
    and 'Option' in v.get('!superclasses', []))
assert opts > 0, 'no Option defs'
print(f'OK — {opts} option defs')
"
```

### 8. Run the code generator

```sh
go run ./internal/clang/cmd/cdbconfiggen \
  -o internal/clang/generated_cdbconfig.go \
  ./internal/clang/scratch/options.json
```

This reads the JSON dump and writes `generated_cdbconfig.go` containing the
`CDBConfig` variable — a `*cdb.Config` with all Clang CDB option patterns
indexed by spelling prefix.
