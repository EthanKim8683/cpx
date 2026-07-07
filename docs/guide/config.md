# Configuration

- **Source**: `.env` (gitignored, human-editable).
- **Loader**: direnv via `.envrc` (`source_up_if_exists` + `dotenv`).
- **Runtime**: Environment variables are the shared API across Go, Make, AWK, and go generate.
- **Testing**: Use `env.Options{Environment: map[string]string{...}}` to isolate tests from the real environment.

See [ADR-0001](../adr/0001-configuration.md) for the full decision history.
