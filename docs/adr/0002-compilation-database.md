# ADR-0002: Compilation Database

- **Status**: Accepted
- **Date**: 2026-06-28
- **Related**: [CP-8: Compilation Database ADR](https://linear.app/ethankim8683/issue/CP-8/compilation-database-adr), [GitHub #5](https://github.com/EthanKim8683/cpx/issues/5)

---

## Context

cpx treats `compile_commands.json` as the single source of truth for per-file compile configuration — including paths, defines, and other flags that drive bundling and related workflows.

cpx needs more than a reader for that file. Workflows must load and save compilation databases, create and update individual entries, and parse compile commands into structured argument lists — then render those lists back when writing. Treating `compile_commands.json` as opaque, read-only input is not enough.

Structured parsing is the hard part. Driver command lines encode thousands of compiler-specific flags (~4000 options per compiler). GCC documents options in `.opt` files processed by AWK scripts; Clang documents them in `.td` files processed by TableGen. Hand-maintaining a parser for that surface area does not scale and drifts from compiler behavior.

Existing tools like Bear, CMake, and Clang's `-MJ` flag can capture real compiler invocations into `compile_commands.json`, but they do not parse, create, or mutate entries programmatically. cpx needs a dedicated compilation database package, backed by generated per-compiler option configs derived from those upstream sources.

---

## Decision

### Package Responsibilities
The compilation database package is cpx's programmatic interface to `compile_commands.json`. It covers:
- **Database I/O**: Load and save compilation database files on disk.
- **Entry Management**: Create, update, and remove individual entries.
- **Structured Commands**: Parse compile command argument lists structurally, not as opaque strings, and render parsed arguments back to command strings when writing entries.

Capture tools may still record real invocations. cpx uses them as one input, then constructs, merges, or rewrites entries itself — for example when bootstrapping a workspace, synthesizing a command from edited flags, or combining captured output with cpx-generated entries.

### Compiler Option Configs
Structured command parsing relies on generated Go source code files (`gccoptions.go` and `clangoptions.go`) containing compiled option metadata maps. This design eliminates runtime parsing overhead, dependencies, and file embedding.
- **Go Option Maps**: Generated maps are built at package initialization time, enabling immediate $O(1)$ lookups based on the compiler named in a compdb entry.
- **Driver-Visible Options Only**: Filter at generation time; options not visible to the driver are omitted.
- **Equivalence Classes**: Options grouped by meaning; matching any spelling in a group matches all spellings in that group. No canonical names.
- **Per-Variant Shape**: Each variant has a spelling, kind, and optional argument count.
- **Option Kinds**: `Flag`, `Joined`, `Separate`, `JoinedOrSeparate`, `CommaJoined`, `MultiArg`, `JoinedAndSeparate`, `JoinedOrMissing`.

At load time, we build a spelling index (exact match for flags and separate options; longest-prefix match for joined spellings).

Per-variant shape is required because spellings in the same equivalence class can still parse argv differently (e.g. `-std=` joined vs `-std` joined-or-separate).

### Config Generation
We generate Go source files containing option maps directly:
- **GCC Option Generation**: Uses an AWK generator (building on GCC's existing `.opt` option scripts) to extract and compile options from `.opt` files directly into Go structures.
- **Clang Option Generation**: Uses `clang-tblgen -dump-json` to extract options from `Driver/Options.td`, processed by a Go generator directly into Go structures.

For Clang, use `-dump-json` and a Go generator in cpx — not a custom TableGen C++ backend inside LLVM.

The GCC generator must also synthesize implicit `no-` forms for `-f`/`-W`/`-m`/`-g` options (not listed in `.opt`) and deduplicate records that appear across multiple `.opt` files.

---

## Alternatives

- **Bear, CMake, `-MJ` as the sole compdb layer**: Useful for capture, but produce files only; rejected because cpx must also parse, create, and mutate entries.
- **Hand-maintained flag parsers**: Rejected; does not scale (~4000 driver options) and drifts from compiler behavior.
- **Custom TableGen C++ backend**: Rejected in favor of `-dump-json` + Go generator; avoids maintaining code inside LLVM.
- **Canonical names as config keys**: Dropped; equivalence classes keyed by shared spellings are sufficient.
- **Visibility in config**: Rejected; filter to driver-visible options at generation instead.
- **Kind at the group level only**: Insufficient when spellings in the same group parse argv differently (e.g. `-std=` vs `-std`).
- **Generated YAML configs**: Originally proposed but rejected in favor of Go source files to eliminate runtime YAML parsing library dependencies and unmarshalling overhead.

---

## Consequences

- **Separate configs per compiler**: Runtime must select the correct config map for the compiler named in a compdb entry.
- **Response files (`@file`)**: Not covered by option configs; expand or pass through before parsing if present in a compile command.
- **Write path fidelity**: Rendering parsed commands back to strings may not preserve the original quoting or token order; consumers should treat parsed structure as authoritative.
- **`clang-tblgen` at generation time**: Clang configs require a pinned LLVM/clang version; regenerate when upgrading the compiler.
- **Generator complexity**: Alias flattening, GCC negated-form synthesis, and prefix-collision handling live in the generators, not the runtime schema.

---

## References

- [CP-8: Compilation Database ADR](https://linear.app/ethankim8683/issue/CP-8/compilation-database-adr)
- [GitHub Issue #5](https://github.com/EthanKim8683/cpx/issues/5)
