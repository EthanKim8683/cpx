# Agent Guidelines: Go Error Handling

Errors in Go are first-class values and must be handled with deliberate care. AI agents and contributors in the `cpx` repository must follow these standards to ensure clear error propagation, debugging context, and codebase consistency.

---

## 1. Error Strings & Formatting

- **Lowercase and Unpunctuated**: Error strings returned by `errors.New()` or `fmt.Errorf()` must begin with a lowercase letter and must not end with punctuation (period, exclamation point, etc.).
  - **Correct**: `errors.New("failed to read configuration")`
  - **Incorrect**: `errors.New("Failed to read configuration.")`
- **Present Participle Phrasing**: Prefix error messages wrapping an underlying action with a present participle verb phrase detailing the operation being performed (e.g. `"fetching source"`, `"reading body"`).
  - **Correct**: `fmt.Errorf("detecting GCC version: %w", err)`
  - **Incorrect**: `fmt.Errorf("gcc version check failed: %w", err)`

---

## 2. Error Wrapping & Context (`%w` vs `%v`)

- **Focus on Action, Not Function Name**: Wrap errors with the higher-level operation being attempted at the call site, rather than stating that a specific function failed.
  - **Good (Context & Action)**:
    ```go
    cfg, err := readConfig()
    if err != nil {
        return nil, fmt.Errorf("loading configuration for user %q: %w", id, err)
    }
    ```
  - **Bad (Function-Name Stutter)**:
    ```go
    cfg, err := readConfig()
    if err != nil {
        return nil, fmt.Errorf("readConfig failed: %w", err)
    }
    ```
- **Wrap to Preserve Causality (`%w`)**: Use `fmt.Errorf("...: %w", err)` when the caller needs to inspect or programmatically handle the underlying root cause.
- **Opaque Formatting (`%v`)**: Use `%v` instead of `%w` only when you want to encapsulate internal implementation details and prevent public callers from depending on internal types.
- **Standard Library Inspection**: Always check wrapped errors using standard library functions:
  - Use `errors.Is(err, target)` to check for specific sentinel errors.
  - Use `errors.As(err, &target)` to unpack custom error types.

---

## 3. Sentinel Errors & Custom Types

- **Sentinel Errors**: Declare package-level sentinel errors with the `Err` prefix.
  ```go
  var ErrNotFound = errors.New("resource not found")
  ```
- **Custom Error Structs**: Use custom struct types (implementing the `error` interface) only when you need to carry extra structured context (such as status codes, file paths, or retryability metadata). Name custom types with an `Error` suffix (e.g. `PathError`).

---

## 4. Avoiding Swallowed or Double-Handled Errors

- **Never Swallow Errors**: Never ignore returned errors. Avoid discarding errors with blank identifiers (`_ = fn()`) unless explicitly documented with a justifying comment.
- **Handle Exactly Once**: An error should be handled (logged, counted, or returned) exactly once in the execution path. Do not log an error and return it simultaneously.
  - **Bad (Double-Handling)**:
    ```go
    log.Printf("failed to parse: %v", err)
    return err
    ```
  - **Good (Propagated)**:
    ```go
    return fmt.Errorf("parsing input: %w", err)
    ```
- **Execution Boundary**: Log or terminate execution (e.g., calling `os.Exit(1)`) only at the main execution boundary (such as the CLI `main()` entry point or HTTP handlers).

---

## Authoritative References

For deeper reading on Go error handling:
- [Go Blog: Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)
- [Effective Go: Errors](https://go.dev/doc/effective_go#errors)
- [Go Code Review Comments: Error Strings](https://github.com/golang/go/wiki/CodeReviewComments#error-strings)
