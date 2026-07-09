TODO: really rewrite. i think it's still a little stale

# GCC CDB Configuration Generator — Agent Bootstrapping

Generate [generated_cdbconfig.go](generated_cdbconfig.go) — a Go source file containing GCC's option configuration as a `*cdb.Config` value for the compilation database.

> **Note:** Before running these steps, check whether the configuration already exists:
> ```bash
> test -f internal/gcc/generated_cdbconfig.go && echo "already bootstrapped"
> ```
> If it does, skip bootstrapping.

The bootstrapping process is automated via the package taskfile.

## Steps

### 1. Run the Bootstrap Task

Execute from the repository root:

```sh
task gcc:bootstrap
```

* **First run**: If `internal/gcc/scripts/bootstrap.py` does not exist, the task copies `bootstrap.example.py` to `bootstrap.py`.
* **Adaptation**: Edit `bootstrap.py`. Set the guard to `if False` after reviewing the script. Verify `GO`, `PYTHON3`, and `GCC` are set in `.env` (loaded by direnv). Set `BASE_URL` to match your GCC version and `OPT_FILES` per the comments in the script (use `$GCC -dumpmachine` for arch-specific `.opt` files).
* **Run**: `task gcc:bootstrap` downloads `.opt` files to `internal/gcc/tmp/` and runs `cdbconfiggen` to write `generated_cdbconfig.go`.

### 2. Manual Regeneration

If `.opt` files are already in `tmp/` and you only need to rerun the generator (for example, after editing `cdbconfiggen`):

Re-run the `go run` invocation at the end of `bootstrap.py`, or run `task gcc:bootstrap` again.
