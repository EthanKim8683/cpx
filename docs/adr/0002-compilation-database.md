# ADR-0002: Compilation Database

- **Status**: Accepted
- **Date**: 2026-06-28
- **Related**: [CP-8: Compilation Database ADR](https://linear.app/ethankim8683/issue/CP-8/compilation-database-adr), [GitHub #5](https://github.com/EthanKim8683/cpx/issues/5), [GitHub #24](https://github.com/EthanKim8683/cpx/issues/24)

---

## Context

cpx treats `compile_commands.json` as the single source of truth for per-file compile configuration — including paths, defines, and other flags that drive bundling and related workflows.

cpx needs more than a reader for that file. Workflows must load and save compilation databases, create and update individual entries, and parse compile commands into structured argument lists — then render those lists back when writing. Treating `compile_commands.json` as opaque, read-only input is not enough.

Structured parsing is the hard part. Driver command lines encode thousands of compiler-specific flags (~4000 options per compiler). GCC documents options in `.opt` files; Clang documents them in `.td` files processed by TableGen. Hand-maintaining a parser for that surface area does not scale and drifts from compiler behavior.

Existing tools like Bear, CMake, and Clang's `-MJ` flag can capture real compiler invocations into `compile_commands.json`, but they do not parse, create, or mutate entries programmatically. cpx needs a dedicated compilation database package, backed by generated per-compiler option configs derived from those upstream sources.

---

## Decision

### Package Structure

The compilation database lives in three packages:

- **`internal/cdb`**: Shared types — `OptionPattern`, `OptionKind`, `Config`, `Option`. Compiler-agnostic.
- **`internal/clang`**: Clang-specific generated config (`CDBConfig`) and agent bootstrapping (`AGENTS.md`).
- **`internal/gcc`**: GCC-specific generated config (`CDBConfig`) and agent bootstrapping (`AGENTS.md`).

Each compiler package contains a `cmd/cdbconfiggen/` directory with the generator tool, and a `generated_cdbconfig.go` file produced by that tool.

### Package Responsibilities
The compilation database package is cpx's programmatic interface to `compile_commands.json`. It covers:
- **Database I/O**: Load and save compilation database files on disk.
- **Entry Management**: Create, update, and remove individual entries.
- **Structured Commands**: Parse compile command argument lists structurally, not as opaque strings, and render parsed arguments back to command strings when writing entries.

Capture tools may still record real invocations. cpx uses them as one input, then constructs, merges, or rewrites entries itself — for example when bootstrapping a workspace, synthesizing a command from edited flags, or combining captured output with cpx-generated entries.

### Compiler Option Configs
Structured command parsing relies on generated Go source files containing compiled option prefix maps. This design eliminates runtime parsing overhead, dependencies, and file embedding.

- **Sorted Patterns with Back-Chain Links**: Configs are flat `[]OptionPattern` slices sorted by spelling. Each joined-kind pattern holds a back-chain pointer to the longest entry whose spelling is a strict prefix, enabling longest-prefix matching:
  ```go
  type Config struct {
      Patterns   []OptionPattern
      BackChains []*OptionPattern
  }
  ```
  Both GCC and Clang match options by longest prefix. A binary search locates the candidate; back-chain traversal resolves the actual match. This mirrors GCC's `cl_option` array and `back_chain` links (built by `optc-gen.awk`), and achieves the same result as Clang's sorted forward iteration.
- **Longest-Prefix Matching**: When parsing `--std=c++17`, the parser must match `--std=` (the longest prefix) rather than `--` or `-`. Binary search lands on the exact match or the insertion point. For exact hits on a joined kind, the back-chain returns the longest joined prefix (e.g., `-std=c++17` → `-std=`). If no joined prefix exists, the match is nil — the option is incomplete without a value. For misses, check `patterns[i-1]` directly if joined, then fall back to `BackChains[i-1]` to find the longest candidate whose spelling is a strict prefix of the argument.
- **Query-Time Resolution (Clang Style)**: Dynamic option behaviors—including resolving overridden arguments, negation flags, and mutual exclusions—are deferred entirely to access/query time. Consumers query the resulting parsed argument list using accessors like `getLastArg` or `hasFlag` to dynamically resolve the final compiler state.
- **Driver-Visible Options Only**: Filter at generation time; options not visible to the driver (e.g. those with `NoDriverOption` in Clang or `RejectDriver` in GCC) are omitted from the configuration.
- **Per-Variant Shape**: Each variant (`OptionPattern`) has a spelling, kind, and optional argument count.
- **Option Kinds**: `Flag`, `Joined`, `Separate`, `MultiArg`, `JoinedAndSeparate`, `RemainingArgs`, `RemainingArgsJoined`.

Kinds that don't map 1:1 to an upstream kind are decomposed at generation time:
- `CommaJoined` (Clang) → `Joined`
- `JoinedOrSeparate` (Clang) → `Joined` + `Separate`
- `JoinedOrMissing` (GCC) → `Flag` + `Joined`
- `NoDriverArg` + `Separate` (GCC) → `Flag`

### Config Generation

Generation follows a two-step pipeline: an agent adapts to the host environment, then a stateless generator reads files from disk and produces Go source.

**Agent-driven bootstrapping** (`AGENTS.md`): Each compiler package includes an `AGENTS.md` document that instructs an agent to detect the installed compiler, locate the upstream source repository, discover which source files are required (by reverse-engineering the upstream build system), and fetch them into a `scratch/` directory. The agent handles all environment-specific decisions — platform differences, version mapping, network access.

**Stateless generator** (`cdbconfiggen`): A Go CLI tool reads source files from `scratch/` and writes a `generated_cdbconfig.go` file containing a `*cdb.Config`. The generator contains no environment detection, no network access, and no compiler-specific knowledge beyond parsing the input format.

#### GCC Option Generation

The GCC generator (`internal/gcc/cmd/cdbconfiggen`) reads `.opt` files from a directory:

1. Parses option records from `.opt` files (name + properties format).
2. Merges records with the same name across files.
3. Translates properties to `OptionPattern` kinds (e.g., `Joined` → `OptionKindJoined`, `Separate` → `OptionKindSeparate`).
4. Synthesizes implicit `no-` negation forms for `-f`/`-W`/`-m` names (not listed in `.opt`).
5. Filters out `RejectDriver` options.
6. Writes the config using Go's `%#v` verb.

The agent discovers which `.opt` files are required by parsing `ALL_OPT_FILES` in `gcc/gcc/Makefile.in`, skipping configure substitutions (`@lang_opt_files@` etc.) that aren't available from source alone.

#### Clang Option Generation

The Clang generator (`internal/clang/cmd/cdbconfiggen`) reads a JSON dump produced by `clang-tblgen --dump-json` or `llvm-tblgen --dump-json` (both produce identical output):

1. Unmarshals the TableGen JSON dump (using `go-json-experiment/json` for `embed` support).
2. Iterates all defs, keeping only those with `Option` in their superclass list.
3. Filters out `NoDriverOption` defs.
4. Translates TableGen kinds to `OptionPattern` kinds (e.g., `KIND_JOINED` → `OptionKindJoined`, `KIND_COMMAJOINED` → `OptionKindJoined`).
5. Expands each prefix × kind into a separate pattern.
6. Writes the config using Go's `%#v` verb.

The agent discovers which `.td` files are required by tracing `tablegen()` rules in `CMakeLists.txt` and following `include` directives.

---

## Alternatives

- **Bear, CMake, `-MJ` as the sole compdb layer**: Useful for capture, but produce files only; rejected because cpx must also parse, create, and mutate entries.
- **Hand-maintained flag parsers**: Rejected; does not scale (~4000 driver options) and drifts from compiler behavior.
- **Custom TableGen C++ backend**: Rejected in favor of `--dump-json` + Go generator; avoids maintaining code inside LLVM.
- **Compile-Time Alias & Negation Flattening**: Rejected; baking alias maps and negation override chains into the compiled config structure was rejected in favor of query-time resolution to keep the schema simple and avoid complex generator logic.
- **Generated YAML configs**: Originally proposed but rejected in favor of Go source files to eliminate runtime YAML parsing library dependencies and unmarshalling overhead.

---

## Consequences

- **Separate configs per compiler**: Runtime must select the correct config map for the compiler named in a compdb entry.
- **Response files (`@file`)**: Not covered by option configs; expand or pass through before parsing if present in a compile command.
- **Write path fidelity**: Rendering parsed commands back to strings may not preserve the original quoting or token order; consumers should treat parsed structure as authoritative.
- **`scratch/` directories**: Intermediate files (`.td` files, `.opt` files, JSON dumps) go in `internal/<compiler>/scratch/`, which is gitignored. Created by the agent during bootstrapping.
- **Generator simplicity**: Since alias resolution and negation override relationships are deferred to query-time helpers, the generator logic remains simple — it only needs to parse source files and synthesize implicit negations.
- **Regeneration**: Regenerate configs when upgrading the compiler version. The `AGENTS.md` documents describe the full bootstrapping pipeline.

---

## References

- [CP-8: Compilation Database ADR](https://linear.app/ethankim8683/issue/CP-8/compilation-database-adr)
- [GitHub Issue #5](https://github.com/EthanKim8683/cpx/issues/5)
- [GitHub Issue #24](https://github.com/EthanKim8683/cpx/issues/24)
- [Go Doc Comments Specification](https://go.dev/doc/comment)
