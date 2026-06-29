# Tests

Derived from [Test Comments](https://go.dev/wiki/TestComments), [Useful Test Failures](https://go.dev/wiki/CodeReviewComments#useful-test-failures), and [TableDrivenTests](https://go.dev/wiki/TableDrivenTests).

## Table-driven

- `t.Run` with descriptive names; include inputs in failure messages
- `got` before `want`: `t.Errorf("Foo(%q) = %v, want %v", in, got, want)`
- `t.Error` in table loops; `t.Helper()` in helpers
- No table-index-only failures; no exact `json.Marshal` bytes when semantics suffice

## go-cmp

[github.com/google/go-cmp](https://github.com/google/go-cmp) — compare whole structs, slices, maps with `cmp.Diff(want, got)`. Not for whole-file output. No `reflect.DeepEqual` or assertion libraries ([assert libraries](https://go.dev/wiki/TestComments#assert-libraries)).

## goldie

[github.com/sebdah/goldie](https://github.com/sebdah/goldie) — large whole-file expectations; `go test -update ./...`, review diffs. Not for small values or build inputs.

## cpx

- Fixtures in `testdata/` (committed)
- Runnable `Example` functions when demonstrating usage
