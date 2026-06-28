# Idiomatic Go

Read the sources below first. They are the authority on idiomatic Go — this doc only records **what cpx adds** on top of them. Do not duplicate their guidance here; link out instead.

When a new cpx-specific convention emerges, add a section below following the same shape.

## Sources

| Source | Covers |
| --- | --- |
| [Effective Go](https://go.dev/doc/effective_go) | Language idioms, naming, formatting, errors, methods, interfaces |
| [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) | Practical review checklist — formatting, errors, interfaces, receivers, doc comments, tests |
| [Organizing a Go module](https://go.dev/doc/modules/layout) | `cmd/`, `internal/`, package boundaries |
| [Go Wiki: TableDrivenTests](https://go.dev/wiki/TableDrivenTests) | Table-driven test structure |
| [google/go-cmp](https://github.com/google/go-cmp) | Semantic equality and diffs in tests (Go wiki–recommended) |
| [Advanced Testing with Go](https://www.youtube.com/watch?v=yszygk1cpEc) (Mitchell Hashimoto) | Golden file testing, integration testing patterns |

## Package layout

cpx-specific choices on top of [Organizing a Go module](https://go.dev/doc/modules/layout) and [Package Names](https://go.dev/wiki/CodeReviewComments#package-names).

**Do**

- Keep related code as sibling `.go` files in one package (`config.go`, `index.go` in `package cdb`) rather than a subfolder that mirrors file names.
- Add a subpackage directory only when another package imports it as a distinct unit (`internal/port/`, `internal/cdb/generate/clang`).

**Do not**

- Mirror JS/Python folder-per-concept layouts (`config/schema.go`, `config/load.go` under one package path).

## Files and data

cpx-specific split; see [Code Review Comments](https://go.dev/wiki/CodeReviewComments) and the `testing` package for general fixture conventions.

**Do**

- Put test fixtures in `testdata/`.
- Put committed runtime or generated data in `data/` within the package that loads it.

**Do not**

- Mix test fixtures and runtime data in the same directory.

## Code generation

cpx-specific; see [ADR-0001](../adr/0001-configuration.md) for env loading.

**Do**

- Wire generation with `//go:generate` (in the file it relates to, or `doc.go` for package-wide directives).
- Run external tools (`clang-tblgen`, AWK) inside `generate/*/main.go` under the parent package.
- Read generation inputs from direnv-loaded env vars.
- Commit generated artifacts consumed at build time (YAML, `zz_generated.*.go`).

**Do not**

- Import `generate/` from the library package.

Makefile or CI may run `go generate ./...` — it does not run on `go build`.

## Naming

Domain abbreviations on top of [Effective Go § Names](https://go.dev/doc/effective_go#names) and [Code Review Comments](https://go.dev/wiki/CodeReviewComments#initialisms).

| Abbreviation | Meaning |
| --- | --- |
| `cdb` | compilation database |
| `oj` | online judge |

## Productivity libraries

Tools cpx adopts to encode idioms. Read each library's own docs for API details — only cpx usage choices are recorded here. Add an entry when cpx adopts a new tool.

### goldie

Golden file testing — pattern from [Advanced Testing with Go](https://www.youtube.com/watch?v=yszygk1cpEc). See [goldie](https://github.com/sebdah/goldie).

**Use for**

- Large or complex expected output (bundled source, JSON, generated YAML) where the PR diff should be reviewable.

**Do not use for**

- Small values — use [go-cmp](#go-cmp) or inline checks.

### go-cmp

**Use for**

- Struct, slice, and map comparisons in tests.

**Do not use for**

- Whole-file byte output — use [goldie](#goldie).

### golangci-lint

**Use for**

- CI and local checks beyond `go vet`.

Configure in `.golangci.yml`; keep enabled linters minimal.

## Tests

cpx-specific choices on top of [Code Review Comments](https://go.dev/wiki/CodeReviewComments#useful-test-failures) and [Examples](https://go.dev/wiki/CodeReviewComments#examples).

**Do**

- Use [goldie](#goldie) and `go test -update ./...` for golden files; review the diff before committing.
- Use [go-cmp](#go-cmp) for value comparisons.

**Do not**

- Use assertion frameworks (e.g. `testify/assert`) when stdlib `t.Errorf` plus go-cmp suffices.
