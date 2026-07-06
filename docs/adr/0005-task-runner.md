# ADR-0005: Task Runner

- **Status**: Accepted
- **Date**: 2026-07-05
- **Related**: [GitHub #24](https://github.com/EthanKim8683/cpx/issues/24)

---

## Context

cpx needs a task runner for common development workflows — generating configs, running tests, building, and other repeated commands. Make is the traditional choice in Go projects, but it is not natively available on Windows. Developers on Windows must use WSL, Cygwin, or MSYS2 to run Makefiles, which adds friction and breaks the out-of-the-box experience.

Cross-platform compatibility matters. cpx should work on macOS, Linux, and Windows without requiring platform-specific shell environments.

---

## Decision

Use [Task](https://taskfile.dev) (Taskfile.yml) as the project task runner.

Task is a Go-based task runner with a YAML task definition format. It is installed via `go install`, `brew`, `scoop`, or direct binary download — no shell environment assumptions. Tasks are defined in `Taskfile.yml` at the project root.

```yaml
version: 3
tasks:
  default:
    cmds:
      - task: gen-env-example
  gen-env-example:
    desc: Generate .env.example file
    cmds:
      - npx -y @dotenvx/dotenvx ext genexample -f .env
```

Tasks use the same shell commands they would in a Makefile, but Task handles invocation cross-platform. Environment variables from direnv are inherited automatically.

---

## Alternatives

- **Make**: Ubiquitous in Go projects, but not natively available on Windows. Requires WSL/Cygwin/MSYS2, which adds setup friction and breaks the cross-platform goal.
- **Just**: A modern command runner, but written in Rust — adds a non-Go dependency to the toolchain.
- **Shell scripts**: Platform-specific (bash vs PowerShell vs cmd). Duplicates logic across platforms.
- **Go generate + go run**: Works for code generation but doesn't provide a named task interface for common workflows.

---

## Consequences

- **Cross-platform**: Task runs on macOS, Linux, and Windows without shell environment assumptions.
- **Go ecosystem**: Installed via `go install`, consistent with the rest of the toolchain.
- **YAML definition**: Task files are YAML, which is widely readable and editable.
- **direnv integration**: Environment variables from `.env` are inherited by task commands automatically.
- **No Make dependency**: Removes the need for Make, simplifying the developer setup.
