// Package cdb provides the compilation database types and command parsing for cpx.
package cdb

import (
	"fmt"
	"sort"
	"strings"
)

// Option represents a parsed option and its arguments.
// Args excludes the flag itself: for joined options it contains the suffix,
// for separate/multi-arg options it contains the consumed arguments.
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

// findPattern locates the option pattern matching arg.
//
// Longest-prefix matching is performed in O(log n) time by binary searching the
// sorted option spellings and traversing pre-computed back-chain links. This algorithm
// is based on GCC's own internal options-parsing implementation, avoiding the complexity
// of a trie while remaining fast and easy to serialize.
func findPattern(cfg *Config, arg string) *OptionPattern {
	i := sort.Search(len(cfg.Patterns), func(i int) bool {
		return cfg.Patterns[i].Spelling >= arg
	})

	// Exact matches: non-joined options are returned directly.
	//
	// Joined options require a non-empty suffix to match. If a compiler option accepts
	// an empty suffix, it is represented as separate Joined and Flag patterns. Therefore,
	// an exact match on a Joined option's spelling (which has no shorter back-chain prefix)
	// returns nil to allow matching an alternate Flag pattern or falling back to a positional argument.
	if i < len(cfg.Patterns) && cfg.Patterns[i].Spelling == arg {
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
// Args for each option exclude the flag: joined options contain the suffix,
// separate/multi-arg options contain the consumed arguments.
func Parse(cfg *Config, argv []string) (Command, error) {
	if len(argv) == 0 {
		return Command{}, fmt.Errorf("argv is empty")
	}

	name := argv[0]
	options := []Option{}
	args := []string{}
	for i := 1; i < len(argv); i++ {
		pattern := findPattern(cfg, argv[i])
		if pattern == nil {
			args = append(args, argv[i])
			continue
		}

		optArgs := []string{}
		if pattern.Kind.IsJoined() {
			optArgs = append(optArgs, strings.TrimPrefix(argv[i], pattern.Spelling))
		}

		var n int
		switch pattern.Kind {
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
			return Command{}, fmt.Errorf(
				"option %s takes %d arguments, but only %d arguments are provided",
				pattern.Spelling, n, len(argv)-i,
			)
		}
		optArgs = append(optArgs, argv[i+1:i+n+1]...)

		options = append(options, Option{
			Pattern: *pattern,
			Args:    optArgs,
		})
		i += n
	}
	return Command{
		Name:    name,
		Options: options,
		Args:    args,
	}, nil
}
