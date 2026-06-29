# Agents Guide for cpx

Welcome to `cpx` — a Go 1.24 tool and library for competitive programming workflows.

## Quick Start Workflow

1. **Track Work**: Follow [issues.md](docs/agents/issues.md) for GitHub issue logging (`gh issue comment`).
2. **Architecture**:
   - [Overview](docs/overview.md) — System scope.
   - [ADRs](docs/adr/README.md) — Architecture decision records.
3. **Topic Guides**: Consult [docs/agents/](docs/agents/README.md).

## Required Verification Commands

Run and verify before completing work:

```bash
go generate ./...
go test ./...
golangci-lint run ./...
```
