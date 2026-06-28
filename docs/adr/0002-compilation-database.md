# ADR-0002: Compilation Database

- **Status:** Accepted
- **Date:** 2026-06-28
- **Related:** [CP-8: Compilation Database ADR](https://linear.app/ethankim8683/issue/CP-8/compilation-database-adr), [GitHub #5](https://github.com/EthanKim8683/cpx/issues/5)

## Context

cpx treats `compile_commands.json` as the single source of truth for per-file compile configuration — include paths, defines, and other flags that drive bundling and related workflows.

cpx needs more than a reader for that file. Workflows must load and save compilation databases, create and update individual entries, and parse compile commands into structured argument lists — then render those lists back when writing. Treating `compile_commands.json` as opaque, read-only input is not enough.

Structured parsing is the hard part. Driver command lines encode thousands of compiler-specific flags (~4000 options per compiler). GCC documents options in `.opt` files processed by AWK scripts; Clang documents them in `.td` files processed by TableGen. Hand-maintaining a parser for that surface area does not scale and drifts from compiler behavior.

Existing tools like Bear, CMake, and Clang's `-MJ` flag can capture real compiler invocations into `compile_commands.json`, but they do not parse, create, or mutate entries programmatically. cpx needs a dedicated compilation database package, backed by generated per-compiler option configs derived from those upstream sources.

## Decision

### Package responsibilities

The compilation database package is cpx's programmatic interface to `compile_commands.json`. It covers:

- **Database I/O** — load and save compilation database files on disk.
- **Entry management** — create, update, and remove individual entries.
- **Structured commands** — parse compile command argument lists structurally, not as opaque strings, and render parsed arguments back to command strings when writing entries.

Capture tools may still record real invocations. cpx uses them as one input, then constructs, merges, or rewrites entries itself — for example when bootstrapping a workspace, synthesizing a command from edited flags, or combining captured output with cpx-generated entries.

### Compiler option configs

Structured command parsing relies on generated YAML configs. These describe **option shape for parsing and grouping**, not validation or semantics — consumers decide meaning.

- **One config per compiler** (GCC, Clang), loaded separately at runtime based on the compiler named in a compdb entry.
- **Driver-visible options only** — filter at generation time; options not visible to the driver are omitted.
- **Equivalence classes** — options grouped by meaning; matching any spelling in a group matches all spellings in that group. No canonical names.
- **Per-variant shape** — each variant has a `spelling`, `kind`, and optional `num_args`:

```yaml
options:
  - variants:
      - {spelling: "-o", kind: Separate}
      - {spelling: "--output", kind: Separate}
  - variants:
      - {spelling: "-std=", kind: Joined}
      - {spelling: "-std", kind: JoinedOrSeparate}
```

**Kinds:** `Flag`, `Joined`, `Separate`, `JoinedOrSeparate`, `CommaJoined`, `MultiArg`, `JoinedAndSeparate`, `JoinedOrMissing`.

At load time, build a spelling index (exact match for flags and separate options; longest-prefix match for joined spellings). ~4000 entries per compiler is not a performance concern for cpx's workload.

Per-variant shape is required because spellings in the same equivalence class can still parse argv differently (e.g. `-std=` joined vs `-std` joined-or-separate).

### Config generation

| Compiler | Source | Pipeline |
|----------|--------|----------|
| **GCC** | `.opt` files | AWK generator (building on GCC's existing option scripts) → YAML |
| **Clang** | `Driver/Options.td` | `clang-tblgen -dump-json` → Go generator → YAML |

For Clang, use `-dump-json` and a Go generator in cpx — not a custom TableGen C++ backend inside LLVM.

The GCC generator must also synthesize implicit `no-` forms for `-f`/`-W`/`-m`/`-g` options (not listed in `.opt`) and deduplicate records that appear across multiple `.opt` files.

## Alternatives

- **Bear, CMake, `-MJ` as the sole compdb layer** — useful for capture, but produce files only; rejected because cpx must also parse, create, and mutate entries.
- **Hand-maintained flag parsers** — rejected; does not scale (~4000 driver options) and drifts from compiler behavior.
- **Custom TableGen C++ backend** — rejected in favor of `-dump-json` + Go generator; avoids maintaining code inside LLVM.
- **Canonical names as config keys** — dropped; equivalence classes keyed by shared spellings are sufficient.
- **Visibility in config** — rejected; filter to driver-visible options at generation instead.
- **Kind at the group level only** — insufficient when spellings in the same group parse argv differently (e.g. `-std=` vs `-std`).

## Consequences

- **Separate configs per compiler** — runtime must select the correct config for the compiler named in a compdb entry.
- **Response files (`@file`)** — not covered by option configs; expand or pass through before parsing if present in a compile command.
- **Write path fidelity** — rendering parsed commands back to strings may not preserve the original quoting or token order; consumers should treat parsed structure as authoritative.
- **`clang-tblgen` at generation time** — Clang configs require a pinned LLVM/clang version; regenerate when upgrading the compiler.
- **Generator complexity** — alias flattening, GCC negated-form synthesis, and prefix-collision handling live in the generators, not the runtime schema.

## References

- [CP-8: Compilation Database ADR](https://linear.app/ethankim8683/issue/CP-8/compilation-database-adr)
- [GitHub #5](https://github.com/EthanKim8683/cpx/issues/5)
