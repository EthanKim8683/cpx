# Agent Guidelines: Go Testing Standards

Testing in `cpx` is built on Go standard library patterns, emphasizing readability, isolation, and robustness. Uniformity is key: all tests must follow these guidelines.

---

## 1. Tool Selection and Library Guidelines

- **Use Testify Assertions**: Use the [testify](https://github.com/stretchr/testify) assertion library (`assert` and `require`) to perform all test assertions. Use `github.com/stretchr/testify/assert` for non-fatal checks and `github.com/stretchr/testify/require` for assertions that must stop execution immediately (e.g., failed setup or nil pointers).
- **No go-cmp**: Avoid introducing additional third-party comparison/assertion dependencies like `github.com/google/go-cmp/cmp`. Testify's `assert.Equal` is sufficient for map, slice, and struct comparisons, and outputs clear visual diffs out of the box.

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

			if tt.wantErr {
				require.Error(t, err)
				assert.Empty(t, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
```

---

## 3. Asserting Correctness: expected-before-actual Ordering (Testify)

When using Testify assertions, it is critical to follow the signature convention of the framework. Always place the expected result (`want`) before the actual result (`got`).

- **Good**: `assert.Equal(t, tt.want, got)`
- **Bad**: `assert.Equal(t, got, tt.want)`

> [!WARNING]
> Mismatching the parameter order in Testify (i.e., putting `got` before `want`) will result in confusing failure output where "Expected" and "Actual" values are inverted in the test report.
>
> If you write manual standard library assertions with `t.Errorf`, use the standard Go `got` before `want` pattern:
> `t.Errorf("parseVersion() = %q; want %q", got, tt.want)`

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

## 6. Directory Structure and File Verification

Verifying file system layouts and file contents must follow industry standards. Depending on the scenario, choose the appropriate strategy:

### A. Primary Standard: In-Memory Verification (Unit Testing)

Whenever possible, decouple your code from the physical storage disk.
- **For Reading**: Implement your functions to take Go's standard `fs.FS` interface. In tests, use Go's standard `testing/fstest.MapFS` to mock file system contents entirely in memory.
- **For Writing**: Use [afero](https://github.com/spf13/afero) (`afero.Fs`) as detailed in [filesystem.md](file:///Users/ethankim8683/Competitive%20Programming/Utilities/cpx/docs/agents/filesystem.md). In tests, swap the disk filesystem with `afero.NewMemMapFs()` and check results in memory.

#### Example using `fstest.MapFS` (Mocking Read I/O)

```go
func ParseConfigs(fds fs.FS) (map[string]string, error) {
	// ... logic reading from fds ...
}

func TestParseConfigs(t *testing.T) {
	mockFS := fstest.MapFS{
		"config.json": &fstest.MapFile{Data: []byte(`{"version": "1.0"}`)},
	}

	got, err := ParseConfigs(mockFS)
	require.NoError(t, err)
	assert.Equal(t, "1.0", got["version"])
}
```

#### Example using `afero.Fs` (Mocking Write I/O)

```go
func WriteConfig(fs afero.Fs, path string, data string) error {
	return afero.WriteFile(fs, path, []byte(data), 0644)
}

func TestWriteConfig(t *testing.T) {
	appFS := afero.NewMemMapFs()

	err := WriteConfig(appFS, "/app/config.json", `{"enabled": true}`)
	require.NoError(t, err)

	exists, err := afero.Exists(appFS, "/app/config.json")
	require.NoError(t, err)
	assert.True(t, exists)

	content, err := afero.ReadFile(appFS, "/app/config.json")
	require.NoError(t, err)
	assert.Equal(t, `{"enabled": true}`, string(content))
}
```

### B. Secondary Standard: Recursive Directory Tree Mapping (Integration Testing)

When testing operations that *must* interact with the physical disk (e.g., compiler generation, code templating output), walk the target directory tree, map it to a structured map, and compare maps using Testify's `assert.Equal`.

1. **Traverse & Map**: Walk the directory using `filepath.WalkDir` and build a `map[string]string` of relative paths to normalized contents.
2. **Assert**: Compare the mapped structure with the expected layout map using `assert.Equal(t, want, got)`.

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

	require.NoError(t, err, "failed to walk directory %s", dir)

	return tree
}

func TestGenerateProject(t *testing.T) {
	tmp := t.TempDir()

	err := generateProjectSkeleton(tmp)
	require.NoError(t, err)

	got := readDirTree(t, tmp)
	want := map[string]string{
		"README.md":   "# Project Name\n",
		"src/main.go": "package main\n\nfunc main() {}\n",
	}

	assert.Equal(t, want, got)
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
> Combining in-memory FS abstractions with map-based Testify assertions is the standard approach used by major Go projects (such as Hugo and the Go compiler itself). This prevents disk test leakage, is cross-platform safe, and delivers clear line-by-line mismatch reports without third-party diff engines.

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
- [github.com/stretchr/testify](https://github.com/stretchr/testify)
- [github.com/sebdah/goldie/v2](https://github.com/sebdah/goldie)
- [filesystem.md](file:///Users/ethankim8683/Competitive%20Programming/Utilities/cpx/docs/agents/filesystem.md)
- [ADR-0003: File System Abstraction](file:///Users/ethankim8683/Competitive%20Programming/Utilities/cpx/docs/adr/0003-filesystem-abstraction.md)
- [ADR-0004: Testing and Assertion Framework](file:///Users/ethankim8683/Competitive%20Programming/Utilities/cpx/docs/adr/0004-testing-framework.md)
- [Go Build Constraints / Tags](https://pkg.go.dev/go/build)
- [Go Test Flags (-short)](https://pkg.go.dev/cmd/go#hdr-Testing_flags)
- [Go toolchain testdir helper (src/cmd/internal/testdir/testdir_test.go)](https://cs.opensource.google/go/go/+/master:src/cmd/internal/testdir/testdir_test.go)
