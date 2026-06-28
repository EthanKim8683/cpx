# ADR-0001: Configuration

- **Status:** Accepted
- **Date:** 2026-06-28
- **Related:** [CP-9: Configuration ADR](https://linear.app/ethankim8683/issue/CP-9/configuration-adr), [GitHub #6](https://github.com/EthanKim8683/cpx/issues/6)

## Context

CPX integrates Go, Make, AWK, go generate, and other tooling. The original implementation never had a single source of truth for configuration shared across these tools. `.env` was the intended source, but loading and sharing it never felt clean — especially for Make, which cannot `include` a `.env` file directly and has its own variable system on top of environment variables.

## Decision

### Source and loader

- **`.env`** is the human-editable configuration source (gitignored).
- **direnv** loads it via `.envrc` (`source_env .env`). This is the single loader — no per-tool `.env` parsing.

### Runtime API

Environment variables are the shared configuration API across all tools.

| Tool | How it reads config |
|------|---------------------|
| **Go** | `internal/config` struct loaded from env at CLI startup |
| **Make** | Environment variables imported as Make variables; use `$(VAR)` in recipes |
| **AWK** | Vars passed in recipes, e.g. `awk -v VAR=$(VAR)` |
| **go generate** | Inherits env from the invoking shell or Make |

### Make caveat

Makefile assignments (`VAR := foo`) override imported environment variables. An optional `include config.mk` can provide Make-native defaults when direnv is not in use.

## References

- [CP-9: Configuration ADR](https://linear.app/ethankim8683/issue/CP-9/configuration-adr)
- [GitHub #6](https://github.com/EthanKim8683/cpx/issues/6)
