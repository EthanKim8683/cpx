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
Execute the top-level task from the root of the repository:

```sh
task clang:bootstrap
```

*   **First Run**: If `internal/clang/scripts/bootstrap.sh` does not exist, the task will copy it from `bootstrap.example.sh` and exit.
*   **Adaptation**: The agent (or user) should verify `CLANG` is set in the environment or `.env` (so `direnv` loads it automatically). Additionally, verify that `BASE_URL` in `bootstrap.sh` points to the correct LLVM repository branch corresponding to the compiler version.
*   **Second Run**: Execute `task clang:bootstrap` again. The script will output the compiler version, configured URL, verify `clang-tblgen` / `llvm-tblgen` is available, prompt for interactive confirmations, download the `.td` files to `internal/clang/tmp/`, generate `options.json`, and run `cdbconfiggen` to write `generated_cdbconfig.go`.

### 2. Manual Regeneration
If you have already generated the `options.json` dump and only want to rerun the generator (for example, after editing the generator codebase):

```sh
task clang:generate
```
