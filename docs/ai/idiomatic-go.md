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

- Put test fixtures in `testdata/` (committed).
- Name directories after what they hold — e.g. `config/` for generated compiler option YAML, not a generic `data/`.
- Keep types and loading logic in sibling `.go` files at the package root (`config.go` loads from `config/*.yaml`); a `config/` directory holds files only, not a separate Go package.

**Do not**

- Use a vague `data/` folder when a specific name describes the contents.
- Create a `config/` **subpackage** (`package config`) when `config.go` in the parent package is enough — the folder is for artifacts, not code.
- Mix test fixtures and generated build inputs in the same directory.

## Code generation

cpx-specific; see [ADR-0001](../adr/0001-configuration.md) for env loading.

**Do**

- Wire generation with `//go:generate` (in the file it relates to, or `doc.go` for package-wide directives).
- Run external tools (`clang-tblgen`, AWK) inside `generate/*/main.go` under the parent package.
- Read generation inputs from direnv-loaded env vars.
- Write generated output to a specifically named directory (e.g. `internal/cdb/config/clang.yaml`).
- Gitignore generated build inputs; regenerate locally with `go generate ./...` before building or testing.
- Run `go generate` in CI before tests — same model as GCC/Clang, which do not commit generated option tables.
- Run `golangci-lint run ./...` in CI and locally before pushing Go changes.

**Do not**

- Commit large generated artifacts (compiler option configs, codegen output). They are reproducible from upstream sources and bloat the repo.
- Import `generate/` from the library package.

Makefile or CI may run `go generate ./...` — it does not run on `go build`.

**Committed vs generated**

| Artifact | Committed? | Example |
| --- | --- | --- |
| Golden test files | Yes — reviewable test expectations | `testdata/*.golden` |
| Generated build inputs | No — regenerate from source | `config/clang.yaml`, `config/gcc.yaml` |

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

- Large or complex **committed test expectations** (bundled source output) where the PR diff should be reviewable.

**Do not use for**

- Generated build inputs (compiler option configs) — those are gitignored and regenerated locally.
- Small values — use [go-cmp](#go-cmp) or inline checks.

### go-cmp

**Use for**

- Struct, slice, and map comparisons in tests.

**Do not use for**

- Whole-file byte output — use [goldie](#goldie).

### golangci-lint

Static analysis beyond `go vet`. See [golangci-lint](https://golangci-lint.run/).

**Use for**

- CI and local lint checks on all Go packages.
- Catching unchecked errors, dead code, and common mistakes — especially in codegen-heavy packages.

**Conventions**

- Configure in [`.golangci.yml`](../../.golangci.yml) at the repo root.
- Default linter set is `standard`; enable additional linters explicitly when needed.
- Run `golangci-lint run ./...`.

**Do not use for**

- Formatting — use `gofmt` / `goimports`.

## Tests

cpx-specific choices on top of [Code Review Comments](https://go.dev/wiki/CodeReviewComments#useful-test-failures) and [Examples](https://go.dev/wiki/CodeReviewComments#examples).

**Do**

- Use [goldie](#goldie) and `go test -update ./...` for golden files; review the diff before committing.
- Use [go-cmp](#go-cmp) for value comparisons.
- Run [golangci-lint](#golangci-lint) before pushing Go changes.

**Do not**

- Use assertion frameworks (e.g. `testify/assert`) when stdlib `t.Errorf` plus go-cmp suffices.
