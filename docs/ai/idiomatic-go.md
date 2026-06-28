# Idiomatic Go

Repo-specific Go conventions for cpx. This is not a Go tutorial — see [Effective Go](https://go.dev/doc/effective_go), [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments), and [Organizing a Go module](https://go.dev/doc/modules/layout) for language and module basics.

Each topic has its own section below. When a new convention area emerges, add a section here following the same shape.

## Formatting

Mechanical style is handled by tools, not by hand.

**Do**

- Run `gofmt` or `goimports` on all Go code before committing.
- Let the formatter resolve indentation, alignment, and import grouping.

**Do not**

- Hand-format to work around `gofmt` — rearrange the code instead.

## Package layout

How to organize code within and across packages.

**Do**

- Put binaries in `cmd/` and non-exported library code in `internal/`.
- Keep multiple `.go` files in one package directory when they share one concern (`config.go`, `index.go` in `package cdb`).
- Create a subpackage (its own directory and import path) only when another package needs to import it as a distinct unit — e.g. `internal/port/`, Prometheus's `config/`.
- Split a large file into sibling files in the same package — not into subfolders that do not define a separate package.

**Do not**

- Nest folders only to mirror file names or JS/Python module trees (`config/schema.go`, `config/load.go` under one `package config` split across paths — use flat files instead).
- Use meaningless package names (`util`, `common`, `misc`, `api`, `types`, `interfaces`).

## Files and data

How to store artifacts alongside code.

**Do**

- Put test fixtures in `testdata/` (Go tooling convention).
- Put committed runtime or generated data in `data/` when the package loads it at runtime.
- Load embedded files with `//go:embed`.

**Do not**

- Put production fixtures in `testdata/` or test fixtures in `data/` — keep the distinction.

## Code generation

How to wire build-time code generation.

**Do**

- Use `//go:generate` directives in the file they relate to, or in `doc.go` when they apply to the whole package.
- Invoke external tools (`clang-tblgen`, AWK) inside `generate/*/main.go`.
- Rely on direnv-loaded env vars for generation inputs (see ADR-0001).
- Commit generated output when it is consumed at build time (YAML, `zz_generated.*.go`).

**Do not**

- Replace `go generate` with ad hoc shell scripts when a `go:generate` directive would suffice.

Makefile or CI targets that run `go generate ./...` are fine — `go generate` does not run automatically on `go build`.

## Build tools

How to structure generator programs.

**Do**

- Place generators as `package main` under `generate/` within the parent package (e.g. `internal/cdb/generate/clang`) or under repo-root `tools/` when shared across packages.
- Import the parent package for shared types.

**Do not**

- Import `generate/` or `tools/` from the library package (creates an import cycle).

## Types and API

How to define types and expose behavior.

**Do**

- Use typed constants (`type Kind int` with `iota`) over stringly-typed enums; prefix enum constants with the type name when it aids clarity.
- Attach behavior to types with methods (`func (c *CompilerConfig) Validate() error`).
- Use pointer receivers when methods mutate the receiver or the receiver is large; be consistent within a type.
- Define interfaces in the consumer package, not the implementor — return concrete types from constructors.
- Name single-method interfaces with the method name plus `-er` (`Reader`, `Stringer`) when the method has a canonical signature.

**Do not**

- Define interfaces on the implementor side "for mocking" — test through the concrete public API.
- Define interfaces before a realistic consumer exists.
- Use free functions when a method on the receiver type is natural (`validateConfig(cfg)` → `cfg.Validate()`).

## Naming

How to name packages, files, and identifiers.

**Do**

- Use short, lowercase, single-word package names without underscores (`cdb`, not `compdb` or `compilation_db`).
- Name files after their concern or behavior (`config.go`, `request.go`, `reload.go`) — lowercase, no hyphens.
- Use abbreviations that match the domain (`cdb` for compilation database, `oj` for online judge).
- Keep initialisms consistent in MixedCaps (`ServeHTTP`, `appID`, `XMLHTTPRequest`).
- Use receiver names that abbreviate the type (`c` for `Client`); keep them consistent across methods.
- Omit `Get` from getter names (`Owner()`, not `GetOwner()`).

**Do not**

- Repeat the package name in exported type names (`cdb.Database`, not `cdb.CDBDatabase`).
- Use `me`, `self`, or `this` as receiver names.

## Errors

How to create, wrap, and check errors.

**Do**

- Return errors from functions; handle them at the boundary (CLI, HTTP handler).
- Wrap errors with context using `%w`: `fmt.Errorf("load config: %w", err)`.
- Write error strings in lowercase with no trailing punctuation — they are usually printed after other context.
- Define sentinel errors as `var ErrFoo = errors.New("…")` for expected conditions callers should branch on.
- Define custom error types with an `Error` suffix when callers need typed fields; implement `Unwrap()` when wrapping is involved.
- Use `errors.Is` and `errors.As` to check wrapped errors.
- Indent error handling — handle errors early and return; avoid unnecessary `else` after error checks.

**Do not**

- Panic for expected error conditions.
- Discard errors with `_`.
- Drop the underlying error when wrapping operational failures — use `%w`, not `%v`.

## Tests

How to write and organize tests.

**Do**

- Write table-driven tests in `*_test.go`.
- Put fixtures in `testdata/`.
- Fail with messages that show input, got, and want: `t.Errorf("Foo(%q) = %d; want %d", in, got, want)`.
- Add runnable `Example` functions when introducing a new exported package API.

**Do not**

- Put test fixtures outside `testdata/` without good reason.

## Doc comments

How to document exported API.

**Do**

- Write doc comments for all exported names and non-trivial unexported declarations.
- Write full sentences starting with the name being described and ending with a period.

**Do not**

- Leave exported types or functions undocumented.
