# Idiomatic Go

Repo-specific Go conventions for cpx. This is not a Go tutorial — see [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) for language basics.

Each topic has its own section below. When a new convention area emerges, add a section here following the same shape.

## Package layout

How to organize code within and across packages.

**Do**

- Keep related code in one flat package (`package cdb`) when it shares one concern and has no separate import boundary.
- Split a large file into sibling files in the same package (`config.go`, `index.go`) — not into subfolders.
- Create a subpackage only when another package needs to import it as a distinct unit.

**Do not**

- Mirror JS/Python folder-per-concept layouts (`config/schema.go`, `config/load.go` under `config/`).
- Create subpackages that exist only to match file names.

## Files and data

How to store committed artifacts alongside code.

**Do**

- Put generated or committed data files in a `data/` directory within the package.
- Load embedded data with `//go:embed` at runtime.

**Do not**

- Use a directory name that implies a separate Go package (e.g. `config/data/` under a `config/` package).

## Code generation

How to wire build-time code generation.

**Do**

- Use `//go:generate` directives in `doc.go`.
- Invoke external tools (`clang-tblgen`, AWK) inside `generate/*/main.go` — not from shell scripts or Makefiles.
- Rely on direnv-loaded env vars for generation inputs (see ADR-0001).

**Do not**

- Add Makefile targets for Go-native codegen when `go generate` suffices.

## Build tools

How to structure generator programs.

**Do**

- Place generators as `package main` under `generate/` within the parent package (e.g. `internal/cdb/generate/clang`).
- Import the parent package for shared types.

**Do not**

- Import `generate/` from the parent package (creates an import cycle).

## Types and API

How to define types and expose behavior.

**Do**

- Use typed constants (`type Kind int` with `iota`) over stringly-typed enums.
- Attach behavior to types with methods (`func (c *CompilerConfig) Validate() error`).
- Write table-driven tests in `*_test.go`.

**Do not**

- Use free functions when a method on the receiver type is natural (`validateConfig(cfg)` → `cfg.Validate()`).

## Naming

How to name packages, files, and identifiers.

**Do**

- Use short, lowercase package names without underscores (`cdb`, not `compdb` or `compilation_db`).
- Name files after their primary type or concern (`config.go`, `database.go`) — lowercase, no hyphens.
- Use abbreviations that match the domain (`cdb` for compilation database, `oj` for online judge).

**Do not**

- Repeat the package name in exported type names (`cdb.Database`, not `cdb.CDBDatabase`).

## Errors

How to create, wrap, and check errors.

**Do**

- Return errors from functions; handle them at the boundary (CLI, HTTP handler).
- Wrap errors with context using `%w`: `fmt.Errorf("load config: %w", err)`.
- Define sentinel errors as package-level variables for expected conditions callers should branch on: `var ErrNotFound = errors.New("not found")`.
- Use `errors.Is` and `errors.As` to check wrapped errors.

**Do not**

- Panic for expected error conditions.
- Drop the underlying error when wrapping — always use `%w`, not `%v`.
