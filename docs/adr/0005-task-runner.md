# ADR-0005: Task Runner

- **Status**: Accepted
- **Date**: 2026-06-28
- **Related**: [CP-7](https://linear.app/ethankim8683/issue/CP-7/task-runner-adr), [GitHub #4](https://github.com/EthanKim8683/cpx/issues/4)

## Context

cpx integrates Go, Make, AWK, and other tools. Without a central task runner, each workflow is invoked differently. The goal is a single `cpx <command>` interface that orchestrates the full pipeline.

## Decision

`cmd/cpx` CLI (currently `main.go` at root, moving to `cmd/cpx/main.go`):

- **`cpx all`**: Run the full pipeline — config load → compdb → clean → generate → bundle → validate.
- **`cpx bundle`**: Skip validation, output only.
- **Individual commands**: `cpx generate`, `cpx clean`, `cpx config`, `cpx compdb`.
- **Execution**: Direct function calls. No goroutines — sequential pipeline, deterministic output.
- **Flag parsing**: Minimal. Use Go stdlib `flag` package for now.
- **AWK integration**: Shell out via `exec.Command("awk", ...)`. AWK scripts live in `scripts/`.

## Alternatives considered

- **Task**: Requires Ruby/Node runtime. Unnecessary dependency.
- **Make-only**: Make is already a dependency for other reasons, but lacks Go-native error handling and struct-based config.
- **`go run main.go`**: Works during development but not distributable.

## Consequences

- **Single binary**: `cpx` is one binary, one entry point.
- **No parallel execution**: Sequential by design. If needed later, can add goroutines with `errgroup`.
- **Make stays**: Makefile remains for non-Go workflows (AWK, shell scripts) and as a secondary interface.

## References

- [GitHub #4](https://github.com/EthanKim8683/cpx/issues/4) — Architecture Pattern ADR discussion

