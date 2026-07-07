# ADR-0004: Testing Frameworks

**Status:** Accepted

## Context

cpx needs a standard testing approach. Current code uses `testing` stdlib with 83 `t.Fatalf`, 59 `t.Errorf`, and 52 `fmt.Printf` across 11 test files. This gives Go-native error reporting and stack traces but verbose assertion code — comparing structs takes 3+ lines vs one assertion call.

## Decision

Use `testify` with `assert` and `require`:

- **`assert` for non-fatal checks**: Continue test execution after failure. Use for independent assertions.
- **`require` for fatal checks**: Stop test immediately when failure makes remaining assertions meaningless.
- **Argument ordering**: `assert.Equal(t, expected, actual)` — expected first, actual second.
- **Golden files**: Use [goldie](https://github.com/sebdah/goldie) for complex output comparisons. Store golden files in `testdata/` directories. Update with `go test -update ./...`.
- **Parallel tests**: Use `t.Parallel()` in subtests where possible. Avoid shared mutable state between parallel tests.
- **No go-cmp**: Use testify for all comparisons. Do not add `go-cmp` as a dependency.
- **Testify mocking**: Use `mock.Mock` when tests need controlled return values or call verification. Prefer interface-based fakes for simple cases.

## Alternatives considered

- **`gocheck`**: Suite-based, YAML-driven tests. Heavier framework, less idiomatic for modern Go.
- **`gomega`**: BDD-style matchers. More verbose than testify for simple comparisons.
- **`go-cmp`**: Deep diff comparisons. Excellent for complex structs but adds a dependency for a single use case.

## Consequences

- Error message style: present participle (`"loading config: %w"`, not `"failed to load config: %w"`).
- Filesystem testing: use `afero.NewMemMapFs()` for read/write, `testing/fstest.MapFS` for read-only.
- One error handler: errors are logged or returned, never both.

## References

- [GitHub #55](https://github.com/EthanKim8683/cpx/issues/55)
