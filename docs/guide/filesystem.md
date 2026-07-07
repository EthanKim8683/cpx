# Filesystem

- **Use afero**: All disk I/O goes through `afero.Fs`. No direct `os` calls in production code.
- **Unit tests**: Use `afero.NewMemMapFs()` for read/write, `testing/fstest.MapFS` for read-only.
- **Integration tests**: Use `afero.NewOsFs()` when real disk output is required.
- **Method syntax**: Wrap with `&afero.Afero{Fs: fs}` and call `afs.ReadFile(path)`.

See [ADR-0003](../adr/0003-filesystem-abstraction.md) for the full decision history.
