package main

import (
	"regexp"
	"strings"

	"github.com/EthanKim8683/cpx/internal/cdb"
)

var (
	negateRE = regexp.MustCompile(`^[Wfgm][^=]+$`)
	aliasRE  = regexp.MustCompile(`\bAlias\(([^)]+)\)`)
)

// hasAttr checks if target is set as a space-separated property in attrs.
func hasAttr(attrs, target string) bool {
	return strings.Contains(" "+attrs+" ", " "+target+" ")
}

// negateName negates a GCC option name by inserting "no-" after the prefix.
// It assumes the option name has already been validated as negatable.
func negateName(name string) string {
	return name[0:1] + "no-" + name[1:]
}

// expandOptRecords synthesizes implicit negative (no-) forms for GCC options
// that lack RejectNegative and do not contain an equals sign.
func expandOptRecords(records []optRecord) []optRecord {
	explicit := make(map[string]bool)
	for _, r := range records {
		explicit[r.name] = true
	}

	expanded := make([]optRecord, 0, len(records)*2)
	expanded = append(expanded, records...)
	for _, r := range records {
		if !negateRE.MatchString(r.name) || hasAttr(r.attrs, "RejectNegative") {
			continue
		}

		neg := negateName(r.name)
		if explicit[neg] {
			continue
		}

		expanded = append(expanded, optRecord{
			name:  neg,
			attrs: "RejectNegative",
		})
	}
	return expanded
}

func parseOptRecord(record optRecord) (cdb.OptionPattern, cdb.OptionAlias) {
	spelling := "-" + record.name
	var kind cdb.OptionKind
	var numArgs int

	if hasAttr(record.attrs, "CommaJoined") {
		kind = cdb.OptionKindCommaJoined
		numArgs = 1
	} else if hasAttr(record.attrs, "Joined") && hasAttr(record.attrs, "Separate") {
		kind = cdb.OptionKindJoinedOrSeparate
		numArgs = 1
	} else if hasAttr(record.attrs, "Joined") {
		kind = cdb.OptionKindJoined
		numArgs = 1
	} else if hasAttr(record.attrs, "Separate") {
		kind = cdb.OptionKindSeparate
		numArgs = 1
	} else {
		kind = cdb.OptionKindFlag
		numArgs = 0
	}

	pattern := cdb.OptionPattern{
		Spelling: spelling,
		Kind:     kind,
		NumArgs:  numArgs,
	}

	// return option pattern and placeholder option alias
	return pattern, cdb.OptionAlias{}
}

func parseOptRecords(records []optRecord) *cdb.Config {
	// expand opt records

	// parse each opt record, building map

	// construct config from map

	// return config
	return nil
}
