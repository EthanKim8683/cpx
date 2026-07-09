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
Execute the top-level task from the root of the repository:

```sh
task gcc:bootstrap
```

*   **First Run**: If `internal/gcc/scripts/bootstrap.sh` does not exist, the task will copy it from `bootstrap.example.sh` and exit.
*   **Adaptation**: The agent (or user) should verify `GCC` is set in the environment or `.env` (so `direnv` loads it automatically). Additionally, verify that `BASE_URL` in `bootstrap.sh` points to the correct GCC mirror branch corresponding to the compiler version.
*   **Second Run**: Execute `task gcc:bootstrap` again. The script will output the compiler version, configured URL, prompt for interactive confirmations, download the `.opt` files to `internal/gcc/tmp/`, and run `cdbconfiggen` to write `generated_cdbconfig.go`.

### 2. Manual Regeneration
If you have already downloaded the `.opt` files and only want to rerun the generator (for example, after editing the generator codebase):

```sh
task gcc:generate
```

