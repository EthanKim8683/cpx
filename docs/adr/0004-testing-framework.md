# ADR-0004: Testing and Assertion Framework

- **Status**: Accepted
- **Date**: 2026-06-30
- **Related**: [GitHub #55](https://github.com/EthanKim8683/cpx/issues/55)

---

## Context

Originally, the `cpx` codebase followed strict Go standard library comments discouraging assertion frameworks. We mandated using standard Go control flow (`if` statements, `t.Errorf`, `t.Fatalf`) and `github.com/google/go-cmp/cmp` for complex data assertions.

While idiomatic, this approach introduced significant developer friction:
- **High Boilerplate**: Basic assertions (like error verification, boolean checks, and empty strings) required multiple lines of control flow.
- **Duplicate Comparison Libraries**: Developers had to manage both `go-cmp` (for deep diffs) and standard comparisons, leading to inconsistent test designs.
- **Reduced Test Velocity**: Writing, reading, and maintaining tests took longer than necessary.

---

## Decision

We adopt **`github.com/stretchr/testify`** (`assert` and `require`) as the primary assertion framework for `cpx`.

### Rules of Engagement:
1.  **Framework Standard**: Use `assert` for non-fatal checks (allowing the test to continue executing other checks) and `require` for fatal assertions (such as setup failures or nil-checks before dereferencing).
2.  **Drop `go-cmp`**: Deprecate the use of `github.com/google/go-cmp/cmp`. Testify's `assert.Equal` and `assert.ElementsMatch` must be used for map, slice, and struct comparisons.
3.  **Strict Argument Ordering**: All Testify assertions must follow the signature `assert.Equal(t, expected, actual)` (expected/want before actual/got) to prevent misleading test failure reports.

---

## Alternatives

-  **Retain Standard Library + `go-cmp`**: (Rejected: high boilerplate and reduced developer velocity).
-  **Use Custom Assertion Helpers**: Writing custom, local check functions. (Rejected: duplicates work and lacks the rich error reporting/diffing provided by Testify).

---

## Consequences

-  **Increased Velocity**: Tests are cleaner, shorter, and faster to write.
-  **Unified Testing Signature**: Standardizes assertion layout and parameters across all test files.
-  **Dependency Simplification**: Keeps the testing dependency tree lean by relying on a single, feature-rich library (Testify) instead of multiple comparison tools.

---

## References

-   [Go Testing Guidelines](file:///Users/ethankim8683/Competitive%20Programming/Utilities/cpx/docs/agents/tests.md)
-   [GitHub Issue #55](https://github.com/EthanKim8683/cpx/issues/55)
