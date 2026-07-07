# ADR-0002: Compilation Database

**Status:** Accepted

## Context

cpx treats `compile_commands.json` as the single source of truth for per-file compile configuration. But reading is not enough — workflows must load, save, create, and update entries, and parse compile commands into structured argument lists.

Hand-maintaining parsers for compiler-specific flags (~4000 options per compiler) does not scale. GCC documents options in `.opt` files; Clang in `.td` files via TableGen. Capture tools (Bear, CMake, `-MJ`) produce files but cannot parse or mutate entries programmatically.

## Decision

**Package structure** (three packages):

- **`internal/cdb`**: Shared types — `OptionPattern`, `OptionKind`, `Config`, `Option`. Compiler-agnostic.
- **`internal/clang`**: Clang-specific generated config (`CDBConfig`) and agent bootstrapping.
- **`internal/gcc`**: GCC-specific generated config (`CDBConfig`) and agent bootstrapping.

**Package responsibilities** — the compilation database is cpx's programmatic interface to `compile_commands.json`:

- **Database I/O**: Load and save compilation database files.
- **Entry Management**: Create, update, and remove individual entries.
- **Structured Commands**: Parse compile command argument lists structurally, and render parsed arguments back to command strings when writing.

**Compiler option configs** — generated Go source files containing compiled prefix maps. Eliminates runtime parsing, dependencies, and file embedding.

- **Sorted patterns with back-chain links**: Flat `[]OptionPattern` slices sorted by spelling. Mirrors GCC's `cl_option` array — `optc-gen.awk` builds a sorted table where each entry has a `back_chain` index pointing to the longest entry whose spelling is a strict prefix. Binary search locates the candidate; back-chain traversal resolves the longest match. Clang achieves the same result via sorted forward iteration.
- **Driver-visible options only**: Filtered at generation time — options with `NoDriverOption` (Clang) or `RejectDriver` (GCC) are omitted.
- **Per-variant shape**: Each `OptionPattern` has a spelling, kind, and optional argument count.
- **Option kinds**: `Flag`, `Joined`, `Separate`, `MultiArg`, `JoinedAndSeparate`, `RemainingArgs`, `RemainingArgsJoined`. This set is atomic — these are the smallest building blocks that compose the behavior of every kind supported by GCC and Clang. In cpx, `Joined` options must have a non-empty argument (e.g. `-std=c++20`, not `-std`).
- **Kinds decomposed at generation time**: Upstream kinds that don't map 1:1 to our atomic set are decomposed:
  - `CommaJoined` (Clang) → `Joined`
  - `JoinedOrSeparate` (Clang) → `Joined` + `Separate`
  - `JoinedOrMissing` (GCC) → `Flag` + `Joined`
  - `NoDriverArg` + `Separate` (GCC) → `Flag`
- **Query-time resolution (Clang-inspired)**: Dynamic behaviors — overridden args, negation, mutual exclusion — are deferred to access time, not baked into the config. This follows Clang's pattern where `InputArgList` collects all options flat, and `getLastArg`/`hasFlag` resolve the final state at query time.

**Config generation** — two-step pipeline:

1. **Agent-driven bootstrapping** (`AGENTS.md`): Detects installed compiler, locates upstream source, discovers required files, fetches into `scratch/`.
2. **Stateless generator** (`cdbconfiggen`): Reads files from `scratch/`, writes `generated_cdbconfig.go`. No environment detection, no network access.

GCC generator reads `.opt` files — parses option records, merges across files, translates properties to kinds, synthesizes implicit `no-` negations for `-f`/`-W`/`-m` names, filters `RejectDriver`.

Clang generator reads a JSON dump from `clang-tblgen --dump-json` or `llvm-tblgen --dump-json` — unmarshals, keeps defs with `Option` superclass, filters `NoDriverOption`, translates kinds, expands prefix × kind into separate patterns.

## Alternatives considered

- **Bear, CMake, `-MJ` as the sole compdb layer**: Useful for capture, but produce files only. (Rejected: cpx must also parse, create, and mutate entries).
- **Hand-maintained flag parsers**: (Rejected: does not scale at ~4000 options and drifts from compiler behavior).
- **Custom TableGen C++ backend**: (Rejected: maintaining code inside LLVM is costly; `--dump-json` + Go generator is simpler).
- **Compile-time alias & negation flattening**: (Rejected: adds complex generator logic; query-time resolution keeps the schema simple).
- **Generated YAML configs**: (Rejected: introduces runtime YAML parsing and unmarshalling overhead).

## Consequences

- **Separate configs per compiler**: Runtime selects the correct config for the compiler named in each entry.
- **Response files (`@file`)**: Not covered by option configs; expand or pass through before parsing.
- **Write path fidelity**: Rendering parsed commands back to strings may not preserve original quoting or token order.
- **`scratch/` directories**: Intermediate files go in `internal/<compiler>/scratch/`, gitignored.
- **Regeneration**: Regenerate configs when upgrading the compiler version.

## References

- [GitHub #5](https://github.com/EthanKim8683/cpx/issues/5)
- [GitHub #24](https://github.com/EthanKim8683/cpx/issues/24)
