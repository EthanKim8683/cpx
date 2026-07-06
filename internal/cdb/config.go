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
	OptionKindFlag                OptionKind = "Flag"
	OptionKindJoined              OptionKind = "Joined"
	OptionKindSeparate            OptionKind = "Separate"
	OptionKindMultiArg            OptionKind = "MultiArg"
	OptionKindJoinedAndSeparate   OptionKind = "JoinedAndSeparate"
	OptionKindRemainingArgs       OptionKind = "RemainingArgs"
	OptionKindRemainingArgsJoined OptionKind = "RemainingArgsJoined"
)

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
	BackChains []*OptionPattern
}

// NewConfig sorts the provided option patterns by spelling and computes
// back-chain links for prefix-based matching.
func NewConfig(patterns []OptionPattern) *Config {
	patterns = slices.Clone(patterns)
	slices.SortFunc(patterns, func(a, b OptionPattern) int {
		return strings.Compare(a.Spelling, b.Spelling)
	})

	// Back-chain: for each joined kind, find the longest joined prefix
	// by scanning backward. Used by findPattern on exact match to
	// resolve to a joined option with a non-empty argument.
	backChains := make([]*OptionPattern, len(patterns))
	for i := range patterns {
		if !patterns[i].Kind.IsJoined() {
			continue
		}
		for j := i - 1; j >= 0; j-- {
			if !strings.HasPrefix(patterns[i].Spelling, patterns[j].Spelling) {
				continue
			}
			if patterns[j].Kind.IsJoined() {
				backChains[i] = &patterns[j]
				break
			}
		}
	}
	return &Config{
		Patterns:   patterns,
		BackChains: backChains,
	}
}
