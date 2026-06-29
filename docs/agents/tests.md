# Testing Standards

What to test and how to assert it in `cpx`. Test mechanics enforced by [golangci-lint](go.md#verification-cheatsheet).

## Authoritative References

* **Principles**: [Go Test Comments](https://go.dev/wiki/TestComments) · [Table-Driven Tests](https://go.dev/wiki/TableDrivenTests)
* **Libraries**: [google/go-cmp](https://pkg.go.dev/github.com/google/go-cmp/cmp) · [goldie (Advanced Testing talk)](https://www.youtube.com/watch?v=yszygk1cpEc)
* **Concurrency**: [Go 1.24 synctest](https://go.dev/blog/synctest)

## Tool Selection Matrix

| Target | Tool / Pattern | Usage |
| --- | --- | --- |
| Structs, slices, maps | `google/go-cmp` | `cmp.Diff(want, got)`. Avoid `reflect.DeepEqual`. |
| Large generated output | `sebdah/goldie` | Fixtures in `testdata/`. Run `go test -update ./...` and review git diff. |
| Multi-valid outputs | Invariants | Assert semantic rules (e.g. topological order), not pinned ordering. |
| Concurrent / time code | `testing/synctest` | Use synthetic clock bubbles instead of real sleeping. |
| Thin wrappers | Smoke test | Smoke test only unless cpx adds logic on top. |

## Principles

* **Test cpx logic**: Parsing, transforms, selection, and errors this repo owns. Do not retest dependencies.
* **Structure**: `t.Run` subtests; descriptive case names; inputs in failure messages; `got` before `want`; `t.Helper()` in helpers.
* **Self-contained**: No network, external processes, or stray goroutines.

## Checklist

- [ ] Tests exercise cpx-owned logic, not third-party behavior.
- [ ] Structs, slices, maps compared with `go-cmp`, not `reflect.DeepEqual`.
- [ ] Large stable output uses `goldie` in `testdata/`.
- [ ] Table-driven with `t.Run` when cases share assertion logic.
- [ ] Failure messages include inputs; `got` before `want`; `t.Helper()` on helpers.
- [ ] Concurrent or time code uses `synctest` over `time.Sleep`.
