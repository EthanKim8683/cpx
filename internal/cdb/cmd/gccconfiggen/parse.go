package main

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/EthanKim8683/cpx/internal/cdb"
)

var (
	negateRE = regexp.MustCompile(`^[Wfgm][^=]+$`)
	aliasRE  = regexp.MustCompile(`\bAlias\(([^)]+)\)`)
	argsRE   = regexp.MustCompile(`\bArgs\((\d+)\)`)
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

	// First parse Args(N) if present
	if m := argsRE.FindStringSubmatch(record.attrs); len(m) > 1 {
		if n, err := strconv.Atoi(m[1]); err == nil && n > 0 {
			numArgs = n
		}
	}

	if numArgs > 1 {
		kind = cdb.OptionKindMultiArg
	} else if hasAttr(record.attrs, "CommaJoined") {
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

	// alias string is:
	// - if no alias attribute, it's the name of the record
	// - else, it's the alias attribute
	//
	// convert alias string to id by hashing

	// determine alias args from attributes if there are any

	// construct option alias from id and args

	// return option pattern and option alias
	return pattern, cdb.OptionAlias{}
}

func parseOptRecords(records []optRecord) *cdb.Config {
	// expand opt records

	// parse each opt record, building map

	// construct config from map

	// return config
	return nil
}
