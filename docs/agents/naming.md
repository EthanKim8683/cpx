# Agent Guidelines: Go Naming Conventions

A core philosophy of the Go language is uniformity: **all Go code across a project should look as if it were written by the exact same developer**. Consistent naming enables contributors (and AI agents) to navigate, understand, and build upon the codebase without cognitive friction.

This document outlines the required Go naming conventions for the `cpx` repository.

---

## 1. Avoid Stuttering & Package Redundancy

Do not repeat package or directory names within variables, structs, or functions. Since Go packages serve as namespaces, the package name qualifies the identifier at the call site.

- **Bad (Stuttering)**:
  - Inside package `gccoptgen`: `gccSourceFiles`, `DetectGCCVersion`, `BuildGCCRawURL`
  - Call site: `gccoptgen.DetectGCCVersion()`
- **Good (Context-Aware)**:
  - Inside package `gccoptgen`: `sourceFiles`, `detectVersion`, `rawURL`
  - Call site: `gccoptgen.DetectVersion()` (if exported)

---

## 2. Scope-Proportional Variable Length

The length of a variable name should be proportional to the size of its lexical scope.
- **Local Scopes**: In small blocks (like loops or short functions of 10-20 lines), use short, concise names. Extensive descriptions add visual noise where context is already obvious.
  - **Good**: `b` (byte slice), `path` (file path), `u` (URL), `err` (error).
  - **Bad**: `bodyBytes`, `relativeFilePath`, `targetDownloadURL`.
- **Global / Exported Scopes**: Use descriptive names for package-level constants, variables, types, and exported functions.

---

## 3. Visibility and Exporting in Commands

For executable command tools (`package main` files under `internal/cdb/cmd/`):
- Default all helper variables, constants, and functions to **unexported** (lowercase).
- Only export symbols if they are explicitly shared across multiple files or test files in the package.

---

## 4. Go File Naming Guidelines

- **Concise & Lowercase**: Filenames must be short, lowercase, and direct.
- **Word Concatenation**: Avoid arbitrary punctuation or formatting. Concatenate multiple words directly (e.g., `reverseproxy.go`).
- **Avoid Arbitrary Underscores (`_`)**: Do not use underscores to separate words in standard filenames. The Go compiler parses suffixes after underscores for build constraints (e.g., `_linux.go`, `_amd64.go`) and tests (`_test.go`). Using arbitrary underscores can lead to unintended compilation issues or file exclusion.

---

## 5. Standard Variable Identifier Lookup Table

Use these standard Go variable names for common types and contexts:

| Type / Context | Standard Identifier | Example Usage |
| :--- | :--- | :--- |
| `context.Context` | `ctx` | `func DoSomething(ctx context.Context)` |
| `error` | `err` | `if err != nil { return err }` |
| `io.Reader` / `*http.Request` | `r` | `r.Read(p)` / `r.URL` |
| `io.Writer` / `http.ResponseWriter` | `w` | `w.Write(p)` |
| `sync.Mutex` | `mu` | `mu.Lock()` |
| Configuration struct | `cfg` | `cfg := config.Load()` |
| Option slice / struct | `opts` | `opts := []Option{...}` |
| Buffer / byte slice | `b` or `buf` | `b, err := io.ReadAll(r)` |
| Method receiver | 1-2 letters matching type | `func (c *Client) Connect()` |

---

## 6. Authoritative References

For deeper reading on idiomatic Go style:
- [Effective Go: Naming](https://go.dev/doc/effective_go#names)
- [Go Code Review Comments: Variable Names](https://github.com/golang/go/wiki/CodeReviewComments#variable-names)
- [Uber Go Style Guide: Naming](https://github.com/uber-go/guide/blob/master/style.md#naming)
