# ADR-0002: Compilation Database

**Status:** Accepted

## Context

cpx treats `compile_commands.json` as the single source of truth for per-file compile configuration. Reading is not enough, though — cpx also needs to load, save, create, and update entries, and parse compile commands into structured argument lists.

The hard part is structured parsing. Driver command lines encode thousands of compiler-specific flags (~4000 per compiler). GCC documents options in `.opt` files; Clang in `.td` files processed by TableGen. Hand-maintaining parsers for that surface area does not scale and drifts from compiler behavior. Capture tools like Bear, CMake, and Clang's `-MJ` flag can produce `compile_commands.json`, but they cannot parse or mutate entries programmatically.

## Decision

The compilation database lives in three packages:

- **`internal/cdb`** — compiler-agnostic shared types (`OptionPattern`, `OptionKind`, `Config`, `Option`).
- **`internal/gcc`** — GCC-specific generated config and agent bootstrapping.
- **`internal/clang`** — Clang-specific generated config and agent bootstrapping.

Each compiler package contains a `cmd/cdbconfiggen/` directory with the generator tool and a `generated_cdbconfig.go` file produced by that tool.

The package provides three capabilities: database I/O (load/save compilation database files), entry management (create/update/remove entries), and structured commands (parse compile command argument lists into structured form, and render them back to strings when writing).

### Option configs

The core design decision is how to parse compiler flags. Instead of hand-maintaining parsers, cpx generates Go source files containing compiled option prefix maps. This eliminates runtime parsing overhead, external dependencies, and file embedding.

Each generated config is a flat `[]OptionPattern` slice sorted by spelling. The structure mirrors GCC's `cl_option` array: `optc-gen.awk` builds a sorted table where each entry has a `back_chain` index pointing to the longest entry whose spelling is a strict prefix. At parse time, binary search locates the candidate and back-chain traversal resolves the longest match. Clang achieves the same result via sorted forward iteration.

Only driver-visible options are included — options with `NoDriverOption` (Clang) or `RejectDriver` (GCC) are filtered out at generation time.

Each `OptionPattern` has a spelling, a kind, and an optional argument count. The kind set is intentionally atomic: `Flag`, `Joined`, `Separate`, `MultiArg`, `JoinedAndSeparate`, `RemainingArgs`, `RemainingArgsJoined`. These are the smallest building blocks — every upstream kind from GCC and Clang can be composed from them. In cpx, `Joined` options must have a non-empty argument (e.g. `-std=c++20`, not `-std`).

Upstream kinds that don't map 1:1 to this atomic set are decomposed at generation time:

- `CommaJoined` (Clang) → `Joined`
- `JoinedOrSeparate` (Clang) → `Joined` + `Separate`
- `JoinedOrMissing` (GCC) → `Flag` + `Joined`
- `NoDriverArg` + `Separate` (GCC) → `Flag`

Dynamic behaviors — overridden arguments, negation, mutual exclusion — are deferred to query time rather than baked into the config. This is inspired by Clang's pattern where `InputArgList` collects all options flat, and accessors like `getLastArg` or `hasFlag` resolve the final compiler state at access time.

### Config generation

Generation follows a two-step pipeline:

1. **Agent-driven bootstrapping** — each compiler package has an `AGENTS.md` that instructs an agent to detect the installed compiler, locate the upstream source repository, discover which files are required, and fetch them into a `scratch/` directory.
2. **Stateless generator** (`cdbconfiggen`) — reads files from `scratch/` and writes `generated_cdbconfig.go`. No environment detection, no network access, no compiler-specific knowledge beyond parsing the input format.

The GCC generator reads `.opt` files: parses option records, merges across files, translates properties to kinds, synthesizes implicit `no-` negations for `-f`/`-W`/`-m` names, and filters `RejectDriver`.

The Clang generator reads a JSON dump from `clang-tblgen --dump-json` or `llvm-tblgen --dump-json`: unmarshals, keeps defs with `Option` in their superclass list, filters `NoDriverOption`, translates kinds, and expands each prefix × kind into a separate pattern.

## Alternatives considered

- **Bear, CMake, `-MJ` as the sole compdb layer**: Useful for capture, but produce files only. cpx must also parse, create, and mutate entries.
- **Hand-maintained flag parsers**: Does not scale at ~4000 options and drifts from compiler behavior.
- **Custom TableGen C++ backend**: Maintaining code inside LLVM is costly; `--dump-json` + Go generator is simpler.
- **Compile-time alias & negation flattening**: Adds complex generator logic; query-time resolution keeps the schema simple.
- **Generated YAML configs**: Introduces runtime YAML parsing and unmarshalling overhead.

## Consequences

- **Separate configs per compiler**: Runtime must select the correct config for the compiler named in each entry.
- **Response files (`@file`)**: Not covered by option configs; expand or pass through before parsing.
- **Write path fidelity**: Rendering parsed commands back to strings may not preserve original quoting or token order.
- **`scratch/` directories**: Intermediate files go in `internal/<compiler>/scratch/`, gitignored.
- **Regeneration**: Regenerate configs when upgrading the compiler version.

## References

- [GitHub #5](https://github.com/EthanKim8683/cpx/issues/5)
- [GitHub #24](https://github.com/EthanKim8683/cpx/issues/24)
