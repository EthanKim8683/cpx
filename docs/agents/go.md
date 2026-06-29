# Go

How to write Go in cpx.

1. Follow [topic docs](#topic-docs) for cpx rules.
2. Run the [checklist](#checklist) before finishing.

Official docs are authoritative — these pages record what cpx adds. Package structure and naming follow issue scope.

## Topic docs

| Doc | Covers |
| --- | --- |
| [comments.md](comments.md) | Comment content |
| [tests.md](tests.md) | Tests, go-cmp, goldie |

## Checklist

Skip items that do not apply.

- [ ] Change aligns with issue scope
- [ ] [comments.md](comments.md) satisfied
- [ ] [tests.md](tests.md) satisfied; golden diffs reviewed if updated
- [ ] Generated build inputs not committed
- [ ] `go generate ./...` · `go test ./...` · `golangci-lint run ./...`

## golangci-lint

[golangci-lint.run](https://golangci-lint.run/) · [`.golangci.yml`](../../.golangci.yml). Use `gofmt` / `goimports` for formatting, not the linter.

Comment **style** — revive `exported` + `package-comments` (comments exclusion preset omitted). Comment **content** — [comments.md](comments.md). Tests — [tests.md](tests.md).
