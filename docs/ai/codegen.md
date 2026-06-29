# Code generation

Derived from [`go generate`](https://go.dev/blog/generate).

Does not run on `go build`. Makefile or CI may run `go generate ./...`.

**Do**

- `//go:generate` in the related file or `doc.go`
- Read inputs from direnv-loaded env
- Write output to named directories; gitignore generated build inputs

**Do not**

- Commit generated build inputs (e.g. `config/clang.yaml`)
- Import generator `main` from library code

## Generator shape

| Need | Use |
| --- | --- |
| Custom transform | `go run` on in-repo `main` ([in-source generators](https://eli.thegreenplace.net/2021/a-comprehensive-guide-to-go-generate/)) |
| Upstream CLI | invoke directly (`clang-tblgen`, AWK) |
| Cleanup or multi-step | shell one-liners or Makefile |

Pin module tools in `go.mod` ([tool dependencies](https://go.dev/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module)). Shell for cleanup only — not core logic ([`go generate` proposal](https://go.googlesource.com/proposal/+/master/design/go-generate.md)).
