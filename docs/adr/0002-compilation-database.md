# ADR-0002: Compilation Database

- **Status**: Accepted
- **Date**: 2026-06-28
- **Related**: [CP-8: Compilation Database ADR](https://linear.app/ethankim8683/issue/CP-8/compilation-database-adr), [GitHub #5](https://github.com/EthanKim8683/cpx/issues/5)

---

## Context

cpx treats `compile_commands.json` as the single source of truth for per-file compile configuration — including paths, defines, and other flags that drive bundling and related workflows.

cpx needs more than a reader for that file. Workflows must load and save compilation databases, create and update individual entries, and parse compile commands into structured argument lists — then render those lists back when writing. Treating `compile_commands.json` as opaque, read-only input is not enough.

Structured parsing is the hard part. Driver command lines encode thousands of compiler-specific flags (~4000 options per compiler). GCC documents options in `.opt` files; Clang documents them in `.td` files processed by TableGen. Hand-maintaining a parser for that surface area does not scale and drifts from compiler behavior.

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
Structured command parsing relies on generated Go source code files (e.g. `config/gcc.go` and `config/clang.go`) containing compiled option prefix maps. This design eliminates runtime parsing overhead, dependencies, and file embedding.
- **Stateless Prefix Registry**: Configs are structured as a flat registry of option patterns indexed directly by spelling prefix:
  ```go
  type Config struct {
      ByPrefix map[string][]OptionPattern
  }
  ```
  This maps option prefixes (like `-std=` or `-o`) to their respective `OptionPattern` specs, allowing the parser to slice command-line arguments statelessly without managing alias tables or exclusion trees during parsing.
- **Query-Time Resolution (Clang Style)**: Dynamic option behaviors—including resolving overridden arguments, negation flags, and mutual exclusions—are deferred entirely to access/query time. Consumers query the resulting parsed argument list (analogous to Clang's `InputArgList`) using accessors like `getLastArg` or `hasFlag` to dynamically resolve the final compiler state.
- **Driver-Visible Options Only**: Filter at generation time; options not visible to the driver (e.g. those with the `RejectDriver` property) are omitted from the configuration.
- **Per-Variant Shape**: Each variant has a spelling, kind, and optional argument count.
- **Option Kinds**: `Flag`, `Joined`, `Separate`, `JoinedOrSeparate`, `CommaJoined`, `MultiArg`, `JoinedAndSeparate`, `JoinedOrMissing`.

### Config Generation
We generate Go source files containing option maps directly:
- **GCC Option Generation**: Uses a Go-based configuration generator (`gccconfiggen` CLI tool) which queries the local GCC driver version, downloads the matching tag-locked `.opt` files from the source mirror, extracts/parses records, and generates the flat prefix map.
- **Clang Option Generation**: Uses `clang-tblgen -dump-json` to extract options from `Driver/Options.td`, processed by a Go generator directly into Go structures.

For Clang, use `-dump-json` and a Go generator in cpx — not a custom TableGen C++ backend inside LLVM.

The GCC generator must also synthesize implicit `no-` forms for `-f`/`-W`/`-m`/`-g` options (not listed in `.opt`) and deduplicate records that appear across multiple `.opt` files.

---

## Alternatives

- **Bear, CMake, `-MJ` as the sole compdb layer**: Useful for capture, but produce files only; rejected because cpx must also parse, create, and mutate entries.
- **Hand-maintained flag parsers**: Rejected; does not scale (~4000 driver options) and drifts from compiler behavior.
- **Custom TableGen C++ backend**: Rejected in favor of `-dump-json` + Go generator; avoids maintaining code inside LLVM.
- **Compile-Time Alias & Negation Flattening**: Rejected baking alias maps and negation override chains (like GCC's `.neg_index` circular chains) into the compiled config structure. This was rejected in favor of Clang's query-time resolution to keep the schema simple and avoid complex generator logic.
- **Visibility in config**: Rejected; filter to driver-visible options at generation instead.
- **Kind at the group level only**: Insufficient when spellings in the same group parse argv differently (e.g. `-std=` vs `-std`).
- **Generated YAML configs**: Originally proposed but rejected in favor of Go source files to eliminate runtime YAML parsing library dependencies and unmarshalling overhead.

---

## Consequences

- **Separate configs per compiler**: Runtime must select the correct config map for the compiler named in a compdb entry.
- **Response files (`@file`)**: Not covered by option configs; expand or pass through before parsing if present in a compile command.
- **Write path fidelity**: Rendering parsed commands back to strings may not preserve the original quoting or token order; consumers should treat parsed structure as authoritative.
- **`clang-tblgen` at generation time**: Clang configs require a pinned LLVM/clang version; regenerate when upgrading the compiler.
- **Generator simplicity**: Since alias resolution and negation override relationships are deferred to query-time helpers, the generator logic remains simple and only needs to synthesize implicit negations.

---

## References

- [CP-8: Compilation Database ADR](https://linear.app/ethankim8683/issue/CP-8/compilation-database-adr)
- [GitHub Issue #5](https://github.com/EthanKim8683/cpx/issues/5)
- [Go Doc Comments Specification](https://go.dev/doc/comment)
