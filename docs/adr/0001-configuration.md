# ADR-0001: Configuration

- **Status**: Accepted
- **Date**: 2026-06-28
- **Related**: [CP-9: Configuration ADR](https://linear.app/ethankim8683/issue/CP-9/configuration-adr), [GitHub #6](https://github.com/EthanKim8683/cpx/issues/6)

---

## Context

cpx integrates Go, Make, AWK, go generate, and other tooling. The original implementation never had a single source of truth for configuration shared across these tools. `.env` was the intended source, but loading and sharing it never felt clean — especially for Make, which cannot `include` a `.env` file directly and has its own variable system on top of environment variables.

---

## Decision

### Source and Loader
- **`.env`**: The human-editable configuration source (gitignored).
- **direnv**: Loads configuration via `.envrc` (`source_up_if_exists` + `dotenv`). This serves as the single loader, eliminating per-tool `.env` parsing.

### Runtime API
Environment variables serve as the shared configuration API across all tools:
- **Go**: Loads configuration into an `internal/config` struct from environment variables at CLI startup.
- **Make**: Imports environment variables as Make variables, referenced as `$(VAR)` in recipes.
- **AWK**: Receives variables explicitly in recipes, e.g., `awk -v VAR=$(VAR)`.
- **go generate**: Inherits environment variables directly from the invoking shell or Make process.

### Make Caveat
Makefile assignments (`VAR := foo`) override imported environment variables. An optional `include config.mk` can provide Make-native defaults when direnv is not in use.

---

## Alternatives

- **Per-tool `.env` parsing**: Each of Go, Make, and AWK reads `.env` independently. (Rejected: duplicates loading logic across languages).
- **Generated bridge artifacts**: A tool generates `config.mk` and other files from a single manifest. (Discussed: not chosen as the primary path once direnv handles loading).
- **`include config.mk` as primary config**: Make-native variables without direnv. (Not chosen as the primary path: kept only as an optional fallback).

---

## Consequences

- **direnv dependency**: The primary development workflow assumes direnv is installed and allowed.
- **Make variable shadowing**: Makefile assignments override imported env vars; easy to accidentally break configuration.
- **AWK requires explicit passing**: AWK scripts do not automatically inherit Make variables; recipes must pass variables via `-v` or `export`.

---

## References

- [CP-9: Configuration ADR](https://linear.app/ethankim8683/issue/CP-9/configuration-adr)
- [GitHub Issue #6](https://github.com/EthanKim8683/cpx/issues/6)
