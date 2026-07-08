// Package cdb provides compiler-agnostic configuration structures and types
// for parsing command-line arguments into discrete options.
//
// Both GCC and Clang match options by longest prefix — e.g., --std=c++17
// matches --std= (not --). NewConfig sorts patterns by spelling and links
// back-chain pointers so that binary search plus back-chain traversal
// achieves longest-prefix matching in O(log n).
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
	// to the option spelling with no separator. The suffix after the spelling
	// is extracted as the option's argument (must be non-empty; e.g., -std=c++17
	// with spelling "-std=" yields argument "c++17", but exactly "-std=" is ignored).
	OptionKindJoined OptionKind = "Joined"

	// OptionKindSeparate represents an option that consumes exactly one
	// subsequent argv element as its argument (e.g., -o out consumes "out").
	OptionKindSeparate OptionKind = "Separate"

	// OptionKindMultiArg represents an option that consumes a fixed number
	// (NumArgs) of subsequent argv elements as its arguments (e.g., -MF a b
	// with NumArgs=2 consumes "a" and "b").
	OptionKindMultiArg OptionKind = "MultiArg"

	// OptionKindJoinedAndSeparate represents an option that accepts both a
	// joined suffix and a separate argument. The suffix is extracted from the
	// same argv element (must be non-empty), and one additional argv element is
	// consumed (e.g., -Ifoo bar with spelling "-I" yields joined suffix "foo" and
	// separate argument "bar").
	OptionKindJoinedAndSeparate OptionKind = "JoinedAndSeparate"

	// OptionKindRemainingArgs represents an option that consumes all remaining
	// argv elements as its arguments. Everything after the flag is captured.
	OptionKindRemainingArgs OptionKind = "RemainingArgs"

	// OptionKindRemainingArgsJoined represents an option that accepts a joined
	// suffix and also consumes all remaining argv elements. The suffix is
	// extracted from the same argv element (must be non-empty), and the rest of
	// argv is appended.
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
	Spelling string
	Kind     OptionKind
	// NumArgs is used only by OptionKindMultiArg; zero for all other kinds.
	NumArgs int
}

// Config holds sorted option entries with back-chain links for prefix matching.
type Config struct {
	Patterns   []OptionPattern
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

	// Back-chain: for each joined kind, find the longest joined prefix
	// by scanning backward. Used by findPattern on exact match to
	// resolve to a joined option with a non-empty argument.
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
