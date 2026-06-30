# Agent Guidelines: Go Testing Standards

Testing in `cpx` is built on Go standard library patterns, emphasizing readability, isolation, and robustness. Uniformity is key: all tests must follow these guidelines.

---

## 1. Tool Selection and Library Guidelines

- **Use the Standard Library**: Rely primarily on the Go standard `testing` package.
- **No Testify**: Do not use third-party assertion libraries such as `github.com/stretchr/testify`. Use simple `if` statements and standard comparisons instead.
- **Compare Complex Data with `go-cmp`**: Use [go-cmp](https://pkg.go.dev/github.com/google/go-cmp/cmp) (`github.com/google/go-cmp/cmp`) to compare structs, slices, maps, and other nested types. Avoid using `reflect.DeepEqual`, as it is less flexible and does not produce human-readable diffs.

---

## 2. Table-Driven Test Patterns

When testing multiple input/output scenarios for a single functional unit, use the table-driven test pattern. Group test cases in a slice of anonymous structs and run each case in isolation using [t.Run](https://pkg.go.dev/testing#T.Run) for subtests.

### Example Pattern

```go
func TestParseVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "valid version",
			input: "gcc version 14.2.0",
			want:  "14.2.0",
		},
		{
			name:    "invalid format",
			input:   "gcc-14.2.0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseVersion(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("parseVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if got != tt.want {
				t.Errorf("parseVersion() = %q; want %q", got, tt.want)
			}
		})
	}
}
```

---

## 3. Asserting Correctness: got-before-want Ordering

Failure messages must be uniform across the codebase. Always present the actual result (`got`) before the expected result (`want`) in assertion messages.

- **Good**: `t.Errorf("parseVersion() = %q; want %q", got, tt.want)`
- **Bad**: `t.Errorf("expected %q, but got %q", tt.want, got)`

Using `got` followed by `want` is standard Go practice and helps maintain consistency, reducing cognitive overhead when scanning failed test outputs.

---

## 4. Test Helpers and Cleanup

### Using `t.Helper()`

When writing helper functions that perform checks or setup on behalf of a test, always call [t.Helper()](https://pkg.go.dev/testing#T.Helper) first. This marks the helper function so that failure line numbers in reports point to the line of the calling test function rather than inside the helper itself.

### Using `t.Cleanup()`

Avoid executing tear-down logic via `defer` inside helper functions, as deferred actions in a helper will run immediately when the helper returns, not at the end of the test. Instead, register teardown logic using [t.Cleanup()](https://pkg.go.dev/testing#T.Cleanup).

- `t.Cleanup` runs after the test (and all its subtests) completes.
- It simplifies setup functions by keeping cleanup logic co-located with resource allocation.

### Example Usage

```go
// newMockServer creates a test HTTP server and registers its teardown via t.Cleanup().
func newMockServer(t *testing.T, responseBody string) string {
	t.Helper() // Align failure reports with the calling test line

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responseBody))
	}))

	t.Cleanup(server.Close) // Automatically closed at the end of the calling test

	return server.URL
}
```

---

## 5. Golden File Testing with Goldie

Golden file testing is used to verify complex, large, or multi-line outputs (e.g., compiler option tables, generated code, large JSON payloads) without cluttering Go source files with massive raw strings.

### Rationale (Based on Mitchell Hashimoto's Talk)

As presented in Mitchell Hashimoto's *Advanced Testing with Go* talk, golden files act as the single source of truth for complex outputs:
1. **Clarity**: Prevents massive hardcoded string constants in `_test.go` files.
2. **Visual Review via Git Diff**: Changes to the expected output are shown clearly in git pull requests, making review simple.
3. **Decoupled Verification**: Separates the test framework assertion logic from the data output.

### Standard Usage with `sebdah/goldie`

We use the [sebdah/goldie/v2](https://github.com/sebdah/goldie) package.

1. Instantiate the goldie client using `goldie.New(t)`.
2. Assert the payload using `g.Assert(t, "fixture_name", []byte(got))`. This compares the output to the contents of `testdata/fixture_name.golden`.

```go
func TestGenerateOutput(t *testing.T) {
	got := generateComplexOutput()

	g := goldie.New(t)
	g.Assert(t, "complex_output", []byte(got))
}
```

### Running and Updating Golden Files

By default, tests will verify actual output against the saved golden files. If you intentionally changed the output and want to update the stored golden files, run:

```bash
go test -update ./...
```

Verify the generated changes in the `testdata/` directory using `git diff` before staging and committing.

---

## 6. Recursive Directory Structure and File Verification

While **Golden File Testing** (Section 5) is ideal for verifying a single, large, or complex output payload, **Recursive Directory Verification** is the standard pattern when your code generates or modifies entire directory skeletons (e.g., project skeleton generators, package bundlers, file system sync utilities) and the file layout itself is part of the correctness contract.

### The Idiomatic Go Pattern

To verify directory structures without external helper libraries:
1. **Traverse & Map**: Use `filepath.WalkDir` to walk the target directory. Map the relative file paths (as keys) to their contents (as values).
2. **Compare**: Use `cmp.Diff` (from `github.com/google/go-cmp/cmp`) to assert the generated actual map structure against the expected (mocked) map structure.

#### Example Implementation

```go
// readDirTree walks dir recursively and returns a map of relative file paths to their string contents.
func readDirTree(t *testing.T, dir string) map[string]string {
	t.Helper()

	tree := make(map[string]string)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Normalize line endings to avoid cross-platform test flakiness
		normalizedContent := strings.ReplaceAll(string(content), "\r\n", "\n")
		tree[rel] = normalizedContent
		return nil
	})

	if err != nil {
		t.Fatalf("failed to walk directory %s: %v", dir, err)
	}

	return tree
}

func TestGenerateProject(t *testing.T) {
	tmp := t.TempDir()

	err := generateProjectSkeleton(tmp)
	if err != nil {
		t.Fatalf("generateProjectSkeleton() error = %v", err)
	}

	got := readDirTree(t, tmp)
	want := map[string]string{
		"README.md":   "# Project Name\n",
		"src/main.go": "package main\n\nfunc main() {}\n",
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("directory structure mismatch (-want +got):\n%s", diff)
	}
}
```

### Key Guidelines and Gotchas for Agents

When implementing recursive directory testing, keep these design considerations in mind:
- **Exclude System/VCS Files**: Do not assume the temp directory will only contain the files you wrote. If necessary, skip OS or VCS files (like `.DS_Store`, `.git`, or temporary locks) during the walk.
- **Normalize Line Endings**: Windows uses `\r\n` (CRLF) while Linux and macOS use `\n` (LF). Always use `strings.ReplaceAll(content, "\r\n", "\n")` on read files to prevent tests from failing on Windows CI environments.
- **Verify Permissions/Modes**: If executable bits or symlink targets are part of your contract, map the relative path to a custom struct containing both the mode and content rather than just a raw string:
  ```go
  type fileInfo struct {
      Mode    fs.FileMode
      Content string
  }
  ```

> [!NOTE]
> This standard-library-first approach (using `filepath.WalkDir` and `cmp.Diff`) is modeled after the Go compiler's own test suite (`src/cmd/internal/testdir/testdir_test.go`). It avoids rigid third-party assertion dependencies, provides transparent control over file traversal, and leverages `go-cmp` to output clean line-by-line diffs on mismatch.

---

## 7. Integration Testing

Integration tests in `cpx` verify the interaction between multiple components or with external services (such as raw compiler binaries, the filesystem, or live network endpoints).

### Hermeticity vs. Live Integration
- **Prefer Hermeticity**: Whenever possible, mock external servers (e.g., using `net/http/httptest`) to keep tests self-contained, deterministic, and fast.
- **Explicit Live Boundaries**: Use live integration tests only when verifying contract compatibility with external third-party systems (such as the GCC GitHub mirror).

### Isolation and Execution Strategies

To prevent integration tests from slowing down everyday development or failing in environments without network/binary access, they must be isolated using one of the following strategies:

#### A. Build Tags (Recommended for Heavyweight/External Tests)
Heavyweight tests (those requiring Docker, databases, or third-party web access) must be placed in a separate file ending in `_integration_test.go` with a build constraint at the very top:

```go
//go:build integration

package mypkg
```

- **Running Unit Tests**: Standard unit tests run via `go test ./...` (tagged files are completely ignored during compilation and execution).
- **Running Integration Tests**: Run with the integration tag explicitly enabled:
  ```bash
  go test -tags=integration ./...
  ```

#### B. Short Mode Skipping (testing.Short)
For lightweight integration tests that are part of the standard test suite but might require network or slow operations, query the `-short` flag at runtime:

```go
func TestFetchSourceLive(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live network integration test in short mode")
	}
	// ...
}
```

- **Running Unit Tests**: Skip these tests locally via `go test -short ./...`.
- **Running Integration Tests**: Run all tests including integration tests via `go test ./...`.

---

## Authoritative References

For deeper reading on Go testing standards and patterns:
- [Go Test Comments](https://go.dev/wiki/TestComments)
- [Table-Driven Tests](https://go.dev/wiki/TableDrivenTests)
- [Advanced Testing with Go (Mitchell Hashimoto)](https://www.youtube.com/watch?v=yszygk1cpEc)
- [google/go-cmp](https://pkg.go.dev/github.com/google/go-cmp/cmp)
- [github.com/sebdah/goldie/v2](https://github.com/sebdah/goldie)
- [Go Build Constraints / Tags](https://pkg.go.dev/go/build)
- [Go Test Flags (-short)](https://pkg.go.dev/cmd/go#hdr-Testing_flags)
- [Go toolchain testdir helper (src/cmd/internal/testdir/testdir_test.go)](https://cs.opensource.google/go/go/+/master:src/cmd/internal/testdir/testdir_test.go)
