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
// It returns an empty string if the option name does not qualify for automated negation.
func negateName(name string) string {
	if !negateRE.MatchString(name) {
		return ""
	}
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
		if hasAttr(r.attrs, "RejectNegative") {
			continue
		}

		neg := negateName(r.name)
		if neg == "" || explicit[neg] {
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
	// determine spelling, kind and num args from attributes

	// alias string is:
	// - if no alias attribute, it's the name of the record
	// - else, it's the alias attribute
	//
	// convert alias string to id by hashing

	// determine alias args from attributes if there are any

	// construct option alias from id and args

	// return option pattern and option alias
	return cdb.OptionPattern{}, cdb.OptionAlias{}
}

func parseOptRecords(records []optRecord) *cdb.Config {
	// expand opt records

	// parse each opt record, building map

	// construct config from map

	// return config
	return nil
}
