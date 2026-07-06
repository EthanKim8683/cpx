package cdb

import (
	"fmt"
	"slices"
	"strings"
)

// Option represents a parsed option and its arguments.
type Option struct {
	Pattern OptionPattern
	Args    []string
}

// Command represents a parsed command with its name, options, and positional arguments.
type Command struct {
	Name    string
	Options []Option
	Args    []string
}

// findPattern locates the option pattern matching arg via binary search and
// back-chain prefix matching.
func findPattern(cfg *Config, arg string) *OptionPattern {
	// Binary searching over a list of strings has the property that the greatest string lexicographically
	// smaller or equal to the target string also has the longest common prefix with the target string
	// among the strings in the list.
	i, ok := slices.BinarySearchFunc(cfg.Patterns, arg, func(e OptionPattern, s string) int {
		return strings.Compare(e.Spelling, s)
	})

	// Exact match: for non-joined kinds, return the pattern directly.
	// For joined kinds, follow the back-chain once to find a proper prefix.
	// One step is sufficient because the exact match is itself a prefix of
	// the argument, so the back-chain target is a proper prefix.
	if ok {
		pattern := cfg.Patterns[i]
		if !pattern.Kind.IsJoined() {
			return &pattern
		}
		if j := cfg.BackChains[i]; j != -1 {
			return &cfg.Patterns[j]
		}
		return nil
	}

	for j := i - 1; j != -1; j = cfg.BackChains[j] {
		pattern := cfg.Patterns[j]
		if pattern.Kind.IsJoined() && strings.HasPrefix(arg, pattern.Spelling) {
			return &pattern
		}
	}
	return nil
}

// Parse parses argv into a Command. The first element is the command name.
// Args for each option include the flag itself followed by any consumed arguments.
func Parse(cfg *Config, argv []string) (Command, error) {
	if len(argv) == 0 {
		return Command{}, fmt.Errorf("argv is empty")
	}

	name := argv[0]
	var options []Option
	var args []string
	for i := 1; i < len(argv); i++ {
		pattern := findPattern(cfg, argv[i])
		if pattern == nil {
			args = append(args, argv[i])
			continue
		}

		var n int
		switch pattern.Kind {
		case OptionKindFlag:
			n = 0
		case OptionKindJoined:
			n = 0
		case OptionKindSeparate:
			n = 1
		case OptionKindMultiArg:
			n = pattern.NumArgs
		case OptionKindJoinedAndSeparate:
			n = 1
		case OptionKindRemainingArgs:
			n = len(argv) - (i + 1)
		case OptionKindRemainingArgsJoined:
			n = len(argv) - (i + 1)
		}
		if i+n+1 > len(argv) {
			return Command{}, fmt.Errorf("option %s takes %d arguments, but only %d arguments are provided", pattern.Spelling, n, len(argv)-i)
		}

		options = append(options, Option{
			Pattern: *pattern,
			Args:    argv[i : i+n+1],
		})
		i += n
	}
	return Command{
		Name:    name,
		Options: options,
		Args:    args,
	}, nil
}
