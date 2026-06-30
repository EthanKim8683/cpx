# Agent Guidelines: Go Commenting Standards

Comments in Go are crucial to understanding design intent. AI agents and contributors in the `cpx` repository must follow these standards to ensure clear package contracts, design rationale documentation, and codebase readability while preventing visual noise.

This document outlines the required Go commenting conventions for the `cpx` repository.

---

## 1. Doc Comments (Black-Box Contract / "What")

Doc comments on packages, types, and functions must focus on **what** the component does and define its public contract.

- **Black-Box Perspective**: Describe the behavior, return values, side effects, concurrency properties, and zero-value behavior of the component from the perspective of a caller who does not read the implementation.
- **Guarantees & Invariants**: Explicitly document behavioral expectations, inputs, outputs, and edge-case guarantees.
- **No Internal Leaks**: Do not detail internal step-by-step code execution, private algorithms, or internal struct manipulations inside doc comments.

### Example
```go
// detectVersion queries the GCC driver binary at path to extract its release version string.
// The path must be non-empty and point to a valid GCC binary.
func detectVersion(path string) (string, error) {
    // ...
}
```

---

## 2. Inline Comments (Design Rationale / "Why")

Code demonstrates *how* logic runs; inline comments must explain **why** a specific design decision, fallback path, or defensive check was implemented.

- **Rationale Over Execution**: Do not write comments that restate what the code is doing. If a comment is needed to explain *what* the code does, simplify or refactor the code instead.
- **Mandatory Logic Annotations**: Any non-obvious conditional logic, defensive workarounds, or parsing functions must be documented inline.
- **Proof & Source Attribution**: Back up complex, non-obvious, or defensive logic with concrete evidence:
  - **Prefer Official Documentation**: Whenever possible, link directly to official documentation, specifications, or issue tracker links that justify the workaround (e.g., linking to GCC changes for legacy flag behavior).
  - **Fallback to Descriptions or Visuals**: If official documentation is unavailable, include a clear descriptive comment or visual snippet showing the shape of the data (e.g., a sample raw compiler error message you are parsing) so readers see why the code is necessary.

### Example
```go
// GCC 7 introduced -dumpfullversion to guarantee a 3-part version string (major.minor.patch) suitable
// for release tag matching (https://gcc.gnu.org/gcc-7/changes.html).
cmd := exec.Command(path, "-dumpfullversion")
```

---

## 3. Sensible Redundancy & Unified Branch Documentation

- **Avoid Repetitive Boilerplate**: Do not repeat identical contextual explanations inside every branch of an `if-else` chain or `switch` block.
- **Grouped Branch Comments**: When multiple conditional branches stem from the same underlying context or decision point, place a single overarching block comment preceding the branching block to explain the shared context and strategy.

### Example
```go
// The compiler binary could be missing, or it might be an older installation (< GCC 7)
// that does not support the -dumpfullversion flag. We check the full version first
// and fall back to -dumpversion if the command fails.
cmd := exec.Command(path, "-dumpfullversion")
out, err := cmd.Output()
if err != nil {
    cmd = exec.Command(path, "-dumpversion")
    out, err = cmd.Output()
    if err != nil {
        return "", fmt.Errorf("detecting GCC version via %s: %w", path, err)
    }
}
```

---

## Authoritative References

For deeper reading on Go commenting:
- [Go Doc Comments Specification](https://go.dev/doc/comment)
- [Effective Go: Commentary](https://go.dev/doc/effective_go#commentary)
- [Google Go Style Guide: Comments](https://google.github.io/styleguide/go/decisions#comment-sentences)
