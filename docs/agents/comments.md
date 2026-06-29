# Commenting Guidelines

Semantic commenting rules and checklists for `cpx`. Style is enforced by [golangci-lint](go.md#verification-cheatsheet).

## Authoritative References

* [Go Doc Comments Specification](https://go.dev/doc/comment)
* [Effective Go — Commentary](https://go.dev/doc/effective_go#commentary)
* [Google Go Style Guide — Comments](https://google.github.io/styleguide/go/decisions#comment-sentences)

## cpx Rules

* **Document exported APIs**: Describe returns, side effects, concurrency, zero-value behavior, and special cases so callers do not read implementation code.
* **Non-trivial internal symbols**: Document complex unexported types and helpers in `internal/`.
* **Inline comments**: Explain *why*, not *what*. If a comment explains what code does, simplify the code.
* **No metadata**: No changelogs, author tags, dates, or commented-out code.

## Checklist

### Doc Comments
- [ ] Callers can learn returns, side effects, concurrency, zero-value behavior, and special cases.
- [ ] `bool` returns use **reports whether**.
- [ ] API documented, not implementation; does not restate the name.
- [ ] Non-trivial unexported types and functions in `internal/` documented.

### Inline Comments
- [ ] Explain why, not what.
- [ ] Intentional error discards and `//nolint` have a reason.
- [ ] No commented-out code, changelog, author, or date metadata.
