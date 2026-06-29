# Comments

Semantic checklist. Style is enforced by [golangci-lint](go.md#golangci-lint).

Derived from [Go Doc Comments](https://go.dev/doc/comment) and [Effective Go — Commentary](https://go.dev/doc/effective_go#commentary).

## Checklist

### Doc comments

- [ ] Callers can learn returns, side effects, concurrency, zero-value behavior, deprecations, and special cases
- [ ] `bool` returns use **reports whether**
- [ ] API documented, not implementation
- [ ] Does not restate the name
- [ ] Non-trivial unexported types and functions documented

### Inline comments

- [ ] Explain why, not what
- [ ] Intentional error discards and `//nolint` have a reason
- [ ] No commented-out code, changelog, author, or date metadata
- [ ] Unclear names fixed instead of commented around

If a comment explains *what* the code does, simplify the code.
