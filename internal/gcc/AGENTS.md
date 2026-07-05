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

Read the `GCC` environment variable (set in `.env` and loaded via direnv)
to find the GCC executable path. Run:

```sh
$GCC --version
```

Parse the output to determine the version number and installation prefix.
The version output typically identifies the distribution (Homebrew, system,
etc.).

If the version cannot be determined, ask the user to clarify.

### 2. Locate the `.opt` source files

GCC option definitions live in `.opt` files within the source tree. The
agent must find or obtain them:

| Source | Path in repo | Notes |
|---|---|---|
| **Upstream GCC** | `gcc/gcc/` | Mirror: `https://github.com/gcc-mirror/gcc` with tag `releases/gcc-<major>` |
| **Homebrew** | Varies by platform | Typically not shipped with the installed GCC — the agent must fetch from upstream. |

The key `.opt` files are:

| File | Description |
|---|---|
| `c.opt` | C and C++ language options |
| `common.opt` | Common driver options shared across frontends |
| `driver` | Driver-level options |
| `objc` | Objective-C language options |
| `go` | Go language options |
| `lto` | Link-time optimization options |

Do not assume a fixed list. Trace which `.opt` files are required from the
upstream source for the detected version.

### 3. Fetch the `.opt` source files

Download the required `.opt` files into `internal/gcc/scratch/`, **preserving
the directory structure from the repository**. This is critical for the
agent to reference the original layout.

For example, if the repository has:
```
gcc/gcc/c.opt
gcc/gcc/common.opt
gcc/gcc/driver
```

Then `scratch/` should contain:
```
internal/gcc/scratch/gcc/c.opt
internal/gcc/scratch/gcc/common.opt
internal/gcc/scratch/gcc/driver
```

Use raw content URLs from the upstream mirror. If network access is
restricted, inform the user and ask them to place the files manually.

### 4. Run the code generator

```sh
go run ./internal/gcc/cmd/cdbconfiggen \
  -o internal/gcc/generated_cdbconfig.go \
  ./internal/gcc/scratch/gcc
```

This reads all `.opt` files from the directory, parses option records, and
writes `generated_cdbconfig.go` containing the `CDBConfig` variable — a
`*cdb.Config` with all GCC CDB option patterns indexed by spelling prefix.
