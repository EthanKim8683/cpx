# Go Guidelines for cpx

Coding standards, references, and verification procedures for writing Go in `cpx`.

## Authoritative References

Official sources for standard Go conventions:

* **Project Layout**: [golang-standards/project-layout](https://github.com/golang-standards/project-layout) · [Organizing a Go Project](https://go.dev/doc/modules/layout)
* **Code Style**: [Effective Go](https://go.dev/doc/effective_go) · [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)

## Topic Guides

| Guide | Scope |
| --- | --- |
| [comments.md](comments.md) | Doc comments and inline commentary |
| [tests.md](tests.md) | Testing mechanics (`go-cmp`, `goldie`) |

## Verification Cheatsheet

Run before completing work:

```bash
go generate ./...
go test ./...
golangci-lint run ./...
```

* Linter configuration: [`.golangci.yml`](../../.golangci.yml). Use `gofmt` / `goimports` for formatting, not the linter.

## Checklist

- [ ] Package structure aligns with Go module conventions (`/internal`, `/cmd`).
- [ ] Code formatted using `gofmt` / `goimports`.
- [ ] [comments.md](comments.md) satisfied.
- [ ] [tests.md](tests.md) satisfied; golden diffs reviewed if updated.
- [ ] Generated build inputs not committed.
- [ ] Verification commands pass cleanly.
