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

	var partials []cdb.OptionPattern
	if record.joined {
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindJoined,
		})
	}
	if record.separate {
		switch {
		case record.noDriverArg:
			// NoDriverArg makes the option behave like a flag to the driver.
			partials = append(partials, cdb.OptionPattern{
				Kind: cdb.OptionKindFlag,
			})
		case record.args != 0:
			partials = append(partials, cdb.OptionPattern{
				Kind:    cdb.OptionKindMultiArg,
				NumArgs: record.args,
			})
		default:
			partials = append(partials, cdb.OptionPattern{
				Kind: cdb.OptionKindSeparate,
			})
		}
	}
	if record.joinedOrMissing {
		// JoinedOrMissing can be decomposed into Joined and Flag.
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindFlag,
		})
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindJoined,
		})
	}
	if partials == nil {
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindFlag,
		})
	}

	var patterns []cdb.OptionPattern
	for _, partial := range partials {
		partial.Spelling = "-" + record.name
		patterns = append(patterns, partial)
	}

	if negateRE.MatchString(record.name) && !record.rejectNegative {
		for _, partial := range partials {
			partial.Spelling = "-" + negative(record.name)
			patterns = append(patterns, partial)
		}
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
