# Testing

- **Assertions**: Use [testify](https://github.com/stretchr/testify) (`assert` for non-fatal, `require` for fatal).
- **Argument ordering**: `assert.Equal(t, expected, actual)`.
- **Golden files**: Use [goldie](https://github.com/sebdah/goldie) for large/complex outputs. Update with `go test -update ./...`.
- **Parallel**: Use `t.Parallel()` in subtests where possible.
- **No go-cmp**: Use testify for all comparisons.

See [ADR-0004](../adr/0004-testing-frameworks.md) for the full decision history.
