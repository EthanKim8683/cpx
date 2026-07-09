package main

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/EthanKim8683/cpx/internal/cdb"
)

// parenRE matches parenthesized groups in property strings.
var parenRE = regexp.MustCompile(`\([^)]+\)`)

// negateRE matches option flags (starting with f, g, W, or m) that GCC automatically
// generates negative forms for (e.g., -ffoo has an implicit -fno-foo negated variant).
var negateRE = regexp.MustCompile(`^[fgWm][^=]+$`)

// mergeOptRecords merges options sharing the same spelling into a single entry.
//
// This replicates the merging behavior of GCC's opt-gather script to consolidate
// options defined across multiple target-independent and language-specific files.
func mergeOptRecords(records []optRecord) []optRecord {
	slices.SortFunc(records, func(a, b optRecord) int {
		return strings.Compare(a.name, b.name)
	})

	merged := make([]optRecord, 0, len(records))
	for i := 0; i < len(records); {
		name := records[i].name
		var b strings.Builder
		b.WriteString(records[i].props)
		for i++; i < len(records); i++ {
			if records[i].name != name {
				break
			}
			b.WriteString(" ")
			b.WriteString(records[i].props)
		}
		merged = append(merged, optRecord{
			name:  name,
			props: b.String(),
		})
	}
	return merged
}

// hasProp emulates flag_set_p from GCC's opt-functions.awk.
func hasProp(prop, props string) bool {
	props = parenRE.ReplaceAllString(props, "")
	return strings.Contains(" "+props+" ", " "+prop+" ")
}

// propArgs emulates opt_args from GCC's opt-functions.awk.
func propArgs(name, props string) string {
	_, s, found := strings.Cut(" "+props, " "+name+"(")
	if !found {
		return ""
	}
	s, found = strings.CutPrefix(s, "{")
	if found {
		s, _, _ = strings.Cut(s, "})")
	} else {
		s, _, _ = strings.Cut(s, ")")
	}
	return s
}

// negative returns the negated option name variant.
//
// GCC's option system inserts "no-" after the first letter of the name (e.g.,
// name "ffoo" yields "fno-foo").
func negative(name string) string {
	return name[0:1] + "no-" + name[1:]
}

// translateOptRecord decomposes a parsed option record into option patterns.
func translateOptRecord(record optRecord) []cdb.OptionPattern {
	if hasProp("RejectDriver", record.props) {
		return nil
	}

	partials := []cdb.OptionPattern{}
	if hasProp("Joined", record.props) {
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindJoined,
		})
	}
	if hasProp("Separate", record.props) {
		switch {
		case hasProp("NoDriverArg", record.props):
			// NoDriverArg induces Flag behavior on the driver.
			partials = append(partials, cdb.OptionPattern{
				Kind: cdb.OptionKindFlag,
			})
		case hasProp("Args", record.props):
			//nolint:errcheck // strconv.Atoi is pre-validated by compiler config rules
			n, _ := strconv.Atoi(propArgs("Args", record.props))
			partials = append(partials, cdb.OptionPattern{
				Kind:    cdb.OptionKindMultiArg,
				NumArgs: n,
			})
		default:
			partials = append(partials, cdb.OptionPattern{
				Kind: cdb.OptionKindSeparate,
			})
		}
	}
	if hasProp("JoinedOrMissing", record.props) {
		// JoinedOrMissing decomposes into Flag and Joined.
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindFlag,
		})
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindJoined,
		})
	}
	// GCC options without any recognized property default to Flag.
	if len(partials) == 0 {
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindFlag,
		})
	}

	patterns := []cdb.OptionPattern{}
	for _, partial := range partials {
		partial.Spelling = "-" + record.name
		patterns = append(patterns, partial)
	}

	if negateRE.MatchString(record.name) && !hasProp("RejectNegative", record.props) {
		for _, partial := range partials {
			partial.Spelling = "-" + negative(record.name)
			patterns = append(patterns, partial)
		}
	}
	return patterns
}

// translateOptRecords translates a slice of parsed option records into CDB option patterns.
func translateOptRecords(records []optRecord) []cdb.OptionPattern {
	patterns := []cdb.OptionPattern{}
	for _, record := range records {
		patterns = append(patterns, translateOptRecord(record)...)
	}
	return patterns
}
