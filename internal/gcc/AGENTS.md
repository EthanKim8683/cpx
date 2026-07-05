# GCC CDB Configuration Generator — Agent Bootstrapping

Generate `internal/gcc/generated_cdbconfig.go` — a Go source file
containing GCC's option configuration as a `*cdb.Config` value for the
compilation database.

The agent handles the full pipeline: detecting the installed GCC, locating
the `.opt` source files, and running the code generator. The upstream GCC
repository is the ground truth for all option definitions — each phase
below mirrors a stage in the upstream build system.

Intermediate files go in `internal/gcc/scratch/`. Create it if it does not
exist. It is gitignored.

## Steps

### 1. Detect the installed GCC

Read the `GCC` environment variable to find the GCC executable path. Run:

```sh
$GCC --version
```

If the variable is not set, work with the user to locate their GCC
installation.

Parse the output to determine the version number and installation prefix.
The version output typically identifies the distribution (Homebrew, system,
etc.).

If the version cannot be determined, ask the user to clarify.

### 2. Locate the upstream source

Find the matching GCC source repository. The upstream repository contains
the `.opt` files and build system definitions that dictate which option
sources are required — all subsequent steps derive from it.

GCC typically has a mirror at `https://github.com/gcc-mirror/gcc`. Find the
branch or tag matching the detected version:

- **Released versions** (e.g., 14.1.0): look for a tag like
  `releases/gcc-<major>`.
- **Development versions** (e.g., 16.1.0 before release): use the `trunk`
  branch, which tracks the current development head.

### 3. Discover which `.opt` files are required

Do not assume a fixed set of files. Reverse-engineer the required `.opt`
sources from the upstream build system:

1. **Find the Make variable.** Locate `gcc/gcc/Makefile.in` in the
   repository. Within it, find the variable that aggregates `.opt` files
   (typically `ALL_OPT_FILES` or similar) and trace the variables it
   references.
2. **Parse the variable assignments.** Resolve the variable chain to
   determine the full set of `.opt` files. Paths are relative to `gcc/gcc/`
   in the repository.
3. **Skip configure substitutions.** Variable assignments containing
   `@`-delimited configure substitutions (e.g., `@lang_opt_files@`) expand
   at build time and are not available from the source alone. Only use the
   concrete file paths listed directly in `Makefile.in`.
4. **Record what you find.** Note each required `.opt` file, its
   repository path, and its role (driver options, language options, etc.).

### 4. Fetch the `.opt` source files

Download the files discovered in step 3 into `internal/gcc/scratch/`,
**preserving the directory structure from the repository**. This is critical
for the code generator to correctly resolve relative paths within `.opt`
files.

Use the raw content base URL for the repository and branch/tag identified
in step 2. Append each relative file path to fetch the content.

If network access is restricted, see step 5.

### 5. Fallback: if `.opt` files cannot be fetched

Work with the user to resolve this. They may have a local source tree,
want to clone the repository, or be able to place the files manually.

Do not install system packages or clone repositories without user permission.

### 6. Run the code generator

```sh
go run ./internal/gcc/cmd/cdbconfiggen \
  -o internal/gcc/generated_cdbconfig.go \
  ./internal/gcc/scratch/gcc
```

This reads all `.opt` files from the directory, parses option records, and
writes `generated_cdbconfig.go` containing the `CDBConfig` variable — a
`*cdb.Config` with all GCC CDB option patterns indexed by spelling prefix.
