# cpx — Agent Setup

## Bootstrapping

> **Note:** Bootstrapping only needs to happen once per machine. Before running
> these steps, check whether configs already exist (e.g.
> `internal/gcc/generated/CDBConfig.go` and
> `internal/clang/generated/CDBConfig.go`). If they do, skip bootstrapping.

The compilation database needs generated option configs that teach it how to
parse the user's GCC and Clang commands. These configs are compiler-specific
and must be generated per-machine.

1. Read `.env.example` and confirm the `GCC` and `CLANG` paths with the user.
   Verify each path exists. If the user doesn't have GCC or Clang installed,
   offer to help them install it (e.g. `brew install gcc` or `brew install
   llvm`).

2. Allow direnv so environment variables load automatically:
   ```
   direnv allow
   ```

3. Bootstrap each compiler config by following the AGENTS.md in its package:
   - `internal/gcc/AGENTS.md`
   - `internal/clang/AGENTS.md`

## Development

> **Mandatory:** Load the `saving-sessions` skill immediately before doing
> anything else, using your preferred method for loading skills.

See [`docs/README.md`](docs/README.md) to get started.
