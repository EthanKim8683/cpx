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
	// OptionKindFlag matches the option spelling exactly with no suffix or arguments.
	OptionKindFlag OptionKind = "Flag"
	// OptionKindJoined matches a spelling prefix followed by a required non-empty
	// suffix within the same argv entry (e.g. -std=c++17).
	OptionKindJoined OptionKind = "Joined"
	// OptionKindSeparate matches the option spelling exactly and consumes one
	// subsequent argv entry as its argument (e.g. -o file).
	OptionKindSeparate OptionKind = "Separate"
	// OptionKindMultiArg matches the option spelling exactly and consumes
	// NumArgs subsequent argv entries.
	OptionKindMultiArg OptionKind = "MultiArg"
	// OptionKindJoinedAndSeparate matches a spelling prefix followed by a
	// required non-empty suffix and additionally consumes one subsequent argv
	// entry (e.g. -MFfoo out.d).
	OptionKindJoinedAndSeparate OptionKind = "JoinedAndSeparate"
	// OptionKindRemainingArgs matches the option spelling exactly and consumes
	// all remaining argv entries.
	OptionKindRemainingArgs OptionKind = "RemainingArgs"
	// OptionKindRemainingArgsJoined matches a spelling prefix followed by an
	// optional non-empty suffix and consumes all remaining argv entries.
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
	// ByPrefix maps each spelling string to its associated OptionPattern(s).
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
