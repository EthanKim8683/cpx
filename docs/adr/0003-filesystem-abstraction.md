# ADR-0003: Filesystem Abstraction

- **Status**: Accepted
- **Date**: 2026-06-28
- **Related**: [CP-6](https://linear.app/ethankim8683/issue/CP-6/filesystem-abstraction-adr)

## Context

cpx reads and writes filesystems — config, caches, bundles, temp dirs. Without abstraction, real I/O is coupled to every workflow. Unit tests hit the real disk, break in CI, and race when parallel.

The requirement is an abstraction layer over disk and in-memory filesystems. cpx currently has 24 `afero.Fs` and 5 `afero.Afero` references, mostly in config, cdb, and gcc. Clang is not yet wired up. The abstraction is in production use, not aspirational.

## Decision

- **Use `afero.Fs` everywhere**: All disk I/O goes through the `afero.Fs` interface.
- **Method syntax**: Wrap with `&afero.Afero{Fs: fs}` and call `afs.ReadFile(path)`. Avoids `afero.ReadFile(fs, path)` call-site noise.
- **Unit tests**: Use `afero.NewMemMapFs()` for read/write, `testing/fstest.MapFS` for read-only.
- **Integration tests**: Use `afero.NewOsFs()` when real disk output is required.
- **Where to apply**: Config loading, cache reading/writing, bundle output, temp directory management, manifest handling, GCC/Clang scratch dirs.

## Alternatives considered

- **`io/fs` + `os` adapter**: Stdlib since Go 1.16, but read-only by design. cpx needs write. Would require an `fs.FS` → `afero.Fs` adapter anyway.
- **`fs.FS` + own interface**: Possible, but invents a new abstraction with no ecosystem.
- **`os` package only**: Breaks unit testing and parallel safety.
- **Billy (src-d)**: Lower-level, no `MemMapFs` equivalent, unmaintained since 2020.

## Consequences

- **`MemMapFs` limitations**: No hardlink/symlink support, no mmap semantics, no real file permissions. Acceptable — cpx does not use these features.
- **Migration scope**: 24 `afero.Fs` + 5 `Afero` references already in cpx. Incremental expansion, not greenfield.
- **Clang gap**: 10 `os` calls in clang still need wiring. GCC is partially wired.
- **`fstest.MapFS` for read-only tests**: Avoids `MemMapFs` write setup when tests only read.

## References

- [CP-6](https://linear.app/ethankim8683/issue/CP-6/filesystem-abstraction-adr)

