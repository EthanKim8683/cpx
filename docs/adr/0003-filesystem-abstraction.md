# ADR-0003: File System Abstraction

- **Status**: Accepted
- **Date**: 2026-06-30
- **Related**: [GitHub #55](https://github.com/EthanKim8683/cpx/issues/55)

---

## Context

`cpx` performs various file system operations, such as reading compiler configurations, generating directory skeletons, and saving build templates. Direct usage of the Go standard library `os` package binds functions to physical storage. This introduces several testing challenges:
- **Flaky & Slow Tests**: Tests performing physical disk I/O are slower, cannot safely run in parallel (risk of path collisions), and require tedious cleanup (`os.RemoveAll`, `defer`).
- **Platform-Specific Bugs**: File separator (`/` vs `\`) and line ending (`\n` vs `\r\n`) variations can lead to tests failing on Windows CI environments.
- **Lack of Mockability**: Hardcoded `os` calls are difficult to stub or mock.

---

## Decision

We adopt **`github.com/spf13/afero`** as the standard file system abstraction layer for the `cpx` repository.

### Rules of Engagement:
1.  **Production Functions**: Code that reads or writes to the file system must not make direct calls to the `os` package. Instead, it must accept an `afero.Fs` interface (or Go's standard `fs.FS` if read-only operations are sufficient) and perform operations through it.
2.  **Unit Testing**: Tests must substitute the physical filesystem with an in-memory implementation:
    -   Use `afero.NewMemMapFs()` for write/read workloads.
    -   Use the standard library `testing/fstest.MapFS` for read-only workloads.
3.  **Integration Testing**: Use `afero.NewOsFs()` when validating actual physical disk outputs is strictly required.
4.  **Method Call Ergonomics**: To maintain clean and readable method-based syntax, wrap `afero.Fs` in `&afero.Afero{Fs: fs}` inside functions and test verification routines. Avoid using Afero package-level function wrappers (e.g., prefer `afs.ReadFile(path)` over `afero.ReadFile(fs, path)`).

---

## Alternatives

-   **Direct OS interactions**: Hardcoding standard `os` calls. (Rejected: prevents parallel unit tests, makes mocking complex, and causes cross-platform flakiness).
-   **Custom Mock File System**: Defining our own interfaces and mocks. (Rejected: duplicates work already handled cleanly by `afero`'s rich ecosystems and helper functions).

---

## Consequences

-   **Test Hermeticity**: Unit tests run fully in-memory, requiring no disk cleanups and preventing file leaks.
-   **Execution Speed**: Test suites execute at CPU/RAM speeds.
-   **Dependency injection**: Production functions require a filesystem instance to be passed, increasing modularity and decoupling.

---

## References

-   [Go Testing Guidelines](../agents/tests.md)
-   [GitHub Issue #55](https://github.com/EthanKim8683/cpx/issues/55)
