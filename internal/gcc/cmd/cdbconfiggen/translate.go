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

// negateRE matches GCC option names with implicit negative forms (e.g. -ffoo → -fno-foo).
var negateRE = regexp.MustCompile(`^[fgWm][^=]+$`)

// mergeOptRecords merges optRecords with the same name into a single record by
// concatenating their props. The input slice is sorted in place.
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

// hasProp emulates flag_set_p from gcc/gcc/opt-functions.awk. It checks whether
// a property like "Joined" or "Separate" appears in the space-separated property
// string, excluding parenthesized groups to avoid matching macro arguments.
func hasProp(prop, props string) bool {
	props = parenRE.ReplaceAllString(props, "")
	return strings.Contains(" "+props+" ", " "+prop+" ")
}

// propArgs emulates opt_args from gcc/gcc/opt-functions.awk. It pulls out the
// value for a named property, handling both parenthesized forms like
// "Var({foo})" and bare forms like "Args(2)". Returns empty string if the
// property isn't present.
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

// negative returns the negated form of a GCC option name by inserting "no-"
// after the first character. For example, "ffoo" becomes "fno-foo".
func negative(name string) string {
	return name[0:1] + "no-" + name[1:]
}

// translateOptRecord decomposes a parsed option record into option patterns.
func translateOptRecord(record optRecord) []cdb.OptionPattern {
	if hasProp("RejectDriver", record.props) {
		return nil
	}

	var partials []cdb.OptionPattern
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

	// Add negated form for negatable names (e.g. -ffoo → -fno-foo).
	if negateRE.MatchString(record.name) && !hasProp("RejectNegative", record.props) {
		for _, partial := range partials {
			partial.Spelling = "-" + negative(record.name)
			patterns = append(patterns, partial)
		}
	}
	return patterns
}

// translateOptRecords translates a slice of parsed option records into a CDB config.
func translateOptRecords(records []optRecord) *cdb.Config {
	var patterns []cdb.OptionPattern
	for _, record := range records {
		patterns = append(patterns, translateOptRecord(record)...)
	}
	return cdb.NewConfig(patterns)
}
