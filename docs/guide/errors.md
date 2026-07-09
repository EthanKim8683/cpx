# Error Handling

- **Wrap with context**: `fmt.Errorf("loading config: %w", err)`.
- **Present participle**: Use "loading", "parsing", "writing" — not "load", "parse", "write".
- **One handler**: Errors are either logged or returned, never both.
