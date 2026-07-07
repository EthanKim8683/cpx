# ADR-0001: Configuration

**Status:** Accepted

## Context

cpx integrates Go, Make, AWK, and other tooling. The original implementation had no single source of truth for configuration. `.env` was the intended source, but loading it per-tool was messy — especially for Make, which cannot `include` a `.env` file directly.

## Decision

- **Source**: `.env` (gitignored, human-editable).
- **Loader**: direnv via `.envrc` (`source_up_if_exists` + `dotenv`). Single loader, no per-tool parsing.
- **Runtime API**: Environment variables are the shared API. Go loads them into `internal/config`. Make imports them as `$(VAR)`. AWK receives them via `-v`. `go generate` inherits them from the shell.
- **Make caveat**: Makefile assignments (`VAR := foo`) override imported env vars. An optional `include config.mk` provides Make-native defaults when direnv is not in use.

## Alternatives considered

- **Per-tool `.env` parsing**: (Rejected: duplicates loading logic across languages).
- **Generated bridge artifacts**: (Rejected: not chosen once direnv handles loading).
- **`include config.mk` as primary config**: (Rejected: kept only as optional fallback).

## Consequences

- direnv is required for the primary development workflow.
- Make variable shadowing can accidentally break configuration.
- AWK scripts must receive variables explicitly via `-v` or `export`.

## References

- [GitHub #6](https://github.com/EthanKim8683/cpx/issues/6)
