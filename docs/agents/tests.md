# Tests

What to test and how to assert it. Test mechanics enforced by [golangci-lint](#golangci-lint).

Derived from [Go Test Comments](https://go.dev/wiki/TestComments), [Table-Driven Tests](https://go.dev/wiki/TableDrivenTests), [go-cmp](https://github.com/google/go-cmp), and [golden-file testing](https://www.youtube.com/watch?v=yszygk1cpEc).

## Tools

### go-cmp

[github.com/google/go-cmp](https://github.com/google/go-cmp) — `cmp.Diff(want, got)` for structs, slices, and maps. Clearer failures than `reflect.DeepEqual`; matches [Go wiki guidance against assertion libraries](https://go.dev/wiki/TestComments#assert-libraries).

```go
if diff := cmp.Diff(want, got); diff != "" {
    t.Errorf("Foo() mismatch (-want +got):\n%s", diff)
}
```

### goldie

[github.com/sebdah/goldie](https://github.com/sebdah/goldie) — golden-file assertions in `testdata/`.

Adopted after [Mitchell Hashimoto — Advanced Testing with Go](https://www.youtube.com/watch?v=yszygk1cpEc): store expected output, compare on run, regenerate with `-update`. Use for large generated output (YAML, JSON, rendered commands); run `go test -update ./...` and review the git diff. Not for small values or build inputs.

## Principles

**Test cpx logic** — parsing, transforms, selection, and errors this repo owns.

**Do not retest dependencies.** A thin `env.ParseAs` wrapper needs one smoke test, not a per-field table that only echoes env vars back (`config.Load` today).

**Prefer semantics over exact output** — use go-cmp for meaning, not raw bytes. When many answers are valid, assert invariants (topological sort: every `u → v` has `u` before `v`), not one pinned ordering.

**Table-driven when it fits** — same assertion logic, varying inputs ([Go wiki](https://go.dev/wiki/TableDrivenTests)). Skip tables when cases need different logic or a single direct test is clearer.

**Keep unit tests self-contained** — no network, external processes, or stray goroutines ([`testing/synctest`](https://pkg.go.dev/testing/synctest) package docs). For concurrent or time-based code, prefer `synctest` over real sleeps ([VictoriaMetrics](https://victoriametrics.com/blog/go-synctest/), [Go blog](https://go.dev/blog/synctest)).

**Structure** — `t.Run` subtests; descriptive case names; inputs in failure messages; `got` before `want`; `t.Helper()` in helpers; fixtures in `testdata/`.

## Checklist

- [ ] Tests exercise cpx-owned logic, not third-party behavior
- [ ] Thin wrappers: smoke test only unless cpx adds logic on top
- [ ] Table-driven with `t.Run` when cases share assertion logic
- [ ] Failure messages include inputs; `got` before `want`; `t.Helper()` on helpers
- [ ] Structs, slices, maps compared with go-cmp, not `reflect.DeepEqual`
- [ ] Multiple valid outputs checked via invariants, not one canonical ordering
- [ ] Large stable output uses goldie; `-update` diffs reviewed
- [ ] Fixtures in `testdata/`; generated build inputs are not golden files

## golangci-lint

On `*_test.go` ([`.golangci.yml`](../../.golangci.yml)): blocks testify and `reflect.DeepEqual`; requires `t.Helper()` in helpers and expected output on `Example` functions. Does not catch over-testing delegated logic — checklist above covers that.
