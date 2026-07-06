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

	// Exact match: for joined kinds (which require a non-empty suffix),
	// follow the back-chain to the longer joined prefix.
	// For non-joined kinds, return the pattern directly.
	if ok {
		if cfg.Patterns[i].Kind.IsJoined() {
			return cfg.BackChains[i]
		}
		return &cfg.Patterns[i]
	}

	if i == 0 {
		return nil
	}
	if pattern := cfg.Patterns[i-1]; pattern.Kind.IsJoined() {
		if strings.HasPrefix(arg, pattern.Spelling) {
			return &pattern
		}
	}
	if pattern := cfg.BackChains[i-1]; pattern != nil {
		if strings.HasPrefix(arg, pattern.Spelling) {
			return pattern
		}
	}
	return nil
}

// Parse parses argv into a Command. The first element is the command name.
// Args for each option include the flag itself followed by any consumed arguments.
func Parse(cfg *Config, argv []string) (Command, error) {
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
