# Agent Guidelines: File System Abstraction

To ensure test isolation, parallelizability, and cross-platform compatibility, all file system operations in `cpx` must be decoupled from the physical operating system using `github.com/spf13/afero`.

---

## 1. The Core Standard: No Direct `os` Calls

Production code must not call package-level file operations from the standard library `os` package (e.g., `os.Create`, `os.Open`, `os.ReadFile`). Instead, operations must be performed using an injected `afero.Fs` interface, preferably wrapped in `afero.Afero` for cleaner method-based syntax.

### Bad: Direct OS Binding
```go
func LoadConfig(path string) ([]byte, error) {
	return os.ReadFile(path) // Hardcoded to real disk, untestable in isolation
}
```

### Good: Dependency Injection (using afero.Afero wrapper)
```go
func LoadConfig(fs afero.Fs, path string) ([]byte, error) {
	afs := &afero.Afero{Fs: fs}
	return afs.ReadFile(path) // Method syntax, fully mockable
}
```

---

## 2. Structured Abstractions

For packages that manage multiple file system interactions, receive the `afero.Fs` dependency at struct initialization, and wrap it internally in `afero.Afero` to keep method calls clean:

```go
type Generator struct {
	fs *afero.Afero
}

func NewGenerator(fs afero.Fs) *Generator {
	return &Generator{fs: &afero.Afero{Fs: fs}}
}

func (g *Generator) CreateSkeleton(dir string) error {
	return g.fs.MkdirAll(dir, 0755)
}
```

---

## 3. Useful Afero Helpers

Always wrap `afero.Fs` in `&afero.Afero{Fs: fs}` to use its rich helper methods directly as method calls rather than using package-level function wrappers:

*   **Read File**: `afs.ReadFile(path)`
*   **Write File**: `afs.WriteFile(path, data, perm)`
*   **Check Existence**: `afs.Exists(path)`
*   **Check Dir**: `afs.IsDir(path)`
*   **Create/Overwrite File**: `afs.Create(path)`
*   **Remove Directory/File**: `afs.RemoveAll(path)`

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

	afs := &afero.Afero{Fs: mockFS}
	exists, err := afs.Exists("/test")
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
