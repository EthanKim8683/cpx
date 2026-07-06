// Package cdb provides compiler-agnostic configuration structures and types
// for parsing command-line arguments into discrete options.
//
// Options are matched by spelling prefix using a flat, prefix-based pattern
// matching pass. Dynamic behaviors — negation, overrides, mutual exclusions —
// are deferred to access time via consumers like getLastArg or hasFlag.
package cdb

// OptionKind defines the parsing behavior for a compiler option pattern.
type OptionKind string

const (
	// OptionKindFlag matches exactly.
	OptionKindFlag OptionKind = "Flag"
	// OptionKindJoined matches a prefix and a non-empty suffix.
	OptionKindJoined OptionKind = "Joined"
	// OptionKindSeparate matches exactly and consumes one subsequent arg.
	OptionKindSeparate OptionKind = "Separate"
	// OptionKindMultiArg matches exactly and consumes NumArgs subsequent args.
	OptionKindMultiArg OptionKind = "MultiArg"
	// OptionKindJoinedAndSeparate matches a prefix and a non-empty suffix and consumes one subsequent arg.
	OptionKindJoinedAndSeparate OptionKind = "JoinedAndSeparate"
	// OptionKindRemainingArgs matches exactly and consumes all remaining args.
	OptionKindRemainingArgs OptionKind = "RemainingArgs"
	// OptionKindRemainingArgsJoined matches a prefix and a non-empty suffix and consumes all remaining args.
	OptionKindRemainingArgsJoined OptionKind = "RemainingArgsJoined"
)

// OptionPattern represents a single command-line spelling variant of an option.
type OptionPattern struct {
	Spelling string
	Kind     OptionKind
	// NumArgs is used only by OptionKindMultiArg; zero for all other kinds.
	NumArgs int
}

// Config maps spelling prefixes to their option patterns.
type Config struct {
	ByPrefix map[string][]OptionPattern
}

// NewConfig constructs a Config by indexing the provided option patterns by
// their spelling prefix.
func NewConfig(patterns []OptionPattern) *Config {
	byPrefix := make(map[string][]OptionPattern, len(patterns))
	for _, pattern := range patterns {
		prefix := pattern.Spelling
		byPrefix[prefix] = append(byPrefix[prefix], pattern)
	}
	return &Config{ByPrefix: byPrefix}
}
