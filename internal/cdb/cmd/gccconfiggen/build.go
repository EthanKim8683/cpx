// build.go compiles parsed option records into compiler-agnostic config patterns
// declared in internal/cdb/config.go.

package main

import (
	"regexp"

	"github.com/EthanKim8683/cpx/internal/cdb"
)

// negateRE matches negatable flags. Options containing "=" (e.g. "-std=") are excluded
// because they require value parameters and do not have an implicit negative form.
var negateRE = regexp.MustCompile(`^[fgWm][^=]+$`)

// negative returns the negated spelling variant of a flag name (e.g. fcommon -> fno-common).
func negative(name string) string {
	return name[0:1] + "no-" + name[1:]
}

// buildOptionPatterns translates a parsed GCC option record into its compile-time spelling patterns.
func buildOptionPatterns(record parsedOptRecord) []cdb.OptionPattern {
	if record.rejectDriver {
		return nil
	}

	var kind cdb.OptionKind
	var numArgs int

	switch {
	case record.noDriverArg:
		// Separate options with NoDriverArg do not pass arguments during the driver stage.
		// They act as simple flags to prevent the parser from consuming subsequent command-line tokens.
		kind = cdb.OptionKindFlag
	case record.args != 0:
		kind = cdb.OptionKindMultiArg
		numArgs = record.args
	case record.joined && record.separate:
		kind = cdb.OptionKindJoinedOrSeparate
	case record.joined:
		kind = cdb.OptionKindJoined
	case record.separate:
		kind = cdb.OptionKindSeparate
	case record.joinedOrMissing:
		kind = cdb.OptionKindJoinedOrMissing
	default:
		kind = cdb.OptionKindFlag
	}

	var patterns []cdb.OptionPattern
	patterns = append(patterns, cdb.OptionPattern{
		Spelling: "-" + record.name,
		Kind:     kind,
		NumArgs:  numArgs,
	})

	if negateRE.MatchString(record.name) && !record.rejectNegative {
		patterns = append(patterns, cdb.OptionPattern{
			Spelling: "-" + negative(record.name),
			Kind:     kind,
			NumArgs:  numArgs,
		})
	}
	return patterns
}

// buildConfig compiles a list of parsed option records into a unified cdb.Config registry.
func buildConfig(records []parsedOptRecord) cdb.Config {
	patterns := make([]cdb.OptionPattern, 0, 2*len(records))
	for _, record := range records {
		patterns = append(patterns, buildOptionPatterns(record)...)
	}
	return cdb.NewConfig(patterns)
}
