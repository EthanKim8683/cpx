TODO: really rewrite. i think it's still a little stale

# Clang CDB Configuration Generator — Agent Bootstrapping

Generate [generated_cdbconfig.go](generated_cdbconfig.go) — a Go source file containing Clang's option configuration as a `*cdb.Config` value for the compilation database.

> **Note:** Before running these steps, check whether the configuration already exists:
> ```bash
> test -f internal/clang/generated_cdbconfig.go && echo "already bootstrapped"
> ```
> If it does, skip bootstrapping.

The bootstrapping process is automated via the package taskfile.

## Steps

### 1. Run the Bootstrap Task

Execute from the repository root:

```sh
task clang:bootstrap
```

* **First run**: If `internal/clang/scripts/bootstrap.py` does not exist, the task copies `bootstrap.example.py` to `bootstrap.py`.
* **Adaptation**: Edit `bootstrap.py`. Set the guard to `if False` after reviewing the script. Verify `GO`, `PYTHON3`, and `CLANG` are set in `.env` (loaded by direnv). Set `BASE_URL`, `TBLGEN`, and `TD_FILES` per the comments in the script.
* **Run**: `task clang:bootstrap` downloads `.td` files to `internal/clang/tmp/`, runs `clang-tblgen` (or `llvm-tblgen`) to produce `options.json`, and runs `cdbconfiggen` to write `generated_cdbconfig.go`.

### 2. Manual Regeneration

If `options.json` is already in `tmp/` and you only need to rerun the generator (for example, after editing `cdbconfiggen`):

Re-run the final `go run` invocation at the end of `bootstrap.py`, or run `task clang:bootstrap` again.
