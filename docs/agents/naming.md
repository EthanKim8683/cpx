# Agent Guidelines: Go Naming Conventions

A core philosophy of the Go language is uniformity: **all Go code across a project should look as if it were written by the exact same developer**. Consistent naming enables contributors (and AI agents) to navigate, understand, and build upon the codebase without cognitive friction.

This document outlines the required Go naming conventions for the `cpx` repository.

---

## 1. Avoid Stuttering & Package Redundancy

Do not repeat package or directory names within variables, structs, or functions. Since Go packages serve as namespaces, the package name qualifies the identifier at the call site.

*   **Bad (Stuttering)**:
    *   Inside package `gccoptgen`: `gccSourceFiles`, `DetectGCCVersion`, `BuildGCCRawURL`
    *   Call site: `gccoptgen.DetectGCCVersion()`
    *   Inside package `config`: `ConfigStruct`, `LoadConfig()`
*   **Good (Context-Aware)**:
    *   Inside package `gccoptgen`: `sourceFiles`, `detectVersion`, `rawURL`
    *   Call site: `gccoptgen.DetectVersion()`
    *   Inside package `config`: `Config`, `Load()`
    *   Call site: `config.Load()`

---

## 2. Scope-Proportional Variable Length

The length of a variable name should be proportional to the size of its lexical scope.
*   **Local Scopes**: In small blocks (like loops or short functions of 10-20 lines), use short, concise names. Extensive descriptions add visual noise where context is already obvious.
    *   **Good**: `b` (byte slice), `path` (file path), `u` (URL), `err` (error).
    *   **Bad**: `bodyBytes`, `relativeFilePath`, `targetDownloadURL`.
*   **Global / Exported Scopes**: Use descriptive names for package-level constants, variables, types, and exported functions.
*   **Scope Progression**: As the distance between a variable's declaration and its usage increases, the name must become more descriptive to maintain clarity.

---

## 3. Visibility and Exporting in Commands

For executable command tools (`package main` files under `internal/cdb/cmd/`):
*   Default all helper variables, constants, and functions to **unexported** (lowercase).
*   Only export symbols if they are explicitly shared across multiple files or test files in the package.

---

## 4. Go File Naming Guidelines

*   **Concise & Lowercase**: Filenames must be short, lowercase, and direct.
*   **Word Concatenation**: Avoid arbitrary punctuation or formatting. Concatenate multiple words directly (e.g., `reverseproxy.go`).
*   **Avoid Arbitrary Underscores (`_`)**: Do not use underscores to separate words in standard filenames. The Go compiler parses suffixes after underscores for build constraints (e.g., `_linux.go`, `_amd64.go`) and tests (`_test.go`).
*   **Case Sensitivity & Cross-Platform Safety**: Always name files in lowercase (e.g., `env.go` instead of `ENV.go` or `Env.go`). Mixed-case or uppercase filenames cause conflicts on case-insensitive filesystems (macOS, Windows) when ported or compiled on case-sensitive filesystems (Linux).

---

## 5. Idiomatic Verb Selection for Functions and Methods

When naming functions or methods that perform actions, use simple, direct, and standard verbs. Avoid overly verbose naming patterns or un-idiomatic verbs.

### Recommended Standard Verbs
Prefer standard library verbs to represent actions:
*   **`Read` / `Write`**: For input/output operations (e.g., `io.Reader.Read`, `io.Writer.Write`).
*   **`Load`**: For retrieving configurations, states, or external assets (e.g., `config.Load()`).
*   **`Get`**: Use **only** for map-like lookups (e.g., `Header.Get("Key")`) or network fetches (e.g., `http.Get()`). Never use `Get` for struct field accessors (see below).
*   **`Prepare`**: Setup or compile resources (e.g., `sql.Stmt.Prepare()`).
*   **`Parse`**: For parsing raw input into structured data (e.g., `time.Parse()`, `url.Parse()`).
*   **`Create` / `Open` / `Close`**: For managing resource lifecycles (e.g., `os.Create()`, `os.Open()`, `Close()`).
*   **`Run`**: For executing processes, tasks, or long-running workers (e.g., `exec.Cmd.Run()`).

### Discouraged and Verbose Naming
*   **Avoid Low-Signal Action Verbs**: Do not use low-signal, catch-all verbs (such as `populate`, `process`, `handle`, or `manage`) for functions or methods when a more precise verb or idiomatic Go pattern exists (e.g., constructors, parsers, or direct assignments).
*   **Avoid Verbose Helper Verbs**: Do not prefix functions with helper verbs like `Calculate`, `Compute`, or `Find` when the noun or property alone is sufficient.

| Bad | Good |
| :--- | :--- |
| `CalculateArea()` | `Area()` |
| `FindUser()` | `User()` |
| `ProcessData()` | `ParseData()` / `WriteData()` |

### When to Omit Verbs
Go commonly omits verbs to keep the API surface concise and noun-focused:
*   **Getter Methods**: Omit the `Get` prefix entirely. Use `user.Email()` instead of `user.GetEmail()`. Setters should still use `Set` (e.g., `user.SetEmail()`).
*   **Properties and Calculations**: For simple properties or computed states, use the noun. E.g., `buf.Len()` instead of `buf.GetLength()`, `rect.Area()` instead of `rect.CalculateArea()`.
*   **Type Conversions**: Omit helper verbs like `To` or `As`. E.g., `String()` instead of `ToString()`, `Bytes()` instead of `ToBytes()`.
*   **Constructor Context**: Omit verbs like `Create` or `New` when the package name itself provides clear context. E.g., `errors.New()` instead of `errors.CreateError()`.

---

## 6. Standard Variable Identifier Lookup Table

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

### Context-Specific & Reserved Identifiers
Certain identifiers carry a very strong semantic meaning in Go. Do not use them for other types or contexts to avoid shadowing and confusion:
*   **`ctx`**: Reserved exclusively for `context.Context`. Do not name general/custom contexts (like AST context or canvas rendering context) as `ctx` (use `astCtx` or `renderCtx`). For web framework contexts (e.g., `gin.Context`), use `c`.
*   **`err`**: Reserved exclusively for `error`. Do not name lists of error structures or other variables as `err`.
*   **`t` and `b`**: Reserved exclusively for `*testing.T` and `*testing.B` in tests/benchmarks. Do not use them for loop counters, time variables, or other values within test files.
*   **`mu`**: Reserved exclusively for mutexes (`sync.Mutex` or `sync.RWMutex`).
*   **`tx`**: Reserved exclusively for database transactions.
*   **`ch`**: Reserved for Go channels (`chan T`).

### Avoid Package Shadowing
Do not name local variables the same as common imported packages (like `env`, `url`, `path`, `log`). Doing so shadows the package import, making it inaccessible within the function scope.

```go
// Bad: local variables shadow the 'log' and 'env' packages
func process(env string) {
    log := logger.New()
    log.Printf("processing in environment: %s", env)
}

// Good: local variables use distinct, descriptive names
func process(envName string) {
    appLogger.Printf("processing in environment: %s", envName)
}
```

---

## Authoritative References

For deeper reading on idiomatic Go style:
- [Effective Go: Naming](https://go.dev/doc/effective_go#names)
- [Go Code Review Comments: Variable Names](https://github.com/golang/go/wiki/CodeReviewComments#variable-names)
- [Uber Go Style Guide: Naming](https://github.com/uber-go/guide/blob/master/style.md#naming)
