// Package cdb provides compiler-agnostic configuration structures, types,
// and execution coordination for building compilation databases.
package cdb

import (
	"slices"
	"strings"
)

// OptionKind defines the parsing behavior for a compiler option pattern.
type OptionKind string

const (
	// OptionKindFlag represents a standalone option that consumes no arguments.
	// The option appears by itself in argv (e.g., -c, -v).
	OptionKindFlag OptionKind = "Flag"

	// OptionKindJoined represents an option whose argument is appended directly
	// to the option spelling with no separator.
	//
	// Joined option patterns strictly require a non-empty suffix. If a joined option
	// can have an empty suffix, it is represented using multiple distinct patterns
	// (e.g., a Joined pattern for a non-empty suffix, and a Flag pattern for an empty suffix).
	OptionKindJoined OptionKind = "Joined"

	// OptionKindSeparate represents an option that consumes exactly one
	// subsequent argv element as its argument (e.g., -o out consumes "out").
	OptionKindSeparate OptionKind = "Separate"

	// OptionKindMultiArg represents an option that consumes a fixed number
	// (NumArgs) of subsequent argv elements as its arguments (e.g., -MF a b
	// with NumArgs=2 consumes "a" and "b").
	OptionKindMultiArg OptionKind = "MultiArg"

	// OptionKindJoinedAndSeparate represents an option that accepts both a
	// joined suffix and a separate argument.
	//
	// The joined suffix is extracted from the same argv element and must be non-empty.
	// If the joined suffix can be empty, it is represented as a separate OptionKindSeparate pattern.
	OptionKindJoinedAndSeparate OptionKind = "JoinedAndSeparate"

	// OptionKindRemainingArgs represents an option that consumes all remaining
	// argv elements as its arguments. Everything after the flag is captured.
	OptionKindRemainingArgs OptionKind = "RemainingArgs"

	// OptionKindRemainingArgsJoined represents an option that accepts a joined
	// suffix and also consumes all remaining argv elements.
	//
	// The joined suffix is extracted from the same argv element and must be non-empty.
	// If the joined suffix can be empty, it is represented as a separate OptionKindRemainingArgs pattern.
	OptionKindRemainingArgsJoined OptionKind = "RemainingArgsJoined"
)

// IsJoined reports whether the option kind accepts a joined suffix — that is,
// an argument appended directly to the option spelling with no separator.
func (k OptionKind) IsJoined() bool {
	return k == OptionKindJoined ||
		k == OptionKindJoinedAndSeparate ||
		k == OptionKindRemainingArgsJoined
}

// OptionPattern represents a single command-line spelling variant of an option.
type OptionPattern struct {
	// Spelling is the exact spelling prefix of the option flag (e.g. "-std=" or "-c").
	Spelling string
	// Kind specifies how the option consumes arguments from the command line.
	Kind     OptionKind
	// NumArgs is used only by OptionKindMultiArg; zero for all other kinds.
	NumArgs int
}

// Config holds sorted option entries with back-chain links for prefix matching.
type Config struct {
	// Patterns is the list of compiler option patterns, sorted lexicographically.
	Patterns   []OptionPattern
	// BackChains maps each index in Patterns to the index of its longest joined proper prefix
	// pattern, or -1 if no such prefix exists.
	BackChains []int
}

// NewConfig sorts the provided option patterns by spelling and kind (separated before
// joined) and computes back-chain links for prefix-based matching.
func NewConfig(patterns []OptionPattern) *Config {
	patterns = slices.Clone(patterns)
	slices.SortFunc(patterns, func(a, b OptionPattern) int {
		if d := strings.Compare(a.Spelling, b.Spelling); d != 0 {
			return d
		}

		if !a.Kind.IsJoined() && b.Kind.IsJoined() {
			return -1
		}
		if a.Kind.IsJoined() && !b.Kind.IsJoined() {
			return 1
		}
		return 0
	})

	// Compute back-chains: for each joined pattern, scan forward to link
	// any patterns that share it as a prefix. As we progress, backChains[j]
	// is overwritten by longer matched prefixes, ensuring it always points
	// to the longest joined proper prefix.
	backChains := make([]int, len(patterns))
	for i := range backChains {
		backChains[i] = -1
	}
	for i := range patterns {
		if !patterns[i].Kind.IsJoined() {
			continue
		}
		for j := i + 1; j < len(patterns); j++ {
			if !strings.HasPrefix(patterns[j].Spelling, patterns[i].Spelling) {
				break
			}
			backChains[j] = i
		}
	}
	return &Config{
		Patterns:   patterns,
		BackChains: backChains,
	}
}
