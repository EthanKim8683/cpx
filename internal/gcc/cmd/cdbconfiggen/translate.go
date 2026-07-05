// translate.go translates parsed option records into CDB parser configs.

package main

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/EthanKim8683/cpx/internal/cdb"
)

var parenRE = regexp.MustCompile(`\([^)]+\)`)

// negateRE matches flags with implicit negative forms.
var negateRE = regexp.MustCompile(`^[fgWm][^=]+$`)

// mergeOptRecords merges optRecords with the same name into a single record by concatenating their props.
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

// hasProp emulates the behavior of flag_set_p:
// https://github.com/gcc-mirror/gcc/blob/releases/gcc-16/gcc/opt-functions.awk
func hasProp(prop, props string) bool {
	props = parenRE.ReplaceAllString(props, "")
	return strings.Contains(" "+props+" ", " "+prop+" ")
}

// propArgs emulates the behavior of opt_args:
// https://github.com/gcc-mirror/gcc/blob/releases/gcc-16/gcc/opt-functions.awk
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

// negative returns the negated unprefixed spelling of a flag name.
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
			// NoDriverArg induces Flag behavior on driver.
			partials = append(partials, cdb.OptionPattern{
				Kind: cdb.OptionKindFlag,
			})
		case hasProp("Args", record.props):
			// n should already be validated prior to parsing.
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

	if negateRE.MatchString(record.name) && !hasProp("RejectNegative", record.props) {
		for _, partial := range partials {
			partial.Spelling = "-" + negative(record.name)
			patterns = append(patterns, partial)
		}
	}
	return patterns
}

func translateOptRecords(records []optRecord) *cdb.Config {
	records = mergeOptRecords(records)
	patterns := make([]cdb.OptionPattern, 0, 2*len(records))
	for _, record := range records {
		patterns = append(patterns, translateOptRecord(record)...)
	}
	return cdb.NewConfig(patterns)
}
