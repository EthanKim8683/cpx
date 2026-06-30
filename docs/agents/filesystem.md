# Agent Guidelines: File System Abstraction

To ensure test isolation, parallelizability, and cross-platform compatibility, all file system operations in `cpx` must be decoupled from the physical operating system using `github.com/spf13/afero`.

---

## 1. The Core Standard: No Direct `os` Calls

Production code must not call package-level file operations from the standard library `os` package (e.g., `os.Create`, `os.Open`, `os.ReadFile`). Instead, operations must be performed using an injected `afero.Fs` interface.

### Bad: Direct OS Binding
```go
func LoadConfig(path string) ([]byte, error) {
	return os.ReadFile(path) // Hardcoded to real disk, untestable in isolation
}
```

### Good: Dependency Injection
```go
func LoadConfig(fs afero.Fs, path string) ([]byte, error) {
	return afero.ReadFile(fs, path) // Fully mockable
}
```

---

## 2. Structured Abstractions

For packages that manage multiple file system interactions, receive the `afero.Fs` dependency at struct initialization:

```go
type Generator struct {
	fs afero.Fs
}

func NewGenerator(fs afero.Fs) *Generator {
	return &Generator{fs: fs}
}

func (g *Generator) CreateSkeleton(dir string) error {
	return g.fs.MkdirAll(dir, 0755)
}
```

---

## 3. Useful Afero Helpers

The `afero` library provides utility functions that mirror standard packages (`os`, `io/ioutil`). Always prefer using these wrapper functions:

*   **Read File**: `afero.ReadFile(fs, path)`
*   **Write File**: `afero.WriteFile(fs, path, data, perm)`
*   **Check Existence**: `afero.Exists(fs, path)`
*   **Create/Overwrite File**: `fs.Create(path)`
*   **Remove Directory/File**: `fs.RemoveAll(path)`

---

## 4. Testing Guidelines

### A. Unit Tests (In-Memory)
Always perform unit testing in memory using `afero.NewMemMapFs()`. This allows tests to run instantly, in parallel, and without cleaning up temporary files.

```go
func TestGenerator(t *testing.T) {
	mockFS := afero.NewMemMapFs()
	gen := NewGenerator(mockFS)

	err := gen.CreateSkeleton("/test")
	require.NoError(t, err)

	exists, err := afero.Exists(mockFS, "/test")
	require.NoError(t, err)
	assert.True(t, exists)
}
```

### B. Integration Tests & Production
Use the physical filesystem `afero.NewOsFs()` only at the application's main entry points (wire-up) or inside heavyweight integration tests.

---

## 5. Cross-Platform Guidelines
- **Path Normalization**: Windows uses `\` as a path separator, while macOS and Linux use `/`. To prevent test failures, use `filepath.ToSlash` to normalize paths before asserting them, or stick to forward slashes in memory (which Afero's `MemMapFs` handles transparently).
- **Line Endings**: When verifying generated text files in tests, normalize CRLF (`\r\n`) to LF (`\n`) to prevent environment differences from breaking assertions.
